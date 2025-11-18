package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/0xdps/api-mock/go/internal/schema"
	"github.com/go-chi/chi/v5"
)

// DynamicHandler handles requests for schema-driven resources
type DynamicHandler struct {
	registry *schema.Registry
}

// NewDynamicHandler creates a new dynamic handler
func NewDynamicHandler(registry *schema.Registry) *DynamicHandler {
	return &DynamicHandler{
		registry: registry,
	}
}

// GetCollection returns a collection of items for a resource
func (h *DynamicHandler) GetCollection(resourceName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get count parameter
		count := getCountParam(r, 10)

		// Generate data using the schema
		data, err := h.registry.GenerateData(resourceName, count)
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
			return
		}

		respondJSON(w, http.StatusOK, data)
	}
}

// GetSingle returns a single item for a resource
func (h *DynamicHandler) GetSingle(resourceName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the ID from URL
		id := chi.URLParam(r, "id")

		// Generate a single item
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

		// Set the ID to the requested ID
		item := data[0]
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
		schema, ok := h.registry.GetSchema(resourceName)
		if !ok {
			respondJSON(w, http.StatusNotFound, map[string]string{
				"error": "Resource not found",
			})
			return
		}

		respondJSON(w, http.StatusOK, map[string]interface{}{
			"name":        schema.Resource.Name,
			"singular":    schema.Resource.Singular,
			"description": schema.Resource.Description,
			"title":       schema.Title,
			"properties":  len(schema.Properties),
		})
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

// respondJSON writes JSON response
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
