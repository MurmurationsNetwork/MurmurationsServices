package node

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/common/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/resterr"
)

type Node struct {
	ID            string              `json:"_id" bson:"_id,omitempty"`
	ProfileUrl    string              `json:"profileUrl" bson:"profileUrl,omitempty"`
	ProfileHash   *string             `json:"profileHash" bson:"profileHash,omitempty"`
	LinkedSchemas []string            `json:"linkedSchemas" bson:"linkedSchemas,omitempty"`
	Status        constant.NodeStatus `json:"status" bson:"status,omitempty"`
	LastChecked int64               `json:"lastChecked" bson:"lastChecked,omitempty"`
	FailedReasons *[]string           `json:"failedReasons" bson:"failedReasons,omitempty"`
	Version       *int32              `json:"-" bson:"version,omitempty"`
}

func (node *Node) Validate() resterr.RestErr {
	if node.ProfileUrl == "" {
		return resterr.NewBadRequestError("profileUrl parameter is missing.")
	}

	if len(node.LinkedSchemas) == 0 {
		return resterr.NewBadRequestError("linkedSchemas parameter is missing.")
	}

	return nil
}

type Nodes []Node

type NodeQuery struct {
	Schema        string `form:"schema"`
	LastChecked int64  `form:"lastChecked"`
}
