package commands

import (
	"fmt"
	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/philippeckel/pair/internal/models"
	"strings"
)

// selectMultipleCoAuthors allows selecting multiple co-authors at once
func selectMultipleCoAuthors(coAuthors []models.CoAuthor, activeCoAuthors []models.CoAuthor) ([]models.CoAuthor, error) {
	userName, userEmail, err := getGitUserInfo()
	if err != nil {
		return nil, fmt.Errorf("failed to get git user info: %w", err)
	}

	// Filter out already active co-authors
	var availableCoAuthors []models.CoAuthor
	for _, author := range coAuthors {
		// Skip yourself - check both name and email
		if strings.EqualFold(author.Email, userEmail) || strings.EqualFold(author.Name, userName) {
			continue
		}
		isActive := false
		for _, active := range activeCoAuthors {
			if author.Email == active.Email {
				isActive = true
				break
			}
		}
		if !isActive {
			availableCoAuthors = append(availableCoAuthors, author)
		}
	}

	// Check if there are any available co-authors left
	if len(availableCoAuthors) == 0 {
		return nil, fmt.Errorf("all co-authors are already active")
	}

	// Run the fuzzy finder and get selected indices
	indices, err := fuzzyfinder.FindMulti(
		availableCoAuthors,
		func(i int) string {
			return fmt.Sprintf("%s (%s) <%s>", availableCoAuthors[i].Name, availableCoAuthors[i].Alias, availableCoAuthors[i].Email)
		},
		fuzzyfinder.WithPromptString("Select co-authors to add using TAB:"),
	)

	if err != nil {
		// User canceled selection
		if err == fuzzyfinder.ErrAbort {
			return nil, fmt.Errorf("selection canceled")
		}
		return nil, fmt.Errorf("fuzzy finder error: %w", err)
	}

	// Convert selected indices to co-authors
	var selectedCoAuthors []models.CoAuthor
	for _, idx := range indices {
		selectedCoAuthors = append(selectedCoAuthors, availableCoAuthors[idx])
	}

	return selectedCoAuthors, nil
}

// selectMultipleCoAuthorsToRemove allows selecting multiple active co-authors for removal
func selectMultipleCoAuthorsToRemove(activeCoAuthors []models.CoAuthor) ([]models.CoAuthor, error) {
	if len(activeCoAuthors) == 0 {
		return nil, fmt.Errorf("no active co-authors to select from")
	}

	// Run the fuzzy finder and get selected indices
	indices, err := fuzzyfinder.FindMulti(
		activeCoAuthors,
		func(i int) string {
			return fmt.Sprintf("%s <%s>", activeCoAuthors[i].Name, activeCoAuthors[i].Email)
		},
		fuzzyfinder.WithPromptString("Select co-authors to remove (TAB/Space to select multiple):"),
		fuzzyfinder.WithPreviewWindow(func(i, _, _ int) string {
			if i == -1 {
				return ""
			}
			return fmt.Sprintf("Name: %s\nEmail: %s",
				activeCoAuthors[i].Name,
				activeCoAuthors[i].Email)
		}),
	)

	if err != nil {
		// User canceled selection
		if err == fuzzyfinder.ErrAbort {
			return nil, fmt.Errorf("Selection canceled")
		}
		return nil, fmt.Errorf("fuzzy finder error: %w", err)
	}

	// Convert selected indices to co-authors
	var selectedCoAuthors []models.CoAuthor
	for _, idx := range indices {
		selectedCoAuthors = append(selectedCoAuthors, activeCoAuthors[idx])
	}

	return selectedCoAuthors, nil
}
