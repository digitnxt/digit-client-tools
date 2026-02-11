package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// ContextConfig represents the structure of the digit config file
type ContextConfig struct {
	APIVersion     string    `yaml:"apiVersion"`
	Kind           string    `yaml:"kind"`
	CurrentContext string    `yaml:"current-context"`
	Contexts       []Context `yaml:"contexts"`
}

// Context represents a single context configuration
type Context struct {
	Name    string        `yaml:"name"`
	Context ContextDetail `yaml:"context"`
}

// ContextDetail contains the actual configuration details
type ContextDetail struct {
	Server       string `yaml:"server"`
	Realm        string `yaml:"realm"`
	ClientID     string `yaml:"client-id"`
	ClientSecret string `yaml:"client-secret"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
}

// LoadContextConfig loads the context configuration from a YAML file
func LoadContextConfig(filePath string) (*ContextConfig, error) {
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file does not exist: %s", filePath)
	}

	// Read the file
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse YAML
	var contextConfig ContextConfig
	if err := yaml.Unmarshal(data, &contextConfig); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &contextConfig, nil
}

// GetCurrentContext returns the current context configuration
func (cc *ContextConfig) GetCurrentContext() (*ContextDetail, error) {
	if cc.CurrentContext == "" {
		return nil, fmt.Errorf("no current context set")
	}

	for _, ctx := range cc.Contexts {
		if ctx.Name == cc.CurrentContext {
			return &ctx.Context, nil
		}
	}

	return nil, fmt.Errorf("current context '%s' not found", cc.CurrentContext)
}

// GetContext returns a specific context by name
func (cc *ContextConfig) GetContext(name string) (*ContextDetail, error) {
	for _, ctx := range cc.Contexts {
		if ctx.Name == name {
			return &ctx.Context, nil
		}
	}

	return nil, fmt.Errorf("context '%s' not found", name)
}

// ListContexts returns all available context names
func (cc *ContextConfig) ListContexts() []string {
	var contexts []string
	for _, ctx := range cc.Contexts {
		contexts = append(contexts, ctx.Name)
	}
	return contexts
}

// SetFilePermissions sets secure permissions on the config file
func SetFilePermissions(filePath string) error {
	return os.Chmod(filePath, 0600) // Read/write for owner only
}

// GetDefaultConfigPath returns the default path for the digit config file
func GetDefaultConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".digit")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create config directory: %w", err)
	}

	return filepath.Join(configDir, "config"), nil
}
