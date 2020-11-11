package query

import (
	"github.com/olivere/elastic"
)

func (q *EsQuery) Build() elastic.Query {
	query := elastic.NewBoolQuery()

	subQueries := make([]elastic.Query, 0)

	if q.Schema != nil {
		subQueries = append(subQueries, elastic.NewMatchQuery("linked_schemas", q.Schema))
	}
	if q.LastValidated != nil {
		subQueries = append(subQueries, elastic.NewRangeQuery("last_validated").Gte(q.LastValidated))
	}
	if q.Locality != nil {
		subQueries = append(subQueries, elastic.NewMatchQuery("location.locality", *q.Locality))
	}
	if q.Region != nil {
		subQueries = append(subQueries, elastic.NewMatchQuery("location.region", *q.Region))
	}
	if q.Country != nil {
		subQueries = append(subQueries, elastic.NewMatchQuery("location.country", *q.Country))
	}

	filters := make([]elastic.Query, 0)
	if q.Lat != nil && q.Lon != nil && q.Range != nil {
		filters = append(filters, elastic.NewGeoDistanceQuery("geolocation").Lat(*q.Lat).Lon(*q.Lon).Distance(*q.Range))
	}

	query.Must(subQueries...).Filter(filters...)

	return query
}
