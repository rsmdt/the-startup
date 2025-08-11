package log

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	// Mutex to protect concurrent file writes
	fileMutex sync.Mutex
)

// WriteSessionLog writes hook data to session-specific JSONL file
func WriteSessionLog(sessionID string, data *HookData) error {
	if sessionID == "" {
		return nil // Skip if no session ID
	}

	projectDir := GetProjectDir()
	startupDir := GetStartupDir(projectDir)

	// Ensure startup directory exists first
	if err := EnsureDirectories(startupDir); err != nil {
		return fmt.Errorf("failed to ensure directories: %w", err)
	}

	// Create session directory
	sessionDir := filepath.Join(startupDir, sessionID)
	if err := os.MkdirAll(sessionDir, 0755); err != nil {
		return fmt.Errorf("failed to create session directory: %w", err)
	}

	// Write to session-specific file
	sessionFile := filepath.Join(sessionDir, "agent-instructions.jsonl")
	return appendJSONL(sessionFile, data)
}

// WriteGlobalLog writes hook data to global JSONL file
func WriteGlobalLog(data *HookData) error {
	projectDir := GetProjectDir()
	startupDir := GetStartupDir(projectDir)

	// Ensure startup directory exists
	if err := EnsureDirectories(startupDir); err != nil {
		return fmt.Errorf("failed to ensure directories: %w", err)
	}

	// Write to global file
	globalFile := filepath.Join(startupDir, "all-agent-instructions.jsonl")
	return appendJSONL(globalFile, data)
}

// EnsureDirectories creates the required directories if they don't exist
func EnsureDirectories(startupDir string) error {
	return os.MkdirAll(startupDir, 0755)
}

// FindLatestSession finds the most recent session directory
func FindLatestSession(projectDir string) string {
	startupDir := GetStartupDir(projectDir)
	
	entries, err := os.ReadDir(startupDir)
	if err != nil {
		return ""
	}

	var latestSession string
	var latestModTime int64

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		
		name := entry.Name()
		if !strings.HasPrefix(name, "dev-") {
			continue
		}

		// Get modification time
		info, err := entry.Info()
		if err != nil {
			continue
		}

		if info.ModTime().Unix() > latestModTime {
			latestModTime = info.ModTime().Unix()
			latestSession = name
		}
	}

	return latestSession
}

// appendJSONL appends a JSON object as a line to the specified file with proper locking
func appendJSONL(filename string, data *HookData) error {
	fileMutex.Lock()
	defer fileMutex.Unlock()

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", filename, err)
	}
	defer file.Close()

	// Marshal to JSON
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Write JSON and newline as a single operation to ensure atomicity
	line := string(jsonBytes) + "\n"
	if _, err := file.WriteString(line); err != nil {
		return fmt.Errorf("failed to write JSON line: %w", err)
	}

	return nil
}