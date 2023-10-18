package streams

import (
	"github.com/nats-io/nats.go"
	"time"
)

const (
	StreamWsConnCache        = "WsConnCache"
	SubjectWsConnCacheForget = "WSCONNCACHE.Forget"

	RetentionPolicyWsConnCache = nats.InterestPolicy

	DuplicateFilterWindowWsConnCache = time.Second * 10
)

var StreamSubjectsWsConnCache = []string{
	SubjectWsConnCacheForget,
}
