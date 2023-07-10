package model

type Updates struct {
	HasError     bool   `json:"has_error"     bson:"has_error,omitempty"`
	ErrorMessage string `json:"error_message" bson:"error_message,omitempty"`
}
