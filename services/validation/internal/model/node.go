package model

// Node represents a node in a network.
type Node struct {
	// ProfileURL is the URL of the node's profile.
	ProfileURL string `json:"profile_url"`
	// Version is the version number of the node.
	Version int32 `json:"version"`
}
