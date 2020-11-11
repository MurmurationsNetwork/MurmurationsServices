package query

type EsQuery struct {
	Schema      *string `form:"schema"`
	LastValidated *int64  `form:"last_validated"`

	Lat    *float64 `form:"lat"`
	Lon    *float64 `form:"lon"`
	Range *string  `form:"range"`

	Locality *string `form:"locality"`
	Region   *string `form:"region"`
	Country  *string `form:"country"`
}

type QueryResult map[string]interface{}

type QueryResults []QueryResult
