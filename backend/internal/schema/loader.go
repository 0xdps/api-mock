package schema

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"math"
	"net/mail"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v7"
)

//go:embed embedded/*.json
var embeddedSchemas embed.FS

// RouteConfig defines custom route configuration
type RouteConfig struct {
	Path    string   `json:"path,omitempty"`    // Custom path (e.g., "/users")
	Methods []string `json:"methods,omitempty"` // Allowed methods ["GET", "POST", etc.]
	Aliases []string `json:"aliases,omitempty"` // Additional paths for this resource
}

// ResourceMetadata contains information about the resource
type ResourceMetadata struct {
	Name        string       `json:"name"`
	Singular    string       `json:"singular"`
	Description string       `json:"description"`
	Group       string       `json:"group,omitempty"`  // Resource group (e.g., "people", "commerce")
	Routes      *RouteConfig `json:"routes,omitempty"` // Custom route configuration
}

// GeneratorParams contains parameters for field generation
type GeneratorParams map[string]interface{}

// PropertySchema defines a single property in the schema
type PropertySchema struct {
	Type            string                    `json:"type"`
	Format          string                    `json:"format,omitempty"`
	Description     string                    `json:"description,omitempty"`
	Generator       string                    `json:"x-generator,omitempty"`
	GeneratorParams GeneratorParams           `json:"x-generator-params,omitempty"`
	Properties      map[string]PropertySchema `json:"properties,omitempty"`
	Items           *PropertySchema           `json:"items,omitempty"`
	GeneratorCount  int                       `json:"x-generator-count,omitempty"`
	
	// Enum constraints
	Enum []interface{} `json:"enum,omitempty"`
	
	// Numeric constraints
	Minimum          *float64 `json:"minimum,omitempty"`
	Maximum          *float64 `json:"maximum,omitempty"`
	ExclusiveMinimum *float64 `json:"exclusiveMinimum,omitempty"`
	ExclusiveMaximum *float64 `json:"exclusiveMaximum,omitempty"`
	MultipleOf       *float64 `json:"multipleOf,omitempty"`
	
	// String constraints
	MinLength *int    `json:"minLength,omitempty"`
	MaxLength *int    `json:"maxLength,omitempty"`
	Pattern   string  `json:"pattern,omitempty"`
	
	// Array constraints
	MinItems    *int `json:"minItems,omitempty"`
	MaxItems    *int `json:"maxItems,omitempty"`
	UniqueItems bool `json:"uniqueItems,omitempty"`
}

// Schema represents a JSON Schema with our custom extensions
type Schema struct {
	SchemaURI            string                    `json:"$schema"`
	Type                 string                    `json:"type"`
	Title                string                    `json:"title"`
	Description          string                    `json:"description,omitempty"`
	Resource             ResourceMetadata          `json:"x-resource"`
	Properties           map[string]PropertySchema `json:"properties"`
	Required             []string                  `json:"required,omitempty"`
	AdditionalProperties interface{}               `json:"additionalProperties,omitempty"`
}

// Field represents a data field to generate
type Field struct {
	Name      string
	Generator string
	Args      map[string]interface{}
}

// Registry holds all loaded schemas
type Registry struct {
	Schemas map[string]*Schema
}

// NewRegistry creates a new schema registry
func NewRegistry() *Registry {
	return &Registry{
		Schemas: make(map[string]*Schema),
	}
}

// LoadEmbeddedSchemas loads all schemas from embedded files
func (r *Registry) LoadEmbeddedSchemas() error {
	entries, err := fs.ReadDir(embeddedSchemas, "embedded")
	if err != nil {
		return fmt.Errorf("failed to read embedded schemas: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		data, err := embeddedSchemas.ReadFile(filepath.Join("embedded", entry.Name()))
		if err != nil {
			return fmt.Errorf("failed to read embedded schema %s: %w", entry.Name(), err)
		}

		if err := r.loadSchemaFromData(data, entry.Name()); err != nil {
			return fmt.Errorf("failed to load schema %s: %w", entry.Name(), err)
		}
	}

	return nil
}

// LoadSchemas loads all schemas from the shared/schemas directory
func (r *Registry) LoadSchemas(schemasDir string) error {
	files, err := os.ReadDir(schemasDir)
	if err != nil {
		return fmt.Errorf("failed to read schemas directory: %w", err)
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		schemaPath := filepath.Join(schemasDir, file.Name())
		if err := r.LoadSchema(schemaPath); err != nil {
			return fmt.Errorf("failed to load schema %s: %w", file.Name(), err)
		}
	}

	return nil
}

// LoadSchema loads a single schema file
func (r *Registry) LoadSchema(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read schema file: %w", err)
	}

	return r.loadSchemaFromData(data, filepath.Base(path))
}

// loadSchemaFromData loads schema from byte data
func (r *Registry) loadSchemaFromData(data []byte, filename string) error {
	var schema Schema
	if err := json.Unmarshal(data, &schema); err != nil {
		return fmt.Errorf("failed to parse schema: %w", err)
	}

	// Validate that resource metadata exists
	if schema.Resource.Name == "" {
		return fmt.Errorf("schema missing x-resource.name")
	}

	r.Schemas[schema.Resource.Name] = &schema
	return nil
}

// GetSchema retrieves a schema by resource name
func (r *Registry) GetSchema(resourceName string) (*Schema, bool) {
	schema, ok := r.Schemas[resourceName]
	return schema, ok
}

// GetAllResourceNames returns all registered resource names
func (r *Registry) GetAllResourceNames() []string {
	names := make([]string, 0, len(r.Schemas))
	for name := range r.Schemas {
		names = append(names, name)
	}
	return names
}

// ValidateItem validates an item against a schema
// ValidateItem validates an item against the schema
// skipFields can be used to skip validation of certain fields (e.g., id during POST)
func (s *Schema) ValidateItem(item map[string]interface{}, skipFields ...string) error {
	// Check required fields
	skipMap := make(map[string]bool)
	for _, field := range skipFields {
		skipMap[field] = true
	}
	
	for _, requiredField := range s.Required {
		if skipMap[requiredField] {
			continue
		}
		if _, exists := item[requiredField]; !exists {
			return fmt.Errorf("missing required field: %s", requiredField)
		}
	}

	// Validate each property in the item
	for propName, propValue := range item {
		propSchema, exists := s.Properties[propName]
		if !exists {
			// Check if additional properties are allowed
			if s.AdditionalProperties != nil {
				if additionalPropsAllowed, ok := s.AdditionalProperties.(bool); ok && !additionalPropsAllowed {
					return fmt.Errorf("additional property '%s' is not allowed by schema", propName)
				}
				// If it's a schema object, we could validate against it, but for now skip
			}
			// Extra properties allowed (lenient mode by default)
			continue
		}

		if err := validateProperty(propName, propValue, propSchema); err != nil {
			return err
		}
	}

	return nil
}

// validateProperty validates a single property value against its schema
func validateProperty(name string, value interface{}, schema PropertySchema) error {
	if value == nil {
		// Null values are allowed unless field is required (checked separately)
		return nil
	}

	// Validate enum constraint (applies to all types)
	if len(schema.Enum) > 0 {
		valid := false
		for _, allowed := range schema.Enum {
			if value == allowed {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("field '%s' must be one of %v, got '%v'", name, schema.Enum, value)
		}
	}

	switch schema.Type {
	case "string":
		str, ok := value.(string)
		if !ok {
			return fmt.Errorf("field '%s' must be a string, got %T", name, value)
		}

		// String length constraints
		if schema.MinLength != nil && len(str) < *schema.MinLength {
			return fmt.Errorf("field '%s' must be at least %d characters, got %d", name, *schema.MinLength, len(str))
		}
		if schema.MaxLength != nil && len(str) > *schema.MaxLength {
			return fmt.Errorf("field '%s' must be at most %d characters, got %d", name, *schema.MaxLength, len(str))
		}

		// Pattern constraint
		if schema.Pattern != "" {
			matched, err := regexp.MatchString(schema.Pattern, str)
			if err != nil {
				return fmt.Errorf("field '%s' has invalid pattern in schema: %v", name, err)
			}
			if !matched {
				return fmt.Errorf("field '%s' does not match required pattern '%s'", name, schema.Pattern)
			}
		}

		// Format validation
		if schema.Format != "" {
			if err := validateFormat(name, str, schema.Format); err != nil {
				return err
			}
		}

	case "number", "integer":
		var numValue float64
		
		switch v := value.(type) {
		case float64:
			numValue = v
			// For integer type, check if it's a whole number
			if schema.Type == "integer" && v != float64(int64(v)) {
				return fmt.Errorf("field '%s' must be an integer, got float %v", name, v)
			}
		case float32:
			numValue = float64(v)
		case int:
			numValue = float64(v)
		case int64:
			numValue = float64(v)
		case int32:
			numValue = float64(v)
		default:
			return fmt.Errorf("field '%s' must be a number, got %T", name, value)
		}

		// Numeric constraints
		if schema.Minimum != nil {
			if numValue < *schema.Minimum {
				return fmt.Errorf("field '%s' must be >= %v, got %v", name, *schema.Minimum, numValue)
			}
		}
		if schema.Maximum != nil {
			if numValue > *schema.Maximum {
				return fmt.Errorf("field '%s' must be <= %v, got %v", name, *schema.Maximum, numValue)
			}
		}
		if schema.ExclusiveMinimum != nil {
			if numValue <= *schema.ExclusiveMinimum {
				return fmt.Errorf("field '%s' must be > %v, got %v", name, *schema.ExclusiveMinimum, numValue)
			}
		}
		if schema.ExclusiveMaximum != nil {
			if numValue >= *schema.ExclusiveMaximum {
				return fmt.Errorf("field '%s' must be < %v, got %v", name, *schema.ExclusiveMaximum, numValue)
			}
		}
		if schema.MultipleOf != nil && *schema.MultipleOf > 0 {
			remainder := math.Mod(numValue, *schema.MultipleOf)
			if math.Abs(remainder) > 1e-10 { // floating point tolerance
				return fmt.Errorf("field '%s' must be a multiple of %v, got %v", name, *schema.MultipleOf, numValue)
			}
		}

	case "boolean":
		if _, ok := value.(bool); !ok {
			return fmt.Errorf("field '%s' must be a boolean, got %T", name, value)
		}

	case "array":
		arr, ok := value.([]interface{})
		if !ok {
			return fmt.Errorf("field '%s' must be an array, got %T", name, value)
		}

		// Array length constraints
		if schema.MinItems != nil && len(arr) < *schema.MinItems {
			return fmt.Errorf("field '%s' must have at least %d items, got %d", name, *schema.MinItems, len(arr))
		}
		if schema.MaxItems != nil && len(arr) > *schema.MaxItems {
			return fmt.Errorf("field '%s' must have at most %d items, got %d", name, *schema.MaxItems, len(arr))
		}

		// Unique items constraint
		if schema.UniqueItems {
			seen := make(map[string]bool)
			for i, item := range arr {
				// Convert to JSON string for comparison
				itemJSON, _ := json.Marshal(item)
				key := string(itemJSON)
				if seen[key] {
					return fmt.Errorf("field '%s' must have unique items, duplicate found at index %d", name, i)
				}
				seen[key] = true
			}
		}

		// Validate array items if schema specifies items type
		if schema.Items != nil {
			for i, item := range arr {
				if err := validateProperty(fmt.Sprintf("%s[%d]", name, i), item, *schema.Items); err != nil {
					return err
				}
			}
		}

	case "object":
		obj, ok := value.(map[string]interface{})
		if !ok {
			return fmt.Errorf("field '%s' must be an object, got %T", name, value)
		}

		// Validate nested properties if schema specifies them
		if len(schema.Properties) > 0 {
			for nestedName, nestedValue := range obj {
				if nestedSchema, exists := schema.Properties[nestedName]; exists {
					if err := validateProperty(fmt.Sprintf("%s.%s", name, nestedName), nestedValue, nestedSchema); err != nil {
						return err
					}
				}
			}
		}

	default:
		// Unknown type or no type specified - accept any value
		return nil
	}

	return nil
}

// validateFormat validates string format constraints
func validateFormat(name string, value string, format string) error {
	switch format {
	case "email":
		if _, err := mail.ParseAddress(value); err != nil {
			return fmt.Errorf("field '%s' must be a valid email address, got '%s'", name, value)
		}

	case "uri", "url":
		if _, err := url.ParseRequestURI(value); err != nil {
			return fmt.Errorf("field '%s' must be a valid URI, got '%s'", name, value)
		}

	case "date":
		if _, err := time.Parse("2006-01-02", value); err != nil {
			return fmt.Errorf("field '%s' must be a valid date in YYYY-MM-DD format, got '%s'", name, value)
		}

	case "date-time":
		// Try RFC3339 format (JSON standard)
		if _, err := time.Parse(time.RFC3339, value); err != nil {
			// Also try RFC3339Nano for more precision
			if _, err2 := time.Parse(time.RFC3339Nano, value); err2 != nil {
				return fmt.Errorf("field '%s' must be a valid date-time in RFC3339 format, got '%s'", name, value)
			}
		}

	case "uuid":
		// UUID v4 pattern: xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx
		uuidPattern := `^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[1-5][0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$`
		matched, err := regexp.MatchString(uuidPattern, value)
		if err != nil || !matched {
			return fmt.Errorf("field '%s' must be a valid UUID, got '%s'", name, value)
		}

	case "hostname":
		// Simple hostname validation
		hostnamePattern := `^([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)*[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?$`
		matched, err := regexp.MatchString(hostnamePattern, value)
		if err != nil || !matched {
			return fmt.Errorf("field '%s' must be a valid hostname, got '%s'", name, value)
		}

	case "ipv4":
		ipPattern := `^(\d{1,3}\.){3}\d{1,3}$`
		matched, err := regexp.MatchString(ipPattern, value)
		if err != nil || !matched {
			return fmt.Errorf("field '%s' must be a valid IPv4 address, got '%s'", name, value)
		}
		// Additional validation for each octet
		parts := strings.Split(value, ".")
		for _, part := range parts {
			var octet int
			fmt.Sscanf(part, "%d", &octet)
			if octet < 0 || octet > 255 {
				return fmt.Errorf("field '%s' must be a valid IPv4 address, got '%s'", name, value)
			}
		}

	// Note: Other formats like ipv6, time, regex, etc. can be added as needed
	}

	return nil
}

// SchemaToFields converts a schema into generator field definitions
func (r *Registry) SchemaToFields(schema *Schema) []Field {
	var fields []Field

	for propName, prop := range schema.Properties {
		field := r.propertyToField(propName, prop)
		if field.Name != "" {
			fields = append(fields, field)
		}
	}

	// Sort fields to ensure dependent fields come after their dependencies
	return r.sortFieldsByDependency(fields)
}

// sortFieldsByDependency orders fields so dependencies are generated first
func (r *Registry) sortFieldsByDependency(fields []Field) []Field {
	// Define field priority (lower number = generated first)
	priority := map[string]int{
		"id":          1,
		"first_name":  2,
		"last_name":   2,
		"name":        2, // Can be product name or person name
		"email":       3, // Depends on first_name, last_name
		"username":    3, // Depends on first_name, last_name
		"avatar":      4, // Depends on id or first_name
		"title":       5, // Post/article title
		"body":        6, // Depends on title
		"description": 6, // Depends on title/name
	}

	// Sort fields by priority, keeping relative order for same priority
	sorted := make([]Field, len(fields))
	copy(sorted, fields)

	// Simple bubble sort by priority
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			priI := priority[sorted[i].Name]
			priJ := priority[sorted[j].Name]

			// Default priority is 10 if not specified
			if priI == 0 {
				priI = 10
			}
			if priJ == 0 {
				priJ = 10
			}

			if priI > priJ {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	return sorted
}

// propertyToField converts a property schema to a field definition
func (r *Registry) propertyToField(name string, prop PropertySchema) Field {
	genType := prop.Generator
	if genType == "" {
		genType = r.inferGeneratorType(prop)
	}

	// Handle special case: autoincrement should use random_int for now
	// (actual autoincrement would need to track state)
	if genType == "autoincrement" {
		genType = "random_int"
		// Set reasonable range for IDs if not specified
		if prop.GeneratorParams == nil {
			return Field{
				Name:      name,
				Generator: genType,
				Args: map[string]interface{}{
					"min": 1,
					"max": 10000,
				},
			}
		}
	}

	return Field{
		Name:      name,
		Generator: genType,
		Args:      prop.GeneratorParams,
	}
}

// inferGeneratorType infers the generator type from JSON Schema type/format
func (r *Registry) inferGeneratorType(prop PropertySchema) string {
	switch prop.Type {
	case "string":
		if prop.Format == "email" {
			return "email"
		}
		if prop.Format == "uri" {
			return "url"
		}
		if prop.Format == "date-time" {
			return "past"
		}
		if prop.Format == "uuid" {
			return "uuid"
		}
		return "word"
	case "integer":
		return "number"
	case "number":
		return "float"
	case "boolean":
		return "bool"
	default:
		return "word"
	}
}

// GenerateData generates fake data for a resource
func (r *Registry) GenerateData(resourceName string, count int) ([]map[string]interface{}, error) {
	schema, ok := r.GetSchema(resourceName)
	if !ok {
		return nil, fmt.Errorf("schema not found: %s", resourceName)
	}

	fields := r.SchemaToFields(schema)
	return r.generateRecords(fields, count)
}

// generateRecords generates multiple records using gofakeit
func (r *Registry) generateRecords(fields []Field, count int) ([]map[string]interface{}, error) {
	records := make([]map[string]interface{}, count)

	for i := 0; i < count; i++ {
		record := make(map[string]interface{})

		// Generate fields, allowing later fields to reference earlier ones
		for _, field := range fields {
			value := r.generateValueWithContext(field, record)
			record[field.Name] = value
		}
		records[i] = record
	}

	return records, nil
}

// generateValueWithContext generates a value based on field type and previously generated fields
func (r *Registry) generateValueWithContext(field Field, record map[string]interface{}) interface{} {
	faker := gofakeit.New(0)

	// Handle derived fields that depend on other fields
	switch field.Generator {
	case "email":
		return r.generateSmartEmail(record, faker)
	case "username":
		return r.generateSmartUsername(record, faker)
	case "avatar":
		return r.generateSmartAvatar(record, faker)
	case "sentence":
		// For title fields, generate more realistic titles
		if field.Name == "title" {
			return r.generateSmartTitle(faker)
		}
	case "paragraph":
		// For body/description, generate context-aware content
		if field.Name == "body" || field.Name == "description" {
			if title, ok := record["title"].(string); ok && title != "" {
				return r.generateSmartBody(title, faker)
			}
		}
	}

	return r.generateValue(field, faker)
}

// generateSmartTitle creates realistic titles for posts/articles
func (r *Registry) generateSmartTitle(faker *gofakeit.Faker) string {
	templates := []string{
		"How to %s in %d Easy Steps",
		"The Ultimate Guide to %s",
		"Understanding %s: A Complete Tutorial",
		"%d Tips for Better %s",
		"Why %s Matters in %d",
		"Getting Started with %s",
		"Advanced %s Techniques",
		"The Future of %s",
		"Common %s Mistakes to Avoid",
		"Best Practices for %s",
	}

	topics := []string{
		"Web Development", "Mobile Apps", "Cloud Computing", "Data Science",
		"Machine Learning", "API Design", "User Experience", "Cybersecurity",
		"DevOps", "Microservices", "React Development", "Node.js",
		"Python Programming", "JavaScript", "System Design", "Database Optimization",
		"Testing Strategies", "Performance Tuning", "Code Review", "Agile Methods",
	}

	template := templates[faker.IntRange(0, len(templates)-1)]
	topic := topics[faker.IntRange(0, len(topics)-1)]

	if strings.Contains(template, "%d") {
		year := faker.IntRange(2020, 2025)
		return fmt.Sprintf(template, topic, year)
	}

	return fmt.Sprintf(template, topic)
}

// generateSmartBody creates contextual body text based on title
func (r *Registry) generateSmartBody(title string, faker *gofakeit.Faker) string {
	// Extract key topic from title if possible
	intro := fmt.Sprintf("In this article, we'll explore %s. ", title)

	paragraphs := []string{intro}

	// Add 2-3 paragraphs
	numParagraphs := faker.IntRange(2, 4)
	for i := 0; i < numParagraphs; i++ {
		paragraphs = append(paragraphs, faker.Paragraph(3, 5, 12, " "))
	}

	return strings.Join(paragraphs, "\n\n")
}

// generateSmartEmail creates an email based on first_name and last_name if available
func (r *Registry) generateSmartEmail(record map[string]interface{}, faker *gofakeit.Faker) string {
	firstName, hasFirst := record["first_name"].(string)
	lastName, hasLast := record["last_name"].(string)

	if hasFirst && hasLast && firstName != "" && lastName != "" {
		// Clean names (remove spaces, lowercase)
		firstName = strings.ToLower(strings.ReplaceAll(firstName, " ", ""))
		lastName = strings.ToLower(strings.ReplaceAll(lastName, " ", ""))

		// Choose email format randomly
		formats := []string{
			"%s.%s@example.com", // john.doe@example.com
			"%s%s@example.com",  // johndoe@example.com
			"%s_%s@example.com", // john_doe@example.com
			"%s.%s@company.com", // john.doe@company.com
			"%s%d@example.com",  // john123@example.com
		}

		format := formats[faker.IntRange(0, len(formats)-1)]

		if strings.Contains(format, "%d") {
			// Format with number (first name only)
			return fmt.Sprintf(format, firstName, faker.IntRange(1, 9999))
		}

		return fmt.Sprintf(format, firstName, lastName)
	}

	// Fallback to random email
	return faker.Email()
}

// generateSmartUsername creates a username based on first_name and last_name if available
func (r *Registry) generateSmartUsername(record map[string]interface{}, faker *gofakeit.Faker) string {
	firstName, hasFirst := record["first_name"].(string)
	lastName, hasLast := record["last_name"].(string)

	if hasFirst && hasLast && firstName != "" && lastName != "" {
		// Clean names
		firstName = strings.ToLower(strings.ReplaceAll(firstName, " ", ""))
		lastName = strings.ToLower(strings.ReplaceAll(lastName, " ", ""))

		// Choose username format randomly
		formats := []string{
			"%s%s",   // johndoe
			"%s_%s",  // john_doe
			"%s.%s",  // john.doe
			"%s%s%d", // johndoe123
			"%s_%d",  // john_123
			"%c%s",   // jdoe (first initial + last name)
		}

		format := formats[faker.IntRange(0, len(formats)-1)]

		switch format {
		case "%s%s":
			return firstName + lastName
		case "%s_%s":
			return firstName + "_" + lastName
		case "%s.%s":
			return firstName + "." + lastName
		case "%s%s%d":
			return firstName + lastName + fmt.Sprintf("%d", faker.IntRange(1, 9999))
		case "%s_%d":
			return firstName + "_" + fmt.Sprintf("%d", faker.IntRange(1, 9999))
		case "%c%s":
			if len(firstName) > 0 {
				return string(firstName[0]) + lastName
			}
			return firstName + lastName
		}
	}

	// Fallback to random username
	return faker.Username()
}

// generateSmartAvatar creates a consistent avatar based on user data
func (r *Registry) generateSmartAvatar(record map[string]interface{}, faker *gofakeit.Faker) string {
	// Use ID if available for consistency
	if id, ok := record["id"].(int); ok && id > 0 {
		// Use ID modulo 70 (pravatar.cc has 70 avatars)
		avatarNum := (id % 70) + 1
		return fmt.Sprintf("https://i.pravatar.cc/300?img=%d", avatarNum)
	}

	// Use first_name for seeding if available
	if firstName, ok := record["first_name"].(string); ok && firstName != "" {
		// Generate a consistent number from the name
		sum := 0
		for _, char := range firstName {
			sum += int(char)
		}
		avatarNum := (sum % 70) + 1
		return fmt.Sprintf("https://i.pravatar.cc/300?img=%d", avatarNum)
	}

	// Fallback to random avatar
	return fmt.Sprintf("https://i.pravatar.cc/300?img=%d", faker.IntRange(1, 70))
}

// generateValue generates a single value based on field type
func (r *Registry) generateValue(field Field, faker *gofakeit.Faker) interface{} {

	switch field.Generator {
	// Numbers
	case "autoincrement", "random_int", "number":
		min := 1
		max := 10000
		if field.Args != nil {
			if m, ok := field.Args["min"].(float64); ok {
				min = int(m)
			}
			if m, ok := field.Args["max"].(float64); ok {
				max = int(m)
			}
		}
		return faker.IntRange(min, max)
	case "float":
		return faker.Float64Range(0, 1000)
	case "price":
		return faker.Price(10, 1000)

	// Personal Info
	case "name":
		return faker.Name()
	case "first_name":
		return faker.FirstName()
	case "last_name":
		return faker.LastName()
	case "email":
		return faker.Email()
	case "username":
		return faker.Username()
	case "password":
		return faker.Password(true, true, true, false, false, 12)
	case "gender":
		return faker.Gender()

	// Address & Location
	case "address":
		return faker.Address().Address
	case "street":
		return faker.Address().Street
	case "city":
		return faker.City()
	case "state":
		return faker.State()
	case "country":
		return faker.Country()
	case "zip", "zip_code":
		return faker.Zip()
	case "latitude":
		return faker.Latitude()
	case "longitude":
		return faker.Longitude()

	// Contact
	case "phone", "phone_number":
		return faker.Phone()

	// Internet
	case "url":
		return faker.URL()
	case "domain", "domain_name":
		return faker.DomainName()
	case "ipv4":
		return faker.IPv4Address()
	case "ipv6":
		return faker.IPv6Address()
	case "uuid":
		return faker.UUID()
	case "mac_address":
		return faker.MacAddress()
	case "user_agent":
		return faker.UserAgent()

	// Dates & Time
	case "date":
		return faker.Date().Format("2006-01-02")
	case "past_date", "past":
		return faker.PastDate().Format("2006-01-02T15:04:05Z07:00")
	case "future_date", "future":
		return faker.FutureDate().Format("2006-01-02T15:04:05Z07:00")
	case "date_time":
		return faker.Date().Format("2006-01-02T15:04:05Z07:00")

	// Text
	case "word":
		return faker.Word()
	case "sentence":
		return faker.Sentence(8)
	case "paragraph":
		return faker.Paragraph(2, 4, 12, " ")
	case "text":
		return faker.Paragraph(3, 5, 15, " ")

	// Company
	case "company":
		return faker.Company()
	case "job", "job_title":
		return faker.JobTitle()
	case "catch_phrase":
		return faker.BuzzWord()

	// E-commerce
	case "currency":
		currencies := []string{"USD", "EUR", "GBP", "JPY", "AUD", "CAD", "CHF", "CNY", "INR"}
		return currencies[faker.IntRange(0, len(currencies)-1)]
	case "category":
		categories := []string{"Electronics", "Clothing", "Books", "Home & Garden", "Sports", "Toys", "Food & Beverage", "Health & Beauty", "Automotive", "Office"}
		return categories[faker.IntRange(0, len(categories)-1)]

	// Weather
	case "temperature":
		return float64(faker.IntRange(-20, 45)) + faker.Float64Range(0, 0.9)
	case "humidity":
		return faker.IntRange(10, 100)
	case "pressure":
		return faker.IntRange(950, 1050)
	case "wind_speed":
		return float64(faker.IntRange(0, 50)) + faker.Float64Range(0, 0.9)
	case "wind_direction":
		directions := []string{"N", "NE", "E", "SE", "S", "SW", "W", "NW"}
		return directions[faker.IntRange(0, len(directions)-1)]
	case "weather_condition":
		conditions := []string{"Sunny", "Partly Cloudy", "Cloudy", "Overcast", "Rainy", "Stormy", "Snowy", "Foggy", "Windy", "Clear"}
		return conditions[faker.IntRange(0, len(conditions)-1)]
	case "weather_description":
		descriptions := []string{
			"Clear skies with plenty of sunshine",
			"Partly cloudy with occasional sun",
			"Overcast skies throughout the day",
			"Light rain expected in the afternoon",
			"Heavy rainfall with possible thunderstorms",
			"Light snow showers",
			"Dense fog reducing visibility",
			"Strong winds with gusts",
		}
		return descriptions[faker.IntRange(0, len(descriptions)-1)]
	case "visibility":
		return faker.IntRange(1000, 50000)
	case "uv_index":
		return faker.IntRange(0, 11)

	// Geographic
	case "country_code":
		codes := []string{"US", "GB", "FR", "DE", "JP", "CN", "IN", "BR", "CA", "AU", "MX", "ES", "IT", "KR", "RU"}
		return codes[faker.IntRange(0, len(codes)-1)]
	case "country_code3":
		codes := []string{"USA", "GBR", "FRA", "DEU", "JPN", "CHN", "IND", "BRA", "CAN", "AUS", "MEX", "ESP", "ITA", "KOR", "RUS"}
		return codes[faker.IntRange(0, len(codes)-1)]
	case "capital":
		capitals := []string{"Washington D.C.", "London", "Paris", "Berlin", "Tokyo", "Beijing", "New Delhi", "Brasilia", "Ottawa", "Canberra"}
		return capitals[faker.IntRange(0, len(capitals)-1)]
	case "region":
		regions := []string{"Africa", "Americas", "Asia", "Europe", "Oceania", "Antarctic"}
		return regions[faker.IntRange(0, len(regions)-1)]
	case "subregion":
		subregions := []string{"Eastern Asia", "Western Europe", "Northern America", "Southern Asia", "South America", "Caribbean", "Middle East"}
		return subregions[faker.IntRange(0, len(subregions)-1)]
	case "population":
		return faker.IntRange(100000, 1500000000)
	case "area":
		return float64(faker.IntRange(1000, 17000000))
	case "currency_code":
		codes := []string{"USD", "EUR", "GBP", "JPY", "CNY", "INR", "BRL", "CAD", "AUD", "CHF", "SEK", "NOK", "DKK", "MXN", "SGD"}
		return codes[faker.IntRange(0, len(codes)-1)]
	case "calling_code":
		codes := []string{"+1", "+44", "+33", "+49", "+81", "+86", "+91", "+55", "+61", "+52", "+34", "+39", "+82", "+7"}
		return codes[faker.IntRange(0, len(codes)-1)]
	case "timezone":
		zones := []string{"UTC", "EST", "PST", "GMT", "CET", "JST", "IST", "AEST", "CST", "MST"}
		return zones[faker.IntRange(0, len(zones)-1)]
	case "flag_url":
		code := faker.IntRange(1, 200)
		return fmt.Sprintf("https://flagcdn.com/w320/%d.png", code)
	case "languages_array", "currency_countries", "language_countries":
		count := faker.IntRange(1, 3)
		arr := make([]string, count)
		for i := 0; i < count; i++ {
			arr[i] = faker.Word()
		}
		return arr
	case "elevation":
		return faker.IntRange(0, 5000)
	case "city_population":
		return faker.IntRange(10000, 20000000)

	// Financial
	case "exchange_rate":
		return faker.Float64Range(0.1, 10.0)
	case "currency_name":
		names := []string{"US Dollar", "Euro", "British Pound", "Japanese Yen", "Chinese Yuan", "Indian Rupee", "Brazilian Real", "Canadian Dollar", "Australian Dollar"}
		return names[faker.IntRange(0, len(names)-1)]
	case "currency_symbol":
		symbols := []string{"$", "€", "£", "¥", "₹", "R$", "C$", "A$", "kr", "₽"}
		return symbols[faker.IntRange(0, len(symbols)-1)]
	case "stock_symbol":
		symbols := []string{"AAPL", "GOOGL", "MSFT", "AMZN", "TSLA", "META", "NVDA", "NFLX", "AMD", "INTC", "ORCL", "CSCO", "IBM", "UBER", "SNAP"}
		return symbols[faker.IntRange(0, len(symbols)-1)]
	case "stock_price":
		return faker.Float64Range(10.0, 500.0)
	case "stock_change":
		return faker.Float64Range(-50.0, 50.0)
	case "stock_change_percent":
		return faker.Float64Range(-10.0, 10.0)
	case "stock_volume":
		return faker.IntRange(1000000, 100000000)
	case "market_cap":
		return float64(faker.IntRange(1000000000, 3000000000000))
	case "stock_exchange":
		exchanges := []string{"NYSE", "NASDAQ", "LSE", "TSE", "SSE", "HKEX", "Euronext"}
		return exchanges[faker.IntRange(0, len(exchanges)-1)]
	case "crypto_symbol":
		symbols := []string{"BTC", "ETH", "USDT", "BNB", "XRP", "ADA", "DOGE", "SOL", "DOT", "MATIC", "AVAX", "LINK", "UNI", "ATOM", "LTC"}
		return symbols[faker.IntRange(0, len(symbols)-1)]
	case "crypto_name":
		names := []string{"Bitcoin", "Ethereum", "Tether", "Binance Coin", "Ripple", "Cardano", "Dogecoin", "Solana", "Polkadot", "Polygon", "Avalanche", "Chainlink", "Uniswap", "Cosmos", "Litecoin"}
		return names[faker.IntRange(0, len(names)-1)]
	case "crypto_price":
		return faker.Float64Range(0.01, 50000.0)
	case "crypto_supply":
		return float64(faker.IntRange(1000000, 100000000000))

	// Business
	case "industry":
		industries := []string{"Technology", "Finance", "Healthcare", "Education", "Manufacturing", "Retail", "Real Estate", "Entertainment", "Transportation", "Energy"}
		return industries[faker.IntRange(0, len(industries)-1)]
	case "company_size":
		sizes := []string{"Small (1-50)", "Medium (51-200)", "Large (201-1000)", "Enterprise (1000+)"}
		return sizes[faker.IntRange(0, len(sizes)-1)]
	case "founded_year":
		return faker.IntRange(1900, 2024)
	case "revenue":
		return float64(faker.IntRange(100000, 1000000000))
	case "employees":
		return faker.IntRange(10, 50000)
	case "salary_min":
		return faker.IntRange(30000, 80000)
	case "salary_max":
		return faker.IntRange(80000, 200000)
	case "job_type":
		types := []string{"Full-time", "Part-time", "Contract", "Internship", "Freelance"}
		return types[faker.IntRange(0, len(types)-1)]
	case "job_level":
		levels := []string{"Entry Level", "Mid Level", "Senior Level", "Lead", "Manager", "Director", "Executive"}
		return levels[faker.IntRange(0, len(levels)-1)]
	case "organization_type":
		types := []string{"Corporation", "Non-Profit", "Government", "Startup", "SME", "Enterprise"}
		return types[faker.IntRange(0, len(types)-1)]

	// Education
	case "course_title":
		titles := []string{
			"Introduction to Programming",
			"Web Development Bootcamp",
			"Data Science Fundamentals",
			"Machine Learning A-Z",
			"Digital Marketing Mastery",
			"Business Strategy",
			"Graphic Design Essentials",
			"Financial Analysis",
		}
		return titles[faker.IntRange(0, len(titles)-1)]
	case "course_category":
		categories := []string{"Programming", "Business", "Design", "Marketing", "Data Science", "Photography", "Music", "Languages"}
		return categories[faker.IntRange(0, len(categories)-1)]
	case "course_level":
		levels := []string{"Beginner", "Intermediate", "Advanced", "Expert"}
		return levels[faker.IntRange(0, len(levels)-1)]
	case "course_duration":
		return faker.IntRange(5, 100)
	case "student_id":
		return fmt.Sprintf("STU%06d", faker.IntRange(1, 999999))
	case "grade":
		grades := []string{"Freshman", "Sophomore", "Junior", "Senior", "Graduate"}
		return grades[faker.IntRange(0, len(grades)-1)]
	case "major":
		majors := []string{"Computer Science", "Business Administration", "Engineering", "Psychology", "Biology", "Mathematics", "English", "History"}
		return majors[faker.IntRange(0, len(majors)-1)]
	case "gpa":
		return faker.Float64Range(2.0, 4.0)

	// Media & Entertainment
	case "movie_title":
		titles := []string{
			"The Last Adventure", "City of Dreams", "Beyond the Horizon", "Lost in Time",
			"The Final Chapter", "Midnight Express", "Rising Sun", "Dark Waters",
		}
		return titles[faker.IntRange(0, len(titles)-1)]
	case "movie_genre":
		genres := []string{"Action", "Comedy", "Drama", "Horror", "Sci-Fi", "Romance", "Thriller", "Documentary", "Animation"}
		return genres[faker.IntRange(0, len(genres)-1)]
	case "movie_year":
		return faker.IntRange(1980, 2024)
	case "movie_duration":
		return faker.IntRange(80, 180)
	case "movie_budget":
		return float64(faker.IntRange(1000000, 300000000))
	case "movie_revenue":
		return float64(faker.IntRange(5000000, 2000000000))
	case "book_title":
		titles := []string{
			"The Art of Programming", "Journey to Success", "Mastering the Mind",
			"The Hidden Truth", "Beyond Imagination", "Stories of Wonder",
		}
		return titles[faker.IntRange(0, len(titles)-1)]
	case "book_genre":
		genres := []string{"Fiction", "Non-Fiction", "Mystery", "Sci-Fi", "Biography", "Self-Help", "History", "Fantasy"}
		return genres[faker.IntRange(0, len(genres)-1)]
	case "book_pages":
		return faker.IntRange(100, 1000)
	case "isbn":
		return fmt.Sprintf("978-%d-%d-%d-%d", faker.IntRange(0, 9), faker.IntRange(10000, 99999), faker.IntRange(100, 999), faker.IntRange(0, 9))
	case "music_genre":
		genres := []string{"Pop", "Rock", "Hip Hop", "Jazz", "Classical", "Electronic", "Country", "R&B", "Blues"}
		return genres[faker.IntRange(0, len(genres)-1)]
	case "album_title":
		return faker.Sentence(3)
	case "album_duration":
		return faker.IntRange(1800, 4800)
	case "video_duration":
		return faker.IntRange(60, 3600)
	case "video_category":
		categories := []string{"Education", "Entertainment", "Music", "Gaming", "News", "Sports", "Technology", "Travel", "Food"}
		return categories[faker.IntRange(0, len(categories)-1)]
	case "podcast_category":
		categories := []string{"Technology", "Business", "Comedy", "News", "Education", "True Crime", "Health", "Sports"}
		return categories[faker.IntRange(0, len(categories)-1)]
	case "podcast_duration":
		return faker.IntRange(15, 180)
	case "news_category":
		categories := []string{"Politics", "Business", "Technology", "Science", "Health", "Entertainment", "Sports", "World"}
		return categories[faker.IntRange(0, len(categories)-1)]

	// Food & Travel
	case "cuisine":
		cuisines := []string{"Italian", "Chinese", "Japanese", "Mexican", "Indian", "French", "Thai", "Mediterranean", "American", "Korean"}
		return cuisines[faker.IntRange(0, len(cuisines)-1)]
	case "restaurant_name":
		return faker.Company() + " " + []string{"Restaurant", "Bistro", "Cafe", "Grill", "Kitchen", "Diner"}[faker.IntRange(0, 5)]
	case "price_range":
		ranges := []string{"$", "$$", "$$$", "$$$$"}
		return ranges[faker.IntRange(0, len(ranges)-1)]
	case "opening_hours":
		return "Mon-Sat: 11:00 AM - 10:00 PM"
	case "hotel_name":
		return faker.Company() + " " + []string{"Hotel", "Resort", "Inn", "Lodge", "Suites"}[faker.IntRange(0, 4)]
	case "hotel_stars":
		return faker.IntRange(1, 5)
	case "hotel_price":
		return faker.Float64Range(50.0, 500.0)
	case "hotel_rooms":
		return faker.IntRange(20, 500)
	case "amenities":
		amenities := []string{"WiFi", "Pool", "Gym", "Restaurant", "Bar", "Spa", "Parking", "Room Service"}
		count := faker.IntRange(3, 6)
		selected := make([]string, count)
		for i := 0; i < count; i++ {
			selected[i] = amenities[faker.IntRange(0, len(amenities)-1)]
		}
		return selected
	case "recipe_title":
		return faker.Sentence(4)
	case "recipe_category":
		categories := []string{"Breakfast", "Lunch", "Dinner", "Dessert", "Snacks", "Appetizers", "Beverages"}
		return categories[faker.IntRange(0, len(categories)-1)]
	case "difficulty":
		levels := []string{"Easy", "Medium", "Hard"}
		return levels[faker.IntRange(0, len(levels)-1)]
	case "prep_time":
		return faker.IntRange(5, 60)
	case "cook_time":
		return faker.IntRange(10, 120)
	case "servings":
		return faker.IntRange(1, 8)
	case "calories":
		return faker.IntRange(100, 1000)
	case "venue":
		return faker.Company() + " Convention Center"
	case "event_title":
		return faker.Sentence(5)
	case "event_category":
		categories := []string{"Conference", "Workshop", "Seminar", "Meetup", "Concert", "Festival", "Networking"}
		return categories[faker.IntRange(0, len(categories)-1)]

	// Transportation
	case "flight_number":
		airlines := []string{"AA", "UA", "DL", "BA", "LH", "AF", "EK"}
		return airlines[faker.IntRange(0, len(airlines)-1)] + fmt.Sprintf("%d", faker.IntRange(100, 9999))
	case "airline":
		airlines := []string{"American Airlines", "United Airlines", "Delta", "British Airways", "Lufthansa", "Air France", "Emirates"}
		return airlines[faker.IntRange(0, len(airlines)-1)]
	case "airport_code":
		codes := []string{"JFK", "LAX", "ORD", "LHR", "CDG", "NRT", "DXB", "SIN", "HND", "SYD"}
		return codes[faker.IntRange(0, len(codes)-1)]
	case "flight_duration":
		return faker.IntRange(60, 900)
	case "aircraft":
		aircraft := []string{"Boeing 737", "Boeing 777", "Boeing 787", "Airbus A320", "Airbus A380", "Airbus A350"}
		return aircraft[faker.IntRange(0, len(aircraft)-1)]
	case "flight_price":
		return faker.Float64Range(100.0, 2000.0)
	case "flight_class":
		classes := []string{"Economy", "Premium Economy", "Business", "First Class"}
		return classes[faker.IntRange(0, len(classes)-1)]
	case "flight_status":
		statuses := []string{"Scheduled", "Boarding", "Departed", "In Air", "Landed", "Delayed", "Cancelled"}
		return statuses[faker.IntRange(0, len(statuses)-1)]

	// Automotive
	case "car_make":
		makes := []string{"Toyota", "Honda", "Ford", "BMW", "Mercedes", "Audi", "Tesla", "Chevrolet", "Nissan", "Volkswagen"}
		return makes[faker.IntRange(0, len(makes)-1)]
	case "car_model":
		models := []string{"Sedan", "SUV", "Coupe", "Hatchback", "Truck", "Van", "Convertible"}
		return models[faker.IntRange(0, len(models)-1)]
	case "car_year":
		return faker.IntRange(2010, 2024)
	case "color":
		colors := []string{"Black", "White", "Silver", "Gray", "Red", "Blue", "Green", "Yellow", "Orange", "Brown"}
		return colors[faker.IntRange(0, len(colors)-1)]
	case "vin":
		return faker.UUID()[0:17]
	case "car_type":
		types := []string{"Sedan", "SUV", "Truck", "Coupe", "Convertible", "Hatchback", "Van", "Wagon"}
		return types[faker.IntRange(0, len(types)-1)]
	case "fuel_type":
		types := []string{"Gasoline", "Diesel", "Electric", "Hybrid", "Plug-in Hybrid"}
		return types[faker.IntRange(0, len(types)-1)]
	case "transmission":
		types := []string{"Automatic", "Manual", "CVT", "Semi-Automatic"}
		return types[faker.IntRange(0, len(types)-1)]
	case "mileage":
		return faker.IntRange(0, 200000)
	case "car_price":
		return faker.Float64Range(5000.0, 100000.0)
	case "condition":
		conditions := []string{"New", "Like New", "Good", "Fair", "Poor"}
		return conditions[faker.IntRange(0, len(conditions)-1)]
	case "car_features":
		features := []string{"Cruise Control", "Backup Camera", "Navigation", "Leather Seats", "Sunroof", "Heated Seats", "Bluetooth", "Parking Sensors"}
		count := faker.IntRange(3, 6)
		selected := make([]string, count)
		for i := 0; i < count; i++ {
			selected[i] = features[faker.IntRange(0, len(features)-1)]
		}
		return selected

	// Real Estate
	case "property_title":
		return faker.Sentence(5)
	case "property_type":
		types := []string{"House", "Apartment", "Condo", "Townhouse", "Villa", "Land", "Commercial"}
		return types[faker.IntRange(0, len(types)-1)]
	case "property_status":
		statuses := []string{"For Sale", "For Rent", "Sold", "Rented", "Pending"}
		return statuses[faker.IntRange(0, len(statuses)-1)]
	case "property_price":
		return faker.Float64Range(50000.0, 5000000.0)
	case "bedrooms":
		return faker.IntRange(1, 6)
	case "bathrooms":
		return faker.IntRange(1, 5)
	case "property_area":
		return float64(faker.IntRange(500, 5000))
	case "year_built":
		return faker.IntRange(1950, 2024)
	case "property_features":
		features := []string{"Garden", "Pool", "Garage", "Balcony", "Fireplace", "Basement", "Attic", "Security System"}
		count := faker.IntRange(2, 5)
		selected := make([]string, count)
		for i := 0; i < count; i++ {
			selected[i] = features[faker.IntRange(0, len(features)-1)]
		}
		return selected

	// Business Operations
	case "invoice_number":
		return fmt.Sprintf("INV-%06d", faker.IntRange(1, 999999))
	case "order_number":
		return fmt.Sprintf("ORD-%08d", faker.IntRange(1, 99999999))
	case "transaction_id":
		return faker.UUID()
	case "tax":
		return faker.Float64Range(5.0, 100.0)
	case "shipping_cost":
		return faker.Float64Range(5.0, 50.0)
	case "invoice_status":
		statuses := []string{"Draft", "Sent", "Paid", "Overdue", "Cancelled"}
		return statuses[faker.IntRange(0, len(statuses)-1)]
	case "order_status":
		statuses := []string{"Pending", "Processing", "Shipped", "Delivered", "Cancelled", "Refunded"}
		return statuses[faker.IntRange(0, len(statuses)-1)]
	case "payment_method":
		methods := []string{"Credit Card", "Debit Card", "PayPal", "Bank Transfer", "Cash", "Cryptocurrency"}
		return methods[faker.IntRange(0, len(methods)-1)]
	case "payment_status":
		statuses := []string{"Pending", "Processing", "Completed", "Failed", "Refunded"}
		return statuses[faker.IntRange(0, len(statuses)-1)]
	case "card_last4":
		return fmt.Sprintf("%04d", faker.IntRange(0, 9999))
	case "card_brand":
		brands := []string{"Visa", "Mastercard", "American Express", "Discover"}
		return brands[faker.IntRange(0, len(brands)-1)]
	case "subscription_plan":
		plans := []string{"Free", "Basic", "Pro", "Premium", "Enterprise"}
		return plans[faker.IntRange(0, len(plans)-1)]
	case "subscription_status":
		statuses := []string{"Active", "Cancelled", "Expired", "Trial", "Suspended"}
		return statuses[faker.IntRange(0, len(statuses)-1)]
	case "billing_cycle":
		cycles := []string{"Monthly", "Quarterly", "Annually"}
		return cycles[faker.IntRange(0, len(cycles)-1)]

	// Project Management
	case "project_status":
		statuses := []string{"Planning", "In Progress", "On Hold", "Completed", "Cancelled"}
		return statuses[faker.IntRange(0, len(statuses)-1)]
	case "task_status":
		statuses := []string{"To Do", "In Progress", "In Review", "Done", "Blocked"}
		return statuses[faker.IntRange(0, len(statuses)-1)]
	case "priority":
		priorities := []string{"Low", "Medium", "High", "Urgent"}
		return priorities[faker.IntRange(0, len(priorities)-1)]
	case "percentage":
		return faker.IntRange(0, 100)
	case "meeting_type":
		types := []string{"Standup", "Planning", "Review", "Retrospective", "One-on-One", "All-Hands"}
		return types[faker.IntRange(0, len(types)-1)]
	case "meeting_status":
		statuses := []string{"Scheduled", "In Progress", "Completed", "Cancelled"}
		return statuses[faker.IntRange(0, len(statuses)-1)]

	// Support & Communication
	case "ticket_number":
		return fmt.Sprintf("TKT-%06d", faker.IntRange(1, 999999))
	case "ticket_status":
		statuses := []string{"Open", "In Progress", "Pending", "Resolved", "Closed"}
		return statuses[faker.IntRange(0, len(statuses)-1)]
	case "ticket_category":
		categories := []string{"Technical", "Billing", "General", "Feature Request", "Bug Report"}
		return categories[faker.IntRange(0, len(categories)-1)]
	case "message_type":
		types := []string{"Text", "Image", "Video", "File", "Link"}
		return types[faker.IntRange(0, len(types)-1)]
	case "notification_type":
		types := []string{"Info", "Success", "Warning", "Error", "Message", "System"}
		return types[faker.IntRange(0, len(types)-1)]

	// Sports
	case "sport":
		sports := []string{"Football", "Basketball", "Baseball", "Soccer", "Hockey", "Tennis", "Golf", "Cricket"}
		return sports[faker.IntRange(0, len(sports)-1)]
	case "league":
		leagues := []string{"Premier League", "La Liga", "Serie A", "Bundesliga", "MLS", "NBA", "NFL", "MLB"}
		return leagues[faker.IntRange(0, len(leagues)-1)]
	case "team_name":
		return faker.City() + " " + []string{"United", "City", "Athletic", "Rangers", "Warriors", "Tigers"}[faker.IntRange(0, 5)]
	case "stadium":
		return faker.City() + " " + []string{"Stadium", "Arena", "Field", "Dome"}[faker.IntRange(0, 3)]
	case "position":
		positions := []string{"Forward", "Midfielder", "Defender", "Goalkeeper", "Guard", "Center", "Pitcher"}
		return positions[faker.IntRange(0, len(positions)-1)]
	case "height":
		return float64(faker.IntRange(160, 210))
	case "weight":
		return float64(faker.IntRange(60, 120))
	case "match_status":
		statuses := []string{"Scheduled", "Live", "Halftime", "Finished", "Postponed", "Cancelled"}
		return statuses[faker.IntRange(0, len(statuses)-1)]

	// Misc
	case "status":
		statuses := []string{"Active", "Inactive", "Pending", "Draft", "Published"}
		return statuses[faker.IntRange(0, len(statuses)-1)]
	case "slug":
		return strings.ToLower(strings.ReplaceAll(faker.Sentence(3), " ", "-"))
	case "tags":
		tags := []string{"technology", "business", "lifestyle", "travel", "food", "health", "education", "entertainment"}
		count := faker.IntRange(1, 4)
		selected := make([]string, count)
		for i := 0; i < count; i++ {
			selected[i] = tags[faker.IntRange(0, len(tags)-1)]
		}
		return selected
	case "article_category", "faq_category", "quote_category":
		categories := []string{"Technology", "Business", "Lifestyle", "Health", "Education", "Entertainment"}
		return categories[faker.IntRange(0, len(categories)-1)]
	case "coupon_code":
		return fmt.Sprintf("SAVE%d", faker.IntRange(10, 50))
	case "discount_type":
		types := []string{"Percentage", "Fixed Amount", "Free Shipping"}
		return types[faker.IntRange(0, len(types)-1)]
	case "image_width":
		return faker.IntRange(800, 4000)
	case "image_height":
		return faker.IntRange(600, 3000)
	case "image_format":
		formats := []string{"JPEG", "PNG", "WebP", "GIF", "SVG"}
		return formats[faker.IntRange(0, len(formats)-1)]
	case "file_size":
		return faker.IntRange(10000, 10000000)
	case "quote":
		quotes := []string{
			"The only way to do great work is to love what you do.",
			"Innovation distinguishes between a leader and a follower.",
			"Your time is limited, don't waste it living someone else's life.",
			"Stay hungry, stay foolish.",
			"The future belongs to those who believe in the beauty of their dreams.",
		}
		return quotes[faker.IntRange(0, len(quotes)-1)]
	case "language":
		languages := []string{"English", "Spanish", "French", "German", "Chinese", "Japanese", "Arabic", "Russian", "Portuguese", "Hindi"}
		return languages[faker.IntRange(0, len(languages)-1)]
	case "language_native":
		return faker.Language()
	case "language_code":
		codes := []string{"en", "es", "fr", "de", "zh", "ja", "ar", "ru", "pt", "hi"}
		return codes[faker.IntRange(0, len(codes)-1)]
	case "language_code3":
		codes := []string{"eng", "spa", "fra", "deu", "zho", "jpn", "ara", "rus", "por", "hin"}
		return codes[faker.IntRange(0, len(codes)-1)]
	case "language_family":
		families := []string{"Indo-European", "Sino-Tibetan", "Afro-Asiatic", "Austronesian", "Niger-Congo"}
		return families[faker.IntRange(0, len(families)-1)]
	case "speakers":
		return faker.IntRange(1000000, 1500000000)

	// Other
	case "bool", "boolean":
		return faker.Bool()
	case "image_url":
		return fmt.Sprintf("https://picsum.photos/400/300?random=%d", faker.IntRange(1, 10000))
	case "avatar":
		return fmt.Sprintf("https://i.pravatar.cc/300?img=%d", faker.IntRange(1, 70))

	default:
		// Fallback to sentence for better-looking data
		return faker.Sentence(5)
	}
}

// GetRoutePath returns the custom path or default path for a resource
func (r *Registry) GetRoutePath(resourceName string) string {
	schema, ok := r.GetSchema(resourceName)
	if !ok {
		return "/" + resourceName
	}

	if schema.Resource.Routes != nil && schema.Resource.Routes.Path != "" {
		return schema.Resource.Routes.Path
	}

	return "/" + resourceName
}

// GetRouteAliases returns additional paths for a resource
func (r *Registry) GetRouteAliases(resourceName string) []string {
	schema, ok := r.GetSchema(resourceName)
	if !ok || schema.Resource.Routes == nil {
		return nil
	}

	return schema.Resource.Routes.Aliases
}

// GetRouteMethods returns allowed methods for a resource (defaults to ["GET"])
func (r *Registry) GetRouteMethods(resourceName string) []string {
	schema, ok := r.GetSchema(resourceName)
	if !ok || schema.Resource.Routes == nil || len(schema.Resource.Routes.Methods) == 0 {
		return []string{"GET"}
	}

	return schema.Resource.Routes.Methods
}

// SupportsMethod checks if a resource supports a specific HTTP method
func (r *Registry) SupportsMethod(resourceName string, method string) bool {
	methods := r.GetRouteMethods(resourceName)
	for _, m := range methods {
		if strings.EqualFold(m, method) {
			return true
		}
	}
	return false
}

// GetAllGroups returns all unique groups from loaded schemas
func (r *Registry) GetAllGroups() map[string][]string {
	groups := make(map[string][]string)

	for name, schema := range r.Schemas {
		group := schema.Resource.Group
		if group == "" {
			group = "other" // Default group for schemas without a group
		}
		groups[group] = append(groups[group], name)
	}

	return groups
}

// GetSchemasByGroup returns all schemas in a specific group
func (r *Registry) GetSchemasByGroup(groupName string) map[string]*Schema {
	schemas := make(map[string]*Schema)

	for name, schema := range r.Schemas {
		group := schema.Resource.Group
		if group == "" {
			group = "other"
		}
		if group == groupName {
			schemas[name] = schema
		}
	}

	return schemas
}

// GetResourceNamesByGroup returns all resource names in a specific group
func (r *Registry) GetResourceNamesByGroup(groupName string) []string {
	var names []string

	for name, schema := range r.Schemas {
		group := schema.Resource.Group
		if group == "" {
			group = "other"
		}
		if group == groupName {
			names = append(names, name)
		}
	}

	return names
}
