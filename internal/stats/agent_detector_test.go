package stats

import (
	"encoding/json"
	"testing"
	"time"
)

// Test Task tool subagent_type extraction with 100% confidence
func TestDetectAgent_TaskToolWith100PercentConfidence(t *testing.T) {
	// Test case: Task tool with subagent_type parameter
	entry := ClaudeLogEntry{
		Type:      "assistant",
		SessionID: "test-session",
		Timestamp: time.Now(),
		Assistant: &AssistantMessage{
			Text: "I'll delegate this to a specialist",
			ToolUses: []ToolUse{
				{
					ID:   "tool_001",
					Name: "Task",
					Parameters: json.RawMessage(`{
						"instructions": "Analyze the code structure",
						"subagent_type": "the-spec-writer"
					}`),
				},
			},
		},
	}

	agent, confidence := DetectAgent(entry)

	if agent != "the-spec-writer" {
		t.Errorf("Expected agent 'the-spec-writer', got '%s'", agent)
	}
	if confidence != 1.0 {
		t.Errorf("Expected confidence 1.0 for Task tool detection, got %f", confidence)
	}
}

// Test Task tool without subagent_type falls back to other methods
func TestDetectAgent_TaskToolWithoutSubagentType(t *testing.T) {
	entry := ClaudeLogEntry{
		Type:      "assistant",
		SessionID: "test-session",
		Timestamp: time.Now(),
		Assistant: &AssistantMessage{
			Text: "Using the-implementor agent for this task",
			ToolUses: []ToolUse{
				{
					ID:         "tool_001",
					Name:       "Task",
					Parameters: json.RawMessage(`{"instructions": "Build the feature"}`),
				},
			},
		},
	}

	agent, confidence := DetectAgent(entry)

	// Should fall back to pattern matching
	if agent != "the-implementor" {
		t.Errorf("Expected agent 'the-implementor' from pattern, got '%s'", agent)
	}
	if confidence != 0.90 {
		t.Errorf("Expected confidence 0.90 for pattern match, got %f", confidence)
	}
}

// Test Task tool in user result message
func TestDetectAgent_TaskToolInUserResult(t *testing.T) {
	entry := ClaudeLogEntry{
		Type:      "user",
		SessionID: "test-session",
		Timestamp: time.Now(),
		User: &UserMessage{
			ToolName:  "Task",
			ToolUseID: "tool_001",
			Result: json.RawMessage(`{
				"status": "completed",
				"subagent_type": "the-test-writer",
				"output": "Tests written successfully"
			}`),
		},
	}

	agent, confidence := DetectAgent(entry)

	if agent != "the-test-writer" {
		t.Errorf("Expected agent 'the-test-writer', got '%s'", agent)
	}
	if confidence != 1.0 {
		t.Errorf("Expected confidence 1.0 for Task result detection, got %f", confidence)
	}
}

// Test pattern matching with "delegating to" (95% confidence)
func TestDetectAgent_DelegatingPattern95Percent(t *testing.T) {
	entry := ClaudeLogEntry{
		Type:      "assistant",
		SessionID: "test-session",
		Timestamp: time.Now(),
		Assistant: &AssistantMessage{
			Text: "I'm delegating to the-reviewer for code quality checks",
		},
	}

	agent, confidence := DetectAgent(entry)

	if agent != "the-reviewer" {
		t.Errorf("Expected agent 'the-reviewer', got '%s'", agent)
	}
	if confidence != 0.95 {
		t.Errorf("Expected confidence 0.95 for 'delegating to' pattern, got %f", confidence)
	}
}

// Test pattern matching with "using agent" (90% confidence)
func TestDetectAgent_UsingAgentPattern90Percent(t *testing.T) {
	entry := ClaudeLogEntry{
		Type:      "assistant",
		SessionID: "test-session",
		Timestamp: time.Now(),
		Assistant: &AssistantMessage{
			Text: "I'll be using the-documenter agent to create the documentation",
		},
	}

	agent, confidence := DetectAgent(entry)

	if agent != "the-documenter" {
		t.Errorf("Expected agent 'the-documenter', got '%s'", agent)
	}
	if confidence != 0.90 {
		t.Errorf("Expected confidence 0.90 for 'using agent' pattern, got %f", confidence)
	}
}

// Test pattern matching without "the-" prefix normalizes correctly
func TestDetectAgent_PatternNormalizesAgentName(t *testing.T) {
	entry := ClaudeLogEntry{
		Type:      "assistant",
		SessionID: "test-session",
		Timestamp: time.Now(),
		Assistant: &AssistantMessage{
			Text: "I'm using spec-writer agent for specifications",
		},
	}

	agent, confidence := DetectAgent(entry)

	if agent != "the-spec-writer" {
		t.Errorf("Expected normalized agent 'the-spec-writer', got '%s'", agent)
	}
	if confidence != 0.90 {
		t.Errorf("Expected confidence 0.90, got %f", confidence)
	}
}

// Test pattern matching with "consulting with specialist" (75% confidence)
func TestDetectAgent_ConsultingPattern75Percent(t *testing.T) {
	entry := ClaudeLogEntry{
		Type:      "assistant",
		SessionID: "test-session",
		Timestamp: time.Now(),
		Assistant: &AssistantMessage{
			Text: "I'm consulting with debugger specialist to resolve this issue",
		},
	}

	agent, confidence := DetectAgent(entry)

	if agent != "the-debugger" {
		t.Errorf("Expected agent 'the-debugger', got '%s'", agent)
	}
	if confidence != 0.75 {
		t.Errorf("Expected confidence 0.75 for 'consulting' pattern, got %f", confidence)
	}
}

// Test pattern matching with "invoke" (80% confidence)
func TestDetectAgent_InvokePattern80Percent(t *testing.T) {
	entry := ClaudeLogEntry{
		Type:      "assistant",
		SessionID: "test-session",
		Timestamp: time.Now(),
		Assistant: &AssistantMessage{
			Text: "Let me invoke the-implementor for this complex feature",
		},
	}

	agent, confidence := DetectAgent(entry)

	if agent != "the-implementor" {
		t.Errorf("Expected agent 'the-implementor', got '%s'", agent)
	}
	if confidence != 0.80 {
		t.Errorf("Expected confidence 0.80 for 'invoke' pattern, got %f", confidence)
	}
}

// Test command context inference for /s:specify (70% confidence)
func TestDetectAgent_CommandContextSpecify(t *testing.T) {
	entry := ClaudeLogEntry{
		Type:      "assistant",
		SessionID: "test-session",
		Timestamp: time.Now(),
		Assistant: &AssistantMessage{
			Text: "I'll run /s:specify to create the specification document",
		},
	}

	agent, confidence := DetectAgent(entry)

	if agent != "the-spec-writer" {
		t.Errorf("Expected agent 'the-spec-writer' from command context, got '%s'", agent)
	}
	if confidence != 0.70 {
		t.Errorf("Expected confidence 0.70 for command context, got %f", confidence)
	}
}

// Test command context inference for /s:implement
func TestDetectAgent_CommandContextImplement(t *testing.T) {
	entry := ClaudeLogEntry{
		Type:      "assistant",
		SessionID: "test-session",
		Timestamp: time.Now(),
		Assistant: &AssistantMessage{
			Text: "Running /s:implement to build this feature",
		},
	}

	agent, confidence := DetectAgent(entry)

	if agent != "the-implementor" {
		t.Errorf("Expected agent 'the-implementor', got '%s'", agent)
	}
	if confidence != 0.70 {
		t.Errorf("Expected confidence 0.70, got %f", confidence)
	}
}

// Test command context inference for /s:test
func TestDetectAgent_CommandContextTest(t *testing.T) {
	entry := ClaudeLogEntry{
		Type:      "assistant",
		SessionID: "test-session",
		Timestamp: time.Now(),
		Assistant: &AssistantMessage{
			Text: "Let's use /s:test to create comprehensive test coverage",
		},
	}

	agent, confidence := DetectAgent(entry)

	if agent != "the-test-writer" {
		t.Errorf("Expected agent 'the-test-writer', got '%s'", agent)
	}
	if confidence != 0.70 {
		t.Errorf("Expected confidence 0.70, got %f", confidence)
	}
}

// Test command context with command tags
func TestDetectAgent_CommandContextWithTags(t *testing.T) {
	entry := ClaudeLogEntry{
		Type:      "assistant",
		SessionID: "test-session",
		Timestamp: time.Now(),
		Assistant: &AssistantMessage{
			Text: "<command-name>/s:review</command-name> Running code review",
		},
	}

	agent, confidence := DetectAgent(entry)

	if agent != "the-reviewer" {
		t.Errorf("Expected agent 'the-reviewer', got '%s'", agent)
	}
	if confidence != 0.70 {
		t.Errorf("Expected confidence 0.70, got %f", confidence)
	}
}

// Test no agent detected returns empty string and 0 confidence
func TestDetectAgent_NoAgentDetected(t *testing.T) {
	entry := ClaudeLogEntry{
		Type:      "assistant",
		SessionID: "test-session",
		Timestamp: time.Now(),
		Assistant: &AssistantMessage{
			Text: "Let me analyze this code for you",
		},
	}

	agent, confidence := DetectAgent(entry)

	if agent != "" {
		t.Errorf("Expected empty agent, got '%s'", agent)
	}
	if confidence != 0.0 {
		t.Errorf("Expected confidence 0.0 for no detection, got %f", confidence)
	}
}

// Test priority: Task tool > Pattern > Context
func TestDetectAgent_PriorityTaskOverPattern(t *testing.T) {
	// Entry has both Task tool and pattern
	entry := ClaudeLogEntry{
		Type:      "assistant",
		SessionID: "test-session",
		Timestamp: time.Now(),
		Assistant: &AssistantMessage{
			Text: "delegating to the-reviewer", // Pattern would give the-reviewer
			ToolUses: []ToolUse{
				{
					ID:   "tool_001",
					Name: "Task",
					Parameters: json.RawMessage(`{
						"subagent_type": "the-spec-writer"
					}`),
				},
			},
		},
	}

	agent, confidence := DetectAgent(entry)

	// Task tool should take priority
	if agent != "the-spec-writer" {
		t.Errorf("Expected 'the-spec-writer' from Task tool priority, got '%s'", agent)
	}
	if confidence != 1.0 {
		t.Errorf("Expected confidence 1.0 from Task tool, got %f", confidence)
	}
}

// Test priority: Pattern > Context
func TestDetectAgent_PriorityPatternOverContext(t *testing.T) {
	// Entry has both pattern and command context
	entry := ClaudeLogEntry{
		Type:      "assistant",
		SessionID: "test-session",
		Timestamp: time.Now(),
		Assistant: &AssistantMessage{
			Text: "Using the-documenter agent with /s:test command", // Has both pattern and command
		},
	}

	agent, confidence := DetectAgent(entry)

	// Pattern should take priority over command context
	if agent != "the-documenter" {
		t.Errorf("Expected 'the-documenter' from pattern priority, got '%s'", agent)
	}
	if confidence != 0.90 {
		t.Errorf("Expected confidence 0.90 from pattern, got %f", confidence)
	}
}

// Test corrupted JSON in Task parameters
func TestDetectAgent_CorruptedTaskParameters(t *testing.T) {
	entry := ClaudeLogEntry{
		Type:      "assistant",
		SessionID: "test-session",
		Timestamp: time.Now(),
		Assistant: &AssistantMessage{
			Text: "Delegating to the-spec-writer",
			ToolUses: []ToolUse{
				{
					ID:         "tool_001",
					Name:       "Task",
					Parameters: json.RawMessage(`{invalid json`),
				},
			},
		},
	}

	agent, confidence := DetectAgent(entry)

	// Should fall back to pattern detection
	if agent != "the-spec-writer" {
		t.Errorf("Expected fallback to pattern detection, got '%s'", agent)
	}
	if confidence != 0.95 {
		t.Errorf("Expected confidence 0.95 from pattern fallback, got %f", confidence)
	}
}

// Test nil parameters in Task tool
func TestDetectAgent_NilTaskParameters(t *testing.T) {
	entry := ClaudeLogEntry{
		Type:      "assistant",
		SessionID: "test-session",
		Timestamp: time.Now(),
		Assistant: &AssistantMessage{
			Text: "Using the-implementor agent",
			ToolUses: []ToolUse{
				{
					ID:         "tool_001",
					Name:       "Task",
					Parameters: nil,
				},
			},
		},
	}

	agent, confidence := DetectAgent(entry)

	// Should fall back to pattern detection
	if agent != "the-implementor" {
		t.Errorf("Expected fallback to pattern, got '%s'", agent)
	}
	if confidence != 0.90 {
		t.Errorf("Expected confidence 0.90, got %f", confidence)
	}
}

// Test empty entry returns no agent
func TestDetectAgent_EmptyEntry(t *testing.T) {
	entry := ClaudeLogEntry{}

	agent, confidence := DetectAgent(entry)

	if agent != "" {
		t.Errorf("Expected empty agent for empty entry, got '%s'", agent)
	}
	if confidence != 0.0 {
		t.Errorf("Expected confidence 0.0 for empty entry, got %f", confidence)
	}
}

// Test nil assistant message
func TestDetectAgent_NilAssistantMessage(t *testing.T) {
	entry := ClaudeLogEntry{
		Type:      "assistant",
		SessionID: "test-session",
		Timestamp: time.Now(),
		Assistant: nil,
	}

	agent, confidence := DetectAgent(entry)

	if agent != "" {
		t.Errorf("Expected empty agent for nil assistant, got '%s'", agent)
	}
	if confidence != 0.0 {
		t.Errorf("Expected confidence 0.0 for nil assistant, got %f", confidence)
	}
}

// Test system message with agent mention
func TestDetectAgent_SystemMessage(t *testing.T) {
	entry := ClaudeLogEntry{
		Type:      "system",
		SessionID: "test-session",
		Timestamp: time.Now(),
		System: &SystemMessage{
			Message: "delegating to the-reviewer for analysis",
			Text:    "System delegating task",
		},
	}

	agent, confidence := DetectAgent(entry)

	if agent != "the-reviewer" {
		t.Errorf("Expected 'the-reviewer' from system message, got '%s'", agent)
	}
	if confidence != 0.95 {
		t.Errorf("Expected confidence 0.95, got %f", confidence)
	}
}

// Test user message with agent mention
func TestDetectAgent_UserMessage(t *testing.T) {
	entry := ClaudeLogEntry{
		Type:      "user",
		SessionID: "test-session",
		Timestamp: time.Now(),
		User: &UserMessage{
			Text: "Please use the-test-writer agent",
		},
	}

	agent, confidence := DetectAgent(entry)

	if agent != "the-test-writer" {
		t.Errorf("Expected 'the-test-writer' from user message, got '%s'", agent)
	}
	if confidence != 0.90 {
		t.Errorf("Expected confidence 0.90, got %f", confidence)
	}
}

// Test ValidateLogEntry with valid entry
func TestValidateLogEntry_ValidEntry(t *testing.T) {
	entry := ClaudeLogEntry{
		Type:      "assistant",
		SessionID: "test-session",
		Timestamp: time.Now(),
		Assistant: &AssistantMessage{
			Text: "Valid message",
		},
	}

	err := ValidateLogEntry(entry)

	if err != nil {
		t.Errorf("Expected no error for valid entry, got %v", err)
	}
}

// Test ValidateLogEntry with missing type
func TestValidateLogEntry_MissingType(t *testing.T) {
	entry := ClaudeLogEntry{
		SessionID: "test-session",
		Timestamp: time.Now(),
	}

	err := ValidateLogEntry(entry)

	if err == nil {
		t.Error("Expected error for missing type")
	}

	if toolErr, ok := err.(*ToolError); ok {
		if toolErr.Code != "INVALID_ENTRY" {
			t.Errorf("Expected error code 'INVALID_ENTRY', got '%s'", toolErr.Code)
		}
	} else {
		t.Error("Expected ToolError type")
	}
}

// Test ValidateLogEntry with nil assistant content
func TestValidateLogEntry_NilAssistantContent(t *testing.T) {
	entry := ClaudeLogEntry{
		Type:      "assistant",
		SessionID: "test-session",
		Timestamp: time.Now(),
		Assistant: nil,
	}

	err := ValidateLogEntry(entry)

	if err == nil {
		t.Error("Expected error for nil assistant content")
	}

	if toolErr, ok := err.(*ToolError); ok {
		if toolErr.Code != "INVALID_ASSISTANT" {
			t.Errorf("Expected error code 'INVALID_ASSISTANT', got '%s'", toolErr.Code)
		}
	} else {
		t.Error("Expected ToolError type")
	}
}

// Test ValidateLogEntry with nil user content
func TestValidateLogEntry_NilUserContent(t *testing.T) {
	entry := ClaudeLogEntry{
		Type:      "user",
		SessionID: "test-session",
		Timestamp: time.Now(),
		User:      nil,
	}

	err := ValidateLogEntry(entry)

	if err == nil {
		t.Error("Expected error for nil user content")
	}

	if toolErr, ok := err.(*ToolError); ok {
		if toolErr.Code != "INVALID_USER" {
			t.Errorf("Expected error code 'INVALID_USER', got '%s'", toolErr.Code)
		}
	} else {
		t.Error("Expected ToolError type")
	}
}

// Test ValidateLogEntry with nil system content
func TestValidateLogEntry_NilSystemContent(t *testing.T) {
	entry := ClaudeLogEntry{
		Type:      "system",
		SessionID: "test-session",
		Timestamp: time.Now(),
		System:    nil,
	}

	err := ValidateLogEntry(entry)

	if err == nil {
		t.Error("Expected error for nil system content")
	}

	if toolErr, ok := err.(*ToolError); ok {
		if toolErr.Code != "INVALID_SYSTEM" {
			t.Errorf("Expected error code 'INVALID_SYSTEM', got '%s'", toolErr.Code)
		}
	} else {
		t.Error("Expected ToolError type")
	}
}

// Test extractSubagentType with valid JSON
func TestExtractSubagentType_ValidJSON(t *testing.T) {
	params := json.RawMessage(`{
		"instructions": "Test instructions",
		"subagent_type": "the-test-writer"
	}`)

	agent := extractSubagentType(params)

	if agent != "the-test-writer" {
		t.Errorf("Expected 'the-test-writer', got '%s'", agent)
	}
}

// Test extractSubagentType with missing field
func TestExtractSubagentType_MissingField(t *testing.T) {
	params := json.RawMessage(`{
		"instructions": "Test instructions"
	}`)

	agent := extractSubagentType(params)

	if agent != "" {
		t.Errorf("Expected empty string for missing field, got '%s'", agent)
	}
}

// Test extractSubagentType with empty string value
func TestExtractSubagentType_EmptyValue(t *testing.T) {
	params := json.RawMessage(`{
		"instructions": "Test instructions",
		"subagent_type": ""
	}`)

	agent := extractSubagentType(params)

	if agent != "" {
		t.Errorf("Expected empty string for empty value, got '%s'", agent)
	}
}

// Test extractSubagentType with malformed JSON
func TestExtractSubagentType_MalformedJSON(t *testing.T) {
	params := json.RawMessage(`{invalid json}`)

	agent := extractSubagentType(params)

	if agent != "" {
		t.Errorf("Expected empty string for malformed JSON, got '%s'", agent)
	}
}

// Test normalizeAgentName with various inputs
func TestNormalizeAgentName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"spec-writer", "the-spec-writer"},
		{"the-spec-writer", "the-spec-writer"},
		{"SPEC-WRITER", "the-spec-writer"},
		{"THE-SPEC-WRITER", "the-spec-writer"},
		{"  spec-writer  ", "the-spec-writer"},
		{"implementor", "the-implementor"},
		{"test-writer", "the-test-writer"},
		{"reviewer", "the-reviewer"},
		{"documenter", "the-documenter"},
		{"debugger", "the-debugger"},
		{"custom-agent", "custom-agent"}, // Unknown agent kept as-is
		{"", ""},
	}

	for _, tt := range tests {
		result := normalizeAgentName(tt.input)
		if result != tt.expected {
			t.Errorf("normalizeAgentName(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

// Test detectAgentWithDetails returns correct method information
func TestDetectAgentWithDetails_MethodInfo(t *testing.T) {
	tests := []struct {
		name           string
		entry          ClaudeLogEntry
		expectedMethod string
	}{
		{
			name: "Task tool method",
			entry: ClaudeLogEntry{
				Type: "assistant",
				Assistant: &AssistantMessage{
					ToolUses: []ToolUse{
						{
							Name: "Task",
							Parameters: json.RawMessage(`{
								"subagent_type": "the-spec-writer"
							}`),
						},
					},
				},
			},
			expectedMethod: "task_tool",
		},
		{
			name: "Pattern method",
			entry: ClaudeLogEntry{
				Type: "assistant",
				Assistant: &AssistantMessage{
					Text: "delegating to the-reviewer",
				},
			},
			expectedMethod: "pattern",
		},
		{
			name: "Context method",
			entry: ClaudeLogEntry{
				Type: "assistant",
				Assistant: &AssistantMessage{
					Text: "Running /s:test command",
				},
			},
			expectedMethod: "context",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			detection := detectAgentWithDetails(tt.entry)
			if detection.Method != tt.expectedMethod {
				t.Errorf("Expected method '%s', got '%s'", tt.expectedMethod, detection.Method)
			}
		})
	}
}

// Benchmark agent detection performance
func BenchmarkDetectAgent_TaskTool(b *testing.B) {
	entry := ClaudeLogEntry{
		Type: "assistant",
		Assistant: &AssistantMessage{
			ToolUses: []ToolUse{
				{
					Name: "Task",
					Parameters: json.RawMessage(`{
						"subagent_type": "the-spec-writer"
					}`),
				},
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DetectAgent(entry)
	}
}

func BenchmarkDetectAgent_Pattern(b *testing.B) {
	entry := ClaudeLogEntry{
		Type: "assistant",
		Assistant: &AssistantMessage{
			Text: "I'm delegating to the-reviewer for code quality checks",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DetectAgent(entry)
	}
}

func BenchmarkDetectAgent_Context(b *testing.B) {
	entry := ClaudeLogEntry{
		Type: "assistant",
		Assistant: &AssistantMessage{
			Text: "Running /s:implement to build this feature",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DetectAgent(entry)
	}
}