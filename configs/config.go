package configs

import (
	"os"

	"gopkg.in/yaml.v3"
)

type MasterDBConfig struct {
	Driver string `yaml:"driver"`
	DSN    string `yaml:"dsn"`
}

type Config struct {
	Port      int              `yaml:"port"`
	MasterDBs []MasterDBConfig `yaml:"master_db"`
	JwtSecret string           `yaml:"jwt_secret"`
}

var config Config

func Init() {
	file, err := os.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(file, &config)
	if err != nil {
		panic(err)
	}
}

func GetEnv() *Config {
	return &config
}
