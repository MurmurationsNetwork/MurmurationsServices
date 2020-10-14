package nodes

import (
	"fmt"
	"time"

	"github.com/MurmurationsNetwork/MurmurationsServices/utils/hash"
	"github.com/MurmurationsNetwork/MurmurationsServices/utils/rest_errors"
)

var (
	nodesDB = make(map[string]*Node)
)

func (node *Node) Add() rest_errors.RestErr {
	id := hash.SHA256(node.ProfileUrl)

	node.ID = id
	node.LastValidated = time.Now().Unix()

	nodesDB[id] = node

	return nil
}

func (node *Node) Get() rest_errors.RestErr {
	result := nodesDB[node.ID]
	if result == nil {
		return rest_errors.NewNotFoundError(fmt.Sprintf("node %s not found", node.ID))
	}

	node.ID = result.ID
	node.ProfileUrl = result.ProfileUrl
	node.LinkedSchemas = result.LinkedSchemas
	node.LastValidated = result.LastValidated

	return nil
}

func (node *Node) Search() rest_errors.RestErr {
	return nil
}

func (node *Node) Delete() rest_errors.RestErr {
	return nil
}
