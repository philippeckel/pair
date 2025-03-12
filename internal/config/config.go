package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/philippeckel/pair/internal/models"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

var (
	ConfigPath string
	Config     models.Config
)

// GetDefaultConfigPath returns the default path for the config file,
// prioritizing a local .pair.json if it exists
func GetConfigPath() string {
	// First check if there's a .pair.json file in the current directory
	localConfig := ".pair.json"
	if _, err := os.Stat(localConfig); err == nil {
		// Local config exists, use it
		return localConfig
	}

	// No local config, use the one in home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// If we can't get home dir, use current directory as fallback
		return localConfig
	}
	return filepath.Join(homeDir, ".pair.json")
}

// LoadConfig loads and parses the config file
func LoadConfig() error {
	data, err := os.ReadFile(ConfigPath)
	if err != nil {
		return fmt.Errorf("could not read config file: %w", err)
	}

	// First parse the structure to get the raw JSON
	var jsonData struct {
		CoAuthors json.RawMessage `json:"coauthors"`
	}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return fmt.Errorf("could not parse config file: %w", err)
	}

	// We'll decode the coauthors map while preserving order
	var orderedAliases []string
	var tempMap map[string]struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	// Decode the map while capturing order
	decoder := json.NewDecoder(bytes.NewReader(jsonData.CoAuthors))
	if err := decoder.Decode(&tempMap); err != nil {
		return fmt.Errorf("could not decode coauthors map: %w", err)
	}

	// Now parse again to get the keys in order
	// This part relies on the decoder preserving token order
	decoder = json.NewDecoder(bytes.NewReader(jsonData.CoAuthors))

	// Skip the opening brace
	_, err = decoder.Token()
	if err != nil {
		return fmt.Errorf("error parsing JSON tokens: %w", err)
	}

	// Read all keys in order
	for decoder.More() {
		token, err := decoder.Token()
		if err != nil {
			return fmt.Errorf("error reading JSON token: %w", err)
		}
		if key, ok := token.(string); ok {
			orderedAliases = append(orderedAliases, key)

			// Skip the value (we already have it in tempMap)
			var v interface{}
			if err := decoder.Decode(&v); err != nil {
				return fmt.Errorf("error skipping JSON value: %w", err)
			}
		}
	}

	// Now build our config data in the correct order
	Config.CoAuthorsMap = make(map[string]models.CoAuthor)
	Config.CoAuthors = make([]models.CoAuthor, 0, len(tempMap))

	for _, alias := range orderedAliases {
		details := tempMap[alias]
		coauthor := models.CoAuthor{
			Name:  details.Name,
			Email: details.Email,
			Alias: alias,
		}

		if err := coauthor.Validate(); err != nil {
			return fmt.Errorf("invalid co-author '%s': %w", alias, err)
		}

		Config.CoAuthorsMap[alias] = coauthor
		Config.CoAuthors = append(Config.CoAuthors, coauthor)
	}

	return nil
}

// SaveConfig saves the config back to the file
func SaveConfig() error {
	// Convert our internal structure back to the expected JSON format
	tempConfig := struct {
		CoAuthorsMap map[string]struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"coauthors"`
	}{
		CoAuthorsMap: make(map[string]struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		}),
	}

	for alias, coauthor := range Config.CoAuthorsMap {
		tempConfig.CoAuthorsMap[alias] = struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		}{
			Name:  coauthor.Name,
			Email: coauthor.Email,
		}
	}

	data, err := json.MarshalIndent(tempConfig, "", "  ")
	if err != nil {
		return fmt.Errorf("error encoding config: %w", err)
	}

	if err := os.WriteFile(ConfigPath, data, 0644); err != nil {
		return fmt.Errorf("error writing config file: %w", err)
	}

	return nil
}

func InitViper() error {
	// Get home directory for config file
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("error getting home directory: %w", err)
	}

	configDir := filepath.Join(home, ".config", "pair")

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("could not create config directory: %w", err)
	}

	// Set default config name and paths
	viper.SetConfigName("config") // Config file name without extension
	viper.SetConfigType("yaml")   // Config file type

	// Search in home directory and current directory
	viper.AddConfigPath(configDir)
	viper.AddConfigPath(".")

	// Set default values
	viper.SetDefault("no_color", false) // Default to using colors
	viper.SetDefault("default_template_path", filepath.Join(home, ".config", "pair", "git_commit_template"))

	// Check environment variables that match the config keys
	// This will override values from config file
	viper.AutomaticEnv()

	// Override with environment variables like NO_COLOR
	if _, exists := os.LookupEnv("NO_COLOR"); exists {
		viper.Set("no_color", true)
	}

	// Load config from file (if it exists)
	if err := viper.ReadInConfig(); err != nil {
		// It's okay if config file doesn't exist
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Return error only if it's not "file not found"
			return fmt.Errorf("error reading config file: %w", err)
		}
	}

	return nil
}

// IsColorDisabled returns true if colors should be disabled
func IsColorDisabled() bool {
	return viper.GetBool("no_color")
}

// GetDebug returns true if debug mode is enabled
func GetDebug() bool {
	return viper.GetBool("debug")
}

// GetTemplatePath returns the configured git template path
func GetTemplatePath() string {
	return viper.GetString("default_template_path")
}
