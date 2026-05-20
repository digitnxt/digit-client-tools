package digit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type ActionInput struct {
	Code      string `json:"code"`
	Label     string `json:"label"`
	NextState string `json:"nextState"`
}

type StateInput struct {
	Code    string        `json:"code"`
	Name    string        `json:"name"`
	Type    string        `json:"type"`
	SLA     int           `json:"sla,omitempty"`
	Actions []ActionInput `json:"actions"`
}

type ProcessDefinitionInput struct {
	Code        string       `json:"code"`
	Name        string       `json:"name"`
	Description string       `json:"description,omitempty"`
	Version     string       `json:"version,omitempty"`
	SLA         int          `json:"sla,omitempty"`
	States      []StateInput `json:"states"`
}

func doRequest(method, url, jwtToken, tenantID string, body []byte) (string, error) {
	var reqBody *bytes.Reader
	if body != nil {
		reqBody = bytes.NewReader(body)
	} else {
		reqBody = bytes.NewReader([]byte{})
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwtToken)
	req.Header.Set("X-Tenant-ID", tenantID)

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody))
	}
	return string(respBody), nil
}

// CreateProcessDefinition creates a complete workflow process definition in a single call.
func CreateProcessDefinition(serverURL, jwtToken, tenantID string, input ProcessDefinitionInput) (string, error) {
	if serverURL == "" {
		return "", fmt.Errorf("server URL is required")
	}
	if tenantID == "" {
		return "", fmt.Errorf("tenant ID is required")
	}

	bodyBytes, err := json.Marshal(input)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	url := strings.TrimSuffix(serverURL, "/") + "/workflow/v3/process/definition"
	return doRequest("POST", url, jwtToken, tenantID, bodyBytes)
}

// SearchProcessDefinition retrieves a workflow process definition by code.
func SearchProcessDefinition(serverURL, jwtToken, tenantID, processCode string) (string, error) {
	if serverURL == "" {
		return "", fmt.Errorf("server URL is required")
	}
	if tenantID == "" {
		return "", fmt.Errorf("tenant ID is required")
	}
	if processCode == "" {
		return "", fmt.Errorf("process code is required")
	}

	url := strings.TrimSuffix(serverURL, "/") + "/workflow/v3/process/definition/" + processCode
	return doRequest("GET", url, jwtToken, tenantID, nil)
}

// DeleteProcess deletes a workflow process definition by code.
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

	url := strings.TrimSuffix(serverURL, "/") + "/workflow/v3/process/definition/" + code
	return doRequest("DELETE", url, jwtToken, tenantID, nil)
}
