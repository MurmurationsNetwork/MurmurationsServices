package node

type Node struct {
	ProfileUrl    string   `json:"profileUrl"`
	LinkedSchemas []string `json:"linkedSchemas"`
	Version       int32    `json:"version"`
}
