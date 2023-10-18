package config

type GiteaConfig struct {
	HostUrl  string `yaml:"host_url"`
	Username string `yaml:"username"`
	Password string `yaml:"password""`
}
