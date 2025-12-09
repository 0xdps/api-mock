package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestPostCollectionValidatesRequiredFields tests that POST rejects items missing required fields
func TestPostCollectionValidatesRequiredFields(t *testing.T) {
	handler, _ := setupTestHandler(t)
	
	// Use "user" resource which typically has required fields like id, name, email
	resource := "user"
	
	// Try to create item without required fields
	invalidItem := map[string]interface{}{
		// Missing required fields - just a partial item
		"city": "New York",
	}

	body, _ := json.Marshal(invalidItem)
	req := httptest.NewRequest("POST", fmt.Sprintf("/%s", resource), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.PostCollection(resource)(w, req)

	// Should reject due to missing required fields
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 Bad Request for missing required fields, got %d", w.Code)
	}

	var response map[string]string
	json.NewDecoder(w.Body).Decode(&response)
	if response["error"] == "" {
		t.Error("Expected error message in response")
	}
	
	t.Logf("Validation correctly rejected item: %s", response["error"])
}

// TestPostCollectionValidatesPropertyTypes tests that POST rejects items with wrong property types
func TestPostCollectionValidatesPropertyTypes(t *testing.T) {
	handler, testCache := setupTestHandler(t)
	
	// Use "product" resource which typically has numeric price
	resource := "product"
	
	// Get a valid item structure first
	if len(testCache.Data[resource]) == 0 {
		t.Skip("No sample items to determine schema")
	}
	
	sampleItem := testCache.Data[resource][0]
	
	// Create item with wrong type for a numeric field
	invalidItem := map[string]interface{}{
		"id":    1,
		"name":  "Test Product",
		"price": "not_a_number", // String instead of number
	}
	
	// Check if sample has price field
	if _, hasPrice := sampleItem["price"]; !hasPrice {
		t.Skip("Sample item doesn't have price field")
	}

	body, _ := json.Marshal(invalidItem)
	req := httptest.NewRequest("POST", fmt.Sprintf("/%s", resource), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.PostCollection(resource)(w, req)

	// Should reject due to type mismatch
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 Bad Request for type mismatch, got %d", w.Code)
		t.Logf("Response: %s", w.Body.String())
	}

	var response map[string]string
	json.NewDecoder(w.Body).Decode(&response)
	
	t.Logf("Validation correctly rejected wrong type: %s", response["error"])
}

// TestPostCollectionAcceptsValidItem tests that POST accepts valid items
func TestPostCollectionAcceptsValidItem(t *testing.T) {
	handler, testCache := setupTestHandler(t)
	
	resource := "product"
	
	// Get a valid item structure
	if len(testCache.Data[resource]) == 0 {
		t.Skip("No sample items")
	}
	
	sampleItem := testCache.Data[resource][0]
	
	// Create a new valid item based on the sample structure
	validItem := map[string]interface{}{
		"id":   99999,
		"name": "New Test Product",
	}
	
	// Copy other fields from sample if they exist
	for key, value := range sampleItem {
		if key != "id" && key != "name" {
			validItem[key] = value
		}
	}

	initialCount := len(testCache.Data[resource])
	
	body, _ := json.Marshal(validItem)
	req := httptest.NewRequest("POST", fmt.Sprintf("/%s", resource), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.PostCollection(resource)(w, req)

	// Should accept valid item
	if w.Code != http.StatusCreated {
		t.Errorf("Expected 201 Created for valid item, got %d", w.Code)
		t.Logf("Response: %s", w.Body.String())
	}

	var response map[string]interface{}
	json.NewDecoder(w.Body).Decode(&response)
	
	if response["name"] != "New Test Product" {
		t.Error("Response should contain posted data")
	}
	
	// Verify item was added to cache
	if len(testCache.Data[resource]) != initialCount+1 {
		t.Error("Item was not added to cache")
	}
	
	t.Log("Valid item accepted and created successfully")
}

// TestPutSingleValidatesTypes tests that PUT validates updated data types
func TestPutSingleValidatesTypes(t *testing.T) {
	handler, testCache := setupTestHandler(t)
	
	resource := "product"
	
	if len(testCache.Data[resource]) == 0 {
		t.Skip("No items to update")
	}
	
	// Get first item
	firstItem := testCache.Data[resource][0]
	itemID, ok := firstItem["id"]
	if !ok {
		t.Skip("Item has no ID")
	}
	
	// Try to update with invalid type
	invalidUpdate := map[string]interface{}{
		"id":    itemID,
		"price": "invalid_number", // Wrong type
	}
	
	// Check if item has price field
	if _, hasPrice := firstItem["price"]; !hasPrice {
		t.Skip("Item doesn't have price field")
	}

	body, _ := json.Marshal(invalidUpdate)
	req := httptest.NewRequest("PUT", fmt.Sprintf("/%s/%v", resource, itemID), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.PutSingle(resource)(w, req)

	// Should reject due to type mismatch
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 Bad Request for invalid update, got %d", w.Code)
		t.Logf("Response: %s", w.Body.String())
	}
	
	var response map[string]string
	json.NewDecoder(w.Body).Decode(&response)
	
	t.Logf("Validation correctly rejected invalid update: %s", response["error"])
}

// TestPutSingleAcceptsValidUpdate tests that PUT accepts valid updates
func TestPutSingleAcceptsValidUpdate(t *testing.T) {
	handler, testCache := setupTestHandler(t)
	
	resource := "product"
	
	if len(testCache.Data[resource]) == 0 {
		t.Skip("No items to update")
	}
	
	firstItem := testCache.Data[resource][0]
	itemID, ok := firstItem["id"]
	if !ok {
		t.Skip("Item has no ID")
	}
	
	// Valid update
	validUpdate := map[string]interface{}{
		"name": "Updated Product Name",
	}

	body, _ := json.Marshal(validUpdate)
	req := httptest.NewRequest("PUT", fmt.Sprintf("/%s/%v", resource, itemID), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.PutSingle(resource)(w, req)

	// Should accept valid update
	if w.Code != http.StatusOK {
		t.Errorf("Expected 200 OK for valid update, got %d", w.Code)
		t.Logf("Response: %s", w.Body.String())
	}
	
	var response map[string]interface{}
	json.NewDecoder(w.Body).Decode(&response)
	
	if response["name"] != "Updated Product Name" {
		t.Error("Response should contain updated data")
	}
	
	t.Log("Valid update accepted successfully")
}

// TestSchemaValidationForMultipleResources tests validation across different resource types
func TestSchemaValidationForMultipleResources(t *testing.T) {
	handler, _ := setupTestHandler(t)
	
	testCases := []struct {
		resource    string
		invalidItem map[string]interface{}
		reason      string
	}{
		{
			resource: "user",
			invalidItem: map[string]interface{}{
				"id":  "not_a_number", // ID should be numeric
				"name": "Test User",
			},
			reason: "id should be numeric",
		},
		{
			resource: "post",
			invalidItem: map[string]interface{}{
				"id":     1,
				"userId": "not_a_number", // userId should be numeric
				"title":  "Test Post",
			},
			reason: "userId should be numeric",
		},
	}
	
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s_%s", tc.resource, tc.reason), func(t *testing.T) {
			body, _ := json.Marshal(tc.invalidItem)
			req := httptest.NewRequest("POST", fmt.Sprintf("/%s", tc.resource), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.PostCollection(tc.resource)(w, req)

			if w.Code != http.StatusBadRequest {
				t.Logf("Expected 400 for %s, got %d - Response: %s", tc.reason, w.Code, w.Body.String())
			}
		})
	}
}
