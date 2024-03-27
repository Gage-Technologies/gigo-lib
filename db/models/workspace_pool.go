package models

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"time"

	"github.com/kisielk/sqlstruct"
)

type WorkspacePoolState int

const (
	WorkspacePoolStateAvailable WorkspacePoolState = iota
	WorkspacePoolStateInUse
	WorkspacePoolStateProvisioning
	WorkspacePoolStateDestroying
)

// WorkspacePool
//
//	Represents a workspace that has been pre-provisioned by the workspace pool
type WorkspacePool struct {
	// ID Unique identifier of the volume
	ID int64 `json:"_id" sql:"_id"`

	// Container Base container of the workspace
	Container string `json:"container" sql:"container"`

	// State Current state of the workspace pool availability
	State WorkspacePoolState `json:"state" sql:"state"`

	// Memory Available memory in GB of the workspace
	Memory int64 `json:"memory" sql:"memory"`

	// CPU Available CPU in cores of the workspace
	CPU int64 `json:"cpu" sql:"cpu"`

	// VolumeSize Size of the volume in gigabytes
	VolumeSize int `json:"size" sql:"size"`

	// StorageClass Name of the storage class that owns the volume
	Secret string `json:"secret" sql:"secret"`

	// AgentID Unique identifier of the agent that owns the workspace
	AgentID int64 `json:"agent_id" sql:"agent_id"`

	// WorkspaceID ID of the workspace that owns the volume
	WorkspaceTableID *int64 `json:"workspace_table_id" sql:"workspace_table_id"`

	// StorageClass Name of the storage class that owns the volume
	StorageClass string `json:"storage_class" sql:"storage_class"`

	// LastStateUpdate Timestamp when the last state update occurred
	LastStateUpdate *time.Time `json:"last_state_update" sql:"last_state_update"`

	// Expiration Timestamp when the workspace pool expires
	Expiration *time.Time `json:"expiration" sql:"expiration"`
}

func CreateWorkspacePool(_id int64, container string, state WorkspacePoolState, memory int64, cpu int64, volumeSize int, secret string, storageClass string, workspaceTableId *int64) *WorkspacePool {
	n := time.Now()
	e := n.Add(time.Hour * 8)
	return &WorkspacePool{
		ID:               _id,
		Container:        container,
		State:            state,
		Memory:           memory,
		CPU:              cpu,
		VolumeSize:       volumeSize,
		Secret:           secret,
		StorageClass:     storageClass,
		WorkspaceTableID: workspaceTableId,
		LastStateUpdate:  &n,
		Expiration:       &e,
	}
}

func WorkspacePoolFromSqlNative(rows *sql.Rows) (*WorkspacePool, error) {
	pool := &WorkspacePool{}
	err := sqlstruct.Scan(pool, rows)
	if err != nil {
		return nil, fmt.Errorf("failed to scan volume: %v", err)
	}
	return pool, nil
}

func (w *WorkspacePool) ToSqlNative() ([]SQLInsertStatement, error) {
	// create fake uuid to placehold if the secret is empty
	secret := w.Secret
	if len(secret) == 0 {
		secret = uuid.New().String()
	}
	return []SQLInsertStatement{
		{
			Statement: `insert into workspace_pool (_id, container, state, memory, cpu, volume_size, secret, agent_id, workspace_table_id, last_state_update, expiration) values (?, ?, ?, ?, ?, ?, uuid_to_bin(?), ?, ?, ?, ?)`,
			Values:    []interface{}{w.ID, w.Container, w.State, w.Memory, w.CPU, w.VolumeSize, secret, w.AgentID, w.WorkspaceTableID, w.LastStateUpdate, w.Expiration},
		},
	}, nil
}
