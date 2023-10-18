package streams

import (
	"github.com/nats-io/nats.go"
)

// this file contains the jetstream configuration for
// broadcasting event

const (
	StreamBroadcastEvent string = "BroadcastEvent"

	SubjectBroadcastEvent               = "BROADCAST_EVENT.Event"
	SubjectBroadcastMessage             = "BROADCAST_EVENT.BroadcastMessage.>"
	SubjectBroadcastMessageDynamic      = "BROADCAST_EVENT.BroadcastMessage.%d"
	SubjectBroadcastNotification        = "BROADCAST_EVENT.BroadcastNotification.>"
	SubjectBroadcastNotificationDynamic = "BROADCAST_EVENT.BroadcastNotification.%d"

	RetentionPolicyBroadcastEvent = nats.InterestPolicy

	DuplicateFilterWindowBroadcastEvent = 0
)

var StreamSubjectsBroadcastEvent = []string{
	SubjectBroadcastEvent,
	SubjectBroadcastMessage,
	SubjectBroadcastNotification,
}
