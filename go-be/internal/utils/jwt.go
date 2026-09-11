package utils

import (
	"golf-booking-go/global"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateTokens(adminID int32) (string, string, error) {
	// generate access token and refresh token

	// 5min
	accessTokenClaims := jwt.MapClaims{
		"admin_id": adminID,
		"exp":      jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
	}

	// 7 days
	refreshTokenClaims := jwt.MapClaims{
		"admin_id": adminID,
		"exp":      jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)

	signedAccessToken, err := accessToken.SignedString([]byte(global.Config.Secret.JwtSecret))
	if err != nil {
		return "", "", err
	}

	signedRefreshToken, err := refreshToken.SignedString([]byte(global.Config.Secret.JwtSecret))
	if err != nil {
		return "", "", err
	}

	return signedAccessToken, signedRefreshToken, nil
}
