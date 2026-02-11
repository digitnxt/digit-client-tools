package cmd

import (
	"fmt"

	"digit-cli/pkg/config"
	"github.com/digitnxt/digit3/code/libraries/digit-library/digit"
	"github.com/spf13/cobra"
)

// configUseContextCmd represents the config use-context command
var configUseContextCmd = &cobra.Command{
	Use:   "use-context",
	Short: "Switch to a different context and authenticate",
	Long: `Switch to a different context from the config file, authenticate with Keycloak,
and update the local configuration.

Examples:
  digit config use-context staging --file ./digit-config.yaml
  digit config use-context production --file ./digit-config.yaml`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		contextName := args[0]
		
		// Get flag values
		filePath, _ := cmd.Flags().GetString("file")

		// Validate required flags
		if filePath == "" {
			return fmt.Errorf("--file flag is required")
		}

		// Load context configuration from file
		fmt.Printf("Loading configuration from: %s\n", filePath)
		contextConfig, err := config.LoadContextConfig(filePath)
		if err != nil {
			return fmt.Errorf("failed to load config file: %w", err)
		}

		// Get the specified context
		ctx, err := contextConfig.GetContext(contextName)
		if err != nil {
			return fmt.Errorf("failed to get context '%s': %w", contextName, err)
		}

		fmt.Printf("Switching to context: %s\n", contextName)
		fmt.Printf("Server: %s\n", ctx.Server)
		fmt.Printf("Realm: %s\n", ctx.Realm)
		fmt.Printf("Username: %s\n", ctx.Username)

		// Get JWT token from Keycloak
		fmt.Println("Authenticating with Keycloak...")
		jwtToken, err := digit.GetJWTToken(
			ctx.Server,
			ctx.Realm,
			ctx.ClientID,
			ctx.ClientSecret,
			ctx.Username,
			ctx.Password,
		)
		if err != nil {
			return fmt.Errorf("failed to authenticate with Keycloak: %w", err)
		}

		fmt.Println("✓ Authentication successful!")

		// Store the configuration locally
		fmt.Println("Updating configuration...")
		err = config.SetServerURL(ctx.Server)
		if err != nil {
			return fmt.Errorf("failed to set server URL: %w", err)
		}

		err = config.SetJWTToken(jwtToken)
		if err != nil {
			return fmt.Errorf("failed to set JWT token: %w", err)
		}

		// Store authentication credentials for auto-refresh
		err = config.SetAuthConfig(ctx.Server, ctx.Realm, ctx.ClientID, ctx.ClientSecret, ctx.Username, ctx.Password)
		if err != nil {
			return fmt.Errorf("failed to store auth config: %w", err)
		}

		fmt.Printf("✓ Switched to context '%s' successfully!\n", contextName)
		fmt.Printf("Server URL: %s\n", ctx.Server)
		fmt.Printf("JWT Token: %s...\n", jwtToken[:50])

		return nil
	},
}

func init() {
	configCmd.AddCommand(configUseContextCmd)
	configUseContextCmd.Flags().StringP("file", "f", "", "Path to the configuration YAML file")
}
