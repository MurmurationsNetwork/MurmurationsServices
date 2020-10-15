package nodes

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/utils/rest_errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Node struct {
	ID            primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	NodeID        string             `json:"nodeId" bson:"nodeId,omitempty"`
	ProfileUrl    string             `json:"profileUrl" bson:"profileUrl,omitempty"`
	LinkedSchemas []string           `json:"linkedSchemas" bson:"linkedSchemas,omitempty"`
	LastValidated int64              `json:"lastValidated" bson:"lastValidated,omitempty"`
}

func (node *Node) Validate() rest_errors.RestErr {
	if node.ProfileUrl == "" {
		return rest_errors.NewBadRequestError("profileUrl parameter is missing.")
	}

	if len(node.LinkedSchemas) == 0 {
		return rest_errors.NewBadRequestError("linkedSchemas parameter is missing.")
	}

	return nil
}
