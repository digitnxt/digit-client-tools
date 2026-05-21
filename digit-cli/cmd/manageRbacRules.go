package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/digitnxt/digit-client-tools/client-libraries/digit-library/digit"
	"github.com/spf13/cobra"
)

// ============================================================================
// LIST RBAC RULES
//
// Usage:
//   digit list-rbac-rules                      → all rules for your tenant
//   digit list-rbac-rules --role SUPERUSER      → filtered by role
//   digit list-rbac-rules --page 0 --size 20    → paginated
//
// Tenant ID is derived from your JWT token (same as create).
// The API returns: { rules: [...], page: 0, size: 50, total: 123 }
// ============================================================================
var listRbacRulesCmd = &cobra.Command{
	Use:   "list-rbac-rules",
	Short: "List RBAC access control rules for your tenant",
	Long: `List RBAC rules configured for your tenant.

Examples:
  digit list-rbac-rules
  digit list-rbac-rules --role SUPERUSER
  digit list-rbac-rules --page 0 --size 20`,

	RunE: func(cmd *cobra.Command, args []string) error {
		// Resolve connection details (same helper from createRbacRule.go)
		serverURL, jwtToken, tenantID, userID, err := resolveServerAndTenant(cmd)
		if err != nil {
			return err
		}

		// Get filter/pagination flags
		roleName, _ := cmd.Flags().GetString("role")
		page, _ := cmd.Flags().GetInt("page")
		size, _ := cmd.Flags().GetInt("size")

		// Call the library function
		responseBody, err := digit.ListRbacRules(serverURL, jwtToken, tenantID, userID, roleName, page, size)
		if err != nil {
			return fmt.Errorf("failed to list RBAC rules: %w", err)
		}

		// Pretty-print the JSON response
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

// ============================================================================
// GET RBAC RULE BY ID
//
// Usage:
//   digit get-rbac-rule --id <uuid>
//
// Fetches a single rule's full details. Useful after list to inspect
// constraints or other fields not shown in the list summary.
// ============================================================================
var getRbacRuleCmd = &cobra.Command{
	Use:   "get-rbac-rule",
	Short: "Get a single RBAC rule by ID",
	Long: `Retrieve the full details of an RBAC rule.

Examples:
  digit get-rbac-rule --id 550e8400-e29b-41d4-a716-446655440000`,

	RunE: func(cmd *cobra.Command, args []string) error {
		serverURL, jwtToken, tenantID, userID, err := resolveServerAndTenant(cmd)
		if err != nil {
			return err
		}

		ruleID, _ := cmd.Flags().GetString("id")
		if ruleID == "" {
			return fmt.Errorf("--id flag is required")
		}

		responseBody, err := digit.GetRbacRule(serverURL, jwtToken, tenantID, userID, ruleID)
		if err != nil {
			return fmt.Errorf("failed to get RBAC rule: %w", err)
		}

		// Pretty-print
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

// ============================================================================
// DELETE RBAC RULE
//
// Usage:
//   digit delete-rbac-rule --id <uuid>
//
// Deletes a single rule by UUID. The API returns 204 No Content on success.
// ============================================================================
var deleteRbacRuleCmd = &cobra.Command{
	Use:   "delete-rbac-rule",
	Short: "Delete an RBAC rule by ID",
	Long: `Delete an RBAC access control rule.

Examples:
  digit delete-rbac-rule --id 550e8400-e29b-41d4-a716-446655440000`,

	RunE: func(cmd *cobra.Command, args []string) error {
		serverURL, jwtToken, tenantID, userID, err := resolveServerAndTenant(cmd)
		if err != nil {
			return err
		}

		ruleID, _ := cmd.Flags().GetString("id")
		if ruleID == "" {
			return fmt.Errorf("--id flag is required")
		}

		err = digit.DeleteRbacRule(serverURL, jwtToken, tenantID, userID, ruleID)
		if err != nil {
			return fmt.Errorf("failed to delete RBAC rule: %w", err)
		}

		fmt.Printf("Rule %s deleted successfully.\n", ruleID)
		return nil
	},
}

// ============================================================================
// DELETE ALL RBAC RULES (for a tenant)
//
// Usage:
//   digit delete-all-rbac-rules
//   digit delete-all-rbac-rules --role SUPERUSER   → only rules for this role
//
// This is the CLI equivalent of the delete_all_rules() function in
// add_manual_rules.py. It lists all rules, then deletes them one by one.
// Useful for resetting a tenant's access control during development/testing.
// ============================================================================
var deleteAllRbacRulesCmd = &cobra.Command{
	Use:   "delete-all-rbac-rules",
	Short: "Delete all RBAC rules for your tenant (use with caution)",
	Long: `Delete all RBAC rules for your tenant. This is a destructive operation.

Examples:
  digit delete-all-rbac-rules
  digit delete-all-rbac-rules --role SUPERUSER`,

	RunE: func(cmd *cobra.Command, args []string) error {
		serverURL, jwtToken, tenantID, userID, err := resolveServerAndTenant(cmd)
		if err != nil {
			return err
		}

		roleName, _ := cmd.Flags().GetString("role")

		// Step 1: List all rules to get their IDs
		// We use a large page size to get everything in one call.
		responseBody, err := digit.ListRbacRules(serverURL, jwtToken, tenantID, userID, roleName, 0, 1000)
		if err != nil {
			return fmt.Errorf("failed to list rules: %w", err)
		}

		// Parse the response to extract rule IDs
		var listResp struct {
			Rules []struct {
				ID          string `json:"id"`
				Description string `json:"description"`
			} `json:"rules"`
			Total int `json:"total"`
		}
		if err := json.Unmarshal([]byte(responseBody), &listResp); err != nil {
			return fmt.Errorf("failed to parse rules list: %w", err)
		}

		if len(listResp.Rules) == 0 {
			fmt.Println("No rules found to delete.")
			return nil
		}

		fmt.Printf("Found %d rules. Deleting...\n", len(listResp.Rules))

		// Step 2: Delete each rule
		successCount := 0
		for _, rule := range listResp.Rules {
			if err := digit.DeleteRbacRule(serverURL, jwtToken, tenantID, userID, rule.ID); err != nil {
				fmt.Printf("  ✗ Failed to delete %s (%s): %v\n", rule.ID, rule.Description, err)
				continue
			}
			fmt.Printf("  ✓ Deleted %s (%s)\n", rule.ID, rule.Description)
			successCount++
		}

		fmt.Printf("\nDone: %d/%d rules deleted.\n", successCount, len(listResp.Rules))
		return nil
	},
}

// ============================================================================
// Flag registration
// ============================================================================
func init() {
	// ── list-rbac-rules ──
	rootCmd.AddCommand(listRbacRulesCmd)
	listRbacRulesCmd.Flags().String("role", "", "Filter rules by role name")
	listRbacRulesCmd.Flags().Int("page", 0, "Page number (default: 0)")
	listRbacRulesCmd.Flags().Int("size", 50, "Page size (default: 50, max: 100)")
	listRbacRulesCmd.Flags().String("server", "", "Server URL (overrides config)")
	listRbacRulesCmd.Flags().String("jwt-token", "", "JWT token (overrides config)")

	// ── get-rbac-rule ──
	rootCmd.AddCommand(getRbacRuleCmd)
	getRbacRuleCmd.Flags().String("id", "", "Rule ID (UUID) to retrieve")
	getRbacRuleCmd.MarkFlagRequired("id")
	getRbacRuleCmd.Flags().String("server", "", "Server URL (overrides config)")
	getRbacRuleCmd.Flags().String("jwt-token", "", "JWT token (overrides config)")

	// ── delete-rbac-rule ──
	rootCmd.AddCommand(deleteRbacRuleCmd)
	deleteRbacRuleCmd.Flags().String("id", "", "Rule ID (UUID) to delete")
	deleteRbacRuleCmd.MarkFlagRequired("id")
	deleteRbacRuleCmd.Flags().String("server", "", "Server URL (overrides config)")
	deleteRbacRuleCmd.Flags().String("jwt-token", "", "JWT token (overrides config)")

	// ── delete-all-rbac-rules ──
	rootCmd.AddCommand(deleteAllRbacRulesCmd)
	deleteAllRbacRulesCmd.Flags().String("role", "", "Only delete rules matching this role")
	deleteAllRbacRulesCmd.Flags().String("server", "", "Server URL (overrides config)")
	deleteAllRbacRulesCmd.Flags().String("jwt-token", "", "JWT token (overrides config)")
}
