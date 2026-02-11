package cmd

import (
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage CLI configuration",
	Long:  `Manage CLI configuration settings such as server URL.`,
}

func init() {
	rootCmd.AddCommand(configCmd)
}
