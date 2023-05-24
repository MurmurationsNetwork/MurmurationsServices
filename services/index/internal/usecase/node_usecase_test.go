package usecase

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/dateutil"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/jsonapi"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/adapter/repository/db"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/entity"
)

var svc = NewNodeService(db.NewRepository())

func TestAddNodeWithoutProfileURL(t *testing.T) {
	node1 := entity.Node{}
	_, err := svc.AddNode(&node1)
	require.Equal(t, "Missing Required Property", err[0].Title)
	require.Equal(t, "The `profile_url` property is required.", err[0].Detail)
}

func TestAddNodeWithProfileURL(t *testing.T) {
	// FIXME: Not able to set version directly.
	version := int32(1)
	node2 := &entity.Node{
		ProfileURL: "https://ic3.dev/test.json",
		Version:    &version,
	}
	_, err := svc.AddNode(node2)
	require.Equal(t, []jsonapi.Error(nil), err)
	require.Equal(
		t,
		"4d62d0d132e2814379c22f2850d7a6c9ae97c16f021c25c975730c6b97b31f2c",
		node2.ID,
	)
	require.Equal(t, constant.NodeStatus.Received, node2.Status)
	require.Equal(t, dateutil.GetNowUnix(), node2.CreatedAt)
}

func TestGetNode(t *testing.T) {
	node, _ := svc.GetNode(
		"40ffe5e7db43150ebbae810b73b19aef318ac93f9191b9009249f62ae4624c69",
	)
	require.Equal(
		t,
		"40ffe5e7db43150ebbae810b73b19aef318ac93f9191b9009249f62ae4624c69",
		node.ID,
	)
}

func TestSetNodeValid(t *testing.T) {
	node := entity.Node{ProfileURL: "https://ic3.dev/test.json"}

	require.NoError(t, svc.SetNodeValid(&node))

	require.Equal(
		t,
		"4d62d0d132e2814379c22f2850d7a6c9ae97c16f021c25c975730c6b97b31f2c",
		node.ID,
	)
	require.Equal(t, constant.NodeStatus.Validated, node.Status)
	require.Equal(t, []jsonapi.Error{}, *node.FailureReasons)
}

func TestSetNodeInvalid(t *testing.T) {
	node := entity.Node{}

	require.NoError(t, svc.SetNodeInvalid(&node))

	require.Equal(t, constant.NodeStatus.ValidationFailed, node.Status)
	require.Equal(t, "", *node.ProfileHash)
	require.Equal(t, dateutil.GetZeroValueUnix(), *node.LastUpdated)
}
