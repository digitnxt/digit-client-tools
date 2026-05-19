package cmd

import (
	"fmt"
	"strings"
	"unicode"

	"digit-cli/pkg/config"
	"digit-cli/pkg/jwt"
	"github.com/digitnxt/digit-client-tools/client-libraries/digit-library/digit"
	"github.com/spf13/cobra"
)

// parseSizeFlag normalises a size string to include a unit.
// Plain numbers get "B" appended; values already ending in B/KB/MB/GB are passed through.
// Examples: "1024" → "1024B", "1MB" → "1MB", "512kb" → "512KB"
func parseSizeFlag(s string) string {
	s = strings.TrimSpace(s)
	upper := strings.ToUpper(s)
	for _, unit := range []string{"GB", "MB", "KB", "B"} {
		if strings.HasSuffix(upper, unit) {
			// normalise unit to uppercase
			num := s[:len(s)-len(unit)]
			return num + unit
		}
	}
	// no unit — check it's all digits and append B
	allDigits := true
	for _, c := range s {
		if !unicode.IsDigit(c) {
			allDigits = false
			break
		}
	}
	if allDigits {
		return s + "B"
	}
	return s
}

// createDocumentCategoryCmd represents the create-document-category command
var createDocumentCategoryCmd = &cobra.Command{
	Use:   "create-filestore-document-category",
	Short: "Create a new document category in filestore",
	Long: `Create a new document category in filestore with the specified configuration.

This command will create a document category that defines allowed file formats,
size limits, and other properties for document uploads.

Example:
  digit create-filestore-document-category --type Identity --code AADHAR --allowed-formats "pdf,jpg,jpeg,xsm" --min-size 1KB --max-size 1MB --sensitive --active
  digit create-filestore-document-category --type Identity --code AADHAR --allowed-formats "pdf,jpg,jpeg,xsm" --min-size 1024MB --max-size 2048MB --sensitive --active
  digit create-filestore-document-category --type Certificate --code BIRTH_CERT --allowed-formats "pdf,jpg" --min-size 512KB --max-size 2MB --description "Birth certificate documents" --server http://localhost:8081`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get flag values
		categoryType, _ := cmd.Flags().GetString("type")
		code, _ := cmd.Flags().GetString("code")
		allowedFormatsStr, _ := cmd.Flags().GetString("allowed-formats")
		minSizeStr, _ := cmd.Flags().GetString("min-size")
		maxSizeStr, _ := cmd.Flags().GetString("max-size")
		isSensitive, _ := cmd.Flags().GetBool("sensitive")
		isActive, _ := cmd.Flags().GetBool("active")
		description, _ := cmd.Flags().GetString("description")
		serverURL, _ := cmd.Flags().GetString("server")
		jwtToken, _ := cmd.Flags().GetString("jwt-token")

		// Get server URL and JWT token from config if not provided
		if serverURL == "" || jwtToken == "" {
			cfg, _ := config.Load()
			if serverURL == "" {
				serverURL = cfg.Server
			}
			if jwtToken == "" {
				jwtToken = cfg.JWTToken
			}
		}

		// Extract tenant ID from JWT token
		tenantID, err := jwt.ExtractTenantID(jwtToken)
		if err != nil {
			return fmt.Errorf("failed to extract tenant ID from JWT token: %w", err)
		}

		// Validate required parameters
		if categoryType == "" {
			return fmt.Errorf("type is required")
		}
		if code == "" {
			return fmt.Errorf("code is required")
		}
		if allowedFormatsStr == "" {
			return fmt.Errorf("allowed-formats is required")
		}
		if jwtToken == "" {
			return fmt.Errorf("JWT token is required (set via config or --jwt-token flag)")
		}
		if serverURL == "" {
			return fmt.Errorf("server URL is required (set via config or --server flag)")
		}

		// Parse allowed formats
		allowedFormats := strings.Split(allowedFormatsStr, ",")
		for i, format := range allowedFormats {
			allowedFormats[i] = strings.TrimSpace(format)
		}

		minSize := parseSizeFlag(minSizeStr)
		maxSize := parseSizeFlag(maxSizeStr)

		// Create the document category
		responseBody, err := digit.CreateDocumentCategory(serverURL, jwtToken, tenantID, categoryType, code, allowedFormats, minSize, maxSize, isSensitive, isActive, description)
		if err != nil {
			return fmt.Errorf("failed to create document category: %w", err)
		}

		// Print response
		fmt.Println("Document category creation response:")
		if responseBody == "" {
			fmt.Println("Document category created successfully (no response body)")
		} else {
			fmt.Println(responseBody)
		}

		return nil
	},
}

// deleteDocumentCategoryCmd represents the delete-document-category command
var deleteDocumentCategoryCmd = &cobra.Command{
	Use:   "delete-filestore-document-category",
	Short: "Delete a document category by code",
	Long: `Delete a document category from filestore by its code.

Example:
  digit delete-filestore-document-category --code AADHAR
  digit delete-filestore-document-category --code AADHAR --server http://localhost:8102`,
	RunE: func(cmd *cobra.Command, args []string) error {
		code, _ := cmd.Flags().GetString("code")
		serverURL, _ := cmd.Flags().GetString("server")
		jwtToken, _ := cmd.Flags().GetString("jwt-token")

		if serverURL == "" || jwtToken == "" {
			cfg, _ := config.Load()
			if serverURL == "" {
				serverURL = cfg.Server
			}
			if jwtToken == "" {
				jwtToken = cfg.JWTToken
			}
		}

		if serverURL == "" {
			return fmt.Errorf("server URL is required (set via config or --server flag)")
		}
		if jwtToken == "" {
			return fmt.Errorf("JWT token is required (set via config or --jwt-token flag)")
		}

		tenantID, err := jwt.ExtractTenantID(jwtToken)
		if err != nil {
			return fmt.Errorf("failed to extract tenant ID from JWT token: %w", err)
		}

		responseBody, err := digit.DeleteDocumentCategory(serverURL, jwtToken, tenantID, code)
		if err != nil {
			return fmt.Errorf("failed to delete document category: %w", err)
		}

		if responseBody == "" {
			fmt.Printf("Document category '%s' deleted successfully\n", code)
		} else {
			fmt.Println(responseBody)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(createDocumentCategoryCmd)

	// Flags for create-document-category command
	createDocumentCategoryCmd.Flags().StringP("type", "", "", "Category type (e.g., Identity, Certificate) (required)")
	createDocumentCategoryCmd.Flags().StringP("code", "", "", "Category code (e.g., AADHAR, BIRTH_CERT) (required)")
	createDocumentCategoryCmd.Flags().StringP("allowed-formats", "", "", "Comma-separated list of allowed file formats (e.g., 'pdf,jpg,jpeg') (required)")
	createDocumentCategoryCmd.Flags().StringP("min-size", "", "1KB", "Minimum file size (e.g. 1024, 1KB, 1MB)")
	createDocumentCategoryCmd.Flags().StringP("max-size", "", "1MB", "Maximum file size (e.g. 512KB, 10MB, 1GB)")
	createDocumentCategoryCmd.Flags().StringP("sensitive", "", "false", "Whether the document category is sensitive")
	createDocumentCategoryCmd.Flags().StringP("active", "", "true", "Whether the document category is active")
	createDocumentCategoryCmd.Flags().StringP("description", "", "", "Description of the document category")
	createDocumentCategoryCmd.Flags().StringP("server", "s", "", "Server URL (overrides config)")
	createDocumentCategoryCmd.Flags().StringP("jwt-token", "t", "", "JWT token for authentication (overrides config)")

	// Mark required flags
	createDocumentCategoryCmd.MarkFlagRequired("type")
	createDocumentCategoryCmd.MarkFlagRequired("code")
	createDocumentCategoryCmd.MarkFlagRequired("allowed-formats")

	rootCmd.AddCommand(deleteDocumentCategoryCmd)
	deleteDocumentCategoryCmd.Flags().StringP("code", "", "", "Category code to delete (required)")
	deleteDocumentCategoryCmd.Flags().StringP("server", "s", "", "Server URL (overrides config)")
	deleteDocumentCategoryCmd.Flags().StringP("jwt-token", "t", "", "JWT token for authentication (overrides config)")
	deleteDocumentCategoryCmd.MarkFlagRequired("code")
}
