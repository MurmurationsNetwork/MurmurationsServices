package event

type NodeCreatedData struct {
	ProfileUrl string `json:"profileUrl"`
	Version    int32  `json:"version"`
}

type NodeValidatedData struct {
	ProfileUrl  string `json:"profileUrl"`
	ProfileHash string `json:"profileHash"`
	ProfileStr  string `json:"profileStr"`
	LastChecked int64  `json:"lastChecked"`
	Version     int32  `json:"version"`
}

type NodeValidationFailedData struct {
	ProfileUrl    string   `json:"profileUrl"`
	FailedReasons []string `json:"failedReasons"`
	Version       int32    `json:"version"`
}
