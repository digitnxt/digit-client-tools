package cmd

import (
	"fmt"

	"digit-cli/pkg/config"
	"github.com/digitnxt/digit3/code/libraries/digit-library/digit"
	"github.com/spf13/cobra"
)

// configSetCmd represents the config set command (authentication-based)
var configSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Authenticate with Keycloak and set configuration",
	Long: `Authenticate with Keycloak using configuration from a YAML file or command-line flags, get JWT token,
and store the configuration locally for use with other commands.

Examples:
  # Using YAML file
  digit config set --file ./digit-config.yaml
  
  # Using command-line flags
  digit config set --server https://digit-lts.digit.org --account CLI --client-id admin-cli --client-secret mysecret --username user@example.com --password mypassword`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get flag values
		filePath, _ := cmd.Flags().GetString("file")
		server, _ := cmd.Flags().GetString("server")
		realm, _ := cmd.Flags().GetString("account")
		clientID, _ := cmd.Flags().GetString("client-id")
		clientSecret, _ := cmd.Flags().GetString("client-secret")
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")

		// Validate that either file or all individual flags are provided
		if filePath == "" && (server == "" || realm == "" || clientID == "" || clientSecret == "" || username == "" || password == "") {
			return fmt.Errorf("either --file flag or all of (--server, --account, --client-id, --client-secret, --username, --password) flags are required")
		}

		if filePath != "" && (server != "" || realm != "" || clientID != "" || clientSecret != "" || username != "" || password != "") {
			return fmt.Errorf("cannot use --file flag together with individual flags (--server, --account, --client-id, --client-secret, --username, --password)")
		}

		var serverURL, realmName, clientIDValue, clientSecretValue, usernameValue, passwordValue string

		if filePath != "" {
			// Load context configuration from file
			fmt.Printf("Loading configuration from: %s\n", filePath)
			contextConfig, err := config.LoadContextConfig(filePath)
			if err != nil {
				return fmt.Errorf("failed to load config file: %w", err)
			}

			// Get current context
			currentCtx, err := contextConfig.GetCurrentContext()
			if err != nil {
				return fmt.Errorf("failed to get current context: %w", err)
			}

			fmt.Printf("Using context: %s\n", contextConfig.CurrentContext)
			serverURL = currentCtx.Server
			realmName = currentCtx.Realm
			clientIDValue = currentCtx.ClientID
			clientSecretValue = currentCtx.ClientSecret
			usernameValue = currentCtx.Username
			passwordValue = currentCtx.Password
		} else {
			// Use individual flags
			fmt.Println("Using configuration from command-line flags")
			serverURL = server
			realmName = realm
			clientIDValue = clientID
			clientSecretValue = clientSecret
			usernameValue = username
			passwordValue = password
		}

		fmt.Printf("Server: %s\n", serverURL)
		fmt.Printf("Account: %s\n", realmName)
		fmt.Printf("Username: %s\n", usernameValue)

		// Get JWT token from Keycloak
		fmt.Println("Authenticating with Keycloak...")
		jwtToken, err := digit.GetJWTToken(
			serverURL,
			realmName,
			clientIDValue,
			clientSecretValue,
			usernameValue,
			passwordValue,
		)
		if err != nil {
			return fmt.Errorf("failed to authenticate with Keycloak: %w", err)
		}

		fmt.Println("✓ Authentication successful!")

		// Store the configuration locally
		fmt.Println("Storing configuration...")
		err = config.SetServerURL(serverURL)
		if err != nil {
			return fmt.Errorf("failed to set server URL: %w", err)
		}

		err = config.SetJWTToken(jwtToken)
		if err != nil {
			return fmt.Errorf("failed to set JWT token: %w", err)
		}

		// Store authentication credentials for auto-refresh
		err = config.SetAuthConfig(serverURL, realmName, clientIDValue, clientSecretValue, usernameValue, passwordValue)
		if err != nil {
			return fmt.Errorf("failed to store auth config: %w", err)
		}

		fmt.Println("✓ Configuration stored successfully!")
		fmt.Printf("Server URL: %s\n", serverURL)
		fmt.Printf("JWT Token: %s...\n", jwtToken[:50])

		return nil
	},
}

func init() {
	configCmd.AddCommand(configSetCmd)
	configSetCmd.Flags().StringP("file", "f", "", "Path to the configuration YAML file")
	configSetCmd.Flags().String("server", "", "Server URL (e.g., https://digit-lts.digit.org)")
	configSetCmd.Flags().String("account", "", "Keycloak account name")
	configSetCmd.Flags().String("client-id", "", "Keycloak client ID")
	configSetCmd.Flags().String("client-secret", "", "Keycloak client secret")
	configSetCmd.Flags().String("username", "", "Username for authentication")
	configSetCmd.Flags().String("password", "", "Password for authentication")
}
