package query

type EsQuery struct {
	Schema        *string `form:"schema"`
	LastValidated *int64  `form:"last_validated,default=0"`

	Lat   *float64 `form:"lat"`
	Lon   *float64 `form:"lon"`
	Range *string  `form:"range"`

	Locality *string `form:"locality"`
	Region   *string `form:"region"`
	Country  *string `form:"country"`

	Page     int64 `form:"page,default=0"`
	PageSize int64 `form:"page_size,default=30"`
}

type QueryResult map[string]interface{}

type QueryResults struct {
	Result          []QueryResult
	NumberOfResults int64
	TotalPages      int64
}
