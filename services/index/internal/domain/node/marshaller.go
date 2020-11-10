package node

import (
	"encoding/json"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/constant"
)

type respond struct {
	Data interface{} `json:"data,omitempty"`
}

type addNodeRespond struct {
	NodeID string `json:"nodeId,omitempty"`
}

type getNodeRespond struct {
	ID            string                  `json:"nodeId,omitempty"`
	ProfileUrl    string                  `json:"profileUrl,omitempty"`
	ProfileHash   *string                 `json:"profileHash,omitempty"`
	Status        constant.NodeStatusType `json:"status,omitempty"`
	LastChecked   *int64                  `json:"lastChecked,omitempty"`
	FailedReasons *[]string               `json:"failedReasons,omitempty"`
}

type searchNodeRespond struct {
	ProfileUrl  string `json:"profileUrl,omitempty"`
	LastChecked *int64 `json:"lastChecked,omitempty"`
}

func (node *Node) AddNodeRespond() interface{} {
	res := addNodeRespond{
		NodeID: node.ID,
	}
	return respond{Data: res}
}

func (node *Node) GetNodeRespond() interface{} {
	if node.Status != constant.NodeStatus().Posted {
		node.ProfileHash = nil
		node.LastChecked = nil
	}
	if node.Status != constant.NodeStatus().ValidationFailed && node.Status != constant.NodeStatus().PostFailed {
		node.FailedReasons = nil
	}

	nodeJson, _ := json.Marshal(node)
	var res getNodeRespond
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
