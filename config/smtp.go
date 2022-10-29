package config

import "os"

type SMTPConfig struct {
	Host string
	Port string
}

func LoadSMTPConfig() SMTPConfig {
	return SMTPConfig{
		Host: os.Getenv("SMTP_HOST"),
		Port: os.Getenv("SMTP_PORT"),
	}
}
