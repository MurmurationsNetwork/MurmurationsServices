package event

type NodeCreatedData struct {
	ProfileURL string `json:"profileUrl"`
	Version    int32  `json:"version"`
}

type NodeValidatedData struct {
	ProfileURL  string `json:"profileUrl"`
	ProfileHash string `json:"profileHash"`
	ProfileStr  string `json:"profileStr"`
	LastChecked int64  `json:"lastChecked"`
	Version     int32  `json:"version"`
}

type NodeValidationFailedData struct {
	ProfileURL    string   `json:"profileUrl"`
	FailedReasons []string `json:"failedReasons"`
	Version       int32    `json:"version"`
}
