package log

import (
	"strings"
	"testing"
)

// Test cases from Python hooks to ensure 100% compatibility
func TestExtractSessionIDCompatibility(t *testing.T) {
	tests := []struct {
		name     string
		prompt   string
		expected string
	}{
		{
			name:     "basic session ID extraction",
			prompt:   "SessionId: dev-session-123\nSome other content",
			expected: "dev-session-123",
		},
		{
			name:     "session ID with comma",
			prompt:   "SessionId: dev-session-456, AgentId: agent-789",
			expected: "dev-session-456",
		},
		{
			name:     "session ID with spaces around colon",
			prompt:   "SessionId : dev-session-789",
			expected: "dev-session-789",
		},
		{
			name:     "no session ID",
			prompt:   "Just some text without session info",
			expected: "",
		},
		{
			name:     "session ID at end of line",
			prompt:   "Some text SessionId: dev-end-123",
			expected: "dev-end-123",
		},
		{
			name:     "session ID with newline after",
			prompt:   "SessionId: dev-newline-456\nNext line",
			expected: "dev-newline-456",
		},
		{
			name:     "session ID case sensitive",
			prompt:   "sessionid: dev-lower-case",
			expected: "", // Should not match lowercase
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractSessionID(tt.prompt)
			if result != tt.expected {
				t.Errorf("ExtractSessionID(%q) = %q, expected %q", tt.prompt, result, tt.expected)
			}
		})
	}
}

func TestExtractAgentIDCompatibility(t *testing.T) {
	tests := []struct {
		name     string
		prompt   string
		expected string
	}{
		{
			name:     "basic agent ID extraction",
			prompt:   "AgentId: arch-001\nSome other content",
			expected: "arch-001",
		},
		{
			name:     "agent ID with comma",
			prompt:   "SessionId: dev-session-456, AgentId: agent-789",
			expected: "agent-789",
		},
		{
			name:     "agent ID with spaces around colon",
			prompt:   "AgentId : dev-agent-789",
			expected: "dev-agent-789",
		},
		{
			name:     "no agent ID",
			prompt:   "Just some text without agent info",
			expected: "",
		},
		{
			name:     "agent ID at end of line",
			prompt:   "Some text AgentId: dev-end-123",
			expected: "dev-end-123",
		},
		{
			name:     "agent ID with newline after",
			prompt:   "AgentId: dev-newline-456\nNext line",
			expected: "dev-newline-456",
		},
		{
			name:     "agent ID case sensitive",
			prompt:   "agentid: dev-lower-case",
			expected: "", // Should not match lowercase
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractAgentID(tt.prompt)
			if result != tt.expected {
				t.Errorf("ExtractAgentID(%q) = %q, expected %q", tt.prompt, result, tt.expected)
			}
		})
	}
}

func TestShouldProcessCompatibility(t *testing.T) {
	tests := []struct {
		name      string
		toolName  string
		toolInput map[string]interface{}
		expected  bool
	}{
		{
			name:     "Task with the- prefix - should process",
			toolName: "Task",
			toolInput: map[string]interface{}{
				"subagent_type": "the-architect",
			},
			expected: true,
		},
		{
			name:     "Task with the- prefix different type - should process",
			toolName: "Task",
			toolInput: map[string]interface{}{
				"subagent_type": "the-developer",
			},
			expected: true,
		},
		{
			name:     "Task without the- prefix - should not process",
			toolName: "Task",
			toolInput: map[string]interface{}{
				"subagent_type": "architect",
			},
			expected: false,
		},
		{
			name:     "Non-Task tool with the- prefix - should not process",
			toolName: "Other",
			toolInput: map[string]interface{}{
				"subagent_type": "the-architect",
			},
			expected: false,
		},
		{
			name:     "Task with missing subagent_type - should not process",
			toolName: "Task",
			toolInput: map[string]interface{}{
				"description": "some task",
			},
			expected: false,
		},
		{
			name:     "Task with non-string subagent_type - should not process",
			toolName: "Task",
			toolInput: map[string]interface{}{
				"subagent_type": 123,
			},
			expected: false,
		},
		{
			name:     "Task with empty subagent_type - should not process",
			toolName: "Task",
			toolInput: map[string]interface{}{
				"subagent_type": "",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ShouldProcess(tt.toolName, tt.toolInput)
			if result != tt.expected {
				t.Errorf("ShouldProcess(%q, %v) = %t, expected %t", tt.toolName, tt.toolInput, result, tt.expected)
			}
		})
	}
}

func TestTruncateOutputCompatibility(t *testing.T) {
	tests := []struct {
		name      string
		output    string
		maxLength int
		expected  string
	}{
		{
			name:      "short output - no truncation",
			output:    "Short response",
			maxLength: 1000,
			expected:  "Short response",
		},
		{
			name:      "exactly max length - no truncation",
			output:    strings.Repeat("a", 1000),
			maxLength: 1000,
			expected:  strings.Repeat("a", 1000),
		},
		{
			name:      "over max length - truncation",
			output:    strings.Repeat("a", 1050),
			maxLength: 1000,
			expected:  strings.Repeat("a", 1000) + "... [truncated 50 chars]",
		},
		{
			name:      "small max length for testing",
			output:    "This is a longer message that should be truncated",
			maxLength: 10,
			expected:  "This is a " + "... [truncated 39 chars]",
		},
		{
			name:      "empty output",
			output:    "",
			maxLength: 1000,
			expected:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TruncateOutput(tt.output, tt.maxLength)
			if result != tt.expected {
				t.Errorf("TruncateOutput(%q, %d) = %q, expected %q", tt.output, tt.maxLength, result, tt.expected)
			}
		})
	}
}

func TestProcessToolCallCompatibility(t *testing.T) {
	tests := []struct {
		name        string
		jsonInput   string
		isPostHook  bool
		expectNil   bool
		expectError bool
		validateFn  func(*testing.T, *HookData)
	}{
		{
			name:      "valid PreToolUse Task with the- prefix",
			jsonInput: `{"tool_name":"Task","tool_input":{"subagent_type":"the-architect","description":"Design system","prompt":"SessionId: dev-123 AgentId: arch-001\nPlease design the system"},"session_id":"dev-123","hook_event_name":"PreToolUse"}`,
			isPostHook: false,
			expectNil:  false,
			validateFn: func(t *testing.T, data *HookData) {
				if data.Event != "agent_start" {
					t.Errorf("Expected event 'agent_start', got %q", data.Event)
				}
				if data.AgentType != "the-architect" {
					t.Errorf("Expected agent_type 'the-architect', got %q", data.AgentType)
				}
				if data.AgentID != "arch-001" {
					t.Errorf("Expected agent_id 'arch-001', got %q", data.AgentID)
				}
				if data.SessionID != "dev-123" {
					t.Errorf("Expected session_id 'dev-123', got %q", data.SessionID)
				}
				if data.Instruction == "" {
					t.Error("Expected instruction to be set for agent_start")
				}
				if data.OutputSummary != "" {
					t.Error("Expected output_summary to be empty for agent_start")
				}
			},
		},
		{
			name:      "valid PostToolUse Task with the- prefix",
			jsonInput: `{"tool_name":"Task","tool_input":{"subagent_type":"the-developer","description":"Implement feature","prompt":"SessionId: dev-456 AgentId: dev-002\nImplement the feature"},"session_id":"dev-456","hook_event_name":"PostToolUse","output":"Feature implemented successfully"}`,
			isPostHook: true,
			expectNil:  false,
			validateFn: func(t *testing.T, data *HookData) {
				if data.Event != "agent_complete" {
					t.Errorf("Expected event 'agent_complete', got %q", data.Event)
				}
				if data.AgentType != "the-developer" {
					t.Errorf("Expected agent_type 'the-developer', got %q", data.AgentType)
				}
				if data.AgentID != "dev-002" {
					t.Errorf("Expected agent_id 'dev-002', got %q", data.AgentID)
				}
				if data.SessionID != "dev-456" {
					t.Errorf("Expected session_id 'dev-456', got %q", data.SessionID)
				}
				if data.Instruction != "" {
					t.Error("Expected instruction to be empty for agent_complete")
				}
				if data.OutputSummary != "Feature implemented successfully" {
					t.Errorf("Expected output_summary 'Feature implemented successfully', got %q", data.OutputSummary)
				}
			},
		},
		{
			name:      "Non-Task tool - should be filtered",
			jsonInput: `{"tool_name":"Other","tool_input":{"subagent_type":"the-architect","description":"Design system","prompt":"Test"},"session_id":"dev-123","hook_event_name":"PreToolUse"}`,
			isPostHook: false,
			expectNil:  true,
		},
		{
			name:      "Task without the- prefix - should be filtered",
			jsonInput: `{"tool_name":"Task","tool_input":{"subagent_type":"architect","description":"Design system","prompt":"Test"},"session_id":"dev-123","hook_event_name":"PreToolUse"}`,
			isPostHook: false,
			expectNil:  true,
		},
		{
			name:        "Invalid JSON - should error",
			jsonInput:   `invalid json`,
			isPostHook:  false,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.jsonInput)
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
				t.Error("Expected non-nil result but got nil")
				return
			}

			if tt.validateFn != nil {
				tt.validateFn(t, result)
			}

			// Validate common fields
			if result.Timestamp == "" {
				t.Error("Expected timestamp to be set")
			}
			if result.Description == "" {
				t.Error("Expected description to be set")
			}
		})
	}
}