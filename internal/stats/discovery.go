package stats

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// LogDiscovery implements the Discovery interface for finding Claude Code log files
type LogDiscovery struct {
	// HomeDir is the user's home directory (can be overridden for testing)
	HomeDir string
}

// NewLogDiscovery creates a new LogDiscovery instance with the user's home directory
func NewLogDiscovery() *LogDiscovery {
	homeDir, _ := os.UserHomeDir()
	return &LogDiscovery{
		HomeDir: homeDir,
	}
}

// NewLogDiscoveryWithHome creates a new LogDiscovery instance with a custom home directory
// This is useful for testing or when Claude's config is in a non-standard location
func NewLogDiscoveryWithHome(homeDir string) *LogDiscovery {
	return &LogDiscovery{
		HomeDir: homeDir,
	}
}

// FindLogFiles discovers Claude Code JSONL log files
// The options.StartTime field supports the --since flag for time-based filtering
func (d *LogDiscovery) FindLogFiles(projectPath string, options FilterOptions) ([]string, error) {
	// If projectPath is empty, try to get current project
	if projectPath == "" {
		projectPath = d.GetCurrentProject()
		if projectPath == "" {
			return nil, fmt.Errorf("no project path specified and could not detect current project")
		}
	} else if strings.HasPrefix(projectPath, "/") || !strings.HasPrefix(projectPath, "-") {
		// If it looks like a filesystem path (starts with / or doesn't start with -), sanitize it
		projectPath = sanitizeProjectPath(projectPath)
	}

	// Get the Claude projects directory
	claudeProjectsDir := filepath.Join(d.HomeDir, ".claude", "projects", projectPath)
	
	// Check if the directory exists
	if _, err := os.Stat(claudeProjectsDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("Claude project directory does not exist: %s", claudeProjectsDir)
	}

	// Find all JSONL files in the directory
	matches, err := filepath.Glob(filepath.Join(claudeProjectsDir, "*.jsonl"))
	if err != nil {
		return nil, fmt.Errorf("error searching for log files: %w", err)
	}

	// Apply time-based filtering if StartTime is specified
	var filtered []string
	for _, match := range matches {
		// Parse the timestamp from the filename if possible
		if options.StartTime != nil {
			fileTime := extractTimestampFromPath(match)
			if fileTime != nil && fileTime.Before(*options.StartTime) {
				continue // Skip files older than StartTime
			}
		}
		
		if options.EndTime != nil {
			fileTime := extractTimestampFromPath(match)
			if fileTime != nil && fileTime.After(*options.EndTime) {
				continue // Skip files newer than EndTime
			}
		}
		
		filtered = append(filtered, match)
	}

	// Sort files by name (which should sort chronologically for timestamp-based names)
	sort.Strings(filtered)

	return filtered, nil
}

// GetCurrentProject attempts to detect the current project from working directory
func (d *LogDiscovery) GetCurrentProject() string {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return ""
	}

	// Sanitize the path to match Claude's format
	sanitized := sanitizeProjectPath(cwd)
	
	// Check if this project exists in Claude's projects directory
	claudeProjectPath := filepath.Join(d.HomeDir, ".claude", "projects", sanitized)
	if _, err := os.Stat(claudeProjectPath); err == nil {
		return sanitized
	}

	return ""
}

// ValidateProjectPath checks if a path contains Claude Code logs
func (d *LogDiscovery) ValidateProjectPath(path string) bool {
	// If the path is already a sanitized project name, check directly
	claudeProjectPath := filepath.Join(d.HomeDir, ".claude", "projects", path)
	if _, err := os.Stat(claudeProjectPath); err == nil {
		// Check if there are any JSONL files
		matches, _ := filepath.Glob(filepath.Join(claudeProjectPath, "*.jsonl"))
		return len(matches) > 0
	}

	// If it's a full path, sanitize it first
	sanitized := sanitizeProjectPath(path)
	claudeProjectPath = filepath.Join(d.HomeDir, ".claude", "projects", sanitized)
	if _, err := os.Stat(claudeProjectPath); err == nil {
		// Check if there are any JSONL files
		matches, _ := filepath.Glob(filepath.Join(claudeProjectPath, "*.jsonl"))
		return len(matches) > 0
	}

	return false
}

// sanitizeProjectPath converts a file system path to Claude's sanitized format
// For example: /Users/irudi/Code/personal/the-startup -> -Users-irudi-Code-personal-the-startup
func sanitizeProjectPath(path string) string {
	// Replace all slashes with hyphens
	// Note: Claude keeps the leading slash as a dash
	sanitized := strings.ReplaceAll(path, "/", "-")
	
	return sanitized
}

// extractTimestampFromPath attempts to extract a timestamp from a log file path
// It looks for common date patterns in the filename
func extractTimestampFromPath(path string) *time.Time {
	// Get just the filename without extension
	base := filepath.Base(path)
	base = strings.TrimSuffix(base, ".jsonl")
	
	// Common timestamp patterns in Claude log filenames
	// Pattern: YYYY-MM-DD-session.jsonl or YYYY-MM-DD-HH-MM-SS.jsonl
	patterns := []string{
		"2006-01-02-15-04-05", // YYYY-MM-DD-HH-MM-SS
		"2006-01-02",          // YYYY-MM-DD
	}
	
	// Try to parse the beginning of the filename as a timestamp
	parts := strings.Split(base, "-")
	if len(parts) >= 3 {
		// Try to reconstruct a timestamp string from the parts
		var timestampStr string
		if len(parts) >= 6 {
			// Try full timestamp with time
			timestampStr = strings.Join(parts[:6], "-")
			if t, err := time.Parse(patterns[0], timestampStr); err == nil {
				return &t
			}
		}
		if len(parts) >= 3 {
			// Try just date
			timestampStr = strings.Join(parts[:3], "-")
			if t, err := time.Parse(patterns[1], timestampStr); err == nil {
				return &t
			}
		}
	}
	
	// Try to get file modification time as fallback
	if info, err := os.Stat(path); err == nil {
		modTime := info.ModTime()
		return &modTime
	}
	
	return nil
}