package event

import "github.com/MurmurationsNetwork/MurmurationsServices/pkg/jsonapi"

type NodeCreatedData struct {
	ProfileURL string `json:"profile_url"`
	Version    int32  `json:"version"`
}

type NodeValidatedData struct {
	ProfileURL  string `json:"profile_url"`
	ProfileHash string `json:"profile_hash"`
	ProfileStr  string `json:"profile_str"`
	LastUpdated int64  `json:"last_updated"`
	Version     int32  `json:"version"`
}

type NodeValidationFailedData struct {
	ProfileURL     string           `json:"profile_url"`
	FailureReasons *[]jsonapi.Error `json:"failure_reasons"`
	Version        int32            `json:"version"`
}
