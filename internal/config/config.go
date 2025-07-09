package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port             string `envconfig:"PORT" required:"true"`
	NodeAPIURL       string `envconfig:"NODE_API_URL" required:"true"`
	NodeAuthURL      string `envconfig:"NODE_AUTH_URL" required:"true"`
	NodeAuthEmail    string `envconfig:"NODE_AUTH_EMAIL" required:"true"`
	NodeAuthPassword string `envconfig:"NODE_AUTH_PASSWORD" required:"true"`
	Environment      string `envconfig:"ENVIRONMENT" default:"development"`
	LogLevel         string `envconfig:"LOG_LEVEL" default:"info"`
	Host             string `envconfig:"HOST" default:"localhost"`
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Could not load .env file. Please ensure it exists in the project root. Error: %v", err)
	}

	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("Failed to load environment configuration: %v", err)
	}
	return &cfg
}
