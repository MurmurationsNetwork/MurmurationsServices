package node

import (
	"encoding/json"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/constant"
)

type respond struct {
	Data interface{} `json:"data,omitempty"`
}

type addNodeRespond struct {
	NodeID string `json:"node_id,omitempty"`
}

type getNodeRespond struct {
	ID             string    `json:"node_id,omitempty"`
	ProfileURL     string    `json:"profile_url,omitempty"`
	ProfileHash    *string   `json:"profile_hash,omitempty"`
	Status         string    `json:"status,omitempty"`
	LastValidated  *int64    `json:"last_validated,omitempty"`
	FailureReasons *[]string `json:"failure_reasons,omitempty"`
}

type searchNodeRespond struct {
	ProfileURL    string `json:"profile_url,omitempty"`
	LastValidated *int64 `json:"last_validated,omitempty"`
}

func (node *Node) AddNodeRespond() interface{} {
	res := addNodeRespond{
		NodeID: node.ID,
	}
	return respond{Data: res}
}

func (node *Node) GetNodeRespond() interface{} {
	if node.Status != constant.NodeStatus.Validated && node.Status != constant.NodeStatus.Posted {
		node.ProfileHash = nil
		node.LastValidated = nil
	}
	if node.Status != constant.NodeStatus.ValidationFailed {
		node.FailureReasons = nil
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
