package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"digit-cli/pkg/config"
	"digit-cli/pkg/jwt"
	"github.com/digitnxt/digit3/code/libraries/digit-library/digit"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// defaultWorkflowYAML contains the embedded default workflow configuration
const defaultWorkflowYAML = `workflow:
  process:
    name: "Application Processing Workflow"
    code: "DEFAULT_CODE"
    description: "A complete workflow for application processing"
    version: "1.0"
    sla: 86400
  states:
    - code: "INIT"
      name: "Init"
      isInitial: true
      isParallel: false
      isJoin: false
      sla: 86400
    - code: "PENDINGFORASSIGNMENT"
      name: "Pendingforassignment"
      isInitial: false
      isParallel: false
      isJoin: false
      sla: 43200
    - code: "PENDINGATLME"
      name: "Pendingatlme"
      isInitial: false
      isParallel: false
      isJoin: false
      sla: 43200
    - code: "PENDINGFORREASSIGNMENT"
      name: "Pendingforreassignment"
      isInitial: false
      isParallel: false
      isJoin: false
      sla: 43200
    - code: "REJECTED"
      name: "Rejected"
      isInitial: false
      isParallel: false
      isJoin: false
      sla: 43200
    - code: "RESOLVED"
      name: "Resolved"
      isInitial: false
      isParallel: false
      isJoin: false
      sla: 43200
    - code: "CLOSEDAFTERREJECTION"
      name: "Closedafterrejection"
      isInitial: false
      isParallel: false
      isJoin: false
      sla: 0
    - code: "CLOSEDAFTERRESOLUTION"
      name: "Closedafterresolution"
      isInitial: false
      isParallel: false
      isJoin: false
      sla: 0
  actions:
    - name: "APPLY"
      currentState: "INIT"
      nextState: "PENDINGFORASSIGNMENT"
      attributeValidation:
        attributes:
          roles: ["CITIZEN", "CSR"]
    - name: "ASSIGN"
      currentState: "PENDINGFORASSIGNMENT"
      nextState: "PENDINGATLME"
      attributeValidation:
        attributes:
          roles: ["GRO"]
    - name: "REJECT"
      currentState: "PENDINGFORASSIGNMENT"
      nextState: "REJECTED"
      attributeValidation:
        attributes:
          roles: ["GRO"]
    - name: "REASSIGN"
      currentState: "PENDINGFORREASSIGNMENT"
      nextState: "PENDINGATLME"
      attributeValidation:
        attributes:
          roles: ["GRO"]
    - name: "REASSIGN"
      currentState: "PENDINGATLME"
      nextState: "PENDINGFORREASSIGNMENT"
      attributeValidation:
        attributes:
          roles: ["LME"]
    - name: "RESOLVE"
      currentState: "PENDINGATLME"
      nextState: "RESOLVED"
      attributeValidation:
        attributes:
          roles: ["LME"]
    - name: "REOPEN"
      currentState: "REJECTED"
      nextState: "PENDINGFORASSIGNMENT"
      attributeValidation:
        attributes:
          roles: ["CITIZEN", "CSR"]
    - name: "RATE"
      currentState: "REJECTED"
      nextState: "CLOSEDAFTERREJECTION"
      attributeValidation:
        attributes:
          roles: ["CITIZEN", "CSR"]
    - name: "REOPEN"
      currentState: "RESOLVED"
      nextState: "PENDINGFORASSIGNMENT"
      attributeValidation:
        attributes:
          roles: ["CITIZEN", "CSR"]
    - name: "RATE"
      currentState: "RESOLVED"
      nextState: "CLOSEDAFTERRESOLUTION"
      attributeValidation:
        attributes:
          roles: ["CITIZEN", "CSR"]`

// WorkflowDefinition represents the YAML structure for workflow definition
type WorkflowDefinition struct {
	Workflow struct {
		Process struct {
			Name        string `yaml:"name"`
			Code        string `yaml:"code"`
			Description string `yaml:"description"`
			Version     string `yaml:"version"`
			SLA         int    `yaml:"sla"`
		} `yaml:"process"`
		States []struct {
			Code       string `yaml:"code"`
			Name       string `yaml:"name"`
			IsInitial  bool   `yaml:"isInitial"`
			IsParallel bool   `yaml:"isParallel"`
			IsJoin     bool   `yaml:"isJoin"`
			SLA        int    `yaml:"sla"`
		} `yaml:"states"`
		Actions []struct {
			Name         string   `yaml:"name"`
			CurrentState string   `yaml:"currentState"`
			NextState    string   `yaml:"nextState"`
			Roles        []string `yaml:"roles"`
			AttributeValidation struct {
				Attributes struct {
					Roles []string `yaml:"roles"`
				} `yaml:"attributes"`
				AssigneeCheck bool `yaml:"assigneeCheck"`
			} `yaml:"attributeValidation"`
		} `yaml:"actions"`
	} `yaml:"workflow"`
}

// createProcessCmd represents the create-process command
var createProcessCmd = &cobra.Command{
	Use:   "create-process",
	Short: "Create a new workflow process",
	Long: `Create a new workflow process with the specified parameters.
	
Examples:
  # Basic process creation
  digit create-process --name "Hello" --code "{{GenProcessId}}" --description "A test process for API validation" --version "1.0" --sla 86400
  
  # With server override
  digit create-process --name "MyWorkflow" --code "WF001" --description "Custom workflow" --version "2.0" --sla 3600 --server http://localhost:9090`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get flag values
		name, _ := cmd.Flags().GetString("name")
		code, _ := cmd.Flags().GetString("code")
		description, _ := cmd.Flags().GetString("description")
		version, _ := cmd.Flags().GetString("version")
		slaStr, _ := cmd.Flags().GetString("sla")
		serverURL, _ := cmd.Flags().GetString("server")
		jwtToken, _ := cmd.Flags().GetString("jwt-token")
		
		// Validate required flags
		if name == "" {
			return fmt.Errorf("--name flag is required")
		}
		if code == "" {
			return fmt.Errorf("--code flag is required")
		}
		if description == "" {
			return fmt.Errorf("--description flag is required")
		}
		if version == "" {
			return fmt.Errorf("--version flag is required")
		}
		if slaStr == "" {
			return fmt.Errorf("--sla flag is required")
		}
		
		// Convert SLA to integer
		sla, err := strconv.Atoi(slaStr)
		if err != nil {
			return fmt.Errorf("invalid SLA value: %w", err)
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
		
		// Call the digit library to create process
		responseBody, err := digit.CreateProcess(serverURL, jwtToken, tenantID, name, code, description, version, sla)
		if err != nil {
			return fmt.Errorf("failed to create process: %w", err)
		}
		
		// Print response body
		fmt.Printf("Process creation response:\n")
		
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

// searchProcessDefinitionCmd represents the search-process-definition command
var searchProcessDefinitionCmd = &cobra.Command{
	Use:   "search-process-definition",
	Short: "Search for a workflow process definition by ID",
	Long: `Search for a workflow process definition by providing the process ID.
	
Examples:
  # Search process definition by ID
  digit search-process-definition --id dd2e8cf5-a53e-44b5-82b9-490ac73c50dd
  
  # With server override
  digit search-process-definition --id dd2e8cf5-a53e-44b5-82b9-490ac73c50dd --server http://localhost:8085`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get flag values
		processID, _ := cmd.Flags().GetString("id")
		serverURL, _ := cmd.Flags().GetString("server")
		jwtToken, _ := cmd.Flags().GetString("jwt-token")
		
		// Validate required flags
		if processID == "" {
			return fmt.Errorf("--id flag is required")
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
		
		// Call the digit library to search process definition
		responseBody, err := digit.SearchProcessDefinition(serverURL, jwtToken, tenantID, processID)
		if err != nil {
			return fmt.Errorf("failed to search process definition: %w", err)
		}
		
		// Print response body
		fmt.Printf("Process definition response:\n")
		
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

// createWorkflowCmd represents the create-workflow command
var createWorkflowCmd = &cobra.Command{
	Use:   "create-workflow",
	Short: "Create a complete workflow from YAML definition",
	Long: `Create a complete workflow (process, states, and actions) from a YAML file definition or using default configuration.
	
Examples:
  # Create workflow from YAML file
  digit create-workflow --file workflow.yaml
  
  # Create workflow using default configuration with custom code
  digit create-workflow --default --code MY_CUSTOM_CODE
  
  # With server override
  digit create-workflow --file workflow.yaml --server http://localhost:9090
  digit create-workflow --default --code MY_CODE --server http://localhost:9090`,
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
		if useDefault && code == "" {
			return fmt.Errorf("--code flag is required when using --default")
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
		
		// Get YAML data - either from file or default configuration
		var yamlData []byte
		if useDefault {
			// Use embedded default configuration and replace the code
			yamlContent := strings.Replace(defaultWorkflowYAML, "DEFAULT_CODE", code, 1)
			yamlData = []byte(yamlContent)
			fmt.Printf("Using default workflow configuration with code: %s\n", code)
		} else {
			// Read YAML file
			var err error
			yamlData, err = os.ReadFile(filePath)
			if err != nil {
				return fmt.Errorf("failed to read YAML file: %w", err)
			}
		}
		
		// Parse YAML
		var workflowDef WorkflowDefinition
		if err := yaml.Unmarshal(yamlData, &workflowDef); err != nil {
			return fmt.Errorf("failed to parse YAML: %w", err)
		}
		
		// Step 1: Create Process
		fmt.Println("Creating workflow process...")
		processResponse, err := digit.CreateProcess(
			serverURL,
			jwtToken,
			tenantID,
			workflowDef.Workflow.Process.Name,
			workflowDef.Workflow.Process.Code,
			workflowDef.Workflow.Process.Description,
			workflowDef.Workflow.Process.Version,
			workflowDef.Workflow.Process.SLA,
		)
		if err != nil {
			return fmt.Errorf("failed to create process: %w", err)
		}
		
		// Parse process response to get process ID
		var processResult map[string]interface{}
		if err := json.Unmarshal([]byte(processResponse), &processResult); err != nil {
			return fmt.Errorf("failed to parse process response: %w", err)
		}
		
		processID, ok := processResult["id"].(string)
		if !ok {
			return fmt.Errorf("failed to extract process ID from response")
		}
		
		fmt.Printf("✓ Process created with ID: %s\n", processID)
		
		// Step 2: Create States
		fmt.Println("Creating workflow states...")
		stateCodeToID := make(map[string]string) // Map state codes to their IDs
		
		for _, state := range workflowDef.Workflow.States {
			stateResponse, err := digit.CreateState(
				serverURL,
				jwtToken,
				tenantID,
				processID,
				state.Code,
				state.Name,
				state.IsInitial,
				state.IsParallel,
				state.IsJoin,
				state.SLA,
			)
			if err != nil {
				return fmt.Errorf("failed to create state %s: %w", state.Code, err)
			}
			
			// Parse state response to get state ID
			var stateResult map[string]interface{}
			if err := json.Unmarshal([]byte(stateResponse), &stateResult); err != nil {
				return fmt.Errorf("failed to parse state response for %s: %w", state.Code, err)
			}
			
			stateID, ok := stateResult["id"].(string)
			if !ok {
				return fmt.Errorf("failed to extract state ID from response for %s", state.Code)
			}
			
			// Store the mapping of state code to state ID
			stateCodeToID[state.Code] = stateID
			
			fmt.Printf("✓ State created: %s (%s) - ID: %s\n", state.Name, state.Code, stateID)
		}
		
		// Step 3: Create Actions
		fmt.Println("Creating workflow actions...")
		for _, action := range workflowDef.Workflow.Actions {
			// Get the state ID for the current state
			currentStateID, exists := stateCodeToID[action.CurrentState]
			if !exists {
				return fmt.Errorf("state ID not found for current state: %s", action.CurrentState)
			}
			
			// Get the state ID for the next state
			nextStateID, exists := stateCodeToID[action.NextState]
			if !exists {
				return fmt.Errorf("state ID not found for next state: %s", action.NextState)
			}
			
			actionResponse, err := digit.CreateAction(
				serverURL,
				jwtToken,
				tenantID,
				currentStateID, // Use the actual state ID as path parameter
				action.Name,
				nextStateID, // Use the actual next state UUID
				action.AttributeValidation.Attributes.Roles,
				action.AttributeValidation.AssigneeCheck,
			)
			if err != nil {
				return fmt.Errorf("failed to create action %s: %w", action.Name, err)
			}
			fmt.Printf("✓ Action created: %s (%s → %s) using state ID: %s\n", action.Name, action.CurrentState, action.NextState, currentStateID)
			_ = actionResponse // We don't need the response for now
		}
		
		fmt.Println("\n🎉 Workflow created successfully!")
		fmt.Printf("Process ID: %s\n", processID)
		fmt.Printf("States: %d\n", len(workflowDef.Workflow.States))
		fmt.Printf("Actions: %d\n", len(workflowDef.Workflow.Actions))
		
		return nil
	},
}

// deleteProcessCmd represents the delete-process command
var deleteProcessCmd = &cobra.Command{
	Use:   "delete-process",
	Short: "Delete an existing workflow process",
	Long: `Delete a workflow process by process code.

This command will delete an existing workflow process from the system.

Example:
  digit delete-process --code TRADE_LICENSE_TEST16
  digit delete-process --code TRADE_LICENSE_TEST16 --server https://digit-lts.digit.org`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get flag values
		code, _ := cmd.Flags().GetString("code")
		serverURL, _ := cmd.Flags().GetString("server")
		jwtToken, _ := cmd.Flags().GetString("jwt-token")

		// Validate required flags
		if code == "" {
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

		// Call the digit library to delete process
		responseBody, err := digit.DeleteProcess(serverURL, jwtToken, tenantID, code)
		if err != nil {
			return fmt.Errorf("failed to delete process: %w", err)
		}

		// Print response body
		fmt.Printf("Process deletion response:\n")

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
			fmt.Println("Process deleted successfully")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(createProcessCmd)
	rootCmd.AddCommand(searchProcessDefinitionCmd)
	rootCmd.AddCommand(createWorkflowCmd)
	rootCmd.AddCommand(deleteProcessCmd)
	
	// Add flags for create-process command
	createProcessCmd.Flags().String("name", "", "Name of the workflow process (required)")
	createProcessCmd.Flags().String("code", "", "Code of the workflow process (required)")
	createProcessCmd.Flags().String("description", "", "Description of the workflow process (required)")
	createProcessCmd.Flags().String("version", "", "Version of the workflow process (required)")
	createProcessCmd.Flags().String("sla", "", "SLA in seconds for the workflow process (required)")
	createProcessCmd.Flags().String("server", "", "Server URL (overrides config)")
	createProcessCmd.Flags().String("jwt-token", "", "JWT token for authentication (overrides config)")
	
	// Mark required flags for create-process
	createProcessCmd.MarkFlagRequired("name")
	createProcessCmd.MarkFlagRequired("code")
	createProcessCmd.MarkFlagRequired("description")
	createProcessCmd.MarkFlagRequired("version")
	createProcessCmd.MarkFlagRequired("sla")
	
	// Add flags for search-process-definition command
	searchProcessDefinitionCmd.Flags().String("id", "", "Process ID to search definition for (required)")
	searchProcessDefinitionCmd.Flags().String("server", "", "Server URL (overrides config)")
	searchProcessDefinitionCmd.Flags().String("jwt-token", "", "JWT token for authentication (overrides config)")
	
	// Mark required flags for search-process-definition
	searchProcessDefinitionCmd.MarkFlagRequired("id")
	
	// Add flags for create-workflow command
	createWorkflowCmd.Flags().String("file", "", "Path to YAML file containing workflow definition")
	createWorkflowCmd.Flags().Bool("default", false, "Use default workflow configuration (requires --code)")
	createWorkflowCmd.Flags().String("code", "", "Process code to use with default configuration (required when using --default)")
	createWorkflowCmd.Flags().String("server", "", "Server URL (overrides config)")
	createWorkflowCmd.Flags().String("jwt-token", "", "JWT token for authentication (overrides config)")
	
	// Note: Required flags are validated conditionally in the command logic
	
	// Add flags for delete-process command
	deleteProcessCmd.Flags().String("code", "", "Process code to delete (required)")
	deleteProcessCmd.Flags().String("server", "", "Server URL (overrides config)")
	deleteProcessCmd.Flags().String("jwt-token", "", "JWT token for authentication (overrides config)")

	deleteProcessCmd.MarkFlagRequired("code")
}
