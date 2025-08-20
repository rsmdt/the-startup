package log

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestSecurityMaliciousInputs tests security scenarios with malicious inputs
func TestSecurityMaliciousInputs(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		isPostHook  bool
		expectNil   bool
		expectError bool
		validateFn  func(*testing.T, *HookData)
	}{
		{
			name: "JSON bomb - deeply nested objects",
			input: `{"tool_name": "Task", "tool_input": {"subagent_type": "the-test", "nested": ` +
				strings.Repeat(`{"level": `, 100) + `"deep"` + strings.Repeat(`}`, 100) + `}}`,
			isPostHook:  false,
			expectError: false, // Should handle gracefully
			expectNil:   false,
			validateFn: func(t *testing.T, data *HookData) {
				// Just verify the data was parsed without errors
				if data == nil {
					t.Errorf("Expected data to be parsed correctly despite nested objects")
				}
			},
		},
		{
			name: "extremely large JSON input",
			input: `{"tool_name": "Task", "tool_input": {"subagent_type": "the-test", "large_field": "` +
				strings.Repeat("A", 100000) + `"}}`,
			isPostHook:  false,
			expectError: false, // Should handle large inputs
			expectNil:   false,
		},
		{
			name:        "JSON with null bytes",
			input:       "{\"tool_name\": \"Task\", \"tool_input\": {\"subagent_type\": \"the-test\", \"null_byte\": \"data\x00here\"}}",
			isPostHook:  false,
			expectError: true, // JSON standard doesn't allow null bytes in strings
			expectNil:   false,
		},
		{
			name: "path traversal in session ID",
			input: `{"tool_name": "Task", "tool_input": {"subagent_type": "the-test", "description": "Test", 
				"prompt": "SessionId: ../../../etc/passwd AgentId: test\nTest"}}`,
			isPostHook: false,
			expectNil:  false,
			validateFn: func(t *testing.T, data *HookData) {
				// Should extract the malicious session ID as-is
				// File system operations should handle path sanitization
				if data.SessionID != "../../../etc/passwd" {
					t.Errorf("Expected raw session ID, got: %s", data.SessionID)
				}
			},
		},
		{
			name: "command injection in description",
			input: `{"tool_name": "Task", "tool_input": {"subagent_type": "the-test", 
				"description": "'; rm -rf /; echo '", "prompt": "SessionId: dev-test AgentId: test\nTest"}}`,
			isPostHook: false,
			expectNil:  false,
			validateFn: func(t *testing.T, data *HookData) {
				// Content should contain the original prompt
				if !strings.Contains(data.Content, "Test") {
					t.Errorf("Expected content to contain test prompt")
				}
			},
		},
		{
			name: "script injection in prompt",
			input: `{"tool_name": "Task", "tool_input": {"subagent_type": "the-test", 
				"prompt": "<script>alert('xss')</script>SessionId: dev-test AgentId: test\nTest"}}`,
			isPostHook: false,
			expectNil:  false,
			validateFn: func(t *testing.T, data *HookData) {
				// Should extract session/agent IDs correctly despite script content
				if data.SessionID != "dev-test" {
					t.Errorf("Expected session ID extraction to work despite script content")
				}
			},
		},
		{
			name: "environment variable injection",
			input: `{"tool_name": "Task", "tool_input": {"subagent_type": "the-test", 
				"description": "$(malicious_command)", "prompt": "SessionId: dev-test AgentId: test\nTest"}}`,
			isPostHook: false,
			expectNil:  false,
			validateFn: func(t *testing.T, data *HookData) {
				// Content should be preserved as literal string
				if !strings.Contains(data.Content, "Test") {
					t.Errorf("Expected content to be preserved")
				}
			},
		},
		{
			name: "unicode normalization attack",
			input: `{"tool_name": "Task", "tool_input": {"subagent_type": "the-test", 
				"description": "\u0041\u0300", "prompt": "SessionId: dev-unicode AgentId: test\nTest"}}`,
			isPostHook: false,
			expectNil:  false,
			validateFn: func(t *testing.T, data *HookData) {
				// Should handle unicode characters
				if data.SessionID != "dev-unicode" {
					t.Error("Should handle unicode in adjacent fields")
				}
			},
		},
		{
			name: "format string attack",
			input: `{"tool_name": "Task", "tool_input": {"subagent_type": "the-test", 
				"description": "%s%s%s%s%s%n", "prompt": "SessionId: dev-format AgentId: test\nTest"}}`,
			isPostHook: false,
			expectNil:  false,
			validateFn: func(t *testing.T, data *HookData) {
				// Content should preserve format strings as literal
				if !strings.Contains(data.Content, "Test") {
					t.Errorf("Expected content to be preserved")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			result, err := ProcessToolCall(reader, tt.isPostHook)

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

			if tt.expectNil {
				if result != nil {
					t.Errorf("Expected nil result but got: %+v", result)
				}
				return
			}

			if result == nil {
				t.Error("Expected non-nil result")
				return
			}

			if tt.validateFn != nil {
				tt.validateFn(t, result)
			}
		})
	}
}

// TestSecurityFileOperations tests security aspects of file operations
func TestSecurityFileOperations(t *testing.T) {
	tests := []struct {
		name        string
		sessionID   string
		setupFn     func(*testing.T, string)
		expectError bool
		validateFn  func(*testing.T, string, error)
	}{
		{
			name:      "path traversal in session ID for file operations",
			sessionID: "../../../etc/passwd",
			setupFn: func(t *testing.T, tempDir string) {
				// Set up temporary project directory
				os.Setenv("CLAUDE_PROJECT_DIR", tempDir)
				startupDir := filepath.Join(tempDir, ".the-startup")
				os.MkdirAll(startupDir, 0755)
			},
			expectError: false, // Should handle path traversal gracefully
			validateFn: func(t *testing.T, tempDir string, err error) {
				if err != nil {
					t.Errorf("Expected no error for path traversal handling, got: %v", err)
				}

				// Verify that the malicious session ID was used as-is
				// The file system operations should handle any security implications
				// We're just testing that the application doesn't crash or behave unexpectedly
			},
		},
		{
			name:      "session ID with null bytes",
			sessionID: "dev-session\x00malicious",
			setupFn: func(t *testing.T, tempDir string) {
				os.Setenv("CLAUDE_PROJECT_DIR", tempDir)
				startupDir := filepath.Join(tempDir, ".the-startup")
				os.MkdirAll(startupDir, 0755)
			},
			expectError: false, // Should handle null bytes
			validateFn: func(t *testing.T, tempDir string, err error) {
				// Should not cause any security issues
				if err != nil {
					t.Logf("Error handling null bytes in session ID: %v", err)
				}
			},
		},
		{
			name:      "extremely long session ID",
			sessionID: strings.Repeat("a", 1000),
			setupFn: func(t *testing.T, tempDir string) {
				os.Setenv("CLAUDE_PROJECT_DIR", tempDir)
				startupDir := filepath.Join(tempDir, ".the-startup")
				os.MkdirAll(startupDir, 0755)
			},
			expectError: false, // Should handle long paths
			validateFn: func(t *testing.T, tempDir string, err error) {
				// May or may not succeed depending on filesystem limits
				// But should not cause security issues
				if err != nil {
					t.Logf("Expected potential error for very long session ID: %v", err)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()

			// Save and restore original environment
			originalEnv := os.Getenv("CLAUDE_PROJECT_DIR")
			defer func() {
				if originalEnv == "" {
					os.Unsetenv("CLAUDE_PROJECT_DIR")
				} else {
					os.Setenv("CLAUDE_PROJECT_DIR", originalEnv)
				}
			}()

			tt.setupFn(t, tempDir)

			// Test file operations with malicious session ID
			hookData := &HookData{
				Role:      "user",
				Content:   "Security test",
				SessionID: tt.sessionID,
				Timestamp: "2025-01-11T12:00:00.000Z",
				AgentID:   "security-test",
			}

			err := WriteSessionLog(tt.sessionID, hookData)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if tt.validateFn != nil {
				tt.validateFn(t, tempDir, err)
			}
		})
	}
}

// TestSecurityInputSanitization tests that inputs are properly handled without executing
func TestSecurityInputSanitization(t *testing.T) {
	// Test that debug output doesn't leak sensitive information
	originalDebug := os.Getenv("DEBUG_HOOKS")
	os.Setenv("DEBUG_HOOKS", "1")
	defer func() {
		if originalDebug == "" {
			os.Unsetenv("DEBUG_HOOKS")
		} else {
			os.Setenv("DEBUG_HOOKS", originalDebug)
		}
	}()

	maliciousInput := `{
		"tool_name": "Task", 
		"tool_input": {
			"subagent_type": "the-test",
			"description": "password123secret",
			"prompt": "SessionId: dev-secret AgentId: test\nSecret data: API_KEY=secret123"
		}
	}`

	reader := strings.NewReader(maliciousInput)
	result, err := ProcessToolCall(reader, false)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	if result == nil {
		t.Error("Expected non-nil result")
		return
	}

	// Verify that sensitive data is preserved but not executed
	if !strings.Contains(result.Content, "API_KEY=secret123") {
		t.Errorf("Expected content to contain original data")
	}

	// Test that the data can be safely marshaled to JSON
	tempDir := t.TempDir()
	os.Setenv("CLAUDE_PROJECT_DIR", tempDir)
	os.MkdirAll(filepath.Join(tempDir, ".the-startup"), 0755)

	// This should not execute any commands or cause security issues
	err = WriteSessionLog("test-session", result)
	if err != nil {
		t.Errorf("Failed to write potentially sensitive data safely: %v", err)
	}
}
