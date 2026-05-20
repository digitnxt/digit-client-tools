package digit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type TenantCreateRequest struct {
	Name                 string            `json:"name"`
	Email                string            `json:"email"`
	Password             string            `json:"password,omitempty"`
	Phone                string            `json:"phone,omitempty"`
	Address              string            `json:"address,omitempty"`
	City                 string            `json:"city,omitempty"`
	State                string            `json:"state,omitempty"`
	Pincode              string            `json:"pincode,omitempty"`
	AdditionalAttributes map[string]string `json:"additionalAttributes,omitempty"`
}

// CreateTenant creates a new tenant account via the admin API.
func CreateTenant(serverURL, jwtToken, tenantID string, req TenantCreateRequest) (string, error) {
	if serverURL == "" {
		return "", fmt.Errorf("server URL is required")
	}
	if req.Name == "" {
		return "", fmt.Errorf("name is required")
	}
	if req.Email == "" {
		return "", fmt.Errorf("email is required")
	}

	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	url := strings.TrimSuffix(serverURL, "/") + "/accounts/v3/tenants"

	httpReq, err := http.NewRequest("POST", url, bytes.NewReader(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+jwtToken)
	httpReq.Header.Set("X-Tenant-ID", tenantID)

	resp, err := (&http.Client{}).Do(httpReq)
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

// SearchTenants lists/searches tenant accounts with optional filters.
func SearchTenants(serverURL, jwtToken, tenantID, name, email string, page, size int) (string, error) {
	if serverURL == "" {
		return "", fmt.Errorf("server URL is required")
	}

	url := strings.TrimSuffix(serverURL, "/") + fmt.Sprintf("/accounts/v3/tenants?page=%d&size=%d", page, size)
	if name != "" {
		url += "&name=" + name
	}
	if email != "" {
		url += "&email=" + email
	}

	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+jwtToken)
	httpReq.Header.Set("X-Tenant-ID", tenantID)

	resp, err := (&http.Client{}).Do(httpReq)
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

// DeleteTenant deletes a tenant account by ID.
func DeleteTenant(serverURL, jwtToken, tenantID, id string) (string, error) {
	if serverURL == "" {
		return "", fmt.Errorf("server URL is required")
	}
	if id == "" {
		return "", fmt.Errorf("tenant ID is required")
	}

	url := strings.TrimSuffix(serverURL, "/") + "/accounts/v3/tenants/" + id

	httpReq, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+jwtToken)
	httpReq.Header.Set("X-Tenant-ID", tenantID)

	resp, err := (&http.Client{}).Do(httpReq)
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
