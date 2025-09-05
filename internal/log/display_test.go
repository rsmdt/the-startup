package log

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestDisplaySummary(t *testing.T) {
	// Create test summary
	summary := &MetricsSummary{
		Period: TimePeriod{
			Start: time.Date(2025, 9, 3, 0, 0, 0, 0, time.UTC),
			End:   time.Date(2025, 9, 3, 23, 59, 59, 0, time.UTC),
		},
		TotalCalls:     100,
		UniqueSessions: 10,
		SuccessRate:    95.0,
		ToolStats: map[string]*ToolStatistics{
			"Bash": {
				Name:          "Bash",
				TotalCalls:    50,
				SuccessCount:  48,
				FailureCount:  2,
				AvgDurationMs: 1234.5,
				MinDurationMs: 100,
				MaxDurationMs: 5000,
			},
			"Edit": {
				Name:          "Edit",
				TotalCalls:    30,
				SuccessCount:  30,
				FailureCount:  0,
				AvgDurationMs: 500.0,
				MinDurationMs: 200,
				MaxDurationMs: 1000,
			},
		},
	}
	
	// Test terminal format
	var buf bytes.Buffer
	opts := &DisplayOptions{
		Format: "terminal",
		Output: &buf,
		TopN:   5,
	}
	
	err := DisplaySummary(summary, opts)
	if err != nil {
		t.Fatalf("DisplaySummary failed: %v", err)
	}
	
	output := buf.String()
	
	// Check for key elements
	if !strings.Contains(output, "Metrics Summary") {
		t.Error("Missing title")
	}
	if !strings.Contains(output, "Total Invocations:     100") {
		t.Error("Missing total invocations")
	}
	if !strings.Contains(output, "Unique Sessions:       10") {
		t.Error("Missing unique sessions")
	}
	if !strings.Contains(output, "Overall Success Rate:  95.0%") {
		t.Error("Missing success rate")
	}
	if !strings.Contains(output, "Bash") {
		t.Error("Missing Bash tool")
	}
	if !strings.Contains(output, "Edit") {
		t.Error("Missing Edit tool")
	}
}

func TestDisplayJSON(t *testing.T) {
	summary := &MetricsSummary{
		TotalCalls:     100,
		UniqueSessions: 10,
		SuccessRate:    95.0,
	}
	
	var buf bytes.Buffer
	opts := &DisplayOptions{
		Format: "json",
		Output: &buf,
	}
	
	err := DisplaySummary(summary, opts)
	if err != nil {
		t.Fatalf("DisplaySummary JSON failed: %v", err)
	}
	
	output := buf.String()
	if !strings.Contains(output, `"total_calls": 100`) {
		t.Error("JSON missing total_calls")
	}
	if !strings.Contains(output, `"unique_sessions": 10`) {
		t.Error("JSON missing unique_sessions")
	}
	if !strings.Contains(output, `"success_rate": 95`) {
		t.Error("JSON missing success_rate")
	}
}

func TestDisplayCSV(t *testing.T) {
	summary := &MetricsSummary{
		ToolStats: map[string]*ToolStatistics{
			"Bash": {
				Name:          "Bash",
				TotalCalls:    50,
				SuccessCount:  48,
				FailureCount:  2,
				AvgDurationMs: 1234.5,
			},
		},
	}
	
	var buf bytes.Buffer
	opts := &DisplayOptions{
		Format: "csv",
		Output: &buf,
	}
	
	err := DisplayTools(summary, opts)
	if err != nil {
		t.Fatalf("DisplayTools CSV failed: %v", err)
	}
	
	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")
	
	// Check header
	if !strings.Contains(lines[0], "Tool,Total Calls,Success Count,Failure Count") {
		t.Errorf("CSV header incorrect. Got: %s", lines[0])
	}
	
	// Check data row
	if len(lines) < 2 {
		t.Fatal("CSV missing data rows")
	}
	if !strings.Contains(lines[1], "Bash,50,48,2") {
		t.Errorf("CSV data incorrect. Got: %s", lines[1])
	}
}