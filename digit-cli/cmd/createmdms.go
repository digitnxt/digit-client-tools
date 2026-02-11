package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"digit-cli/pkg/config"
	"digit-cli/pkg/jwt"
	"github.com/digitnxt/digit3/code/libraries/digit-library/digit"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// SchemaDefinition represents the YAML structure for MDMS schema definition
type SchemaDefinition struct {
	Schema struct {
		Code        string      `yaml:"code"`
		Description string      `yaml:"description"`
		Definition  interface{} `yaml:"definition"`
		IsActive    bool        `yaml:"isActive"`
	} `yaml:"schema"`
}

// MdmsDataDefinition represents the YAML structure for MDMS data entries
type MdmsDataDefinition struct {
	Mdms []struct {
		SchemaCode       string                 `yaml:"schemaCode" json:"schemaCode"`
		UniqueIdentifier string                 `yaml:"uniqueIdentifier" json:"uniqueIdentifier"`
		Data             map[string]interface{} `yaml:"data" json:"data"`
		IsActive         bool                   `yaml:"isActive" json:"isActive"`
	} `yaml:"mdms"`
}



// createSchemaCmd represents the create-schema command
var createSchemaCmd = &cobra.Command{
	Use:   "create-schema",
	Short: "Create a new MDMS schema from YAML file",
	Long: `Create a new MDMS schema from a YAML file definition.
	
Examples:
  # Create schema from YAML file
  digit create-schema --file schema.yaml
  
  # With server override
  digit create-schema --file schema.yaml --server http://localhost:8094`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get flag values
		filePath, _ := cmd.Flags().GetString("file")
		serverURL, _ := cmd.Flags().GetString("server")
		jwtToken, _ := cmd.Flags().GetString("jwt-token")
		
		// Validate required flags
		if filePath == "" {
			return fmt.Errorf("--file flag is required")
		}
		
		// Read YAML file
		yamlData, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read YAML file: %w", err)
		}
		
		// Parse YAML
		var schemaDef SchemaDefinition
		if err := yaml.Unmarshal(yamlData, &schemaDef); err != nil {
			return fmt.Errorf("failed to parse YAML: %w", err)
		}
		
		// Validate required fields from YAML
		if schemaDef.Schema.Code == "" {
			return fmt.Errorf("code is required in YAML file")
		}
		if schemaDef.Schema.Description == "" {
			return fmt.Errorf("description is required in YAML file")
		}
		
		// Convert definition to JSON string
		definitionBytes, err := json.Marshal(schemaDef.Schema.Definition)
		if err != nil {
			return fmt.Errorf("failed to marshal definition to JSON: %w", err)
		}
		definition := string(definitionBytes)
		
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
				// JWT token is optional for MDMS, so we don't error if it's missing
			}
		}

		// Extract tenant ID, client ID, and unique ID from JWT token
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
		
		// Call the digit library to create schema
		responseBody, err := digit.CreateSchema(serverURL, jwtToken, tenantID, clientID, schemaDef.Schema.Code, schemaDef.Schema.Description, definition, schemaDef.Schema.IsActive)
		if err != nil {
			return fmt.Errorf("failed to create schema: %w", err)
		}
		
		// Print response body
		fmt.Printf("Schema creation response:\n")
		
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

// createMdmsDataCmd represents the create-mdms-data command
var createMdmsDataCmd = &cobra.Command{
	Use:   "create-mdms-data",
	Short: "Create MDMS data entries from YAML file",
	Long: `Create MDMS data entries from a YAML file definition.
	
Examples:
  # Create MDMS data from YAML file
  digit create-mdms-data --file mdms-data.yaml
  
  # With server override
  digit create-mdms-data --file mdms-data.yaml --server http://localhost:8081`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get flag values
		filePath, _ := cmd.Flags().GetString("file")
		serverURL, _ := cmd.Flags().GetString("server")
		jwtToken, _ := cmd.Flags().GetString("jwt-token")
		
		// Validate required flags
		if filePath == "" {
			return fmt.Errorf("--file flag is required")
		}
		
		// Read YAML file
		yamlData, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read YAML file: %w", err)
		}
		
		// Parse YAML
		var mdmsDataDef MdmsDataDefinition
		if err := yaml.Unmarshal(yamlData, &mdmsDataDef); err != nil {
			return fmt.Errorf("failed to parse YAML: %w", err)
		}
		
		// Validate required fields from YAML
		if len(mdmsDataDef.Mdms) == 0 {
			return fmt.Errorf("at least one MDMS entry is required in YAML file")
		}
		
		// Convert MDMS data to JSON string - just the Mdms array without the top-level wrapper
		mdmsDataBytes, err := json.Marshal(mdmsDataDef.Mdms)
		if err != nil {
			return fmt.Errorf("failed to marshal MDMS data to JSON: %w", err)
		}
		mdmsData := string(mdmsDataBytes)
		
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
		
		// Call the digit library to create MDMS data
		responseBody, err := digit.CreateMdmsData(serverURL, jwtToken, tenantID, clientID, mdmsData)
		if err != nil {
			return fmt.Errorf("failed to create MDMS data: %w", err)
		}
		
		// Print response body
		fmt.Printf("MDMS data creation response:\n")
		
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

// searchSchemaCmd represents the search-schema command
var searchSchemaCmd = &cobra.Command{
	Use:   "search-schema",
	Short: "Search for an MDMS schema by code",
	Long: `Search for an MDMS schema by providing the schema code.
	
Examples:
  # Search schema by code
  digit search-schema --code "RAINMAKER-PGR.ServiceDefsns"
  
  # With server override
  digit search-schema --code "RAINMAKER-PGR.ServiceDefsns" --server http://localhost:8099`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get flag values
		schemaCode, _ := cmd.Flags().GetString("code")
		serverURL, _ := cmd.Flags().GetString("server")
		jwtToken, _ := cmd.Flags().GetString("jwt-token")
		
		// Validate required flags
		if schemaCode == "" {
			return fmt.Errorf("--code flag is required")
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
		
		// Call the digit library to search schema
		responseBody, err := digit.SearchSchema(serverURL, jwtToken, tenantID, clientID, schemaCode)
		if err != nil {
			return fmt.Errorf("failed to search schema: %w", err)
		}
		
		// Print response body
		fmt.Printf("Schema response:\n")
		
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

// searchMdmsDataCmd represents the search-mdms-data command
var searchMdmsDataCmd = &cobra.Command{
	Use:   "search-mdms-data",
	Short: "Search MDMS data by schema code",
	Long: `Search MDMS data by providing the schema code and optional unique identifiers.
	
Examples:
  # Search MDMS data by schema code
  digit search-mdms-data --code "common-masters.abcd"
  
  # Search MDMS data with unique identifiers
  digit search-mdms-data --code "common-masters.abcd" --unique-identifiers "Alice1,Alice3"
  
  # With server override
  digit search-mdms-data --code "common-masters.abcd" --server http://localhost:8099`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get flag values
		schemaCode, _ := cmd.Flags().GetString("code")
		uniqueIdentifiers, _ := cmd.Flags().GetString("unique-identifiers")
		serverURL, _ := cmd.Flags().GetString("server")
		jwtToken, _ := cmd.Flags().GetString("jwt-token")
		
		// Validate required flags
		if schemaCode == "" {
			return fmt.Errorf("--code flag is required")
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
		
		// Call the digit library to search MDMS data
		responseBody, err := digit.SearchMdmsData(serverURL, jwtToken, tenantID, clientID, schemaCode, uniqueIdentifiers)
		if err != nil {
			return fmt.Errorf("failed to search MDMS data: %w", err)
		}
		
		// Print response body
		fmt.Printf("MDMS data response:\n")
		
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
	rootCmd.AddCommand(createSchemaCmd)
	rootCmd.AddCommand(createMdmsDataCmd)
	rootCmd.AddCommand(searchSchemaCmd)
	rootCmd.AddCommand(searchMdmsDataCmd)
	
	// Add flags for create-schema command
	createSchemaCmd.Flags().String("file", "", "Path to YAML file containing schema definition (required)")
	createSchemaCmd.Flags().String("server", "", "Server URL (overrides config)")
	createSchemaCmd.Flags().String("jwt-token", "", "JWT token for authentication (overrides config)")
	
	// Mark required flags
	createSchemaCmd.MarkFlagRequired("file")
	
	// Add flags for create-mdms-data command
	createMdmsDataCmd.Flags().String("file", "", "Path to YAML file containing MDMS data (required)")
	createMdmsDataCmd.Flags().String("server", "", "Server URL (overrides config)")
	createMdmsDataCmd.Flags().String("jwt-token", "", "JWT token for authentication (overrides config)")
	
	// Mark required flags
	createMdmsDataCmd.MarkFlagRequired("file")
	
	// Add flags for search-schema command
	searchSchemaCmd.Flags().String("code", "", "Schema code to search for (required)")
	searchSchemaCmd.Flags().String("server", "", "Server URL (overrides config)")
	searchSchemaCmd.Flags().String("jwt-token", "", "JWT token for authentication (overrides config)")
	
	// Mark required flags
	searchSchemaCmd.MarkFlagRequired("code")
	
	// Add flags for search-mdms-data command
	searchMdmsDataCmd.Flags().String("code", "", "Schema code to search MDMS data for (required)")
	searchMdmsDataCmd.Flags().String("unique-identifiers", "", "Comma-separated unique identifiers to filter data (optional)")
	searchMdmsDataCmd.Flags().String("server", "", "Server URL (overrides config)")
	searchMdmsDataCmd.Flags().String("jwt-token", "", "JWT token for authentication (overrides config)")
	
	// Mark required flags
	searchMdmsDataCmd.MarkFlagRequired("code")
}
