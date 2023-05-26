package http

import (
	"encoding/json"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/entity"
)

type Respond struct {
	Data interface{} `json:"data,omitempty"`
}

type AddNodeVO struct {
	NodeID string `json:"node_id,omitempty"`
}

type GetNodeVO struct {
	ID             string    `json:"node_id,omitempty"`
	ProfileURL     string    `json:"profile_url,omitempty"`
	ProfileHash    *string   `json:"profile_hash,omitempty"`
	Status         string    `json:"status,omitempty"`
	LastUpdated    *int64    `json:"last_updated,omitempty"`
	FailureReasons *[]string `json:"failure_reasons,omitempty"`
}

type SearchNodeVO struct {
	ProfileURL  string `json:"profile_url,omitempty"`
	LastUpdated *int64 `json:"last_updated,omitempty"`
}

func (handler *nodeHandler) ToAddNodeVO(node *entity.Node) interface{} {
	return AddNodeVO{
		NodeID: node.ID,
	}
}

func (handler *nodeHandler) ToGetNodeVO(node *entity.Node) interface{} {
	if node.Status != constant.NodeStatus.Validated &&
		node.Status != constant.NodeStatus.Posted &&
		node.Status != constant.NodeStatus.Deleted &&
		node.Status != constant.NodeStatus.PostFailed {
		node.ProfileHash = nil
		node.LastUpdated = nil
	}
	if node.Status != constant.NodeStatus.ValidationFailed {
		node.FailureReasons = nil
	}

	nodeJSON, _ := json.Marshal(toDTO(node))
	var res GetNodeVO
	_ = json.Unmarshal(nodeJSON, &res)
	return res
}

func (handler *nodeHandler) ToSearchNodeVO(nodes entity.Nodes) interface{} {
	data := make([]interface{}, len(nodes))
	for index, node := range nodes {
		nodeJSON, _ := json.Marshal(node)
		var res SearchNodeVO
		_ = json.Unmarshal(nodeJSON, &res)
		data[index] = res
	}
	return Respond{Data: data}
}
