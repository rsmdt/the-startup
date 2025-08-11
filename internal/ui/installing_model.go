package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// InstallingModel handles installation progress display
type InstallingModel struct {
	context         *Context
	installProgress InstallProgressMsg
	installComplete bool
}

// NewInstallingModel creates a new installing model
func NewInstallingModel(context *Context) *InstallingModel {
	return &InstallingModel{
		context: context,
	}
}

// Init initializes the installing model and starts installation
func (m *InstallingModel) Init() tea.Cmd {
	return m.performInstallation()
}

// performInstallation starts the installation process
func (m *InstallingModel) performInstallation() tea.Cmd {
	return func() tea.Msg {
		// Perform actual installation
		if err := m.context.Installer.Install(); err != nil {
			return InstallCompleteMsg{Success: false, Error: err}
		}
		return InstallCompleteMsg{Success: true}
	}
}

// Update handles messages for installing
func (m *InstallingModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case InstallProgressMsg:
		m.installProgress = msg
		return m, nil
	
	case InstallCompleteMsg:
		m.installComplete = true
		if msg.Success {
			return m, func() tea.Msg {
				return ViewTransitionMsg{NextView: StateComplete}
			}
		} else {
			return m, func() tea.Msg {
				return ViewTransitionMsg{NextView: StateError, Data: map[string]interface{}{
					"error": msg.Error,
					"context": "during installation",
				}}
			}
		}
	
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	
	return m, nil
}

// View renders the installing screen
func (m *InstallingModel) View() string {
	var s strings.Builder
	
	// ASCII art banner
	s.WriteString(m.context.Styles.Title.Render(WelcomeBanner))
	s.WriteString("\n\n")
	
	// Title
	s.WriteString(m.context.Renderer.RenderTitle("Installing The (Agentic) Startup"))
	
	// Progress information
	if m.installProgress.Stage != "" {
		s.WriteString(m.context.Styles.Info.Render(fmt.Sprintf("Stage: %s", m.installProgress.Stage)))
		s.WriteString("\n")
	}
	
	if m.installProgress.Message != "" {
		s.WriteString(m.context.Styles.Normal.Render(m.installProgress.Message))
		s.WriteString("\n")
	}
	
	// Simple progress indicator (spinner could be added here)
	s.WriteString("\n")
	s.WriteString(m.context.Styles.Info.Render("⚡ Installing components..."))
	s.WriteString("\n\n")
	
	// Progress bar could be added here based on m.installProgress.Percent
	
	return s.String()
}