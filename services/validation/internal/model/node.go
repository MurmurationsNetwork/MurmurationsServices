package model

// Node represents a node stored in the index.
type Node struct {
	// ProfileURL is the URL of the node's profile.
	ProfileURL string `json:"profile_url"`
	// Version is the version vector of the node.
	// https://en.wikipedia.org/wiki/Version_vector
	Version int32 `json:"version"`
}
