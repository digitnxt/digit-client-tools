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

// ============================================================================
// YAML structs
// ============================================================================
type RbacRuleYAML struct {
	RoleNames   []string          `yaml:"roleNames"`
	HTTPMethod  string            `yaml:"httpMethod"`
	Path        string            `yaml:"path"`
	Effect      *string           `yaml:"effect"`
	Priority    *int              `yaml:"priority"`
	Enabled     *bool             `yaml:"enabled"`
	Description string            `yaml:"description"`
	CRUD        bool              `yaml:"crud"`
	Constraints []digit.Constraint `yaml:"constraints"`
}

type RbacRulesFileYAML struct {
	Defaults *RbacRuleYAML  `yaml:"defaults"`
	Rules    []RbacRuleYAML `yaml:"rules"`
}

var (
	defaultEffect   = "ALLOW"
	defaultPriority = 100
	defaultEnabled  = true
)

var crudMethods = []string{"GET", "POST", "PUT", "DELETE", "PATCH"}

// ============================================================================
// Helper: resolve server URL, JWT token, and tenant ID
// ============================================================================
func resolveServerAndTenant(cmd *cobra.Command) (serverURL, jwtToken, tenantID, userID string, err error) {
	serverURL, _ = cmd.Flags().GetString("server")
	jwtToken, _ = cmd.Flags().GetString("jwt-token")

	if serverURL == "" || jwtToken == "" {
		cfg, cfgErr := config.Load()
		if cfgErr != nil {
			return "", "", "", "", fmt.Errorf("failed to load config: %w", cfgErr)
		}

		if serverURL == "" {
			serverURL = cfg.GetServer()
		}
		if jwtToken == "" {
			jwtToken = cfg.GetJWTToken()
		}
	}

	tenantID, err = jwt.ExtractTenantID(jwtToken)
	if err != nil {
		return
	}
	userID, err = jwt.ExtractClientID(jwtToken)
	return
}

// ============================================================================
// Constraint parsing helper (CLI → library model)
// Format: attribute:operator:value
// Example: boundary:EQ:KA.BLR.WARD1
// ============================================================================
func parseConstraints(strs []string) ([]digit.Constraint, error) {
	var out []digit.Constraint

	for _, s := range strs {
		parts := strings.Split(s, ":")
		if len(parts) < 4 {
			return nil, fmt.Errorf("invalid constraint format: %s (expected type:key:op:value)", s)
		}

		out = append(out, digit.Constraint{
			Type:     parts[0],
			Key:      parts[1],
			Operator: parts[2],
			Value:    parts[3],
		})
	}

	return out, nil
}

// ============================================================================
// Helper: send one rule
// ============================================================================
func sendRbacRule(serverURL, jwtToken, tenantID, userID string, rule digit.RbacRulePayload) error {
	responseBody, err := digit.CreateRbacRule(serverURL, jwtToken, tenantID, userID, rule)
	if err != nil {
		return err
	}

	var jsonResponse interface{}
	if err := json.Unmarshal([]byte(responseBody), &jsonResponse); err == nil {
		prettyJSON, _ := json.MarshalIndent(jsonResponse, "", "  ")
		fmt.Println(string(prettyJSON))
	} else {
		fmt.Println(responseBody)
	}
	return nil
}

// ============================================================================
// Helper: auto description
// ============================================================================
func autoDescription(roles []string, method, path, effect string) string {
	return fmt.Sprintf("%s %s %s %s", effect, strings.Join(roles, ","), method, path)
}

// ============================================================================
// MODE 1: flags
// ============================================================================
func handleFlagMode(cmd *cobra.Command, serverURL, jwtToken, tenantID, userID string) error {
	rolesStr, _ := cmd.Flags().GetString("roles")
	method, _ := cmd.Flags().GetString("method")
	path, _ := cmd.Flags().GetString("path")
	effect, _ := cmd.Flags().GetString("effect")
	priority, _ := cmd.Flags().GetInt("priority")
	description, _ := cmd.Flags().GetString("description")
	crud, _ := cmd.Flags().GetBool("crud")

	constraintStrs, _ := cmd.Flags().GetStringArray("constraint")
	constraints, err := parseConstraints(constraintStrs)
	if err != nil {
		return err
	}

	if rolesStr == "" || path == "" {
		return fmt.Errorf("--roles and --path required")
	}

	roles := strings.Split(rolesStr, ",")
	methods := []string{strings.ToUpper(method)}
	if crud {
		methods = crudMethods
	}

	for _, m := range methods {
		desc := description
		if desc == "" {
			desc = autoDescription(roles, m, path, effect)
		}

		rule := digit.RbacRulePayload{
			RoleNames:   roles,
			HTTPMethod:  m,
			Path:        path,
			Effect:      effect,
			Priority:    priority,
			Enabled:     true,
			Description: desc,
			Constraints: constraints,
		}

		if err := sendRbacRule(serverURL, jwtToken, tenantID, userID, rule); err != nil {
			return err
		}
	}

	return nil
}

// ============================================================================
// MODE 2: YAML
// ============================================================================
func handleFileMode(filePath, serverURL, jwtToken, tenantID, userID string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var rulesFile RbacRulesFileYAML
	if err := yaml.Unmarshal(data, &rulesFile); err != nil {
		return err
	}

	for _, r := range rulesFile.Rules {
		methods := []string{strings.ToUpper(r.HTTPMethod)}
		if r.CRUD {
			methods = crudMethods
		}

		for _, m := range methods {
			rule := digit.RbacRulePayload{
				RoleNames:   r.RoleNames,
				HTTPMethod:  m,
				Path:        r.Path,
				Effect:      defaultEffect,
				Priority:    defaultPriority,
				Enabled:     defaultEnabled,
				Description: r.Description,
				Constraints: r.Constraints,
			}

			if err := sendRbacRule(serverURL, jwtToken, tenantID, userID, rule); err != nil {
				return err
			}
		}
	}

	return nil
}

// ============================================================================
// Cobra command
// ============================================================================
var createRbacRuleCmd = &cobra.Command{
	Use:   "create-rbac-rule",
	Short: "Create RBAC/JBAC rules",
	RunE: func(cmd *cobra.Command, args []string) error {
		serverURL, jwtToken, tenantID, userID, err := resolveServerAndTenant(cmd)
		if err != nil {
			return err
		}

		filePath, _ := cmd.Flags().GetString("file")
		if filePath != "" {
			return handleFileMode(filePath, serverURL, jwtToken, tenantID, userID)
		}
		return handleFlagMode(cmd, serverURL, jwtToken, tenantID, userID)
	},
}

func init() {
	rootCmd.AddCommand(createRbacRuleCmd)

	createRbacRuleCmd.Flags().String("file", "", "YAML file")
	createRbacRuleCmd.Flags().String("roles", "", "Roles")
	createRbacRuleCmd.Flags().String("method", "", "HTTP method")
	createRbacRuleCmd.Flags().String("path", "", "API path")
	createRbacRuleCmd.Flags().String("effect", "ALLOW", "Effect")
	createRbacRuleCmd.Flags().Int("priority", 100, "Priority")
	createRbacRuleCmd.Flags().String("description", "", "Description")
	createRbacRuleCmd.Flags().Bool("crud", false, "CRUD shortcut")

	createRbacRuleCmd.Flags().StringArray("constraint", []string{}, "attribute:operator:value")

	createRbacRuleCmd.Flags().String("server", "", "Server URL")
	createRbacRuleCmd.Flags().String("jwt-token", "", "JWT token")
}
