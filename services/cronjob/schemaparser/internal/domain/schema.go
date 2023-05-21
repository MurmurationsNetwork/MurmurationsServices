package domain

import "go.mongodb.org/mongo-driver/bson"

type Schema struct {
	Title       string `bson:"title,omitempty"`
	Description string `bson:"description,omitempty"`
	Name        string `bson:"name,omitempty"`
	URL         string `bson:"url,omitempty"`
	FullSchema  bson.D `bson:"full_schema,omitempty"`
}

type BranchInfo struct {
	Commit struct {
		Sha         string `json:"sha"`
		InnerCommit struct {
			Author struct {
				Date string `json:"date"`
			} `json:"author"`
		} `json:"commit"`
	} `json:"commit"`
}

type SchemaJSON struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Metadata    struct {
		Schema struct {
			Name    string `json:"name"`
			Version int    `json:"version"`
			URL     string `json:"url"`
		} `json:"schema"`
	} `json:"metadata"`
}
