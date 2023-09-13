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

// WithURLSchemas sets the base URL and the schemas for the ProfileValidator to be built.
func (b *Builder) WithURLSchemas(baseURL string, schemas []string) *Builder {
	b.profilevalidator.Schemas = schemas
	b.profilevalidator.SchemaLoader = &URLSchemaLoader{BaseURL: baseURL}
	return b
}

// WithStrSchemas sets the schemas for the ProfileValidator to be built.
func (b *Builder) WithStrSchemas(schemas []string) *Builder {
	b.profilevalidator.Schemas = schemas
	b.profilevalidator.SchemaLoader = &StrSchemaLoader{}
	return b
}

// WithURLProfile sets the URL for loading the data to be validated.
func (b *Builder) WithURLProfile(dataURL string) *Builder {
	b.profilevalidator.ProfileLoader = &URLProfileLoader{dataURL: dataURL}
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

// WithCustomValidation enables to custom validation.
func (b *Builder) WithCustomValidation() *Builder {
	b.profilevalidator.CustomValidation = true
	return b
}

// Build validates the builder state and returns the built ProfileValidator.
func (b *Builder) Build() (*ProfileValidator, error) {
	// Check that required fields are set.
	if b.profilevalidator.Schemas == nil {
		return nil, fmt.Errorf("schemas must be provided")
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

	b.profilevalidator.JSON = jsonMap

	return b.profilevalidator, nil
}
