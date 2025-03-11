package config

import (
	"encoding/json"
	"fmt"
	"github.com/philippeckel/pair/internal/models"
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

	// Define a temporary struct for unmarshaling with the right structure
	var tempConfig struct {
		CoAuthorsMap map[string]struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"coauthors"`
	}

	if err := json.Unmarshal(data, &tempConfig); err != nil {
		return fmt.Errorf("could not parse config file: %w", err)
	}

	// Convert the map to our internal structure
	Config.CoAuthorsMap = make(map[string]models.CoAuthor)
	Config.CoAuthors = make([]models.CoAuthor, 0, len(tempConfig.CoAuthorsMap))

	for alias, details := range tempConfig.CoAuthorsMap {
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
