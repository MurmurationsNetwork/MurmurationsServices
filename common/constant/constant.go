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

type ESIndexType string

type esIndex struct {
	Node ESIndexType
}

func ESIndex() *esIndex {
	return &esIndex{
		Node: "nodes",
	}
}
