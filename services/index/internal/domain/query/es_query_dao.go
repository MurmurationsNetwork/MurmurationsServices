package query

import (
	"github.com/olivere/elastic"
)

func (q *EsQuery) Build() elastic.Query {
	query := elastic.NewBoolQuery()

	subQueries := make([]elastic.Query, 0)

	if q.Schema != nil {
		subQueries = append(subQueries, elastic.NewMatchQuery("linkedSchemas", q.Schema))
	}
	if q.LastChecked != nil {
		subQueries = append(subQueries, elastic.NewRangeQuery("lastChecked").Gte(q.LastChecked))
	}
	if q.Locality != nil {
		subQueries = append(subQueries, elastic.NewMatchQuery("maplocation.locality", *q.Locality))
	}
	if q.Region != nil {
		subQueries = append(subQueries, elastic.NewMatchQuery("maplocation.region", *q.Region))
	}
	if q.Country != nil {
		subQueries = append(subQueries, elastic.NewMatchQuery("maplocation.country", *q.Country))
	}

	filters := make([]elastic.Query, 0)
	if q.Lat != nil && q.Lon != nil && q.Radius != nil {
		filters = append(filters, elastic.NewGeoDistanceQuery("geolocation").Lat(*q.Lat).Lon(*q.Lon).Distance(*q.Radius))
	}

	query.Must(subQueries...).Filter(filters...)

	return query
}
