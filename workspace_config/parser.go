package workspace_config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

type GigoVSCodeConfig struct {
	Enabled    bool     `yaml:"enabled"`
	Extensions []string `yaml:"extensions"`
}

type GigoExecConfig struct {
	Name    string `yaml:"name" json:"name"`
	Init    bool   `yaml:"init" json:"init"`
	Command string `yaml:"command" json:"command"`
}

type GigoPortForwardConfig struct {
	Name string `yaml:"name"`
	Port uint16 `yaml:"port"`
}

type GigoWorkspaceConfig struct {
	Version   float64 `yaml:"version"`
	Resources struct {
		CPU  int `yaml:"cpu"`
		Mem  int `yaml:"mem"`
		Disk int `yaml:"disk"`
		GPU  struct {
			Count int    `yaml:"count"`
			Class string `yaml:"class"`
		} `yaml:"gpu"`
	} `yaml:"resources"`
	BaseContainer    string                  `yaml:"base_container"`
	WorkingDirectory string                  `yaml:"working_directory"`
	Environment      map[string]string       `yaml:"environment"`
	Containers       map[string]interface{}  `yaml:"containers"`
	VSCode           GigoVSCodeConfig        `yaml:"vscode"`
	PortForward      []GigoPortForwardConfig `yaml:"port_forward"`
	Exec             []GigoExecConfig        `yaml:"exec"`
}

type GigoWorkspaceConfigAgent struct {
	Version          float64                 `yaml:"version"`
	WorkingDirectory string                  `json:"working_directory"`
	Environment      map[string]string       `json:"environment"`
	Containers       map[string]interface{}  `json:"containers"`
	VSCode           GigoVSCodeConfig        `yaml:"vscode"`
	PortForward      []GigoPortForwardConfig `json:"port_forward"`
	Exec             []GigoExecConfig        `json:"exec"`
}

func ParseWorkspaceConfig(dir string) (*GigoWorkspaceConfig, error) {
	// read the provided config file and returns a byte array
	filename, err := filepath.Abs(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to form the absolute path: %v", err)
	}
	f, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	// create empty Config struct
	var c GigoWorkspaceConfig

	// decode document and assigns to cfg
	err = yaml.Unmarshal(f, &c)
	if err != nil {
		return nil, fmt.Errorf("failed decode config file: %v", err)
	}

	return &c, nil
}

func (c GigoWorkspaceConfig) ToAgent() GigoWorkspaceConfigAgent {
	return GigoWorkspaceConfigAgent{
		Version:          c.Version,
		WorkingDirectory: c.WorkingDirectory,
		Environment:      c.Environment,
		Containers:       c.Containers,
		VSCode:           c.VSCode,
		PortForward:      c.PortForward,
		Exec:             c.Exec,
	}
}
