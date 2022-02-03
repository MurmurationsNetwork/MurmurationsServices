package query

import (
	"fmt"
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

	if q.Tags != nil {
		tags := strings.Replace(*q.Tags, ",", " ", -1)
		fmt.Println(tags)
		if q.TagsFilter != nil && *q.TagsFilter == "and" {
			subQueries = append(subQueries, elastic.NewMatchQuery("tags", tags).Operator("AND").Fuzziness(config.Conf.Server.TagsFuzziness))
		} else {
			subQueries = append(subQueries, elastic.NewMatchQuery("tags", tags).Fuzziness(config.Conf.Server.TagsFuzziness))
		}
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
