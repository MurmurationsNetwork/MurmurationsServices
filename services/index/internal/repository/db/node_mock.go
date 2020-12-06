package db

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/common/resterr"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/domain/node"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/domain/query"
)

type mockNodeRepository struct{}

func (r *mockNodeRepository) Add(node *node.Node) resterr.RestErr {
	return nil
}

func (r *mockNodeRepository) Get(node *node.Node) resterr.RestErr {
	return nil
}

func (r *mockNodeRepository) Update(node *node.Node) error {
	return nil
}

func (r *mockNodeRepository) Search(q *query.EsQuery) (*query.QueryResults, resterr.RestErr) {
	return nil, nil
}

func (r *mockNodeRepository) Delete(node *node.Node) resterr.RestErr {
	return nil
}
