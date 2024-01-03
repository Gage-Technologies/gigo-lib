package models

import (
	"database/sql"
	"fmt"

	"github.com/kisielk/sqlstruct"
)

type WorkspacePoolState int

const (
	WorkspacePoolStateAvailable VolumeState = iota
	WorkspacePoolStateStateInUse
)

// Volume
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

	// Storage Available storage in GB of the workspace
	Storage int64 `json:"storage" sql:"storage"`

	// StorageClass Name of the storage class that owns the volume
	Secret string `json:"secret" sql:"secret"`

	// WorkspaceID ID of the workspace that owns the volume
	WorkspaceTableID *int64 `json:"workspace_table_id" sql:"workspace_table_id"`
}

func CreateWorkspacePool(_id int64, container string, state WorkspacePoolState, memory int64, cpu int64, storage int64, secret string, workspaceTableId *int64) *WorkspacePool {
	return &WorkspacePool{
		ID:               _id,
		Container:        container,
		State:            state,
		Memory:           memory,
		CPU:              cpu,
		Storage:          storage,
		Secret:           secret,
		WorkspaceTableID: workspaceTableId,
	}
}

func WorkspacePoolFromSqlNative(rows *sql.Rows) (*WorkspacePool, error) {
	volume := &WorkspacePool{}
	err := sqlstruct.Scan(volume, rows)
	if err != nil {
		return nil, fmt.Errorf("failed to scan volume: %v", err)
	}
	return volume, nil
}

func (w *WorkspacePool) ToSqlNative() ([]SQLInsertStatement, error) {
	return []SQLInsertStatement{
		{
			Statement: `insert into workspace_pool (_id, container, state, memory, cpu, storage, secret, workspace_table_id) values (?, ?, ?, ?, ?, ?, ?, ?)`,
			Values:    []interface{}{w.ID, w.Container, w.State, w.Memory, w.CPU, w.Storage, w.Secret, w.WorkspaceTableID},
		},
	}, nil
}
