package domain

type Schema struct {
	Title       string `bson:"title,omitempty"`
	Description string `bson:"description,omitempty"`
	Name        string `bson:"name,omitempty"`
	Version     int    `bson:"version,omitempty"`
	URL         string `bson:"url,omitempty"`
}

type DnsInfo struct {
	LastCommit string   `json:"last_commit"`
	SchemaList []string `json:"schema_list"`
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
