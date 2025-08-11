package ui

import (
	"embed"
	"strings"
	"testing"
)

//go:embed test_assets/*
var testAssetsRefinements embed.FS

// TestSingleSelectNoCircles verifies that single-select lists don't show circles
func TestSingleSelectNoCircles(t *testing.T) {
	model := NewInstallerModel(&testAssetsRefinements, &testAssetsRefinements, &testAssetsRefinements, &testAssetsRefinements)

	// Test tool selection (single-select) - should NOT have circles
	model.transitionToState(StateToolSelection)
	view := model.View()
	
	// Should not contain circle symbols in tool selection
	if strings.Contains(view, "○") || strings.Contains(view, "●") {
		t.Error("Tool selection (single-select) should not contain circle symbols")
	}
	
	// Should still have cursor indicator
	if !strings.Contains(view, "> ") {
		t.Error("Tool selection should still have cursor indicator")
	}

	// Test path selection (single-select) - should NOT have circles
	model.selectedTool = "claude-code"
	model.transitionToState(StatePathSelection)
	view = model.View()
	
	if strings.Contains(view, "○") || strings.Contains(view, "●") {
		t.Error("Path selection (single-select) should not contain circle symbols")
	}

	// Test file selection (single-select) - should NOT have circles
	model.selectedPath = "~/.config/the-startup"
	model.selectedFiles = []string{"agents/test.md"}
	model.transitionToState(StateFileSelection)
	view = model.View()
	
	if strings.Contains(view, "○") || strings.Contains(view, "●") {
		t.Error("File selection (single-select) should not contain circle symbols")
	}

	// Test confirmation (single-select) - should NOT have circles
	model.transitionToState(StateConfirmation)
	view = model.View()
	
	if strings.Contains(view, "○") || strings.Contains(view, "●") {
		t.Error("Confirmation selection (single-select) should not contain circle symbols")
	}
}

// TestFileSelectionView verifies that file selection works as expected  
func TestFileSelectionView(t *testing.T) {
	model := NewInstallerModel(&testAssetsRefinements, &testAssetsRefinements, &testAssetsRefinements, &testAssetsRefinements)

	// Setup for file selection - go through proper transitions
	model.transitionToState(StateToolSelection)
	model.selectedTool = "claude-code"
	model.transitionToState(StatePathSelection)
	model.selectedPath = "~/.config/the-startup"
	model.transitionToState(StateFileSelection)
	
	view := model.View()
	
	// Should show informative message about files being moved
	if !strings.Contains(view, "files will be moved") {
		t.Error("File selection should show informative message about files being moved")
	}
}

// TestWelcomeAutoAdvance verifies that welcome screen auto-advances after 2 seconds
func TestWelcomeAutoAdvance(t *testing.T) {
	model := NewInstallerModel(&testAssetsRefinements, &testAssetsRefinements, &testAssetsRefinements, &testAssetsRefinements)
	
	// Should start in welcome state
	if model.state != StateWelcome {
		t.Fatalf("Expected initial state to be Welcome, got %s", model.state.String())
	}
	
	// Init should return a tick command for auto-advance
	cmd := model.Init()
	if cmd == nil {
		t.Error("Expected Init to return a tick command for auto-advance")
	}
	
	// Simulate waiting 2+ seconds by sending a tick message
	tickMsg := TickMsg{} // This will use the existing TickMsg from progress.go
	updatedModel, _ := model.Update(tickMsg)
	m := updatedModel.(*InstallerModel)
	
	// Should have transitioned to tool selection
	if m.state != StateToolSelection {
		t.Errorf("Expected state to transition to ToolSelection after tick, got %s", m.state.String())
	}
}

// TestSingleScreenViews verifies that views show one screen at a time with banner
func TestSingleScreenViews(t *testing.T) {
	model := NewInstallerModel(&testAssetsRefinements, &testAssetsRefinements, &testAssetsRefinements, &testAssetsRefinements)
	
	// Get welcome view
	welcomeView := model.View()
	if !strings.Contains(welcomeView, "████████") || !strings.Contains(welcomeView, "Welcome to The Agentic Startup") {
		t.Error("Expected welcome view to contain ASCII banner")
	}
	
	// Transition to tool selection
	model.transitionToState(StateToolSelection)
	toolView := model.View()
	
	// Should contain banner at top of screen
	if !strings.Contains(toolView, "████████") {
		t.Error("Tool selection view should contain ASCII banner")
	}
	
	// Should contain tool selection content
	if !strings.Contains(toolView, "Select your development tool") {
		t.Error("Tool selection view should contain tool selection content")
	}
	
	// Transition to path selection
	model.selectedTool = "claude-code"
	model.transitionToState(StatePathSelection)
	pathView := model.View()
	
	// Should contain banner at top of screen
	if !strings.Contains(pathView, "████████") {
		t.Error("Path selection view should contain ASCII banner")
	}
	
	// Should show previous tool selection in header (but not welcome message)
	if !strings.Contains(pathView, "Tool: claude-code") {
		t.Error("Path selection view should show previous tool selection in header")
	}
	
	if !strings.Contains(pathView, "Select installation location") {
		t.Error("Path selection view should contain current path selection content")
	}
}

// TestRenderChoiceMethodSignature verifies the RenderChoice method supports multi-select parameter
func TestRenderChoiceMethodSignature(t *testing.T) {
	renderer := NewProgressiveDisclosureRenderer()
	
	// Test single-select rendering (no circles)
	singleSelectChoice := renderer.RenderChoiceWithMultiSelect("Option 1", false, false, false)
	if strings.Contains(singleSelectChoice, "○") || strings.Contains(singleSelectChoice, "●") {
		t.Error("Single-select choice should not contain circles")
	}
	
	// Test multi-select rendering (with circles)
	multiSelectChoice := renderer.RenderChoiceWithMultiSelect("Option 1", false, false, true)
	if !strings.Contains(multiSelectChoice, "○") {
		t.Error("Multi-select choice should contain circle")
	}
	
	// Test multi-select selected rendering
	multiSelectSelected := renderer.RenderChoiceWithMultiSelect("Option 1", false, true, true)
	if !strings.Contains(multiSelectSelected, "●") {
		t.Error("Multi-select selected choice should contain filled circle")
	}
}

// TestCompletionScreenShowsSingleScreen verifies that completion screen shows its own complete screen with banner
func TestCompletionScreenShowsSingleScreen(t *testing.T) {
	model := NewInstallerModel(&testAssetsRefinements, &testAssetsRefinements, &testAssetsRefinements, &testAssetsRefinements)
	
	// Set up state for completion
	model.transitionToState(StateToolSelection)
	model.selectedTool = "claude-code"
	model.transitionToState(StatePathSelection)
	model.selectedPath = "~/.config/the-startup"
	model.selectedFiles = []string{"agents/test.md", "commands/test.md"}
	model.transitionToState(StateFileSelection)
	model.transitionToState(StateHuhConfirmation)
	model.transitionToState(StateConfirmation)
	
	// Simulate completion (must go through Installing state)
	model.transitionToState(StateInstalling)
	model.transitionToState(StateComplete)
	completeView := model.View()
	
	// Debug: log the complete view
	t.Logf("Complete view: %s", completeView)
	
	// The completion view should contain banner at top
	if !strings.Contains(completeView, "████████") {
		t.Error("Completion screen should contain ASCII banner")
	}
	
	// Should have completion message
	if !strings.Contains(completeView, "✅ Installation Complete!") {
		t.Error("Completion screen should contain completion message")
	}
	
	// Should show installation path
	if !strings.Contains(completeView, "has been installed") {
		t.Error("Completion screen should show installation details")
	}
}

// TestFileSelectionShowsStaticTreeDisplay verifies that file selection shows static tree
func TestFileSelectionShowsStaticTreeDisplay(t *testing.T) {
	model := NewInstallerModel(&testAssetsRefinements, &testAssetsRefinements, &testAssetsRefinements, &testAssetsRefinements)
	
	// Navigate to file selection
	model.transitionToState(StateToolSelection)
	model.selectedTool = "claude-code"
	model.transitionToState(StatePathSelection)
	model.selectedPath = "~/.config/the-startup"
	model.selectedFiles = []string{"agents/test.md", "commands/test.md"}
	model.transitionToState(StateFileSelection)
	
	view := model.View()
	
	// Should show static tree preview
	if strings.Contains(view, "space: toggle") || strings.Contains(view, "a: toggle all") {
		t.Error("File selection should show static tree preview, not interactive controls")
	}
	
	// Should show informative message
	if !strings.Contains(view, "files will be moved") {
		t.Error("File selection should show informative message about files being moved")
	}
}

// TestFileSelectionToConfirmation verifies that file selection goes to confirmation
func TestFileSelectionToConfirmation(t *testing.T) {
	model := NewInstallerModel(&testAssetsRefinements, &testAssetsRefinements, &testAssetsRefinements, &testAssetsRefinements)
	
	// Navigate to file selection
	model.transitionToState(StateToolSelection)
	model.selectedTool = "claude-code"
	model.transitionToState(StatePathSelection)
	model.selectedPath = "~/.config/the-startup"
	model.selectedFiles = []string{"agents/test.md"}
	model.transitionToState(StateFileSelection)
	
	// Press enter to continue
	model.cursor = 0 // "Continue" option
	updatedModel, _ := model.handleEnter()
	m := updatedModel.(*InstallerModel)
	
	// Debug: check for error
	if m.err != nil {
		t.Logf("Error occurred: %v (context: %s)", m.err, m.errContext)
	}
	
	// Should transition to huh confirmation first
	if m.state != StateHuhConfirmation {
		t.Errorf("File selection should transition to StateHuhConfirmation, got %s", m.state.String())
	}
	
	// Should have files selected
	if len(m.selectedFiles) == 0 {
		t.Error("Should have files selected for confirmation")
	}
}