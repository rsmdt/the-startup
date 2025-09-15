package agent

import (
	"fmt"
	"testing"
	"time"
)

// TestTimelineOutputVisual is a visual test to see actual graph output
func TestTimelineOutputVisual(t *testing.T) {
	// Create timeline component
	timeline := NewMessageTimelineComponent()
	timeline.MaxWidth = 100
	timeline.MaxHeight = 10

	// Add test events simulating a conversation
	now := time.Now()
	events := []MessageEvent{
		{
			ID:        "1",
			Timestamp: now.Add(-3 * time.Hour),
			Role:      UserMessage,
			Content:   "Can you help me implement a stats dashboard?",
		},
		{
			ID:        "2",
			Timestamp: now.Add(-2*time.Hour - 55*time.Minute),
			Role:      AssistantMessage,
			Content:   "I'll help you implement the stats dashboard...",
		},
		{
			ID:        "3",
			Timestamp: now.Add(-2*time.Hour - 30*time.Minute),
			Role:      UserMessage,
			Content:   "The layout seems incorrect",
		},
		{
			ID:        "4",
			Timestamp: now.Add(-2*time.Hour - 25*time.Minute),
			Role:      AssistantMessage,
			Content:   "Let me fix that layout issue...",
		},
		{
			ID:        "5",
			Timestamp: now.Add(-1 * time.Hour),
			Role:      UserMessage,
			Content:   "Now I need a timeline graph",
		},
		{
			ID:        "6",
			Timestamp: now.Add(-55 * time.Minute),
			Role:      AssistantMessage,
			Content:   "I'll create a btop-style graph...",
		},
		{
			ID:        "7",
			Timestamp: now.Add(-30 * time.Minute),
			Role:      UserMessage,
			Content:   "The graph is not visible",
		},
		{
			ID:        "8",
			Timestamp: now.Add(-25 * time.Minute),
			Role:      AssistantMessage,
			Content:   "Let me fix the graph rendering...",
		},
	}

	timeline.SetEvents(events)

	// Render timeline
	output := timeline.RenderTimeline(true, "📊 Message Timeline")

	// Print for visual inspection
	fmt.Println("\n=== TIMELINE GRAPH OUTPUT ===")
	fmt.Println(output)
	fmt.Println("=== END OF OUTPUT ===")

	// Basic checks
	if output == "" {
		t.Error("Timeline output is empty")
	}
}