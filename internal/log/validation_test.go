package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateLogCommand(t *testing.T) {
	testCases := []struct {
		name        string
		flags       LogFlags
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid assistant mode",
			flags: LogFlags{
				Assistant: true,
				User:      false,
				Read:      false,
			},
			expectError: false,
		},
		{
			name: "Valid user mode",
			flags: LogFlags{
				Assistant: false,
				User:      true,
				Read:      false,
			},
			expectError: false,
		},
		{
			name: "Valid read mode with required fields",
			flags: LogFlags{
				Assistant: false,
				User:      false,
				Read:      true,
				AgentID:   "arch-001",
				Lines:     50,
				Format:    "json",
			},
			expectError: false,
		},
		{
			name: "Multiple modes selected - assistant and user",
			flags: LogFlags{
				Assistant: true,
				User:      true,
				Read:      false,
			},
			expectError: true,
			errorMsg:    "exactly one mode required: --assistant, --user, or --read",
		},
		{
			name: "Multiple modes selected - all three",
			flags: LogFlags{
				Assistant: true,
				User:      true,
				Read:      true,
				AgentID:   "arch-001",
			},
			expectError: true,
			errorMsg:    "exactly one mode required: --assistant, --user, or --read",
		},
		{
			name: "No mode selected",
			flags: LogFlags{
				Assistant: false,
				User:      false,
				Read:      false,
			},
			expectError: true,
			errorMsg:    "exactly one mode required: --assistant, --user, or --read",
		},
		{
			name: "Read mode without agent-id",
			flags: LogFlags{
				Assistant: false,
				User:      false,
				Read:      true,
				AgentID:   "",
			},
			expectError: true,
			errorMsg:    "--agent-id required when using --read",
		},
		{
			name: "Read mode with invalid agent-id format",
			flags: LogFlags{
				Assistant: false,
				User:      false,
				Read:      true,
				AgentID:   "invalid@id",
			},
			expectError: true,
			errorMsg:    "invalid agent-id format",
		},
		{
			name: "Read mode with invalid agent-id - reserved word",
			flags: LogFlags{
				Assistant: false,
				User:      false,
				Read:      true,
				AgentID:   "main",
				Lines:     50,
				Format:    "json",
			},
			expectError: true,
			errorMsg:    "invalid agent-id format",
		},
		{
			name: "Read mode with lines too small",
			flags: LogFlags{
				Assistant: false,
				User:      false,
				Read:      true,
				AgentID:   "arch-001",
				Lines:     0,
			},
			expectError: true,
			errorMsg:    "--lines must be between 1 and 1000",
		},
		{
			name: "Read mode with lines too large",
			flags: LogFlags{
				Assistant: false,
				User:      false,
				Read:      true,
				AgentID:   "arch-001",
				Lines:     1001,
			},
			expectError: true,
			errorMsg:    "--lines must be between 1 and 1000",
		},
		{
			name: "Read mode with maximum valid lines",
			flags: LogFlags{
				Assistant: false,
				User:      false,
				Read:      true,
				AgentID:   "arch-001",
				Lines:     1000,
				Format:    "json",
			},
			expectError: false,
		},
		{
			name: "Read mode with minimum valid lines",
			flags: LogFlags{
				Assistant: false,
				User:      false,
				Read:      true,
				AgentID:   "arch-001",
				Lines:     1,
				Format:    "json",
			},
			expectError: false,
		},
		{
			name: "Read mode with invalid format",
			flags: LogFlags{
				Assistant: false,
				User:      false,
				Read:      true,
				AgentID:   "arch-001",
				Lines:     50,
				Format:    "xml",
			},
			expectError: true,
			errorMsg:    "invalid format: must be 'json' or 'text'",
		},
		{
			name: "Read mode with text format",
			flags: LogFlags{
				Assistant: false,
				User:      false,
				Read:      true,
				AgentID:   "arch-001",
				Lines:     50,
				Format:    "text",
			},
			expectError: false,
		},
		{
			name: "Read mode with empty format (should use default)",
			flags: LogFlags{
				Assistant: false,
				User:      false,
				Read:      true,
				AgentID:   "arch-001",
				Lines:     50,
				Format:    "",
			},
			expectError: false, // Empty format should use default "json"
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateLogCommand(tc.flags)
			
			if tc.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateLogCommand_EdgeCases(t *testing.T) {
	testCases := []struct {
		name        string
		flags       LogFlags
		expectError bool
		errorMsg    string
	}{
		{
			name: "Read mode with complex valid agent-id",
			flags: LogFlags{
				Read:    true,
				AgentID: "arch-auth-system-v2",
				Lines:   50,
				Format:  "json",
			},
			expectError: false,
		},
		{
			name: "Read mode with underscores in agent-id",
			flags: LogFlags{
				Read:    true,
				AgentID: "dev_frontend_2024",
				Lines:   50,
				Format:  "json",
			},
			expectError: false,
		},
		{
			name: "Read mode with mixed case agent-id",
			flags: LogFlags{
				Read:    true,
				AgentID: "Test-Agent-001",
				Lines:   50,
				Format:  "json",
			},
			expectError: false,
		},
		{
			name: "Read mode with agent-id too short",
			flags: LogFlags{
				Read:    true,
				AgentID: "a",
				Lines:   50,
				Format:  "json",
			},
			expectError: true,
			errorMsg:    "invalid agent-id format",
		},
		{
			name: "Read mode with agent-id too long",
			flags: LogFlags{
				Read:    true,
				AgentID: "a" + string(make([]byte, 65)), // 66 characters
				Lines:   50,
				Format:  "json",
			},
			expectError: true,
			errorMsg:    "invalid agent-id format",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateLogCommand(tc.flags)
			
			if tc.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}