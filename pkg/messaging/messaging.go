package messaging

import "github.com/MurmurationsNetwork/MurmurationsServices/pkg/jsonapi"

type NodeCreatedData struct {
	ProfileURL string `json:"profile_url"`
	Version    int32  `json:"version"`
}

// NodeValidatedData represents the validated data of a node.
type NodeValidatedData struct {
	// ProfileURL is the URL of the profile associated with the node.
	ProfileURL string `json:"profile_url"`

	// ProfileHash is the hash of the profile data.
	ProfileHash string `json:"profile_hash"`

	// ProfileStr is a string representation of the profile data.
	ProfileStr string `json:"profile_str"`

	// LastUpdated is a Unix timestamp indicating when the node data was last updated.
	LastUpdated int64 `json:"last_updated"`

	// Version is the version vector of the node.
	// https://en.wikipedia.org/wiki/Version_vector
	Version int32 `json:"version"`
}

type NodeValidationFailedData struct {
	ProfileURL     string           `json:"profile_url"`
	FailureReasons *[]jsonapi.Error `json:"failure_reasons"`
	Version        int32            `json:"version"`
}
