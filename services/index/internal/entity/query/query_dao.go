package query

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/common/elastic"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/pagination"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/config"
	"strings"
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
	if q.Status != nil {
		subQueries = append(subQueries, elastic.NewMatchQuery("status", *q.Status))
	}
	if q.PrimaryUrl != nil {
		subQueries = append(subQueries, elastic.NewMatchQuery("primary_url", *q.PrimaryUrl))
	}

	if q.Tags != nil {
		tagsQueries := elastic.NewQueries()
		tags := strings.Replace(*q.Tags, ",", " ", -1)
		tagQuery := elastic.NewMatchQuery("tags", tags)
		if q.TagsFilter != nil && *q.TagsFilter == "and" {
			tagQuery = tagQuery.Operator("AND")
			tagsQueries = append(tagsQueries, tagQuery.Fuzziness("0").Boost(3))
			if !(q.TagsExact != nil && *q.TagsExact == "true") {
				tagsQueries = append(tagsQueries, tagQuery.Fuzziness(config.Conf.Server.TagsFuzziness))
			}
		} else {
			tagsQueries = append(tagsQueries, tagQuery.Fuzziness("0").Boost(3))
			if !(q.TagsExact != nil && *q.TagsExact == "true") {
				tagsQueries = append(tagsQueries, tagQuery.Fuzziness(config.Conf.Server.TagsFuzziness))
			}
		}
		query.Should(tagsQueries...).MinimumShouldMatch("1")
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
