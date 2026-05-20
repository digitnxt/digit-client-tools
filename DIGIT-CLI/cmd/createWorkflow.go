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
      type: "INITIAL"
      sla: 86400
      actions:
        - code: "APPLY"
          label: "Apply"
          nextState: "PENDINGFORASSIGNMENT"
    - code: "PENDINGFORASSIGNMENT"
      name: "Pending For Assignment"
      type: "INTERMEDIATE"
      sla: 43200
      actions:
        - code: "ASSIGN"
          label: "Assign"
          nextState: "PENDINGATLME"
        - code: "REJECT"
          label: "Reject"
          nextState: "REJECTED"
    - code: "PENDINGATLME"
      name: "Pending At LME"
      type: "INTERMEDIATE"
      sla: 43200
      actions:
        - code: "RESOLVE"
          label: "Resolve"
          nextState: "RESOLVED"
        - code: "REASSIGN"
          label: "Reassign"
          nextState: "PENDINGFORREASSIGNMENT"
    - code: "PENDINGFORREASSIGNMENT"
      name: "Pending For Reassignment"
      type: "INTERMEDIATE"
      sla: 43200
      actions:
        - code: "REASSIGN"
          label: "Reassign"
          nextState: "PENDINGATLME"
    - code: "REJECTED"
      name: "Rejected"
      type: "INTERMEDIATE"
      sla: 43200
      actions:
        - code: "REOPEN"
          label: "Reopen"
          nextState: "PENDINGFORASSIGNMENT"
        - code: "RATE"
          label: "Rate"
          nextState: "CLOSEDAFTERREJECTION"
    - code: "RESOLVED"
      name: "Resolved"
      type: "INTERMEDIATE"
      sla: 43200
      actions:
        - code: "REOPEN"
          label: "Reopen"
          nextState: "PENDINGFORASSIGNMENT"
        - code: "RATE"
          label: "Rate"
          nextState: "CLOSEDAFTERRESOLUTION"
    - code: "CLOSEDAFTERREJECTION"
      name: "Closed After Rejection"
      type: "TERMINAL_FAILURE"
      sla: 0
      actions: []
    - code: "CLOSEDAFTERRESOLUTION"
      name: "Closed After Resolution"
      type: "TERMINAL_SUCCESS"
      sla: 0
      actions: []`

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
			Code    string `yaml:"code"`
			Name    string `yaml:"name"`
			Type    string `yaml:"type"`
			SLA     int    `yaml:"sla"`
			Actions []struct {
				Code      string `yaml:"code"`
				Label     string `yaml:"label"`
				NextState string `yaml:"nextState"`
			} `yaml:"actions"`
		} `yaml:"states"`
	} `yaml:"workflow"`
}

var searchWorkflowCmd = &cobra.Command{
	Use:   "search-workflow",
	Short: "Get a workflow process definition by code",
	Long: `Get a workflow process definition by providing the process code.

Examples:
  digit search-workflow --code MY_WORKFLOW_CODE
  digit search-workflow --code MY_WORKFLOW_CODE --server http://localhost:8085`,
	RunE: func(cmd *cobra.Command, args []string) error {
		processCode, _ := cmd.Flags().GetString("code")
		serverURL, _ := cmd.Flags().GetString("server")
		jwtToken, _ := cmd.Flags().GetString("jwt-token")

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

		tenantID, err := jwt.ExtractTenantID(jwtToken)
		if err != nil {
			return fmt.Errorf("failed to extract tenant ID from JWT token: %w", err)
		}

		responseBody, err := digit.SearchProcessDefinition(serverURL, jwtToken, tenantID, processCode)
		if err != nil {
			return errors.HandleAPIError(err)
		}

		var jsonResponse interface{}
		if err := json.Unmarshal([]byte(responseBody), &jsonResponse); err == nil {
			prettyJSON, err := json.MarshalIndent(jsonResponse, "", "  ")
			if err == nil {
				fmt.Println(string(prettyJSON))
				return nil
			}
		}
		fmt.Println(responseBody)
		return nil
	},
}

var createWorkflowCmd = &cobra.Command{
	Use:   "create-workflow",
	Short: "Create a complete workflow from a YAML definition",
	Long: `Create a complete workflow (process, states, and actions) from a YAML file or using the built-in default.

Examples:
  digit create-workflow --file workflow.yaml
  digit create-workflow --default --code MY_CODE
  digit create-workflow --file workflow.yaml --server http://localhost:9090`,
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath, _ := cmd.Flags().GetString("file")
		useDefault, _ := cmd.Flags().GetBool("default")
		code, _ := cmd.Flags().GetString("code")
		serverURL, _ := cmd.Flags().GetString("server")
		jwtToken, _ := cmd.Flags().GetString("jwt-token")

		if !useDefault && filePath == "" {
			return fmt.Errorf("either --file or --default flag is required")
		}
		if useDefault && filePath != "" {
			return fmt.Errorf("cannot use both --file and --default flags together")
		}
		if useDefault && code == "" {
			return fmt.Errorf("--code flag is required when using --default")
		}

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

		tenantID, err := jwt.ExtractTenantID(jwtToken)
		if err != nil {
			return fmt.Errorf("failed to extract tenant ID from JWT token: %w", err)
		}

		var yamlData []byte
		if useDefault {
			yamlData = []byte(strings.Replace(defaultWorkflowYAML, "DEFAULT_CODE", code, 1))
			fmt.Printf("Using default workflow configuration with code: %s\n", code)
		} else {
			yamlData, err = os.ReadFile(filePath)
			if err != nil {
				return fmt.Errorf("failed to read YAML file: %w", err)
			}
		}

		var workflowDef WorkflowDefinition
		if err := yaml.Unmarshal(yamlData, &workflowDef); err != nil {
			return fmt.Errorf("failed to parse YAML: %w", err)
		}

		input := digit.ProcessDefinitionInput{
			Code:        workflowDef.Workflow.Process.Code,
			Name:        workflowDef.Workflow.Process.Name,
			Description: workflowDef.Workflow.Process.Description,
			Version:     workflowDef.Workflow.Process.Version,
			SLA:         workflowDef.Workflow.Process.SLA,
		}

		for _, s := range workflowDef.Workflow.States {
			state := digit.StateInput{
				Code:    s.Code,
				Name:    s.Name,
				Type:    s.Type,
				SLA:     s.SLA,
				Actions: []digit.ActionInput{},
			}
			for _, a := range s.Actions {
				state.Actions = append(state.Actions, digit.ActionInput{
					Code:      a.Code,
					Label:     a.Label,
					NextState: a.NextState,
				})
			}
			input.States = append(input.States, state)
		}

		fmt.Println("Creating workflow...")
		responseBody, err := digit.CreateProcessDefinition(serverURL, jwtToken, tenantID, input)
		if err != nil {
			return errors.HandleAPIError(err)
		}

		fmt.Println("Workflow created successfully!")
		var jsonResponse interface{}
		if err := json.Unmarshal([]byte(responseBody), &jsonResponse); err == nil {
			prettyJSON, err := json.MarshalIndent(jsonResponse, "", "  ")
			if err == nil {
				fmt.Println(string(prettyJSON))
				return nil
			}
		}
		fmt.Println(responseBody)
		return nil
	},
}

var deleteProcessCmd = &cobra.Command{
	Use:   "delete-workflow",
	Short: "Delete a workflow process definition",
	Long: `Delete a workflow process definition by its process code.

Examples:
  digit delete-workflow --code TRADE_LICENSE_TEST16
  digit delete-workflow --code TRADE_LICENSE_TEST16 --server https://digit-lts.digit.org`,
	RunE: func(cmd *cobra.Command, args []string) error {
		code, _ := cmd.Flags().GetString("code")
		serverURL, _ := cmd.Flags().GetString("server")
		jwtToken, _ := cmd.Flags().GetString("jwt-token")

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

		tenantID, err := jwt.ExtractTenantID(jwtToken)
		if err != nil {
			return fmt.Errorf("failed to extract tenant ID from JWT token: %w", err)
		}

		responseBody, err := digit.DeleteProcess(serverURL, jwtToken, tenantID, code)
		if err != nil {
			return errors.HandleAPIError(err)
		}

		if responseBody != "" {
			var jsonResponse interface{}
			if err := json.Unmarshal([]byte(responseBody), &jsonResponse); err == nil {
				prettyJSON, err := json.MarshalIndent(jsonResponse, "", "  ")
				if err == nil {
					fmt.Println(string(prettyJSON))
					return nil
				}
			}
			fmt.Println(responseBody)
		} else {
			fmt.Println("Process deleted successfully")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(searchWorkflowCmd)
	rootCmd.AddCommand(createWorkflowCmd)
	rootCmd.AddCommand(deleteProcessCmd)

	searchWorkflowCmd.Flags().String("code", "", "Process code (required)")
	searchWorkflowCmd.Flags().String("server", "", "Server URL (overrides config)")
	searchWorkflowCmd.Flags().String("jwt-token", "", "JWT token (overrides config)")
	searchWorkflowCmd.MarkFlagRequired("code")

	createWorkflowCmd.Flags().String("file", "", "Path to YAML workflow definition file")
	createWorkflowCmd.Flags().Bool("default", false, "Use built-in default workflow (requires --code)")
	createWorkflowCmd.Flags().String("code", "", "Process code for default workflow")
	createWorkflowCmd.Flags().String("server", "", "Server URL (overrides config)")
	createWorkflowCmd.Flags().String("jwt-token", "", "JWT token (overrides config)")

	deleteProcessCmd.Flags().String("code", "", "Process code to delete (required)")
	deleteProcessCmd.Flags().String("server", "", "Server URL (overrides config)")
	deleteProcessCmd.Flags().String("jwt-token", "", "JWT token (overrides config)")
	deleteProcessCmd.MarkFlagRequired("code")
}
