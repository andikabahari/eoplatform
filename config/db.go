package config

import "os"

type DBConfig struct {
	Driver string
	User   string
	Pass   string
	Name   string
	Host   string
	Port   string
}

func LoadDBConfig() DBConfig {
	return DBConfig{
		Driver: os.Getenv("DB_DRIVER"),
		User:   os.Getenv("DB_USER"),
		Pass:   os.Getenv("DB_PASS"),
		Name:   os.Getenv("DB_NAME"),
		Host:   os.Getenv("DB_HOST"),
		Port:   os.Getenv("DB_PORT"),
	}
}
