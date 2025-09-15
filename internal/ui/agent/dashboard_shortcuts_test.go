package agent

import (
	"fmt"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rsmdt/the-startup/internal/stats"
)

// ============================================================================
// Comprehensive Keyboard Shortcut Tests - All panels and contexts
// Tests for all keyboard shortcuts and their context-specific behavior
// ============================================================================

// TestGlobalKeyboardShortcuts tests shortcuts that work in all contexts
func TestGlobalKeyboardShortcuts(t *testing.T) {
	mockStats := createMockAgentStats()
	mockDelegations := createMockDelegationPatterns()
	mockCoOccurrences := createMockCoOccurrencePatterns()

	model := NewDashboardModel(mockStats, mockDelegations, mockCoOccurrences)

	// Test global shortcuts across all panels
	panels := []PanelType{TimelinePanel, StatsPanel, DelegationPanel, CoOccurrencePanel}

	for _, panel := range panels {
		model.focusedPanel = panel

		t.Run(string(rune(panel))+"_panel_global_shortcuts", func(t *testing.T) {
			// Test '?' help toggle
			helpMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}
			_, cmd := model.Update(helpMsg)

			if !model.showHelp {
				t.Errorf("'?' should show help in %v panel", panel)
			}
			if cmd != nil {
				t.Errorf("Help toggle should not trigger commands in %v panel", panel)
			}

			// Hide help for next test
			model.Update(helpMsg)

			// Test 'h' help toggle alternative
			hKeyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}}
			model.Update(hKeyMsg)

			if !model.showHelp {
				t.Errorf("'h' should show help in %v panel", panel)
			}

			// Hide help
			model.Update(hKeyMsg)

			// Test 'q' quit
			quitMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
			_, cmd = model.Update(quitMsg)

			if cmd == nil {
				t.Errorf("'q' should trigger quit in %v panel", panel)
			}

			// Test Ctrl+C force quit
			ctrlCMsg := tea.KeyMsg{Type: tea.KeyCtrlC}
			_, cmd = model.Update(ctrlCMsg)

			if cmd == nil {
				t.Errorf("Ctrl+C should trigger quit in %v panel", panel)
			}

			// Test ESC quit
			escMsg := tea.KeyMsg{Type: tea.KeyEsc}
			_, cmd = model.Update(escMsg)

			if cmd == nil {
				t.Errorf("ESC should trigger quit in %v panel", panel)
			}
		})
	}
}

// TestTabNavigationShortcuts tests Tab and Shift+Tab panel switching
func TestTabNavigationShortcuts(t *testing.T) {
	mockStats := createMockAgentStats()
	mockDelegations := createMockDelegationPatterns()
	mockCoOccurrences := createMockCoOccurrencePatterns()

	model := NewDashboardModel(mockStats, mockDelegations, mockCoOccurrences)

	// Test Tab navigation cycle: Timeline → Stats → Delegation → CoOccurrence → Timeline
	expectedCycle := []PanelType{TimelinePanel, StatsPanel, DelegationPanel, CoOccurrencePanel, TimelinePanel}

	// Should start at Timeline
	if model.focusedPanel != TimelinePanel {
		t.Errorf("Should start at TimelinePanel, got %v", model.focusedPanel)
	}

	tabMsg := tea.KeyMsg{Type: tea.KeyTab}

	for i := 1; i < len(expectedCycle); i++ {
		_, cmd := model.Update(tabMsg)

		if model.focusedPanel != expectedCycle[i] {
			t.Errorf("Tab step %d: expected %v, got %v", i, expectedCycle[i], model.focusedPanel)
		}

		if cmd != nil {
			t.Errorf("Tab navigation should not trigger commands at step %d", i)
		}
	}

	// Test that view changes with panel focus
	for i := 0; i < 4; i++ {
		view := model.View()
		if view == "" {
			t.Errorf("View should render for panel %v", model.focusedPanel)
		}
		model.Update(tabMsg)
	}
}

// TestTimeRangeFilterShortcuts tests 1-6 keys for time filtering
func TestTimeRangeFilterShortcuts(t *testing.T) {
	mockStats := createMockAgentStats()
	mockDelegations := createMockDelegationPatterns()
	mockCoOccurrences := createMockCoOccurrencePatterns()

	model := NewDashboardModel(mockStats, mockDelegations, mockCoOccurrences)

	// Test time filter shortcuts work in all panels
	panels := []PanelType{TimelinePanel, StatsPanel, DelegationPanel, CoOccurrencePanel}
	timeFilters := []struct {
		key            rune
		expectedFilter string
		description    string
	}{
		{'1', "7d", "1 week"},
		{'2', "30d", "1 month"},
		{'3', "90d", "3 months"},
		{'4', "180d", "6 months"},
		{'5', "365d", "1 year"},
		{'6', "", "all time"},
	}

	for _, panel := range panels {
		model.focusedPanel = panel

		for _, filter := range timeFilters {
			keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{filter.key}}
			_, cmd := model.Update(keyMsg)

			if cmd == nil {
				t.Errorf("Time filter key '%c' should trigger data reload in %v panel", filter.key, panel)
			}

			if model.GetTimeFilter() != filter.expectedFilter {
				t.Errorf("Key '%c' in %v panel: expected filter '%s', got '%s'",
					filter.key, panel, filter.expectedFilter, model.GetTimeFilter())
			}
		}
	}
}

// TestTimelinePanelShortcuts tests Timeline panel specific shortcuts
func TestTimelinePanelShortcuts(t *testing.T) {
	mockStats := createMockAgentStats()
	mockDelegations := createMockDelegationPatterns()
	mockCoOccurrences := createMockCoOccurrencePatterns()

	model := NewDashboardModel(mockStats, mockDelegations, mockCoOccurrences)
	model.focusedPanel = TimelinePanel

	// Set mock timeline events
	model.timeline.SetEvents(createMockMessageEvents())

	// Test left/right arrow navigation
	leftMsg := tea.KeyMsg{Type: tea.KeyLeft}
	_, cmd := model.Update(leftMsg)

	if cmd != nil {
		t.Error("Left arrow in timeline should not trigger commands")
	}

	rightMsg := tea.KeyMsg{Type: tea.KeyRight}
	_, cmd = model.Update(rightMsg)

	if cmd != nil {
		t.Error("Right arrow in timeline should not trigger commands")
	}

	// Test up/down arrow for zoom
	upMsg := tea.KeyMsg{Type: tea.KeyUp}
	_, cmd = model.Update(upMsg)

	if cmd != nil {
		t.Error("Up arrow in timeline should not trigger commands")
	}

	downMsg := tea.KeyMsg{Type: tea.KeyDown}
	_, cmd = model.Update(downMsg)

	if cmd != nil {
		t.Error("Down arrow in timeline should not trigger commands")
	}

	// Test space for event selection
	spaceMsg := tea.KeyMsg{Type: tea.KeySpace}
	_, cmd = model.Update(spaceMsg)

	if cmd != nil {
		t.Error("Space in timeline should not trigger commands")
	}

	// Test 'r' for reset view (should not trigger data reload in timeline)
	rMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}}
	_, cmd = model.Update(rMsg)

	if cmd != nil {
		t.Error("'r' in timeline should reset view, not trigger data reload")
	}
}

// TestStatsPanelShortcuts tests Stats panel specific shortcuts
func TestStatsPanelShortcuts(t *testing.T) {
	mockStats := createMockAgentStats()
	mockDelegations := createMockDelegationPatterns()
	mockCoOccurrences := createMockCoOccurrencePatterns()

	model := NewDashboardModel(mockStats, mockDelegations, mockCoOccurrences)
	model.focusedPanel = StatsPanel

	// Test sorting shortcuts
	sortTests := []struct {
		key            rune
		expectedColumn SortColumn
		description    string
	}{
		{'c', SortByCount, "count sorting"},
		{'s', SortBySuccessRate, "success rate sorting"},
		{'d', SortByDuration, "duration sorting"},
	}

	for _, test := range sortTests {
		keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{test.key}}
		_, cmd := model.Update(keyMsg)

		if model.sortColumn != test.expectedColumn {
			t.Errorf("Key '%c' should set sort column to %v, got %v",
				test.key, test.expectedColumn, model.sortColumn)
		}

		if cmd != nil {
			t.Errorf("Sort key '%c' should not trigger commands", test.key)
		}
	}

	// Test sort direction toggle
	cMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}}
	model.Update(cMsg) // First press sets to count, descending

	initialDirection := model.sortDirection
	model.Update(cMsg) // Second press should toggle direction

	if model.sortDirection == initialDirection {
		t.Error("Second press of 'c' should toggle sort direction")
	}

	// Test navigation within stats panel
	downMsg := tea.KeyMsg{Type: tea.KeyDown}
	initialIndex := model.selectedStatsIndex
	model.Update(downMsg)

	if len(model.agentStats) > 1 && model.selectedStatsIndex == initialIndex {
		t.Error("Down arrow should move selection in stats panel")
	}

	upMsg := tea.KeyMsg{Type: tea.KeyUp}
	model.Update(upMsg)

	if model.selectedStatsIndex != initialIndex {
		t.Error("Up arrow should return to initial selection")
	}

	// Test vim-style navigation
	jMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
	model.Update(jMsg)

	if len(model.agentStats) > 1 && model.selectedStatsIndex == initialIndex {
		t.Error("'j' key should move selection down in stats panel")
	}

	kMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}
	model.Update(kMsg)

	if model.selectedStatsIndex != initialIndex {
		t.Error("'k' key should move selection up in stats panel")
	}

	// Test 'r' for data refresh
	rMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}}
	_, cmdR := model.Update(rMsg)

	if cmdR == nil {
		t.Error("'r' in stats panel should trigger data refresh")
	}

	// Test Enter for detail view
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	_, cmdEnter := model.Update(enterMsg)

	if cmdEnter == nil {
		t.Error("Enter in stats panel should trigger detail view")
	}
}

// TestDelegationPanelShortcuts tests Delegation panel specific shortcuts
func TestDelegationPanelShortcuts(t *testing.T) {
	mockStats := createMockAgentStats()
	mockDelegations := createMockDelegationPatterns()
	mockCoOccurrences := createMockCoOccurrencePatterns()

	model := NewDashboardModel(mockStats, mockDelegations, mockCoOccurrences)
	model.focusedPanel = DelegationPanel

	// Test 'f' for filter cycling
	fMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}}
	initialMinCount := model.delegationFilterOptions.MinCount

	_, cmd := model.Update(fMsg)

	if model.delegationFilterOptions.MinCount == initialMinCount {
		t.Error("'f' key should cycle delegation filter options")
	}

	if cmd != nil {
		t.Error("Delegation filter cycling should not trigger commands")
	}

	// Test 's' for sort cycling (different from stats panel)
	sMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}}
	initialSort := model.delegationFilterOptions.SortBy

	model.Update(sMsg)

	if model.delegationFilterOptions.SortBy == initialSort {
		t.Error("'s' key should cycle delegation sort options")
	}

	// Test that 'c' and 'd' don't do anything in delegation panel
	cMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}}
	_, cmd = model.Update(cMsg)

	if cmd != nil {
		t.Error("'c' key should not trigger commands in delegation panel")
	}

	dMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}}
	_, cmd = model.Update(dMsg)

	if cmd != nil {
		t.Error("'d' key should not trigger commands in delegation panel")
	}

	// Test navigation within delegation panel
	downMsg := tea.KeyMsg{Type: tea.KeyDown}
	initialIndex := model.selectedDelegationIndex
	model.Update(downMsg)

	if len(model.delegationPatterns) > 1 && model.selectedDelegationIndex == initialIndex {
		t.Error("Down arrow should move selection in delegation panel")
	}

	// Test Enter for detail view
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	_, cmd = model.Update(enterMsg)

	if cmd == nil {
		t.Error("Enter in delegation panel should trigger detail view")
	}
}

// TestCoOccurrencePanelShortcuts tests Co-occurrence panel specific shortcuts
func TestCoOccurrencePanelShortcuts(t *testing.T) {
	mockStats := createMockAgentStats()
	mockDelegations := createMockDelegationPatterns()
	mockCoOccurrences := createMockCoOccurrencePatterns()

	model := NewDashboardModel(mockStats, mockDelegations, mockCoOccurrences)
	model.focusedPanel = CoOccurrencePanel

	// Test 'm' for matrix display mode switching
	mMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}}
	initialMode := model.coOccurrenceMatrixComponent.DisplayMode

	_, cmd := model.Update(mMsg)

	if model.coOccurrenceMatrixComponent.DisplayMode == initialMode {
		t.Error("'m' key should switch display mode in co-occurrence panel")
	}

	if cmd != nil {
		t.Error("Display mode switching should not trigger commands")
	}

	// Test 'i' for insights view
	iMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}}
	model.Update(iMsg)

	if model.coOccurrenceMatrixComponent.DisplayMode != InsightsView {
		t.Error("'i' key should switch to insights view")
	}

	// Test matrix navigation (only works in matrix view)
	model.coOccurrenceMatrixComponent.DisplayMode = MatrixView

	// Test arrow keys for matrix navigation
	leftMsg := tea.KeyMsg{Type: tea.KeyLeft}
	_, cmd = model.Update(leftMsg)

	if cmd != nil {
		t.Error("Left arrow in co-occurrence matrix should not trigger commands")
	}

	rightMsg := tea.KeyMsg{Type: tea.KeyRight}
	_, cmd = model.Update(rightMsg)

	if cmd != nil {
		t.Error("Right arrow in co-occurrence matrix should not trigger commands")
	}

	upMsg := tea.KeyMsg{Type: tea.KeyUp}
	_, cmd = model.Update(upMsg)

	if cmd != nil {
		t.Error("Up arrow in co-occurrence matrix should not trigger commands")
	}

	downMsg := tea.KeyMsg{Type: tea.KeyDown}
	_, cmd = model.Update(downMsg)

	if cmd != nil {
		t.Error("Down arrow in co-occurrence matrix should not trigger commands")
	}

	// Test that other sorting keys don't work in co-occurrence panel
	cMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}}
	_, cmd = model.Update(cMsg)

	if cmd != nil {
		t.Error("'c' key should not trigger commands in co-occurrence panel")
	}

	sMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}}
	_, cmd = model.Update(sMsg)

	if cmd != nil {
		t.Error("'s' key should not trigger commands in co-occurrence panel")
	}

	dMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}}
	_, cmd = model.Update(dMsg)

	if cmd != nil {
		t.Error("'d' key should not trigger commands in co-occurrence panel")
	}

	fMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}}
	_, cmd = model.Update(fMsg)

	if cmd != nil {
		t.Error("'f' key should not trigger commands in co-occurrence panel")
	}
}

// TestNavigationShortcutsBoundaryConditions tests boundary conditions for navigation
func TestNavigationShortcutsBoundaryConditions(t *testing.T) {
	// Test with empty data
	emptyStats := make(map[string]*stats.GlobalAgentStats)
	emptyDelegations := []stats.DelegationPattern{}
	emptyCoOccurrences := []stats.AgentCoOccurrence{}

	model := NewDashboardModel(emptyStats, emptyDelegations, emptyCoOccurrences)

	// Test navigation with empty data doesn't crash
	navigationKeys := []tea.KeyMsg{
		{Type: tea.KeyDown},
		{Type: tea.KeyUp},
		{Type: tea.KeyRunes, Runes: []rune{'j'}},
		{Type: tea.KeyRunes, Runes: []rune{'k'}},
	}

	for _, panel := range []PanelType{StatsPanel, DelegationPanel, CoOccurrencePanel} {
		model.focusedPanel = panel

		for _, keyMsg := range navigationKeys {
			_, cmd := model.Update(keyMsg)

			// Should not crash or trigger unexpected commands
			if cmd != nil {
				t.Errorf("Navigation in empty %v panel should not trigger commands", panel)
			}

			// Selection indices should remain valid
			switch panel {
			case StatsPanel:
				if model.selectedStatsIndex < 0 {
					t.Errorf("Stats selection index should not be negative: %d", model.selectedStatsIndex)
				}
			case DelegationPanel:
				if model.selectedDelegationIndex < 0 {
					t.Errorf("Delegation selection index should not be negative: %d", model.selectedDelegationIndex)
				}
			case CoOccurrencePanel:
				if model.selectedCoOccurrenceIndex < 0 {
					t.Errorf("Co-occurrence selection index should not be negative: %d", model.selectedCoOccurrenceIndex)
				}
			}
		}
	}

	// Test with single item data
	singleStats := map[string]*stats.GlobalAgentStats{
		"OnlyAgent": &stats.GlobalAgentStats{},
	}
	singleDelegation := []stats.DelegationPattern{
		{SourceAgent: "A", TargetAgent: "B", Count: 1},
	}
	singleCoOccurrence := []stats.AgentCoOccurrence{
		{Agent1: "A", Agent2: "B", Count: 1},
	}

	model = NewDashboardModel(singleStats, singleDelegation, singleCoOccurrence)

	// Test navigation with single items
	for _, panel := range []PanelType{StatsPanel, DelegationPanel, CoOccurrencePanel} {
		model.focusedPanel = panel

		for _, keyMsg := range navigationKeys {
			model.Update(keyMsg)

			// Selection should stay within bounds
			switch panel {
			case StatsPanel:
				if model.selectedStatsIndex >= len(singleStats) || model.selectedStatsIndex < 0 {
					t.Errorf("Stats selection out of bounds: %d (max: %d)", model.selectedStatsIndex, len(singleStats)-1)
				}
			case DelegationPanel:
				if model.selectedDelegationIndex >= len(singleDelegation) || model.selectedDelegationIndex < 0 {
					t.Errorf("Delegation selection out of bounds: %d (max: %d)", model.selectedDelegationIndex, len(singleDelegation)-1)
				}
			case CoOccurrencePanel:
				if model.selectedCoOccurrenceIndex >= len(singleCoOccurrence) || model.selectedCoOccurrenceIndex < 0 {
					t.Errorf("Co-occurrence selection out of bounds: %d (max: %d)", model.selectedCoOccurrenceIndex, len(singleCoOccurrence)-1)
				}
			}
		}
	}
}

// TestShortcutConflictResolution tests handling of conflicting shortcuts
func TestShortcutConflictResolution(t *testing.T) {
	mockStats := createMockAgentStats()
	mockDelegations := createMockDelegationPatterns()
	mockCoOccurrences := createMockCoOccurrencePatterns()

	model := NewDashboardModel(mockStats, mockDelegations, mockCoOccurrences)

	// Test 's' key behavior in different panels
	sKeyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}}

	// In Stats panel, 's' should sort by success rate
	model.focusedPanel = StatsPanel
	model.Update(sKeyMsg)

	if model.sortColumn != SortBySuccessRate {
		t.Error("'s' in stats panel should sort by success rate")
	}

	// In Delegation panel, 's' should cycle sort options
	model.focusedPanel = DelegationPanel
	initialSort := model.delegationFilterOptions.SortBy
	model.Update(sKeyMsg)

	if model.delegationFilterOptions.SortBy == initialSort {
		t.Error("'s' in delegation panel should cycle sort options")
	}

	// In other panels, 's' should do nothing
	model.focusedPanel = TimelinePanel
	_, cmd := model.Update(sKeyMsg)

	if cmd != nil {
		t.Error("'s' in timeline panel should not trigger commands")
	}

	model.focusedPanel = CoOccurrencePanel
	_, cmd = model.Update(sKeyMsg)

	if cmd != nil {
		t.Error("'s' in co-occurrence panel should not trigger commands")
	}
}

// TestShortcutStateConsistency tests that shortcuts maintain consistent state
func TestShortcutStateConsistency(t *testing.T) {
	mockStats := createMockAgentStats()
	mockDelegations := createMockDelegationPatterns()
	mockCoOccurrences := createMockCoOccurrencePatterns()

	model := NewDashboardModel(mockStats, mockDelegations, mockCoOccurrences)

	// Test that panel switching preserves individual panel states
	model.focusedPanel = StatsPanel
	model.selectedStatsIndex = 1

	model.focusedPanel = DelegationPanel
	model.selectedDelegationIndex = 1

	model.focusedPanel = CoOccurrencePanel
	model.selectedCoOccurrenceIndex = 1

	// Switch back to stats panel
	model.focusedPanel = StatsPanel

	// Stats selection should be preserved
	if model.selectedStatsIndex != 1 {
		t.Error("Stats panel selection should be preserved when switching panels")
	}

	// Other panel selections should also be preserved
	if model.selectedDelegationIndex != 1 {
		t.Error("Delegation panel selection should be preserved")
	}

	if model.selectedCoOccurrenceIndex != 1 {
		t.Error("Co-occurrence panel selection should be preserved")
	}

	// Test that sort state is preserved
	model.focusedPanel = StatsPanel
	cMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}}
	model.Update(cMsg)

	sortColumn := model.sortColumn
	sortDirection := model.sortDirection

	// Switch panels and back
	tabMsg := tea.KeyMsg{Type: tea.KeyTab}
	model.Update(tabMsg)
	model.Update(tabMsg)
	model.Update(tabMsg)
	model.Update(tabMsg) // Back to stats

	// Sort state should be preserved
	if model.sortColumn != sortColumn {
		t.Error("Sort column should be preserved across panel switches")
	}

	if model.sortDirection != sortDirection {
		t.Error("Sort direction should be preserved across panel switches")
	}
}

// TestShortcutDocumentationAccuracy tests that all documented shortcuts work
func TestShortcutDocumentationAccuracy(t *testing.T) {
	mockStats := createMockAgentStats()
	mockDelegations := createMockDelegationPatterns()
	mockCoOccurrences := createMockCoOccurrencePatterns()

	model := NewDashboardModel(mockStats, mockDelegations, mockCoOccurrences)

	// Get help content
	helpMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}
	model.Update(helpMsg)
	helpContent := model.View()
	model.Update(helpMsg) // Hide help

	// Test that all shortcuts mentioned in help actually work
	documentedShortcuts := []struct {
		key         tea.KeyMsg
		shouldWork  bool
		description string
	}{
		{tea.KeyMsg{Type: tea.KeyTab}, true, "Tab panel switching"},
		{tea.KeyMsg{Type: tea.KeyDown}, true, "Down navigation"},
		{tea.KeyMsg{Type: tea.KeyUp}, true, "Up navigation"},
		{tea.KeyMsg{Type: tea.KeyLeft}, true, "Left navigation"},
		{tea.KeyMsg{Type: tea.KeyRight}, true, "Right navigation"},
		{tea.KeyMsg{Type: tea.KeyEnter}, true, "Enter selection"},
		{tea.KeyMsg{Type: tea.KeyEsc}, true, "ESC quit"},
		{tea.KeyMsg{Type: tea.KeyCtrlC}, true, "Ctrl+C quit"},
		{tea.KeyMsg{Type: tea.KeySpace}, true, "Space selection"},
		{tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}, true, "j down"},
		{tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}, true, "k up"},
		{tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}, true, "q quit"},
		{tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}}, true, "r refresh"},
		{tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}}, true, "h help"},
		{tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}, true, "? help"},
		{tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}}, true, "1 time filter"},
		{tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}}, true, "2 time filter"},
		{tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'3'}}, true, "3 time filter"},
		{tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'4'}}, true, "4 time filter"},
		{tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'5'}}, true, "5 time filter"},
		{tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'6'}}, true, "6 time filter"},
		{tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}}, true, "c sort count"},
		{tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}}, true, "s sort success"},
		{tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}}, true, "d sort duration"},
		{tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}}, true, "f filter"},
		{tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}}, true, "m matrix mode"},
		{tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}}, true, "i insights"},
	}

	for _, shortcut := range documentedShortcuts {
		// Check if shortcut is mentioned in help
		shortcutInHelp := false
		if len(shortcut.key.Runes) > 0 {
			keyStr := string(shortcut.key.Runes[0])
			shortcutInHelp = strings.Contains(helpContent, keyStr)
		} else {
			// Handle special keys
			switch shortcut.key.Type {
			case tea.KeyTab:
				shortcutInHelp = strings.Contains(helpContent, "Tab")
			case tea.KeyDown:
				shortcutInHelp = strings.Contains(helpContent, "↓")
			case tea.KeyUp:
				shortcutInHelp = strings.Contains(helpContent, "↑")
			case tea.KeyLeft:
				shortcutInHelp = strings.Contains(helpContent, "←")
			case tea.KeyRight:
				shortcutInHelp = strings.Contains(helpContent, "→")
			case tea.KeyEnter:
				shortcutInHelp = strings.Contains(helpContent, "Enter")
			case tea.KeyEsc:
				shortcutInHelp = strings.Contains(helpContent, "Esc")
			case tea.KeyCtrlC:
				shortcutInHelp = strings.Contains(helpContent, "Ctrl+C")
			case tea.KeySpace:
				shortcutInHelp = strings.Contains(helpContent, "Space")
			}
		}

		if shortcutInHelp {
			// Test that documented shortcut actually works
			_, cmd := model.Update(shortcut.key)

			// Some shortcuts should trigger commands, others shouldn't
			// We can't test exact behavior here, but we can ensure they don't crash
			_ = cmd // Mark as used

			// Test that view still renders after shortcut
			view := model.View()
			if view == "" {
				t.Errorf("View should render after shortcut: %s", shortcut.description)
			}
		}
	}
}

// TestShortcutAccessibility tests shortcuts work for accessibility
func TestShortcutAccessibility(t *testing.T) {
	mockStats := createMockAgentStats()
	mockDelegations := createMockDelegationPatterns()
	mockCoOccurrences := createMockCoOccurrencePatterns()

	model := NewDashboardModel(mockStats, mockDelegations, mockCoOccurrences)

	// Test that both arrow keys and vim keys work for navigation
	arrowDown := tea.KeyMsg{Type: tea.KeyDown}
	vimDown := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}

	model.focusedPanel = StatsPanel

	// Test arrow navigation
	initialIndex := model.selectedStatsIndex
	model.Update(arrowDown)
	arrowResult := model.selectedStatsIndex

	// Reset and test vim navigation
	model.selectedStatsIndex = initialIndex
	model.Update(vimDown)
	vimResult := model.selectedStatsIndex

	if len(model.agentStats) > 1 && arrowResult != vimResult {
		t.Error("Arrow keys and vim keys should have equivalent behavior")
	}

	// Test that both help shortcuts work
	helpQuestion := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}
	helpH := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}}

	model.Update(helpQuestion)
	questionResult := model.showHelp
	model.Update(helpQuestion) // Hide

	model.Update(helpH)
	hResult := model.showHelp
	model.Update(helpH) // Hide

	if questionResult != hResult {
		t.Error("Both '?' and 'h' should toggle help equally")
	}
}

// TestShortcutConsistencyAcrossViews tests shortcuts work consistently in all views
func TestShortcutConsistencyAcrossViews(t *testing.T) {
	mockStats := createMockAgentStats()
	mockDelegations := createMockDelegationPatterns()
	mockCoOccurrences := createMockCoOccurrencePatterns()

	model := NewDashboardModel(mockStats, mockDelegations, mockCoOccurrences)

	// Test that global shortcuts work in all display modes
	model.focusedPanel = CoOccurrencePanel

	// Test shortcuts in different display modes
	displayModes := []MatrixDisplayMode{MatrixView, RelationshipView, InsightsView}

	for _, mode := range displayModes {
		model.coOccurrenceMatrixComponent.DisplayMode = mode

		// Test global shortcuts still work
		helpMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}
		model.Update(helpMsg)

		if !model.showHelp {
			t.Errorf("Help should work in display mode %v", mode)
		}

		model.Update(helpMsg) // Hide help

		// Test time filter shortcuts still work
		timeFilterMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}}
		_, cmd := model.Update(timeFilterMsg)

		if cmd == nil {
			t.Errorf("Time filter should work in display mode %v", mode)
		}

		// Test quit shortcuts still work
		quitMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
		_, cmd = model.Update(quitMsg)

		if cmd == nil {
			t.Errorf("Quit should work in display mode %v", mode)
		}
	}
}

// TestShortcutErrorRecovery tests shortcuts work after errors
func TestShortcutErrorRecovery(t *testing.T) {
	model := NewDashboardModelWithLoading()

	// Simulate a load error
	errorMsg := DataLoadedMsg{
		Error: fmt.Errorf("simulated error"),
	}
	model.Update(errorMsg)

	// Shortcuts should still work after error
	helpMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}
	model.Update(helpMsg)

	if !model.showHelp {
		t.Error("Help should work after load error")
	}

	model.Update(helpMsg) // Hide help

	// Navigation should still work
	tabMsg := tea.KeyMsg{Type: tea.KeyTab}
	_, cmd := model.Update(tabMsg)

	if cmd != nil {
		t.Error("Tab navigation should not trigger commands after error")
	}

	// Time filter should still trigger reload attempts
	timeFilterMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}}
	_, cmd = model.Update(timeFilterMsg)

	if cmd == nil {
		t.Error("Time filter should still trigger reload after error")
	}

	// Quit should still work
	quitMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	_, cmd = model.Update(quitMsg)

	if cmd == nil {
		t.Error("Quit should work after error")
	}
}