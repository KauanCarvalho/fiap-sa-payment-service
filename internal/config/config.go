package config

import (
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type Config struct {
	AppName  string
	AppEnv   string
	MongoURI string
	Port     string
}

func Load() *Config {
	environment := getEnv("APP_ENV", "development")

	err := godotenv.Load()
	if err != nil && environment == "development" {
		zap.L().Info("No .env file found, using environment variables")
	}

	config := &Config{
		AppName:  getEnv("APP_NAME", "fiap-sa-payment-service"),
		AppEnv:   environment,
		MongoURI: fetchEnv("DATABASE_URI"),
		Port:     getEnv("PORT", "8080"),
	}

	return config
}

func (cfg Config) IsDevelopment() bool {
	return cfg.AppEnv == "development"
}

func (cfg Config) IsProduction() bool {
	return cfg.AppEnv == "production"
}
