package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// ConfirmationModel handles final installation confirmation
type ConfirmationModel struct {
	context *Context
	cursor  int
	choices []string
}

// NewConfirmationModel creates a new confirmation model
func NewConfirmationModel(context *Context) *ConfirmationModel {
	return &ConfirmationModel{
		context: context,
		cursor:  0,
		choices: []string{"Install", "Go back"},
	}
}

// Init initializes the confirmation model
func (m *ConfirmationModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for confirmation
func (m *ConfirmationModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		
		case "esc":
			// Go back to huh confirmation
			return m, func() tea.Msg {
				return ViewTransitionMsg{NextView: StateHuhConfirmation}
			}
		
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		
		case "enter":
			if m.cursor == 0 { // "Install" option
				return m, func() tea.Msg {
					return ViewTransitionMsg{NextView: StateInstalling}
				}
			} else { // "Go back" option
				return m, func() tea.Msg {
					return ViewTransitionMsg{NextView: StateHuhConfirmation}
				}
			}
		}
	}
	
	return m, nil
}

// View renders the confirmation screen
func (m *ConfirmationModel) View() string {
	var s strings.Builder
	
	// ASCII art banner
	s.WriteString(m.context.Styles.Title.Render(WelcomeBanner))
	s.WriteString("\n\n")
	
	// Progressive disclosure header
	displayPath := m.context.SelectedPath
	if displayPath == "" {
		displayPath = "~/.config/the-startup"
	}
	s.WriteString(m.context.Renderer.RenderSelections(m.context.SelectedTool, displayPath, len(m.context.SelectedFiles)))
	
	// Installation summary
	updatingFiles := []string{}
	if len(m.context.SelectedFiles) > 0 {
		for _, file := range m.context.SelectedFiles {
			if m.context.Installer.CheckFileExists(file) {
				updatingFiles = append(updatingFiles, file)
			}
		}
	}
	
	s.WriteString(m.context.Renderer.RenderInstallationSummary(
		m.context.SelectedTool, displayPath, m.context.SelectedFiles, updatingFiles))
	
	// Confirmation choices
	for i, choice := range m.choices {
		s.WriteString(m.context.Renderer.RenderChoiceWithMultiSelect(choice, i == m.cursor, false, false))
		s.WriteString("\n")
	}
	
	// Help
	s.WriteString(m.context.Renderer.RenderHelp("↑↓ navigate • Enter: select • Escape: back"))
	
	return s.String()
}