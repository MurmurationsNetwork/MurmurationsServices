package elastic

import (
	"github.com/olivere/elastic/v7"
	"os"
	"time"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/backoff"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
)

var (
	Client esClientInterface
)

type esClientInterface interface {
	CreateMappings([]Index) error
	Index(string, interface{}) (*elastic.IndexResponse, error)
	IndexWithID(string, string, interface{}) (*elastic.IndexResponse, error)
	Search(string, *Query) (*elastic.SearchResult, error)
	Update(string, string, map[string]interface{}) error
	Delete(string, string) error
	DeleteMany(string, *Query) error

	setClient(*elastic.Client)
}

func init() {
	if os.Getenv("ENV") == "test" {
		Client = &mockClient{}
		return
	}
	Client = &esClient{}
}

func NewClient(url string) error {
	var client *elastic.Client

	if os.Getenv("ENV") != "test" {
		operation := func() error {
			log := logger.GetLogger()

			var err error
			client, err = elastic.NewClient(
				elastic.SetURL(url),
				elastic.SetHealthcheckInterval(10*time.Second),
				elastic.SetErrorLog(log),
				elastic.SetInfoLog(log),
				// If you found any errors in ES, uncomment the following line to see the request and response
				//elastic.SetTraceLog(log),
			)
			if err != nil {
				return err
			}

			return nil
		}
		err := backoff.NewBackoff(operation, "Trying to re-connect ElasticSearch")
		if err != nil {
			return err
		}
	}

	Client.setClient(client)

	return nil
}
