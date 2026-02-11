package digit

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

// CreateProcess creates a new workflow process
func CreateProcess(serverURL, jwtToken, tenantID, name, code, description, version string, sla int) (string, error) {
	if serverURL == "" {
		return "", fmt.Errorf("server URL is required")
	}
	if tenantID == "" {
		return "", fmt.Errorf("tenant ID is required")
	}

	url := strings.TrimSuffix(serverURL, "/") + "/workflow/v1/process"

	// Prepare request body
	requestBody := fmt.Sprintf(`{
  "name": "%s",
  "code": "%s",
  "description": "%s",
  "version": "%s",
  "sla": %d
}`, name, code, description, version, sla)

	// Create request
	req, err := http.NewRequest("POST", url, strings.NewReader(requestBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwtToken)
	req.Header.Set("X-Tenant-ID", tenantID)

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

// CreateState creates a new state for a workflow process
func CreateState(serverURL, jwtToken, tenantID, processID, code, name string, isInitial, isParallel, isJoin bool, sla int) (string, error) {
	if serverURL == "" {
		return "", fmt.Errorf("server URL is required")
	}
	if tenantID == "" {
		return "", fmt.Errorf("tenant ID is required")
	}
	if processID == "" {
		return "", fmt.Errorf("process ID is required")
	}

	url := strings.TrimSuffix(serverURL, "/") + fmt.Sprintf("/workflow/v1/process/%s/state", processID)

	// Prepare request body
	requestBody := fmt.Sprintf(`{
  "code": "%s",
  "name": "%s",
  "isInitial": %t,
  "isParallel": %t,
  "isJoin": %t,
  "sla": %d
}`, code, name, isInitial, isParallel, isJoin, sla)

	// Create request
	req, err := http.NewRequest("POST", url, strings.NewReader(requestBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwtToken)
	req.Header.Set("X-Tenant-ID", tenantID)

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

// CreateAction creates a new action for workflow states
func CreateAction(serverURL, jwtToken, tenantID, stateID, name, nextState string, roles []string, assigneeCheck bool) (string, error) {
	if serverURL == "" {
		return "", fmt.Errorf("server URL is required")
	}
	if tenantID == "" {
		return "", fmt.Errorf("tenant ID is required")
	}
	if stateID == "" {
		return "", fmt.Errorf("state ID is required")
	}

	url := strings.TrimSuffix(serverURL, "/") + fmt.Sprintf("/workflow/v1/state/%s/action", stateID)

	// Prepare roles array as JSON
	rolesJSON := "["
	for i, role := range roles {
		if i > 0 {
			rolesJSON += ", "
		}
		rolesJSON += fmt.Sprintf(`"%s"`, role)
	}
	rolesJSON += "]"

	// Prepare request body
	requestBody := fmt.Sprintf(`{
  "name": "%s",
  "nextState": "%s",
  "attributeValidation": {
    "attributes": {
      "roles": %s
    },
    "assigneeCheck": %t
  }
}`, name, nextState, rolesJSON, assigneeCheck)

	// Create request
	req, err := http.NewRequest("POST", url, strings.NewReader(requestBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwtToken)
	req.Header.Set("X-Tenant-ID", tenantID)

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

// SearchProcessDefinition searches for a workflow process definition by ID
func SearchProcessDefinition(serverURL, jwtToken, tenantID, processID string) (string, error) {
	if serverURL == "" {
		return "", fmt.Errorf("server URL is required")
	}
	if tenantID == "" {
		return "", fmt.Errorf("tenant ID is required")
	}
	if processID == "" {
		return "", fmt.Errorf("process ID is required")
	}

	url := strings.TrimSuffix(serverURL, "/") + "/workflow/v1/process/definition?id=" + processID

	// Create request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("X-Tenant-ID", tenantID)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+jwtToken)

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
// DeleteProcess deletes a workflow process by code
// Returns the raw response body as string and any error encountered
func DeleteProcess(serverURL, jwtToken, tenantID, code string) (string, error) {
	if serverURL == "" {
		return "", fmt.Errorf("server URL is required")
	}
	if tenantID == "" {
		return "", fmt.Errorf("tenant ID is required")
	}
	if code == "" {
		return "", fmt.Errorf("process code is required")
	}

	url := strings.TrimSuffix(serverURL, "/") + "/workflow/v1/process?code=" + code

	// Create request
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("X-Tenant-ID", tenantID)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwtToken)

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
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return string(body), nil
}