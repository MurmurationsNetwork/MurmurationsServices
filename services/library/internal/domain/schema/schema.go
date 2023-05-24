package schema

type Schema struct {
	Title       string `json:"title"       bson:"title,omitempty"`
	Description string `json:"description" bson:"description,omitempty"`
	Name        string `json:"name"        bson:"name,omitempty"`
	URL         string `json:"url"         bson:"url,omitempty"`
}

type Schemas []*Schema

func (schema *Schema) Marshall() interface{} {
	return schema
}

func (schemas Schemas) Marshall() interface{} {
	data := make([]interface{}, len(schemas))
	for index, schema := range schemas {
		data[index] = schema.Marshall()
	}
	return data
}
