package agent

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/rsmdt/the-startup/internal/stats"
)

// LeaderboardEntry represents a single entry in the leaderboard
type LeaderboardEntry struct {
	Name         string
	Count        int
	SuccessRate  float64
	AvgDuration  float64
	IsSelected   bool
	Rank         int
}

// LeaderboardComponentStyle defines styling for leaderboard components
type LeaderboardComponentStyle struct {
	Title            lipgloss.Style
	FocusedTitle     lipgloss.Style
	Header           lipgloss.Style
	Entry            lipgloss.Style
	SelectedEntry    lipgloss.Style
	EmptyState       lipgloss.Style
	Indicator        lipgloss.Style
	SortIndicator    lipgloss.Style
	HighSuccessRate  lipgloss.Style
	MidSuccessRate   lipgloss.Style
	LowSuccessRate   lipgloss.Style
	FastDuration     lipgloss.Style
	SlowDuration     lipgloss.Style
	Row              lipgloss.Style   // For matrix rows
	SelectedRow      lipgloss.Style   // For selected matrix rows
	Success          lipgloss.Style   // Generic success style
	Warning          lipgloss.Style   // Generic warning style
	Error            lipgloss.Style   // Generic error style
	Secondary        lipgloss.Style   // Secondary text style
	MatrixCellStyles [5]lipgloss.Style // Matrix cell strength colors
}

// DefaultLeaderboardStyle creates professional styling for leaderboard components
func DefaultLeaderboardStyle() LeaderboardComponentStyle {
	return LeaderboardComponentStyle{
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
		Header: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("99")).
			Background(lipgloss.Color("235")).
			Padding(0, 1).
			Margin(0, 0, 1, 0),
		Entry: lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")).
			Padding(0, 1),
		SelectedEntry: lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Background(lipgloss.Color("57")).
			Bold(true).
			Padding(0, 1),
		EmptyState: lipgloss.NewStyle().
			Italic(true).
			Foreground(lipgloss.Color("240")).
			Align(lipgloss.Center).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(2, 4).
			Margin(1, 0),
		Indicator: lipgloss.NewStyle().
			Foreground(lipgloss.Color("46")).
			Bold(true),
		SortIndicator: lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Bold(true),
		HighSuccessRate: lipgloss.NewStyle().
			Foreground(lipgloss.Color("46")).
			Bold(true), // Green for high success
		MidSuccessRate: lipgloss.NewStyle().
			Foreground(lipgloss.Color("220")).
			Bold(true), // Yellow for mid success
		LowSuccessRate: lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true), // Red for low success
		FastDuration: lipgloss.NewStyle().
			Foreground(lipgloss.Color("46")).
			Bold(true), // Green for fast
		SlowDuration: lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true), // Red for slow
		Row: lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")).
			Padding(0, 1),
		SelectedRow: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("15")).
			Background(lipgloss.Color("57")).
			Padding(0, 1),
		Success: lipgloss.NewStyle().
			Foreground(lipgloss.Color("46")).
			Bold(true),
		Warning: lipgloss.NewStyle().
			Foreground(lipgloss.Color("220")).
			Bold(true),
		Error: lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true),
		Secondary: lipgloss.NewStyle().
			Foreground(lipgloss.Color("244")).
			Italic(true),
		MatrixCellStyles: [5]lipgloss.Style{
			lipgloss.NewStyle().Foreground(lipgloss.Color("238")), // Minimal
			lipgloss.NewStyle().Foreground(lipgloss.Color("244")), // Weak
			lipgloss.NewStyle().Foreground(lipgloss.Color("220")), // Medium
			lipgloss.NewStyle().Foreground(lipgloss.Color("205")), // Strong
			lipgloss.NewStyle().Foreground(lipgloss.Color("46")),  // Very strong
		},
	}
}

// LeaderboardComponent provides reusable leaderboard rendering
type LeaderboardComponent struct {
	Title    string
	Style    LeaderboardComponentStyle
	MaxWidth int
}

// NewLeaderboardComponent creates a new leaderboard component
func NewLeaderboardComponent(title string) *LeaderboardComponent {
	return &LeaderboardComponent{
		Title:    title,
		Style:    DefaultLeaderboardStyle(),
		MaxWidth: 60,
	}
}


// RenderAgentLeaderboard renders agent statistics as a leaderboard
func (lc *LeaderboardComponent) RenderAgentLeaderboard(
	agentStats map[string]*stats.GlobalAgentStats,
	selectedIndex int,
	isFocused bool,
) string {
	var content strings.Builder

	// Title
	titleStyle := lc.Style.Title
	if isFocused {
		titleStyle = lc.Style.FocusedTitle
	}
	content.WriteString(titleStyle.Render(lc.Title) + "\n\n")

	if len(agentStats) == 0 {
		content.WriteString(lc.Style.EmptyState.Render("No agents found"))
		return content.String()
	}

	// Convert and sort entries
	entries := lc.convertAgentStatsToEntries(agentStats, selectedIndex)

	// Render header
	header := fmt.Sprintf("%-3s %-35s %8s", "Rank", "Agent", "Calls")
	content.WriteString(lc.Style.Header.Render(header) + "\n")

	// Render entries
	for i, entry := range entries {
		line := lc.renderLeaderboardEntry(entry, i == selectedIndex && isFocused)
		content.WriteString(line + "\n")
	}

	return content.String()
}

// RenderAgentLeaderboardWithSort renders agent statistics with professional styling and custom sorting
func (lc *LeaderboardComponent) RenderAgentLeaderboardWithSort(
	agentStats map[string]*stats.GlobalAgentStats,
	selectedIndex int,
	isFocused bool,
	sortColumn SortColumn,
	sortDirection SortDirection,
	sortIndicators map[SortColumn]string,
) string {
	var content strings.Builder

	// Enhanced title with focus indicator
	titleStyle := lc.Style.Title
	if isFocused {
		titleStyle = lc.Style.FocusedTitle
	}
	content.WriteString(titleStyle.Render(lc.Title) + "\n")

	if len(agentStats) == 0 {
		content.WriteString(lc.Style.EmptyState.Render("No agent data available\nRun some commands to see statistics"))
		return content.String()
	}

	// Convert and sort entries with custom sorting
	entries := lc.convertAndSortAgentStats(agentStats, sortColumn, sortDirection)

	// Enhanced header with sort indicators and better formatting
	headerText := fmt.Sprintf("%-5s %-35s %8s%s",
		"Rank", "Agent",
		"Calls", sortIndicators[SortByCount])
	content.WriteString(lc.Style.Header.Render(headerText) + "\n")

	// Add visual separator
	separator := strings.Repeat("─", len(headerText)-10) // Adjust for color codes
	separatorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("99")).Faint(true)
	content.WriteString(separatorStyle.Render(separator) + "\n")

	// Render entries with enhanced formatting and visual feedback
	for i, entry := range entries {
		line := lc.renderEnhancedEntry(entry, i == selectedIndex && isFocused)
		content.WriteString(line + "\n")
	}

	// Add summary footer if focused
	if isFocused && len(entries) > 0 {
		totalCalls := 0
		totalSuccess := 0
		for _, entry := range entries {
			totalCalls += entry.Count
			if entry.SuccessRate > 0 {
				totalSuccess += int(float64(entry.Count) * entry.SuccessRate / 100)
			}
		}
		avgSuccessRate := 0.0
		if totalCalls > 0 {
			avgSuccessRate = float64(totalSuccess) / float64(totalCalls) * 100
		}

		summaryStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Italic(true).
			Border(lipgloss.NormalBorder(), true, false, false, false).
			BorderForeground(lipgloss.Color("99")).
			Padding(1, 0, 0, 0)

		summary := fmt.Sprintf("Total: %d agents, %d calls, %.1f%% avg success",
			len(entries), totalCalls, avgSuccessRate)
		content.WriteString("\n" + summaryStyle.Render(summary))
	}

	return content.String()
}

// RenderDelegationLeaderboard renders delegation patterns as a leaderboard
func (lc *LeaderboardComponent) RenderDelegationLeaderboard(
	patterns []stats.DelegationPattern,
	selectedIndex int,
	isFocused bool,
	title string,
) string {
	var content strings.Builder

	// Title
	titleStyle := lc.Style.Title
	if isFocused {
		titleStyle = lc.Style.FocusedTitle
	}
	content.WriteString(titleStyle.Render(title) + "\n\n")

	if len(patterns) == 0 {
		content.WriteString(lc.Style.EmptyState.Render("No delegation patterns found"))
		return content.String()
	}

	// Sort patterns by count (descending)
	sortedPatterns := make([]stats.DelegationPattern, len(patterns))
	copy(sortedPatterns, patterns)
	sort.Slice(sortedPatterns, func(i, j int) bool {
		return sortedPatterns[i].Count > sortedPatterns[j].Count
	})

	// Render entries
	for i, pattern := range sortedPatterns {
		line := lc.renderDelegationEntry(pattern, i == selectedIndex && isFocused)
		content.WriteString(line + "\n")
	}

	return content.String()
}

// RenderCoOccurrenceLeaderboard renders co-occurrence patterns as a leaderboard
func (lc *LeaderboardComponent) RenderCoOccurrenceLeaderboard(
	coOccurrences []stats.AgentCoOccurrence,
	selectedIndex int,
	isFocused bool,
	title string,
) string {
	var content strings.Builder

	// Title
	titleStyle := lc.Style.Title
	if isFocused {
		titleStyle = lc.Style.FocusedTitle
	}
	content.WriteString(titleStyle.Render(title) + "\n\n")

	if len(coOccurrences) == 0 {
		content.WriteString(lc.Style.EmptyState.Render("No co-occurrence patterns found"))
		return content.String()
	}

	// Sort co-occurrences by count (descending)
	sortedCoOccurrences := make([]stats.AgentCoOccurrence, len(coOccurrences))
	copy(sortedCoOccurrences, coOccurrences)
	sort.Slice(sortedCoOccurrences, func(i, j int) bool {
		return sortedCoOccurrences[i].Count > sortedCoOccurrences[j].Count
	})

	// Render entries
	for i, coOcc := range sortedCoOccurrences {
		line := lc.renderCoOccurrenceEntry(coOcc, i == selectedIndex && isFocused)
		content.WriteString(line + "\n")
	}

	return content.String()
}

// convertAgentStatsToEntries converts agent stats map to sorted entries
func (lc *LeaderboardComponent) convertAgentStatsToEntries(
	agentStats map[string]*stats.GlobalAgentStats,
	selectedIndex int,
) []LeaderboardEntry {
	entries := make([]LeaderboardEntry, 0, len(agentStats))

	for name, stats := range agentStats {
		successRate := float64(0)
		if stats.Count > 0 {
			successRate = float64(stats.SuccessCount) / float64(stats.Count) * 100
		}

		avgDuration := float64(0)
		if stats.Count > 0 {
			avgDuration = stats.Mean()
		}

		entries = append(entries, LeaderboardEntry{
			Name:        name,
			Count:       stats.Count,
			SuccessRate: successRate,
			AvgDuration: avgDuration,
		})
	}

	// Sort by count (descending)
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Count > entries[j].Count
	})

	// Assign ranks
	for i := range entries {
		entries[i].Rank = i + 1
		entries[i].IsSelected = i == selectedIndex
	}

	return entries
}

// renderLeaderboardEntry renders a single leaderboard entry
func (lc *LeaderboardComponent) renderLeaderboardEntry(entry LeaderboardEntry, isSelected bool) string {
	var line strings.Builder

	// Selection indicator
	if isSelected {
		line.WriteString(lc.Style.Indicator.Render("▶ "))
	} else {
		line.WriteString("  ")
	}

	// Format entry - now showing only rank, agent name, and calls
	// Agent name is color-coded based on success rate
	agentName := lc.truncateString(entry.Name, 35)
	var agentText string
	if entry.SuccessRate >= 90 {
		agentText = lc.Style.HighSuccessRate.Render(agentName)
	} else if entry.SuccessRate >= 70 {
		agentText = lc.Style.MidSuccessRate.Render(agentName)
	} else {
		agentText = lc.Style.LowSuccessRate.Render(agentName)
	}

	entryText := fmt.Sprintf("%-3d %-35s %8d",
		entry.Rank,
		agentText,
		entry.Count)

	// Apply style based on selection
	if isSelected {
		entryText = lc.Style.SelectedEntry.Render(entryText)
	} else {
		entryText = lc.Style.Entry.Render(entryText)
	}

	line.WriteString(entryText)
	return line.String()
}

// renderDelegationEntry renders a single delegation pattern entry
func (lc *LeaderboardComponent) renderDelegationEntry(pattern stats.DelegationPattern, isSelected bool) string {
	var line strings.Builder

	// Selection indicator
	if isSelected {
		line.WriteString(lc.Style.Indicator.Render("▶ "))
	} else {
		line.WriteString("  ")
	}

	// Format delegation pattern with more space for agent names
	entryText := fmt.Sprintf("%s → %s (%d times)",
		lc.truncateString(pattern.SourceAgent, 20),
		lc.truncateString(pattern.TargetAgent, 20),
		pattern.Count)

	// Apply style based on selection
	if isSelected {
		entryText = lc.Style.SelectedEntry.Render(entryText)
	} else {
		entryText = lc.Style.Entry.Render(entryText)
	}

	line.WriteString(entryText)
	return line.String()
}

// renderCoOccurrenceEntry renders a single co-occurrence pattern entry
func (lc *LeaderboardComponent) renderCoOccurrenceEntry(coOcc stats.AgentCoOccurrence, isSelected bool) string {
	var line strings.Builder

	// Selection indicator
	if isSelected {
		line.WriteString(lc.Style.Indicator.Render("▶ "))
	} else {
		line.WriteString("  ")
	}

	// Format co-occurrence pattern with more space for agent names
	entryText := fmt.Sprintf("%s & %s (%d sessions)",
		lc.truncateString(coOcc.Agent1, 20),
		lc.truncateString(coOcc.Agent2, 20),
		coOcc.Count)

	// Apply style based on selection
	if isSelected {
		entryText = lc.Style.SelectedEntry.Render(entryText)
	} else {
		entryText = lc.Style.Entry.Render(entryText)
	}

	line.WriteString(entryText)
	return line.String()
}

// truncateString truncates string to specified length with ellipsis
func (lc *LeaderboardComponent) truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

// convertAndSortAgentStats converts and sorts agent stats based on specified criteria
func (lc *LeaderboardComponent) convertAndSortAgentStats(
	agentStats map[string]*stats.GlobalAgentStats,
	sortColumn SortColumn,
	sortDirection SortDirection,
) []LeaderboardEntry {
	entries := make([]LeaderboardEntry, 0, len(agentStats))

	for name, stats := range agentStats {
		successRate := float64(0)
		if stats.Count > 0 {
			successRate = float64(stats.SuccessCount) / float64(stats.Count) * 100
		}

		avgDuration := float64(0)
		if stats.Count > 0 {
			avgDuration = stats.Mean()
		}

		entries = append(entries, LeaderboardEntry{
			Name:        name,
			Count:       stats.Count,
			SuccessRate: successRate,
			AvgDuration: avgDuration,
		})
	}

	// Sort based on the specified column and direction
	sort.Slice(entries, func(i, j int) bool {
		var less bool
		switch sortColumn {
		case SortByCount:
			less = entries[i].Count < entries[j].Count
		case SortBySuccessRate:
			less = entries[i].SuccessRate < entries[j].SuccessRate
		case SortByDuration:
			less = entries[i].AvgDuration < entries[j].AvgDuration
		default:
			less = entries[i].Count < entries[j].Count
		}

		if sortDirection == SortDescending {
			return !less
		}
		return less
	})

	// Assign ranks
	for i := range entries {
		entries[i].Rank = i + 1
	}

	return entries
}

// renderEnhancedEntry renders a leaderboard entry with professional color coding and visual indicators
func (lc *LeaderboardComponent) renderEnhancedEntry(entry LeaderboardEntry, isSelected bool) string {
	var line strings.Builder

	// Enhanced selection indicator with visual feedback
	if isSelected {
		indicator := lc.Style.Indicator.Render("▶ ")
		line.WriteString(indicator)
	} else {
		line.WriteString("  ")
	}

	// Rank with medal indicators for top performers
	rankText := fmt.Sprintf("%-3d", entry.Rank)
	if entry.Rank == 1 {
		rankText = lipgloss.NewStyle().Foreground(lipgloss.Color("220")).Bold(true).Render("🥇 1")
	} else if entry.Rank == 2 {
		rankText = lipgloss.NewStyle().Foreground(lipgloss.Color("250")).Bold(true).Render("🥈 2")
	} else if entry.Rank == 3 {
		rankText = lipgloss.NewStyle().Foreground(lipgloss.Color("173")).Bold(true).Render("🥉 3")
	}

	// Agent name with truncation and color coding based on performance
	agentName := lc.truncateString(entry.Name, 35)

	// Apply color coding based on success rate for better visual feedback
	var agentText string
	if entry.SuccessRate >= 90 {
		agentText = lc.Style.HighSuccessRate.Render(fmt.Sprintf("%-35s", agentName))
	} else if entry.SuccessRate >= 70 {
		agentText = lc.Style.MidSuccessRate.Render(fmt.Sprintf("%-35s", agentName))
	} else {
		agentText = lc.Style.LowSuccessRate.Render(fmt.Sprintf("%-35s", agentName))
	}

	// Count with visual emphasis for high usage
	countText := fmt.Sprintf("%8d", entry.Count)
	if entry.Count >= 100 {
		countText = lipgloss.NewStyle().Foreground(lipgloss.Color("46")).Bold(true).Render(countText)
	} else if entry.Count >= 50 {
		countText = lipgloss.NewStyle().Foreground(lipgloss.Color("220")).Render(countText)
	}

	// Combine all text with proper spacing - simplified to rank, agent, calls
	entryText := fmt.Sprintf("%s %s %s",
		rankText, agentText, countText)

	// Apply overall style based on selection with enhanced visual feedback
	if isSelected {
		entryText = lc.Style.SelectedEntry.Render(entryText)
	} else {
		entryText = lc.Style.Entry.Render(entryText)
	}

	line.WriteString(entryText)
	return line.String()
}

// SetStyle allows customizing the leaderboard style
func (lc *LeaderboardComponent) SetStyle(style LeaderboardComponentStyle) {
	lc.Style = style
}

// SetMaxWidth sets the maximum width for the leaderboard
func (lc *LeaderboardComponent) SetMaxWidth(width int) {
	lc.MaxWidth = width
}

// CompactLeaderboardStyle provides a more compact styling option
func CompactLeaderboardStyle() LeaderboardComponentStyle {
	return LeaderboardComponentStyle{
		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("240")),
		FocusedTitle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")),
		Header: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("99")),
		Entry: lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")),
		SelectedEntry: lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")),
		EmptyState: lipgloss.NewStyle().
			Italic(true).
			Foreground(lipgloss.Color("240")),
		Indicator: lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")),
	}
}

// DelegationFilterOptions defines filtering options for delegation patterns
type DelegationFilterOptions struct {
	SortBy      DelegationSortType
	MinCount    int
	SourceAgent string
	TargetAgent string
	MaxResults  int
}

// DelegationSortType defines sorting options for delegation patterns
type DelegationSortType int

const (
	SortByFrequency DelegationSortType = iota
	SortByRecency
	SortBySource
	SortByTarget
)

// DelegationFlowComponent provides enhanced delegation pattern visualization
type DelegationFlowComponent struct {
	Style    LeaderboardComponentStyle
	MaxWidth int
}

// NewDelegationFlowComponent creates a new delegation flow component
func NewDelegationFlowComponent() *DelegationFlowComponent {
	return &DelegationFlowComponent{
		Style:    DefaultLeaderboardStyle(),
		MaxWidth: 80,
	}
}

// RenderDelegationFlows renders delegation patterns with flow visualization
func (dfc *DelegationFlowComponent) RenderDelegationFlows(
	patterns []stats.DelegationPattern,
	filterOptions DelegationFilterOptions,
	selectedIndex int,
	isFocused bool,
	title string,
) string {
	var content strings.Builder

	// Title
	titleStyle := dfc.Style.Title
	if isFocused {
		titleStyle = dfc.Style.FocusedTitle
	}
	content.WriteString(titleStyle.Render(title) + "\n\n")

	if len(patterns) == 0 {
		content.WriteString(dfc.Style.EmptyState.Render("No delegation patterns found"))
		return content.String()
	}

	// Filter and sort patterns
	filteredPatterns := dfc.filterPatterns(patterns, filterOptions)

	if len(filteredPatterns) == 0 {
		content.WriteString(dfc.Style.EmptyState.Render("No patterns match current filters"))
		return content.String()
	}

	// Render filter summary
	filterSummary := dfc.renderFilterSummary(filterOptions, len(patterns), len(filteredPatterns))
	if filterSummary != "" {
		content.WriteString(filterSummary + "\n\n")
	}

	// Render flow visualization
	flowViz := dfc.renderFlowVisualization(filteredPatterns[:min(5, len(filteredPatterns))])
	content.WriteString(flowViz + "\n")

	// Render detailed pattern list
	for i, pattern := range filteredPatterns {
		line := dfc.renderDelegationFlowEntry(pattern, i == selectedIndex && isFocused)
		content.WriteString(line + "\n")
	}

	return content.String()
}

// filterPatterns applies filtering and sorting to delegation patterns
func (dfc *DelegationFlowComponent) filterPatterns(
	patterns []stats.DelegationPattern,
	options DelegationFilterOptions,
) []stats.DelegationPattern {
	filteredPatterns := make([]stats.DelegationPattern, 0)

	// Apply filters
	for _, pattern := range patterns {
		// Min count filter
		if options.MinCount > 0 && pattern.Count < options.MinCount {
			continue
		}

		// Source agent filter
		if options.SourceAgent != "" && pattern.SourceAgent != options.SourceAgent {
			continue
		}

		// Target agent filter
		if options.TargetAgent != "" && pattern.TargetAgent != options.TargetAgent {
			continue
		}

		filteredPatterns = append(filteredPatterns, pattern)
	}

	// Sort patterns
	switch options.SortBy {
	case SortByFrequency:
		sort.Slice(filteredPatterns, func(i, j int) bool {
			return filteredPatterns[i].Count > filteredPatterns[j].Count
		})
	case SortByRecency:
		sort.Slice(filteredPatterns, func(i, j int) bool {
			return filteredPatterns[i].LastSeen.After(filteredPatterns[j].LastSeen)
		})
	case SortBySource:
		sort.Slice(filteredPatterns, func(i, j int) bool {
			return filteredPatterns[i].SourceAgent < filteredPatterns[j].SourceAgent
		})
	case SortByTarget:
		sort.Slice(filteredPatterns, func(i, j int) bool {
			return filteredPatterns[i].TargetAgent < filteredPatterns[j].TargetAgent
		})
	}

	// Limit results
	if options.MaxResults > 0 && len(filteredPatterns) > options.MaxResults {
		filteredPatterns = filteredPatterns[:options.MaxResults]
	}

	return filteredPatterns
}

// renderFilterSummary creates a summary of active filters
func (dfc *DelegationFlowComponent) renderFilterSummary(
	options DelegationFilterOptions,
	totalPatterns, filteredPatterns int,
) string {
	filters := []string{}

	if options.MinCount > 0 {
		filters = append(filters, fmt.Sprintf("min count: %d", options.MinCount))
	}
	if options.SourceAgent != "" {
		filters = append(filters, fmt.Sprintf("source: %s", options.SourceAgent))
	}
	if options.TargetAgent != "" {
		filters = append(filters, fmt.Sprintf("target: %s", options.TargetAgent))
	}

	if len(filters) == 0 {
		return ""
	}

	sortName := map[DelegationSortType]string{
		SortByFrequency: "frequency",
		SortByRecency:   "recency",
		SortBySource:    "source",
		SortByTarget:    "target",
	}[options.SortBy]

	filterText := fmt.Sprintf("Filters: %s | Sort: %s | Showing %d/%d patterns",
		strings.Join(filters, ", "), sortName, filteredPatterns, totalPatterns)

	return dfc.Style.Header.Render(filterText)
}

// renderFlowVisualization creates ASCII flow diagram of top patterns
func (dfc *DelegationFlowComponent) renderFlowVisualization(patterns []stats.DelegationPattern) string {
	if len(patterns) == 0 {
		return ""
	}

	var content strings.Builder
	content.WriteString(dfc.Style.Header.Render("Top Flow Patterns") + "\n")

	// Group patterns by flow
	flows := make(map[string][]stats.DelegationPattern)
	for _, pattern := range patterns {
		flowKey := fmt.Sprintf("%s→%s", pattern.SourceAgent, pattern.TargetAgent)
		flows[flowKey] = append(flows[flowKey], pattern)
	}

	// Render each flow
	for _, flowPatterns := range flows {
		if len(flowPatterns) == 0 {
			continue
		}

		pattern := flowPatterns[0] // Use first pattern for display
		flowLine := dfc.renderSingleFlow(pattern)
		content.WriteString(flowLine + "\n")
	}

	return content.String()
}

// renderSingleFlow creates ASCII representation of a single delegation flow
func (dfc *DelegationFlowComponent) renderSingleFlow(pattern stats.DelegationPattern) string {
	// Truncate agent names for display
	sourceAgent := dfc.truncateString(pattern.SourceAgent, 12)
	targetAgent := dfc.truncateString(pattern.TargetAgent, 12)

	// Create flow visualization with frequency indicator
	frequencyBar := dfc.createFrequencyBar(pattern.Count, 20)
	timeSince := dfc.formatTimeSince(pattern.LastSeen)

	// Format: Source ──[freq]──> Target (count times, last: time)
	flowText := fmt.Sprintf("%-12s ──[%s]──> %-12s (%d times, last: %s)",
		sourceAgent, frequencyBar, targetAgent, pattern.Count, timeSince)

	return dfc.Style.Entry.Render(flowText)
}

// createFrequencyBar creates a visual frequency indicator
func (dfc *DelegationFlowComponent) createFrequencyBar(count int, maxWidth int) string {
	if count <= 0 {
		return strings.Repeat("·", maxWidth)
	}

	// Scale based on count (simple linear scaling)
	barLength := min(maxWidth, max(1, count/10+1))
	bar := strings.Repeat("█", barLength)
	padding := strings.Repeat("·", maxWidth-barLength)

	return bar + padding
}

// formatTimeSince formats time duration since last seen
func (dfc *DelegationFlowComponent) formatTimeSince(lastSeen time.Time) string {
	duration := time.Since(lastSeen)

	if duration < time.Hour {
		return fmt.Sprintf("%dm ago", int(duration.Minutes()))
	} else if duration < 24*time.Hour {
		return fmt.Sprintf("%dh ago", int(duration.Hours()))
	} else {
		return fmt.Sprintf("%dd ago", int(duration.Hours()/24))
	}
}

// renderDelegationFlowEntry renders a single delegation pattern with enhanced details
func (dfc *DelegationFlowComponent) renderDelegationFlowEntry(pattern stats.DelegationPattern, isSelected bool) string {
	var line strings.Builder

	// Selection indicator
	if isSelected {
		line.WriteString(dfc.Style.Indicator.Render("▶ "))
	} else {
		line.WriteString("  ")
	}

	// Enhanced delegation pattern with frequency and timing info
	sourceAgent := dfc.truncateString(pattern.SourceAgent, 15)
	targetAgent := dfc.truncateString(pattern.TargetAgent, 15)
	timeSince := dfc.formatTimeSince(pattern.LastSeen)

	entryText := fmt.Sprintf("%-15s → %-15s │ %3d times │ %s",
		sourceAgent, targetAgent, pattern.Count, timeSince)

	// Apply style based on selection
	if isSelected {
		entryText = dfc.Style.SelectedEntry.Render(entryText)
	} else {
		entryText = dfc.Style.Entry.Render(entryText)
	}

	line.WriteString(entryText)
	return line.String()
}

// truncateString truncates string to specified length with ellipsis
func (dfc *DelegationFlowComponent) truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

// CoOccurrencePattern represents different types of agent co-occurrence patterns
type CoOccurrencePattern string

const (
	Sequential   CoOccurrencePattern = "sequential"
	Concurrent   CoOccurrencePattern = "concurrent"
	Alternative  CoOccurrencePattern = "alternative"
)

// AgentRelationship represents enhanced agent relationship data for matrix display
type AgentRelationship struct {
	Agent1       string
	Agent2       string
	Sessions     []string
	Frequency    int
	Correlation  float64
	Pattern      CoOccurrencePattern
	Strength     RelationshipStrength
	LastSeen     time.Time
}

// RelationshipStrength represents the strength of agent relationships
type RelationshipStrength string

const (
	StrongRelationship   RelationshipStrength = "strong"
	MediumRelationship   RelationshipStrength = "medium"
	WeakRelationship     RelationshipStrength = "weak"
	MinimalRelationship  RelationshipStrength = "minimal"
)

// MatrixDisplayMode represents different matrix visualization modes
type MatrixDisplayMode string

const (
	MatrixView     MatrixDisplayMode = "matrix"
	RelationshipView MatrixDisplayMode = "relationships"
	InsightsView     MatrixDisplayMode = "insights"
)

// CoOccurrenceMatrixComponent provides enhanced co-occurrence pattern visualization
type CoOccurrenceMatrixComponent struct {
	Style         LeaderboardComponentStyle
	MaxWidth      int
	DisplayMode   MatrixDisplayMode
	SelectedRow   int
	SelectedCol   int
	MatrixData    [][]AgentRelationship
	AgentNames    []string
	currentTheme  ColorTheme
}

// NewCoOccurrenceMatrixComponent creates a new co-occurrence matrix component
func NewCoOccurrenceMatrixComponent() *CoOccurrenceMatrixComponent {
	return &CoOccurrenceMatrixComponent{
		Style:       DefaultLeaderboardStyle(),
		MaxWidth:    80,
		DisplayMode: MatrixView,
		SelectedRow: 0,
		SelectedCol: 0,
		MatrixData:  [][]AgentRelationship{},
		AgentNames:  []string{},
	}
}

// ApplyTheme applies a color theme to the co-occurrence matrix component
func (cmc *CoOccurrenceMatrixComponent) ApplyTheme(theme ColorTheme) {
	cmc.currentTheme = theme
	cmc.Style = LeaderboardComponentStyle{
		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(theme.Primary).
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(theme.Border).
			Padding(0, 0, 1, 0),
		FocusedTitle: lipgloss.NewStyle().
			Bold(true).
			Foreground(theme.Highlight).
			Background(theme.Border).
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(theme.Highlight).
			Padding(0, 1, 1, 0),
		Header: lipgloss.NewStyle().
			Bold(true).
			Foreground(theme.Secondary).
			Background(theme.Border).
			Padding(0, 1),
		Row: lipgloss.NewStyle().
			Foreground(theme.Foreground).
			Padding(0, 1),
		SelectedRow: lipgloss.NewStyle().
			Bold(true).
			Foreground(theme.Foreground).
			Background(theme.Secondary).
			Padding(0, 1),
		Success: lipgloss.NewStyle().
			Foreground(theme.Success).
			Bold(true),
		Warning: lipgloss.NewStyle().
			Foreground(theme.Warning).
			Bold(true),
		Error: lipgloss.NewStyle().
			Foreground(theme.Error).
			Bold(true),
		Secondary: lipgloss.NewStyle().
			Foreground(theme.Info).
			Italic(true),
		EmptyState: lipgloss.NewStyle().
			Italic(true).
			Foreground(theme.Muted).
			Align(lipgloss.Center).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(theme.Border).
			Padding(2, 4).
			Margin(1, 0),
		// Matrix strength colors using theme's matrix colors
		MatrixCellStyles: [5]lipgloss.Style{
			lipgloss.NewStyle().Foreground(theme.Matrix[0]),
			lipgloss.NewStyle().Foreground(theme.Matrix[1]),
			lipgloss.NewStyle().Foreground(theme.Matrix[2]),
			lipgloss.NewStyle().Foreground(theme.Matrix[3]),
			lipgloss.NewStyle().Foreground(theme.Matrix[4]),
		},
	}
}

// RenderCoOccurrenceMatrixContent renders just the matrix content without title (for btop-style borders)
func (cmc *CoOccurrenceMatrixComponent) RenderCoOccurrenceMatrixContent(
	coOccurrences []stats.AgentCoOccurrence,
	selectedIndex int,
	isFocused bool,
) string {
	var content strings.Builder

	// Add mode indicator
	modeStyle := cmc.Style.Secondary
	modeText := fmt.Sprintf("Mode: %s", string(cmc.DisplayMode))
	content.WriteString(modeStyle.Render(modeText) + "\n\n")

	if len(coOccurrences) == 0 {
		content.WriteString(cmc.Style.EmptyState.Render("No co-occurrence patterns found"))
		return content.String()
	}

	// Build enhanced agent relationships with correlation analysis
	relationships := cmc.buildAgentRelationships(coOccurrences)

	// Render based on display mode
	switch cmc.DisplayMode {
	case MatrixView:
		content.WriteString(cmc.renderMatrixView(relationships, isFocused))
	case RelationshipView:
		content.WriteString(cmc.renderRelationshipView(relationships, selectedIndex, isFocused))
	case InsightsView:
		content.WriteString(cmc.renderInsightsView(relationships))
	default:
		content.WriteString(cmc.renderMatrixView(relationships, isFocused))
	}

	return content.String()
}

// RenderCoOccurrenceMatrix renders co-occurrence patterns with enhanced matrix visualization (DEPRECATED - use RenderCoOccurrenceMatrixContent for btop-style)
func (cmc *CoOccurrenceMatrixComponent) RenderCoOccurrenceMatrix(
	coOccurrences []stats.AgentCoOccurrence,
	selectedIndex int,
	isFocused bool,
	title string,
) string {
	var content strings.Builder

	// Title with display mode indicator
	titleStyle := cmc.Style.Title
	if isFocused {
		titleStyle = cmc.Style.FocusedTitle
	}
	fullTitle := fmt.Sprintf("%s [%s]", title, string(cmc.DisplayMode))
	content.WriteString(titleStyle.Render(fullTitle) + "\n\n")

	if len(coOccurrences) == 0 {
		content.WriteString(cmc.Style.EmptyState.Render("No co-occurrence patterns found"))
		return content.String()
	}

	// Build enhanced agent relationships with correlation analysis
	relationships := cmc.buildAgentRelationships(coOccurrences)

	// Render based on display mode
	switch cmc.DisplayMode {
	case MatrixView:
		content.WriteString(cmc.renderMatrixView(relationships, isFocused))
	case RelationshipView:
		content.WriteString(cmc.renderRelationshipView(relationships, selectedIndex, isFocused))
	case InsightsView:
		content.WriteString(cmc.renderInsightsView(relationships))
	default:
		content.WriteString(cmc.renderMatrixView(relationships, isFocused))
	}

	// Navigation hints
	if isFocused {
		hints := "\n" + cmc.Style.Header.Render("Navigation: m=matrix, r=relationships, i=insights, ↑↓←→=navigate")
		content.WriteString(hints)
	}

	return content.String()
}

// buildAgentRelationships converts co-occurrence data to enhanced relationship objects
func (cmc *CoOccurrenceMatrixComponent) buildAgentRelationships(coOccurrences []stats.AgentCoOccurrence) []AgentRelationship {
	relationships := make([]AgentRelationship, 0, len(coOccurrences))

	for _, coOcc := range coOccurrences {
		// Calculate correlation based on session frequency and recency
		correlation := cmc.calculateCorrelation(coOcc)

		// Determine pattern based on agent names and frequency
		pattern := cmc.determinePattern(coOcc)

		// Calculate relationship strength
		strength := cmc.calculateRelationshipStrength(coOcc.Count)

		relationships = append(relationships, AgentRelationship{
			Agent1:      coOcc.Agent1,
			Agent2:      coOcc.Agent2,
			Sessions:    coOcc.Sessions,
			Frequency:   coOcc.Count,
			Correlation: correlation,
			Pattern:     pattern,
			Strength:    strength,
			LastSeen:    coOcc.LastSeen,
		})
	}

	// Sort by correlation strength
	sort.Slice(relationships, func(i, j int) bool {
		return relationships[i].Correlation > relationships[j].Correlation
	})

	return relationships
}

// calculateCorrelation computes correlation strength based on frequency and recency
func (cmc *CoOccurrenceMatrixComponent) calculateCorrelation(coOcc stats.AgentCoOccurrence) float64 {
	// Base correlation on frequency (normalized to 0-1 scale)
	frequencyScore := math.Min(float64(coOcc.Count)/10.0, 1.0)

	// Recency factor (sessions within last week get bonus)
	recencyScore := 0.0
	if time.Since(coOcc.LastSeen) < 7*24*time.Hour {
		recencyScore = 0.2 // 20% bonus for recent activity
	}

	// Combined correlation score
	correlation := (frequencyScore * 0.8) + recencyScore
	return math.Min(correlation, 1.0)
}

// determinePattern analyzes agent usage patterns
func (cmc *CoOccurrenceMatrixComponent) determinePattern(coOcc stats.AgentCoOccurrence) CoOccurrencePattern {
	// Simple pattern determination based on agent types and frequency
	if coOcc.Count >= 8 {
		return Concurrent // High frequency suggests concurrent usage
	} else if coOcc.Count >= 3 {
		return Sequential // Medium frequency suggests sequential workflow
	} else {
		return Alternative // Low frequency suggests alternative usage
	}
}

// calculateRelationshipStrength determines relationship strength category
func (cmc *CoOccurrenceMatrixComponent) calculateRelationshipStrength(count int) RelationshipStrength {
	if count >= 10 {
		return StrongRelationship
	} else if count >= 5 {
		return MediumRelationship
	} else if count >= 2 {
		return WeakRelationship
	} else {
		return MinimalRelationship
	}
}

// renderMatrixView renders agent relationships in matrix format
func (cmc *CoOccurrenceMatrixComponent) renderMatrixView(relationships []AgentRelationship, isFocused bool) string {
	var content strings.Builder

	// Build unique agent list
	agentSet := make(map[string]bool)
	for _, rel := range relationships {
		agentSet[rel.Agent1] = true
		agentSet[rel.Agent2] = true
	}

	agents := make([]string, 0, len(agentSet))
	for agent := range agentSet {
		agents = append(agents, agent)
	}
	sort.Strings(agents)

	// Limit to top 6 agents for matrix display
	if len(agents) > 6 {
		agents = agents[:6]
	}

	// Build relationship lookup
	relLookup := make(map[string]AgentRelationship)
	for _, rel := range relationships {
		key1 := fmt.Sprintf("%s_%s", rel.Agent1, rel.Agent2)
		key2 := fmt.Sprintf("%s_%s", rel.Agent2, rel.Agent1)
		relLookup[key1] = rel
		relLookup[key2] = rel
	}

	// Render matrix header
	content.WriteString(cmc.Style.Header.Render("Agent Collaboration Matrix") + "\n")
	content.WriteString(cmc.Style.Entry.Render("Legend: █=Strong ▓=Medium ▒=Weak ░=Minimal ·=None") + "\n\n")

	// Render column headers - use full agent names
	headerLine := "                        " // Space for row labels
	for _, agent := range agents {
		// Use full name up to 20 chars for better readability
		headerLine += fmt.Sprintf("%-20s", cmc.truncateString(agent, 18))
	}
	content.WriteString(cmc.Style.Header.Render(headerLine) + "\n")

	// Render matrix rows
	for i, agent1 := range agents {
		var line strings.Builder

		// Row label - use full agent name
		rowLabel := fmt.Sprintf("%-22s ", cmc.truncateString(agent1, 20))
		line.WriteString(rowLabel)

		for j, agent2 := range agents {
			if i == j {
				// Diagonal - same agent
				cell := "        ■           " // Wider cell for new column width
				if isFocused && cmc.SelectedRow == i && cmc.SelectedCol == j {
					cell = cmc.Style.SelectedEntry.Render(cell)
				}
				line.WriteString(cell)
			} else {
				// Check for relationship
				key := fmt.Sprintf("%s_%s", agent1, agent2)
				if rel, exists := relLookup[key]; exists {
					cell := cmc.renderMatrixCell(rel, i == cmc.SelectedRow && j == cmc.SelectedCol && isFocused)
					line.WriteString(cell)
				} else {
					cell := "        ·           " // Wider cell for new column width
					if isFocused && cmc.SelectedRow == i && cmc.SelectedCol == j {
						cell = cmc.Style.SelectedEntry.Render(cell)
					}
					line.WriteString(cell)
				}
			}
		}
		content.WriteString(line.String() + "\n")
	}

	// Show selected relationship details
	if isFocused && cmc.SelectedRow < len(agents) && cmc.SelectedCol < len(agents) {
		agent1 := agents[cmc.SelectedRow]
		agent2 := agents[cmc.SelectedCol]
		key := fmt.Sprintf("%s_%s", agent1, agent2)
		if rel, exists := relLookup[key]; exists {
			details := cmc.renderRelationshipDetails(rel)
			content.WriteString("\n" + details)
		}
	}

	return content.String()
}

// renderMatrixCell renders a single cell in the matrix with relationship strength indicator
func (cmc *CoOccurrenceMatrixComponent) renderMatrixCell(rel AgentRelationship, isSelected bool) string {
	var symbol string
	switch rel.Strength {
	case StrongRelationship:
		symbol = "        █           " // Full block with padding
	case MediumRelationship:
		symbol = "        ▓           " // Medium shade with padding
	case WeakRelationship:
		symbol = "        ▒           " // Light shade with padding
	case MinimalRelationship:
		symbol = "        ░           " // Very light shade with padding
	default:
		symbol = "        ·           " // No relationship with padding
	}

	if isSelected {
		symbol = cmc.Style.SelectedEntry.Render(symbol)
	}

	return symbol
}

// renderRelationshipDetails renders detailed information about a selected relationship
func (cmc *CoOccurrenceMatrixComponent) renderRelationshipDetails(rel AgentRelationship) string {
	var content strings.Builder

	content.WriteString(cmc.Style.Header.Render("Relationship Details") + "\n")
	details := fmt.Sprintf("Agents: %s ⟷ %s\n", rel.Agent1, rel.Agent2)
	details += fmt.Sprintf("Frequency: %d sessions\n", rel.Frequency)
	details += fmt.Sprintf("Correlation: %.2f\n", rel.Correlation)
	details += fmt.Sprintf("Pattern: %s\n", string(rel.Pattern))
	details += fmt.Sprintf("Strength: %s\n", string(rel.Strength))
	details += fmt.Sprintf("Last Seen: %s", cmc.formatTimeSince(rel.LastSeen))

	content.WriteString(cmc.Style.Entry.Render(details))
	return content.String()
}

// renderRelationshipView renders relationships in list format with correlation analysis
func (cmc *CoOccurrenceMatrixComponent) renderRelationshipView(relationships []AgentRelationship, selectedIndex int, isFocused bool) string {
	var content strings.Builder

	content.WriteString(cmc.Style.Header.Render("Agent Relationships & Correlation Analysis") + "\n\n")

	if len(relationships) == 0 {
		content.WriteString(cmc.Style.EmptyState.Render("No relationships found"))
		return content.String()
	}

	// Render top relationships
	topRelationships := min(10, len(relationships))
	for i := 0; i < topRelationships; i++ {
		rel := relationships[i]
		line := cmc.renderRelationshipEntry(rel, i == selectedIndex && isFocused)
		content.WriteString(line + "\n")
	}

	return content.String()
}

// renderRelationshipEntry renders a single relationship with correlation info
func (cmc *CoOccurrenceMatrixComponent) renderRelationshipEntry(rel AgentRelationship, isSelected bool) string {
	var line strings.Builder

	// Selection indicator
	if isSelected {
		line.WriteString(cmc.Style.Indicator.Render("▶ "))
	} else {
		line.WriteString("  ")
	}

	// Format agent names
	agent1 := cmc.truncateString(rel.Agent1, 12)
	agent2 := cmc.truncateString(rel.Agent2, 12)

	// Create pattern indicator
	patternSymbol := map[CoOccurrencePattern]string{
		Sequential:  "→",
		Concurrent:  "⟷",
		Alternative: "↔",
	}[rel.Pattern]

	// Create strength and correlation display
	strengthIndicator := cmc.createVisualStrengthIndicator(rel.Strength)
	correlationText := fmt.Sprintf("%.2f", rel.Correlation)
	timeSince := cmc.formatTimeSince(rel.LastSeen)

	// Compose entry line
	entryText := fmt.Sprintf("%-12s %s %-12s │ %s │ %s │ r=%-4s │ %s",
		agent1, patternSymbol, agent2, strengthIndicator, string(rel.Pattern)[:3], correlationText, timeSince)

	// Apply style based on selection
	if isSelected {
		entryText = cmc.Style.SelectedEntry.Render(entryText)
	} else {
		entryText = cmc.Style.Entry.Render(entryText)
	}

	line.WriteString(entryText)
	return line.String()
}

// createVisualStrengthIndicator creates a visual indicator of relationship strength
func (cmc *CoOccurrenceMatrixComponent) createVisualStrengthIndicator(strength RelationshipStrength) string {
	switch strength {
	case StrongRelationship:
		return "████ Strong "
	case MediumRelationship:
		return "███· Medium "
	case WeakRelationship:
		return "██·· Weak   "
	case MinimalRelationship:
		return "█··· Minimal"
	default:
		return "···· None   "
	}
}

// renderInsightsView renders analytical insights about agent collaboration patterns
func (cmc *CoOccurrenceMatrixComponent) renderInsightsView(relationships []AgentRelationship) string {
	var content strings.Builder

	content.WriteString(cmc.Style.Header.Render("Agent Collaboration Insights") + "\n\n")

	if len(relationships) == 0 {
		content.WriteString(cmc.Style.EmptyState.Render("No collaboration data available"))
		return content.String()
	}

	// Pattern analysis
	patternCounts := make(map[CoOccurrencePattern]int)
	strengthCounts := make(map[RelationshipStrength]int)
	totalCorrelation := 0.0

	for _, rel := range relationships {
		patternCounts[rel.Pattern]++
		strengthCounts[rel.Strength]++
		totalCorrelation += rel.Correlation
	}

	avgCorrelation := totalCorrelation / float64(len(relationships))

	// Render pattern distribution
	content.WriteString(cmc.Style.Header.Render("Collaboration Patterns") + "\n")
	for pattern, count := range patternCounts {
		percentage := float64(count) / float64(len(relationships)) * 100
		patternLine := fmt.Sprintf("  %s: %d relationships (%.1f%%)",
			string(pattern), count, percentage)
		content.WriteString(cmc.Style.Entry.Render(patternLine) + "\n")
	}

	content.WriteString("\n")

	// Render strength distribution
	content.WriteString(cmc.Style.Header.Render("Relationship Strength Distribution") + "\n")
	for strength, count := range strengthCounts {
		percentage := float64(count) / float64(len(relationships)) * 100
		strengthLine := fmt.Sprintf("  %s: %d relationships (%.1f%%)",
			string(strength), count, percentage)
		content.WriteString(cmc.Style.Entry.Render(strengthLine) + "\n")
	}

	content.WriteString("\n")

	// Render correlation insights
	content.WriteString(cmc.Style.Header.Render("Correlation Analysis") + "\n")
	correlationLine := fmt.Sprintf("  Average Correlation: %.3f", avgCorrelation)
	content.WriteString(cmc.Style.Entry.Render(correlationLine) + "\n")

	// Find strongest relationships
	strongRelationships := 0
	for _, rel := range relationships {
		if rel.Correlation >= 0.8 {
			strongRelationships++
		}
	}
	strongLine := fmt.Sprintf("  High Correlation (≥0.8): %d relationships", strongRelationships)
	content.WriteString(cmc.Style.Entry.Render(strongLine) + "\n")

	// Top collaborative pairs
	content.WriteString("\n" + cmc.Style.Header.Render("Top Collaborative Pairs") + "\n")
	topPairs := min(5, len(relationships))
	for i := 0; i < topPairs; i++ {
		rel := relationships[i]
		pairLine := fmt.Sprintf("  %s ⟷ %s (r=%.3f, %d sessions)",
			cmc.truncateString(rel.Agent1, 12),
			cmc.truncateString(rel.Agent2, 12),
			rel.Correlation, rel.Frequency)
		content.WriteString(cmc.Style.Entry.Render(pairLine) + "\n")
	}

	return content.String()
}

// SwitchDisplayMode cycles through available display modes
func (cmc *CoOccurrenceMatrixComponent) SwitchDisplayMode() {
	switch cmc.DisplayMode {
	case MatrixView:
		cmc.DisplayMode = RelationshipView
	case RelationshipView:
		cmc.DisplayMode = InsightsView
	case InsightsView:
		cmc.DisplayMode = MatrixView
	}
}

// NavigateMatrix handles matrix navigation
func (cmc *CoOccurrenceMatrixComponent) NavigateMatrix(direction string) {
	if cmc.DisplayMode != MatrixView {
		return
	}

	maxSize := 6 // Assuming 6x6 matrix max
	switch direction {
	case "up":
		if cmc.SelectedRow > 0 {
			cmc.SelectedRow--
		}
	case "down":
		if cmc.SelectedRow < maxSize-1 {
			cmc.SelectedRow++
		}
	case "left":
		if cmc.SelectedCol > 0 {
			cmc.SelectedCol--
		}
	case "right":
		if cmc.SelectedCol < maxSize-1 {
			cmc.SelectedCol++
		}
	}
}

// formatTimeSince formats time duration since last seen (reuse from DelegationFlowComponent)
func (cmc *CoOccurrenceMatrixComponent) formatTimeSince(lastSeen time.Time) string {
	duration := time.Since(lastSeen)

	if duration < time.Hour {
		return fmt.Sprintf("%dm ago", int(duration.Minutes()))
	} else if duration < 24*time.Hour {
		return fmt.Sprintf("%dh ago", int(duration.Hours()))
	} else {
		return fmt.Sprintf("%dd ago", int(duration.Hours()/24))
	}
}

// truncateString truncates string to specified length with ellipsis (reuse from DelegationFlowComponent)
func (cmc *CoOccurrenceMatrixComponent) truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

// Helper functions
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}