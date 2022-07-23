package entity

type Profile struct {
	Cuid string `json:"cuid" bson:"cuid,omitempty"`
	Oid  string `json:"oid" bson:"oid,omitempty"`
}
