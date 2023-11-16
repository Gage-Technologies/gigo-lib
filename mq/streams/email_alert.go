package streams

import "github.com/nats-io/nats.go"

// this file contains the jetstream work ques for email alerts

const (
	StreamEmail string = "Email"

	SubjectEmailUserInactiveWeek    = "EMAIL.UserInactiveWeek"
	SubjectEmailUserInactiveMonth   = "EMAIL.UserInactiveMonth"
	SubjectEmailUserStreakEnd       = "EMAIL.UserStreakEnd"
	SubjectEmailUserReceivedMessage = "EMAIL.UserReceivedMessage"

	RetentionPolicyEmail = nats.WorkQueuePolicy

	DuplicateFilterWindowEmail = 0
)

var StreamSubjectsEmail = []string{
	SubjectEmailUserInactiveWeek,
	SubjectEmailUserInactiveMonth,
	SubjectEmailUserStreakEnd,
	SubjectEmailUserReceivedMessage,
}
