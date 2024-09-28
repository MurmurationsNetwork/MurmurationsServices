package profilevalidator

import "fmt"

// Builder is a builder type for creating ProfileValidator instances.
type Builder struct {
	// profilevalidator is the ProfileValidator instance being constructed.
	profilevalidator *ProfileValidator
}

// NewBuilder creates and returns a new Builder.
func NewBuilder() *Builder {
	return &Builder{
		profilevalidator: &ProfileValidator{},
	}
}

// WithURLSchemas configures the ProfileValidator to use schema references
// that are accessible via URLs. The base URL is used to fetch the schemas.
func (b *Builder) WithURLSchemas(baseURL string, schemaReferences []string) *Builder {
	// Assign schema references, which will be used to fetch schemas from the base URL.
	b.profilevalidator.SchemaReferences = schemaReferences

	// Set schema names to match the schema references, as they represent the same identifiers.
	b.profilevalidator.SchemaNames = schemaReferences

	// Use a URL-based loader for retrieving schema content.
	b.profilevalidator.SchemaLoader = &URLSchemaLoader{BaseURL: baseURL}
	return b
}

// WithJSONSchemas configures the ProfileValidator to use preloaded JSON schemas.
// The schemaNames correspond to the names of the loaded JSON schemas.
func (b *Builder) WithJSONSchemas(schemaNames []string, loadedSchemas []string) *Builder {
	// Assign the actual JSON schema content that will be used for validation.
	b.profilevalidator.LoadedSchemas = loadedSchemas

	// Assign the schema names to provide context for each JSON schema.
	b.profilevalidator.SchemaNames = schemaNames

	// Use a string-based loader since the schemas are already loaded into memory.
	b.profilevalidator.SchemaLoader = &StrSchemaLoader{}
	return b
}

// WithStrProfile sets the data string to be validated.
func (b *Builder) WithStrProfile(dataString string) *Builder {
	b.profilevalidator.ProfileLoader = &StrProfileLoader{dataString: dataString}
	return b
}

// WithMapProfile sets the data map to be validated.
func (b *Builder) WithMapProfile(
	dataMap map[string]interface{},
) *Builder {
	b.profilevalidator.ProfileLoader = &MapProfileLoader{dataMap: dataMap}
	return b
}

// Build validates the builder state and returns the built ProfileValidator.
func (b *Builder) Build() (*ProfileValidator, error) {
	// Check that required fields are set.
	if b.profilevalidator.SchemaNames == nil {
		return nil, fmt.Errorf("schema names must be provided")
	}
	if b.profilevalidator.SchemaLoader == nil {
		return nil, fmt.Errorf("a schema loader must be provided")
	}
	if b.profilevalidator.ProfileLoader == nil {
		return nil, fmt.Errorf("a data loader must be provided")
	}

	profileData, err := b.profilevalidator.ProfileLoader.Load().LoadJSON()
	if err != nil {
		return nil, fmt.Errorf(
			"failed to load JSON from the profile URL: %w",
			err,
		)
	}

	jsonMap, ok := profileData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf(
			"invalid JSON format: expected a map[string]interface{}",
		)
	}

	b.profilevalidator.ProfileJSON = jsonMap

	return b.profilevalidator, nil
}
