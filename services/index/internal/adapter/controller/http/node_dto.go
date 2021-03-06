package http

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/common/resterr"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/entity"
)

type nodeDTO struct {
	ID             string    `json:"node_id" `
	ProfileURL     string    `json:"profile_url" `
	ProfileHash    *string   `json:"profile_hash" `
	Status         string    `json:"status" `
	LastValidated  *int64    `json:"last_validated" `
	FailureReasons *[]string `json:"failure_reasons" `
}

func (dto *nodeDTO) Validate() resterr.RestErr {
	if dto.ProfileURL == "" {
		return resterr.NewBadRequestError("The profile_url parameter is missing.")
	}
	return nil
}

func toDTO(node *entity.Node) *nodeDTO {
	return &nodeDTO{
		ID:             node.ID,
		ProfileURL:     node.ProfileURL,
		ProfileHash:    node.ProfileHash,
		Status:         node.Status,
		LastValidated:  node.LastValidated,
		FailureReasons: node.FailureReasons,
	}
}

func (dto *nodeDTO) toEntity() *entity.Node {
	return &entity.Node{
		ProfileURL: dto.ProfileURL,
	}
}
