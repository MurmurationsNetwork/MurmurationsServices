package query

import "github.com/MurmurationsNetwork/MurmurationsServices/common/elastic"

func (q *EsQuery) Build() elastic.Query {
	query := elastic.NewBoolQuery()

	subQueries := elastic.NewQueries()

	if q.Schema != nil {
		subQueries = append(subQueries, elastic.NewMatchQuery("linked_schemas", *q.Schema))
	}
	if q.LastValidated != nil {
		subQueries = append(subQueries, elastic.NewRangeQuery("last_validated").Gte(q.LastValidated))
	}
	if q.Locality != nil {
		subQueries = append(subQueries, elastic.NewTextQuery("location.locality", *q.Locality))
	}
	if q.Region != nil {
		subQueries = append(subQueries, elastic.NewTextQuery("location.region", *q.Region))
	}
	if q.Country != nil {
		subQueries = append(subQueries, elastic.NewTextQuery("location.country", *q.Country))
	}

	filters := elastic.NewQueries()
	if q.Lat != nil && q.Lon != nil && q.Range != nil {
		filters = append(filters, elastic.NewGeoDistanceQuery("geolocation").Lat(*q.Lat).Lon(*q.Lon).Distance(*q.Range))
	}

	query.Must(subQueries...).Filter(filters...)

	return query
}
