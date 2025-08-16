package log

import (
	"fmt"
	"io"
	"strings"
	"testing"
)

// TestProcessToolCallEdgeCases tests edge cases and error scenarios
func TestProcessToolCallEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		isPostHook  bool
		expectNil   bool
		expectError bool
		validateFn  func(*testing.T, *HookData, error)
	}{
		{
			name:        "empty input",
			input:       "",
			isPostHook:  false,
			expectError: true,
		},
		{
			name:        "invalid JSON - missing quotes",
			input:       `{tool_name: Task, tool_input: {}}`,
			isPostHook:  false,
			expectError: true,
		},
		{
			name:        "invalid JSON - trailing comma",
			input:       `{"tool_name": "Task", "tool_input": {},}`,
			isPostHook:  false,
			expectError: true,
		},
		{
			name:        "malformed JSON - unclosed brace",
			input:       `{"tool_name": "Task", "tool_input": {`,
			isPostHook:  false,
			expectError: true,
		},
		{
			name:       "missing tool_input field",
			input:      `{"tool_name": "Task"}`,
			isPostHook: false,
			expectNil:  true, // Should be filtered out
		},
		{
			name:       "null tool_input",
			input:      `{"tool_name": "Task", "tool_input": null}`,
			isPostHook: false,
			expectNil:  true,
		},
		{
			name:       "empty tool_input object",
			input:      `{"tool_name": "Task", "tool_input": {}}`,
			isPostHook: false,
			expectNil:  true, // No subagent_type
		},
		{
			name:       "tool_input with wrong subagent_type type",
			input:      `{"tool_name": "Task", "tool_input": {"subagent_type": 123}}`,
			isPostHook: false,
			expectNil:  true,
		},
		{
			name:       "tool_input with null subagent_type",
			input:      `{"tool_name": "Task", "tool_input": {"subagent_type": null}}`,
			isPostHook: false,
			expectNil:  true,
		},
		{
			name:       "very large JSON input",
			input:      fmt.Sprintf(`{"tool_name": "Task", "tool_input": {"subagent_type": "the-test", "description": "%s", "prompt": "SessionId: dev-large AgentId: test-001\nTest"}}`, strings.Repeat("x", 10000)),
			isPostHook: false,
			expectNil:  false,
			validateFn: func(t *testing.T, data *HookData, err error) {
				if err != nil {
					t.Errorf("Large input should be handled: %v", err)
				}
				if data == nil {
					t.Error("Expected data for large input")
				}
			},
		},
		{
			name:       "unicode characters in JSON",
			input:      `{"tool_name": "Task", "tool_input": {"subagent_type": "the-测试", "description": "测试描述", "prompt": "SessionId: dev-unicode AgentId: test-unicode\n测试提示"}}`,
			isPostHook: false,
			expectNil:  false,
			validateFn: func(t *testing.T, data *HookData, err error) {
				if err != nil {
					t.Errorf("Unicode input should be handled: %v", err)
				}
				if data.AgentType != "the-测试" {
					t.Errorf("Expected unicode agent type, got: %s", data.AgentType)
				}
				if data.Description != "测试描述" {
					t.Errorf("Expected unicode description, got: %s", data.Description)
				}
			},
		},
		{
			name:       "PostToolUse with very large output",
			input:      fmt.Sprintf(`{"tool_name": "Task", "tool_input": {"subagent_type": "the-test", "description": "Test", "prompt": "SessionId: dev-large-output AgentId: test-002\nTest"}, "output": "%s"}`, strings.Repeat("A", 5000)),
			isPostHook: true,
			expectNil:  false,
			validateFn: func(t *testing.T, data *HookData, err error) {
				if err != nil {
					t.Errorf("Large output should be handled: %v", err)
				}
				if len(data.OutputSummary) != 1000+len("... [truncated 4000 chars]") {
					t.Errorf("Expected truncated output, got length: %d", len(data.OutputSummary))
				}
				if !strings.Contains(data.OutputSummary, "... [truncated 4000 chars]") {
					t.Error("Expected truncation message in output")
				}
			},
		},
		{
			name:       "special characters in session and agent IDs",
			input:      `{"tool_name": "Task", "tool_input": {"subagent_type": "the-special", "description": "Test", "prompt": "SessionId: dev-session_with-special.chars AgentId: agent_with-special.chars\nTest"}}`,
			isPostHook: false,
			expectNil:  false,
			validateFn: func(t *testing.T, data *HookData, err error) {
				if data.SessionID != "dev-session_with-special.chars" {
					t.Errorf("Expected special characters in session ID, got: %s", data.SessionID)
				}
				// AgentID with dots is invalid, so should generate fallback
				if !strings.HasPrefix(data.AgentID, "the-special-") {
					t.Errorf("Expected fallback AgentID starting with 'the-special-', got: %s", data.AgentID)
				}
			},
		},
		{
			name:       "missing session and agent IDs - should find latest session",
			input:      `{"tool_name": "Task", "tool_input": {"subagent_type": "the-fallback", "description": "Test", "prompt": "Just a prompt without IDs"}}`,
			isPostHook: false,
			expectNil:  false,
			validateFn: func(t *testing.T, data *HookData, err error) {
				// Should generate fallback agent ID since no explicit AgentID provided
				if !strings.HasPrefix(data.AgentID, "the-fallback-") {
					t.Errorf("Expected fallback AgentID starting with 'the-fallback-', got: %s", data.AgentID)
				}
				// Session ID might be empty or found from latest session
				// Don't test specific value as it depends on environment
			},
		},
		{
			name:       "case sensitivity in tool name",
			input:      `{"tool_name": "task", "tool_input": {"subagent_type": "the-case", "description": "Test"}}`,
			isPostHook: false,
			expectNil:  true, // Should be filtered - only "Task" is valid
		},
		{
			name:       "case sensitivity in subagent_type prefix",
			input:      `{"tool_name": "Task", "tool_input": {"subagent_type": "The-case", "description": "Test"}}`,
			isPostHook: false,
			expectNil:  true, // Should be filtered - only "the-" prefix is valid
		},
		{
			name:       "nested objects in tool_input",
			input:      `{"tool_name": "Task", "tool_input": {"subagent_type": "the-nested", "description": "Test", "nested": {"key": "value"}, "prompt": "SessionId: dev-nested AgentId: nested-001\nTest"}}`,
			isPostHook: false,
			expectNil:  false,
			validateFn: func(t *testing.T, data *HookData, err error) {
				if data.AgentType != "the-nested" {
					t.Errorf("Expected nested agent type, got: %s", data.AgentType)
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

			if tt.expectNil {
				if result != nil {
					t.Errorf("Expected nil result but got: %+v", result)
				}
				return
			}

			if tt.validateFn != nil {
				tt.validateFn(t, result, err)
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result == nil {
					t.Error("Expected non-nil result")
				}
			}
		})
	}
}

// TestProcessToolCallIOErrors tests I/O related errors
func TestProcessToolCallIOErrors(t *testing.T) {
	tests := []struct {
		name        string
		readerFn    func() io.Reader
		expectError bool
	}{
		{
			name: "read error from reader",
			readerFn: func() io.Reader {
				return &errorReader{err: fmt.Errorf("read error")}
			},
			expectError: true,
		},
		{
			name: "partial read",
			readerFn: func() io.Reader {
				return &partialReader{data: `{"tool_name": "Task"`}
			},
			expectError: true, // Invalid JSON
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := tt.readerFn()
			_, err := ProcessToolCall(reader, false)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

// errorReader always returns an error on Read
type errorReader struct {
	err error
}

func (r *errorReader) Read(p []byte) (n int, err error) {
	return 0, r.err
}

// partialReader returns partial data then EOF
type partialReader struct {
	data string
	read bool
}

func (r *partialReader) Read(p []byte) (n int, err error) {
	if r.read {
		return 0, io.EOF
	}
	r.read = true
	n = copy(p, []byte(r.data))
	return n, nil
}

// TestRegexPatterns tests edge cases in regex pattern matching
func TestRegexPatterns(t *testing.T) {
	sessionTests := []struct {
		name     string
		prompt   string
		expected string
	}{
		{
			name:     "session ID with underscores and hyphens",
			prompt:   "SessionId: dev_session-with_mixed-chars",
			expected: "dev_session-with_mixed-chars",
		},
		{
			name:     "session ID with numbers",
			prompt:   "SessionId: dev-session-12345-test",
			expected: "dev-session-12345-test",
		},
		{
			name:     "session ID with dots",
			prompt:   "SessionId: dev.session.with.dots",
			expected: "dev.session.with.dots",
		},
		{
			name:     "multiple session IDs - should match first",
			prompt:   "SessionId: first-session SessionId: second-session",
			expected: "first-session",
		},
		{
			name:     "session ID at start of string",
			prompt:   "SessionId: at-start\nContent follows",
			expected: "at-start",
		},
		{
			name:     "session ID with tab separator",
			prompt:   "SessionId:\tdev-tab-separated",
			expected: "dev-tab-separated",
		},
		{
			name:     "session ID with multiple spaces",
			prompt:   "SessionId:     dev-many-spaces",
			expected: "dev-many-spaces",
		},
		{
			name:     "malformed - no space after colon",
			prompt:   "SessionId:no-space-after-colon",
			expected: "no-space-after-colon", // Regex should still match
		},
		{
			name:     "partial match - should not match",
			prompt:   "NotSessionId: should-not-match",
			expected: "",
		},
	}

	for _, tt := range sessionTests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractSessionID(tt.prompt)
			if result != tt.expected {
				t.Errorf("ExtractSessionID(%q) = %q, expected %q", tt.prompt, result, tt.expected)
			}
		})
	}

	agentTests := []struct {
		name     string
		prompt   string
		expected string
	}{
		{
			name:     "agent ID with complex characters",
			prompt:   "AgentId: the-complex_agent.v1-test",
			expected: "the-complex_agent.v1-test",
		},
		{
			name:     "agent ID with version number",
			prompt:   "AgentId: agent-v2.1.0",
			expected: "agent-v2.1.0",
		},
		{
			name:     "multiple agent IDs - should match first",
			prompt:   "AgentId: first-agent AgentId: second-agent",
			expected: "first-agent",
		},
		{
			name:     "agent ID with UUID format",
			prompt:   "AgentId: agent-12345678-1234-5678-9abc-123456789012",
			expected: "agent-12345678-1234-5678-9abc-123456789012",
		},
	}

	for _, tt := range agentTests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractAgentID(tt.prompt)
			if result != tt.expected {
				t.Errorf("ExtractAgentID(%q) = %q, expected %q", tt.prompt, result, tt.expected)
			}
		})
	}
}

// TestTruncateOutputEdgeCases tests edge cases in output truncation
func TestTruncateOutputEdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		output    string
		maxLength int
		expected  string
	}{
		{
			name:      "zero max length",
			output:    "Some output",
			maxLength: 0,
			expected:  "... [truncated 11 chars]",
		},
		{
			name:      "negative max length",
			output:    "Some output",
			maxLength: -1,
			expected:  "... [truncated 11 chars]",
		},
		{
			name:      "max length of 1",
			output:    "AB",
			maxLength: 1,
			expected:  "A... [truncated 1 chars]",
		},
		{
			name:      "unicode characters in output",
			output:    "测试输出内容",
			maxLength: 3,
			expected:  "测... [truncated 15 chars]", // Go's string slicing uses bytes
		},
		{
			name:      "newlines and special characters",
			output:    "Line 1\nLine 2\tTabbed\r\nWindows line ending",
			maxLength: 10,
			expected:  "Line 1\nLin... [truncated 31 chars]",
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
