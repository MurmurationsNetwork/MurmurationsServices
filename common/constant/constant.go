package constant

type NodeStatus string

const (
	Received         NodeStatus = "received"
	Validated        NodeStatus = "validated"
	ValidationFailed NodeStatus = "validation_failed"
	PostFailed       NodeStatus = "post_failed"
	Posted           NodeStatus = "posted"
)
