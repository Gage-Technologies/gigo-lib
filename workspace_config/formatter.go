package workspace_config

import (
	"embed"
	_ "embed"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"strings"
)

var (
	//go:embed resources
	gigoWorkspaceConfigResources embed.FS
)

func FormatTerraformTemplate(wsConfig *GigoWorkspaceConfig) (string, error) {
	// load bootstrap script from embedded filesystem
	bootstrapScriptBytes, err := gigoWorkspaceConfigResources.ReadFile("resources/gigo_ws_init.py")
	if err != nil {
		return "", fmt.Errorf("failed to load bootstrap script from embeded file system: %v", err)
	}

	// ensure that the script is valid
	if len(bootstrapScriptBytes) == 0 {
		return "", fmt.Errorf("empty bootstrap script loaded from embeded file system")
	}

	// format bootstrap script to a string
	bootstrapScript := string(bootstrapScriptBytes)

	// format WORKING_DIRECTORY
	bootstrapScript = strings.ReplaceAll(bootstrapScript, "<working_directory>", wsConfig.WorkingDirectory)

	// conditionally format CONTAINER_COMPOSE
	if len(wsConfig.Containers) > 0 {
		// marshall docker compose map to yaml buffer
		buf, err := yaml.Marshal(wsConfig.Containers)
		if err != nil {
			return "", fmt.Errorf("failed to marshall docker compose map: %v", err)
		}
		// replace <containers> placeholder with docker compose yaml
		bootstrapScript = strings.Replace(bootstrapScript, "<containers>", string(buf), 1)
	}

	// conditionally format SHELL_EXECUTIONS
	if len(wsConfig.Exec) > 0 {
		// marshall shell commands map to json buffer
		buf, err := json.Marshal(wsConfig.Exec)
		if err != nil {
			return "", fmt.Errorf("failed to marshall shell commands map: %v", err)
		}

		// convert buffer to string
		shellExecs := string(buf)

		// replace backwards slashes with double slashes for compatibility with python
		shellExecs = strings.ReplaceAll(shellExecs, "\\", "\\\\")

		// replace SHELL_EXECUTIONS placeholder with shell commands json
		bootstrapScript = strings.ReplaceAll(
			bootstrapScript,
			"SHELL_EXECUTIONS = \"\"\"[]\"\"\"",
			fmt.Sprintf("SHELL_EXECUTIONS = \"\"\"%s\"\"\"", shellExecs),
		)
	}

	// conditionally format vscode
	if wsConfig.VSCode.Enabled {
		// append code-tour extension install
		wsConfig.VSCode.Extensions = append(wsConfig.VSCode.Extensions, "vsls-contrib.codetour")

		// marshall vscode extensions to json buffer
		buf, err := json.Marshal(wsConfig.VSCode.Extensions)
		if err != nil {
			return "", fmt.Errorf("failed to marshall vscode extensions slice: %v", err)
		}
		// replace VSCODE_EXTENSIONS placeholder with vscode extensions json
		bootstrapScript = strings.ReplaceAll(
			bootstrapScript,
			"VSCODE_EXTENSIONS = []",
			fmt.Sprintf("VSCODE_EXTENSIONS = %s", string(buf)),
		)
	} else {
		// disable vscode
		bootstrapScript = strings.ReplaceAll(
			bootstrapScript,
			"USE_VSCODE = True",
			fmt.Sprintf("USE_VSCODE = False"),
		)
	}

	// load terraform template from embedded file system
	terraformTemplateBytes, err := gigoWorkspaceConfigResources.ReadFile("resources/gigo_ws_template.tf")
	if err != nil {
		return "", fmt.Errorf("failed to load terraform template from embeded file system: %v", err)
	}

	// ensure that the template is valid
	if len(terraformTemplateBytes) == 0 {
		return "", fmt.Errorf("empty terraform template loaded from embeded file system")
	}

	// format terraform template to a string
	terraformTemplate := string(terraformTemplateBytes)

	// conditionally format <environment>
	if len(wsConfig.Environment) > 0 {
		// create string to hold the environment variables
		env := ""

		// iterate environment variables formatting them to terraform format
		for k, v := range wsConfig.Environment {
			// conditionally new line and indentation
			if len(env) > 0 {
				env += "\n    "
			}

			// format variable
			env += fmt.Sprintf("%s = \"%s\"", k, v)
		}

		// format environment variables into terraform template
		terraformTemplate = strings.ReplaceAll(
			terraformTemplate,
			"<environment>",
			env,
		)
	} else {
		// remove <environment> place holder from terraform template
		terraformTemplate = strings.ReplaceAll(terraformTemplate, "<environment>", "")
	}

	// format working directory
	terraformTemplate = strings.ReplaceAll(terraformTemplate, "<working_directory>", wsConfig.WorkingDirectory)

	// format base container
	terraformTemplate = strings.ReplaceAll(terraformTemplate, "<base_container>", wsConfig.BaseContainer)

	// format resources
	terraformTemplate = strings.ReplaceAll(terraformTemplate, "<resources.cpu>", fmt.Sprintf("%d", wsConfig.Resources.CPU))
	terraformTemplate = strings.ReplaceAll(terraformTemplate, "<resources.mem>", fmt.Sprintf("%d", wsConfig.Resources.Mem))
	terraformTemplate = strings.ReplaceAll(terraformTemplate, "<resources.disk>", fmt.Sprintf("%d", wsConfig.Resources.Disk))

	// format indented bootstrap script into terraform template
	terraformTemplate = strings.ReplaceAll(terraformTemplate, "<bootstrap_script>", bootstrapScript)

	return terraformTemplate, nil
}
