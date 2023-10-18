package ti

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateDatabase(t *testing.T) {
	_, err := CreateDatabase("gigo-dev-tidb", "4000", "mysql", "gigo-dev",
		"gigo-dev",
		"gigo_dev_test")
	assert.NoError(t, err)
}
