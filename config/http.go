package config

import "os"

type HTTPConfig struct {
	Port string
}

func LoadHTTPConfig() HTTPConfig {
	return HTTPConfig{
		Port: os.Getenv("HTTP_PORT"),
	}
}
