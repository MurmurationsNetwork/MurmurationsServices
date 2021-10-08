package schema

type respond struct {
	Data interface{} `json:"data,omitempty"`
}

type SchemaJSON struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Name        string `json:"name"`
	Version     int    `json:"version"`
	URL         string `json:"url"`
}

func (schema Schema) Marshall() interface{} {
	return SchemaJSON{
		Title:       schema.Title,
		Description: schema.Description,
		Name:        schema.Name,
		Version:     schema.Version,
		URL:         schema.URL,
	}
}

func (schemas Schemas) Marshall() interface{} {
	data := make([]interface{}, len(schemas))
	for index, schema := range schemas {
		data[index] = schema.Marshall()
	}
	return respond{Data: data}
}
