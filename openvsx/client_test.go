package openvsx

import (
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
)

func TestClient_GetMetadata(t *testing.T) {
	client := NewClient("", nil)

	meta, err := client.GetMetadata("ms-python.python", "")
	assert.NoError(t, err)

	assert.Equal(t, "ms-python", meta.Namespace)
	assert.Equal(t, "python", meta.Name)
	dl, ok := meta.Files["download"]
	assert.True(t, ok)
	assert.True(t, dl != "")
}

func TestClient_DownloadExtension(t *testing.T) {
	client := NewClient("", nil)

	buf, cl, err := client.DownloadExtension("ms-python.python", "")
	assert.NoError(t, err)

	f, err := os.CreateTemp("", "gigo-openvsx-client-test-*")
	assert.NoError(t, err)
	defer os.Remove(f.Name())

	_, err = io.Copy(f, buf)
	assert.NoError(t, err)

	_ = f.Close()

	stat, err := os.Stat(f.Name())
	assert.NoError(t, err)
	assert.Greater(t, cl, int64(0))
	assert.Equal(t, stat.Size(), cl)
}

func TestClient_DownloadExtension_Version(t *testing.T) {
	client := NewClient("", nil)

	buf, cl, err := client.DownloadExtension("ms-python.python", "2023.10.1")
	assert.NoError(t, err)

	f, err := os.CreateTemp("", "gigo-openvsx-client-test-*")
	assert.NoError(t, err)
	defer os.Remove(f.Name())

	_, err = io.Copy(f, buf)
	assert.NoError(t, err)

	_ = f.Close()

	stat, err := os.Stat(f.Name())
	assert.NoError(t, err)
	assert.Greater(t, cl, int64(0))
	assert.Equal(t, stat.Size(), cl)
}
