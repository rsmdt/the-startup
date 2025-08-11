package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewMainModel(t *testing.T) {
	// Test that NewMainModel creates a valid model
	model := NewMainModel(&testAssets, &testAssets, &testAssets, &testAssets)
	
	if model == nil {
		t.Fatal("NewMainModel returned nil")
	}
	
	if model.state != StateWelcome {
		t.Errorf("Expected initial state to be StateWelcome, got %v", model.state)
	}
	
	if model.context == nil {
		t.Error("Expected context to be initialized")
	}
	
	if model.context.Installer == nil {
		t.Error("Expected installer to be initialized")
	}
	
	if model.currentView == nil {
		t.Error("Expected current view to be initialized")
	}
}

func TestMainModelInit(t *testing.T) {
	model := NewMainModel(&testAssets, &testAssets, &testAssets, &testAssets)
	
	// Init should delegate to current view (WelcomeModel)
	cmd := model.Init()
	if cmd == nil {
		t.Error("Expected Init() to return a command from WelcomeModel")
	}
}

func TestMainModelViewTransitions(t *testing.T) {
	model := NewMainModel(&testAssets, &testAssets, &testAssets, &testAssets)
	
	// Test transition from Welcome to Tool Selection
	transitionMsg := ViewTransitionMsg{NextView: StateToolSelection}
	updatedModel, cmd := model.Update(transitionMsg)
	
	if updatedModel == nil {
		t.Fatal("Expected model to be returned")
	}
	
	mainModel := updatedModel.(*MainModel)
	if mainModel.state != StateToolSelection {
		t.Errorf("Expected state to be StateToolSelection, got %v", mainModel.state)
	}
	
	// cmd should be from the new view's Init()
	_ = cmd // cmd can be nil, that's fine
}

func TestMainModelInvalidTransition(t *testing.T) {
	model := NewMainModel(&testAssets, &testAssets, &testAssets, &testAssets)
	
	// Try invalid transition from Welcome to Complete
	transitionMsg := ViewTransitionMsg{NextView: StateComplete}
	updatedModel, _ := model.Update(transitionMsg)
	
	mainModel := updatedModel.(*MainModel)
	if mainModel.state != StateError {
		t.Errorf("Expected state to transition to StateError for invalid transition, got %v", mainModel.state)
	}
}

func TestMainModelSelectionMadeMsg(t *testing.T) {
	model := NewMainModel(&testAssets, &testAssets, &testAssets, &testAssets)
	
	// Test tool selection
	selectionMsg := SelectionMadeMsg{Tool: "claude-code"}
	updatedModel, _ := model.Update(selectionMsg)
	
	mainModel := updatedModel.(*MainModel)
	if mainModel.context.SelectedTool != "claude-code" {
		t.Errorf("Expected SelectedTool to be 'claude-code', got '%s'", mainModel.context.SelectedTool)
	}
}

func TestMainModelWindowSizeMsg(t *testing.T) {
	model := NewMainModel(&testAssets, &testAssets, &testAssets, &testAssets)
	
	// Test window resize
	sizeMsg := tea.WindowSizeMsg{Width: 120, Height: 40}
	updatedModel, _ := model.Update(sizeMsg)
	
	mainModel := updatedModel.(*MainModel)
	if mainModel.context.Width != 120 || mainModel.context.Height != 40 {
		t.Errorf("Expected dimensions to be 120x40, got %dx%d", 
			mainModel.context.Width, mainModel.context.Height)
	}
}

func TestMainModelView(t *testing.T) {
	model := NewMainModel(&testAssets, &testAssets, &testAssets, &testAssets)
	
	// Test that view delegates to current view
	view := model.View()
	if view == "" {
		t.Error("Expected non-empty view")
	}
	
	// Should contain welcome content since we start in StateWelcome
	if !containsString(view, "Startup") {
		t.Error("Expected view to contain welcome content")
	}
}

func TestMainModelErrorHandling(t *testing.T) {
	model := NewMainModel(&testAssets, &testAssets, &testAssets, &testAssets)
	
	// Test error transition with data
	errorData := map[string]interface{}{
		"error": &TestError{"test error"},
		"context": "test context",
	}
	
	transitionMsg := ViewTransitionMsg{NextView: StateError, Data: errorData}
	updatedModel, _ := model.Update(transitionMsg)
	
	mainModel := updatedModel.(*MainModel)
	if mainModel.state != StateError {
		t.Errorf("Expected state to be StateError, got %v", mainModel.state)
	}
	
	// Check that error view received the error data
	errorView, ok := mainModel.currentView.(*ErrorModel)
	if !ok {
		t.Fatal("Expected current view to be ErrorModel")
	}
	
	if errorView.err == nil {
		t.Error("Expected error to be set in ErrorModel")
	}
}

// TestError implements error interface for testing
type TestError struct {
	message string
}

func (e *TestError) Error() string {
	return e.message
}