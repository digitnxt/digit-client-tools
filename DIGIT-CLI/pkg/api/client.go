package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"digit-cli/pkg/auth"
	"digit-cli/pkg/config"
)

// Client represents an API client with automatic token refresh
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient creates a new API client
func NewClient() (*Client, error) {
	serverURL, err := config.GetServerURL()
	if err != nil {
		return nil, fmt.Errorf("failed to get server URL: %w", err)
	}

	if serverURL == "" {
		return nil, fmt.Errorf("no server URL configured. Please run 'digit config set' to authenticate")
	}

	return &Client{
		BaseURL:    serverURL,
		HTTPClient: &http.Client{},
	}, nil
}

// Get performs a GET request with automatic token refresh
func (c *Client) Get(endpoint string) (*http.Response, error) {
	return c.makeRequest("GET", endpoint, nil)
}

// Post performs a POST request with automatic token refresh
func (c *Client) Post(endpoint string, body interface{}) (*http.Response, error) {
	return c.makeRequest("POST", endpoint, body)
}

// Put performs a PUT request with automatic token refresh
func (c *Client) Put(endpoint string, body interface{}) (*http.Response, error) {
	return c.makeRequest("PUT", endpoint, body)
}

// Delete performs a DELETE request with automatic token refresh
func (c *Client) Delete(endpoint string) (*http.Response, error) {
	return c.makeRequest("DELETE", endpoint, nil)
}

// makeRequest creates and executes an HTTP request with automatic token refresh
func (c *Client) makeRequest(method, endpoint string, body interface{}) (*http.Response, error) {
	// Prepare request body
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	// Create request
	url := c.BaseURL + endpoint
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set content type for POST/PUT requests
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Get valid authorization header (with automatic refresh)
	authHeader, err := auth.GetAuthorizationHeader()
	if err != nil {
		return nil, fmt.Errorf("failed to get authorization header: %w", err)
	}
	req.Header.Set("Authorization", authHeader)

	// Execute request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	return resp, nil
}

// GetJSON performs a GET request and unmarshals the response into the provided interface
func (c *Client) GetJSON(endpoint string, result interface{}) error {
	resp, err := c.Get(endpoint)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return json.NewDecoder(resp.Body).Decode(result)
}

// PostJSON performs a POST request and unmarshals the response into the provided interface
func (c *Client) PostJSON(endpoint string, requestBody interface{}, result interface{}) error {
	resp, err := c.Post(endpoint, requestBody)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	if result != nil {
		return json.NewDecoder(resp.Body).Decode(result)
	}

	return nil
}
