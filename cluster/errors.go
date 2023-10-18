package cluster

import (
	"errors"
)

var (
	ErrClusterTickDisagreement = errors.New("cluster tick disagreement")
	ErrNoLease                 = errors.New("no lease")
)
