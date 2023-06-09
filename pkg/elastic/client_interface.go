package elastic

import (
	"os"
	"time"

	elastic "github.com/olivere/elastic/v7"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/retry"
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
	Export(string, *Query, []interface{}) (*elastic.SearchResult, error)
	GetNodes(string, *Query) (*elastic.SearchResult, error)

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
		err := retry.Do(
			operation,
			"Trying to re-connect Elasticsearch",
		)
		if err != nil {
			return err
		}
	}

	Client.setClient(client)

	return nil
}
