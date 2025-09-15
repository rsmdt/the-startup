package agent

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// MessageEvent represents a single message event in the timeline
type MessageEvent struct {
	ID          string
	Timestamp   time.Time
	SessionID   string
	Role        MessageRole
	Content     string
	TokenCount  int
	AgentUsed   string
	ToolsUsed   []string
	Success     bool
	Duration    time.Duration
}

// MessageRole defines the role of a message
type MessageRole string

const (
	UserMessage      MessageRole = "user"
	AssistantMessage MessageRole = "assistant"
	ToolMessage      MessageRole = "tool"
	SystemMessage    MessageRole = "system"
)

// TimelineRange represents a time range for filtering
type TimelineRange struct {
	Start    time.Time
	End      time.Time
	Duration time.Duration
}

// TimelineNavigationState tracks navigation state
type TimelineNavigationState struct {
	CurrentRange    TimelineRange
	ZoomLevel       int // 0=hour, 1=day, 2=week, 3=month
	SelectedEventID string
	ViewportOffset  int
}

// TimelineComponentStyle defines styling for timeline components
type TimelineComponentStyle struct {
	Title           lipgloss.Style
	FocusedTitle    lipgloss.Style
	TimeAxis        lipgloss.Style
	UserMessage     lipgloss.Style
	AssistantMessage lipgloss.Style
	ToolMessage     lipgloss.Style
	SystemMessage   lipgloss.Style
	SuccessMessage  lipgloss.Style
	FailureMessage  lipgloss.Style
	SelectedMessage lipgloss.Style
	TimeMarker      lipgloss.Style
	NavigationHelp  lipgloss.Style
	EmptyState      lipgloss.Style
	AgentIndicator  lipgloss.Style
}

// DefaultTimelineStyle creates professional styling for timeline components with visual consistency
func DefaultTimelineStyle() TimelineComponentStyle {
	return TimelineComponentStyle{
		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("99")).
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(lipgloss.Color("99")).
			Padding(0, 0, 1, 0),
		FocusedTitle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("46")).
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(lipgloss.Color("46")).
			Padding(0, 0, 1, 0),
		TimeAxis: lipgloss.NewStyle().
			Foreground(lipgloss.Color("99")).
			Bold(true).
			Background(lipgloss.Color("235")).
			Padding(0, 1),
		UserMessage: lipgloss.NewStyle().
			Foreground(lipgloss.Color("39")).  // Blue for user
			Background(lipgloss.Color("17")).
			Bold(true),
		AssistantMessage: lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")). // Pink for assistant
			Background(lipgloss.Color("53")).
			Bold(true),
		ToolMessage: lipgloss.NewStyle().
			Foreground(lipgloss.Color("220")). // Yellow for tools
			Background(lipgloss.Color("94")).
			Bold(true),
		SystemMessage: lipgloss.NewStyle().
			Foreground(lipgloss.Color("245")). // Gray for system
			Background(lipgloss.Color("235")).
			Bold(true),
		SuccessMessage: lipgloss.NewStyle().
			Foreground(lipgloss.Color("46")). // Green for success
			Bold(true),
		FailureMessage: lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")). // Red for failure
			Bold(true),
		SelectedMessage: lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Background(lipgloss.Color("57")).
			Bold(true).
			Blink(true),
		TimeMarker: lipgloss.NewStyle().
			Foreground(lipgloss.Color("99")).
			Faint(true),
		NavigationHelp: lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Italic(true).
			Border(lipgloss.NormalBorder(), true, false, false, false).
			BorderForeground(lipgloss.Color("99")).
			Padding(1, 0, 0, 0),
		EmptyState: lipgloss.NewStyle().
			Italic(true).
			Foreground(lipgloss.Color("240")).
			Align(lipgloss.Center).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(2, 4).
			Margin(1, 0),
		AgentIndicator: lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Bold(true).
			Background(lipgloss.Color("235")).
			Padding(0, 1),
	}
}

// MessageTimelineComponent provides horizontal timeline visualization
type MessageTimelineComponent struct {
	Style           TimelineComponentStyle
	MaxWidth        int
	MaxHeight       int
	NavigationState TimelineNavigationState
	Events          []MessageEvent
	currentTheme    ColorTheme

	// Graph data for btop-style visualization
	graphData       *TimelineGraphData
	graphWidth      int
	graphHeight     int
	scrollOffset    int
}

// TimelineGraphData stores processed data for graph rendering
type TimelineGraphData struct {
	UserActivity      []float64 // Normalized activity levels (0-1)
	AssistantActivity []float64 // Normalized activity levels (0-1)
	TimeLabels        []string  // Time labels for x-axis
	MaxValue          int       // Maximum message count in any bucket
	BucketDuration    time.Duration
	StartTime         time.Time
	EndTime           time.Time
}

// NewMessageTimelineComponent creates a new timeline component
func NewMessageTimelineComponent() *MessageTimelineComponent {
	now := time.Now()
	return &MessageTimelineComponent{
		Style:    DefaultTimelineStyle(),
		MaxWidth: 120,
		MaxHeight: 10,
		NavigationState: TimelineNavigationState{
			CurrentRange: TimelineRange{
				Start:    now.Add(-24 * time.Hour),
				End:      now,
				Duration: 24 * time.Hour,
			},
			ZoomLevel:       1, // Start at day level
			SelectedEventID: "",
			ViewportOffset:  0,
		},
		Events:       []MessageEvent{},
		graphWidth:   80,
		graphHeight:  8,
		scrollOffset: 0,
	}
}

// SetEvents updates the timeline events
func (tc *MessageTimelineComponent) SetEvents(events []MessageEvent) {
	tc.Events = events
	tc.sortEventsByTime()
	tc.processGraphData()
}

// ApplyTheme applies a color theme to the timeline component
func (tc *MessageTimelineComponent) ApplyTheme(theme ColorTheme) {
	tc.currentTheme = theme
	tc.Style = TimelineComponentStyle{
		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(theme.Primary).
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(theme.Border).
			Padding(0, 0, 1, 0),
		FocusedTitle: lipgloss.NewStyle().
			Bold(true).
			Foreground(theme.Highlight).
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(theme.Highlight).
			Padding(0, 0, 1, 0),
		TimeAxis: lipgloss.NewStyle().
			Foreground(theme.Info).
			Bold(true).
			Padding(0, 1),
		UserMessage: lipgloss.NewStyle().
			Foreground(theme.GraphUser).
			Bold(true),
		AssistantMessage: lipgloss.NewStyle().
			Foreground(theme.GraphAsst).
			Bold(true),
		ToolMessage: lipgloss.NewStyle().
			Foreground(theme.Warning).
			Bold(true),
		SystemMessage: lipgloss.NewStyle().
			Foreground(theme.Muted).
			Bold(true),
		SuccessMessage: lipgloss.NewStyle().
			Foreground(theme.Success).
			Bold(true),
		FailureMessage: lipgloss.NewStyle().
			Foreground(theme.Error).
			Bold(true),
		SelectedMessage: lipgloss.NewStyle().
			Foreground(theme.Foreground).
			Background(theme.Secondary).
			Bold(true),
		TimeMarker: lipgloss.NewStyle().
			Foreground(theme.Muted).
			Faint(true),
		NavigationHelp: lipgloss.NewStyle().
			Foreground(theme.Muted).
			Italic(true).
			Border(lipgloss.NormalBorder(), true, false, false, false).
			BorderForeground(theme.Border).
			Padding(1, 0, 0, 0),
		EmptyState: lipgloss.NewStyle().
			Italic(true).
			Foreground(theme.Muted).
			Align(lipgloss.Center).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(theme.Border).
			Padding(2, 4).
			Margin(1, 0),
		AgentIndicator: lipgloss.NewStyle().
			Foreground(theme.Secondary).
			Bold(true).
			Padding(0, 1),
	}
}

// sortEventsByTime sorts events chronologically
func (tc *MessageTimelineComponent) sortEventsByTime() {
	sort.Slice(tc.Events, func(i, j int) bool {
		return tc.Events[i].Timestamp.Before(tc.Events[j].Timestamp)
	})
}

// RenderTimelineContent renders just the timeline content without title (for btop-style borders)
func (tc *MessageTimelineComponent) RenderTimelineContent(isFocused bool) string {
	var content strings.Builder

	// Process graph data if not already done
	if tc.graphData == nil {
		tc.processGraphData()
	}

	// Check if we have events
	if len(tc.Events) == 0 {
		content.WriteString(tc.Style.EmptyState.Render("No message events found") + "\n")
		return content.String()
	}

	// Render the btop-style graph
	content.WriteString(tc.renderBtopGraph())

	return content.String()
}

// RenderTimeline renders the horizontal timeline view (DEPRECATED - use RenderTimelineContent for btop-style)
func (tc *MessageTimelineComponent) RenderTimeline(
	isFocused bool,
	title string,
) string {
	var content strings.Builder

	// Title with graph indicator
	titleStyle := tc.Style.Title
	if isFocused {
		titleStyle = tc.Style.FocusedTitle
	}
	content.WriteString(titleStyle.Render(title) + "\n")

	// Check if we have events
	if len(tc.Events) == 0 {
		content.WriteString("\n" + tc.Style.EmptyState.Render("No message events found") + "\n")
		return content.String()
	}

	// Process graph data if not already done
	if tc.graphData == nil {
		tc.processGraphData()
	}

	// Render btop-style graph
	graphContent := tc.renderBtopGraph()
	content.WriteString(graphContent)

	// Render navigation help inline
	tc.renderNavigationHelp(&content)

	return content.String()
}

// filterEventsInRange returns events within the current time range
func (tc *MessageTimelineComponent) filterEventsInRange() []MessageEvent {
	var filtered []MessageEvent

	for _, event := range tc.Events {
		if event.Timestamp.After(tc.NavigationState.CurrentRange.Start) &&
		   event.Timestamp.Before(tc.NavigationState.CurrentRange.End) {
			filtered = append(filtered, event)
		}
	}

	return filtered
}

// renderTimeAxisHeader creates the time axis with markers
func (tc *MessageTimelineComponent) renderTimeAxisHeader() string {
	var header strings.Builder

	// Create time markers across the width
	timelineWidth := tc.MaxWidth - 20 // Leave space for labels
	rangeStart := tc.NavigationState.CurrentRange.Start
	rangeEnd := tc.NavigationState.CurrentRange.End
	rangeDuration := rangeEnd.Sub(rangeStart)

	// Determine time marker spacing based on zoom level
	var markerInterval time.Duration
	var timeFormat string

	switch tc.NavigationState.ZoomLevel {
	case 0: // Hour view
		markerInterval = 10 * time.Minute
		timeFormat = "15:04"
	case 1: // Day view
		markerInterval = 2 * time.Hour
		timeFormat = "15:04"
	case 2: // Week view
		markerInterval = 12 * time.Hour
		timeFormat = "Mon 15:04"
	case 3: // Month view
		markerInterval = 3 * 24 * time.Hour
		timeFormat = "Jan 2"
	default:
		markerInterval = 2 * time.Hour
		timeFormat = "15:04"
	}

	// Create time axis line
	axisLine := strings.Repeat("─", timelineWidth)
	header.WriteString(tc.Style.TimeAxis.Render("Time: " + axisLine) + "\n")

	// Add time markers
	var markerLine strings.Builder
	markerLine.WriteString("      ") // Indent to align with axis

	currentTime := rangeStart
	for currentTime.Before(rangeEnd) {
		position := int(float64(currentTime.Sub(rangeStart)) / float64(rangeDuration) * float64(timelineWidth))
		if position < timelineWidth {
			// Add marker at position
			timeLabel := currentTime.Format(timeFormat)
			marker := "┬" + timeLabel

			// Add spacing to reach position
			for markerLine.Len()-6 < position {
				markerLine.WriteRune(' ')
			}
			markerLine.WriteString(marker)
		}
		currentTime = currentTime.Add(markerInterval)
	}

	header.WriteString(tc.Style.TimeMarker.Render(markerLine.String()))

	return header.String()
}

// renderTimelineBars creates horizontal bars showing message events
func (tc *MessageTimelineComponent) renderTimelineBars(events []MessageEvent) string {
	var content strings.Builder

	// Group events by role for separate timeline tracks
	roleEvents := map[MessageRole][]MessageEvent{
		UserMessage:      {},
		AssistantMessage: {},
		ToolMessage:      {},
		SystemMessage:    {},
	}

	for _, event := range events {
		roleEvents[event.Role] = append(roleEvents[event.Role], event)
	}

	// Render timeline track for each role
	roles := []MessageRole{UserMessage, AssistantMessage, ToolMessage, SystemMessage}
	roleLabels := map[MessageRole]string{
		UserMessage:      "User    ",
		AssistantMessage: "Assistant",
		ToolMessage:      "Tools   ",
		SystemMessage:    "System  ",
	}

	for _, role := range roles {
		if len(roleEvents[role]) == 0 {
			continue // Skip empty tracks
		}

		content.WriteString("\n")
		trackLine := tc.renderTimelineTrack(role, roleEvents[role], roleLabels[role])
		content.WriteString(trackLine)
	}

	return content.String()
}

// renderTimelineTrack renders a single timeline track for a message role
func (tc *MessageTimelineComponent) renderTimelineTrack(role MessageRole, events []MessageEvent, label string) string {
	var track strings.Builder

	timelineWidth := tc.MaxWidth - 20
	rangeStart := tc.NavigationState.CurrentRange.Start
	rangeEnd := tc.NavigationState.CurrentRange.End
	rangeDuration := rangeEnd.Sub(rangeStart)

	// Create base timeline
	timeline := make([]rune, timelineWidth)
	for i := range timeline {
		timeline[i] = '·'
	}

	// Add events to timeline
	for _, event := range events {
		position := int(float64(event.Timestamp.Sub(rangeStart)) / float64(rangeDuration) * float64(timelineWidth))
		if position >= 0 && position < timelineWidth {
			// Choose symbol based on success/failure and selection
			var symbol rune
			if event.ID == tc.NavigationState.SelectedEventID {
				symbol = '◆'
			} else if event.Success {
				symbol = '●'
			} else {
				symbol = '○'
			}
			timeline[position] = symbol
		}
	}

	// Apply role-specific styling
	timelineStr := string(timeline)
	var styledTimeline string

	switch role {
	case UserMessage:
		styledTimeline = tc.Style.UserMessage.Render(timelineStr)
	case AssistantMessage:
		styledTimeline = tc.Style.AssistantMessage.Render(timelineStr)
	case ToolMessage:
		styledTimeline = tc.Style.ToolMessage.Render(timelineStr)
	case SystemMessage:
		styledTimeline = tc.Style.SystemMessage.Render(timelineStr)
	default:
		styledTimeline = timelineStr
	}

	track.WriteString(fmt.Sprintf("%-9s %s", label, styledTimeline))

	return track.String()
}

// renderSelectedEventDetails shows enhanced details for the selected event
func (tc *MessageTimelineComponent) renderSelectedEventDetails() string {
	if tc.NavigationState.SelectedEventID == "" {
		return ""
	}

	// Find selected event
	var selectedEvent *MessageEvent
	for _, event := range tc.Events {
		if event.ID == tc.NavigationState.SelectedEventID {
			selectedEvent = &event
			break
		}
	}

	if selectedEvent == nil {
		return ""
	}

	var details strings.Builder

	// Enhanced title with visual indicator
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("46")).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("46")).
		Padding(0, 2).
		Margin(1, 0)

	details.WriteString(titleStyle.Render("🔍 Selected Event Details") + "\n")

	// Enhanced event details with better formatting
	timeStr := selectedEvent.Timestamp.Format("2006-01-02 15:04:05")
	roleStr := string(selectedEvent.Role)

	// Success indicator with visual feedback
	var successStr string
	if selectedEvent.Success {
		successStr = lipgloss.NewStyle().Foreground(lipgloss.Color("46")).Bold(true).Render("✓ Success")
	} else {
		successStr = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true).Render("✗ Failed")
	}

	// Duration with performance indicator
	durationStr := selectedEvent.Duration.String()
	if selectedEvent.Duration < 500*time.Millisecond {
		durationStr = lipgloss.NewStyle().Foreground(lipgloss.Color("46")).Render("⚡ " + durationStr)
	} else if selectedEvent.Duration > 5*time.Second {
		durationStr = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render("🐌 " + durationStr)
	}

	// Main details section
	mainDetailsStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252")).
		Padding(1).
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("99"))

	mainDetails := fmt.Sprintf("🕰️ Time: %s\n💬 Role: %s\n%s\n⏱️ Duration: %s",
		timeStr, roleStr, successStr, durationStr)

	details.WriteString(mainDetailsStyle.Render(mainDetails) + "\n")

	// Agent and tools section
	if selectedEvent.AgentUsed != "" || len(selectedEvent.ToolsUsed) > 0 {
		var agentToolsDetails strings.Builder
		if selectedEvent.AgentUsed != "" {
			agentStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
			agentToolsDetails.WriteString(fmt.Sprintf("🤖 Agent: %s\n", agentStyle.Render(selectedEvent.AgentUsed)))
		}
		if len(selectedEvent.ToolsUsed) > 0 {
			toolsStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("220"))
			agentToolsDetails.WriteString(fmt.Sprintf("🛠️ Tools: %s", toolsStyle.Render(strings.Join(selectedEvent.ToolsUsed, ", "))))
		}

		agentToolsStyle := lipgloss.NewStyle().
			Padding(1).
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("99"))

		details.WriteString(agentToolsStyle.Render(agentToolsDetails.String()) + "\n")
	}

	// Content preview section
	if selectedEvent.Content != "" {
		contentPreview := selectedEvent.Content
		if len(contentPreview) > 100 {
			contentPreview = contentPreview[:97] + "..."
		}

		contentStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Italic(true).
			Padding(1).
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240"))

		contentText := fmt.Sprintf("📝 Content Preview:\n%s", contentPreview)
		details.WriteString(contentStyle.Render(contentText))
	}

	return details.String()
}

// renderNavigationHelp shows enhanced navigation instructions with visual indicators
func (tc *MessageTimelineComponent) renderNavigationHelp(content *strings.Builder) {
	zoomLevels := []string{"Hour", "Day", "Week", "Month"}
	currentZoom := "Day"
	if tc.NavigationState.ZoomLevel >= 0 && tc.NavigationState.ZoomLevel < len(zoomLevels) {
		currentZoom = zoomLevels[tc.NavigationState.ZoomLevel]
	}

	// Enhanced help with visual indicators and better organization
	zoomIndicator := lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true).Render(currentZoom)
	navCommands := lipgloss.NewStyle().Foreground(lipgloss.Color("99")).Render("←→ time nav")
	zoomCommands := lipgloss.NewStyle().Foreground(lipgloss.Color("99")).Render("↑↓ zoom")
	selectCommand := lipgloss.NewStyle().Foreground(lipgloss.Color("99")).Render("Space select")
	resetCommand := lipgloss.NewStyle().Foreground(lipgloss.Color("99")).Render("r reset")

	helpText := fmt.Sprintf("🔍 Zoom: %s | %s | %s | %s | %s",
		zoomIndicator, navCommands, zoomCommands, selectCommand, resetCommand)

	// Add event count if available
	visibleEvents := len(tc.filterEventsInRange())
	if visibleEvents > 0 {
		eventCountStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Italic(true)
		eventCount := eventCountStyle.Render(fmt.Sprintf(" | %d events visible", visibleEvents))
		helpText += eventCount
	}

	content.WriteString(tc.Style.NavigationHelp.Render(helpText))
}

// NavigateTime moves the timeline view forward or backward
func (tc *MessageTimelineComponent) NavigateTime(direction int) {
	// For btop-style graph, scroll the data window
	if tc.graphData != nil {
		// Calculate scroll amount (10% of visible range)
		scrollAmount := tc.graphData.BucketDuration * time.Duration(tc.graphWidth/10)
		if scrollAmount < tc.graphData.BucketDuration {
			scrollAmount = tc.graphData.BucketDuration
		}

		if direction < 0 {
			scrollAmount = -scrollAmount
		}

		// Update the scroll offset for smooth scrolling
		if direction > 0 {
			tc.scrollOffset += 5
		} else {
			tc.scrollOffset -= 5
			if tc.scrollOffset < 0 {
				tc.scrollOffset = 0
			}
		}

		// Reprocess graph data with new window
		tc.processGraphData()
	}

	// Also update navigation state for compatibility
	var moveAmount time.Duration
	switch tc.NavigationState.ZoomLevel {
	case 0: // Hour
		moveAmount = 15 * time.Minute
	case 1: // Day
		moveAmount = 2 * time.Hour
	case 2: // Week
		moveAmount = 12 * time.Hour
	case 3: // Month
		moveAmount = 3 * 24 * time.Hour
	default:
		moveAmount = 2 * time.Hour
	}

	if direction < 0 {
		moveAmount = -moveAmount
	}

	tc.NavigationState.CurrentRange.Start = tc.NavigationState.CurrentRange.Start.Add(moveAmount)
	tc.NavigationState.CurrentRange.End = tc.NavigationState.CurrentRange.End.Add(moveAmount)
}

// ZoomTimeline changes the zoom level
func (tc *MessageTimelineComponent) ZoomTimeline(direction int) {
	newZoomLevel := tc.NavigationState.ZoomLevel + direction

	// Clamp zoom level
	if newZoomLevel < 0 {
		newZoomLevel = 0
	}
	if newZoomLevel > 3 {
		newZoomLevel = 3
	}

	if newZoomLevel != tc.NavigationState.ZoomLevel {
		tc.NavigationState.ZoomLevel = newZoomLevel

		// Adjust graph resolution based on zoom
		switch newZoomLevel {
		case 0: // Hour - fine detail
			tc.graphWidth = 120
		case 1: // Day - normal detail
			tc.graphWidth = 80
		case 2: // Week - compressed
			tc.graphWidth = 60
		case 3: // Month - overview
			tc.graphWidth = 40
		}

		// Reprocess graph with new resolution
		tc.processGraphData()

		// Adjust time range based on new zoom level
		center := tc.NavigationState.CurrentRange.Start.Add(tc.NavigationState.CurrentRange.Duration / 2)
		var newDuration time.Duration

		switch newZoomLevel {
		case 0: // Hour
			newDuration = time.Hour
		case 1: // Day
			newDuration = 24 * time.Hour
		case 2: // Week
			newDuration = 7 * 24 * time.Hour
		case 3: // Month
			newDuration = 30 * 24 * time.Hour
		}

		tc.NavigationState.CurrentRange.Start = center.Add(-newDuration / 2)
		tc.NavigationState.CurrentRange.End = center.Add(newDuration / 2)
		tc.NavigationState.CurrentRange.Duration = newDuration
	}
}

// SelectNearestEvent selects the event closest to the current viewport center
func (tc *MessageTimelineComponent) SelectNearestEvent() {
	if len(tc.Events) == 0 {
		return
	}

	// Find center of current time range
	center := tc.NavigationState.CurrentRange.Start.Add(tc.NavigationState.CurrentRange.Duration / 2)

	// Find nearest event
	var nearestEvent *MessageEvent
	var minDistance time.Duration

	for _, event := range tc.Events {
		if event.Timestamp.After(tc.NavigationState.CurrentRange.Start) &&
		   event.Timestamp.Before(tc.NavigationState.CurrentRange.End) {

			distance := event.Timestamp.Sub(center)
			if distance < 0 {
				distance = -distance
			}

			if nearestEvent == nil || distance < minDistance {
				nearestEvent = &event
				minDistance = distance
			}
		}
	}

	if nearestEvent != nil {
		tc.NavigationState.SelectedEventID = nearestEvent.ID
	}
}

// ResetView resets the timeline to show recent events
func (tc *MessageTimelineComponent) ResetView() {
	now := time.Now()
	tc.NavigationState.CurrentRange = TimelineRange{
		Start:    now.Add(-24 * time.Hour),
		End:      now,
		Duration: 24 * time.Hour,
	}
	tc.NavigationState.ZoomLevel = 1
	tc.NavigationState.SelectedEventID = ""
	tc.NavigationState.ViewportOffset = 0
	tc.scrollOffset = 0
	tc.graphWidth = 80

	// Reprocess graph data for reset view
	tc.processGraphData()
}

// GetEventCount returns the number of events in the current time range
func (tc *MessageTimelineComponent) GetEventCount() int {
	return len(tc.filterEventsInRange())
}

// GetEventsByAgent returns events grouped by agent
func (tc *MessageTimelineComponent) GetEventsByAgent() map[string][]MessageEvent {
	eventsByAgent := make(map[string][]MessageEvent)

	for _, event := range tc.filterEventsInRange() {
		if event.AgentUsed != "" {
			eventsByAgent[event.AgentUsed] = append(eventsByAgent[event.AgentUsed], event)
		}
	}

	return eventsByAgent
}

// Braille characters for smooth graph rendering (btop-style)
var brailleChars = []rune{' ', '⡀', '⡄', '⡆', '⡇', '⣇', '⣧', '⣷', '⣿'}

// processGraphData processes events into graph buckets for btop-style visualization
func (tc *MessageTimelineComponent) processGraphData() {
	if len(tc.Events) == 0 {
		tc.graphData = nil
		return
	}

	// Determine time range from events
	startTime := tc.Events[0].Timestamp
	endTime := tc.Events[len(tc.Events)-1].Timestamp

	// Ensure we have at least an hour of data
	if endTime.Sub(startTime) < time.Hour {
		startTime = endTime.Add(-time.Hour)
	}

	// Calculate number of buckets based on graph width
	numBuckets := tc.graphWidth
	if numBuckets > 120 {
		numBuckets = 120
	}
	if numBuckets < 20 {
		numBuckets = 20
	}

	// Calculate bucket duration
	totalDuration := endTime.Sub(startTime)
	bucketDuration := totalDuration / time.Duration(numBuckets)

	// Initialize buckets
	userBuckets := make([]int, numBuckets)
	assistantBuckets := make([]int, numBuckets)

	// Fill buckets with event counts
	maxCount := 0
	for _, event := range tc.Events {
		// Calculate bucket index
		bucketIdx := int(event.Timestamp.Sub(startTime) / bucketDuration)
		if bucketIdx < 0 {
			bucketIdx = 0
		}
		if bucketIdx >= numBuckets {
			bucketIdx = numBuckets - 1
		}

		// Increment appropriate bucket
		switch event.Role {
		case UserMessage:
			userBuckets[bucketIdx]++
			if userBuckets[bucketIdx] > maxCount {
				maxCount = userBuckets[bucketIdx]
			}
		case AssistantMessage:
			assistantBuckets[bucketIdx]++
			if assistantBuckets[bucketIdx] > maxCount {
				maxCount = assistantBuckets[bucketIdx]
			}
		}
	}

	// Normalize to 0-1 range
	userActivity := make([]float64, numBuckets)
	assistantActivity := make([]float64, numBuckets)

	if maxCount > 0 {
		for i := 0; i < numBuckets; i++ {
			userActivity[i] = float64(userBuckets[i]) / float64(maxCount)
			assistantActivity[i] = float64(assistantBuckets[i]) / float64(maxCount)
		}
	}

	// Generate time labels
	timeLabels := tc.generateTimeLabels(startTime, endTime, numBuckets, bucketDuration)

	tc.graphData = &TimelineGraphData{
		UserActivity:      userActivity,
		AssistantActivity: assistantActivity,
		TimeLabels:        timeLabels,
		MaxValue:          maxCount,
		BucketDuration:    bucketDuration,
		StartTime:         startTime,
		EndTime:           endTime,
	}
}

// generateTimeLabels creates time labels for the x-axis
func (tc *MessageTimelineComponent) generateTimeLabels(startTime, endTime time.Time, numBuckets int, bucketDuration time.Duration) []string {
	labels := make([]string, numBuckets)

	// Determine label format based on duration
	totalDuration := endTime.Sub(startTime)
	var labelFormat string
	var labelInterval int

	if totalDuration <= time.Hour {
		labelFormat = "15:04"
		labelInterval = numBuckets / 4 // Show 4 labels
	} else if totalDuration <= 24*time.Hour {
		labelFormat = "15:04"
		labelInterval = numBuckets / 6 // Show 6 labels
	} else if totalDuration <= 7*24*time.Hour {
		labelFormat = "Mon 15:04"
		labelInterval = numBuckets / 7 // Show 7 labels
	} else {
		labelFormat = "Jan 2"
		labelInterval = numBuckets / 5 // Show 5 labels
	}

	if labelInterval < 1 {
		labelInterval = 1
	}

	for i := 0; i < numBuckets; i++ {
		if i%labelInterval == 0 || i == numBuckets-1 {
			bucketTime := startTime.Add(time.Duration(i) * bucketDuration)
			labels[i] = bucketTime.Format(labelFormat)
		} else {
			labels[i] = ""
		}
	}

	return labels
}

// renderBtopGraph renders a btop-style graph with smooth curves
func (tc *MessageTimelineComponent) renderBtopGraph() string {
	if tc.graphData == nil {
		return tc.Style.EmptyState.Render("No graph data available\n")
	}

	var content strings.Builder

	// Calculate actual graph dimensions
	graphWidth := tc.MaxWidth - 20 // Leave space for labels and values
	if graphWidth > 100 {
		graphWidth = 100 // Cap at reasonable width
	}
	if graphWidth < 40 {
		graphWidth = 40 // Minimum width for visibility
	}

	// Sample data to fit the width
	userSampled := tc.sampleData(tc.graphData.UserActivity, graphWidth)
	assistantSampled := tc.sampleData(tc.graphData.AssistantActivity, graphWidth)

	// Render graph with colors
	userStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("39"))      // Blue
	assistantStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("205")) // Pink/Magenta
	axisStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))      // Gray
	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("245"))     // Light gray
	valueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("252"))     // White

	// Calculate current values (last bucket or average of last few)
	userValue := 0
	assistantValue := 0
	if len(tc.graphData.UserActivity) > 0 {
		// Get average of last 3 buckets for smoother current value
		count := 0
		for i := len(tc.graphData.UserActivity) - 1; i >= 0 && count < 3; i-- {
			userValue += int(tc.graphData.UserActivity[i])
			assistantValue += int(tc.graphData.AssistantActivity[i])
			count++
		}
		if count > 0 {
			userValue /= count
			assistantValue /= count
		}
	}

	// Render horizontal timeline graph (btop-style)
	// User messages on top line
	content.WriteString(userStyle.Render("User "))
	content.WriteString(axisStyle.Render("│"))
	for i := 0; i < graphWidth && i < len(userSampled); i++ {
		height := userSampled[i]
		char := tc.getHorizontalBrailleChar(height)
		content.WriteString(userStyle.Render(string(char)))
	}
	content.WriteString(valueStyle.Render(fmt.Sprintf(" %d msg/hr", userValue)))
	content.WriteString("\n")

	// Assistant messages on bottom line
	content.WriteString(assistantStyle.Render("Asst "))
	content.WriteString(axisStyle.Render("│"))
	for i := 0; i < graphWidth && i < len(assistantSampled); i++ {
		height := assistantSampled[i]
		char := tc.getHorizontalBrailleChar(height)
		content.WriteString(assistantStyle.Render(string(char)))
	}
	content.WriteString(valueStyle.Render(fmt.Sprintf(" %d msg/hr", assistantValue)))
	content.WriteString("\n")

	// X-axis line
	content.WriteString(axisStyle.Render("     └"))
	content.WriteString(axisStyle.Render(strings.Repeat("─", graphWidth)))
	content.WriteString("\n")

	// Time scale markers
	content.WriteString("      ")
	timeMarkers := tc.createTimeMarkers(graphWidth)
	content.WriteString(labelStyle.Render(timeMarkers))
	content.WriteString("\n")

	return content.String()
}

// getHorizontalBrailleChar returns a braille character for horizontal graph based on value (0.0 to 1.0)
func (tc *MessageTimelineComponent) getHorizontalBrailleChar(value float64) rune {
	// Use different braille characters to show intensity
	if value <= 0.0 {
		return ' '
	} else if value < 0.1 {
		return '⠁'
	} else if value < 0.2 {
		return '⠃'
	} else if value < 0.3 {
		return '⠇'
	} else if value < 0.4 {
		return '⡇'
	} else if value < 0.5 {
		return '⣇'
	} else if value < 0.6 {
		return '⣧'
	} else if value < 0.7 {
		return '⣷'
	} else if value < 0.8 {
		return '⣿'
	} else if value < 0.9 {
		return '⣿'
	} else {
		return '█'
	}
}

// createGraphLine creates a single graph line using braille characters
func (tc *MessageTimelineComponent) createGraphLine(data []float64, width, height int) [][]rune {
	// Initialize graph matrix
	graph := make([][]rune, height)
	for i := range graph {
		graph[i] = make([]rune, width)
		for j := range graph[i] {
			graph[i][j] = ' '
		}
	}

	// Sample or interpolate data to match width
	sampledData := tc.sampleData(data, width)

	// Plot the graph line with braille characters
	for x := 0; x < width && x < len(sampledData); x++ {
		// Calculate Y position (inverted since terminal draws top-down)
		yPos := int((1.0 - sampledData[x]) * float64(height-1))
		if yPos < 0 {
			yPos = 0
		}
		if yPos >= height {
			yPos = height - 1
		}

		// Use braille character based on intensity
		intensity := sampledData[x]
		charIndex := int(intensity * float64(len(brailleChars)-1))
		if charIndex < 0 {
			charIndex = 0
		}
		if charIndex >= len(brailleChars) {
			charIndex = len(brailleChars) - 1
		}

		graph[yPos][x] = brailleChars[charIndex]

		// Add vertical fill for better visualization
		if intensity > 0.1 {
			for y := yPos + 1; y < height; y++ {
				fillIntensity := intensity * (1.0 - float64(y-yPos)/float64(height-yPos))
				fillCharIndex := int(fillIntensity * float64(len(brailleChars)-1))
				if fillCharIndex > 0 && graph[y][x] == ' ' {
					graph[y][x] = brailleChars[fillCharIndex]
				}
			}
		}
	}

	return graph
}

// sampleData samples or interpolates data to match target width
func (tc *MessageTimelineComponent) sampleData(data []float64, targetWidth int) []float64{
	if len(data) == 0 {
		return make([]float64, targetWidth)
	}

	if len(data) == targetWidth {
		return data
	}

	result := make([]float64, targetWidth)

	if len(data) < targetWidth {
		// Interpolate to expand
		ratio := float64(len(data)-1) / float64(targetWidth-1)
		for i := 0; i < targetWidth; i++ {
			srcIdx := float64(i) * ratio
			lowIdx := int(math.Floor(srcIdx))
			highIdx := int(math.Ceil(srcIdx))

			if highIdx >= len(data) {
				highIdx = len(data) - 1
			}

			if lowIdx == highIdx {
				result[i] = data[lowIdx]
			} else {
				// Linear interpolation
				t := srcIdx - float64(lowIdx)
				result[i] = data[lowIdx]*(1-t) + data[highIdx]*t
			}
		}
	} else {
		// Sample to compress
		ratio := float64(len(data)) / float64(targetWidth)
		for i := 0; i < targetWidth; i++ {
			startIdx := int(float64(i) * ratio)
			endIdx := int(float64(i+1) * ratio)
			if endIdx > len(data) {
				endIdx = len(data)
			}

			// Take average of samples in range
			sum := 0.0
			count := 0
			for j := startIdx; j < endIdx; j++ {
				sum += data[j]
				count++
			}
			if count > 0 {
				result[i] = sum / float64(count)
			}
		}
	}

	// Apply smoothing for better visualization
	smoothed := make([]float64, len(result))
	for i := range result {
		sum := result[i]
		count := 1.0

		if i > 0 {
			sum += result[i-1] * 0.3
			count += 0.3
		}
		if i < len(result)-1 {
			sum += result[i+1] * 0.3
			count += 0.3
		}

		smoothed[i] = sum / count
	}

	return smoothed
}

// createTimeMarkers creates time scale markers for x-axis
func (tc *MessageTimelineComponent) createTimeMarkers(width int) string {
	if tc.graphData == nil || len(tc.graphData.TimeLabels) == 0 {
		return strings.Repeat(" ", width)
	}

	markers := make([]rune, width)
	for i := range markers {
		markers[i] = ' '
	}

	// Calculate positions for time labels
	labelPositions := []int{0, width/4, width/2, 3*width/4, width-1}

	for _, pos := range labelPositions {
		if pos >= 0 && pos < width {
			// Map position to data index
			dataIdx := (pos * len(tc.graphData.TimeLabels)) / width
			if dataIdx >= len(tc.graphData.TimeLabels) {
				dataIdx = len(tc.graphData.TimeLabels) - 1
			}

			label := tc.graphData.TimeLabels[dataIdx]
			if label != "" {
				// Place label
				for i, ch := range label {
					targetPos := pos - len(label)/2 + i
					if targetPos >= 0 && targetPos < width {
						markers[targetPos] = ch
					}
				}
			}
		}
	}

	return string(markers)
}