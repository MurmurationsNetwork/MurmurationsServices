package constant

var NodeStatus = struct {
	Received         string
	Validated        string
	ValidationFailed string
	PostFailed       string
	Posted           string
	Deleted          string
}{
	Received:         "received",
	Validated:        "validated",
	ValidationFailed: "validation_failed",
	PostFailed:       "post_failed",
	Posted:           "posted",
	Deleted:          "deleted",
}
