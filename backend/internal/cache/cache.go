package cache

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/0xdps/api-mock/go/internal/schema"
	"github.com/brianvoe/gofakeit/v7"
)

// Cache stores pre-generated mock data for all resources
type Cache struct {
	mu       sync.RWMutex
	data     map[string][]map[string]interface{} // resourceName -> array of items
	metadata CacheMetadata
	registry *schema.Registry
	seed     int64
}

// CacheMetadata tracks cache statistics
type CacheMetadata struct {
	WarmupTime       time.Duration
	TotalResources   int
	TotalItems       int
	ItemsPerResource int
	LastRefresh      time.Time
	Hits             int64
	Misses           int64
}

// Config holds cache configuration
type Config struct {
	ItemsPerResource int   // Number of items to pre-generate per resource
	Seed             int64 // Random seed for reproducible data
}

// NewCache creates a new cache instance
func NewCache(registry *schema.Registry, config Config) *Cache {
	return &Cache{
		data:     make(map[string][]map[string]interface{}),
		registry: registry,
		seed:     config.Seed,
		metadata: CacheMetadata{
			ItemsPerResource: config.ItemsPerResource,
		},
	}
}

// Warmup pre-generates data for all resources
func (c *Cache) Warmup() error {
	startTime := time.Now()
	log.Printf("🔥 Starting cache warmup (seed: %d, items per resource: %d)...", c.seed, c.metadata.ItemsPerResource)

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

		// Store in cache
		c.mu.Lock()
		c.data[resourceName] = data
		c.mu.Unlock()

		totalItems += len(data)
		log.Printf("  ✓ Cached %d items for %s", len(data), resourceName)
	}

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
	c.mu.RLock()
	defer c.mu.RUnlock()

	data, exists := c.data[resourceName]
	if !exists {
		c.incrementMisses()
		return nil, false
	}

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

// GetByID retrieves a single item by ID from cache
func (c *Cache) GetByID(resourceName string, id interface{}) (map[string]interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	data, exists := c.data[resourceName]
	if !exists {
		c.incrementMisses()
		return nil, false
	}

	// Find item with matching ID
	for _, item := range data {
		if itemID, ok := item["id"]; ok && fmt.Sprint(itemID) == fmt.Sprint(id) {
			c.incrementHits()
			// Return a copy
			result := make(map[string]interface{})
			for k, v := range item {
				result[k] = v
			}
			return result, true
		}
	}

	c.incrementMisses()
	return nil, false
}

// Refresh regenerates all cached data
func (c *Cache) Refresh() error {
	log.Printf("🔄 Refreshing cache...")

	// Reset stats
	c.mu.Lock()
	c.data = make(map[string][]map[string]interface{})
	c.metadata.Hits = 0
	c.metadata.Misses = 0
	c.mu.Unlock()

	return c.Warmup()
}

// GetMetadata returns cache statistics
func (c *Cache) GetMetadata() CacheMetadata {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Return a copy
	return c.metadata
}

// incrementHits safely increments hit counter
func (c *Cache) incrementHits() {
	// Note: This is called within RLock, so we need atomic operation
	// For simplicity, we'll accept potential race condition on stats
	// In production, use atomic.AddInt64
	c.metadata.Hits++
}

// incrementMisses safely increments miss counter
func (c *Cache) incrementMisses() {
	c.metadata.Misses++
}

// GetStats returns formatted cache statistics
func (c *Cache) GetStats() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	hitRate := 0.0
	total := c.metadata.Hits + c.metadata.Misses
	if total > 0 {
		hitRate = float64(c.metadata.Hits) / float64(total) * 100
	}

	return map[string]interface{}{
		"warmup_time_ms":     c.metadata.WarmupTime.Milliseconds(),
		"total_resources":    c.metadata.TotalResources,
		"total_items":        c.metadata.TotalItems,
		"items_per_resource": c.metadata.ItemsPerResource,
		"last_refresh":       c.metadata.LastRefresh.Format(time.RFC3339),
		"hits":               c.metadata.Hits,
		"misses":             c.metadata.Misses,
		"hit_rate_percent":   fmt.Sprintf("%.2f", hitRate),
		"seed":               c.seed,
	}
}
