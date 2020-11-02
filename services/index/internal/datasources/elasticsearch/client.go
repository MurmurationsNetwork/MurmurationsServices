package elasticsearch

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"github.com/cenkalti/backoff"
	"github.com/olivere/elastic"
)

const docType = "_doc"

var (
	Client esClientInterface = &esClient{}
)

type esClientInterface interface {
	Index(string, interface{}) (*elastic.IndexResponse, error)
	IndexWithID(string, string, interface{}) (*elastic.IndexResponse, error)
	Search(string, elastic.Query) (*elastic.SearchResult, error)
	setClient(*elastic.Client)
}

type esClient struct {
	client *elastic.Client
}

func Init() {
	op := func() error {
		log := logger.GetLogger()

		client, err := elastic.NewClient(
			elastic.SetURL(os.Getenv("ELASTICSEARCH_URL")),
			elastic.SetHealthcheckInterval(10*time.Second),
			elastic.SetErrorLog(log),
			elastic.SetInfoLog(log),
		)
		if err != nil {
			return err
		}

		Client.setClient(client)
		return nil
	}
	notify := func(err error, time time.Duration) {
		logger.Info(fmt.Sprintf("trying to re-connect ElasticSearch %s \n", err))
	}
	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = 2 * time.Minute

	err := backoff.RetryNotify(op, b, notify)
	if err != nil {
		logger.Panic("error when trying to ping the ElasticSearch", err)
	}
}

func (c *esClient) setClient(client *elastic.Client) {
	c.client = client
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
		logger.Error(fmt.Sprintf("error when trying to index document in index %s", index), err)
		return nil, err
	}

	return result, nil
}

func (c *esClient) Get(index string, id string) (*elastic.GetResult, error) {
	ctx := context.Background()
	result, err := c.client.Get().
		Index(index).
		Type(docType).
		Id(id).
		Do(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("error when trying to get id %s", id), err)
		return nil, err
	}

	return result, nil
}

func (c *esClient) Search(index string, query elastic.Query) (*elastic.SearchResult, error) {
	ctx := context.Background()
	result, err := c.client.Search(index).
		Query(query).
		RestTotalHitsAsInt(true).
		Do(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("error when trying to search documents in index %s", index), err)
		return nil, err
	}

	return result, nil
}
