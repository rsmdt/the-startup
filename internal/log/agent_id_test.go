package log

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestExtractAgentIDFromPrompt tests basic agent ID extraction
func TestExtractAgentIDFromPrompt(t *testing.T) {
	tests := []struct {
		name     string
		prompt   string
		expected string
	}{
		{
			name:     "valid agent ID",
			prompt:   "AgentId: test-agent\nSome content",
			expected: "test-agent",
		},
		{
			name:     "case insensitive",
			prompt:   "agentid: test-agent\nSome content",
			expected: "test-agent",
		},
		{
			name:     "nested agent ID",
			prompt:   "AgentId: the-architect/system-design\nSome content",
			expected: "the-architect/system-design",
		},
		{
			name:     "nested with underscores and hyphens",
			prompt:   "AgentId: the-backend_engineer/api-design\nSome content",
			expected: "the-backend_engineer/api-design",
		},
		{
			name:     "no agent ID",
			prompt:   "Some content without agent ID",
			expected: "",
		},
		{
			name:     "invalid characters",
			prompt:   "AgentId: test@agent\nSome content",
			expected: "", // @ is not allowed
		},
		{
			name:     "too many nesting levels",
			prompt:   "AgentId: the-architect/system/design\nSome content",
			expected: "", // Max depth is 2
		},
		{
			name:     "invalid nested - slash at start",
			prompt:   "AgentId: /system-design\nSome content",
			expected: "", // Must not start with slash
		},
		{
			name:     "invalid nested - slash at end",
			prompt:   "AgentId: the-architect/\nSome content",
			expected: "", // Must not end with slash
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractAgentIDFromPrompt(tt.prompt)
			if result != tt.expected {
				t.Errorf("ExtractAgentIDFromPrompt(%q) = %q, expected %q", tt.prompt, result, tt.expected)
			}
		})
	}
}

// TestGenerateAgentID tests basic ID generation
func TestGenerateAgentID(t *testing.T) {
	tests := []struct {
		name         string
		agentType    string
		expectedEnd  string
		description  string
	}{
		{
			name:        "flat agent type",
			agentType:   "the-architect",
			expectedEnd: "-the-architect",
			description: "Should generate ID with flat agent type",
		},
		{
			name:        "nested agent type",
			agentType:   "the-architect/system-design",
			expectedEnd: "-the-architect-system-design", // Slashes replaced with hyphens for filesystem
			description: "Should generate ID with nested agent type (slashes to hyphens)",
		},
		{
			name:        "nested with underscores",
			agentType:   "the-backend_engineer/api-design",
			expectedEnd: "-the-backend_engineer-api-design",
			description: "Should handle underscores in nested type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := GenerateAgentID("session-123", tt.agentType, "some prompt")

			// Should be format: {8-char-id}-{agent-type}
			if !strings.HasSuffix(id, tt.expectedEnd) {
				t.Errorf("%s: generated ID should end with %s, got %s", tt.description, tt.expectedEnd, id)
			}

			// Should have an 8-character prefix
			parts := strings.SplitN(id, "-", 2)
			if len(parts[0]) != 8 {
				t.Errorf("ID prefix should be 8 characters, got %d: %s", len(parts[0]), parts[0])
			}
		})
	}
}

// TestEnsureUniqueness tests that file collisions are handled
func TestEnsureUniqueness(t *testing.T) {
	tempDir := t.TempDir()
	sessionID := "test-session"
	baseID := "test-agent"

	// Create session directory
	sessionDir := filepath.Join(tempDir, sessionID)
	os.MkdirAll(sessionDir, 0755)

	// First call should return base ID
	result1 := ensureUniqueness(baseID, sessionID, tempDir)
	if result1 != baseID {
		t.Errorf("First call should return base ID: %s", result1)
	}

	// Create a file with that name to simulate collision
	os.WriteFile(filepath.Join(sessionDir, baseID+".jsonl"), []byte("test"), 0644)

	// Second call should return disambiguated ID
	result2 := ensureUniqueness(baseID, sessionID, tempDir)
	if result2 == baseID {
		t.Errorf("Second call should return different ID due to collision")
	}
	if !strings.HasPrefix(result2, baseID+"-") {
		t.Errorf("Disambiguated ID should start with base ID: %s", result2)
	}
}

// TestExtractOrGenerateAgentID tests the main function that's actually used
func TestExtractOrGenerateAgentID(t *testing.T) {
	tests := []struct {
		name      string
		prompt    string
		agentType string
		checkFunc func(string) bool
	}{
		{
			name:      "extracts when present",
			prompt:    "AgentId: explicit-id\nContent",
			agentType: "the-architect",
			checkFunc: func(id string) bool { return id == "explicit-id" },
		},
		{
			name:      "generates when absent",
			prompt:    "No agent ID here",
			agentType: "the-developer",
			checkFunc: func(id string) bool { return strings.HasSuffix(id, "-the-developer") },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractOrGenerateAgentID(tt.prompt, tt.agentType, "session-123")
			if !tt.checkFunc(result) {
				t.Errorf("ExtractOrGenerateAgentID failed check for %s, got: %s", tt.name, result)
			}
		})
	}
}
