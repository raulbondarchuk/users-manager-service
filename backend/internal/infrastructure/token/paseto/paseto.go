package paseto

import (
	"app/pkg/config"
	"crypto/sha256"
	"errors"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/o1egl/paseto"
)

// PasetoManager manages PASETO tokens
type PasetoManager struct {
	baseKey        string
	expirationTime time.Duration
}

var (
	instance *PasetoManager
	once     sync.Once

	BearerPrefix = "bearer "
)

// Paseto returns the singleton PasetoManager
func Paseto() *PasetoManager {
	once.Do(func() {
		expirationTime, err := time.ParseDuration(config.ENV().PASETO_EXPIRATION_TIME)
		if err != nil {
			expirationTime = 6 * time.Hour
		}

		instance = &PasetoManager{
			baseKey:        config.ENV().PASETO_SK,
			expirationTime: expirationTime,
		}
	})
	return instance
}

// GenerateToken creates a new PASETO token and returns it along with the claims
func (p *PasetoManager) GenerateToken(claims PasetoClaims) (string, *PasetoClaims, error) {
	if claims.Username == "" {
		return "", nil, errors.New("missing username in claims")
	}

	// Set expiration and issued times if not set
	if claims.ExpiresAt.IsZero() {
		claims.ExpiresAt = time.Now().Add(p.expirationTime)
	}
	if claims.IssuedAt.IsZero() {
		claims.IssuedAt = time.Now()
	}

	// Prepare JSONToken from paseto
	jsonToken := paseto.JSONToken{
		Subject:    claims.Username,
		IssuedAt:   claims.IssuedAt,
		Expiration: claims.ExpiresAt,
	}

	// Set custom fields
	jsonToken.Set("companyId", strconv.Itoa(claims.CompanyID))
	jsonToken.Set("companyName", claims.CompanyName)
	jsonToken.Set("roles", claims.Roles)
	jsonToken.Set("isPrimary", strconv.FormatBool(claims.IsPrimary))

	if claims.OwnerUsername != "" {
		jsonToken.Set("ownerUsername", claims.OwnerUsername)
	}

	// Form the key (symmetricKey)
	key := sha256.Sum256([]byte(p.baseKey))
	symmetricKey := key[:]

	// Encrypt
	token, err := paseto.NewV2().Encrypt(symmetricKey, jsonToken, nil)
	if err != nil {
		return "", nil, errors.New("error generating token")
	}

	return token, &claims, nil
}

// ValidateToken validates a PASETO token and checks expiration
func (p *PasetoManager) ValidateToken(tokenStr string) (*PasetoClaims, error) {
	return p.validateTokenInternal(tokenStr, true)
}

// ValidateTokenWithoutExpirationCheck validates a PASETO token without checking expiration
func (p *PasetoManager) ValidateTokenWithoutExpirationCheck(tokenStr string) (*PasetoClaims, error) {
	return p.validateTokenInternal(tokenStr, false)
}

// Internal method to validate token with optional expiration check
func (p *PasetoManager) validateTokenInternal(tokenStr string, checkExpiration bool) (*PasetoClaims, error) {
	tokenStr = strings.TrimSpace(tokenStr)
	if strings.HasPrefix(strings.ToLower(tokenStr), BearerPrefix) {
		tokenStr = tokenStr[len(BearerPrefix):]
	}

	var jsonToken paseto.JSONToken

	// Get key
	key := sha256.Sum256([]byte(p.baseKey))
	symmetricKey := key[:]

	err := paseto.NewV2().Decrypt(tokenStr, symmetricKey, &jsonToken, nil)
	if err != nil {
		return nil, errors.New("token decrypt error: " + err.Error())
	}

	// Check if the token has expired
	if checkExpiration && time.Now().After(jsonToken.Expiration) {
		return nil, errors.New("token expired")
	}

	// Collect PasetoClaims
	claims := &PasetoClaims{
		Username:  jsonToken.Subject,
		IssuedAt:  jsonToken.IssuedAt,
		ExpiresAt: jsonToken.Expiration,
	}

	if ownerUsername := jsonToken.Get("ownerUsername"); ownerUsername != "" {
		claims.OwnerUsername = ownerUsername
	}

	// Get other fields
	if companyIdStr := jsonToken.Get("companyId"); companyIdStr != "" {
		companyID, err := strconv.Atoi(companyIdStr)
		if err == nil {
			claims.CompanyID = companyID
		}
	}
	claims.CompanyName = jsonToken.Get("companyName")
	claims.Roles = jsonToken.Get("roles")
	if isPrimaryStr := jsonToken.Get("isPrimary"); isPrimaryStr != "" {
		claims.IsPrimary, _ = strconv.ParseBool(isPrimaryStr)
	}

	// Check required fields
	if claims.Username == "" {
		return nil, errors.New("missing username in token")
	}

	return claims, nil
}
