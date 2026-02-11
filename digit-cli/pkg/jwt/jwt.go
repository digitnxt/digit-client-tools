package jwt

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// JWTClaims represents the JWT payload structure
type JWTClaims struct {
	Sub               string `json:"sub"`
	Iss               string `json:"iss"`
	PreferredUsername string `json:"preferred_username"`
	Email             string `json:"email"`
	Name              string `json:"name"`
	GivenName         string `json:"given_name"`
	FamilyName        string `json:"family_name"`
	Exp               int64  `json:"exp"`
	Iat               int64  `json:"iat"`
}

// DecodeJWT decodes a JWT token and extracts claims
func DecodeJWT(token string) (*JWTClaims, error) {
	// Split the token into parts
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid JWT token format")
	}

	// Decode the payload (second part)
	payload := parts[1]
	
	// Add padding if needed
	for len(payload)%4 != 0 {
		payload += "="
	}

	// Base64 decode
	decoded, err := base64.URLEncoding.DecodeString(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to decode JWT payload: %w", err)
	}

	// Parse JSON
	var claims JWTClaims
	if err := json.Unmarshal(decoded, &claims); err != nil {
		return nil, fmt.Errorf("failed to parse JWT claims: %w", err)
	}

	return &claims, nil
}

// ExtractTenantID extracts tenant ID from JWT token
// Tenant ID is derived from the realm name in the issuer URL
func ExtractTenantID(token string) (string, error) {
	claims, err := DecodeJWT(token)
	if err != nil {
		return "", err
	}

	// Extract realm from issuer URL
	// Format: https://domain/keycloak-test/realms/REALM_NAME
	parts := strings.Split(claims.Iss, "/")
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid issuer format in JWT")
	}

	// Get the last part which should be the realm name
	realm := parts[len(parts)-1]
	if realm == "" {
		return "", fmt.Errorf("could not extract realm from JWT issuer")
	}

	return realm, nil
}

// ExtractClientID extracts client ID from JWT token
// Client ID is the user ID (sub claim)
func ExtractClientID(token string) (string, error) {
	claims, err := DecodeJWT(token)
	if err != nil {
		return "", err
	}

	if claims.Sub == "" {
		return "", fmt.Errorf("no subject (sub) found in JWT token")
	}

	return claims.Sub, nil
}

// IsTokenExpired checks if a JWT token is expired
func IsTokenExpired(token string) (bool, error) {
	claims, err := DecodeJWT(token)
	if err != nil {
		return true, err
	}

	// Check if token is expired (with 30 second buffer)
	expiryTime := time.Unix(claims.Exp, 0)
	return time.Now().Add(30*time.Second).After(expiryTime), nil
}

// GetTokenExpiryTime returns the expiry time of a JWT token
func GetTokenExpiryTime(token string) (time.Time, error) {
	claims, err := DecodeJWT(token)
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(claims.Exp, 0), nil
}
