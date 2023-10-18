package streams

import (
	"github.com/nats-io/nats.go"
)

// this file contains the jetstream configuration for
// session based work queues

const (
	StreamMisc string = "Misc"

	SubjectMiscSessionCleanKeys = "MISC.SessionCleanKeys"
	SubjectMiscUserFreePremium  = "MISC.UserFreePremium"

	RetentionPolicyMisc = nats.WorkQueuePolicy

	DuplicateFilterWindowMisc = 0
)

var StreamSubjectsMisc = []string{
	SubjectMiscSessionCleanKeys,
	SubjectMiscUserFreePremium,
}
