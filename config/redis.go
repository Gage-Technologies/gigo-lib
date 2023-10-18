package config

type RedisConfig struct {
	RedisNodes    []string `yaml:"redis_nodes"`
	RedisDatabase int      `yaml:"redis_db"`
	RedisType     string   `yaml:"redis_type"`
	RedisPassword string   `yaml:"redis_pass"`
}
