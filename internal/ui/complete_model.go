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
	mode         OperationMode
}

// autoExitMsg is sent after a delay to automatically exit
type autoExitMsg struct{}

func NewCompleteModel(selectedTool string, installer *installer.Installer) CompleteModel {
	return NewCompleteModelWithMode(selectedTool, installer, ModeInstall)
}

func NewCompleteModelWithMode(selectedTool string, installer *installer.Installer, mode OperationMode) CompleteModel {
	return CompleteModel{
		styles:       GetStyles(),
		installer:    installer,
		selectedTool: selectedTool,
		ready:        false,
		mode:         mode,
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
	if m.mode == ModeUninstall {
		s.WriteString(m.styles.Success.Render("✅ Uninstallation Complete!"))
	} else {
		s.WriteString(m.styles.Success.Render("✅ Installation Complete!"))
	}
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

	if m.mode == ModeUninstall {
		s.WriteString(m.styles.Normal.Render("Uninstallation locations:"))
		s.WriteString("\n")
		s.WriteString(m.styles.Warning.Render("  Claude files: " + displayClaudePath))
		s.WriteString("\n")
		s.WriteString(m.styles.Warning.Render("  Startup files: " + displayStartupPath))
		s.WriteString("\n\n")
	} else {
		s.WriteString(m.styles.Normal.Render("Installation locations:"))
		s.WriteString("\n")
		s.WriteString(m.styles.Info.Render("  Claude files: " + displayClaudePath))
		s.WriteString("\n")
		s.WriteString(m.styles.Info.Render("  Startup files: " + displayStartupPath))
		s.WriteString("\n\n")
	}

	// Display commands (installed or removed)
	commands := m.installer.GetInstalledCommands()
	if len(commands) > 0 {
		sort.Strings(commands)
		if m.mode == ModeUninstall {
			s.WriteString(m.styles.Warning.Render("  Commands removed:"))
		} else {
			s.WriteString(m.styles.Info.Render("  Commands:"))
		}
		s.WriteString("\n")
		for _, cmd := range commands {
			if m.mode == ModeUninstall {
				s.WriteString(m.styles.Normal.Render("    ✗ " + cmd))
			} else {
				s.WriteString(m.styles.Normal.Render("    • " + cmd))
			}
			s.WriteString("\n")
		}
		s.WriteString("\n")
	}

	// Display agents (installed or removed)
	agents := m.installer.GetInstalledAgents()
	if len(agents) > 0 {
		sort.Strings(agents)
		if m.mode == ModeUninstall {
			s.WriteString(m.styles.Warning.Render("  Agents removed:"))
		} else {
			s.WriteString(m.styles.Info.Render("  Agents:"))
		}
		s.WriteString("\n")
		for _, agent := range agents {
			if m.mode == ModeUninstall {
				s.WriteString(m.styles.Normal.Render("    ✗ " + agent))
			} else {
				s.WriteString(m.styles.Normal.Render("    • " + agent))
			}
			s.WriteString("\n")
		}
		s.WriteString("\n")
	}

	// Display removed files if any
	removedFiles := m.installer.GetDeprecatedFilesList()
	if len(removedFiles) > 0 {
		sort.Strings(removedFiles)
		s.WriteString(m.styles.Error.Render("  Removed deprecated files:"))
		s.WriteString("\n")
		for _, file := range removedFiles {
			s.WriteString(m.styles.Normal.Render("    ✗ " + file))
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
