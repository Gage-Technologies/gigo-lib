package config

type StorageEngine string

const (
	StorageEngineS3 StorageEngine = "s3"
	StorageEngineFS StorageEngine = "fs"
)

type StorageS3Config struct {
	Bucket    string `yaml:"bucket"`
	Region    string `yaml:"region"`
	AccessKey string `yaml:"access_key"`
	SecretKey string `yaml:"secret_key"`
	Endpoint  string `yaml:"endpoint"`
	UseSSL    bool   `yaml:"use_ssl"`
}

type StorageFSConfig struct {
	Root string `yaml:"root"`
}

type StorageConfig struct {
	Engine StorageEngine   `yaml:"engine"`
	S3     StorageS3Config `yaml:"s3"`
	FS     StorageFSConfig `yaml:"fs"`
}
