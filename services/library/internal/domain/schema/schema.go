package schema

type Schema struct {
	Title       string `bson:"title,omitempty"`
	Description string `bson:"description,omitempty"`
	Name        string `bson:"name,omitempty"`
	URL         string `bson:"url,omitempty"`
}

type Schemas []*Schema
