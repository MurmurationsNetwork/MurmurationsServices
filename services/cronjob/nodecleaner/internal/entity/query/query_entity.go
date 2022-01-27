package query

type EsQuery struct {
	Status     *string
	TimeBefore *int64
}

type QueryResult map[string]interface{}

type QueryResults struct {
	Result          []QueryResult
	NumberOfResults int64
	TotalPages      int64
}
