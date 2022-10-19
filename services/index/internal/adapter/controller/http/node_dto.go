package http

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/common/jsonapi"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/entity"
	"net/http"
	"net/url"
	"strings"
)

type nodeDTO struct {
	ID             string           `json:"node_id" `
	ProfileURL     string           `json:"profile_url" `
	ProfileHash    *string          `json:"profile_hash" `
	Status         string           `json:"status" `
	LastUpdated    *int64           `json:"last_updated" `
	FailureReasons *[]jsonapi.Error `json:"failure_reasons" `
}

func (dto *nodeDTO) Validate() []jsonapi.Error {
	if dto.ProfileURL == "" {
		return jsonapi.NewError([]string{"Missing Required Property"}, []string{"The `profile_url` property is required."}, nil, []int{http.StatusBadRequest})
	}
	u, err := url.Parse(dto.ProfileURL)
	// count '.' in the hostname to filter invalid hostname with zero dot, for example: https://blah is invalid
	uCount := strings.Count(u.Host, ".")
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" || uCount == 0 {
		return jsonapi.NewError([]string{"Invalid Profile URL"}, []string{"The `profile_url` is not a valid URL."}, nil, []int{http.StatusBadRequest})
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
