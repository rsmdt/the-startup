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
			expectedErr:   "exactly one of --assistant or --user must be specified",
		},
		{
			name:          "both flags provided",
			args:          []string{"--assistant", "--user"},
			expectSuccess: false,
			expectedErr:   "exactly one of --assistant or --user must be specified",
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

	// Check that flags are properly defined
	assistantFlag := cmd.Flags().Lookup("assistant")
	if assistantFlag == nil {
		t.Error("--assistant flag not found")
	}

	userFlag := cmd.Flags().Lookup("user")
	if userFlag == nil {
		t.Error("--user flag not found")
	}

	// Check short versions
	if assistantFlag.Shorthand != "a" {
		t.Errorf("Expected --assistant short flag to be 'a', got '%s'", assistantFlag.Shorthand)
	}

	if userFlag.Shorthand != "u" {
		t.Errorf("Expected --user short flag to be 'u', got '%s'", userFlag.Shorthand)
	}
}