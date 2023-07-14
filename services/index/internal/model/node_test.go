package model_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/jsonapi"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/model"
)

func TestNodeStatus(t *testing.T) {
	node := &model.Node{
		FailureReasons: &[]jsonapi.Error{
			{
				Status: 400,
				Title:  "Invalid URL",
				Detail: "The URL provided is invalid",
			},
		},
	}

	// Test setting the status to Validated.
	node.SetStatusValidated()
	require.Equal(t, constant.NodeStatus.Validated, node.Status)

	// Test setting the status to PostFailed.
	node.SetStatusPostFailed()
	require.Equal(t, constant.NodeStatus.PostFailed, node.Status)

	// Test setting the status to Posted.
	node.SetStatusPosted()
	require.Equal(t, constant.NodeStatus.Posted, node.Status)

	// Test resetting failure reasons.
	node.ResetFailureReasons()
	require.Empty(t, *node.FailureReasons)
}
