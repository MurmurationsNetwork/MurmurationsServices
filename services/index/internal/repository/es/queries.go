package es

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/elastic"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/pagination"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/config"
)

// Query defines the parameters that can be used to filter and search profiles
// in Elasticsearch.
type Query struct {
	// Name is used to match profiles based on the "name" field.
	Name *string `form:"name"`

	// Schema is used to match profiles linked to a specific schema pattern.
	Schema *string `form:"schema"`

	// LastUpdated is used to filter profiles based on when they were last updated.
	LastUpdated *int64 `form:"last_updated,default=0"`

	// Lat and Lon, along with Range, are used for geo-based queries.
	Lat   *float64 `form:"lat"`
	Lon   *float64 `form:"lon"`
	Range *string  `form:"range"`

	// Locality, Region, and Country are used to filter profiles based on
	// their associated geographical metadata.
	Locality *string `form:"locality"`
	Region   *string `form:"region"`
	Country  *string `form:"country"`

	// Status is used to match profiles based on their "status" field.
	Status *string `form:"status"`

	// Tags, TagsFilter and TagsExact are used to filter profiles based on
	// the "tags" field.
	Tags *string `form:"tags"`
	// TagsFilter, if set to "and", indicates that the "tags" field filter should
	// perform an "AND" operation (i.e., match all supplied tags).
	TagsFilter *string `form:"tags_filter"`
	// TagsExact, if set to true, indicates that the "tags" field filter should perform
	// an exact match.
	TagsExact *string `form:"tags_exact"`

	// PrimaryURL is used to match profiles based on the "primary_url" field.
	PrimaryURL *string `form:"primary_url"`

	// Page and PageSize are used to control the pagination of the search
	// results.
	Page     int64 `form:"page,default=0"`
	PageSize int64 `form:"page_size,default=30"`
}

func (q *Query) Build(isMap bool) *elastic.Query {
	builder := &elastic.QueryBuilder{}

	builder.BuildTextQuery("name", q.Name)
	builder.BuildWildcardQuery("linked_schemas", q.Schema)
	builder.BuildRangeQuery("last_updated", q.LastUpdated)
	builder.BuildTextQuery("locality", q.Locality)
	builder.BuildTextQuery("region", q.Region)
	builder.BuildTextQuery("country", q.Country)
	builder.BuildMatchQuery("status", q.Status)
	builder.BuildMatchQuery("primary_url", q.PrimaryURL)
	builder.BuildGeoQuery(q.Lat, q.Lon, q.Range)

	if q.Tags != nil {
		tagQuery := elastic.NewMatchQuery("tags", *q.Tags)
		if q.TagsFilter != nil && *q.TagsFilter == "and" {
			tagQuery = tagQuery.Operator("AND")
		}
		if q.TagsExact != nil && *q.TagsExact == "true" {
			builder.AddSubQuery(tagQuery.Fuzziness("0"))
		} else {
			builder.AddSubQuery(tagQuery.Fuzziness(config.Values.Server.TagsFuzziness))
		}
	}

	if isMap {
		builder.AddSubQuery(elastic.NewExistQuery("geolocation"))
	}

	query := elastic.NewBoolQuery().
		Must(builder.GetSubQueries()...).
		Filter(builder.GetFilters()...)

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

type QueryResult map[string]interface{}

type QueryResults struct {
	Result          []QueryResult
	NumberOfResults int64
	TotalPages      int64
}

// BlockQuery defines the parameters that can be used to search for blocks in
// Elasticsearch.
type BlockQuery struct {
	// Schema is used to match blocks linked to a specific schema.
	Schema *string `json:"schema,omitempty"`

	// PageSize controls the number of results per page.
	PageSize int64 `json:"page_size"`

	// SearchAfter defines the cursor for pagination and is used in conjunction with PageSize.
	SearchAfter []interface{} `json:"search_after,omitempty"`
}

// BuildBlock constructs an Elasticsearch query based on the BlockQuery parameters.
func (q *BlockQuery) BuildBlock() *elastic.Query {
	builder := &elastic.QueryBuilder{}

	builder.BuildWildcardQuery("linked_schemas", q.Schema)

	query := elastic.NewBoolQuery().Must(builder.GetSubQueries()...)

	return &elastic.Query{
		Query: query,
		From:  0,
		Size:  pagination.Size(q.PageSize),
	}
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
