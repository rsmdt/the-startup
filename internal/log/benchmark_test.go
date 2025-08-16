package log

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"
)

// BenchmarkProcessToolCall benchmarks the main processing function
func BenchmarkProcessToolCall(b *testing.B) {
	testInput := `{
		"tool_name": "Task",
		"tool_input": {
			"subagent_type": "the-architect",
			"description": "Design system architecture",
			"prompt": "SessionId: dev-benchmark-123 AgentId: arch-001\nPlease design a scalable microservices architecture"
		}
	}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := strings.NewReader(testInput)
		_, err := ProcessToolCall(reader, false)
		if err != nil {
			b.Fatalf("Unexpected error: %v", err)
		}
	}
}

// BenchmarkProcessToolCallLarge benchmarks with large input
func BenchmarkProcessToolCallLarge(b *testing.B) {
	largeDescription := strings.Repeat("This is a very detailed description that contains a lot of text. ", 1000)
	largePrompt := strings.Repeat("This is a comprehensive prompt with extensive details about the task. ", 1000)

	testInput := fmt.Sprintf(`{
		"tool_name": "Task",
		"tool_input": {
			"subagent_type": "the-architect",
			"description": "%s",
			"prompt": "SessionId: dev-large-123 AgentId: arch-002\n%s"
		}
	}`, largeDescription, largePrompt)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := strings.NewReader(testInput)
		_, err := ProcessToolCall(reader, false)
		if err != nil {
			b.Fatalf("Unexpected error: %v", err)
		}
	}
}

// BenchmarkExtractSessionID benchmarks session ID extraction
func BenchmarkExtractSessionID(b *testing.B) {
	testPrompts := []string{
		"SessionId: dev-session-123\nSome content here",
		"A longer prompt with SessionId: dev-session-456 somewhere in the middle and more content after",
		"SessionId: dev-session-with-very-long-id-that-contains-many-characters-789\nContent",
		"No session ID in this prompt at all, just regular content",
		"Multiple SessionId: first-session and SessionId: second-session patterns",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		prompt := testPrompts[i%len(testPrompts)]
		_ = ExtractSessionID(prompt)
	}
}

// BenchmarkExtractAgentID benchmarks agent ID extraction
func BenchmarkExtractAgentID(b *testing.B) {
	testPrompts := []string{
		"AgentId: arch-001\nContent",
		"Longer prompt with AgentId: dev-002 in middle",
		"AgentId: agent-with-long-identifier-name-123\nContent",
		"No agent ID here",
		"Multiple AgentId: first-agent and AgentId: second-agent",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		prompt := testPrompts[i%len(testPrompts)]
		_ = ExtractAgentID(prompt)
	}
}

// BenchmarkShouldProcess benchmarks filtering logic
func BenchmarkShouldProcess(b *testing.B) {
	testCases := []struct {
		toolName  string
		toolInput map[string]interface{}
	}{
		{"Task", map[string]interface{}{"subagent_type": "the-architect"}},
		{"Task", map[string]interface{}{"subagent_type": "the-developer"}},
		{"Task", map[string]interface{}{"subagent_type": "architect"}}, // Should be filtered
		{"Other", map[string]interface{}{"subagent_type": "the-test"}}, // Should be filtered
		{"Task", map[string]interface{}{"other_field": "value"}},       // Should be filtered
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tc := testCases[i%len(testCases)]
		_ = ShouldProcess(tc.toolName, tc.toolInput)
	}
}

// BenchmarkTruncateOutput benchmarks output truncation
func BenchmarkTruncateOutput(b *testing.B) {
	// Create test outputs of various sizes
	smallOutput := "Short output"
	mediumOutput := strings.Repeat("Medium length output. ", 50)
	largeOutput := strings.Repeat("This is a large output that will be truncated. ", 100)

	testOutputs := []string{smallOutput, mediumOutput, largeOutput}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		output := testOutputs[i%len(testOutputs)]
		_ = TruncateOutput(output, 1000)
	}
}

// BenchmarkTimestamp benchmarks timestamp generation
func BenchmarkTimestamp(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	}
}

// BenchmarkRegexPerformance benchmarks regex performance with various inputs
func BenchmarkRegexPerformance(b *testing.B) {
	// Test with different prompt lengths and patterns
	testPrompts := []string{
		"SessionId: simple",
		strings.Repeat("No match here. ", 100) + "SessionId: found-it" + strings.Repeat(" More content. ", 100),
		"AgentId: agent-123 SessionId: session-456",
		strings.Repeat("Very long prompt without any matching patterns. ", 500),
		"Multiple SessionId: first SessionId: second SessionId: third patterns",
	}

	b.Run("SessionID", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			prompt := testPrompts[i%len(testPrompts)]
			_ = ExtractSessionID(prompt)
		}
	})

	b.Run("AgentID", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			prompt := testPrompts[i%len(testPrompts)]
			_ = ExtractAgentID(prompt)
		}
	})
}

// BenchmarkMemoryUsage benchmarks memory allocation patterns
func BenchmarkMemoryUsage(b *testing.B) {
	testInput := `{
		"tool_name": "Task",
		"tool_input": {
			"subagent_type": "the-memory-test",
			"description": "Memory usage test",
			"prompt": "SessionId: dev-memory-123 AgentId: mem-001\nTest memory usage patterns"
		}
	}`

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		reader := strings.NewReader(testInput)
		result, err := ProcessToolCall(reader, false)
		if err != nil {
			b.Fatalf("Unexpected error: %v", err)
		}
		// Prevent compiler optimizations
		_ = result
	}
}

// BenchmarkConcurrentProcessing benchmarks concurrent hook processing
func BenchmarkConcurrentProcessing(b *testing.B) {
	testInput := `{
		"tool_name": "Task",
		"tool_input": {
			"subagent_type": "the-concurrent",
			"description": "Concurrent processing test",
			"prompt": "SessionId: dev-concurrent-123 AgentId: conc-001\nTest concurrent processing"
		}
	}`

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			reader := strings.NewReader(testInput)
			_, err := ProcessToolCall(reader, false)
			if err != nil {
				b.Errorf("Unexpected error: %v", err)
			}
		}
	})
}

// BenchmarkJSONMarshalUnmarshal benchmarks JSON operations
func BenchmarkJSONMarshalUnmarshal(b *testing.B) {
	testInput := `{
		"tool_name": "Task",
		"tool_input": {
			"subagent_type": "the-json-test",
			"description": "JSON performance test",
			"prompt": "SessionId: dev-json-123 AgentId: json-001\nTest JSON performance"
		},
		"session_id": "dev-json-123",
		"transcript_path": "/path/to/transcript.jsonl",
		"cwd": "/current/working/directory",
		"hook_event_name": "PreToolUse"
	}`

	b.Run("Unmarshal", func(b *testing.B) {
		inputBytes := []byte(testInput)
		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			var hookInput HookInput
			err := json.Unmarshal(inputBytes, &hookInput)
			if err != nil {
				b.Fatalf("Unmarshal error: %v", err)
			}
			_ = hookInput
		}
	})

	b.Run("Marshal", func(b *testing.B) {
		hookData := &HookData{
			Event:       "agent_start",
			AgentType:   "the-json-test",
			AgentID:     "json-001",
			Description: "JSON performance test",
			SessionID:   "dev-json-123",
			Timestamp:   "2025-01-11T12:00:00.000Z",
			Instruction: "Test JSON performance",
		}

		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			_, err := json.Marshal(hookData)
			if err != nil {
				b.Fatalf("Marshal error: %v", err)
			}
		}
	})
}
