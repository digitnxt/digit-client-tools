package cmd

import (
	"fmt"

	"digit-cli/pkg/config"
	"github.com/spf13/cobra"
)

// configGetContextsCmd represents the config get-contexts command
var configGetContextsCmd = &cobra.Command{
	Use:   "get-contexts",
	Short: "List all available contexts from config file",
	Long: `Display all available contexts from the specified configuration file.

Examples:
  digit config get-contexts --file ./digit-config.yaml`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get flag values
		filePath, _ := cmd.Flags().GetString("file")

		// Validate required flags
		if filePath == "" {
			return fmt.Errorf("--file flag is required")
		}

		// Load context configuration from file
		contextConfig, err := config.LoadContextConfig(filePath)
		if err != nil {
			return fmt.Errorf("failed to load config file: %w", err)
		}

		// List all contexts
		contexts := contextConfig.ListContexts()
		if len(contexts) == 0 {
			fmt.Println("No contexts found in config file")
			return nil
		}

		fmt.Printf("Available contexts (current: %s):\n", contextConfig.CurrentContext)
		for _, ctx := range contexts {
			if ctx == contextConfig.CurrentContext {
				fmt.Printf("* %s (current)\n", ctx)
			} else {
				fmt.Printf("  %s\n", ctx)
			}
		}

		return nil
	},
}

func init() {
	configCmd.AddCommand(configGetContextsCmd)
	configGetContextsCmd.Flags().StringP("file", "f", "", "Path to the configuration YAML file")
}
