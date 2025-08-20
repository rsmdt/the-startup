package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// StatuslineInput represents the JSON structure from Claude Code
type StatuslineInput struct {
	HookEventName  string `json:"hook_event_name"`
	SessionID      string `json:"session_id"`
	TranscriptPath string `json:"transcript_path"`
	CWD            string `json:"cwd"`
	Model          struct {
		ID          string `json:"id"`
		DisplayName string `json:"display_name"`
	} `json:"model"`
	Workspace struct {
		CurrentDir string `json:"current_dir"`
		ProjectDir string `json:"project_dir"`
	} `json:"workspace"`
	Version     string `json:"version"`
	OutputStyle struct {
		Name string `json:"name"`
	} `json:"output_style"`
}

// NewStatuslineCommand creates the statusline command
func NewStatuslineCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "statusline",
		Short: "Generate a retro terminal status line for Claude Code",
		Long: `Generate a retro terminal status line that displays workspace info and model details.
Reads JSON input from stdin and outputs a formatted status line.

The status line shows:
- Current directory (with ~ for home)
- Git branch (if in a git repo)
- Model name and output style
- Help text

Example:
  echo '{"workspace":{"current_dir":"/path/to/project"},"model":{"display_name":"Claude"},"output_style":{"name":"default"}}' | the-startup statusline`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runStatusline(cmd.InOrStdin(), cmd.OutOrStdout())
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}
}

func runStatusline(input io.Reader, output io.Writer) error {
	// Read JSON input
	var data StatuslineInput
	if err := json.NewDecoder(input).Decode(&data); err != nil {
		return nil // Silent fail for hook compatibility
	}

	// Get terminal width
	termWidth := getTermWidth()

	// Build and output status line
	statusLine := buildStatusLine(data, termWidth)
	fmt.Fprintln(output, statusLine)
	return nil
}

func buildStatusLine(data StatuslineInput, termWidth int) string {
	// Process current directory
	currentDir := data.Workspace.CurrentDir
	if homeDir, _ := os.UserHomeDir(); homeDir != "" && strings.HasPrefix(currentDir, homeDir) {
		currentDir = "~" + strings.TrimPrefix(currentDir, homeDir)
	}

	// Build directory part with git info if available
	dirPart := fmt.Sprintf("📁 %s", currentDir)
	if gitInfo := getGitInfo(data.Workspace.CurrentDir); gitInfo != "" {
		dirPart = fmt.Sprintf("📁 %s %s", currentDir, gitInfo)
	}

	// Build parts of the status line
	parts := []string{
		dirPart,
		fmt.Sprintf("🤖 %s (%s)", data.Model.DisplayName, data.OutputStyle.Name),
	}

	// Add help text as the last element with styling
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#606060")).
		Italic(true)
	parts = append(parts, helpStyle.Render("? for shortcuts"))

	// Join all parts with lipgloss padding
	content := lipgloss.JoinHorizontal(lipgloss.Left, 
		lipgloss.NewStyle().PaddingRight(2).Render(parts[0]),
		lipgloss.NewStyle().PaddingRight(2).Render(parts[1]),
		parts[2])
	
	// Apply main style with color and max width for truncation
	mainStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FAFAFA")).
		MaxWidth(termWidth)
	
	return mainStyle.Render(content)
}

func getGitInfo(workingDir string) string {
	// Check if in git repo
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	cmd.Dir = workingDir
	if err := cmd.Run(); err != nil {
		return ""
	}

	// Get current branch
	cmd = exec.Command("git", "branch", "--show-current")
	cmd.Dir = workingDir
	output, err := cmd.Output()
	if err != nil {
		return "⎇ HEAD"
	}

	branch := strings.TrimSpace(string(output))
	if branch == "" {
		branch = "HEAD"
	}
	return fmt.Sprintf("⎇ %s", branch)
}

func getTermWidth() int {
	// Check COLUMNS env var first (most reliable in hooks/scripts)
	if cols := os.Getenv("COLUMNS"); cols != "" {
		var width int
		if _, err := fmt.Sscanf(cols, "%d", &width); err == nil && width > 0 {
			return width
		}
	}

	// Use golang.org/x/term for terminal size detection
	if width, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil && width > 0 {
		return width
	}

	if width, _, err := term.GetSize(int(os.Stderr.Fd())); err == nil && width > 0 {
		return width
	}

	// Default fallback
	return 120
}
