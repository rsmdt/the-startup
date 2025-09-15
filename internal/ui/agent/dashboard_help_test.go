package agent

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// ============================================================================
// Help System Tests - Testing help overlay display and navigation
// Tests for h/? key toggles, comprehensive help content, and context-sensitive help
// ============================================================================

// TestHelpSystemToggle tests help overlay show/hide functionality
func TestHelpSystemToggle(t *testing.T) {
	mockStats := createMockAgentStats()
	mockDelegations := createMockDelegationPatterns()
	mockCoOccurrences := createMockCoOccurrencePatterns()

	model := NewDashboardModel(mockStats, mockDelegations, mockCoOccurrences)

	// Initially help should be hidden
	if model.showHelp {
		t.Error("Help should be hidden by default")
	}

	// Test '?' key shows help
	helpMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}
	_, cmd := model.Update(helpMsg)

	if !model.showHelp {
		t.Error("Help should be shown after '?' key")
	}

	// Should not trigger additional commands
	if cmd != nil {
		t.Error("Help toggle should not trigger commands")
	}

	// Test view contains help content
	view := model.View()
	if !strings.Contains(view, "Help") {
		t.Error("View should contain help content when help is shown")
	}

	// Test '?' key again hides help
	model.Update(helpMsg)
	if model.showHelp {
		t.Error("Help should be hidden after second '?' key")
	}

	// View should show normal dashboard
	view = model.View()
	if strings.Contains(view, "Agent Dashboard Help") {
		t.Error("View should not contain help content when help is hidden")
	}
}

// TestHelpSystemHKeyAlternative tests 'h' key as alternative help toggle
func TestHelpSystemHKeyAlternative(t *testing.T) {
	mockStats := createMockAgentStats()
	mockDelegations := createMockDelegationPatterns()
	mockCoOccurrences := createMockCoOccurrencePatterns()

	model := NewDashboardModel(mockStats, mockDelegations, mockCoOccurrences)

	// Test 'h' key shows help
	hKeyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}}
	_, cmd := model.Update(hKeyMsg)

	if !model.showHelp {
		t.Error("Help should be shown after 'h' key")
	}

	if cmd != nil {
		t.Error("Help toggle with 'h' should not trigger commands")
	}

	// Test 'h' key again hides help
	model.Update(hKeyMsg)
	if model.showHelp {
		t.Error("Help should be hidden after second 'h' key")
	}

	// Test mixed toggle (h then ? then h)
	model.Update(hKeyMsg) // Show with h
	if !model.showHelp {
		t.Error("Help should be shown with 'h' key")
	}

	questionMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}
	model.Update(questionMsg) // Hide with ?
	if model.showHelp {
		t.Error("Help should be hidden with '?' key")
	}

	model.Update(hKeyMsg) // Show with h again
	if !model.showHelp {
		t.Error("Help should be shown again with 'h' key")
	}
}

// TestHelpSystemESCKey tests ESC key to close help
func TestHelpSystemESCKey(t *testing.T) {
	mockStats := createMockAgentStats()
	mockDelegations := createMockDelegationPatterns()
	mockCoOccurrences := createMockCoOccurrencePatterns()

	model := NewDashboardModel(mockStats, mockDelegations, mockCoOccurrences)

	// Show help
	helpMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}
	model.Update(helpMsg)

	if !model.showHelp {
		t.Error("Help should be shown")
	}

	// ESC should close help
	escMsg := tea.KeyMsg{Type: tea.KeyEsc}
	_, cmd := model.Update(escMsg)

	if model.showHelp {
		t.Error("Help should be hidden after ESC key")
	}

	// Should not trigger quit when help is shown
	if cmd != nil {
		t.Error("ESC should close help, not quit dashboard")
	}

	// ESC again should quit dashboard
	_, cmd = model.Update(escMsg)
	if cmd == nil {
		t.Error("ESC should quit dashboard when help is not shown")
	}
}

// TestHelpSystemNavigationBlocking tests that help blocks other navigation
func TestHelpSystemNavigationBlocking(t *testing.T) {
	mockStats := createMockAgentStats()
	mockDelegations := createMockDelegationPatterns()
	mockCoOccurrences := createMockCoOccurrencePatterns()

	model := NewDashboardModel(mockStats, mockDelegations, mockCoOccurrences)

	// Show help
	helpMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}
	model.Update(helpMsg)

	// Remember initial state
	initialPanel := model.focusedPanel
	initialStatsIndex := model.selectedStatsIndex
	initialTimeFilter := model.GetTimeFilter()

	// Navigation keys should be blocked when help is shown
	navigationKeys := []tea.KeyMsg{
		{Type: tea.KeyTab},
		{Type: tea.KeyDown},
		{Type: tea.KeyUp},
		{Type: tea.KeyLeft},
		{Type: tea.KeyRight},
		{Type: tea.KeyEnter},
		{Type: tea.KeyRunes, Runes: []rune{'j'}},
		{Type: tea.KeyRunes, Runes: []rune{'k'}},
		{Type: tea.KeyRunes, Runes: []rune{'c'}},
		{Type: tea.KeyRunes, Runes: []rune{'s'}},
		{Type: tea.KeyRunes, Runes: []rune{'d'}},
		{Type: tea.KeyRunes, Runes: []rune{'f'}},
		{Type: tea.KeyRunes, Runes: []rune{'r'}},
		{Type: tea.KeyRunes, Runes: []rune{'1'}},
		{Type: tea.KeyRunes, Runes: []rune{'2'}},
		{Type: tea.KeyRunes, Runes: []rune{'3'}},
		{Type: tea.KeyRunes, Runes: []rune{'4'}},
		{Type: tea.KeyRunes, Runes: []rune{'5'}},
		{Type: tea.KeyRunes, Runes: []rune{'6'}},
	}

	for _, keyMsg := range navigationKeys {
		_, cmd := model.Update(keyMsg)

		// State should not change
		if model.focusedPanel != initialPanel {
			t.Errorf("Panel focus should not change when help is shown (key: %v)", keyMsg)
		}
		if model.selectedStatsIndex != initialStatsIndex {
			t.Errorf("Selection should not change when help is shown (key: %v)", keyMsg)
		}
		if model.GetTimeFilter() != initialTimeFilter {
			t.Errorf("Time filter should not change when help is shown (key: %v)", keyMsg)
		}

		// Should not trigger commands
		if cmd != nil {
			t.Errorf("Navigation keys should not trigger commands when help is shown (key: %v)", keyMsg)
		}
	}

	// Help should still be shown
	if !model.showHelp {
		t.Error("Help should still be shown after blocked navigation attempts")
	}
}

// TestHelpSystemGlobalKeysStillWork tests that global keys work even with help shown
func TestHelpSystemGlobalKeysStillWork(t *testing.T) {
	mockStats := createMockAgentStats()
	mockDelegations := createMockDelegationPatterns()
	mockCoOccurrences := createMockCoOccurrencePatterns()

	model := NewDashboardModel(mockStats, mockDelegations, mockCoOccurrences)

	// Show help
	helpMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}
	model.Update(helpMsg)

	// Global keys should still work
	ctrlCMsg := tea.KeyMsg{Type: tea.KeyCtrlC}
	_, cmd := model.Update(ctrlCMsg)

	if cmd == nil {
		t.Error("Ctrl+C should still work when help is shown")
	}

	// Reset for next test
	model.showHelp = true

	quitMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	_, cmd = model.Update(quitMsg)

	if cmd == nil {
		t.Error("'q' key should still work when help is shown")
	}
}

// TestHelpSystemContent tests comprehensive help content
func TestHelpSystemContent(t *testing.T) {
	mockStats := createMockAgentStats()
	mockDelegations := createMockDelegationPatterns()
	mockCoOccurrences := createMockCoOccurrencePatterns()

	model := NewDashboardModel(mockStats, mockDelegations, mockCoOccurrences)

	// Show help
	helpMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}
	model.Update(helpMsg)

	view := model.View()

	// Test that help contains main sections
	expectedSections := []string{
		"Agent Dashboard Help",
		"Navigation Commands",
		"Timeline Panel Commands",
		"Sorting Commands",
		"Delegation Panel Commands",
		"Co-occurrence Panel Commands",
		"Data Commands",
		"Display Commands",
		"Exit Commands",
		"Panels",
	}

	for _, section := range expectedSections {
		if !strings.Contains(view, section) {
			t.Errorf("Help should contain section: %s", section)
		}
	}

	// Test navigation commands are documented
	navigationCommands := []string{
		"Tab",
		"↑/↓ or j/k",
		"Enter",
	}

	for _, command := range navigationCommands {
		if !strings.Contains(view, command) {
			t.Errorf("Help should document navigation command: %s", command)
		}
	}

	// Test timeline commands are documented
	timelineCommands := []string{
		"←/→",
		"Space",
		"r",
	}

	for _, command := range timelineCommands {
		if !strings.Contains(view, command) {
			t.Errorf("Help should document timeline command: %s", command)
		}
	}

	// Test sorting commands are documented
	sortingCommands := []string{
		"c",
		"s",
		"d",
	}

	for _, command := range sortingCommands {
		if !strings.Contains(view, command) {
			t.Errorf("Help should document sorting command: %s", command)
		}
	}

	// Test filter commands are documented
	filterCommands := []string{
		"f",
		"m",
		"i",
	}

	for _, command := range filterCommands {
		if !strings.Contains(view, command) {
			t.Errorf("Help should document filter command: %s", command)
		}
	}

	// Test exit commands are documented
	exitCommands := []string{
		"q",
		"Esc",
		"Ctrl+C",
	}

	for _, command := range exitCommands {
		if !strings.Contains(view, command) {
			t.Errorf("Help should document exit command: %s", command)
		}
	}

	// Test time filter shortcuts are documented
	timeFilterShortcuts := []string{
		"1-6",
	}

	for _, shortcut := range timeFilterShortcuts {
		if !strings.Contains(view, shortcut) {
			t.Errorf("Help should document time filter shortcut: %s", shortcut)
		}
	}
}

// TestHelpSystemPanelDescriptions tests panel descriptions in help
func TestHelpSystemPanelDescriptions(t *testing.T) {
	mockStats := createMockAgentStats()
	mockDelegations := createMockDelegationPatterns()
	mockCoOccurrences := createMockCoOccurrencePatterns()

	model := NewDashboardModel(mockStats, mockDelegations, mockCoOccurrences)

	// Show help
	helpMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}
	model.Update(helpMsg)

	view := model.View()

	// Test panel descriptions
	panelDescriptions := []string{
		"Message Timeline",
		"Agent Statistics",
		"Delegation Flow Patterns",
		"Agent Collaboration Matrix",
	}

	for _, description := range panelDescriptions {
		if !strings.Contains(view, description) {
			t.Errorf("Help should describe panel: %s", description)
		}
	}

	// Test timeline panel details
	timelineDetails := []string{
		"Horizontal timeline",
		"Color-coded by role",
		"Success indicators",
		"Interactive time navigation",
	}

	for _, detail := range timelineDetails {
		if !strings.Contains(view, detail) {
			t.Errorf("Help should contain timeline detail: %s", detail)
		}
	}

	// Test agent statistics details
	statsDetails := []string{
		"Leaderboard",
		"success rates",
		"performance",
	}

	for _, detail := range statsDetails {
		if !strings.Contains(view, detail) {
			t.Errorf("Help should contain stats detail: %s", detail)
		}
	}

	// Test delegation panel details
	delegationDetails := []string{
		"A → B flows",
		"frequency",
		"filtering",
		"sorting",
	}

	for _, detail := range delegationDetails {
		if !strings.Contains(view, detail) {
			t.Errorf("Help should contain delegation detail: %s", detail)
		}
	}

	// Test co-occurrence panel details
	coOccurrenceDetails := []string{
		"relationship visualization",
		"Matrix view",
		"Relationship view",
		"Insights view",
		"strength indicators",
	}

	for _, detail := range coOccurrenceDetails {
		if !strings.Contains(view, detail) {
			t.Errorf("Help should contain co-occurrence detail: %s", detail)
		}
	}
}

// TestHelpSystemStyling tests help screen styling and layout
func TestHelpSystemStyling(t *testing.T) {
	mockStats := createMockAgentStats()
	mockDelegations := createMockDelegationPatterns()
	mockCoOccurrences := createMockCoOccurrencePatterns()

	model := NewDashboardModel(mockStats, mockDelegations, mockCoOccurrences)

	// Show help
	helpMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}
	model.Update(helpMsg)

	view := model.View()

	// Help should be properly formatted
	if len(view) < 100 {
		t.Error("Help content should be comprehensive")
	}

	// Should contain border or visual separation
	if !strings.Contains(view, "─") && !strings.Contains(view, "━") && !strings.Contains(view, "═") {
		t.Error("Help should have visual formatting/borders")
	}

	// Should end with continuation prompt
	if !strings.Contains(view, "Press any key to continue") && !strings.Contains(view, "continue") {
		t.Error("Help should contain continuation prompt")
	}

	// Should be visually distinct from main dashboard
	model.showHelp = false
	dashboardView := model.View()

	if view == dashboardView {
		t.Error("Help view should be visually distinct from dashboard view")
	}
}

// TestHelpSystemContextSensitive tests context-sensitive help based on focused panel
func TestHelpSystemContextSensitive(t *testing.T) {
	mockStats := createMockAgentStats()
	mockDelegations := createMockDelegationPatterns()
	mockCoOccurrences := createMockCoOccurrencePatterns()

	model := NewDashboardModel(mockStats, mockDelegations, mockCoOccurrences)

	// Test help content for each panel focus
	panels := []struct {
		panel       PanelType
		expectedKey string
		description string
	}{
		{TimelinePanel, "←/→", "timeline navigation"},
		{StatsPanel, "c/s/d", "sorting commands"},
		{DelegationPanel, "f", "filter commands"},
		{CoOccurrencePanel, "m", "matrix mode commands"},
	}

	for _, test := range panels {
		model.focusedPanel = test.panel

		// Show help
		helpMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}
		model.Update(helpMsg)

		view := model.View()

		// Should contain panel-specific commands
		if !strings.Contains(view, test.expectedKey) {
			t.Errorf("Help for %v panel should mention %s for %s", test.panel, test.expectedKey, test.description)
		}

		// Hide help for next test
		model.Update(helpMsg)
	}
}

// TestHelpSystemKeyboardShortcutDocumentation tests all keyboard shortcuts are documented
func TestHelpSystemKeyboardShortcutDocumentation(t *testing.T) {
	mockStats := createMockAgentStats()
	mockDelegations := createMockDelegationPatterns()
	mockCoOccurrences := createMockCoOccurrencePatterns()

	model := NewDashboardModel(mockStats, mockDelegations, mockCoOccurrences)

	// Show help
	helpMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}
	model.Update(helpMsg)

	view := model.View()

	// Test all keyboard shortcuts from specification are documented
	keyboardShortcuts := []struct {
		key         string
		description string
	}{
		{"Tab", "panel navigation"},
		{"Shift+Tab", "reverse panel navigation"},
		{"↑↓", "vertical navigation"},
		{"←→", "horizontal navigation"},
		{"j/k", "vim-style navigation"},
		{"Enter", "select/drill down"},
		{"Esc", "go back/exit"},
		{"q", "quit"},
		{"r", "refresh"},
		{"h/?", "help toggle"},
		{"1-6", "time range filters"},
		{"c", "sort by count"},
		{"s", "sort by success rate"},
		{"d", "sort by duration"},
		{"f", "filter menu"},
		{"m", "matrix mode"},
		{"i", "insights view"},
		{"Space", "select event"},
		{"Ctrl+C", "force quit"},
	}

	for _, shortcut := range keyboardShortcuts {
		if !strings.Contains(view, shortcut.key) {
			t.Errorf("Help should document keyboard shortcut: %s for %s", shortcut.key, shortcut.description)
		}
	}
}

// TestHelpSystemTimeRangeFilterDocumentation tests time filter documentation
func TestHelpSystemTimeRangeFilterDocumentation(t *testing.T) {
	mockStats := createMockAgentStats()
	mockDelegations := createMockDelegationPatterns()
	mockCoOccurrences := createMockCoOccurrencePatterns()

	model := NewDashboardModel(mockStats, mockDelegations, mockCoOccurrences)

	// Show help
	helpMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}
	model.Update(helpMsg)

	view := model.View()

	// Test time range filters are documented with their meanings
	timeRangeFilters := []struct {
		key    string
		period string
	}{
		{"1", "1w"}, // 1 week
		{"2", "1m"}, // 1 month
		{"3", "3m"}, // 3 months
		{"4", "6m"}, // 6 months
		{"5", "1y"}, // 1 year
		{"6", "all"}, // all time
	}

	for _, filter := range timeRangeFilters {
		// Should document the key
		if !strings.Contains(view, filter.key) {
			t.Errorf("Help should document time filter key: %s", filter.key)
		}

		// Should document the time period
		if !strings.Contains(view, filter.period) {
			t.Errorf("Help should document time period: %s for key %s", filter.period, filter.key)
		}
	}

	// Should explain the filter functionality
	filterExplanations := []string{
		"Quick filter by time range",
		"1w, 1m, 3m, 6m, 1y, all",
	}

	for _, explanation := range filterExplanations {
		if !strings.Contains(view, explanation) {
			t.Errorf("Help should contain filter explanation: %s", explanation)
		}
	}
}

// TestHelpSystemColorCodingDocumentation tests color coding documentation
func TestHelpSystemColorCodingDocumentation(t *testing.T) {
	mockStats := createMockAgentStats()
	mockDelegations := createMockDelegationPatterns()
	mockCoOccurrences := createMockCoOccurrencePatterns()

	model := NewDashboardModel(mockStats, mockDelegations, mockCoOccurrences)

	// Show help
	helpMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}
	model.Update(helpMsg)

	view := model.View()

	// Test color coding documentation
	colorCoding := []string{
		"Blue (User)",
		"Pink (Assistant)",
		"Yellow (Tool)",
		"Gray (System)",
		"Green (≥90%)",
		"Yellow (≥70%)",
		"Red (<70%)",
		"● (success)",
		"○ (failure)",
		"◆ (selected)",
	}

	for _, coding := range colorCoding {
		if !strings.Contains(view, coding) {
			t.Errorf("Help should document color coding: %s", coding)
		}
	}

	// Test success rate indicators
	successIndicators := []string{
		"success rates",
		"Success indicators",
	}

	for _, indicator := range successIndicators {
		if !strings.Contains(view, indicator) {
			t.Errorf("Help should document success indicators: %s", indicator)
		}
	}
}

// TestHelpSystemDisplayModeDocumentation tests display mode documentation
func TestHelpSystemDisplayModeDocumentation(t *testing.T) {
	mockStats := createMockAgentStats()
	mockDelegations := createMockDelegationPatterns()
	mockCoOccurrences := createMockCoOccurrencePatterns()

	model := NewDashboardModel(mockStats, mockDelegations, mockCoOccurrences)

	// Show help
	helpMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}
	model.Update(helpMsg)

	view := model.View()

	// Test co-occurrence display modes are documented
	displayModes := []string{
		"Matrix view",
		"Relationship view",
		"Insights view",
		"matrix → relationships → insights → matrix",
	}

	for _, mode := range displayModes {
		if !strings.Contains(view, mode) {
			t.Errorf("Help should document display mode: %s", mode)
		}
	}

	// Test relationship strength documentation
	strengthIndicators := []string{
		"Strong/Medium/Weak/Minimal",
		"strength indicators",
	}

	for _, indicator := range strengthIndicators {
		if !strings.Contains(view, indicator) {
			t.Errorf("Help should document strength indicators: %s", indicator)
		}
	}
}

// TestHelpSystemQuickStartGuide tests quick start guidance for new users
func TestHelpSystemQuickStartGuide(t *testing.T) {
	mockStats := createMockAgentStats()
	mockDelegations := createMockDelegationPatterns()
	mockCoOccurrences := createMockCoOccurrencePatterns()

	model := NewDashboardModel(mockStats, mockDelegations, mockCoOccurrences)

	// Show help
	helpMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}
	model.Update(helpMsg)

	view := model.View()

	// Test quick start elements
	quickStartElements := []string{
		"Navigation Commands",
		"Tab",
		"switch panels",
		"navigate",
		"Enter",
		"View details",
	}

	for _, element := range quickStartElements {
		if !strings.Contains(view, element) {
			t.Errorf("Help should contain quick start element: %s", element)
		}
	}

	// Test that help is accessible and user-friendly
	userFriendlyElements := []string{
		"Commands",
		"Panel",
		"switch",
		"navigate",
		"Sort",
		"Filter",
	}

	for _, element := range userFriendlyElements {
		if !strings.Contains(view, element) {
			t.Errorf("Help should contain user-friendly element: %s", element)
		}
	}
}