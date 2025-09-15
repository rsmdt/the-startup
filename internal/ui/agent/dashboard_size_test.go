package agent

import (
	"fmt"
	"strings"
	"testing"
	"time"
	"github.com/rsmdt/the-startup/internal/stats"
)

func TestDashboardRenderingWithSizes(t *testing.T) {
	// Create test data
	now := time.Now()
	events := []MessageEvent{
		{ID: "1", Timestamp: now.Add(-3 * time.Hour), Role: UserMessage, Content: "Test user 1"},
		{ID: "2", Timestamp: now.Add(-2 * time.Hour), Role: AssistantMessage, Content: "Test assistant 1"},
		{ID: "3", Timestamp: now.Add(-1 * time.Hour), Role: UserMessage, Content: "Test user 2"},
		{ID: "4", Timestamp: now.Add(-30 * time.Minute), Role: AssistantMessage, Content: "Test assistant 2"},
	}

	agentStats := map[string]*stats.GlobalAgentStats{
		"the-chief": {
			Count:           100,
			SuccessCount:    90,
			FailureCount:    10,
			TotalDurationMs: 50000,
		},
	}

	// Test different terminal sizes
	sizes := []struct {
		name   string
		width  int
		height int
	}{
		{"Small Terminal", 80, 24},
		{"Medium Terminal", 120, 40},
		{"Large Terminal", 160, 60},
	}

	for _, size := range sizes {
		t.Run(size.name, func(t *testing.T) {
			// Create dashboard
			dashboard := NewDashboardModel(agentStats, nil, nil)
			dashboard.width = size.width
			dashboard.height = size.height
			dashboard.timeline.SetEvents(events)

			// Render
			output := dashboard.View()

			// Count lines
			lines := strings.Split(output, "\n")
			lineCount := len(lines)

			// Check timeline presence
			hasTimeline := strings.Contains(output, "Timeline") || strings.Contains(output, "User") || strings.Contains(output, "Asst")

			fmt.Printf("\n=== %s (%dx%d) ===\n", size.name, size.width, size.height)
			fmt.Printf("Output lines: %d\n", lineCount)
			fmt.Printf("Has timeline: %v\n", hasTimeline)

			// Calculate what heights should be
			availableHeight := size.height - 3
			timelineHeight := int(float64(availableHeight) * 0.25)
			if timelineHeight < 8 {
				timelineHeight = 8
			}
			fmt.Printf("Calculated timeline height: %d\n", timelineHeight)

			// Show first 10 lines to see what's rendering
			for i := 0; i < 10 && i < len(lines); i++ {
				if len(lines[i]) > 60 {
					fmt.Printf("Line %d: %s...\n", i, lines[i][:60])
				} else {
					fmt.Printf("Line %d: %s\n", i, lines[i])
				}
			}
		})
	}
}