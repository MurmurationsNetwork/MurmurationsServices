package es

import (
	"context"
	"fmt"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/elastic"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/nodecleaner/internal/model/query"
)

// NodeRepository defines the interface for operations that can be performed on
// nodes in an Elasticsearch repository.
type NodeRepository interface {
	// Remove deletes nodes with the specified status and creation time earlier
	// than the given timeBefore.
	Remove(ctx context.Context, status string, timeBefore int64) error
	// UpdateStatusByExpiration updates the status of nodes with expired status before the given time.
	UpdateStatusByExpiration(
		ctx context.Context,
		status string,
		timeBefore int64,
	) error
}

type nodeRepository struct {
}

// NewNodeRepository initializes and returns a new NodeRepository instance for
// interacting with Elasticsearch.
func NewNodeRepository() NodeRepository {
	return &nodeRepository{}
}

// Remove deletes nodes from Elasticsearch that have the specified status and
// were created or updated before the given time.
func (r *nodeRepository) Remove(
	_ context.Context,
	status string,
	timeBefore int64,
) error {
	q := query.EsQuery{Status: &status, TimeBefore: &timeBefore}

	err := elastic.Client.DeleteMany(constant.ESIndex.Node, q.Build())
	if err != nil {
		return fmt.Errorf(
			"error removing nodes with status %s and timeBefore %d from Elasticsearch: %v",
			status,
			timeBefore,
			err,
		)
	}

	return nil
}

// UpdateStatusByExpiration updates the status of nodes with expired status before the given time.
func (r *nodeRepository) UpdateStatusByExpiration(
	_ context.Context,
	status string,
	timeBefore int64,
) error {
	q := query.EsQuery{Status: &status, Expires: &timeBefore}

	update := map[string]interface{}{
		"status": "deleted",
	}

	err := elastic.Client.UpdateMany(constant.ESIndex.Node, q.Build(), update)
	if err != nil {
		return fmt.Errorf(
			"error updating nodes status with status %s and expiresBefore %d in Elasticsearch: %v",
			status,
			timeBefore,
			err,
		)
	}

	return nil
}
