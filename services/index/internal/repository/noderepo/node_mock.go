package noderepo

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/common/resterr"
	model "github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/domain/node"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/domain/query"
)

type mockNode struct{}

func (dao *mockNode) Add(node *model.Node) resterr.RestErr {
	return nil
}

func (dao *mockNode) Get(node *model.Node) resterr.RestErr {
	return nil
}

func (dao *mockNode) Update(node *model.Node) error {
	return nil
}

func (dao *mockNode) Search(q *query.EsQuery) (*query.QueryResults, resterr.RestErr) {
	return nil, nil
}

func (dao *mockNode) Delete(node *model.Node) resterr.RestErr {
	return nil
}
