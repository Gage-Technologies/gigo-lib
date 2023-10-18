package config

type EtcdConfig struct {
	Hosts    []string `yaml:"hosts"`
	Username string   `yaml:"username"`
	Password string   `yaml:"password"`
}
