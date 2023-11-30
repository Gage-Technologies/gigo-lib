package zitimesh

import (
	"github.com/gage-technologies/gigo-lib/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestManager(t *testing.T) {
	// Set up Ziti configuration
	zitiConfig := config.ZitiConfig{
		ManagementUser: "gigo-dev",
		ManagementPass: "gigo-dev",
		EdgeHost:       "gigo-dev-ziti-controller:1280",
		EdgeBasePath:   "/",
		EdgeSchemes:    []string{"https"},
	}

	// Create a new Manager instance
	manager, err := NewManager(zitiConfig)
	assert.NoError(t, err)

	// Test creating a workspace service policy
	defer manager.DeleteWorkspaceServicePolicy()
	err = manager.CreateWorkspaceServicePolicy()
	assert.NoError(t, err)

	// Test creating a server
	defer manager.DeleteServer(69)
	serverId, serverToken, err := manager.CreateServer(69)
	assert.NoError(t, err)
	assert.NotEmpty(t, serverToken)
	assert.NotEmpty(t, serverId)

	// Test creating an agent
	defer manager.DeleteAgent(420)
	agentId, agentToken, err := manager.CreateAgent(420)
	assert.NoError(t, err)
	assert.NotEmpty(t, agentToken)
	assert.NotEmpty(t, agentId)

	// Test creating a workspace service
	defer manager.DeleteWorkspaceService()
	_, err = manager.CreateWorkspaceService()
	assert.NoError(t, err)
}
