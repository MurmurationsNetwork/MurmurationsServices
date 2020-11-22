package elastic

import (
	"fmt"
	"time"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"github.com/cenkalti/backoff"
	"github.com/olivere/elastic"
)

var (
	Client esClientInterface
)

type esClientInterface interface {
	setClient(*elastic.Client)

	CreateMappings([]Index) error
	Index(string, interface{}) (*elastic.IndexResponse, error)
	IndexWithID(string, string, interface{}) (*elastic.IndexResponse, error)
	Search(string, *Query) (*elastic.SearchResult, error)
	Delete(string, string) error
}

func init() {
	// TODO: Create test environment.
	if true {
		Client = &esClient{}
	} else {
		Client = &mockClient{}
	}
}

func NewClient(url string) error {
	var client *elastic.Client

	if true {
		op := func() error {
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
		notify := func(err error, time time.Duration) {
			logger.Info(fmt.Sprintf("trying to re-connect ElasticSearch %s \n", err))
		}
		b := backoff.NewExponentialBackOff()
		b.MaxElapsedTime = 2 * time.Minute

		err := backoff.RetryNotify(op, b, notify)
		if err != nil {
			return err
		}
	}

	Client.setClient(client)

	return nil
}
