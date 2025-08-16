package log

import (
	"os"
	"path/filepath"
)

// GetProjectDir returns the project directory from environment or current directory
func GetProjectDir() string {
	projectDir := os.Getenv("CLAUDE_PROJECT_DIR")
	if projectDir == "" {
		projectDir = "."
	}
	return projectDir
}

// GetStartupDir determines the startup directory for logs
// Follows Python logic: check project-local .the-startup first, fallback to ~/.the-startup
func GetStartupDir(projectDir string) string {
	// First check project local directory
	localStartup := filepath.Join(projectDir, ".the-startup")
	if dirExists(localStartup) {
		return localStartup
	}

	// Otherwise use home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fallback to current directory if home dir not accessible
		return filepath.Join(".", ".the-startup")
	}

	homeStartup := filepath.Join(homeDir, ".the-startup")

	// Don't auto-create the home directory - let the caller decide
	// This prevents test pollution and gives more control
	return homeStartup
}

// IsDebugEnabled checks if debug output is enabled via DEBUG_HOOKS environment variable
func IsDebugEnabled() bool {
	return os.Getenv("DEBUG_HOOKS") != ""
}

// dirExists checks if a directory exists
func dirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}
