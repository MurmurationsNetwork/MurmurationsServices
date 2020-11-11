package node

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/common/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/resterr"
)

type Node struct {
	ID             string                  `json:"node_id" bson:"_id,omitempty"`
	ProfileURL     string                  `json:"profile_url" bson:"profile_url,omitempty"`
	ProfileHash    *string                 `json:"profile_hash" bson:"profile_hash,omitempty"`
	Status         constant.NodeStatusType `json:"status" bson:"status,omitempty"`
	LastValidated  *int64                  `json:"last_validated" bson:"last_validated,omitempty"`
	FailureReasons *[]string               `json:"failure_reasons" bson:"failure_reasons,omitempty"`

	Version *int32 `json:"-" bson:"version,omitempty"`

	ProfileStr string `json:"-" bson:"-"`
}

func (node *Node) Validate() resterr.RestErr {
	if node.ProfileURL == "" {
		return resterr.NewBadRequestError("The profile_url parameter is missing.")
	}
	return nil
}

type Nodes []Node
