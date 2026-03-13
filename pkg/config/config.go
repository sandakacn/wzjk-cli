package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config holds the CLI configuration
type Config struct {
	APIURL string `mapstructure:"api_url"`
	Token  string `mapstructure:"token"`
	User   *User  `mapstructure:"user"`
}

// User holds the authenticated user information
type User struct {
	ID    string `mapstructure:"id"`
	Name  string `mapstructure:"name"`
	Email string `mapstructure:"email"`
}

// GetConfigDir returns the configuration directory path
func GetConfigDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	return filepath.Join(home, ".wzjk-cli")
}

// GetConfigPath returns the full path to the config file
func GetConfigPath() string {
	return filepath.Join(GetConfigDir(), "config.yaml")
}

// Load reads the configuration from file
func Load() (*Config, error) {
	configPath := GetConfigPath()

	viper.SetConfigFile(configPath)
	viper.SetDefault("api_url", "http://localhost:3000")

	// Try to read the config file
	if err := viper.ReadInConfig(); err != nil {
		// If the file doesn't exist, that's okay - we'll use defaults
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config: %w", err)
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error parsing config: %w", err)
	}

	return &cfg, nil
}

// Save writes the configuration to file
func Save(cfg *Config) error {
	dir := GetConfigDir()
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("error creating config directory: %w", err)
	}

	viper.Set("api_url", cfg.APIURL)
	viper.Set("token", cfg.Token)
	viper.Set("user", cfg.User)

	configPath := GetConfigPath()
	if err := viper.WriteConfigAs(configPath); err != nil {
		return fmt.Errorf("error writing config: %w", err)
	}

	// Ensure the config file has restricted permissions
	if err := os.Chmod(configPath, 0600); err != nil {
		return fmt.Errorf("error setting config permissions: %w", err)
	}

	return nil
}

// Clear removes the configuration directory
func Clear() error {
	dir := GetConfigDir()
	return os.RemoveAll(dir)
}

// IsLoggedIn checks if the user is currently logged in
func IsLoggedIn() bool {
	cfg, err := Load()
	if err != nil {
		return false
	}
	return cfg.Token != ""
}
