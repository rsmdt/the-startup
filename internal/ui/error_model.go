package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type ErrorModel struct {
	styles     Styles
	renderer   *ProgressiveDisclosureRenderer
	err        error
	errContext string
	ready      bool
}

func NewErrorModel(err error, context string) ErrorModel {
	return ErrorModel{
		styles:     GetStyles(),
		renderer:   NewProgressiveDisclosureRenderer(),
		err:        err,
		errContext: context,
		ready:      false,
	}
}

func (m ErrorModel) Init() tea.Cmd {
	return nil
}

func (m ErrorModel) Update(msg tea.Msg) (ErrorModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.ready = true
		}
	}
	return m, nil
}

func (m ErrorModel) View() string {
	var s strings.Builder
	
	s.WriteString(m.styles.Title.Render(AppBanner))
	s.WriteString("\n\n")
	
	s.WriteString(m.renderer.RenderError(m.err, m.errContext))
	s.WriteString("\n")
	
	s.WriteString(m.styles.Help.Render("Press Escape to go back"))
	
	return s.String()
}

func (m ErrorModel) Ready() bool {
	return m.ready
}

func (m ErrorModel) SetError(err error, context string) ErrorModel {
	m.err = err
	m.errContext = context
	m.ready = false
	return m
}