package handlers

import (
	"net/http"
	"strconv"

	"github.com/0xdps/api-mock/lib/schema"

	"github.com/gin-gonic/gin"
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
func (h *DynamicHandler) GetCollection(resourceName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get count parameter
		count := getCountParam(c, 10)

		// Generate data using the schema
		data, err := h.registry.GenerateData(resourceName, count)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, data)
	}
}

// GetSingle returns a single item for a resource
func (h *DynamicHandler) GetSingle(resourceName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the ID from URL
		id := c.Param("id")

		// Generate a single item
		data, err := h.registry.GenerateData(resourceName, 1)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		if len(data) == 0 {
			c.JSON(http.StatusNotFound, gin.H{
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

		c.JSON(http.StatusOK, item)
	}
}

// GetResourceMetadata returns metadata about a resource
func (h *DynamicHandler) GetResourceMetadata(resourceName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		schema, ok := h.registry.GetSchema(resourceName)
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Resource not found",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"name":        schema.Resource.Name,
			"singular":    schema.Resource.Singular,
			"description": schema.Resource.Description,
			"title":       schema.Title,
			"properties":  len(schema.Properties),
		})
	}
}

// getCountParam extracts and validates count parameter
func getCountParam(c *gin.Context, defaultCount int) int {
	countStr := c.Query("count")
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
