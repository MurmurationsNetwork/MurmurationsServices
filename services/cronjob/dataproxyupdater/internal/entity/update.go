package entity

type Update struct {
	Schema       string `json:"schema" bson:"schema,omitempty"`
	LastUpdated  int64  `json:"last_updated" bson:"last_updated,omitempty"`
	HasError     bool   `json:"has_error" bson:"has_error,omitempty"`
	ApiEntry     string `json:"api_entry" bson:"api_entry,omitempty"`
	ErrorMessage string `json:"error_message" bson:"error_message,omitempty"`
	Version      *int32 `json:"-" bson:"__v,omitempty"`
}
