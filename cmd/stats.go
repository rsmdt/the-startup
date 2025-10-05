package cmd

import (
	"fmt"

	"github.com/charmbracelet/bubbletea"
	"github.com/rsmdt/the-startup/internal/ui/agent"
	"github.com/spf13/cobra"
)

// NewStatsCommand creates the stats command for launching the interactive dashboard
func NewStatsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stats",
		Short: "Launch interactive agent analytics dashboard",
		Long: `Launch an interactive dashboard to analyze Claude Code usage statistics.

The dashboard provides real-time analytics including:
  • Agent usage leaderboard with performance metrics
  • Agent delegation patterns and workflows
  • Co-occurrence analysis of agents used together
  • Interactive navigation and data exploration

Examples:
  # Launch dashboard for current project
  the-startup stats

  # Dashboard shows current project statistics by default
  # Use arrow keys or j/k to navigate, Tab to switch panels`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInteractiveDashboard()
		},
	}

	return cmd
}

// runInteractiveDashboard launches the interactive TUI dashboard
func runInteractiveDashboard() error {
	// Create dashboard with dynamic data loading
	// The dashboard will load data itself based on filters
	model := agent.NewDashboardModelWithLoading()

	program := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		return fmt.Errorf("failed to run interactive dashboard: %w", err)
	}

	return nil
}

