package agent

import (
	"strings"
	"testing"
	"time"
)

// TestTimelineGraphRendering tests the btop-style graph rendering
func TestTimelineGraphRendering(t *testing.T) {
	// Create timeline component
	timeline := NewMessageTimelineComponent()
	timeline.MaxWidth = 100
	timeline.MaxHeight = 12

	// Add test events
	now := time.Now()
	events := []MessageEvent{
		{
			ID:        "1",
			Timestamp: now.Add(-2 * time.Hour),
			Role:      UserMessage,
			Content:   "Test user message 1",
		},
		{
			ID:        "2",
			Timestamp: now.Add(-90 * time.Minute),
			Role:      AssistantMessage,
			Content:   "Test assistant response 1",
		},
		{
			ID:        "3",
			Timestamp: now.Add(-1 * time.Hour),
			Role:      UserMessage,
			Content:   "Test user message 2",
		},
		{
			ID:        "4",
			Timestamp: now.Add(-30 * time.Minute),
			Role:      AssistantMessage,
			Content:   "Test assistant response 2",
		},
		{
			ID:        "5",
			Timestamp: now.Add(-15 * time.Minute),
			Role:      UserMessage,
			Content:   "Test user message 3",
		},
		{
			ID:        "6",
			Timestamp: now.Add(-5 * time.Minute),
			Role:      AssistantMessage,
			Content:   "Test assistant response 3",
		},
	}

	timeline.SetEvents(events)

	// Render timeline
	output := timeline.RenderTimeline(true, "Message Timeline Test")

	// Check for key elements of horizontal btop-style graph
	tests := []struct {
		name     string
		expected string
	}{
		{"Has Y-axis", "│"},
		{"Has X-axis", "└"},
		{"Has axis lines", "─"},
		{"Has user label", "User"},
		{"Has assistant label", "Asst"},
		{"Has msg/hr indicator", "msg/hr"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !strings.Contains(output, tt.expected) {
				t.Errorf("Graph output missing %s: expected to contain '%s'\nGot:\n%s",
					tt.name, tt.expected, output)
			}
		})
	}
}

// TestBrailleCharacterUsage tests that braille characters are used for smooth graphs
func TestBrailleCharacterUsage(t *testing.T) {
	timeline := NewMessageTimelineComponent()

	// Process some test data
	timeline.graphData = &TimelineGraphData{
		UserActivity:      []float64{0.2, 0.4, 0.6, 0.8, 1.0, 0.8, 0.6, 0.4, 0.2},
		AssistantActivity: []float64{0.1, 0.3, 0.5, 0.7, 0.9, 0.7, 0.5, 0.3, 0.1},
		TimeLabels:        []string{"00:00", "", "", "06:00", "", "", "12:00", "", ""},
		MaxValue:          10,
		BucketDuration:    time.Hour,
		StartTime:         time.Now().Add(-24 * time.Hour),
		EndTime:           time.Now(),
	}

	// Render the graph
	graphOutput := timeline.renderBtopGraph()

	// Check for braille characters
	brailleFound := false
	for _, char := range brailleChars {
		if char != ' ' && strings.ContainsRune(graphOutput, char) {
			brailleFound = true
			break
		}
	}

	if !brailleFound {
		t.Error("Expected graph to use braille characters for smooth visualization")
	}
}

// TestGraphDataProcessing tests the graph data processing logic
func TestGraphDataProcessing(t *testing.T) {
	timeline := NewMessageTimelineComponent()

	// Create events with known distribution
	now := time.Now()
	events := make([]MessageEvent, 0)

	// Create burst of user messages
	for i := 0; i < 5; i++ {
		events = append(events, MessageEvent{
			ID:        string(rune('a' + i)),
			Timestamp: now.Add(time.Duration(-120+i) * time.Minute),
			Role:      UserMessage,
		})
	}

	// Create burst of assistant messages
	for i := 0; i < 3; i++ {
		events = append(events, MessageEvent{
			ID:        string(rune('x' + i)),
			Timestamp: now.Add(time.Duration(-60+i) * time.Minute),
			Role:      AssistantMessage,
		})
	}

	timeline.SetEvents(events)

	// Verify graph data was processed
	if timeline.graphData == nil {
		t.Fatal("Expected graph data to be processed")
	}

	if timeline.graphData.MaxValue == 0 {
		t.Error("Expected MaxValue to be greater than 0")
	}

	if len(timeline.graphData.UserActivity) == 0 {
		t.Error("Expected UserActivity data to be populated")
	}

	if len(timeline.graphData.AssistantActivity) == 0 {
		t.Error("Expected AssistantActivity data to be populated")
	}
}

// TestTimelineNavigation tests navigation in btop-style graph
func TestTimelineNavigation(t *testing.T) {
	timeline := NewMessageTimelineComponent()

	// Add some events
	now := time.Now()
	events := []MessageEvent{
		{ID: "1", Timestamp: now.Add(-2 * time.Hour), Role: UserMessage},
		{ID: "2", Timestamp: now.Add(-1 * time.Hour), Role: AssistantMessage},
		{ID: "3", Timestamp: now, Role: UserMessage},
	}
	timeline.SetEvents(events)

	// Test zoom
	initialWidth := timeline.graphWidth
	timeline.ZoomTimeline(1) // Zoom out
	if timeline.graphWidth >= initialWidth {
		t.Error("Expected graph width to decrease when zooming out")
	}

	timeline.ZoomTimeline(-1) // Zoom in
	timeline.ZoomTimeline(-1) // Zoom in more
	if timeline.graphWidth <= initialWidth {
		t.Error("Expected graph width to increase when zooming in")
	}

	// Test reset
	timeline.ResetView()
	if timeline.scrollOffset != 0 {
		t.Error("Expected scroll offset to be reset to 0")
	}
	if timeline.graphWidth != 80 {
		t.Errorf("Expected graph width to be reset to 80, got %d", timeline.graphWidth)
	}
}

// TestGraphInterpolation tests data interpolation for smooth curves
func TestGraphInterpolation(t *testing.T) {
	timeline := NewMessageTimelineComponent()

	// Test data interpolation
	data := []float64{0.0, 0.5, 1.0}
	interpolated := timeline.sampleData(data, 5)

	if len(interpolated) != 5 {
		t.Errorf("Expected interpolated data to have 5 points, got %d", len(interpolated))
	}

	// Check that interpolation creates smooth transitions
	if interpolated[0] > interpolated[1] || interpolated[1] > interpolated[2] {
		t.Error("Expected interpolated data to have smooth transitions")
	}
}

// TestHorizontalGraphFormat tests that the graph renders horizontally with user on top and assistant on bottom
func TestHorizontalGraphFormat(t *testing.T) {
	timeline := NewMessageTimelineComponent()
	timeline.MaxWidth = 100
	timeline.MaxHeight = 12

	// Add test events
	now := time.Now()
	events := []MessageEvent{
		{ID: "1", Timestamp: now.Add(-1 * time.Hour), Role: UserMessage, Content: "User msg"},
		{ID: "2", Timestamp: now.Add(-30 * time.Minute), Role: AssistantMessage, Content: "Assistant msg"},
		{ID: "3", Timestamp: now, Role: UserMessage, Content: "Another user msg"},
	}
	timeline.SetEvents(events)

	// Process and render the graph
	timeline.processGraphData()
	graphOutput := timeline.renderBtopGraph()

	// Split into lines
	lines := strings.Split(graphOutput, "\n")

	// Check that we have at least 4 lines (User, Asst, X-axis, time markers)
	if len(lines) < 4 {
		t.Fatalf("Expected at least 4 lines in graph output, got %d", len(lines))
	}

	// Check that first line starts with "User"
	if !strings.HasPrefix(lines[0], "User") {
		t.Errorf("First line should start with 'User', got: %s", lines[0])
	}

	// Check that second line starts with "Asst"
	if !strings.HasPrefix(lines[1], "Asst") {
		t.Errorf("Second line should start with 'Asst', got: %s", lines[1])
	}

	// Check that third line has X-axis
	if !strings.Contains(lines[2], "└") || !strings.Contains(lines[2], "─") {
		t.Errorf("Third line should contain X-axis with └ and ─, got: %s", lines[2])
	}

	// Check for msg/hr indicators
	if !strings.Contains(lines[0], "msg/hr") {
		t.Error("User line should contain 'msg/hr' indicator")
	}
	if !strings.Contains(lines[1], "msg/hr") {
		t.Error("Assistant line should contain 'msg/hr' indicator")
	}
}

// TestTimeMarkersGeneration tests time marker generation for x-axis
func TestTimeMarkersGeneration(t *testing.T) {
	timeline := NewMessageTimelineComponent()

	// Test label generation for different time ranges
	start := time.Now().Add(-24 * time.Hour)
	end := time.Now()
	labels := timeline.generateTimeLabels(start, end, 24, time.Hour)

	if len(labels) != 24 {
		t.Errorf("Expected 24 labels for 24-hour period, got %d", len(labels))
	}

	// Check that some labels are populated
	hasLabels := false
	for _, label := range labels {
		if label != "" {
			hasLabels = true
			break
		}
	}

	if !hasLabels {
		t.Error("Expected at least some time labels to be generated")
	}
}