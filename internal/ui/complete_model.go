package ui

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/the-startup/the-startup/internal/installer"
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
	// Automatically exit after 3 seconds
	return tea.Tick(3*time.Second, func(t time.Time) tea.Msg {
		return autoExitMsg{}
	})
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
	
	// Installation location
	installLocation := m.installer.GetInstallPath()
	if installLocation == "" {
		installLocation = "~/.config/the-startup/.claude"
	}
	
	// Simplify the path for display
	displayPath := installLocation
	if strings.Contains(installLocation, ".config/the-startup") {
		displayPath = "~/.config/the-startup/.claude"
	} else if strings.Contains(installLocation, ".the-startup") {
		displayPath = ".the-startup/.claude (local)"
	}
	
	s.WriteString(m.styles.Normal.Render("Installation location:"))
	s.WriteString("\n")
	s.WriteString(m.styles.Info.Render("  " + displayPath))
	s.WriteString("\n\n")
	
	// Next steps
	s.WriteString(m.styles.Normal.Render("Next steps:"))
	s.WriteString("\n\n")
	
	// Step 1
	s.WriteString(m.styles.Normal.Render("  1. Restart your Claude Code session"))
	s.WriteString("\n")
	
	// Step 2
	isLocal := strings.Contains(installLocation, ".the-startup") && !strings.Contains(installLocation, ".config")
	if isLocal {
		s.WriteString(m.styles.Normal.Render("  2. Use /develop command in this project"))
	} else {
		s.WriteString(m.styles.Normal.Render("  2. Use /develop command in any project"))
	}
	s.WriteString("\n")
	
	// Step 3
	s.WriteString(m.styles.Normal.Render("  3. Check logs in ~/.the-startup/"))
	s.WriteString("\n\n")
	
	// Auto-exit message
	s.WriteString(m.styles.Help.Render("Exiting in 3 seconds..."))
	
	return s.String()
}

func (m CompleteModel) Ready() bool {
	return m.ready
}