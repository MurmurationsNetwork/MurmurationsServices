package query

type EsQuery struct {
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

type Result map[string]interface{}

type Results struct {
	Result          []Result
	NumberOfResults int64
	TotalPages      int64
}

type EsBlockQuery struct {
	Schema *string `json:"schema,omitempty"`

	PageSize int64 `json:"page_size"`

	SearchAfter []interface{} `json:"search_after,omitempty"`
}

type BlockQueryResults struct {
	Result []Result
	Sort   []interface{}
}

type MapQueryResults struct {
	Result          [][]interface{}
	NumberOfResults int64
	TotalPages      int64
}
