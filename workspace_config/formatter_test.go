package workspace_config

import (
	"fmt"
	"github.com/gage-technologies/gigo-lib/utils"
	"testing"
)

func TestFormatTerraformTemplate(t *testing.T) {
	// _, b, _, _ := runtime.Caller(0)
	// basepath := strings.Replace(filepath.Dir(b), "/src/gigo/workspace_config", "", -1)
	//
	// wsConfig, err := ParseWorkspaceConfig(basepath + "/test_data/gigoconfig.yaml")
	// if err != nil {
	// 	t.Fatalf("\nFormatTerraformTemplate failed\n    Error: %v", err)
	// }

	// /tmp/1611617066815586304/.gigo/workspace.yaml

	wsConfig, err := ParseWorkspaceConfig("/tmp/1611617066815586304/.gigo/workspace.yaml")
	if err != nil {
		t.Fatalf("\nFormatTerraformTemplate failed\n    Error: %v", err)
	}

	tfTemplate, err := FormatTerraformTemplate(wsConfig)
	if err != nil {
		t.Fatalf("\nFormatTerraformTemplate failed\n    Error: %v", err)
	}

	h, err := utils.HashData([]byte(tfTemplate))
	if err != nil {
		t.Fatalf("\nFormatTerraformTemplate failed\n    Error: %v", err)
	}

	fmt.Println(h)

	fmt.Println(tfTemplate)
}
