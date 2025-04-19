package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type DatabaseConfig struct {
	Host     string `env:"DB_HOST"`
	Port     int    `env:"DB_PORT"`
	User     string `env:"DB_USER"`
	Password string `env:"DB_PASSWORD"`
	Name     string `env:"DB_NAME"`
}

type MinioConfig struct {
	Endpoint        string `env:"MINIO_ENDPOINT"`
	PublicEndpoint  string `env:"MINIO_PUBLIC_ENDPOINT"`
	AccessKeyID     string `env:"MINIO_ROOT_USER"`
	SecretAccessKey string `env:"MINIO_ROOT_PASSWORD"`
	UseSSL          bool   `env:"MINIO_USE_SSL"`
	BucketName      string `env:"MINIO_BUCKET_NAME"`
	Port            int    `env:"MINIO_PORT"`
}
type JWTConfig struct {
	Secret string `env:"JWT_SECRET"`
}

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	ServerPort    int    `env:"SERVER_PORT"`

	Database DatabaseConfig
	Minio    MinioConfig
	JWT      JWTConfig
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("failed to load .env file: %w", err)
	}

	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse environment variables: %w", err)
	}

	return cfg, nil
}
