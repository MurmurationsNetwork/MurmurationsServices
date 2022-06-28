package elastic

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
)

type esClient struct {
	client *elastic.Client
}

func (c *esClient) setClient(client *elastic.Client) {
	c.client = client
}

func (c *esClient) CreateMappings(indices []Index) error {
	for _, index := range indices {
		exists, err := c.client.IndexExists(string(index.Name)).Do(context.Background())
		if err != nil {
			return err
		}
		if !exists {
			createIndex, err := c.client.CreateIndex(string(index.Name)).BodyString(index.Body).Do(context.Background())
			if err != nil {
				return err
			}
			if !createIndex.Acknowledged {
				return err
			}
		}
	}
	return nil
}

func (c *esClient) Index(index string, doc interface{}) (*elastic.IndexResponse, error) {
	ctx := context.Background()
	result, err := c.client.Index().
		Index(index).
		Type(docType).
		BodyJson(doc).
		Do(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("error when trying to index document in index %s", index), err)
		return nil, err
	}

	return result, nil
}

func (c *esClient) IndexWithID(index string, id string, doc interface{}) (*elastic.IndexResponse, error) {
	ctx := context.Background()
	result, err := c.client.Index().
		Index(index).
		Id(id).
		Type("_doc").
		BodyJson(doc).
		Do(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("error when trying to index a document in index %s", index), err)
		return nil, err
	}

	return result, nil
}

func (c *esClient) Search(index string, q *Query) (*elastic.SearchResult, error) {
	ctx := context.Background()
	result, err := c.client.Search(index).
		TrackTotalHits(true).
		Query(q.Query).
		From(int(q.From)).
		Size(int(q.Size)).
		RestTotalHitsAsInt(true).
		Do(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("error when trying to search documents in index %s", index), err)
		return nil, err
	}

	return result, nil
}

func (c *esClient) Update(index string, id string, update map[string]interface{}) error {
	ctx := context.Background()
	_, err := c.client.Update().
		Index(index).
		Type(docType).
		Id(id).
		Doc(update).
		Do(ctx)
	if err != nil {
		// Don't need to tell the client data doesn't exist.
		if elastic.IsNotFound(err) {
			return nil
		}
		logger.Error(fmt.Sprintf("error when trying to delete a document in index %s", index), err)
		return err
	}
	return nil
}

func (c *esClient) Delete(index string, id string) error {
	ctx := context.Background()
	_, err := c.client.Delete().
		Index(index).
		Type(docType).
		Id(id).
		Do(ctx)
	if err != nil {
		// Don't need to tell the client data doesn't exist.
		if elastic.IsNotFound(err) {
			return nil
		}
		logger.Error(fmt.Sprintf("error when trying to delete a document in index %s", index), err)
		return err
	}
	return nil
}

func (c *esClient) DeleteMany(index string, q *Query) error {
	ctx := context.Background()
	_, err := c.client.DeleteByQuery().
		Index(index).
		Query(q.Query).
		Do(ctx)
	if err != nil {
		// Don't need to tell the client data doesn't exist.
		if elastic.IsNotFound(err) {
			return nil
		}
		logger.Error(fmt.Sprintf("error when trying to delete documents"), err)
		return err
	}
	return nil
}
