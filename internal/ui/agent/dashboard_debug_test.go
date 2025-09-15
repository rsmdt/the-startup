package agent

import (
	"fmt"
	"testing"
	"github.com/rsmdt/the-startup/internal/stats"
)

func TestDashboardDataLoading(t *testing.T) {
	// Load some real log entries
	discovery := stats.NewLogDiscovery()
	projectPath := discovery.GetCurrentProject()
	if projectPath == "" {
		projectPath = "/Users/irudi/Code/personal/the-startup"
	}

	logFiles, err := discovery.FindLogFiles(projectPath, stats.FilterOptions{})
	if err != nil {
		t.Skipf("Could not find log files: %v", err)
	}

	if len(logFiles) == 0 {
		t.Skip("No log files found")
	}

	// Find a log file with sufficient data
	var selectedFile string
	for _, file := range logFiles {
		selectedFile = file
		// Try to use a larger file if available
		if len(file) > 100 {
			break
		}
	}

	// Parse log file
	parser := stats.NewJSONLParser()
	var entries []stats.ClaudeLogEntry

	// Process entries
	entryChan, errorChan := parser.ParseFile(selectedFile)

	count := 0
	for entry := range entryChan {
		entries = append(entries, entry)
		count++
		if count >= 100 {
			break
		}
	}

	// Drain error channel
	go func() {
		for range errorChan {
		}
	}()

	fmt.Printf("Loaded %d entries from %s\n", len(entries), selectedFile)

	// Show entry types
	typeCounts := make(map[string]int)
	for _, entry := range entries {
		typeCounts[entry.Type]++
	}
	fmt.Printf("Entry types: %+v\n", typeCounts)

	// Extract message events
	events := extractMessageEvents(entries)
	fmt.Printf("Created %d message events from %d entries\n", len(events), len(entries))

	// Count roles
	roleCounts := make(map[MessageRole]int)
	for _, event := range events {
		roleCounts[event.Role]++
	}
	fmt.Printf("Event roles: User=%d, Assistant=%d\n", roleCounts[UserMessage], roleCounts[AssistantMessage])

	// Create timeline and test rendering
	if len(events) > 0 {
		timeline := NewMessageTimelineComponent()
		timeline.MaxWidth = 100
		timeline.MaxHeight = 10
		timeline.SetEvents(events)

		output := timeline.RenderTimeline(true, "Test Timeline")
		fmt.Printf("\nTimeline Output:\n%s\n", output)

		// Check if graph appears
		if !contains(output, "User") || !contains(output, "Asst") {
			t.Error("Timeline graph not rendering properly")
		}
	}
}

func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || contains(s[1:], substr)))
}