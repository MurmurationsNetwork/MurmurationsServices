package db

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/countries"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/elastic"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/jsonutil"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/mongo"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/pagination"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/tagsfilter"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/validateurl"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/config"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/index"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/model"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/model/query"
)

type NodeRepository interface {
	Add(node *model.Node) error
	GetNode(nodeID string) (*model.Node, error)
	Get(nodeID string) (*model.Node, error)
	Update(node *model.Node) error
	Search(q *query.EsQuery) (*query.Results, error)
	Delete(node *model.Node) error
	SoftDelete(node *model.Node) error
	Export(q *query.EsBlockQuery) (*query.BlockQueryResults, error)
	GetNodes(q *query.EsQuery) (*query.MapQueryResults, error)
}

func NewRepository() NodeRepository {
	return &nodeRepository{}
}

type nodeRepository struct {
}

func (r *nodeRepository) Add(node *model.Node) error {
	filter := bson.M{"_id": node.ID}
	update := bson.M{"$set": node}
	opt := options.FindOneAndUpdate().SetUpsert(true)

	result, err := mongo.Client.FindOneAndUpdate(
		constant.MongoIndex.Node,
		filter,
		update,
		opt,
	)
	if err != nil {
		return index.DatabaseError{
			Message: "Error occurred during node upsert operation",
			Err:     err,
		}
	}

	var updated model.Node
	err = result.Decode(&updated)
	if err != nil {
		return index.DatabaseError{
			Message: "Error occurred during decoding of updated node",
			Err:     err,
		}
	}

	node.Version = updated.Version

	return nil
}

func (r *nodeRepository) GetNode(
	nodeID string,
) (*model.Node, error) {
	filter := bson.M{"_id": nodeID}

	result := mongo.Client.FindOne(constant.MongoIndex.Node, filter)
	if err := result.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, index.NotFoundError{
				Err: err,
			}
		}
		return nil, index.DatabaseError{
			Message: "Error when trying to find a node",
			Err:     err,
		}
	}

	var node model.Node
	err := result.Decode(&node)
	if err != nil {
		return nil, index.DatabaseError{
			Message: "Error when trying to find a node",
			Err:     err,
		}
	}

	return &node, nil
}

func (r *nodeRepository) Get(nodeID string) (*model.Node, error) {
	filter := bson.M{"_id": nodeID}

	result := mongo.Client.FindOne(constant.MongoIndex.Node, filter)
	if err := result.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, index.NotFoundError{
				Err: err,
			}
		}
		logger.Error("Error when trying to find a node", result.Err())
		return nil, index.DatabaseError{
			Message: "Error when trying to find a node",
			Err:     err,
		}
	}

	var node model.Node
	err := result.Decode(&node)
	if err != nil {
		logger.Error("Error when trying to parse database response", err)
		return nil, index.DatabaseError{
			Message: "Error when trying to find a node",
			Err:     err,
		}
	}

	return &node, nil
}

func (r *nodeRepository) Update(node *model.Node) error {
	filter := bson.M{"_id": node.ID, "__v": node.Version}
	// Unset the version to prevent setting it.
	node.Version = nil
	update := bson.M{"$set": node}

	_, err := mongo.Client.FindOneAndUpdate(
		constant.MongoIndex.Node,
		filter,
		update,
	)
	if err != nil {
		// Update the document only if the version matches.
		// If the version does not match, it's an expected concurrent issue.
		if err == mongo.ErrNoDocuments {
			return nil
		}
		return index.DatabaseError{
			Message: "Error when trying to update a node",
			Err:     err,
		}
	}

	// NOTE: Maybe it's better to convert into another event?
	if node.Status == constant.NodeStatus.Validated {
		profileJSON := jsonutil.ToJSON(node.ProfileStr)
		profileJSON["profile_url"] = node.ProfileURL
		profileJSON["last_updated"] = node.LastUpdated

		// if the geolocation is array type, make it as object type for consistent [#208]
		if _, ok := profileJSON["geolocation"].(string); ok {
			g := strings.Split(profileJSON["geolocation"].(string), ",")
			profileJSON["latitude"], err = strconv.ParseFloat(g[0], 64)
			if err != nil {
				return err
			}
			profileJSON["longitude"], err = strconv.ParseFloat(g[1], 64)
			if err != nil {
				return err
			}
		}

		// if we can find latitude and longitude in the root, move them into geolocation [#208]
		if profileJSON["latitude"] != nil || profileJSON["longitude"] != nil {
			geoLocation := make(map[string]interface{})
			if profileJSON["latitude"] != nil {
				geoLocation["lat"] = profileJSON["latitude"]
			} else {
				geoLocation["lat"] = 0
			}
			if profileJSON["longitude"] != nil {
				geoLocation["lon"] = profileJSON["longitude"]
			} else {
				geoLocation["lon"] = 0
			}
			profileJSON["geolocation"] = geoLocation
		}

		if profileJSON["country_iso_3166"] != nil ||
			profileJSON["country_name"] != nil ||
			profileJSON["country"] != nil {
			if profileJSON["country_iso_3166"] != nil {
				profileJSON["country"] = profileJSON["country_iso_3166"]
				delete(profileJSON, "country_iso_3166")
			} else if profileJSON["country"] == nil && profileJSON["country_name"] != nil {
				countryCode, err := countries.FindAlpha2ByName(config.Values.Library.InternalURL+"/v2/countries", profileJSON["country_name"])
				if err != nil {
					return err
				}
				countryStr := fmt.Sprintf("%v", profileJSON["country_name"])
				profileURLStr := fmt.Sprintf("%v", profileJSON["profile_url"])
				if countryCode != "undefined" {
					profileJSON["country"] = countryCode
					fmt.Println("Country code matched: " + countryStr + " = " + countryCode + " --- profile_url: " + profileURLStr)
				} else {
					// can't find countryCode, log to server
					fmt.Println("Country code not found: " + countryStr + " --- profile_url: " + profileURLStr)
				}
			}
		}

		// Default node's status is posted [#217]
		profileJSON["status"] = "posted"

		// Deal with tags [#227]
		arraySize, _ := strconv.Atoi(config.Values.Server.TagsArraySize)
		stringLength, _ := strconv.Atoi(config.Values.Server.TagsStringLength)
		tags, err := tagsfilter.Filter(arraySize, stringLength, node.ProfileStr)
		if err != nil {
			return err
		}

		if tags != nil {
			profileJSON["tags"] = tags
		}

		// validate primary_url [#238]
		if profileJSON["primary_url"] != nil {
			profileJSON["primary_url"], err = validateurl.Validate(
				profileJSON["primary_url"].(string),
			)
			if err != nil {
				return err
			}
		}

		_, err = elastic.Client.IndexWithID(
			constant.ESIndex.Node,
			node.ID,
			profileJSON,
		)
		if err != nil {
			// Fail to parse into Elasticsearch, set the status to 'post_failed'.
			err = r.setPostFailed(node)
			if err != nil {
				return err
			}
		} else {
			// Successfully parse into Elasticsearch, set the status to 'posted'.
			err = r.setPosted(node)
			if err != nil {
				return err
			}
		}
	}

	if node.Status == constant.NodeStatus.ValidationFailed {
		err := elastic.Client.Delete(constant.ESIndex.Node, node.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *nodeRepository) setPostFailed(node *model.Node) error {
	node.Version = nil
	node.Status = constant.NodeStatus.PostFailed

	filter := bson.M{"_id": node.ID}
	update := bson.M{"$set": node}

	_, err := mongo.Client.FindOneAndUpdate(
		constant.MongoIndex.Node,
		filter,
		update,
	)
	if err != nil {
		logger.Error("Error when trying to update a node", err)
		return err
	}

	return nil
}

func (r *nodeRepository) setPosted(node *model.Node) error {
	node.Version = nil
	node.Status = constant.NodeStatus.Posted

	filter := bson.M{"_id": node.ID}
	update := bson.M{"$set": node}

	_, err := mongo.Client.FindOneAndUpdate(
		constant.MongoIndex.Node,
		filter,
		update,
	)
	if err != nil {
		logger.Error("Error when trying to update a node", err)
		return err
	}

	return nil
}

func (r *nodeRepository) Search(
	q *query.EsQuery,
) (*query.Results, error) {
	result, err := elastic.Client.Search(constant.ESIndex.Node, q.Build(false))
	if err != nil {
		return nil, index.DatabaseError{
			Err: err,
		}
	}

	queryResults := make([]query.Result, 0)
	for _, hit := range result.Hits.Hits {
		bytes, _ := hit.Source.MarshalJSON()
		var result query.Result
		if err := json.Unmarshal(bytes, &result); err != nil {
			return nil, index.DatabaseError{
				Err: err,
			}
		}
		queryResults = append(queryResults, result)
	}

	return &query.Results{
		Result:          queryResults,
		NumberOfResults: result.Hits.TotalHits.Value,
		TotalPages: pagination.TotalPages(
			result.Hits.TotalHits.Value,
			q.PageSize,
		),
	}, nil
}

func (r *nodeRepository) Delete(node *model.Node) error {
	filter := bson.M{"_id": node.ID}

	err := mongo.Client.DeleteOne(constant.MongoIndex.Node, filter)
	if err != nil {
		return index.DatabaseError{
			Err: err,
		}
	}
	err = elastic.Client.Delete(constant.ESIndex.Node, node.ID)
	if err != nil {
		return index.DatabaseError{
			Err: err,
		}
	}

	return nil
}

func (r *nodeRepository) SoftDelete(node *model.Node) error {
	err := r.setDeleted(node)
	if err != nil {
		return index.DatabaseError{
			Err: err,
		}
	}

	err = elastic.Client.Update(
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

func (r *nodeRepository) setDeleted(node *model.Node) error {
	node.Version = nil
	node.Status = constant.NodeStatus.Deleted
	currentTime := time.Now().Unix()
	node.LastUpdated = &currentTime

	filter := bson.M{"_id": node.ID}
	update := bson.M{"$set": node}

	_, err := mongo.Client.FindOneAndUpdate(
		constant.MongoIndex.Node,
		filter,
		update,
	)
	if err != nil {
		logger.Error("Error when trying to update a node", err)
		return err
	}

	return nil
}

func (r *nodeRepository) Export(
	q *query.EsBlockQuery,
) (*query.BlockQueryResults, error) {
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

	queryResults := make([]query.Result, 0)
	hitLength := len(result.Hits.Hits)
	var sort []interface{}
	for i, hit := range result.Hits.Hits {
		bytes, _ := hit.Source.MarshalJSON()
		var result query.Result
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

	return &query.BlockQueryResults{
		Result: queryResults,
		Sort:   sort,
	}, nil
}

func (r *nodeRepository) GetNodes(
	q *query.EsQuery,
) (*query.MapQueryResults, error) {
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

	return &query.MapQueryResults{
		Result:          queryResults,
		NumberOfResults: result.Hits.TotalHits.Value,
		TotalPages: pagination.TotalPages(
			result.Hits.TotalHits.Value,
			q.PageSize,
		),
	}, nil
}
