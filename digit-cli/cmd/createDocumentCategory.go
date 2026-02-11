package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"digit-cli/pkg/config"
	"digit-cli/pkg/jwt"
	"github.com/digitnxt/digit3/code/libraries/digit-library/digit"
	"github.com/spf13/cobra"
)

// createDocumentCategoryCmd represents the create-document-category command
var createDocumentCategoryCmd = &cobra.Command{
	Use:   "create-document-category",
	Short: "Create a new document category in filestore",
	Long: `Create a new document category in filestore with the specified configuration.

This command will create a document category that defines allowed file formats,
size limits, and other properties for document uploads.

Example:
  digit create-document-category --type Identity --code AADHAR --allowed-formats "pdf,jpg,jpeg,xsm" --min-size 1024 --max-size 1024000 --sensitive --active
  digit create-document-category --type Certificate --code BIRTH_CERT --allowed-formats "pdf,jpg" --min-size 512 --max-size 2048000 --description "Birth certificate documents" --server http://localhost:8081`,
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

		// Parse numeric parameters with defaults
		minSize := 1024
		if minSizeStr != "" {
			var err error
			minSize, err = strconv.Atoi(minSizeStr)
			if err != nil {
				return fmt.Errorf("invalid min-size value: %s", minSizeStr)
			}
		}

		maxSize := 1024000
		if maxSizeStr != "" {
			var err error
			maxSize, err = strconv.Atoi(maxSizeStr)
			if err != nil {
				return fmt.Errorf("invalid max-size value: %s", maxSizeStr)
			}
		}

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

func init() {
	rootCmd.AddCommand(createDocumentCategoryCmd)

	// Flags for create-document-category command
	createDocumentCategoryCmd.Flags().StringP("type", "", "", "Category type (e.g., Identity, Certificate) (required)")
	createDocumentCategoryCmd.Flags().StringP("code", "", "", "Category code (e.g., AADHAR, BIRTH_CERT) (required)")
	createDocumentCategoryCmd.Flags().StringP("allowed-formats", "", "", "Comma-separated list of allowed file formats (e.g., 'pdf,jpg,jpeg') (required)")
	createDocumentCategoryCmd.Flags().StringP("min-size", "", "1024", "Minimum file size in bytes")
	createDocumentCategoryCmd.Flags().StringP("max-size", "", "1024000", "Maximum file size in bytes")
	createDocumentCategoryCmd.Flags().StringP("sensitive", "", "false", "Whether the document category is sensitive")
	createDocumentCategoryCmd.Flags().StringP("active", "", "true", "Whether the document category is active")
	createDocumentCategoryCmd.Flags().StringP("description", "", "", "Description of the document category")
	createDocumentCategoryCmd.Flags().StringP("server", "s", "", "Server URL (overrides config)")
	createDocumentCategoryCmd.Flags().StringP("jwt-token", "t", "", "JWT token for authentication (overrides config)")

	// Mark required flags
	createDocumentCategoryCmd.MarkFlagRequired("type")
	createDocumentCategoryCmd.MarkFlagRequired("code")
	createDocumentCategoryCmd.MarkFlagRequired("allowed-formats")
}
