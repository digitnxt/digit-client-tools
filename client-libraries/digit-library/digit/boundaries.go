package digit

import (
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
)

func boundaryClient(serverURL, jwtToken, tenantID string) *resty.Request {
	return resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", jwtToken)).
		SetHeader("X-Tenant-ID", tenantID)
}

// CreateBoundaries creates boundaries via POST /boundary/v3/boundaries
func CreateBoundaries(serverURL, jwtToken, tenantID, clientID string, boundaryData []map[string]interface{}) (string, error) {
	resp, err := boundaryClient(serverURL, jwtToken, tenantID).
		SetHeader("X-User-ID", clientID).
		SetBody(map[string]interface{}{"boundary": boundaryData}).
		Post(serverURL + "/boundary/v3/boundaries")
	if err != nil {
		return "", fmt.Errorf("failed to make API request: %w", err)
	}
	return string(resp.Body()), nil
}

// CreateBoundaryHierarchy creates a boundary hierarchy via POST /boundary/v3/hierarchy
func CreateBoundaryHierarchy(serverURL, jwtToken, tenantID, clientID string, boundaryHierarchy map[string]interface{}) (string, error) {
	resp, err := boundaryClient(serverURL, jwtToken, tenantID).
		SetHeader("X-User-ID", clientID).
		SetBody(map[string]interface{}{"hierarchy": boundaryHierarchy}).
		Post(serverURL + "/boundary/v3/hierarchy")
	if err != nil {
		return "", fmt.Errorf("failed to make API request: %w", err)
	}
	return string(resp.Body()), nil
}

// SearchBoundaryHierarchy searches boundary hierarchy via GET /boundary/v3/hierarchy?hierarchyType=...
func SearchBoundaryHierarchy(serverURL, jwtToken, tenantID, clientID, hierarchyType string) (string, error) {
	resp, err := boundaryClient(serverURL, jwtToken, tenantID).
		SetHeader("X-User-ID", clientID).
		SetQueryParam("hierarchyType", hierarchyType).
		Get(serverURL + "/boundary/v3/hierarchy")
	if err != nil {
		return "", fmt.Errorf("failed to make API request: %w", err)
	}
	return string(resp.Body()), nil
}

// CreateBoundaryRelationships creates boundary relationships via POST /boundary/v3/relationship
func CreateBoundaryRelationships(serverURL, jwtToken, tenantID, clientID string, boundaryRelationship map[string]interface{}) (string, error) {
	resp, err := boundaryClient(serverURL, jwtToken, tenantID).
		SetHeader("X-User-ID", clientID).
		SetBody(map[string]interface{}{"relationship": boundaryRelationship}).
		Post(serverURL + "/boundary/v3/relationship")
	if err != nil {
		return "", fmt.Errorf("failed to make API request: %w", err)
	}
	return string(resp.Body()), nil
}

// SearchBoundaryRelationships searches boundary relationships via GET /boundary/v3/relationship
func SearchBoundaryRelationships(serverURL, jwtToken, tenantID, clientID, hierarchyType, boundaryType, codes string, includeChildren bool) (string, error) {
	req := boundaryClient(serverURL, jwtToken, tenantID).
		SetHeader("X-User-ID", clientID)

	if hierarchyType != "" {
		req.SetQueryParam("hierarchyType", hierarchyType)
	}
	if boundaryType != "" {
		req.SetQueryParam("boundaryType", boundaryType)
	}
	if codes != "" {
		req.SetQueryParam("codes", strings.Join(strings.Split(codes, ","), ","))
	}
	if includeChildren {
		req.SetQueryParam("includeChildren", "true")
	}

	resp, err := req.Get(serverURL + "/boundary/v3/relationship")
	if err != nil {
		return "", fmt.Errorf("failed to make API request: %w", err)
	}
	return string(resp.Body()), nil
}
