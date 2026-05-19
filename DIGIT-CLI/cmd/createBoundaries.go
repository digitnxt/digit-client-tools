package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"digit-cli/pkg/config"
	"digit-cli/pkg/errors"
	"digit-cli/pkg/jwt"
	"github.com/digitnxt/digit-client-tools/client-libraries/digit-library/digit"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// defaultBoundaryHierarchyYAML contains the embedded default boundary hierarchy configuration
const defaultBoundaryHierarchyYAML = `boundaryHierarchy:
  hierarchyType: "state-district-hierarchy"
  boundaryHierarchy:
    - boundaryType: "state"
      parentBoundaryType: null
      active: true
    - boundaryType: "district"
      parentBoundaryType: "state"
      active: true`
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

// BoundaryHierarchyDefinition represents the YAML structure for boundary hierarchy data
type BoundaryHierarchyDefinition struct {
	BoundaryHierarchy map[string]interface{} `yaml:"boundaryHierarchy"`
}

// BoundaryRelationshipDefinition represents the YAML structure for boundary relationship data
type BoundaryRelationshipDefinition struct {
	BoundaryRelationship map[string]interface{} `yaml:"boundaryRelationship"`
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
			return errors.HandleAPIError(err)
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

// createBoundaryHierarchyCmd represents the create-boundary-hierarchy command
var createBoundaryHierarchyCmd = &cobra.Command{
	Use:   "create-boundary-hierarchy",
	Short: "Create boundary hierarchy from YAML file or using default configuration",
	Long: `Create boundary hierarchy from a YAML file definition or using default configuration.
	
Examples:
  # Create boundary hierarchy from YAML file
  digit create-boundary-hierarchy --file boundary-hierarchy.yaml
  
  # Create boundary hierarchy using default configuration
  digit create-boundary-hierarchy --default
  
  # Create boundary hierarchy using default configuration with custom hierarchy type
  digit create-boundary-hierarchy --default --hierarchy-type "custom-hierarchy"
  
  # With server override
  digit create-boundary-hierarchy --file boundary-hierarchy.yaml --server http://localhost:8080`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get flag values
		filePath, _ := cmd.Flags().GetString("file")
		useDefault, _ := cmd.Flags().GetBool("default")
		hierarchyType, _ := cmd.Flags().GetString("hierarchy-type")
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
			// Set default hierarchy type if not provided
			defaultHierarchyType := "state-district-hierarchy"
			if hierarchyType != "" {
				defaultHierarchyType = hierarchyType
			}
			
			// Use embedded default configuration and replace the hierarchy type if needed
			yamlContent := defaultBoundaryHierarchyYAML
			if hierarchyType != "" && hierarchyType != "state-district-hierarchy" {
				yamlContent = strings.ReplaceAll(yamlContent, "state-district-hierarchy", hierarchyType)
			}
			
			yamlData = []byte(yamlContent)
			fmt.Printf("Using default boundary hierarchy configuration with hierarchy type: %s\n", defaultHierarchyType)
		} else {
			// Read YAML file
			var err error
			yamlData, err = os.ReadFile(filePath)
			if err != nil {
				return fmt.Errorf("failed to read YAML file: %w", err)
			}
		}
		
		// Parse YAML
		var boundaryHierarchyDef BoundaryHierarchyDefinition
		if err := yaml.Unmarshal(yamlData, &boundaryHierarchyDef); err != nil {
			return fmt.Errorf("failed to parse YAML: %w", err)
		}
		
		// Validate required fields from YAML
		if boundaryHierarchyDef.BoundaryHierarchy == nil {
			return fmt.Errorf("boundaryHierarchy is required in YAML file")
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
		
		// Call the digit library to create boundary hierarchy
		responseBody, err := digit.CreateBoundaryHierarchy(serverURL, jwtToken, tenantID, clientID, boundaryHierarchyDef.BoundaryHierarchy)
		if err != nil {
			return errors.HandleAPIError(err)
		}
		
		// Print response body
		fmt.Printf("Boundary hierarchy creation response:\n")
		
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

// searchBoundaryHierarchyCmd represents the search-boundary-hierarchy command
var searchBoundaryHierarchyCmd = &cobra.Command{
	Use:   "search-boundary-hierarchy",
	Short: "Search for boundary hierarchy by hierarchy type",
	Long: `Search for boundary hierarchy by providing the hierarchy type.
	
Examples:
  # Search boundary hierarchy by hierarchy type
  digit search-boundary-hierarchy --hierarchy-type "state-district-hierarchy"
  
  # With server override
  digit search-boundary-hierarchy --hierarchy-type "state-district-hierarchy" --server http://localhost:8080`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get flag values
		hierarchyType, _ := cmd.Flags().GetString("hierarchy-type")
		serverURL, _ := cmd.Flags().GetString("server")
		jwtToken, _ := cmd.Flags().GetString("jwt-token")
		
		// Validate required flags
		if hierarchyType == "" {
			return fmt.Errorf("--hierarchy-type flag is required")
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
		
		// Call the digit library to search boundary hierarchy
		responseBody, err := digit.SearchBoundaryHierarchy(serverURL, jwtToken, tenantID, clientID, hierarchyType)
		if err != nil {
			return errors.HandleAPIError(err)
		}
		
		// Print response body
		fmt.Printf("Boundary hierarchy search response:\n")
		
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

// createBoundaryRelationshipsCmd represents the create-boundary-relationships command
var createBoundaryRelationshipsCmd = &cobra.Command{
	Use:   "create-boundary-relationships",
	Short: "Create boundary relationships using command flags",
	Long: `Create boundary relationships by providing the required parameters as flags.
	
Examples:
  # Create boundary relationships with flags
  digit create-boundary-relationships --code "STATE1" --hierarchy-type "state-district-hierarchy" --boundary-type "state"
  
  # Create boundary relationships with parent
  digit create-boundary-relationships --code "DISTRICT1" --hierarchy-type "state-district-hierarchy" --boundary-type "district" --parent "STATE1"
  
  # With server override
  digit create-boundary-relationships --code "STATE1" --hierarchy-type "state-district-hierarchy" --boundary-type "state" --server http://localhost:8093`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get flag values
		code, _ := cmd.Flags().GetString("code")
		hierarchyType, _ := cmd.Flags().GetString("hierarchy-type")
		boundaryType, _ := cmd.Flags().GetString("boundary-type")
		parent, _ := cmd.Flags().GetString("parent")
		serverURL, _ := cmd.Flags().GetString("server")
		jwtToken, _ := cmd.Flags().GetString("jwt-token")
		
		// Validate required flags
		if code == "" {
			return fmt.Errorf("--code flag is required")
		}
		if hierarchyType == "" {
			return fmt.Errorf("--hierarchy-type flag is required")
		}
		if boundaryType == "" {
			return fmt.Errorf("--boundary-type flag is required")
		}
		
		// Build boundary relationship data
		boundaryRelationship := map[string]interface{}{
			"code":          code,
			"hierarchyType": hierarchyType,
			"boundaryType":  boundaryType,
		}
		
		// Add parent if provided, otherwise set to null
		if parent != "" {
			boundaryRelationship["parent"] = parent
		} else {
			boundaryRelationship["parent"] = nil
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
					serverURL = "http://localhost:8093" // Default server URL for boundary relationships
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
		
		// Call the digit library to create boundary relationships
		responseBody, err := digit.CreateBoundaryRelationships(serverURL, jwtToken, tenantID, clientID, boundaryRelationship)
		if err != nil {
			return errors.HandleAPIError(err)
		}
		
		// Print response body
		fmt.Printf("Boundary relationships creation response:\n")
		
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

// searchBoundaryRelationshipsCmd represents the search-boundary-relationships command
var searchBoundaryRelationshipsCmd = &cobra.Command{
	Use:   "search-boundary-relationships",
	Short: "Search for boundary relationships using query parameters",
	Long: `Search for boundary relationships by providing hierarchy type, boundary type, and optional parameters.
	
Examples:
  # Search boundary relationships by hierarchy and boundary type
  digit search-boundary-relationships --hierarchy-type "state-district" --boundary-type "state"
  
  # Search with specific codes
  digit search-boundary-relationships --hierarchy-type "state-district" --boundary-type "state" --codes "STATE1"
  
  # Search with children included
  digit search-boundary-relationships --hierarchy-type "state-district" --boundary-type "state" --codes "STATE1" --include-children
  
  # With server override
  digit search-boundary-relationships --hierarchy-type "state-district" --boundary-type "state" --server http://localhost:8080`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get flag values
		hierarchyType, _ := cmd.Flags().GetString("hierarchy-type")
		boundaryType, _ := cmd.Flags().GetString("boundary-type")
		codes, _ := cmd.Flags().GetString("codes")
		includeChildren, _ := cmd.Flags().GetBool("include-children")
		serverURL, _ := cmd.Flags().GetString("server")
		jwtToken, _ := cmd.Flags().GetString("jwt-token")
		
		// Validate required flags
		if hierarchyType == "" {
			return fmt.Errorf("--hierarchy-type flag is required")
		}
		if boundaryType == "" {
			return fmt.Errorf("--boundary-type flag is required")
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
		
		// Call the digit library to search boundary relationships
		responseBody, err := digit.SearchBoundaryRelationships(serverURL, jwtToken, tenantID, clientID, hierarchyType, boundaryType, codes, includeChildren)
		if err != nil {
			return errors.HandleAPIError(err)
		}
		
		// Print response body
		fmt.Printf("Boundary relationships search response:\n")
		
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
	rootCmd.AddCommand(createBoundaryHierarchyCmd)
	rootCmd.AddCommand(searchBoundaryHierarchyCmd)
	rootCmd.AddCommand(createBoundaryRelationshipsCmd)
	rootCmd.AddCommand(searchBoundaryRelationshipsCmd)
	
	// Add flags for create-boundaries command
	createBoundariesCmd.Flags().String("file", "", "Path to YAML file containing boundary data")
	createBoundariesCmd.Flags().Bool("default", false, "Use default boundary configuration")
	createBoundariesCmd.Flags().String("code-prefix", "", "Code prefix for default boundaries (default: 'DEFAULT')")
	createBoundariesCmd.Flags().String("server", "", "Server URL (overrides config)")
	createBoundariesCmd.Flags().String("jwt-token", "", "JWT token for authentication (overrides config)")
	
	// Add flags for create-boundary-hierarchy command
	createBoundaryHierarchyCmd.Flags().String("file", "", "Path to YAML file containing boundary hierarchy data")
	createBoundaryHierarchyCmd.Flags().Bool("default", false, "Use default boundary hierarchy configuration")
	createBoundaryHierarchyCmd.Flags().String("hierarchy-type", "", "Hierarchy type for default configuration (default: 'state-district-hierarchy')")
	createBoundaryHierarchyCmd.Flags().String("server", "", "Server URL (overrides config)")
	createBoundaryHierarchyCmd.Flags().String("jwt-token", "", "JWT token for authentication (overrides config)")
	
	// Add flags for search-boundary-hierarchy command
	searchBoundaryHierarchyCmd.Flags().String("hierarchy-type", "", "Hierarchy type to search for (required)")
	searchBoundaryHierarchyCmd.Flags().String("server", "", "Server URL (overrides config)")
	searchBoundaryHierarchyCmd.Flags().String("jwt-token", "", "JWT token for authentication (overrides config)")
	
	// Add flags for create-boundary-relationships command
	createBoundaryRelationshipsCmd.Flags().String("code", "", "Boundary code (required)")
	createBoundaryRelationshipsCmd.Flags().String("hierarchy-type", "", "Hierarchy type (required)")
	createBoundaryRelationshipsCmd.Flags().String("boundary-type", "", "Boundary type (required)")
	createBoundaryRelationshipsCmd.Flags().String("parent", "", "Parent boundary code (optional, null if not provided)")
	createBoundaryRelationshipsCmd.Flags().String("server", "", "Server URL (overrides config)")
	createBoundaryRelationshipsCmd.Flags().String("jwt-token", "", "JWT token for authentication (overrides config)")
	
	// Add flags for search-boundary-relationships command
	searchBoundaryRelationshipsCmd.Flags().String("hierarchy-type", "", "Hierarchy type (required)")
	searchBoundaryRelationshipsCmd.Flags().String("boundary-type", "", "Boundary type (required)")
	searchBoundaryRelationshipsCmd.Flags().String("codes", "", "Comma-separated boundary codes (optional)")
	searchBoundaryRelationshipsCmd.Flags().Bool("include-children", false, "Include children in the response (optional)")
	searchBoundaryRelationshipsCmd.Flags().String("server", "", "Server URL (overrides config)")
	searchBoundaryRelationshipsCmd.Flags().String("jwt-token", "", "JWT token for authentication (overrides config)")
	
	// Mark required flags
	searchBoundaryHierarchyCmd.MarkFlagRequired("hierarchy-type")
	createBoundaryRelationshipsCmd.MarkFlagRequired("code")
	createBoundaryRelationshipsCmd.MarkFlagRequired("hierarchy-type")
	createBoundaryRelationshipsCmd.MarkFlagRequired("boundary-type")
	searchBoundaryRelationshipsCmd.MarkFlagRequired("hierarchy-type")
	searchBoundaryRelationshipsCmd.MarkFlagRequired("boundary-type")
	
	// Note: Required flags are validated conditionally in the command logic
}