package config

type ElasticConfig struct {
	ESNodes         []string `yaml:"elastic_nodes"`
	Username        string   `yaml:"username"`
	ESPass          string   `yaml:"elastic_pass"`
	Index           string   `yaml:"index"`
	BatchSize       int      `yaml:"batch_size"`
	BatchTimeMillis int      `yaml:"batch_time_millis"`
}
