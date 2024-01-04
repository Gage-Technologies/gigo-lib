package models

import (
	"database/sql"
	"fmt"

	"github.com/kisielk/sqlstruct"
)

type WorkspacePoolState int

const (
	WorkspacePoolStateAvailable WorkspacePoolState = iota
	WorkspacePoolStateInUse
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

	// VolumeSize Size of the volume in gigabytes
	VolumeSize int `json:"size" sql:"size"`

	// StorageClass Name of the storage class that owns the volume
	Secret string `json:"secret" sql:"secret"`

	// WorkspaceID ID of the workspace that owns the volume
	WorkspaceTableID *int64 `json:"workspace_table_id" sql:"workspace_table_id"`

	// StorageClass Name of the storage class that owns the volume
	StorageClass string `json:"storage_class" sql:"storage_class"`
}

func CreateWorkspacePool(_id int64, container string, state WorkspacePoolState, memory int64, cpu int64, volumeSize int, secret string, storageClass string, workspaceTableId *int64) *WorkspacePool {
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
	return []SQLInsertStatement{
		{
			Statement: `insert into workspace_pool (_id, container, state, memory, cpu, volume_size, secret, workspace_table_id) values (?, ?, ?, ?, ?, ?, ?, ?)`,
			Values:    []interface{}{w.ID, w.Container, w.State, w.Memory, w.CPU, w.VolumeSize, w.Secret, w.WorkspaceTableID},
		},
	}, nil
}
