package schemavalidator

import "fmt"

// Builder is a builder type for creating SchemaValidator instances.
type Builder struct {
	// schemaValidator is the SchemaValidator instance being constructed.
	schemaValidator *SchemaValidator
}

// NewBuilder creates and returns a new Builder.
func NewBuilder() *Builder {
	return &Builder{
		schemaValidator: &SchemaValidator{},
	}
}

// WithURLSchemas sets the base URL and the schemas for the SchemaValidator to be built.
func (b *Builder) WithURLSchemas(baseURL string, schemas []string) *Builder {
	b.schemaValidator.Schemas = schemas
	b.schemaValidator.SchemaLoader = &URLSchemaLoader{BaseURL: baseURL}
	return b
}

// WithStrSchemas sets the schemas for the SchemaValidator to be built.
func (b *Builder) WithStrSchemas(schemas []string) *Builder {
	b.schemaValidator.Schemas = schemas
	b.schemaValidator.SchemaLoader = &StrSchemaLoader{}
	return b
}

// WithURLProfileLoader sets the URL for loading the data to be validated.
func (b *Builder) WithURLProfileLoader(dataURL string) *Builder {
	b.schemaValidator.ProfileLoader = &URLProfileLoader{dataURL: dataURL}
	return b
}

// WithStrProfileLoader sets the data string to be validated.
func (b *Builder) WithStrProfileLoader(dataString string) *Builder {
	b.schemaValidator.ProfileLoader = &StrProfileLoader{dataString: dataString}
	return b
}

// WithMapProfileLoader sets the data map to be validated.
func (b *Builder) WithMapProfileLoader(
	dataMap map[string]interface{},
) *Builder {
	b.schemaValidator.ProfileLoader = &MapProfileLoader{dataMap: dataMap}
	return b
}

// Build validates the builder state and returns the built SchemaValidator.
func (b *Builder) Build() (*SchemaValidator, error) {
	// Check that required fields are set.
	if b.schemaValidator.Schemas == nil {
		return nil, fmt.Errorf("schemas must be provided")
	}
	if b.schemaValidator.SchemaLoader == nil {
		return nil, fmt.Errorf("a schema loader must be provided")
	}
	if b.schemaValidator.ProfileLoader == nil {
		return nil, fmt.Errorf("a data loader must be provided")
	}

	return b.schemaValidator, nil
}
