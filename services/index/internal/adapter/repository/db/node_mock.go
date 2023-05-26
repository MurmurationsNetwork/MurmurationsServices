package db

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/common/jsonapi"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/entity"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/entity/query"
)

type mockNodeRepository struct{}

func (*mockNodeRepository) Add(_ *entity.Node) []jsonapi.Error {
	return nil
}

func (*mockNodeRepository) GetNode(
	nodeID string,
) (*entity.Node, []jsonapi.Error) {
	return &entity.Node{
		ID: nodeID,
	}, nil
}

func (*mockNodeRepository) Get(
	nodeID string,
) (*entity.Node, []jsonapi.Error) {
	return &entity.Node{
		ID: nodeID,
	}, nil
}

func (*mockNodeRepository) Update(_ *entity.Node) error {
	return nil
}

func (*mockNodeRepository) Search(
	_ *query.EsQuery,
) (*query.Results, []jsonapi.Error) {
	return nil, nil
}

func (*mockNodeRepository) Delete(_ *entity.Node) []jsonapi.Error {
	return nil
}

func (*mockNodeRepository) SoftDelete(_ *entity.Node) []jsonapi.Error {
	return nil
}

func (*mockNodeRepository) Export(
	_ *query.EsBlockQuery,
) (*query.BlockQueryResults, []jsonapi.Error) {
	return nil, nil
}

func (*mockNodeRepository) GetNodes(
	_ *query.EsQuery,
) (*query.MapQueryResults, []jsonapi.Error) {
	return nil, nil
}
