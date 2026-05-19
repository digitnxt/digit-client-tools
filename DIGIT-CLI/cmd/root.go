package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "digit",
	Short: "A CLI tool for interacting with DIGIT services",
	Long: `digit is a command-line interface for interacting with DIGIT services.
It provides commands for account management, user creation, role assignment, and more.

Configuration Commands:
  digit config set --server <url> --jwt-token <token>     # Set server and JWT token manually
  digit config export --file <config.yaml>               # Import config from YAML file and authenticate
  digit config show                                       # Show current configuration
  digit config get-contexts --file <config.yaml>         # List available contexts from config file
  digit config use-context <name> --file <config.yaml>   # Switch to different context

Examples:
  digit config export --file sample-digit-config.yaml
  digit create-account --name kongnew1 --email test@example.com
  digit create-user --username johndoe --password pass123 --email john@example.com --realm CLI`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
}
