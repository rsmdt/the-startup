package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type ToolSelectionModel struct {
	styles   Styles
	renderer *ProgressiveDisclosureRenderer
	choices  []string
	cursor   int
	selected string
	ready    bool
}

func NewToolSelectionModel() ToolSelectionModel {
	return ToolSelectionModel{
		styles:   GetStyles(),
		renderer: NewProgressiveDisclosureRenderer(),
		choices:  []string{"claude-code", "Cancel"},
		cursor:   0,
		ready:    false,
	}
}

func (m ToolSelectionModel) Init() tea.Cmd {
	return nil
}

func (m ToolSelectionModel) Update(msg tea.Msg) (ToolSelectionModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter":
			if m.cursor < len(m.choices) {
				m.selected = m.choices[m.cursor]
				m.ready = true
			}
		}
	}
	return m, nil
}

func (m ToolSelectionModel) View() string {
	var s strings.Builder

	s.WriteString(m.styles.Title.Render(AppBanner))
	s.WriteString("\n\n")

	s.WriteString(m.renderer.RenderSelections("", "", 0))

	s.WriteString(m.renderer.RenderTitle("Select your development tool"))

	for i, choice := range m.choices {
		displayText := choice
		if choice != "Cancel" {
			displayText = "Claude Code"
		}
		s.WriteString(m.renderer.RenderChoiceWithMultiSelect(displayText, i == m.cursor, false, false))
		s.WriteString("\n")
	}

	s.WriteString(m.renderer.RenderHelp("↑↓ navigate • Enter: select • Escape: back"))

	return s.String()
}

func (m ToolSelectionModel) Ready() bool {
	return m.ready
}

func (m ToolSelectionModel) Selected() string {
	return m.selected
}

func (m ToolSelectionModel) Reset() ToolSelectionModel {
	m.ready = false
	m.cursor = 0
	return m
}
