package model

type Batch struct {
	UserID         string   `json:"user_id,omitempty"         bson:"user_id,omitempty"`
	Title          string   `json:"title,omitempty"           bson:"title,omitempty"`
	BatchID        string   `json:"batch_id,omitempty"        bson:"batch_id,omitempty"`
	Schemas        []string `json:"schemas,omitempty"         bson:"schemas,omitempty"`
	Status         string   `json:"status,omitempty"          bson:"status,omitempty"`
	TotalNodes     int      `json:"total_nodes,omitempty"     bson:"total_nodes,omitempty"`
	ProcessedNodes int      `json:"processed_nodes,omitempty" bson:"processed_nodes,omitempty"`
	Error          string   `json:"error,omitempty"           bson:"error,omitempty"`
}
