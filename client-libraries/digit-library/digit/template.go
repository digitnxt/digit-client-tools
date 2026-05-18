package digit

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

// TemplateRequest represents the request structure for creating a notification template
type TemplateRequest struct {
	TemplateID string `json:"templateId"`
	Version    string `json:"version"`
	Type       string `json:"type"`
	Subject    string `json:"subject"`
	Content    string `json:"content"`
	IsHTML     bool   `json:"isHTML"`
}

// CreateTemplate creates a new notification template
// Returns the raw response body as string and any error encountered
func CreateTemplate(serverURL, jwtToken, tenantID, templateID, version, templateType, subject, content string, isHTML bool) (string, error) {
	// Validate required parameters
	if serverURL == "" {
		return "", fmt.Errorf("serverURL cannot be empty")
	}
	if tenantID == "" {
		return "", fmt.Errorf("tenantID cannot be empty")
	}
	if templateID == "" {
		return "", fmt.Errorf("templateID cannot be empty")
	}
	if version == "" {
		return "", fmt.Errorf("version cannot be empty")
	}
	if templateType == "" {
		return "", fmt.Errorf("templateType cannot be empty")
	}
	if subject == "" {
		return "", fmt.Errorf("subject cannot be empty")
	}
	if content == "" {
		return "", fmt.Errorf("content cannot be empty")
	}

	// Create the request payload
	templateReq := TemplateRequest{
		TemplateID: templateID,
		Version:    version,
		Type:       templateType,
		Subject:    subject,
		Content:    content,
		IsHTML:     isHTML,
	}

	// Create HTTP client
	client := resty.New()

	// Make the API request
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+jwtToken).
		SetHeader("X-Tenant-ID", tenantID).
		SetBody(templateReq).
		Post(serverURL + "/notification/v1/template")

	if err != nil {
		return "", fmt.Errorf("failed to make API request: %w", err)
	}

	// Return the raw response body as string
	return string(resp.Body()), nil
}

// SearchNotificationTemplate searches for notification templates by template ID
// Returns the raw response body as string and any error encountered
func SearchNotificationTemplate(serverURL, jwtToken, tenantID, templateID string) (string, error) {
	// Validate required parameters
	if serverURL == "" {
		return "", fmt.Errorf("serverURL cannot be empty")
	}
	if tenantID == "" {
		return "", fmt.Errorf("tenantID cannot be empty")
	}
	if templateID == "" {
		return "", fmt.Errorf("templateID cannot be empty")
	}

	// Create HTTP client
	client := resty.New()

	// Make the API request
	resp, err := client.R().
		SetHeader("X-Tenant-ID", tenantID).
		SetHeader("Authorization", "Bearer "+jwtToken).
		SetQueryParam("templateId", templateID).
		Get(serverURL + "/notification/v1/template")

	if err != nil {
		return "", fmt.Errorf("failed to make API request: %w", err)
	}

	// Return the raw response body as string
	return string(resp.Body()), nil
}
// DeleteNotificationTemplate deletes a notification template by templateId and version
// Returns the raw response body as string and any error encountered
func DeleteNotificationTemplate(serverURL, jwtToken, tenantID, templateID, version string) (string, error) {
	// Validate required parameters
	if serverURL == "" {
		return "", fmt.Errorf("serverURL cannot be empty")
	}
	if tenantID == "" {
		return "", fmt.Errorf("tenantID cannot be empty")
	}
	if templateID == "" {
		return "", fmt.Errorf("templateID cannot be empty")
	}
	if version == "" {
		return "", fmt.Errorf("version cannot be empty")
	}

	// Create HTTP client
	client := resty.New()

	// Make the API request
	resp, err := client.R().
		SetHeader("X-Tenant-ID", tenantID).
		SetHeader("Authorization", "Bearer "+jwtToken).
		SetQueryParam("templateId", templateID).
		SetQueryParam("version", version).
		Delete(serverURL + "/notification/v1/template")

	if err != nil {
		return "", fmt.Errorf("failed to delete notification template: %w", err)
	}

	// Check for successful response
	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusNoContent {
		return "", fmt.Errorf("failed to delete notification template: HTTP %d - %s", resp.StatusCode(), string(resp.Body()))
	}

	// Return the raw response body as string
	return string(resp.Body()), nil
}