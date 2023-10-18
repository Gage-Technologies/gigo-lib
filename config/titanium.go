package config

type TitaniumConfig struct {
	TitaniumHost       string   `yaml:"db_host"`
	TitaniumPort       string   `yaml:"db_port"`
	TitaniumPDHosts    []string `yaml:"pd_hosts"`
	TitaniumName       string   `yaml:"db_name"`
	TitaniumUser       string   `yaml:"db_user"`
	TitaniumPassword   string   `yaml:"db_password"`
	TitaniumBackupPath string   `yaml:"db_backup_path"`
}
