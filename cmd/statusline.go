package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"unsafe"

	"github.com/spf13/cobra"
)

// StatuslineInput represents the JSON structure from Claude Code
type StatuslineInput struct {
	HookEventName string `json:"hook_event_name"`
	SessionID     string `json:"session_id"`
	TranscriptPath string `json:"transcript_path"`
	CWD           string `json:"cwd"`
	Model struct {
		ID          string `json:"id"`
		DisplayName string `json:"display_name"`
	} `json:"model"`
	Workspace struct {
		CurrentDir  string `json:"current_dir"`
		ProjectDir  string `json:"project_dir"`
	} `json:"workspace"`
	Version string `json:"version"`
	OutputStyle struct {
		Name string `json:"name"`
	} `json:"output_style"`
}

// NewStatuslineCommand creates the statusline command
func NewStatuslineCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "statusline",
		Short: "Generate a retro terminal status line for Claude Code",
		Long: `Generate a retro terminal status line that displays workspace info, model details, 
and a progress bar. Reads JSON input from stdin and outputs a formatted status line.

The status line shows:
- Current directory (with ~ for home)
- Git branch (if in a git repo)
- Model name and output style
- Context window progress bar

Example:
  echo '{"workspace":{"current_dir":"/path/to/project"},"model":{"display_name":"Claude"},"output_style":{"name":"default"},"session_id":"abc123"}' | the-startup statusline`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runStatusline(cmd.InOrStdin(), cmd.OutOrStdout())
		},
		// Silent to match hook behavior
		SilenceUsage:  true,
		SilenceErrors: true,
	}
}

func runStatusline(input io.Reader, output io.Writer) error {
	// Read JSON input from stdin
	var data StatuslineInput
	decoder := json.NewDecoder(input)
	if err := decoder.Decode(&data); err != nil {
		// Return error silently like the shell script would
		return nil
	}

	// Process current directory (replace home with ~)
	currentDir := data.Workspace.CurrentDir
	homeDir, _ := os.UserHomeDir()
	if homeDir != "" && strings.HasPrefix(currentDir, homeDir) {
		currentDir = "~" + strings.TrimPrefix(currentDir, homeDir)
	}

	// Get git branch info if in a git repo
	gitInfo := getGitInfo(data.Workspace.CurrentDir)

	// Build the main content
	var mainContent string
	if gitInfo != "" {
		mainContent = fmt.Sprintf("📁 %s %s | 🤖 %s (%s)",
			currentDir,
			gitInfo,
			data.Model.DisplayName,
			data.OutputStyle.Name,
		)
	} else {
		mainContent = fmt.Sprintf("📁 %s | 🤖 %s (%s)",
			currentDir,
			data.Model.DisplayName,
			data.OutputStyle.Name,
		)
	}

	// Help text to right-align
	helpText := "? for shortcuts"

	// Get terminal width
	termWidth := getTerminalWidth()
	
	// If we couldn't detect terminal width, just add a couple spaces
	if termWidth == 0 {
		fmt.Fprintf(output, "%s  %s\n", mainContent, helpText)
		return nil
	}
	
	// Subtract 2 columns as a safety margin to prevent wrapping
	termWidth = termWidth - 4
	if termWidth < 40 {
		termWidth = 40 // Minimum reasonable width
	}
	
	// Calculate the visible length (accounting for emoji and color codes)
	mainLen := getVisibleLength(mainContent)
	helpLen := len(helpText)
	
	// Calculate padding needed (no separator, just spaces)
	// We need at least 2 spaces between main content and help text
	minSpacing := 2
	totalNeeded := mainLen + minSpacing + helpLen
	
	// If content is too long for terminal, don't try to right-align
	if totalNeeded > termWidth {
		fmt.Fprintf(output, "%s  %s\n", mainContent, helpText)
		return nil
	}
	
	// Calculate padding for right alignment
	padding := termWidth - mainLen - helpLen
	
	// Build and output the status line with right-aligned help
	fmt.Fprintf(output, "%s%s%s\n",
		mainContent,
		strings.Repeat(" ", padding),
		helpText,
	)

	return nil
}

// getGitInfo returns git branch information if in a git repository
func getGitInfo(workingDir string) string {
	// Check if we're in a git repo
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
		// Fallback to HEAD if branch command fails
		return "⎇ HEAD"
	}

	branch := strings.TrimSpace(string(output))
	if branch == "" {
		branch = "HEAD"
	}

	return fmt.Sprintf("⎇ %s", branch)
}

// getTerminalWidth returns the width of the terminal
func getTerminalWidth() int {
	const TIOCGWINSZ = 0x40087468 // macOS/Darwin value
	
	type winsize struct {
		Row    uint16
		Col    uint16
		Xpixel uint16
		Ypixel uint16
	}

	// Try all file descriptors to find one connected to a terminal
	for _, fd := range []uintptr{uintptr(syscall.Stderr), uintptr(syscall.Stdout), uintptr(syscall.Stdin)} {
		ws := &winsize{}
		retCode, _, _ := syscall.Syscall(syscall.SYS_IOCTL,
			fd,
			uintptr(TIOCGWINSZ),
			uintptr(unsafe.Pointer(ws)))
		
		if int(retCode) != -1 && ws.Col > 0 {
			return int(ws.Col)
		}
	}
	
	// Try /dev/tty directly - this often works even when stdio is redirected
	if tty, err := os.Open("/dev/tty"); err == nil {
		defer tty.Close()
		ws := &winsize{}
		retCode, _, _ := syscall.Syscall(syscall.SYS_IOCTL,
			tty.Fd(),
			uintptr(TIOCGWINSZ),
			uintptr(unsafe.Pointer(ws)))
		
		if int(retCode) != -1 && ws.Col > 0 {
			return int(ws.Col)
		}
	}
	
	// Check if COLUMNS env var is set (often set by shells)
	if cols := os.Getenv("COLUMNS"); cols != "" {
		var width int
		if _, err := fmt.Sscanf(cols, "%d", &width); err == nil && width > 0 {
			return width
		}
	}
	
	// Try using tput if available
	if output, err := exec.Command("tput", "cols").Output(); err == nil {
		var width int
		if _, err := fmt.Sscanf(strings.TrimSpace(string(output)), "%d", &width); err == nil && width > 0 {
			return width
		}
	}
	
	// Default width - don't assume, just use minimal alignment
	// We'll just put one space between main content and help text
	return 0 // Signal that we couldn't detect width
}

// getVisibleLength calculates the visible length of a string, accounting for emojis
func getVisibleLength(s string) int {
	// Count runes instead of bytes to handle multi-byte characters
	// This is a simplified version - emojis are typically 2 display columns wide
	length := 0
	for _, r := range s {
		if r >= 0x1F300 && r <= 0x1F9FF { // Emoji range (simplified)
			length += 2 // Emojis typically take 2 columns
		} else if r >= 0x2600 && r <= 0x26FF { // Miscellaneous Symbols
			length += 2
		} else if r >= 0x2700 && r <= 0x27BF { // Dingbats
			length += 2
		} else if r == '📁' || r == '🤖' || r == '⎇' { // Our specific emojis
			length += 2
		} else {
			length += 1
		}
	}
	return length
}

