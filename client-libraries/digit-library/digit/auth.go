package digit

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
)

// TokenResponse represents the response from Keycloak token endpoint
type TokenResponse struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	NotBeforePolicy  int    `json:"not-before-policy"`
	SessionState     string `json:"session_state"`
	Scope            string `json:"scope"`
}

// GetJWTToken retrieves JWT token from Keycloak using username/password with client secret
// Returns the JWT token string and any error encountered
func GetJWTToken(server, realm, clientID, clientSecret, username, password string) (string, error) {
	// Validate required parameters
	if server == "" {
		return "", fmt.Errorf("server cannot be empty")
	}
	if realm == "" {
		return "", fmt.Errorf("realm cannot be empty")
	}
	if clientID == "" {
		return "", fmt.Errorf("clientID cannot be empty")
	}
	if clientSecret == "" {
		return "", fmt.Errorf("clientSecret cannot be empty")
	}
	if username == "" {
		return "", fmt.Errorf("username cannot be empty")
	}
	if password == "" {
		return "", fmt.Errorf("password cannot be empty")
	}

	// Create HTTP client
	client := resty.New()

	// Prepare form data for token request
	formData := map[string]string{
		"grant_type":    "password",
		"client_id":     clientID,
		"client_secret": clientSecret,
		"username":      username,
		"password":      password,
	}

	// Make the token request
	resp, err := client.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(formData).
		Post(server + "/keycloak/realms/" + realm + "/protocol/openid-connect/token")

	if err != nil {
		return "", fmt.Errorf("failed to make token request: %w", err)
	}

	// Check for successful response
	if resp.StatusCode() != 200 {
		return "", fmt.Errorf("failed to get token: %s", string(resp.Body()))
	}

	// Parse the response
	var tokenResp TokenResponse
	if err := json.Unmarshal(resp.Body(), &tokenResp); err != nil {
		return "", fmt.Errorf("failed to parse token response: %w", err)
	}

	return tokenResp.AccessToken, nil
}
