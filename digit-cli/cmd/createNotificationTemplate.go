package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"digit-cli/pkg/config"
	"digit-cli/pkg/jwt"
	"github.com/digitnxt/digit3/code/libraries/digit-library/digit"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// defaultNotificationTemplateYAML contains the embedded default notification template configuration
const defaultNotificationTemplateYAML = `template-id: "DEFAULT_TEMPLATE_ID"
version: "1.0.0"
type: "EMAIL"
subject: "Welcome to DIGIT Services"
content: |
  <html>
  <body>
    <h1>Welcome to DIGIT Services!</h1>
    <p>Dear User,</p>
    <p>Thank you for registering with DIGIT Services. Your account has been successfully created.</p>
    <p>Account Details:</p>
    <ul>
      <li>Username: [USERNAME]</li>
      <li>Email: [EMAIL]</li>
      <li>Registration Date: [DATE]</li>
    </ul>
    <p>If you have any questions, please contact our support team.</p>
    <p>Best regards,<br>DIGIT Team</p>
  </body>
  </html>
html: true`

// TemplateConfig represents the YAML configuration for template creation
type TemplateConfig struct {
	TemplateID  string `yaml:"template-id"`
	Version     string `yaml:"version"`
	Type        string `yaml:"type"`
	Subject     string `yaml:"subject"`
	Content     string `yaml:"content"`
	ContentFile string `yaml:"content-file"`
	HTML        bool   `yaml:"html"`
	ServerURL   string `yaml:"server"`
	JWTToken    string `yaml:"jwt-token"`
}

// createNotificationTemplateCmd represents the create-notification-template command
var createNotificationTemplateCmd = &cobra.Command{
	Use:   "create-notification-template",
	Short: "Create a new notification template",
	Long: `Create a new notification template with the specified parameters.
	
Examples:
  # Using direct content
  digit create-notification-template --template-id "my-template" --version "1.0.0" --type "EMAIL" --subject "Test Subject" --content "Test Content"
  
  # Using content from file
  digit create-notification-template --template-id "my-template" --version "1.0.0" --type "EMAIL" --subject "Test Subject" --content-file "./template.html" --html=true
  
  # Using YAML configuration file
  digit create-notification-template --file template-config.yaml
  
  # Using default configuration with custom template ID
  digit create-notification-template --default --template-id "my-custom-template"
  
  # With server override
  digit create-notification-template --template-id "my-template" --version "1.0.0" --type "SMS" --subject "Test Subject" --content "Test Content" --server http://localhost:8081`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get flag values
		filePath, _ := cmd.Flags().GetString("file")
		useDefault, _ := cmd.Flags().GetBool("default")
		templateID, _ := cmd.Flags().GetString("template-id")
		version, _ := cmd.Flags().GetString("version")
		templateType, _ := cmd.Flags().GetString("type")
		subject, _ := cmd.Flags().GetString("subject")
		content, _ := cmd.Flags().GetString("content")
		contentFile, _ := cmd.Flags().GetString("content-file")
		isHTML, _ := cmd.Flags().GetBool("html")
		serverURL, _ := cmd.Flags().GetString("server")
		jwtToken, _ := cmd.Flags().GetString("jwt-token")
		
		// Validate flags - either file, default, or individual flags must be specified
		if useDefault && filePath != "" {
			return fmt.Errorf("cannot use both --file and --default flags together")
		}
		if useDefault && templateID == "" {
			return fmt.Errorf("--template-id flag is required when using --default")
		}
		
		// Check if using default configuration
		if useDefault {
			// Use embedded default configuration and replace the template ID
			yamlContent := strings.Replace(defaultNotificationTemplateYAML, "DEFAULT_TEMPLATE_ID", templateID, 1)
			yamlData := []byte(yamlContent)
			
			var templateConfig TemplateConfig
			if err := yaml.Unmarshal(yamlData, &templateConfig); err != nil {
				return fmt.Errorf("failed to parse default template configuration: %w", err)
			}
			
			// Use values from default configuration
			version = templateConfig.Version
			templateType = templateConfig.Type
			subject = templateConfig.Subject
			content = templateConfig.Content
			contentFile = templateConfig.ContentFile
			isHTML = templateConfig.HTML
			if templateConfig.ServerURL != "" {
				serverURL = templateConfig.ServerURL
			}
			if templateConfig.JWTToken != "" {
				jwtToken = templateConfig.JWTToken
			}
			
			fmt.Printf("Using default template configuration with template ID: %s\n", templateID)
		} else if filePath != "" {
			// Read and parse YAML file
			yamlData, err := os.ReadFile(filePath)
			if err != nil {
				return fmt.Errorf("failed to read YAML file: %w", err)
			}
			
			var templateConfig TemplateConfig
			if err := yaml.Unmarshal(yamlData, &templateConfig); err != nil {
				return fmt.Errorf("failed to parse YAML file: %w", err)
			}
			
			// Override with YAML values
			templateID = templateConfig.TemplateID
			version = templateConfig.Version
			templateType = templateConfig.Type
			subject = templateConfig.Subject
			content = templateConfig.Content
			contentFile = templateConfig.ContentFile
			isHTML = templateConfig.HTML
			if templateConfig.ServerURL != "" {
				serverURL = templateConfig.ServerURL
			}
			if templateConfig.JWTToken != "" {
				jwtToken = templateConfig.JWTToken
			}
		} else {
			// Validate required flags when not using YAML file or default
			if templateID == "" {
				return fmt.Errorf("--template-id flag is required")
			}
			if version == "" {
				return fmt.Errorf("--version flag is required")
			}
			if templateType == "" {
				return fmt.Errorf("--type flag is required")
			}
			if subject == "" {
				return fmt.Errorf("--subject flag is required")
			}
			if content == "" && contentFile == "" {
				return fmt.Errorf("either --content or --content-file flag is required")
			}
			if content != "" && contentFile != "" {
				return fmt.Errorf("cannot use both --content and --content-file flags together")
			}
		}
		
		// Validate required fields from YAML, default, or flags
		if templateID == "" {
			return fmt.Errorf("template-id is required")
		}
		if version == "" {
			return fmt.Errorf("version is required")
		}
		if templateType == "" {
			return fmt.Errorf("type is required")
		}
		if subject == "" {
			return fmt.Errorf("subject is required")
		}
		if content == "" && contentFile == "" {
			return fmt.Errorf("either content or content-file is required")
		}
		if content != "" && contentFile != "" {
			return fmt.Errorf("cannot use both content and content-file together")
		}
		
		// Read content from file if content-file is provided
		if contentFile != "" {
			fileContent, err := os.ReadFile(contentFile)
			if err != nil {
				return fmt.Errorf("failed to read content file: %w", err)
			}
			content = string(fileContent)
		}
		
		// Get server URL and JWT token from config if not provided
		if serverURL == "" || jwtToken == "" {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}
			
			if serverURL == "" {
				serverURL = cfg.GetServer()
				if serverURL == "" {
					return fmt.Errorf("server URL not configured. Use 'digit config set --server <url>' or provide --server flag")
				}
			}
			
			if jwtToken == "" {
				jwtToken = cfg.GetJWTToken()
				if jwtToken == "" {
					return fmt.Errorf("JWT token not configured. Use 'digit config set --jwt-token <token>' or provide --jwt-token flag")
				}
			}
		}

		// Extract tenant ID from JWT token
		tenantID, err := jwt.ExtractTenantID(jwtToken)
		if err != nil {
			return fmt.Errorf("failed to extract tenant ID from JWT token: %w", err)
		}
		
		// Call the digit library to create template
		responseBody, err := digit.CreateTemplate(serverURL, jwtToken, tenantID, templateID, version, templateType, subject, content, isHTML)
		if err != nil {
			return fmt.Errorf("failed to create template: %w", err)
		}
		
		// Print response body
		fmt.Printf("Template creation response:\n")
		
		// Try to pretty print JSON response
		var jsonResponse interface{}
		if err := json.Unmarshal([]byte(responseBody), &jsonResponse); err == nil {
			prettyJSON, err := json.MarshalIndent(jsonResponse, "", "  ")
			if err == nil {
				fmt.Println(string(prettyJSON))
			} else {
				fmt.Println(responseBody)
			}
		} else {
			fmt.Println(responseBody)
		}
		
		return nil
	},
}

// searchNotificationTemplateCmd represents the search-notification-template command
var searchNotificationTemplateCmd = &cobra.Command{
	Use:   "search-notification-template",
	Short: "Search for notification templates by template ID",
	Long: `Search for notification templates by template ID.
	
Examples:
  digit search-notification-template --template-id "user-notify"
  digit search-notification-template --template-id "user-notify" --server http://localhost:8091`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get flag values
		templateID, _ := cmd.Flags().GetString("template-id")
		serverURL, _ := cmd.Flags().GetString("server")
		jwtToken, _ := cmd.Flags().GetString("jwt-token")
		
		// Validate required flags
		if templateID == "" {
			return fmt.Errorf("--template-id flag is required")
		}
		
		// Get server URL and JWT token from config if not provided
		if serverURL == "" || jwtToken == "" {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}
			
			if serverURL == "" {
				serverURL = cfg.GetServer()
				if serverURL == "" {
					return fmt.Errorf("server URL not configured. Use 'digit config set --server <url>' or provide --server flag")
				}
			}
			
			if jwtToken == "" {
				jwtToken = cfg.GetJWTToken()
				if jwtToken == "" {
					return fmt.Errorf("JWT token not configured. Use 'digit config set --jwt-token <token>' or provide --jwt-token flag")
				}
			}
		}

		// Extract tenant ID from JWT token
		tenantID, err := jwt.ExtractTenantID(jwtToken)
		if err != nil {
			return fmt.Errorf("failed to extract tenant ID from JWT token: %w", err)
		}
		
		// Call the digit library to search template
		responseBody, err := digit.SearchNotificationTemplate(serverURL, jwtToken, tenantID, templateID)
		if err != nil {
			return fmt.Errorf("failed to search notification template: %w", err)
		}
		
		// Print response body
		fmt.Printf("Search notification template response:\n")
		
		// Try to pretty print JSON response
		var jsonResponse interface{}
		if err := json.Unmarshal([]byte(responseBody), &jsonResponse); err == nil {
			prettyJSON, err := json.MarshalIndent(jsonResponse, "", "  ")
			if err == nil {
				fmt.Println(string(prettyJSON))
			} else {
				fmt.Println(responseBody)
			}
		} else {
			fmt.Println(responseBody)
		}
		
		return nil
	},
}

// deleteNotificationTemplateCmd represents the delete-notification-template command
var deleteNotificationTemplateCmd = &cobra.Command{
	Use:   "delete-notification-template",
	Short: "Delete an existing notification template",
	Long: `Delete a notification template by template ID and version.

This command will delete an existing notification template from the system.

Example:
  digit delete-notification-template --template-id user-notify --version v1
  digit delete-notification-template --template-id user-notify --version v1 --server http://localhost:8091`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get flag values
		templateID, _ := cmd.Flags().GetString("template-id")
		version, _ := cmd.Flags().GetString("version")
		serverURL, _ := cmd.Flags().GetString("server")
		jwtToken, _ := cmd.Flags().GetString("jwt-token")

		// Validate required flags
		if templateID == "" {
			return fmt.Errorf("--template-id flag is required")
		}
		if version == "" {
			return fmt.Errorf("--version flag is required")
		}

		// Get server URL and JWT token from config if not provided
		if serverURL == "" || jwtToken == "" {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			if serverURL == "" {
				serverURL = cfg.GetServer()
				if serverURL == "" {
					return fmt.Errorf("server URL not configured. Use 'digit config set --server <url>' or provide --server flag")
				}
			}

			if jwtToken == "" {
				jwtToken = cfg.GetJWTToken()
				if jwtToken == "" {
					return fmt.Errorf("JWT token not configured. Use 'digit config set --jwt-token <token>' or provide --jwt-token flag")
				}
			}
		}

		// Extract tenant ID from JWT token
		tenantID, err := jwt.ExtractTenantID(jwtToken)
		if err != nil {
			return fmt.Errorf("failed to extract tenant ID from JWT token: %w", err)
		}

		// Call the digit library to delete template
		responseBody, err := digit.DeleteNotificationTemplate(serverURL, jwtToken, tenantID, templateID, version)
		if err != nil {
			return fmt.Errorf("failed to delete notification template: %w", err)
		}

		// Print response body
		fmt.Printf("Notification template deletion response:\n")

		// Try to pretty print JSON response
		if responseBody != "" {
			var jsonResponse interface{}
			if err := json.Unmarshal([]byte(responseBody), &jsonResponse); err == nil {
				prettyJSON, err := json.MarshalIndent(jsonResponse, "", "  ")
				if err == nil {
					fmt.Println(string(prettyJSON))
				} else {
					fmt.Println(responseBody)
				}
			} else {
				fmt.Println(responseBody)
			}
		} else {
			fmt.Println("Notification template deleted successfully")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(createNotificationTemplateCmd)
	rootCmd.AddCommand(searchNotificationTemplateCmd)
	rootCmd.AddCommand(deleteNotificationTemplateCmd)
	
	// Add flags for create-notification-template command
	createNotificationTemplateCmd.Flags().String("file", "", "Path to YAML file containing all template configuration")
	createNotificationTemplateCmd.Flags().Bool("default", false, "Use default template configuration from template-config.yaml (requires --template-id)")
	createNotificationTemplateCmd.Flags().String("template-id", "", "Template ID for the notification template (required if not using --file)")
	createNotificationTemplateCmd.Flags().String("version", "", "Version of the template (required if not using --file or --default)")
	createNotificationTemplateCmd.Flags().String("type", "", "Type of template (EMAIL, SMS, etc.) (required if not using --file or --default)")
	createNotificationTemplateCmd.Flags().String("subject", "", "Subject of the template (required if not using --file or --default)")
	createNotificationTemplateCmd.Flags().String("content", "", "Content of the template (use either --content or --content-file)")
	createNotificationTemplateCmd.Flags().String("content-file", "", "Path to file containing template content (use either --content or --content-file)")
	createNotificationTemplateCmd.Flags().Bool("html", false, "Whether the content is HTML (default: false)")
	createNotificationTemplateCmd.Flags().String("server", "", "Server URL (overrides config)")
	createNotificationTemplateCmd.Flags().String("jwt-token", "", "JWT token for authentication (overrides config)")
	
	// Add flags for search-notification-template command
	searchNotificationTemplateCmd.Flags().String("template-id", "", "Template ID to search for (required)")
	searchNotificationTemplateCmd.Flags().String("server", "", "Server URL (overrides config)")
	searchNotificationTemplateCmd.Flags().String("jwt-token", "", "JWT token for authentication (overrides config)")
	
	// Add flags for delete-notification-template command
	deleteNotificationTemplateCmd.Flags().String("template-id", "", "Template ID to delete (required)")
	deleteNotificationTemplateCmd.Flags().String("version", "", "Version of the template to delete (required)")
	deleteNotificationTemplateCmd.Flags().String("server", "", "Server URL (overrides config)")
	deleteNotificationTemplateCmd.Flags().String("jwt-token", "", "JWT token for authentication (overrides config)")

	deleteNotificationTemplateCmd.MarkFlagRequired("template-id")
	deleteNotificationTemplateCmd.MarkFlagRequired("version")
	
	// Note: Required flags are validated conditionally in the command logic based on --file, --default, or individual flags
}