package workspace_config

import (
	"encoding/json"
	"github.com/gage-technologies/gigo-lib/utils"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestParseWorkspaceConfig(t *testing.T) {
	_, b, _, _ := runtime.Caller(0)
	basepath := strings.Replace(filepath.Dir(b), "/src/gigo/workspace_config", "", -1)

	data, err := ParseWorkspaceConfig(basepath + "/test_data/gigoconfig.yaml")
	if err != nil {
		t.Fatalf("\nParseWorkspaceConfig failed\n    Error: %v", err)
	}

	if data == nil {
		t.Fatalf("\nParseWorkspaceConfig failed\n    Error: config returned nil")
	}

	buf, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		t.Fatalf("\nParseWorkspaceConfig failed\n    Error: %v", err)
	}

	h, err := utils.HashData(buf)
	if err != nil {
		t.Fatalf("\nParseWorkspaceConfig failed\n    Error: %v", err)
	}

	if h != "273abe63a07915b2ff18c4b1de7c34b669b0412a0eccc179f2ec71c1cd482e65" {
		t.Fatalf("\nParseWorkspaceConfig failed\n    Error: got %v want %v\n%s", h, "273abe63a07915b2ff18c4b1de7c34b669b0412a0eccc179f2ec71c1", string(buf))
	}

	t.Log("\nParseWorkspaceConfig succeeded")
}
