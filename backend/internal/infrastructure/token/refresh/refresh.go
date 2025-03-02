package refresh

import (
	"app/internal/infrastructure/token"
	"app/pkg/config"
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"
)

var (
	errTokenExpirationTime = token.NewError(token.ErrConfiguration, nil, "refresh token expiration time")
	errTokenValidation     = token.NewError(token.ErrTokenValidation, nil, "invalid expiration time format")
	errTokenGeneration     = token.NewError(token.ErrTokenGeneration, nil, "refresh token")
	errTokenExpired        = token.NewError(token.ErrRefreshTokenExpired, nil)
	errTokenInvalid        = token.NewError(token.ErrRefreshTokenInvalid, nil)
)

// Generate refresh-token and save it in the database
func GenerateRefreshToken() (string, string, error) {
	// Generate refreshToken
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", "", errTokenGeneration
	}

	// Codificar en base64
	refreshToken := base64.URLEncoding.EncodeToString(b)
	timeExp, err := time.ParseDuration(config.ENV().REFRESH_EXPIRATION_TIME)
	if err != nil {
		return "", "", errTokenExpirationTime
	}
	refreshTokenExp := time.Now().Add(timeExp).Format("2006-01-02 15:04:05")

	return refreshToken, refreshTokenExp, nil
}

// ValidateRefreshToken comprueba la validez del refreshToken
func ValidateRefreshToken(refreshToken, refreshTokenExp string) (int, error) {
	if refreshToken == "" {
		return http.StatusUnauthorized, errTokenInvalid
	}

	// Parsing expiration time
	expiresAt, err := time.Parse("2006-01-02 15:04:05", refreshTokenExp)
	if err != nil {
		return http.StatusInternalServerError, errTokenValidation
	}

	// Check if the refreshToken has expired
	if time.Now().UTC().After(expiresAt) {
		return http.StatusUnauthorized, errTokenExpired
	}

	return http.StatusOK, nil
}
