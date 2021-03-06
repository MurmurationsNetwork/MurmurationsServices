package elastic

import (
	"os"
	"time"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/backoff"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"github.com/olivere/elastic"
)

var (
	Client esClientInterface
)

type esClientInterface interface {
	CreateMappings([]Index) error
	Index(string, interface{}) (*elastic.IndexResponse, error)
	IndexWithID(string, string, interface{}) (*elastic.IndexResponse, error)
	Search(string, *Query) (*elastic.SearchResult, error)
	Delete(string, string) error

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
