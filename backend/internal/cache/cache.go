package cache

import (
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/0xdps/api-mock/go/internal/redisstore"
	"github.com/0xdps/api-mock/go/internal/schema"
	"github.com/brianvoe/gofakeit/v7"
)

// CacheMode defines the caching strategy
type CacheMode string

const (
	// CacheModeOff disables all caching (generates data on-demand)
	CacheModeOff CacheMode = "off"
	// CacheModeLocal uses only in-memory cache
	CacheModeLocal CacheMode = "local"
	// CacheModeRemote uses only Redis cache
	CacheModeRemote CacheMode = "remote"
	// CacheModeAll uses both local and Redis cache
	CacheModeAll CacheMode = "all"
)

// Cache stores pre-generated mock data for all resources
type Cache struct {
	mu        sync.RWMutex
	Data      map[string][]map[string]interface{} // resourceName -> array of items (in-memory) - exported for testing
	metaCache map[string]interface{}              // resourceName -> meta response
	metadata  CacheMetadata
	registry  *schema.Registry
	seed      int64
	redis     *redisstore.Store // Redis for persistence
	config    Config
	mode      CacheMode          // Current cache mode
}

// CacheMetadata tracks cache statistics
// Note: Hits and Misses use atomic operations for thread-safety
type CacheMetadata struct {
	WarmupTime       time.Duration
	TotalResources   int
	TotalItems       int
	ItemsPerResource int
	LastRefresh      time.Time
	Hits             int64 // Use atomic.AddInt64 and atomic.LoadInt64
	Misses           int64 // Use atomic.AddInt64 and atomic.LoadInt64
}

// Config holds cache configuration
type Config struct {
	ItemsPerResource    int       // Number of items to pre-generate per resource
	Seed                int64     // Random seed for reproducible data
	MaxItemsPerResource int       // Maximum items allowed per resource
	Mode                CacheMode // Cache mode: off, local, remote, or all
}

// NewCache creates a new cache instance
func NewCache(registry *schema.Registry, config Config, redis *redisstore.Store) *Cache {
	// Default to "all" if mode not specified
	mode := config.Mode
	if mode == "" {
		mode = CacheModeAll
	}
	
	// Validate cache mode
	if mode == CacheModeRemote || mode == CacheModeAll {
		if redis == nil {
			log.Printf("⚠️  Cache mode '%s' requires Redis, falling back to 'local'", mode)
			mode = CacheModeLocal
		}
	}
	
	return &Cache{
		Data:      make(map[string][]map[string]interface{}),
		metaCache: make(map[string]interface{}),
		registry:  registry,
		seed:      config.Seed,
		redis:     redis,
		config:    config,
		mode:      mode,
		metadata: CacheMetadata{
			ItemsPerResource: config.ItemsPerResource,
		},
	}
}

// Warmup pre-generates data for all resources and saves according to cache mode
func (c *Cache) Warmup() error {
	if c.mode == CacheModeOff {
		log.Printf("⚠️  Cache mode is 'off' - skipping warmup")
		return nil
	}
	
	startTime := time.Now()
	log.Printf("🔥 Starting cache warmup (mode: %s, seed: %d, items per resource: %d)...", c.mode, c.seed, c.metadata.ItemsPerResource)

	// Set global seed for reproducibility
	gofakeit.Seed(c.seed)

	resourceNames := c.registry.GetAllResourceNames()
	totalItems := 0

	for _, resourceName := range resourceNames {
		// Generate data for this resource
		data, err := c.registry.GenerateData(resourceName, c.metadata.ItemsPerResource)
		if err != nil {
			log.Printf("⚠️  Failed to generate data for %s: %v", resourceName, err)
			continue
		}

		saved := len(data)
		
		// Save to Redis if mode includes remote
		if c.mode == CacheModeRemote || c.mode == CacheModeAll {
			if c.redis != nil {
				count, err := c.redis.SaveItems(resourceName, data)
				if err != nil {
					log.Printf("⚠️  Failed to save %s to Redis: %v", resourceName, err)
					continue
				}
				saved = count
			}
		}

		// Store in in-memory cache if mode includes local
		if c.mode == CacheModeLocal || c.mode == CacheModeAll {
			c.mu.Lock()
			c.Data[resourceName] = data
			c.mu.Unlock()
		}

		totalItems += saved
	}

	// Pre-generate and cache meta responses
	log.Printf("📋 Pre-generating meta responses...")
	for _, resourceName := range resourceNames {
		schema, ok := c.registry.GetSchema(resourceName)
		if !ok {
			continue
		}

		// Use root-level description if available, otherwise fall back to x-resource description
		description := schema.Description
		if description == "" {
			description = schema.Resource.Description
		}

		metaResponse := map[string]interface{}{
			"$schema":       schema.SchemaURI,
			"title":         schema.Title,
			"type":          schema.Type,
			"description":   description,
			"name":          schema.Resource.Name,
			"singular":      schema.Resource.Singular,
			"group":         schema.Resource.Group,
			"properties":    schema.Properties,
			"required":      schema.Required,
			"property_count": len(schema.Properties),
		}

		// Save meta to Redis if mode includes remote
		if c.mode == CacheModeRemote || c.mode == CacheModeAll {
			if c.redis != nil {
				if err := c.redis.SaveMeta(resourceName, metaResponse); err != nil {
					log.Printf("⚠️  Failed to save meta for %s to Redis: %v", resourceName, err)
					continue
				}
			}
		}

		// Store in in-memory cache if mode includes local
		if c.mode == CacheModeLocal || c.mode == CacheModeAll {
			c.mu.Lock()
			c.metaCache[resourceName] = metaResponse
			c.mu.Unlock()
		}
	}
	log.Printf("  ✓ Cached meta responses for %d resources", len(resourceNames))

	c.metadata.WarmupTime = time.Since(startTime)
	c.metadata.TotalResources = len(resourceNames)
	c.metadata.TotalItems = totalItems
	c.metadata.LastRefresh = time.Now()

	log.Printf("✅ Cache warmup completed in %v", c.metadata.WarmupTime)
	log.Printf("   📊 Resources: %d, Total items: %d", c.metadata.TotalResources, c.metadata.TotalItems)

	return nil
}

// Get retrieves items from cache by resource name
func (c *Cache) Get(resourceName string, count int) ([]map[string]interface{}, bool) {
	// If cache is off, generate data on-demand
	if c.mode == CacheModeOff {
		data, err := c.registry.GenerateData(resourceName, count)
		if err != nil {
			return nil, false
		}
		return data, true
	}
	
	// Try local cache first if mode includes local
	if c.mode == CacheModeLocal || c.mode == CacheModeAll {
		c.mu.RLock()
		data, exists := c.Data[resourceName]
		c.mu.RUnlock()
		
		if exists {
			c.incrementHits()
			
			// Return requested count or all if count exceeds available
			if count > len(data) {
				count = len(data)
			}

			// Return a copy to prevent external modification
			result := make([]map[string]interface{}, count)
			copy(result, data[:count])
			return result, true
		}
		
		// If mode is local only and not found, miss
		if c.mode == CacheModeLocal {
			c.incrementMisses()
			return nil, false
		}
	}
	
	// Try Redis if mode includes remote
	if (c.mode == CacheModeRemote || c.mode == CacheModeAll) && c.redis != nil {
		data, err := c.redis.GetItems(resourceName, 0, count)
		if err == nil && len(data) > 0 {
			c.incrementHits()
			return data, true
		}
	}
	
	c.incrementMisses()
	return nil, false
}

// GetByID retrieves a single item by ID from cache
func (c *Cache) GetByID(resourceName string, id interface{}) (map[string]interface{}, bool) {
	// If cache is off, we can't get by ID without cache
	if c.mode == CacheModeOff {
		return nil, false
	}
	
	// Try local cache first if mode includes local
	if c.mode == CacheModeLocal || c.mode == CacheModeAll {
		c.mu.RLock()
		data, exists := c.Data[resourceName]
		
		if exists {
			// Find item with matching ID (keep lock during iteration)
			for _, item := range data {
				if itemID, ok := item["id"]; ok && fmt.Sprint(itemID) == fmt.Sprint(id) {
					// Make a copy while holding the lock
					result := make(map[string]interface{})
					for k, v := range item {
						result[k] = v
					}
					c.mu.RUnlock()
					c.incrementHits()
					return result, true
				}
			}
		}
		c.mu.RUnlock()
	}
	
	// Try Redis if mode includes remote  
	if (c.mode == CacheModeRemote || c.mode == CacheModeAll) && c.redis != nil {
		item, err := c.redis.GetItemByID(resourceName, id)
		if err == nil && item != nil {
			c.incrementHits()
			return item, true
		}
	}

	c.incrementMisses()
	return nil, false
}

// GetMeta retrieves cached meta response for a resource
func (c *Cache) GetMeta(resourceName string) (interface{}, bool) {
	// If cache is off, generate meta on-demand from schema
	if c.mode == CacheModeOff {
		schema, ok := c.registry.GetSchema(resourceName)
		if !ok {
			return nil, false
		}
		
		description := schema.Description
		if description == "" {
			description = schema.Resource.Description
		}
		
		return map[string]interface{}{
			"$schema":       schema.SchemaURI,
			"title":         schema.Title,
			"type":          schema.Type,
			"description":   description,
			"name":          schema.Resource.Name,
			"singular":      schema.Resource.Singular,
			"group":         schema.Resource.Group,
			"properties":    schema.Properties,
			"required":      schema.Required,
			"property_count": len(schema.Properties),
		}, true
	}
	
	// Try local cache if mode includes local
	if c.mode == CacheModeLocal || c.mode == CacheModeAll {
		c.mu.RLock()
		meta, exists := c.metaCache[resourceName]
		c.mu.RUnlock()
		
		if exists {
			return meta, true
		}
	}
	
	// Try Redis if mode includes remote
	if (c.mode == CacheModeRemote || c.mode == CacheModeAll) && c.redis != nil {
		meta, err := c.redis.GetMeta(resourceName)
		if err == nil && meta != nil {
			return meta, true
		}
	}

	return nil, false
}

// Refresh regenerates all cached data
func (c *Cache) Refresh() error {
	log.Printf("🔄 Refreshing cache...")

	// Reset stats
	c.mu.Lock()
	c.Data = make(map[string][]map[string]interface{})
	c.metaCache = make(map[string]interface{})
	c.mu.Unlock()
	
	// Reset atomic counters
	atomic.StoreInt64(&c.metadata.Hits, 0)
	atomic.StoreInt64(&c.metadata.Misses, 0)

	return c.Warmup()
}

// GetMetadata returns cache statistics
func (c *Cache) GetMetadata() CacheMetadata {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Return a copy with atomic-loaded values
	metadata := c.metadata
	metadata.Hits = atomic.LoadInt64(&c.metadata.Hits)
	metadata.Misses = atomic.LoadInt64(&c.metadata.Misses)
	return metadata
}

// incrementHits safely increments hit counter using atomic operations
func (c *Cache) incrementHits() {
	atomic.AddInt64(&c.metadata.Hits, 1)
}

// incrementMisses safely increments miss counter using atomic operations
func (c *Cache) incrementMisses() {
	atomic.AddInt64(&c.metadata.Misses, 1)
}

// GetStats returns formatted cache statistics
func (c *Cache) GetStats() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Use atomic loads for thread-safe reads
	hits := atomic.LoadInt64(&c.metadata.Hits)
	misses := atomic.LoadInt64(&c.metadata.Misses)

	hitRate := 0.0
	total := hits + misses
	if total > 0 {
		hitRate = float64(hits) / float64(total) * 100
	}

	return map[string]interface{}{
		"warmup_time_ms":     c.metadata.WarmupTime.Milliseconds(),
		"total_resources":    c.metadata.TotalResources,
		"total_items":        c.metadata.TotalItems,
		"items_per_resource": c.metadata.ItemsPerResource,
		"last_refresh":       c.metadata.LastRefresh.Format(time.RFC3339),
		"hits":               hits,
		"misses":             misses,
		"hit_rate_percent":   fmt.Sprintf("%.2f", hitRate),
		"seed":               c.seed,
	}
}

// AddItem adds a new item to a resource (CRUD: CREATE)
func (c *Cache) AddItem(resourceName string, item map[string]interface{}) error {
	// Check max limit constraint
	currentCount := 0
	
	if c.mode == CacheModeLocal || c.mode == CacheModeAll {
		c.mu.Lock()
		currentCount = len(c.Data[resourceName])
		c.mu.Unlock()
	}

	if (c.mode == CacheModeRemote || c.mode == CacheModeAll) && c.redis != nil {
		count, err := c.redis.GetItemCount(resourceName)
		if err != nil {
			return fmt.Errorf("failed to check item count: %w", err)
		}
		currentCount = int(count)
	}

	if currentCount >= c.config.MaxItemsPerResource {
		return fmt.Errorf("cannot exceed maximum items limit (%d) for resource", c.config.MaxItemsPerResource)
	}

	// Save to Redis if mode includes remote
	if (c.mode == CacheModeRemote || c.mode == CacheModeAll) && c.redis != nil {
		if err := c.redis.AddItem(resourceName, item); err != nil {
			return fmt.Errorf("failed to save item to Redis: %w", err)
		}
	}

	// Update in-memory cache if mode includes local
	if c.mode == CacheModeLocal || c.mode == CacheModeAll {
		c.mu.Lock()
		defer c.mu.Unlock()
		c.Data[resourceName] = append(c.Data[resourceName], item)
	}
	
	return nil
}

// UpdateItemByID updates an item by ID (CRUD: UPDATE)
func (c *Cache) UpdateItemByID(resourceName string, id interface{}, updates map[string]interface{}) error {
	// Update in Redis if mode includes remote
	if (c.mode == CacheModeRemote || c.mode == CacheModeAll) && c.redis != nil {
		if err := c.redis.UpdateItem(resourceName, id, updates); err != nil {
			return err
		}
	}

	// Update in-memory cache if mode includes local
	if c.mode == CacheModeLocal || c.mode == CacheModeAll {
		c.mu.Lock()
		defer c.mu.Unlock()

		if items, ok := c.Data[resourceName]; ok {
			for i, item := range items {
				if itemID, ok := item["id"]; ok && fmt.Sprint(itemID) == fmt.Sprint(id) {
					// Merge updates
					for k, v := range updates {
						item[k] = v
					}
					c.Data[resourceName][i] = item
					return nil
				}
			}
		}
		return fmt.Errorf("item not found")
	}

	return nil
}

// DeleteItemByID deletes an item by ID (CRUD: DELETE)
func (c *Cache) DeleteItemByID(resourceName string, id interface{}) error {
	// Check minimum constraint - can't delete if only 1 item
	currentCount := 0
	
	if c.mode == CacheModeLocal || c.mode == CacheModeAll {
		c.mu.Lock()
		currentCount = len(c.Data[resourceName])
		c.mu.Unlock()
	}
	
	if (c.mode == CacheModeRemote || c.mode == CacheModeAll) && c.redis != nil {
		count, err := c.redis.GetItemCount(resourceName)
		if err != nil {
			return fmt.Errorf("failed to check item count: %w", err)
		}
		currentCount = int(count)
	}

	if currentCount <= 1 {
		return fmt.Errorf("cannot delete item: resource must have at least 1 item")
	}

	// Delete from Redis if mode includes remote
	if (c.mode == CacheModeRemote || c.mode == CacheModeAll) && c.redis != nil {
		if err := c.redis.DeleteItem(resourceName, id); err != nil {
			return err
		}
	}

	// Delete from in-memory cache if mode includes local
	if c.mode == CacheModeLocal || c.mode == CacheModeAll {
		c.mu.Lock()
		defer c.mu.Unlock()

		if items, ok := c.Data[resourceName]; ok {
			for i, item := range items {
				if itemID, ok := item["id"]; ok && fmt.Sprint(itemID) == fmt.Sprint(id) {
					// Remove item from slice
					c.Data[resourceName] = append(items[:i], items[i+1:]...)
					return nil
				}
			}
		}
		return fmt.Errorf("item not found")
	}

	return nil
}

// LoadFromRedis loads all data from Redis into in-memory cache
func (c *Cache) LoadFromRedis() error {
	log.Printf("📥 Loading data from Redis into memory...")

	resourceNames := c.registry.GetAllResourceNames()
	totalItems := 0

	c.mu.Lock()
	defer c.mu.Unlock()

	for _, resourceName := range resourceNames {
		// Get all items from Redis
		items, err := c.redis.GetItems(resourceName, 0, 999999) // Get all
		if err != nil {
			log.Printf("⚠️  Failed to load %s from Redis: %v", resourceName, err)
			continue
		}

		c.Data[resourceName] = items
		totalItems += len(items)
		log.Printf("  ✓ Loaded %d items for %s from Redis", len(items), resourceName)

		// Load meta from Redis
		meta, err := c.redis.GetMeta(resourceName)
		if err == nil {
			c.metaCache[resourceName] = meta
		}
	}

	c.metadata.TotalResources = len(resourceNames)
	c.metadata.TotalItems = totalItems
	c.metadata.LastRefresh = time.Now()

	log.Printf("✅ Loaded from Redis: %d resources, %d total items", c.metadata.TotalResources, totalItems)
	return nil
}
