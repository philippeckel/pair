package models

import (
	"fmt"
	"strings"
)

// CoAuthor represents a contributor that can be added to commits
type CoAuthor struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Alias string `json:"-"` // Not in JSON, filled from the map key
}

// Validate checks if the CoAuthor has valid fields
func (c *CoAuthor) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}

	if c.Email == "" {
		return fmt.Errorf("email cannot be empty")
	}

	if !strings.Contains(c.Email, "@") {
		return fmt.Errorf("email must contain '@' character")
	}

	return nil
}

// Config holds all available co-authors
type Config struct {
	CoAuthorsMap map[string]CoAuthor `json:"coauthors"`
	CoAuthors    []CoAuthor          `json:"-"` // This will be filled after loading
}
