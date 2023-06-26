package query

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/elastic"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/pagination"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/config"
)

func (q *EsQuery) Build(isMap bool) *elastic.Query {
	query := elastic.NewBoolQuery()

	subQueries := elastic.NewQueries()

	if q.Name != nil {
		subQueries = append(
			subQueries,
			elastic.NewTextQuery("name", *q.Name),
		)
	}
	if q.Schema != nil {
		subQueries = append(
			subQueries,
			elastic.NewWildcardQuery("linked_schemas", *q.Schema+"*"),
		)
	}
	if q.LastUpdated != nil {
		subQueries = append(
			subQueries,
			elastic.NewRangeQuery("last_updated").Gte(q.LastUpdated),
		)
	}
	if q.Locality != nil {
		subQueries = append(
			subQueries,
			elastic.NewTextQuery("locality", *q.Locality),
		)
	}
	if q.Region != nil {
		subQueries = append(
			subQueries,
			elastic.NewTextQuery("region", *q.Region),
		)
	}
	if q.Country != nil {
		subQueries = append(
			subQueries,
			elastic.NewTextQuery("country", *q.Country),
		)
	}
	if q.Status != nil {
		subQueries = append(
			subQueries,
			elastic.NewMatchQuery("status", *q.Status),
		)
	}
	if q.PrimaryURL != nil {
		subQueries = append(
			subQueries,
			elastic.NewMatchQuery("primary_url", *q.PrimaryURL),
		)
	}

	if q.Tags != nil {
		tagQuery := elastic.NewMatchQuery("tags", *q.Tags)
		if q.TagsFilter != nil && *q.TagsFilter == "and" {
			tagQuery = tagQuery.Operator("AND")
		}
		if q.TagsExact != nil && *q.TagsExact == "true" {
			subQueries = append(subQueries, tagQuery.Fuzziness("0"))
		} else {
			subQueries = append(subQueries, tagQuery.Fuzziness(config.Conf.Server.TagsFuzziness))
		}
	}

	filters := elastic.NewQueries()
	if q.Lat != nil && q.Lon != nil && q.Range != nil {
		filters = append(
			filters,
			elastic.NewGeoDistanceQuery("geolocation").
				Lat(*q.Lat).
				Lon(*q.Lon).
				Distance(*q.Range),
		)
	}

	if isMap {
		subQueries = append(subQueries, elastic.NewExistQuery("geolocation"))
	}

	query.Must(subQueries...).Filter(filters...)

	if isMap {
		return &elastic.Query{
			Query: query,
			From:  pagination.From(q.Page, q.PageSize),
			Size:  pagination.MaximumSize(q.PageSize),
		}
	}

	return &elastic.Query{
		Query: query,
		From:  pagination.From(q.Page, q.PageSize),
		Size:  pagination.Size(q.PageSize),
	}
}

func (q *EsBlockQuery) BuildBlock() *elastic.Query {
	query := elastic.NewBoolQuery()

	subQueries := elastic.NewQueries()

	if q.Schema != nil {
		subQueries = append(
			subQueries,
			elastic.NewWildcardQuery("linked_schemas", *q.Schema+"*"),
		)
	}

	query.Must(subQueries...)

	return &elastic.Query{
		Query: query,
		From:  0,
		Size:  pagination.Size(q.PageSize),
	}
}
