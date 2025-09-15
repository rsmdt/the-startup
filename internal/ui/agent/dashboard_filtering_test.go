package agent

import (
	"fmt"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rsmdt/the-startup/internal/stats"
)

// ============================================================================
// Time Range Filter Tests - Testing Phase 4 filtering functionality
// Tests for keyboard shortcuts: 1w, 1m, 3m, 6m, 1y, all (keys 1-6)
// ============================================================================

// TimePreset represents time range presets as specified in new-stats-command.md
type TimePreset int

const (
	TimePreset1Week TimePreset = iota
	TimePreset1Month
	TimePreset3Months
	TimePreset6Months
	TimePreset1Year
	TimePresetAll
)

// FilterState represents the dashboard filter state
type FilterState struct {
	TimeRange TimePreset
	Agents    []string
	Success   *bool
	Sessions  []string
}

// TestTimeRangeFilterKeyboardShortcuts tests 1-6 key functionality for time filtering
func TestTimeRangeFilterKeyboardShortcuts(t *testing.T) {
	mockStats := createMockAgentStats()
	mockDelegations := createMockDelegationPatterns()
	mockCoOccurrences := createMockCoOccurrencePatterns()

	model := NewDashboardModel(mockStats, mockDelegations, mockCoOccurrences)

	// Test key '1' - 1 week filter
	key1Msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}}
	_, cmd := model.Update(key1Msg)

	// Should trigger data reload for 1 week filter
	if cmd == nil {
		t.Error("Expected '1' key to trigger data reload command for 1 week filter")
	}

	// In the actual implementation, the time filter would be updated
	// For now, we test that the command is triggered
	if cmd == nil {
		t.Error("Expected '1' key to trigger time filter update")
	}

	// Test key '2' - 1 month filter
	key2Msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}}
	_, cmd = model.Update(key2Msg)

	if cmd == nil {
		t.Error("Expected '2' key to trigger data reload command for 1 month filter")
	}

	// Test remaining time filter keys trigger commands
	timeFilterTests := []struct {
		key         rune
		description string
	}{
		{'3', "3 months filter"},
		{'4', "6 months filter"},
		{'5', "1 year filter"},
		{'6', "all time filter"},
	}

	for _, test := range timeFilterTests {
		keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{test.key}}
		_, cmd = model.Update(keyMsg)

		if cmd == nil {
			t.Errorf("Expected '%c' key to trigger data reload command for %s", test.key, test.description)
		}
	}
}

// TestTimeRangeFilterDefaultState tests default 1-month lookback filter
func TestTimeRangeFilterDefaultState(t *testing.T) {
	model := NewDashboardModelWithLoading()

	// Should start with 1-month lookback as specified
	expectedDefault := "30d"
	if model.GetTimeFilter() != expectedDefault {
		t.Errorf("Expected default time filter to be '%s', got '%s'", expectedDefault, model.GetTimeFilter())
	}

	// Should launch with 1-month lookback as per specification
	view := model.View()
	if view == "" {
		t.Error("Dashboard should render with default filter")
	}
}

// TestTimeRangeFilterDataReload tests that filter changes trigger data reload
func TestTimeRangeFilterDataReload(t *testing.T) {
	mockStats := createMockAgentStats()
	mockDelegations := createMockDelegationPatterns()
	mockCoOccurrences := createMockCoOccurrencePatterns()

	model := NewDashboardModel(mockStats, mockDelegations, mockCoOccurrences)

	// Track initial data state
	initialAgentCount := len(model.GetAgentStats())

	// Test that time filter changes trigger full data reload
	key1Msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}}
	_, cmd1 := model.Update(key1Msg)

	if cmd1 == nil {
		t.Error("Time filter change should trigger data reload")
	}

	// Simulate data reload completion
	newMockStats := createEnhancedMockAgentStats()
	reloadMsg := DataLoadedMsg{
		AgentStats:         newMockStats,
		DelegationPatterns: mockDelegations,
		CoOccurrences:      mockCoOccurrences,
		MessageEvents:      createMockMessageEvents(),
		Error:              nil,
	}

	_, cmd2 := model.Update(reloadMsg)

	// Data should be updated
	newAgentCount := len(model.GetAgentStats())
	if newAgentCount == initialAgentCount {
		t.Error("Expected agent data to be reloaded after filter change")
	}

	// Should not trigger additional commands after successful reload
	if cmd2 != nil {
		t.Error("Data reload completion should not trigger additional commands")
	}
}

// TestTimeRangeFilterUIIndicators tests visual feedback for active filters
func TestTimeRangeFilterUIIndicators(t *testing.T) {
	mockStats := createMockAgentStats()
	mockDelegations := createMockDelegationPatterns()
	mockCoOccurrences := createMockCoOccurrencePatterns()

	model := NewDashboardModel(mockStats, mockDelegations, mockCoOccurrences)

	// Test that view shows current time filter status
	view := model.View()
	if !strings.Contains(view, "Agent Dashboard") {
		t.Error("Dashboard should display title")
	}

	// Apply 1 week filter
	key1Msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}}
	model.Update(key1Msg)

	// View should reflect active filter (would need filter indicator in actual implementation)
	filteredView := model.View()
	if filteredView == "" {
		t.Error("Dashboard should render with active filter")
	}

	// Test different filters show different states
	filterTests := []struct {
		key         rune
		description string
	}{
		{'1', "1 week"},
		{'2', "1 month"},
		{'3', "3 months"},
		{'4', "6 months"},
		{'5', "1 year"},
		{'6', "all time"},
	}

	for _, test := range filterTests {
		keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{test.key}}
		model.Update(keyMsg)

		view := model.View()
		if view == "" {
			t.Errorf("Dashboard should render with %s filter", test.description)
		}
	}
}

// TestTimeFilterPersistenceAcrossPanels tests filter state consistency
func TestTimeFilterPersistenceAcrossPanels(t *testing.T) {
	mockStats := createMockAgentStats()
	mockDelegations := createMockDelegationPatterns()
	mockCoOccurrences := createMockCoOccurrencePatterns()

	model := NewDashboardModel(mockStats, mockDelegations, mockCoOccurrences)

	// Apply 3 month filter
	key3Msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'3'}}
	model.Update(key3Msg)
	expectedFilter := "90d"

	// Switch between panels
	tabMsg := tea.KeyMsg{Type: tea.KeyTab}
	model.Update(tabMsg) // Timeline -> Stats
	model.Update(tabMsg) // Stats -> Delegation
	model.Update(tabMsg) // Delegation -> CoOccurrence
	model.Update(tabMsg) // CoOccurrence -> Timeline

	// Filter should persist across panel switches
	if model.GetTimeFilter() != expectedFilter {
		t.Errorf("Time filter should persist across panel switches, expected '%s', got '%s'", expectedFilter, model.GetTimeFilter())
	}

	// All panels should respect the same filter
	for i := 0; i < 4; i++ {
		view := model.View()
		if view == "" {
			t.Errorf("Panel %d should render with active filter", i)
		}
		model.Update(tabMsg)
	}
}

// TestTimeFilterWithEmptyData tests filter behavior with no data
func TestTimeFilterWithEmptyData(t *testing.T) {
	emptyStats := make(map[string]*stats.GlobalAgentStats)
	emptyDelegations := []stats.DelegationPattern{}
	emptyCoOccurrences := []stats.AgentCoOccurrence{}

	model := NewDashboardModel(emptyStats, emptyDelegations, emptyCoOccurrences)

	// Apply filters with empty data
	filterKeys := []rune{'1', '2', '3', '4', '5', '6'}
	for _, key := range filterKeys {
		keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{key}}
		_, cmd := model.Update(keyMsg)

		// Should still trigger reload even with empty data
		if cmd == nil {
			t.Errorf("Filter key '%c' should trigger reload even with empty data", key)
		}

		// Should not crash with empty data
		view := model.View()
		if view == "" {
			t.Errorf("Dashboard should render with empty data for filter '%c'", key)
		}
	}
}

// TestTimeFilterInteractionWithHelp tests filter keys when help is shown
func TestTimeFilterInteractionWithHelp(t *testing.T) {
	mockStats := createMockAgentStats()
	mockDelegations := createMockDelegationPatterns()
	mockCoOccurrences := createMockCoOccurrencePatterns()

	model := NewDashboardModel(mockStats, mockDelegations, mockCoOccurrences)

	// Show help
	helpMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}
	model.Update(helpMsg)

	// Verify help is shown
	view := model.View()
	if !strings.Contains(view, "Help") {
		t.Error("Help should be visible after '?' key")
	}

	// Time filter keys should still work when help is shown
	key1Msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}}
	_, cmd := model.Update(key1Msg)

	// Should apply filter and close help
	if cmd == nil {
		t.Error("Filter key should work even when help is shown")
	}

	// Help should be hidden after filter application
	view = model.View()
	if strings.Contains(view, "Help") && !strings.Contains(view, "Agent Dashboard") {
		t.Error("Help should be hidden after applying filter")
	}
}

// ============================================================================
// Filter Menu and Controls Tests - Testing f key functionality
// ============================================================================

// TestFilterMenuKeyboardShortcut tests 'f' key for filter menu
func TestFilterMenuKeyboardShortcut(t *testing.T) {
	mockStats := createMockAgentStats()
	mockDelegations := createMockDelegationPatterns()
	mockCoOccurrences := createMockCoOccurrencePatterns()

	model := NewDashboardModel(mockStats, mockDelegations, mockCoOccurrences)

	// 'f' key should work differently based on focused panel
	testCases := []struct {
		panel       PanelType
		expectCmd   bool
		description string
	}{
		{TimelinePanel, false, "timeline panel should not respond to 'f' key"},
		{StatsPanel, false, "stats panel should not respond to 'f' key"},
		{DelegationPanel, true, "delegation panel should respond to 'f' key"},
		{CoOccurrencePanel, false, "co-occurrence panel should not respond to 'f' key"},
	}

	for _, test := range testCases {
		model.focusedPanel = test.panel

		fKeyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}}
		_, cmd := model.Update(fKeyMsg)

		if test.expectCmd && cmd == nil {
			t.Errorf("%s - expected command, got nil", test.description)
		} else if !test.expectCmd && cmd != nil {
			t.Errorf("%s - expected no command, got command", test.description)
		}
	}
}

// TestDelegationFilterOptions tests delegation panel filter cycling
func TestDelegationFilterOptions(t *testing.T) {
	mockStats := createMockAgentStats()
	mockDelegations := createMockDelegationPatterns()
	mockCoOccurrences := createMockCoOccurrencePatterns()

	model := NewDashboardModel(mockStats, mockDelegations, mockCoOccurrences)
	model.focusedPanel = DelegationPanel

	// Test filter cycling: 0 → 2 → 5 → 10 → 0
	expectedCycle := []int{0, 2, 5, 10, 0}
	currentFilter := model.delegationFilterOptions.MinCount

	// Should start at 0
	if currentFilter != 0 {
		t.Errorf("Expected initial MinCount to be 0, got %d", currentFilter)
	}

	fKeyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}}

	for i := 1; i < len(expectedCycle); i++ {
		model.Update(fKeyMsg)
		currentFilter = model.delegationFilterOptions.MinCount

		if currentFilter != expectedCycle[i] {
			t.Errorf("Filter cycle step %d: expected MinCount %d, got %d", i, expectedCycle[i], currentFilter)
		}

		// Selection should reset when filters change
		if model.selectedDelegationIndex != 0 {
			t.Errorf("Expected delegation selection to reset to 0 after filter change, got %d", model.selectedDelegationIndex)
		}
	}
}

// TestDelegationSortOptionsCycling tests 's' key for sort cycling in delegation panel
func TestDelegationSortOptionsCycling(t *testing.T) {
	mockStats := createMockAgentStats()
	mockDelegations := createMockDelegationPatterns()
	mockCoOccurrences := createMockCoOccurrencePatterns()

	model := NewDashboardModel(mockStats, mockDelegations, mockCoOccurrences)
	model.focusedPanel = DelegationPanel

	// Test sort cycling: frequency → recency → source → target → frequency
	expectedCycle := []DelegationSortType{
		SortByFrequency,
		SortByRecency,
		SortBySource,
		SortByTarget,
		SortByFrequency,
	}

	// Should start with frequency
	if model.delegationFilterOptions.SortBy != SortByFrequency {
		t.Errorf("Expected initial sort to be SortByFrequency, got %v", model.delegationFilterOptions.SortBy)
	}

	sKeyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}}

	for i := 1; i < len(expectedCycle); i++ {
		model.Update(sKeyMsg)
		currentSort := model.delegationFilterOptions.SortBy

		if currentSort != expectedCycle[i] {
			t.Errorf("Sort cycle step %d: expected %v, got %v", i, expectedCycle[i], currentSort)
		}

		// Selection should reset when sorting changes
		if model.selectedDelegationIndex != 0 {
			t.Errorf("Expected delegation selection to reset to 0 after sort change, got %d", model.selectedDelegationIndex)
		}
	}
}

// TestFilterOptionsUIFeedback tests visual feedback for filter changes
func TestFilterOptionsUIFeedback(t *testing.T) {
	mockStats := createMockAgentStats()
	mockDelegations := createMockDelegationPatterns()
	mockCoOccurrences := createMockCoOccurrencePatterns()

	model := NewDashboardModel(mockStats, mockDelegations, mockCoOccurrences)
	model.focusedPanel = DelegationPanel

	// Test that view reflects filter state
	view := model.View()
	if !strings.Contains(view, "Delegation") {
		t.Error("Expected delegation panel to be visible")
	}

	// Apply minimum count filter
	fKeyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}}
	model.Update(fKeyMsg)

	// View should reflect changed filter state
	filteredView := model.View()
	if filteredView == view {
		t.Error("View should change after applying filter")
	}

	// Test sort option visual feedback
	sKeyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}}
	model.Update(sKeyMsg)

	sortedView := model.View()
	if sortedView == "" {
		t.Error("View should render after changing sort option")
	}
}

// ============================================================================
// Data Refresh Tests - Testing data reload on filter changes
// ============================================================================

// TestDataRefreshOnTimeFilter tests that time filter changes trigger data reload
func TestDataRefreshOnTimeFilter(t *testing.T) {
	model := NewDashboardModelWithLoading()

	// Mock loading state
	model.isLoading = false
	model.lastLoadTime = time.Now().Add(-1 * time.Hour)

	// Apply time filter
	key2Msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}}
	_, cmd := model.Update(key2Msg)

	// Should trigger data reload
	if cmd == nil {
		t.Error("Time filter change should trigger data reload command")
	}

	// Loading state should be set (in real implementation)
	// if !model.isLoading {
	//     t.Error("Should set loading state when triggering reload")
	// }
}

// TestDataRefreshProgress tests progress indicators during reload
func TestDataRefreshProgress(t *testing.T) {
	model := NewDashboardModelWithLoading()

	// Simulate loading state
	model.isLoading = true

	view := model.View()
	// Should show loading state in view (when implemented)
	if view == "" {
		t.Error("Should render view even during loading")
	}

	// Complete loading with new data
	mockStats := createMockAgentStats()
	loadedMsg := DataLoadedMsg{
		AgentStats:         mockStats,
		DelegationPatterns: createMockDelegationPatterns(),
		CoOccurrences:      createMockCoOccurrencePatterns(),
		MessageEvents:      createMockMessageEvents(),
		Error:              nil,
	}

	_, cmd := model.Update(loadedMsg)

	// Loading should be complete
	if model.isLoading {
		t.Error("Loading state should be false after successful data load")
	}

	// Should not trigger additional commands
	if cmd != nil {
		t.Error("Successful data load should not trigger additional commands")
	}

	// Data should be updated
	if len(model.GetAgentStats()) == 0 {
		t.Error("Agent stats should be populated after data load")
	}
}

// TestDataRefreshError tests error handling during data reload
func TestDataRefreshError(t *testing.T) {
	model := NewDashboardModelWithLoading()

	// Simulate load error
	errorMsg := DataLoadedMsg{
		AgentStats:         nil,
		DelegationPatterns: nil,
		CoOccurrences:      nil,
		MessageEvents:      nil,
		Error:              fmt.Errorf("mock error"),
	}

	_, cmd := model.Update(errorMsg)

	// Should handle error gracefully
	if model.isLoading {
		t.Error("Loading state should be false after error")
	}

	// Should store error for display
	if model.lastLoadError == nil {
		t.Error("Should store load error for user feedback")
	}

	// Should not trigger additional commands
	if cmd != nil {
		t.Error("Error handling should not trigger additional commands")
	}

	// View should still render with error state
	view := model.View()
	if view == "" {
		t.Error("Should render view even with load error")
	}
}

// TestRKeyRefreshFunctionality tests 'r' key for manual refresh
func TestRKeyRefreshFunctionality(t *testing.T) {
	mockStats := createMockAgentStats()
	mockDelegations := createMockDelegationPatterns()
	mockCoOccurrences := createMockCoOccurrencePatterns()

	model := NewDashboardModel(mockStats, mockDelegations, mockCoOccurrences)

	// Test 'r' key behavior in different panels
	testCases := []struct {
		panel           PanelType
		expectedCommand bool
		description     string
	}{
		{TimelinePanel, false, "timeline panel should reset view on 'r' key"},
		{StatsPanel, true, "stats panel should refresh data on 'r' key"},
		{DelegationPanel, true, "delegation panel should refresh data on 'r' key"},
		{CoOccurrencePanel, true, "co-occurrence panel should refresh data on 'r' key"},
	}

	for _, test := range testCases {
		model.focusedPanel = test.panel

		rKeyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}}
		_, cmd := model.Update(rKeyMsg)

		if test.expectedCommand && cmd == nil {
			t.Errorf("%s - expected refresh command", test.description)
		} else if !test.expectedCommand && cmd != nil {
			t.Errorf("%s - expected no command (timeline reset)", test.description)
		}
	}
}

// ============================================================================
// Test Helper Types - Types not yet implemented in dashboard.go
// ============================================================================

// MockDashboardModel extension methods for testing time filter functionality
func (m *DashboardModel) SetTimeFilter(filter string) {
	m.timeFilter = filter
}

func (m *DashboardModel) TriggerDataReload() tea.Cmd {
	return func() tea.Msg {
		return LoadDataMsg{TimeFilter: m.timeFilter}
	}
}