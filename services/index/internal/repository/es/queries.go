package es

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/elastic"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/pagination"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/config"
)

type Query struct {
	Name        *string `form:"name"`
	Schema      *string `form:"schema"`
	LastUpdated *int64  `form:"last_updated,default=0"`

	Lat   *float64 `form:"lat"`
	Lon   *float64 `form:"lon"`
	Range *string  `form:"range"`

	Locality *string `form:"locality"`
	Region   *string `form:"region"`
	Country  *string `form:"country"`

	Status *string `form:"status"`

	Tags       *string `form:"tags"`
	TagsFilter *string `form:"tags_filter"`
	TagsExact  *string `form:"tags_exact"`

	PrimaryURL *string `form:"primary_url"`

	Page     int64 `form:"page,default=0"`
	PageSize int64 `form:"page_size,default=30"`
}

func (q *Query) Build(isMap bool) *elastic.Query {
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
			subQueries = append(subQueries, tagQuery.Fuzziness(config.Values.Server.TagsFuzziness))
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

func (q *BlockQuery) BuildBlock() *elastic.Query {
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

type QueryResult map[string]interface{}

type QueryResults struct {
	Result          []QueryResult
	NumberOfResults int64
	TotalPages      int64
}

type BlockQuery struct {
	Schema *string `json:"schema,omitempty"`

	PageSize int64 `json:"page_size"`

	SearchAfter []interface{} `json:"search_after,omitempty"`
}

type BlockQueryResults struct {
	Result []QueryResult
	Sort   []interface{}
}

type MapQueryResults struct {
	Result          [][]interface{}
	NumberOfResults int64
	TotalPages      int64
}
