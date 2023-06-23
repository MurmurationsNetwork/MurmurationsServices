package rest

import (
	"encoding/json"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/model"
)

// Respond struct is used to format the API response data.
type Respond struct {
	Data interface{} `json:"data,omitempty"`
}

// AddNodeResponse struct is used to format response of AddNode operation.
type AddNodeResponse struct {
	NodeID string `json:"node_id,omitempty"`
}

// GetNodeResponse struct is used to format the GetNode operation response.
type GetNodeResponse struct {
	ID             string    `json:"node_id,omitempty"`
	ProfileURL     string    `json:"profile_url,omitempty"`
	ProfileHash    *string   `json:"profile_hash,omitempty"`
	Status         string    `json:"status,omitempty"`
	LastUpdated    *int64    `json:"last_updated,omitempty"`
	FailureReasons *[]string `json:"failure_reasons,omitempty"`
}

// SearchNodeResponse struct is used to format the SearchNode operation response.
type SearchNodeResponse struct {
	ProfileURL  string `json:"profile_url,omitempty"`
	LastUpdated *int64 `json:"last_updated,omitempty"`
}

// ToAddNodeResponse converts the node model to AddNodeResponse format.
func ToAddNodeResponse(node *model.Node) interface{} {
	return AddNodeResponse{
		NodeID: node.ID,
	}
}

// ToGetNodeResponse converts the node model to GetNodeResponse format.
func ToGetNodeResponse(node *model.Node) interface{} {
	// Modify node properties based on its status.
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

	// Convert the node model to GetNodeResponse.
	nodeJSON, _ := json.Marshal(toDTO(node))
	var res GetNodeResponse
	_ = json.Unmarshal(nodeJSON, &res)
	return res
}

// ToSearchNodeResponse converts the nodes model to SearchNodeResponse format.
func ToSearchNodeResponse(nodes model.Nodes) interface{} {
	data := make([]interface{}, len(nodes))
	for index, node := range nodes {
		nodeJSON, _ := json.Marshal(node)
		var res SearchNodeResponse
		_ = json.Unmarshal(nodeJSON, &res)
		data[index] = res
	}
	return Respond{Data: data}
}
