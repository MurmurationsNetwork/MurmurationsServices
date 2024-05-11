package elastic

import (
	"context"
	"fmt"

	"github.com/olivere/elastic/v7"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
)

type esClient struct {
	client *elastic.Client
}

func (c *esClient) setClient(client *elastic.Client) {
	c.client = client
}

func (c *esClient) GetClient() *elastic.Client {
	return c.client
}

func (c *esClient) Ping() error {
	respond, err := c.client.ClusterHealth().
		Index(constant.ESIndex.Node).
		Do(context.Background())
	if err != nil || respond.Status != "green" {
		return fmt.Errorf("error pinging Elasticsearch: %w", err)
	}
	return nil
}

func (c *esClient) CreateMappings(indices []Index) error {
	for _, index := range indices {
		exists, err := c.client.IndexExists(index.Name).
			Do(context.Background())
		if err != nil {
			return err
		}
		if !exists {
			createIndex, err := c.client.CreateIndex(index.Name).
				BodyString(index.Body).
				Do(context.Background())
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

func (c *esClient) Index(
	index string,
	doc interface{},
) (*elastic.IndexResponse, error) {
	ctx := context.Background()
	result, err := c.client.Index().
		Index(index).
		BodyJson(doc).
		Do(ctx)
	if err != nil {
		logger.Error(
			fmt.Sprintf(
				"Error when trying to index document in Index: %s",
				index,
			),
			err,
		)
		return nil, err
	}

	return result, nil
}

func (c *esClient) IndexWithID(
	index string,
	id string,
	doc interface{},
) (*elastic.IndexResponse, error) {
	ctx := context.Background()
	result, err := c.client.Index().
		Index(index).
		Id(id).
		BodyJson(doc).
		Do(ctx)
	if err != nil {
		logger.Error(
			fmt.Sprintf(
				"Error when trying to index a document in Index: %s",
				index,
			),
			err,
		)
		return nil, err
	}

	return result, nil
}

func (c *esClient) Search(
	index string,
	q *Query,
) (*elastic.SearchResult, error) {
	ctx := context.Background()

	// sort strategy - 1. _score 2. primary_url
	sortQuery1 := elastic.NewFieldSort("_score").Desc()
	sortQuery2 := elastic.NewFieldSort("primary_url")

	result, err := c.client.Search(index).
		TrackTotalHits(true).
		Query(q.Query).
		From(int(q.From)).
		Size(int(q.Size)).
		RestTotalHitsAsInt(true).
		SortBy(sortQuery1, sortQuery2).
		Do(ctx)
	if err != nil {
		logger.Error(
			fmt.Sprintf(
				"Error when trying to search documents in Index: %s",
				index,
			),
			err,
		)
		return nil, err
	}

	return result, nil
}

func (c *esClient) Update(
	index string,
	id string,
	update map[string]interface{},
) error {
	ctx := context.Background()
	_, err := c.client.Update().
		Index(index).
		Id(id).
		Doc(update).
		Do(ctx)
	if err != nil {
		// Don't need to tell the client data doesn't exist.
		if elastic.IsNotFound(err) {
			return nil
		}
		logger.Error(
			fmt.Sprintf(
				"Error when trying to delete a document in Index: %s",
				index,
			),
			err,
		)
		return err
	}
	return nil
}

func (c *esClient) Delete(index string, id string) error {
	ctx := context.Background()
	_, err := c.client.Delete().
		Index(index).
		Id(id).
		Do(ctx)
	if err != nil {
		// Don't need to tell the client data doesn't exist.
		if elastic.IsNotFound(err) {
			return nil
		}
		logger.Error(
			fmt.Sprintf(
				"Error when trying to delete a document in Index: %s",
				index,
			),
			err,
		)
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
		logger.Error("error when trying to delete documents", err)
		return err
	}
	return nil
}

func (c *esClient) Export(
	index string,
	q *Query,
	searchAfter []interface{},
) (*elastic.SearchResult, error) {
	ctx := context.Background()

	sortQuery1 := elastic.NewFieldSort("last_updated")
	sortQuery2 := elastic.NewFieldSort("profile_url")

	result, err := c.client.Search(index).
		TrackTotalHits(true).
		Query(q.Query).
		SearchAfter(searchAfter...).
		From(int(q.From)).
		Size(int(q.Size)).
		RestTotalHitsAsInt(true).
		SortBy(sortQuery1, sortQuery2).
		Do(ctx)
	if err != nil {
		logger.Error(
			fmt.Sprintf(
				"Error when trying to search documents in Index: %s",
				index,
			),
			err,
		)
		return nil, err
	}

	return result, nil
}

func (c *esClient) GetNodes(
	index string,
	q *Query,
) (*elastic.SearchResult, error) {
	ctx := context.Background()

	source := elastic.NewFetchSourceContext(true).
		Include("geolocation", "profile_url")

	// sort strategy - 1. _score 2. primary_url
	sortQuery1 := elastic.NewFieldSort("_score").Desc()
	sortQuery2 := elastic.NewFieldSort("primary_url")

	result, err := c.client.Search(index).
		TrackTotalHits(true).
		Query(q.Query).
		FetchSourceContext(source).
		From(int(q.From)).
		Size(int(q.Size)).
		RestTotalHitsAsInt(true).
		SortBy(sortQuery1, sortQuery2).
		Do(ctx)
	if err != nil {
		logger.Error(
			fmt.Sprintf(
				"Error when trying to search documents in Index: %s",
				index,
			),
			err,
		)
		return nil, err
	}

	return result, nil
}
