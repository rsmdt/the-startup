package stats

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestNewDisplayFormatter(t *testing.T) {
	formatter := NewDisplayFormatter()
	if formatter == nil {
		t.Fatal("NewDisplayFormatter returned nil")
	}
	if formatter.outputFormat != "table" {
		t.Errorf("Expected default format 'table', got %s", formatter.outputFormat)
	}
}

func TestSetOutputFormat(t *testing.T) {
	formatter := NewDisplayFormatter()
	
	tests := []struct {
		name        string
		format      string
		expectError bool
	}{
		{"Valid table format", "table", false},
		{"Valid json format", "json", false},
		{"Valid csv format", "csv", false},
		{"Invalid format", "xml", true},
		{"Empty format", "", true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := formatter.SetOutputFormat(tt.format)
			if tt.expectError && err == nil {
				t.Errorf("Expected error for format %s, got nil", tt.format)
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error for format %s: %v", tt.format, err)
			}
			if !tt.expectError && formatter.outputFormat != tt.format {
				t.Errorf("Format not set correctly: expected %s, got %s", tt.format, formatter.outputFormat)
			}
		})
	}
}

func TestWrite(t *testing.T) {
	formatter := NewDisplayFormatter()
	var buf bytes.Buffer
	
	content := "test content"
	err := formatter.Write(&buf, content)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}
	
	if buf.String() != content {
		t.Errorf("Write output mismatch: expected %s, got %s", content, buf.String())
	}
}

func TestFormatStats(t *testing.T) {
	formatter := NewDisplayFormatter()
	
	// Test with nil stats
	result := formatter.FormatStats(nil)
	if result != "" {
		t.Error("Expected empty string for nil stats")
	}
	
	// Create test stats
	stats := &AggregatedStatistics{
		Period: TimePeriod{
			Start: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			End:   time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
		},
		TotalSessions:      5,
		ActiveSessions:     []string{"session1", "session2"},
		AvgSessionDuration: 30 * time.Minute,
		TotalToolCalls:     100,
		OverallSuccessRate: 0.95,
		PerformanceMetrics: PerformanceMetrics{
			AvgResponseTimeMs:  150,
			P50ResponseTimeMs:  100,
			P95ResponseTimeMs:  500,
			P99ResponseTimeMs:  1000,
			ToolCallsPerMinute: 10.5,
		},
	}
	
	// Test table format
	formatter.SetOutputFormat("table")
	result = formatter.FormatStats(stats)
	if !strings.Contains(result, "AGGREGATED STATISTICS") {
		t.Error("Table format missing header")
	}
	if !strings.Contains(result, "Total Sessions:") {
		t.Error("Table format missing session count")
	}
	if !strings.Contains(result, "95.0%") {
		t.Error("Table format missing success rate")
	}
	
	// Test JSON format
	formatter.SetOutputFormat("json")
	result = formatter.FormatStats(stats)
	var jsonStats AggregatedStatistics
	if err := json.Unmarshal([]byte(result), &jsonStats); err != nil {
		t.Errorf("Invalid JSON output: %v", err)
	}
	if jsonStats.TotalSessions != 5 {
		t.Errorf("JSON format incorrect: expected 5 sessions, got %d", jsonStats.TotalSessions)
	}
	
	// Test CSV format
	formatter.SetOutputFormat("csv")
	result = formatter.FormatStats(stats)
	if !strings.Contains(result, "Metric,Value") {
		t.Error("CSV format missing headers")
	}
	if !strings.Contains(result, "Total Sessions,5") {
		t.Error("CSV format missing session count")
	}
}

func TestFormatSession(t *testing.T) {
	formatter := NewDisplayFormatter()
	
	// Test with nil session
	result := formatter.FormatSession(nil)
	if result != "" {
		t.Error("Expected empty string for nil session")
	}
	
	// Create test session
	endTime := time.Date(2024, 1, 1, 1, 0, 0, 0, time.UTC)
	session := &SessionStatistics{
		SessionID:         "test-session",
		StartTime:         time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		EndTime:          &endTime,
		Duration:         time.Hour,
		UserMessages:     10,
		AssistantMessages: 15,
		SystemMessages:   2,
		TotalToolCalls:   20,
		SuccessfulCalls:  18,
		FailedCalls:      2,
		ErrorRate:        0.1,
		TotalTokens: TokenUsage{
			Input:  1000,
			Output: 2000,
			Total:  3000,
		},
		ToolStats: map[string]*ToolSessionStats{
			"Bash": {
				Name:          "Bash",
				CallCount:     10,
				SuccessCount:  9,
				FailureCount:  1,
				AvgDurationMs: 250,
			},
			"Read": {
				Name:          "Read",
				CallCount:     10,
				SuccessCount:  9,
				FailureCount:  1,
				AvgDurationMs: 50,
			},
		},
	}
	
	// Test table format
	formatter.SetOutputFormat("table")
	result = formatter.FormatSession(session)
	if !strings.Contains(result, "SESSION: test-session") {
		t.Error("Table format missing session ID")
	}
	if !strings.Contains(result, "User Messages:") {
		t.Error("Table format missing message counts")
	}
	if !strings.Contains(result, "TOOL BREAKDOWN") {
		t.Error("Table format missing tool breakdown")
	}
	
	// Test JSON format
	formatter.SetOutputFormat("json")
	result = formatter.FormatSession(session)
	var jsonSession SessionStatistics
	if err := json.Unmarshal([]byte(result), &jsonSession); err != nil {
		t.Errorf("Invalid JSON output: %v", err)
	}
	if jsonSession.SessionID != "test-session" {
		t.Errorf("JSON format incorrect session ID: %s", jsonSession.SessionID)
	}
	
	// Test CSV format
	formatter.SetOutputFormat("csv")
	result = formatter.FormatSession(session)
	if !strings.Contains(result, "Session ID,test-session") {
		t.Error("CSV format missing session ID")
	}
	if !strings.Contains(result, "Tool Breakdown") {
		t.Error("CSV format missing tool breakdown section")
	}
}

func TestFormatToolLeaderboard(t *testing.T) {
	formatter := NewDisplayFormatter()
	
	// Test with nil/empty map
	result := formatter.FormatToolLeaderboard(nil, 0)
	if result != "" {
		t.Error("Expected empty string for nil tool stats")
	}
	
	result = formatter.FormatToolLeaderboard(map[string]*GlobalToolStats{}, 0)
	if result != "" {
		t.Error("Expected empty string for empty tool stats")
	}
	
	// Create test tool stats
	toolStats := map[string]*GlobalToolStats{
		"Bash": {
			Name:             "Bash",
			TotalCalls:       100,
			SuccessRate:      0.95,
			AvgDurationMs:    250,
			MinDurationMs:    10,
			MaxDurationMs:    5000,
			P95DurationMs:    1000,
			P99DurationMs:    2000,
		},
		"Read": {
			Name:             "Read",
			TotalCalls:       50,
			SuccessRate:      0.98,
			AvgDurationMs:    50,
			MinDurationMs:    5,
			MaxDurationMs:    500,
			P95DurationMs:    100,
			P99DurationMs:    200,
		},
		"Edit": {
			Name:             "Edit",
			TotalCalls:       30,
			SuccessRate:      0.90,
			AvgDurationMs:    150,
			MinDurationMs:    20,
			MaxDurationMs:    1000,
			P95DurationMs:    500,
			P99DurationMs:    800,
		},
	}
	
	// Test table format with no limit
	formatter.SetOutputFormat("table")
	result = formatter.FormatToolLeaderboard(toolStats, 0)
	if !strings.Contains(result, "TOOL USAGE LEADERBOARD") {
		t.Error("Table format missing header")
	}
	if !strings.Contains(result, "Bash") {
		t.Error("Table format missing Bash tool")
	}
	if !strings.Contains(result, "Sparkline") {
		t.Error("Table format missing sparkline column")
	}
	
	// Check order (Bash should be first with 100 calls)
	lines := strings.Split(result, "\n")
	bashFound := false
	readFound := false
	for _, line := range lines {
		// tabwriter formats with spaces, not literal tabs
		if strings.Contains(line, "1") && strings.Contains(line, "Bash") && strings.Contains(line, "100") {
			bashFound = true
		}
		if strings.Contains(line, "2") && strings.Contains(line, "Read") && strings.Contains(line, "50") {
			readFound = true
		}
	}
	if !bashFound {
		t.Error("Bash should be ranked #1 with 100 calls")
	}
	if !readFound {
		t.Error("Read should be ranked #2 with 50 calls")
	}
	
	// Test with limit
	result = formatter.FormatToolLeaderboard(toolStats, 2)
	if strings.Contains(result, "Edit") {
		t.Error("Edit should not appear when limit is 2")
	}
	
	// Test JSON format
	formatter.SetOutputFormat("json")
	result = formatter.FormatToolLeaderboard(toolStats, 0)
	var jsonResult map[string]interface{}
	if err := json.Unmarshal([]byte(result), &jsonResult); err != nil {
		t.Errorf("Invalid JSON output: %v", err)
	}
	if jsonResult["total_tools"].(float64) != 3 {
		t.Errorf("JSON format incorrect total_tools: %v", jsonResult["total_tools"])
	}
	
	// Test CSV format
	formatter.SetOutputFormat("csv")
	result = formatter.FormatToolLeaderboard(toolStats, 0)
	if !strings.Contains(result, "Rank,Tool,Total Calls") {
		t.Error("CSV format missing headers")
	}
	if !strings.Contains(result, "1,Bash,100") {
		t.Error("CSV format missing Bash entry")
	}
}

func TestFormatErrors(t *testing.T) {
	formatter := NewDisplayFormatter()
	
	// Test with empty slices
	result := formatter.FormatErrors([]ErrorPattern{}, []ErrorFrequency{})
	if result != "" {
		t.Error("Expected empty string for empty errors")
	}
	
	// Create test error data
	patterns := []ErrorPattern{
		{
			Pattern:   "Connection timeout",
			Count:     10,
			Tools:     []string{"Bash", "WebFetch"},
			FirstSeen: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			LastSeen:  time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
		},
		{
			Pattern:   "File not found",
			Count:     5,
			Tools:     []string{"Read"},
			FirstSeen: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			LastSeen:  time.Date(2024, 1, 1, 18, 0, 0, 0, time.UTC),
		},
	}
	
	frequencies := []ErrorFrequency{
		{
			ErrorCode:    "TIMEOUT",
			ErrorMessage: "Operation timed out",
			Count:        10,
			Percentage:   66.7,
			ToolName:     "Bash",
		},
		{
			ErrorCode:    "NOT_FOUND",
			ErrorMessage: "File not found",
			Count:        5,
			Percentage:   33.3,
			ToolName:     "Read",
		},
	}
	
	// Test table format
	formatter.SetOutputFormat("table")
	result = formatter.FormatErrors(patterns, frequencies)
	if !strings.Contains(result, "ERROR ANALYSIS") {
		t.Error("Table format missing header")
	}
	if !strings.Contains(result, "ERROR PATTERNS") {
		t.Error("Table format missing patterns section")
	}
	if !strings.Contains(result, "TOP ERRORS BY FREQUENCY") {
		t.Error("Table format missing frequency section")
	}
	
	// Test JSON format
	formatter.SetOutputFormat("json")
	result = formatter.FormatErrors(patterns, frequencies)
	var jsonResult map[string]interface{}
	if err := json.Unmarshal([]byte(result), &jsonResult); err != nil {
		t.Errorf("Invalid JSON output: %v", err)
	}
	if _, ok := jsonResult["patterns"]; !ok {
		t.Error("JSON format missing patterns")
	}
	if _, ok := jsonResult["frequencies"]; !ok {
		t.Error("JSON format missing frequencies")
	}
	
	// Test CSV format
	formatter.SetOutputFormat("csv")
	result = formatter.FormatErrors(patterns, frequencies)
	if !strings.Contains(result, "Error Patterns") {
		t.Error("CSV format missing patterns section")
	}
	if !strings.Contains(result, "Error Frequencies") {
		t.Error("CSV format missing frequencies section")
	}
}

func TestFormatTimeline(t *testing.T) {
	formatter := NewDisplayFormatter()
	
	// Test with empty slices
	result := formatter.FormatTimeline([]HourlyActivity{}, []DailyActivity{})
	if result != "" {
		t.Error("Expected empty string for empty timeline")
	}
	
	// Create test timeline data
	hourly := []HourlyActivity{
		{
			Hour:           time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
			ToolCalls:      50,
			UniqueTools:    5,
			UniqueSessions: 3,
			SuccessCount:   45,
			FailureCount:   5,
			AvgDurationMs:  200,
		},
		{
			Hour:           time.Date(2024, 1, 1, 11, 0, 0, 0, time.UTC),
			ToolCalls:      60,
			UniqueTools:    6,
			UniqueSessions: 4,
			SuccessCount:   58,
			FailureCount:   2,
			AvgDurationMs:  180,
		},
	}
	
	daily := []DailyActivity{
		{
			Date:            time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			Sessions:        10,
			ToolCalls:       500,
			UniqueTools:     15,
			SuccessRate:     0.92,
			TotalDurationMs: 100000,
			Errors:          40,
		},
		{
			Date:            time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
			Sessions:        12,
			ToolCalls:       600,
			UniqueTools:     18,
			SuccessRate:     0.95,
			TotalDurationMs: 120000,
			Errors:          30,
		},
	}
	
	// Test table format
	formatter.SetOutputFormat("table")
	result = formatter.FormatTimeline(hourly, daily)
	if !strings.Contains(result, "TEMPORAL ACTIVITY") {
		t.Error("Table format missing header")
	}
	if !strings.Contains(result, "DAILY ACTIVITY") {
		t.Error("Table format missing daily section")
	}
	if !strings.Contains(result, "HOURLY ACTIVITY") {
		t.Error("Table format missing hourly section")
	}
	
	// Test JSON format
	formatter.SetOutputFormat("json")
	result = formatter.FormatTimeline(hourly, daily)
	var jsonResult map[string]interface{}
	if err := json.Unmarshal([]byte(result), &jsonResult); err != nil {
		t.Errorf("Invalid JSON output: %v", err)
	}
	if _, ok := jsonResult["hourly"]; !ok {
		t.Error("JSON format missing hourly data")
	}
	if _, ok := jsonResult["daily"]; !ok {
		t.Error("JSON format missing daily data")
	}
	
	// Test CSV format
	formatter.SetOutputFormat("csv")
	result = formatter.FormatTimeline(hourly, daily)
	if !strings.Contains(result, "Daily Activity") {
		t.Error("CSV format missing daily section")
	}
	if !strings.Contains(result, "Hourly Activity") {
		t.Error("CSV format missing hourly section")
	}
}

func TestFormatCommandLeaderboard(t *testing.T) {
	formatter := NewDisplayFormatter()
	
	// Test with nil/empty map
	result := formatter.FormatCommandLeaderboard(nil)
	if result != "" {
		t.Error("Expected empty string for nil commands")
	}
	
	result = formatter.FormatCommandLeaderboard(map[string]int{})
	if result != "" {
		t.Error("Expected empty string for empty commands")
	}
	
	// Create test command data
	commands := map[string]int{
		"specify":  50,
		"build":    30,
		"test":     20,
		"analyze":  15,
		"optimize": 10,
	}
	
	// Test table format
	formatter.SetOutputFormat("table")
	result = formatter.FormatCommandLeaderboard(commands)
	if !strings.Contains(result, "COMMAND/AGENT LEADERBOARD") {
		t.Error("Table format missing header")
	}
	if !strings.Contains(result, "specify") {
		t.Error("Table format missing specify command")
	}
	if !strings.Contains(result, "Sparkline") {
		t.Error("Table format missing sparkline column")
	}
	
	// Check order (specify should be first with 50 uses)
	// tabwriter formats with spaces, not literal tabs
	specifyFound := false
	lines := strings.Split(result, "\n")
	for _, line := range lines {
		if strings.Contains(line, "1") && strings.Contains(line, "specify") && strings.Contains(line, "50") {
			specifyFound = true
			break
		}
	}
	if !specifyFound {
		t.Error("specify should be ranked #1 with 50 uses")
	}
	
	// Test JSON format
	formatter.SetOutputFormat("json")
	result = formatter.FormatCommandLeaderboard(commands)
	var jsonResult map[string]interface{}
	if err := json.Unmarshal([]byte(result), &jsonResult); err != nil {
		t.Errorf("Invalid JSON output: %v", err)
	}
	if jsonResult["total_commands"].(float64) != 5 {
		t.Errorf("JSON format incorrect total_commands: %v", jsonResult["total_commands"])
	}
	
	// Test CSV format
	formatter.SetOutputFormat("csv")
	result = formatter.FormatCommandLeaderboard(commands)
	if !strings.Contains(result, "Rank,Command/Agent,Usage Count") {
		t.Error("CSV format missing headers")
	}
	if !strings.Contains(result, "1,specify,50") {
		t.Error("CSV format missing specify entry")
	}
}

func TestHelperFunctions(t *testing.T) {
	// Test formatDuration
	tests := []struct {
		duration time.Duration
		expected string
	}{
		{30 * time.Second, "30s"},
		{90 * time.Second, "2m"},
		{90 * time.Minute, "1.5h"},
	}
	
	for _, tt := range tests {
		result := formatDuration(tt.duration)
		if result != tt.expected {
			t.Errorf("formatDuration(%v) = %s, expected %s", tt.duration, result, tt.expected)
		}
	}
	
	// Test truncateString
	truncTests := []struct {
		input    string
		maxLen   int
		expected string
	}{
		{"short", 10, "short"},
		{"this is a very long string", 10, "this is..."},
		{"exact", 5, "exact"},
	}
	
	for _, tt := range truncTests {
		result := truncateString(tt.input, tt.maxLen)
		if result != tt.expected {
			t.Errorf("truncateString(%s, %d) = %s, expected %s", tt.input, tt.maxLen, result, tt.expected)
		}
	}
	
	// Test generateSparkline
	sparkTests := []struct {
		value    int
		max      int
		width    int
		hasChars bool
	}{
		{50, 100, 10, true},
		{0, 100, 10, false},
		{100, 100, 10, true},
		{25, 0, 10, false}, // max = 0 case
	}
	
	for _, tt := range sparkTests {
		result := generateSparkline(tt.value, tt.max, tt.width)
		if len(result) == 0 {
			t.Errorf("generateSparkline(%d, %d, %d) returned empty string", tt.value, tt.max, tt.width)
		}
		// Check that sparkline has expected width (unicode characters count as multiple bytes)
		runeCount := len([]rune(result))
		if runeCount != tt.width {
			t.Errorf("generateSparkline(%d, %d, %d) width = %d, expected %d", tt.value, tt.max, tt.width, runeCount, tt.width)
		}
	}
}

func TestTableFormatting(t *testing.T) {
	formatter := NewDisplayFormatter()
	formatter.SetOutputFormat("table")
	
	// Test that table output contains box drawing characters
	stats := &AggregatedStatistics{
		Period: TimePeriod{
			Start: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			End:   time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
		},
		TotalSessions:      1,
		OverallSuccessRate: 1.0,
	}
	
	result := formatter.FormatStats(stats)
	if !strings.Contains(result, "╔") || !strings.Contains(result, "╗") {
		t.Error("Table format missing box drawing characters")
	}
	if !strings.Contains(result, "═") {
		t.Error("Table format missing horizontal lines")
	}
}

func TestCSVFormatting(t *testing.T) {
	formatter := NewDisplayFormatter()
	formatter.SetOutputFormat("csv")
	
	// Test that CSV output is properly formatted
	stats := &AggregatedStatistics{
		Period: TimePeriod{
			Start: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			End:   time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
		},
		TotalSessions: 5,
	}
	
	result := formatter.FormatStats(stats)
	lines := strings.Split(strings.TrimSpace(result), "\n")
	
	// Check header
	if lines[0] != "Metric,Value" {
		t.Errorf("CSV header incorrect: %s", lines[0])
	}
	
	// Check that each line has correct number of columns
	for i, line := range lines {
		fields := strings.Split(line, ",")
		if len(fields) != 2 {
			t.Errorf("CSV line %d has %d columns, expected 2: %s", i, len(fields), line)
		}
	}
}

func TestJSONFormatting(t *testing.T) {
	formatter := NewDisplayFormatter()
	formatter.SetOutputFormat("json")
	
	// Test that JSON output is valid and properly indented
	stats := &AggregatedStatistics{
		Period: TimePeriod{
			Start: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			End:   time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
		},
		TotalSessions: 5,
	}
	
	result := formatter.FormatStats(stats)
	
	// Check valid JSON
	var parsed interface{}
	if err := json.Unmarshal([]byte(result), &parsed); err != nil {
		t.Errorf("Invalid JSON output: %v", err)
	}
	
	// Check indentation (should have newlines and spaces)
	if !strings.Contains(result, "\n") {
		t.Error("JSON output not properly formatted with newlines")
	}
	if !strings.Contains(result, "  ") {
		t.Error("JSON output not properly indented")
	}
}

func TestEdgeCases(t *testing.T) {
	formatter := NewDisplayFormatter()
	
	// Test with zero values
	stats := &AggregatedStatistics{
		Period: TimePeriod{},
		PerformanceMetrics: PerformanceMetrics{
			AvgResponseTimeMs: 0,
		},
	}
	
	// Should not panic
	result := formatter.FormatStats(stats)
	if result == "" {
		t.Error("Should return formatted output even with zero values")
	}
	
	// Test with very large numbers
	session := &SessionStatistics{
		SessionID:       "test",
		TotalToolCalls:  999999999,
		SuccessfulCalls: 888888888,
		TotalTokens: TokenUsage{
			Total: 999999999,
		},
	}
	
	result = formatter.FormatSession(session)
	if !strings.Contains(result, "999999999") {
		t.Error("Should handle large numbers correctly")
	}
	
	// Test with special characters in strings
	patterns := []ErrorPattern{
		{
			Pattern: "Error with \"quotes\" and 'apostrophes'",
			Count:   1,
			Tools:   []string{"Tool,with,commas"},
		},
	}
	
	formatter.SetOutputFormat("csv")
	result = formatter.FormatErrors(patterns, nil)
	if result == "" {
		t.Error("Should handle special characters in CSV format")
	}
}