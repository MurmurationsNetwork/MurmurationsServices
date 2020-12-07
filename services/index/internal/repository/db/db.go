package db

import (
	"os"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/resterr"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/domain/node"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/domain/query"
)

type NodeRepository interface {
	Add(node *node.Node) resterr.RestErr
	Get(node *node.Node) resterr.RestErr
	Update(node *node.Node) error
	Search(q *query.EsQuery) (*query.QueryResults, resterr.RestErr)
	Delete(node *node.Node) resterr.RestErr
}

func NewRepository() NodeRepository {
	if os.Getenv("ENV") == "test" {
		return &mockNodeRepository{}
	}
	return &nodeRepository{}
}
