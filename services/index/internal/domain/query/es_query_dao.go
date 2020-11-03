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
	subQueries = append(subQueries, elastic.NewRangeQuery("locality").Gte(q.Locality))
	subQueries = append(subQueries, elastic.NewRangeQuery("region").Gte(q.Region))
	subQueries = append(subQueries, elastic.NewRangeQuery("country").Gte(q.Country))

	filters := make([]elastic.Query, 0)
	filters = append(filters, elastic.NewGeoDistanceQuery("geolocation").Lat(q.Lat).Lon(q.Lon).Distance(q.Radius))

	query.Must(subQueries...).Filter(filters...)

	fmt.Println("==================================")
	fmt.Printf("query %+v \n", query)
	fmt.Println("==================================")

	return query
}
