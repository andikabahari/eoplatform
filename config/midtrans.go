package config

import "os"

type MidtransConfig struct {
	BaseURL   string
	ServerKey string
}

func LoadMidtransConfig() MidtransConfig {
	return MidtransConfig{
		BaseURL:   os.Getenv("MIDTRANS_BASE_URL"),
		ServerKey: os.Getenv("MIDTRANS_SERVER_KEY"),
	}
}
