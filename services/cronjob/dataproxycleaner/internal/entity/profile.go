package entity

type Profile struct {
	Cuid   string `json:"cuid" bson:"cuid,omitempty"`
	NodeId string `json:"node_id" bson:"node_id,omitempty"`
}
