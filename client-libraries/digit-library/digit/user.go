package digit

import (
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
)

// UserRequest represents the request structure for creating a user in Keycloak
type UserRequest struct {
	Username      string                 `json:"username"`
	Email         string                 `json:"email"`
	Enabled       bool                   `json:"enabled"`
	EmailVerified bool                   `json:"emailVerified"`
	Credentials   []UserCredential       `json:"credentials"`
	Attributes    map[string]interface{} `json:"attributes,omitempty"`
}

// UserCredential represents the password credential for a user
type UserCredential struct {
	Type      string `json:"type"`
	Value     string `json:"value"`
	Temporary bool   `json:"temporary"`
}

// CreateUser creates a new user in Keycloak
// Returns the raw response body as string and any error encountered
func CreateUser(serverURL, jwtToken, realm, username, password, email string) (string, error) {
	// Validate required parameters
	if serverURL == "" {
		return "", fmt.Errorf("serverURL cannot be empty")
	}
	if jwtToken == "" {
		return "", fmt.Errorf("jwtToken cannot be empty")
	}
	if realm == "" {
		return "", fmt.Errorf("realm cannot be empty")
	}
	if username == "" {
		return "", fmt.Errorf("username cannot be empty")
	}
	if password == "" {
		return "", fmt.Errorf("password cannot be empty")
	}
	if email == "" {
		return "", fmt.Errorf("email cannot be empty")
	}

	// Create the request payload
	userReq := UserRequest{
		Username:      username,
		Email:         email,
		Enabled:       true,
		EmailVerified: true,
		Credentials: []UserCredential{
			{
				Type:      "password",
				Value:     password,
				Temporary: false,
			},
		},
		Attributes: make(map[string]interface{}),
	}

	// Create HTTP client
	client := resty.New()

	// Make the API request
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+jwtToken).
		SetBody(userReq).
		Post(serverURL + "/keycloak/admin/realms/" + realm + "/users")

	if err != nil {
		return "", fmt.Errorf("failed to make API request: %w", err)
	}

	// Return the raw response body as string
	return string(resp.Body()), nil
}

// ResetPassword resets a user's password in Keycloak
// Returns the raw response body as string and any error encountered
func ResetPassword(serverURL, jwtToken, realm, username, newPassword string) (string, error) {
	// Validate required parameters
	if serverURL == "" {
		return "", fmt.Errorf("serverURL cannot be empty")
	}
	if jwtToken == "" {
		return "", fmt.Errorf("jwtToken cannot be empty")
	}
	if realm == "" {
		return "", fmt.Errorf("realm cannot be empty")
	}
	if username == "" {
		return "", fmt.Errorf("username cannot be empty")
	}
	if newPassword == "" {
		return "", fmt.Errorf("newPassword cannot be empty")
	}

	// Create HTTP client
	client := resty.New()

	// First, get the user ID by username
	getUserResp, err := client.R().
		SetHeader("Authorization", "Bearer "+jwtToken).
		SetQueryParam("username", username).
		Get(serverURL + "/keycloak/admin/realms/" + realm + "/users")

	if err != nil {
		return "", fmt.Errorf("failed to get user: %w", err)
	}

	// Parse the response to extract user ID (simplified - assumes first user in array)
	// In a real implementation, you'd want to parse JSON properly
	userRespBody := string(getUserResp.Body())
	if getUserResp.StatusCode() != 200 {
		return "", fmt.Errorf("failed to find user %s: %s", username, userRespBody)
	}

	// For simplicity, we'll extract the user ID from the response
	// This is a simplified approach - in production you'd use proper JSON parsing
	if len(userRespBody) < 10 || userRespBody == "[]" {
		return "", fmt.Errorf("user %s not found", username)
	}

	// Extract user ID from response (this is a simplified extraction)
	// In production, you should use proper JSON parsing
	start := strings.Index(userRespBody, `"id":"`) + 6
	if start < 6 {
		return "", fmt.Errorf("could not extract user ID from response")
	}
	end := strings.Index(userRespBody[start:], `"`) + start
	if end <= start {
		return "", fmt.Errorf("could not extract user ID from response")
	}
	userID := userRespBody[start:end]

	// Create the password reset payload
	passwordResetReq := UserCredential{
		Type:      "password",
		Value:     newPassword,
		Temporary: false,
	}

	// Reset the password
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+jwtToken).
		SetBody(passwordResetReq).
		Put(serverURL + "/keycloak/admin/realms/" + realm + "/users/" + userID + "/reset-password")

	if err != nil {
		return "", fmt.Errorf("failed to reset password: %w", err)
	}

	// Return the raw response body as string
	return string(resp.Body()), nil
}

// DeleteUser deletes a user from Keycloak
// Returns the raw response body as string and any error encountered
func DeleteUser(serverURL, jwtToken, realm, username string) (string, error) {
	// Validate required parameters
	if serverURL == "" {
		return "", fmt.Errorf("serverURL cannot be empty")
	}
	if jwtToken == "" {
		return "", fmt.Errorf("jwtToken cannot be empty")
	}
	if realm == "" {
		return "", fmt.Errorf("realm cannot be empty")
	}
	if username == "" {
		return "", fmt.Errorf("username cannot be empty")
	}

	// Create HTTP client
	client := resty.New()

	// First, get the user ID by username
	getUserResp, err := client.R().
		SetHeader("Authorization", "Bearer "+jwtToken).
		SetQueryParam("username", username).
		Get(serverURL + "/keycloak/admin/realms/" + realm + "/users")

	if err != nil {
		return "", fmt.Errorf("failed to get user: %w", err)
	}

	// Parse the response to extract user ID
	userRespBody := string(getUserResp.Body())
	if getUserResp.StatusCode() != 200 {
		return "", fmt.Errorf("failed to find user %s: %s", username, userRespBody)
	}

	// Check if user exists
	if len(userRespBody) < 10 || userRespBody == "[]" {
		return "", fmt.Errorf("user %s not found", username)
	}

	// Extract user ID from response (simplified extraction)
	start := strings.Index(userRespBody, `"id":"`) + 6
	if start < 6 {
		return "", fmt.Errorf("could not extract user ID from response")
	}
	end := strings.Index(userRespBody[start:], `"`) + start
	if end <= start {
		return "", fmt.Errorf("could not extract user ID from response")
	}
	userID := userRespBody[start:end]

	// Delete the user
	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+jwtToken).
		Delete(serverURL + "/keycloak/admin/realms/" + realm + "/users/" + userID)

	if err != nil {
		return "", fmt.Errorf("failed to delete user: %w", err)
	}

	// Return the raw response body as string
	return string(resp.Body()), nil
}

// SearchUser searches for users in Keycloak
// If username is provided, searches for that specific user
// If username is empty, returns all users in the realm
// Returns the raw response body as string and any error encountered
func SearchUser(serverURL, jwtToken, realm, username string) (string, error) {
	// Validate required parameters
	if serverURL == "" {
		return "", fmt.Errorf("serverURL cannot be empty")
	}
	if jwtToken == "" {
		return "", fmt.Errorf("jwtToken cannot be empty")
	}
	if realm == "" {
		return "", fmt.Errorf("realm cannot be empty")
	}

	// Create HTTP client
	client := resty.New()

	// Build the request
	req := client.R().
		SetHeader("Authorization", "Bearer "+jwtToken)

	// Add username query parameter if provided
	if username != "" {
		req = req.SetQueryParam("username", username)
	}

	// Search for users
	resp, err := req.Get(serverURL + "/keycloak/admin/realms/" + realm + "/users")

	if err != nil {
		return "", fmt.Errorf("failed to search users: %w", err)
	}

	// Check for successful response
	if resp.StatusCode() != 200 {
		return "", fmt.Errorf("failed to search users: %s", string(resp.Body()))
	}

	// Return the raw response body as string
	return string(resp.Body()), nil
}

// UserUpdateRequest represents the request structure for updating a user in Keycloak
type UserUpdateRequest struct {
	Username      string                 `json:"username,omitempty"`
	Email         string                 `json:"email,omitempty"`
	FirstName     string                 `json:"firstName,omitempty"`
	LastName      string                 `json:"lastName,omitempty"`
	Enabled       *bool                  `json:"enabled,omitempty"`
	EmailVerified *bool                  `json:"emailVerified,omitempty"`
	Attributes    map[string]interface{} `json:"attributes,omitempty"`
}

// UpdateUser updates a user in Keycloak
// Returns the raw response body as string and any error encountered
func UpdateUser(serverURL, jwtToken, realm, username, email, firstName, lastName string, enabled *bool) (string, error) {
	// Validate required parameters
	if serverURL == "" {
		return "", fmt.Errorf("serverURL cannot be empty")
	}
	if jwtToken == "" {
		return "", fmt.Errorf("jwtToken cannot be empty")
	}
	if realm == "" {
		return "", fmt.Errorf("realm cannot be empty")
	}
	if username == "" {
		return "", fmt.Errorf("username cannot be empty")
	}

	// Create HTTP client
	client := resty.New()

	// First, get the user ID by username
	getUserResp, err := client.R().
		SetHeader("Authorization", "Bearer "+jwtToken).
		SetQueryParam("username", username).
		Get(serverURL + "/keycloak/admin/realms/" + realm + "/users")

	if err != nil {
		return "", fmt.Errorf("failed to get user: %w", err)
	}

	// Parse the response to extract user ID
	userRespBody := string(getUserResp.Body())
	if getUserResp.StatusCode() != 200 {
		return "", fmt.Errorf("failed to find user %s: %s", username, userRespBody)
	}

	// Check if user exists
	if len(userRespBody) < 10 || userRespBody == "[]" {
		return "", fmt.Errorf("user %s not found", username)
	}

	// Extract user ID from response (simplified extraction)
	start := strings.Index(userRespBody, `"id":"`) + 6
	if start < 6 {
		return "", fmt.Errorf("could not extract user ID from response")
	}
	end := strings.Index(userRespBody[start:], `"`) + start
	if end <= start {
		return "", fmt.Errorf("could not extract user ID from response")
	}
	userID := userRespBody[start:end]

	// Build update request with only provided fields
	updateReq := UserUpdateRequest{}
	
	if email != "" {
		updateReq.Email = email
	}
	if firstName != "" {
		updateReq.FirstName = firstName
	}
	if lastName != "" {
		updateReq.LastName = lastName
	}
	if enabled != nil {
		updateReq.Enabled = enabled
	}

	// Update the user
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+jwtToken).
		SetBody(updateReq).
		Put(serverURL + "/keycloak/admin/realms/" + realm + "/users/" + userID)

	if err != nil {
		return "", fmt.Errorf("failed to update user: %w", err)
	}

	// Return the raw response body as string
	return string(resp.Body()), nil
}

// RoleRequest represents the request structure for creating a role in Keycloak
type RoleRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Composite   bool   `json:"composite,omitempty"`
}

// CreateRole creates a new role in Keycloak
// Returns the raw response body as string and any error encountered
func CreateRole(serverURL, jwtToken, realm, roleName, description string) (string, error) {
	// Validate required parameters
	if serverURL == "" {
		return "", fmt.Errorf("serverURL cannot be empty")
	}
	if jwtToken == "" {
		return "", fmt.Errorf("jwtToken cannot be empty")
	}
	if realm == "" {
		return "", fmt.Errorf("realm cannot be empty")
	}
	if roleName == "" {
		return "", fmt.Errorf("roleName cannot be empty")
	}

	// Create the request payload
	roleReq := RoleRequest{
		Name:        roleName,
		Description: description,
		Composite:   false,
	}

	// Create HTTP client
	client := resty.New()

	// Make the API request
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+jwtToken).
		SetBody(roleReq).
		Post(serverURL + "/keycloak/admin/realms/" + realm + "/roles")

	if err != nil {
		return "", fmt.Errorf("failed to make API request: %w", err)
	}

	// Return the raw response body as string
	return string(resp.Body()), nil
}

// AssignRoleToUser assigns a role to a user in Keycloak
// Returns the raw response body as string and any error encountered
func AssignRoleToUser(serverURL, jwtToken, realm, username, roleName string) (string, error) {
	// Validate required parameters
	if serverURL == "" {
		return "", fmt.Errorf("serverURL cannot be empty")
	}
	if jwtToken == "" {
		return "", fmt.Errorf("jwtToken cannot be empty")
	}
	if realm == "" {
		return "", fmt.Errorf("realm cannot be empty")
	}
	if username == "" {
		return "", fmt.Errorf("username cannot be empty")
	}
	if roleName == "" {
		return "", fmt.Errorf("roleName cannot be empty")
	}

	// Create HTTP client
	client := resty.New()

	// First, get the user ID by username
	getUserResp, err := client.R().
		SetHeader("Authorization", "Bearer "+jwtToken).
		SetQueryParam("username", username).
		Get(serverURL + "/keycloak/admin/realms/" + realm + "/users")

	if err != nil {
		return "", fmt.Errorf("failed to get user: %w", err)
	}

	// Parse the response to extract user ID
	userRespBody := string(getUserResp.Body())
	if getUserResp.StatusCode() != 200 {
		return "", fmt.Errorf("failed to find user %s: %s", username, userRespBody)
	}

	// Check if user exists
	if len(userRespBody) < 10 || userRespBody == "[]" {
		return "", fmt.Errorf("user %s not found", username)
	}

	// Extract user ID from response (simplified extraction)
	start := strings.Index(userRespBody, `"id":"`) + 6
	if start < 6 {
		return "", fmt.Errorf("could not extract user ID from response")
	}
	end := strings.Index(userRespBody[start:], `"`) + start
	if end <= start {
		return "", fmt.Errorf("could not extract user ID from response")
	}
	userID := userRespBody[start:end]

	// Get the role details by name
	getRoleResp, err := client.R().
		SetHeader("Authorization", "Bearer "+jwtToken).
		Get(serverURL + "/keycloak/admin/realms/" + realm + "/roles/" + roleName)

	if err != nil {
		return "", fmt.Errorf("failed to get role: %w", err)
	}

	// Check if role exists
	roleRespBody := string(getRoleResp.Body())
	if getRoleResp.StatusCode() != 200 {
		return "", fmt.Errorf("role %s not found: %s", roleName, roleRespBody)
	}

	// Create role mapping payload (array of role objects)
	// We need to send the complete role object, not just the name
	roleMappingPayload := "[" + roleRespBody + "]"

	// Assign the role to the user
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+jwtToken).
		SetBody(roleMappingPayload).
		Post(serverURL + "/keycloak/admin/realms/" + realm + "/users/" + userID + "/role-mappings/realm")

	if err != nil {
		return "", fmt.Errorf("failed to assign role to user: %w", err)
	}

	// Return the raw response body as string
	return string(resp.Body()), nil
}
