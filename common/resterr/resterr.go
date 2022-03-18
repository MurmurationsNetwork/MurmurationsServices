package resterr

import "net/http"

type RestErr interface {
	Message() string
	Status() int
	StatusCode() int
}

type restErr struct {
	ErrMessage    string        `json:"message,omitempty"`
	ErrStatus     int           `json:"status,omitempty"`
	ErrStatusCode int           `json:"-"`
	ErrCauses     []interface{} `json:"causes,omitempty"`
}

func (e restErr) Message() string {
	return e.ErrMessage
}

func (e restErr) Status() int {
	return e.ErrStatus
}

func (e restErr) StatusCode() int {
	return e.ErrStatusCode
}

func NewBadRequestError(message string) RestErr {
	return restErr{
		ErrMessage:    message,
		ErrStatus:     http.StatusBadRequest,
		ErrStatusCode: http.StatusOK,
	}
}

func NewTooManyRequestsError(message string) RestErr {
	return restErr{
		ErrMessage:    message,
		ErrStatus:     http.StatusTooManyRequests,
		ErrStatusCode: http.StatusTooManyRequests,
	}
}

func NewInternalServerError(message string, err error) RestErr {
	result := restErr{
		ErrMessage:    message,
		ErrStatus:     http.StatusInternalServerError,
		ErrStatusCode: http.StatusInternalServerError,
	}
	if err != nil {
		result.ErrCauses = append(result.ErrCauses, err.Error())
	}
	return result
}

func NewNotFoundError(message string) RestErr {
	return restErr{
		ErrMessage:    message,
		ErrStatus:     http.StatusNotFound,
		ErrStatusCode: http.StatusOK,
	}
}
