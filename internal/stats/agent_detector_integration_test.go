package stats

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

// Integration test simulating real log processing with agent detection
func TestAgentDetection_IntegrationScenario(t *testing.T) {
	// Simulate a series of log entries with various agent interactions
	logEntries := []struct {
		name             string
		jsonLine         string
		expectedAgent    string
		expectedConfidence float64
		shouldValidate   bool
	}{
		{
			name: "Task tool with subagent_type",
			jsonLine: `{
				"type": "assistant",
				"session_id": "sess_001",
				"timestamp": "2025-01-15T10:00:00Z",
				"content": {
					"text": "I'll delegate this task",
					"tool_uses": [{
						"id": "tool_001",
						"name": "Task",
						"parameters": {
							"instructions": "Write comprehensive tests",
							"subagent_type": "the-test-writer"
						}
					}]
				}
			}`,
			expectedAgent:      "the-test-writer",
			expectedConfidence: 1.0,
			shouldValidate:     true,
		},
		{
			name: "Pattern matching - delegating",
			jsonLine: `{
				"type": "assistant",
				"session_id": "sess_001",
				"timestamp": "2025-01-15T10:01:00Z",
				"content": {
					"text": "I'm delegating to the-spec-writer for the specification"
				}
			}`,
			expectedAgent:      "the-spec-writer",
			expectedConfidence: 0.95,
			shouldValidate:     true,
		},
		{
			name: "Command context inference",
			jsonLine: `{
				"type": "assistant",
				"session_id": "sess_001",
				"timestamp": "2025-01-15T10:02:00Z",
				"content": {
					"text": "Let me run /s:implement to build this feature"
				}
			}`,
			expectedAgent:      "the-implementor",
			expectedConfidence: 0.70,
			shouldValidate:     true,
		},
		{
			name: "No agent detected",
			jsonLine: `{
				"type": "assistant",
				"session_id": "sess_001",
				"timestamp": "2025-01-15T10:03:00Z",
				"content": {
					"text": "Let me analyze this code for you"
				}
			}`,
			expectedAgent:      "",
			expectedConfidence: 0.0,
			shouldValidate:     true,
		},
		{
			name: "Corrupted JSON - should skip",
			jsonLine: `{
				"type": "assistant",
				invalid json here
			}`,
			expectedAgent:      "",
			expectedConfidence: 0.0,
			shouldValidate:     false,
		},
		{
			name: "Task result in user message",
			jsonLine: `{
				"type": "user",
				"session_id": "sess_001",
				"timestamp": "2025-01-15T10:04:00Z",
				"content": {
					"tool_use_id": "tool_001",
					"tool_name": "Task",
					"result": {
						"subagent_type": "the-reviewer",
						"status": "completed"
					}
				}
			}`,
			expectedAgent:      "the-reviewer",
			expectedConfidence: 1.0,
			shouldValidate:     true,
		},
		{
			name: "Missing type field - validation error",
			jsonLine: `{
				"session_id": "sess_001",
				"timestamp": "2025-01-15T10:05:00Z",
				"content": {}
			}`,
			expectedAgent:      "",
			expectedConfidence: 0.0,
			shouldValidate:     true,
		},
	}

	for _, tc := range logEntries {
		t.Run(tc.name, func(t *testing.T) {
			// Parse the JSON line
			var entry ClaudeLogEntry
			err := json.Unmarshal([]byte(tc.jsonLine), &entry)

			if !tc.shouldValidate && err != nil {
				// Expected parse error for corrupted JSON
				t.Logf("Expected parse error: %v", err)
				return
			}

			if err != nil {
				t.Fatalf("Failed to parse JSON: %v", err)
			}

			// Parse content based on type
			if entry.Content != nil {
				switch entry.Type {
				case "assistant":
					var content AssistantMessage
					if err := json.Unmarshal(entry.Content, &content); err == nil {
						entry.Assistant = &content
					}
				case "user":
					var content UserMessage
					if err := json.Unmarshal(entry.Content, &content); err == nil {
						entry.User = &content
					}
				}
			}

			// Validate the entry
			validationErr := ValidateLogEntry(entry)
			if tc.name == "Missing type field - validation error" {
				if validationErr == nil {
					t.Error("Expected validation error for missing type")
				}
				return
			}

			// Detect agent
			agent, confidence := DetectAgent(entry)

			// Check results
			if agent != tc.expectedAgent {
				t.Errorf("Expected agent '%s', got '%s'", tc.expectedAgent, agent)
			}
			if confidence != tc.expectedConfidence {
				t.Errorf("Expected confidence %f, got %f", tc.expectedConfidence, confidence)
			}
		})
	}
}

// Test handling of malformed and edge case data
func TestAgentDetection_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		entry       ClaudeLogEntry
		expectAgent string
		expectConf  float64
	}{
		{
			name:        "Completely empty entry",
			entry:       ClaudeLogEntry{},
			expectAgent: "",
			expectConf:  0.0,
		},
		{
			name: "Entry with only type",
			entry: ClaudeLogEntry{
				Type: "assistant",
			},
			expectAgent: "",
			expectConf:  0.0,
		},
		{
			name: "Tool use with empty parameters",
			entry: ClaudeLogEntry{
				Type: "assistant",
				Assistant: &AssistantMessage{
					ToolUses: []ToolUse{
						{
							Name:       "Task",
							Parameters: json.RawMessage(`{}`),
						},
					},
				},
			},
			expectAgent: "",
			expectConf:  0.0,
		},
		{
			name: "Multiple patterns in same text - first wins",
			entry: ClaudeLogEntry{
				Type: "assistant",
				Assistant: &AssistantMessage{
					Text: "delegating to the-spec-writer and using the-implementor agent",
				},
			},
			expectAgent: "the-spec-writer",
			expectConf:  0.95, // delegating pattern has higher priority
		},
		{
			name: "Case insensitive pattern matching",
			entry: ClaudeLogEntry{
				Type: "assistant",
				Assistant: &AssistantMessage{
					Text: "DELEGATING TO THE-REVIEWER",
				},
			},
			expectAgent: "the-reviewer",
			expectConf:  0.95,
		},
		{
			name: "Agent name with hyphens and underscores",
			entry: ClaudeLogEntry{
				Type: "assistant",
				Assistant: &AssistantMessage{
					Text: "delegating to custom-agent_v2",
				},
			},
			expectAgent: "custom-agent_v2",
			expectConf:  0.95,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agent, confidence := DetectAgent(tt.entry)

			if agent != tt.expectAgent {
				t.Errorf("Expected agent '%s', got '%s'", tt.expectAgent, agent)
			}
			if confidence != tt.expectConf {
				t.Errorf("Expected confidence %f, got %f", tt.expectConf, confidence)
			}
		})
	}
}

// Test performance with realistic log data
func TestAgentDetection_Performance(t *testing.T) {
	// Create a large batch of entries
	entries := make([]ClaudeLogEntry, 1000)

	for i := 0; i < 1000; i++ {
		switch i % 4 {
		case 0: // Task tool entry
			entries[i] = ClaudeLogEntry{
				Type: "assistant",
				Assistant: &AssistantMessage{
					ToolUses: []ToolUse{
						{
							Name: "Task",
							Parameters: json.RawMessage(`{"subagent_type": "the-test-writer"}`),
						},
					},
				},
			}
		case 1: // Pattern matching entry
			entries[i] = ClaudeLogEntry{
				Type: "assistant",
				Assistant: &AssistantMessage{
					Text: "delegating to the-spec-writer for specifications",
				},
			}
		case 2: // Command context entry
			entries[i] = ClaudeLogEntry{
				Type: "assistant",
				Assistant: &AssistantMessage{
					Text: "Running /s:implement command",
				},
			}
		case 3: // No agent entry
			entries[i] = ClaudeLogEntry{
				Type: "assistant",
				Assistant: &AssistantMessage{
					Text: "Processing your request",
				},
			}
		}
	}

	// Process all entries and measure time
	start := time.Now()
	detectionCount := 0

	for _, entry := range entries {
		agent, _ := DetectAgent(entry)
		if agent != "" {
			detectionCount++
		}
	}

	elapsed := time.Since(start)

	// Should process 1000 entries quickly
	if elapsed > 100*time.Millisecond {
		t.Errorf("Processing took too long: %v", elapsed)
	}

	// Should detect agents in 75% of entries (3 out of 4 have agents)
	expectedDetections := 750
	if detectionCount < expectedDetections-10 || detectionCount > expectedDetections+10 {
		t.Errorf("Expected around %d detections, got %d", expectedDetections, detectionCount)
	}

	t.Logf("Processed %d entries in %v, detected %d agents", len(entries), elapsed, detectionCount)
}

// Test concurrent detection safety
func TestAgentDetection_Concurrent(t *testing.T) {
	entry := ClaudeLogEntry{
		Type: "assistant",
		Assistant: &AssistantMessage{
			Text: "delegating to the-reviewer",
		},
	}

	// Run detection concurrently
	results := make(chan struct {
		agent      string
		confidence float64
	}, 100)

	for i := 0; i < 100; i++ {
		go func() {
			agent, confidence := DetectAgent(entry)
			results <- struct {
				agent      string
				confidence float64
			}{agent, confidence}
		}()
	}

	// Collect results
	for i := 0; i < 100; i++ {
		result := <-results
		if result.agent != "the-reviewer" {
			t.Errorf("Concurrent detection failed: expected 'the-reviewer', got '%s'", result.agent)
		}
		if result.confidence != 0.95 {
			t.Errorf("Concurrent detection failed: expected confidence 0.95, got %f", result.confidence)
		}
	}
}

// Test that detection works with real JSONL parsing
func TestAgentDetection_WithJSONLParser(t *testing.T) {
	jsonlData := strings.Join([]string{
		`{"type":"assistant","sessionId":"s1","timestamp":"2025-01-15T10:00:00Z","message":{"role":"assistant","content":[{"type":"text","text":"Let me delegate to the-spec-writer"},{"type":"tool_use","id":"t1","name":"Task","input":{"subagent_type":"the-spec-writer"}}]}}`,
		`{"type":"user","sessionId":"s1","timestamp":"2025-01-15T10:00:05Z","message":{"role":"user","content":[{"type":"tool_result","tool_use_id":"t1","tool_name":"Task","output":"Task completed","result":{"subagent_type":"the-spec-writer","status":"completed"}}]}}`,
		`{"type":"assistant","sessionId":"s1","timestamp":"2025-01-15T10:00:10Z","message":{"role":"assistant","content":[{"type":"text","text":"Now using the-implementor agent"}]}}`,
		`{"type":"assistant","sessionId":"s1","timestamp":"2025-01-15T10:00:15Z","message":{"role":"assistant","content":[{"type":"text","text":"Running /s:test to verify"}]}}`,
	}, "\n")

	parser := NewJSONLParser()
	reader := strings.NewReader(jsonlData)

	entries, errors := parser.ParseStream(reader)

	// Consume any errors
	go func() {
		for err := range errors {
			t.Logf("Parser error (expected for this test): %v", err)
		}
	}()

	agentDetections := []struct {
		expectedAgent string
		expectedConf  float64
	}{
		{"the-spec-writer", 1.0},  // Task tool
		{"the-spec-writer", 1.0},  // Task result
		{"the-implementor", 0.90}, // Pattern
		{"the-test-writer", 0.70}, // Command context
	}

	i := 0
	for entry := range entries {
		if i >= len(agentDetections) {
			break
		}

		// The parser populates the specific message types
		// so we should be able to detect agents directly
		agent, confidence := DetectAgent(entry)
		expected := agentDetections[i]

		// Log the entry details for debugging
		t.Logf("Entry %d type: %s, has assistant: %v, has user: %v",
			i, entry.Type, entry.Assistant != nil, entry.User != nil)

		if agent != expected.expectedAgent {
			t.Errorf("Entry %d: expected agent '%s', got '%s'", i, expected.expectedAgent, agent)
		}
		if confidence != expected.expectedConf {
			t.Errorf("Entry %d: expected confidence %f, got %f", i, expected.expectedConf, confidence)
		}

		i++
	}

	if i != len(agentDetections) {
		t.Errorf("Expected %d entries, but only processed %d", len(agentDetections), i)
	}
}