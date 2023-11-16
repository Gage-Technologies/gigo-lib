package models

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/kisielk/sqlstruct"
	"time"
)

type WorkspaceAgentState int

const (
	WorkspaceAgentStateUnknown  WorkspaceAgentState = 0
	WorkspaceAgentStateStarting WorkspaceAgentState = 1
	WorkspaceAgentStateRunning  WorkspaceAgentState = 2
	WorkspaceAgentStateStopping WorkspaceAgentState = 3
	WorkspaceAgentStateStopped  WorkspaceAgentState = 4
	WorkspaceAgentStateFailed   WorkspaceAgentState = 5
	WorkspaceAgentStateTimeout  WorkspaceAgentState = 6
)

func (s WorkspaceAgentState) String() string {
	switch s {
	case WorkspaceAgentStateStarting:
		return "Starting"
	case WorkspaceAgentStateRunning:
		return "Running"
	case WorkspaceAgentStateStopping:
		return "Stopping"
	case WorkspaceAgentStateStopped:
		return "Stopped"
	case WorkspaceAgentStateFailed:
		return "Failed"
	case WorkspaceAgentStateTimeout:
		return "Timeout"
	case WorkspaceAgentStateUnknown:
		return "Unknown"
	}
	return "Invalid"
}

type WorkspaceAgent struct {
	ID                int64               `json:"_id" sql:"_id"`
	CreatedAt         time.Time           `json:"created_at" sql:"created_at"`
	UpdatedAt         time.Time           `json:"updated_at" sql:"updated_at"`
	FirstConnect      *time.Time          `json:"first_connect" sql:"firs_connect"`
	LastConnect       *time.Time          `json:"last_connect" sql:"last_connect"`
	LastDisconnect    *time.Time          `json:"last_disconnect" sql:"last_disconnect"`
	LastConnectedNode int64               `json:"last_connected_node" sql:"last_connected_node"`
	DisconnectCount   int                 `json:"disconnect_count" sql:"disconnect_count"`
	State             WorkspaceAgentState `json:"state" sql:"state"`
	WorkspaceID       int64               `json:"workspace_id" sql:"workspace_id"`
	Version           string              `json:"version" sql:"version"`
	OwnerID           int64               `json:"owner_id" sql:"owner_id"`
	Secret            uuid.UUID           `json:"secret" sql:"secret"`
	ZitiID            string              `json:"ziti_id" sql:"ziti_id"`
	ZitiToken         string              `json:"ziti_token" sql:"ziti_token"`
}

func CreateWorkspaceAgent(id int64, workspace int64, version string, ownerID int64, secret uuid.UUID, zitiID string, zitiToken string) *WorkspaceAgent {
	return &WorkspaceAgent{
		ID:          id,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		State:       WorkspaceAgentStateUnknown,
		WorkspaceID: workspace,
		Version:     version,
		OwnerID:     ownerID,
		Secret:      secret,
		ZitiID:      zitiID,
		ZitiToken:   zitiToken,
	}
}

func WorkspaceAgentFromSQLNative(rows *sql.Rows) (*WorkspaceAgent, error) {
	var agent WorkspaceAgent
	err := sqlstruct.Scan(&agent, rows)
	if err != nil {
		return nil, fmt.Errorf("failed to scan workspace agent from rows: %v", err)
	}
	return &agent, nil
}

func (a *WorkspaceAgent) ToSQLNative() []*SQLInsertStatement {
	// create slice to hold insertion statements for this workspace config and initialize the slice with the main insertion statement
	return []*SQLInsertStatement{
		{
			Statement: "insert ignore into workspace_agent(_id, created_at, updated_at, first_connect, last_connect, last_disconnect, last_connected_node, disconnect_count, state, workspace_id, version, owner_id, secret, ziti_id, ziti_token) values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, uuid_to_bin(?), ?, ?);",
			Values: []interface{}{
				a.ID, a.CreatedAt, a.UpdatedAt, a.FirstConnect, a.LastConnect, a.LastDisconnect, a.LastConnectedNode, a.DisconnectCount, a.State, a.WorkspaceID, a.Version, a.OwnerID, a.Secret, a.ZitiID, a.ZitiToken,
			},
		},
	}
}
