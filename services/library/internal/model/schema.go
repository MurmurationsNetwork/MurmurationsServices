package model

import (
	"github.com/iancoleman/orderedmap"
	"go.mongodb.org/mongo-driver/bson"
)

// Schema defines the structure for a schema.
type Schema struct {
	Title       string `json:"title"       bson:"title,omitempty"`
	Description string `json:"description" bson:"description,omitempty"`
	Name        string `json:"name"        bson:"name,omitempty"`
	URL         string `json:"url"         bson:"url,omitempty"`
}

// Marshall transforms the Schema instance to an interface.
func (schema *Schema) Marshall() interface{} {
	return schema
}

// Schemas is a slice of Schema instances.
type Schemas []*Schema

func (schemas Schemas) Marshall() interface{} {
	data := make([]interface{}, len(schemas))
	for index, schema := range schemas {
		data[index] = schema.Marshall()
	}
	return data
}

// SingleSchema represents a schema with its description and full schema.
type SingleSchema struct {
	Description string `bson:"description"`
	FullSchema  bson.D `bson:"full_schema"`
}

// ToMap transforms the full schema into an ordered map.
func (s *SingleSchema) ToMap() *orderedmap.OrderedMap {
	result := orderedmap.New()
	for _, element := range s.FullSchema {
		key := element.Key
		value := element.Value

		if innerDoc, ok := value.(bson.D); ok {
			singleSchema := &SingleSchema{FullSchema: innerDoc}
			result.Set(key, singleSchema.ToMap())
		} else {
			result.Set(key, value)
		}
	}
	return result
}
