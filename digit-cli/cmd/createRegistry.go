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

// defaultPGRSchemaYAML contains the embedded PGR registry schema configuration
const defaultPGRSchemaYAML = `schemaCode: "pgr"
definition:
  $schema: "https://json-schema.org/draft/2020-12/schema"
  type: "object"
  additionalProperties: false
  properties:
    serviceRequestId:
      type: "string"
      description: "Unique identifier for the service request"
    tenantId:
      type: "string"
      description: "Tenant identifier"
    serviceCode:
      type: "string"
      description: "Code identifying the type of service"
    description:
      type: "string"
      description: "Description of the service request"
    accountId:
      type: "string"
      description: "Account identifier of the requester"
    source:
      type: "string"
      description: "Source of the service request"
    applicationStatus:
      type: "string"
      description: "Current status of the application"
    action:
      type: "string"
      description: "Action to be performed"
    fileStoreId:
      type: "string"
      description: "File store identifier for attachments"
    boundaryCode:
      type: "string"
      description: "Boundary/ward code where service is requested"
    individualId:
      type: "string"
      description: "Individual identifier of the requester"
    email:
      type: "string"
      format: "email"
      description: "Email address of the requester"
    mobile:
      type: "string"
      pattern: "^[0-9]{10}$"
      description: "Mobile number of the requester"
    processId:
      type: "string"
      description: "Process identifier for workflow"
    workflowInstanceId:
      type: "string"
      description: "Workflow instance identifier"
    auditDetails:
      type: "object"
      description: "Audit information for the record"
      properties:
        createdBy:
          type: "string"
        createdTime:
          type: "integer"
          format: "int64"
        lastModifiedBy:
          type: "string"
        lastModifiedTime:
          type: "integer"
          format: "int64"
    address:
      type: "object"
      description: "Address information for the service request"
      properties:
        id:
          type: "string"
        serviceRequestId:
          type: "string"
        address:
          type: "string"
        city:
          type: "string"
        pincode:
          type: "string"
          pattern: "^[0-9]{6}$"
        latitude:
          type: "number"
          format: "double"
          minimum: -90
          maximum: 90
        longitude:
          type: "number"
          format: "double"
          minimum: -180
          maximum: 180
        auditDetails:
          type: "object"
          properties:
            createdBy:
              type: "string"
            createdTime:
              type: "integer"
              format: "int64"
            lastModifiedBy:
              type: "string"
            lastModifiedTime:
              type: "integer"
              format: "int64"
  required:
    - serviceRequestId
    - tenantId
  x-indexes:
    - name: "idx_pgr2_service_request_id"
      fieldPath: "serviceRequestId"
      method: "btree"
    - name: "idx_pgr2_tenant_id"
      fieldPath: "tenantId"
      method: "btree"
    - name: "idx_pgr2_application_status"
      fieldPath: "applicationStatus"
      method: "btree"
    - fieldPath: "boundaryCode"
      method: "gin"`

// RegistrySchemaDefinition represents the YAML structure for registry schema data
type RegistrySchemaDefinition struct {
	SchemaCode string                 `yaml:"schemaCode"`
	Definition map[string]interface{} `yaml:"definition"`
}

// getDefaultSchemaYAML returns the PGR default schema YAML
func getDefaultSchemaYAML(schemaCode string) (string, error) {
	// Always return PGR schema, but allow custom schema code
	yamlContent := defaultPGRSchemaYAML
	if schemaCode != "" && schemaCode != "pgr" {
		yamlContent = strings.ReplaceAll(yamlContent, "pgr", schemaCode)
	}
	return yamlContent, nil
}

// createRegistrySchemaCmd represents the create-registry-schema command
var createRegistrySchemaCmd = &cobra.Command{
	Use:   "create-registry-schema",
	Short: "Create registry schema from YAML file",
	Long: `Create registry schema from a YAML file definition or using default PGR configuration.
	
Examples:
  # Create registry schema from YAML file
  digit create-registry-schema --file registry-schema.yaml
  
  # Create registry schema using default PGR configuration with custom schema code
  digit create-registry-schema --default --schema-code "custom-pgr"
  
  # Create registry schema using default PGR configuration
  digit create-registry-schema --default
  
  # With server override
  digit create-registry-schema --file registry-schema.yaml --server http://localhost:8085`,
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
			// Set default schema code to pgr if not provided
			defaultSchemaCode := "pgr"
			if schemaCode != "" {
				defaultSchemaCode = schemaCode
			}
			
			// Get the PGR default schema
			yamlContent, err := getDefaultSchemaYAML(defaultSchemaCode)
			if err != nil {
				return err
			}
			
			yamlData = []byte(yamlContent)
			fmt.Printf("Using default PGR registry schema configuration with schema code: %s\n", defaultSchemaCode)
		} else {
			// Read YAML file
			var err error
			yamlData, err = os.ReadFile(filePath)
			if err != nil {
				return fmt.Errorf("failed to read YAML file: %w", err)
			}
		}
		
		// Parse YAML
		var registryDef RegistrySchemaDefinition
		if err := yaml.Unmarshal(yamlData, &registryDef); err != nil {
			return fmt.Errorf("failed to parse YAML: %w", err)
		}
		
		// Validate required fields from YAML
		if registryDef.SchemaCode == "" {
			return fmt.Errorf("schemaCode is required in YAML file")
		}
		if registryDef.Definition == nil {
			return fmt.Errorf("definition is required in YAML file")
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
					serverURL = "http://localhost:8085" // Default server URL for registry
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
		
		// Call the digit library to create registry schema
		responseBody, err := digit.CreateRegistrySchema(serverURL, jwtToken, tenantID, clientID, registryDef.SchemaCode, registryDef.Definition)
		if err != nil {
			return errors.HandleAPIError(err)
		}
		
		// Print response body
		fmt.Printf("Registry schema creation response:\n")
		
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

// searchRegistrySchemaCmd represents the search-registry-schema command
var searchRegistrySchemaCmd = &cobra.Command{
	Use:   "search-registry-schema",
	Short: "Search for a registry schema by schema code",
	Long: `Search for a registry schema by schema code and optional version.
	
Examples:
  # Search registry schema by code
  digit search-registry-schema --schema-code "license-registry"
  
  # Search registry schema by code and version
  digit search-registry-schema --schema-code "license-registry" --version "1"
  
  # With server override
  digit search-registry-schema --schema-code "license-registry" --server http://localhost:8085`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get flag values
		schemaCode, _ := cmd.Flags().GetString("schema-code")
		version, _ := cmd.Flags().GetString("version")
		serverURL, _ := cmd.Flags().GetString("server")
		jwtToken, _ := cmd.Flags().GetString("jwt-token")
		
		// Validate required flags
		if schemaCode == "" {
			return fmt.Errorf("--schema-code flag is required")
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
					serverURL = "http://localhost:8085" // Default server URL for registry
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
		
		// Call the digit library to search registry schema
		responseBody, err := digit.SearchRegistrySchema(serverURL, jwtToken, tenantID, clientID, schemaCode, version)
		if err != nil {
			return errors.HandleAPIError(err)
		}
		
		// Print response body
		fmt.Printf("Registry schema search response:\n")
		
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

// deleteRegistrySchemaCmd represents the delete-registry-schema command
var deleteRegistrySchemaCmd = &cobra.Command{
	Use:   "delete-registry-schema",
	Short: "Delete a registry schema by schema code",
	Long: `Delete a registry schema by schema code.
	
Examples:
  # Delete registry schema by code
  digit delete-registry-schema --schema-code "license-registry"
  
  # With server override
  digit delete-registry-schema --schema-code "license-registry" --server http://localhost:8085`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get flag values
		schemaCode, _ := cmd.Flags().GetString("schema-code")
		serverURL, _ := cmd.Flags().GetString("server")
		jwtToken, _ := cmd.Flags().GetString("jwt-token")
		
		// Validate required flags
		if schemaCode == "" {
			return fmt.Errorf("--schema-code flag is required")
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
					serverURL = "http://localhost:8085" // Default server URL for registry
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
		
		// Call the digit library to delete registry schema
		responseBody, err := digit.DeleteRegistrySchema(serverURL, jwtToken, tenantID, clientID, schemaCode)
		if err != nil {
			return errors.HandleAPIError(err)
		}
		
		// Print response body
		fmt.Printf("Registry schema deletion response:\n")
		
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
	rootCmd.AddCommand(createRegistrySchemaCmd)
	rootCmd.AddCommand(searchRegistrySchemaCmd)
	rootCmd.AddCommand(deleteRegistrySchemaCmd)

	createRegistrySchemaCmd.Flags().String("file", "", "Path to YAML file containing registry schema data")
	createRegistrySchemaCmd.Flags().Bool("default", false, "Use default registry schema configuration")
	createRegistrySchemaCmd.Flags().String("schema-code", "", "Schema code for default PGR registry schema (default: 'pgr')")
	createRegistrySchemaCmd.Flags().String("server", "", "Server URL (overrides config, default: http://localhost:8085)")
	createRegistrySchemaCmd.Flags().String("jwt-token", "", "JWT token for authentication (overrides config)")

	searchRegistrySchemaCmd.Flags().String("schema-code", "", "Schema code to search for (required)")
	searchRegistrySchemaCmd.Flags().String("version", "", "Version of the schema (optional)")
	searchRegistrySchemaCmd.Flags().String("server", "", "Server URL (overrides config, default: http://localhost:8085)")
	searchRegistrySchemaCmd.Flags().String("jwt-token", "", "JWT token for authentication (overrides config)")
	searchRegistrySchemaCmd.MarkFlagRequired("schema-code")

	deleteRegistrySchemaCmd.Flags().String("schema-code", "", "Schema code to delete (required)")
	deleteRegistrySchemaCmd.Flags().String("server", "", "Server URL (overrides config, default: http://localhost:8085)")
	deleteRegistrySchemaCmd.Flags().String("jwt-token", "", "JWT token for authentication (overrides config)")
	deleteRegistrySchemaCmd.MarkFlagRequired("schema-code")
}