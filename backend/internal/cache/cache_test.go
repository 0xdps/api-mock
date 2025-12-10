package cache

import (
	"sync"
	"testing"

	"github.com/0xdps/api-mock/go/internal/schema"
)

// ============================================================================
// TEST HELPERS
// ============================================================================

// setupTestCache creates a test cache without Redis
func setupTestCache(t *testing.T, mode CacheMode) (*Cache, *schema.Registry) {
	t.Helper()

	registry := schema.NewRegistry()
	if err := registry.LoadEmbeddedSchemas(); err != nil {
		t.Fatalf("Failed to load schemas: %v", err)
	}

	config := Config{
		ItemsPerResource:    10,
		Seed:                42,
		MaxItemsPerResource: 100,
		Mode:                mode,
	}

	cache := NewCache(registry, config, nil)
	return cache, registry
}

// ============================================================================
// CACHE MODE TESTS
// ============================================================================

func TestNewCache_LocalMode(t *testing.T) {
	cache, _ := setupTestCache(t, CacheModeLocal)

	if cache == nil {
		t.Fatal("Expected cache to be created")
	}

	if cache.mode != CacheModeLocal {
		t.Errorf("Expected mode %s, got %s", CacheModeLocal, cache.mode)
	}

	if cache.Data == nil {
		t.Error("Expected Data map to be initialized")
	}
}

func TestNewCache_OffMode(t *testing.T) {
	cache, _ := setupTestCache(t, CacheModeOff)

	if cache.mode != CacheModeOff {
		t.Errorf("Expected mode %s, got %s", CacheModeOff, cache.mode)
	}
}

func TestNewCache_ModeFallback(t *testing.T) {
	// When Redis mode is specified but no Redis client provided
	registry := schema.NewRegistry()
	if err := registry.LoadEmbeddedSchemas(); err != nil {
		t.Fatalf("Failed to load schemas: %v", err)
	}

	config := Config{
		ItemsPerResource:    10,
		Mode:                CacheModeRemote, // Request Redis mode
		MaxItemsPerResource: 100,
	}

	cache := NewCache(registry, config, nil) // But no Redis client

	// Should fallback to local mode
	if cache.mode != CacheModeLocal {
		t.Errorf("Expected fallback to %s when Redis not available, got %s", CacheModeLocal, cache.mode)
	}
}

// ============================================================================
// WARMUP TESTS
// ============================================================================

func TestCache_Warmup_LocalMode(t *testing.T) {
	cache, registry := setupTestCache(t, CacheModeLocal)

	if err := cache.Warmup(); err != nil {
		t.Fatalf("Warmup failed: %v", err)
	}

	// Verify data was generated
	resourceNames := registry.GetAllResourceNames()
	if len(resourceNames) == 0 {
		t.Fatal("No resources loaded")
	}

	for _, name := range resourceNames {
		if len(cache.Data[name]) == 0 {
			t.Errorf("No data generated for resource %s", name)
		}

		if len(cache.Data[name]) != cache.config.ItemsPerResource {
			t.Errorf("Expected %d items for %s, got %d",
				cache.config.ItemsPerResource, name, len(cache.Data[name]))
		}
	}

	// Verify metadata
	metadata := cache.GetMetadata()
	if metadata.TotalResources != len(resourceNames) {
		t.Errorf("Expected %d resources in metadata, got %d",
			len(resourceNames), metadata.TotalResources)
	}

	if metadata.TotalItems == 0 {
		t.Error("Expected total items > 0")
	}
}

func TestCache_Warmup_OffMode(t *testing.T) {
	cache, _ := setupTestCache(t, CacheModeOff)

	// Warmup should be skipped in off mode
	if err := cache.Warmup(); err != nil {
		t.Fatalf("Warmup should not fail in off mode: %v", err)
	}

	// Data should still be empty
	if len(cache.Data) > 0 {
		t.Error("Data should be empty in off mode")
	}
}

// ============================================================================
// GET TESTS (Local Mode)
// ============================================================================

func TestCache_Get_LocalMode_Hit(t *testing.T) {
	cache, _ := setupTestCache(t, CacheModeLocal)

	if err := cache.Warmup(); err != nil {
		t.Fatalf("Warmup failed: %v", err)
	}

	resource := "users"
	count := 5

	data, found := cache.Get(resource, count)

	if !found {
		t.Fatal("Expected to find data")
	}

	if len(data) != count {
		t.Errorf("Expected %d items, got %d", count, len(data))
	}

	// Check hit counter increased
	metadata := cache.GetMetadata()
	if metadata.Hits == 0 {
		t.Error("Expected hits > 0")
	}
}

func TestCache_Get_LocalMode_Miss(t *testing.T) {
	cache, _ := setupTestCache(t, CacheModeLocal)

	// Don't warmup - cache is empty

	resource := "users"
	data, found := cache.Get(resource, 5)

	if found {
		t.Error("Expected miss for empty cache")
	}

	if len(data) > 0 {
		t.Error("Expected no data")
	}

	// Check miss counter
	metadata := cache.GetMetadata()
	if metadata.Misses == 0 {
		t.Error("Expected misses > 0")
	}
}

func TestCache_Get_OffMode_GeneratesOnDemand(t *testing.T) {
	cache, _ := setupTestCache(t, CacheModeOff)

	resource := "users"
	count := 5

	// Should generate data on-demand even without warmup
	data, found := cache.Get(resource, count)

	if !found {
		t.Fatal("Expected to find data (generated on-demand)")
	}

	if len(data) != count {
		t.Errorf("Expected %d items, got %d", count, len(data))
	}

	// Cache should still be empty
	if len(cache.Data[resource]) > 0 {
		t.Error("Off mode should not cache data")
	}
}

// ============================================================================
// GET BY ID TESTS
// ============================================================================

func TestCache_GetByID_LocalMode_Found(t *testing.T) {
	cache, _ := setupTestCache(t, CacheModeLocal)

	if err := cache.Warmup(); err != nil {
		t.Fatalf("Warmup failed: %v", err)
	}

	resource := "users"

	// Get first item's ID
	if len(cache.Data[resource]) == 0 {
		t.Fatal("No data in cache")
	}

	firstItem := cache.Data[resource][0]
	itemID := firstItem["id"]

	// Retrieve by ID
	item, found := cache.GetByID(resource, itemID)

	if !found {
		t.Fatal("Expected to find item by ID")
	}

	if item["id"] != itemID {
		t.Errorf("Expected ID %v, got %v", itemID, item["id"])
	}
}

func TestCache_GetByID_LocalMode_NotFound(t *testing.T) {
	cache, _ := setupTestCache(t, CacheModeLocal)

	if err := cache.Warmup(); err != nil {
		t.Fatalf("Warmup failed: %v", err)
	}

	resource := "users"
	nonExistentID := 999999

	item, found := cache.GetByID(resource, nonExistentID)

	if found {
		t.Error("Expected not to find item")
	}

	if item != nil {
		t.Error("Expected nil item")
	}
}

// ============================================================================
// CRUD OPERATIONS TESTS
// ============================================================================

func TestCache_AddItem_LocalMode(t *testing.T) {
	cache, _ := setupTestCache(t, CacheModeLocal)

	if err := cache.Warmup(); err != nil {
		t.Fatalf("Warmup failed: %v", err)
	}

	resource := "users"
	initialCount := len(cache.Data[resource])

	newItem := map[string]interface{}{
		"id":   99999,
		"name": "Test User",
	}

	if err := cache.AddItem(resource, newItem); err != nil {
		t.Fatalf("AddItem failed: %v", err)
	}

	// Verify item was added
	if len(cache.Data[resource]) != initialCount+1 {
		t.Errorf("Expected %d items, got %d", initialCount+1, len(cache.Data[resource]))
	}

	// Verify item exists
	item, found := cache.GetByID(resource, 99999)
	if !found {
		t.Error("Expected to find newly added item")
	}

	if item["name"] != "Test User" {
		t.Errorf("Expected name 'Test User', got %v", item["name"])
	}
}

func TestCache_AddItem_MaxConstraint(t *testing.T) {
	// Create cache with low max limit
	registry := schema.NewRegistry()
	if err := registry.LoadEmbeddedSchemas(); err != nil {
		t.Fatalf("Failed to load schemas: %v", err)
	}

	config := Config{
		ItemsPerResource:    5,
		Seed:                42,
		MaxItemsPerResource: 10, // Low limit
		Mode:                CacheModeLocal,
	}

	cache := NewCache(registry, config, nil)
	if err := cache.Warmup(); err != nil {
		t.Fatalf("Warmup failed: %v", err)
	}

	resource := "users"

	// Add items up to max
	for i := 0; i < 5; i++ {
		newItem := map[string]interface{}{
			"id":   10000 + i,
			"name": "Test User",
		}
		if err := cache.AddItem(resource, newItem); err != nil {
			t.Errorf("Failed to add item %d: %v", i, err)
		}
	}

	// Now at max, next add should fail
	newItem := map[string]interface{}{
		"id":   99999,
		"name": "Too Many",
	}

	err := cache.AddItem(resource, newItem)
	if err == nil {
		t.Error("Expected error when exceeding max items, got nil")
	}
}

func TestCache_UpdateItemByID_LocalMode(t *testing.T) {
	cache, _ := setupTestCache(t, CacheModeLocal)

	if err := cache.Warmup(); err != nil {
		t.Fatalf("Warmup failed: %v", err)
	}

	resource := "users"

	if len(cache.Data[resource]) == 0 {
		t.Fatal("No data in cache")
	}

	firstItem := cache.Data[resource][0]
	itemID := firstItem["id"]

	updates := map[string]interface{}{
		"name":    "Updated Name",
		"updated": true,
	}

	if err := cache.UpdateItemByID(resource, itemID, updates); err != nil {
		t.Fatalf("UpdateItemByID failed: %v", err)
	}

	// Verify updates applied
	item, found := cache.GetByID(resource, itemID)
	if !found {
		t.Fatal("Expected to find updated item")
	}

	if item["name"] != "Updated Name" {
		t.Errorf("Expected name 'Updated Name', got %v", item["name"])
	}

	if item["updated"] != true {
		t.Error("Expected updated field to be true")
	}

	// Verify ID didn't change
	if item["id"] != itemID {
		t.Error("Item ID should not change during update")
	}
}

func TestCache_DeleteItemByID_LocalMode(t *testing.T) {
	cache, _ := setupTestCache(t, CacheModeLocal)

	if err := cache.Warmup(); err != nil {
		t.Fatalf("Warmup failed: %v", err)
	}

	resource := "users"
	initialCount := len(cache.Data[resource])

	if initialCount <= 1 {
		t.Skip("Need at least 2 items for delete test")
	}

	itemToDelete := cache.Data[resource][0]
	itemID := itemToDelete["id"]

	if err := cache.DeleteItemByID(resource, itemID); err != nil {
		t.Fatalf("DeleteItemByID failed: %v", err)
	}

	// Verify count decreased
	if len(cache.Data[resource]) != initialCount-1 {
		t.Errorf("Expected %d items, got %d", initialCount-1, len(cache.Data[resource]))
	}

	// Verify item no longer exists
	_, found := cache.GetByID(resource, itemID)
	if found {
		t.Error("Expected item to be deleted")
	}
}

func TestCache_DeleteItemByID_MinConstraint(t *testing.T) {
	// Create cache with minimal items
	registry := schema.NewRegistry()
	if err := registry.LoadEmbeddedSchemas(); err != nil {
		t.Fatalf("Failed to load schemas: %v", err)
	}

	config := Config{
		ItemsPerResource:    1, // Only 1 item
		Seed:                42,
		MaxItemsPerResource: 100,
		Mode:                CacheModeLocal,
	}

	cache := NewCache(registry, config, nil)
	if err := cache.Warmup(); err != nil {
		t.Fatalf("Warmup failed: %v", err)
	}

	resource := "users"

	if len(cache.Data[resource]) != 1 {
		t.Fatalf("Expected 1 item, got %d", len(cache.Data[resource]))
	}

	itemID := cache.Data[resource][0]["id"]

	// Try to delete the only item - should fail
	err := cache.DeleteItemByID(resource, itemID)

	if err == nil {
		t.Error("Expected error when deleting last item")
	}

	// Verify item still exists
	if len(cache.Data[resource]) != 1 {
		t.Error("Item should not be deleted due to minimum constraint")
	}
}

// ============================================================================
// CONCURRENT ACCESS TESTS
// ============================================================================

func TestCache_ConcurrentReads(t *testing.T) {
	cache, _ := setupTestCache(t, CacheModeLocal)

	if err := cache.Warmup(); err != nil {
		t.Fatalf("Warmup failed: %v", err)
	}

	resource := "users"
	goroutines := 50
	iterations := 100

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				data, found := cache.Get(resource, 5)
				if !found {
					t.Error("Expected to find data")
				}
				if len(data) == 0 {
					t.Error("Expected non-empty data")
				}
			}
		}()
	}

	wg.Wait()

	// Verify no data corruption
	data, found := cache.Get(resource, 5)
	if !found || len(data) != 5 {
		t.Error("Data may have been corrupted by concurrent access")
	}
}

// ============================================================================
// STATS TESTS
// ============================================================================

func TestCache_GetStats(t *testing.T) {
	cache, _ := setupTestCache(t, CacheModeLocal)

	if err := cache.Warmup(); err != nil {
		t.Fatalf("Warmup failed: %v", err)
	}

	// Perform some operations
	cache.Get("users", 5)
	cache.Get("posts", 3)
	cache.Get("nonexistent", 1) // Should miss

	stats := cache.GetStats()

	// Check required fields
	requiredFields := []string{
		"warmup_time_ms",
		"total_resources",
		"total_items",
		"items_per_resource",
		"last_refresh",
		"hits",
		"misses",
		"hit_rate_percent",
		"seed",
	}

	for _, field := range requiredFields {
		if _, exists := stats[field]; !exists {
			t.Errorf("Expected stats to have field %s", field)
		}
	}

	// Verify hit/miss counts
	hits, ok := stats["hits"].(int64)
	if !ok {
		t.Fatal("Expected hits to be int64")
	}

	if hits == 0 {
		t.Error("Expected hits > 0")
	}
}
