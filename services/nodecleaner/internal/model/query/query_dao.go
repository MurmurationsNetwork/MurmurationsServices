package query

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/elastic"
)

func (q *EsQuery) Build() *elastic.Query {
	query := elastic.NewBoolQuery()

	subQueries := elastic.NewQueries()

	if q.Status != nil {
		subQueries = append(
			subQueries,
			elastic.NewMatchQuery("status", *q.Status),
		)
	}

	filters := elastic.NewQueries()
	if q.TimeBefore != nil {
		filters = append(
			filters,
			elastic.NewRangeQuery("last_updated").Lte(*q.TimeBefore),
		)
	}

	if q.Expires != nil {
		filters = append(
			filters,
			elastic.NewRangeQuery("expires").Lte(*q.Expires),
		)
	}

	query.Must(subQueries...).Filter(filters...)

	return &elastic.Query{
		Query: query,
	}
}
