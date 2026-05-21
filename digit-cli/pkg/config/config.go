package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the CLI configuration
type Config struct {
	Server   string `yaml:"server"`
	JWTToken string `yaml:"jwt_token"`
	// Authentication credentials for auto-refresh
	AuthConfig *AuthConfig `yaml:"auth_config,omitempty"`
}

// AuthConfig stores authentication details for token refresh
type AuthConfig struct {
	ServerURL    string `yaml:"server_url"`
	Realm        string `yaml:"realm"`
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
}

// SetJWTToken updates the JWT token in the config
func (c *Config) SetJWTToken(token string) {
	c.JWTToken = token
}

// GetJWTToken returns the configured JWT token
func (c *Config) GetJWTToken() string {
	return c.JWTToken
}

// ConfigPath returns the path to the config file
func ConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	
	configDir := filepath.Join(homeDir, ".digit")
	configFile := filepath.Join(configDir, "config.yaml")
	
	return configFile, nil
}

// Load reads the configuration from the config file
func Load() (*Config, error) {
	configFile, err := ConfigPath()
	if err != nil {
		return nil, err
	}
	
	// If config file doesn't exist, return default config
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return &Config{}, nil
	}
	
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}
	
	return &config, nil
}

// Save writes the configuration to the config file
func (c *Config) Save() error {
	configFile, err := ConfigPath()
	if err != nil {
		return err
	}
	
	// Create config directory if it doesn't exist
	configDir := filepath.Dir(configFile)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}
	
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	
	if err := os.WriteFile(configFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	
	return nil
}

// SetServer updates the server URL in the config
func (c *Config) SetServer(server string) {
	c.Server = server
}

// GetServer returns the configured server URL
func (c *Config) GetServer() string {
	return c.Server
}

// Global functions for easier access

// SetServerURL sets the server URL in the global config
func SetServerURL(serverURL string) error {
	config, err := Load()
	if err != nil {
		return err
	}
	config.SetServer(serverURL)
	return config.Save()
}

// GetServerURL gets the server URL from the global config
func GetServerURL() (string, error) {
	config, err := Load()
	if err != nil {
		return "", err
	}
	return config.GetServer(), nil
}

// SetJWTToken sets the JWT token in the global config
func SetJWTToken(token string) error {
	config, err := Load()
	if err != nil {
		return err
	}
	config.SetJWTToken(token)
	return config.Save()
}

// GetJWTToken gets the JWT token from the global config
func GetJWTToken() (string, error) {
	config, err := Load()
	if err != nil {
		return "", err
	}
	return config.GetJWTToken(), nil
}

// SetAuthConfig sets the authentication configuration for token refresh
func SetAuthConfig(serverURL, realm, clientID, clientSecret, username, password string) error {
	config, err := Load()
	if err != nil {
		return err
	}
	config.AuthConfig = &AuthConfig{
		ServerURL:    serverURL,
		Realm:        realm,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Username:     username,
		Password:     password,
	}
	return config.Save()
}

// GetAuthConfig gets the authentication configuration
func GetAuthConfig() (*AuthConfig, error) {
	config, err := Load()
	if err != nil {
		return nil, err
	}
	return config.AuthConfig, nil
}

// GetRealm gets the realm from the authentication configuration
func GetRealm() (string, error) {
	config, err := Load()
	if err != nil {
		return "", err
	}
	if config.AuthConfig != nil {
		return config.AuthConfig.Realm, nil
	}
	return "", nil
}
