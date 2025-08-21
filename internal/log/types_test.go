package log

import (
	"encoding/json"
	"testing"
)

func TestHookInputMarshaling(t *testing.T) {
	// Test data based on SDD documentation
	input := HookInput{
		SessionID:      "abc123",
		TranscriptPath: "/path/to/transcript.jsonl",
		Cwd:            "/current/working/directory",
		HookEventName:  "PreToolUse",
		ToolName:       "Task",
		ToolInput: map[string]interface{}{
			"description":   "Short task description",
			"prompt":        "Full prompt for the agent",
			"subagent_type": "the-architect",
		},
	}

	// Marshal to JSON
	jsonBytes, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("Failed to marshal HookInput: %v", err)
	}

	// Unmarshal back
	var unmarshaled HookInput
	if err := json.Unmarshal(jsonBytes, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal HookInput: %v", err)
	}

	// Verify fields
	if unmarshaled.SessionID != input.SessionID {
		t.Errorf("SessionID mismatch: expected %s, got %s", input.SessionID, unmarshaled.SessionID)
	}
	if unmarshaled.ToolName != input.ToolName {
		t.Errorf("ToolName mismatch: expected %s, got %s", input.ToolName, unmarshaled.ToolName)
	}
	if unmarshaled.HookEventName != input.HookEventName {
		t.Errorf("HookEventName mismatch: expected %s, got %s", input.HookEventName, unmarshaled.HookEventName)
	}

	// Verify tool_input fields
	subagentType, ok := unmarshaled.ToolInput["subagent_type"].(string)
	if !ok || subagentType != "the-architect" {
		t.Errorf("ToolInput subagent_type mismatch: expected the-architect, got %v", subagentType)
	}
}

func TestHookInputWithOutput(t *testing.T) {
	// Test PostToolUse case with output
	input := HookInput{
		ToolName:      "Task",
		HookEventName: "PostToolUse",
		ToolInput: map[string]interface{}{
			"subagent_type": "the-developer",
		},
		Output: "Agent response here",
	}

	jsonBytes, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("Failed to marshal HookInput: %v", err)
	}

	var unmarshaled HookInput
	if err := json.Unmarshal(jsonBytes, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal HookInput: %v", err)
	}

	if unmarshaled.Output != input.Output {
		t.Errorf("Output mismatch: expected %s, got %s", input.Output, unmarshaled.Output)
	}
}

func TestHookDataMarshaling(t *testing.T) {
	// Test user role (PreToolUse)
	userData := HookData{
		Role:      "user",
		Content:   "Please design a scalable architecture",
		Timestamp: "2025-01-11T10:30:00.000Z",
		SessionID: "dev-session-123", // Internal field, not serialized
		AgentID:   "arch-001",        // Internal field, not serialized
	}

	jsonBytes, err := json.Marshal(userData)
	if err != nil {
		t.Fatalf("Failed to marshal HookData: %v", err)
	}

	// Verify JSON only contains serialized fields
	var jsonMap map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &jsonMap); err != nil {
		t.Fatalf("Failed to unmarshal to map: %v", err)
	}

	// Check that only role, content, and timestamp are in JSON
	if role, exists := jsonMap["role"]; !exists || role != "user" {
		t.Errorf("Expected role 'user' in JSON, got %v", role)
	}
	if content, exists := jsonMap["content"]; !exists || content != userData.Content {
		t.Errorf("Expected content in JSON, got %v", content)
	}
	if timestamp, exists := jsonMap["timestamp"]; !exists || timestamp != userData.Timestamp {
		t.Errorf("Expected timestamp in JSON, got %v", timestamp)
	}

	// SessionID and AgentID should NOT be in JSON (marked with json:"-")
	if _, exists := jsonMap["session_id"]; exists {
		t.Error("SessionID should not be serialized to JSON")
	}
	if _, exists := jsonMap["agent_id"]; exists {
		t.Error("AgentID should not be serialized to JSON")
	}
}

func TestHookDataAssistantRole(t *testing.T) {
	// Test assistant role (PostToolUse)
	assistantData := HookData{
		Role:      "assistant",
		Content:   "Successfully implemented the feature with tests",
		Timestamp: "2025-01-11T10:35:00.000Z",
		SessionID: "dev-session-123",
		AgentID:   "dev-001",
	}

	jsonBytes, err := json.Marshal(assistantData)
	if err != nil {
		t.Fatalf("Failed to marshal HookData: %v", err)
	}

	// Unmarshal and verify
	var unmarshaled HookData
	if err := json.Unmarshal(jsonBytes, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal HookData: %v", err)
	}

	if unmarshaled.Role != assistantData.Role {
		t.Errorf("Role mismatch: expected %s, got %s", assistantData.Role, unmarshaled.Role)
	}
	if unmarshaled.Content != assistantData.Content {
		t.Errorf("Content mismatch: expected %s, got %s", assistantData.Content, unmarshaled.Content)
	}
	if unmarshaled.Timestamp != assistantData.Timestamp {
		t.Errorf("Timestamp mismatch: expected %s, got %s", assistantData.Timestamp, unmarshaled.Timestamp)
	}
}

func TestHookDataInternalFields(t *testing.T) {
	// Test that internal fields are preserved in struct but not serialized
	data := HookData{
		Role:      "user",
		Content:   "Test content",
		Timestamp: "2025-01-11T10:40:00.000Z",
		SessionID: "internal-session",
		AgentID:   "internal-agent",
	}

	// Internal fields should be accessible on the struct
	if data.SessionID != "internal-session" {
		t.Errorf("SessionID should be accessible on struct")
	}
	if data.AgentID != "internal-agent" {
		t.Errorf("AgentID should be accessible on struct")
	}

	// But not serialized to JSON
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Failed to marshal HookData: %v", err)
	}

	jsonStr := string(jsonBytes)
	if contains := "internal-session"; contains != "" && len(jsonStr) > 0 {
		// Check that internal fields are not in JSON
		var jsonMap map[string]interface{}
		json.Unmarshal(jsonBytes, &jsonMap)
		if _, exists := jsonMap["SessionID"]; exists {
			t.Error("SessionID should not appear in JSON")
		}
		if _, exists := jsonMap["AgentID"]; exists {
			t.Error("AgentID should not appear in JSON")
		}
	}
}
