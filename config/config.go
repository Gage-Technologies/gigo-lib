package config

import "github.com/sirupsen/logrus"

type LoggerConfig struct {
	Name     string        `yaml:"name"`
	WorkerId string        `yaml:"worker_id"`
	Level    logrus.Level  `yaml:"level"`
	File     string        `yaml:"file"`
	ESConfig ElasticConfig `yaml:"es_config"`
}
