package db

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/jsonapi"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/entity"
)

type nodeDAO struct {
	ID             string           `bson:"_id,omitempty"`
	ProfileURL     string           `bson:"profile_url,omitempty"`
	ProfileHash    *string          `bson:"profile_hash,omitempty"`
	Status         string           `bson:"status,omitempty"`
	LastUpdated    *int64           `bson:"last_updated,omitempty"`
	FailureReasons *[]jsonapi.Error `bson:"failure_reasons,omitempty"`
	Version        *int32           `bson:"__v,omitempty"`
	CreatedAt      int64            `bson:"createdAt,omitempty"`
}

func (r *nodeRepository) toDAO(node *entity.Node) *nodeDAO {
	return &nodeDAO{
		ID:             node.ID,
		ProfileURL:     node.ProfileURL,
		ProfileHash:    node.ProfileHash,
		Status:         node.Status,
		LastUpdated:    node.LastUpdated,
		FailureReasons: node.FailureReasons,
		Version:        node.Version,
		CreatedAt:      node.CreatedAt,
	}
}

func (dao *nodeDAO) toEntity() *entity.Node {
	return &entity.Node{
		ID:             dao.ID,
		ProfileURL:     dao.ProfileURL,
		ProfileHash:    dao.ProfileHash,
		Status:         dao.Status,
		LastUpdated:    dao.LastUpdated,
		FailureReasons: dao.FailureReasons,
		Version:        dao.Version,
		CreatedAt:      dao.CreatedAt,
	}
}
