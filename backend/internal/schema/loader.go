package schema

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/brianvoe/gofakeit/v7"
)

//go:embed embedded/*.json
var embeddedSchemas embed.FS

// RouteConfig defines custom route configuration
type RouteConfig struct {
	Path    string   `json:"path,omitempty"`    // Custom path (e.g., "/v1/users")
	Methods []string `json:"methods,omitempty"` // Allowed methods ["GET", "POST", etc.]
	Aliases []string `json:"aliases,omitempty"` // Additional paths for this resource
}

// ResourceMetadata contains information about the resource
type ResourceMetadata struct {
	Name        string       `json:"name"`
	Singular    string       `json:"singular"`
	Description string       `json:"description"`
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
}

// Schema represents a JSON Schema with our custom extensions
type Schema struct {
	SchemaURI   string                    `json:"$schema"`
	Type        string                    `json:"type"`
	Title       string                    `json:"title"`
	Description string                    `json:"description,omitempty"`
	Resource    ResourceMetadata          `json:"x-resource"`
	Properties  map[string]PropertySchema `json:"properties"`
	Required    []string                  `json:"required,omitempty"`
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
			"%s.%s@example.com",           // john.doe@example.com
			"%s%s@example.com",             // johndoe@example.com
			"%s_%s@example.com",            // john_doe@example.com
			"%s.%s@company.com",            // john.doe@company.com
			"%s%d@example.com",             // john123@example.com
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
			"%s%s",              // johndoe
			"%s_%s",             // john_doe
			"%s.%s",             // john.doe
			"%s%s%d",            // johndoe123
			"%s_%d",             // john_123
			"%c%s",              // jdoe (first initial + last name)
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
