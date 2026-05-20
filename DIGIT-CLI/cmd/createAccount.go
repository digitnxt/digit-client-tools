package cmd

import (
	"encoding/json"
	"fmt"

	"digit-cli/pkg/config"
	"digit-cli/pkg/jwt"
	"github.com/digitnxt/digit-client-tools/client-libraries/digit-library/digit"
	"github.com/spf13/cobra"
)

func getServerURL(cmd *cobra.Command) (string, error) {
	serverURL, _ := cmd.Flags().GetString("server")
	if serverURL != "" {
		return serverURL, nil
	}
	cfg, err := config.Load()
	if err != nil {
		return "", fmt.Errorf("failed to load config: %w", err)
	}
	serverURL = cfg.GetServer()
	if serverURL == "" {
		return "", fmt.Errorf("server URL not configured. Use 'digit config set --server <url>' or provide --server flag")
	}
	return serverURL, nil
}

func printJSON(responseBody string) {
	var jsonResponse interface{}
	if err := json.Unmarshal([]byte(responseBody), &jsonResponse); err == nil {
		prettyJSON, err := json.MarshalIndent(jsonResponse, "", "  ")
		if err == nil {
			fmt.Println(string(prettyJSON))
			return
		}
	}
	fmt.Println(responseBody)
}

var createAccountCmd = &cobra.Command{
	Use:   "create-account",
	Short: "Create a new tenant account",
	Long: `Create a new tenant account via the admin API.
If --password is omitted, the server generates one and emails it to the tenant.

Examples:
  digit create-account --name "Nairobi City" --email admin@nairobi.go.ke
  digit create-account --name "Nairobi City" --email admin@nairobi.go.ke --password Changeme1
  digit create-account --name "Nairobi City" --email admin@nairobi.go.ke --phone "+254202229000" --city Nairobi
  digit create-account --name "Nairobi City" --email admin@nairobi.go.ke --server http://localhost:8080`,
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		email, _ := cmd.Flags().GetString("email")
		password, _ := cmd.Flags().GetString("password")
		phone, _ := cmd.Flags().GetString("phone")
		address, _ := cmd.Flags().GetString("address")
		city, _ := cmd.Flags().GetString("city")
		state, _ := cmd.Flags().GetString("state")
		pincode, _ := cmd.Flags().GetString("pincode")
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
					return fmt.Errorf("JWT token not configured. Use 'digit config set' or provide --jwt-token flag")
				}
			}
		}

		tenantID, err := jwt.ExtractTenantID(jwtToken)
		if err != nil {
			return fmt.Errorf("failed to extract tenant ID from JWT token: %w", err)
		}

		req := digit.TenantCreateRequest{
			Name:     name,
			Email:    email,
			Password: password,
			Phone:    phone,
			Address:  address,
			City:     city,
			State:    state,
			Pincode:  pincode,
		}

		responseBody, err := digit.CreateTenant(serverURL, jwtToken, tenantID, req)
		if err != nil {
			return fmt.Errorf("failed to create account: %w", err)
		}

		fmt.Println("Account created successfully:")
		printJSON(responseBody)
		return nil
	},
}

var searchAccountCmd = &cobra.Command{
	Use:   "search-account",
	Short: "Search or list tenant accounts",
	Long: `Search tenant accounts with optional filters on name or email.

Examples:
  digit search-account
  digit search-account --name "Nairobi"
  digit search-account --email admin@nairobi.go.ke
  digit search-account --page 2 --size 10
  digit search-account --name "Nairobi" --server http://localhost:8080`,
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		email, _ := cmd.Flags().GetString("email")
		page, _ := cmd.Flags().GetInt("page")
		size, _ := cmd.Flags().GetInt("size")
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
					return fmt.Errorf("JWT token not configured. Use 'digit config set' or provide --jwt-token flag")
				}
			}
		}

		tenantID, err := jwt.ExtractTenantID(jwtToken)
		if err != nil {
			return fmt.Errorf("failed to extract tenant ID from JWT token: %w", err)
		}

		responseBody, err := digit.SearchTenants(serverURL, jwtToken, tenantID, name, email, page, size)
		if err != nil {
			return fmt.Errorf("failed to search accounts: %w", err)
		}

		printJSON(responseBody)
		return nil
	},
}

var deleteAccountCmd = &cobra.Command{
	Use:   "delete-account",
	Short: "Delete a tenant account by ID",
	Long: `Permanently delete a tenant account by its ID.

Examples:
  digit delete-account --id 3fa85f64-5717-4562-b3fc-2c963f66afa6
  digit delete-account --id 3fa85f64-5717-4562-b3fc-2c963f66afa6 --server http://localhost:8080`,
	RunE: func(cmd *cobra.Command, args []string) error {
		id, _ := cmd.Flags().GetString("id")
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
					return fmt.Errorf("JWT token not configured. Use 'digit config set' or provide --jwt-token flag")
				}
			}
		}

		tenantID, err := jwt.ExtractTenantID(jwtToken)
		if err != nil {
			return fmt.Errorf("failed to extract tenant ID from JWT token: %w", err)
		}

		responseBody, err := digit.DeleteTenant(serverURL, jwtToken, tenantID, id)
		if err != nil {
			return fmt.Errorf("failed to delete account: %w", err)
		}

		if responseBody != "" {
			printJSON(responseBody)
		} else {
			fmt.Println("Account deleted successfully")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(createAccountCmd)
	rootCmd.AddCommand(searchAccountCmd)
	rootCmd.AddCommand(deleteAccountCmd)

	createAccountCmd.Flags().String("name", "", "Tenant name (required)")
	createAccountCmd.Flags().String("email", "", "Tenant email (required)")
	createAccountCmd.Flags().String("password", "", "Password (optional — server generates one if omitted)")
	createAccountCmd.Flags().String("phone", "", "Phone number (optional)")
	createAccountCmd.Flags().String("address", "", "Address (optional)")
	createAccountCmd.Flags().String("city", "", "City (optional)")
	createAccountCmd.Flags().String("state", "", "State (optional)")
	createAccountCmd.Flags().String("pincode", "", "Pincode (optional)")
	createAccountCmd.Flags().String("server", "", "Server URL (overrides config)")
	createAccountCmd.Flags().String("jwt-token", "", "JWT token (overrides config)")

	createAccountCmd.MarkFlagRequired("name")
	createAccountCmd.MarkFlagRequired("email")

	searchAccountCmd.Flags().String("name", "", "Filter by tenant name (partial match)")
	searchAccountCmd.Flags().String("email", "", "Filter by tenant email (partial match)")
	searchAccountCmd.Flags().Int("page", 1, "Page number (1-indexed)")
	searchAccountCmd.Flags().Int("size", 20, "Number of results per page")
	searchAccountCmd.Flags().String("server", "", "Server URL (overrides config)")
	searchAccountCmd.Flags().String("jwt-token", "", "JWT token (overrides config)")

	deleteAccountCmd.Flags().String("id", "", "Tenant account ID to delete (required)")
	deleteAccountCmd.Flags().String("server", "", "Server URL (overrides config)")
	deleteAccountCmd.Flags().String("jwt-token", "", "JWT token (overrides config)")
	deleteAccountCmd.MarkFlagRequired("id")
}
