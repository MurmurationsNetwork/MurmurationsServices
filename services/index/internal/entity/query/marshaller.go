package query

type meta struct {
	NumberOfResults int64 `json:"number_of_results"`
	TotalPages      int64 `json:"total_pages"`
}

type respond struct {
	Data interface{} `json:"data,omitempty"`
	Meta meta        `json:"meta,omitempty"`
}

func (results QueryResults) Marshall() interface{} {
	return respond{
		Data: results.Result,
		Meta: meta{
			TotalPages:      results.TotalPages,
			NumberOfResults: results.NumberOfResults,
		},
	}
}
