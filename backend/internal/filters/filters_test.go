package filters

import (
"testing"
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
