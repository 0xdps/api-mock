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

	return fields
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
		for _, field := range fields {
			value := r.generateValue(field)
			record[field.Name] = value
		}
		records[i] = record
	}
	
	return records, nil
}

// generateValue generates a single value based on field type
func (r *Registry) generateValue(field Field) interface{} {
	faker := gofakeit.New(0)
	
	switch field.Generator {
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
	case "email":
		return faker.Email()
	case "first_name":
		return faker.FirstName()
	case "last_name":
		return faker.LastName()
	case "username":
		return faker.Username()
	case "word", "sentence":
		return faker.Sentence(5)
	case "text", "paragraph":
		return faker.Paragraph(1, 3, 10, " ")
	case "bool", "boolean":
		return faker.Bool()
	case "past", "date_time":
		return faker.PastDate()
	case "url":
		return faker.URL()
	case "uuid":
		return faker.UUID()
	case "city":
		return faker.City()
	case "country":
		return faker.Country()
	case "image_url":
		return fmt.Sprintf("https://picsum.photos/400/300?random=%d", faker.IntRange(1, 10000))
	case "price":
		return faker.Price(10, 1000)
	case "float":
		return faker.Float64Range(0, 1000)
	default:
		return faker.Word()
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
