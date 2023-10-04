package model

// Node represents a node in the system with relevant attributes.
type Node struct {
	// URL of the node's profile.
	ProfileURL string `json:"profile_url" bson:"profile_url,omitempty"`
	// Status of the node.
	Status string `json:"status"      bson:"status,omitempty"`
	// Version is the version vector of the node.
	// https://en.wikipedia.org/wiki/Version_vector
	Version *int32 `json:"-"           bson:"__v,omitempty"`
}
