package constant

type NodeStatusType string

type nodeStatus struct {
	Received         NodeStatusType
	Validated        NodeStatusType
	ValidationFailed NodeStatusType
	PostFailed       NodeStatusType
	Posted           NodeStatusType
}

func NodeStatus() *nodeStatus {
	return &nodeStatus{
		Received:         "received",
		Validated:        "validated",
		ValidationFailed: "validation_failed",
		PostFailed:       "post_failed",
		Posted:           "posted",
	}
}

type esIndex struct {
	Node string
}

func ESIndex() *esIndex {
	return &esIndex{
		Node: "nodes",
	}
}
