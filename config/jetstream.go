package config

type JetstreamConfig struct {
	Host        string `yaml:"host"`
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
	MaxPubQueue int    `yaml:"max_pub_queue"`
}
