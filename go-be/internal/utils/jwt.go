package utils

import (
	"golf-booking-go/global"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(adminID int32) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"admin_id": adminID,
	})

	// Sign the token with a secret key (replace "your_secret_key" with your actual secret)
	tokenString, err := token.SignedString([]byte(global.Config.SecretJwtKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
