package cmd

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/rsmdt/the-startup/internal/stats"
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

// extractAgentAnalytics processes log entries to extract agent statistics and patterns
func extractAgentAnalytics(entries []stats.ClaudeLogEntry) (
	map[string]*stats.GlobalAgentStats,
	[]stats.DelegationPattern,
	[]stats.AgentCoOccurrence,
	error,
) {
	// Create analyzers
	agentAnalyzer := stats.NewAgentAnalyzer()
	delegationAnalyzer := stats.NewDelegationAnalyzer()

	// Process entries to detect agents and build patterns
	sessionAgents := make(map[string][]string)

	for _, entry := range entries {
		// Detect agent in the entry using existing function
		agentType, confidence := stats.DetectAgent(entry)

		if agentType != "" && confidence > 0.5 { // Only process high confidence detections
			// Create agent invocation
			inv := stats.AgentInvocation{
				AgentType:       agentType,
				SessionID:       entry.SessionID,
				Timestamp:       entry.Timestamp,
				DurationMs:      1000, // Default duration, could be extracted from logs
				Success:         true, // Simplified, could be improved with error detection
				Confidence:      confidence,
				DetectionMethod: "pattern_detection",
			}

			// Process with agent analyzer
			agentAnalyzer.ProcessAgentInvocation(inv)

			// Track session agents for delegation analysis
			if sessionAgents[entry.SessionID] == nil {
				sessionAgents[entry.SessionID] = []string{}
			}
			sessionAgents[entry.SessionID] = append(sessionAgents[entry.SessionID], agentType)
		}
	}

	// Build delegation patterns from session agents
	for sessionID, agents := range sessionAgents {
		// Process agent transitions
		for i := 0; i < len(agents)-1; i++ {
			delegationAnalyzer.ProcessAgentTransition(
				sessionID,
				agents[i],
				agents[i+1],
				time.Now(), // Simplified timestamp
			)
		}

		// Process co-occurrence patterns manually by creating transitions for all pairs
		for i := 0; i < len(agents); i++ {
			for j := i + 1; j < len(agents); j++ {
				// Create bidirectional transitions for co-occurrence
				delegationAnalyzer.ProcessAgentTransition(sessionID, agents[i], agents[j], time.Now())
				delegationAnalyzer.ProcessAgentTransition(sessionID, agents[j], agents[i], time.Now())
			}
		}
	}

	// Get results
	agentStats, err := agentAnalyzer.GetAllAgentStats()
	if err != nil {
		return nil, nil, nil, err
	}

	delegationPatterns, err := delegationAnalyzer.GetTransitionPatterns()
	if err != nil {
		return nil, nil, nil, err
	}

	coOccurrences, err := delegationAnalyzer.GetCoOccurrencePatterns()
	if err != nil {
		return nil, nil, nil, err
	}

	return agentStats, delegationPatterns, coOccurrences, nil
}

// Helper function to get project name from path
func getProjectName(projectPath string) string {
	if projectPath == "" {
		return "unknown"
	}
	return filepath.Base(projectPath)
}