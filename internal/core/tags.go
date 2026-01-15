package core

import "strings"

// Returns modules whose Labels match any of the provided tags.
func FilterModulesByTags(tags []string) []ModuleEntry {
	if len(tags) == 0 {
		return AllModules
	}
	result := []ModuleEntry{}
	for _, entry := range AllModules {
		for _, t := range entry.Metadata.Labels {
			for _, tag := range tags {
				if strings.EqualFold(t, tag) {
					result = append(result, entry)
					break
				}
			}
		}
	}
	return result
}
