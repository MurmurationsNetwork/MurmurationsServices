package elastic

import (
	"strings"

	elastic "github.com/olivere/elastic/v7"
)

type Query struct {
	Query elastic.Query
	From  int64
	Size  int64
}

func NewQueries() []elastic.Query {
	return make([]elastic.Query, 0)
}

func NewBoolQuery() *elastic.BoolQuery {
	return elastic.NewBoolQuery()
}

func NewMatchQuery(name, text string) *elastic.MatchQuery {
	return elastic.NewMatchQuery(name, text)
}

func NewRangeQuery(name string) *elastic.RangeQuery {
	return elastic.NewRangeQuery(name)
}

func NewGeoDistanceQuery(name string) *elastic.GeoDistanceQuery {
	return elastic.NewGeoDistanceQuery(name)
}

func NewTextQuery(name, text string) *elastic.BoolQuery {
	q := elastic.NewBoolQuery()
	q.Should(elastic.NewMatchQuery(name, text).Fuzziness("AUTO"))
	q.Should(elastic.NewRegexpQuery(name, ".*"+strings.ToLower(text)+".*"))
	return q
}

func NewWildcardQuery(name, wildcard string) *elastic.WildcardQuery {
	q := elastic.NewWildcardQuery(name, wildcard)
	q.CaseInsensitive(true)
	return q
}

func NewExistQuery(name string) *elastic.ExistsQuery {
	return elastic.NewExistsQuery(name)
}
