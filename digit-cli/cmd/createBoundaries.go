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

// defaultBoundaryYAML contains the embedded default boundary configuration
const defaultBoundaryYAML = `boundary:
  - code: "DEFAULT_BOUNDARY_001"
    geometry:
      type: "Polygon"
      coordinates:
        - - [77.0, 28.5]
          - [77.1, 28.5]
          - [77.1, 28.6]
          - [77.0, 28.6]
          - [77.0, 28.5]
    additionalDetails: {}
  - code: "DEFAULT_BOUNDARY_002"
    geometry:
      type: "Point"
      coordinates: [77.0, 28.5]
    additionalDetails: {}
  - code: "DEFAULT_BOUNDARY_003"
    geometry:
      type: "Polygon"
      coordinates:
        - - [77.2, 28.7]
          - [77.3, 28.7]
          - [77.3, 28.8]
          - [77.2, 28.8]
          - [77.2, 28.7]
    additionalDetails: {}`

// BoundaryDefinition represents the YAML structure for boundary data
type BoundaryDefinition struct {
	Boundary []map[string]interface{} `yaml:"boundary"`
}

// createBoundariesCmd represents the create-boundaries command
var createBoundariesCmd = &cobra.Command{
	Use:   "create-boundaries",
	Short: "Create boundaries from YAML file",
	Long: `Create boundaries from a YAML file definition or using default configuration.
	
Examples:
  # Create boundaries from YAML file
  digit create-boundaries --file boundaries.yaml
  
  # Create boundaries using default configuration with custom prefix
  digit create-boundaries --default --code-prefix "CUSTOM"
  
  # Create boundaries using default configuration (uses DEFAULT prefix)
  digit create-boundaries --default
  
  # With server override
  digit create-boundaries --file boundaries.yaml --server http://localhost:8080`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get flag values
		filePath, _ := cmd.Flags().GetString("file")
		useDefault, _ := cmd.Flags().GetBool("default")
		codePrefix, _ := cmd.Flags().GetString("code-prefix")
		serverURL, _ := cmd.Flags().GetString("server")
		jwtToken, _ := cmd.Flags().GetString("jwt-token")
		
		// Validate flags - either file or default must be specified
		if !useDefault && filePath == "" {
			return fmt.Errorf("either --file or --default flag is required")
		}
		if useDefault && filePath != "" {
			return fmt.Errorf("cannot use both --file and --default flags together")
		}
		
		// Get YAML data - either from file or default configuration
		var yamlData []byte
		if useDefault {
			// Use embedded default configuration and replace the code prefix
			if codePrefix == "" {
				codePrefix = "DEFAULT"
			}
			yamlContent := strings.ReplaceAll(defaultBoundaryYAML, "DEFAULT_BOUNDARY", codePrefix+"_BOUNDARY")
			yamlData = []byte(yamlContent)
			fmt.Printf("Using default boundary configuration with code prefix: %s\n", codePrefix)
		} else {
			// Read YAML file
			var err error
			yamlData, err = os.ReadFile(filePath)
			if err != nil {
				return fmt.Errorf("failed to read YAML file: %w", err)
			}
		}
		
		// Parse YAML
		var boundaryDef BoundaryDefinition
		if err := yaml.Unmarshal(yamlData, &boundaryDef); err != nil {
			return fmt.Errorf("failed to parse YAML: %w", err)
		}
		
		// Validate required fields from YAML
		if len(boundaryDef.Boundary) == 0 {
			return fmt.Errorf("at least one boundary entry is required in YAML file")
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
					serverURL = "http://localhost:8080" // Default server URL
				}
			}
			
			if jwtToken == "" {
				jwtToken = cfg.GetJWTToken()
			}
		}

		// Extract tenant ID and client ID from JWT token
		if jwtToken == "" {
			return fmt.Errorf("JWT token is required")
		}
		
		tenantID, err := jwt.ExtractTenantID(jwtToken)
		if err != nil {
			return fmt.Errorf("failed to extract tenant ID from JWT token: %w", err)
		}
		
		clientID, err := jwt.ExtractClientID(jwtToken)
		if err != nil {
			return fmt.Errorf("failed to extract client ID from JWT token: %w", err)
		}
		
		// Call the digit library to create boundaries
		responseBody, err := digit.CreateBoundaries(serverURL, jwtToken, tenantID, clientID, boundaryDef.Boundary)
		if err != nil {
			return fmt.Errorf("failed to create boundaries: %w", err)
		}
		
		// Print response body
		fmt.Printf("Boundary creation response:\n")
		
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

func init() {
	rootCmd.AddCommand(createBoundariesCmd)
	
	// Add flags for create-boundaries command
	createBoundariesCmd.Flags().String("file", "", "Path to YAML file containing boundary data")
	createBoundariesCmd.Flags().Bool("default", false, "Use default boundary configuration")
	createBoundariesCmd.Flags().String("code-prefix", "", "Code prefix for default boundaries (default: 'DEFAULT')")
	createBoundariesCmd.Flags().String("server", "", "Server URL (overrides config)")
	createBoundariesCmd.Flags().String("jwt-token", "", "JWT token for authentication (overrides config)")
	
	// Note: Required flags are validated conditionally in the command logic
}