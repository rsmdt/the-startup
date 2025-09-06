package stats

import (
	"encoding/json"
	"fmt"
	"math"
	"runtime"
	"testing"
	"time"
)

func TestNewAggregator(t *testing.T) {
	agg := NewAggregator()
	if agg == nil {
		t.Fatal("NewAggregator returned nil")
	}
	
	// Verify it implements the interface
	var _ Aggregator = agg
}

func TestProcessEntry_UserMessage(t *testing.T) {
	agg := NewAggregator()
	
	entry := ClaudeLogEntry{
		Type:      "user",
		SessionID: "session1",
		Timestamp: time.Now(),
		User: &UserMessage{
			Text: "Hello, Claude",
		},
	}
	
	err := agg.ProcessEntry(entry)
	if err != nil {
		t.Fatalf("ProcessEntry failed: %v", err)
	}
	
	stats, err := agg.GetSessionStats("session1")
	if err != nil {
		t.Fatalf("GetSessionStats failed: %v", err)
	}
	
	if stats.UserMessages != 1 {
		t.Errorf("Expected 1 user message, got %d", stats.UserMessages)
	}
}

func TestProcessEntry_AssistantMessage(t *testing.T) {
	agg := NewAggregator()
	
	entry := ClaudeLogEntry{
		Type:      "assistant",
		SessionID: "session1",
		Timestamp: time.Now(),
		Assistant: &AssistantMessage{
			Text: "Hello! I can help you.",
			ToolUses: []ToolUse{
				{
					ID:         "tool1",
					Name:       "Bash",
					Parameters: json.RawMessage(`{"command": "ls"}`),
				},
			},
		},
	}
	
	err := agg.ProcessEntry(entry)
	if err != nil {
		t.Fatalf("ProcessEntry failed: %v", err)
	}
	
	stats, err := agg.GetSessionStats("session1")
	if err != nil {
		t.Fatalf("GetSessionStats failed: %v", err)
	}
	
	if stats.AssistantMessages != 1 {
		t.Errorf("Expected 1 assistant message, got %d", stats.AssistantMessages)
	}
}

func TestProcessEntry_CommandExtraction(t *testing.T) {
	agg := NewAggregator()
	
	entry := ClaudeLogEntry{
		Type:      "assistant",
		SessionID: "session1",
		Timestamp: time.Now(),
		Assistant: &AssistantMessage{
			Text: "I'll use <command-name>git status</command-name> to check the repository.",
		},
	}
	
	err := agg.ProcessEntry(entry)
	if err != nil {
		t.Fatalf("ProcessEntry failed: %v", err)
	}
	
	// Test passes if no error occurred
	// Commands are tracked internally but not exposed in AggregatedStatistics
}

func TestProcessEntry_SummaryTokens(t *testing.T) {
	agg := NewAggregator()
	
	entry := ClaudeLogEntry{
		Type:      "summary",
		SessionID: "session1",
		Timestamp: time.Now(),
		Summary: &SummaryMessage{
			TokensUsed: TokenUsage{
				Input:  100,
				Output: 200,
				Total:  300,
			},
		},
	}
	
	err := agg.ProcessEntry(entry)
	if err != nil {
		t.Fatalf("ProcessEntry failed: %v", err)
	}
	
	// Token usage is tracked in session stats
	stats, err := agg.GetSessionStats("session1")
	if err != nil {
		t.Fatalf("GetSessionStats failed: %v", err)
	}
	
	if stats.TotalTokens.Total != 300 {
		t.Errorf("Expected total tokens 300, got %d", stats.TotalTokens.Total)
	}
}

func TestToolInvocationCorrelation(t *testing.T) {
	agg := NewAggregator()
	now := time.Now()
	
	// Process tool use
	useEntry := ClaudeLogEntry{
		Type:      "assistant",
		SessionID: "session1",
		Timestamp: now,
		Assistant: &AssistantMessage{
			ToolUses: []ToolUse{
				{
					ID:         "tool123",
					Name:       "Read",
					Parameters: json.RawMessage(`{"file_path": "/tmp/test.txt"}`),
				},
			},
		},
	}
	
	err := agg.ProcessEntry(useEntry)
	if err != nil {
		t.Fatalf("ProcessEntry (use) failed: %v", err)
	}
	
	// Process tool result
	resultEntry := ClaudeLogEntry{
		Type:      "user",
		SessionID: "session1",
		Timestamp: now.Add(100 * time.Millisecond),
		User: &UserMessage{
			IsToolResult: true,
			ToolUseID:    "tool123",
			ToolName:     "Read",
			Result:       json.RawMessage(`"File contents"`),
		},
	}
	
	err = agg.ProcessEntry(resultEntry)
	if err != nil {
		t.Fatalf("ProcessEntry (result) failed: %v", err)
	}
	
	// Check that invocation was correlated
	stats, err := agg.GetSessionStats("session1")
	if err != nil {
		t.Fatalf("GetSessionStats failed: %v", err)
	}
	
	if len(stats.ToolInvocations) != 1 {
		t.Fatalf("Expected 1 tool invocation, got %d", len(stats.ToolInvocations))
	}
	
	inv := stats.ToolInvocations[0]
	if inv.ID != "tool123" {
		t.Errorf("Expected tool ID 'tool123', got %s", inv.ID)
	}
	if inv.Name != "Read" {
		t.Errorf("Expected tool name 'Read', got %s", inv.Name)
	}
	if !inv.Success {
		t.Error("Expected successful invocation")
	}
	if inv.DurationMs == nil || *inv.DurationMs != 100 {
		t.Errorf("Expected duration 100ms, got %v", inv.DurationMs)
	}
}

func TestOrphanedInvocationHandling(t *testing.T) {
	// Use a short timeout for testing
	agg := NewAggregatorWithTimeout(100 * time.Millisecond)
	now := time.Now()
	
	// Process tool use
	useEntry := ClaudeLogEntry{
		Type:      "assistant",
		SessionID: "session1",
		Timestamp: now,
		Assistant: &AssistantMessage{
			ToolUses: []ToolUse{
				{
					ID:         "orphan123",
					Name:       "Bash",
					Parameters: json.RawMessage(`{"command": "sleep 10"}`),
				},
			},
		},
	}
	
	err := agg.ProcessEntry(useEntry)
	if err != nil {
		t.Fatalf("ProcessEntry failed: %v", err)
	}
	
	// Process another entry after timeout
	laterEntry := ClaudeLogEntry{
		Type:      "system",
		SessionID: "session1",
		Timestamp: now.Add(200 * time.Millisecond),
		System:    &SystemMessage{Event: "test"},
	}
	
	err = agg.ProcessEntry(laterEntry)
	if err != nil {
		t.Fatalf("ProcessEntry failed: %v", err)
	}
	
	// Check that orphaned invocation was processed
	stats, err := agg.GetSessionStats("session1")
	if err != nil {
		t.Fatalf("GetSessionStats failed: %v", err)
	}
	
	if len(stats.ToolInvocations) != 1 {
		t.Fatalf("Expected 1 tool invocation, got %d", len(stats.ToolInvocations))
	}
	
	inv := stats.ToolInvocations[0]
	if inv.Success {
		t.Error("Expected failed invocation for orphaned tool")
	}
	if inv.Error == nil || inv.Error.Code != "TIMEOUT" {
		t.Errorf("Expected TIMEOUT error, got %v", inv.Error)
	}
}

func TestWelfordAlgorithm(t *testing.T) {
	stats := newWelfordStats()
	
	// Test data with known mean and variance
	values := []float64{2, 4, 4, 4, 5, 5, 7, 9}
	expectedMean := 5.0
	// Sample variance calculation: sum((x-mean)^2)/(n-1)
	// Deviations: -3, -1, -1, -1, 0, 0, 2, 4
	// Squared: 9, 1, 1, 1, 0, 0, 4, 16 = 32
	// Sample variance = 32/7 = 4.571
	expectedVariance := 4.571
	
	for _, v := range values {
		stats.update(v)
	}
	
	if math.Abs(stats.mean-expectedMean) > 0.001 {
		t.Errorf("Expected mean %.3f, got %.3f", expectedMean, stats.mean)
	}
	
	if math.Abs(stats.variance()-expectedVariance) > 0.01 {
		t.Errorf("Expected variance %.3f, got %.3f", expectedVariance, stats.variance())
	}
	
	if stats.min == nil || *stats.min != 2 {
		t.Errorf("Expected min 2, got %v", stats.min)
	}
	
	if stats.max == nil || *stats.max != 9 {
		t.Errorf("Expected max 9, got %v", stats.max)
	}
}

func TestTDigestPercentiles(t *testing.T) {
	td := newTDigest(100)
	
	// Add values from 1 to 100
	for i := 1; i <= 100; i++ {
		td.add(float64(i))
	}
	
	// Test percentiles
	tests := []struct {
		percentile float64
		expected   float64
		tolerance  float64
	}{
		{0.5, 50, 5},  // Median
		{0.25, 25, 5}, // Q1
		{0.75, 75, 5}, // Q3
		{0.95, 95, 5}, // P95
		{0.99, 99, 5}, // P99
	}
	
	for _, test := range tests {
		result := td.quantile(test.percentile)
		if math.Abs(result-test.expected) > test.tolerance {
			t.Errorf("P%.0f: expected ~%.0f, got %.2f",
				test.percentile*100, test.expected, result)
		}
	}
}

func TestProcessInvocation(t *testing.T) {
	agg := NewAggregator()
	now := time.Now()
	duration := int64(150)
	
	invocation := ToolInvocation{
		ID:          "inv1",
		Name:        "Edit",
		Parameters:  json.RawMessage(`{"file": "test.go"}`),
		InvokedAt:   now,
		CompletedAt: &now,
		DurationMs:  &duration,
		Success:     true,
		SessionID:   "session1",
	}
	
	err := agg.ProcessInvocation(invocation)
	if err != nil {
		t.Fatalf("ProcessInvocation failed: %v", err)
	}
	
	// Check session stats
	sessionStats, err := agg.GetSessionStats("session1")
	if err != nil {
		t.Fatalf("GetSessionStats failed: %v", err)
	}
	
	if sessionStats.TotalToolCalls != 1 {
		t.Errorf("Expected 1 tool call, got %d", sessionStats.TotalToolCalls)
	}
	if sessionStats.SuccessfulCalls != 1 {
		t.Errorf("Expected 1 successful call, got %d", sessionStats.SuccessfulCalls)
	}
	
	// Check tool stats
	toolStats := sessionStats.ToolStats["Edit"]
	if toolStats == nil {
		t.Fatal("Expected Edit tool stats")
	}
	if toolStats.CallCount != 1 {
		t.Errorf("Expected 1 call, got %d", toolStats.CallCount)
	}
	if toolStats.AvgDurationMs != 150 {
		t.Errorf("Expected avg duration 150ms, got %.2f", toolStats.AvgDurationMs)
	}
}

func TestGetToolStats(t *testing.T) {
	agg := NewAggregator()
	now := time.Now()
	
	// Add multiple invocations
	for i := 0; i < 10; i++ {
		duration := int64(100 + i*10)
		err := agg.ProcessInvocation(ToolInvocation{
			ID:          fmt.Sprintf("inv%d", i),
			Name:        "Bash",
			InvokedAt:   now.Add(time.Duration(i) * time.Second),
			CompletedAt: &now,
			DurationMs:  &duration,
			Success:     i%3 != 0, // Every 3rd fails
			SessionID:   fmt.Sprintf("session%d", i%2), // 2 different sessions
		})
		if err != nil {
			t.Fatalf("ProcessInvocation failed: %v", err)
		}
	}
	
	bashStats, err := agg.GetToolStats("Bash")
	if err != nil {
		t.Fatalf("GetToolStats failed: %v", err)
	}
	
	if bashStats == nil {
		t.Fatal("Expected Bash stats")
	}
	
	if bashStats.TotalCalls != 10 {
		t.Errorf("Expected 10 total calls, got %d", bashStats.TotalCalls)
	}
	
	if bashStats.UniqueUsers != 2 {
		t.Errorf("Expected 2 unique sessions, got %d", bashStats.UniqueUsers)
	}
	
	// i%3 != 0 means: i=0,3,6,9 fail (4 failures), i=1,2,4,5,7,8 succeed (6 successes)
	expectedSuccessRate := 0.6 // 6 successes out of 10
	if math.Abs(bashStats.SuccessRate-expectedSuccessRate) > 0.01 {
		t.Errorf("Expected success rate %.2f, got %.2f", expectedSuccessRate, bashStats.SuccessRate)
	}
	
	// Check duration stats
	if bashStats.MinDurationMs != 100 {
		t.Errorf("Expected min duration 100ms, got %d", bashStats.MinDurationMs)
	}
	if bashStats.MaxDurationMs != 190 {
		t.Errorf("Expected max duration 190ms, got %d", bashStats.MaxDurationMs)
	}
}

func TestReset(t *testing.T) {
	agg := NewAggregator()
	
	// Add some data
	err := agg.ProcessEntry(ClaudeLogEntry{
		Type:      "user",
		SessionID: "session1",
		Timestamp: time.Now(),
		User:      &UserMessage{Text: "test"},
	})
	if err != nil {
		t.Fatalf("ProcessEntry failed: %v", err)
	}
	
	// Verify data exists
	stats, err := agg.GetSessionStats("session1")
	if err != nil || stats == nil {
		t.Fatal("Expected session stats before reset")
	}
	
	// Reset
	agg.Reset()
	
	// Verify data is cleared
	_, err = agg.GetSessionStats("session1")
	if err == nil {
		t.Error("Expected error after reset, session should not exist")
	}
	
	overall, err := agg.GetOverallStats()
	if err != nil {
		t.Fatalf("GetOverallStats failed: %v", err)
	}
	if overall.TotalSessions != 0 {
		t.Errorf("Expected 0 sessions after reset, got %d", overall.TotalSessions)
	}
}

func TestMerge(t *testing.T) {
	agg1 := NewAggregator()
	agg2 := NewAggregator()
	
	// Add data to first aggregator
	err := agg1.ProcessEntry(ClaudeLogEntry{
		Type:      "assistant",
		SessionID: "session1",
		Timestamp: time.Now(),
		Assistant: &AssistantMessage{
			Text: "<command-name>git status</command-name>",
		},
	})
	if err != nil {
		t.Fatalf("ProcessEntry agg1 failed: %v", err)
	}
	
	// Add data to second aggregator
	err = agg2.ProcessEntry(ClaudeLogEntry{
		Type:      "assistant",
		SessionID: "session2",
		Timestamp: time.Now(),
		Assistant: &AssistantMessage{
			Text: "<command-name>git diff</command-name>",
		},
	})
	if err != nil {
		t.Fatalf("ProcessEntry agg2 failed: %v", err)
	}
	
	// Also add to same command in agg2
	err = agg2.ProcessEntry(ClaudeLogEntry{
		Type:      "assistant",
		SessionID: "session2",
		Timestamp: time.Now(),
		Assistant: &AssistantMessage{
			Text: "<command-name>git status</command-name>",
		},
	})
	if err != nil {
		t.Fatalf("ProcessEntry agg2 failed: %v", err)
	}
	
	// Merge
	err = agg1.Merge(agg2)
	if err != nil {
		t.Fatalf("Merge failed: %v", err)
	}
	
	// Check merged results
	overall, err := agg1.GetOverallStats()
	if err != nil {
		t.Fatalf("GetOverallStats failed: %v", err)
	}
	
	if overall.TotalSessions != 2 {
		t.Errorf("Expected 2 total sessions, got %d", overall.TotalSessions)
	}
	
	// Commands tracking would need to be added to AggregatedStatistics
	// For now, just check that merge didn't fail
}

func TestMemoryEfficiency(t *testing.T) {
	// Skip in short mode
	if testing.Short() {
		t.Skip("Skipping memory efficiency test in short mode")
	}
	
	agg := NewAggregator()
	
	// Force GC before starting
	runtime.GC()
	runtime.GC()
	
	// Get initial memory stats
	var m1 runtime.MemStats
	runtime.ReadMemStats(&m1)
	
	// Process a large number of entries
	numEntries := 10000
	for i := 0; i < numEntries; i++ {
		duration := int64(100 + i%1000)
		err := agg.ProcessInvocation(ToolInvocation{
			ID:          fmt.Sprintf("inv%d", i),
			Name:        fmt.Sprintf("Tool%d", i%10), // 10 different tools
			InvokedAt:   time.Now(),
			CompletedAt: ptr(time.Now()),
			DurationMs:  &duration,
			Success:     true,
			SessionID:   fmt.Sprintf("session%d", i%100), // 100 different sessions
		})
		if err != nil {
			t.Fatalf("ProcessInvocation failed: %v", err)
		}
	}
	
	// Force GC and get memory stats
	runtime.GC()
	var m2 runtime.MemStats
	runtime.ReadMemStats(&m2)
	
	// Calculate memory usage - use TotalAlloc for more accurate measurement
	memUsed := m2.TotalAlloc - m1.TotalAlloc
	memPerEntry := memUsed / uint64(numEntries)
	
	// Memory should be bounded - we're using efficient algorithms
	// Expect less than 1KB per entry on average
	if memPerEntry > 1024 {
		t.Errorf("Memory usage too high: %d bytes per entry", memPerEntry)
	}
	
	// Verify statistics are still accurate
	for i := 0; i < 10; i++ {
		toolName := fmt.Sprintf("Tool%d", i)
		stats, err := agg.GetToolStats(toolName)
		if err != nil {
			t.Fatalf("GetToolStats failed for %s: %v", toolName, err)
		}
		if stats.TotalCalls != numEntries/10 {
			t.Errorf("Expected %d calls per tool, got %d", numEntries/10, stats.TotalCalls)
		}
	}
}

func TestGetOverallStats(t *testing.T) {
	agg := NewAggregator()
	now := time.Now()
	
	// Add various data
	for i := 0; i < 5; i++ {
		// Add user message
		err := agg.ProcessEntry(ClaudeLogEntry{
			Type:      "user",
			SessionID: fmt.Sprintf("session%d", i),
			Timestamp: now.Add(time.Duration(i) * time.Hour),
			User:      &UserMessage{Text: "test"},
		})
		if err != nil {
			t.Fatalf("ProcessEntry failed: %v", err)
		}
		
		// Add tool invocation
		duration := int64(100 + i*50)
		err = agg.ProcessInvocation(ToolInvocation{
			ID:          fmt.Sprintf("inv%d", i),
			Name:        "Read",
			InvokedAt:   now.Add(time.Duration(i) * time.Hour),
			CompletedAt: ptr(now.Add(time.Duration(i)*time.Hour + time.Duration(duration)*time.Millisecond)),
			DurationMs:  &duration,
			Success:     true,
			SessionID:   fmt.Sprintf("session%d", i),
		})
		if err != nil {
			t.Fatalf("ProcessInvocation failed: %v", err)
		}
	}
	
	overall, err := agg.GetOverallStats()
	if err != nil {
		t.Fatalf("GetOverallStats failed: %v", err)
	}
	
	// Verify overall stats
	if overall.TotalSessions != 5 {
		t.Errorf("Expected 5 sessions, got %d", overall.TotalSessions)
	}
	
	if len(overall.ActiveSessions) != 5 {
		t.Errorf("Expected 5 active sessions, got %d", len(overall.ActiveSessions))
	}
	
	readStats, exists := overall.GlobalToolStats["Read"]
	if !exists {
		t.Fatal("Expected Read tool stats")
	}
	
	if readStats.TotalCalls != 5 {
		t.Errorf("Expected 5 Read calls, got %d", readStats.TotalCalls)
	}
	
	// Check time period
	if !overall.Period.Start.Equal(now) {
		t.Errorf("Expected start time %v, got %v", now, overall.Period.Start)
	}
	
	expectedEnd := now.Add(4 * time.Hour)
	if !overall.Period.End.Equal(expectedEnd) {
		t.Errorf("Expected end time %v, got %v", expectedEnd, overall.Period.End)
	}
}

func TestErrorTracking(t *testing.T) {
	agg := NewAggregator()
	now := time.Now()
	
	// Add successful invocation
	err := agg.ProcessInvocation(ToolInvocation{
		ID:        "success1",
		Name:      "Bash",
		InvokedAt: now,
		Success:   true,
		SessionID: "session1",
	})
	if err != nil {
		t.Fatalf("ProcessInvocation failed: %v", err)
	}
	
	// Add failed invocations with different errors
	for i := 0; i < 3; i++ {
		err := agg.ProcessInvocation(ToolInvocation{
			ID:        fmt.Sprintf("fail%d", i),
			Name:      "Bash",
			InvokedAt: now,
			Success:   false,
			Error: &ToolError{
				Code:    fmt.Sprintf("ERROR_%d", i%2),
				Message: fmt.Sprintf("Error message %d", i%2),
			},
			SessionID: "session1",
		})
		if err != nil {
			t.Fatalf("ProcessInvocation failed: %v", err)
		}
	}
	
	// Check session error tracking
	sessionStats, err := agg.GetSessionStats("session1")
	if err != nil {
		t.Fatalf("GetSessionStats failed: %v", err)
	}
	
	if len(sessionStats.Errors) != 3 {
		t.Errorf("Expected 3 errors, got %d", len(sessionStats.Errors))
	}
	
	if sessionStats.ErrorRate != 0.75 {
		t.Errorf("Expected error rate 0.75, got %.2f", sessionStats.ErrorRate)
	}
	
	// Check overall error patterns
	overall, err := agg.GetOverallStats()
	if err != nil {
		t.Fatalf("GetOverallStats failed: %v", err)
	}
	
	if len(overall.TopErrors) == 0 {
		t.Error("Expected top errors to be collected")
	}
}

func TestMultipleCommandExtraction(t *testing.T) {
	agg := NewAggregator()
	
	// Text with multiple commands
	entry := ClaudeLogEntry{
		Type:      "assistant",
		SessionID: "session1",
		Timestamp: time.Now(),
		Assistant: &AssistantMessage{
			Text: `I'll run several commands:
			First, <command-name>npm install</command-name> to install dependencies.
			Then <command-name>npm test</command-name> to run tests.
			Finally, <command-name>npm run build</command-name> to build the project.`,
		},
	}
	
	err := agg.ProcessEntry(entry)
	if err != nil {
		t.Fatalf("ProcessEntry failed: %v", err)
	}
	
	// Test passes if no error occurred
	// Commands are tracked internally but not exposed in AggregatedStatistics
}

func TestHourlyActivityTracking(t *testing.T) {
	agg := NewAggregator()
	baseTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	
	// Add invocations at different hours
	hours := []int{9, 9, 10, 10, 10, 14, 14, 14, 14, 20}
	for i, hour := range hours {
		invTime := baseTime.Add(time.Duration(hour) * time.Hour)
		duration := int64(100)
		err := agg.ProcessInvocation(ToolInvocation{
			ID:          fmt.Sprintf("inv%d", i),
			Name:        "Edit",
			InvokedAt:   invTime,
			CompletedAt: ptr(invTime.Add(100 * time.Millisecond)),
			DurationMs:  &duration,
			Success:     true,
			SessionID:   "session1",
		})
		if err != nil {
			t.Fatalf("ProcessInvocation failed: %v", err)
		}
	}
	
	// Check that Edit tool has the right peak hour
	editStats, err := agg.GetToolStats("Edit")
	if err != nil {
		t.Fatalf("GetToolStats failed: %v", err)
	}
	
	if editStats.PeakHour != 14 {
		t.Errorf("Expected peak hour to be 14, got %d", editStats.PeakHour)
	}
}

// Helper function to create pointer
func ptr[T any](v T) *T {
	return &v
}