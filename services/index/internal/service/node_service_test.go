package service

import (
	"testing"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/dateutil"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/domain/node"
	"github.com/stretchr/testify/assert"
)

func TestAddNode(t *testing.T) {
	// Without profile_url
	node1 := node.Node{}
	_, err := NodeService.AddNode(&node1)
	assert.Equal(t, "The profile_url parameter is missing.", err.Message())

	// With profile_url
	// FIXME: Not able to set version directly.
	version := int32(1)
	node2 := &node.Node{ProfileURL: "https://ic3.dev/test.json", Version: &version}
	_, err = NodeService.AddNode(node2)
	assert.Equal(t, nil, err)
	assert.Equal(t, "4d62d0d132e2814379c22f2850d7a6c9ae97c16f021c25c975730c6b97b31f2c", node2.ID)
	assert.Equal(t, constant.NodeStatus.Received, node2.Status)
}

func TestGetNode(t *testing.T) {
	node, _ := NodeService.GetNode("40ffe5e7db43150ebbae810b73b19aef318ac93f9191b9009249f62ae4624c69")
	assert.Equal(t, "40ffe5e7db43150ebbae810b73b19aef318ac93f9191b9009249f62ae4624c69", node.ID)
}

func TestSetNodeValid(t *testing.T) {
	node := node.Node{ProfileURL: "https://ic3.dev/test.json"}

	NodeService.SetNodeValid(&node)

	assert.Equal(t, "4d62d0d132e2814379c22f2850d7a6c9ae97c16f021c25c975730c6b97b31f2c", node.ID)
	assert.Equal(t, constant.NodeStatus.Validated, node.Status)
	assert.Equal(t, []string{}, *node.FailureReasons)
}

func TestSetNodeInvalid(t *testing.T) {
	node := node.Node{}

	NodeService.SetNodeInvalid(&node)

	assert.Equal(t, constant.NodeStatus.ValidationFailed, node.Status)
	assert.Equal(t, "", *node.ProfileHash)
	assert.Equal(t, dateutil.GetZeroValueUnix(), *node.LastValidated)
}
