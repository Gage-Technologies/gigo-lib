package streams

import (
	"github.com/nats-io/nats.go"
	"time"
)

// this file contains the jetstream configuration for
// workspace based work queues

const (
	StreamWorkspace string = "Workspace"

	SubjectWorkspaceCreate  = "WORKSPACE.Create"
	SubjectWorkspaceStart   = "WORKSPACE.Start"
	SubjectWorkspaceStop    = "WORKSPACE.Stop"
	SubjectWorkspaceDestroy = "WORKSPACE.Destroy"
	SubjectWorkspaceDelete  = "WORKSPACE.Delete"

	RetentionPolicyWorkspace = nats.WorkQueuePolicy

	DuplicateFilterWindowWorkspace = time.Second * 10
)

var StreamSubjectsWorkspace = []string{
	SubjectWorkspaceCreate,
	SubjectWorkspaceStart,
	SubjectWorkspaceStop,
	SubjectWorkspaceDestroy,
	SubjectWorkspaceDelete,
}
