package models

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"math"
	"reflect"
	"testing"
	"time"

	ti "github.com/gage-technologies/gigo-lib/db"
)

func TestCreateWorkspace(t *testing.T) {
	id := int64(12345)
	csId := int64(785734)
	repoId := int64(420)
	createdAt := time.Unix(1677485673, 0)
	ownerId := int64(69)
	templateId := int64(42069)
	expiration := time.Unix(1677485673, 0)
	commit := hex.EncodeToString(sha1.New().Sum([]byte("test")))
	ports := []WorkspacePort{
		{
			Name: "t1",
			Port: 3000,
		},
		{
			Name: "t2",
			Port: 4000,
		},
	}

	testWorkspace, err := CreateWorkspace(id, repoId, csId, CodeSourcePost, createdAt, ownerId, templateId, expiration, commit, &DefaultWorkspaceSettings, nil, ports)
	if err != nil {
		t.Errorf("failed to create workspace, err: %v", err)
		return
	}

	if testWorkspace == nil {
		t.Error("\nCreate Workspace failed\n    Error: CreateWorkspace object returned none")
		return
	}

	if testWorkspace.ID != id {
		t.Error("\nCreate Workspace failed\n    Error: incorrect workspace id returned")
		return
	}

	if testWorkspace.RepoID != repoId {
		t.Error("\nCreate Workspace failed\n    Error: incorrect repo id returned")
		return
	}

	if testWorkspace.CodeSourceID != csId {
		t.Error("\nCreate Workspace failed\n    Error: incorrect code source id returned")
		return
	}

	if testWorkspace.CodeSourceType != CodeSourcePost {
		t.Error("\nCreate Workspace failed\n    Error: incorrect code source type returned")
		return
	}

	if testWorkspace.OwnerID != ownerId {
		t.Error("\nCreate Workspace failed\n    Error: incorrect workspace owner id returned")
		return
	}

	if testWorkspace.CreatedAt != createdAt {
		t.Error("\nCreate Workspace failed\n    Error: incorrect created at time returned")
		return
	}

	if testWorkspace.TemplateID != templateId {
		t.Error("\nCreate Workspace failed\n    Error: incorrect workspace template id returned")
		return
	}

	if testWorkspace.Expiration != expiration {
		t.Error("\nCreate Workspace failed\n    Error: incorrect expiration time returned")
		return
	}

	if testWorkspace.Commit != commit {
		t.Error("\nCreate Workspace failed\n    Error: incorrect commit returned")
		return
	}

	if !reflect.DeepEqual(*testWorkspace.WorkspaceSettings, DefaultWorkspaceSettings) {
		t.Error("\nCreate Workspace failed\n    Error: incorrect workspace settings returned")
		return
	}

	if !reflect.DeepEqual(testWorkspace.Ports, ports) {
		t.Error("\nCreate Workspace failed\n    Error: incorrect ports returned")
		return
	}

	t.Log("\nCreate Workspace succeeded")
}

func TestWorkspace_ToSQLNative(t *testing.T) {
	db, err := ti.CreateDatabase("gigo-dev-tidb", "4000", "mysql", "gigo-dev", "gigo-dev", "gigo_dev_test")
	if err != nil {
		t.Error("\nToSQlNative Workspace Table\n    Error: ", err)
		return
	}

	defer db.DB.Exec("delete from workspaces")

	id := int64(12345)
	csId := int64(8685)
	repoId := int64(420)
	createdAt := time.Unix(1677485673, 0)
	ownerId := int64(69)
	templateId := int64(42069)
	expiration := time.Unix(1677485673, 0)
	commit := hex.EncodeToString(sha1.New().Sum([]byte("test")))
	ports := []WorkspacePort{
		{
			Name: "t1",
			Port: 3000,
		},
		{
			Name: "t2",
			Port: 4000,
		},
	}

	workspace, err := CreateWorkspace(id, repoId, csId, CodeSourcePost, createdAt, ownerId, templateId, expiration, commit, &DefaultWorkspaceSettings, nil, ports)

	statements, err := workspace.ToSQLNative()
	if err != nil {
		t.Error("\nToSQLNative failed\n    Error: ", err)
		return
	}

	if len(statements) < 1 {
		t.Error("failed workspace to sql native, err: no statements returned")
		return
	}

	if statements[0].Statement != "insert ignore into workspaces(_id, code_source_id, code_source_type, repo_id, created_at, owner_id, template_id, expiration, commit, state, init_state, init_failure, last_state_update, workspace_settings, over_allocated, ports) values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);" {
		t.Errorf("failed workspace to sql native, err: workspace insert statement was incorrect: %v", statements[0].Statement)
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

	t.Logf("Workspace To SQL Native Succeeded")
}

func TestWorkspaceFromSQLNative(t *testing.T) {
	db, err := ti.CreateDatabase("gigo-dev-tidb", "4000", "mysql", "gigo-dev", "gigo-dev", "gigo_dev_test")
	if err != nil {
		t.Error("\nToSQlNative Workspace Table\n    Error: ", err)
		return
	}

	defer db.DB.Exec("delete from workspaces")

	id := int64(12345)
	csId := int64(47478568)
	repoId := int64(420)
	createdAt := time.Unix(1677485673, 0)
	ownerId := int64(69)
	templateId := int64(42069)
	expiration := time.Unix(1677485673, 0)
	commit := hex.EncodeToString(sha1.New().Sum([]byte("test")))
	ports := []WorkspacePort{
		{
			Name: "t1",
			Port: 3000,
		},
		{
			Name: "t2",
			Port: 4000,
		},
	}

	workspace, err := CreateWorkspace(id, repoId, csId, CodeSourcePost, createdAt, ownerId, templateId, expiration, commit, &DefaultWorkspaceSettings, nil, ports)

	secret := make([]byte, 64)
	_, err = rand.Read(secret)
	if err != nil {
		t.Error("\nWorkspaceFromSQLNative failed\n    Error: ", err)
		return
	}
	workspace.State = WorkspaceActive
	workspace.InitState = WorkspaceInitContainerComposeUp
	workspace.InitFailure = &WorkspaceInitFailure{
		Command: "test",
		Status:  255,
		Stdout:  "testing",
		Stderr:  "testing error",
	}

	statements, err := workspace.ToSQLNative()
	if err != nil {
		t.Error("\nTo sql native failed\n    Error: ", err)
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

	res, err := db.DB.Query("select _id, code_source_id, code_source_type, repo_id, created_at, owner_id, template_id, expiration, commit, state, init_state, init_failure, last_state_update, workspace_settings, ports from workspaces where _id = ?", workspace.ID)
	if err != nil {
		t.Errorf("\nWorkspaceFromSQLNative failed\n    Error: failed to query for workspace: %v", err)
	}

	res.Next()

	ws, err := WorkspaceFromSQLNative(res)
	if err != nil {
		t.Errorf("\nWorkspaceFromSQLNative failed\n    Error: failed to load workspace: %v", err)
	}

	if ws == nil {
		t.Error("\nWorkspaceFromSQLNative failed\n    Error: CreateWorkspace object returned none")
		return
	}

	if ws.ID != workspace.ID {
		t.Error("\nWorkspaceFromSQLNative failed\n    Error: incorrect workspace id returned")
		return
	}

	if ws.RepoID != workspace.RepoID {
		t.Error("\nWorkspaceFromSQLNative failed\n    Error: incorrect repo id returned")
		return
	}

	if ws.CodeSourceID != workspace.CodeSourceID {
		t.Error("\nWorkspaceFromSQLNative failed\n    Error: incorrect code source id returned")
		return
	}

	if ws.CodeSourceType != workspace.CodeSourceType {
		t.Error("\nWorkspaceFromSQLNative failed\n    Error: incorrect code source type returned")
		return
	}

	if ws.OwnerID != workspace.OwnerID {
		t.Error("\nWorkspaceFromSQLNative failed\n    Error: incorrect workspace owner id returned")
		return
	}

	if time.Since(ws.CreatedAt)-time.Since(workspace.CreatedAt) > time.Second*3 {
		fmt.Println(ws.CreatedAt)
		fmt.Println(workspace.CreatedAt.UTC())
		t.Error("\nWorkspaceFromSQLNative failed\n    Error: incorrect created at time returned")
		return
	}

	if ws.TemplateID != workspace.TemplateID {
		t.Error("\nWorkspaceFromSQLNative failed\n    Error: incorrect workspace template id returned")
		return
	}

	if time.Since(ws.Expiration)-time.Since(workspace.Expiration) > time.Second*3 {
		t.Error("\nWorkspaceFromSQLNative failed\n    Error: incorrect expiration time returned")
		return
	}

	if ws.Commit != workspace.Commit {
		fmt.Println(ws.Commit)
		fmt.Println(workspace.Commit)
		t.Error("\nWorkspaceFromSQLNative failed\n    Error: incorrect commit returned")
		return
	}

	if ws.State != workspace.State {
		t.Error("\nWorkspaceFromSQLNative failed\n    Error: incorrect state returned")
		return
	}

	if ws.InitState != workspace.InitState {
		t.Error("\nWorkspaceFromSQLNative failed\n    Error: incorrect init state returned")
		return
	}

	if !reflect.DeepEqual(ws.InitFailure, workspace.InitFailure) {
		t.Error("\nWorkspaceFromSQLNative failed\n    Error: incorrect init failure returned")
		return
	}

	if math.Abs(float64(ws.LastStateUpdate.Unix()-workspace.LastStateUpdate.Unix())) > 1 {
		fmt.Println(ws.LastStateUpdate)
		fmt.Println(workspace.LastStateUpdate)
		t.Error("\nWorkspaceFromSQLNative failed\n    Error: incorrect last state update returned")
		return
	}

	if !reflect.DeepEqual(*ws.WorkspaceSettings, *workspace.WorkspaceSettings) {
		t.Error("\nWorkspaceFromSQLNative failed\n    Error: incorrect workspace settings returned")
		return
	}

	if !reflect.DeepEqual(ws.Ports, workspace.Ports) {
		t.Error("\nWorkspaceFromSQLNative failed\n    Error: incorrect ports returned")
		return
	}

	t.Logf("WorkspaceFromSQLNative Succeeded")
}
