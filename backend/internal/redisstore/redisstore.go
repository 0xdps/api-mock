package redisstore

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/redis/go-redis/v9"
)

// Store manages Redis operations for cache data
type Store struct {
	client *redis.Client
	ctx    context.Context
}

// Config holds Redis connection configuration
type Config struct {
	// URL is the Redis connection URL.
	// Formats:
	//   redis://:password@host:port/db        (plain)
	//   rediss://:password@host:port/db       (TLS)
	URL string
}

// NewStore creates a new Redis store
func NewStore(config Config) (*Store, error) {
	opts, err := redis.ParseURL(config.URL)
	if err != nil {
		return nil, fmt.Errorf("invalid REDIS_URL: %w", err)
	}

	client := redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test connection
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	// Log host without exposing credentials
	redisHost := opts.Addr
	if u, err := url.Parse(config.URL); err == nil {
		redisHost = u.Host
	}
	log.Printf("✅ Connected to Redis at %s", redisHost)

	return &Store{
		client: client,
		ctx:    context.Background(),
	}, nil
}

// SaveItems stores array of items as a Redis list
// Returns number of items saved
func (s *Store) SaveItems(resourceName string, items []map[string]interface{}) (int, error) {
	key := fmt.Sprintf("mockly:resource:%s", resourceName)

	// Clear existing list
	if err := s.client.Del(s.ctx, key).Err(); err != nil {
		return 0, fmt.Errorf("failed to clear resource list: %w", err)
	}

	// Add items to list
	for _, item := range items {
		jsonData, err := json.Marshal(item)
		if err != nil {
			return 0, fmt.Errorf("failed to marshal item: %w", err)
		}

		if err := s.client.RPush(s.ctx, key, string(jsonData)).Err(); err != nil {
			return 0, fmt.Errorf("failed to save item: %w", err)
		}
	}

	log.Printf("  ✓ Saved %d items for %s to Redis", len(items), resourceName)
	return len(items), nil
}

// GetItems retrieves items from Redis list with pagination
// offset: starting position (0-based), limit: number of items to return
func (s *Store) GetItems(resourceName string, offset int, limit int) ([]map[string]interface{}, error) {
	key := fmt.Sprintf("mockly:resource:%s", resourceName)

	// Get items from list
	values, err := s.client.LRange(s.ctx, key, int64(offset), int64(offset+limit-1)).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get items: %w", err)
	}

	items := make([]map[string]interface{}, 0, len(values))
	for _, value := range values {
		var item map[string]interface{}
		if err := json.Unmarshal([]byte(value), &item); err != nil {
			return nil, fmt.Errorf("failed to unmarshal item: %w", err)
		}
		items = append(items, item)
	}

	return items, nil
}

// GetItemByID retrieves a single item by ID from the list
func (s *Store) GetItemByID(resourceName string, id interface{}) (map[string]interface{}, error) {
	key := fmt.Sprintf("mockly:resource:%s", resourceName)

	// Get all items and find by ID (optimization: could use a hash)
	values, err := s.client.LRange(s.ctx, key, 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get items: %w", err)
	}

	for _, value := range values {
		var item map[string]interface{}
		if err := json.Unmarshal([]byte(value), &item); err != nil {
			continue
		}

		if itemID, ok := item["id"]; ok && fmt.Sprint(itemID) == fmt.Sprint(id) {
			return item, nil
		}
	}

	return nil, fmt.Errorf("item not found")
}

// SaveMeta stores metadata for a resource
func (s *Store) SaveMeta(resourceName string, meta interface{}) error {
	key := fmt.Sprintf("mockly:meta:%s", resourceName)

	jsonData, err := json.Marshal(meta)
	if err != nil {
		return fmt.Errorf("failed to marshal meta: %w", err)
	}

	if err := s.client.Set(s.ctx, key, string(jsonData), 0).Err(); err != nil {
		return fmt.Errorf("failed to save meta: %w", err)
	}

	return nil
}

// GetMeta retrieves metadata for a resource
func (s *Store) GetMeta(resourceName string) (interface{}, error) {
	key := fmt.Sprintf("mockly:meta:%s", resourceName)

	value, err := s.client.Get(s.ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("meta not found")
		}
		return nil, fmt.Errorf("failed to get meta: %w", err)
	}

	var meta interface{}
	if err := json.Unmarshal([]byte(value), &meta); err != nil {
		return nil, fmt.Errorf("failed to unmarshal meta: %w", err)
	}

	return meta, nil
}

// GetItemCount returns the number of items in a resource list
func (s *Store) GetItemCount(resourceName string) (int64, error) {
	key := fmt.Sprintf("mockly:resource:%s", resourceName)

	count, err := s.client.LLen(s.ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get item count: %w", err)
	}

	return count, nil
}

// AddItem appends an item to the resource list
func (s *Store) AddItem(resourceName string, item map[string]interface{}) error {
	key := fmt.Sprintf("mockly:resource:%s", resourceName)

	jsonData, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("failed to marshal item: %w", err)
	}

	if err := s.client.RPush(s.ctx, key, string(jsonData)).Err(); err != nil {
		return fmt.Errorf("failed to add item: %w", err)
	}

	return nil
}

// UpdateItem updates an item by ID in the list
func (s *Store) UpdateItem(resourceName string, id interface{}, updates map[string]interface{}) error {
	key := fmt.Sprintf("mockly:resource:%s", resourceName)

	// Get all items
	values, err := s.client.LRange(s.ctx, key, 0, -1).Result()
	if err != nil {
		return fmt.Errorf("failed to get items: %w", err)
	}

	// Find and update item
	for i, value := range values {
		var item map[string]interface{}
		if err := json.Unmarshal([]byte(value), &item); err != nil {
			continue
		}

		if itemID, ok := item["id"]; ok && fmt.Sprint(itemID) == fmt.Sprint(id) {
			// Merge updates
			for k, v := range updates {
				item[k] = v
			}

			// Update in Redis list
			jsonData, err := json.Marshal(item)
			if err != nil {
				return fmt.Errorf("failed to marshal item: %w", err)
			}

			if err := s.client.LSet(s.ctx, key, int64(i), string(jsonData)).Err(); err != nil {
				return fmt.Errorf("failed to update item: %w", err)
			}

			return nil
		}
	}

	return fmt.Errorf("item not found")
}

// DeleteItem removes an item by ID from the list
func (s *Store) DeleteItem(resourceName string, id interface{}) error {
	key := fmt.Sprintf("mockly:resource:%s", resourceName)

	// Get all items
	values, err := s.client.LRange(s.ctx, key, 0, -1).Result()
	if err != nil {
		return fmt.Errorf("failed to get items: %w", err)
	}

	// Find item index and remove
	for i, value := range values {
		var item map[string]interface{}
		if err := json.Unmarshal([]byte(value), &item); err != nil {
			continue
		}

		if itemID, ok := item["id"]; ok && fmt.Sprint(itemID) == fmt.Sprint(id) {
			// Remove by replacing with a sentinel and trimming
			if err := s.client.LSet(s.ctx, key, int64(i), "__DELETE__").Err(); err != nil {
				return fmt.Errorf("failed to mark item for deletion: %w", err)
			}

			// Remove all sentinel values
			if err := s.client.LRem(s.ctx, key, 0, "__DELETE__").Err(); err != nil {
				return fmt.Errorf("failed to remove item: %w", err)
			}

			return nil
		}
	}

	return fmt.Errorf("item not found")
}

// HasResources checks if any resources are stored (checks if at least one resource list exists)
func (s *Store) HasResources(resourceNames []string) (bool, error) {
	for _, resourceName := range resourceNames {
		key := fmt.Sprintf("mockly:resource:%s", resourceName)
		count, err := s.client.LLen(s.ctx, key).Result()
		if err != nil {
			return false, fmt.Errorf("failed to check resource: %w", err)
		}

		if count > 0 {
			return true, nil
		}
	}

	return false, nil
}

// Client returns the underlying *redis.Client for use by other packages
// (e.g. auth middleware, rate limiter).
func (s *Store) Client() *redis.Client {
	return s.client
}

// Close closes the Redis connection
func (s *Store) Close() error {
	return s.client.Close()
}
