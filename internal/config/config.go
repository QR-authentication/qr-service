package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Service  Service
	Postgres Postgres
	Metrics  Metrics
	Platform Platform
	Security Security
}

type Service struct {
	Port string `env:"QR_SERVICE_PORT"`
	Name string `env:"QR_SERVICE_NAME"`
}

type Postgres struct {
	User     string `env:"QR_SERVICE_POSTGRES_USER"`
	Password string `env:"QR_SERVICE_POSTGRES_PASSWORD"`
	Database string `env:"QR_SERVICE_POSTGRES_DB"`
	Host     string `env:"QR_SERVICE_POSTGRES_HOST"`
	Port     string `env:"QR_SERVICE_POSTGRES_PORT"`
}

type Metrics struct {
	Host string `env:"GRAFANA_HOST"`
	Port int    `env:"GRAFANA_PORT"`
}

type Platform struct {
	Env string `env:"ENV"`
}

type Security struct {
	SigningKey string `env:"SECURITY_QR_SIGNING_KEY"`
}

func MustLoad() *Config {
	cfg := &Config{}
	err := cleanenv.ReadEnv(cfg)

	if err != nil {
		log.Fatalf("failed to read env variables: %s", err)
	}

	return cfg
}
