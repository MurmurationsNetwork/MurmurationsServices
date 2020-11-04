package query

type EsQuery struct {
	Schema      *string `form:"schema"`
	LastChecked *int64  `form:"lastChecked"`

	Lat    *float64 `form:"lat"`
	Lon    *float64 `form:"lon"`
	Radius *string  `form:"radius"`

	Locality *string `form:"locality"`
	Region   *string `form:"region"`
	Country  *string `form:"country"`
}

type QueryResult map[string]interface{}

type QueryResults []QueryResult
