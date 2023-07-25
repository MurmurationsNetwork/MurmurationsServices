package elastic

import (
	elastic "github.com/olivere/elastic/v7"
)

// QueryBuilder is a utility to help build ElasticSearch queries.
type QueryBuilder struct {
	// Holds the subqueries used in building the main query.
	subQueries []elastic.Query
	// Holds the filters used in building the main query.
	filters []elastic.Query
}

// GetSubQueries returns the slice of sub queries.
func (b *QueryBuilder) GetSubQueries() []elastic.Query {
	return b.subQueries
}

// GetFilters returns the slice of filters.
func (b *QueryBuilder) GetFilters() []elastic.Query {
	return b.filters
}

// AddSubQuery appends a new query to the QueryBuilder's subQueries.
func (b *QueryBuilder) AddSubQuery(query ...elastic.Query) {
	b.subQueries = append(b.subQueries, query...)
}

// AddFilter appends a new filter to the QueryBuilder's filters.
func (b *QueryBuilder) AddFilter(query ...elastic.Query) {
	b.filters = append(b.filters, query...)
}

// BuildTextQuery generates a text query with the given field.
func (b *QueryBuilder) BuildTextQuery(field string, value *string) {
	if value != nil {
		b.AddSubQuery(NewTextQuery(field, *value))
	}
}

// BuildWildcardQuery generates a wildcard query with the given field.
func (b *QueryBuilder) BuildWildcardQuery(field string, value *string) {
	if value != nil {
		b.AddSubQuery(NewWildcardQuery(field, *value+"*"))
	}
}

// BuildRangeQuery generates a range query with the given field.
func (b *QueryBuilder) BuildRangeQuery(field string, value *int64) {
	if value != nil {
		b.AddSubQuery(NewRangeQuery(field).Gte(*value))
	}
}

// BuildMatchQuery generates a match query with the given field.
func (b *QueryBuilder) BuildMatchQuery(field string, value *string) {
	if value != nil {
		b.AddSubQuery(NewMatchQuery(field, *value))
	}
}

// BuildGeoQuery generates a geolocation query with the given latitude, longitude,
// and distance.
func (b *QueryBuilder) BuildGeoQuery(
	lat *float64,
	lon *float64,
	distance *string,
) {
	if lat != nil && lon != nil && distance != nil {
		b.AddFilter(
			NewGeoDistanceQuery(
				"geolocation",
			).Lat(*lat).
				Lon(*lon).
				Distance(*distance),
		)
	}
}
