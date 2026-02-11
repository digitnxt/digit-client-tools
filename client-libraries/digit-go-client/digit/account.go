package digit

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

// TenantRequest represents the request structure for creating an account
type TenantRequest struct {
	Tenant Tenant `json:"tenant"`
}

// Tenant represents the tenant information
type Tenant struct {
	Name                 string                 `json:"name"`
	Email                string                 `json:"email"`
	IsActive             bool                   `json:"isActive"`
	AdditionalAttributes map[string]interface{} `json:"additionalAttributes"`
}

// CreateAccount creates a new account in DIGIT services
// Returns the raw response body as string and any error encountered
func CreateAccount(serverURL, clientID, name, email string, active bool) (string, error) {
	// Validate required parameters
	if serverURL == "" {
		return "", fmt.Errorf("serverURL cannot be empty")
	}
	if name == "" {
		return "", fmt.Errorf("name cannot be empty")
	}
	if email == "" {
		return "", fmt.Errorf("email cannot be empty")
	}
	if clientID == "" {
		return "", fmt.Errorf("clientID cannot be empty")
	}

	// Create the request payload
	tenantReq := TenantRequest{
		Tenant: Tenant{
			Name:                 name,
			Email:                email,
			IsActive:             active,
			AdditionalAttributes: make(map[string]interface{}),
		},
	}

	// Create HTTP client
	client := resty.New()

	// Make the API request
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("X-Client-Id", clientID).
		SetBody(tenantReq).
		Post(serverURL + "/account/v1")

	if err != nil {
		return "", fmt.Errorf("failed to make API request: %w", err)
	}

	// Return the raw response body as string
	return string(resp.Body()), nil
}
