package config

import (
	"log"
	"os"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

type AuthConfig struct {
	Secret string
	Cost   int
}

func LoadAuthConfig() AuthConfig {
	cost, err := strconv.Atoi(os.Getenv("AUTH_COST"))
	if err != nil {
		log.Print("Invalid bcrypt cost. Default value will be used!")
		cost = bcrypt.DefaultCost
	}

	return AuthConfig{
		Secret: os.Getenv("AUTH_SECRET"),
		Cost:   cost,
	}
}
