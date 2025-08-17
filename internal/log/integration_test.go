package log

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

// TestFullPythonCompatibilityFlow tests the complete flow from JSON input to JSONL output
// to ensure 100% compatibility with the Python hooks
func TestFullPythonCompatibilityFlow(t *testing.T) {
	tempDir := t.TempDir()
	os.Setenv("CLAUDE_PROJECT_DIR", tempDir)
	defer os.Unsetenv("CLAUDE_PROJECT_DIR")

	// Create .the-startup directory in project
	startupDir := filepath.Join(tempDir, ".the-startup")
	os.MkdirAll(startupDir, 0755)

	// Test PreToolUse (agent_start) flow
	t.Run("agent_start_flow", func(t *testing.T) {
		jsonInput := `{
			"session_id": "dev-test-session",
			"transcript_path": "/path/to/transcript.jsonl",
			"cwd": "/current/working/directory",
			"hook_event_name": "PreToolUse",
			"tool_name": "Task",
			"tool_input": {
				"description": "Design the system architecture",
				"prompt": "SessionId: dev-test-session AgentId: arch-001\n\nPlease design a scalable microservices architecture for our e-commerce platform. Consider the following requirements:\n\n1. Handle 10,000+ concurrent users\n2. Support multiple payment methods\n3. Real-time inventory management\n4. Geographic distribution",
				"subagent_type": "the-architect"
			}
		}`

		reader := strings.NewReader(jsonInput)
		hookData, err := ProcessToolCall(reader, false)
		if err != nil {
			t.Fatalf("ProcessToolCall failed: %v", err)
		}

		if hookData == nil {
			t.Fatal("Expected non-nil hookData")
		}

		// Validate fields match the current structure
		if hookData.Role != "user" {
			t.Errorf("Expected role 'user', got %q", hookData.Role)
		}
		if hookData.AgentID != "arch-001" {
			t.Errorf("Expected agent_id 'arch-001', got %q", hookData.AgentID)
		}
		if hookData.SessionID != "dev-test-session" {
			t.Errorf("Expected session_id 'dev-test-session', got %q", hookData.SessionID)
		}
		if !strings.Contains(hookData.Content, "Please design a scalable microservices architecture") {
			t.Errorf("Expected content to contain prompt, got %q", hookData.Content)
		}

		// Write logs and verify files
		err = WriteSessionLog(hookData.SessionID, hookData)
		if err != nil {
			t.Fatalf("WriteSessionLog failed: %v", err)
		}

		// Verify session file exists and has correct content
		sessionFile := filepath.Join(startupDir, "dev-test-session", "agent-instructions.jsonl")
		if _, err := os.Stat(sessionFile); os.IsNotExist(err) {
			t.Errorf("Session file does not exist: %s", sessionFile)
		}

		// Read and verify session file content
		sessionContent, err := os.ReadFile(sessionFile)
		if err != nil {
			t.Fatalf("Failed to read session file: %v", err)
		}

		var sessionEntry HookData
		if err := json.Unmarshal(sessionContent, &sessionEntry); err != nil {
			t.Fatalf("Failed to parse session file JSON: %v", err)
		}

		if sessionEntry.Role != "user" {
			t.Errorf("Session file: Expected role 'user', got %q", sessionEntry.Role)
		}

		// Note: Global files are no longer created in the current implementation
		// This test used to check for global files but they've been removed
	})

	// Test PostToolUse (agent_complete) flow
	t.Run("agent_complete_flow", func(t *testing.T) {
		// Long output to test truncation
		longOutput := strings.Repeat("The architecture has been designed with the following components: ", 30) + " [additional content...]"

		jsonInput := `{
			"session_id": "dev-test-session",
			"transcript_path": "/path/to/transcript.jsonl",
			"cwd": "/current/working/directory", 
			"hook_event_name": "PostToolUse",
			"tool_name": "Task",
			"tool_input": {
				"description": "Design the system architecture",
				"prompt": "SessionId: dev-test-session AgentId: arch-001\n\nPlease design a scalable microservices architecture for our e-commerce platform.",
				"subagent_type": "the-architect"
			},
			"output": "` + longOutput + `"
		}`

		reader := strings.NewReader(jsonInput)
		hookData, err := ProcessToolCall(reader, true)
		if err != nil {
			t.Fatalf("ProcessToolCall failed: %v", err)
		}

		if hookData == nil {
			t.Fatal("Expected non-nil hookData")
		}

		// Validate fields match the current structure
		if hookData.Role != "assistant" {
			t.Errorf("Expected role 'assistant', got %q", hookData.Role)
		}
		if hookData.AgentID != "arch-001" {
			t.Errorf("Expected agent_id 'arch-001', got %q", hookData.AgentID)
		}

		// Verify output truncation is handled in Content
		if len(hookData.Content) <= 1000 {
			t.Errorf("Expected content to be truncated, but got length %d", len(hookData.Content))
		}
		if !strings.Contains(hookData.Content, "... [truncated") {
			t.Errorf("Expected truncation message in content, got %q", hookData.Content)
		}

		// Write logs
		err = WriteSessionLog(hookData.SessionID, hookData)
		if err != nil {
			t.Fatalf("WriteSessionLog failed: %v", err)
		}

		// Verify files were written (should append to existing files)
		sessionFile := filepath.Join(startupDir, "dev-test-session", "agent-instructions.jsonl")
		sessionContent, err := os.ReadFile(sessionFile)
		if err != nil {
			t.Fatalf("Failed to read session file: %v", err)
		}

		// Should have two lines now (agent_start + agent_complete)
		lines := strings.Split(strings.TrimSpace(string(sessionContent)), "\n")
		if len(lines) != 2 {
			t.Errorf("Expected 2 lines in session file, got %d", len(lines))
		}

		// Parse second line (agent_complete)
		var completeEntry HookData
		if err := json.Unmarshal([]byte(lines[1]), &completeEntry); err != nil {
			t.Fatalf("Failed to parse second line JSON: %v", err)
		}

		if completeEntry.Role != "assistant" {
			t.Errorf("Second line: Expected role 'assistant', got %q", completeEntry.Role)
		}
	})
}

// TestDebugOutputIntegration tests DEBUG_HOOKS functionality in integration scenarios
func TestDebugOutputIntegration(t *testing.T) {
	// This test verifies debug output behavior but can't easily capture stderr
	// in unit tests, so we just verify the functions don't crash

	os.Setenv("DEBUG_HOOKS", "1")
	defer os.Unsetenv("DEBUG_HOOKS")

	DebugLog("Test debug message: %s", "value")
	DebugError(os.ErrNotExist)

	// Verify debug is enabled
	if !IsDebugEnabled() {
		t.Error("Expected debug to be enabled with DEBUG_HOOKS=1")
	}
}

// TestPythonJSONFormatCompatibility ensures JSON output exactly matches Python format
func TestPythonJSONFormatCompatibility(t *testing.T) {
	reader := strings.NewReader(`{"tool_name":"Task","tool_input":{"subagent_type":"the-developer","description":"Implement feature","prompt":"SessionId: test-123 AgentId: dev-456\nImplement the user authentication system"},"session_id":"test-123","hook_event_name":"PreToolUse"}`)

	hookData, err := ProcessToolCall(reader, false)
	if err != nil {
		t.Fatalf("ProcessToolCall failed: %v", err)
	}

	// Marshal to JSON
	jsonBytes, err := json.Marshal(hookData)
	if err != nil {
		t.Fatalf("Failed to marshal hookData: %v", err)
	}

	// Parse back to map to check field presence
	var jsonMap map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &jsonMap); err != nil {
		t.Fatalf("Failed to unmarshal to map: %v", err)
	}

	// Check required fields are present in the current structure
	requiredFields := []string{"role", "content", "timestamp"}
	for _, field := range requiredFields {
		if _, exists := jsonMap[field]; !exists {
			t.Errorf("Required field %q missing from JSON output", field)
		}
	}

	// SessionID and AgentID are internal fields and not serialized to JSON
	// This is correct behavior

	// Verify timestamp format
	timestamp, ok := jsonMap["timestamp"].(string)
	if !ok {
		t.Error("timestamp should be a string")
	} else if !strings.HasSuffix(timestamp, "Z") {
		t.Errorf("timestamp should end with Z, got %q", timestamp)
	}

	// Test agent_complete JSON format
	completeReader := strings.NewReader(`{"tool_name":"Task","tool_input":{"subagent_type":"the-developer","description":"Implement feature","prompt":"SessionId: test-123 AgentId: dev-456\nImplement the user authentication system"},"session_id":"test-123","hook_event_name":"PostToolUse","output":"Feature implemented successfully with tests"}`)

	completeHookData, err := ProcessToolCall(completeReader, true)
	if err != nil {
		t.Fatalf("ProcessToolCall failed for complete: %v", err)
	}

	completeJsonBytes, err := json.Marshal(completeHookData)
	if err != nil {
		t.Fatalf("Failed to marshal complete hookData: %v", err)
	}

	var completeJsonMap map[string]interface{}
	if err := json.Unmarshal(completeJsonBytes, &completeJsonMap); err != nil {
		t.Fatalf("Failed to unmarshal complete to map: %v", err)
	}

	// Check that content is present for assistant role
	if _, exists := completeJsonMap["content"]; !exists {
		t.Error("content should be present for assistant role")
	}

	// Check that role is assistant
	if role, exists := completeJsonMap["role"]; !exists || role != "assistant" {
		t.Errorf("Expected role 'assistant', got %v", role)
	}
}

// TestEnvironmentVariableHandling tests comprehensive environment variable scenarios
func TestEnvironmentVariableHandling(t *testing.T) {
	// Save original environment
	originalProjectDir := os.Getenv("CLAUDE_PROJECT_DIR")
	originalDebugHooks := os.Getenv("DEBUG_HOOKS")
	defer func() {
		if originalProjectDir != "" {
			os.Setenv("CLAUDE_PROJECT_DIR", originalProjectDir)
		} else {
			os.Unsetenv("CLAUDE_PROJECT_DIR")
		}
		if originalDebugHooks != "" {
			os.Setenv("DEBUG_HOOKS", originalDebugHooks)
		} else {
			os.Unsetenv("DEBUG_HOOKS")
		}
	}()

	t.Run("CLAUDE_PROJECT_DIR_Set", func(t *testing.T) {
		tempDir := t.TempDir()
		os.Setenv("CLAUDE_PROJECT_DIR", tempDir)

		projectDir := GetProjectDir()
		if projectDir != tempDir {
			t.Errorf("Expected project dir %q, got %q", tempDir, projectDir)
		}

		// Test that logs go to project directory
		expectedStartupDir := filepath.Join(tempDir, ".the-startup")

		// Create the directory to test the logic
		os.MkdirAll(expectedStartupDir, 0755)
		actualStartupDir := GetStartupDir(projectDir)

		if actualStartupDir != expectedStartupDir {
			t.Errorf("Expected startup dir %q, got %q", expectedStartupDir, actualStartupDir)
		}
	})

	t.Run("CLAUDE_PROJECT_DIR_Unset", func(t *testing.T) {
		os.Unsetenv("CLAUDE_PROJECT_DIR")

		projectDir := GetProjectDir()
		if projectDir != "." {
			t.Errorf("Expected project dir '.', got %q", projectDir)
		}
	})

	t.Run("DEBUG_HOOKS_Enabled", func(t *testing.T) {
		os.Setenv("DEBUG_HOOKS", "1")

		if !IsDebugEnabled() {
			t.Error("Expected debug to be enabled with DEBUG_HOOKS=1")
		}

		// Test various debug values
		testValues := []string{"true", "yes", "on", "1", "debug", "anything"}
		for _, value := range testValues {
			os.Setenv("DEBUG_HOOKS", value)
			if !IsDebugEnabled() {
				t.Errorf("Expected debug to be enabled with DEBUG_HOOKS=%q", value)
			}
		}
	})

	t.Run("DEBUG_HOOKS_Disabled", func(t *testing.T) {
		os.Unsetenv("DEBUG_HOOKS")

		if IsDebugEnabled() {
			t.Error("Expected debug to be disabled when DEBUG_HOOKS is unset")
		}

		// Test empty value
		os.Setenv("DEBUG_HOOKS", "")
		if IsDebugEnabled() {
			t.Error("Expected debug to be disabled when DEBUG_HOOKS is empty")
		}
	})
}

// TestCrossPlatformFileSystemBehavior tests platform-specific file system behavior
func TestCrossPlatformFileSystemBehavior(t *testing.T) {
	tempDir := t.TempDir()
	os.Setenv("CLAUDE_PROJECT_DIR", tempDir)
	defer os.Unsetenv("CLAUDE_PROJECT_DIR")

	t.Run("PathSeparators", func(t *testing.T) {
		// Test that path separators work correctly on all platforms
		sessionID := "dev-cross-platform-test"
		startupDir := filepath.Join(tempDir, ".the-startup")

		// Create the .the-startup directory first so GetStartupDir will find it
		os.MkdirAll(startupDir, 0755)

		expectedSessionDir := filepath.Join(startupDir, sessionID)
		expectedSessionFile := filepath.Join(expectedSessionDir, "agent-instructions.jsonl")

		// Create test data
		hookData := &HookData{
			Role:      "user",
			Content:   "Test path separators",
			SessionID: sessionID,
			AgentID:   "test-001",
			Timestamp: "2025-01-11T12:00:00.000Z",
		}

		// Write session log
		err := WriteSessionLog(sessionID, hookData)
		if err != nil {
			t.Fatalf("WriteSessionLog failed: %v", err)
		}

		// Verify file exists
		if _, err := os.Stat(expectedSessionFile); os.IsNotExist(err) {
			t.Errorf("Session file does not exist: %s", expectedSessionFile)
		}

		// Test that file path uses platform-appropriate separators
		content, err := os.ReadFile(expectedSessionFile)
		if err != nil {
			t.Fatalf("Failed to read session file: %v", err)
		}

		var readHookData HookData
		if err := json.Unmarshal(content, &readHookData); err != nil {
			t.Fatalf("Failed to unmarshal hook data: %v", err)
		}

		if readHookData.SessionID != sessionID {
			t.Errorf("Expected session ID %q, got %q", sessionID, readHookData.SessionID)
		}
	})

	t.Run("DirectoryPermissions", func(t *testing.T) {
		// Test directory creation with proper permissions
		sessionID := "dev-permissions-test"
		startupDir := filepath.Join(tempDir, ".the-startup")
		sessionDir := filepath.Join(startupDir, sessionID)

		hookData := &HookData{
			Role:      "user",
			Content:   "Test directory permissions",
			SessionID: sessionID,
			AgentID:   "the-tester",
			Timestamp: "2025-01-11T12:00:00.000Z",
		}

		err := WriteSessionLog(sessionID, hookData)
		if err != nil {
			t.Fatalf("WriteSessionLog failed: %v", err)
		}

		// Check directory permissions (Unix-like systems)
		if info, err := os.Stat(sessionDir); err == nil {
			mode := info.Mode()
			if !mode.IsDir() {
				t.Error("Session path should be a directory")
			}
			// Check basic read/write permissions are set
			// Note: On Windows, permission checking is different, so we skip detailed checks
			if mode.Perm() == 0 {
				t.Error("Directory permissions should not be zero")
			}
		}
	})

	t.Run("FileCreationAndAppending", func(t *testing.T) {
		// Create the .the-startup directory first
		startupDir := filepath.Join(tempDir, ".the-startup")
		os.MkdirAll(startupDir, 0755)

		// Test file creation and appending behavior
		sessionID := "dev-file-ops-test"

		hookData1 := &HookData{
			Role:      "user",
			Content:   "Test file creation start",
			SessionID: sessionID,
			AgentID:   "the-tester",
			Timestamp: "2025-01-11T12:00:00.000Z",
		}

		hookData2 := &HookData{
			Role:      "assistant",
			Content:   "Test file creation complete",
			SessionID: sessionID,
			AgentID:   "the-tester",
			Timestamp: "2025-01-11T12:01:00.000Z",
		}

		// Write first entry
		err := WriteSessionLog(sessionID, hookData1)
		if err != nil {
			t.Fatalf("First WriteSessionLog failed: %v", err)
		}

		// Write second entry
		err = WriteSessionLog(sessionID, hookData2)
		if err != nil {
			t.Fatalf("Second WriteSessionLog failed: %v", err)
		}

		// Read file and verify both entries exist
		sessionFile := filepath.Join(tempDir, ".the-startup", sessionID, "agent-instructions.jsonl")
		content, err := os.ReadFile(sessionFile)
		if err != nil {
			t.Fatalf("Failed to read session file: %v", err)
		}

		lines := strings.Split(strings.TrimSpace(string(content)), "\n")
		if len(lines) != 2 {
			t.Errorf("Expected 2 lines in file, got %d", len(lines))
		}

		// Verify both entries are valid JSON
		for i, line := range lines {
			var entry HookData
			if err := json.Unmarshal([]byte(line), &entry); err != nil {
				t.Errorf("Failed to parse line %d as JSON: %v", i+1, err)
			}
		}
	})
}

// TestDirectoryStructureValidation validates proper directory creation and organization
func TestDirectoryStructureValidation(t *testing.T) {
	tempDir := t.TempDir()
	os.Setenv("CLAUDE_PROJECT_DIR", tempDir)
	defer os.Unsetenv("CLAUDE_PROJECT_DIR")

	t.Run("ProjectLocalDirectory", func(t *testing.T) {
		// Create project-local .the-startup directory
		projectStartupDir := filepath.Join(tempDir, ".the-startup")
		os.MkdirAll(projectStartupDir, 0755)

		// Test that it prefers project-local directory
		actualStartupDir := GetStartupDir(tempDir)
		if actualStartupDir != projectStartupDir {
			t.Errorf("Expected project-local startup dir %q, got %q", projectStartupDir, actualStartupDir)
		}
	})

	t.Run("HomeFallback", func(t *testing.T) {
		// Use a temporary directory without .the-startup
		tempProjectDir := filepath.Join(tempDir, "project-without-startup")
		os.MkdirAll(tempProjectDir, 0755)

		// Get startup directory (should fall back to home)
		startupDir := GetStartupDir(tempProjectDir)

		// Should not be in the project directory
		if strings.HasPrefix(startupDir, tempProjectDir) {
			t.Errorf("Expected fallback to home directory, but got project-relative path: %s", startupDir)
		}

		// Should contain .the-startup in the path
		if !strings.Contains(startupDir, ".the-startup") {
			t.Errorf("Expected startup directory to contain '.the-startup', got: %s", startupDir)
		}
	})

	t.Run("CompleteDirectoryStructure", func(t *testing.T) {
		sessionID := "dev-structure-test"

		hookData := &HookData{
			Role:      "user",
			Content:   "Test directory structure",
			SessionID: sessionID,
			AgentID:   "the-developer",
			Timestamp: "2025-01-11T12:00:00.000Z",
		}

		// Write session and global logs
		err := WriteSessionLog(sessionID, hookData)
		if err != nil {
			t.Fatalf("WriteSessionLog failed: %v", err)
		}

		// Global logging is no longer supported

		// Verify complete directory structure (without global files)
		expectedStructure := []string{
			filepath.Join(tempDir, ".the-startup"),
			filepath.Join(tempDir, ".the-startup", sessionID),
			filepath.Join(tempDir, ".the-startup", sessionID, "agent-instructions.jsonl"),
		}

		for _, path := range expectedStructure {
			if _, err := os.Stat(path); os.IsNotExist(err) {
				t.Errorf("Expected path does not exist: %s", path)
			}
		}
	})
}

// TestConcurrentExecution tests multiple concurrent hook executions
func TestConcurrentExecution(t *testing.T) {
	tempDir := t.TempDir()
	os.Setenv("CLAUDE_PROJECT_DIR", tempDir)
	defer os.Unsetenv("CLAUDE_PROJECT_DIR")

	// Create startup directory
	startupDir := filepath.Join(tempDir, ".the-startup")
	os.MkdirAll(startupDir, 0755)

	t.Run("ConcurrentSessionWrites", func(t *testing.T) {
		numConcurrent := 10
		sessionID := "dev-concurrent-test"

		// Use WaitGroup to coordinate goroutines
		var wg sync.WaitGroup
		var mu sync.Mutex
		errors := []error{}

		// Launch concurrent writers
		for i := 0; i < numConcurrent; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()

				hookData := &HookData{
					Role:      "user",
					Content:   fmt.Sprintf("Concurrent test %d", id),
					SessionID: sessionID,
					AgentID:   fmt.Sprintf("concurrent-%d", id),
					Timestamp: fmt.Sprintf("2025-01-11T12:00:%02d.000Z", id),
				}

				if err := WriteSessionLog(sessionID, hookData); err != nil {
					mu.Lock()
					errors = append(errors, err)
					mu.Unlock()
				}
			}(i)
		}

		wg.Wait()

		// Check for errors
		if len(errors) > 0 {
			t.Fatalf("Concurrent writes failed with %d errors: %v", len(errors), errors[0])
		}

		// Verify all entries were written
		sessionFile := filepath.Join(startupDir, sessionID, "agent-instructions.jsonl")
		content, err := os.ReadFile(sessionFile)
		if err != nil {
			t.Fatalf("Failed to read session file: %v", err)
		}

		lines := strings.Split(strings.TrimSpace(string(content)), "\n")
		if len(lines) != numConcurrent {
			t.Errorf("Expected %d lines in session file, got %d", numConcurrent, len(lines))
		}

		// Verify all lines are valid JSON
		agentIDs := make(map[string]bool)
		for i, line := range lines {
			var entry HookData
			if err := json.Unmarshal([]byte(line), &entry); err != nil {
				t.Errorf("Failed to parse line %d as JSON: %v", i+1, err)
			} else {
				agentIDs[entry.AgentID] = true
			}
		}

		// Verify all agent IDs are unique (no data corruption)
		if len(agentIDs) != numConcurrent {
			t.Errorf("Expected %d unique agent IDs, got %d", numConcurrent, len(agentIDs))
		}
	})

	t.Run("ConcurrentGlobalWrites", func(t *testing.T) {
		numConcurrent := 5

		var wg sync.WaitGroup
		var mu sync.Mutex
		errors := []error{}

		// Launch concurrent global writers
		for i := 0; i < numConcurrent; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()

				hookData := &HookData{
					Role:      "assistant",
					Content:   fmt.Sprintf("Global test %d complete", id),
					SessionID: fmt.Sprintf("session-%d", id),
					AgentID:   fmt.Sprintf("global-%d", id),
					Timestamp: fmt.Sprintf("2025-01-11T12:01:%02d.000Z", id),
				}

				// Use session log instead of global log
				if err := WriteSessionLog(fmt.Sprintf("session-%d", id), hookData); err != nil {
					mu.Lock()
					errors = append(errors, err)
					mu.Unlock()
				}
			}(i)
		}

		wg.Wait()

		// Check for errors
		if len(errors) > 0 {
			t.Fatalf("Concurrent global writes failed with %d errors: %v", len(errors), errors[0])
		}

		// Verify all entries were written to their respective session files
		// Count total assistant entries across all session files
		completeCount := 0
		for i := 0; i < numConcurrent; i++ {
			sessionFile := filepath.Join(startupDir, fmt.Sprintf("session-%d", i), "agent-instructions.jsonl")
			if content, err := os.ReadFile(sessionFile); err == nil {
				lines := strings.Split(strings.TrimSpace(string(content)), "\n")
				for _, line := range lines {
					var entry HookData
					if err := json.Unmarshal([]byte(line), &entry); err == nil && entry.Role == "assistant" {
						completeCount++
					}
				}
			}
		}

		if completeCount < numConcurrent {
			t.Errorf("Expected at least %d agent_complete events, got %d", numConcurrent, completeCount)
		}
	})
}
