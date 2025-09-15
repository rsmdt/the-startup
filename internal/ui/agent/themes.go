package agent

import "github.com/charmbracelet/lipgloss"

// ColorTheme represents a color scheme inspired by bat themes
type ColorTheme struct {
	Name        string
	Background  lipgloss.Color
	Foreground  lipgloss.Color
	Primary     lipgloss.Color   // Main accent color (timeline, headers)
	Secondary   lipgloss.Color   // Secondary accent (selected items)
	Success     lipgloss.Color   // Success/User messages
	Warning     lipgloss.Color   // Warnings/Assistant messages
	Error       lipgloss.Color   // Errors
	Info        lipgloss.Color   // Info text
	Muted       lipgloss.Color   // Muted/disabled text
	Border      lipgloss.Color   // Border colors
	Highlight   lipgloss.Color   // Highlighted/focused items
	GraphUser   lipgloss.Color   // User graph line
	GraphAsst   lipgloss.Color   // Assistant graph line
	Matrix      [5]lipgloss.Color // Matrix strength colors (weak to strong)
}

// Predefined themes inspired by popular bat colorschemes
var Themes = map[string]ColorTheme{
	"dracula": {
		Name:       "Dracula",
		Background: lipgloss.Color("#282a36"),
		Foreground: lipgloss.Color("#f8f8f2"),
		Primary:    lipgloss.Color("#bd93f9"), // Purple
		Secondary:  lipgloss.Color("#ff79c6"), // Pink
		Success:    lipgloss.Color("#50fa7b"), // Green
		Warning:    lipgloss.Color("#ffb86c"), // Orange
		Error:      lipgloss.Color("#ff5555"), // Red
		Info:       lipgloss.Color("#8be9fd"), // Cyan
		Muted:      lipgloss.Color("#6272a4"), // Comment
		Border:     lipgloss.Color("#44475a"), // Current line
		Highlight:  lipgloss.Color("#f1fa8c"), // Yellow
		GraphUser:  lipgloss.Color("#8be9fd"), // Cyan for user
		GraphAsst:  lipgloss.Color("#ff79c6"), // Pink for assistant
		Matrix: [5]lipgloss.Color{
			lipgloss.Color("#44475a"), // Minimal
			lipgloss.Color("#6272a4"), // Weak
			lipgloss.Color("#bd93f9"), // Medium
			lipgloss.Color("#ff79c6"), // Strong
			lipgloss.Color("#50fa7b"), // Very strong
		},
	},
	"nord": {
		Name:       "Nord",
		Background: lipgloss.Color("#2e3440"),
		Foreground: lipgloss.Color("#eceff4"),
		Primary:    lipgloss.Color("#88c0d0"), // Frost blue
		Secondary:  lipgloss.Color("#81a1c1"), // Frost blue dark
		Success:    lipgloss.Color("#a3be8c"), // Green
		Warning:    lipgloss.Color("#ebcb8b"), // Yellow
		Error:      lipgloss.Color("#bf616a"), // Red
		Info:       lipgloss.Color("#5e81ac"), // Blue
		Muted:      lipgloss.Color("#4c566a"), // Polar night
		Border:     lipgloss.Color("#434c5e"), // Polar night light
		Highlight:  lipgloss.Color("#d8dee9"), // Snow storm
		GraphUser:  lipgloss.Color("#88c0d0"), // Frost for user
		GraphAsst:  lipgloss.Color("#b48ead"), // Aurora purple for assistant
		Matrix: [5]lipgloss.Color{
			lipgloss.Color("#3b4252"), // Minimal
			lipgloss.Color("#434c5e"), // Weak
			lipgloss.Color("#81a1c1"), // Medium
			lipgloss.Color("#88c0d0"), // Strong
			lipgloss.Color("#a3be8c"), // Very strong
		},
	},
	"monokai": {
		Name:       "Monokai",
		Background: lipgloss.Color("#272822"),
		Foreground: lipgloss.Color("#f8f8f2"),
		Primary:    lipgloss.Color("#66d9ef"), // Blue
		Secondary:  lipgloss.Color("#a6e22e"), // Green
		Success:    lipgloss.Color("#a6e22e"), // Green
		Warning:    lipgloss.Color("#e6db74"), // Yellow
		Error:      lipgloss.Color("#f92672"), // Red
		Info:       lipgloss.Color("#66d9ef"), // Blue
		Muted:      lipgloss.Color("#75715e"), // Comment
		Border:     lipgloss.Color("#3e3d32"), // Line
		Highlight:  lipgloss.Color("#fd971f"), // Orange
		GraphUser:  lipgloss.Color("#66d9ef"), // Blue for user
		GraphAsst:  lipgloss.Color("#f92672"), // Red/pink for assistant
		Matrix: [5]lipgloss.Color{
			lipgloss.Color("#3e3d32"), // Minimal
			lipgloss.Color("#75715e"), // Weak
			lipgloss.Color("#e6db74"), // Medium
			lipgloss.Color("#fd971f"), // Strong
			lipgloss.Color("#a6e22e"), // Very strong
		},
	},
	"github": {
		Name:       "GitHub",
		Background: lipgloss.Color("#ffffff"),
		Foreground: lipgloss.Color("#24292e"),
		Primary:    lipgloss.Color("#0366d6"), // Blue
		Secondary:  lipgloss.Color("#28a745"), // Green
		Success:    lipgloss.Color("#28a745"), // Green
		Warning:    lipgloss.Color("#ffd33d"), // Yellow
		Error:      lipgloss.Color("#d73a49"), // Red
		Info:       lipgloss.Color("#0366d6"), // Blue
		Muted:      lipgloss.Color("#6a737d"), // Gray
		Border:     lipgloss.Color("#e1e4e8"), // Border
		Highlight:  lipgloss.Color("#f6f8fa"), // Light gray
		GraphUser:  lipgloss.Color("#0366d6"), // Blue for user
		GraphAsst:  lipgloss.Color("#6f42c1"), // Purple for assistant
		Matrix: [5]lipgloss.Color{
			lipgloss.Color("#f6f8fa"), // Minimal
			lipgloss.Color("#e1e4e8"), // Weak
			lipgloss.Color("#ffd33d"), // Medium
			lipgloss.Color("#f9826c"), // Strong
			lipgloss.Color("#28a745"), // Very strong
		},
	},
	"solarized-dark": {
		Name:       "Solarized Dark",
		Background: lipgloss.Color("#002b36"),
		Foreground: lipgloss.Color("#839496"),
		Primary:    lipgloss.Color("#268bd2"), // Blue
		Secondary:  lipgloss.Color("#2aa198"), // Cyan
		Success:    lipgloss.Color("#859900"), // Green
		Warning:    lipgloss.Color("#b58900"), // Yellow
		Error:      lipgloss.Color("#dc322f"), // Red
		Info:       lipgloss.Color("#268bd2"), // Blue
		Muted:      lipgloss.Color("#586e75"), // Base01
		Border:     lipgloss.Color("#073642"), // Base02
		Highlight:  lipgloss.Color("#93a1a1"), // Base1
		GraphUser:  lipgloss.Color("#268bd2"), // Blue for user
		GraphAsst:  lipgloss.Color("#d33682"), // Magenta for assistant
		Matrix: [5]lipgloss.Color{
			lipgloss.Color("#073642"), // Minimal
			lipgloss.Color("#586e75"), // Weak
			lipgloss.Color("#b58900"), // Medium
			lipgloss.Color("#cb4b16"), // Strong
			lipgloss.Color("#859900"), // Very strong
		},
	},
	"one-dark": {
		Name:       "One Dark",
		Background: lipgloss.Color("#282c34"),
		Foreground: lipgloss.Color("#abb2bf"),
		Primary:    lipgloss.Color("#61afef"), // Blue
		Secondary:  lipgloss.Color("#c678dd"), // Purple
		Success:    lipgloss.Color("#98c379"), // Green
		Warning:    lipgloss.Color("#e5c07b"), // Yellow
		Error:      lipgloss.Color("#e06c75"), // Red
		Info:       lipgloss.Color("#56b6c2"), // Cyan
		Muted:      lipgloss.Color("#5c6370"), // Comment
		Border:     lipgloss.Color("#3b4048"), // Guide
		Highlight:  lipgloss.Color("#d19a66"), // Orange
		GraphUser:  lipgloss.Color("#61afef"), // Blue for user
		GraphAsst:  lipgloss.Color("#c678dd"), // Purple for assistant
		Matrix: [5]lipgloss.Color{
			lipgloss.Color("#3b4048"), // Minimal
			lipgloss.Color("#5c6370"), // Weak
			lipgloss.Color("#e5c07b"), // Medium
			lipgloss.Color("#d19a66"), // Strong
			lipgloss.Color("#98c379"), // Very strong
		},
	},
}

// GetTheme returns a theme by name, defaulting to dracula if not found
func GetTheme(name string) ColorTheme {
	if theme, ok := Themes[name]; ok {
		return theme
	}
	return Themes["dracula"]
}

// ApplyTheme applies a color theme to all dashboard styles
func ApplyTheme(theme ColorTheme) *DashboardStyles {
	return &DashboardStyles{
		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(theme.Primary),

		BorderNormal: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(theme.Border),

		BorderFocused: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(theme.Highlight),

		PanelTitle: lipgloss.NewStyle().
			Bold(true).
			Foreground(theme.Secondary),

		PanelTitleFocused: lipgloss.NewStyle().
			Bold(true).
			Foreground(theme.Primary).
			Background(theme.Border),

		GraphUser: lipgloss.NewStyle().
			Foreground(theme.GraphUser),

		GraphAssistant: lipgloss.NewStyle().
			Foreground(theme.GraphAsst),

		Success: lipgloss.NewStyle().
			Foreground(theme.Success),

		Warning: lipgloss.NewStyle().
			Foreground(theme.Warning),

		Error: lipgloss.NewStyle().
			Foreground(theme.Error),

		Info: lipgloss.NewStyle().
			Foreground(theme.Info),

		Muted: lipgloss.NewStyle().
			Foreground(theme.Muted),

		EmptyState: lipgloss.NewStyle().
			Italic(true).
			Foreground(theme.Muted).
			Align(lipgloss.Center).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(theme.Border).
			Padding(2, 4).
			Margin(1, 0),

		MatrixStrength: [5]lipgloss.Style{
			lipgloss.NewStyle().Foreground(theme.Matrix[0]),
			lipgloss.NewStyle().Foreground(theme.Matrix[1]),
			lipgloss.NewStyle().Foreground(theme.Matrix[2]),
			lipgloss.NewStyle().Foreground(theme.Matrix[3]),
			lipgloss.NewStyle().Foreground(theme.Matrix[4]),
		},
	}
}

// DashboardStyles contains all styled components
type DashboardStyles struct {
	Title             lipgloss.Style
	BorderNormal      lipgloss.Style
	BorderFocused     lipgloss.Style
	PanelTitle        lipgloss.Style
	PanelTitleFocused lipgloss.Style
	GraphUser         lipgloss.Style
	GraphAssistant    lipgloss.Style
	Success           lipgloss.Style
	Warning           lipgloss.Style
	Error             lipgloss.Style
	Info              lipgloss.Style
	Muted             lipgloss.Style
	EmptyState        lipgloss.Style
	MatrixStrength    [5]lipgloss.Style
}