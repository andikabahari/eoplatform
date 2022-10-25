package helper

import (
	"time"

	"github.com/andikabahari/eoplatform/config"
	"github.com/golang-jwt/jwt"
)

type JWTCustomClaims struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	jwt.StandardClaims
}

func CreateToken(id uint, name string) (string, error) {
	authConfig := config.LoadAuthConfig()

	claims := JWTCustomClaims{
		ID:   id,
		Name: name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(authConfig.Secret))
	if err != nil {
		return "", err
	}

	return token, nil
}
