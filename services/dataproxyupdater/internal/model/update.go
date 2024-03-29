package model

const (
	ErrorStatusOK             = 0
	ErrorStatusAPIUnavailable = 1
)

type Update struct {
	Schema       string `json:"schema"        bson:"schema,omitempty"`
	LastUpdated  int64  `json:"last_updated"  bson:"last_updated,omitempty"`
	HasError     bool   `json:"has_error"     bson:"has_error,omitempty"`
	APIEntry     string `json:"api_entry"     bson:"api_entry,omitempty"`
	ErrorMessage string `json:"error_message" bson:"error_message,omitempty"`
	ErrorStatus  int    `json:"error_status"  bson:"error_status,omitempty"`
	Version      *int32 `json:"-"             bson:"__v,omitempty"`
}
