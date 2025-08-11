package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

// HuhConfirmationModel handles huh confirmation dialog
type HuhConfirmationModel struct {
	context      *Context
	huhForm      *huh.Form
	huhConfirmed bool
}

// NewHuhConfirmationModel creates a new huh confirmation model
func NewHuhConfirmationModel(context *Context) *HuhConfirmationModel {
	model := &HuhConfirmationModel{
		context: context,
	}
	
	// Create huh confirmation form
	model.huhForm = huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Ready to install?").
				Description("This will install The (Agentic) Startup to your .claude directory.").
				Affirmative("Yes, give me awesome").
				Negative("Huh? I did not sign up for this").
				Value(&model.huhConfirmed),
		),
	)
	
	return model
}

// Init initializes the huh confirmation model
func (m *HuhConfirmationModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for huh confirmation
func (m *HuhConfirmationModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		
		case "esc":
			// Go back to file selection
			return m, func() tea.Msg {
				return ViewTransitionMsg{NextView: StateFileSelection}
			}
		}
	}
	
	// Handle huh form updates
	if m.huhForm != nil {
		form, cmd := m.huhForm.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.huhForm = f
			
			// Check if form is complete
			if m.huhForm.State == huh.StateCompleted {
				// Process the confirmation result
				if m.huhConfirmed {
					return m, func() tea.Msg {
						return ViewTransitionMsg{NextView: StateConfirmation}
					}
				} else {
					return m, func() tea.Msg {
						return ViewTransitionMsg{NextView: StateFileSelection}
					}
				}
			}
		}
		return m, cmd
	}
	
	return m, nil
}

// View renders the huh confirmation screen
func (m *HuhConfirmationModel) View() string {
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
	
	// Render the huh form
	if m.huhForm != nil {
		s.WriteString(m.huhForm.View())
	}
	
	return s.String()
}