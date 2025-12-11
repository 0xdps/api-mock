package middleware

import (
	"context"
	"math"
	"net/http"
	"strconv"
	"strings"
)

// Context keys for query parameters
const (
	PaginationKey contextKey = "pagination"
	SortingKey    contextKey = "sorting"
	SearchKey     contextKey = "search"
)

// PaginationParams holds pagination information
type PaginationParams struct {
	Page   int // 1-based page number
	Limit  int // Items per page
	Offset int // Calculated offset
}

// SortingParams holds sorting information
type SortingParams struct {
	Field string // Field to sort by
	Order string // "asc" or "desc"
}

// SearchParams holds search information
type SearchParams struct {
	Query  string   // Search query
	Fields []string // Fields to search in (optional)
}

// PaginationMiddleware extracts pagination parameters from query string
// Supports: page, limit, offset
// Defaults: page=1, limit=10
func PaginationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		// Default values
		page := 1
		limit := 10
		offset := 0

		// Parse page
		if pageStr := query.Get("page"); pageStr != "" {
			if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
				page = p
			}
		}

		// Parse limit
		if limitStr := query.Get("limit"); limitStr != "" {
			if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
				// Cap at 100 to prevent abuse
				if l > 100 {
					l = 100
				}
				limit = l
			}
		}

		// Parse offset (takes precedence over page if both provided)
		if offsetStr := query.Get("offset"); offsetStr != "" {
			if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
				offset = o
			}
		} else {
			// Calculate offset from page
			offset = (page - 1) * limit
		}

		params := &PaginationParams{
			Page:   page,
			Limit:  limit,
			Offset: offset,
		}

		ctx := context.WithValue(r.Context(), PaginationKey, params)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// SortingMiddleware extracts sorting parameters from query string
// Supports: sort (field name), order (asc/desc)
func SortingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		params := &SortingParams{
			Field: query.Get("sort"),
			Order: "asc", // Default order
		}

		// Parse order
		if orderStr := query.Get("order"); orderStr != "" {
			orderStr = strings.ToLower(orderStr)
			if orderStr == "desc" || orderStr == "descending" {
				params.Order = "desc"
			} else if orderStr == "asc" || orderStr == "ascending" {
				params.Order = "asc"
			}
		}

		// Store in context even if no sort field (handlers can check)
		ctx := context.WithValue(r.Context(), SortingKey, params)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// SearchMiddleware extracts search parameters from query string
// Supports: q or search (query string), fields (comma-separated fields to search)
func SearchMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		// Check both 'q' and 'search' parameters
		searchQuery := query.Get("q")
		if searchQuery == "" {
			searchQuery = query.Get("search")
		}

		params := &SearchParams{
			Query: searchQuery,
		}

		// Parse search fields if provided
		if fieldsStr := query.Get("search_fields"); fieldsStr != "" {
			fields := strings.Split(fieldsStr, ",")
			for i, field := range fields {
				fields[i] = strings.TrimSpace(field)
			}
			params.Fields = fields
		}

		// Store in context even if no search query (handlers can check)
		ctx := context.WithValue(r.Context(), SearchKey, params)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// Helper functions to extract values from context

// GetPagination retrieves pagination params from context
func GetPagination(ctx context.Context) (*PaginationParams, bool) {
	params, ok := ctx.Value(PaginationKey).(*PaginationParams)
	return params, ok
}

// GetSorting retrieves sorting params from context
func GetSorting(ctx context.Context) (*SortingParams, bool) {
	params, ok := ctx.Value(SortingKey).(*SortingParams)
	return params, ok
}

// GetSearch retrieves search params from context
func GetSearch(ctx context.Context) (*SearchParams, bool) {
	params, ok := ctx.Value(SearchKey).(*SearchParams)
	return params, ok
}

// PaginationMeta holds metadata for paginated responses
type PaginationMeta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// NewPaginationMeta creates pagination metadata
func NewPaginationMeta(params *PaginationParams, total int) *PaginationMeta {
	totalPages := int(math.Ceil(float64(total) / float64(params.Limit)))
	if totalPages == 0 {
		totalPages = 1
	}

	return &PaginationMeta{
		Page:       params.Page,
		Limit:      params.Limit,
		Total:      total,
		TotalPages: totalPages,
	}
}

// PaginatedResponse wraps data with pagination metadata
type PaginatedResponse struct {
	Data       interface{}     `json:"data"`
	Pagination *PaginationMeta `json:"pagination"`
}

// NewPaginatedResponse creates a paginated response
func NewPaginatedResponse(data interface{}, params *PaginationParams, total int) *PaginatedResponse {
	return &PaginatedResponse{
		Data:       data,
		Pagination: NewPaginationMeta(params, total),
	}
}

// ApplyPagination applies pagination to a slice
func ApplyPagination(items []map[string]interface{}, params *PaginationParams) []map[string]interface{} {
	start := params.Offset
	end := start + params.Limit

	if start >= len(items) {
		return []map[string]interface{}{}
	}

	if end > len(items) {
		end = len(items)
	}

	return items[start:end]
}

// ApplySorting sorts a slice of maps by a field
func ApplySorting(items []map[string]interface{}, params *SortingParams) []map[string]interface{} {
	if params.Field == "" {
		return items
	}

	// Create a copy to avoid modifying original
	sorted := make([]map[string]interface{}, len(items))
	copy(sorted, items)

	// Simple bubble sort (for production, use more efficient sorting)
	for i := 0; i < len(sorted)-1; i++ {
		for j := 0; j < len(sorted)-i-1; j++ {
			val1, ok1 := sorted[j][params.Field]
			val2, ok2 := sorted[j+1][params.Field]

			if !ok1 || !ok2 {
				continue
			}

			shouldSwap := false
			if params.Order == "asc" {
				shouldSwap = compareValues(val1, val2) > 0
			} else {
				shouldSwap = compareValues(val1, val2) < 0
			}

			if shouldSwap {
				sorted[j], sorted[j+1] = sorted[j+1], sorted[j]
			}
		}
	}

	return sorted
}

// compareValues compares two values for sorting
func compareValues(a, b interface{}) int {
	// Try numeric comparison
	if aNum, aOk := toFloat64(a); aOk {
		if bNum, bOk := toFloat64(b); bOk {
			if aNum < bNum {
				return -1
			} else if aNum > bNum {
				return 1
			}
			return 0
		}
	}

	// Fall back to string comparison
	aStr := toString(a)
	bStr := toString(b)
	if aStr < bStr {
		return -1
	} else if aStr > bStr {
		return 1
	}
	return 0
}

// toFloat64 converts interface{} to float64 if possible
func toFloat64(v interface{}) (float64, bool) {
	switch val := v.(type) {
	case int:
		return float64(val), true
	case int32:
		return float64(val), true
	case int64:
		return float64(val), true
	case float32:
		return float64(val), true
	case float64:
		return val, true
	default:
		return 0, false
	}
}

// toString converts interface{} to string
func toString(v interface{}) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

// ApplySearch filters items based on search query
func ApplySearch(items []map[string]interface{}, params *SearchParams) []map[string]interface{} {
	if params.Query == "" {
		return items
	}

	query := strings.ToLower(params.Query)
	filtered := make([]map[string]interface{}, 0)

	for _, item := range items {
		if matchesSearch(item, query, params.Fields) {
			filtered = append(filtered, item)
		}
	}

	return filtered
}

// matchesSearch checks if an item matches the search query
func matchesSearch(item map[string]interface{}, query string, fields []string) bool {
	// If specific fields provided, search only in those
	if len(fields) > 0 {
		for _, field := range fields {
			if val, exists := item[field]; exists {
				if strings.Contains(strings.ToLower(toString(val)), query) {
					return true
				}
			}
		}
		return false
	}

	// Otherwise, search in all string fields
	for _, val := range item {
		if strings.Contains(strings.ToLower(toString(val)), query) {
			return true
		}
	}

	return false
}
