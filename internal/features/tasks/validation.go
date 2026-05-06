package tasks

import (
	"fmt"
	"strings"
)

const (
	maxTitleLength       = 200
	maxDescriptionLength = 2000
	defaultLimit         = 50
	maxLimit             = 100
)

func validateTitle(title string) error {
	if strings.TrimSpace(title) == "" {
		return fmt.Errorf("title is required")
	}
	if len(title) > maxTitleLength {
		return fmt.Errorf("title must be at most %d characters", maxTitleLength)
	}
	return nil
}

func validateDescription(description string) error {
	if len(description) > maxDescriptionLength {
		return fmt.Errorf("description must be at most %d characters", maxDescriptionLength)
	}
	return nil
}

func normalizeLimitOffset(limit, offset int) (int, int) {
	if limit <= 0 {
		limit = defaultLimit
	}
	if limit > maxLimit {
		limit = maxLimit
	}
	if offset < 0 {
		offset = 0
	}
	return limit, offset
}
