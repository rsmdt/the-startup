package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// TestMainModelIntegration tests the complete flow through the MainModel
func TestMainModelIntegration(t *testing.T) {
	// Create MainModel
	model := NewMainModel(&testAssets, &testAssets, &testAssets, &testAssets)
	
	// Should start in Welcome state
	if model.state != StateWelcome {
		t.Fatalf("Expected initial state to be StateWelcome, got %v", model.state)
	}
	
	// Simulate complete flow: Welcome -> Tool Selection -> Path Selection -> File Selection -> Confirmation
	
	// 1. Transition from Welcome to Tool Selection
	transitionMsg := ViewTransitionMsg{NextView: StateToolSelection}
	updatedModel, _ := model.Update(transitionMsg)
	model = updatedModel.(*MainModel)
	
	if model.state != StateToolSelection {
		t.Fatalf("Expected state to be StateToolSelection, got %v", model.state)
	}
	
	// 2. Make tool selection
	selectionMsg := SelectionMadeMsg{Tool: "claude-code"}
	updatedModel, _ = model.Update(selectionMsg)
	model = updatedModel.(*MainModel)
	
	if model.context.SelectedTool != "claude-code" {
		t.Fatalf("Expected SelectedTool to be 'claude-code', got '%s'", model.context.SelectedTool)
	}
	
	// 3. Transition to Path Selection
	transitionMsg = ViewTransitionMsg{NextView: StatePathSelection}
	updatedModel, _ = model.Update(transitionMsg)
	model = updatedModel.(*MainModel)
	
	if model.state != StatePathSelection {
		t.Fatalf("Expected state to be StatePathSelection, got %v", model.state)
	}
	
	// 4. Make path and file selection
	selectionMsg = SelectionMadeMsg{Path: "/test/path", Files: []string{"test1.md", "test2.md"}}
	updatedModel, _ = model.Update(selectionMsg)
	model = updatedModel.(*MainModel)
	
	if model.context.SelectedPath != "/test/path" {
		t.Fatalf("Expected SelectedPath to be '/test/path', got '%s'", model.context.SelectedPath)
	}
	
	if len(model.context.SelectedFiles) != 2 {
		t.Fatalf("Expected 2 selected files, got %d", len(model.context.SelectedFiles))
	}
	
	// 5. Transition to File Selection
	transitionMsg = ViewTransitionMsg{NextView: StateFileSelection}
	updatedModel, _ = model.Update(transitionMsg)
	model = updatedModel.(*MainModel)
	
	if model.state != StateFileSelection {
		t.Fatalf("Expected state to be StateFileSelection, got %v", model.state)
	}
	
	// 6. Transition to HuhConfirmation first
	transitionMsg = ViewTransitionMsg{NextView: StateHuhConfirmation}
	updatedModel, _ = model.Update(transitionMsg)
	model = updatedModel.(*MainModel)
	
	if model.state != StateHuhConfirmation {
		t.Fatalf("Expected state to be StateHuhConfirmation, got %v", model.state)
	}
	
	// 7. Then transition to Confirmation
	transitionMsg = ViewTransitionMsg{NextView: StateConfirmation}
	updatedModel, _ = model.Update(transitionMsg)
	model = updatedModel.(*MainModel)
	
	if model.state != StateConfirmation {
		t.Fatalf("Expected state to be StateConfirmation, got %v", model.state)
	}
	
	// 8. Test View delegation works at each stage
	view := model.View()
	if view == "" {
		t.Error("Expected non-empty view from ConfirmationModel")
	}
	
	// All state and context should be preserved throughout
	if model.context.SelectedTool != "claude-code" {
		t.Error("SelectedTool was lost during transitions")
	}
	
	if model.context.SelectedPath != "/test/path" {
		t.Error("SelectedPath was lost during transitions")
	}
	
	if len(model.context.SelectedFiles) != 2 {
		t.Error("SelectedFiles were lost during transitions")
	}
}

// TestMainModelIntegrationErrorHandling tests error scenarios
func TestMainModelIntegrationErrorHandling(t *testing.T) {
	model := NewMainModel(&testAssets, &testAssets, &testAssets, &testAssets)
	
	// Test invalid transition
	invalidMsg := ViewTransitionMsg{NextView: StateComplete} // Invalid from Welcome
	updatedModel, _ := model.Update(invalidMsg)
	model = updatedModel.(*MainModel)
	
	if model.state != StateError {
		t.Fatalf("Expected state to be StateError for invalid transition, got %v", model.state)
	}
	
	// Test error recovery
	recoveryMsg := ViewTransitionMsg{NextView: StateWelcome}
	updatedModel, _ = model.Update(recoveryMsg)
	model = updatedModel.(*MainModel)
	
	if model.state != StateWelcome {
		t.Fatalf("Expected state to recover to StateWelcome, got %v", model.state)
	}
}

// TestMainModelViewDelegation tests that views are properly delegated
func TestMainModelViewDelegation(t *testing.T) {
	model := NewMainModel(&testAssets, &testAssets, &testAssets, &testAssets)
	
	// Test Welcome view
	view := model.View()
	if !containsString(view, "Startup") {
		t.Error("Welcome view should contain 'Startup'")
	}
	
	// Transition to Tool Selection
	transitionMsg := ViewTransitionMsg{NextView: StateToolSelection}
	updatedModel, _ := model.Update(transitionMsg)
	model = updatedModel.(*MainModel)
	
	// Test Tool Selection view
	view = model.View()
	if !containsString(view, "development tool") {
		t.Error("Tool selection view should contain selection prompt")
	}
}

// TestMainModelWindowResize tests window resize handling
func TestMainModelWindowResize(t *testing.T) {
	model := NewMainModel(&testAssets, &testAssets, &testAssets, &testAssets)
	
	// Test window resize
	resizeMsg := tea.WindowSizeMsg{Width: 120, Height: 40}
	updatedModel, _ := model.Update(resizeMsg)
	model = updatedModel.(*MainModel)
	
	if model.context.Width != 120 {
		t.Errorf("Expected width to be 120, got %d", model.context.Width)
	}
	
	if model.context.Height != 40 {
		t.Errorf("Expected height to be 40, got %d", model.context.Height)
	}
}