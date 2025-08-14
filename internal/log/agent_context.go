package log

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// AgentFileManager handles agent-specific file operations with thread-safe access.
// It encapsulates the file paths and operations for a specific agent within a session.
type AgentFileManager struct {
	SessionID string      // Current Claude Code session identifier (e.g., "dev-20250812-143022")
	AgentID   string      // Unique agent instance identifier (e.g., "arch-001")
	BaseDir   string      // Base directory path (typically ".the-startup")
	mutex     sync.Mutex  // File operation synchronization for concurrent access safety
}

// GetAgentFilePath returns the path to the agent's JSONL file
func (afm *AgentFileManager) GetAgentFilePath() string {
	return filepath.Join(afm.BaseDir, afm.SessionID, afm.AgentID+".jsonl")
}

// GetSessionDir returns the session directory path
func (afm *AgentFileManager) GetSessionDir() string {
	return filepath.Join(afm.BaseDir, afm.SessionID)
}

// EnsureSessionDir creates the session directory with proper permissions
func (afm *AgentFileManager) EnsureSessionDir() error {
	sessionDir := afm.GetSessionDir()
	return os.MkdirAll(sessionDir, 0755)
}

// WriteAgentContext writes context to agent-specific file with atomic operations
func WriteAgentContext(sessionID, agentID string, data *HookData) error {
	// Validation
	if sessionID == "" || agentID == "" || data == nil {
		return fmt.Errorf("invalid parameters: sessionID, agentID, and data required")
	}
	
	projectDir := GetProjectDir()
	startupDir := GetStartupDir(projectDir)
	
	// File path generation
	agentFile := filepath.Join(startupDir, sessionID, agentID+".jsonl")
	
	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(agentFile), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	
	// Use existing appendJSONL function which has atomic operations and mutex protection
	return appendJSONL(agentFile, data)
}


// findLatestSessionWithAgent locates the most recent session containing the agent
func findLatestSessionWithAgent(agentID string) (string, error) {
	projectDir := GetProjectDir()
	baseDir := GetStartupDir(projectDir)
	
	entries, err := os.ReadDir(baseDir)
	if err != nil {
		return "", err
	}
	
	var latestSession string
	var latestTime time.Time
	
	for _, entry := range entries {
		if !entry.IsDir() || !strings.HasPrefix(entry.Name(), "dev-") {
			continue
		}
		
		agentFile := filepath.Join(baseDir, entry.Name(), agentID+".jsonl")
		if stat, err := os.Stat(agentFile); err == nil {
			if stat.ModTime().After(latestTime) {
				latestTime = stat.ModTime()
				latestSession = entry.Name()
			}
		}
	}
	
	if latestSession == "" {
		return "", fmt.Errorf("no session found for agent %s", agentID)
	}
	
	return latestSession, nil
}

// readLastNLines efficiently reads last N lines from large files
func readLastNLines(filename string, n int) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	
	size := stat.Size()
	if size == 0 {
		return []string{}, nil
	}
	
	// For small files (< 1MB), read entire content for simplicity and performance
	const smallFileThreshold = 1024 * 1024 // 1MB
	if size < smallFileThreshold {
		scanner := bufio.NewScanner(file)
		var allLines []string
		for scanner.Scan() {
			allLines = append(allLines, scanner.Text())
		}
		
		start := len(allLines) - n
		if start < 0 {
			start = 0
		}
		return allLines[start:], nil
	}
	
	// For large files, use reverse reading with buffer
	return readTailWithBuffer(file, size, n)
}

// readTailWithBuffer implements efficient tail reading for large files using reverse buffered reading
func readTailWithBuffer(file *os.File, size int64, n int) ([]string, error) {
	const bufferSize = 8192 // 8KB buffer size for optimal I/O performance
	bufSize := int64(bufferSize)
	var lines []string
	var remaining []byte
	
	for offset := size; offset > 0 && len(lines) < n; {
		readSize := bufSize
		if offset < bufSize {
			readSize = offset
		}
		
		offset -= readSize
		buffer := make([]byte, readSize)
		
		_, err := file.ReadAt(buffer, offset)
		if err != nil && err != io.EOF {
			return nil, err
		}
		
		// Combine with remaining bytes from previous iteration
		if len(remaining) > 0 {
			buffer = append(buffer, remaining...)
		}
		
		// Split by newlines and process
		parts := bytes.Split(buffer, []byte{'\n'})
		remaining = parts[0] // First part might be incomplete
		
		// Process complete lines in reverse order
		for i := len(parts) - 1; i >= 1; i-- {
			if len(parts[i]) > 0 && len(lines) < n {
				lines = append([]string{string(parts[i])}, lines...)
			}
		}
	}
	
	// Add the remaining incomplete line if it exists
	if len(remaining) > 0 && len(lines) < n {
		lines = append([]string{string(remaining)}, lines...)
	}
	
	return lines, nil
}

// ReadAgentContextRaw reads raw JSONL lines from agent context file
func ReadAgentContextRaw(sessionID, agentID string, maxLines int) ([]string, error) {
	// Validation
	if agentID == "" {
		return nil, fmt.Errorf("agentID is required")
	}
	
	// Set default maxLines
	if maxLines <= 0 || maxLines > 1000 {
		maxLines = 50
	}
	
	// Determine session if not provided
	if sessionID == "" {
		var err error
		sessionID, err = findLatestSessionWithAgent(agentID)
		if err != nil || sessionID == "" {
			return []string{}, nil // Return empty if no session found
		}
	}
	
	// Build file path
	projectDir := GetProjectDir()
	startupDir := GetStartupDir(projectDir)
	agentFile := filepath.Join(startupDir, sessionID, agentID+".jsonl")
	
	// Read the last N lines from the file
	lines, err := readLastNLines(agentFile, maxLines)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil // Return empty for non-existent files
		}
		return nil, err
	}
	
	return lines, nil
}
