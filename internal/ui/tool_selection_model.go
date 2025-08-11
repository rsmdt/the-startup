package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// ToolSelectionModel handles tool selection
type ToolSelectionModel struct {
	context *Context
	cursor  int
	choices []string
}

// NewToolSelectionModel creates a new tool selection model
func NewToolSelectionModel(context *Context) *ToolSelectionModel {
	return &ToolSelectionModel{
		context: context,
		cursor:  0,
		choices: []string{"claude-code", "Cancel"},
	}
}

// Init initializes the tool selection model
func (m *ToolSelectionModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for tool selection
func (m *ToolSelectionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			// Go back to welcome or quit
			return m, func() tea.Msg {
				return ViewTransitionMsg{NextView: StateWelcome}
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
			choice := m.choices[m.cursor]
			if choice == "Cancel" {
				return m, tea.Quit
			}
			
			// Selection made - update context and advance
			return m, tea.Batch(
				func() tea.Msg {
					return SelectionMadeMsg{Tool: choice}
				},
				func() tea.Msg {
					return ViewTransitionMsg{NextView: StatePathSelection}
				},
			)
		}
	}
	
	return m, nil
}

// View renders the tool selection screen
func (m *ToolSelectionModel) View() string {
	var s strings.Builder
	
	// ASCII art banner
	s.WriteString(m.context.Styles.Title.Render(WelcomeBanner))
	s.WriteString("\n\n")
	
	// Progressive disclosure header (empty for first step)
	s.WriteString(m.context.Renderer.RenderSelections("", "", 0))
	
	// Title
	s.WriteString(m.context.Renderer.RenderTitle("Select your development tool"))
	
	// Choices
	for i, choice := range m.choices {
		displayText := choice
		if choice != "Cancel" {
			displayText = "Claude Code"
		}
		s.WriteString(m.context.Renderer.RenderChoiceWithMultiSelect(displayText, i == m.cursor, false, false))
		s.WriteString("\n")
	}
	
	// Help
	s.WriteString(m.context.Renderer.RenderHelp("↑↓ navigate • Enter: select • Escape: back"))
	
	return s.String()
}