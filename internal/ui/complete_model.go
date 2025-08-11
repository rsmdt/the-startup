package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// CompleteModel handles installation completion display
type CompleteModel struct {
	context *Context
}

// NewCompleteModel creates a new complete model
func NewCompleteModel(context *Context) *CompleteModel {
	return &CompleteModel{
		context: context,
	}
}

// Init initializes the complete model
func (m *CompleteModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for completion screen
func (m *CompleteModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		// Any key exits
		return m, tea.Quit
	}
	
	return m, nil
}

// View renders the completion screen
func (m *CompleteModel) View() string {
	var s strings.Builder
	
	// ASCII art banner
	s.WriteString(m.context.Styles.Title.Render(WelcomeBanner))
	s.WriteString("\n\n")
	
	// Success message
	s.WriteString(m.context.Styles.Success.Render("✅ Installation Complete!"))
	s.WriteString("\n\n")
	
	// Installation details
	installLocation := m.context.Installer.GetInstallPath()
	s.WriteString(m.context.Styles.Normal.Render(fmt.Sprintf("The Startup has been installed to: %s", installLocation)))
	s.WriteString("\n\n")
	
	// Next steps
	s.WriteString(m.context.Styles.Info.Render("Next steps:"))
	s.WriteString("\n")
	s.WriteString(m.context.Styles.Normal.Render("1. Restart your " + m.context.SelectedTool + " session"))
	s.WriteString("\n")
	
	isLocal := strings.Contains(installLocation, ".the-startup") && !strings.Contains(installLocation, ".config")
	if isLocal {
		s.WriteString(m.context.Styles.Normal.Render("2. Use /develop command in this project"))
	} else {
		s.WriteString(m.context.Styles.Normal.Render("2. Use /develop command in any project"))
	}
	s.WriteString("\n")
	s.WriteString(m.context.Styles.Normal.Render("3. Check logs in ~/.the-startup/"))
	s.WriteString("\n\n")
	
	// Exit instruction
	s.WriteString(m.context.Styles.Help.Render("Press any key to exit"))
	
	return s.String()
}