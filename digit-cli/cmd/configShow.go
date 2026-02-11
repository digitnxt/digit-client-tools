package cmd

import (
	"fmt"

	"digit-cli/pkg/config"
	"github.com/spf13/cobra"
)

// configShowCmd represents the config show command
var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Long: `Display the current configuration including server URL and JWT token status.

Examples:
  digit config show`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get current configuration
		serverURL, err := config.GetServerURL()
		if err != nil {
			fmt.Println("Server URL: Not set")
		} else {
			fmt.Printf("Server URL: %s\n", serverURL)
		}

		jwtToken, err := config.GetJWTToken()
		if err != nil {
			fmt.Println("JWT Token: Not set")
		} else {
			if len(jwtToken) > 50 {
				fmt.Printf("JWT Token: %s...\n", jwtToken[:50])
			} else {
				fmt.Printf("JWT Token: %s\n", jwtToken)
			}
		}

		return nil
	},
}

func init() {
	configCmd.AddCommand(configShowCmd)
}
