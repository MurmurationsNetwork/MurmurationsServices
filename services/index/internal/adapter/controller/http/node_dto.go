package http

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/common/resterr"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/entity"
	"net/url"
	"strings"
)

type nodeDTO struct {
	ID             string    `json:"node_id" `
	ProfileURL     string    `json:"profile_url" `
	ProfileHash    *string   `json:"profile_hash" `
	Status         string    `json:"status" `
	LastUpdated    *int64    `json:"last_updated" `
	FailureReasons *[]string `json:"failure_reasons" `
}

func (dto *nodeDTO) Validate() resterr.RestErr {
	if dto.ProfileURL == "" {
		return resterr.NewBadRequestError("The profile_url parameter is missing.")
	}
	u, err := url.Parse(dto.ProfileURL)
	// count '.' in the hostname to filter invalid hostname with zero dot, for example: https://blah is invalid
	uCount := strings.Count(u.Host, ".")
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" || uCount == 0 {
		return resterr.NewBadRequestError("The profile_url is invalid.")
	}
	return nil
}

func toDTO(node *entity.Node) *nodeDTO {
	return &nodeDTO{
		ID:             node.ID,
		ProfileURL:     node.ProfileURL,
		ProfileHash:    node.ProfileHash,
		Status:         node.Status,
		LastUpdated:    node.LastUpdated,
		FailureReasons: node.FailureReasons,
	}
}

func (dto *nodeDTO) toEntity() *entity.Node {
	return &entity.Node{
		ProfileURL: dto.ProfileURL,
	}
}
