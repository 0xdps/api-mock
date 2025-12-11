package filters

import (
	"fmt"
	"strconv"
	"strings"
)

// FilterOperator represents the type of comparison operation
type FilterOperator string

const (
	OpEqual         FilterOperator = "="
	OpNotEqual      FilterOperator = "!="
	OpGreater       FilterOperator = ">"
	OpGreaterEqual  FilterOperator = ">="
	OpLess          FilterOperator = "<"
	OpLessEqual     FilterOperator = "<="
	OpContains      FilterOperator = "contains"
	OpStartsWith    FilterOperator = "startsWith"
	OpEndsWith      FilterOperator = "endsWith"
)

// Filter represents a single filter condition
type Filter struct {
	Field    string
	Operator FilterOperator
	Value    interface{}
}

// ParseFilters extracts and parses filter conditions from query parameters
// Format: ?name=John&age>30&email~contains@gmail
// Operators: = (default), !=, >, <, >=, <=, ~contains, ~startsWith, ~endsWith
func ParseFilters(queryParams map[string][]string) []Filter {
	filters := []Filter{}

	for key, values := range queryParams {
		if len(values) == 0 {
			continue
		}

		// Skip reserved parameters (middleware and system params)
		if key == "count" || key == "nocache" || key == "fresh" ||
			key == "page" || key == "limit" || key == "offset" || // Pagination
			key == "sort" || key == "order" || // Sorting
			key == "q" || key == "search" || // Search
			key == "fields" { // Field filtering
			continue
		}

		value := values[0]

		// Parse operator from key
		// Format: fieldName, fieldName>, fieldName<, fieldName>=, fieldName<=, fieldName!=, fieldName~op
		field := key
		operator := OpEqual

		if strings.Contains(key, "~") {
			// e.g., email~contains, name~startsWith
			parts := strings.Split(key, "~")
			if len(parts) == 2 {
				field = parts[0]
				opStr := parts[1]
				switch opStr {
				case "contains":
					operator = OpContains
				case "startsWith":
					operator = OpStartsWith
				case "endsWith":
					operator = OpEndsWith
				}
			}
		} else if strings.HasSuffix(key, "!=") {
			field = strings.TrimSuffix(key, "!=")
			operator = OpNotEqual
		} else if strings.HasSuffix(key, ">=") {
			field = strings.TrimSuffix(key, ">=")
			operator = OpGreaterEqual
		} else if strings.HasSuffix(key, "<=") {
			field = strings.TrimSuffix(key, "<=")
			operator = OpLessEqual
		} else if strings.HasSuffix(key, ">") {
			field = strings.TrimSuffix(key, ">")
			operator = OpGreater
		} else if strings.HasSuffix(key, "<") {
			field = strings.TrimSuffix(key, "<")
			operator = OpLess
		}

		filters = append(filters, Filter{
			Field:    field,
			Operator: operator,
			Value:    value,
		})
	}

	return filters
}

// ApplyFilters applies all filters to a list of items
func ApplyFilters(items []map[string]interface{}, filters []Filter) []map[string]interface{} {
	if len(filters) == 0 {
		return items
	}

	result := []map[string]interface{}{}

	for _, item := range items {
		if matchesAllFilters(item, filters) {
			result = append(result, item)
		}
	}

	return result
}

// matchesAllFilters checks if an item matches all filter conditions
func matchesAllFilters(item map[string]interface{}, filters []Filter) bool {
	for _, filter := range filters {
		if !matchesFilter(item, filter) {
			return false
		}
	}
	return true
}

// matchesFilter checks if an item matches a single filter condition
func matchesFilter(item map[string]interface{}, filter Filter) bool {
	itemValue, ok := item[filter.Field]
	if !ok {
		return false
	}

	switch filter.Operator {
	case OpEqual:
		return compareEqual(itemValue, filter.Value)
	case OpNotEqual:
		return !compareEqual(itemValue, filter.Value)
	case OpGreater:
		return compareGreater(itemValue, filter.Value)
	case OpGreaterEqual:
		return compareGreaterEqual(itemValue, filter.Value)
	case OpLess:
		return compareLess(itemValue, filter.Value)
	case OpLessEqual:
		return compareLessEqual(itemValue, filter.Value)
	case OpContains:
		return compareContains(itemValue, filter.Value)
	case OpStartsWith:
		return compareStartsWith(itemValue, filter.Value)
	case OpEndsWith:
		return compareEndsWith(itemValue, filter.Value)
	default:
		return false
	}
}

// Comparison helpers

func compareEqual(a, b interface{}) bool {
	return fmt.Sprint(a) == fmt.Sprint(b)
}

func compareGreater(a interface{}, b interface{}) bool {
	aNum := toFloat(a)
	bNum := toFloat(b)
	return aNum > bNum
}

func compareGreaterEqual(a interface{}, b interface{}) bool {
	aNum := toFloat(a)
	bNum := toFloat(b)
	return aNum >= bNum
}

func compareLess(a interface{}, b interface{}) bool {
	aNum := toFloat(a)
	bNum := toFloat(b)
	return aNum < bNum
}

func compareLessEqual(a interface{}, b interface{}) bool {
	aNum := toFloat(a)
	bNum := toFloat(b)
	return aNum <= bNum
}

func compareContains(a interface{}, b interface{}) bool {
	return strings.Contains(strings.ToLower(fmt.Sprint(a)), strings.ToLower(fmt.Sprint(b)))
}

func compareStartsWith(a interface{}, b interface{}) bool {
	return strings.HasPrefix(strings.ToLower(fmt.Sprint(a)), strings.ToLower(fmt.Sprint(b)))
}

func compareEndsWith(a interface{}, b interface{}) bool {
	return strings.HasSuffix(strings.ToLower(fmt.Sprint(a)), strings.ToLower(fmt.Sprint(b)))
}

// toFloat converts a value to float64 for numeric comparisons
func toFloat(v interface{}) float64 {
	switch val := v.(type) {
	case float64:
		return val
	case float32:
		return float64(val)
	case int:
		return float64(val)
	case int64:
		return float64(val)
	case string:
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return 0
		}
		return f
	default:
		return 0
	}
}
