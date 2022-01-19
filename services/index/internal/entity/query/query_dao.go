package query

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/common/elastic"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/pagination"
)

func (q *EsQuery) Build() *elastic.Query {
	query := elastic.NewBoolQuery()

	subQueries := elastic.NewQueries()

	if q.Schema != nil {
		subQueries = append(subQueries, elastic.NewMatchQuery("linked_schemas", *q.Schema))
	}
	if q.LastUpdated != nil {
		subQueries = append(subQueries, elastic.NewRangeQuery("last_updated").Gte(q.LastUpdated))
	}
	if q.Locality != nil {
		subQueries = append(subQueries, elastic.NewTextQuery("locality", *q.Locality))
	}
	if q.Region != nil {
		subQueries = append(subQueries, elastic.NewTextQuery("region", *q.Region))
	}
	if q.Country != nil {
		subQueries = append(subQueries, elastic.NewTextQuery("country", *q.Country))
	}

	filters := elastic.NewQueries()
	if q.Lat != nil && q.Lon != nil && q.Range != nil {
		filters = append(filters, elastic.NewGeoDistanceQuery("geolocation").Lat(*q.Lat).Lon(*q.Lon).Distance(*q.Range))
	}

	query.Must(subQueries...).Filter(filters...)

	return &elastic.Query{
		Query: query,
		From:  pagination.From(q.Page, q.PageSize),
		Size:  pagination.Size(q.PageSize),
	}
}
