package tailnettest_test

import (
	"testing"

	"go.uber.org/goleak"

	"github.com/gage-technologies/gigo-lib/coder/tailnet/tailnettest"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestRunDERPAndSTUN(t *testing.T) {
	t.Parallel()
	_ = tailnettest.RunDERPAndSTUN(t)
}
