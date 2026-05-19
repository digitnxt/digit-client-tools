package errors

import (
	"fmt"
	"strings"
)

// HandleAPIError checks for 403 access denied errors and provides helpful messages
func HandleAPIError(err error) error {
	if err == nil {
		return nil
	}

	errMsg := err.Error()
	
	// Check for 403 Access Denied errors
	if strings.Contains(errMsg, "API request failed with status 403") ||
		strings.Contains(errMsg, "Access denied") ||
		strings.Contains(errMsg, "authorization service") {
		return fmt.Errorf("%s\n\nIt looks like your authentication token may be invalid or expired.\nPlease run the following command with your credentials:\n\ndigit config set --server http://localhost:8095 --account AMARAVATI --client-id auth-server --client-secret changeme --username test@example.com --password default\n\nNote: Change the username and password according to your actual credentials.\n\nTo get new credentials, please contact your system administrator or check your authentication service.", errMsg)
	}
	
	// Return original error if no specific handling is needed
	return err
}