package nodes

import (
	"encoding/json"
)

type respond struct {
	Data interface{} `json:"data"`
}

type addNodeRespond struct {
	NodeID        string `json:"nodeId"`
	LastValidated int64  `json:"lastValidated"`
}

type searchNodeRespond struct {
	ProfileUrl    string `json:"profileUrl"`
	LastValidated int64  `json:"lastValidated"`
}

func (node *Node) Marshall() interface{} {
	nodeJson, _ := json.Marshal(node)
	var res addNodeRespond
	json.Unmarshal(nodeJson, &res)
	return respond{Data: res}
}

func (nodes Nodes) Marshall() interface{} {
	data := make([]interface{}, len(nodes))
	for index, node := range nodes {
		nodeJson, _ := json.Marshal(node)
		var res searchNodeRespond
		json.Unmarshal(nodeJson, &res)
		data[index] = res
	}
	return respond{Data: data}
}
