package ui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// containsString checks if a string contains a substring
func containsString(s, substr string) bool {
	return strings.Contains(s, substr)
}

func TestNewMainModel(t *testing.T) {
	// Test that NewMainModel creates a valid model
	model := NewMainModel(&testAssets, &testAssets, &testAssets, &testAssets)
	
	if model == nil {
		t.Fatal("NewMainModel returned nil")
	}
	
	// Should start directly with ToolSelection state
	if model.state != StateToolSelection {
		t.Errorf("Expected initial state to be StateToolSelection, got %v", model.state)
	}
	
	if model.installer == nil {
		t.Error("Expected installer to be initialized")
	}
	
	// Check that sub-models are initialized
	if model.toolSelectionModel.renderer == nil {
		t.Error("Expected tool selection model renderer to be initialized")
	}
}

func TestMainModelInit(t *testing.T) {
	model := NewMainModel(&testAssets, &testAssets, &testAssets, &testAssets)
	
	// Init should return nil for consolidated model
	cmd := model.Init()
	if cmd != nil {
		t.Error("Expected Init() to return nil for consolidated model")
	}
}

func TestMainModelDirectStates(t *testing.T) {
	model := NewMainModel(&testAssets, &testAssets, &testAssets, &testAssets)
	
	// Test that we start directly with ToolSelection
	if model.state != StateToolSelection {
		t.Errorf("Expected StateToolSelection, got %v", model.state)
	}
	
	// Test view renders
	view := model.View()
	if len(view) == 0 {
		t.Error("Expected tool selection view to have content")
	}
	
	// Test state transition
	model.transitionToState(StatePathSelection)
	if model.state != StatePathSelection {
		t.Errorf("Expected StatePathSelection, got %v", model.state)
	}
	
	view = model.View()
	if len(view) == 0 {
		t.Error("Expected path selection view to have content")
	}
}

func TestMainModelKeyHandling(t *testing.T) {
	model := NewMainModel(&testAssets, &testAssets, &testAssets, &testAssets)
	
	// Test ESC key - should quit from tool selection
	escMsg := tea.KeyMsg{Type: tea.KeyEsc}
	_, cmd := model.Update(escMsg)
	
	// Should return tea.Quit command when going back from tool selection to welcome
	if cmd == nil {
		t.Error("Expected ESC to return a quit command from initial state")
	}
}

func TestMainModelFileSelectionIntegration(t *testing.T) {
	model := NewMainModel(&testAssets, &testAssets, &testAssets, &testAssets)
	
	// Navigate through proper flow to file selection
	model.selectedTool = "claude-code"
	model.transitionToState(StateFileSelection)
	
	// Should have files auto-selected after proper path selection
	if len(model.selectedFiles) == 0 {
		// Debug: check what files are available
		allFiles := model.getAllAvailableFiles()
		t.Logf("Available files: %v", allFiles)
		t.Error("Expected files to be auto-selected when entering file selection state properly")
	}
	
	// Should have confirmation choices initialized in file selection model
	if len(model.fileSelectionModel.choices) != 2 {
		t.Error("Expected confirmation choices to be initialized in file selection state")
	}
	
	// Test view renders
	view := model.View()
	if len(view) == 0 {
		t.Error("Expected file selection view to have content")
	}
	
	// Verify the view contains the tree structure
	if !strings.Contains(view, "agents") {
		t.Error("Expected view to contain 'agents' folder in tree")
	}
	if !strings.Contains(view, "commands") {
		t.Error("Expected view to contain 'commands' folder in tree")
	}
	if !strings.Contains(view, "hooks") {
		t.Error("Expected view to contain 'hooks' folder in tree")
	}
	if !strings.Contains(view, "templates") {
		t.Error("Expected view to contain 'templates' folder in tree")
	}
	
	// Verify confirmation options are shown
	if !strings.Contains(view, "Yes, give me awesome") {
		t.Error("Expected view to contain confirmation option 'Yes, give me awesome'")
	}
	if !strings.Contains(view, "Huh? I did not sign up for this") {
		t.Error("Expected view to contain confirmation option 'Huh? I did not sign up for this'")
	}
	
	// Verify the ready to install prompt
	if !strings.Contains(view, "Ready to install?") {
		t.Error("Expected view to contain 'Ready to install?' prompt")
	}
}

func TestMainModelHuhIntegration(t *testing.T) {
	model := NewMainModel(&testAssets, &testAssets, &testAssets, &testAssets)
	
	// Setup file selection state
	model.selectedTool = "claude-code"
	model.transitionToState(StateFileSelection)
	
	// Test view renders with huh form always visible
	view := model.View()
	if len(view) == 0 {
		t.Error("Expected file selection view with huh form to have content")
	}
	
	// Test that state is FileSelection with huh form always shown
	if model.state != StateFileSelection {
		t.Errorf("Expected state to be StateFileSelection, got %v", model.state)
	}
	
	// Confirmation choices should be initialized in file selection model
	if len(model.fileSelectionModel.choices) != 2 {
		t.Error("Expected confirmation choices to be initialized")
	}
}

func TestMainModelChoicesInitialization(t *testing.T) {
	model := NewMainModel(&testAssets, &testAssets, &testAssets, &testAssets)
	
	// Should start with tool selection choices
	expectedChoices := []string{"claude-code", "Cancel"}
	if len(model.toolSelectionModel.choices) != len(expectedChoices) {
		t.Errorf("Expected %d choices for tool selection, got %d", len(expectedChoices), len(model.toolSelectionModel.choices))
	}
	
	// Test path selection choices
	model.transitionToState(StatePathSelection)
	if len(model.pathSelectionModel.choices) < 3 { // Should have at least recommended, local, custom, cancel
		t.Errorf("Expected at least 3 choices for path selection, got %d", len(model.pathSelectionModel.choices))
	}
	
	// Test file selection choices - should have confirmation options
	model.transitionToState(StateFileSelection)
	expectedFileChoices := []string{
		"Yes, give me awesome",
		"Huh? I did not sign up for this",
	}
	if len(model.fileSelectionModel.choices) != len(expectedFileChoices) {
		t.Errorf("Expected %d choices for file selection, got %d", len(expectedFileChoices), len(model.fileSelectionModel.choices))
	}
}