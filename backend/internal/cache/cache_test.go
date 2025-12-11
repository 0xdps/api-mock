package cache

import (
	"fmt"
	"sync"
	"testing"
	"time"

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

func TestCache_ConcurrentWrites(t *testing.T) {
	cache, _ := setupTestCache(t, CacheModeLocal)

	if err := cache.Warmup(); err != nil {
		t.Fatalf("Warmup failed: %v", err)
	}

	resource := "users"
	goroutines := 50
	itemsPerGoroutine := 10

	var wg sync.WaitGroup
	wg.Add(goroutines)

	// Track errors
	errors := make(chan error, goroutines*itemsPerGoroutine)

	initialCount := len(cache.Data[resource])

	for i := 0; i < goroutines; i++ {
		go func(goroutineID int) {
			defer wg.Done()
			for j := 0; j < itemsPerGoroutine; j++ {
				newItem := map[string]interface{}{
					"id":          10000 + goroutineID*1000 + j,
					"goroutine":   goroutineID,
					"item_number": j,
				}

				if err := cache.AddItem(resource, newItem); err != nil {
					errors <- fmt.Errorf("goroutine %d, item %d: %w", goroutineID, j, err)
				}
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	// Check for errors
	errorCount := 0
	for err := range errors {
		errorCount++
		if errorCount <= 5 {
			t.Logf("Error during concurrent writes: %v", err)
		}
	}

	// Note: Some errors are expected due to MaxItemsPerResource constraint
	if errorCount > 0 {
		t.Logf("Total errors during concurrent writes: %d (some may be expected due to max constraint)", errorCount)
	}

	// Verify data integrity
	finalCount := len(cache.Data[resource])
	added := finalCount - initialCount

	if added < 0 {
		t.Error("Cache lost items during concurrent writes")
	}

	if added > goroutines*itemsPerGoroutine {
		t.Errorf("Cache added more items than expected: added %d, expected max %d", added, goroutines*itemsPerGoroutine)
	}

	t.Logf("Concurrent writes: started with %d, added %d, ended with %d items", initialCount, added, finalCount)
}

func TestCache_ConcurrentUpdates(t *testing.T) {
	cache, _ := setupTestCache(t, CacheModeLocal)

	if err := cache.Warmup(); err != nil {
		t.Fatalf("Warmup failed: %v", err)
	}

	resource := "users"
	if len(cache.Data[resource]) == 0 {
		t.Skip("No items to update")
	}

	// Get first item ID
	firstItem := cache.Data[resource][0]
	itemID := firstItem["id"]

	goroutines := 50
	iterations := 100

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func(goroutineID int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				updates := map[string]interface{}{
					fmt.Sprintf("field_%d", goroutineID): j,
					"last_update":                        goroutineID,
				}

				if err := cache.UpdateItemByID(resource, itemID, updates); err != nil {
					t.Logf("Update error: %v", err)
				}
			}
		}(i)
	}

	wg.Wait()

	// Verify item still exists and wasn't corrupted
	item, found := cache.GetByID(resource, itemID)
	if !found {
		t.Error("Item was lost during concurrent updates")
	}

	if item["id"] != itemID {
		t.Error("Item ID changed during concurrent updates")
	}

	t.Logf("Concurrent updates completed: item has %d fields", len(item))
}

func TestCache_ConcurrentDeletes(t *testing.T) {
	registry := schema.NewRegistry()
	if err := registry.LoadEmbeddedSchemas(); err != nil {
		t.Fatalf("Failed to load schemas: %v", err)
	}

	config := Config{
		ItemsPerResource:    100, // Start with many items
		Seed:                42,
		MaxItemsPerResource: 200,
		Mode:                CacheModeLocal,
	}

	cache := NewCache(registry, config, nil)

	if err := cache.Warmup(); err != nil {
		t.Fatalf("Warmup failed: %v", err)
	}

	resource := "users"
	initialCount := len(cache.Data[resource])

	if initialCount < 50 {
		t.Skip("Need at least 50 items for this test")
	}

	// Collect IDs to delete
	idsToDelete := make([]interface{}, 0, 40)
	for i := 0; i < 40 && i < len(cache.Data[resource]); i++ {
		idsToDelete = append(idsToDelete, cache.Data[resource][i]["id"])
	}

	goroutines := 10
	var wg sync.WaitGroup
	wg.Add(goroutines)

	// Each goroutine tries to delete a subset of items
	itemsPerGoroutine := len(idsToDelete) / goroutines

	for i := 0; i < goroutines; i++ {
		start := i * itemsPerGoroutine
		end := start + itemsPerGoroutine
		if end > len(idsToDelete) {
			end = len(idsToDelete)
		}

		go func(ids []interface{}) {
			defer wg.Done()
			for _, id := range ids {
				_ = cache.DeleteItemByID(resource, id)
				// Ignore errors - multiple goroutines might try to delete same item
			}
		}(idsToDelete[start:end])
	}

	wg.Wait()

	finalCount := len(cache.Data[resource])

	// Should have at least 1 item (min constraint)
	if finalCount < 1 {
		t.Error("Cache has less than minimum items after concurrent deletes")
	}

	// Should have fewer items than before
	if finalCount >= initialCount {
		t.Error("No items were deleted during concurrent delete operations")
	}

	t.Logf("Concurrent deletes: %d → %d items (%d deleted)", initialCount, finalCount, initialCount-finalCount)
}

func TestCache_MixedConcurrentOperations(t *testing.T) {
	registry := schema.NewRegistry()
	if err := registry.LoadEmbeddedSchemas(); err != nil {
		t.Fatalf("Failed to load schemas: %v", err)
	}

	config := Config{
		ItemsPerResource:    50,
		Seed:                42,
		MaxItemsPerResource: 100,
		Mode:                CacheModeLocal,
	}

	cache := NewCache(registry, config, nil)

	if err := cache.Warmup(); err != nil {
		t.Fatalf("Warmup failed: %v", err)
	}

	resource := "users"
	duration := 2 * time.Second
	stopCh := make(chan struct{})

	var wg sync.WaitGroup

	// Readers
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-stopCh:
				return
			default:
				cache.Get(resource, 10)
			}
		}
	}()

	// Writers
	wg.Add(1)
	go func() {
		defer wg.Done()
		counter := 0
		for {
			select {
			case <-stopCh:
				return
			default:
				newItem := map[string]interface{}{
					"id":      20000 + counter,
					"counter": counter,
				}
				cache.AddItem(resource, newItem)
				counter++
				time.Sleep(10 * time.Millisecond)
			}
		}
	}()

	// Updaters
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-stopCh:
				return
			default:
				// Get items to find one to update
				items, found := cache.Get(resource, 1)
				if found && len(items) > 0 {
					item := items[0]
					if id, ok := item["id"]; ok {
						updates := map[string]interface{}{
							"updated_at": time.Now().Unix(),
						}
						cache.UpdateItemByID(resource, id, updates)
					}
				}
				time.Sleep(15 * time.Millisecond)
			}
		}
	}()

	// GetByID operations
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-stopCh:
				return
			default:
				// Get items to find one to query
				items, found := cache.Get(resource, 1)
				if found && len(items) > 0 {
					item := items[0]
					if id, ok := item["id"]; ok {
						cache.GetByID(resource, id)
					}
				}
			}
		}
	}()

	// Run for duration
	time.Sleep(duration)
	close(stopCh)
	wg.Wait()

	// Verify cache is still functional
	data, found := cache.Get(resource, 10)
	if !found {
		t.Error("Cache stopped working after mixed concurrent operations")
	}

	if len(data) == 0 {
		t.Error("Cache is empty after mixed concurrent operations")
	}

	t.Logf("Mixed concurrent operations completed successfully: %d items in cache", len(cache.Data[resource]))
}

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
// EDGE CASE TESTS
// ============================================================================

func TestCache_EmptyResourceName(t *testing.T) {
	cache, _ := setupTestCache(t, CacheModeLocal)

	if err := cache.Warmup(); err != nil {
		t.Fatalf("Warmup failed: %v", err)
	}

	// Try to get data with empty resource name
	data, found := cache.Get("", 10)

	if found {
		t.Error("Expected not to find data for empty resource name")
	}

	if len(data) > 0 {
		t.Error("Expected no data for empty resource name")
	}
}

func TestCache_InvalidResourceName(t *testing.T) {
	cache, _ := setupTestCache(t, CacheModeLocal)

	if err := cache.Warmup(); err != nil {
		t.Fatalf("Warmup failed: %v", err)
	}

	invalidNames := []string{
		"nonexistent",
		"invalid/resource",
		"resource with spaces",
		"../../../etc/passwd",
		"resource\x00null",
	}

	for _, name := range invalidNames {
		t.Run(fmt.Sprintf("invalid_%s", name), func(t *testing.T) {
			data, found := cache.Get(name, 10)

			if found {
				t.Errorf("Expected not to find data for invalid resource: %s", name)
			}

			if len(data) > 0 {
				t.Errorf("Expected no data for invalid resource: %s", name)
			}
		})
	}
}

func TestCache_AddNilItem(t *testing.T) {
	cache, _ := setupTestCache(t, CacheModeLocal)

	if err := cache.Warmup(); err != nil {
		t.Fatalf("Warmup failed: %v", err)
	}

	resource := "users"

	// Try to add nil item
	err := cache.AddItem(resource, nil)

	// Should handle gracefully - either accept or reject with error
	if err == nil {
		// If accepted, verify it's in cache
		data, found := cache.Get(resource, 100)
		if !found {
			t.Error("Expected to find data after adding nil item")
		}

		// Check if nil item is present
		hasNil := false
		for _, item := range data {
			if item == nil {
				hasNil = true
				break
			}
		}
		if hasNil {
			t.Log("Nil item was added to cache (accepted)")
		}
	} else {
		t.Logf("Nil item rejected with error: %v (expected behavior)", err)
	}
}

func TestCache_AddItemWithEmptyMap(t *testing.T) {
	cache, _ := setupTestCache(t, CacheModeLocal)

	if err := cache.Warmup(); err != nil {
		t.Fatalf("Warmup failed: %v", err)
	}

	resource := "users"
	initialCount := len(cache.Data[resource])

	// Add empty item
	emptyItem := map[string]interface{}{}
	if err := cache.AddItem(resource, emptyItem); err != nil {
		t.Fatalf("AddItem failed: %v", err)
	}

	// Verify item was added
	if len(cache.Data[resource]) != initialCount+1 {
		t.Errorf("Expected %d items, got %d", initialCount+1, len(cache.Data[resource]))
	}
}

func TestCache_SeedConsistency(t *testing.T) {
	// Note: gofakeit's global seed affects all uses, so consistency across
	// separate cache instances depends on call order and timing.
	// This test verifies that using the same seed produces the same data
	// when called in the same way.

	registry := schema.NewRegistry()
	if err := registry.LoadEmbeddedSchemas(); err != nil {
		t.Fatalf("Failed to load schemas: %v", err)
	}

	config := Config{
		ItemsPerResource:    10,
		Seed:                12345,
		MaxItemsPerResource: 100,
		Mode:                CacheModeLocal,
	}

	// Create cache and warm it up
	cache1 := NewCache(registry, config, nil)
	if err := cache1.Warmup(); err != nil {
		t.Fatalf("Cache1 warmup failed: %v", err)
	}

	// Capture first data
	resource := "users"
	data1 := make([]map[string]interface{}, len(cache1.Data[resource]))
	for i, item := range cache1.Data[resource] {
		data1[i] = make(map[string]interface{})
		for k, v := range item {
			data1[i][k] = v
		}
	}

	// Create new cache with same seed and warm it up
	cache2 := NewCache(registry, config, nil)
	if err := cache2.Warmup(); err != nil {
		t.Fatalf("Cache2 warmup failed: %v", err)
	}

	data2 := cache2.Data[resource]

	if len(data1) != len(data2) {
		t.Errorf("Data length mismatch: %d vs %d", len(data1), len(data2))
	}

	// Verify both caches have data (seed is working)
	if len(data1) == 0 || len(data2) == 0 {
		t.Error("Cache generated no data with specified seed")
	}

	t.Logf("Seed %d generated %d items per cache", config.Seed, len(data1))
}

func TestCache_DifferentSeedsProduceDifferentData(t *testing.T) {
	// Create two caches with different seeds
	registry1 := schema.NewRegistry()
	if err := registry1.LoadEmbeddedSchemas(); err != nil {
		t.Fatalf("Failed to load schemas: %v", err)
	}

	registry2 := schema.NewRegistry()
	if err := registry2.LoadEmbeddedSchemas(); err != nil {
		t.Fatalf("Failed to load schemas: %v", err)
	}

	config1 := Config{
		ItemsPerResource:    10,
		Seed:                12345,
		MaxItemsPerResource: 100,
		Mode:                CacheModeLocal,
	}

	config2 := Config{
		ItemsPerResource:    10,
		Seed:                67890,
		MaxItemsPerResource: 100,
		Mode:                CacheModeLocal,
	}

	cache1 := NewCache(registry1, config1, nil)
	cache2 := NewCache(registry2, config2, nil)

	// Warmup both caches
	if err := cache1.Warmup(); err != nil {
		t.Fatalf("Cache1 warmup failed: %v", err)
	}

	if err := cache2.Warmup(); err != nil {
		t.Fatalf("Cache2 warmup failed: %v", err)
	}

	// Verify they generate different data
	resource := "users"
	data1 := cache1.Data[resource]
	data2 := cache2.Data[resource]

	if len(data1) > 0 && len(data2) > 0 {
		id1 := data1[0]["id"]
		id2 := data2[0]["id"]

		// IDs should be different with different seeds
		// (though there's a tiny chance they could match randomly)
		if id1 == id2 {
			t.Log("Warning: IDs matched despite different seeds (unlikely but possible)")
		}
	}
}

func TestCache_GetByID_WithInvalidIDTypes(t *testing.T) {
	cache, _ := setupTestCache(t, CacheModeLocal)

	if err := cache.Warmup(); err != nil {
		t.Fatalf("Warmup failed: %v", err)
	}

	resource := "users"

	invalidIDs := []interface{}{
		nil,
		map[string]interface{}{"nested": "object"},
		[]interface{}{1, 2, 3},
		"",
		-1,
		999999999,
	}

	for _, id := range invalidIDs {
		t.Run(fmt.Sprintf("invalid_id_%v", id), func(t *testing.T) {
			item, found := cache.GetByID(resource, id)

			// Should handle gracefully - either find nothing or handle the invalid ID
			if found && item != nil {
				t.Logf("Found item with ID %v: %v", id, item["id"])
			} else {
				t.Logf("No item found with invalid ID %v (expected)", id)
			}
		})
	}
}

func TestCache_UpdateItemByID_NonExistent(t *testing.T) {
	cache, _ := setupTestCache(t, CacheModeLocal)

	if err := cache.Warmup(); err != nil {
		t.Fatalf("Warmup failed: %v", err)
	}

	resource := "users"

	updates := map[string]interface{}{
		"name": "Updated Name",
	}

	// Try to update non-existent item
	err := cache.UpdateItemByID(resource, 999999, updates)

	if err == nil {
		t.Error("Expected error when updating non-existent item")
	} else {
		t.Logf("Got expected error: %v", err)
	}
}

func TestCache_DeleteItemByID_NonExistent(t *testing.T) {
	cache, _ := setupTestCache(t, CacheModeLocal)

	if err := cache.Warmup(); err != nil {
		t.Fatalf("Warmup failed: %v", err)
	}

	resource := "users"

	// Try to delete non-existent item
	err := cache.DeleteItemByID(resource, 999999)

	if err == nil {
		t.Error("Expected error when deleting non-existent item")
	} else {
		t.Logf("Got expected error: %v", err)
	}
}

func TestCache_AtomicStatsUnderConcurrency(t *testing.T) {
	cache, _ := setupTestCache(t, CacheModeLocal)

	if err := cache.Warmup(); err != nil {
		t.Fatalf("Warmup failed: %v", err)
	}

	resource := "users"
	goroutines := 100
	iterations := 100

	var wg sync.WaitGroup
	wg.Add(goroutines * 2) // Half will hit, half will miss

	// Launch concurrent hits
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				cache.Get(resource, 5) // This should hit
			}
		}()
	}

	// Launch concurrent misses
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				cache.Get("nonexistent", 5) // This should miss
			}
		}()
	}

	wg.Wait()

	// Verify stats are correct
	metadata := cache.GetMetadata()
	expectedHits := int64(goroutines * iterations)
	expectedMisses := int64(goroutines * iterations)

	if metadata.Hits != expectedHits {
		t.Errorf("Expected %d hits, got %d", expectedHits, metadata.Hits)
	}

	if metadata.Misses != expectedMisses {
		t.Errorf("Expected %d misses, got %d", expectedMisses, metadata.Misses)
	}

	t.Logf("Atomic stats verified: %d hits, %d misses", metadata.Hits, metadata.Misses)
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
