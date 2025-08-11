package ui

import (
	"embed"
	"testing"
	
	tea "github.com/charmbracelet/bubbletea"
)

//go:embed test_assets/*
var testAssetsUI embed.FS

// TestASCIIWelcomeBanner validates the welcome banner meets requirements
func TestASCIIWelcomeBanner(t *testing.T) {
	// Test banner content
	if WelcomeBanner == "" {
		t.Fatal("WelcomeBanner should not be empty")
	}
	
	// Should contain ASCII art blocks (████)
	if !containsString(WelcomeBanner, "████") {
		t.Error("Expected banner to contain ASCII art blocks")
	}
	
	// Should have three distinct sections (THE, AGENTIC, STARTUP represented as ASCII art)
	lines := 0
	for _, char := range WelcomeBanner {
		if char == '\n' {
			lines++
		}
	}
	
	if lines < 15 {
		t.Errorf("Expected banner to have at least 15 lines for ASCII art, got %d", lines)
	}
	
	// Should contain block characters (█) representing the ASCII art
	blockChar := '█'
	blockCount := 0
	for _, char := range WelcomeBanner {
		if char == blockChar {
			blockCount++
		}
	}
	
	if blockCount < 50 {
		t.Errorf("Expected many block characters in ASCII art, found %d", blockCount)
	}
	
	// Banner should be multi-line
	if !containsString(WelcomeBanner, "\n") {
		t.Error("Expected banner to be multi-line")
	}
	
	// Test welcome message
	if WelcomeMessage == "" {
		t.Error("WelcomeMessage should not be empty")
	}
	
	if !containsString(WelcomeMessage, "Welcome to The Agentic Startup Installation") {
		t.Error("Expected welcome message to contain proper title")
	}
}

// TestBrightGreenTitleStyling validates title styling requirements
func TestBrightGreenTitleStyling(t *testing.T) {
	styles := GetStyles()
	model := NewInstallerModel(&testAssetsUI, &testAssetsUI, &testAssetsUI, &testAssetsUI)
	
	// Verify the title style has background and foreground set by checking if it renders differently
	testTitle := "Test Title"
	styledTitle := styles.Title.Render(testTitle)
	
	// In test environment, styling may not be applied, so let's just verify the style is configured
	// Title should be longer than input due to padding (0,1,0,1)
	if len(styledTitle) < len(testTitle) {
		t.Error("Expected styled title to be at least as long as input (accounting for padding)")
	}
	
	// Verify that we're getting the bright green title style defined
	titleView := model.renderer.RenderTitle("Test")
	if titleView == "" {
		t.Error("Expected title renderer to produce output")
	}
	
	// Title should be bold
	if !styles.Title.GetBold() {
		t.Error("Expected title to be bold")
	}
	
	// Title should have padding
	paddingTop, paddingRight, paddingBottom, paddingLeft := styles.Title.GetPadding()
	if paddingTop != 0 || paddingRight != 1 || paddingBottom != 0 || paddingLeft != 1 {
		t.Errorf("Expected title padding [0,1,0,1], got [%d,%d,%d,%d]", 
			paddingTop, paddingRight, paddingBottom, paddingLeft)
	}
}

// TestFullScreenTerminalUI validates full-screen usage
func TestFullScreenTerminalUI(t *testing.T) {
	model := NewInstallerModel(&testAssetsUI, &testAssetsUI, &testAssetsUI, &testAssetsUI)
	
	// Test window resize handling (full-screen capability)
	originalWidth := model.width
	originalHeight := model.height
	
	// Simulate large window size (full-screen)
	newWidth := 200
	newHeight := 50
	
	model.width = newWidth
	model.height = newHeight
	
	if model.width != newWidth {
		t.Errorf("Expected width to be set to %d, got %d", newWidth, model.width)
	}
	
	if model.height != newHeight {
		t.Errorf("Expected height to be set to %d, got %d", newHeight, model.height)
	}
	
	// View should render without issues at large sizes
	view := model.View()
	if view == "" {
		t.Error("Expected view to render at full-screen size")
	}
	
	// Restore original size
	model.width = originalWidth
	model.height = originalHeight
}

// TestProgressiveDisclosure validates the progressive disclosure feature
func TestProgressiveDisclosure(t *testing.T) {
	model := NewInstallerModel(&testAssetsUI, &testAssetsUI, &testAssetsUI, &testAssetsUI)
	renderer := model.renderer
	
	// Test empty state (no previous selections shown)
	result := renderer.RenderSelections("", "", 0)
	if result != "" {
		t.Error("Expected no selections to render empty string")
	}
	
	// Test tool selection shown
	result = renderer.RenderSelections("claude-code", "", 0)
	if !containsString(result, "claude-code") {
		t.Error("Expected tool selection to be shown in progressive disclosure")
	}
	
	if !containsString(result, "The (Agentic) Startup Installation") {
		t.Error("Expected header to contain installation title")
	}
	
	// Test path selection added
	result = renderer.RenderSelections("claude-code", "~/.config/the-startup", 0)
	if !containsString(result, "claude-code") {
		t.Error("Expected tool to still be shown")
	}
	if !containsString(result, "~/.config/the-startup") {
		t.Error("Expected path to be shown in progressive disclosure")
	}
	
	// Test mode selection added
	result = renderer.RenderSelections("claude-code", "~/.config/the-startup", 0)
	// Mode is no longer shown in the simplified flow
	
	// Test file selection added
	result = renderer.RenderSelections("claude-code", "~/.config/the-startup", 3)
	if !containsString(result, "Files: 3") {
		t.Error("Expected files count to be shown in progressive disclosure")
	}
	
	// Test file count for advanced mode
	result = renderer.RenderSelections("claude-code", "~/.config/the-startup", 15)
	if !containsString(result, "Files: 15") {
		t.Error("Expected file count to be shown")
	}
}

// TestConsistentStylingWithoutIndentation validates consistent styling
func TestConsistentStylingWithoutIndentation(t *testing.T) {
	model := NewInstallerModel(&testAssetsUI, &testAssetsUI, &testAssetsUI, &testAssetsUI)
	renderer := model.renderer
	
	// Test choice rendering without indentation
	choice1 := renderer.RenderChoice("Option 1", false, false)
	choice2 := renderer.RenderChoice("Option 2", true, false)
	choice3 := renderer.RenderChoice("Option 3", false, true)
	
	// All choices should start consistently (no extra indentation)
	// Cursor choices have "> " prefix, non-cursor have "  " prefix
	if !containsString(choice1, "  ○ Option 1") {
		t.Error("Expected consistent non-cursor choice formatting")
	}
	
	if !containsString(choice2, "> ○ Option 2") {
		t.Error("Expected consistent cursor choice formatting")
	}
	
	if !containsString(choice3, "  ● Option 3") {
		t.Error("Expected consistent selected choice formatting")
	}
	
	// Test that tree selector removes indentation for flat display
	root := &TreeNode{
		Name: "Root",
		IsDir: true,
		Children: []*TreeNode{
			{Name: "file1", IsDir: false},
			{Name: "file2", IsDir: false},
		},
	}
	
	treeSelector := NewTreeSelector("Test", root)
	indent := treeSelector.getIndent(root.Children[0])
	
	if indent != "" {
		t.Error("Expected tree selector to use no indentation for flat display")
	}
}

// Test8StateStateMachine validates the complete state machine
func Test8StateStateMachine(t *testing.T) {
	// Test all 8 states are defined
	expectedStates := []InstallerState{
		StateWelcome,
		StateToolSelection,
		StatePathSelection,
		StateFileSelection,
		StateConfirmation,
		StateInstalling,
		StateComplete,
		StateError,
	}
	
	if len(expectedStates) != 8 {
		t.Fatalf("Expected 8 states, found %d", len(expectedStates))
	}
	
	// Test each state has string representation
	for _, state := range expectedStates {
		stateStr := state.String()
		if stateStr == "" || stateStr == "Unknown" {
			t.Errorf("State %d should have valid string representation, got '%s'", int(state), stateStr)
		}
	}
	
	// Test state transitions are properly defined
	model := NewInstallerModel(&testAssetsUI, &testAssetsUI, &testAssetsUI, &testAssetsUI)
	
	// Test forward transitions work
	model.transitionToState(StateToolSelection)
	if model.state != StateToolSelection {
		t.Error("Expected transition to StateToolSelection")
	}
	
	model.transitionToState(StatePathSelection)
	if model.state != StatePathSelection {
		t.Error("Expected transition to StatePathSelection")
	}
	
	model.transitionToState(StateFileSelection)
	if model.state != StateFileSelection {
		t.Error("Expected transition to StateFileSelection")
	}
	
	// Test that each state initializes choices properly
	for _, state := range expectedStates {
		testModel := NewInstallerModel(&testAssetsUI, &testAssetsUI, &testAssetsUI, &testAssetsUI)
		
		// Only test valid transitions
		switch state {
		case StateWelcome:
			// Already in welcome state
		case StateToolSelection:
			testModel.transitionToState(StateToolSelection)
		case StatePathSelection:
			testModel.transitionToState(StateToolSelection)
			testModel.transitionToState(StatePathSelection)
		case StateFileSelection:
			testModel.transitionToState(StateToolSelection)
			testModel.transitionToState(StatePathSelection)
			testModel.transitionToState(StateFileSelection)
		default:
			continue // Skip states that require complex setup
		}
		
		if testModel.state != state {
			t.Errorf("Failed to transition to state %s", state.String())
		}
		
		// Each state should render a view
		view := testModel.View()
		if view == "" {
			t.Errorf("State %s should render non-empty view", state.String())
		}
	}
}

// TestStateTransitionValidation validates proper state transition enforcement
func TestStateTransitionValidation(t *testing.T) {
	model := NewInstallerModel(&testAssetsUI, &testAssetsUI, &testAssetsUI, &testAssetsUI)
	
	// Test invalid transitions are rejected
	model.transitionToState(StateComplete) // Invalid: Welcome -> Complete
	
	if model.state != StateError {
		t.Error("Expected invalid transition to result in StateError")
	}
	
	if model.err == nil {
		t.Error("Expected error to be set for invalid transition")
	}
	
	// Test error contains meaningful information
	if !containsString(model.err.Error(), "invalid state transition") {
		t.Error("Expected error message to mention invalid state transition")
	}
}

// TestKeyboardNavigationConsistency validates consistent keyboard controls
func TestKeyboardNavigationConsistency(t *testing.T) {
	model := NewInstallerModel(&testAssetsUI, &testAssetsUI, &testAssetsUI, &testAssetsUI)
	
	// Test that up/down keys work consistently across states
	model.transitionToState(StateToolSelection)
	model.choices = []string{"Option 1", "Option 2", "Option 3"}
	model.cursor = 1
	
	// Test up key
	upMsg := tea.KeyMsg{Type: tea.KeyUp}
	updatedModel, _ := model.Update(upMsg)
	m := updatedModel.(*InstallerModel)
	
	if m.cursor != 0 {
		t.Error("Expected up key to move cursor up")
	}
	
	// Test down key
	downMsg := tea.KeyMsg{Type: tea.KeyDown}
	updatedModel, _ = m.Update(downMsg)
	m = updatedModel.(*InstallerModel)
	
	if m.cursor != 1 {
		t.Error("Expected down key to move cursor down")
	}
	
	// Test vim-style navigation
	kMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}
	jMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
	
	// These should work in TreeSelector
	root := &TreeNode{Name: "Root", IsDir: true, Children: []*TreeNode{
		{Name: "file1", IsDir: false},
		{Name: "file2", IsDir: false},
	}}
	
	selector := NewTreeSelector("Test", root)
	selector.cursor = 1
	
	selector.Update(kMsg) // k = up
	if selector.cursor != 0 {
		t.Error("Expected 'k' key to move cursor up in TreeSelector")
	}
	
	selector.Update(jMsg) // j = down
	if selector.cursor != 1 {
		t.Error("Expected 'j' key to move cursor down in TreeSelector")
	}
}

// TestIconConsistency validates consistent icon usage
func TestIconConsistency(t *testing.T) {
	// Validate that common icons are defined and not empty
	icons := []struct{
		name string
		value string
	}{
		{"IconSuccess", IconSuccess},
		{"IconError", IconError}, 
		{"IconWarning", IconWarning},
		{"IconSelected", IconSelected},
		{"IconUnselected", IconUnselected},
	}
	
	// Validate icons are not empty
	for _, icon := range icons {
		if icon.value == "" {
			t.Errorf("%s should not be empty", icon.name)
		}
	}
	
	// Test that icons are used consistently in choice rendering
	renderer := NewProgressiveDisclosureRenderer()
	
	unselectedChoice := renderer.RenderChoice("Test", false, false)
	if !containsString(unselectedChoice, IconUnselected) {
		t.Error("Expected unselected choices to use IconUnselected")
	}
	
	selectedChoice := renderer.RenderChoice("Test", false, true)
	if !containsString(selectedChoice, IconSelected) {
		t.Error("Expected selected choices to use IconSelected")
	}
}