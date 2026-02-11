package digit

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-resty/resty/v2"
)

// DocumentCategoryRequest represents the request structure for creating a document category
type DocumentCategoryRequest struct {
	Type           string   `json:"type"`
	Code           string   `json:"code"`
	AllowedFormats []string `json:"allowedFormats"`
	MinSize        string   `json:"minSize"`
	MaxSize        string   `json:"maxSize"`
	IsSensitive    bool     `json:"isSensitive"`
	Description    string   `json:"description"`
	IsActive       bool     `json:"isActive"`
}

// CreateDocumentCategory creates a new document category in filestore
// Returns the raw response body as string and any error encountered
func CreateDocumentCategory(serverURL, jwtToken, tenantID, categoryType, code string, allowedFormats []string, minSize, maxSize int, isSensitive, isActive bool, description string) (string, error) {
	// Validate required parameters
	if serverURL == "" {
		return "", fmt.Errorf("serverURL cannot be empty")
	}
	if tenantID == "" {
		return "", fmt.Errorf("tenantID cannot be empty")
	}
	if categoryType == "" {
		return "", fmt.Errorf("type cannot be empty")
	}
	if code == "" {
		return "", fmt.Errorf("code cannot be empty")
	}
	if len(allowedFormats) == 0 {
		return "", fmt.Errorf("allowedFormats cannot be empty")
	}

	// Create HTTP client
	client := resty.New()

	// Build the request payload
	categoryReq := DocumentCategoryRequest{
		Type:           categoryType,
		Code:           code,
		AllowedFormats: allowedFormats,
		MinSize:        strconv.Itoa(minSize),
		MaxSize:        strconv.Itoa(maxSize),
		IsSensitive:    isSensitive,
		Description:    description,
		IsActive:       isActive,
	}

	// Make the API call
	resp, err := client.R().
		SetHeader("X-Tenant-ID", tenantID).
		SetHeader("Authorization", "Bearer "+jwtToken).
		SetHeader("Content-Type", "application/json").
		SetBody(categoryReq).
		Post(serverURL + "/filestore/v1/files/document-categories")

	if err != nil {
		return "", fmt.Errorf("failed to create document category: %w", err)
	}

	// Check for successful response
	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusCreated {
		return "", fmt.Errorf("failed to create document category: HTTP %d - %s", resp.StatusCode(), string(resp.Body()))
	}

	// Return the raw response body as string
	return string(resp.Body()), nil
}
