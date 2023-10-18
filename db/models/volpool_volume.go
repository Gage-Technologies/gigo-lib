package models

import (
	"database/sql"
	"fmt"

	"github.com/kisielk/sqlstruct"
)

type VolumeState int

const (
	VolumeStateAvailable VolumeState = iota
	VolumeStateInUse
)

// Volume
//
//	Represents a valume that has been pre-provisioned by the volume pool
type VolpoolVolume struct {
	// ID Unique identifier of the volume
	ID int64 `json:"_id" sql:"_id"`

	// Size Size of the volume in gigabytes
	Size int `json:"size" sql:"size"`

	// State Current state of the volume
	State VolumeState `json:"state" sql:"state"`

	// PVCName Name of the PVC that owns the volume
	PVCName string `json:"pvc_name" sql:"pvc_name"`

	// StorageClass Name of the storage class that owns the volume
	StorageClass string `json:"storage_class" sql:"storage_class"`

	// WorkspaceID ID of the workspace that owns the volume
	WorkspaceID *int64 `json:"workspace_id" sql:"workspace_id"`
}

func CreateVolpoolVolume(_id int64, size int, state VolumeState, pvcName string, storageClass string, workspaceId *int64) *VolpoolVolume {
	return &VolpoolVolume{
		ID:           _id,
		Size:         size,
		State:        state,
		PVCName:      pvcName,
		StorageClass: storageClass,
		WorkspaceID:  workspaceId,
	}
}

func VolpoolVolumeFromSqlNative(rows *sql.Rows) (*VolpoolVolume, error) {
	volume := &VolpoolVolume{}
	err := sqlstruct.Scan(volume, rows)
	if err != nil {
		return nil, fmt.Errorf("failed to scan volume: %v", err)
	}
	return volume, nil
}

func (v *VolpoolVolume) ToSqlNative() ([]SQLInsertStatement, error) {
	return []SQLInsertStatement{
		{
			Statement: `insert into volpool_volume (_id, size, state, pvc_name, storage_class, workspace_id) values (?, ?, ?, ?, ?, ?)`,
			Values:    []interface{}{v.ID, v.Size, v.State, v.PVCName, v.StorageClass, v.WorkspaceID},
		},
	}, nil
}
