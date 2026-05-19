package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"digit-cli/pkg/config"
	"digit-cli/pkg/jwt"
	"github.com/digitnxt/digit-client-tools/client-libraries/digit-library/digit"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// defaultMdmsSchemaYAML contains the embedded default MDMS schema configuration
const defaultMdmsSchemaYAML = `schema:
  code: "RAINMAKER_PGR_ServiceDefs"
  description: "Schema for PGR Service Definitions"
  definition:
    $schema: "http://json-schema.org/draft-07/schema#"
    type: "object"
    properties:
      serviceCode:
        type: "string"
      name:
        type: "string"
      keywords:
        type: "string"
      department:
        type: "string"
      slaHours:
        type: "number"
      order:
        type: "number"
      active:
        type: "boolean"
      menuPath:
        type: "string"
    required: ["serviceCode", "name", "department", "slaHours", "active"]
    x-unique: ["serviceCode"]
    x-ref-schema: []
  isActive: true`

// SchemaDefinition represents the YAML structure for MDMS schema definition
type SchemaDefinition struct {
	Schema struct {
		Code        string      `yaml:"code"`
		Description string      `yaml:"description"`
		Definition  interface{} `yaml:"definition"`
		IsActive    bool        `yaml:"isActive"`
	} `yaml:"schema"`
}

// defaultMdmsDataYAML contains the embedded default MDMS data configuration
const defaultMdmsDataYAML = `mdms:
  - schemaCode: "RAINMAKER_PGR_ServiceDefs"
    uniqueIdentifier: "WATER_SUPPLY"
    data:
      serviceCode: "WATER_SUPPLY"
      name: "Water Supply Issue"
      keywords: "water, supply, shortage, quality"
      department: "Water Department"
      slaHours: 72
      order: 1
      active: true
      menuPath: "Complaints.Water"
    isActive: true
  - schemaCode: "RAINMAKER_PGR_ServiceDefs"
    uniqueIdentifier: "ROAD_REPAIR"
    data:
      serviceCode: "ROAD_REPAIR"
      name: "Road Repair Request"
      keywords: "road, repair, pothole, maintenance"
      department: "Public Works Department"
      slaHours: 168
      order: 2
      active: true
      menuPath: "Complaints.Infrastructure"
    isActive: true
  - schemaCode: "RAINMAKER_PGR_ServiceDefs"
    uniqueIdentifier: "GARBAGE_COLLECTION"
    data:
      serviceCode: "GARBAGE_COLLECTION"
      name: "Garbage Collection Issue"
      keywords: "garbage, waste, collection, sanitation"
      department: "Sanitation Department"
      slaHours: 24
      order: 3
      active: true
      menuPath: "Complaints.Sanitation"
    isActive: true`
type MdmsDataDefinition struct {
	Mdms []struct {
		SchemaCode       string                 `yaml:"schemaCode" json:"schemaCode"`
		UniqueIdentifier string                 `yaml:"uniqueIdentifier" json:"uniqueIdentifier"`
		Data             map[string]interface{} `yaml:"data" json:"data"`
		IsActive         bool                   `yaml:"isActive" json:"isActive"`
	} `yaml:"mdms"`
}



// createMdmsSchemaCmd represents the create-mdms-schema command
var createMdmsSchemaCmd = &cobra.Command{
	Use:   "create-mdms-schema",
	Short: "Create a new MDMS schema from YAML file or using default configuration",
	Long: `Create a new MDMS schema from a YAML file definition or using default configuration.
	
Examples:
  # Create MDMS schema from YAML file
  digit create-mdms-schema --file schema.yaml
  
  # Create MDMS schema using default configuration with custom code
  digit create-mdms-schema --default --code "MY_CUSTOM_SCHEMA"
  
  # Create MDMS schema using default configuration (uses RAINMAKER_PGR_ServiceDefs)
  digit create-mdms-schema --default
  
  # With server override
  digit create-mdms-schema --file schema.yaml --server http://localhost:8094`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get flag values
		filePath, _ := cmd.Flags().GetString("file")
		useDefault, _ := cmd.Flags().GetBool("default")
		code, _ := cmd.Flags().GetString("code")
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
			// Set default schema code if not provided
			defaultSchemaCode := "RAINMAKER_PGR_ServiceDefs"
			if code != "" {
				defaultSchemaCode = code
			}
			
			// Use embedded default configuration and replace the code if needed
			yamlContent := defaultMdmsSchemaYAML
			if code != "" && code != "RAINMAKER_PGR_ServiceDefs" {
				yamlContent = strings.ReplaceAll(yamlContent, "RAINMAKER_PGR_ServiceDefs", code)
			}
			
			yamlData = []byte(yamlContent)
			fmt.Printf("Using default MDMS schema configuration with code: %s\n", defaultSchemaCode)
		} else {
			// Read YAML file
			var err error
			yamlData, err = os.ReadFile(filePath)
			if err != nil {
				return fmt.Errorf("failed to read YAML file: %w", err)
			}
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
	Short: "Create MDMS data entries from YAML file or using default configuration",
	Long: `Create MDMS data entries from a YAML file definition or using default configuration.
	
Examples:
  # Create MDMS data from YAML file
  digit create-mdms-data --file mdms-data.yaml
  
  # Create MDMS data using default configuration
  digit create-mdms-data --default
  
  # Create MDMS data using default configuration with custom schema code
  digit create-mdms-data --default --schema-code "MY_CUSTOM_SCHEMA"
  
  # With server override
  digit create-mdms-data --file mdms-data.yaml --server http://localhost:8081`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get flag values
		filePath, _ := cmd.Flags().GetString("file")
		useDefault, _ := cmd.Flags().GetBool("default")
		schemaCode, _ := cmd.Flags().GetString("schema-code")
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
			// Set default schema code if not provided
			defaultSchemaCode := "RAINMAKER_PGR_ServiceDefs"
			if schemaCode != "" {
				defaultSchemaCode = schemaCode
			}
			
			// Use embedded default configuration and replace the schema code if needed
			yamlContent := defaultMdmsDataYAML
			if schemaCode != "" && schemaCode != "RAINMAKER_PGR_ServiceDefs" {
				yamlContent = strings.ReplaceAll(yamlContent, "RAINMAKER_PGR_ServiceDefs", schemaCode)
			}
			
			yamlData = []byte(yamlContent)
			fmt.Printf("Using default MDMS data configuration with schema code: %s\n", defaultSchemaCode)
		} else {
			// Read YAML file
			var err error
			yamlData, err = os.ReadFile(filePath)
			if err != nil {
				return fmt.Errorf("failed to read YAML file: %w", err)
			}
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

// searchMdmsSchemaCmd represents the search-mdms-schema command
var searchMdmsSchemaCmd = &cobra.Command{
	Use:   "search-mdms-schema",
	Short: "Search for an MDMS schema by code",
	Long: `Search for an MDMS schema by providing the schema code.
	
Examples:
  # Search MDMS schema by code
  digit search-mdms-schema --code "RAINMAKER-PGR.ServiceDefsns"
  
  # With server override
  digit search-mdms-schema --code "RAINMAKER-PGR.ServiceDefsns" --server http://localhost:8099`,
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
	rootCmd.AddCommand(createMdmsSchemaCmd)
	rootCmd.AddCommand(createMdmsDataCmd)
	rootCmd.AddCommand(searchMdmsSchemaCmd)
	rootCmd.AddCommand(searchMdmsDataCmd)
	
	// Add flags for create-mdms-schema command
	createMdmsSchemaCmd.Flags().String("file", "", "Path to YAML file containing schema definition")
	createMdmsSchemaCmd.Flags().Bool("default", false, "Use default MDMS schema configuration")
	createMdmsSchemaCmd.Flags().String("code", "", "Schema code for default configuration (default: 'RAINMAKER_PGR_ServiceDefs')")
	createMdmsSchemaCmd.Flags().String("server", "", "Server URL (overrides config)")
	createMdmsSchemaCmd.Flags().String("jwt-token", "", "JWT token for authentication (overrides config)")
	
	// Note: Required flags are validated conditionally in the command logic
	
	// Add flags for create-mdms-data command
	createMdmsDataCmd.Flags().String("file", "", "Path to YAML file containing MDMS data")
	createMdmsDataCmd.Flags().Bool("default", false, "Use default MDMS data configuration")
	createMdmsDataCmd.Flags().String("schema-code", "", "Schema code for default configuration (default: 'RAINMAKER_PGR_ServiceDefs')")
	createMdmsDataCmd.Flags().String("server", "", "Server URL (overrides config)")
	createMdmsDataCmd.Flags().String("jwt-token", "", "JWT token for authentication (overrides config)")
	
	// Note: Required flags are validated conditionally in the command logic
	
	// Add flags for search-mdms-schema command
	searchMdmsSchemaCmd.Flags().String("code", "", "Schema code to search for (required)")
	searchMdmsSchemaCmd.Flags().String("server", "", "Server URL (overrides config)")
	searchMdmsSchemaCmd.Flags().String("jwt-token", "", "JWT token for authentication (overrides config)")
	
	// Mark required flags
	searchMdmsSchemaCmd.MarkFlagRequired("code")
	
	// Add flags for search-mdms-data command
	searchMdmsDataCmd.Flags().String("code", "", "Schema code to search MDMS data for (required)")
	searchMdmsDataCmd.Flags().String("unique-identifiers", "", "Comma-separated unique identifiers to filter data (optional)")
	searchMdmsDataCmd.Flags().String("server", "", "Server URL (overrides config)")
	searchMdmsDataCmd.Flags().String("jwt-token", "", "JWT token for authentication (overrides config)")
	
	// Mark required flags
	searchMdmsDataCmd.MarkFlagRequired("code")
}
