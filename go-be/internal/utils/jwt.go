package utils

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"golf-booking-go/global"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const RefreshTokenTTL = 7 * 24 * time.Hour

var ErrInvalidToken = errors.New("invalid or expired token")

type TokenClaims struct {
	AdminID   int32  `json:"admin_id"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

func GenerateTokens(adminID int32) (string, string, error) {
	if global.Config == nil || global.Config.Secret.JwtSecret == "" || adminID <= 0 {
		return "", "", errors.New("invalid JWT configuration or admin ID")
	}
	var ids [64]byte
	if _, err := rand.Read(ids[:]); err != nil {
		return "", "", err
	}
	now := time.Now()
	secret := []byte(global.Config.Secret.JwtSecret)
	accessClaims := TokenClaims{
		AdminID:   adminID,
		TokenType: "access",
		ID:        hex.EncodeToString(ids[:32]),
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(5 * time.Minute)),
	}
	access, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(secret)
	if err != nil {
		return "", "", err
	}

	refreshClaims := TokenClaims{
		AdminID:   adminID,
		TokenType: "refresh",
		ID:        hex.EncodeToString(ids[32:]),
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(RefreshTokenTTL)),
	}
	refresh, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(secret)
	if err != nil {
		return "", "", err
	}
	return access, refresh, nil
}

// ParseToken requires the expected type so refresh tokens cannot authenticate API requests.
func ParseToken(raw, expectedType string) (*TokenClaims, error) {
	if global.Config == nil || global.Config.Secret.JwtSecret == "" {
		return nil, errors.New("JWT secret is not configured")
	}
	if expectedType != "access" && expectedType != "refresh" {
		return nil, ErrInvalidToken
	}
	claims := new(TokenClaims)
	token, err := jwt.ParseWithClaims(
		raw,
		claims,
		func(*jwt.Token) (interface{}, error) {
			return []byte(global.Config.Secret.JwtSecret), nil
		},
		jwt.WithValidMethods([]string{"HS256"}),
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt())
	if err != nil || !token.Valid || claims.AdminID <= 0 || claims.ID == "" || claims.IssuedAt == nil || claims.TokenType != expectedType {
		return nil, ErrInvalidToken
	}
	return claims, nil
}
