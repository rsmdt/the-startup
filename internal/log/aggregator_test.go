package log

import (
	"testing"
	"time"
)

func TestAggregateMetrics(t *testing.T) {
	// Create test entries
	entries := make(chan MetricsEntry, 10)
	
	// Add some test data
	go func() {
		defer close(entries)
		
		// Tool 1: Successful Edit with duration
		entries <- MetricsEntry{
			ToolID:    "edit_20250903T120000Z_session1",
			ToolName:  "Edit",
			HookEvent: "PreToolUse",
			Timestamp: time.Date(2025, 9, 3, 12, 0, 0, 0, time.UTC),
			SessionID: "session1",
		}
		entries <- MetricsEntry{
			ToolID:    "edit_20250903T120000Z_session1",
			ToolName:  "Edit",
			HookEvent: "PostToolUse",
			Timestamp: time.Date(2025, 9, 3, 12, 0, 1, 0, time.UTC),
			SessionID: "session1",
			Success:   boolPtr(true),
		}
		
		// Tool 2: Failed Bash command
		entries <- MetricsEntry{
			ToolID:    "bash_20250903T120100Z_session1",
			ToolName:  "Bash",
			HookEvent: "PreToolUse",
			Timestamp: time.Date(2025, 9, 3, 12, 1, 0, 0, time.UTC),
			SessionID: "session1",
		}
		entries <- MetricsEntry{
			ToolID:    "bash_20250903T120100Z_session1",
			ToolName:  "Bash",
			HookEvent: "PostToolUse",
			Timestamp: time.Date(2025, 9, 3, 12, 1, 2, 0, time.UTC),
			SessionID: "session1",
			Success:   boolPtr(false),
			Error:     "command failed",
			ErrorType: "exit_code_1",
		}
		
		// Tool 3: Another session
		entries <- MetricsEntry{
			ToolID:    "read_20250903T130000Z_session2",
			ToolName:  "Read",
			HookEvent: "PreToolUse",
			Timestamp: time.Date(2025, 9, 3, 13, 0, 0, 0, time.UTC),
			SessionID: "session2",
		}
		entries <- MetricsEntry{
			ToolID:    "read_20250903T130000Z_session2",
			ToolName:  "Read",
			HookEvent: "PostToolUse",
			Timestamp: time.Date(2025, 9, 3, 13, 0, 0, 500000000, time.UTC),
			SessionID: "session2",
			Success:   boolPtr(true),
		}
	}()
	
	// Aggregate the metrics
	summary := AggregateMetrics(entries, nil)
	
	// Verify results
	if summary.TotalCalls != 3 {
		t.Errorf("TotalCalls = %d, want 3", summary.TotalCalls)
	}
	
	if summary.UniqueSessions != 2 {
		t.Errorf("UniqueSessions = %d, want 2", summary.UniqueSessions)
	}
	
	// Check tool statistics
	editStats := summary.ToolStats["Edit"]
	if editStats == nil {
		t.Fatal("Edit stats not found")
	}
	if editStats.TotalCalls != 1 {
		t.Errorf("Edit TotalCalls = %d, want 1", editStats.TotalCalls)
	}
	if editStats.SuccessCount != 1 {
		t.Errorf("Edit SuccessCount = %d, want 1", editStats.SuccessCount)
	}
	if editStats.AvgDurationMs != 1000 { // 1 second
		t.Errorf("Edit AvgDurationMs = %.0f, want 1000", editStats.AvgDurationMs)
	}
	
	bashStats := summary.ToolStats["Bash"]
	if bashStats == nil {
		t.Fatal("Bash stats not found")
	}
	if bashStats.FailureCount != 1 {
		t.Errorf("Bash FailureCount = %d, want 1", bashStats.FailureCount)
	}
	if bashStats.ErrorTypes["exit_code_1"] != 1 {
		t.Errorf("Bash error count = %d, want 1", bashStats.ErrorTypes["exit_code_1"])
	}
	
	readStats := summary.ToolStats["Read"]
	if readStats == nil {
		t.Fatal("Read stats not found")
	}
	if readStats.AvgDurationMs != 500 { // 0.5 seconds
		t.Errorf("Read AvgDurationMs = %.0f, want 500", readStats.AvgDurationMs)
	}
	
	// Check success rate
	expectedRate := 2.0 / 3.0 * 100
	if summary.SuccessRate < expectedRate-1 || summary.SuccessRate > expectedRate+1 {
		t.Errorf("SuccessRate = %.1f%%, want ~%.1f%%", summary.SuccessRate, expectedRate)
	}
	
	// Check error patterns
	if len(summary.TopErrors) == 0 {
		t.Error("No error patterns found")
	}
}

func TestCorrelation(t *testing.T) {
	entries := make(chan MetricsEntry, 4)
	
	go func() {
		defer close(entries)
		
		// Out-of-order events
		entries <- MetricsEntry{
			ToolID:    "tool1",
			ToolName:  "Test",
			HookEvent: "PostToolUse",
			Timestamp: time.Now(),
			Success:   boolPtr(true),
		}
		entries <- MetricsEntry{
			ToolID:    "tool1",
			ToolName:  "Test",
			HookEvent: "PreToolUse",
			Timestamp: time.Now().Add(-1 * time.Second),
		}
		
		// Orphaned PreToolUse
		entries <- MetricsEntry{
			ToolID:    "tool2",
			ToolName:  "Test",
			HookEvent: "PreToolUse",
			Timestamp: time.Now(),
		}
		
		// Orphaned PostToolUse
		entries <- MetricsEntry{
			ToolID:    "tool3",
			ToolName:  "Test",
			HookEvent: "PostToolUse",
			Timestamp: time.Now(),
			DurationMs: intPtr(500),
		}
	}()
	
	summary := AggregateMetrics(entries, nil)
	
	stats := summary.ToolStats["Test"]
	// We have 2 PostToolUse events (counted as calls) and 2 PreToolUse events (one orphaned)
	// The aggregator counts PostToolUse as the primary call indicator
	if stats.TotalCalls != 2 {
		t.Errorf("TotalCalls = %d, want 2", stats.TotalCalls)
	}
	
	// Should have duration for tool1 (correlated) and tool3 (provided)
	if stats.TotalDurationMs == 0 {
		t.Error("No duration calculated for correlated events")
	}
}

func intPtr(i int64) *int64 {
	return &i
}