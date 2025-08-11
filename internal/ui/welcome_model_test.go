package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestWelcomeModel(t *testing.T) {
	context := &Context{
		Styles: GetStyles(),
	}
	model := NewWelcomeModel(context)
	
	if model == nil {
		t.Fatal("NewWelcomeModel returned nil")
	}
	
	if model.context != context {
		t.Error("Expected context to be set")
	}
}

func TestWelcomeModelInit(t *testing.T) {
	context := &Context{
		Styles: GetStyles(),
	}
	model := NewWelcomeModel(context)
	
	// Init should return a tick command for auto-advance
	cmd := model.Init()
	if cmd == nil {
		t.Error("Expected Init() to return a tick command for auto-advance")
	}
}

func TestWelcomeModelAutoAdvance(t *testing.T) {
	context := &Context{
		Styles: GetStyles(),
	}
	model := NewWelcomeModel(context)
	
	// Test auto-advance on tick
	tickMsg := AutoAdvanceMsg{}
	updatedModel, cmd := model.Update(tickMsg)
	
	if updatedModel == nil {
		t.Fatal("Expected model to be returned")
	}
	
	if cmd == nil {
		t.Error("Expected tick to generate ViewTransitionMsg command")
	}
	
	// Execute the command to get the message
	msg := cmd()
	if transitionMsg, ok := msg.(ViewTransitionMsg); ok {
		if transitionMsg.NextView != StateToolSelection {
			t.Errorf("Expected transition to StateToolSelection, got %v", transitionMsg.NextView)
		}
	} else {
		t.Error("Expected ViewTransitionMsg from tick")
	}
}

func TestWelcomeModelKeyHandling(t *testing.T) {
	context := &Context{
		Styles: GetStyles(),
	}
	model := NewWelcomeModel(context)
	
	// Test Enter key
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, cmd := model.Update(enterMsg)
	
	if updatedModel == nil {
		t.Fatal("Expected model to be returned")
	}
	
	if cmd == nil {
		t.Error("Expected Enter to generate ViewTransitionMsg command")
	}
	
	// Execute the command to get the message
	msg := cmd()
	if transitionMsg, ok := msg.(ViewTransitionMsg); ok {
		if transitionMsg.NextView != StateToolSelection {
			t.Errorf("Expected transition to StateToolSelection, got %v", transitionMsg.NextView)
		}
	} else {
		t.Error("Expected ViewTransitionMsg from Enter key")
	}
}

func TestWelcomeModelQuit(t *testing.T) {
	context := &Context{
		Styles: GetStyles(),
	}
	model := NewWelcomeModel(context)
	
	// Test Ctrl+C
	quitMsg := tea.KeyMsg{Type: tea.KeyCtrlC}
	updatedModel, cmd := model.Update(quitMsg)
	
	if updatedModel == nil {
		t.Fatal("Expected model to be returned")
	}
	
	if cmd == nil {
		t.Error("Expected Ctrl+C to generate quit command")
	}
	
	// The command should be tea.Quit - we can't easily test this
	// but we can verify a command was returned
}

func TestWelcomeModelView(t *testing.T) {
	context := &Context{
		Styles: GetStyles(),
	}
	model := NewWelcomeModel(context)
	
	view := model.View()
	if view == "" {
		t.Error("Expected non-empty view")
	}
	
	// Should contain banner and welcome content
	if !containsString(view, "Startup") {
		t.Error("Expected view to contain 'Startup'")
	}
	
	if !containsString(view, "Enter to begin") {
		t.Error("Expected view to contain navigation instructions")
	}
}