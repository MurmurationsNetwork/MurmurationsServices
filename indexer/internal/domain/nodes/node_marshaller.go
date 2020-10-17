package nodes

import (
	"encoding/json"
)

type AddNodeRespond struct {
	NodeID        string `json:"nodeId"`
	LastValidated int64  `json:"lastValidated"`
}

func (node *Node) AddNodeRespond() interface{} {
	nodeJson, _ := json.Marshal(node)
	var addNodeRespond AddNodeRespond
	json.Unmarshal(nodeJson, &addNodeRespond)
	return addNodeRespond
}
