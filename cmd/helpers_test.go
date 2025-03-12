package commands

import (
	"testing"

	"github.com/philippeckel/pair/internal/config"
	"github.com/philippeckel/pair/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestFindCoAuthorByAliasOrIndex(t *testing.T) {
	// Setup test config
	config.Config = models.Config{
		CoAuthorsMap: map[string]models.CoAuthor{
			"john": {
				Name:  "John Doe",
				Email: "john@example.com",
				Alias: "john",
			},
			"jane": {
				Name:  "Jane Smith",
				Email: "jane@example.com",
				Alias: "jane",
			},
			"sam": {
				Name:  "Sam Johnson",
				Email: "sam@example.com",
				Alias: "sam",
			},
		},
	}

	// Fill the CoAuthors slice from the map
	config.Config.CoAuthors = make([]models.CoAuthor, 0, len(config.Config.CoAuthorsMap))
	for _, author := range config.Config.CoAuthorsMap {
		config.Config.CoAuthors = append(config.Config.CoAuthors, author)
	}

	// Test cases
	tests := []struct {
		name           string
		identifier     string
		expectAuthor   models.CoAuthor
		expectIndex    int
		expectErr      bool
		expectErrMatch string
	}{
		{
			name:       "Find by valid alias",
			identifier: "john",
			expectAuthor: models.CoAuthor{
				Name:  "John Doe",
				Email: "john@example.com",
				Alias: "john",
			},
			expectIndex: -1,
			expectErr:   false,
		},
		{
			name:         "Find by valid index",
			identifier:   "0",
			expectAuthor: config.Config.CoAuthors[0],
			expectIndex:  0,
			expectErr:    false,
		},
		{
			name:           "Invalid alias",
			identifier:     "unknown",
			expectAuthor:   models.CoAuthor{},
			expectIndex:    -1,
			expectErr:      true,
			expectErrMatch: "no co-author found with alias 'unknown'",
		},
		{
			name:           "Invalid index - negative",
			identifier:     "-1",
			expectAuthor:   models.CoAuthor{},
			expectIndex:    -1,
			expectErr:      true,
			expectErrMatch: "invalid co-author index",
		},
		{
			name:           "Invalid index - too large",
			identifier:     "999",
			expectAuthor:   models.CoAuthor{},
			expectIndex:    -1,
			expectErr:      true,
			expectErrMatch: "invalid co-author index",
		},
		{
			name:           "Non-numeric string, not an alias",
			identifier:     "abc123",
			expectAuthor:   models.CoAuthor{},
			expectIndex:    -1,
			expectErr:      true,
			expectErrMatch: "no co-author found with alias 'abc123'",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			author, idx, err := findCoAuthorByAliasOrIndex(tc.identifier)

			if tc.expectErr {
				assert.Error(t, err)
				if tc.expectErrMatch != "" {
					assert.Contains(t, err.Error(), tc.expectErrMatch)
				}
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.expectIndex, idx)
			assert.Equal(t, tc.expectAuthor, author)
		})
	}
}
