package config

import (
	"log"
	"os"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

type AuthConfig struct {
	Secret   string
	Cost     int
	ExpHours int
}

func LoadAuthConfig() AuthConfig {
	cost, err := strconv.Atoi(os.Getenv("AUTH_COST"))
	if err != nil {
		log.Print("Invalid bcrypt cost. Default value will be used!")
		cost = bcrypt.DefaultCost
	}

	expHours, err := strconv.Atoi(os.Getenv("AUTH_EXP_HOURS"))
	if err != nil {
		log.Print("Invalid expiration hours. Default value will be used!")
		expHours = 1
	}

	return AuthConfig{
		Secret:   os.Getenv("AUTH_SECRET"),
		Cost:     cost,
		ExpHours: expHours,
	}
}
