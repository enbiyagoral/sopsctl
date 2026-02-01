package config

import (
	"os"
	"path/filepath"
	"strings"
)

// ResolveProfile attempts to find a matching profile based on the current directory.
// Returns the profile name and whether auto-selection is enabled.
func (c *Config) ResolveProfile(cwd string) (profileName string, auto bool, found bool) {
	if c.Directories == nil {
		return "", false, false
	}

	// Expand home directory in cwd
	cwd = expandPath(cwd)

	// Check for exact match or parent directory match
	for dirPath, dirMap := range c.Directories {
		expandedDir := expandPath(dirPath)

		// Check if cwd is the directory or a subdirectory
		if cwd == expandedDir || strings.HasPrefix(cwd, expandedDir+string(filepath.Separator)) {
			return dirMap.Profile, dirMap.Auto, true
		}
	}

	return "", false, false
}

// expandPath expands ~ to home directory.
func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err == nil {
			path = filepath.Join(home, path[2:])
		}
	}
	return filepath.Clean(path)
}
