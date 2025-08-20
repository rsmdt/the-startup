package log

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestProcessToolCall tests the complete flow from JSON input to HookData output
func TestProcessToolCall(t *testing.T) {
	tests := []struct {
		name        string
		jsonInput   string
		isPostHook  bool
		expectNil   bool
		expectError bool
		validateFn  func(*testing.T, *HookData)
	}{
		{
			name: "PreToolUse with valid agent",
			jsonInput: `{
				"tool_name":"Task",
				"tool_input":{
					"subagent_type":"the-architect",
					"description":"Design system",
					"prompt":"SessionId: dev-123 AgentId: arch-001\nPlease design the system"
				},
				"session_id":"dev-123",
				"hook_event_name":"PreToolUse"
			}`,
			isPostHook: false,
			expectNil:  false,
			validateFn: func(t *testing.T, data *HookData) {
				if data.Role != "user" {
					t.Errorf("Expected role 'user', got %q", data.Role)
				}
				if data.AgentID != "arch-001" {
					t.Errorf("Expected agent_id 'arch-001', got %q", data.AgentID)
				}
				if data.SessionID != "dev-123" {
					t.Errorf("Expected session_id 'dev-123', got %q", data.SessionID)
				}
				if data.Content == "" {
					t.Error("Expected content to be set for PreToolUse")
				}
				// Verify timestamp format (ISO 8601 with Z suffix)
				if !strings.HasSuffix(data.Timestamp, "Z") {
					t.Errorf("Timestamp should end with Z, got %q", data.Timestamp)
				}
			},
		},
		{
			name: "PostToolUse with output",
			jsonInput: `{
				"tool_name":"Task",
				"tool_input":{
					"subagent_type":"the-developer",
					"description":"Implement feature",
					"prompt":"SessionId: dev-456 AgentId: dev-002\nImplement the feature"
				},
				"session_id":"dev-456",
				"hook_event_name":"PostToolUse",
				"output":"Feature implemented successfully"
			}`,
			isPostHook: true,
			expectNil:  false,
			validateFn: func(t *testing.T, data *HookData) {
				if data.Role != "assistant" {
					t.Errorf("Expected role 'assistant', got %q", data.Role)
				}
				if data.AgentID != "dev-002" {
					t.Errorf("Expected agent_id 'dev-002', got %q", data.AgentID)
				}
				if data.Content != "Feature implemented successfully" {
					t.Errorf("Expected content to match output, got %q", data.Content)
				}
			},
		},
		{
			name: "Long output truncation",
			jsonInput: `{
				"tool_name":"Task",
				"tool_input":{
					"subagent_type":"the-architect",
					"prompt":"SessionId: dev-789 AgentId: arch-002\nDesign"
				},
				"hook_event_name":"PostToolUse",
				"output":"` + strings.Repeat("The architecture has been designed with the following components: ", 50) + `"
			}`,
			isPostHook: true,
			expectNil:  false,
			validateFn: func(t *testing.T, data *HookData) {
				// Output should be preserved in full (no truncation in the data itself)
				if !strings.Contains(data.Content, "The architecture has been designed") {
					t.Error("Expected content to contain the output")
				}
			},
		},
		{
			name:       "Non-Task tool should be filtered",
			jsonInput:  `{"tool_name":"Edit","tool_input":{"file":"test.go"},"session_id":"dev-789"}`,
			isPostHook: false,
			expectNil:  true,
		},
		{
			name:       "Task without the- prefix should be filtered",
			jsonInput:  `{"tool_name":"Task","tool_input":{"subagent_type":"architect"},"session_id":"dev-000"}`,
			isPostHook: false,
			expectNil:  true,
		},
		{
			name:        "Invalid JSON",
			jsonInput:   `{invalid json}`,
			isPostHook:  false,
			expectError: true,
		},
		{
			name: "Missing prompt extracts from session",
			jsonInput: `{
				"tool_name":"Task",
				"tool_input":{"subagent_type":"the-tester"},
				"session_id":"dev-222"
			}`,
			isPostHook: false,
			expectNil:  false,
			validateFn: func(t *testing.T, data *HookData) {
				if data.Role != "user" {
					t.Errorf("Expected role 'user', got %q", data.Role)
				}
				// AgentID should be generated since it's not in the prompt
				if data.AgentID == "" {
					t.Error("Expected generated AgentID")
				}
			},
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

// TestJSONSerialization verifies HookData serializes correctly
func TestJSONSerialization(t *testing.T) {
	hookData := &HookData{
		Role:      "user",
		Content:   "Test content",
		Timestamp: "2025-01-11T10:30:00.000Z",
		SessionID: "internal-session", // Internal field
		AgentID:   "internal-agent",   // Internal field
	}

	jsonBytes, err := json.Marshal(hookData)
	if err != nil {
		t.Fatalf("Failed to marshal HookData: %v", err)
	}

	// Parse back to verify structure
	var jsonMap map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &jsonMap); err != nil {
		t.Fatalf("Failed to unmarshal to map: %v", err)
	}

	// Check required fields are present
	requiredFields := []string{"role", "content", "timestamp"}
	for _, field := range requiredFields {
		if _, exists := jsonMap[field]; !exists {
			t.Errorf("Required field %q missing from JSON output", field)
		}
	}

	// Internal fields should NOT be in JSON
	internalFields := []string{"session_id", "agent_id", "SessionID", "AgentID"}
	for _, field := range internalFields {
		if _, exists := jsonMap[field]; exists {
			t.Errorf("Internal field %q should not be in JSON output", field)
		}
	}
}

// TestShouldProcess verifies filtering logic
func TestShouldProcess(t *testing.T) {
	tests := []struct {
		name          string
		toolName      string
		toolInput     map[string]interface{}
		shouldProcess bool
	}{
		{
			name:          "Task with the- prefix",
			toolName:      "Task",
			toolInput:     map[string]interface{}{"subagent_type": "the-architect"},
			shouldProcess: true,
		},
		{
			name:          "Task without the- prefix",
			toolName:      "Task",
			toolInput:     map[string]interface{}{"subagent_type": "architect"},
			shouldProcess: false,
		},
		{
			name:          "Non-Task tool",
			toolName:      "Edit",
			toolInput:     map[string]interface{}{"file": "test.go"},
			shouldProcess: false,
		},
		{
			name:          "Task with nil tool_input",
			toolName:      "Task",
			toolInput:     nil,
			shouldProcess: false,
		},
		{
			name:          "Case sensitive Task",
			toolName:      "task",
			toolInput:     map[string]interface{}{"subagent_type": "the-test"},
			shouldProcess: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ShouldProcess(tt.toolName, tt.toolInput)
			if result != tt.shouldProcess {
				t.Errorf("ShouldProcess(%q, %v) = %v, expected %v",
					tt.toolName, tt.toolInput, result, tt.shouldProcess)
			}
		})
	}
}

// TestExtractSessionID tests session ID extraction
func TestExtractSessionID(t *testing.T) {
	tests := []struct {
		name     string
		prompt   string
		expected string
	}{
		{
			name:     "standard format",
			prompt:   "SessionId: dev-123\nSome prompt",
			expected: "dev-123",
		},
		{
			name:     "case insensitive",
			prompt:   "sessionid: dev-789\nSome prompt",
			expected: "dev-789",
		},
		{
			name:     "no session ID",
			prompt:   "Just a regular prompt",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractSessionID(tt.prompt)
			if result != tt.expected {
				t.Errorf("ExtractSessionID(%q) = %q, expected %q",
					tt.prompt, result, tt.expected)
			}
		})
	}
}

// TestExtractAgentID tests agent ID extraction
func TestExtractAgentID(t *testing.T) {
	tests := []struct {
		name     string
		prompt   string
		expected string
	}{
		{
			name:     "standard format",
			prompt:   "AgentId: arch-001\nContent",
			expected: "arch-001",
		},
		{
			name:     "case insensitive",
			prompt:   "agentid: test-003\nContent",
			expected: "test-003",
		},
		{
			name:     "no agent ID",
			prompt:   "Just a regular prompt",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractAgentID(tt.prompt)
			if result != tt.expected {
				t.Errorf("ExtractAgentID(%q) = %q, expected %q",
					tt.prompt, result, tt.expected)
			}
		})
	}
}

// TestFileWriting tests that logs are written correctly
func TestFileWriting(t *testing.T) {
	tempDir := t.TempDir()
	os.Setenv("CLAUDE_PROJECT_DIR", tempDir)
	defer os.Unsetenv("CLAUDE_PROJECT_DIR")

	// Create .the-startup directory
	startupDir := filepath.Join(tempDir, ".the-startup")
	os.MkdirAll(startupDir, 0755)

	sessionID := "test-session"
	agentID := "test-agent"
	
	hookData := &HookData{
		Role:      "user",
		Content:   "Test task content",
		Timestamp: "2025-01-11T10:30:00.000Z",
		SessionID: sessionID,
		AgentID:   agentID,
	}

	// Write session log
	err := WriteSessionLog(sessionID, hookData)
	if err != nil {
		t.Fatalf("Failed to write session log: %v", err)
	}

	// Verify file was created
	sessionFile := filepath.Join(startupDir, sessionID, "agent-instructions.jsonl")
	if _, err := os.Stat(sessionFile); os.IsNotExist(err) {
		t.Error("Session log file was not created")
	}

	// Write agent context
	err = WriteAgentContext(sessionID, agentID, hookData)
	if err != nil {
		t.Fatalf("Failed to write agent context: %v", err)
	}

	// Verify agent context file was created
	agentFile := filepath.Join(startupDir, sessionID, agentID+".jsonl")
	if _, err := os.Stat(agentFile); os.IsNotExist(err) {
		t.Error("Agent context file was not created")
	}
}