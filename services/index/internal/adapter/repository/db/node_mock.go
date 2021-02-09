package db

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/common/resterr"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/entity"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/entity/query"
)

type mockNodeRepository struct{}

func (r *mockNodeRepository) Add(node *entity.Node) resterr.RestErr {
	return nil
}

func (r *mockNodeRepository) Get(nodeID string) (*entity.Node, resterr.RestErr) {
	return &entity.Node{
		ID: nodeID,
	}, nil
}

func (r *mockNodeRepository) Update(node *entity.Node) error {
	return nil
}

func (r *mockNodeRepository) Search(q *query.EsQuery) (*query.QueryResults, resterr.RestErr) {
	return nil, nil
}

func (r *mockNodeRepository) Delete(node *entity.Node) resterr.RestErr {
	return nil
}
