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

// CreateBoundaryHierarchy creates boundary hierarchy using the boundary hierarchy definition API
func CreateBoundaryHierarchy(serverURL, jwtToken, tenantID, clientID string, boundaryHierarchy map[string]interface{}) (string, error) {
	// Prepare the request payload
	payload := map[string]interface{}{
		"boundaryHierarchy": boundaryHierarchy,
	}
	
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal boundary hierarchy data to JSON: %w", err)
	}
	
	// Make HTTP request to boundary hierarchy definition API
	url := fmt.Sprintf("%s/boundary/v1/boundary-hierarchy-definition", serverURL)
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

// SearchBoundaryHierarchy searches for boundary hierarchy by hierarchy type
func SearchBoundaryHierarchy(serverURL, jwtToken, tenantID, clientID, hierarchyType string) (string, error) {
	// Make HTTP request to boundary hierarchy definition API
	url := fmt.Sprintf("%s/boundary/v1/boundary-hierarchy-definition?hierarchyType=%s", serverURL, hierarchyType)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}
	
	// Set headers
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

// CreateBoundaryRelationships creates boundary relationships using the boundary relationships API
func CreateBoundaryRelationships(serverURL, jwtToken, tenantID, clientID string, boundaryRelationship map[string]interface{}) (string, error) {
	// Prepare the request payload
	payload := map[string]interface{}{
		"boundaryRelationship": boundaryRelationship,
	}
	
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal boundary relationship data to JSON: %w", err)
	}
	
	// Make HTTP request to boundary relationships API
	url := fmt.Sprintf("%s/boundary/v1/boundary-relationships", serverURL)
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

// SearchBoundaryRelationships searches for boundary relationships with query parameters
func SearchBoundaryRelationships(serverURL, jwtToken, tenantID, clientID, hierarchyType, boundaryType, codes string, includeChildren bool) (string, error) {
	// Build query parameters
	queryParams := fmt.Sprintf("hierarchyType=%s&boundaryType=%s", hierarchyType, boundaryType)
	
	if codes != "" {
		queryParams += fmt.Sprintf("&codes=%s", codes)
	}
	
	if includeChildren {
		queryParams += "&includeChildren=true"
	}
	
	// Make HTTP request to boundary relationships search API
	url := fmt.Sprintf("%s/boundary/v1/boundary-relationships?%s", serverURL, queryParams)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}
	
	// Set headers
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