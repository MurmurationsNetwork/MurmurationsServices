package entity

type Profile struct {
	Cuid           string `json:"cuid"             bson:"cuid,omitempty"`
	Oid            string `json:"oid"              bson:"oid,omitempty"`
	NodeId         string `json:"node_id"          bson:"node_id,omitempty"`
	SourceDataHash string `json:"source_data_hash" bson:"source_data_hash,omitempty"`
}
