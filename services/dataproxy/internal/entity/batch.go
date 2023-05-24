package entity

type Batch struct {
	UserId  string   `json:"user_id,omitempty"  bson:"user_id,omitempty"`
	Title   string   `json:"title,omitempty"    bson:"title,omitempty"`
	BatchId string   `json:"batch_id,omitempty" bson:"batch_id,omitempty"`
	Schemas []string `json:"schemas,omitempty"  bson:"schemas,omitempty"`
}
