package ui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Spinner frames for animation
var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// ProgressModel represents a progress indicator
type ProgressModel struct {
	spinner  int
	progress float64
	total    float64
	message  string
	done     bool
	styles   Styles
}

// TickMsg is sent to animate the spinner
type TickMsg time.Time

// ProgressMsg updates the progress
type ProgressMsg struct {
	Current float64
	Total   float64
	Message string
}

// DoneMsg signals completion
type DoneMsg struct {
	Success bool
	Message string
}

// NewProgressModel creates a new progress indicator
func NewProgressModel(message string) ProgressModel {
	return ProgressModel{
		message: message,
		styles:  GetStyles(),
	}
}

func (m ProgressModel) Init() tea.Cmd {
	return tickCmd()
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (m ProgressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case TickMsg:
		m.spinner++
		if m.spinner >= len(spinnerFrames) {
			m.spinner = 0
		}
		return m, tickCmd()

	case ProgressMsg:
		m.progress = msg.Current
		m.total = msg.Total
		if msg.Message != "" {
			m.message = msg.Message
		}
		return m, nil

	case DoneMsg:
		m.done = true
		m.message = msg.Message
		return m, tea.Quit
	}

	return m, nil
}

func (m ProgressModel) View() string {
	if m.done {
		return ""
	}

	// Spinner
	spinner := m.styles.Spinner.Render(spinnerFrames[m.spinner])

	// Progress bar (if we have total)
	progressBar := ""
	if m.total > 0 {
		percentage := m.progress / m.total
		barWidth := 30
		filled := int(percentage * float64(barWidth))
		empty := barWidth - filled

		filledBar := strings.Repeat("█", filled)
		emptyBar := strings.Repeat("░", empty)

		progressBar = fmt.Sprintf(" [%s%s] %.0f%%",
			m.styles.ProgressFilled.Render(filledBar),
			m.styles.ProgressEmpty.Render(emptyBar),
			percentage*100,
		)
	}

	return fmt.Sprintf("%s %s%s", spinner, m.message, progressBar)
}

// SimpleSpinner creates a simple spinner without progress tracking
type SimpleSpinner struct {
	frame   int
	message string
	styles  Styles
}

// NewSimpleSpinner creates a new simple spinner
func NewSimpleSpinner(message string) *SimpleSpinner {
	return &SimpleSpinner{
		message: message,
		styles:  GetStyles(),
	}
}

// Next advances the spinner animation
func (s *SimpleSpinner) Next() string {
	s.frame++
	if s.frame >= len(spinnerFrames) {
		s.frame = 0
	}
	return s.Render()
}

// Render returns the current spinner frame with message
func (s *SimpleSpinner) Render() string {
	spinner := s.styles.Spinner.Render(spinnerFrames[s.frame])
	return fmt.Sprintf("%s %s", spinner, s.message)
}

// SetMessage updates the spinner message
func (s *SimpleSpinner) SetMessage(message string) {
	s.message = message
}

// Success shows a success message
func (s *SimpleSpinner) Success(message string) string {
	return s.styles.Success.Render(fmt.Sprintf("%s %s", IconSuccess, message))
}

// Error shows an error message
func (s *SimpleSpinner) Error(message string) string {
	return s.styles.Error.Render(fmt.Sprintf("%s %s", IconError, message))
}

// ProgressBar creates a simple progress bar
func ProgressBar(current, total int, width int) string {
	styles := GetStyles()
	if total == 0 {
		return ""
	}

	percentage := float64(current) / float64(total)
	filled := int(percentage * float64(width))
	empty := width - filled

	filledBar := strings.Repeat("█", filled)
	emptyBar := strings.Repeat("░", empty)

	return fmt.Sprintf("[%s%s] %d/%d",
		styles.ProgressFilled.Render(filledBar),
		styles.ProgressEmpty.Render(emptyBar),
		current,
		total,
	)
}

// InstallProgress tracks installation progress
type InstallProgress struct {
	spinner      *SimpleSpinner
	currentItem  string
	itemsTotal   int
	itemsCurrent int
	styles       Styles
}

// NewInstallProgress creates a new installation progress tracker
func NewInstallProgress(total int) *InstallProgress {
	return &InstallProgress{
		spinner:    NewSimpleSpinner("Preparing installation..."),
		itemsTotal: total,
		styles:     GetStyles(),
	}
}

// StartItem begins processing a new item
func (ip *InstallProgress) StartItem(name string) {
	ip.itemsCurrent++
	ip.currentItem = name
	message := fmt.Sprintf("Installing %s... %s", name, ProgressBar(ip.itemsCurrent, ip.itemsTotal, 20))
	ip.spinner.SetMessage(message)
}

// CompleteItem marks an item as complete
func (ip *InstallProgress) CompleteItem(name string) string {
	return ip.styles.Success.Render(fmt.Sprintf("  %s %s installed", IconSuccess, name))
}

// Error shows an error for an item
func (ip *InstallProgress) Error(name string, err error) string {
	return ip.styles.Error.Render(fmt.Sprintf("  %s %s: %v", IconError, name, err))
}

// Complete shows the final completion message
func (ip *InstallProgress) Complete() string {
	return ip.styles.Success.Render(fmt.Sprintf("%s All components installed successfully!", IconSuccess))
}

// RenderSpinner returns the current spinner state
func (ip *InstallProgress) RenderSpinner() string {
	return ip.spinner.Next()
}