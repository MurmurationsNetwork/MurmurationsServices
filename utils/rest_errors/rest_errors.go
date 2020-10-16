package rest_errors

import "net/http"

type RestErr interface {
	Message() string
	Status() int
}

type restErr struct {
	ErrMessage string `json:"message"`
	ErrStatus  int    `json:"status"`
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

func NewInternalServerError(message string) RestErr {
	result := restErr{
		ErrMessage: message,
		ErrStatus:  http.StatusInternalServerError,
	}
	return result
}

func NewNotFoundError(message string) RestErr {
	return restErr{
		ErrMessage: message,
		ErrStatus:  http.StatusNotFound,
	}
}
