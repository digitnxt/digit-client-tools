package digit

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

// CreateSchema creates a new MDMS schema
func CreateSchema(serverURL, jwtToken, tenantID, clientID, code, description, definition string, isActive bool) (string, error) {
	if serverURL == "" {
		return "", fmt.Errorf("server URL is required")
	}
	if tenantID == "" {
		return "", fmt.Errorf("tenant ID is required")
	}
	if clientID == "" {
		return "", fmt.Errorf("client ID is required")
	}
	if code == "" {
		return "", fmt.Errorf("code is required")
	}
	if description == "" {
		return "", fmt.Errorf("description is required")
	}
	if definition == "" {
		return "", fmt.Errorf("definition is required")
	}

	url := strings.TrimSuffix(serverURL, "/") + "/mdms-v2/v1/schema"

	// Prepare request body
	requestBody := fmt.Sprintf(`{
  "SchemaDefinition": {
    "code": "%s",
    "description": "%s",
    "definition": %s,
    "isActive": %t
  }
}`, code, description, definition, isActive)

	// Create request
	req, err := http.NewRequest("POST", url, strings.NewReader(requestBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-tenant-id", tenantID)
	req.Header.Set("x-client-id", clientID)
	if jwtToken != "" {
		req.Header.Set("Authorization", "Bearer "+jwtToken)
	}

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	// Check status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return string(body), nil
}

// CreateMdmsData creates MDMS data entries
func CreateMdmsData(serverURL, jwtToken, tenantID, clientID, mdmsData string) (string, error) {
	if serverURL == "" {
		return "", fmt.Errorf("server URL is required")
	}
	if tenantID == "" {
		return "", fmt.Errorf("tenant ID is required")
	}
	if clientID == "" {
		return "", fmt.Errorf("client ID is required")
	}
	if mdmsData == "" {
		return "", fmt.Errorf("MDMS data is required")
	}

	url := strings.TrimSuffix(serverURL, "/") + "/mdms-v2/v2"

	// Debug log the input data
	fmt.Printf("Debug - MDMS data input: %s\n", mdmsData)

	// Prepare request body
	requestBody := fmt.Sprintf(`{
  "Mdms": %s
}`, mdmsData)

	// Debug log the final request body
	fmt.Printf("Debug - Final request body: %s\n", requestBody)

	// Create request
	req, err := http.NewRequest("POST", url, strings.NewReader(requestBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-tenant-id", tenantID)
	req.Header.Set("x-client-id", clientID)
	if jwtToken != "" {
		req.Header.Set("Authorization", "Bearer "+jwtToken)
	}

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	// Check status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return string(body), nil
}

// SearchSchema searches for an MDMS schema by code
func SearchSchema(serverURL, jwtToken, tenantID, clientID, schemaCode string) (string, error) {
	if serverURL == "" {
		return "", fmt.Errorf("server URL is required")
	}
	if tenantID == "" {
		return "", fmt.Errorf("tenant ID is required")
	}
	if clientID == "" {
		return "", fmt.Errorf("client ID is required")
	}
	if schemaCode == "" {
		return "", fmt.Errorf("schema code is required")
	}

	url := strings.TrimSuffix(serverURL, "/") + "/mdms-v2/v1/schema?code=" + schemaCode

	// Create request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("X-tenant-ID", tenantID)
	req.Header.Set("X-Client-ID", clientID)
	if jwtToken != "" {
		req.Header.Set("Authorization", "Bearer "+jwtToken)
	}

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	// Check status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return string(body), nil
}

// SearchMdmsData searches for MDMS data by schema code and optional unique identifiers
func SearchMdmsData(serverURL, jwtToken, tenantID, clientID, schemaCode, uniqueIdentifiers string) (string, error) {
	if serverURL == "" {
		return "", fmt.Errorf("server URL is required")
	}
	if tenantID == "" {
		return "", fmt.Errorf("tenant ID is required")
	}
	if clientID == "" {
		return "", fmt.Errorf("client ID is required")
	}
	if schemaCode == "" {
		return "", fmt.Errorf("schema code is required")
	}

	url := strings.TrimSuffix(serverURL, "/") + "/mdms-v2/v2?schemaCode=" + schemaCode
	
	// Add uniqueIdentifiers parameter if provided
	if uniqueIdentifiers != "" {
		url += "&uniqueIdentifiers=" + uniqueIdentifiers
	}

	// Create request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("X-Tenant-ID", tenantID)
	req.Header.Set("X-Client-Id", clientID)
	if jwtToken != "" {
		req.Header.Set("Authorization", "Bearer "+jwtToken)
	}

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	// Check status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return string(body), nil
}
