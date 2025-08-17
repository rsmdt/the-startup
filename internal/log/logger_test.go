package log

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestWriteSessionLog(t *testing.T) {
	tests := []struct {
		name        string
		sessionID   string
		hookData    *HookData
		expectError bool
		validateFn  func(*testing.T, string)
	}{
		{
			name:      "valid session log write",
			sessionID: "dev-session-123",
			hookData: &HookData{
				Role:      "user",
				Content:   "Design system",
				SessionID: "dev-session-123",
				AgentID:   "arch-001",
				Timestamp: "2025-01-11T10:00:00.000Z",
			},
			expectError: false,
			validateFn: func(t *testing.T, sessionDir string) {
				// Verify session directory was created
				if _, err := os.Stat(sessionDir); os.IsNotExist(err) {
					t.Errorf("Session directory was not created: %s", sessionDir)
				}

				// Verify JSONL file was created and contains correct data
				jsonlFile := filepath.Join(sessionDir, "agent-instructions.jsonl")
				content, err := os.ReadFile(jsonlFile)
				if err != nil {
					t.Errorf("Failed to read JSONL file: %v", err)
					return
				}

				// Parse the JSON line
				var logEntry HookData
				if err := json.Unmarshal(content[:len(content)-1], &logEntry); err != nil { // Remove trailing newline
					t.Errorf("Failed to parse JSONL content: %v", err)
					return
				}

				if logEntry.Role != "user" {
					t.Errorf("Expected role 'user', got %s", logEntry.Role)
				}
				if !strings.Contains(logEntry.Content, "Design system") {
					t.Errorf("Expected content to contain 'Design system', got %s", logEntry.Content)
				}
			},
		},
		{
			name:        "empty session ID - should skip",
			sessionID:   "",
			hookData:    &HookData{Role: "user", Content: "test", Timestamp: "2025-01-11T10:00:00.000Z"},
			expectError: false,
			validateFn: func(t *testing.T, sessionDir string) {
				// For empty session ID, sessionDir will be empty string
				// No directory operations should occur
			},
		},
		{
			name:      "append to existing session file",
			sessionID: "dev-session-456",
			hookData: &HookData{
				Role:      "assistant",
				Content:   "Task completed",
				SessionID: "dev-session-456",
				AgentID:   "dev-001",
				Timestamp: "2025-01-11T10:05:00.000Z",
			},
			expectError: false,
			validateFn: func(t *testing.T, sessionDir string) {
				jsonlFile := filepath.Join(sessionDir, "agent-instructions.jsonl")

				// Write first entry
				firstData := &HookData{
					Role:      "user",
					Content:   "Start task",
					SessionID: "dev-session-456",
					AgentID:   "dev-001",
					Timestamp: "2025-01-11T10:00:00.000Z",
				}
				if err := WriteSessionLog("dev-session-456", firstData); err != nil {
					t.Fatalf("Failed to write first entry: %v", err)
				}

				// Second entry should append
				content, err := os.ReadFile(jsonlFile)
				if err != nil {
					t.Errorf("Failed to read JSONL file: %v", err)
					return
				}

				lines := strings.Split(strings.TrimSpace(string(content)), "\n")
				if len(lines) != 2 {
					t.Errorf("Expected 2 lines in JSONL file, got %d", len(lines))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up temporary project directory
			tempDir := t.TempDir()
			originalEnv := os.Getenv("CLAUDE_PROJECT_DIR")
			defer func() {
				if originalEnv == "" {
					os.Unsetenv("CLAUDE_PROJECT_DIR")
				} else {
					os.Setenv("CLAUDE_PROJECT_DIR", originalEnv)
				}
			}()
			os.Setenv("CLAUDE_PROJECT_DIR", tempDir)

			// Create .the-startup directory in temp dir
			startupDir := filepath.Join(tempDir, ".the-startup")
			os.MkdirAll(startupDir, 0755)

			err := WriteSessionLog(tt.sessionID, tt.hookData)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if tt.validateFn != nil {
				sessionDir := filepath.Join(startupDir, tt.sessionID)
				tt.validateFn(t, sessionDir)
			}
		})
	}
}

// TestWriteGlobalLog has been removed because WriteGlobalLog function no longer exists.
// Only WriteSessionLog function is available in the current implementation.
// Global logging has been removed from the current architecture.

func TestAppendJSONL(t *testing.T) {
	tests := []struct {
		name        string
		setupFn     func(*testing.T, string) string
		hookData    *HookData
		expectError bool
		validateFn  func(*testing.T, string)
	}{
		{
			name: "create new file",
			setupFn: func(t *testing.T, tempDir string) string {
				return filepath.Join(tempDir, "new-file.jsonl")
			},
			hookData: &HookData{
				Role:      "user",
				Content:   "Start new task",
				Timestamp: "2025-01-11T12:00:00.000Z",
				SessionID: "test-session",
				AgentID:   "the-new",
			},
			expectError: false,
			validateFn: func(t *testing.T, filename string) {
				// Verify file was created with correct permissions
				info, err := os.Stat(filename)
				if err != nil {
					t.Errorf("File was not created: %v", err)
					return
				}

				// Check file permissions (0644)
				if info.Mode().Perm() != 0644 {
					t.Errorf("Expected file permissions 0644, got %o", info.Mode().Perm())
				}

				// Verify content
				content, err := os.ReadFile(filename)
				if err != nil {
					t.Errorf("Failed to read file: %v", err)
					return
				}

				if !strings.HasSuffix(string(content), "\n") {
					t.Error("JSONL file should end with newline")
				}
			},
		},
		{
			name: "permission denied scenario",
			setupFn: func(t *testing.T, tempDir string) string {
				// Create read-only directory
				readOnlyDir := filepath.Join(tempDir, "readonly")
				os.MkdirAll(readOnlyDir, 0555) // read and execute only
				return filepath.Join(readOnlyDir, "test.jsonl")
			},
			hookData: &HookData{
				Role:      "user",
				Content:   "Test content",
				Timestamp: "2025-01-11T12:05:00.000Z",
				SessionID: "test-session",
				AgentID:   "the-test",
			},
			expectError: true,
		},
		{
			name: "valid data - nil case",
			setupFn: func(t *testing.T, tempDir string) string {
				return filepath.Join(tempDir, "nil-data.jsonl")
			},
			hookData:    nil,
			expectError: false, // nil marshals to "null" in JSON
			validateFn: func(t *testing.T, filename string) {
				content, err := os.ReadFile(filename)
				if err != nil {
					t.Errorf("Failed to read file: %v", err)
					return
				}
				if string(content) != "null\n" {
					t.Errorf("Expected 'null\\n', got: %q", string(content))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			filename := tt.setupFn(t, tempDir)

			err := appendJSONL(filename, tt.hookData)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if tt.validateFn != nil {
				tt.validateFn(t, filename)
			}
		})
	}
}

func TestFindLatestSessionExtensive(t *testing.T) {
	tests := []struct {
		name     string
		setupFn  func(*testing.T, string)
		expectFn func(*testing.T, string)
	}{
		{
			name: "multiple sessions with timestamps",
			setupFn: func(t *testing.T, tempDir string) {
				startupDir := filepath.Join(tempDir, ".the-startup")
				os.MkdirAll(startupDir, 0755)

				// Create sessions with artificial timestamps
				sessions := []string{"dev-session-old", "dev-session-new", "dev-session-newest"}
				for i, session := range sessions {
					sessionPath := filepath.Join(startupDir, session)
					os.MkdirAll(sessionPath, 0755)

					// Modify the directory timestamp
					modTime := time.Now().Add(time.Duration(i) * time.Minute)
					os.Chtimes(sessionPath, modTime, modTime)
				}

				// Also create non-dev directories that should be ignored
				os.MkdirAll(filepath.Join(startupDir, "not-dev-session"), 0755)
				os.MkdirAll(filepath.Join(startupDir, "other-directory"), 0755)
			},
			expectFn: func(t *testing.T, result string) {
				if !strings.HasPrefix(result, "dev-") {
					t.Errorf("Expected session to start with 'dev-', got: %s", result)
				}
				if result == "" {
					t.Error("Expected to find a session, got empty string")
				}
			},
		},
		{
			name: "no dev- sessions but other directories exist",
			setupFn: func(t *testing.T, tempDir string) {
				startupDir := filepath.Join(tempDir, ".the-startup")
				os.MkdirAll(startupDir, 0755)

				// Create non-dev directories
				os.MkdirAll(filepath.Join(startupDir, "other-session"), 0755)
				os.MkdirAll(filepath.Join(startupDir, "random-dir"), 0755)

				// Create files (not directories) - should be ignored
				file, _ := os.Create(filepath.Join(startupDir, "dev-session-file.txt"))
				file.Close()
			},
			expectFn: func(t *testing.T, result string) {
				if result != "" {
					t.Errorf("Expected empty result when no dev- sessions exist, got: %s", result)
				}
			},
		},
		{
			name: "startup directory doesn't exist",
			setupFn: func(t *testing.T, tempDir string) {
				// Don't create .the-startup directory
				// Also ensure we're not finding the home directory
				// by creating a unique temp home
				tempHome := filepath.Join(tempDir, "fake-home")
				os.MkdirAll(tempHome, 0755)
				os.Setenv("HOME", tempHome)
			},
			expectFn: func(t *testing.T, result string) {
				if result != "" {
					t.Errorf("Expected empty result when startup dir doesn't exist, got: %s", result)
				}
			},
		},
		{
			name: "permission denied on startup directory",
			setupFn: func(t *testing.T, tempDir string) {
				startupDir := filepath.Join(tempDir, ".the-startup")
				os.MkdirAll(startupDir, 0000) // No permissions
			},
			expectFn: func(t *testing.T, result string) {
				if result != "" {
					t.Errorf("Expected empty result when permission denied, got: %s", result)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()

			// Set CLAUDE_PROJECT_DIR to isolate test from real home directory
			oldProjectDir := os.Getenv("CLAUDE_PROJECT_DIR")
			oldHome := os.Getenv("HOME")
			os.Setenv("CLAUDE_PROJECT_DIR", tempDir)
			defer func() {
				os.Setenv("CLAUDE_PROJECT_DIR", oldProjectDir)
				os.Setenv("HOME", oldHome)
			}()

			tt.setupFn(t, tempDir)

			result := FindLatestSession(tempDir)
			tt.expectFn(t, result)
		})
	}
}

func TestEnsureDirectories(t *testing.T) {
	tests := []struct {
		name        string
		setupFn     func(*testing.T, string) string
		expectError bool
		validateFn  func(*testing.T, string)
	}{
		{
			name: "create new directory structure",
			setupFn: func(t *testing.T, tempDir string) string {
				return filepath.Join(tempDir, "new", "nested", "directories")
			},
			expectError: false,
			validateFn: func(t *testing.T, dirPath string) {
				if _, err := os.Stat(dirPath); os.IsNotExist(err) {
					t.Errorf("Directory structure was not created: %s", dirPath)
				}

				// Check permissions
				info, err := os.Stat(dirPath)
				if err != nil {
					t.Errorf("Failed to stat directory: %v", err)
					return
				}

				if info.Mode().Perm() != 0755 {
					t.Errorf("Expected directory permissions 0755, got %o", info.Mode().Perm())
				}
			},
		},
		{
			name: "directory already exists",
			setupFn: func(t *testing.T, tempDir string) string {
				existingDir := filepath.Join(tempDir, "existing")
				os.MkdirAll(existingDir, 0755)
				return existingDir
			},
			expectError: false,
			validateFn: func(t *testing.T, dirPath string) {
				if _, err := os.Stat(dirPath); os.IsNotExist(err) {
					t.Errorf("Existing directory should still exist: %s", dirPath)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			dirPath := tt.setupFn(t, tempDir)

			err := EnsureDirectories(dirPath)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if tt.validateFn != nil {
				tt.validateFn(t, dirPath)
			}
		})
	}
}
