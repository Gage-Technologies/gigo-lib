package streams

import (
	"github.com/nats-io/nats.go"
	"time"
)

const (
	StreamStreakXP = "Streak"

	SubjectStreakAddXP      = "STREAK.AddXp"
	SubjectStreakExpiration = "STREAK.ExpirationRemoval"
	SubjectDayRollover      = "STREAK.DayRollover"
	SubjectPremiumFreeze    = "STREAK.PremiumFreeze"

	RetentionPolicyStreak = nats.WorkQueuePolicy

	DuplicateFilterWindowStreak = time.Second * 10
)

var StreamSubjectsStreak = []string{
	SubjectStreakAddXP,
	SubjectStreakExpiration,
	SubjectDayRollover,
	SubjectPremiumFreeze,
}
