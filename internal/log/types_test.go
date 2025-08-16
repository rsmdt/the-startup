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
	// Test agent_start event
	startData := HookData{
		Event:       "agent_start",
		AgentType:   "the-architect",
		AgentID:     "arch-001",
		Description: "Design the system architecture",
		Instruction: "Please design a scalable architecture",
		SessionID:   "dev-session-123",
		Timestamp:   "2025-01-11T10:30:00.000Z",
	}

	jsonBytes, err := json.Marshal(startData)
	if err != nil {
		t.Fatalf("Failed to marshal HookData: %v", err)
	}

	var unmarshaled HookData
	if err := json.Unmarshal(jsonBytes, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal HookData: %v", err)
	}

	// Verify all fields
	if unmarshaled.Event != startData.Event {
		t.Errorf("Event mismatch: expected %s, got %s", startData.Event, unmarshaled.Event)
	}
	if unmarshaled.AgentType != startData.AgentType {
		t.Errorf("AgentType mismatch: expected %s, got %s", startData.AgentType, unmarshaled.AgentType)
	}
	if unmarshaled.Instruction != startData.Instruction {
		t.Errorf("Instruction mismatch: expected %s, got %s", startData.Instruction, unmarshaled.Instruction)
	}
	if unmarshaled.SessionID != startData.SessionID {
		t.Errorf("SessionID mismatch: expected %s, got %s", startData.SessionID, unmarshaled.SessionID)
	}
}

func TestHookDataAgentComplete(t *testing.T) {
	// Test agent_complete event
	completeData := HookData{
		Event:         "agent_complete",
		AgentType:     "the-developer",
		AgentID:       "dev-001",
		Description:   "Implement the feature",
		OutputSummary: "Successfully implemented the feature with tests",
		SessionID:     "dev-session-123",
		Timestamp:     "2025-01-11T10:35:00.000Z",
	}

	jsonBytes, err := json.Marshal(completeData)
	if err != nil {
		t.Fatalf("Failed to marshal HookData: %v", err)
	}

	// Verify it contains the expected JSON fields
	var jsonMap map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &jsonMap); err != nil {
		t.Fatalf("Failed to unmarshal to map: %v", err)
	}

	// Should have output_summary but not instruction for complete events
	if _, exists := jsonMap["output_summary"]; !exists {
		t.Error("Expected output_summary field in agent_complete JSON")
	}

	// Instruction should be empty and omitted from JSON
	if instruction, exists := jsonMap["instruction"]; exists && instruction != "" {
		t.Error("Instruction should be omitted in agent_complete JSON")
	}
}

func TestJSONOmitEmpty(t *testing.T) {
	// Test that omitempty works correctly
	data := HookData{
		Event:     "agent_start",
		AgentType: "the-tester",
		Timestamp: "2025-01-11T10:40:00.000Z",
		// OutputSummary intentionally empty - should be omitted
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Failed to marshal HookData: %v", err)
	}

	var jsonMap map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &jsonMap); err != nil {
		t.Fatalf("Failed to unmarshal to map: %v", err)
	}

	// OutputSummary should be omitted when empty
	if _, exists := jsonMap["output_summary"]; exists {
		t.Error("Empty output_summary should be omitted from JSON")
	}

	// But Event should always be present
	if _, exists := jsonMap["event"]; !exists {
		t.Error("Event field should always be present in JSON")
	}
}
