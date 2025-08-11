package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// Theme represents a color theme for the application
type Theme struct {
	Name        string
	Primary     lipgloss.Color // Brand color, accents
	Success     lipgloss.Color // Success messages, checkmarks
	Error       lipgloss.Color // Errors, failures
	Warning     lipgloss.Color // Warnings, updates
	Info        lipgloss.Color // Information, hints
	Text        lipgloss.Color // Regular text
	TextMuted   lipgloss.Color // Help text, descriptions
	TextBright  lipgloss.Color // Emphasized text, cursor lines
	Background  lipgloss.Color // Background color
}

// Predefined themes
var (
	// CharmTheme is the default Charm color scheme
	CharmTheme = Theme{
		Name:        "charm",
		Primary:     lipgloss.Color("#FF06B7"), // Pink/Magenta
		Success:     lipgloss.Color("#04B575"), // Green
		Error:       lipgloss.Color("#FF4444"), // Red
		Warning:     lipgloss.Color("#FFA500"), // Orange
		Info:        lipgloss.Color("#3C7EFF"), // Blue
		Text:        lipgloss.Color("#FAFAFA"), // Light gray
		TextMuted:   lipgloss.Color("#606060"), // Dark gray
		TextBright:  lipgloss.Color("#42FF76"), // Bright green
		Background:  lipgloss.Color("#000000"), // Black
	}

	// CatppuccinTheme is the popular Catppuccin Mocha theme
	CatppuccinTheme = Theme{
		Name:        "catppuccin",
		Primary:     lipgloss.Color("#CBA6F7"), // Mauve
		Success:     lipgloss.Color("#A6E3A1"), // Green
		Error:       lipgloss.Color("#F38BA8"), // Red
		Warning:     lipgloss.Color("#FAB387"), // Peach
		Info:        lipgloss.Color("#89B4FA"), // Blue
		Text:        lipgloss.Color("#CDD6F4"), // Text
		TextMuted:   lipgloss.Color("#6C7086"), // Overlay
		TextBright:  lipgloss.Color("#A6E3A1"), // Green
		Background:  lipgloss.Color("#1E1E2E"), // Base
	}
)

// CurrentTheme is the active theme (default to Charm)
var CurrentTheme = CharmTheme

// Styles based on current theme
type Styles struct {
	Title          lipgloss.Style
	Success        lipgloss.Style
	Error          lipgloss.Style
	Warning        lipgloss.Style
	Info           lipgloss.Style
	Help           lipgloss.Style
	Cursor         lipgloss.Style
	CursorLine     lipgloss.Style
	Normal         lipgloss.Style
	Selected       lipgloss.Style
	Spinner        lipgloss.Style
	ProgressBar    lipgloss.Style
	ProgressFilled lipgloss.Style
	ProgressEmpty  lipgloss.Style
}

// GetStyles returns styles for the current theme
func GetStyles() Styles {
	return Styles{
		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#000000")). // Black text
			Background(lipgloss.Color("#42FF76")). // Bright green background
			Padding(0, 1).
			MarginTop(1).
			MarginBottom(1),

		Success: lipgloss.NewStyle().
			Foreground(CurrentTheme.Success),

		Error: lipgloss.NewStyle().
			Foreground(CurrentTheme.Error),

		Warning: lipgloss.NewStyle().
			Foreground(CurrentTheme.Warning),

		Info: lipgloss.NewStyle().
			Foreground(CurrentTheme.Info),

		Help: lipgloss.NewStyle().
			Foreground(CurrentTheme.TextMuted),

		Cursor: lipgloss.NewStyle().
			Foreground(CurrentTheme.Primary),

		CursorLine: lipgloss.NewStyle().
			Foreground(CurrentTheme.TextBright),

		Normal: lipgloss.NewStyle().
			Foreground(CurrentTheme.Text),

		Selected: lipgloss.NewStyle().
			Foreground(CurrentTheme.Primary).
			Bold(true),

		Spinner: lipgloss.NewStyle().
			Foreground(CurrentTheme.Primary),

		ProgressBar: lipgloss.NewStyle().
			Foreground(CurrentTheme.TextMuted),

		ProgressFilled: lipgloss.NewStyle().
			Foreground(CurrentTheme.Primary).
			Background(CurrentTheme.Primary),

		ProgressEmpty: lipgloss.NewStyle().
			Foreground(CurrentTheme.TextMuted).
			Background(CurrentTheme.Background),
	}
}

// SetTheme changes the current theme
func SetTheme(theme Theme) {
	CurrentTheme = theme
}

// Icons for consistent UI elements
const (
	IconSuccess   = "✓"
	IconError     = "✗"
	IconWarning   = "⚠"
	IconInfo      = "ℹ"
	IconBullet    = "•"
	IconArrow     = "→"
	IconUpdate    = "↻"
	IconSelected  = "●"
	IconUnselected = "○"
	IconFolder    = "📁"
	IconFile      = "📄"
	IconRocket    = "🚀"
)