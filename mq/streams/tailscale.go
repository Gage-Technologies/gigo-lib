package streams

import (
	"github.com/nats-io/nats.go"
)

// this file contains the jetstream configuration for
// cluster node communication regarding the tailscale network

const (
	StreamTailscale string = "Tailscale"

	SubjectTailscaleAgent      = "TAILSCALE.Agent"
	SubjectTailscaleConnection = "TAILSCALE.Connection"

	RetentionPolicyTailscale = nats.InterestPolicy

	DuplicateFilterWindowTailscale = 0
)

var StreamSubjectsTailscale = []string{
	SubjectTailscaleAgent,
	SubjectTailscaleConnection,
}
