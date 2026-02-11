package digit

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

// SequenceConfig represents the sequence configuration for ID generation
type SequenceConfig struct {
	Scope   string         `json:"scope"`
	Start   int            `json:"start"`
	Padding PaddingConfig  `json:"padding"`
}

// PaddingConfig represents the padding configuration
type PaddingConfig struct {
	Length int    `json:"length"`
	Char   string `json:"char"`
}

// RandomConfig represents the random configuration for ID generation
type RandomConfig struct {
	Length  int    `json:"length"`
	Charset string `json:"charset"`
}

// IdGenConfig represents the configuration for ID generation template
type IdGenConfig struct {
	Template string         `json:"template"`
	Sequence SequenceConfig `json:"sequence"`
	Random   RandomConfig   `json:"random"`
}

// IdGenTemplateRequest represents the request structure for creating an ID generation template
type IdGenTemplateRequest struct {
	TemplateCode string      `json:"templateCode"`
	Config       IdGenConfig `json:"config"`
}


// SearchIdGenTemplate searches for an ID generation template by templateCode
func SearchIdGenTemplate(serverURL, jwtToken, clientID, tenantID, templateCode string) (string, error) {
	if serverURL == "" {
		return "", fmt.Errorf("serverURL cannot be empty")
	}
	if clientID == "" {
		return "", fmt.Errorf("clientID cannot be empty")
	}
	if tenantID == "" {
		return "", fmt.Errorf("tenantID cannot be empty")
	}
	if templateCode == "" {
		return "", fmt.Errorf("templateCode cannot be empty")
	}

	client := resty.New()

	resp, err := client.R().
		SetHeader("X-Client-ID", clientID).
		SetHeader("X-Tenant-ID", tenantID).
		SetHeader("Authorization", "Bearer "+jwtToken).
		Get(serverURL + "/idgen/v1/template?templateCode=" + templateCode)

	if err != nil {
		return "", fmt.Errorf("failed to fetch ID generation template: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("fetch failed: HTTP %d - %s",
			resp.StatusCode(), string(resp.Body()))
	}

	return string(resp.Body()), nil
}


// CreateIdGenTemplate creates a new ID generation template
// Returns the raw response body as string and any error encountered
func CreateIdGenTemplate(serverURL, jwtToken, clientID, tenantID, templateCode, template, scope string, start, paddingLength int, paddingChar string, randomLength int, randomCharset string) (string, error) {
	// Validate required parameters
	if serverURL == "" {
		return "", fmt.Errorf("serverURL cannot be empty")
	}
	if clientID == "" {
		return "", fmt.Errorf("clientID cannot be empty")
	}
	if templateCode == "" {
		return "", fmt.Errorf("templateCode cannot be empty")
	}
	if template == "" {
		return "", fmt.Errorf("template cannot be empty")
	}

	// Create HTTP client
	client := resty.New()

	// Build the request payload
	templateReq := IdGenTemplateRequest{
		TemplateCode: templateCode,
		Config: IdGenConfig{
			Template: template,
			Sequence: SequenceConfig{
				Scope: scope,
				Start: start,
				Padding: PaddingConfig{
					Length: paddingLength,
					Char:   paddingChar,
				},
			},
			Random: RandomConfig{
				Length:  randomLength,
				Charset: randomCharset,
			},
		},
	}

	// Make the API call
	resp, err := client.R().
		SetHeader("X-Client-ID", clientID).
		SetHeader("X-Tenant-ID", tenantID).
		SetHeader("Authorization", "Bearer "+jwtToken).
		SetHeader("Content-Type", "application/json").
		SetBody(templateReq).
		Post(serverURL + "/idgen/v1/template")

	if err != nil {
		return "", fmt.Errorf("failed to create ID generation template: %w", err)
	}

	// Check for successful response
	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusCreated {
		return "", fmt.Errorf("failed to create ID generation template: HTTP %d - %s", resp.StatusCode(), string(resp.Body()))
	}

	// Return the raw response body as string
	return string(resp.Body()), nil
}

// DeleteIdGenTemplate deletes an ID generation template by templateCode and version
// Returns the raw response body as string and any error encountered
func DeleteIdGenTemplate(serverURL, jwtToken, clientID, tenantID, templateCode, version string) (string, error) {
	// Validate required parameters
	if serverURL == "" {
		return "", fmt.Errorf("serverURL cannot be empty")
	}
	if clientID == "" {
		return "", fmt.Errorf("clientID cannot be empty")
	}
	if tenantID == "" {
		return "", fmt.Errorf("tenantID cannot be empty")
	}
	if templateCode == "" {
		return "", fmt.Errorf("templateCode cannot be empty")
	}
	if version == "" {
		return "", fmt.Errorf("version cannot be empty")
	}

	// Create HTTP client
	client := resty.New()

	// Make the API call
	resp, err := client.R().
		SetHeader("X-Client-ID", clientID).
		SetHeader("X-Tenant-ID", tenantID).
		SetHeader("Authorization", "Bearer "+jwtToken).
		SetQueryParam("templateCode", templateCode).
		SetQueryParam("version", version).
		Delete(serverURL + "/idgen/v1/template")

	if err != nil {
		return "", fmt.Errorf("failed to delete ID generation template: %w", err)
	}

	// Check for successful response
	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusNoContent {
		return "", fmt.Errorf("failed to delete ID generation template: HTTP %d - %s", resp.StatusCode(), string(resp.Body()))
	}

	// Return the raw response body as string
	return string(resp.Body()), nil
}
