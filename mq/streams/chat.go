package streams

import (
	"github.com/nats-io/nats.go"
)

// this file contains the jetstream configuration for
// all chat based streaming

const (
	StreamChat string = "Chat"

	// we set the actual value for the subject to `>` so that we can dynamically create
	// new subjects using the chat id that we are subscribing to - the `>` means
	// all subjects within the stream `CHAT`
	SubjectChatMessages        = "CHAT.MESSAGES.>"
	SubjectChatMessagesDynamic = "CHAT.MESSAGES.%d"

	SubjectChatNewChat        = "CHAT.NEW_CHAT.>"
	SubjectChatNewChatDynamic = "CHAT.NEW_CHAT.%d"

	SubjectChatKick        = "CHAT.KICK.>"
	SubjectChatKickDynamic = "CHAT.KICK.%d"

	SubjectChatUpdated        = "CHAT.UPDATED.>"
	SubjectChatUpdatedDynamic = "CHAT.UPDATED.%d"

	RetentionPolicyChat = nats.InterestPolicy

	DuplicateFilterWindowChat = 0
)

var StreamSubjectsChat = []string{
	SubjectChatMessages,
	SubjectChatNewChat,
	SubjectChatKick,
	SubjectChatUpdated,
}
