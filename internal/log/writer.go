package log

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// fileMutex provides thread-safe access to JSONL files
var fileMutex sync.Mutex

// AppendMetrics writes a MetricsEntry to the daily JSONL file
// Creates the logs directory and daily file if they don't exist
// Returns nil on any error (silent failure mode)
func AppendMetrics(entry *MetricsEntry) error {
	if entry == nil {
		return nil
	}

	// Lock for thread-safe file operations
	fileMutex.Lock()
	defer fileMutex.Unlock()

	// Get the startup installation path
	startupPath, err := getStartupPath()
	if err != nil {
		// Silent failure
		return nil
	}

	// Create logs directory path
	logsDir := filepath.Join(startupPath, "logs")
	
	// Ensure logs directory exists
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		// Silent failure
		return nil
	}

	// Generate daily filename (YYYYMMDD.jsonl)
	filename := fmt.Sprintf("%s.jsonl", entry.Timestamp.Format("20060102"))
	filePath := filepath.Join(logsDir, filename)

	// Open file in append mode (creates if doesn't exist)
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		// Silent failure
		return nil
	}
	defer file.Close()

	// Marshal entry to JSON
	data, err := json.Marshal(entry)
	if err != nil {
		// Silent failure
		return nil
	}

	// Write JSON line with newline
	if _, err := file.Write(append(data, '\n')); err != nil {
		// Silent failure
		return nil
	}

	// Sync to ensure data is written to disk
	if err := file.Sync(); err != nil {
		// Silent failure
		return nil
	}

	return nil
}

// getStartupPath finds the installation directory for the-startup
// Tries multiple strategies in order of preference
func getStartupPath() (string, error) {
	// Strategy 0: Check environment variable (for testing)
	if envPath := os.Getenv("THE_STARTUP_PATH"); envPath != "" {
		return envPath, nil
	}
	
	// Strategy 1: Check for local .the-startup directory first (project-local takes precedence)
	if _, err := os.Stat(".the-startup"); err == nil {
		return ".the-startup", nil
	}
	
	// Strategy 2: Check for lockfile in common locations
	lockfilePaths := []string{
		expandHome("~/.the-startup/the-startup.lock"), // home directory
		expandHome("~/.config/the-startup/the-startup.lock"), // XDG config
	}

	for _, lockfilePath := range lockfilePaths {
		if lockfilePath == "" {
			continue
		}
		if _, err := os.Stat(lockfilePath); err == nil {
			// Found lockfile, extract directory
			return filepath.Dir(lockfilePath), nil
		}
	}

	// Strategy 3: Check for evidence of installation (bin directory)
	installPaths := []string{
		expandHome("~/.the-startup"),
		expandHome("~/.config/the-startup"),
	}

	for _, path := range installPaths {
		if path == "" {
			continue
		}
		// Check for bin/the-startup executable
		binaryPath := filepath.Join(path, "bin", "the-startup")
		if _, err := os.Stat(binaryPath); err == nil {
			return path, nil
		}
		// Check for templates directory
		templatesPath := filepath.Join(path, "templates")
		if _, err := os.Stat(templatesPath); err == nil {
			return path, nil
		}
	}

	// Strategy 4: Default to project-local installation
	// Create the directory if it doesn't exist
	defaultPath := ".the-startup"
	if err := os.MkdirAll(defaultPath, 0755); err != nil {
		return "", fmt.Errorf("could not find or create startup directory")
	}
	return defaultPath, nil
}

// expandHome expands ~ to the user's home directory
// Returns empty string on error (for silent failure)
func expandHome(path string) string {
	if !strings.HasPrefix(path, "~/") {
		return path
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	return filepath.Join(homeDir, path[2:])
}