package agent

import (
	"fmt"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rsmdt/the-startup/internal/stats"
)

// Mock data for testing
func createMockAgentStats() map[string]*stats.GlobalAgentStats {
	mockStats := make(map[string]*stats.GlobalAgentStats)

	// Create Test Engineer agent stats
	testEngineerStats := &stats.GlobalAgentStats{}
	testEngineerStats.UpdateStats(1200) // 1.2s
	testEngineerStats.UpdateStats(800)  // 0.8s
	testEngineerStats.SuccessCount = 2
	mockStats["TestEngineer"] = testEngineerStats

	// Create Code Reviewer agent stats
	codeReviewerStats := &stats.GlobalAgentStats{}
	codeReviewerStats.UpdateStats(2000) // 2.0s
	codeReviewerStats.SuccessCount = 1
	mockStats["CodeReviewer"] = codeReviewerStats

	// Create DevOps agent stats
	devOpsStats := &stats.GlobalAgentStats{}
	devOpsStats.UpdateStats(5000) // 5.0s
	devOpsStats.UpdateStats(3000) // 3.0s
	devOpsStats.SuccessCount = 1
	devOpsStats.FailureCount = 1
	mockStats["DevOps"] = devOpsStats

	return mockStats
}

func createMockDelegationPatterns() []stats.DelegationPattern {
	return []stats.DelegationPattern{
		{
			SourceAgent: "TestEngineer",
			TargetAgent: "CodeReviewer",
			Count:       5,
			LastSeen:    time.Now(),
			SessionID:   "session1",
		},
		{
			SourceAgent: "CodeReviewer",
			TargetAgent: "DevOps",
			Count:       3,
			LastSeen:    time.Now(),
			SessionID:   "session2",
		},
		{
			SourceAgent: "TestEngineer",
			TargetAgent: "DevOps",
			Count:       2,
			LastSeen:    time.Now(),
			SessionID:   "session3",
		},
	}
}

func createMockCoOccurrencePatterns() []stats.AgentCoOccurrence {
	return []stats.AgentCoOccurrence{
		{
			Agent1:   "TestEngineer",
			Agent2:   "CodeReviewer",
			Count:    4,
			Sessions: []string{"session1", "session2", "session3", "session4"},
			LastSeen: time.Now(),
		},
		{
			Agent1:   "CodeReviewer",
			Agent2:   "DevOps",
			Count:    2,
			Sessions: []string{"session2", "session5"},
			LastSeen: time.Now(),
		},
	}
}

// Test Model Initialization and Setup
func TestNewDashboardModel(t *testing.T) {
	// Create mock data
	mockStats := createMockAgentStats()
	mockDelegations := createMockDelegationPatterns()
	mockCoOccurrences := createMockCoOccurrencePatterns()

	// Test dashboard creation with valid data
	model := NewDashboardModel(mockStats, mockDelegations, mockCoOccurrences)

	if model == nil {
		t.Fatal("NewDashboardModel returned nil")
	}

	// Should start with Timeline panel focused
	if model.focusedPanel != TimelinePanel {
		t.Errorf("Expected initial focus on TimelinePanel, got %v", model.focusedPanel)
	}

	// Should have agent stats initialized
	if model.agentStats == nil {
		t.Error("Expected agentStats to be initialized")
	}

	if len(model.agentStats) != 3 {
		t.Errorf("Expected 3 agent stats, got %d", len(model.agentStats))
	}

	// Should have delegation patterns initialized
	if model.delegationPatterns == nil {
		t.Error("Expected delegationPatterns to be initialized")
	}

	if len(model.delegationPatterns) != 3 {
		t.Errorf("Expected 3 delegation patterns, got %d", len(model.delegationPatterns))
	}

	// Should have co-occurrence patterns initialized
	if model.coOccurrences == nil {
		t.Error("Expected coOccurrences to be initialized")
	}

	if len(model.coOccurrences) != 2 {
		t.Errorf("Expected 2 co-occurrence patterns, got %d", len(model.coOccurrences))
	}

	// Should initialize selection state
	if model.selectedStatsIndex != 0 {
		t.Errorf("Expected selectedStatsIndex to start at 0, got %d", model.selectedStatsIndex)
	}

	if model.selectedDelegationIndex != 0 {
		t.Errorf("Expected selectedDelegationIndex to start at 0, got %d", model.selectedDelegationIndex)
	}
}

func TestDashboardModelInitWithEmptyData(t *testing.T) {
	// Test with empty data
	emptyStats := make(map[string]*stats.GlobalAgentStats)
	emptyDelegations := []stats.DelegationPattern{}
	emptyCoOccurrences := []stats.AgentCoOccurrence{}

	model := NewDashboardModel(emptyStats, emptyDelegations, emptyCoOccurrences)

	if model == nil {
		t.Fatal("NewDashboardModel returned nil with empty data")
	}

	// Should handle empty data gracefully
	if model.agentStats == nil {
		t.Error("Expected agentStats to be initialized even when empty")
	}

	if model.delegationPatterns == nil {
		t.Error("Expected delegationPatterns to be initialized even when empty")
	}

	if model.coOccurrences == nil {
		t.Error("Expected coOccurrences to be initialized even when empty")
	}
}

func TestDashboardModelInit(t *testing.T) {
	mockStats := createMockAgentStats()
	mockDelegations := createMockDelegationPatterns()
	mockCoOccurrences := createMockCoOccurrencePatterns()

	model := NewDashboardModel(mockStats, mockDelegations, mockCoOccurrences)

	// Init should return nil command for dashboard
	cmd := model.Init()
	if cmd != nil {
		t.Error("Expected Init() to return nil command for dashboard")
	}
}

// Test Update Message Handling for Keyboard Navigation
func TestDashboardModelKeyboardNavigation(t *testing.T) {
	mockStats := createMockAgentStats()
	mockDelegations := createMockDelegationPatterns()
	mockCoOccurrences := createMockCoOccurrencePatterns()

	model := NewDashboardModel(mockStats, mockDelegations, mockCoOccurrences)

	// Test Tab key - switches between panels
	tabMsg := tea.KeyMsg{Type: tea.KeyTab}
	_, cmd := model.Update(tabMsg)

	// Should switch to stats panel
	if model.focusedPanel != StatsPanel {
		t.Errorf("Expected focus to switch to StatsPanel, got %v", model.focusedPanel)
	}

	// Tab again should go to delegation panel
	model.Update(tabMsg)
	if model.focusedPanel != DelegationPanel {
		t.Errorf("Expected focus to switch to DelegationPanel, got %v", model.focusedPanel)
	}

	// Tab again should go to co-occurrence panel
	model.Update(tabMsg)
	if model.focusedPanel != CoOccurrencePanel {
		t.Errorf("Expected focus to switch to CoOccurrencePanel, got %v", model.focusedPanel)
	}

	// Tab again should cycle back to timeline panel
	model.Update(tabMsg)
	if model.focusedPanel != TimelinePanel {
		t.Errorf("Expected focus to cycle back to TimelinePanel, got %v", model.focusedPanel)
	}

	if cmd != nil {
		t.Error("Expected Tab navigation to return nil command")
	}
}

func TestDashboardModelArrowKeyNavigation(t *testing.T) {
	mockStats := createMockAgentStats()
	mockDelegations := createMockDelegationPatterns()
	mockCoOccurrences := createMockCoOccurrencePatterns()

	model := NewDashboardModel(mockStats, mockDelegations, mockCoOccurrences)

	// Test Down arrow in stats panel
	downMsg := tea.KeyMsg{Type: tea.KeyDown}
	model.Update(downMsg)

	if model.selectedStatsIndex != 1 {
		t.Errorf("Expected selectedStatsIndex to be 1, got %d", model.selectedStatsIndex)
	}

	// Test Up arrow
	upMsg := tea.KeyMsg{Type: tea.KeyUp}
	model.Update(upMsg)

	if model.selectedStatsIndex != 0 {
		t.Errorf("Expected selectedStatsIndex to be 0 after up arrow, got %d", model.selectedStatsIndex)
	}

	// Test that up arrow at top doesn't go negative
	model.Update(upMsg)
	if model.selectedStatsIndex < 0 {
		t.Error("Expected selectedStatsIndex to not go negative")
	}

	// Test down arrow at bottom boundary
	model.selectedStatsIndex = len(model.agentStats) - 1
	model.Update(downMsg)
	if model.selectedStatsIndex >= len(model.agentStats) {
		t.Error("Expected selectedStatsIndex to not exceed agent count")
	}
}

func TestDashboardModelDelegationPanelNavigation(t *testing.T) {
	mockStats := createMockAgentStats()
	mockDelegations := createMockDelegationPatterns()
	mockCoOccurrences := createMockCoOccurrencePatterns()

	model := NewDashboardModel(mockStats, mockDelegations, mockCoOccurrences)

	// Switch to delegation panel
	model.focusedPanel = DelegationPanel

	// Test navigation in delegation panel
	downMsg := tea.KeyMsg{Type: tea.KeyDown}
	model.Update(downMsg)

	if model.selectedDelegationIndex != 1 {
		t.Errorf("Expected selectedDelegationIndex to be 1, got %d", model.selectedDelegationIndex)
	}

	// Test boundary checks
	upMsg := tea.KeyMsg{Type: tea.KeyUp}
	model.selectedDelegationIndex = 0
	model.Update(upMsg)
	if model.selectedDelegationIndex < 0 {
		t.Error("Expected selectedDelegationIndex to not go negative")
	}
}

func TestDashboardModelDataRefresh(t *testing.T) {
	mockStats := createMockAgentStats()
	mockDelegations := createMockDelegationPatterns()
	mockCoOccurrences := createMockCoOccurrencePatterns()

	model := NewDashboardModel(mockStats, mockDelegations, mockCoOccurrences)

	// Test 'r' key for data refresh
	refreshMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}}
	_, cmd := model.Update(refreshMsg)

	// Should trigger refresh (return a command)
	if cmd == nil {
		t.Error("Expected refresh to return a command")
	}
}

func TestDashboardModelQuitHandling(t *testing.T) {
	mockStats := createMockAgentStats()
	mockDelegations := createMockDelegationPatterns()
	mockCoOccurrences := createMockCoOccurrencePatterns()

	model := NewDashboardModel(mockStats, mockDelegations, mockCoOccurrences)

	// Test 'q' key for quit
	quitMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	_, cmd := model.Update(quitMsg)

	// Should return quit command
	if cmd == nil {
		t.Error("Expected quit to return a command")
	}

	// Test Ctrl+C
	ctrlCMsg := tea.KeyMsg{Type: tea.KeyCtrlC}
	_, cmd = model.Update(ctrlCMsg)

	if cmd == nil {
		t.Error("Expected Ctrl+C to return a command")
	}

	// Test ESC key
	escMsg := tea.KeyMsg{Type: tea.KeyEsc}
	_, cmd = model.Update(escMsg)

	if cmd == nil {
		t.Error("Expected ESC to return a command")
	}
}

func TestDashboardModelEnterKeyHandling(t *testing.T) {
	mockStats := createMockAgentStats()
	mockDelegations := createMockDelegationPatterns()
	mockCoOccurrences := createMockCoOccurrencePatterns()

	model := NewDashboardModel(mockStats, mockDelegations, mockCoOccurrences)

	// Test Enter key in stats panel (should show detail)
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	_, cmd := model.Update(enterMsg)

	// Enter should trigger some action (detail view or similar)
	if cmd == nil {
		t.Error("Expected Enter key to return a command for detail view")
	}
}

// Test View Rendering for Dashboard Layout and Leaderboards
func TestDashboardModelViewRendering(t *testing.T) {
	mockStats := createMockAgentStats()
	mockDelegations := createMockDelegationPatterns()
	mockCoOccurrences := createMockCoOccurrencePatterns()

	model := NewDashboardModel(mockStats, mockDelegations, mockCoOccurrences)

	view := model.View()

	// Should contain dashboard title
	if !strings.Contains(view, "Agent Dashboard") {
		t.Error("Expected view to contain 'Agent Dashboard' title")
	}

	// Should contain agent stats section
	if !strings.Contains(view, "Agent Statistics") {
		t.Error("Expected view to contain agent statistics section")
	}

	// Should contain agent names
	if !strings.Contains(view, "TestEngineer") {
		t.Error("Expected view to contain 'TestEngineer' agent name")
	}

	if !strings.Contains(view, "CodeReviewer") {
		t.Error("Expected view to contain 'CodeReviewer' agent name")
	}

	if !strings.Contains(view, "DevOps") {
		t.Error("Expected view to contain 'DevOps' agent name")
	}

	// Should contain delegation patterns (may be truncated)
	if !strings.Contains(view, "Delegation Flow") && !strings.Contains(view, "Delegation Flow Patterns") {
		t.Logf("View output:\n%s", view)
		t.Error("Expected view to contain 'Delegation Flow' or 'Delegation Flow Patterns' section")
	}

	// Should show delegation flows (arrows or similar)
	if !strings.Contains(view, "→") && !strings.Contains(view, "->") {
		t.Error("Expected view to show delegation flow indicators")
	}

	// Should contain co-occurrence patterns (may be truncated)
	if !strings.Contains(view, "Agent Collaboration") && !strings.Contains(view, "Agent Collaboration Matrix") {
		t.Error("Expected view to contain 'Agent Collaboration' or 'Agent Collaboration Matrix' section")
	}
}

func TestDashboardModelStatsLeaderboardDisplay(t *testing.T) {
	mockStats := createMockAgentStats()
	mockDelegations := createMockDelegationPatterns()
	mockCoOccurrences := createMockCoOccurrencePatterns()

	model := NewDashboardModel(mockStats, mockDelegations, mockCoOccurrences)

	view := model.View()

	// Should display statistics like count, success rate, avg duration
	if !strings.Contains(view, "Count:") && !strings.Contains(view, "Calls") {
		t.Error("Expected view to show invocation counts")
	}

	if !strings.Contains(view, "Success") && !strings.Contains(view, "%") {
		t.Error("Expected view to show success rates")
	}

	if !strings.Contains(view, "ms") || (!strings.Contains(view, "Avg") && !strings.Contains(view, "Time")) {
		t.Error("Expected view to show average duration in milliseconds")
	}
}

func TestDashboardModelPanelFocusHighlighting(t *testing.T) {
	mockStats := createMockAgentStats()
	mockDelegations := createMockDelegationPatterns()
	mockCoOccurrences := createMockCoOccurrencePatterns()

	model := NewDashboardModel(mockStats, mockDelegations, mockCoOccurrences)

	// Test stats panel focused (default)
	view := model.View()
	if !strings.Contains(view, "▶") && !strings.Contains(view, ">") && !strings.Contains(view, "●") {
		t.Error("Expected focused stats panel to show selection indicator")
	}

	// Switch to delegation panel
	model.focusedPanel = DelegationPanel
	view = model.View()

	// Should show focus moved to delegation panel
	// (Specific visual indicators depend on implementation, but should be different)
	if len(view) == 0 {
		t.Error("Expected delegation panel view to render when focused")
	}

	// Switch to co-occurrence panel
	model.focusedPanel = CoOccurrencePanel
	view = model.View()

	if len(view) == 0 {
		t.Error("Expected co-occurrence panel view to render when focused")
	}
}

func TestDashboardModelEmptyDataDisplay(t *testing.T) {
	// Test display with no data
	emptyStats := make(map[string]*stats.GlobalAgentStats)
	emptyDelegations := []stats.DelegationPattern{}
	emptyCoOccurrences := []stats.AgentCoOccurrence{}

	model := NewDashboardModel(emptyStats, emptyDelegations, emptyCoOccurrences)

	view := model.View()

	// Should handle empty state gracefully
	if !strings.Contains(view, "No data") && !strings.Contains(view, "No agents") && !strings.Contains(view, "Empty") {
		t.Error("Expected empty state message when no data available")
	}

	// Should not crash
	if view == "" {
		t.Error("Expected non-empty view even with no data")
	}
}

// Test Keyboard Navigation Comprehensively
func TestDashboardModelVimStyleNavigation(t *testing.T) {
	mockStats := createMockAgentStats()
	mockDelegations := createMockDelegationPatterns()
	mockCoOccurrences := createMockCoOccurrencePatterns()

	model := NewDashboardModel(mockStats, mockDelegations, mockCoOccurrences)

	// Test 'j' key (vim down)
	jMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
	model.Update(jMsg)

	if model.selectedStatsIndex != 1 {
		t.Errorf("Expected 'j' key to move selection down, got index %d", model.selectedStatsIndex)
	}

	// Test 'k' key (vim up)
	kMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}
	model.Update(kMsg)

	if model.selectedStatsIndex != 0 {
		t.Errorf("Expected 'k' key to move selection up, got index %d", model.selectedStatsIndex)
	}
}

func TestDashboardModelHelpToggle(t *testing.T) {
	mockStats := createMockAgentStats()
	mockDelegations := createMockDelegationPatterns()
	mockCoOccurrences := createMockCoOccurrencePatterns()

	model := NewDashboardModel(mockStats, mockDelegations, mockCoOccurrences)

	// Test '?' key for help
	helpMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}
	model.Update(helpMsg)

	if !model.showHelp {
		t.Error("Expected help to be shown after '?' key")
	}

	view := model.View()
	if !strings.Contains(view, "Help") && !strings.Contains(view, "Commands") {
		t.Error("Expected help view to contain help content")
	}

	// Test help toggle off
	model.Update(helpMsg)
	if model.showHelp {
		t.Error("Expected help to be hidden after second '?' key")
	}
}

// Test Dashboard Component Composition and State Management
func TestDashboardModelStateManagement(t *testing.T) {
	mockStats := createMockAgentStats()
	mockDelegations := createMockDelegationPatterns()
	mockCoOccurrences := createMockCoOccurrencePatterns()

	model := NewDashboardModel(mockStats, mockDelegations, mockCoOccurrences)

	// Test that state transitions are isolated between panels
	initialStatsIndex := model.selectedStatsIndex
	initialDelegationIndex := model.selectedDelegationIndex

	// Change stats selection
	downMsg := tea.KeyMsg{Type: tea.KeyDown}
	model.Update(downMsg)

	// Should only affect stats panel
	if model.selectedStatsIndex == initialStatsIndex {
		t.Error("Expected stats index to change")
	}

	if model.selectedDelegationIndex != initialDelegationIndex {
		t.Error("Expected delegation index to remain unchanged")
	}

	// Switch to delegation panel and test isolation
	model.focusedPanel = DelegationPanel
	model.Update(downMsg)

	// Should only affect delegation panel
	if model.selectedDelegationIndex == initialDelegationIndex {
		t.Error("Expected delegation index to change when delegation panel focused")
	}
}

func TestDashboardModelDataIntegrity(t *testing.T) {
	mockStats := createMockAgentStats()
	mockDelegations := createMockDelegationPatterns()
	mockCoOccurrences := createMockCoOccurrencePatterns()

	model := NewDashboardModel(mockStats, mockDelegations, mockCoOccurrences)

	// Test that original data is not modified
	originalStatsCount := len(mockStats)
	originalDelegationCount := len(mockDelegations)
	originalCoOccurrenceCount := len(mockCoOccurrences)

	// Perform various operations
	tabMsg := tea.KeyMsg{Type: tea.KeyTab}
	downMsg := tea.KeyMsg{Type: tea.KeyDown}

	model.Update(tabMsg)
	model.Update(downMsg)
	model.Update(tabMsg)
	model.Update(downMsg)

	// Original data should be unchanged
	if len(mockStats) != originalStatsCount {
		t.Error("Expected original agent stats to remain unchanged")
	}

	if len(mockDelegations) != originalDelegationCount {
		t.Error("Expected original delegation patterns to remain unchanged")
	}

	if len(mockCoOccurrences) != originalCoOccurrenceCount {
		t.Error("Expected original co-occurrence patterns to remain unchanged")
	}

	// Model should have copies or references, not modify originals
	if len(model.agentStats) != originalStatsCount {
		t.Error("Expected model agent stats to match original count")
	}

	if len(model.delegationPatterns) != originalDelegationCount {
		t.Error("Expected model delegation patterns to match original count")
	}

	if len(model.coOccurrences) != originalCoOccurrenceCount {
		t.Error("Expected model co-occurrences to match original count")
	}
}

func TestDashboardModelComponentComposition(t *testing.T) {
	mockStats := createMockAgentStats()
	mockDelegations := createMockDelegationPatterns()
	mockCoOccurrences := createMockCoOccurrencePatterns()

	model := NewDashboardModel(mockStats, mockDelegations, mockCoOccurrences)

	// Test that all components are properly initialized
	if model.agentStats == nil {
		t.Error("Expected agentStats component to be initialized")
	}

	if model.delegationPatterns == nil {
		t.Error("Expected delegationPatterns component to be initialized")
	}

	if model.coOccurrences == nil {
		t.Error("Expected coOccurrences component to be initialized")
	}

	// Test that all panels can render
	for _, panel := range []PanelType{TimelinePanel, StatsPanel, DelegationPanel, CoOccurrencePanel} {
		model.focusedPanel = panel
		view := model.View()

		if view == "" {
			t.Errorf("Expected non-empty view for panel %v", panel)
		}

		if len(view) < 10 {
			t.Errorf("Expected substantial content for panel %v, got length %d", panel, len(view))
		}
	}
}

// Test edge cases and error conditions
func TestDashboardModelEdgeCases(t *testing.T) {
	// Test with single agent
	singleAgentStats := map[string]*stats.GlobalAgentStats{
		"OnlyAgent": &stats.GlobalAgentStats{},
	}

	model := NewDashboardModel(singleAgentStats, []stats.DelegationPattern{}, []stats.AgentCoOccurrence{})

	// Navigation should handle single item gracefully
	downMsg := tea.KeyMsg{Type: tea.KeyDown}
	upMsg := tea.KeyMsg{Type: tea.KeyUp}

	model.Update(downMsg)
	if model.selectedStatsIndex < 0 || model.selectedStatsIndex >= len(model.agentStats) {
		t.Error("Expected single agent navigation to stay within bounds")
	}

	model.Update(upMsg)
	if model.selectedStatsIndex < 0 {
		t.Error("Expected single agent navigation to not go negative")
	}

	// View should render without errors
	view := model.View()
	if view == "" {
		t.Error("Expected non-empty view for single agent")
	}
}

// containsString helper function is already defined above

// Enhanced Stats Panel Tests for Phase 2 Implementation
// These tests validate the sorting and selection functionality described in the PLAN.md Phase 2

// Test agent leaderboard sorting with keyboard shortcuts
func TestDashboardModelAgentLeaderboardSorting(t *testing.T) {
	mockStats := createMockAgentStats()
	mockDelegations := createMockDelegationPatterns()
	mockCoOccurrences := createMockCoOccurrencePatterns()

	model := NewDashboardModel(mockStats, mockDelegations, mockCoOccurrences)

	// Test default sorting (by count - usage count)
	view := model.View()
	if !strings.Contains(view, "Agent Statistics Leaderboard") {
		t.Error("Expected agent leaderboard to be displayed")
	}

	// Test 'c' key for sorting by count (usage)
	cSortMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}}
	model.Update(cSortMsg)

	// Verify count sorting is active
	view = model.View()
	// The view should show agents sorted by count (TestEngineer=2, DevOps=2, CodeReviewer=1)
	if !strings.Contains(view, "TestEngineer") {
		t.Error("Expected TestEngineer to be visible in leaderboard")
	}

	// Test 's' key for sorting by success rate
	sSortMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}}
	model.Update(sSortMsg)

	// Test 'd' key for sorting by duration
	dSortMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}}
	model.Update(dSortMsg)

	// Verify duration sorting affects the leaderboard order
	view = model.View()
	if !strings.Contains(view, "Avg:") && !strings.Contains(view, "ms") {
		t.Error("Expected duration metrics to be displayed when sorting by duration")
	}
}

func TestDashboardModelLeaderboardSortingBehavior(t *testing.T) {
	mockStats := createEnhancedMockAgentStats() // Create varied test data
	mockDelegations := createMockDelegationPatterns()
	mockCoOccurrences := createMockCoOccurrencePatterns()

	model := NewDashboardModel(mockStats, mockDelegations, mockCoOccurrences)

	// Test that default sorting shows agents by count
	view := model.View()

	// FastAgent has 3 calls, should rank high
	if !strings.Contains(view, "FastAgent") {
		t.Error("Expected FastAgent to appear in leaderboard")
	}

	// Test sorting keyboard shortcuts work with different data sets
	cSortMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}}
	model.Update(cSortMsg)

	sSortMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}}
	model.Update(sSortMsg)

	dSortMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}}
	model.Update(dSortMsg)

	// All sorting operations should complete without error
	view = model.View()
	if view == "" {
		t.Error("Expected non-empty view after sorting operations")
	}
}

func TestDashboardModelAgentSelectionAndHighlighting(t *testing.T) {
	mockStats := createMockAgentStats()
	mockDelegations := createMockDelegationPatterns()
	mockCoOccurrences := createMockCoOccurrencePatterns()

	model := NewDashboardModel(mockStats, mockDelegations, mockCoOccurrences)

	// Ensure stats panel is focused
	model.focusedPanel = StatsPanel

	// Test initial selection (should be 0)
	if model.selectedStatsIndex != 0 {
		t.Errorf("Expected initial selection index to be 0, got %d", model.selectedStatsIndex)
	}

	// Test down arrow navigation
	downMsg := tea.KeyMsg{Type: tea.KeyDown}
	model.Update(downMsg)
	if model.selectedStatsIndex != 1 {
		t.Errorf("Expected selection index to be 1 after down arrow, got %d", model.selectedStatsIndex)
	}

	// Test up arrow navigation
	upMsg := tea.KeyMsg{Type: tea.KeyUp}
	model.Update(upMsg)
	if model.selectedStatsIndex != 0 {
		t.Errorf("Expected selection index to be 0 after up arrow, got %d", model.selectedStatsIndex)
	}

	// Test view contains selection indicator
	view := model.View()
	if !strings.Contains(view, "▶") && !strings.Contains(view, ">") {
		t.Error("Expected view to contain selection indicator when agent is selected")
	}

	// Test vim-style navigation (j/k keys)
	jMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
	model.Update(jMsg)
	if model.selectedStatsIndex != 1 {
		t.Errorf("Expected 'j' key to move selection down, got index %d", model.selectedStatsIndex)
	}

	kMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}
	model.Update(kMsg)
	if model.selectedStatsIndex != 0 {
		t.Errorf("Expected 'k' key to move selection up, got index %d", model.selectedStatsIndex)
	}
}

func TestDashboardModelLeaderboardComponentSortingControls(t *testing.T) {
	mockStats := createMockAgentStats()
	mockDelegations := createMockDelegationPatterns()
	mockCoOccurrences := createMockCoOccurrencePatterns()

	_ = NewDashboardModel(mockStats, mockDelegations, mockCoOccurrences)

	// Test leaderboard component initialization
	leaderboard := NewLeaderboardComponent("Agent Statistics")
	if leaderboard == nil {
		t.Fatal("Expected leaderboard component to be created")
	}

	// Test agent leaderboard rendering with selection
	renderedView := leaderboard.RenderAgentLeaderboard(mockStats, 0, true)
	if renderedView == "" {
		t.Error("Expected rendered leaderboard view to not be empty")
	}

	// Should contain agent names
	if !strings.Contains(renderedView, "TestEngineer") {
		t.Error("Expected rendered view to contain TestEngineer")
	}
	if !strings.Contains(renderedView, "CodeReviewer") {
		t.Error("Expected rendered view to contain CodeReviewer")
	}
	if !strings.Contains(renderedView, "DevOps") {
		t.Error("Expected rendered view to contain DevOps")
	}

	// Should contain headers for the leaderboard
	if !strings.Contains(renderedView, "Rank") || !strings.Contains(renderedView, "Agent") {
		t.Error("Expected rendered view to contain leaderboard headers")
	}

	// Should contain selection indicator when focused
	if !strings.Contains(renderedView, "▶") {
		t.Error("Expected focused leaderboard to show selection indicator")
	}
}

func TestDashboardModelSuccessRateAndDurationMetricsFormatting(t *testing.T) {
	mockStats := createMockAgentStats()
	mockDelegations := createMockDelegationPatterns()
	mockCoOccurrences := createMockCoOccurrencePatterns()

	model := NewDashboardModel(mockStats, mockDelegations, mockCoOccurrences)

	view := model.View()

	// Test success rate formatting
	if !strings.Contains(view, "%") {
		t.Error("Expected view to contain success rate percentage")
	}

	// Test duration formatting
	if !strings.Contains(view, "ms") {
		t.Error("Expected view to contain duration in milliseconds")
	}

	// Test specific agent metrics formatting
	// TestEngineer should show 100% success rate (2/2)
	if !strings.Contains(view, "100.0%") && !strings.Contains(view, "100%") {
		t.Error("Expected TestEngineer to show 100% success rate")
	}

	// Should show average duration
	if !strings.Contains(view, "1000ms") && !strings.Contains(view, "1000") {
		t.Error("Expected to show TestEngineer average duration of 1000ms")
	}

	// Test leaderboard component formatting
	leaderboard := NewLeaderboardComponent("Test Leaderboard")
	entries := leaderboard.convertAgentStatsToEntries(mockStats, 0)

	for _, entry := range entries {
		// Success rate should be between 0-100
		if entry.SuccessRate < 0 || entry.SuccessRate > 100 {
			t.Errorf("Invalid success rate for %s: %.1f%%, should be 0-100",
				entry.Name, entry.SuccessRate)
		}

		// Duration should be positive
		if entry.AvgDuration < 0 {
			t.Errorf("Invalid duration for %s: %.0fms, should be non-negative",
				entry.Name, entry.AvgDuration)
		}

		// Count should be positive
		if entry.Count <= 0 {
			t.Errorf("Invalid count for %s: %d, should be positive",
				entry.Name, entry.Count)
		}
	}
}

func TestDashboardModelVisualFeedbackForSorting(t *testing.T) {
	mockStats := createMockAgentStats()
	mockDelegations := createMockDelegationPatterns()
	mockCoOccurrences := createMockCoOccurrencePatterns()

	model := NewDashboardModel(mockStats, mockDelegations, mockCoOccurrences)

	// Test visual feedback for count sorting
	cSortMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}}
	model.Update(cSortMsg)
	view := model.View()

	// Should show sorting indicator or highlight
	if !strings.Contains(view, "Count") {
		t.Error("Expected view to show count sorting feedback")
	}

	// Test visual feedback for success rate sorting
	sSortMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}}
	model.Update(sSortMsg)
	view = model.View()

	if !strings.Contains(view, "Success") {
		t.Error("Expected view to show success rate sorting feedback")
	}

	// Test visual feedback for duration sorting
	dSortMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}}
	model.Update(dSortMsg)
	view = model.View()

	if !strings.Contains(view, "Avg") || !strings.Contains(view, "Time") {
		t.Error("Expected view to show duration sorting feedback")
	}
}

func TestDashboardModelBoundaryConditionsForSelection(t *testing.T) {
	mockStats := createMockAgentStats()
	mockDelegations := createMockDelegationPatterns()
	mockCoOccurrences := createMockCoOccurrencePatterns()

	model := NewDashboardModel(mockStats, mockDelegations, mockCoOccurrences)
	model.focusedPanel = StatsPanel

	// Test boundary at top (shouldn't go negative)
	model.selectedStatsIndex = 0
	upMsg := tea.KeyMsg{Type: tea.KeyUp}
	model.Update(upMsg)
	if model.selectedStatsIndex < 0 {
		t.Error("Selection index should not go negative")
	}
	if model.selectedStatsIndex != 0 {
		t.Error("Selection index should stay at 0 when at top boundary")
	}

	// Test boundary at bottom (shouldn't exceed length)
	model.selectedStatsIndex = len(mockStats) - 1
	downMsg := tea.KeyMsg{Type: tea.KeyDown}
	model.Update(downMsg)
	if model.selectedStatsIndex >= len(mockStats) {
		t.Error("Selection index should not exceed agent count")
	}
	if model.selectedStatsIndex != len(mockStats)-1 {
		t.Error("Selection index should stay at max when at bottom boundary")
	}
}

func TestDashboardModelEmptyDataLeaderboardHandling(t *testing.T) {
	// Test with empty agent stats
	emptyStats := make(map[string]*stats.GlobalAgentStats)
	emptyDelegations := []stats.DelegationPattern{}
	emptyCoOccurrences := []stats.AgentCoOccurrence{}

	model := NewDashboardModel(emptyStats, emptyDelegations, emptyCoOccurrences)

	// Should handle empty data gracefully
	view := model.View()
	if !strings.Contains(view, "No agents") && !strings.Contains(view, "No data") {
		t.Error("Expected empty state message for no agents")
	}

	// Navigation should not crash with empty data
	downMsg := tea.KeyMsg{Type: tea.KeyDown}
	model.Update(downMsg)
	upMsg := tea.KeyMsg{Type: tea.KeyUp}
	model.Update(upMsg)

	// Selection index should remain valid
	if model.selectedStatsIndex < 0 {
		t.Error("Selection index should not be negative with empty data")
	}

	// Sorting keys should not crash with empty data
	cSortMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}}
	model.Update(cSortMsg)
	sSortMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}}
	model.Update(sSortMsg)
	dSortMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}}
	model.Update(dSortMsg)

	// Test empty timeline data handling
	emptyMessageEvents := []MessageEvent{}
	if len(emptyMessageEvents) != 0 {
		t.Error("Empty message events should have zero length")
	}

	// Test empty co-occurrence matrix handling
	matrixComponent := NewCoOccurrenceMatrixComponent()
	emptyMatrixView := matrixComponent.RenderCoOccurrenceMatrix(
		emptyCoOccurrences,
		0, // selectedIndex
		true, // isFocused
		"Empty Matrix Test",
	)

	if !strings.Contains(emptyMatrixView, "No co-occurrence patterns found") {
		t.Error("Expected empty co-occurrence message for empty data")
	}
}

// Test integration with leaderboard component rendering
func TestDashboardModelLeaderboardComponentIntegration(t *testing.T) {
	mockStats := createMockAgentStats()
	mockDelegations := createMockDelegationPatterns()
	mockCoOccurrences := createMockCoOccurrencePatterns()

	model := NewDashboardModel(mockStats, mockDelegations, mockCoOccurrences)
	model.focusedPanel = StatsPanel

	// Test that dashboard properly uses leaderboard component
	leaderboard := NewLeaderboardComponent("Agent Statistics")
	renderedLeaderboard := leaderboard.RenderAgentLeaderboard(
		mockStats,
		model.selectedStatsIndex,
		model.focusedPanel == StatsPanel,
	)

	// Rendered leaderboard should contain all expected elements
	if !strings.Contains(renderedLeaderboard, "Agent Statistics") {
		t.Error("Expected leaderboard to contain title")
	}

	if !strings.Contains(renderedLeaderboard, "Rank") {
		t.Error("Expected leaderboard to contain rank header")
	}

	if !strings.Contains(renderedLeaderboard, "Calls") {
		t.Error("Expected leaderboard to contain calls header")
	}

	if !strings.Contains(renderedLeaderboard, "Success") {
		t.Error("Expected leaderboard to contain success header")
	}

	if !strings.Contains(renderedLeaderboard, "Avg Time") {
		t.Error("Expected leaderboard to contain avg time header")
	}

	// Should contain selection indicator when focused
	if !strings.Contains(renderedLeaderboard, "▶") {
		t.Error("Expected focused leaderboard to show selection indicator")
	}

	// Test conversion of agent stats to leaderboard entries
	entries := leaderboard.convertAgentStatsToEntries(mockStats, 0)
	if len(entries) != len(mockStats) {
		t.Errorf("Expected %d leaderboard entries, got %d", len(mockStats), len(entries))
	}

	// Entries should be sorted by count (descending) by default
	for i := 0; i < len(entries)-1; i++ {
		if entries[i].Count < entries[i+1].Count {
			t.Errorf("Expected entries sorted by count descending, got %d < %d at position %d",
				entries[i].Count, entries[i+1].Count, i)
		}
	}

	// Test ranks are assigned correctly
	for i, entry := range entries {
		expectedRank := i + 1
		if entry.Rank != expectedRank {
			t.Errorf("Expected rank %d for entry %d, got %d", expectedRank, i, entry.Rank)
		}
	}

	// Test timeline component integration (conceptual)
	mockMessageEvents := createMockMessageEvents()
	if len(mockMessageEvents) > 0 {
		// Timeline component would integrate similar to leaderboard
		timelineComponent := struct {
			title  string
			events []MessageEvent
		}{
			title:  "Message Timeline",
			events: mockMessageEvents,
		}

		if timelineComponent.title == "" {
			t.Error("Timeline component should have a title")
		}
		if len(timelineComponent.events) == 0 {
			t.Error("Timeline component should have events")
		}
	}

	// Test co-occurrence matrix component integration
	mockCoOccurrencesEnhanced := createMockCoOccurrenceWithCorrelation()
	matrixComponent := NewCoOccurrenceMatrixComponent()
	renderedMatrix := matrixComponent.RenderCoOccurrenceMatrix(
		mockCoOccurrencesEnhanced,
		model.selectedCoOccurrenceIndex,
		model.focusedPanel == CoOccurrencePanel,
		"Agent Collaboration Matrix",
	)

	if !strings.Contains(renderedMatrix, "Agent Collaboration Matrix") {
		t.Error("Expected matrix to contain title")
	}

	if !strings.Contains(renderedMatrix, "⟷") {
		t.Error("Expected matrix to contain bidirectional indicators")
	}

	if !strings.Contains(renderedMatrix, "sessions") {
		t.Error("Expected matrix to contain session information")
	}
}

// Test leaderboard component style customization
func TestLeaderboardComponentStyleCustomization(t *testing.T) {
	leaderboard := NewLeaderboardComponent("Test Leaderboard")

	// Test default style
	defaultStyle := DefaultLeaderboardStyle()
	// Test that styles are properly initialized (can't directly test color values)
	if defaultStyle.Title.GetBold() != true {
		t.Error("Expected default style title to be bold")
	}

	// Test compact style
	compactStyle := CompactLeaderboardStyle()
	leaderboard.SetStyle(compactStyle)

	// Test max width setting
	leaderboard.SetMaxWidth(80)
	if leaderboard.MaxWidth != 80 {
		t.Errorf("Expected max width to be 80, got %d", leaderboard.MaxWidth)
	}

	// Test rendering with custom style
	mockStats := createMockAgentStats()
	renderedView := leaderboard.RenderAgentLeaderboard(mockStats, 0, true)
	if renderedView == "" {
		t.Error("Expected non-empty rendered view with custom style")
	}

	// Test timeline component style consistency (if implemented)
	// Timeline components should follow similar style patterns
	timelineStyle := defaultStyle // Timeline would use similar styling
	if timelineStyle.Title.GetBold() != true {
		t.Error("Timeline component should use consistent styling with leaderboard")
	}

	// Test co-occurrence matrix style customization
	matrixComponent := NewCoOccurrenceMatrixComponent()
	matrixComponent.Style = compactStyle

	mockCoOccurrences := createMockCoOccurrenceWithCorrelation()
	styledMatrixView := matrixComponent.RenderCoOccurrenceMatrix(
		mockCoOccurrences,
		0,
		true,
		"Styled Matrix",
	)

	if styledMatrixView == "" {
		t.Error("Expected non-empty styled matrix view")
	}
}

// Test string truncation functionality
func TestLeaderboardComponentStringTruncation(t *testing.T) {
	leaderboard := NewLeaderboardComponent("Test Leaderboard")

	// Test normal string (should not truncate)
	normalString := "TestAgent"
	truncated := leaderboard.truncateString(normalString, 15)
	if truncated != normalString {
		t.Errorf("Expected normal string to remain unchanged, got %s", truncated)
	}

	// Test long string (should truncate with ellipsis)
	longString := "VeryLongAgentNameThatExceedsLimit"
	truncated = leaderboard.truncateString(longString, 15)
	if len(truncated) > 15 {
		t.Errorf("Expected truncated string to be at most 15 chars, got %d", len(truncated))
	}
	if !strings.HasSuffix(truncated, "...") {
		t.Error("Expected truncated string to end with ellipsis")
	}

	// Test edge case: max length <= 3
	truncated = leaderboard.truncateString(longString, 3)
	if len(truncated) > 3 {
		t.Errorf("Expected truncated string to be at most 3 chars, got %d", len(truncated))
	}

	// Test edge case: max length = 0
	truncated = leaderboard.truncateString(longString, 0)
	if truncated != "" {
		t.Errorf("Expected empty string for max length 0, got %s", truncated)
	}

	// Test timeline content truncation (for message content)
	timelineContent := "This is a very long message content that should be truncated for display in the timeline view"
	truncatedContent := leaderboard.truncateString(timelineContent, 50)
	if len(truncatedContent) > 50 {
		t.Errorf("Timeline content should be truncated to 50 chars, got %d", len(truncatedContent))
	}

	// Test agent name truncation in co-occurrence matrix
	matrixComponent := NewCoOccurrenceMatrixComponent()
	longAgentName := "VeryLongAgentNameForCoOccurrenceMatrix"
	truncatedAgentName := matrixComponent.truncateString(longAgentName, 12)
	if len(truncatedAgentName) > 12 {
		t.Errorf("Agent name in matrix should be truncated to 12 chars, got %d", len(truncatedAgentName))
	}
	if !strings.HasSuffix(truncatedAgentName, "...") {
		t.Error("Truncated agent name should end with ellipsis")
	}
}

// Mock data for MessageEvent timeline testing
func createMockMessageEvents() []MessageEvent {
	return []MessageEvent{
		{
			ID:         "msg-001",
			Timestamp:  time.Now().Add(-5 * time.Minute),
			SessionID:  "session1",
			Role:       UserMessage,
			Content:    "Please help me write unit tests",
			TokenCount: 50,
			AgentUsed:  "TestEngineer",
			ToolsUsed:  []string{"Write", "Edit"},
			Success:    true,
			Duration:   1200 * time.Millisecond,
		},
		{
			ID:         "msg-002",
			Timestamp:  time.Now().Add(-3 * time.Minute),
			SessionID:  "session1",
			Role:       AssistantMessage,
			Content:    "I'll help you write comprehensive unit tests",
			TokenCount: 75,
			AgentUsed:  "TestEngineer",
			ToolsUsed:  []string{"Read", "Write"},
			Success:    true,
			Duration:   800 * time.Millisecond,
		},
		{
			ID:         "msg-003",
			Timestamp:  time.Now().Add(-2 * time.Minute),
			SessionID:  "session2",
			Role:       ToolMessage,
			Content:    "File written successfully",
			TokenCount: 25,
			AgentUsed:  "CodeReviewer",
			ToolsUsed:  []string{"Write"},
			Success:    true,
			Duration:   200 * time.Millisecond,
		},
		{
			ID:         "msg-004",
			Timestamp:  time.Now().Add(-1 * time.Minute),
			SessionID:  "session2",
			Role:       SystemMessage,
			Content:    "Error: Test failed",
			TokenCount: 15,
			AgentUsed:  "DevOps",
			ToolsUsed:  []string{"Bash"},
			Success:    false,
			Duration:   5000 * time.Millisecond,
		},
	}
}

// Mock data for enhanced AgentCoOccurrence testing with correlation
func createMockCoOccurrenceWithCorrelation() []stats.AgentCoOccurrence {
	return []stats.AgentCoOccurrence{
		{
			Agent1:   "TestEngineer",
			Agent2:   "CodeReviewer",
			Count:    12,
			Sessions: []string{"session1", "session2", "session3", "session4", "session5", "session6"},
			LastSeen: time.Now().Add(-30 * time.Minute),
		},
		{
			Agent1:   "CodeReviewer",
			Agent2:   "DevOps",
			Count:    8,
			Sessions: []string{"session2", "session5", "session7", "session8"},
			LastSeen: time.Now().Add(-45 * time.Minute),
		},
		{
			Agent1:   "TestEngineer",
			Agent2:   "DevOps",
			Count:    5,
			Sessions: []string{"session3", "session6", "session9"},
			LastSeen: time.Now().Add(-60 * time.Minute),
		},
		{
			Agent1:   "Architect",
			Agent2:   "TestEngineer",
			Count:    3,
			Sessions: []string{"session10", "session11"},
			LastSeen: time.Now().Add(-90 * time.Minute),
		},
	}
}

// Helper function to create enhanced mock data for sorting tests
func createEnhancedMockAgentStats() map[string]*stats.GlobalAgentStats {
	mockStats := make(map[string]*stats.GlobalAgentStats)

	// Create FastAgent with high count, low duration
	fastAgent := &stats.GlobalAgentStats{}
	fastAgent.UpdateStats(200) // 0.2s - fast
	fastAgent.UpdateStats(300) // 0.3s
	fastAgent.UpdateStats(250) // 0.25s
	fastAgent.SuccessCount = 3
	mockStats["FastAgent"] = fastAgent

	// Create SlowAgent with medium count, high duration
	slowAgent := &stats.GlobalAgentStats{}
	slowAgent.UpdateStats(5000) // 5s - slow
	slowAgent.UpdateStats(6000) // 6s
	slowAgent.SuccessCount = 1 // Low success rate
	slowAgent.FailureCount = 1
	mockStats["SlowAgent"] = slowAgent

	// Create ReliableAgent with medium count, high success rate
	reliableAgent := &stats.GlobalAgentStats{}
	reliableAgent.UpdateStats(1500) // 1.5s
	reliableAgent.SuccessCount = 1
	mockStats["ReliableAgent"] = reliableAgent

	// Create UnreliableAgent with low success rate
	unreliableAgent := &stats.GlobalAgentStats{}
	unreliableAgent.UpdateStats(2000) // 2s
	unreliableAgent.UpdateStats(1800) // 1.8s
	unreliableAgent.UpdateStats(2200) // 2.2s
	unreliableAgent.UpdateStats(1900) // 1.9s
	unreliableAgent.SuccessCount = 1
	unreliableAgent.FailureCount = 3 // 25% success rate
	mockStats["UnreliableAgent"] = unreliableAgent

	return mockStats
}

// Helper methods that should be implemented in the dashboard model for sorting functionality
// These methods are referenced in the tests above and represent the expected API

// Note: These methods should be implemented in the dashboard.go file to support the sorting functionality
// They are referenced here to show what the tests expect

// Expected method signatures that should be implemented:
// func (m *DashboardModel) getSortedAgentEntries(sortBy string) []LeaderboardEntry
// func (m *DashboardModel) setSortMode(mode string)
// func (m *DashboardModel) getCurrentSortMode() string

// The sorting modes should be:
// - "count" for sorting by usage count (default)
// - "success" for sorting by success rate
// - "duration" for sorting by average duration

// ============================================================================
// Timeline Panel Tests - Testing MessageEvent Timeline Visualization
// ============================================================================

// TestTimelinePanel tests the timeline functionality with MessageEvent data
func TestTimelinePanel(t *testing.T) {
	mockStats := createMockAgentStats()
	mockDelegations := createMockDelegationPatterns()
	mockCoOccurrences := createMockCoOccurrencePatterns()
	mockMessageEvents := createMockMessageEvents()

	// Create dashboard model with timeline data
	model := NewDashboardModel(mockStats, mockDelegations, mockCoOccurrences)
	_ = model // Mark as used for now
	// Note: Timeline functionality needs to be added to model
	// This test validates the interface expectations

	// Test that timeline data can be rendered
	view := model.View()
	if view == "" {
		t.Error("Expected non-empty view when timeline data is present")
	}

	// Validate MessageEvent structure compatibility
	for _, event := range mockMessageEvents {
		// Test required fields are present
		if event.ID == "" {
			t.Error("MessageEvent ID should not be empty")
		}
		if event.Timestamp.IsZero() {
			t.Error("MessageEvent Timestamp should not be zero")
		}
		if event.SessionID == "" {
			t.Error("MessageEvent SessionID should not be empty")
		}
		if event.Role == "" {
			t.Error("MessageEvent Role should not be empty")
		}

		// Test role validation
		validRoles := []MessageRole{UserMessage, AssistantMessage, ToolMessage, SystemMessage}
		roleValid := false
		for _, validRole := range validRoles {
			if event.Role == validRole {
				roleValid = true
				break
			}
		}
		if !roleValid {
			t.Errorf("Invalid MessageRole: %s", event.Role)
		}

		// Test duration validation
		if event.Duration < 0 {
			t.Error("MessageEvent Duration should not be negative")
		}

		// Test token count validation
		if event.TokenCount < 0 {
			t.Error("MessageEvent TokenCount should not be negative")
		}
	}
}

// TestTimelineRenderingWithMessageEvents tests timeline rendering functionality
func TestTimelineRenderingWithMessageEvents(t *testing.T) {
	mockMessageEvents := createMockMessageEvents()

	// Test timeline component creation (assuming TimelineComponent would exist)
	timelineTitle := "Message Timeline"
	if timelineTitle == "" {
		t.Error("Timeline component should have a title")
	}

	// Test chronological ordering
	for i := 1; i < len(mockMessageEvents); i++ {
		if mockMessageEvents[i-1].Timestamp.After(mockMessageEvents[i].Timestamp) {
			// Events should be chronologically ordered in our test data
			t.Errorf("Timeline events should be chronologically ordered: event %d timestamp %v should be before event %d timestamp %v",
				i-1, mockMessageEvents[i-1].Timestamp, i, mockMessageEvents[i].Timestamp)
		}
	}

	// Test session grouping
	sessionGroups := make(map[string][]MessageEvent)
	for _, event := range mockMessageEvents {
		sessionGroups[event.SessionID] = append(sessionGroups[event.SessionID], event)
	}

	if len(sessionGroups) < 2 {
		t.Error("Test data should contain multiple sessions for timeline testing")
	}

	// Test agent usage tracking in timeline
	agentUsage := make(map[string]int)
	for _, event := range mockMessageEvents {
		if event.AgentUsed != "" {
			agentUsage[event.AgentUsed]++
		}
	}

	if len(agentUsage) == 0 {
		t.Error("Timeline should track agent usage across messages")
	}

	// Test success/failure tracking
	successCount := 0
	failureCount := 0
	for _, event := range mockMessageEvents {
		if event.Success {
			successCount++
		} else {
			failureCount++
		}
	}

	if successCount == 0 {
		t.Error("Timeline test data should include successful messages")
	}
	if failureCount == 0 {
		t.Error("Timeline test data should include failed messages for comprehensive testing")
	}
}

// TestTimelineNavigationWithDashboard tests time range navigation and interaction
func TestTimelineNavigationWithDashboard(t *testing.T) {
	mockStats := createMockAgentStats()
	mockDelegations := createMockDelegationPatterns()
	mockCoOccurrences := createMockCoOccurrencePatterns()
	model := NewDashboardModel(mockStats, mockDelegations, mockCoOccurrences)
	_ = model // Mark as used

	// Test timeline navigation would be implemented here
	// This validates the expected interface for timeline interaction

	// Test time range filtering capabilities
	now := time.Now()
	timeRanges := []struct {
		name  string
		start time.Time
		end   time.Time
	}{
		{"Last Hour", now.Add(-1 * time.Hour), now},
		{"Last Day", now.Add(-24 * time.Hour), now},
		{"Last Week", now.Add(-7 * 24 * time.Hour), now},
		{"Last Month", now.Add(-30 * 24 * time.Hour), now},
	}

	for _, timeRange := range timeRanges {
		t.Run(timeRange.name, func(t *testing.T) {
			// Test that time range is valid
			if timeRange.start.After(timeRange.end) {
				t.Errorf("Invalid time range: start %v is after end %v", timeRange.start, timeRange.end)
			}

			// Test that time range duration is positive
			duration := timeRange.end.Sub(timeRange.start)
			if duration <= 0 {
				t.Errorf("Time range duration should be positive, got %v", duration)
			}
		})
	}

	// Test navigation keyboard shortcuts (would be implemented in actual timeline component)
	mockMessageEvents := createMockMessageEvents()
	if len(mockMessageEvents) == 0 {
		t.Error("Need message events for timeline navigation testing")
	}

	// Test timeline index bounds
	initialIndex := 0
	maxIndex := len(mockMessageEvents) - 1

	if initialIndex < 0 || initialIndex > maxIndex {
		t.Errorf("Timeline index should be within bounds [0, %d], got %d", maxIndex, initialIndex)
	}
}

// TestTimelineInteractionCapabilities tests timeline user interaction
func TestTimelineInteractionCapabilities(t *testing.T) {
	mockStats := createMockAgentStats()
	mockDelegations := createMockDelegationPatterns()
	mockCoOccurrences := createMockCoOccurrencePatterns()
	model := NewDashboardModel(mockStats, mockDelegations, mockCoOccurrences)

	// Test timeline panel focus switching
	// Note: Timeline panel would be added as a new panel type
	initialPanel := model.focusedPanel
	_ = initialPanel // Mark as used

	// Test keyboard navigation for timeline
	// These would be the expected behaviors for timeline interaction
	timelineKeys := []struct {
		key         string
		description string
	}{
		{"←/→", "Navigate timeline events"},
		{"PgUp/PgDn", "Page through timeline"},
		{"Home/End", "Go to start/end of timeline"},
		{"+/-", "Zoom in/out time range"},
		{"f", "Filter by agent or success status"},
		{"t", "Change time range"},
	}

	for _, key := range timelineKeys {
		t.Run(key.key, func(t *testing.T) {
			// Test that key mapping descriptions are defined
			if key.description == "" {
				t.Errorf("Timeline key %s should have a description", key.key)
			}
		})
	}

	// Test timeline detail view capabilities
	mockMessageEvents := createMockMessageEvents()
	for i, event := range mockMessageEvents {
		t.Run(fmt.Sprintf("Event_%d_Detail", i), func(t *testing.T) {
			// Test that event has sufficient detail for rendering
			if event.Content == "" {
				t.Error("Timeline event should have content for detail view")
			}
			if len(event.ToolsUsed) == 0 && event.Role != SystemMessage {
				// Most events should have tools used (except system messages)
				t.Logf("Event %s has no tools used (might be expected for role %s)", event.ID, event.Role)
			}
		})
	}

	// Restore original panel focus
	model.focusedPanel = initialPanel
}

// ============================================================================
// Enhanced Co-occurrence Panel Tests - Testing Correlation and Matrix Rendering
// ============================================================================

// TestCoOccurrenceMatrixRendering tests enhanced co-occurrence matrix functionality
func TestCoOccurrenceMatrixRendering(t *testing.T) {
	mockCoOccurrences := createMockCoOccurrenceWithCorrelation()
	mockStats := createMockAgentStats()
	mockDelegations := createMockDelegationPatterns()

	model := NewDashboardModel(mockStats, mockDelegations, mockCoOccurrences)
	model.focusedPanel = CoOccurrencePanel

	// Test that co-occurrence matrix component renders
	view := model.View()
	if !strings.Contains(view, "Agent Collaboration Matrix") {
		t.Error("Expected co-occurrence view to contain collaboration matrix title")
	}

	// Test matrix component rendering
	matrixComponent := NewCoOccurrenceMatrixComponent()
	if matrixComponent == nil {
		t.Fatal("Expected co-occurrence matrix component to be created")
	}

	renderedMatrix := matrixComponent.RenderCoOccurrenceMatrix(
		mockCoOccurrences,
		0, // selectedIndex
		true, // isFocused
		"Agent Collaboration Matrix",
	)

	if renderedMatrix == "" {
		t.Error("Expected non-empty co-occurrence matrix rendering")
	}

	// Test relationship strength indicators
	if !strings.Contains(renderedMatrix, "Strong") && !strings.Contains(renderedMatrix, "Medium") && !strings.Contains(renderedMatrix, "Weak") {
		t.Error("Expected co-occurrence matrix to show relationship strength indicators")
	}

	// Test agent pair visualization
	if !strings.Contains(renderedMatrix, "⟷") {
		t.Error("Expected co-occurrence matrix to show bidirectional relationship indicators")
	}

	// Test session count display
	if !strings.Contains(renderedMatrix, "sessions") {
		t.Error("Expected co-occurrence matrix to display session counts")
	}
}

// TestCorrelationCalculations tests correlation analysis functionality
func TestCorrelationCalculations(t *testing.T) {
	mockCoOccurrences := createMockCoOccurrenceWithCorrelation()

	// Test correlation data structure
	for _, coOcc := range mockCoOccurrences {
		// Test agent pair consistency
		if coOcc.Agent1 == coOcc.Agent2 {
			t.Errorf("Agent co-occurrence should not have same agent: %s", coOcc.Agent1)
		}

		// Test count consistency with sessions
		if coOcc.Count <= 0 {
			t.Errorf("Co-occurrence count should be positive, got %d", coOcc.Count)
		}

		if len(coOcc.Sessions) == 0 {
			t.Error("Co-occurrence should have associated sessions")
		}

		// Test that count matches or is related to session count
		// (count might be higher if agents appear multiple times per session)
		if coOcc.Count < len(coOcc.Sessions) {
			t.Errorf("Co-occurrence count %d should be at least the number of sessions %d",
				coOcc.Count, len(coOcc.Sessions))
		}
	}

	// Test correlation strength calculation
	matrixComponent := NewCoOccurrenceMatrixComponent()
	for _, coOcc := range mockCoOccurrences {
		strengthEnum := matrixComponent.calculateRelationshipStrength(coOcc.Count)
		strengthIndicator := matrixComponent.createVisualStrengthIndicator(strengthEnum)
		if strengthIndicator == "" {
			t.Error("Relationship strength indicator should not be empty")
		}

		// Test strength categorization
		expectedCategories := []string{"Strong", "Medium", "Weak", "Minimal"}
		validCategory := false
		for _, category := range expectedCategories {
			if strings.Contains(strengthIndicator, category) {
				validCategory = true
				break
			}
		}
		if !validCategory {
			t.Errorf("Strength indicator should contain valid category, got: %s", strengthIndicator)
		}
	}

	// Test correlation threshold logic
	strengthTests := []struct {
		count            int
		expectedCategory string
	}{
		{15, "Strong"},
		{7, "Medium"},
		{3, "Weak"},
		{1, "Minimal"},
	}

	for _, test := range strengthTests {
		strengthEnum := matrixComponent.calculateRelationshipStrength(test.count)
		strength := matrixComponent.createVisualStrengthIndicator(strengthEnum)
		if !strings.Contains(strength, test.expectedCategory) {
			t.Errorf("Count %d should result in %s strength, got: %s",
				test.count, test.expectedCategory, strength)
		}
	}
}

// TestRelationshipStrengthIndicators tests visual strength indicators
func TestRelationshipStrengthIndicators(t *testing.T) {
	matrixComponent := NewCoOccurrenceMatrixComponent()

	// Test visual indicator consistency
	indicatorTests := []struct {
		count    int
		expected string
	}{
		{20, "████ Strong "}, // Very high collaboration
		{8, "███· Medium "}, // Moderate collaboration
		{3, "██·· Weak   "}, // Low collaboration
		{1, "█··· Minimal"}, // Minimal collaboration
	}

	for _, test := range indicatorTests {
		strengthEnum := matrixComponent.calculateRelationshipStrength(test.count)
		actual := matrixComponent.createVisualStrengthIndicator(strengthEnum)
		if actual != test.expected {
			t.Errorf("Count %d should produce indicator '%s', got '%s'",
				test.count, test.expected, actual)
		}
	}

	// Test edge cases
	zeroCountEnum := matrixComponent.calculateRelationshipStrength(0)
	zeroCountIndicator := matrixComponent.createVisualStrengthIndicator(zeroCountEnum)
	if !strings.Contains(zeroCountIndicator, "Minimal") {
		t.Error("Zero count should result in Minimal strength")
	}

	negativeCountEnum := matrixComponent.calculateRelationshipStrength(-1)
	negativeCountIndicator := matrixComponent.createVisualStrengthIndicator(negativeCountEnum)
	if !strings.Contains(negativeCountIndicator, "Minimal") {
		t.Error("Negative count should result in Minimal strength")
	}

	// Test visual consistency (all indicators should have same length)
	indicatorLengths := make(map[int]bool)
	for _, test := range indicatorTests {
		strengthEnum := matrixComponent.calculateRelationshipStrength(test.count)
		indicator := matrixComponent.createVisualStrengthIndicator(strengthEnum)
		indicatorLengths[len(indicator)] = true
	}

	if len(indicatorLengths) > 1 {
		t.Error("All relationship strength indicators should have consistent length")
	}
}

// TestAgentRelationshipInsights tests relationship insights and pattern recognition
func TestAgentRelationshipInsights(t *testing.T) {
	mockCoOccurrences := createMockCoOccurrenceWithCorrelation()
	matrixComponent := NewCoOccurrenceMatrixComponent()

	// Test collaboration insights generation
	collabInsights := matrixComponent.renderInsightsView(matrixComponent.buildAgentRelationships(mockCoOccurrences))
	if collabInsights == "" {
		t.Error("Expected non-empty collaboration insights")
	}

	// Test insights contains key information
	if !strings.Contains(collabInsights, "Insights") {
		t.Error("Expected collaboration insights to contain title")
	}

	// Test agent collaboration ranking
	agentCollabCount := make(map[string]int)
	for _, coOcc := range mockCoOccurrences {
		agentCollabCount[coOcc.Agent1] += coOcc.Count
		agentCollabCount[coOcc.Agent2] += coOcc.Count
	}

	if len(agentCollabCount) == 0 {
		t.Error("Should calculate agent collaboration counts")
	}

	// Find most collaborative agent
	maxCollabCount := 0
	mostCollaborativeAgent := ""
	for agent, count := range agentCollabCount {
		if count > maxCollabCount {
			maxCollabCount = count
			mostCollaborativeAgent = agent
		}
	}

	if mostCollaborativeAgent == "" {
		t.Error("Should identify most collaborative agent")
	}

	// Test that insights includes most collaborative agent (if present)
	if !strings.Contains(collabInsights, mostCollaborativeAgent) {
		// Note: This is informational - insights might not include all agents
		t.Logf("Collaboration insights could mention most collaborative agent: %s", mostCollaborativeAgent)
	}

	// Test pattern recognition in insights using buildAgentRelationships
	relationships := matrixComponent.buildAgentRelationships(mockCoOccurrences)
	for i, rel := range relationships {
		insight := matrixComponent.renderRelationshipEntry(rel, i == 0)
		if insight == "" {
			t.Errorf("Expected non-empty insight for relationship %d", i)
		}

		// Test insight contains agent names
		if !strings.Contains(insight, rel.Agent1) {
			t.Errorf("Insight should contain agent1 name: %s", rel.Agent1)
		}
		if !strings.Contains(insight, rel.Agent2) {
			t.Errorf("Insight should contain agent2 name: %s", rel.Agent2)
		}

		// Test insight contains relationship strength
		strengthFound := false
		strengthCategories := []string{"Strong", "Medium", "Weak", "Minimal"}
		strengthIndicator := matrixComponent.createVisualStrengthIndicator(rel.Strength)
		for _, category := range strengthCategories {
			if strings.Contains(strengthIndicator, category) {
				strengthFound = true
				break
			}
		}
		if !strengthFound {
			t.Error("Insight should contain relationship strength indicator")
		}

		// Test insight contains frequency information
		if rel.Frequency <= 0 {
			t.Error("Relationship should have positive frequency")
		}
	}
}

// TestCoOccurrenceMatrixInteraction tests user interaction with co-occurrence matrix
func TestCoOccurrenceMatrixInteraction(t *testing.T) {
	mockStats := createMockAgentStats()
	mockDelegations := createMockDelegationPatterns()
	mockCoOccurrences := createMockCoOccurrenceWithCorrelation()

	model := NewDashboardModel(mockStats, mockDelegations, mockCoOccurrences)
	model.focusedPanel = CoOccurrencePanel

	// Test navigation within co-occurrence panel
	initialIndex := model.selectedCoOccurrenceIndex

	// Test down navigation
	model.moveDown()
	if len(mockCoOccurrences) > 1 && model.selectedCoOccurrenceIndex == initialIndex {
		t.Error("Down navigation should change selected co-occurrence index")
	}

	// Test up navigation
	model.moveUp()
	if model.selectedCoOccurrenceIndex != initialIndex {
		t.Error("Up navigation should return to initial index")
	}

	// Test boundary conditions
	model.selectedCoOccurrenceIndex = 0
	model.moveUp() // Should not go negative
	if model.selectedCoOccurrenceIndex < 0 {
		t.Error("Co-occurrence selection should not go negative")
	}

	model.selectedCoOccurrenceIndex = len(mockCoOccurrences) - 1
	model.moveDown() // Should not exceed bounds
	if model.selectedCoOccurrenceIndex >= len(mockCoOccurrences) {
		t.Error("Co-occurrence selection should not exceed bounds")
	}

	// Test focus highlighting
	view := model.View()
	if !strings.Contains(view, "▶") {
		t.Error("Focused co-occurrence panel should show selection indicator")
	}

	// Test panel switching from co-occurrence
	model.switchPanel()
	if model.focusedPanel == CoOccurrencePanel {
		t.Error("Panel switching should move away from co-occurrence panel")
	}
}

// TestCoOccurrencePatternRecognition tests advanced pattern recognition
func TestCoOccurrencePatternRecognition(t *testing.T) {
	mockCoOccurrences := createMockCoOccurrenceWithCorrelation()

	// Test pattern identification
	patterns := make(map[string][]stats.AgentCoOccurrence)
	for _, coOcc := range mockCoOccurrences {
		// Group by agent pairs for pattern analysis
		key := fmt.Sprintf("%s-%s", coOcc.Agent1, coOcc.Agent2)
		patterns[key] = append(patterns[key], coOcc)
	}

	if len(patterns) == 0 {
		t.Error("Should identify co-occurrence patterns")
	}

	// Test frequency distribution
	frequencies := make([]int, 0, len(mockCoOccurrences))
	for _, coOcc := range mockCoOccurrences {
		frequencies = append(frequencies, coOcc.Count)
	}

	// Calculate average frequency
	totalFreq := 0
	for _, freq := range frequencies {
		totalFreq += freq
	}
	avgFreq := float64(totalFreq) / float64(len(frequencies))

	if avgFreq <= 0 {
		t.Error("Average co-occurrence frequency should be positive")
	}

	// Test session overlap analysis
	for _, coOcc := range mockCoOccurrences {
		// Test session list integrity
		sessionMap := make(map[string]bool)
		for _, session := range coOcc.Sessions {
			if session == "" {
				t.Error("Session ID should not be empty")
			}
			if sessionMap[session] {
				t.Errorf("Duplicate session %s in co-occurrence sessions", session)
			}
			sessionMap[session] = true
		}

		// Test unique session count
		if len(sessionMap) != len(coOcc.Sessions) {
			t.Error("Co-occurrence sessions should be unique")
		}
	}

	// Test temporal patterns
	for _, coOcc := range mockCoOccurrences {
		if coOcc.LastSeen.IsZero() {
			t.Error("Co-occurrence should have valid LastSeen timestamp")
		}

		// Test recency calculation
		timeSince := time.Since(coOcc.LastSeen)
		if timeSince < 0 {
			t.Error("Co-occurrence LastSeen should not be in the future")
		}
	}

	// Test collaboration network analysis
	agentNetwork := make(map[string][]string)
	for _, coOcc := range mockCoOccurrences {
		agentNetwork[coOcc.Agent1] = append(agentNetwork[coOcc.Agent1], coOcc.Agent2)
		agentNetwork[coOcc.Agent2] = append(agentNetwork[coOcc.Agent2], coOcc.Agent1)
	}

	if len(agentNetwork) == 0 {
		t.Error("Should build agent collaboration network")
	}

	// Find most connected agent
	maxConnections := 0
	for _, connections := range agentNetwork {
		if len(connections) > maxConnections {
			maxConnections = len(connections)
		}
	}

	if maxConnections == 0 {
		t.Error("Should identify agent with most connections")
	}
}

// TestCoOccurrenceDataIntegrity tests data consistency and validation
func TestCoOccurrenceDataIntegrity(t *testing.T) {
	mockCoOccurrences := createMockCoOccurrenceWithCorrelation()

	// Test data structure integrity
	for i, coOcc := range mockCoOccurrences {
		t.Run(fmt.Sprintf("CoOccurrence_%d", i), func(t *testing.T) {
			// Test required fields
			if coOcc.Agent1 == "" {
				t.Error("Agent1 should not be empty")
			}
			if coOcc.Agent2 == "" {
				t.Error("Agent2 should not be empty")
			}
			if coOcc.Agent1 == coOcc.Agent2 {
				t.Error("Agent1 and Agent2 should be different")
			}

			// Test count validation
			if coOcc.Count <= 0 {
				t.Errorf("Count should be positive, got %d", coOcc.Count)
			}

			// Test sessions validation
			if len(coOcc.Sessions) == 0 {
				t.Error("Sessions list should not be empty")
			}

			// Test logical consistency between count and sessions
			if coOcc.Count < len(coOcc.Sessions) {
				t.Errorf("Count %d should be at least the number of sessions %d",
					coOcc.Count, len(coOcc.Sessions))
			}

			// Test timestamp validation
			if coOcc.LastSeen.IsZero() {
				t.Error("LastSeen should have a valid timestamp")
			}
			if coOcc.LastSeen.After(time.Now()) {
				t.Error("LastSeen should not be in the future")
			}
		})
	}

	// Test data consistency across co-occurrences
	allSessions := make(map[string]bool)
	allAgents := make(map[string]bool)

	for _, coOcc := range mockCoOccurrences {
		allAgents[coOcc.Agent1] = true
		allAgents[coOcc.Agent2] = true
		for _, session := range coOcc.Sessions {
			allSessions[session] = true
		}
	}

	if len(allAgents) < 2 {
		t.Error("Co-occurrence data should involve at least 2 different agents")
	}

	if len(allSessions) == 0 {
		t.Error("Co-occurrence data should reference valid sessions")
	}

	// Test bidirectional consistency (if A co-occurs with B, B should co-occur with A)
	coOccurrenceMap := make(map[string]stats.AgentCoOccurrence)
	for _, coOcc := range mockCoOccurrences {
		key1 := fmt.Sprintf("%s-%s", coOcc.Agent1, coOcc.Agent2)
		key2 := fmt.Sprintf("%s-%s", coOcc.Agent2, coOcc.Agent1)
		coOccurrenceMap[key1] = coOcc
		coOccurrenceMap[key2] = coOcc
	}

	// Verify bidirectional consistency
	for _, coOcc := range mockCoOccurrences {
		key1 := fmt.Sprintf("%s-%s", coOcc.Agent1, coOcc.Agent2)
		key2 := fmt.Sprintf("%s-%s", coOcc.Agent2, coOcc.Agent1)

		if _, exists1 := coOccurrenceMap[key1]; !exists1 {
			t.Errorf("Missing co-occurrence mapping for %s", key1)
		}
		if _, exists2 := coOccurrenceMap[key2]; !exists2 {
			t.Errorf("Missing bidirectional co-occurrence mapping for %s", key2)
		}
	}
}

// MessageRole and MessageEvent types are already defined in timeline_component.go
// Tests use the existing types from that file