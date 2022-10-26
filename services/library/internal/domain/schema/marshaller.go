package schema

type SchemaJSON struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Name        string `json:"name"`
	URL         string `json:"url"`
}

func (schema Schema) Marshall() interface{} {
	return SchemaJSON{
		Title:       schema.Title,
		Description: schema.Description,
		Name:        schema.Name,
		URL:         schema.URL,
	}
}

func (schemas Schemas) Marshall() interface{} {
	data := make([]interface{}, len(schemas))
	for index, schema := range schemas {
		data[index] = schema.Marshall()
	}
	return data
}
