package event

type NodeCreatedData struct {
	ProfileUrl    string   `json:"profileUrl"`
	LinkedSchemas []string `json:"linkedSchemas"`
}

type NodeValidatedData struct {
	ProfileUrl    string `json:"profileUrl"`
	LastValidated int64  `json:"lastValidated"`
}

type NodeValidationFailedData struct {
	ProfileUrl    string   `json:"profileUrl"`
	FailedReasons []string `json:"failedReasons"`
}
