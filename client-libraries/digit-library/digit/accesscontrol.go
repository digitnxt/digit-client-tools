package digit

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

//
// ============================================================================
// CONSTRAINT MODEL (JBAC EXTENSION)
// ============================================================================
// This is used by Kong JBAC plugin to evaluate contextual access.
// Example:
//   boundary:tenantId:EQ:KA.BLR
//
type Constraint struct {
	Type     string            `json:"type"`
	Key      string            `json:"key"`
	Operator string            `json:"op"`
	Value    string            `json:"value"`
	Extra    map[string]string `json:"extra,omitempty"`
}

//
// ============================================================================
// RBAC RULE PAYLOAD (RBAC + JBAC)
// ============================================================================
// NOTE:
// Constraints are optional → pure RBAC if empty
//
type RbacRulePayload struct {
	RoleNames   []string     `json:"roleNames"`
	HTTPMethod  string       `json:"httpMethod"`
	Path        string       `json:"path"`
	Effect      string       `json:"effect"`
	Priority    int          `json:"priority"`
	Enabled     bool         `json:"enabled"`
	Description string       `json:"description,omitempty"`
	Constraints []Constraint `json:"constraints,omitempty"`
}

//
// ============================================================================
// CREATE RULE
// ============================================================================
func CreateRbacRule(serverURL, jwtToken, tenantID string, rule RbacRulePayload) (string, error) {
	if serverURL == "" {
		return "", fmt.Errorf("server URL is required")
	}
	if tenantID == "" {
		return "", fmt.Errorf("tenant ID is required")
	}

	url := strings.TrimSuffix(serverURL, "/") + "/access/v3/rbac/rules/"

	body, err := json.Marshal(rule)
	if err != nil {
		return "", fmt.Errorf("failed to marshal rule: %w", err)
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(body)))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwtToken)
	req.Header.Set("X-Tenant-ID", tenantID)

	client := &http.Client{}
	resp, err := client.Do(req)
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

//
// ============================================================================
// LIST RULES
// ============================================================================
func ListRbacRules(serverURL, jwtToken, tenantID, roleName string, page, size int) (string, error) {
	if serverURL == "" {
		return "", fmt.Errorf("server URL is required")
	}
	if tenantID == "" {
		return "", fmt.Errorf("tenant ID is required")
	}

	url := fmt.Sprintf("%s/access/v3/rbac/rules/?page=%d&size=%d",
		strings.TrimSuffix(serverURL, "/"), page, size)

	if roleName != "" {
		url += "&roleName=" + roleName
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+jwtToken)
	req.Header.Set("X-Tenant-ID", tenantID)

	client := &http.Client{}
	resp, err := client.Do(req)
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

//
// ============================================================================
// GET RULE
// ============================================================================
func GetRbacRule(serverURL, jwtToken, tenantID, ruleID string) (string, error) {
	if ruleID == "" {
		return "", fmt.Errorf("rule ID is required")
	}

	url := fmt.Sprintf("%s/access/v3/rbac/rules/%s/",
		strings.TrimSuffix(serverURL, "/"), ruleID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+jwtToken)
	req.Header.Set("X-Tenant-ID", tenantID)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	return string(respBody), nil
}

//
// ============================================================================
// DELETE RULE
// ============================================================================
func DeleteRbacRule(serverURL, jwtToken, tenantID, ruleID string) error {
	if ruleID == "" {
		return fmt.Errorf("rule ID is required")
	}

	url := fmt.Sprintf("%s/access/v3/rbac/rules/%s/",
		strings.TrimSuffix(serverURL, "/"), ruleID)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+jwtToken)
	req.Header.Set("X-Tenant-ID", tenantID)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}
