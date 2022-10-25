package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB   DBConfig
	Auth AuthConfig
	HTTP HTTPConfig
}

func NewConfig() *Config {
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	return &Config{
		DB:   LoadDBConfig(),
		Auth: LoadAuthConfig(),
		HTTP: LoadHTTPConfig(),
	}
}
