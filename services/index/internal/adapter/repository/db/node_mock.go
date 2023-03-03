package db

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/common/jsonapi"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/entity"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/entity/query"
)

type mockNodeRepository struct{}

func (r *mockNodeRepository) Add(node *entity.Node) []jsonapi.Error {
	return nil
}

func (r *mockNodeRepository) GetNode(nodeID string) (*entity.Node, []jsonapi.Error) {
	return &entity.Node{
		ID: nodeID,
	}, nil
}

func (r *mockNodeRepository) Get(nodeID string) (*entity.Node, []jsonapi.Error) {
	return &entity.Node{
		ID: nodeID,
	}, nil
}

func (r *mockNodeRepository) Update(node *entity.Node) error {
	return nil
}

func (r *mockNodeRepository) Search(q *query.EsQuery) (*query.QueryResults, []jsonapi.Error) {
	return nil, nil
}

func (r *mockNodeRepository) Delete(node *entity.Node) []jsonapi.Error {
	return nil
}

func (r *mockNodeRepository) SoftDelete(node *entity.Node) []jsonapi.Error {
	return nil
}

func (r *mockNodeRepository) Export(q *query.EsBlockQuery) (*query.BlockQueryResults, []jsonapi.Error) {
	return nil, nil
}

func (r *mockNodeRepository) GetNodes(q *query.EsQuery) (*query.MapQueryResults, []jsonapi.Error) {
	return nil, nil
}
