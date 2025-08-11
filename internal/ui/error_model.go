package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// ErrorModel handles error display
type ErrorModel struct {
	context    *Context
	err        error
	errContext string
}

// NewErrorModel creates a new error model
func NewErrorModel(context *Context) *ErrorModel {
	return &ErrorModel{
		context: context,
	}
}

// SetError sets the error and context for display
func (m *ErrorModel) SetError(err error, context string) {
	m.err = err
	m.errContext = context
}

// Init initializes the error model
func (m *ErrorModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the error screen
func (m *ErrorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		
		case "esc":
			// Return to welcome screen for recovery
			return m, func() tea.Msg {
				return ViewTransitionMsg{NextView: StateWelcome}
			}
		}
	}
	
	return m, nil
}

// View renders the error screen
func (m *ErrorModel) View() string {
	var s strings.Builder
	
	// ASCII art banner
	s.WriteString(m.context.Styles.Title.Render(WelcomeBanner))
	s.WriteString("\n\n")
	
	// Error message
	s.WriteString(m.context.Renderer.RenderError(m.err, m.errContext))
	s.WriteString("\n")
	
	// Recovery instructions
	s.WriteString(m.context.Styles.Help.Render("Press Escape to return to welcome screen"))
	
	return s.String()
}