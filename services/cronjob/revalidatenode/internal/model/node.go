package model

type Node struct {
	ProfileURL string `json:"profile_url" bson:"profile_url,omitempty"`
	Status     string `json:"status"      bson:"status,omitempty"`
	Version    *int32 `json:"-"           bson:"__v,omitempty"`
}

type Nodes []*Node
