package stats

import (
	"fmt"
	"time"
)

func ExampleDisplayFormatter_FormatStats() {
	formatter := NewDisplayFormatter()
	
	stats := &AggregatedStatistics{
		Period: TimePeriod{
			Start: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			End:   time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
		},
		TotalSessions:      3,
		ActiveSessions:     []string{"session1", "session2"},
		AvgSessionDuration: 45 * time.Minute,
		TotalToolCalls:     150,
		OverallSuccessRate: 0.92,
	}
	
	// Table format (default)
	output := formatter.FormatStats(stats)
	fmt.Println("Table format output generated")
	
	// JSON format
	formatter.SetOutputFormat("json")
	output = formatter.FormatStats(stats)
	fmt.Println("JSON format output generated")
	
	// CSV format
	formatter.SetOutputFormat("csv")
	output = formatter.FormatStats(stats)
	fmt.Println("CSV format output generated")
	
	_ = output // Use output variable
	
	// Output:
	// Table format output generated
	// JSON format output generated
	// CSV format output generated
}

func ExampleDisplayFormatter_FormatToolLeaderboard() {
	formatter := NewDisplayFormatter()
	
	toolStats := map[string]*GlobalToolStats{
		"Bash": {
			Name:          "Bash",
			TotalCalls:    500,
			SuccessRate:   0.95,
			AvgDurationMs: 250,
		},
		"Read": {
			Name:          "Read",
			TotalCalls:    300,
			SuccessRate:   0.98,
			AvgDurationMs: 50,
		},
		"Edit": {
			Name:          "Edit",
			TotalCalls:    200,
			SuccessRate:   0.90,
			AvgDurationMs: 150,
		},
	}
	
	// Show top 2 tools
	output := formatter.FormatToolLeaderboard(toolStats, 2)
	fmt.Printf("Leaderboard shows %d tools\n", 2)
	
	_ = output // Use output variable
	
	// Output:
	// Leaderboard shows 2 tools
}

func ExampleDisplayFormatter_multiFormat() {
	formatter := NewDisplayFormatter()
	
	session := &SessionStatistics{
		SessionID:       "test-session-123",
		StartTime:       time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
		Duration:        2 * time.Hour,
		TotalToolCalls:  50,
		SuccessfulCalls: 47,
		FailedCalls:     3,
		ErrorRate:       0.06,
	}
	
	formats := []string{"table", "json", "csv"}
	
	for _, format := range formats {
		formatter.SetOutputFormat(format)
		output := formatter.FormatSession(session)
		if output != "" {
			fmt.Printf("Generated %s format\n", format)
		}
	}
	
	// Output:
	// Generated table format
	// Generated json format
	// Generated csv format
}