package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type PathSelectionModel struct {
	styles       Styles
	renderer     *ProgressiveDisclosureRenderer
	choices      []string
	cursor       int
	selectedTool string
	selectedPath string
	ready        bool
}

func NewPathSelectionModel(selectedTool string) PathSelectionModel {
	cwd, _ := os.Getwd()
	localPath := fmt.Sprintf(".the-startup (local to %s)", cwd)

	return PathSelectionModel{
		styles:       GetStyles(),
		renderer:     NewProgressiveDisclosureRenderer(),
		selectedTool: selectedTool,
		choices: []string{
			"~/.config/the-startup (recommended)",
			localPath,
			"Custom location",
			"Cancel",
		},
		cursor: 0,
		ready:  false,
	}
}

func (m PathSelectionModel) Init() tea.Cmd {
	return nil
}

func (m PathSelectionModel) Update(msg tea.Msg) (PathSelectionModel, tea.Cmd) {
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
				choice := m.choices[m.cursor]
				switch choice {
				case "Cancel":
					m.selectedPath = "CANCEL"
					m.ready = true
				case "~/.config/the-startup (recommended)":
					m.selectedPath = ""
					m.ready = true
				case "Custom location":
					m.selectedPath = "CUSTOM"
					m.ready = true
				default:
					cwd, _ := os.Getwd()
					m.selectedPath = filepath.Join(cwd, ".the-startup")
					m.ready = true
				}
			}
		}
	}
	return m, nil
}

func (m PathSelectionModel) View() string {
	var s strings.Builder

	s.WriteString(m.styles.Title.Render(AppBanner))
	s.WriteString("\n\n")

	s.WriteString(m.renderer.RenderSelections(m.selectedTool, "", 0))

	s.WriteString(m.renderer.RenderTitle("Select installation location"))

	for i, choice := range m.choices {
		s.WriteString(m.renderer.RenderChoiceWithMultiSelect(choice, i == m.cursor, false, false))
		s.WriteString("\n")
	}

	s.WriteString(m.renderer.RenderHelp("↑↓ navigate • Enter: select • Escape: back"))

	return s.String()
}

func (m PathSelectionModel) Ready() bool {
	return m.ready
}

func (m PathSelectionModel) SelectedPath() string {
	return m.selectedPath
}

func (m PathSelectionModel) Reset() PathSelectionModel {
	m.ready = false
	m.cursor = 0
	return m
}
