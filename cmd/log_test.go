package cmd

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestLogCommand(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		stdin         string
		expectSuccess bool
		expectedErr   string
	}{
		{
			name:          "assistant flag with valid JSON",
			args:          []string{"--assistant"},
			stdin:         `{"tool_name":"Task","tool_input":{"subagent_type":"the-developer","description":"test task","prompt":"SessionId: dev-123 AgentId: agent-456\nTest instruction"},"session_id":"dev-123","hook_event_name":"PreToolUse"}`,
			expectSuccess: true,
		},
		{
			name:          "user flag with valid JSON",
			args:          []string{"--user"},
			stdin:         `{"tool_name":"Task","tool_input":{"subagent_type":"the-developer","description":"test task","prompt":"SessionId: dev-123 AgentId: agent-456\nTest instruction"},"session_id":"dev-123","hook_event_name":"PostToolUse","output":"Task completed successfully"}`,
			expectSuccess: true,
		},
		{
			name:          "assistant short flag",
			args:          []string{"-a"},
			stdin:         `{"tool_name":"Task","tool_input":{"subagent_type":"the-developer","description":"test task","prompt":"SessionId: dev-123 AgentId: agent-456\nTest instruction"},"session_id":"dev-123","hook_event_name":"PreToolUse"}`,
			expectSuccess: true,
		},
		{
			name:          "user short flag",
			args:          []string{"-u"},
			stdin:         `{"tool_name":"Task","tool_input":{"subagent_type":"the-developer","description":"test task","prompt":"SessionId: dev-123 AgentId: agent-456\nTest instruction"},"session_id":"dev-123","hook_event_name":"PostToolUse","output":"Task completed successfully"}`,
			expectSuccess: true,
		},
		{
			name:          "invalid JSON input - should exit silently",
			args:          []string{"--assistant"},
			stdin:         `invalid json`,
			expectSuccess: true, // Should exit silently (code 0)
		},
		{
			name:          "non-Task tool - should exit silently",
			args:          []string{"--assistant"},
			stdin:         `{"tool_name":"Other","tool_input":{"subagent_type":"the-developer","description":"test task","prompt":"Test instruction"},"session_id":"dev-123","hook_event_name":"PreToolUse"}`,
			expectSuccess: true, // Should exit silently (code 0)
		},
		{
			name:          "no flags provided",
			args:          []string{},
			expectSuccess: false,
			expectedErr:   "exactly one mode required: --assistant, --user, or --read",
		},
		{
			name:          "both flags provided",
			args:          []string{"--assistant", "--user"},
			expectSuccess: false,
			expectedErr:   "exactly one mode required: --assistant, --user, or --read",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary directory for test logs
			tempDir := t.TempDir()
			os.Setenv("CLAUDE_PROJECT_DIR", tempDir)
			defer os.Unsetenv("CLAUDE_PROJECT_DIR")

			// Create the log command
			cmd := NewLogCommand()

			// Set up stdin
			if tt.stdin != "" {
				cmd.SetIn(strings.NewReader(tt.stdin))
			}

			// Capture output
			var outBuf, errBuf bytes.Buffer
			cmd.SetOut(&outBuf)
			cmd.SetErr(&errBuf)

			// Set args
			cmd.SetArgs(tt.args)

			// Execute command
			err := cmd.Execute()

			if tt.expectSuccess {
				if err != nil {
					t.Errorf("Expected success but got error: %v", err)
				}
			} else {
				if err == nil {
					t.Errorf("Expected error but command succeeded")
				} else if tt.expectedErr != "" && !strings.Contains(err.Error(), tt.expectedErr) {
					t.Errorf("Expected error containing '%s', got: %v", tt.expectedErr, err)
				}
			}

			// Verify stderr is empty (no debug output unless DEBUG_HOOKS is set)
			if errBuf.Len() > 0 && os.Getenv("DEBUG_HOOKS") == "" {
				t.Errorf("Unexpected stderr output: %s", errBuf.String())
			}
		})
	}
}

func TestLogCommandFlags(t *testing.T) {
	cmd := NewLogCommand()

	// Check that existing flags are properly defined
	assistantFlag := cmd.Flags().Lookup("assistant")
	if assistantFlag == nil {
		t.Error("--assistant flag not found")
	}

	userFlag := cmd.Flags().Lookup("user")
	if userFlag == nil {
		t.Error("--user flag not found")
	}

	// Check new read mode flags
	readFlag := cmd.Flags().Lookup("read")
	if readFlag == nil {
		t.Error("--read flag not found")
	}

	agentIDFlag := cmd.Flags().Lookup("agent-id")
	if agentIDFlag == nil {
		t.Error("--agent-id flag not found")
	}

	linesFlag := cmd.Flags().Lookup("lines")
	if linesFlag == nil {
		t.Error("--lines flag not found")
	}

	sessionFlag := cmd.Flags().Lookup("session")
	if sessionFlag == nil {
		t.Error("--session flag not found")
	}

	formatFlag := cmd.Flags().Lookup("format")
	if formatFlag == nil {
		t.Error("--format flag not found")
	}

	metadataFlag := cmd.Flags().Lookup("include-metadata")
	if metadataFlag == nil {
		t.Error("--include-metadata flag not found")
	}

	// Check short versions
	if assistantFlag.Shorthand != "a" {
		t.Errorf("Expected --assistant short flag to be 'a', got '%s'", assistantFlag.Shorthand)
	}

	if userFlag.Shorthand != "u" {
		t.Errorf("Expected --user short flag to be 'u', got '%s'", userFlag.Shorthand)
	}

	if readFlag.Shorthand != "r" {
		t.Errorf("Expected --read short flag to be 'r', got '%s'", readFlag.Shorthand)
	}
}

func TestLogCommand_ReadModeValidation(t *testing.T) {
	testCases := []struct {
		name        string
		args        []string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Read mode without agent-id",
			args:        []string{"--read"},
			expectError: true,
			errorMsg:    "--agent-id required",
		},
		{
			name:        "Read mode with invalid agent-id format",
			args:        []string{"--read", "--agent-id", "invalid@id"},
			expectError: true,
			errorMsg:    "invalid agent-id format",
		},
		{
			name:        "Read mode with reserved word agent-id",
			args:        []string{"--read", "--agent-id", "main"},
			expectError: true,
			errorMsg:    "invalid agent-id format",
		},
		{
			name:        "Read mode with lines too small",
			args:        []string{"--read", "--agent-id", "arch-001", "--lines", "0"},
			expectError: true,
			errorMsg:    "--lines must be between 1 and 1000",
		},
		{
			name:        "Read mode with lines too large",
			args:        []string{"--read", "--agent-id", "arch-001", "--lines", "1001"},
			expectError: true,
			errorMsg:    "--lines must be between 1 and 1000",
		},
		{
			name:        "Read mode with invalid format",
			args:        []string{"--read", "--agent-id", "arch-001", "--format", "xml"},
			expectError: true,
			errorMsg:    "invalid format",
		},
		{
			name:        "Valid read mode",
			args:        []string{"--read", "--agent-id", "arch-001"},
			expectError: false,
		},
		{
			name:        "Valid read mode with all options",
			args:        []string{"--read", "--agent-id", "test-agent", "--lines", "100", "--session", "dev-test", "--format", "text", "--include-metadata"},
			expectError: false,
		},
		{
			name:        "Multiple modes - assistant and read",
			args:        []string{"--assistant", "--read", "--agent-id", "arch-001"},
			expectError: true,
			errorMsg:    "exactly one mode required",
		},
		{
			name:        "Multiple modes - user and read",
			args:        []string{"--user", "--read", "--agent-id", "arch-001"},
			expectError: true,
			errorMsg:    "exactly one mode required",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tempDir := t.TempDir()
			os.Setenv("CLAUDE_PROJECT_DIR", tempDir)
			defer os.Unsetenv("CLAUDE_PROJECT_DIR")

			cmd := NewLogCommand()
			cmd.SetArgs(tc.args)

			// Capture output to prevent test pollution
			var stdout, stderr bytes.Buffer
			cmd.SetOut(&stdout)
			cmd.SetErr(&stderr)

			err := cmd.Execute()

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error but command succeeded")
				} else if !strings.Contains(err.Error(), tc.errorMsg) {
					t.Errorf("Expected error containing '%s', got: %v", tc.errorMsg, err)
				}
			} else {
				// For valid read mode, we expect it to succeed even with no data
				if err != nil {
					t.Errorf("Expected success but got error: %v", err)
				}
			}
		})
	}
}
