package log

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestCoverageBoost adds simple tests to boost coverage for remaining lines
func TestCoverageBoost(t *testing.T) {
	// Test home directory fallback in GetStartupDir
	t.Run("GetStartupDir_home_fallback", func(t *testing.T) {
		tempDir := t.TempDir()
		// Don't create local .the-startup, so it should use home directory fallback
		result := GetStartupDir(tempDir)

		// Should not be the local directory since it doesn't exist
		localStartup := filepath.Join(tempDir, ".the-startup")
		if result == localStartup {
			t.Errorf("Expected home directory fallback, got local directory: %s", result)
		}
	})

	// Test CLAUDE_PROJECT_DIR environment variable
	t.Run("GetProjectDir_with_env", func(t *testing.T) {
		original := os.Getenv("CLAUDE_PROJECT_DIR")
		defer func() {
			if original == "" {
				os.Unsetenv("CLAUDE_PROJECT_DIR")
			} else {
				os.Setenv("CLAUDE_PROJECT_DIR", original)
			}
		}()

		testDir := "/test/project/dir"
		os.Setenv("CLAUDE_PROJECT_DIR", testDir)

		result := GetProjectDir()
		if result != testDir {
			t.Errorf("Expected %s, got %s", testDir, result)
		}
	})

	// Test WriteSessionLog with empty session ID
	t.Run("WriteSessionLog_empty_session", func(t *testing.T) {
		hookData := &HookData{
			Event:     "agent_start",
			AgentType: "the-test",
			Timestamp: "2025-01-11T12:00:00.000Z",
		}

		err := WriteSessionLog("", hookData)
		if err != nil {
			t.Errorf("WriteSessionLog with empty session should not error: %v", err)
		}
	})

	// Test simple JSON marshaling/unmarshaling scenarios
	t.Run("JSON_operations", func(t *testing.T) {
		hookData := &HookData{
			Event:         "agent_complete",
			AgentType:     "the-coverage",
			AgentID:       "cov-001",
			Description:   "Coverage test",
			OutputSummary: "Test completed",
			SessionID:     "dev-coverage-123",
			Timestamp:     "2025-01-11T12:00:00.000Z",
		}

		// This exercises the JSON marshaling code path
		tempDir := t.TempDir()
		original := os.Getenv("CLAUDE_PROJECT_DIR")
		defer func() {
			if original == "" {
				os.Unsetenv("CLAUDE_PROJECT_DIR")
			} else {
				os.Setenv("CLAUDE_PROJECT_DIR", original)
			}
		}()
		os.Setenv("CLAUDE_PROJECT_DIR", tempDir)

		// Test WriteGlobalLog - this will exercise the JSON marshaling and file writing
		err := WriteGlobalLog(hookData)
		if err != nil {
			t.Errorf("WriteGlobalLog failed: %v", err)
		}
	})

	// Test FindLatestSession with startup directory but no sessions
	t.Run("FindLatestSession_no_sessions", func(t *testing.T) {
		tempDir := t.TempDir()
		startupDir := filepath.Join(tempDir, ".the-startup")
		os.MkdirAll(startupDir, 0755)
		// Create non-dev directories that should be ignored
		os.MkdirAll(filepath.Join(startupDir, "not-dev-session"), 0755)

		result := FindLatestSession(tempDir)
		if result != "" {
			t.Errorf("Expected empty string when no dev- sessions exist, got: %s", result)
		}
	})

	// Test debug functions
	t.Run("Debug_functions", func(t *testing.T) {
		original := os.Getenv("DEBUG_HOOKS")
		defer func() {
			if original == "" {
				os.Unsetenv("DEBUG_HOOKS")
			} else {
				os.Setenv("DEBUG_HOOKS", original)
			}
		}()

		// Test with debug enabled
		os.Setenv("DEBUG_HOOKS", "1")

		// These should not cause any errors, just exercise the code paths
		DebugLog("test message: %s", "value")
		DebugError(fmt.Errorf("test error"))

		// Test with debug disabled
		os.Unsetenv("DEBUG_HOOKS")
		DebugLog("should not output")
		DebugError(fmt.Errorf("should not output"))
	})

	// Test error handling paths in ProcessToolCall
	t.Run("ProcessToolCall_error_paths", func(t *testing.T) {
		// Test when GetProjectDir and FindLatestSession are called
		testInput := `{
			"tool_name": "Task",
			"tool_input": {
				"subagent_type": "the-test-no-session",
				"description": "Test without session ID",
				"prompt": "No session ID in this prompt"
			}
		}`

		// Temporarily unset CLAUDE_PROJECT_DIR to test fallback
		original := os.Getenv("CLAUDE_PROJECT_DIR")
		os.Unsetenv("CLAUDE_PROJECT_DIR")
		defer func() {
			if original != "" {
				os.Setenv("CLAUDE_PROJECT_DIR", original)
			}
		}()

		reader := strings.NewReader(testInput)
		result, err := ProcessToolCall(reader, false)
		if err != nil {
			t.Errorf("ProcessToolCall should handle missing session ID gracefully: %v", err)
		}
		if result == nil {
			t.Error("Expected non-nil result")
		}
	})
}
