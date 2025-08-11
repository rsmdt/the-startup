package ui

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// WelcomeModel handles the welcome screen
type WelcomeModel struct {
	context *Context
}

// NewWelcomeModel creates a new welcome model
func NewWelcomeModel(context *Context) *WelcomeModel {
	return &WelcomeModel{
		context: context,
	}
}

// Init initializes the welcome model with auto-advance timer
func (m *WelcomeModel) Init() tea.Cmd {
	// Auto-advance from welcome screen after 2 seconds
	return tea.Tick(2*time.Second, func(time.Time) tea.Msg {
		return AutoAdvanceMsg{}
	})
}

// Update handles messages for the welcome screen
func (m *WelcomeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case AutoAdvanceMsg:
		// Auto-advance to tool selection
		return m, func() tea.Msg {
			return ViewTransitionMsg{NextView: StateToolSelection}
		}
		
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			// Manual advance to tool selection
			return m, func() tea.Msg {
				return ViewTransitionMsg{NextView: StateToolSelection}
			}
		}
	}
	
	return m, nil
}

// View renders the welcome screen
func (m *WelcomeModel) View() string {
	var s strings.Builder
	
	// ASCII art banner
	s.WriteString(m.context.Styles.Title.Render(WelcomeBanner))
	s.WriteString("\n")
	
	// Welcome message
	s.WriteString(m.context.Styles.Normal.Render(WelcomeMessage))
	s.WriteString("\n\n")
	
	// Instructions
	s.WriteString(m.context.Styles.Help.Render("Press Enter to begin • Escape or Ctrl+C to exit"))
	
	return s.String()
}