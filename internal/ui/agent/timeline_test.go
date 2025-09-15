package agent

import (
	"testing"
	"time"
)

func TestMessageTimelineComponent(t *testing.T) {
	// Test timeline component creation
	timeline := NewMessageTimelineComponent()
	if timeline == nil {
		t.Fatal("NewMessageTimelineComponent returned nil")
	}

	// Test initial state
	if timeline.MaxWidth != 120 {
		t.Errorf("Expected MaxWidth to be 120, got %d", timeline.MaxWidth)
	}

	if timeline.MaxHeight != 10 {
		t.Errorf("Expected MaxHeight to be 10, got %d", timeline.MaxHeight)
	}

	if timeline.NavigationState.ZoomLevel != 1 {
		t.Errorf("Expected initial zoom level to be 1, got %d", timeline.NavigationState.ZoomLevel)
	}
}

func TestMessageEventCreation(t *testing.T) {
	// Test MessageEvent creation
	event := MessageEvent{
		ID:          "test-001",
		Timestamp:   time.Now(),
		SessionID:   "session1",
		Role:        UserMessage,
		Content:     "Test message",
		TokenCount:  25,
		AgentUsed:   "TestAgent",
		ToolsUsed:   []string{"Read", "Write"},
		Success:     true,
		Duration:    1000 * time.Millisecond,
	}

	// Validate event structure
	if event.ID == "" {
		t.Error("MessageEvent ID should not be empty")
	}

	if event.Role != UserMessage {
		t.Errorf("Expected role to be UserMessage, got %s", event.Role)
	}

	if event.Success != true {
		t.Error("Expected success to be true")
	}

	if len(event.ToolsUsed) != 2 {
		t.Errorf("Expected 2 tools used, got %d", len(event.ToolsUsed))
	}
}

func TestTimelineEventManagement(t *testing.T) {
	timeline := NewMessageTimelineComponent()

	// Create test events
	events := []MessageEvent{
		{
			ID:        "msg-001",
			Timestamp: time.Now().Add(-5 * time.Minute),
			SessionID: "session1",
			Role:      UserMessage,
			Content:   "First message",
			Success:   true,
		},
		{
			ID:        "msg-002",
			Timestamp: time.Now().Add(-3 * time.Minute),
			SessionID: "session1",
			Role:      AssistantMessage,
			Content:   "Second message",
			Success:   true,
		},
	}

	// Set events
	timeline.SetEvents(events)

	// Verify events were set
	if len(timeline.Events) != 2 {
		t.Errorf("Expected 2 events, got %d", len(timeline.Events))
	}

	// Test event count
	count := timeline.GetEventCount()
	if count != 2 {
		t.Errorf("Expected event count to be 2, got %d", count)
	}
}

func TestTimelineNavigationBasic(t *testing.T) {
	timeline := NewMessageTimelineComponent()

	// Test zoom functionality
	initialZoom := timeline.NavigationState.ZoomLevel
	timeline.ZoomTimeline(1) // Zoom out
	if timeline.NavigationState.ZoomLevel <= initialZoom {
		t.Error("Expected zoom level to increase when zooming out")
	}

	timeline.ZoomTimeline(-1) // Zoom in
	if timeline.NavigationState.ZoomLevel != initialZoom {
		t.Error("Expected zoom level to return to initial after zoom in")
	}

	// Test time navigation
	initialStart := timeline.NavigationState.CurrentRange.Start
	timeline.NavigateTime(1) // Move forward
	if timeline.NavigationState.CurrentRange.Start.Equal(initialStart) {
		t.Error("Expected time range to change after navigation")
	}

	// Test reset view
	timeline.ResetView()
	if timeline.NavigationState.ZoomLevel != 1 {
		t.Error("Expected zoom level to reset to 1")
	}

	if timeline.NavigationState.SelectedEventID != "" {
		t.Error("Expected selected event ID to be cleared after reset")
	}
}

func TestTimelineRendering(t *testing.T) {
	timeline := NewMessageTimelineComponent()

	// Create test events
	events := []MessageEvent{
		{
			ID:        "msg-001",
			Timestamp: time.Now().Add(-5 * time.Minute),
			SessionID: "session1",
			Role:      UserMessage,
			Content:   "Test message",
			Success:   true,
		},
	}

	timeline.SetEvents(events)

	// Test rendering
	rendered := timeline.RenderTimeline(true, "Test Timeline")
	if rendered == "" {
		t.Error("Expected non-empty rendered timeline")
	}

	// Should contain title
	if !timelineContainsString(rendered, "Test Timeline") {
		t.Error("Expected rendered timeline to contain title")
	}

	// Test rendering when focused
	renderedFocused := timeline.RenderTimeline(true, "Focused Timeline")
	if renderedFocused == "" {
		t.Error("Expected non-empty rendered timeline when focused")
	}

	// Test rendering when not focused
	renderedNotFocused := timeline.RenderTimeline(false, "Unfocused Timeline")
	if renderedNotFocused == "" {
		t.Error("Expected non-empty rendered timeline when not focused")
	}
}

// Helper function for string containment check
func timelineContainsString(s, substr string) bool {
	return len(s) >= len(substr) &&
		   (len(substr) == 0 || stringContains(s, substr))
}

func stringContains(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	if len(s) < len(substr) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}