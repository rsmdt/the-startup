package ui

import (
	"embed"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

//go:embed test_assets/*
var testAssets embed.FS

func TestNewInstallerModel(t *testing.T) {
	// Test that NewInstallerModel creates a valid model
	model := NewInstallerModel(&testAssets, &testAssets, &testAssets, &testAssets)
	
	if model == nil {
		t.Fatal("NewInstallerModel returned nil")
	}
	
	if model.state != StateWelcome {
		t.Errorf("Expected initial state to be StateWelcome, got %v", model.state)
	}
	
	if model.installer == nil {
		t.Error("Expected installer to be initialized")
	}
}

func TestStateTransitions(t *testing.T) {
	tests := []struct {
		name  string
		from  InstallerState
		to    InstallerState
		valid bool
	}{
		{"Welcome to Tool Selection", StateWelcome, StateToolSelection, true},
		{"Tool Selection to Path Selection", StateToolSelection, StatePathSelection, true},
		{"Path to File Selection", StatePathSelection, StateFileSelection, true},
		{"Invalid transition", StateWelcome, StateComplete, false},
		{"Back navigation", StateToolSelection, StateWelcome, true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := IsValidTransition(tt.from, tt.to)
			if valid != tt.valid {
				t.Errorf("IsValidTransition(%v, %v) = %v, want %v", 
					tt.from, tt.to, valid, tt.valid)
			}
		})
	}
}

func TestProgressiveDisclosureRenderer(t *testing.T) {
	renderer := NewProgressiveDisclosureRenderer()
	
	// Test empty selections
	result := renderer.RenderSelections("", "", 0)
	if result != "" {
		t.Error("Expected empty string for empty selections")
	}
	
	// Test with selections
	result = renderer.RenderSelections("claude-code", "~/.config/the-startup", 5)
	if result == "" {
		t.Error("Expected non-empty string for valid selections")
	}
	
	if !containsString(result, "claude-code") {
		t.Error("Expected result to contain tool selection")
	}
}

func TestInstallerModelInit(t *testing.T) {
	model := NewInstallerModel(&testAssets, &testAssets, &testAssets, &testAssets)
	
	// Since model starts in StateWelcome, Init() should return a tick command for auto-advance
	cmd := model.Init()
	if cmd == nil {
		t.Error("Expected Init() to return a tick command for auto-advance from welcome screen")
	}
	
	// After transitioning to another state, Init() should return nil
	model.transitionToState(StateToolSelection)
	cmd = model.Init()
	if cmd != nil {
		t.Error("Expected Init() to return nil command for non-welcome states")
	}
}

func TestInstallerModelView(t *testing.T) {
	model := NewInstallerModel(&testAssets, &testAssets, &testAssets, &testAssets)
	
	// Test welcome view
	view := model.View()
	if view == "" {
		t.Error("Expected non-empty view")
	}
	
	if !containsString(view, "Startup") {
		maxLen := len(view)
		if maxLen > 100 {
			maxLen = 100
		}
		t.Errorf("Expected welcome view to contain 'Startup', got: %q", view[:maxLen])
	}
}

func TestModelKeyHandling(t *testing.T) {
	model := NewInstallerModel(&testAssets, &testAssets, &testAssets, &testAssets)
	
	// Test enter key on welcome screen
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, cmd := model.Update(enterMsg)
	
	if updatedModel == nil {
		t.Error("Expected model to be returned")
	}
	
	m := updatedModel.(*InstallerModel)
	if m.state != StateToolSelection {
		t.Errorf("Expected state to transition to StateToolSelection, got %v", m.state)
	}
	
	// cmd can be nil, that's fine
	_ = cmd
}

func TestWelcomeBanner(t *testing.T) {
	if WelcomeBanner == "" {
		t.Error("WelcomeBanner should not be empty")
	}
	
	if !containsString(WelcomeBanner, "████") {
		t.Error("Expected banner to contain ASCII art")
	}
}

// Helper function to check if string contains substring
func containsString(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && 
		   len(s) >= len(substr) && 
		   findSubstring(s, substr)
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			if s[i+j] != substr[j] {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}