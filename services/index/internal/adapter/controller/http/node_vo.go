package http

import (
	"encoding/json"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/entity"
)

type respond struct {
	Data interface{} `json:"data,omitempty"`
}

type addNodeVO struct {
	NodeID string `json:"node_id,omitempty"`
}

type getNodeVO struct {
	ID             string    `json:"node_id,omitempty"`
	ProfileURL     string    `json:"profile_url,omitempty"`
	ProfileHash    *string   `json:"profile_hash,omitempty"`
	Status         string    `json:"status,omitempty"`
	LastUpdated    *int64    `json:"last_updated,omitempty"`
	FailureReasons *[]string `json:"failure_reasons,omitempty"`
}

type searchNodeVO struct {
	ProfileURL  string `json:"profile_url,omitempty"`
	LastUpdated *int64 `json:"last_updated,omitempty"`
}

func (handler *nodeHandler) toAddNodeVO(node *entity.Node) interface{} {
	return addNodeVO{
		NodeID: node.ID,
	}
}

func (handler *nodeHandler) toGetNodeVO(node *entity.Node) interface{} {
	if node.Status != constant.NodeStatus.Validated && node.Status != constant.NodeStatus.Posted && node.Status != constant.NodeStatus.Deleted && node.Status != constant.NodeStatus.PostFailed {
		node.ProfileHash = nil
		node.LastUpdated = nil
	}
	if node.Status != constant.NodeStatus.ValidationFailed {
		node.FailureReasons = nil
	}

	nodeJSON, _ := json.Marshal(toDTO(node))
	var res getNodeVO
	json.Unmarshal(nodeJSON, &res)
	return respond{Data: res}
}

func (handler *nodeHandler) toSearchNodeVO(nodes entity.Nodes) interface{} {
	data := make([]interface{}, len(nodes))
	for index, node := range nodes {
		nodeJson, _ := json.Marshal(node)
		var res searchNodeVO
		json.Unmarshal(nodeJson, &res)
		data[index] = res
	}
	return respond{Data: data}
}
