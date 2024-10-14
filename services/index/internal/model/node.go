package model

import (
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/jsonapi"
)

// Node represents a node stored in the index.
type Node struct {
	// ID is the unique identifier for the Node.
	ID string `bson:"_id,omitempty"`

	// ProfileURL is the URL of the node's profile.
	ProfileURL string `bson:"profile_url,omitempty"`

	// ProfileHash stores the hash of the node's profile.
	ProfileHash *string `bson:"profile_hash,omitempty"`

	// Status represents the current status of the node.
	Status string `bson:"status,omitempty"`

	// LastUpdated stores the Unix timestamp of the last update made to the node.
	LastUpdated *int64 `bson:"last_updated,omitempty"`

	// FailureReasons stores a list of errors
	// encountered during the node's operation.
	FailureReasons *[]jsonapi.Error `bson:"failure_reasons,omitempty"`

	// Version is the version vector of the node.
	// https://en.wikipedia.org/wiki/Version_vector
	Version *int32 `bson:"__v,omitempty"`

	// CreatedAt stores the Unix timestamp when the node was created.
	CreatedAt int64 `bson:"createdAt,omitempty"`

	// ProfileStr stores the node's profile in string format.
	// It won't be stored in MongoDB.
	ProfileStr string `bson:"-"`
	
	// Expires stores the Unix timestamp when the node expires.
	Expires *int64 `bson:"expires,omitempty"`
}

func (n *Node) SetStatusValidated() {
	n.Status = constant.NodeStatus.Validated
}

func (n *Node) SetStatusPostFailed() {
	n.Status = constant.NodeStatus.PostFailed
}

func (n *Node) SetStatusPosted() {
	n.Status = constant.NodeStatus.Posted
}

func (n *Node) ResetFailureReasons() {
	n.FailureReasons = &[]jsonapi.Error{}
}

// ClearLastUpdated sets the LastUpdated field to nil.
func (n *Node) ClearLastUpdated() {
	n.LastUpdated = nil
}
