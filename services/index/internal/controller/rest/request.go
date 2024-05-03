package rest

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/jsonapi"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/model"
)

// NodeCreateRequest is a structure representing the request to create a new node.
type NodeCreateRequest struct {
	ID             string           `json:"node_id"`
	ProfileURL     string           `json:"profile_url"`
	ProfileHash    *string          `json:"profile_hash"`
	Status         string           `json:"status"`
	LastUpdated    *int64           `json:"last_updated"`
	FailureReasons *[]jsonapi.Error `json:"failure_reasons"`
}

// Validate is a method of NodeCreateRequest that validates the request fields.
func (n *NodeCreateRequest) Validate() []jsonapi.Error {
	// Check if ProfileURL is provided.
	if n.ProfileURL == "" {
		return jsonapi.NewError(
			[]string{"Missing Required Property"},
			[]string{"The `profile_url` property is required."},
			nil,
			[]int{http.StatusBadRequest},
		)
	}

	// Check if ProfileURL is a valid URL.
	u, err := url.Parse(n.ProfileURL)
	if err != nil || !isValidURL(u) {
		return jsonapi.NewError(
			[]string{"Invalid Profile URL"},
			[]string{"The `profile_url` is not a valid URL."},
			nil,
			[]int{http.StatusBadRequest},
		)
	}

	return nil
}

// isValidURL is a helper function that checks whether a URL is valid.
func isValidURL(u *url.URL) bool {
	return (u.Scheme == "http" || u.Scheme == "https") && u.Host != "" &&
		validHostname(u.Host)
}

// validHostname checks if the hostname of the URL has at least one dot or the
// URL is data-proxy-app.
func validHostname(host string) bool {
	uCount := strings.Count(host, ".")
	return uCount > 0 || strings.HasPrefix(host, "data-proxy-app")
}

// toDTO is a function that converts the model entity to a NodeCreateRequest DTO.
func toDTO(node *model.Node) *NodeCreateRequest {
	return &NodeCreateRequest{
		ID:             node.ID,
		ProfileURL:     node.ProfileURL,
		ProfileHash:    node.ProfileHash,
		Status:         node.Status,
		LastUpdated:    node.LastUpdated,
		FailureReasons: node.FailureReasons,
	}
}
