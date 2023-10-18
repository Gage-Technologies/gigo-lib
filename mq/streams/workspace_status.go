package streams

import (
	"github.com/nats-io/nats.go"
)

// this file contains the jetstream configuration for
// workspace status updates

const (
	StreamWorkspaceStatus string = "WorkspaceStatus"

	// we set the actual value for the subject to `>` so that we can dynamically create
	// new subjects using the workspace id that we are subscribing to - the `>` means
	// all subjects within the stream `WORKSPACE_STATUS`
	SubjectWorkspaceStatusUpdate        = "WORKSPACE_STATUS.>"
	SubjectWorkspaceStatusUpdateDynamic = "WORKSPACE_STATUS.%d"

	RetentionPolicyWorkspaceStatus = nats.InterestPolicy

	DuplicateFilterWindowWorkspaceStatus = 0
)

var StreamSubjectsWorkspaceStatus = []string{
	SubjectWorkspaceStatusUpdate,
}
