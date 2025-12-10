package filters

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestParseFilters_Equal(t *testing.T) {
	queryParams := map[string][]string{
		"name": {"John"},
		"age":  {"30"},
	}

	filters := ParseFilters(queryParams)

	if len(filters) != 2 {
		t.Errorf("Expected 2 filters, got %d", len(filters))
	}

	nameFilter := findFilter(filters, "name")
	if nameFilter == nil {
		t.Fatal("Expected to find 'name' filter")
	}
	if nameFilter.Operator != OpEqual {
		t.Errorf("Expected operator '=', got '%s'", nameFilter.Operator)
	}
	if nameFilter.Value != "John" {
		t.Errorf("Expected value 'John', got '%v'", nameFilter.Value)
	}
}

func TestParseFilters_NotEqual(t *testing.T) {
	queryParams := map[string][]string{
		"status!=": {"inactive"},
	}

	filters := ParseFilters(queryParams)

	if len(filters) != 1 {
		t.Fatalf("Expected 1 filter, got %d", len(filters))
	}

	filter := filters[0]
	if filter.Field != "status" {
		t.Errorf("Expected field 'status', got '%s'", filter.Field)
	}
	if filter.Operator != OpNotEqual {
		t.Errorf("Expected operator '!=', got '%s'", filter.Operator)
	}
}

func TestParseFilters_NumericOperators(t *testing.T) {
	queryParams := map[string][]string{
		"age>":    {"18"},
		"price<":  {"100"},
		"year>=":  {"2020"},
		"score<=": {"50"},
	}

	filters := ParseFilters(queryParams)

	if len(filters) != 4 {
		t.Fatalf("Expected 4 filters, got %d", len(filters))
	}

	ageFilter := findFilter(filters, "age")
	if ageFilter == nil || ageFilter.Operator != OpGreater {
		t.Error("Expected age filter with '>' operator")
	}

	priceFilter := findFilter(filters, "price")
	if priceFilter == nil || priceFilter.Operator != OpLess {
		t.Error("Expected price filter with '<' operator")
	}
}

func TestParseFilters_StringOperators(t *testing.T) {
	queryParams := map[string][]string{
		"email~contains":  {"@gmail"},
		"name~startsWith": {"John"},
		"domain~endsWith": {".com"},
	}

	filters := ParseFilters(queryParams)

	if len(filters) != 3 {
		t.Fatalf("Expected 3 filters, got %d", len(filters))
	}

	emailFilter := findFilter(filters, "email")
	if emailFilter == nil || emailFilter.Operator != OpContains {
		t.Error("Expected email filter with 'contains' operator")
	}
}

func TestParseFilters_SkipsReservedParams(t *testing.T) {
	queryParams := map[string][]string{
		"name":    {"John"},
		"count":   {"10"},
		"nocache": {"true"},
		"fresh":   {"true"},
	}

	filters := ParseFilters(queryParams)

	if len(filters) != 1 {
		t.Errorf("Expected 1 filter, got %d", len(filters))
	}

	if filters[0].Field != "name" {
		t.Errorf("Expected field 'name', got '%s'", filters[0].Field)
	}
}

func TestApplyFilters_NoFilters(t *testing.T) {
	items := []map[string]interface{}{
		{"id": 1, "name": "John"},
		{"id": 2, "name": "Jane"},
	}

	result := ApplyFilters(items, []Filter{})

	if len(result) != 2 {
		t.Errorf("Expected 2 items, got %d", len(result))
	}
}

func TestApplyFilters_Equal(t *testing.T) {
	items := []map[string]interface{}{
		{"id": 1, "name": "John", "age": 30},
		{"id": 2, "name": "Jane", "age": 25},
		{"id": 3, "name": "John", "age": 35},
	}

	filters := []Filter{
		{Field: "name", Operator: OpEqual, Value: "John"},
	}

	result := ApplyFilters(items, filters)

	if len(result) != 2 {
		t.Errorf("Expected 2 items, got %d", len(result))
	}

	for _, item := range result {
		if item["name"] != "John" {
			t.Errorf("Expected name 'John', got '%v'", item["name"])
		}
	}
}

func TestApplyFilters_NotEqual(t *testing.T) {
	items := []map[string]interface{}{
		{"id": 1, "status": "active"},
		{"id": 2, "status": "inactive"},
		{"id": 3, "status": "active"},
	}

	filters := []Filter{
		{Field: "status", Operator: OpNotEqual, Value: "inactive"},
	}

	result := ApplyFilters(items, filters)

	if len(result) != 2 {
		t.Errorf("Expected 2 items, got %d", len(result))
	}
}

func TestApplyFilters_NumericGreater(t *testing.T) {
	items := []map[string]interface{}{
		{"id": 1, "age": 18},
		{"id": 2, "age": 25},
		{"id": 3, "age": 30},
		{"id": 4, "age": 15},
	}

	filters := []Filter{
		{Field: "age", Operator: OpGreater, Value: "20"},
	}

	result := ApplyFilters(items, filters)

	if len(result) != 2 {
		t.Errorf("Expected 2 items, got %d", len(result))
	}
}

func TestApplyFilters_NumericLess(t *testing.T) {
	items := []map[string]interface{}{
		{"id": 1, "price": 50.0},
		{"id": 2, "price": 100.0},
		{"id": 3, "price": 75.0},
		{"id": 4, "price": 120.0},
	}

	filters := []Filter{
		{Field: "price", Operator: OpLess, Value: "100"},
	}

	result := ApplyFilters(items, filters)

	if len(result) != 2 {
		t.Errorf("Expected 2 items, got %d", len(result))
	}
}

func TestApplyFilters_Contains(t *testing.T) {
	items := []map[string]interface{}{
		{"id": 1, "email": "john@gmail.com"},
		{"id": 2, "email": "jane@yahoo.com"},
		{"id": 3, "email": "bob@gmail.com"},
	}

	filters := []Filter{
		{Field: "email", Operator: OpContains, Value: "gmail"},
	}

	result := ApplyFilters(items, filters)

	if len(result) != 2 {
		t.Errorf("Expected 2 items, got %d", len(result))
	}
}

func TestApplyFilters_StartsWith(t *testing.T) {
	items := []map[string]interface{}{
		{"id": 1, "name": "John Doe"},
		{"id": 2, "name": "Jane Smith"},
		{"id": 3, "name": "John Adams"},
	}

	filters := []Filter{
		{Field: "name", Operator: OpStartsWith, Value: "John"},
	}

	result := ApplyFilters(items, filters)

	if len(result) != 2 {
		t.Errorf("Expected 2 items, got %d", len(result))
	}
}

func TestApplyFilters_EndsWith(t *testing.T) {
	items := []map[string]interface{}{
		{"id": 1, "domain": "example.com"},
		{"id": 2, "domain": "test.org"},
		{"id": 3, "domain": "demo.com"},
	}

	filters := []Filter{
		{Field: "domain", Operator: OpEndsWith, Value: ".com"},
	}

	result := ApplyFilters(items, filters)

	if len(result) != 2 {
		t.Errorf("Expected 2 items, got %d", len(result))
	}
}

func TestApplyFilters_MultipleFilters(t *testing.T) {
	items := []map[string]interface{}{
		{"id": 1, "name": "John", "age": 25, "city": "NYC"},
		{"id": 2, "name": "Jane", "age": 30, "city": "LA"},
		{"id": 3, "name": "John", "age": 35, "city": "NYC"},
		{"id": 4, "name": "Bob", "age": 28, "city": "NYC"},
	}

	filters := []Filter{
		{Field: "city", Operator: OpEqual, Value: "NYC"},
		{Field: "age", Operator: OpGreater, Value: "26"},
	}

	result := ApplyFilters(items, filters)

	if len(result) != 2 {
		t.Errorf("Expected 2 items, got %d", len(result))
	}
}

func TestApplyFilters_NoMatches(t *testing.T) {
	items := []map[string]interface{}{
		{"id": 1, "name": "John"},
		{"id": 2, "name": "Jane"},
	}

	filters := []Filter{
		{Field: "name", Operator: OpEqual, Value: "Bob"},
	}

	result := ApplyFilters(items, filters)

	if len(result) != 0 {
		t.Errorf("Expected 0 items, got %d", len(result))
	}
}

func TestApplyFilters_MissingField(t *testing.T) {
	items := []map[string]interface{}{
		{"id": 1, "name": "John"},
		{"id": 2, "name": "Jane"},
	}

	filters := []Filter{
		{Field: "age", Operator: OpEqual, Value: "30"},
	}

	result := ApplyFilters(items, filters)

	if len(result) != 0 {
		t.Errorf("Expected 0 items, got %d", len(result))
	}
}

func TestCompareEqual_DifferentTypes(t *testing.T) {
	testCases := []struct {
		a        interface{}
		b        interface{}
		expected bool
	}{
		{30, "30", true},
		{"hello", "hello", true},
		{30, "31", false},
	}

	for _, tc := range testCases {
		result := compareEqual(tc.a, tc.b)
		if result != tc.expected {
			t.Errorf("compareEqual(%v, %v) = %v, expected %v", tc.a, tc.b, result, tc.expected)
		}
	}
}

func TestToFloat_Conversions(t *testing.T) {
	testCases := []struct {
		input    interface{}
		expected float64
	}{
		{42, 42.0},
		{int64(42), 42.0},
		{float32(42.5), 42.5},
		{float64(42.5), 42.5},
		{"42.5", 42.5},
		{"invalid", 0.0},
	}

	for _, tc := range testCases {
		result := toFloat(tc.input)
		if result != tc.expected {
			t.Errorf("toFloat(%v) = %v, expected %v", tc.input, result, tc.expected)
		}
	}
}

func TestMatchesFilter_CaseInsensitive(t *testing.T) {
	item := map[string]interface{}{
		"email": "John@Example.COM",
	}

	filter := Filter{
		Field:    "email",
		Operator: OpContains,
		Value:    "example",
	}
	if !matchesFilter(item, filter) {
		t.Error("Expected case-insensitive contains to match")
	}

	filter = Filter{
		Field:    "email",
		Operator: OpStartsWith,
		Value:    "john",
	}
	if !matchesFilter(item, filter) {
		t.Error("Expected case-insensitive startsWith to match")
	}
}

func findFilter(filters []Filter, field string) *Filter {
	for _, f := range filters {
		if f.Field == field {
			return &f
		}
	}
	return nil
}

// ============================================================================
// EDGE CASE TESTS
// ============================================================================

func TestParseFilters_EmptyValues(t *testing.T) {
	queryParams := map[string][]string{
		"name":  {""},
		"age>":  {""},
		"email": {},
	}

	filters := ParseFilters(queryParams)

	// Should handle empty values gracefully
	for _, filter := range filters {
		if filter.Value == "" {
			t.Logf("Filter with empty value: %s %s '%v'", filter.Field, filter.Operator, filter.Value)
		}
	}
}

func TestParseFilters_SpecialCharacters(t *testing.T) {
	queryParams := map[string][]string{
		"field_with_underscore": {"value"},
		"field-with-dash":       {"value"},
		"field.with.dots":       {"value"},
		"field[with]brackets":   {"value"},
		"field/with/slash":      {"value"},
	}

	filters := ParseFilters(queryParams)

	// Should parse fields with special characters
	if len(filters) == 0 {
		t.Error("Expected to parse fields with special characters")
	}

	t.Logf("Parsed %d filters with special characters", len(filters))
}

func TestParseFilters_MalformedOperators(t *testing.T) {
	queryParams := map[string][]string{
		"field>>":   {"10"}, // Double operator
		"field<<<":  {"10"}, // Triple operator
		"field~":    {"test"}, // Incomplete string operator
		"field~xyz": {"test"}, // Invalid string operator
	}

	filters := ParseFilters(queryParams)

	// Should handle malformed operators gracefully (may treat as field names)
	for _, filter := range filters {
		t.Logf("Malformed: Field='%s', Operator='%s', Value='%v'", filter.Field, filter.Operator, filter.Value)
	}
}

func TestParseFilters_UnicodeCharacters(t *testing.T) {
	queryParams := map[string][]string{
		"name":    {"José"},
		"city":    {"北京"},
		"country": {"🇺🇸"},
	}

	filters := ParseFilters(queryParams)

	if len(filters) != 3 {
		t.Errorf("Expected 3 filters, got %d", len(filters))
	}

	for _, filter := range filters {
		if filter.Value == "" {
			t.Errorf("Unicode value lost for field %s", filter.Field)
		}
	}
}

func TestApplyFilters_NilValues(t *testing.T) {
	items := []map[string]interface{}{
		{"id": 1, "name": "John", "email": nil},
		{"id": 2, "name": nil, "email": "jane@example.com"},
		{"id": 3, "name": "Bob", "email": "bob@example.com"},
	}

	filters := []Filter{
		{Field: "email", Operator: OpEqual, Value: "jane@example.com"},
	}

	result := ApplyFilters(items, filters)

	if len(result) != 1 {
		t.Errorf("Expected 1 item, got %d", len(result))
	}
}

func TestApplyFilters_NestedObjects(t *testing.T) {
	items := []map[string]interface{}{
		{
			"id":   1,
			"user": map[string]interface{}{"name": "John", "age": 30},
		},
		{
			"id":   2,
			"user": map[string]interface{}{"name": "Jane", "age": 25},
		},
	}

	// Filters don't support nested fields currently
	filters := []Filter{
		{Field: "user", Operator: OpEqual, Value: "something"},
	}

	result := ApplyFilters(items, filters)

	// Should handle gracefully (likely return empty or handle as string comparison)
	t.Logf("Nested object filter result: %d items", len(result))
}

func TestApplyFilters_LargeDataset(t *testing.T) {
	// Create large dataset
	items := make([]map[string]interface{}, 10000)
	for i := 0; i < 10000; i++ {
		items[i] = map[string]interface{}{
			"id":   i,
			"age":  20 + (i % 50),
			"name": fmt.Sprintf("User%d", i),
		}
	}

	filters := []Filter{
		{Field: "age", Operator: OpGreater, Value: "30"},
		{Field: "age", Operator: OpLess, Value: "40"},
	}

	start := time.Now()
	result := ApplyFilters(items, filters)
	duration := time.Since(start)

	expectedCount := 0
	for i := 0; i < 10000; i++ {
		age := 20 + (i % 50)
		if age > 30 && age < 40 {
			expectedCount++
		}
	}

	if len(result) != expectedCount {
		t.Errorf("Expected %d items, got %d", expectedCount, len(result))
	}

	t.Logf("Filtered 10000 items to %d in %v", len(result), duration)

	// Performance check - should complete in reasonable time
	if duration > 100*time.Millisecond {
		t.Logf("Warning: Filter performance may need optimization (took %v)", duration)
	}
}

func TestApplyFilters_MultipleComplexFilters(t *testing.T) {
	items := []map[string]interface{}{
		{"id": 1, "name": "John Smith", "age": 30, "email": "john@gmail.com", "city": "NYC"},
		{"id": 2, "name": "Jane Doe", "age": 25, "email": "jane@yahoo.com", "city": "LA"},
		{"id": 3, "name": "John Adams", "age": 35, "email": "john@gmail.com", "city": "NYC"},
		{"id": 4, "name": "Bob Smith", "age": 28, "email": "bob@gmail.com", "city": "SF"},
	}

	filters := []Filter{
		{Field: "name", Operator: OpContains, Value: "John"},
		{Field: "age", Operator: OpGreaterEqual, Value: "30"},
		{Field: "email", Operator: OpContains, Value: "gmail"},
		{Field: "city", Operator: OpEqual, Value: "NYC"},
	}

	result := ApplyFilters(items, filters)

	// Should match: John Smith (30, NYC) and John Adams (35, NYC)
	if len(result) != 2 {
		t.Errorf("Expected 2 items, got %d", len(result))
	}

	for _, item := range result {
		if !strings.Contains(item["name"].(string), "John") {
			t.Error("Result should contain 'John' in name")
		}
		if item["city"] != "NYC" {
			t.Error("Result should be from NYC")
		}
	}
}

func TestApplyFilters_EmptyFilterValue(t *testing.T) {
	items := []map[string]interface{}{
		{"id": 1, "name": "John"},
		{"id": 2, "name": ""},
		{"id": 3, "name": "Bob"},
	}

	filters := []Filter{
		{Field: "name", Operator: OpEqual, Value: ""},
	}

	result := ApplyFilters(items, filters)

	// Should match item with empty name
	if len(result) != 1 {
		t.Errorf("Expected 1 item with empty name, got %d", len(result))
	}
}

func TestApplyFilters_NumericStringMixedComparison(t *testing.T) {
	items := []map[string]interface{}{
		{"id": 1, "count": "100"},  // String
		{"id": 2, "count": 200},    // Int
		{"id": 3, "count": 50.5},   // Float
		{"id": 4, "count": "75"},   // String
	}

	filters := []Filter{
		{Field: "count", Operator: OpGreater, Value: "80"},
	}

	result := ApplyFilters(items, filters)

	// Should match: "100", 200, and potentially 50.5 depending on string conversion
	if len(result) < 2 {
		t.Errorf("Expected at least 2 items, got %d", len(result))
	}

	t.Logf("Numeric string comparison result: %d items", len(result))
}

func TestCompareEqual_EdgeCases(t *testing.T) {
	testCases := []struct {
		name     string
		a        interface{}
		b        interface{}
		expected bool
	}{
		{"nil vs nil", nil, nil, true},
		{"nil vs string", nil, "test", false},
		{"number vs string number", 42, "42", true},
		{"float vs int", 42.0, 42, true},
		{"empty string vs empty string", "", "", true},
		{"boolean vs string", true, "true", true},
		{"array (as string)", []int{1, 2}, "[1 2]", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := compareEqual(tc.a, tc.b)
			if result != tc.expected {
				t.Errorf("compareEqual(%v, %v) = %v, expected %v", tc.a, tc.b, result, tc.expected)
			}
		})
	}
}

func TestToFloat_EdgeCases(t *testing.T) {
	testCases := []struct {
		name     string
		input    interface{}
		expected float64
	}{
		{"nil", nil, 0.0},
		{"empty string", "", 0.0},
		{"non-numeric string", "abc", 0.0},
		{"negative", -42.5, -42.5},
		{"scientific notation", "1.5e2", 150.0},
		{"hex string", "0xFF", 0.0}, // Not supported
		{"boolean true", true, 0.0},
		{"boolean false", false, 0.0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := toFloat(tc.input)
			if result != tc.expected {
				t.Errorf("toFloat(%v) = %v, expected %v", tc.input, result, tc.expected)
			}
		})
	}
}
