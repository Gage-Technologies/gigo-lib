package streams

import (
	"github.com/nats-io/nats.go"
	"time"
)

const (
	StreamNemesis            = "Nemesis"
	SubjectNemesisStatChange = "NEMESIS.StatChange"

	RetentionPolicyNemesis = nats.WorkQueuePolicy

	DuplicateFilterWindowNemesis = time.Second * 10
)

var StreamSubjectsNemesis = []string{
	SubjectNemesisStatChange,
}
