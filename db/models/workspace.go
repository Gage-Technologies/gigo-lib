package models

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/kisielk/sqlstruct"
	"time"
)

type WorkspaceState int

const (
	WorkspaceStarting WorkspaceState = iota
	WorkspaceActive
	WorkspaceStopping
	WorkspaceSuspended
	WorkspaceRemoving
	WorkspaceFailed
	WorkspaceDeleted
)

func (w WorkspaceState) String() string {
	switch w {
	case WorkspaceStarting:
		return "Starting"
	case WorkspaceActive:
		return "Active"
	case WorkspaceStopping:
		return "Stopping"
	case WorkspaceSuspended:
		return "Suspended"
	case WorkspaceRemoving:
		return "Removing"
	case WorkspaceFailed:
		return "Failed"
	case WorkspaceDeleted:
		return "Deleted"
	}
	return "Unknown"
}

type WorkspaceInitFailure struct {
	Command string `json:"command"`
	Status  int    `json:"status"`
	Stdout  string `json:"stdout"`
	Stderr  string `json:"stderr"`
}

type WorkspaceInitState int

const (
	WorkspaceInitProvisioning WorkspaceInitState = iota
	WorkspaceInitRemoteInitialization
	WorkspaceInitWriteGitConfig
	WorkspaceInitWriteWorkspaceConfig
	WorkspaceInitGitClone
	WorkspaceInitGitCheckout
	WorkspaceInitCreateContainerDirectory
	WorkspaceInitWriteContainerCompose
	WorkspaceInitContainerComposeUp
	WorkspaceInitVSCodeInstall
	WorkspaceInitVSCodeExtensionInstall
	WorkspaceInitShellExecutions
	WorkspaceInitVSCodeLaunch
	WorkspaceInitCompleted
	WorkspaceInitReadExistingWorkspaceConfig
)

func (w WorkspaceInitState) String() string {
	switch w {
	case WorkspaceInitProvisioning:
		return "Provisioning"
	case WorkspaceInitRemoteInitialization:
		return "RemoteInitialization"
	case WorkspaceInitWriteGitConfig:
		return "WriteGitConfig"
	case WorkspaceInitWriteWorkspaceConfig:
		return "WriteWorkspaceConfig"
	case WorkspaceInitGitClone:
		return "GitClone"
	case WorkspaceInitGitCheckout:
		return "GitCheckout"
	case WorkspaceInitCreateContainerDirectory:
		return "CreateContainerDirectory"
	case WorkspaceInitWriteContainerCompose:
		return "WriteContainerCompose"
	case WorkspaceInitContainerComposeUp:
		return "ContainerComposeUp"
	case WorkspaceInitShellExecutions:
		return "ShellExecutions"
	case WorkspaceInitVSCodeInstall:
		return "VSCodeInstall"
	case WorkspaceInitVSCodeExtensionInstall:
		return "VSCodeExtensionInstall"
	case WorkspaceInitVSCodeLaunch:
		return "VSCodeLaunch"
	case WorkspaceInitCompleted:
		return "Completed"
	case WorkspaceInitReadExistingWorkspaceConfig:
		return "ReadExistingWorkspaceConfig"
	default:
		return "Unknown"
	}
}

type OverAllocated struct {
	CPU  int `json:"cpu"`
	RAM  int `json:"ram"`
	DISK int `json:"disk"`
}

type WorkspacePort struct {
	Name       string `json:"name"`
	Port       uint16 `json:"port"`
	Configured bool   `json:"configured"`
	Active     bool   `json:"active"`
}

type WorkspacePortFrontend struct {
	Name       string `json:"name"`
	Port       uint16 `json:"port"`
	Url        string `json:"url"`
	Configured bool   `json:"configured"`
	Active     bool   `json:"active"`
}

func (p *WorkspacePort) ToFrontend(userId int64, workspaceId int64, hostname string, https bool) *WorkspacePortFrontend {
	scheme := "http"
	if https {
		scheme = "https"
	}
	return &WorkspacePortFrontend{
		Name: p.Name,
		Port: p.Port,
		Url: fmt.Sprintf("%s://%d-%d-%d.%s",
			scheme,
			userId,
			workspaceId,
			p.Port,
			hostname,
		),
	}
}

type Workspace struct {
	ID                int64                 `json:"_id" sql:"_id"`
	CodeSourceID      int64                 `json:"code_source_id" sql:"code_source_id"`
	CodeSourceType    CodeSource            `json:"code_source_type" sql:"code_source_type"`
	RepoID            int64                 `json:"repo_id" sql:"repo_id"`
	CreatedAt         time.Time             `json:"created_at" sql:"created_at"`
	OwnerID           int64                 `json:"owner_id" sql:"owner_id"`
	TemplateID        int64                 `json:"template_id" sql:"template_id"`
	Expiration        time.Time             `json:"expiration" sql:"expiration"`
	Commit            string                `json:"commit" sql:"commit"`
	State             WorkspaceState        `json:"state" sql:"state"`
	InitState         WorkspaceInitState    `json:"init_state" sql:"init_state"`
	InitFailure       *WorkspaceInitFailure `json:"init_failure" sql:"init_failure"`
	LastStateUpdate   time.Time             `json:"last_state_update" sql:"last_state_update"`
	WorkspaceSettings *WorkspaceSettings    `json:"workspace_settings" sql:"workspace_settings"`
	OverAllocated     *OverAllocated        `json:"over_allocated" sql:"over_allocated"`
	Ports             []WorkspacePort       `json:"ports" sql:"ports"`
	IsEphemeral       bool                  `json:"is_ephemeral" sql:"is_ephemeral"`
}

type WorkspaceSQL struct {
	ID                int64              `json:"_id" sql:"_id"`
	CodeSourceID      int64              `json:"code_source_id" sql:"code_source_id"`
	CodeSourceType    CodeSource         `json:"code_source_type" sql:"code_source_type"`
	RepoID            int64              `json:"repo_id" sql:"repo_id"`
	CreatedAt         time.Time          `json:"created_at" sql:"created_at"`
	OwnerID           int64              `json:"owner_id" sql:"owner_id"`
	TemplateID        int64              `json:"template_id" sql:"template_id"`
	Expiration        time.Time          `json:"expiration" sql:"expiration"`
	Commit            string             `json:"commit" sql:"commit"`
	State             WorkspaceState     `json:"state" sql:"state"`
	InitState         WorkspaceInitState `json:"init_state" sql:"init_state"`
	InitFailure       []byte             `json:"init_failure" sql:"init_failure"`
	LastStateUpdate   time.Time          `json:"last_state_update" sql:"last_state_update"`
	WorkspaceSettings []byte             `json:"workspace_settings" sql:"workspace_settings"`
	OverAllocated     []byte             `json:"over_allocated" sql:"over_allocated"`
	Ports             []byte             `json:"ports" sql:"ports"`
	IsEphemeral       bool               `json:"is_ephemeral" sql:"is_ephemeral"`
}

type WorkspaceFrontend struct {
	ID                   string                  `json:"_id"`
	CodeSourceID         string                  `json:"code_source_id"`
	CodeSourceType       CodeSource              `json:"code_source_type"`
	CodeSourceTypeString string                  `json:"code_source_type_string"`
	RepoID               string                  `json:"repo_id"`
	CreatedAt            string                  `json:"created_at"`
	OwnerID              string                  `json:"owner_id"`
	Expiration           string                  `json:"expiration"`
	Commit               string                  `json:"commit"`
	State                WorkspaceState          `json:"state"`
	StateString          string                  `json:"state_string"`
	InitState            WorkspaceInitState      `json:"init_state"`
	InitStateString      string                  `json:"init_state_string"`
	InitFailure          *WorkspaceInitFailure   `json:"init_failure"`
	WorkspaceSettings    *WorkspaceSettings      `json:"workspace_settings"`
	OverAllocated        *OverAllocated          `json:"over_allocated" sql:"over_allocated"`
	Ports                []WorkspacePortFrontend `json:"ports"`
	IsEphemeral          bool                    `json:"is_ephemeral" sql:"is_ephemeral"`
}

func CreateWorkspace(id int64, repoId int64, codeSourceId int64, codeSourceType CodeSource, createdAt time.Time,
	ownerId int64, templateId int64, expiration time.Time, commit string, settings *WorkspaceSettings, overAllocated *OverAllocated,
	ports []WorkspacePort) (*Workspace, error) {
	// create secret for new workspace
	secret := make([]byte, 64)
	_, err := rand.Read(secret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate random secret for workspace: %v", err)
	}

	return &Workspace{
		ID:                id,
		CodeSourceID:      codeSourceId,
		CodeSourceType:    codeSourceType,
		RepoID:            repoId,
		CreatedAt:         createdAt,
		OwnerID:           ownerId,
		TemplateID:        templateId,
		Expiration:        expiration,
		Commit:            commit,
		LastStateUpdate:   time.Now(),
		WorkspaceSettings: settings,
		OverAllocated:     overAllocated,
		Ports:             ports,
	}, nil
}

func WorkspaceFromSQLNative(rows *sql.Rows) (*Workspace, error) {
	// create new workspace object to load into
	workspaceSQL := new(WorkspaceSQL)

	// scan row into templates object
	err := sqlstruct.Scan(workspaceSQL, rows)
	if err != nil {
		return nil, err
	}

	// create empty variable to hold workspace initialization failure data
	var workspaceInitFailure WorkspaceInitFailure

	// conditionally unmarshall json for workspace initialization failure
	if workspaceSQL.InitFailure != nil {
		err = json.Unmarshal(workspaceSQL.InitFailure, &workspaceInitFailure)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshall initialization failure: %v", err)
		}
	}

	// create empty variable to hold over allocated workspace resources
	var overAllocated OverAllocated

	if workspaceSQL.OverAllocated != nil {
		err = json.Unmarshal(workspaceSQL.OverAllocated, &overAllocated)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshall over allocated workspace resources: %v", err)
		}
	}

	// create empty variable to hold workspace settings data
	var workspaceSettings WorkspaceSettings

	// conditionally unmarshall json for workspace settings
	if workspaceSQL.WorkspaceSettings != nil {
		err = json.Unmarshal(workspaceSQL.WorkspaceSettings, &workspaceSettings)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshall workspace settings: %v", err)
		}
	}

	// create empty variable to hold workspace ports data
	var workspacePorts []WorkspacePort

	// conditionally unmarshall json for workspace settings
	if workspaceSQL.Ports != nil {
		err = json.Unmarshal(workspaceSQL.Ports, &workspacePorts)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshall ports: %v", err)
		}
	}

	workspace := &Workspace{
		ID:                workspaceSQL.ID,
		CodeSourceID:      workspaceSQL.CodeSourceID,
		CodeSourceType:    workspaceSQL.CodeSourceType,
		RepoID:            workspaceSQL.RepoID,
		CreatedAt:         workspaceSQL.CreatedAt,
		OwnerID:           workspaceSQL.OwnerID,
		TemplateID:        workspaceSQL.TemplateID,
		Expiration:        workspaceSQL.Expiration,
		Commit:            workspaceSQL.Commit,
		State:             workspaceSQL.State,
		InitState:         workspaceSQL.InitState,
		InitFailure:       &workspaceInitFailure,
		LastStateUpdate:   workspaceSQL.LastStateUpdate,
		WorkspaceSettings: &workspaceSettings,
		OverAllocated:     &overAllocated,
		Ports:             workspacePorts,
		IsEphemeral:       workspaceSQL.IsEphemeral,
	}

	return workspace, nil
}

func (w *Workspace) ToFrontend(hostname string, https bool) *WorkspaceFrontend {
	// format port to frontend
	frontendPorts := make([]WorkspacePortFrontend, len(w.Ports))
	for i, port := range w.Ports {
		frontendPorts[i] = *port.ToFrontend(w.OwnerID, w.ID, hostname, https)
	}

	// create new workspace to frontend
	workspaceFront := &WorkspaceFrontend{
		ID:                   fmt.Sprintf("%d", w.ID),
		CodeSourceID:         fmt.Sprintf("%d", w.CodeSourceID),
		CodeSourceType:       w.CodeSourceType,
		CodeSourceTypeString: w.CodeSourceType.String(),
		RepoID:               fmt.Sprintf("%d", w.RepoID),
		CreatedAt:            w.CreatedAt.String(),
		OwnerID:              fmt.Sprintf("%d", w.OwnerID),
		Expiration:           w.Expiration.String(),
		Commit:               w.Commit,
		State:                w.State,
		StateString:          w.State.String(),
		InitState:            w.InitState,
		InitStateString:      w.InitState.String(),
		InitFailure:          w.InitFailure,
		WorkspaceSettings:    w.WorkspaceSettings,
		OverAllocated:        w.OverAllocated,
		Ports:                frontendPorts,
		IsEphemeral:          w.IsEphemeral,
	}

	return workspaceFront
}

func (w *Workspace) ToSQLNative() ([]*SQLInsertStatement, error) {
	// create byte slice variable to hold serialized initialization failure json
	var initFailureBytes []byte

	// conditionally marshall initialization failure json
	if w.InitFailure != nil {
		buf, err := json.Marshal(w.InitFailure)
		if err != nil {
			return nil, fmt.Errorf("failed to marshall initialization failure: %v", err)
		}
		initFailureBytes = buf
	}

	// create byte slice variable to hold serialized workspace settings json
	var wsOverAllocatedBytes []byte

	// conditionally marshall workspace settings json
	if w.OverAllocated != nil {
		buf, err := json.Marshal(w.OverAllocated)
		if err != nil {
			return nil, fmt.Errorf("failed to marshall over allocated workspace resources: %v", err)
		}
		wsOverAllocatedBytes = buf
	}

	// create byte slice variable to hold serialized workspace settings json
	var wsSettingsBytes []byte

	// conditionally marshall workspace settings json
	if w.WorkspaceSettings != nil {
		buf, err := json.Marshal(w.WorkspaceSettings)
		if err != nil {
			return nil, fmt.Errorf("failed to marshall workspace settings: %v", err)
		}
		wsSettingsBytes = buf
	}

	// create byte slice variable to hold serialized ports json
	var wsPorts []byte

	// conditionally marshall ports json
	if w.Ports != nil {
		buf, err := json.Marshal(w.Ports)
		if err != nil {
			return nil, fmt.Errorf("failed to marshall ports: %v", err)
		}
		wsPorts = buf
	}

	return []*SQLInsertStatement{
		{
			Statement: "insert ignore into workspaces(_id, code_source_id, code_source_type, repo_id, created_at, owner_id, template_id, expiration, commit, state, init_state, init_failure, last_state_update, workspace_settings, over_allocated, ports, is_ephemeral) " +
				"values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);",
			Values: []interface{}{w.ID, w.CodeSourceID, w.CodeSourceType, w.RepoID, w.CreatedAt, w.OwnerID, w.TemplateID, w.Expiration, w.Commit, w.State, w.InitState, initFailureBytes, w.LastStateUpdate, wsSettingsBytes, wsOverAllocatedBytes, wsPorts, w.IsEphemeral},
		},
	}, nil
}
