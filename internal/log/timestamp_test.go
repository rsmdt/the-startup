package log

import (
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestTimestampFormatValidation(t *testing.T) {
	// Test that timestamp format matches Python's datetime.utcnow().isoformat() + 'Z'
	// Python format: 2025-01-11T10:30:00.123456Z (with microseconds)
	// Go format should be: 2025-01-11T10:30:00.123Z (with milliseconds, but close enough)

	reader := strings.NewReader(`{"tool_name":"Task","tool_input":{"subagent_type":"the-test","description":"test","prompt":"SessionId: test-123\nTest"},"session_id":"test-123","hook_event_name":"PreToolUse"}`)
	
	hookData, err := ProcessToolCall(reader, false)
	if err != nil {
		t.Fatalf("ProcessToolCall failed: %v", err)
	}
	
	if hookData == nil {
		t.Fatal("Expected non-nil hookData")
	}
	
	// Check timestamp format: YYYY-MM-DDTHH:MM:SS.sssZ
	timestampRegex := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}Z$`)
	if !timestampRegex.MatchString(hookData.Timestamp) {
		t.Errorf("Timestamp format doesn't match Python format. Got: %s", hookData.Timestamp)
	}
	
	// Verify it's a valid time
	parsedTime, err := time.Parse("2006-01-02T15:04:05.000Z", hookData.Timestamp)
	if err != nil {
		t.Errorf("Failed to parse timestamp %s: %v", hookData.Timestamp, err)
	}
	
	// Verify it's recent (within last minute)
	now := time.Now().UTC()
	if now.Sub(parsedTime) > time.Minute {
		t.Errorf("Timestamp %s is not recent enough (more than 1 minute old)", hookData.Timestamp)
	}
	
	// Verify it's UTC (ends with Z)
	if !strings.HasSuffix(hookData.Timestamp, "Z") {
		t.Errorf("Timestamp should end with 'Z' to indicate UTC: %s", hookData.Timestamp)
	}
}

func TestTimestampConsistencyCompatibility(t *testing.T) {
	// Test that multiple calls produce timestamps in the correct format
	for i := 0; i < 5; i++ {
		reader := strings.NewReader(`{"tool_name":"Task","tool_input":{"subagent_type":"the-test","description":"test","prompt":"SessionId: test-123\nTest"},"session_id":"test-123","hook_event_name":"PreToolUse"}`)
		
		hookData, err := ProcessToolCall(reader, false)
		if err != nil {
			t.Fatalf("ProcessToolCall failed on iteration %d: %v", i, err)
		}
		
		if hookData == nil {
			t.Fatalf("Expected non-nil hookData on iteration %d", i)
		}
		
		// Check format consistency
		timestampRegex := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}Z$`)
		if !timestampRegex.MatchString(hookData.Timestamp) {
			t.Errorf("Iteration %d: Timestamp format inconsistent. Got: %s", i, hookData.Timestamp)
		}
	}
}