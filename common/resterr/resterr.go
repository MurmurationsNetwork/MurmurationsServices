package resterr

import "net/http"

type RestErr interface {
	Message() string
	Status() int
}

type restErr struct {
	ErrMessage string        `json:"message,omitempty"`
	ErrStatus  int           `json:"status,omitempty"`
	ErrCauses  []interface{} `json:"causes,omitempty"`
}

func (e restErr) Message() string {
	return e.ErrMessage
}

func (e restErr) Status() int {
	return e.ErrStatus
}

func NewBadRequestError(message string) RestErr {
	return restErr{
		ErrMessage: message,
		ErrStatus:  http.StatusBadRequest,
	}
}

func NewTooManyRequestsError(message string) RestErr {
	return restErr{
		ErrMessage: message,
		ErrStatus:  http.StatusTooManyRequests,
	}
}

func NewInternalServerError(message string, err error) RestErr {
	result := restErr{
		ErrMessage: message,
		ErrStatus:  http.StatusInternalServerError,
	}
	if err != nil {
		result.ErrCauses = append(result.ErrCauses, err.Error())
	}
	return result
}

func NewNotFoundError(message string) RestErr {
	return restErr{
		ErrMessage: message,
	}
}
