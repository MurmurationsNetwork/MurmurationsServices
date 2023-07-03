package es

import (
	"encoding/json"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/elastic"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/pagination"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/index"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/model"
)

type NodeRepository interface {
	IndexByID(id string, json interface{}) error
	GetNodes(q *Query) (*MapQueryResults, error)
	Search(q *Query) (*QueryResults, error)
	DeleteByID(id string) error
	SoftDelete(node *model.Node) error
	Export(q *BlockQuery) (*BlockQueryResults, error)
}

func NewNodeRepository() NodeRepository {
	return &nodeRepository{}
}

type nodeRepository struct {
}

func (r *nodeRepository) IndexByID(id string, json interface{}) error {
	_, err := elastic.Client.IndexWithID(
		constant.ESIndex.Node,
		id,
		json,
	)
	return err
}

func (r *nodeRepository) GetNodes(q *Query) (*MapQueryResults, error) {
	result, err := elastic.Client.GetNodes(constant.ESIndex.Node, q.Build(true))
	if err != nil {
		return nil, index.DatabaseError{
			Err: err,
		}
	}

	queryResults := make([][]interface{}, 0)
	for _, hit := range result.Hits.Hits {
		bytes, _ := hit.Source.MarshalJSON()
		var result map[string]interface{}
		if err := json.Unmarshal(bytes, &result); err != nil {
			return nil, index.DatabaseError{
				Err: err,
			}
		}
		// create specific format for map (issue-405)
		// [lon, lat, profile_url]
		geolocation := result["geolocation"].(map[string]interface{})
		mapResult := []interface{}{
			geolocation["lon"],
			geolocation["lat"],
			result["profile_url"],
		}
		queryResults = append(queryResults, mapResult)
	}

	return &MapQueryResults{
		Result:          queryResults,
		NumberOfResults: result.Hits.TotalHits.Value,
		TotalPages: pagination.TotalPages(
			result.Hits.TotalHits.Value,
			q.PageSize,
		),
	}, nil
}

func (r *nodeRepository) Search(q *Query) (*QueryResults, error) {
	result, err := elastic.Client.Search(constant.ESIndex.Node, q.Build(false))
	if err != nil {
		return nil, index.DatabaseError{
			Err: err,
		}
	}

	queryResults := make([]QueryResult, 0)
	for _, hit := range result.Hits.Hits {
		bytes, _ := hit.Source.MarshalJSON()
		var result QueryResult
		if err := json.Unmarshal(bytes, &result); err != nil {
			return nil, index.DatabaseError{
				Err: err,
			}
		}
		queryResults = append(queryResults, result)
	}

	return &QueryResults{
		Result:          queryResults,
		NumberOfResults: result.Hits.TotalHits.Value,
		TotalPages: pagination.TotalPages(
			result.Hits.TotalHits.Value,
			q.PageSize,
		),
	}, nil
}

func (r *nodeRepository) DeleteByID(id string) error {
	return elastic.Client.Delete(constant.ESIndex.Node, id)
}

func (r *nodeRepository) SoftDelete(node *model.Node) error {
	err := elastic.Client.Update(
		constant.ESIndex.Node,
		node.ID,
		map[string]interface{}{
			"status":       "deleted",
			"last_updated": node.LastUpdated,
		},
	)
	if err != nil {
		return index.DatabaseError{
			Err: err,
		}
	}
	return nil
}

func (r *nodeRepository) Export(q *BlockQuery) (*BlockQueryResults, error) {
	result, err := elastic.Client.Export(
		constant.ESIndex.Node,
		q.BuildBlock(),
		q.SearchAfter,
	)
	if err != nil {
		return nil, index.DatabaseError{
			Err: err,
		}
	}

	queryResults := make([]QueryResult, 0)
	hitLength := len(result.Hits.Hits)
	var sort []interface{}
	for i, hit := range result.Hits.Hits {
		bytes, _ := hit.Source.MarshalJSON()
		var result QueryResult
		if err := json.Unmarshal(bytes, &result); err != nil {
			return nil, index.DatabaseError{
				Err: err,
			}
		}
		queryResults = append(queryResults, result)
		// get sort: only get the last item
		if i == hitLength-1 {
			sort = hit.Sort
		}
	}

	return &BlockQueryResults{
		Result: queryResults,
		Sort:   sort,
	}, nil
}
