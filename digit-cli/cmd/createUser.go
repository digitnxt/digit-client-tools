package cmd

import (
	"encoding/json"
	"fmt"

	"digit-cli/pkg/config"
	"github.com/digitnxt/digit3/code/libraries/digit-library/digit"
	"github.com/spf13/cobra"
)

// createUserCmd represents the create-user command
var createUserCmd = &cobra.Command{
	Use:   "create-user",
	Short: "Create a new user in Keycloak",
	Long: `Create a new user in Keycloak with the specified username, password, email, and account.
	
Examples:
  digit create-user --username johndoe --password mypassword --email john@example.com --account master
  digit create-user --username adminuser --password adminpass --email admin@example.com --account master --server http://localhost:8080
  digit create-user --username johndoe --password mypassword --email john@example.com --account master --jwt-token <token>`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get flag values
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")
		email, _ := cmd.Flags().GetString("email")
		realm, _ := cmd.Flags().GetString("account")
		serverURL, _ := cmd.Flags().GetString("server")
		jwtToken, _ := cmd.Flags().GetString("jwt-token")
		
		// Validate required flags
		if username == "" {
			return fmt.Errorf("--username flag is required")
		}
		if password == "" {
			return fmt.Errorf("--password flag is required")
		}
		if email == "" {
			return fmt.Errorf("--email flag is required")
		}
		
		// Load configuration if server URL, JWT token, or account not provided as flags
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
		
		if realm == "" {
			realm, err = config.GetRealm()
			if err != nil {
				return fmt.Errorf("failed to get realm from config: %w", err)
			}
			if realm == "" {
				return fmt.Errorf("account not configured. Use 'digit config set' or provide --account flag")
			}
		}
		
		// Call the digit library to create user
		responseBody, err := digit.CreateUser(serverURL, jwtToken, realm, username, password, email)
		if err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}
		
		// Print response body
		fmt.Printf("User creation response:\n")
		
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

// resetPasswordCmd represents the reset-password command
var resetPasswordCmd = &cobra.Command{
	Use:   "reset-password",
	Short: "Reset a user's password in Keycloak",
	Long: `Reset a user's password in Keycloak with the specified username and new password.
	
Examples:
  digit reset-password --username johndoe --new-password newpassword123 --account master
  digit reset-password --username johndoe --new-password newpassword123 --account myrealm --server https://keycloak.example.com
  digit reset-password --username johndoe --new-password newpassword123 --account master --jwt-token <token>`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get flag values
		username, _ := cmd.Flags().GetString("username")
		newPassword, _ := cmd.Flags().GetString("new-password")
		realm, _ := cmd.Flags().GetString("account")
		serverURL, _ := cmd.Flags().GetString("server")
		jwtToken, _ := cmd.Flags().GetString("jwt-token")
		
		// Validate required flags
		if username == "" {
			return fmt.Errorf("--username flag is required")
		}
		if newPassword == "" {
			return fmt.Errorf("--new-password flag is required")
		}
		
		// Load configuration if server URL, JWT token, or account not provided as flags
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
		
		if realm == "" {
			realm, err = config.GetRealm()
			if err != nil {
				return fmt.Errorf("failed to get realm from config: %w", err)
			}
			if realm == "" {
				return fmt.Errorf("account not configured. Use 'digit config set' or provide --account flag")
			}
		}
		
		// Call the digit library to reset password
		responseBody, err := digit.ResetPassword(serverURL, jwtToken, realm, username, newPassword)
		if err != nil {
			return fmt.Errorf("failed to reset password: %w", err)
		}
		
		// Print response body
		fmt.Printf("Password reset response:\n")
		
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
			if responseBody == "" {
				fmt.Println("Password reset successful (no response body)")
			} else {
				fmt.Println(responseBody)
			}
		}
		
		return nil
	},
}

// deleteUserCmd represents the delete-user command
var deleteUserCmd = &cobra.Command{
	Use:   "delete-user",
	Short: "Delete a user from Keycloak",
	Long: `Delete a user from Keycloak by username.

This command will:
1. Find the user by username in the specified account
2. Delete the user from Keycloak

Example:
  digit delete-user --username johndoe --account master
  digit delete-user --username johndoe --account myrealm --server https://keycloak.example.com --jwt-token <token>`,
	RunE: func(cmd *cobra.Command, args []string) error {
		username, _ := cmd.Flags().GetString("username")
		realm, _ := cmd.Flags().GetString("account")
		serverURL, _ := cmd.Flags().GetString("server")
		jwtToken, _ := cmd.Flags().GetString("jwt-token")

		// Load configuration
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		// Get server URL from config if not provided
		if serverURL == "" {
			serverURL = cfg.Server
		}

		// Get JWT token from config if not provided
		if jwtToken == "" {
			jwtToken = cfg.JWTToken
		}

		// Get realm from config if not provided
		if realm == "" {
			realm, err = config.GetRealm()
			if err != nil {
				return fmt.Errorf("failed to get realm from config: %w", err)
			}
		}

		// Validate required parameters
		if username == "" {
			return fmt.Errorf("username is required")
		}
		if realm == "" {
			return fmt.Errorf("account is required (set via 'digit config set' or --account flag)")
		}
		if serverURL == "" {
			return fmt.Errorf("server URL is required (set via 'digit config set' or --server flag)")
		}
		if jwtToken == "" {
			return fmt.Errorf("JWT token is required (set via 'digit config set' or --jwt-token flag)")
		}

		// Delete the user
		responseBody, err := digit.DeleteUser(serverURL, jwtToken, realm, username)
		if err != nil {
			return fmt.Errorf("failed to delete user: %w", err)
		}

		// Print response
		fmt.Println("User deletion response:")
		if responseBody == "" {
			fmt.Println("User deleted successfully (no response body)")
		} else {
			fmt.Println(responseBody)
		}

		return nil
	},
}

// searchUserCmd represents the search-user command
var searchUserCmd = &cobra.Command{
	Use:   "search-user",
	Short: "Search for users in Keycloak",
	Long: `Search for users in Keycloak by username or list all users.

This command will:
1. If username is provided, search for that specific user
2. If no username is provided, return all users in the account

Example:
  digit search-user --account master                    # List all users
  digit search-user --username johndoe --account master # Search specific user
  digit search-user --username johndoe --account myrealm --server https://keycloak.example.com --jwt-token <token>`,
	RunE: func(cmd *cobra.Command, args []string) error {
		username, _ := cmd.Flags().GetString("username")
		realm, _ := cmd.Flags().GetString("account")
		serverURL, _ := cmd.Flags().GetString("server")
		jwtToken, _ := cmd.Flags().GetString("jwt-token")

		// Load configuration
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		// Get server URL from config if not provided
		if serverURL == "" {
			serverURL = cfg.Server
		}

		// Get JWT token from config if not provided
		if jwtToken == "" {
			jwtToken = cfg.JWTToken
		}

		// Get realm from config if not provided
		if realm == "" {
			realm, err = config.GetRealm()
			if err != nil {
				return fmt.Errorf("failed to get realm from config: %w", err)
			}
		}

		// Validate required parameters
		if realm == "" {
			return fmt.Errorf("account is required (set via 'digit config set' or --account flag)")
		}
		if serverURL == "" {
			return fmt.Errorf("server URL is required (set via 'digit config set' or --server flag)")
		}
		if jwtToken == "" {
			return fmt.Errorf("JWT token is required (set via 'digit config set' or --jwt-token flag)")
		}

		// Search for users
		responseBody, err := digit.SearchUser(serverURL, jwtToken, realm, username)
		if err != nil {
			return fmt.Errorf("failed to search users: %w", err)
		}

		// Print response
		if username != "" {
			fmt.Printf("Search results for user '%s':\n", username)
		} else {
			fmt.Println("All users in realm:")
		}
		
		if responseBody == "[]" {
			if username != "" {
				fmt.Printf("No user found with username '%s'\n", username)
			} else {
				fmt.Println("No users found in realm")
			}
		} else {
			fmt.Println(responseBody)
		}

		return nil
	},
}

// updateUserCmd represents the update-user command
var updateUserCmd = &cobra.Command{
	Use:   "update-user",
	Short: "Update a user in Keycloak",
	Long: `Update a user's information in Keycloak.

This command will:
1. Find the user by username in the specified account
2. Update the user's information with provided fields

Example:
  digit update-user --username johndoe --first-name John --last-name Doe --account master
  digit update-user --username johndoe --enabled=false --account master --server https://keycloak.example.com --jwt-token <token>`,
	RunE: func(cmd *cobra.Command, args []string) error {
		username, _ := cmd.Flags().GetString("username")
		email, _ := cmd.Flags().GetString("email")
		firstName, _ := cmd.Flags().GetString("first-name")
		lastName, _ := cmd.Flags().GetString("last-name")
		enabledStr, _ := cmd.Flags().GetString("enabled")
		realm, _ := cmd.Flags().GetString("account")
		serverURL, _ := cmd.Flags().GetString("server")
		jwtToken, _ := cmd.Flags().GetString("jwt-token")

		// Load configuration
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		// Get server URL from config if not provided
		if serverURL == "" {
			serverURL = cfg.Server
		}

		// Get JWT token from config if not provided
		if jwtToken == "" {
			jwtToken = cfg.JWTToken
		}

		// Get realm from config if not provided
		if realm == "" {
			realm, err = config.GetRealm()
			if err != nil {
				return fmt.Errorf("failed to get realm from config: %w", err)
			}
		}

		// Parse enabled flag
		var enabled *bool
		if enabledStr != "" {
			if enabledStr == "true" {
				val := true
				enabled = &val
			} else if enabledStr == "false" {
				val := false
				enabled = &val
			} else {
				return fmt.Errorf("enabled flag must be 'true' or 'false'")
			}
		}

		// Validate required parameters
		if username == "" {
			return fmt.Errorf("username is required")
		}
		if realm == "" {
			return fmt.Errorf("account is required (set via 'digit config set' or --account flag)")
		}
		if serverURL == "" {
			return fmt.Errorf("server URL is required (set via 'digit config set' or --server flag)")
		}
		if jwtToken == "" {
			return fmt.Errorf("JWT token is required (set via 'digit config set' or --jwt-token flag)")
		}

		// Check if at least one field to update is provided
		if email == "" && firstName == "" && lastName == "" && enabled == nil {
			return fmt.Errorf("at least one field to update must be provided (email, first-name, last-name, or enabled)")
		}

		// Update the user
		responseBody, err := digit.UpdateUser(serverURL, jwtToken, realm, username, email, firstName, lastName, enabled)
		if err != nil {
			return fmt.Errorf("failed to update user: %w", err)
		}

		// Print response
		fmt.Println("User update response:")
		if responseBody == "" {
			fmt.Println("User updated successfully (no response body)")
		} else {
			fmt.Println(responseBody)
		}

		return nil
	},
}

// createRoleCmd represents the create-role command
var createRoleCmd = &cobra.Command{
	Use:   "create-role",
	Short: "Create a new role in Keycloak",
	Long: `Create a new role in Keycloak with the specified role name and optional description.
	
Examples:
  digit create-role --role-name admin --description "Administrator role" --account master
  digit create-role --role-name user --account master --server http://localhost:8080
  digit create-role --role-name manager --description "Manager role" --account master --jwt-token <token>`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get flag values
		roleName, _ := cmd.Flags().GetString("role-name")
		description, _ := cmd.Flags().GetString("description")
		realm, _ := cmd.Flags().GetString("account")
		serverURL, _ := cmd.Flags().GetString("server")
		jwtToken, _ := cmd.Flags().GetString("jwt-token")
		
		// Validate required flags
		if roleName == "" {
			return fmt.Errorf("--role-name flag is required")
		}
		
		// Load configuration if server URL, JWT token, or account not provided as flags
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
		
		if realm == "" {
			realm, err = config.GetRealm()
			if err != nil {
				return fmt.Errorf("failed to get realm from config: %w", err)
			}
			if realm == "" {
				return fmt.Errorf("account not configured. Use 'digit config set' or provide --account flag")
			}
		}
		
		// Call the digit library to create role
		responseBody, err := digit.CreateRole(serverURL, jwtToken, realm, roleName, description)
		if err != nil {
			return fmt.Errorf("failed to create role: %w", err)
		}
		
		// Print response body
		fmt.Printf("Role creation response:\n")
		
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
			if responseBody == "" {
				fmt.Printf("Role '%s' created successfully\n", roleName)
			} else {
				fmt.Println(responseBody)
			}
		}
		
		return nil
	},
}

// assignRoleCmd represents the assign-role command
var assignRoleCmd = &cobra.Command{
	Use:   "assign-role",
	Short: "Assign a role to a user in Keycloak",
	Long: `Assign an existing role to a user in Keycloak.
	
Examples:
  digit assign-role --username johndoe --role-name admin --account master
  digit assign-role --username johndoe --role-name user --account master --server http://localhost:8080
  digit assign-role --username johndoe --role-name manager --account master --jwt-token <token>`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get flag values
		username, _ := cmd.Flags().GetString("username")
		roleName, _ := cmd.Flags().GetString("role-name")
		realm, _ := cmd.Flags().GetString("account")
		serverURL, _ := cmd.Flags().GetString("server")
		jwtToken, _ := cmd.Flags().GetString("jwt-token")
		
		// Validate required flags
		if username == "" {
			return fmt.Errorf("--username flag is required")
		}
		if roleName == "" {
			return fmt.Errorf("--role-name flag is required")
		}
		
		// Load configuration if server URL, JWT token, or account not provided as flags
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
		
		if realm == "" {
			realm, err = config.GetRealm()
			if err != nil {
				return fmt.Errorf("failed to get realm from config: %w", err)
			}
			if realm == "" {
				return fmt.Errorf("account not configured. Use 'digit config set' or provide --account flag")
			}
		}
		
		// Call the digit library to assign role to user
		responseBody, err := digit.AssignRoleToUser(serverURL, jwtToken, realm, username, roleName)
		if err != nil {
			return fmt.Errorf("failed to assign role to user: %w", err)
		}
		
		// Print response body
		fmt.Printf("Role assignment response:\n")
		
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
			if responseBody == "" {
				fmt.Printf("Role '%s' assigned to user '%s' successfully\n", roleName, username)
			} else {
				fmt.Println(responseBody)
			}
		}
		
		return nil
	},
}

func init() {
	rootCmd.AddCommand(createUserCmd)
	rootCmd.AddCommand(resetPasswordCmd)
	rootCmd.AddCommand(deleteUserCmd)
	rootCmd.AddCommand(searchUserCmd)
	rootCmd.AddCommand(updateUserCmd)
	rootCmd.AddCommand(createRoleCmd)
	rootCmd.AddCommand(assignRoleCmd)

	// Flags for create-user command
	createUserCmd.Flags().StringP("username", "u", "", "Username for the new user")
	createUserCmd.Flags().StringP("password", "p", "", "Password for the new user")
	createUserCmd.Flags().StringP("email", "e", "", "Email for the new user")
	createUserCmd.Flags().StringP("account", "a", "", "Keycloak account")
	createUserCmd.Flags().StringP("server", "s", "", "Keycloak server URL (overrides config)")
	createUserCmd.Flags().StringP("jwt-token", "t", "", "JWT token for authentication (overrides config)")

	// Mark required flags
	createUserCmd.MarkFlagRequired("username")
	createUserCmd.MarkFlagRequired("password")
	createUserCmd.MarkFlagRequired("email")

	// Flags for reset-password command
	resetPasswordCmd.Flags().StringP("username", "u", "", "Username of the user to reset password")
	resetPasswordCmd.Flags().StringP("new-password", "p", "", "New password for the user")
	resetPasswordCmd.Flags().StringP("account", "a", "", "Keycloak account")
	resetPasswordCmd.Flags().StringP("server", "s", "", "Keycloak server URL (overrides config)")
	resetPasswordCmd.Flags().StringP("jwt-token", "t", "", "JWT token for authentication (overrides config)")

	// Mark required flags for reset-password
	resetPasswordCmd.MarkFlagRequired("username")
	resetPasswordCmd.MarkFlagRequired("new-password")

	// Flags for delete-user command
	deleteUserCmd.Flags().StringP("username", "u", "", "Username of the user to delete")
	deleteUserCmd.Flags().StringP("account", "a", "", "Keycloak account")
	deleteUserCmd.Flags().StringP("server", "s", "", "Keycloak server URL (overrides config)")
	deleteUserCmd.Flags().StringP("jwt-token", "t", "", "JWT token for authentication (overrides config)")

	// Mark required flags for delete-user
	deleteUserCmd.MarkFlagRequired("username")

	// Flags for search-user command
	searchUserCmd.Flags().StringP("username", "u", "", "Username to search for (optional - if not provided, lists all users)")
	searchUserCmd.Flags().StringP("account", "a", "", "Keycloak account")
	searchUserCmd.Flags().StringP("server", "s", "", "Keycloak server URL (overrides config)")
	searchUserCmd.Flags().StringP("jwt-token", "t", "", "JWT token for authentication (overrides config)")

	// Mark required flags for search-user (none - all optional)

	// Flags for update-user command
	updateUserCmd.Flags().StringP("username", "u", "", "Username of the user to update")
	updateUserCmd.Flags().StringP("email", "e", "", "New email for the user")
	updateUserCmd.Flags().StringP("first-name", "f", "", "New first name for the user")
	updateUserCmd.Flags().StringP("last-name", "l", "", "New last name for the user")
	updateUserCmd.Flags().StringP("enabled", "", "", "Enable/disable user (true/false)")
	updateUserCmd.Flags().StringP("account", "a", "", "Keycloak account")
	updateUserCmd.Flags().StringP("server", "s", "", "Keycloak server URL (overrides config)")
	updateUserCmd.Flags().StringP("jwt-token", "t", "", "JWT token for authentication (overrides config)")

	// Mark required flags for update-user
	updateUserCmd.MarkFlagRequired("username")

	// Flags for create-role command
	createRoleCmd.Flags().StringP("role-name", "n", "", "Name for the new role")
	createRoleCmd.Flags().StringP("description", "d", "", "Description for the new role")
	createRoleCmd.Flags().StringP("account", "a", "", "Keycloak account")
	createRoleCmd.Flags().StringP("server", "s", "", "Keycloak server URL (overrides config)")
	createRoleCmd.Flags().StringP("jwt-token", "t", "", "JWT token for authentication (overrides config)")

	// Mark required flags for create-role
	createRoleCmd.MarkFlagRequired("role-name")

	// Flags for assign-role command
	assignRoleCmd.Flags().StringP("username", "u", "", "Username of the user to assign role to")
	assignRoleCmd.Flags().StringP("role-name", "n", "", "Name of the role to assign")
	assignRoleCmd.Flags().StringP("account", "a", "", "Keycloak account")
	assignRoleCmd.Flags().StringP("server", "s", "", "Keycloak server URL (overrides config)")
	assignRoleCmd.Flags().StringP("jwt-token", "t", "", "JWT token for authentication (overrides config)")

	// Mark required flags for assign-role
	assignRoleCmd.MarkFlagRequired("username")
	assignRoleCmd.MarkFlagRequired("role-name")
}
