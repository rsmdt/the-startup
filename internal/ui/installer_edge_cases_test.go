package ui

import (
	"embed"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

//go:embed test_assets/*
var testAssetsEdge embed.FS

func TestInvalidStateTransitions(t *testing.T) {
	model := NewInstallerModel(&testAssetsEdge, &testAssetsEdge, &testAssetsEdge, &testAssetsEdge)
	
	// Try to transition from Welcome directly to Complete (invalid)
	model.transitionToState(StateComplete)
	
	// Should have moved to Error state instead
	if model.state != StateError {
		t.Errorf("Expected state to be StateError after invalid transition, got %v", model.state)
	}
	
	if model.err == nil {
		t.Error("Expected error to be set after invalid transition")
	}
}

func TestEmptyChoicesHandling(t *testing.T) {
	model := NewInstallerModel(&testAssetsEdge, &testAssetsEdge, &testAssetsEdge, &testAssetsEdge)
	
	// Set cursor beyond choices length
	model.cursor = 100
	
	// Try to handle enter - should not crash
	updatedModel, _ := model.handleEnter()
	m := updatedModel.(*InstallerModel)
	
	// Should not have crashed
	if m == nil {
		t.Error("Expected model to be returned")
	}
}

func TestBackNavigationFromAllStates(t *testing.T) {
	tests := []struct {
		from     InstallerState
		expected InstallerState
	}{
		{StateWelcome, StateWelcome},           // Should quit instead, but stays in Welcome
		{StateToolSelection, StateWelcome},
		{StatePathSelection, StateToolSelection},
		{StateFileSelection, StatePathSelection},
		{StateError, StateWelcome},
	}
	
	for _, tt := range tests {
		t.Run(tt.from.String()+"_back", func(t *testing.T) {
			model := NewInstallerModel(&testAssetsEdge, &testAssetsEdge, &testAssetsEdge, &testAssetsEdge)
			model.state = tt.from
			
			// Simulate ESC key
			escMsg := tea.KeyMsg{Type: tea.KeyEsc}
			updatedModel, _ := model.Update(escMsg)
			m := updatedModel.(*InstallerModel)
			
			if m.state != tt.expected {
				t.Errorf("Back from %v: expected %v, got %v", tt.from, tt.expected, m.state)
			}
		})
	}
}

func TestWindowResizeHandling(t *testing.T) {
	model := NewInstallerModel(&testAssetsEdge, &testAssetsEdge, &testAssetsEdge, &testAssetsEdge)
	
	// Simulate window resize
	resizeMsg := tea.WindowSizeMsg{Width: 120, Height: 30}
	updatedModel, _ := model.Update(resizeMsg)
	m := updatedModel.(*InstallerModel)
	
	if m.width != 120 {
		t.Errorf("Expected width to be 120, got %d", m.width)
	}
	
	if m.height != 30 {
		t.Errorf("Expected height to be 30, got %d", m.height)
	}
}

func TestTreeSelectorWithNoFiles(t *testing.T) {
	// Create empty tree
	root := &TreeNode{
		Name:     "Empty",
		IsDir:    true,
		Children: []*TreeNode{},
	}
	
	selector := NewTreeSelector("Test", root)
	
	// Should not crash with empty tree
	if selector == nil {
		t.Error("Expected TreeSelector to be created")
	}
	
	view := selector.View()
	if view == "" {
		t.Error("Expected non-empty view for empty tree")
	}
}

func TestFileSelectionFlow(t *testing.T) {
	model := NewInstallerModel(&testAssetsEdge, &testAssetsEdge, &testAssetsEdge, &testAssetsEdge)
	
	// Go through proper state sequence: Welcome -> ToolSelection -> PathSelection -> FileSelection
	model.transitionToState(StateToolSelection)
	model.transitionToState(StatePathSelection)  
	model.transitionToState(StateFileSelection)
	
	// Files should be automatically selected (though may be 0 in test environment)
	allFiles := model.getAllAvailableFiles()
	if len(allFiles) > 0 && len(model.selectedFiles) == 0 {
		t.Error("Expected files to be selected in FileSelection state when files are available")
	}
	
	// Should be able to proceed to confirmation
	model.cursor = 0 // "Continue" option
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, _ := model.Update(enterMsg)
	m := updatedModel.(*InstallerModel)
	
	// Should transition to huh confirmation first
	if m.state != StateHuhConfirmation {
		t.Errorf("Expected state to transition to StateHuhConfirmation, got %v", m.state)
	}
}

func TestMalformedUserInput(t *testing.T) {
	model := NewInstallerModel(&testAssetsEdge, &testAssetsEdge, &testAssetsEdge, &testAssetsEdge)
	
	// Test with malformed key message
	malformedMsg := tea.KeyMsg{Type: tea.KeyType(999), Runes: []rune{0x00, 0xFF}}
	
	// Should not crash
	updatedModel, _ := model.Update(malformedMsg)
	if updatedModel == nil {
		t.Error("Expected model to handle malformed input gracefully")
	}
}

func TestProgressiveDisclosureRendererEdgeCases(t *testing.T) {
	renderer := NewProgressiveDisclosureRenderer()
	
	// Test with empty parameters
	result := renderer.RenderSelections("", "", 0)
	if result != "" {
		t.Error("Expected empty result for empty parameters")
	}
	
	// Test with partial parameters
	result = renderer.RenderSelections("claude-code", "", 0)
	if !containsString(result, "claude-code") {
		t.Error("Expected result to contain tool name")
	}
	
	// Test error rendering
	testErr := &testError{"test error"}
	errorResult := renderer.RenderError(testErr, "test context")
	if !containsString(errorResult, "test error") {
		t.Error("Expected error result to contain error message")
	}
}

func TestInstallationSummaryRendering(t *testing.T) {
	renderer := NewProgressiveDisclosureRenderer()
	
	// Test file mode summary
	files := []string{"agents/the-architect.md", "commands/develop.md"}
	summary := renderer.RenderInstallationSummary(
		"claude-code",
		"~/.config/the-startup",
		files,
		nil,
	)
	
	if !containsString(summary, "claude-code") {
		t.Error("Expected summary to contain tool name")
	}
	
	if !containsString(summary, "2") {
		t.Error("Expected summary to show file count")
	}
	
	// Test with updating files
	updatingFiles := []string{"agents/the-architect.md"}
	summary = renderer.RenderInstallationSummary(
		"claude-code",
		"~/.config/the-startup",
		files,
		updatingFiles,
	)
	
	if !containsString(summary, "⚠") {
		t.Error("Expected summary to show warning for updating files")
	}
}

func TestTreeSelectorKeyboardNavigation(t *testing.T) {
	root := &TreeNode{
		Name:     "Root",
		IsDir:    true,
		Children: []*TreeNode{
			{Name: "file1.txt", IsDir: false, Selected: true},
			{Name: "file2.txt", IsDir: false, Selected: false},
		},
	}
	
	selector := NewTreeSelector("Test", root)
	
	// Test cursor movement
	if selector.cursor != 0 {
		t.Error("Expected initial cursor position to be 0")
	}
	
	// Move down
	downMsg := tea.KeyMsg{Type: tea.KeyDown}
	selector.Update(downMsg)
	
	if selector.cursor != 1 {
		t.Errorf("Expected cursor to be 1 after down, got %d", selector.cursor)
	}
	
	// Move up
	upMsg := tea.KeyMsg{Type: tea.KeyUp}
	selector.Update(upMsg)
	
	if selector.cursor != 0 {
		t.Errorf("Expected cursor to be 0 after up, got %d", selector.cursor)
	}
	
	// Test space toggle
	spaceMsg := tea.KeyMsg{Type: tea.KeySpace}
	selector.Update(spaceMsg)
	
	// Should have toggled selection of first file
	if selector.nodes[0].Selected {
		t.Error("Expected first file to be deselected after space")
	}
}

func TestTreeSelectorHelpOverlay(t *testing.T) {
	root := &TreeNode{
		Name:     "Root",
		IsDir:    true,
		Children: []*TreeNode{},
	}
	
	selector := NewTreeSelector("Test", root)
	
	// Test help toggle
	helpMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}
	selector.Update(helpMsg)
	
	if !selector.showHelp {
		t.Error("Expected help to be shown after '?' key")
	}
	
	// Test help view
	helpView := selector.View()
	if !containsString(helpView, "Keyboard Shortcuts") {
		t.Error("Expected help view to contain shortcuts")
	}
	
	// Test closing help
	escMsg := tea.KeyMsg{Type: tea.KeyEsc}
	selector.Update(escMsg)
	
	if selector.showHelp {
		t.Error("Expected help to be hidden after ESC")
	}
}

func TestInstallProgressHandling(t *testing.T) {
	model := NewInstallerModel(&testAssetsEdge, &testAssetsEdge, &testAssetsEdge, &testAssetsEdge)
	
	// Test progress message
	progressMsg := InstallProgressMsg{
		Stage:   "Installing",
		Message: "Installing agents...",
		Percent: 50.0,
	}
	
	updatedModel, _ := model.Update(progressMsg)
	m := updatedModel.(*InstallerModel)
	
	if m.installProgress.Stage != "Installing" {
		t.Error("Expected progress stage to be set")
	}
	
	// Test completion message (success)
	completeMsg := InstallCompleteMsg{
		Success: true,
		Error:   nil,
	}
	
	updatedModel, _ = model.Update(completeMsg)
	m = updatedModel.(*InstallerModel)
	
	if m.state != StateComplete {
		t.Errorf("Expected state to be Complete after successful installation, got %v", m.state)
	}
	
	// Test completion message (error)
	errorCompleteMsg := InstallCompleteMsg{
		Success: false,
		Error:   &testError{"installation failed"},
	}
	
	updatedModel, _ = model.Update(errorCompleteMsg)
	m = updatedModel.(*InstallerModel)
	
	if m.state != StateError {
		t.Errorf("Expected state to be Error after failed installation, got %v", m.state)
	}
}

func TestGetAllAvailableFiles(t *testing.T) {
	model := NewInstallerModel(&testAssetsEdge, &testAssetsEdge, &testAssetsEdge, &testAssetsEdge)
	
	// Get all available files for advanced mode
	allFiles := model.getAllAvailableFiles()
	
	// In test environment, there might not be actual files in embed.FS
	// This is expected, so we just test that the method doesn't panic and returns a slice
	if allFiles == nil {
		t.Error("Expected getAllAvailableFiles to return a non-nil slice")
	}
	
	t.Logf("Found %d files in test environment", len(allFiles))
	
	// If files are found, verify they have correct prefixes
	if len(allFiles) > 0 {
		// Check that files from different components are included
		validPrefixes := []string{"agents/", "commands/", "hooks/", "templates/"}
		
		for _, file := range allFiles {
			hasValidPrefix := false
			for _, prefix := range validPrefixes {
				if strings.HasPrefix(file, prefix) {
					hasValidPrefix = true
					break
				}
			}
			if !hasValidPrefix {
				t.Errorf("File %s does not have a valid component prefix", file)
			}
		}
	}
}

func TestCustomPathSelection(t *testing.T) {
	model := NewInstallerModel(&testAssetsEdge, &testAssetsEdge, &testAssetsEdge, &testAssetsEdge)
	model.transitionToState(StatePathSelection)
	
	// Select "Custom location" option (should be index 2)
	model.cursor = 2
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, _ := model.Update(enterMsg)
	m := updatedModel.(*InstallerModel)
	
	// Should transition to error state (not implemented)
	if m.state != StateError {
		t.Errorf("Expected state to be Error for custom path, got %v", m.state)
	}
	
	if m.err == nil {
		t.Error("Expected error to be set for custom path selection")
	}
}

// Test error type for testing
type testError struct {
	message string
}

func (e *testError) Error() string {
	return e.message
}