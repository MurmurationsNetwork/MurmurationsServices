package query

import (
	"fmt"

	"github.com/olivere/elastic"
)

func (q *EsQuery) Build() elastic.Query {
	query := elastic.NewBoolQuery()

	fmt.Println("==================================")
	fmt.Printf("q %+v \n", q)
	fmt.Println("==================================")

	subQueries := make([]elastic.Query, 0)
	subQueries = append(subQueries, elastic.NewMatchQuery("linkedSchemas", q.Schema))
	subQueries = append(subQueries, elastic.NewRangeQuery("lastChecked").Gte(q.LastChecked))

	query.Must(subQueries...)

	return query
}
