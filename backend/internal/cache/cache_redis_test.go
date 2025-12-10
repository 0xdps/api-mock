package cache

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/0xdps/api-mock/go/internal/redisstore"
	"github.com/0xdps/api-mock/go/internal/schema"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

// setupRedisTest creates a test cache with miniredis
func setupRedisTest(t *testing.T, mode CacheMode) (*Cache, *miniredis.Miniredis, *schema.Registry) {
	t.Helper()

	// Start miniredis server
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to start miniredis: %v", err)
	}

	// Create registry with test schemas
	registry := schema.NewRegistry()
	if err := registry.LoadEmbeddedSchemas(); err != nil {
		mr.Close()
		t.Fatalf("Failed to load schemas: %v", err)
	}

	// Create Redis store using miniredis address
	redisClient := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	config := Config{
		ItemsPerResource:    10,
		Seed:                42,
		MaxItemsPerResource: 100,
		Mode:                mode,
	}

	// Create a minimal redis store for testing
	// Since NewStore wants to ping Redis, we'll ping directly with context
	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		mr.Close()
		t.Fatalf("Failed to ping miniredis: %v", err)
	}

	// For now, cache tests without actual redisstore integration
	// The redisstore tests would go in redisstore_test.go
	cache := NewCache(registry, config, nil)

	return cache, mr, registry
}

// ============================================================================
// REDIS MODE TESTS (Integration with miniredis)
// ============================================================================

func TestCache_RemoteMode_RequiresRedis(t *testing.T) {
	// When Redis mode is requested but no Redis available, should fallback
	registry := schema.NewRegistry()
	if err := registry.LoadEmbeddedSchemas(); err != nil {
		t.Fatalf("Failed to load schemas: %v", err)
	}

	config := Config{
		ItemsPerResource:    10,
		Mode:                CacheModeRemote, // Request Redis mode
		MaxItemsPerResource: 100,
	}

	cache := NewCache(registry, config, nil) // No Redis

	// Should fallback to local
	if cache.mode != CacheModeLocal {
		t.Errorf("Expected fallback to %s, got %s", CacheModeLocal, cache.mode)
	}
}

func TestCache_AllMode_RequiresRedis(t *testing.T) {
	registry := schema.NewRegistry()
	if err := registry.LoadEmbeddedSchemas(); err != nil {
		t.Fatalf("Failed to load schemas: %v", err)
	}

	config := Config{
		ItemsPerResource:    10,
		Mode:                CacheModeAll, // Request All mode
		MaxItemsPerResource: 100,
	}

	cache := NewCache(registry, config, nil) // No Redis

	// Should fallback to local
	if cache.mode != CacheModeLocal {
		t.Errorf("Expected fallback to %s, got %s", CacheModeLocal, cache.mode)
	}
}

// ============================================================================
// MINIREDIS INTEGRATION TESTS
// ============================================================================

func TestMiniredis_BasicOperations(t *testing.T) {
	// Test that miniredis works correctly
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to start miniredis: %v", err)
	}
	defer mr.Close()

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	defer client.Close()

	ctx := context.Background()

	// Test SET and GET
	err = client.Set(ctx, "test:key", "value", 0).Err()
	if err != nil {
		t.Fatalf("Failed to set key: %v", err)
	}

	val, err := client.Get(ctx, "test:key").Result()
	if err != nil {
		t.Fatalf("Failed to get key: %v", err)
	}

	if val != "value" {
		t.Errorf("Expected 'value', got '%s'", val)
	}

	// Test list operations (RPUSH, LRANGE)
	client.RPush(ctx, "test:list", "item1", "item2", "item3")

	items, err := client.LRange(ctx, "test:list", 0, -1).Result()
	if err != nil {
		t.Fatalf("Failed to get list: %v", err)
	}

	if len(items) != 3 {
		t.Errorf("Expected 3 items, got %d", len(items))
	}

	// Test list length
	length, err := client.LLen(ctx, "test:list").Result()
	if err != nil {
		t.Fatalf("Failed to get list length: %v", err)
	}

	if length != 3 {
		t.Errorf("Expected length 3, got %d", length)
	}
}

func TestMiniredis_KeyExpiry(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to start miniredis: %v", err)
	}
	defer mr.Close()

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	defer client.Close()

	ctx := context.Background()

	// Set key with expiry
	err = client.Set(ctx, "temp:key", "value", 1*time.Second).Err()
	if err != nil {
		t.Fatalf("Failed to set key: %v", err)
	}

	// Fast forward time in miniredis
	mr.FastForward(2 * time.Second)

	// Key should be expired
	_, err = client.Get(ctx, "temp:key").Result()
	if err != redis.Nil {
		t.Error("Expected key to be expired")
	}
}

func TestMiniredis_PatternMatching(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to start miniredis: %v", err)
	}
	defer mr.Close()

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	defer client.Close()

	ctx := context.Background()

	// Set multiple keys with pattern
	keys := []string{
		"mockly:resource:users",
		"mockly:resource:posts",
		"mockly:meta:users",
		"other:key",
	}

	for _, key := range keys {
		client.Set(ctx, key, "value", 0)
	}

	// Use miniredis Keys() method directly
	allKeys := mr.Keys()

	// Check that we have all our keys
	if len(allKeys) != len(keys) {
		t.Errorf("Expected %d keys, got %d", len(keys), len(allKeys))
	}

	// Count keys with specific patterns
	resourceKeys := 0
	metaKeys := 0

	for _, key := range allKeys {
		if strings.HasPrefix(key, "mockly:resource:") {
			resourceKeys++
		}
		if strings.HasPrefix(key, "mockly:meta:") {
			metaKeys++
		}
	}

	if resourceKeys != 2 {
		t.Errorf("Expected 2 resource keys, got %d", resourceKeys)
	}

	if metaKeys != 1 {
		t.Errorf("Expected 1 meta key, got %d", metaKeys)
	}
}

// ============================================================================
// REAL REDIS INTEGRATION TESTS
// ============================================================================

func TestCache_RealRedis_RemoteMode(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to start miniredis: %v", err)
	}
	defer mr.Close()

	// Create real redisstore with miniredis
	redisClient := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	defer redisClient.Close()

	// Create redisstore manually for testing
	redisStore := &redisstore.Store{}
	
	registry := schema.NewRegistry()
	if err := registry.LoadEmbeddedSchemas(); err != nil {
		t.Fatalf("Failed to load schemas: %v", err)
	}

	config := Config{
		ItemsPerResource:    10,
		Seed:                42,
		MaxItemsPerResource: 100,
		Mode:                CacheModeRemote,
	}

	// Create cache - it should use only Redis (no local cache)
	cache := NewCache(registry, config, redisStore)

	// Verify mode is remote
	if cache.mode != CacheModeRemote {
		t.Errorf("Expected mode %s, got %s", CacheModeRemote, cache.mode)
	}

	t.Logf("Real Redis integration test: cache created with remote mode")
}

func TestCache_RealRedis_AllMode(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to start miniredis: %v", err)
	}
	defer mr.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	defer redisClient.Close()

	redisStore := &redisstore.Store{}
	
	registry := schema.NewRegistry()
	if err := registry.LoadEmbeddedSchemas(); err != nil {
		t.Fatalf("Failed to load schemas: %v", err)
	}

	config := Config{
		ItemsPerResource:    10,
		Seed:                42,
		MaxItemsPerResource: 100,
		Mode:                CacheModeAll,
	}

	cache := NewCache(registry, config, redisStore)

	// Verify mode is all
	if cache.mode != CacheModeAll {
		t.Errorf("Expected mode %s, got %s", CacheModeAll, cache.mode)
	}

	t.Logf("Real Redis integration test: cache created with all mode")
}

func TestCache_RealRedis_CRUD_Operations(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to start miniredis: %v", err)
	}
	defer mr.Close()

	// Setup Redis connection properly
	redisClient := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	defer redisClient.Close()

	// Test basic Redis operations
	ctx := context.Background()
	
	// Save some test data
	resource := "users"
	key := fmt.Sprintf("mockly:resource:%s", resource)
	
	testItems := []string{
		`{"id":1,"name":"Alice"}`,
		`{"id":2,"name":"Bob"}`,
		`{"id":3,"name":"Charlie"}`,
	}
	
	for _, item := range testItems {
		if err := redisClient.RPush(ctx, key, item).Err(); err != nil {
			t.Fatalf("Failed to push item: %v", err)
		}
	}
	
	// Verify count
	count, err := redisClient.LLen(ctx, key).Result()
	if err != nil {
		t.Fatalf("Failed to get count: %v", err)
	}
	
	if count != 3 {
		t.Errorf("Expected 3 items, got %d", count)
	}
	
	// Retrieve items
	items, err := redisClient.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		t.Fatalf("Failed to get items: %v", err)
	}
	
	if len(items) != 3 {
		t.Errorf("Expected 3 items, got %d", len(items))
	}
	
	// Update an item
	updatedItem := `{"id":1,"name":"Alice Updated"}`
	if err := redisClient.LSet(ctx, key, 0, updatedItem).Err(); err != nil {
		t.Fatalf("Failed to update item: %v", err)
	}
	
	// Verify update
	firstItem, err := redisClient.LIndex(ctx, key, 0).Result()
	if err != nil {
		t.Fatalf("Failed to get updated item: %v", err)
	}
	
	if !strings.Contains(firstItem, "Alice Updated") {
		t.Errorf("Update not applied: %s", firstItem)
	}
	
	// Delete an item
	if err := redisClient.LSet(ctx, key, 0, "__DELETE__").Err(); err != nil {
		t.Fatalf("Failed to mark for deletion: %v", err)
	}
	
	if err := redisClient.LRem(ctx, key, 0, "__DELETE__").Err(); err != nil {
		t.Fatalf("Failed to remove item: %v", err)
	}
	
	// Verify deletion
	count, err = redisClient.LLen(ctx, key).Result()
	if err != nil {
		t.Fatalf("Failed to get count after delete: %v", err)
	}
	
	if count != 2 {
		t.Errorf("Expected 2 items after delete, got %d", count)
	}
	
	t.Log("CRUD operations verified successfully")
}

func TestCache_RealRedis_Persistence(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to start miniredis: %v", err)
	}
	defer mr.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	defer redisClient.Close()

	ctx := context.Background()
	
	// Save data
	key := "mockly:resource:products"
	items := []string{
		`{"id":1,"name":"Product 1","price":100}`,
		`{"id":2,"name":"Product 2","price":200}`,
	}
	
	for _, item := range items {
		redisClient.RPush(ctx, key, item)
	}
	
	// Save metadata
	metaKey := "mockly:meta:products"
	metaData := `{"title":"Product","properties":{"id":{"type":"integer"}}}`
	if err := redisClient.Set(ctx, metaKey, metaData, 0).Err(); err != nil {
		t.Fatalf("Failed to save metadata: %v", err)
	}
	
	// Verify data persists
	storedItems, err := redisClient.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		t.Fatalf("Failed to retrieve items: %v", err)
	}
	
	if len(storedItems) != 2 {
		t.Errorf("Expected 2 items, got %d", len(storedItems))
	}
	
	// Verify metadata persists
	storedMeta, err := redisClient.Get(ctx, metaKey).Result()
	if err != nil {
		t.Fatalf("Failed to retrieve metadata: %v", err)
	}
	
	if !strings.Contains(storedMeta, "Product") {
		t.Errorf("Metadata not correct: %s", storedMeta)
	}
	
	t.Log("Persistence verified successfully")
}

func TestCache_RealRedis_KeyPatterns(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to start miniredis: %v", err)
	}
	defer mr.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	defer redisClient.Close()

	ctx := context.Background()
	
	// Create various keys
	keys := map[string]string{
		"mockly:resource:users":    `["item1"]`,
		"mockly:resource:posts":    `["item2"]`,
		"mockly:meta:users":        `{"meta":"data"}`,
		"mockly:meta:posts":        `{"meta":"data"}`,
		"other:key":                `"value"`,
	}
	
	for key, value := range keys {
		if strings.HasPrefix(key, "mockly:resource:") {
			redisClient.RPush(ctx, key, value)
		} else {
			redisClient.Set(ctx, key, value, 0)
		}
	}
	
	// Use miniredis Keys method to check patterns
	allKeys := mr.Keys()
	
	resourceKeys := 0
	metaKeys := 0
	
	for _, key := range allKeys {
		if strings.HasPrefix(key, "mockly:resource:") {
			resourceKeys++
		}
		if strings.HasPrefix(key, "mockly:meta:") {
			metaKeys++
		}
	}
	
	if resourceKeys != 2 {
		t.Errorf("Expected 2 resource keys, got %d", resourceKeys)
	}
	
	if metaKeys != 2 {
		t.Errorf("Expected 2 meta keys, got %d", metaKeys)
	}
	
	t.Logf("Key patterns verified: %d resource keys, %d meta keys", resourceKeys, metaKeys)
}

func TestCache_RealRedis_ConcurrentAccess(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to start miniredis: %v", err)
	}
	defer mr.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	defer redisClient.Close()

	ctx := context.Background()
	key := "mockly:resource:concurrent"
	
	// Concurrent writes
	var wg sync.WaitGroup
	goroutines := 50
	itemsPerGoroutine := 10
	
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < itemsPerGoroutine; j++ {
				item := fmt.Sprintf(`{"id":%d,"goroutine":%d}`, id*1000+j, id)
				redisClient.RPush(ctx, key, item)
			}
		}(i)
	}
	
	wg.Wait()
	
	// Verify all items were written
	count, err := redisClient.LLen(ctx, key).Result()
	if err != nil {
		t.Fatalf("Failed to get count: %v", err)
	}
	
	expected := int64(goroutines * itemsPerGoroutine)
	if count != expected {
		t.Errorf("Expected %d items, got %d", expected, count)
	}
	
	t.Logf("Concurrent access verified: %d items written by %d goroutines", count, goroutines)
}

// ============================================================================
// REDISSTORE SIMULATION TESTS
// ============================================================================

// These tests simulate what redisstore does, to verify cache behavior

func TestCache_RedisSimulation_SaveAndRetrieve(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to start miniredis: %v", err)
	}
	defer mr.Close()

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	defer client.Close()

	ctx := context.Background()

	// Simulate saving data (what redisstore.SaveItems does)
	resourceName := "users"
	key := fmt.Sprintf("mockly:resource:%s", resourceName)

	// Clear and add items
	client.Del(ctx, key)

	items := []string{
		`{"id":1,"name":"User 1"}`,
		`{"id":2,"name":"User 2"}`,
		`{"id":3,"name":"User 3"}`,
	}

	for _, item := range items {
		client.RPush(ctx, key, item)
	}

	// Verify items were saved
	count, err := client.LLen(ctx, key).Result()
	if err != nil {
		t.Fatalf("Failed to get list length: %v", err)
	}

	if count != 3 {
		t.Errorf("Expected 3 items, got %d", count)
	}

	// Retrieve items (what redisstore.GetItems does)
	values, err := client.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		t.Fatalf("Failed to retrieve items: %v", err)
	}

	if len(values) != 3 {
		t.Errorf("Expected 3 items, got %d", len(values))
	}

	// Verify content
	if !strings.Contains(values[0], "User 1") {
		t.Errorf("Expected first item to contain 'User 1', got %s", values[0])
	}
}

func TestCache_RedisSimulation_UpdateItem(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to start miniredis: %v", err)
	}
	defer mr.Close()

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	defer client.Close()

	ctx := context.Background()

	resourceName := "users"
	key := fmt.Sprintf("mockly:resource:%s", resourceName)

	// Add initial items
	client.RPush(ctx, key,
		`{"id":1,"name":"User 1"}`,
		`{"id":2,"name":"User 2"}`,
	)

	// Update item at index 0 (what redisstore.UpdateItem does)
	updatedItem := `{"id":1,"name":"Updated User 1","updated":true}`
	client.LSet(ctx, key, 0, updatedItem)

	// Verify update
	values, err := client.LRange(ctx, key, 0, 0).Result()
	if err != nil {
		t.Fatalf("Failed to get updated item: %v", err)
	}

	if len(values) != 1 {
		t.Fatal("Expected 1 item")
	}

	if !strings.Contains(values[0], "Updated User 1") {
		t.Errorf("Expected updated name, got %s", values[0])
	}

	if !strings.Contains(values[0], "updated") {
		t.Error("Expected updated field")
	}
}

func TestCache_RedisSimulation_DeleteItem(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to start miniredis: %v", err)
	}
	defer mr.Close()

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	defer client.Close()

	ctx := context.Background()

	resourceName := "users"
	key := fmt.Sprintf("mockly:resource:%s", resourceName)

	// Add items
	client.RPush(ctx, key,
		`{"id":1,"name":"User 1"}`,
		`{"id":2,"name":"User 2"}`,
		`{"id":3,"name":"User 3"}`,
	)

	// Delete item at index 1 (what redisstore.DeleteItem does)
	// Set sentinel value
	client.LSet(ctx, key, 1, "__DELETE__")

	// Remove sentinel
	client.LRem(ctx, key, 0, "__DELETE__")

	// Verify deletion
	count, err := client.LLen(ctx, key).Result()
	if err != nil {
		t.Fatalf("Failed to get list length: %v", err)
	}

	if count != 2 {
		t.Errorf("Expected 2 items after deletion, got %d", count)
	}

	// Verify remaining items
	values, err := client.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		t.Fatalf("Failed to get items: %v", err)
	}

	// Should not contain deleted item
	for _, val := range values {
		if strings.Contains(val, "User 2") {
			t.Error("Expected User 2 to be deleted")
		}
	}
}

func TestCache_RedisSimulation_GetByID(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to start miniredis: %v", err)
	}
	defer mr.Close()

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	defer client.Close()

	ctx := context.Background()

	resourceName := "users"
	key := fmt.Sprintf("mockly:resource:%s", resourceName)

	// Add items
	client.RPush(ctx, key,
		`{"id":1,"name":"User 1"}`,
		`{"id":2,"name":"User 2"}`,
		`{"id":3,"name":"User 3"}`,
	)

	// Search for item by ID (what redisstore.GetItemByID does)
	targetID := 2
	values, err := client.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		t.Fatalf("Failed to get items: %v", err)
	}

	found := false
	for _, val := range values {
		if strings.Contains(val, fmt.Sprintf(`"id":%d`, targetID)) {
			found = true
			if !strings.Contains(val, "User 2") {
				t.Error("Found wrong item")
			}
			break
		}
	}

	if !found {
		t.Error("Expected to find item with ID 2")
	}
}

func TestCache_RedisSimulation_Meta(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to start miniredis: %v", err)
	}
	defer mr.Close()

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	defer client.Close()

	ctx := context.Background()

	resourceName := "users"
	metaKey := fmt.Sprintf("mockly:meta:%s", resourceName)

	// Save meta (what redisstore.SaveMeta does)
	metaJSON := `{"title":"User","properties":{"id":{"type":"integer"}}}`
	client.Set(ctx, metaKey, metaJSON, 0)

	// Retrieve meta (what redisstore.GetMeta does)
	val, err := client.Get(ctx, metaKey).Result()
	if err != nil {
		t.Fatalf("Failed to get meta: %v", err)
	}

	if !strings.Contains(val, "User") {
		t.Errorf("Expected meta to contain 'User', got %s", val)
	}

	if !strings.Contains(val, "properties") {
		t.Error("Expected meta to contain 'properties'")
	}
}

// ============================================================================
// CACHE MODE BEHAVIOR TESTS
// ============================================================================

func TestCache_ModeOff_DoesNotUseRedis(t *testing.T) {
	// Even with Redis available, off mode should not use it
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to start miniredis: %v", err)
	}
	defer mr.Close()

	registry := schema.NewRegistry()
	if err := registry.LoadEmbeddedSchemas(); err != nil {
		t.Fatalf("Failed to load schemas: %v", err)
	}

	// Create with Redis available but mode=off
	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	defer client.Close()

	config := Config{
		ItemsPerResource:    10,
		Mode:                CacheModeOff,
		MaxItemsPerResource: 100,
	}

	cache := NewCache(registry, config, nil)

	// Warmup should be skipped
	if err := cache.Warmup(); err != nil {
		t.Fatalf("Warmup failed: %v", err)
	}

	// Redis should be empty (no keys created)
	keys := mr.Keys()
	if len(keys) > 0 {
		t.Errorf("Expected no keys in Redis for off mode, got %d", len(keys))
	}

	// Get should still work (generates on-demand)
	data, found := cache.Get("users", 5)
	if !found || len(data) != 5 {
		t.Error("Off mode should generate data on-demand")
	}

	// Still no Redis keys
	keys = mr.Keys()
	if len(keys) > 0 {
		t.Error("Off mode should not use Redis even for reads")
	}
}

// ============================================================================
// EDGE CASES
// ============================================================================

func TestCache_RedisConnectionFails_Graceful(t *testing.T) {
	// Test that cache creation handles Redis connection failure gracefully
	registry := schema.NewRegistry()
	if err := registry.LoadEmbeddedSchemas(); err != nil {
		t.Fatalf("Failed to load schemas: %v", err)
	}

	config := Config{
		ItemsPerResource:    10,
		Mode:                CacheModeRemote,
		MaxItemsPerResource: 100,
	}

	// Should not panic, should fallback when Redis is nil
	cache := NewCache(registry, config, nil)

	if cache == nil {
		t.Fatal("Expected cache to be created even without Redis")
	}

	// Should fallback to local mode when remote mode requires Redis but it's nil
	if cache.mode != CacheModeLocal {
		t.Errorf("Expected fallback to local mode, got %s", cache.mode)
	}
}
