package ui

import (
	"embed"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rsmdt/the-startup/internal/installer"
)

// ViewModel represents a composable view in the installer
type ViewModel interface {
	tea.Model
}

// ViewTransitionMsg requests a transition to a different view
type ViewTransitionMsg struct {
	NextView InstallerState
	Data     interface{} // Optional data to pass to next view
}

// SelectionMadeMsg represents a user selection being made
type SelectionMadeMsg struct {
	Tool  string
	Path  string
	Files []string
}

// Context contains shared dependencies and state for all view models
type Context struct {
	// Dependencies
	Installer     *installer.Installer
	AgentFiles    *embed.FS
	CommandFiles  *embed.FS
	HookFiles     *embed.FS
	TemplateFiles *embed.FS
	Styles        Styles
	Renderer      *ProgressiveDisclosureRenderer
	
	// Shared state
	SelectedTool  string
	SelectedPath  string
	SelectedFiles []string
	Width         int
	Height        int
}

// AutoAdvanceMsg represents a timer tick for auto-advance
type AutoAdvanceMsg struct{}

// BackMsg represents a back navigation request
type BackMsg struct{}