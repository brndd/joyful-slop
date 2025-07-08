package config

import "slices"

// validateModes checks the provided modes against a larger subset of modes (usually all defined ones)
// and returns false if any of the modes are not defined.
func validateModes(modes []string, allModes []string) bool {
	if len(modes) == 0 {
		return true
	}

	for _, mode := range modes {
		if !slices.Contains(allModes, mode) {
			return false
		}
	}

	return true
}
