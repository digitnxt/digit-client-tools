package digit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// CreateBoundaries creates boundaries using the boundary API
func CreateBoundaries(serverURL, jwtToken, tenantID, clientID string, boundaryData []map[string]interface{}) (string, error) {
	// Prepare the request payload
	payload := map[string]interface{}{
		"boundary": boundaryData,
	}
	
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal boundary data to JSON: %w", err)
	}
	
	// Make HTTP request to boundary API
	url := fmt.Sprintf("%s/boundary/v1", serverURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}
	
	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-ID", tenantID)
	req.Header.Set("X-Client-Id", clientID)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwtToken))
	
	// Execute request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute HTTP request: %w", err)
	}
	defer resp.Body.Close()
	
	// Read response
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}
	
	// Check status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(responseBody))
	}
	
	return string(responseBody), nil
}