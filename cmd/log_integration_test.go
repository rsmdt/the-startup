package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rsmdt/the-startup/internal/log"
)

func TestLogCommandIntegration(t *testing.T) {
	// Create a temporary directory for integration test
	tempDir := t.TempDir()
	os.Setenv("CLAUDE_PROJECT_DIR", tempDir)
	defer os.Unsetenv("CLAUDE_PROJECT_DIR")

	// Create the .the-startup directory to ensure logs go to project directory
	startupDir := filepath.Join(tempDir, ".the-startup")
	if err := os.MkdirAll(startupDir, 0755); err != nil {
		t.Fatalf("Failed to create startup directory: %v", err)
	}

	// Test data
	testJSON := `{
		"tool_name": "Task",
		"tool_input": {
			"subagent_type": "the-developer",
			"description": "Integration test task",
			"prompt": "SessionId: dev-integration-test AgentId: agent-456\nThis is an integration test instruction"
		},
		"session_id": "dev-integration-test",
		"hook_event_name": "PreToolUse"
	}`

	// Create and execute the log command
	cmd := NewLogCommand()
	cmd.SetIn(strings.NewReader(testJSON))

	// Set the --assistant flag
	cmd.SetArgs([]string{"--assistant"})

	// Execute command
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Command execution failed: %v", err)
	}

	// Verify global log file was created
	globalLogPath := filepath.Join(startupDir, "all-agent-instructions.jsonl")
	
	if _, err := os.Stat(globalLogPath); os.IsNotExist(err) {
		t.Errorf("Global log file was not created at %s", globalLogPath)
	}

	// Verify session log file was created
	sessionLogPath := filepath.Join(startupDir, "dev-integration-test", "agent-instructions.jsonl")
	
	if _, err := os.Stat(sessionLogPath); os.IsNotExist(err) {
		t.Errorf("Session log file was not created at %s", sessionLogPath)
	}

	// Read and verify global log content
	globalContent, err := os.ReadFile(globalLogPath)
	if err != nil {
		t.Fatalf("Failed to read global log file: %v", err)
	}

	var logEntry log.HookData
	if err := json.Unmarshal([]byte(strings.TrimSpace(string(globalContent))), &logEntry); err != nil {
		t.Fatalf("Failed to parse log entry JSON: %v", err)
	}

	// Verify log entry content
	if logEntry.Event != "agent_start" {
		t.Errorf("Expected event 'agent_start', got '%s'", logEntry.Event)
	}

	if logEntry.AgentType != "the-developer" {
		t.Errorf("Expected agent_type 'the-developer', got '%s'", logEntry.AgentType)
	}

	if logEntry.AgentID != "agent-456" {
		t.Errorf("Expected agent_id 'agent-456', got '%s'", logEntry.AgentID)
	}

	if logEntry.SessionID != "dev-integration-test" {
		t.Errorf("Expected session_id 'dev-integration-test', got '%s'", logEntry.SessionID)
	}

	if logEntry.Description != "Integration test task" {
		t.Errorf("Expected description 'Integration test task', got '%s'", logEntry.Description)
	}

	if logEntry.Instruction == "" {
		t.Error("Expected instruction to be populated for agent_start event")
	}

	if logEntry.OutputSummary != "" {
		t.Errorf("Expected empty output_summary for agent_start event, got '%s'", logEntry.OutputSummary)
	}
}

func TestLogCommandUserIntegration(t *testing.T) {
	// Create a temporary directory for integration test
	tempDir := t.TempDir()
	os.Setenv("CLAUDE_PROJECT_DIR", tempDir)
	defer os.Unsetenv("CLAUDE_PROJECT_DIR")

	// Create the .the-startup directory to ensure logs go to project directory
	startupDir := filepath.Join(tempDir, ".the-startup")
	if err := os.MkdirAll(startupDir, 0755); err != nil {
		t.Fatalf("Failed to create startup directory: %v", err)
	}

	// Test data for PostToolUse (user flag)
	testJSON := `{
		"tool_name": "Task",
		"tool_input": {
			"subagent_type": "the-tester",
			"description": "User integration test task",
			"prompt": "SessionId: dev-user-test AgentId: agent-789\nThis is a user integration test instruction"
		},
		"session_id": "dev-user-test",
		"hook_event_name": "PostToolUse",
		"output": "Task completed successfully with comprehensive test results and validation."
	}`

	// Create and execute the log command
	cmd := NewLogCommand()
	cmd.SetIn(strings.NewReader(testJSON))

	// Set the --user flag
	cmd.SetArgs([]string{"--user"})

	// Execute command
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Command execution failed: %v", err)
	}

	// Read and verify global log content
	globalLogPath := filepath.Join(tempDir, ".the-startup", "all-agent-instructions.jsonl")
	globalContent, err := os.ReadFile(globalLogPath)
	if err != nil {
		t.Fatalf("Failed to read global log file: %v", err)
	}

	var logEntry log.HookData
	if err := json.Unmarshal([]byte(strings.TrimSpace(string(globalContent))), &logEntry); err != nil {
		t.Fatalf("Failed to parse log entry JSON: %v", err)
	}

	// Verify log entry content for agent_complete event
	if logEntry.Event != "agent_complete" {
		t.Errorf("Expected event 'agent_complete', got '%s'", logEntry.Event)
	}

	if logEntry.AgentType != "the-tester" {
		t.Errorf("Expected agent_type 'the-tester', got '%s'", logEntry.AgentType)
	}

	if logEntry.Instruction != "" {
		t.Errorf("Expected empty instruction for agent_complete event, got '%s'", logEntry.Instruction)
	}

	if logEntry.OutputSummary == "" {
		t.Error("Expected output_summary to be populated for agent_complete event")
	}

	expectedOutput := "Task completed successfully with comprehensive test results and validation."
	if logEntry.OutputSummary != expectedOutput {
		t.Errorf("Expected output_summary '%s', got '%s'", expectedOutput, logEntry.OutputSummary)
	}
}

// TestLogCommandFlagValidation tests CLI flag parsing and validation edge cases
func TestLogCommandFlagValidation(t *testing.T) {
	t.Run("BothFlagsSpecified", func(t *testing.T) {
		cmd := NewLogCommand()
		cmd.SetArgs([]string{"--assistant", "--user"})
		
		err := cmd.Execute()
		if err == nil {
			t.Error("Expected error when both --assistant and --user flags are specified")
		}
	})

	t.Run("NoFlagsSpecified", func(t *testing.T) {
		cmd := NewLogCommand()
		cmd.SetArgs([]string{})
		
		err := cmd.Execute()
		if err == nil {
			t.Error("Expected error when no flags are specified")
		}
	})

	t.Run("ShortFlagAssistant", func(t *testing.T) {
		tempDir := t.TempDir()
		os.Setenv("CLAUDE_PROJECT_DIR", tempDir)
		defer os.Unsetenv("CLAUDE_PROJECT_DIR")

		os.MkdirAll(filepath.Join(tempDir, ".the-startup"), 0755)

		testJSON := `{
			"tool_name": "Task",
			"tool_input": {
				"subagent_type": "the-developer",
				"description": "Short flag test",
				"prompt": "SessionId: dev-short-test AgentId: agent-short\nShort flag test instruction"
			},
			"session_id": "dev-short-test",
			"hook_event_name": "PreToolUse"
		}`

		cmd := NewLogCommand()
		cmd.SetIn(strings.NewReader(testJSON))
		cmd.SetArgs([]string{"-a"})

		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Command execution failed with short flag: %v", err)
		}

		// Verify log was created
		globalLogPath := filepath.Join(tempDir, ".the-startup", "all-agent-instructions.jsonl")
		if _, err := os.Stat(globalLogPath); os.IsNotExist(err) {
			t.Error("Global log file was not created with short flag")
		}
	})

	t.Run("ShortFlagUser", func(t *testing.T) {
		tempDir := t.TempDir()
		os.Setenv("CLAUDE_PROJECT_DIR", tempDir)
		defer os.Unsetenv("CLAUDE_PROJECT_DIR")

		os.MkdirAll(filepath.Join(tempDir, ".the-startup"), 0755)

		testJSON := `{
			"tool_name": "Task",
			"tool_input": {
				"subagent_type": "the-developer",
				"description": "Short flag user test",
				"prompt": "SessionId: dev-user-short-test AgentId: agent-user-short\nShort flag user test instruction"
			},
			"session_id": "dev-user-short-test",
			"hook_event_name": "PostToolUse",
			"output": "Short flag user test completed"
		}`

		cmd := NewLogCommand()
		cmd.SetIn(strings.NewReader(testJSON))
		cmd.SetArgs([]string{"-u"})

		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Command execution failed with short user flag: %v", err)
		}

		// Verify log was created
		globalLogPath := filepath.Join(tempDir, ".the-startup", "all-agent-instructions.jsonl")
		if _, err := os.Stat(globalLogPath); os.IsNotExist(err) {
			t.Error("Global log file was not created with short user flag")
		}
	})
}

// TestLogCommandStdinProcessing tests various stdin processing scenarios
func TestLogCommandStdinProcessing(t *testing.T) {
	t.Run("EmptyStdin", func(t *testing.T) {
		tempDir := t.TempDir()
		os.Setenv("CLAUDE_PROJECT_DIR", tempDir)
		defer os.Unsetenv("CLAUDE_PROJECT_DIR")

		cmd := NewLogCommand()
		cmd.SetIn(strings.NewReader(""))
		cmd.SetArgs([]string{"--assistant"})

		// Should exit silently (matching Python behavior)
		err := cmd.Execute()
		if err != nil {
			t.Errorf("Expected silent exit on empty stdin, got error: %v", err)
		}
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		tempDir := t.TempDir()
		os.Setenv("CLAUDE_PROJECT_DIR", tempDir)
		defer os.Unsetenv("CLAUDE_PROJECT_DIR")

		cmd := NewLogCommand()
		cmd.SetIn(strings.NewReader("{invalid json"))
		cmd.SetArgs([]string{"--assistant"})

		// Should exit silently (matching Python behavior)
		err := cmd.Execute()
		if err != nil {
			t.Errorf("Expected silent exit on invalid JSON, got error: %v", err)
		}
	})

	t.Run("FilteredOutTool", func(t *testing.T) {
		tempDir := t.TempDir()
		os.Setenv("CLAUDE_PROJECT_DIR", tempDir)
		defer os.Unsetenv("CLAUDE_PROJECT_DIR")

		// Non-Task tool should be filtered out
		testJSON := `{
			"tool_name": "NotTask",
			"tool_input": {
				"subagent_type": "the-developer"
			},
			"session_id": "dev-filtered-test",
			"hook_event_name": "PreToolUse"
		}`

		cmd := NewLogCommand()
		cmd.SetIn(strings.NewReader(testJSON))
		cmd.SetArgs([]string{"--assistant"})

		err := cmd.Execute()
		if err != nil {
			t.Errorf("Expected silent exit for filtered tool, got error: %v", err)
		}

		// Verify no files were created
		globalLogPath := filepath.Join(tempDir, ".the-startup", "all-agent-instructions.jsonl")
		if _, err := os.Stat(globalLogPath); !os.IsNotExist(err) {
			t.Error("Global log file should not be created for filtered tools")
		}
	})

	t.Run("FilteredOutSubagent", func(t *testing.T) {
		tempDir := t.TempDir()
		os.Setenv("CLAUDE_PROJECT_DIR", tempDir)
		defer os.Unsetenv("CLAUDE_PROJECT_DIR")

		// Non-"the-" subagent should be filtered out
		testJSON := `{
			"tool_name": "Task",
			"tool_input": {
				"subagent_type": "not-the-developer"
			},
			"session_id": "dev-filtered-subagent-test",
			"hook_event_name": "PreToolUse"
		}`

		cmd := NewLogCommand()
		cmd.SetIn(strings.NewReader(testJSON))
		cmd.SetArgs([]string{"--assistant"})

		err := cmd.Execute()
		if err != nil {
			t.Errorf("Expected silent exit for filtered subagent, got error: %v", err)
		}

		// Verify no files were created
		globalLogPath := filepath.Join(tempDir, ".the-startup", "all-agent-instructions.jsonl")
		if _, err := os.Stat(globalLogPath); !os.IsNotExist(err) {
			t.Error("Global log file should not be created for filtered subagents")
		}
	})

	t.Run("LargeJSON", func(t *testing.T) {
		tempDir := t.TempDir()
		os.Setenv("CLAUDE_PROJECT_DIR", tempDir)
		defer os.Unsetenv("CLAUDE_PROJECT_DIR")

		os.MkdirAll(filepath.Join(tempDir, ".the-startup"), 0755)

		// Create large JSON with long output
		largeOutput := strings.Repeat("This is a very long output that should be truncated. ", 100)
		testJSON := fmt.Sprintf(`{
			"tool_name": "Task",
			"tool_input": {
				"subagent_type": "the-developer",
				"description": "Large JSON test",
				"prompt": "SessionId: dev-large-test AgentId: agent-large\nLarge JSON test instruction"
			},
			"session_id": "dev-large-test",
			"hook_event_name": "PostToolUse",
			"output": "%s"
		}`, largeOutput)

		cmd := NewLogCommand()
		cmd.SetIn(strings.NewReader(testJSON))
		cmd.SetArgs([]string{"--user"})

		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Command execution failed with large JSON: %v", err)
		}

		// Read and verify truncation
		globalLogPath := filepath.Join(tempDir, ".the-startup", "all-agent-instructions.jsonl")
		globalContent, err := os.ReadFile(globalLogPath)
		if err != nil {
			t.Fatalf("Failed to read global log file: %v", err)
		}

		var logEntry log.HookData
		if err := json.Unmarshal([]byte(strings.TrimSpace(string(globalContent))), &logEntry); err != nil {
			t.Fatalf("Failed to parse log entry JSON: %v", err)
		}

		// Output should be truncated
		if len(logEntry.OutputSummary) <= 1000 {
			t.Errorf("Expected output to be truncated, but got length %d", len(logEntry.OutputSummary))
		}
		if !strings.Contains(logEntry.OutputSummary, "... [truncated") {
			t.Error("Expected truncation message in output")
		}
	})
}

// TestLogCommandErrorHandling tests file system error handling and fallback scenarios
func TestLogCommandErrorHandling(t *testing.T) {
	t.Run("ReadOnlyDirectory", func(t *testing.T) {
		// This test is platform-specific and may not work on all systems
		if testing.Short() {
			t.Skip("Skipping read-only directory test in short mode")
		}

		tempDir := t.TempDir()
		os.Setenv("CLAUDE_PROJECT_DIR", tempDir)
		defer os.Unsetenv("CLAUDE_PROJECT_DIR")

		// Create startup directory and make it read-only
		startupDir := filepath.Join(tempDir, ".the-startup")
		os.MkdirAll(startupDir, 0755)
		
		// Try to make directory read-only (may not work on all systems)
		os.Chmod(startupDir, 0444)
		defer os.Chmod(startupDir, 0755) // Restore permissions for cleanup

		testJSON := `{
			"tool_name": "Task",
			"tool_input": {
				"subagent_type": "the-developer",
				"description": "Read-only test",
				"prompt": "SessionId: dev-readonly-test AgentId: agent-readonly\nRead-only test instruction"
			},
			"session_id": "dev-readonly-test",
			"hook_event_name": "PreToolUse"
		}`

		cmd := NewLogCommand()
		cmd.SetIn(strings.NewReader(testJSON))
		cmd.SetArgs([]string{"--assistant"})

		// Should exit silently even on permission errors
		err := cmd.Execute()
		if err != nil {
			t.Errorf("Expected silent exit on permission error, got: %v", err)
		}
	})

	t.Run("MissingSessionID", func(t *testing.T) {
		tempDir := t.TempDir()
		os.Setenv("CLAUDE_PROJECT_DIR", tempDir)
		defer os.Unsetenv("CLAUDE_PROJECT_DIR")

		os.MkdirAll(filepath.Join(tempDir, ".the-startup"), 0755)

		// No session ID in prompt
		testJSON := `{
			"tool_name": "Task",
			"tool_input": {
				"subagent_type": "the-developer",
				"description": "No session ID test",
				"prompt": "No session ID in this prompt\nTest instruction"
			},
			"session_id": "dev-no-session-test",
			"hook_event_name": "PreToolUse"
		}`

		cmd := NewLogCommand()
		cmd.SetIn(strings.NewReader(testJSON))
		cmd.SetArgs([]string{"--assistant"})

		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Command execution failed: %v", err)
		}

		// Should still create global log
		globalLogPath := filepath.Join(tempDir, ".the-startup", "all-agent-instructions.jsonl")
		if _, err := os.Stat(globalLogPath); os.IsNotExist(err) {
			t.Error("Global log file should be created even without session ID")
		}
	})
}