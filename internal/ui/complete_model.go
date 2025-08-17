package ui

import (
	"os"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rsmdt/the-startup/internal/installer"
)

type CompleteModel struct {
	styles       Styles
	installer    *installer.Installer
	selectedTool string
	ready        bool
}

// autoExitMsg is sent after a delay to automatically exit
type autoExitMsg struct{}

func NewCompleteModel(selectedTool string, installer *installer.Installer) CompleteModel {
	return CompleteModel{
		styles:       GetStyles(),
		installer:    installer,
		selectedTool: selectedTool,
		ready:        false,
	}
}

func (m CompleteModel) Init() tea.Cmd {
	// Exit immediately
	return func() tea.Msg {
		return autoExitMsg{}
	}
}

func (m CompleteModel) Update(msg tea.Msg) (CompleteModel, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		// Allow immediate exit on any key press
		m.ready = true
	case autoExitMsg:
		// Auto-exit after delay
		m.ready = true
	}
	return m, nil
}

func (m CompleteModel) View() string {
	var s strings.Builder

	// Banner
	s.WriteString(m.styles.Title.Render(AppBanner))
	s.WriteString("\n\n")

	// Success message
	s.WriteString(m.styles.Success.Render("✅ Installation Complete!"))
	s.WriteString("\n\n")

	// Installation locations
	claudePath := m.installer.GetClaudePath()
	startupPath := m.installer.GetInstallPath()

	// Simplify paths for display
	displayClaudePath := claudePath
	displayStartupPath := startupPath

	home := os.Getenv("HOME")
	if home != "" {
		if strings.HasPrefix(claudePath, home) {
			displayClaudePath = strings.Replace(claudePath, home, "~", 1)
		}
		if strings.HasPrefix(startupPath, home) {
			displayStartupPath = strings.Replace(startupPath, home, "~", 1)
		}
	}

	s.WriteString(m.styles.Normal.Render("Installation locations:"))
	s.WriteString("\n")
	s.WriteString(m.styles.Info.Render("  Claude files: " + displayClaudePath))
	s.WriteString("\n")
	s.WriteString(m.styles.Info.Render("  Startup files: " + displayStartupPath))
	s.WriteString("\n\n")

	// Available commands and agents
	s.WriteString(m.styles.Normal.Render("Now available in Claude Code:"))
	s.WriteString("\n\n")

	// Display installed commands
	commands := m.installer.GetInstalledCommands()
	if len(commands) > 0 {
		sort.Strings(commands)
		s.WriteString(m.styles.Info.Render("  Commands:"))
		s.WriteString("\n")
		for _, cmd := range commands {
			s.WriteString(m.styles.Normal.Render("    • " + cmd))
			s.WriteString("\n")
		}
		s.WriteString("\n")
	}

	// Display installed agents
	agents := m.installer.GetInstalledAgents()
	if len(agents) > 0 {
		sort.Strings(agents)
		s.WriteString(m.styles.Info.Render("  Agents:"))
		s.WriteString("\n")
		for _, agent := range agents {
			s.WriteString(m.styles.Normal.Render("    • " + agent))
			s.WriteString("\n")
		}
		s.WriteString("\n")
	}

	// Add repository link
	s.WriteString(m.styles.Info.Render("See https://github.com/rsmdt/the-startup for details"))
	s.WriteString("\n")

	return s.String()
}

func (m CompleteModel) Ready() bool {
	return m.ready
}
