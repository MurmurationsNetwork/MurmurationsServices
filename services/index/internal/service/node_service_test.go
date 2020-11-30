package service

import (
	"testing"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/dateutil"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/internal/domain/node"
	"github.com/stretchr/testify/assert"
)

func TestSetNodeInvalid(t *testing.T) {
	node := node.Node{}

	NodeService.SetNodeInvalid(&node)

	assert.Equal(t, constant.NodeStatus.ValidationFailed, node.Status)
	assert.Equal(t, "", *node.ProfileHash)
	assert.Equal(t, dateutil.GetZeroValueUnix(), *node.LastValidated)
}
