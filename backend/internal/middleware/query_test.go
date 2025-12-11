package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPaginationMiddleware(t *testing.T) {
	tests := []struct {
		name         string
		queryString  string
		expectedPage int
		expectedLimit int
		expectedOffset int
	}{
		{"Default values", "", 1, 10, 0},
		{"Custom page", "page=3", 3, 10, 20},
		{"Custom limit", "limit=25", 1, 25, 0},
		{"Page and limit", "page=2&limit=50", 2, 50, 50},
		{"Explicit offset", "offset=100", 1, 10, 100},
		{"Offset takes precedence", "page=3&offset=75", 3, 10, 75},
		{"Limit capped at 100", "limit=200", 1, 100, 0},
		{"Invalid page defaults to 1", "page=-5", 1, 10, 0},
		{"Invalid limit defaults to 10", "limit=invalid", 1, 10, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := PaginationMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				params, ok := GetPagination(r.Context())
				if !ok {
					t.Fatal("Expected pagination params in context")
				}
				if params.Page != tt.expectedPage {
					t.Errorf("Expected page %d, got %d", tt.expectedPage, params.Page)
				}
				if params.Limit != tt.expectedLimit {
					t.Errorf("Expected limit %d, got %d", tt.expectedLimit, params.Limit)
				}
				if params.Offset != tt.expectedOffset {
					t.Errorf("Expected offset %d, got %d", tt.expectedOffset, params.Offset)
				}
				w.WriteHeader(http.StatusOK)
			}))

			req := httptest.NewRequest("GET", "/?"+tt.queryString, nil)
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)
		})
	}
}

func TestSortingMiddleware(t *testing.T) {
	tests := []struct {
		name          string
		queryString   string
		expectedField string
		expectedOrder string
	}{
		{"No sorting", "", "", "asc"},
		{"Sort by name ascending", "sort=name", "name", "asc"},
		{"Sort by age descending", "sort=age&order=desc", "age", "desc"},
		{"Sort with descending variant", "sort=created_at&order=descending", "created_at", "desc"},
		{"Sort with ascending variant", "sort=updated_at&order=ascending", "updated_at", "asc"},
		{"Invalid order defaults to asc", "sort=name&order=invalid", "name", "asc"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := SortingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				params, ok := GetSorting(r.Context())
				if !ok {
					t.Fatal("Expected sorting params in context")
				}
				if params.Field != tt.expectedField {
					t.Errorf("Expected field %s, got %s", tt.expectedField, params.Field)
				}
				if params.Order != tt.expectedOrder {
					t.Errorf("Expected order %s, got %s", tt.expectedOrder, params.Order)
				}
				w.WriteHeader(http.StatusOK)
			}))

			req := httptest.NewRequest("GET", "/?"+tt.queryString, nil)
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)
		})
	}
}

func TestSearchMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		queryString    string
		expectedQuery  string
		expectedFields []string
	}{
		{"No search", "", "", nil},
		{"Search with q parameter", "q=john", "john", nil},
		{"Search with search parameter", "search=alice", "alice", nil},
		{"q takes precedence", "q=bob&search=alice", "bob", nil},
		{"Search with fields", "q=test&search_fields=name,email", "test", []string{"name", "email"}},
		{"Fields with spaces", "q=test&search_fields=id, name , email", "test", []string{"id", "name", "email"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := SearchMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				params, ok := GetSearch(r.Context())
				if !ok {
					t.Fatal("Expected search params in context")
				}
				if params.Query != tt.expectedQuery {
					t.Errorf("Expected query %s, got %s", tt.expectedQuery, params.Query)
				}
				if len(params.Fields) != len(tt.expectedFields) {
					t.Errorf("Expected %d fields, got %d", len(tt.expectedFields), len(params.Fields))
				}
				for i, field := range params.Fields {
					if i < len(tt.expectedFields) && field != tt.expectedFields[i] {
						t.Errorf("Expected field %s, got %s", tt.expectedFields[i], field)
					}
				}
				w.WriteHeader(http.StatusOK)
			}))

			req := httptest.NewRequest("GET", "/?"+tt.queryString, nil)
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)
		})
	}
}

func TestNewPaginationMeta(t *testing.T) {
	tests := []struct {
		name              string
		params            *PaginationParams
		total             int
		expectedPage      int
		expectedLimit     int
		expectedTotal     int
		expectedTotalPages int
	}{
		{"Basic pagination", &PaginationParams{Page: 1, Limit: 10, Offset: 0}, 50, 1, 10, 50, 5},
		{"Second page", &PaginationParams{Page: 2, Limit: 10, Offset: 10}, 50, 2, 10, 50, 5},
		{"Partial last page", &PaginationParams{Page: 1, Limit: 10, Offset: 0}, 25, 1, 10, 25, 3},
		{"Single page", &PaginationParams{Page: 1, Limit: 100, Offset: 0}, 50, 1, 100, 50, 1},
		{"Empty results", &PaginationParams{Page: 1, Limit: 10, Offset: 0}, 0, 1, 10, 0, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			meta := NewPaginationMeta(tt.params, tt.total)
			if meta.Page != tt.expectedPage {
				t.Errorf("Expected page %d, got %d", tt.expectedPage, meta.Page)
			}
			if meta.Limit != tt.expectedLimit {
				t.Errorf("Expected limit %d, got %d", tt.expectedLimit, meta.Limit)
			}
			if meta.Total != tt.expectedTotal {
				t.Errorf("Expected total %d, got %d", tt.expectedTotal, meta.Total)
			}
			if meta.TotalPages != tt.expectedTotalPages {
				t.Errorf("Expected total pages %d, got %d", tt.expectedTotalPages, meta.TotalPages)
			}
		})
	}
}

func TestApplyPagination(t *testing.T) {
	items := []map[string]interface{}{
		{"id": 1, "name": "Alice"},
		{"id": 2, "name": "Bob"},
		{"id": 3, "name": "Charlie"},
		{"id": 4, "name": "Diana"},
		{"id": 5, "name": "Eve"},
	}

	tests := []struct {
		name        string
		params      *PaginationParams
		expectedIDs []int
	}{
		{"First page", &PaginationParams{Page: 1, Limit: 2, Offset: 0}, []int{1, 2}},
		{"Second page", &PaginationParams{Page: 2, Limit: 2, Offset: 2}, []int{3, 4}},
		{"Last page partial", &PaginationParams{Page: 3, Limit: 2, Offset: 4}, []int{5}},
		{"Beyond last page", &PaginationParams{Page: 10, Limit: 2, Offset: 100}, []int{}},
		{"Larger limit", &PaginationParams{Page: 1, Limit: 10, Offset: 0}, []int{1, 2, 3, 4, 5}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ApplyPagination(items, tt.params)
			if len(result) != len(tt.expectedIDs) {
				t.Fatalf("Expected %d items, got %d", len(tt.expectedIDs), len(result))
			}
			for i, id := range tt.expectedIDs {
				if result[i]["id"] != id {
					t.Errorf("Item %d: expected id %d, got %v", i, id, result[i]["id"])
				}
			}
		})
	}
}

func TestApplySorting(t *testing.T) {
	items := []map[string]interface{}{
		{"id": 3, "name": "Charlie", "age": 30},
		{"id": 1, "name": "Alice", "age": 25},
		{"id": 2, "name": "Bob", "age": 35},
	}

	tests := []struct {
		name        string
		params      *SortingParams
		expectedIDs []int
	}{
		{"No sorting", &SortingParams{Field: "", Order: "asc"}, []int{3, 1, 2}},
		{"Sort by id ascending", &SortingParams{Field: "id", Order: "asc"}, []int{1, 2, 3}},
		{"Sort by id descending", &SortingParams{Field: "id", Order: "desc"}, []int{3, 2, 1}},
		{"Sort by name ascending", &SortingParams{Field: "name", Order: "asc"}, []int{1, 2, 3}},
		{"Sort by name descending", &SortingParams{Field: "name", Order: "desc"}, []int{3, 2, 1}},
		{"Sort by age ascending", &SortingParams{Field: "age", Order: "asc"}, []int{1, 3, 2}},
		{"Sort by age descending", &SortingParams{Field: "age", Order: "desc"}, []int{2, 3, 1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ApplySorting(items, tt.params)
			if len(result) != len(tt.expectedIDs) {
				t.Fatalf("Expected %d items, got %d", len(tt.expectedIDs), len(result))
			}
			for i, id := range tt.expectedIDs {
				if result[i]["id"] != id {
					t.Errorf("Item %d: expected id %d, got %v", i, id, result[i]["id"])
				}
			}
		})
	}
}

func TestApplySearch(t *testing.T) {
	items := []map[string]interface{}{
		{"id": 1, "name": "Alice Smith", "email": "alice@example.com", "city": "New York"},
		{"id": 2, "name": "Bob Johnson", "email": "bob@test.com", "city": "Los Angeles"},
		{"id": 3, "name": "Charlie Brown", "email": "charlie@example.com", "city": "Chicago"},
	}

	tests := []struct {
		name        string
		params      *SearchParams
		expectedIDs []int
	}{
		{"No search query", &SearchParams{Query: ""}, []int{1, 2, 3}},
		{"Search all fields - alice", &SearchParams{Query: "alice"}, []int{1}},
		{"Search all fields - example", &SearchParams{Query: "example"}, []int{1, 3}},
		{"Search specific field - name contains bob", &SearchParams{Query: "bob", Fields: []string{"name"}}, []int{2}},
		{"Search specific field - email contains test", &SearchParams{Query: "test", Fields: []string{"email"}}, []int{2}},
		{"Search multiple fields", &SearchParams{Query: "chicago", Fields: []string{"name", "city"}}, []int{3}},
		{"No matches", &SearchParams{Query: "nonexistent"}, []int{}},
		{"Case insensitive", &SearchParams{Query: "ALICE"}, []int{1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ApplySearch(items, tt.params)
			if len(result) != len(tt.expectedIDs) {
				t.Fatalf("Expected %d items, got %d", len(tt.expectedIDs), len(result))
			}
			for i, id := range tt.expectedIDs {
				if result[i]["id"] != id {
					t.Errorf("Item %d: expected id %d, got %v", i, id, result[i]["id"])
				}
			}
		})
	}
}

func TestCompareValues(t *testing.T) {
	tests := []struct {
		name     string
		a        interface{}
		b        interface{}
		expected int
	}{
		{"Int comparison - less", 5, 10, -1},
		{"Int comparison - greater", 10, 5, 1},
		{"Int comparison - equal", 5, 5, 0},
		{"Float comparison - less", 3.14, 6.28, -1},
		{"Float comparison - greater", 6.28, 3.14, 1},
		{"String comparison - less", "apple", "banana", -1},
		{"String comparison - greater", "banana", "apple", 1},
		{"String comparison - equal", "test", "test", 0},
		{"Mixed types - number vs string", 5, "apple", -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := compareValues(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestToFloat64(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected float64
		ok       bool
	}{
		{"Int", 42, 42.0, true},
		{"Int32", int32(42), 42.0, true},
		{"Int64", int64(42), 42.0, true},
		{"Float32", float32(3.14), 3.14, true},
		{"Float64", float64(3.14159), 3.14159, true},
		{"String", "42", 0.0, false},
		{"Bool", true, 0.0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := toFloat64(tt.value)
			if ok != tt.ok {
				t.Errorf("Expected ok=%v, got %v", tt.ok, ok)
			}
			if ok && result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestToString(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected string
	}{
		{"String", "hello", "hello"},
		{"Nil", nil, ""},
		{"Int", 42, ""},
		{"Bool", true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := toString(tt.value)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}
