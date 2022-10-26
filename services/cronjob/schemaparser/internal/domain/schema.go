package domain

type Schema struct {
	Title       string                 `bson:"title,omitempty"`
	Description string                 `bson:"description,omitempty"`
	Name        string                 `bson:"name,omitempty"`
	URL         string                 `bson:"url,omitempty"`
	FullSchema  map[string]interface{} `bson:"full_schema,omitempty"`
}

type DnsInfo struct {
	LastCommit string   `json:"last_commit"`
	SchemaList []string `json:"schema_list"`
	Error      string   `json:"error"`
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
