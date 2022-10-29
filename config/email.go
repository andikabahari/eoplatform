package config

import "os"

type EmailConfig struct {
	Address  string
	Password string
}

func LoadEmailConfig() EmailConfig {
	return EmailConfig{
		Address:  os.Getenv("EMAIL_ADDRESS"),
		Password: os.Getenv("EMAIL_PASSWORD"),
	}
}
