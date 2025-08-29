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
	model := NewMainModel(nil, nil)

	if model == nil {
		t.Fatal("NewMainModel returned nil")
	}

	// Should start with StartupPath state
	if model.state != StateStartupPath {
		t.Errorf("Expected initial state to be StateStartupPath, got %v", model.state)
	}

	if model.installer == nil {
		t.Error("Expected installer to be initialized")
	}

	// Check that sub-models are initialized
	// (textInput is a struct, not a pointer, so we just verify the model exists)
}

func TestMainModelInit(t *testing.T) {
	model := NewMainModel(nil, nil)

	// Init should return nil for consolidated model
	cmd := model.Init()
	if cmd != nil {
		t.Error("Expected Init() to return nil for consolidated model")
	}
}

func TestMainModelDirectStates(t *testing.T) {
	model := NewMainModel(nil, nil)

	// Test that we start with StartupPath
	if model.state != StateStartupPath {
		t.Errorf("Expected StateStartupPath, got %v", model.state)
	}

	// Test view renders
	view := model.View()
	if len(view) == 0 {
		t.Error("Expected tool selection view to have content")
	}

	// Test state transition
	model.transitionToState(StateClaudePath)
	if model.state != StateClaudePath {
		t.Errorf("Expected StateClaudePath, got %v", model.state)
	}

	view = model.View()
	if len(view) == 0 {
		t.Error("Expected path selection view to have content")
	}
}

func TestMainModelKeyHandling(t *testing.T) {
	model := NewMainModel(nil, nil)

	// Test ESC key - should quit from tool selection
	escMsg := tea.KeyMsg{Type: tea.KeyEsc}
	_, cmd := model.Update(escMsg)

	// Should return tea.Quit command when going back from tool selection to welcome
	if cmd == nil {
		t.Error("Expected ESC to return a quit command from initial state")
	}
}

func TestMainModelFileSelectionIntegration(t *testing.T) {
	model := NewMainModel(nil, nil)

	// Navigate through proper flow to file selection
	model.startupPath = "~/.config/the-startup"
	model.claudePath = "~/.claude"
	model.transitionToState(StateFileSelection)

	// With nil assets, no files will be auto-selected (expected in test environment)
	allFiles := model.getAllAvailableFiles()
	if len(allFiles) == 0 {
		t.Logf("No files available (expected with nil assets in test environment)")
	} else {
		// If files are available, they should be auto-selected
		if len(model.selectedFiles) == 0 {
			t.Error("Expected files to be auto-selected when files are available")
		}
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

	// Verify the view contains the tree structure (if assets are available)
	// The test uses nil assets, so tree structure may not be present
	hasAgents := strings.Contains(view, "agents") || strings.Contains(view, "agents/")
	hasCommands := strings.Contains(view, "commands") || strings.Contains(view, "commands/")
	
	if !hasAgents && !strings.Contains(view, "Files to be installed") {
		t.Log("Note: View doesn't contain 'agents' folder, expected with nil assets")
	}
	if !hasCommands && !strings.Contains(view, "Files to be installed") {
		t.Log("Note: View doesn't contain 'commands' folder, expected with nil assets")
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
	model := NewMainModel(nil, nil)

	// Setup file selection state
	model.startupPath = "~/.config/the-startup"
	model.claudePath = "~/.claude"
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
	model := NewMainModel(nil, nil)

	// Test that startup path model is initialized
	// (textInput is a struct, not a pointer, so we just verify the model exists)

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
