package models

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	ti "github.com/gage-technologies/gigo-lib/db"
	"github.com/google/uuid"
)

func TestCreateWorkspaceAgent(t *testing.T) {
	id := int64(12345)
	ownerId := int64(420)
	wsId := int64(69)
	version := "test"
	secret := uuid.New()

	agent := CreateWorkspaceAgent(id, wsId, version, ownerId, secret, "", "")

	if agent == nil {
		t.Error("\nCreateWorkspaceAgent failed\n    Error: nil agennt")
		return
	}

	if agent.ID != id {
		t.Error("\nCreateWorkspaceAgent failed\n    Error: incorrect workspace id returned")
		return
	}

	if agent.OwnerID != ownerId {
		t.Error("\nCreateWorkspaceAgent failed\n    Error: incorrect workspace owner id returned")
		return
	}

	if agent.WorkspaceID != wsId {
		t.Error("\nCreateWorkspaceAgent failed\n    Error: incorrect workspace id returned")
		return
	}

	if agent.Secret != secret {
		t.Error("\nCreateWorkspaceAgent failed\n    Error: incorrect workspace secret returned")
		return
	}

	t.Log("\nCreateWorkspaceAgent succeeded")
}

func TestWorkspaceAgentFromSQLNative(t *testing.T) {
	db, err := ti.CreateDatabase("gigo-dev-tidb", "4000", "mysql", "gigo-dev", "gigo-dev", "gigo_dev_test")
	if err != nil {
		t.Error("\nToSQlNative WorkspaceAgent Table\n    Error: ", err)
		return
	}

	defer db.DB.Exec("DROP TABLE workspace_agent")

	id := int64(12345)
	ownerId := int64(420)
	wsId := int64(69)
	version := "test"
	secret := uuid.New()

	agent := CreateWorkspaceAgent(id, wsId, version, ownerId, secret, "", "")

	agent.CreatedAt = time.Now().Add(-time.Hour)
	agent.UpdatedAt = time.Now().Add(-time.Minute)
	fc := time.Now().Add(-time.Minute * 58)
	agent.FirstConnect = &fc
	lc := time.Now().Add(-time.Minute * 9)
	agent.LastConnect = &lc
	ld := time.Now().Add(-time.Minute * 10)
	agent.LastConnectedNode = 69420
	agent.LastDisconnect = &ld
	agent.State = WorkspaceAgentStateRunning

	statements := agent.ToSQLNative()

	if statements[0].Statement != "insert ignore into workspace_agent(_id, created_at, updated_at, first_connect, last_connect, last_disconnect, last_connected_node, disconnect_count, state, workspace_id, version, owner_id, secret) values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, uuid_to_bin(?));" {
		t.Errorf("failed workspace agent to sql native, err: workspace insert statement was incorrect: %v", statements[0].Statement)
		return
	}

	tx, err := db.DB.Begin()
	if err != nil {
		t.Error("\nTo sql native failed\n    Error: incorrect values returned for user badges table")
		return
	}

	for _, s := range statements {
		_, err = tx.Exec(s.Statement, s.Values...)
		if err != nil {
			_ = tx.Rollback()
			t.Errorf("\nTo sql native failed\n    Error: %v", err)
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		t.Errorf(fmt.Sprintf("failed to commit transaction, err: %v", err))
		return
	}

	res, err := db.DB.Query("select * from workspace_agent where _id = ?", agent.ID)
	if err != nil {
		t.Errorf("\nWorkspaceAgentFromSQLNative failed\n    Error: failed to query for workspace: %v", err)
	}

	res.Next()

	a2, err := WorkspaceAgentFromSQLNative(res)
	if err != nil {
		t.Errorf("\nWorkspaceAgentFromSQLNative failed\n    Error: failed to load workspace: %v", err)
	}

	if a2 == nil {
		t.Error("\nWorkspaceAgentFromSQLNative failed\n    Error: load from sql is nil")
		return
	}

	agent.CreatedAt = time.Unix(0, 0)
	agent.UpdatedAt = time.Unix(0, 0)
	agent.FirstConnect = nil
	agent.LastConnect = nil
	agent.LastDisconnect = nil

	a2.CreatedAt = time.Unix(0, 0)
	a2.UpdatedAt = time.Unix(0, 0)
	a2.FirstConnect = nil
	a2.LastConnect = nil
	a2.LastDisconnect = nil

	if !reflect.DeepEqual(*agent, *a2) {
		t.Errorf("\nWorkspaceAgentFromSQLNative failed\n    Error: incorrect data returned\n%+v\n%+v", agent, a2)
		return
	}

	t.Logf("WorkspaceAgentFromSQLNative Succeeded")
}
