package digit

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

// SignupTenantRequest is the request body for /account/v3/signup
type SignupTenantRequest struct {
	Tenant SignupTenant `json:"tenant"`
}

// SignupTenant holds name, email, and password for v3 signup
type SignupTenant struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// VerifyAccountRequest is the request body for /account/v3/signup/verify
type VerifyAccountRequest struct {
	RequestID string `json:"requestId"`
	OTP       string `json:"otp"`
}

// SignupAccount calls /account/v3/signup to initiate account creation.
// Returns the raw response body (which includes a requestId for OTP verification).
func SignupAccount(serverURL, name, email, password string) (string, error) {
	if serverURL == "" {
		return "", fmt.Errorf("serverURL cannot be empty")
	}
	if name == "" {
		return "", fmt.Errorf("name cannot be empty")
	}
	if email == "" {
		return "", fmt.Errorf("email cannot be empty")
	}
	if password == "" {
		return "", fmt.Errorf("password cannot be empty")
	}

	req := SignupTenantRequest{
		Tenant: SignupTenant{
			Name:     name,
			Email:    email,
			Password: password,
		},
	}

	client := resty.New()

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(req).
		Post(serverURL + "/account/v3/signup")

	if err != nil {
		return "", fmt.Errorf("failed to make API request: %w", err)
	}

	return string(resp.Body()), nil
}

// VerifyAccount calls /account/v3/signup/verify to confirm the OTP sent after signup.
func VerifyAccount(serverURL, requestID, otp string) (string, error) {
	if serverURL == "" {
		return "", fmt.Errorf("serverURL cannot be empty")
	}
	if requestID == "" {
		return "", fmt.Errorf("requestId cannot be empty")
	}
	if otp == "" {
		return "", fmt.Errorf("otp cannot be empty")
	}

	req := VerifyAccountRequest{
		RequestID: requestID,
		OTP:       otp,
	}

	client := resty.New()

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(req).
		Post(serverURL + "/account/v3/signup/verify")

	if err != nil {
		return "", fmt.Errorf("failed to make API request: %w", err)
	}

	return string(resp.Body()), nil
}
