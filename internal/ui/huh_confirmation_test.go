package ui

import (
	"embed"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

//go:embed test_assets/*
var testAssetsHuh embed.FS

// TestHuhConfirmationStateTransition tests the huh confirmation state transition
func TestHuhConfirmationStateTransition(t *testing.T) {
	model := NewInstallerModel(&testAssetsHuh, &testAssetsHuh, &testAssetsHuh, &testAssetsHuh)
	
	// Set up for huh confirmation state
	model.transitionToState(StateToolSelection)
	model.selectedTool = "claude-code"
	model.transitionToState(StatePathSelection)
	model.selectedPath = "~/.config/the-startup"
	model.transitionToState(StateFileSelection)
	model.selectedFiles = []string{"agents/test.md", "commands/test.md"}
	model.transitionToState(StateHuhConfirmation)
	
	// Verify state
	if model.state != StateHuhConfirmation {
		t.Errorf("Expected StateHuhConfirmation, got %v", model.state)
	}
	
	// Verify huh form is created
	if model.huhForm == nil {
		t.Error("Expected huh form to be created")
	}
	
	// Test view renders without error
	view := model.View()
	if len(view) == 0 {
		t.Error("Expected huh confirmation view to have content")
	}
	
	// Should contain the welcome banner
	if !containsString(view, "████████ ██   ██ ███████") {
		t.Error("Huh confirmation view should contain banner")
	}
}

// TestHuhConfirmationToConfirmationFlow tests confirming the huh dialog
func TestHuhConfirmationToConfirmationFlow(t *testing.T) {
	model := NewInstallerModel(&testAssetsHuh, &testAssetsHuh, &testAssetsHuh, &testAssetsHuh)
	
	// Set up proper state progression
	model.transitionToState(StateToolSelection)
	model.selectedTool = "claude-code"
	model.transitionToState(StatePathSelection)
	model.selectedPath = "~/.config/the-startup"
	model.transitionToState(StateFileSelection)
	model.selectedFiles = []string{"agents/test.md", "commands/test.md"}
	model.transitionToState(StateHuhConfirmation)
	
	// Simulate confirming (yes)
	confirmMsg := HuhConfirmCompleteMsg{Confirmed: true}
	updatedModel, _ := model.Update(confirmMsg)
	m := updatedModel.(*InstallerModel)
	
	// Should transition to regular confirmation
	if m.state != StateConfirmation {
		t.Errorf("Expected StateConfirmation after confirming, got %v", m.state)
	}
}

// TestHuhConfirmationToFileSelectionFlow tests declining the huh dialog
func TestHuhConfirmationToFileSelectionFlow(t *testing.T) {
	model := NewInstallerModel(&testAssetsHuh, &testAssetsHuh, &testAssetsHuh, &testAssetsHuh)
	
	// Set up proper state progression
	model.transitionToState(StateToolSelection)
	model.selectedTool = "claude-code"
	model.transitionToState(StatePathSelection)
	model.selectedPath = "~/.config/the-startup"
	model.transitionToState(StateFileSelection)
	model.selectedFiles = []string{"agents/test.md", "commands/test.md"}
	model.transitionToState(StateHuhConfirmation)
	
	// Simulate declining (no)
	declineMsg := HuhConfirmCompleteMsg{Confirmed: false}
	updatedModel, _ := model.Update(declineMsg)
	m := updatedModel.(*InstallerModel)
	
	// Should go back to file selection
	if m.state != StateFileSelection {
		t.Errorf("Expected StateFileSelection after declining, got %v", m.state)
	}
}

// TestHuhConfirmationBackNavigation tests ESC key from huh confirmation
func TestHuhConfirmationBackNavigation(t *testing.T) {
	model := NewInstallerModel(&testAssetsHuh, &testAssetsHuh, &testAssetsHuh, &testAssetsHuh)
	
	// Set up proper state progression
	model.transitionToState(StateToolSelection)
	model.selectedTool = "claude-code"
	model.transitionToState(StatePathSelection)
	model.selectedPath = "~/.config/the-startup"
	model.transitionToState(StateFileSelection)
	model.selectedFiles = []string{"agents/test.md", "commands/test.md"}
	model.transitionToState(StateHuhConfirmation)
	
	// Press escape
	escMsg := tea.KeyMsg{Type: tea.KeyEsc}
	updatedModel, _ := model.Update(escMsg)
	m := updatedModel.(*InstallerModel)
	
	// Should go back to file selection
	if m.state != StateFileSelection {
		t.Errorf("Expected StateFileSelection after ESC, got %v", m.state)
	}
}

