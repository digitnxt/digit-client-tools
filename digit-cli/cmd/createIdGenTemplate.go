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

// defaultIdGenConfig contains the embedded default IdGen template configuration
const defaultIdGenConfig = `template-code: "DEFAULT_TEMPLATE_CODE"
template: "{ORG}-{DATE:yyyyMMdd}-{SEQ}-{RAND}"
scope: "daily"
start: 1
padding-length: 4
padding-char: "0"
random-length: 2
random-charset: "A-Z0-9"`

// createIdGenTemplateCmd represents the create-idgen-template command
var createIdGenTemplateCmd = &cobra.Command{
	Use:   "create-idgen-template",
	Short: "Create a new ID generation template",
	Long: `Create a new ID generation template with the specified configuration.

This command will create an ID generation template that can be used to generate unique IDs
with custom patterns including sequences, random characters, and date formatting.

Examples:
  # Using individual flags
  digit create-idgen-template --template-code orgId --template "{ORG}-{DATE:yyyyMMdd}-{SEQ}-{RAND}" --scope daily --start 1 --padding-length 4 --padding-char "0" --random-length 2 --random-charset "A-Z0-9"
  
  # Using default configuration with custom template code
  digit create-idgen-template --default --template-code "my-custom-template"
  
  # With server override
  digit create-idgen-template --template-code userId --template "USER-{SEQ}-{RAND}" --scope global --start 100 --server http://localhost:8080`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get flag values
		useDefault, _ := cmd.Flags().GetBool("default")
		templateCode, _ := cmd.Flags().GetString("template-code")
		template, _ := cmd.Flags().GetString("template")
		scope, _ := cmd.Flags().GetString("scope")
		startStr, _ := cmd.Flags().GetString("start")
		paddingLengthStr, _ := cmd.Flags().GetString("padding-length")
		paddingChar, _ := cmd.Flags().GetString("padding-char")
		randomLengthStr, _ := cmd.Flags().GetString("random-length")
		randomCharset, _ := cmd.Flags().GetString("random-charset")
		serverURL, _ := cmd.Flags().GetString("server")
		jwtToken, _ := cmd.Flags().GetString("jwt-token")
		
		// Validate flags - template-code is always required
		if templateCode == "" {
			return fmt.Errorf("--template-code flag is required")
		}
		
		// Check if using default configuration
		if useDefault {
			// Use embedded default configuration and replace the template code
			configContent := strings.Replace(defaultIdGenConfig, "DEFAULT_TEMPLATE_CODE", templateCode, 1)
			
			// Parse the default config (simple key-value parsing)
			lines := strings.Split(configContent, "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line == "" || strings.HasPrefix(line, "#") {
					continue
				}
				parts := strings.SplitN(line, ":", 2)
				if len(parts) != 2 {
					continue
				}
				key := strings.TrimSpace(parts[0])
				value := strings.Trim(strings.TrimSpace(parts[1]), `"`)
				
				switch key {
				case "template":
					template = value
				case "scope":
					scope = value
				case "start":
					startStr = value
				case "padding-length":
					paddingLengthStr = value
				case "padding-char":
					paddingChar = value
				case "random-length":
					randomLengthStr = value
				case "random-charset":
					randomCharset = value
				}
			}
			
			fmt.Printf("Using default IdGen template configuration with template code: %s\n", templateCode)
		}

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

		// Extract client ID and tenant ID from JWT token
		clientID, err := jwt.ExtractClientID(jwtToken)
		if err != nil {
			return fmt.Errorf("failed to extract client ID from JWT token: %w", err)
		}

		tenantID, err := jwt.ExtractTenantID(jwtToken)
		if err != nil {
			return fmt.Errorf("failed to extract tenant ID from JWT token: %w", err)
		}

		// Validate required parameters
		if template == "" {
			return fmt.Errorf("template is required")
		}
		if jwtToken == "" {
			return fmt.Errorf("JWT token is required (set via config or --jwt-token flag)")
		}
		if serverURL == "" {
			return fmt.Errorf("server URL is required (set via config or --server flag)")
		}

		// Parse numeric parameters with defaults
		start := 1
		if startStr != "" {
			var err error
			start, err = strconv.Atoi(startStr)
			if err != nil {
				return fmt.Errorf("invalid start value: %s", startStr)
			}
		}

		paddingLength := 4
		if paddingLengthStr != "" {
			var err error
			paddingLength, err = strconv.Atoi(paddingLengthStr)
			if err != nil {
				return fmt.Errorf("invalid padding-length value: %s", paddingLengthStr)
			}
		}

		randomLength := 2
		if randomLengthStr != "" {
			var err error
			randomLength, err = strconv.Atoi(randomLengthStr)
			if err != nil {
				return fmt.Errorf("invalid random-length value: %s", randomLengthStr)
			}
		}

		// Set defaults for optional parameters
		if scope == "" {
			scope = "daily"
		}
		if paddingChar == "" {
			paddingChar = "0"
		}
		if randomCharset == "" {
			randomCharset = "A-Z0-9"
		}

		// Create the ID generation template
		responseBody, err := digit.CreateIdGenTemplate(serverURL, jwtToken, clientID, tenantID, templateCode, template, scope, start, paddingLength, paddingChar, randomLength, randomCharset)
		if err != nil {
			return fmt.Errorf("failed to create ID generation template: %w", err)
		}

		// Print response
		fmt.Println("ID generation template creation response:")
		if responseBody == "" {
			fmt.Println("ID generation template created successfully (no response body)")
		} else {
			fmt.Println(responseBody)
		}

		return nil
	},
}

var searchIdGenTemplateCmd = &cobra.Command{
	Use:   "search-idgen-template",
	Short: "Search for an existing IDGen Template",
	Long: `Search for an ID generation template using templateCode.

Example:
  digit search-idgen-template --template-code orgId
`,
	RunE: func(cmd *cobra.Command, args []string) error {

		templateCode, _ := cmd.Flags().GetString("template-code")
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

		// Extract metadata from JWT
		clientID, err := jwt.ExtractClientID(jwtToken)
		if err != nil {
			return fmt.Errorf("invalid JWT: %w", err)
		}

		tenantID, err := jwt.ExtractTenantID(jwtToken)
		if err != nil {
			return fmt.Errorf("invalid JWT: %w", err)
		}

		resp, err := digit.SearchIdGenTemplate(serverURL, jwtToken, clientID, tenantID, templateCode)
		if err != nil {
			return err
		}

		fmt.Println("IDGen Template Response:")
		fmt.Println(resp)
		return nil
	},
}


// deleteIdGenTemplateCmd represents the delete-idgen-template command
var deleteIdGenTemplateCmd = &cobra.Command{
	Use:   "delete-idgen-template",
	Short: "Delete an existing ID generation template",
	Long: `Delete an ID generation template by template code and version.

This command will delete an existing ID generation template from the system.

Example:
  digit delete-idgen-template --template-code orgId --version v2
  digit delete-idgen-template --template-code userId --version v1 --server http://localhost:8100`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get flag values
		templateCode, _ := cmd.Flags().GetString("template-code")
		version, _ := cmd.Flags().GetString("version")
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

		// Extract client ID and tenant ID from JWT token
		clientID, err := jwt.ExtractClientID(jwtToken)
		if err != nil {
			return fmt.Errorf("failed to extract client ID from JWT token: %w", err)
		}

		tenantID, err := jwt.ExtractTenantID(jwtToken)
		if err != nil {
			return fmt.Errorf("failed to extract tenant ID from JWT token: %w", err)
		}

		// Validate required parameters
		if templateCode == "" {
			return fmt.Errorf("template-code is required")
		}
		if version == "" {
			return fmt.Errorf("version is required")
		}
		if jwtToken == "" {
			return fmt.Errorf("JWT token is required (set via config or --jwt-token flag)")
		}
		if serverURL == "" {
			return fmt.Errorf("server URL is required (set via config or --server flag)")
		}

		// Delete the ID generation template
		responseBody, err := digit.DeleteIdGenTemplate(serverURL, jwtToken, clientID, tenantID, templateCode, version)
		if err != nil {
			return fmt.Errorf("failed to delete ID generation template: %w", err)
		}

		// Print response
		fmt.Println("ID generation template deletion response:")
		if responseBody == "" {
			fmt.Println("ID generation template deleted successfully")
		} else {
			fmt.Println(responseBody)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(createIdGenTemplateCmd)
	rootCmd.AddCommand(searchIdGenTemplateCmd)
	rootCmd.AddCommand(deleteIdGenTemplateCmd)

// Flags for create-idgen-template command
createIdGenTemplateCmd.Flags().Bool("default", false, "Use default IdGen template configuration (requires --template-code)")
createIdGenTemplateCmd.Flags().StringP("template-code", "", "", "Template code for the ID generation template (required)")
createIdGenTemplateCmd.Flags().StringP("template", "", "", "Template pattern (e.g., '{ORG}-{DATE:yyyyMMdd}-{SEQ}-{RAND}') (required if not using --default)")
createIdGenTemplateCmd.Flags().StringP("scope", "", "daily", "Sequence scope (daily, monthly, yearly, global)")
createIdGenTemplateCmd.Flags().StringP("start", "", "1", "Starting number for sequence")
createIdGenTemplateCmd.Flags().StringP("padding-length", "", "4", "Padding length for sequence numbers")
createIdGenTemplateCmd.Flags().StringP("padding-char", "", "0", "Padding character for sequence numbers")
createIdGenTemplateCmd.Flags().StringP("random-length", "", "2", "Length of random string")
createIdGenTemplateCmd.Flags().StringP("random-charset", "", "A-Z0-9", "Character set for random string")
createIdGenTemplateCmd.Flags().StringP("server", "s", "", "Server URL (overrides config)")
createIdGenTemplateCmd.Flags().StringP("jwt-token", "t", "", "JWT token for authentication (overrides config)")

createIdGenTemplateCmd.MarkFlagRequired("template-code")
// Note: template is validated conditionally in the command logic

// Flags for search-idgen-template command
searchIdGenTemplateCmd.Flags().StringP("template-code", "", "", "Template code (required)")
searchIdGenTemplateCmd.Flags().StringP("server", "s", "", "Server URL")
searchIdGenTemplateCmd.Flags().StringP("jwt-token", "t", "", "JWT Token")

searchIdGenTemplateCmd.MarkFlagRequired("template-code")

// Flags for delete-idgen-template command
deleteIdGenTemplateCmd.Flags().StringP("template-code", "", "", "Template code for the ID generation template (required)")
deleteIdGenTemplateCmd.Flags().StringP("version", "", "", "Version of the template to delete (required)")
deleteIdGenTemplateCmd.Flags().StringP("server", "s", "", "Server URL (overrides config)")
deleteIdGenTemplateCmd.Flags().StringP("jwt-token", "t", "", "JWT token for authentication (overrides config)")

deleteIdGenTemplateCmd.MarkFlagRequired("template-code")
deleteIdGenTemplateCmd.MarkFlagRequired("version")
}
