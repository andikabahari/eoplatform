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
	claims := JWTCustomClaims{
		id,
		name,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString([]byte(config.LoadAuthConfig().Secret))
	if err != nil {
		return "", err
	}

	return signed, nil
}
