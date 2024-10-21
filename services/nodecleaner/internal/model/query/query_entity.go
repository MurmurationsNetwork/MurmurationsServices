package query

type EsQuery struct {
	Status     *string
	TimeBefore *int64
	Expires    *int64
}

type Result map[string]interface{}

type Results struct {
	Result          []Result
	NumberOfResults int64
	TotalPages      int64
}
