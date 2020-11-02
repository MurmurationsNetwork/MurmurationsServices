package query

type respond struct {
	Data interface{} `json:"data,omitempty"`
}

func (results QueryResults) Marshall() interface{} {
	return respond{Data: results}
}
