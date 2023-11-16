package config

type ZitiConfig struct {
	ManagementUser string   `yaml:"management_user"`
	ManagementPass string   `yaml:"management_pass"`
	EdgeHost       string   `yaml:"edge_host"`
	EdgeBasePath   string   `yaml:"edge_base_path"`
	EdgeSchemes    []string `yaml:"edge_schemes"`
}
