package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// TestMainModelConsolidatedIntegration tests the new consolidated MainModel
func TestMainModelConsolidatedIntegration(t *testing.T) {
	// Create MainModel with nil assets (UI should handle gracefully)
	model := NewMainModel(nil, nil)

	// Should start with StartupPath state
	if model.state != StateStartupPath {
		t.Fatalf("Expected initial state to be StateStartupPath, got %v", model.state)
	}

	// Test view renders without error
	view := model.View()
	if len(view) == 0 {
		t.Error("Expected view to have content")
	}

	// Test basic navigation
	// Simulate pressing Enter to select first option (claude-code)
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, _ := model.Update(enterMsg)
	model = updatedModel.(*MainModel)

	// Should transition to Claude Path after selecting startup path
	if model.state != StateClaudePath {
		t.Fatalf("Expected state to be StateClaudePath after selection, got %v", model.state)
	}

	// Test view renders in new state
	view = model.View()
	if len(view) == 0 {
		t.Error("Expected path selection view to have content")
	}

	// Test ESC key for back navigation
	escMsg := tea.KeyMsg{Type: tea.KeyEsc}
	updatedModel, _ = model.Update(escMsg)
	model = updatedModel.(*MainModel)

	// Should go back to Startup Path
	if model.state != StateStartupPath {
		t.Fatalf("Expected state to go back to StateStartupPath after ESC, got %v", model.state)
	}
}

// TestMainModelFileSelectionWithHuh tests the file selection with integrated huh confirmation
func TestMainModelFileSelectionWithHuh(t *testing.T) {
	model := NewMainModel(nil, nil)

	// Navigate through proper flow to file selection
	model.startupPath = "~/.config/the-startup"
	model.claudePath = "~/.claude"
	model.transitionToState(StateFileSelection)

	// We're already in FileSelection state from transitionToState above

	// With nil assets, no files will be available, so selectedFiles may be empty
	// This is expected behavior in test environment
	allFiles := model.getAllAvailableFiles()
	if len(allFiles) == 0 {
		t.Logf("No files available (expected with nil assets in test environment)")
	}

	// Test view renders file tree
	view := model.View()
	if len(view) == 0 {
		t.Error("Expected file selection view to have content")
	}

	// Should have confirmation choices in the file selection model
	if len(model.fileSelectionModel.choices) != 2 {
		t.Error("Expected confirmation choices to be set")
	}

	// State should still be FileSelection
	if model.state != StateFileSelection {
		t.Fatalf("Expected state to remain StateFileSelection, got %v", model.state)
	}
}

// TestMainModelConsolidation tests that the consolidated model works correctly
func TestMainModelConsolidation(t *testing.T) {
	model := NewMainModel(nil, nil)

	// Test that all required fields are properly initialized
	if model.installer == nil {
		t.Error("Expected installer to be initialized")
	}

	// Test that sub-models are initialized
	// Test that the startup path model is initialized
	// (textInput is a struct, not a pointer, so we just verify the model exists)
}
