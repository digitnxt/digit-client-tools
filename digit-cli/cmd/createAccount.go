package cmd

import (
	"encoding/json"
	"fmt"

	"digit-cli/pkg/config"
	"github.com/digitnxt/digit3/code/libraries/digit-library/digit"
	"github.com/spf13/cobra"
)


// createAccountCmd represents the create-account command
var createAccountCmd = &cobra.Command{
	Use:   "create-account",
	Short: "Create a new account in DIGIT services",
	Long: `Create a new account in DIGIT services with the specified name, email, and status.
	
Examples:
  digit create-account --name kongnew1 --email test@example.com
  digit create-account --name kongnew1 --email test@example.com --active=false
  digit create-account --name kongnew1 --email test@example.com`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get flag values
		name, _ := cmd.Flags().GetString("name")
		email, _ := cmd.Flags().GetString("email")
		active, _ := cmd.Flags().GetBool("active")
		clientID, _ := cmd.Flags().GetString("client-id")
		serverURL, _ := cmd.Flags().GetString("server")
		
		// Validate required flags
		if name == "" {
			return fmt.Errorf("--name flag is required")
		}
		if email == "" {
			return fmt.Errorf("--email flag is required")
		}
		
		// Load server URL from config if not provided as flag
		if serverURL == "" {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}
			serverURL = cfg.GetServer()
			
			if serverURL == "" {
				return fmt.Errorf("server URL not configured. Use 'digit config set --server <url>' or provide --server flag")
			}
		}
		
		// Call the digit library to create account
		responseBody, err := digit.CreateAccount(serverURL, clientID, name, email, active)
		if err != nil {
			return fmt.Errorf("failed to create account: %w", err)
		}
		
		// Print response body
		fmt.Printf("Account creation response:\n")
		
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
	rootCmd.AddCommand(createAccountCmd)
	
	// Add flags for create-account command
	createAccountCmd.Flags().String("name", "", "Name of the tenant (required)")
	createAccountCmd.Flags().String("email", "", "Email of the tenant (required)")
	createAccountCmd.Flags().Bool("active", true, "Whether the tenant is active (default: true)")
	createAccountCmd.Flags().String("client-id", "test-client", "Client ID for the request (default: test-client)")
	createAccountCmd.Flags().String("server", "", "Server URL (overrides config)")
	
	// Mark required flags
	createAccountCmd.MarkFlagRequired("name")
	createAccountCmd.MarkFlagRequired("email")
}
