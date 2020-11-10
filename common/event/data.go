package event

type NodeCreatedData struct {
	ProfileURL string `json:"profile_url"`
	Version    int32  `json:"version"`
}

type NodeValidatedData struct {
	ProfileURL    string `json:"profile_url"`
	ProfileHash   string `json:"profile_hash"`
	ProfileStr    string `json:"profile_str"`
	LastValidated int64  `json:"last_validated"`
	Version       int32  `json:"version"`
}

type NodeValidationFailedData struct {
	ProfileURL    string   `json:"profile_url"`
	FailedReasons []string `json:"failure_reasons"`
	Version       int32    `json:"version"`
}
