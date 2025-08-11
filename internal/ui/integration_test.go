package ui

import (
	"embed"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

//go:embed test_assets/assets/*
var testAssets embed.FS

// TestMainModelConsolidatedIntegration tests the new consolidated MainModel
func TestMainModelConsolidatedIntegration(t *testing.T) {
	// Create MainModel
	model := NewMainModel(&testAssets, &testAssets, &testAssets, &testAssets)
	
	// Should start directly with ToolSelection state
	if model.state != StateToolSelection {
		t.Fatalf("Expected initial state to be StateToolSelection, got %v", model.state)
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
	
	// Should transition to Path Selection and have tool selected
	if model.state != StatePathSelection {
		t.Fatalf("Expected state to be StatePathSelection after selection, got %v", model.state)
	}
	
	if model.selectedTool != "claude-code" {
		t.Fatalf("Expected SelectedTool to be 'claude-code', got '%s'", model.selectedTool)
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
	
	// Should go back to Tool Selection
	if model.state != StateToolSelection {
		t.Fatalf("Expected state to go back to StateToolSelection after ESC, got %v", model.state)
	}
}

// TestMainModelFileSelectionWithHuh tests the file selection with integrated huh confirmation
func TestMainModelFileSelectionWithHuh(t *testing.T) {
	model := NewMainModel(&testAssets, &testAssets, &testAssets, &testAssets)
	
	// Navigate through proper flow to file selection
	model.selectedTool = "claude-code"
	model.transitionToState(StatePathSelection)
	
	// Select recommended path (first option) - this will auto-select files
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, _ := model.Update(enterMsg)
	model = updatedModel.(*MainModel)
	
	// Should transition to File Selection
	if model.state != StateFileSelection {
		t.Fatalf("Expected state to be StateFileSelection, got %v", model.state)
	}
	
	// Should have files selected automatically
	if len(model.selectedFiles) == 0 {
		// Debug: check what files are available
		allFiles := model.getAllAvailableFiles()
		t.Errorf("Expected files to be selected automatically. Available files: %v", allFiles)
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
	model := NewMainModel(&testAssets, &testAssets, &testAssets, &testAssets)
	
	// Test that all required fields are properly initialized
	if model.installer == nil {
		t.Error("Expected installer to be initialized")
	}
	
	// Test that sub-models are initialized
	if model.toolSelectionModel.renderer == nil {
		t.Error("Expected tool selection model renderer to be initialized")
	}
	
	// Test that choices are set up correctly for tool selection
	if len(model.toolSelectionModel.choices) == 0 {
		t.Error("Expected choices to be initialized for tool selection")
	}
	
	expectedChoices := []string{"claude-code", "Cancel"}
	if len(model.toolSelectionModel.choices) != len(expectedChoices) {
		t.Errorf("Expected %d choices, got %d", len(expectedChoices), len(model.toolSelectionModel.choices))
	}
}