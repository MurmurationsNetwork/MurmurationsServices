package service

import (
	"testing"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/dateutil"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/domain/node"
	"github.com/stretchr/testify/assert"
)

func TestSetNodeValid(t *testing.T) {
	node := node.Node{
		ProfileURL: "https://ic3.dev/test3d.json",
	}

	NodeService.SetNodeValid(&node)

	assert.Equal(t, "be27585485bec6808bbd768061af1fa903800dcdfd93493f7aca50d2118798d9", node.ID)
	// TODO: We combine two status update in one SetNodeValid method.
	assert.Equal(t, constant.NodeStatus.Posted, node.Status)
	assert.Equal(t, []string{}, *node.FailureReasons)
}

func TestSetNodeInvalid(t *testing.T) {
	node := node.Node{}

	NodeService.SetNodeInvalid(&node)

	assert.Equal(t, constant.NodeStatus.ValidationFailed, node.Status)
	assert.Equal(t, "", *node.ProfileHash)
	assert.Equal(t, dateutil.GetZeroValueUnix(), *node.LastValidated)
}
