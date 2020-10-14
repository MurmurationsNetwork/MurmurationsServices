package nodes

import "github.com/MurmurationsNetwork/MurmurationsServices/utils/rest_errors"

type Node struct {
	ID            string   `json:"nodeId"`
	ProfileUrl    string   `json:"profileUrl"`
	LinkedSchemas []string `json:"linkedSchemas"`
	LastValidated int64    `json:"lastValidated"`
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
