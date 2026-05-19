package cmd

import (
	"encoding/json"
	"fmt"

	"digit-cli/pkg/config"
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

// createAccountCmd represents the create-account command
var createAccountCmd = &cobra.Command{
	Use:   "create-account",
	Short: "Create a new account in DIGIT services",
	Long: `Initiate account creation via /account/v3/signup.
The response includes a requestId; use verify-account to confirm the OTP.

Examples:
  digit create-account --name kongnew1 --email test@example.com --password default
  digit create-account --name kongnew1 --email test@example.com --password mypass --server http://localhost:8094`,
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		email, _ := cmd.Flags().GetString("email")
		password, _ := cmd.Flags().GetString("password")

		serverURL, err := getServerURL(cmd)
		if err != nil {
			return err
		}

		responseBody, err := digit.SignupAccount(serverURL, name, email, password)
		if err != nil {
			return fmt.Errorf("failed to create account: %w", err)
		}

		fmt.Println("Account signup response:")
		printJSON(responseBody)
		return nil
	},
}

// verifyAccountCmd represents the verify-account command
var verifyAccountCmd = &cobra.Command{
	Use:   "verify-account",
	Short: "Verify account signup using OTP",
	Long: `Confirm account creation via /account/v3/signup/verify using the requestId and OTP received after signup.

Examples:
  digit verify-account --request-id 6a7ac916-6e49-49de-b129-9630655a37b1 --otp 245332
  digit verify-account --request-id <id> --otp <otp> --server http://localhost:8094`,
	RunE: func(cmd *cobra.Command, args []string) error {
		requestID, _ := cmd.Flags().GetString("request-id")
		otp, _ := cmd.Flags().GetString("otp")

		serverURL, err := getServerURL(cmd)
		if err != nil {
			return err
		}

		responseBody, err := digit.VerifyAccount(serverURL, requestID, otp)
		if err != nil {
			return fmt.Errorf("failed to verify account: %w", err)
		}

		fmt.Println("Account verification response:")
		printJSON(responseBody)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(createAccountCmd)
	createAccountCmd.Flags().String("name", "", "Name of the tenant (required)")
	createAccountCmd.Flags().String("email", "", "Email of the tenant (required)")
	createAccountCmd.Flags().String("password", "default", "Password for the tenant (default: default)")
	createAccountCmd.Flags().String("server", "", "Server URL (overrides config)")
	createAccountCmd.MarkFlagRequired("name")
	createAccountCmd.MarkFlagRequired("email")

	rootCmd.AddCommand(verifyAccountCmd)
	verifyAccountCmd.Flags().String("request-id", "", "Request ID returned from create-account (required)")
	verifyAccountCmd.Flags().String("otp", "", "OTP received via email (required)")
	verifyAccountCmd.Flags().String("server", "", "Server URL (overrides config)")
	verifyAccountCmd.MarkFlagRequired("request-id")
	verifyAccountCmd.MarkFlagRequired("otp")
}
