package helper

import (
	"time"

	"github.com/andikabahari/eoplatform/config"
	"github.com/golang-jwt/jwt"
)

type JWTCustomClaims struct {
	ID   uint   `json:"id"`
	Role string `json:"role"`
	jwt.StandardClaims
}

func CreateToken(id uint, role string) (string, error) {
	exp := time.Duration(config.LoadAuthConfig().ExpHours) * time.Hour
	claims := JWTCustomClaims{
		id,
		role,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(exp).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString([]byte(config.LoadAuthConfig().Secret))
	if err != nil {
		return "", err
	}

	return signed, nil
}
