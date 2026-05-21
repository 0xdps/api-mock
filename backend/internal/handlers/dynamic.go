package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/0xdps/api-mock/go/internal/cache"
	"github.com/0xdps/api-mock/go/internal/filters"
	"github.com/0xdps/api-mock/go/internal/middleware"
	"github.com/0xdps/api-mock/go/internal/schema"
	"github.com/go-chi/chi/v5"
)

// DynamicHandler handles requests for schema-driven resources
type DynamicHandler struct {
	registry *schema.Registry
	cache    *cache.Cache
}

// NewDynamicHandler creates a new dynamic handler
func NewDynamicHandler(registry *schema.Registry, c *cache.Cache) *DynamicHandler {
	return &DynamicHandler{
		registry: registry,
		cache:    c,
	}
}

// GetCollection returns a collection of items for a resource
func (h *DynamicHandler) GetCollection(resourceName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get context values
		ctx := r.Context()

		// Check if cache should be bypassed (check both context and request for robustness)
		skipCache := middleware.ShouldSkipCache(ctx) || shouldSkipCache(r)

		// locale overrides name/location/phone generators; bypass cache when set
		locale := r.URL.Query().Get("locale")
		if locale != "" {
			skipCache = true
		}

		// Get pagination params (defaults: page=1, limit=10)
		pagination, hasPagination := middleware.GetPagination(ctx)
		count := 10
		if hasPagination {
			count = pagination.Limit
		} else {
			count = getCountParam(r, 10)
		}

		// Get sorting params
		sorting, _ := middleware.GetSorting(ctx)

		// Get search params
		search, _ := middleware.GetSearch(ctx)

		// Parse filters from query parameters
		filterList := filters.ParseFilters(r.URL.Query())

		var data []map[string]interface{}
		var err error

		// Generate more data to handle pagination properly
		// For production, you'd fetch from DB with proper LIMIT/OFFSET
		generateCount := count * 10 // Generate extra for pagination/search
		if hasPagination && pagination.Page > 1 {
			generateCount = pagination.Offset + count
		}

		if skipCache {
			// Generate fresh data (bypass cache), applying locale overrides if requested
			if locale != "" {
				data, err = h.registry.GenerateDataWithLocale(resourceName, generateCount, locale)
			} else {
				data, err = h.registry.GenerateData(resourceName, generateCount)
			}
			if err != nil {
				respondJSON(w, http.StatusInternalServerError, map[string]string{
					"error": err.Error(),
				})
				return
			}
			w.Header().Set("X-Cache", "BYPASS")
		} else {
			// Try to get data from cache
			var found bool
			data, found = h.cache.Get(resourceName, generateCount)
			if !found {
				// Fallback: generate data on the fly (shouldn't happen after warmup)
				data, err = h.registry.GenerateData(resourceName, generateCount)
				if err != nil {
					respondJSON(w, http.StatusInternalServerError, map[string]string{
						"error": err.Error(),
					})
					return
				}
				w.Header().Set("X-Cache", "MISS")
			} else {
				w.Header().Set("X-Cache", "HIT")
			}
		}

		// Apply search if present
		if search != nil && search.Query != "" {
			data = middleware.ApplySearch(data, search)
		}

		// Apply filters if present
		if len(filterList) > 0 {
			data = filters.ApplyFilters(data, filterList)
		}

		// Apply sorting if present
		if sorting != nil && sorting.Field != "" {
			data = middleware.ApplySorting(data, sorting)
		}

		// Store total after filtering for pagination
		total := len(data)

		// Apply field filtering if present
		if fields, ok := middleware.GetFields(ctx); ok {
			data = middleware.FilterFieldsSlice(data, fields)
		}

		// Apply pagination if present
		if hasPagination {
			data = middleware.ApplyPagination(data, pagination)

			// Return paginated response with metadata
			response := middleware.NewPaginatedResponse(data, pagination, total)
			respondJSON(w, http.StatusOK, response)
			return
		}

		// Return simple response if no pagination (limit to requested count)
		if len(data) > count {
			data = data[:count]
		}
		respondJSON(w, http.StatusOK, data)
	}
}

// GetSingle returns a single item for a resource
func (h *DynamicHandler) GetSingle(resourceName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the ID from URL
		id := chi.URLParam(r, "id")

		// Check if cache should be bypassed
		skipCache := shouldSkipCache(r)

		var item map[string]interface{}

		if skipCache {
			// Generate fresh data (bypass cache)
			data, err := h.registry.GenerateData(resourceName, 1)
			if err != nil {
				respondJSON(w, http.StatusInternalServerError, map[string]string{
					"error": err.Error(),
				})
				return
			}
			if len(data) == 0 {
				respondJSON(w, http.StatusNotFound, map[string]string{
					"error": "Resource not found",
				})
				return
			}
			item = data[0]
			w.Header().Set("X-Cache", "BYPASS")
		} else {
			// Try to get item by ID from cache
			var found bool
			item, found = h.cache.GetByID(resourceName, id)
			if found {
				w.Header().Set("X-Cache", "HIT")
				// Set the ID to the requested ID before returning
				if idNum, err := strconv.Atoi(id); err == nil {
					item["id"] = idNum
				} else {
					item["id"] = id
				}
				respondJSON(w, http.StatusOK, item)
				return
			}

			// If not found in cache, get first item and set the requested ID
			data, found := h.cache.Get(resourceName, 1)
			if !found || len(data) == 0 {
				// Fallback: generate on the fly
				var err error
				data, err = h.registry.GenerateData(resourceName, 1)
				if err != nil {
					respondJSON(w, http.StatusInternalServerError, map[string]string{
						"error": err.Error(),
					})
					return
				}
				w.Header().Set("X-Cache", "MISS")
			} else {
				w.Header().Set("X-Cache", "PARTIAL")
			}

			if len(data) == 0 {
				respondJSON(w, http.StatusNotFound, map[string]string{
					"error": "Resource not found",
				})
				return
			}

			item = data[0]
		}

		// Set the ID to the requested ID
		if idNum, err := strconv.Atoi(id); err == nil {
			item["id"] = idNum
		} else {
			item["id"] = id
		}

		respondJSON(w, http.StatusOK, item)
	}
}

// GetResourceMetadata returns metadata about a resource
func (h *DynamicHandler) GetResourceMetadata(resourceName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Try to get from cache first
		meta, found := h.cache.GetMeta(resourceName)
		if found {
			// Set cache headers for 1 hour
			w.Header().Set("Cache-Control", "public, max-age=3600")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(meta)
			return
		}

		// Fallback: generate on the fly if not in cache (shouldn't happen after warmup)
		schema, ok := h.registry.GetSchema(resourceName)
		if !ok {
			respondJSON(w, http.StatusNotFound, map[string]string{
				"error": "Resource not found",
			})
			return
		}

		// Use root-level description if available, otherwise fall back to x-resource description
		description := schema.Description
		if description == "" {
			description = schema.Resource.Description
		}

		// Return full schema with metadata
		metaResponse := map[string]interface{}{
			"$schema":        schema.SchemaURI,
			"title":          schema.Title,
			"type":           schema.Type,
			"description":    description,
			"name":           schema.Resource.Name,
			"singular":       schema.Resource.Singular,
			"group":          schema.Resource.Group,
			"properties":     schema.Properties,
			"required":       schema.Required,
			"property_count": len(schema.Properties),
		}

		// Set cache headers for 1 hour
		w.Header().Set("Cache-Control", "public, max-age=3600")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(metaResponse)
	}
}

// GetGroupInfo returns metadata about a group (no data)
func (h *DynamicHandler) GetGroupInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		groupName := chi.URLParam(r, "group")

		resourceNames := h.registry.GetResourceNamesByGroup(groupName)
		if len(resourceNames) == 0 {
			respondJSON(w, http.StatusNotFound, map[string]string{
				"error": "Group not found",
			})
			return
		}

		respondJSON(w, http.StatusOK, map[string]interface{}{
			"group":     groupName,
			"resources": resourceNames,
			"count":     len(resourceNames),
		})
	}
}

// GetGroupResourceCollection returns collection data for a resource within a group
func (h *DynamicHandler) GetGroupResourceCollection() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		groupName := chi.URLParam(r, "group")
		resourceName := chi.URLParam(r, "resource")

		// Verify the resource belongs to this group
		groupResources := h.registry.GetResourceNamesByGroup(groupName)
		found := false
		for _, res := range groupResources {
			if res == resourceName {
				found = true
				break
			}
		}

		if !found {
			respondJSON(w, http.StatusNotFound, map[string]string{
				"error": "Resource not found in this group",
			})
			return
		}

		// Get count parameter
		count := getCountParam(r, 10)

		// Check if cache should be bypassed
		skipCache := shouldSkipCache(r)

		var data []map[string]interface{}
		var err error

		if skipCache {
			// Generate fresh data (bypass cache)
			data, err = h.registry.GenerateData(resourceName, count)
			if err != nil {
				respondJSON(w, http.StatusInternalServerError, map[string]string{
					"error": err.Error(),
				})
				return
			}
			w.Header().Set("X-Cache", "BYPASS")
		} else {
			// Get data from cache
			var found bool
			data, found = h.cache.Get(resourceName, count)
			if !found {
				// Fallback: generate data on the fly
				data, err = h.registry.GenerateData(resourceName, count)
				if err != nil {
					respondJSON(w, http.StatusInternalServerError, map[string]string{
						"error": err.Error(),
					})
					return
				}
				w.Header().Set("X-Cache", "MISS")
			} else {
				w.Header().Set("X-Cache", "HIT")
			}
		}

		respondJSON(w, http.StatusOK, data)
	}
}

// GetGroupResourceSingle returns a single item for a resource within a group
func (h *DynamicHandler) GetGroupResourceSingle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		groupName := chi.URLParam(r, "group")
		resourceName := chi.URLParam(r, "resource")
		id := chi.URLParam(r, "id")

		// Verify the resource belongs to this group
		groupResources := h.registry.GetResourceNamesByGroup(groupName)
		found := false
		for _, res := range groupResources {
			if res == resourceName {
				found = true
				break
			}
		}

		if !found {
			respondJSON(w, http.StatusNotFound, map[string]string{
				"error": "Resource not found in this group",
			})
			return
		}

		// Check if cache should be bypassed
		skipCache := shouldSkipCache(r)

		var item map[string]interface{}

		if skipCache {
			// Generate fresh data (bypass cache)
			data, err := h.registry.GenerateData(resourceName, 1)
			if err != nil {
				respondJSON(w, http.StatusInternalServerError, map[string]string{
					"error": err.Error(),
				})
				return
			}
			if len(data) == 0 {
				respondJSON(w, http.StatusNotFound, map[string]string{
					"error": "Resource not found",
				})
				return
			}
			item = data[0]
			w.Header().Set("X-Cache", "BYPASS")
		} else {
			// Try to get item by ID from cache
			var found bool
			item, found = h.cache.GetByID(resourceName, id)
			if found {
				w.Header().Set("X-Cache", "HIT")
				// Set the ID to the requested ID before returning
				if idNum, err := strconv.Atoi(id); err == nil {
					item["id"] = idNum
				} else {
					item["id"] = id
				}
				respondJSON(w, http.StatusOK, item)
				return
			}

			// If not found in cache, get first item and set the requested ID
			data, found := h.cache.Get(resourceName, 1)
			if !found || len(data) == 0 {
				// Fallback: generate on the fly
				var err error
				data, err = h.registry.GenerateData(resourceName, 1)
				if err != nil {
					respondJSON(w, http.StatusInternalServerError, map[string]string{
						"error": err.Error(),
					})
					return
				}
				w.Header().Set("X-Cache", "MISS")
			} else {
				w.Header().Set("X-Cache", "PARTIAL")
			}

			if len(data) == 0 {
				respondJSON(w, http.StatusNotFound, map[string]string{
					"error": "Resource not found",
				})
				return
			}

			item = data[0]
		}

		// Set the ID to the requested ID
		if idNum, err := strconv.Atoi(id); err == nil {
			item["id"] = idNum
		} else {
			item["id"] = id
		}

		respondJSON(w, http.StatusOK, item)
	}
}

// PostCollection creates a new item in a resource collection
func (h *DynamicHandler) PostCollection(resourceName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newItem map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&newItem); err != nil {
			respondJSON(w, http.StatusBadRequest, map[string]string{
				"error": "Invalid request body",
			})
			return
		}

		// Reject if user tries to provide an ID (IDs are auto-generated)
		if _, hasID := newItem["id"]; hasID {
			respondJSON(w, http.StatusBadRequest, map[string]string{
				"error": "Cannot specify 'id' field in POST request. IDs are auto-generated.",
			})
			return
		}

		// Validate against schema (skip 'id' field as it will be auto-generated)
		schema, ok := h.registry.GetSchema(resourceName)
		if ok && schema != nil {
			if err := schema.ValidateItem(newItem, "id"); err != nil {
				respondJSON(w, http.StatusBadRequest, map[string]string{
					"error": fmt.Sprintf("Validation failed: %s", err.Error()),
				})
				return
			}
		}

		// Auto-generate ID
		newItem["id"] = int(time.Now().UnixNano() / 1000000) // milliseconds

		// Add item to cache (includes constraint validation)
		if err := h.cache.AddItem(resourceName, newItem); err != nil {
			respondJSON(w, http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
			return
		}

		respondJSON(w, http.StatusCreated, newItem)
	}
}

// PutSingle updates an item in a resource
func (h *DynamicHandler) PutSingle(resourceName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		var updates map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
			respondJSON(w, http.StatusBadRequest, map[string]string{
				"error": "Invalid request body",
			})
			return
		}

		// Check if updates is nil or empty (e.g., body was "null" or {})
		if len(updates) == 0 {
			respondJSON(w, http.StatusBadRequest, map[string]string{
				"error": "Request body must contain fields to update",
			})
			return
		}

		// Get existing item to merge with updates for validation
		existingItem, found := h.cache.GetByID(resourceName, id)
		if !found {
			respondJSON(w, http.StatusNotFound, map[string]string{
				"error": "Item not found",
			})
			return
		}

		// Merge updates with existing item for full validation
		mergedItem := make(map[string]interface{})
		for k, v := range existingItem {
			mergedItem[k] = v
		}
		for k, v := range updates {
			mergedItem[k] = v
		}

		// Validate merged item against schema
		schema, ok := h.registry.GetSchema(resourceName)
		if ok && schema != nil {
			if err := schema.ValidateItem(mergedItem); err != nil {
				respondJSON(w, http.StatusBadRequest, map[string]string{
					"error": fmt.Sprintf("Validation failed: %s", err.Error()),
				})
				return
			}
		}

		// Update item in cache
		if err := h.cache.UpdateItemByID(resourceName, id, updates); err != nil {
			respondJSON(w, http.StatusNotFound, map[string]string{
				"error": err.Error(),
			})
			return
		}

		// Get updated item
		item, found := h.cache.GetByID(resourceName, id)
		if !found {
			respondJSON(w, http.StatusNotFound, map[string]string{
				"error": "Item not found after update",
			})
			return
		}

		respondJSON(w, http.StatusOK, item)
	}
}

// DeleteSingle deletes an item from a resource
func (h *DynamicHandler) DeleteSingle(resourceName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		// Delete item from cache (includes constraint validation)
		if err := h.cache.DeleteItemByID(resourceName, id); err != nil {
			// Check if error is "item not found"
			if err.Error() == "item not found" {
				respondJSON(w, http.StatusNotFound, map[string]string{
					"error": err.Error(),
				})
				return
			}
			respondJSON(w, http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
			return
		}

		respondJSON(w, http.StatusNoContent, nil)
	}
}

// getCountParam extracts and validates count parameter
func getCountParam(r *http.Request, defaultCount int) int {
	countStr := r.URL.Query().Get("count")
	if countStr == "" {
		return defaultCount
	}

	count, err := strconv.Atoi(countStr)
	if err != nil || count <= 0 {
		return defaultCount
	}

	// Cap at 100
	if count > 100 {
		return 100
	}

	return count
}

// shouldSkipCache checks if cache should be bypassed for this request
func shouldSkipCache(r *http.Request) bool {
	// Check query parameter: ?nocache=true or ?fresh=true
	if r.URL.Query().Get("nocache") == "true" || r.URL.Query().Get("fresh") == "true" {
		return true
	}

	// Check headers: X-No-Cache: true or Cache-Control: no-cache
	if r.Header.Get("X-No-Cache") == "true" {
		return true
	}

	cacheControl := r.Header.Get("Cache-Control")
	if cacheControl == "no-cache" || cacheControl == "no-store" {
		return true
	}

	return false
}

// respondJSON writes JSON response
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
