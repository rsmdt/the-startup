package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type InstallingModel struct {
	styles          Styles
	renderer        *ProgressiveDisclosureRenderer
	installProgress InstallProgressMsg
	installComplete bool
	success         bool
	err             error
}

func NewInstallingModel() InstallingModel {
	return InstallingModel{
		styles:   GetStyles(),
		renderer: NewProgressiveDisclosureRenderer(),
	}
}

func (m InstallingModel) Init() tea.Cmd {
	return nil
}

func (m InstallingModel) Update(msg tea.Msg) (InstallingModel, tea.Cmd) {
	switch msg := msg.(type) {
	case InstallProgressMsg:
		m.installProgress = msg
	case InstallCompleteMsg:
		m.installComplete = true
		m.success = msg.Success
		m.err = msg.Error
	}
	return m, nil
}

func (m InstallingModel) View() string {
	var s strings.Builder
	
	s.WriteString(m.styles.Title.Render(AppBanner))
	s.WriteString("\n\n")
	
	s.WriteString(m.renderer.RenderTitle("Installing The (Agentic) Startup"))
	
	if m.installProgress.Stage != "" {
		s.WriteString(m.styles.Info.Render(fmt.Sprintf("Stage: %s", m.installProgress.Stage)))
		s.WriteString("\n")
	}
	
	if m.installProgress.Message != "" {
		s.WriteString(m.styles.Normal.Render(m.installProgress.Message))
		s.WriteString("\n")
	}
	
	s.WriteString("\n")
	s.WriteString(m.styles.Info.Render("⚡ Installing components..."))
	s.WriteString("\n\n")
	
	return s.String()
}

func (m InstallingModel) IsComplete() bool {
	return m.installComplete
}

func (m InstallingModel) IsSuccess() bool {
	return m.success
}

func (m InstallingModel) Error() error {
	return m.err
}