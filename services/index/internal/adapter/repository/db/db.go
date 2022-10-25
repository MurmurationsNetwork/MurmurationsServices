package db

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/jsonapi"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/tagsfilter"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/validateurl"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/config"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/countries"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/elastic"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/jsonutil"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/mongo"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/pagination"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/entity"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/entity/query"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type NodeRepository interface {
	Add(node *entity.Node) []jsonapi.Error
	Get(nodeID string) (*entity.Node, []jsonapi.Error)
	Update(node *entity.Node) error
	Search(q *query.EsQuery) (*query.QueryResults, []jsonapi.Error)
	Delete(node *entity.Node) []jsonapi.Error
	SoftDelete(node *entity.Node) []jsonapi.Error
}

func NewRepository() NodeRepository {
	if os.Getenv("ENV") == "test" {
		return &mockNodeRepository{}
	}
	return &nodeRepository{}
}

type nodeRepository struct {
}

func (r *nodeRepository) Add(node *entity.Node) []jsonapi.Error {
	filter := bson.M{"_id": node.ID}
	update := bson.M{"$set": r.toDAO(node)}
	opt := options.FindOneAndUpdate().SetUpsert(true)

	result, err := mongo.Client.FindOneAndUpdate(constant.MongoIndex.Node, filter, update, opt)
	if err != nil {
		logger.Error("Error when trying to create a node", err)
		return jsonapi.NewError([]string{"Database Error"}, []string{"Error when trying to add a node."}, nil, []int{http.StatusInternalServerError})
	}

	var updated nodeDAO
	result.Decode(&updated)
	node.Version = updated.Version

	return nil
}

func (r *nodeRepository) Get(nodeID string) (*entity.Node, []jsonapi.Error) {
	filter := bson.M{"_id": nodeID}

	result := mongo.Client.FindOne(constant.MongoIndex.Node, filter)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return nil, jsonapi.NewError([]string{"Node Not Found"}, []string{fmt.Sprintf("Could not locate the following node_id in the Index: %s", nodeID)}, nil, []int{http.StatusNotFound})
		}
		logger.Error("Error when trying to find a node", result.Err())
		return nil, jsonapi.NewError([]string{"Database Error"}, []string{"Error when trying to find a node."}, nil, []int{http.StatusInternalServerError})
	}

	var node nodeDAO
	err := result.Decode(&node)
	if err != nil {
		logger.Error("Error when trying to parse database response", result.Err())
		return nil, jsonapi.NewError([]string{"Database Error"}, []string{"Error when trying to find a node."}, nil, []int{http.StatusInternalServerError})
	}

	return node.toEntity(), nil
}

func (r *nodeRepository) Update(node *entity.Node) error {
	filter := bson.M{"_id": node.ID, "__v": node.Version}
	// Unset the version to prevent setting it.
	node.Version = nil
	update := bson.M{"$set": r.toDAO(node)}

	_, err := mongo.Client.FindOneAndUpdate(constant.MongoIndex.Node, filter, update)
	if err != nil {
		// Update the document only if the version matches.
		// If the version does not match, it's an expected concurrent issue.
		if err == mongo.ErrNoDocuments {
			return nil
		}
		logger.Error("Error when trying to update a node", err)
		return ErrUpdate
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

		if profileJSON["country_iso_3166"] != nil || profileJSON["country_name"] != nil || profileJSON["country"] != nil {
			if profileJSON["country_iso_3166"] != nil {
				profileJSON["country"] = profileJSON["country_iso_3166"]
				delete(profileJSON, "country_iso_3166")
			} else if profileJSON["country"] == nil && profileJSON["country_name"] != nil {
				countryCode, err := countries.FindAlpha2ByName(profileJSON["country_name"])
				if err != nil {
					return err
				}
				countryStr := fmt.Sprintf("%v", profileJSON["country_name"])
				profileUrlStr := fmt.Sprintf("%v", profileJSON["profile_url"])
				if countryCode != "undefined" {
					profileJSON["country"] = countryCode
					fmt.Println("Country code matched: " + countryStr + " = " + countryCode + " --- profile_url: " + profileUrlStr)
				} else {
					// can't find countryCode, log to server
					fmt.Println("Country code not found: " + countryStr + " --- profile_url: " + profileUrlStr)
				}
			}
		}

		// Default node's status is posted [#217]
		profileJSON["status"] = "posted"

		// Deal with tags [#227]
		arraySize, _ := strconv.Atoi(config.Conf.Server.TagsArraySize)
		stringLength, _ := strconv.Atoi(config.Conf.Server.TagsStringLength)
		tags, err := tagsfilter.Filter(arraySize, stringLength, node.ProfileStr)
		if err != nil {
			return err
		}

		if tags != nil {
			profileJSON["tags"] = tags
		}

		// validate primary_url [#238]
		if profileJSON["primary_url"] != nil {
			profileJSON["primary_url"], err = validateurl.Validate(profileJSON["primary_url"].(string))
			if err != nil {
				return err
			}
		}

		_, err = elastic.Client.IndexWithID(constant.ESIndex.Node, node.ID, profileJSON)
		if err != nil {
			// Fail to parse into ElasticSearch, set the status to 'post_failed'.
			err = r.setPostFailed(node)
			if err != nil {
				return err
			}
		} else {
			// Successfully parse into ElasticSearch, set the statue to 'posted'.
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

func (r *nodeRepository) setPostFailed(node *entity.Node) error {
	node.Version = nil
	node.Status = constant.NodeStatus.PostFailed

	filter := bson.M{"_id": node.ID}
	update := bson.M{"$set": r.toDAO(node)}

	_, err := mongo.Client.FindOneAndUpdate(constant.MongoIndex.Node, filter, update)
	if err != nil {
		logger.Error("Error when trying to update a node", err)
		return err
	}

	return nil
}

func (r *nodeRepository) setPosted(node *entity.Node) error {
	node.Version = nil
	node.Status = constant.NodeStatus.Posted

	filter := bson.M{"_id": node.ID}
	update := bson.M{"$set": r.toDAO(node)}

	_, err := mongo.Client.FindOneAndUpdate(constant.MongoIndex.Node, filter, update)
	if err != nil {
		logger.Error("Error when trying to update a node", err)
		return err
	}

	return nil
}

func (r *nodeRepository) Search(q *query.EsQuery) (*query.QueryResults, []jsonapi.Error) {
	result, err := elastic.Client.Search(constant.ESIndex.Node, q.Build())
	if err != nil {
		return nil, jsonapi.NewError([]string{"Database Error"}, []string{"Error when trying to search documents."}, nil, []int{http.StatusInternalServerError})
	}

	queryResults := make([]query.QueryResult, 0)
	for _, hit := range result.Hits.Hits {
		bytes, _ := hit.Source.MarshalJSON()
		var result query.QueryResult
		if err := json.Unmarshal(bytes, &result); err != nil {
			return nil, jsonapi.NewError([]string{"Database Error"}, []string{"Error when trying to search documents."}, nil, []int{http.StatusInternalServerError})
		}
		queryResults = append(queryResults, result)
	}

	return &query.QueryResults{
		Result:          queryResults,
		NumberOfResults: result.Hits.TotalHits.Value,
		TotalPages:      pagination.TotalPages(result.Hits.TotalHits.Value, q.PageSize),
	}, nil
}

func (r *nodeRepository) Delete(node *entity.Node) []jsonapi.Error {
	filter := bson.M{"_id": node.ID}

	err := mongo.Client.DeleteOne(constant.MongoIndex.Node, filter)
	if err != nil {
		return jsonapi.NewError([]string{"Database Error"}, []string{"Error when trying to delete a node."}, nil, []int{http.StatusInternalServerError})
	}
	err = elastic.Client.Delete(constant.ESIndex.Node, node.ID)
	if err != nil {
		return jsonapi.NewError([]string{"Database Error"}, []string{"Error when trying to delete a node."}, nil, []int{http.StatusInternalServerError})
	}

	return nil
}

func (r *nodeRepository) SoftDelete(node *entity.Node) []jsonapi.Error {
	err := r.setDeleted(node)
	if err != nil {
		return jsonapi.NewError([]string{"Database Error"}, []string{"Error when trying to delete a node."}, nil, []int{http.StatusInternalServerError})
	}

	err = elastic.Client.Update(constant.ESIndex.Node, node.ID, map[string]interface{}{"status": "deleted", "last_updated": node.LastUpdated})
	if err != nil {
		return jsonapi.NewError([]string{"Database Error"}, []string{"Error when trying to delete a node."}, nil, []int{http.StatusInternalServerError})
	}

	return nil
}

func (r *nodeRepository) setDeleted(node *entity.Node) error {
	node.Version = nil
	node.Status = constant.NodeStatus.Deleted
	currentTime := time.Now().Unix()
	node.LastUpdated = &currentTime

	filter := bson.M{"_id": node.ID}
	update := bson.M{"$set": r.toDAO(node)}

	_, err := mongo.Client.FindOneAndUpdate(constant.MongoIndex.Node, filter, update)
	if err != nil {
		logger.Error("Error when trying to update a node", err)
		return err
	}

	return nil
}
