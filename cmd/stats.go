package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rsmdt/the-startup/internal/stats"
	"github.com/spf13/cobra"
)

// NewStatsCommand creates the stats command for analyzing Claude Code usage
func NewStatsCommand() *cobra.Command {
	var (
		global bool
		since  string
		format string
	)

	cmd := &cobra.Command{
		Use:   "stats",
		Short: "Analyze Claude Code usage statistics",
		Long: `Analyze and display statistics from Claude Code usage logs.

By default, shows statistics for the current project. Use -g flag for global statistics
across all projects.

Examples:
  # Show current project stats
  the-startup stats
  
  # Show global stats across all projects
  the-startup stats -g
  
  # Show stats for the last 24 hours
  the-startup stats --since 24h
  
  # Output as JSON
  the-startup stats --format json
  
  # Show tool-specific statistics
  the-startup stats tools
  
  # Show agent usage
  the-startup stats agents`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Run default stats (summary view)
			return runStats(global, since, format, "summary")
		},
	}

	// Add persistent flags (available to all subcommands)
	cmd.PersistentFlags().BoolVarP(&global, "global", "g", false, "Show global statistics across all projects")
	cmd.PersistentFlags().StringVar(&since, "since", "", "Filter by time: 1h, 24h, 7d, 30d, or YYYY-MM-DD")
	cmd.PersistentFlags().StringVar(&format, "format", "table", "Output format: table, json, or csv")

	// Add subcommands
	cmd.AddCommand(newToolsStatsCommand(&global, &since, &format))
	cmd.AddCommand(newAgentsStatsCommand(&global, &since, &format))
	cmd.AddCommand(newCommandsStatsCommand(&global, &since, &format))
	cmd.AddCommand(newSessionsStatsCommand(&global, &since, &format))

	return cmd
}

// newToolsStatsCommand creates the tools subcommand
func newToolsStatsCommand(global *bool, since *string, format *string) *cobra.Command {
	return &cobra.Command{
		Use:   "tools",
		Short: "Display tool usage statistics",
		Long: `Display detailed statistics for tool usage.

Shows per-tool metrics including:
  - Total invocations
  - Success/failure rates
  - Execution times (min, max, average)
  - Error patterns`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runStats(*global, *since, *format, "tools")
		},
	}
}

// newAgentsStatsCommand creates the agents subcommand
func newAgentsStatsCommand(global *bool, since *string, format *string) *cobra.Command {
	return &cobra.Command{
		Use:   "agents",
		Short: "Display agent usage statistics",
		Long: `Display statistics for agent invocations.

Shows agent-specific metrics including:
  - Invocation counts by agent type
  - Success/failure rates
  - Most frequently used agents
  - Agent execution patterns`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runStats(*global, *since, *format, "agents")
		},
	}
}

// newCommandsStatsCommand creates the commands subcommand
func newCommandsStatsCommand(global *bool, since *string, format *string) *cobra.Command {
	return &cobra.Command{
		Use:   "commands",
		Short: "Display command usage statistics",
		Long: `Display statistics for command executions.

Shows command-specific metrics including:
  - Command invocation counts
  - Success/failure rates
  - Most frequently used commands
  - Command execution patterns`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runStats(*global, *since, *format, "commands")
		},
	}
}

// newSessionsStatsCommand creates the sessions subcommand
func newSessionsStatsCommand(global *bool, since *string, format *string) *cobra.Command {
	return &cobra.Command{
		Use:   "sessions",
		Short: "Display session statistics",
		Long: `Display statistics grouped by Claude Code sessions.

Shows session-level metrics including:
  - Session durations
  - Tools used per session
  - Session activity patterns
  - Most active sessions`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runStats(*global, *since, *format, "sessions")
		},
	}
}

// runStats executes the stats pipeline
func runStats(global bool, since string, format string, viewType string) error {
	// Step 1: Discovery - find log files
	discovery := stats.NewLogDiscovery()
	
	// Parse time filter
	filter := stats.FilterOptions{}
	if since != "" {
		t, err := parseTimeFilter(since)
		if err != nil {
			return fmt.Errorf("invalid --since value: %w", err)
		}
		filter.StartTime = &t
	}

	var logFiles []string
	var err error

	if global {
		// For global stats, we need to find all projects
		// This would require enumerating all project directories
		// For now, we'll return an error
		return fmt.Errorf("global stats not yet implemented - please run from a project directory")
	} else {
		// Get current project logs
		projectPath := discovery.GetCurrentProject()
		if projectPath == "" {
			// Try to use the current directory as project
			// The discovery component will handle path sanitization
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}
			projectPath = cwd
		}
		
		logFiles, err = discovery.FindLogFiles(projectPath, filter)
		if err != nil {
			return fmt.Errorf("failed to discover log files: %w", err)
		}
	}

	if len(logFiles) == 0 {
		scope := "current project"
		if global {
			scope = "all projects"
		}
		fmt.Printf("No log files found in %s.\n", scope)
		fmt.Println("\nMake sure you've used Claude Code in this project.")
		return nil
	}


	// Step 2: Stream and parse log entries
	parser := stats.NewJSONLParser()

	// Aggregate statistics using the streaming approach
	aggregator := stats.NewAggregator()
	entriesProcessed := 0
	var entries []stats.ClaudeLogEntry // Store entries for command/agent extraction

	for _, logFile := range logFiles {
		// Parse the file and get channels
		entryChan, errorChan := parser.ParseFile(logFile)
		
		// Process entries from the channel
		for entry := range entryChan {
			entries = append(entries, entry) // Store for later analysis
			// Process the entry through the aggregator
			if err := aggregator.ProcessEntry(entry); err != nil {
				// Skip entries with processing errors
				continue
			}
			entriesProcessed++
		}
		
		// Drain error channel (we're ignoring errors for now)
		for range errorChan {
			// Errors are logged but don't stop processing
		}
	}

	if entriesProcessed == 0 {
		fmt.Println("No log entries found matching the criteria.")
		return nil
	}

	// Get the aggregated statistics
	statistics, err := aggregator.GetOverallStats()
	if err != nil {
		return fmt.Errorf("failed to get aggregated statistics: %w", err)
	}

	// Step 5: Display results
	display := stats.NewDisplayFormatter()
	display.SetOutputFormat(format)

	var output string

	switch viewType {
	case "summary":
		output = display.FormatStats(statistics)
	case "tools":
		// Get all tool stats and format as leaderboard
		toolStatsMap := make(map[string]*stats.GlobalToolStats)
		// Use the GlobalToolStats from the aggregated statistics
		if statistics != nil && statistics.GlobalToolStats != nil {
			toolStatsMap = statistics.GlobalToolStats
		}
		output = display.FormatToolLeaderboard(toolStatsMap, 20)
	case "agents":
		// Extract agent stats from tools (Task tool with subagent_type)
		agentStats := make(map[string]int)
		for _, entry := range entries {
			if entry.Type == "assistant" {
				var msg stats.AssistantMessage
				if err := json.Unmarshal(entry.Content, &msg); err == nil {
					for _, tool := range msg.ToolUses {
						if tool.Name == "Task" {
							// Parse the parameters to extract subagent_type
							var params map[string]interface{}
							if err := json.Unmarshal(tool.Parameters, &params); err == nil {
								if subagent, ok := params["subagent_type"].(string); ok {
									agentStats[subagent]++
								}
							}
						}
					}
				}
			}
		}
		output = display.FormatCommandLeaderboard(agentStats)
	case "commands":
		// Extract command stats from the log entries
		commandStats := make(map[string]int)
		// For now, parse entries to extract commands
		// In the future, this should come from the aggregator
		output = display.FormatCommandLeaderboard(commandStats)
	case "sessions":
		// Format individual session statistics
		// For now, show overall stats as session view
		output = display.FormatStats(statistics)
	default:
		return fmt.Errorf("unknown view type: %s", viewType)
	}

	fmt.Print(output)
	return nil
}


// parseTimeFilter converts a time string to a time.Time
func parseTimeFilter(since string) (time.Time, error) {
	now := time.Now()

	// Handle special keywords
	switch strings.ToLower(since) {
	case "today":
		return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()), nil
	case "yesterday":
		yesterday := now.AddDate(0, 0, -1)
		return time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, yesterday.Location()), nil
	}

	// Handle duration formats (1h, 24h, 7d, 30d)
	if strings.HasSuffix(since, "h") {
		hours := strings.TrimSuffix(since, "h")
		if duration, err := time.ParseDuration(hours + "h"); err == nil {
			return now.Add(-duration), nil
		}
	}

	if strings.HasSuffix(since, "d") {
		days := strings.TrimSuffix(since, "d")
		var d int
		if _, err := fmt.Sscanf(days, "%d", &d); err == nil {
			return now.AddDate(0, 0, -d), nil
		}
	}

	// Try parsing as date (YYYY-MM-DD)
	if t, err := time.Parse("2006-01-02", since); err == nil {
		return t, nil
	}

	// Try parsing as datetime (YYYY-MM-DD HH:MM:SS)
	if t, err := time.Parse("2006-01-02 15:04:05", since); err == nil {
		return t, nil
	}

	return time.Time{}, fmt.Errorf("unrecognized time format: %s (use 1h, 24h, 7d, 30d, or YYYY-MM-DD)", since)
}

// Helper function to get project name from path
func getProjectName(projectPath string) string {
	if projectPath == "" {
		return "unknown"
	}
	return filepath.Base(projectPath)
}