package node

import (
	"encoding/json"
)

type respond struct {
	Data interface{} `json:"data,omitempty"`
}

type addNodeRespond struct {
	NodeID        string `json:"nodeId,omitempty"`
	LastChecked int64  `json:"lastChecked,omitempty"`
}

type searchNodeRespond struct {
	ProfileUrl    string `json:"profileUrl,omitempty"`
	LastChecked int64  `json:"lastChecked,omitempty"`
}

func (node *Node) Marshall() interface{} {
	res := addNodeRespond{
		NodeID:        node.ID,
		LastChecked: node.LastChecked,
	}
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
