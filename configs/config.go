package configs

import (
	"os"

	"gopkg.in/yaml.v3"
)

type MasterDBConfig struct {
	Driver string `yaml:"driver"`
	DSN    string `yaml:"dsn"`
}

type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type GoogleOAuthConfig struct {
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	RedirectURL  string `yaml:"redirect_url"`
}

type OAuthConfig struct {
	Google GoogleOAuthConfig `yaml:"google"`
}

type AWSS3Config struct {
	Bucket string `yaml:"bucket"`
}

type AWSConfig struct {
	Region string      `yaml:"region"`
	S3     AWSS3Config `yaml:"s3"`
}

type RAGServerConfig struct {
	BaseURL           string `yaml:"base_url"`
	UploadDocumentURL string `yaml:"upload_document_url"`
	ChatURL           string `yaml:"chat_url"`
}
type Config struct {
	Port         int              `yaml:"port"`
	MasterDBs    []MasterDBConfig `yaml:"master_db"`
	Redis        RedisConfig      `yaml:"redis"`
	OAuth        OAuthConfig      `yaml:"oauth"`
	JwtSecret    string           `yaml:"jwt_secret"`
	WebClientURL string           `yaml:"web_client_url"`
	AllowOrigins []string         `yaml:"allow_origins"`
	AWS          AWSConfig        `yaml:"aws"`
	RAGServer    RAGServerConfig  `yaml:"rag_server"`
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
