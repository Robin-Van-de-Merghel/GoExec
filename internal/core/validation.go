package core

import (
	"fmt"
	"strings"
)

// ValidateModules checks for duplicate module names and other issues
func ValidateModules() error {
	names := make(map[string]int)
	nameLowerMap := make(map[string][]string)

	for _, entry := range AllModules {
		name := entry.Metadata.UniqueName
		lower := strings.ToLower(name)

		names[name]++
		nameLowerMap[lower] = append(nameLowerMap[lower], name)
	}

	// Check for exact duplicates
	for name, count := range names {
		if count > 1 {
			return fmt.Errorf("duplicate module name (exact match): '%s' found %d times", name, count)
		}
	}

	// Check for case-insensitive duplicates
	for lower, variants := range nameLowerMap {
		if len(variants) > 1 {
			return fmt.Errorf("duplicate module name (case-insensitive): '%s' matches %v", lower, variants)
		}
	}

	return nil
}

// InitializeModules should be called at startup to validate modules
func InitializeModules() error {
	if err := ValidateModules(); err != nil {
		return err
	}
	return nil
}
