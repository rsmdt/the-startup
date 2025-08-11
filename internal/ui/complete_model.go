package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/the-startup/the-startup/internal/installer"
)

type CompleteModel struct {
	styles       Styles
	installer    *installer.Installer
	selectedTool string
	ready        bool
}

func NewCompleteModel(selectedTool string, installer *installer.Installer) CompleteModel {
	return CompleteModel{
		styles:       GetStyles(),
		installer:    installer,
		selectedTool: selectedTool,
		ready:        false,
	}
}

func (m CompleteModel) Init() tea.Cmd {
	return nil
}

func (m CompleteModel) Update(msg tea.Msg) (CompleteModel, tea.Cmd) {
	if _, ok := msg.(tea.KeyMsg); ok {
		m.ready = true
	}
	return m, nil
}

func (m CompleteModel) View() string {
	var s strings.Builder
	
	s.WriteString(m.styles.Title.Render(AppBanner))
	s.WriteString("\n\n")
	
	s.WriteString(m.styles.Success.Render("✅ Installation Complete!"))
	s.WriteString("\n\n")
	
	installLocation := m.installer.GetInstallPath()
	s.WriteString(m.styles.Normal.Render(fmt.Sprintf("The Startup has been installed to: %s", installLocation)))
	s.WriteString("\n\n")
	
	s.WriteString(m.styles.Info.Render("Next steps:"))
	s.WriteString("\n")
	s.WriteString(m.styles.Normal.Render("1. Restart your " + m.selectedTool + " session"))
	s.WriteString("\n")
	
	isLocal := strings.Contains(installLocation, ".the-startup") && !strings.Contains(installLocation, ".config")
	if isLocal {
		s.WriteString(m.styles.Normal.Render("2. Use /develop command in this project"))
	} else {
		s.WriteString(m.styles.Normal.Render("2. Use /develop command in any project"))
	}
	s.WriteString("\n")
	s.WriteString(m.styles.Normal.Render("3. Check logs in ~/.the-startup/"))
	s.WriteString("\n\n")
	
	s.WriteString(m.styles.Help.Render("Press any key to exit"))
	
	return s.String()
}

func (m CompleteModel) Ready() bool {
	return m.ready
}