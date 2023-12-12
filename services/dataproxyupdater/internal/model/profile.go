package model

type Profile struct {
	NodeID   string `json:"node_id"   bson:"node_id,omitempty"`
	IsPosted bool   `json:"is_posted" bson:"is_posted,omitempty"`
	Cuid     string `json:"cuid"      bson:"cuid,omitempty"`
	Oid      string `json:"oid"       bson:"oid,omitempty"`
}
