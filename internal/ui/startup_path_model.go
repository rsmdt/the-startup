package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textinput"
)

type StartupPathModel struct {
	styles          Styles
	renderer        *ProgressiveDisclosureRenderer
	choices         []string
	cursor          int
	selectedPath    string
	ready           bool
	inputMode       bool
	textInput       textinput.Model
	suggestions     []string
	suggestionIndex int
}

func NewStartupPathModel() StartupPathModel {
	cwd, _ := os.Getwd()
	homeDir, _ := os.UserHomeDir()
	
	// Format local path with tilde notation
	localFullPath := filepath.Join(cwd, ".the-startup")
	if strings.HasPrefix(localFullPath, homeDir) {
		localFullPath = "~" + strings.TrimPrefix(localFullPath, homeDir)
	}
	localPath := fmt.Sprintf("%s (local)", localFullPath)
	
	ti := textinput.New()
	ti.Placeholder = "Enter custom path (Tab for autocomplete)"
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 60
	
	return StartupPathModel{
		styles:       GetStyles(),
		renderer:     NewProgressiveDisclosureRenderer(),
		choices: []string{
			"~/.config/the-startup (recommended)",
			localPath,
			"Custom location",
			"Cancel",
		},
		cursor:          0,
		ready:           false,
		inputMode:       false,
		textInput:       ti,
		suggestions:     []string{},
		suggestionIndex: -1,
	}
}

func (m StartupPathModel) Init() tea.Cmd {
	return nil
}

func (m StartupPathModel) Update(msg tea.Msg) (StartupPathModel, tea.Cmd) {
	if m.inputMode {
		return m.updateInputMode(msg)
	}
	
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
				switch {
				case choice == "Cancel":
					m.selectedPath = "CANCEL"
					m.ready = true
				case choice == "~/.config/the-startup (recommended)":
					homeDir, _ := os.UserHomeDir()
					m.selectedPath = filepath.Join(homeDir, ".config", "the-startup")
					m.ready = true
				case strings.Contains(choice, "/.the-startup (local)"):
					// Local option
					cwd, _ := os.Getwd()
					m.selectedPath = filepath.Join(cwd, ".the-startup")
					m.ready = true
				case choice == "Custom location":
					m.inputMode = true
					m.textInput.SetValue("")
					m.textInput.Focus()
					return m, textinput.Blink
				}
			}
		}
	}
	return m, nil
}

func (m StartupPathModel) updateInputMode(msg tea.Msg) (StartupPathModel, tea.Cmd) {
	var cmd tea.Cmd
	
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.inputMode = false
			m.suggestions = []string{}
			m.suggestionIndex = -1
			return m, nil
			
		case "tab":
			// Handle autocomplete
			currentPath := m.textInput.Value()
			m.suggestions = m.getPathSuggestions(currentPath)
			
			if len(m.suggestions) > 0 {
				m.suggestionIndex = (m.suggestionIndex + 1) % len(m.suggestions)
				m.textInput.SetValue(m.suggestions[m.suggestionIndex])
			}
			return m, nil
			
		case "enter":
			path := m.textInput.Value()
			if path != "" {
				// Expand tilde
				if strings.HasPrefix(path, "~/") {
					homeDir, _ := os.UserHomeDir()
					path = filepath.Join(homeDir, path[2:])
				}
				
				// Ensure it ends with .the-startup
				if !strings.HasSuffix(path, ".the-startup") {
					path = filepath.Join(path, ".the-startup")
				}
				
				m.selectedPath = path
				m.ready = true
				m.inputMode = false
			}
			return m, nil
		
		default:
			// Reset suggestions when typing
			if msg.String() != "tab" {
				m.suggestions = []string{}
				m.suggestionIndex = -1
			}
		}
	}
	
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m StartupPathModel) getPathSuggestions(input string) []string {
	suggestions := []string{}
	
	// Expand tilde for suggestions
	expandedInput := input
	if strings.HasPrefix(input, "~/") {
		homeDir, _ := os.UserHomeDir()
		expandedInput = filepath.Join(homeDir, input[2:])
	}
	
	// Get directory to search
	dir := filepath.Dir(expandedInput)
	base := filepath.Base(expandedInput)
	
	// If directory doesn't exist, try parent
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		dir = filepath.Dir(dir)
		base = ""
	}
	
	// Read directory entries
	entries, err := os.ReadDir(dir)
	if err != nil {
		return suggestions
	}
	
	// Find matching directories
	for _, entry := range entries {
		if entry.IsDir() && strings.HasPrefix(entry.Name(), base) {
			fullPath := filepath.Join(dir, entry.Name())
			
			// Convert back to tilde notation if in home directory
			homeDir, _ := os.UserHomeDir()
			if strings.HasPrefix(fullPath, homeDir) {
				fullPath = "~" + strings.TrimPrefix(fullPath, homeDir)
			}
			
			suggestions = append(suggestions, fullPath)
			
			// Limit suggestions
			if len(suggestions) >= 5 {
				break
			}
		}
	}
	
	return suggestions
}

func (m StartupPathModel) View() string {
	var s strings.Builder
	
	s.WriteString(m.styles.Title.Render(AppBanner))
	s.WriteString("\n\n")
	
	s.WriteString(m.renderer.RenderTitle("Select .the-startup installation location"))
	s.WriteString(m.styles.Info.Render("This is where The Startup's binary and templates will be installed"))
	s.WriteString("\n\n")
	
	if m.inputMode {
		s.WriteString(m.styles.Normal.Render("Enter custom path:"))
		s.WriteString("\n")
		s.WriteString(m.textInput.View())
		s.WriteString("\n\n")
		
		if len(m.suggestions) > 0 {
			s.WriteString(m.styles.Help.Render("Suggestions (Tab to cycle):"))
			s.WriteString("\n")
			for i, suggestion := range m.suggestions {
				if i == m.suggestionIndex {
					s.WriteString(m.styles.Selected.Render("  → " + suggestion))
				} else {
					s.WriteString(m.styles.Normal.Render("    " + suggestion))
				}
				s.WriteString("\n")
			}
			s.WriteString("\n")
		}
		
		s.WriteString(m.renderer.RenderHelp("Tab: autocomplete • Enter: confirm • Escape: cancel"))
	} else {
		for i, choice := range m.choices {
			s.WriteString(m.renderer.RenderChoiceWithMultiSelect(choice, i == m.cursor, false, false))
			s.WriteString("\n")
		}
		
		s.WriteString(m.renderer.RenderHelp("↑↓ navigate • Enter: select • Escape: quit"))
	}
	
	return s.String()
}

func (m StartupPathModel) Ready() bool {
	return m.ready
}

func (m StartupPathModel) SelectedPath() string {
	return m.selectedPath
}

func (m StartupPathModel) Reset() StartupPathModel {
	m.ready = false
	m.cursor = 0
	m.inputMode = false
	m.textInput.SetValue("")
	m.suggestions = []string{}
	m.suggestionIndex = -1
	return m
}