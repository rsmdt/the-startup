package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rsmdt/the-startup/internal/log"
	"github.com/spf13/cobra"
)

// NewLogCommand creates the log command for hook processing and metrics analysis
func NewLogCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "log",
		Short: "Process Claude Code hook data or analyze tool usage metrics",
		Long: `Process JSON hook data from Claude Code via stdin or analyze collected metrics.

Default behavior (no subcommand):
  Processes hook data from stdin and writes to metrics log files.
  Automatically detects PreToolUse or PostToolUse events via hook_event_name field.

Examples:
  # Process hook data from Claude Code
  echo '{"hook_event_name":"PreToolUse","tool_name":"Bash"}' | the-startup log
  
  # View today's metrics summary
  the-startup log summary
  
  # View metrics for specific tools
  the-startup log tools --tool Bash --tool Read
  
  # Show errors from the last hour
  the-startup log errors --since 1h`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Default behavior: process hook from stdin
			if err := log.ProcessHook(os.Stdin); err != nil {
				// Silent failure for hook processing (exit code 0)
				return nil
			}
			return nil
		},
		// Silent usage and errors to match hook behavior
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// Add subcommands
	cmd.AddCommand(newSummaryCommand())
	cmd.AddCommand(newToolsCommand())
	cmd.AddCommand(newErrorsCommand())
	cmd.AddCommand(newTimelineCommand())

	return cmd
}

// newSummaryCommand creates the summary subcommand
func newSummaryCommand() *cobra.Command {
	var (
		since   string
		format  string
		output  string
		session string
	)

	cmd := &cobra.Command{
		Use:   "summary",
		Short: "Display metrics summary dashboard",
		Long: `Display an aggregated summary of tool usage metrics.

Shows overall statistics including:
  - Total tool invocations
  - Success/failure rates
  - Most used tools
  - Common errors
  - Activity patterns`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Parse time filter
			filter, err := parseTimeFilter(since)
			if err != nil {
				return fmt.Errorf("invalid --since value: %w", err)
			}

			// Add session filter if provided
			if session != "" {
				filter.SessionIDs = []string{session}
			}

			// Stream metrics and aggregate
			entries, err := log.StreamMetrics(*filter)
			if err != nil {
				return fmt.Errorf("failed to stream metrics: %w", err)
			}

			summary := log.AggregateMetrics(entries, filter)

			// Configure display options
			opts := log.DefaultDisplayOptions()
			opts.Format = format
			
			// Handle output file
			if output != "" && output != "-" {
				file, err := os.Create(output)
				if err != nil {
					return fmt.Errorf("failed to create output file: %w", err)
				}
				defer file.Close()
				opts.Output = file
			}

			// Display the summary
			return log.DisplaySummary(summary, opts)
		},
	}

	cmd.Flags().StringVar(&since, "since", "7d", "Time range: today, yesterday, 1h, 24h, 7d, 30d, or YYYY-MM-DD")
	cmd.Flags().StringVar(&format, "format", "terminal", "Output format: terminal, json, or csv")
	cmd.Flags().StringVar(&output, "output", "-", "Output file (- for stdout)")
	cmd.Flags().StringVar(&session, "session", "", "Filter by session ID")

	return cmd
}

// newToolsCommand creates the tools subcommand
func newToolsCommand() *cobra.Command {
	var (
		since   string
		format  string
		output  string
		tools   []string
		session string
		topN    int
	)

	cmd := &cobra.Command{
		Use:   "tools",
		Short: "Display detailed tool usage statistics",
		Long: `Display detailed statistics for tool usage.

Shows per-tool metrics including:
  - Total invocations
  - Success/failure counts
  - Average, min, max execution times
  - Error patterns`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Parse time filter
			filter, err := parseTimeFilter(since)
			if err != nil {
				return fmt.Errorf("invalid --since value: %w", err)
			}

			// Add tool filter if provided
			if len(tools) > 0 {
				filter.ToolNames = tools
			}

			// Add session filter if provided
			if session != "" {
				filter.SessionIDs = []string{session}
			}

			// Stream metrics and aggregate
			entries, err := log.StreamMetrics(*filter)
			if err != nil {
				return fmt.Errorf("failed to stream metrics: %w", err)
			}

			summary := log.AggregateMetrics(entries, filter)

			// Configure display options
			opts := log.DefaultDisplayOptions()
			opts.Format = format
			opts.TopN = topN
			
			// Handle output file
			if output != "" && output != "-" {
				file, err := os.Create(output)
				if err != nil {
					return fmt.Errorf("failed to create output file: %w", err)
				}
				defer file.Close()
				opts.Output = file
			}

			// Display the tools statistics
			return log.DisplayTools(summary, opts)
		},
	}

	cmd.Flags().StringVar(&since, "since", "7d", "Time range: today, yesterday, 1h, 24h, 7d, 30d, or YYYY-MM-DD")
	cmd.Flags().StringVar(&format, "format", "terminal", "Output format: terminal, json, or csv")
	cmd.Flags().StringVar(&output, "output", "-", "Output file (- for stdout)")
	cmd.Flags().StringArrayVar(&tools, "tool", nil, "Filter by tool name (can be specified multiple times)")
	cmd.Flags().StringVar(&session, "session", "", "Filter by session ID")
	cmd.Flags().IntVar(&topN, "top", 10, "Number of top tools to show")

	return cmd
}

// newErrorsCommand creates the errors subcommand
func newErrorsCommand() *cobra.Command {
	var (
		since   string
		format  string
		output  string
		tool    string
		session string
		topN    int
	)

	cmd := &cobra.Command{
		Use:   "errors",
		Short: "Display error analysis and patterns",
		Long: `Display analysis of tool execution errors.

Shows error patterns including:
  - Most common error types
  - Error frequency
  - Affected tools
  - Error messages`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Parse time filter
			filter, err := parseTimeFilter(since)
			if err != nil {
				return fmt.Errorf("invalid --since value: %w", err)
			}

			// Add tool filter if provided
			if tool != "" {
				filter.ToolNames = []string{tool}
			}

			// Add session filter if provided
			if session != "" {
				filter.SessionIDs = []string{session}
			}

			// Only show failures
			filter.FailuresOnly = true

			// Stream metrics and aggregate
			entries, err := log.StreamMetrics(*filter)
			if err != nil {
				return fmt.Errorf("failed to stream metrics: %w", err)
			}

			summary := log.AggregateMetrics(entries, filter)

			// Configure display options
			opts := log.DefaultDisplayOptions()
			opts.Format = format
			opts.TopN = topN
			
			// Handle output file
			if output != "" && output != "-" {
				file, err := os.Create(output)
				if err != nil {
					return fmt.Errorf("failed to create output file: %w", err)
				}
				defer file.Close()
				opts.Output = file
			}

			// Display the error analysis
			return log.DisplayErrors(summary, opts)
		},
	}

	cmd.Flags().StringVar(&since, "since", "7d", "Time range: today, yesterday, 1h, 24h, 7d, 30d, or YYYY-MM-DD")
	cmd.Flags().StringVar(&format, "format", "terminal", "Output format: terminal, json, or csv")
	cmd.Flags().StringVar(&output, "output", "-", "Output file (- for stdout)")
	cmd.Flags().StringVar(&tool, "tool", "", "Filter by specific tool")
	cmd.Flags().StringVar(&session, "session", "", "Filter by session ID")
	cmd.Flags().IntVar(&topN, "top", 10, "Number of top errors to show")

	return cmd
}

// newTimelineCommand creates the timeline subcommand
func newTimelineCommand() *cobra.Command {
	var (
		since   string
		format  string
		output  string
		session string
	)

	cmd := &cobra.Command{
		Use:   "timeline",
		Short: "Display tool usage timeline and activity patterns",
		Long: `Display hourly timeline of tool usage activity.

Shows temporal patterns including:
  - Hourly activity levels
  - Tool invocations over time
  - Success/failure trends
  - Peak usage periods`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Parse time filter
			filter, err := parseTimeFilter(since)
			if err != nil {
				return fmt.Errorf("invalid --since value: %w", err)
			}

			// Add session filter if provided
			if session != "" {
				filter.SessionIDs = []string{session}
			}

			// Stream metrics and aggregate
			entries, err := log.StreamMetrics(*filter)
			if err != nil {
				return fmt.Errorf("failed to stream metrics: %w", err)
			}

			summary := log.AggregateMetrics(entries, filter)

			// Configure display options
			opts := log.DefaultDisplayOptions()
			opts.Format = format
			
			// Handle output file
			if output != "" && output != "-" {
				file, err := os.Create(output)
				if err != nil {
					return fmt.Errorf("failed to create output file: %w", err)
				}
				defer file.Close()
				opts.Output = file
			}

			// Display the timeline
			return log.DisplayTimeline(summary, opts)
		},
	}

	cmd.Flags().StringVar(&since, "since", "7d", "Time range: today, yesterday, 1h, 24h, 7d, 30d, or YYYY-MM-DD")
	cmd.Flags().StringVar(&format, "format", "terminal", "Output format: terminal, json, or csv")
	cmd.Flags().StringVar(&output, "output", "-", "Output file (- for stdout)")
	cmd.Flags().StringVar(&session, "session", "", "Filter by session ID")

	return cmd
}

// parseTimeFilter parses a time string into a MetricsFilter
func parseTimeFilter(since string) (*log.MetricsFilter, error) {
	now := time.Now()
	var startTime, endTime time.Time

	// Default to last 7 days if no filter is provided
	if since == "" {
		since = "7d"
	}

	switch strings.ToLower(since) {
	case "today":
		startTime = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		endTime = now
	case "yesterday":
		yesterday := now.AddDate(0, 0, -1)
		startTime = time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, yesterday.Location())
		endTime = time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 23, 59, 59, 999999999, yesterday.Location())
	case "1h":
		startTime = now.Add(-1 * time.Hour)
		endTime = now
	case "24h":
		startTime = now.Add(-24 * time.Hour)
		endTime = now
	case "7d":
		startTime = now.AddDate(0, 0, -7)
		endTime = now
	case "30d":
		startTime = now.AddDate(0, 0, -30)
		endTime = now
	default:
		// Try parsing as date (YYYY-MM-DD)
		if t, err := time.Parse("2006-01-02", since); err == nil {
			startTime = t
			endTime = t.Add(24*time.Hour).Add(-time.Nanosecond)
		} else {
			// Try parsing as duration (e.g., "3h", "2d")
			if strings.HasSuffix(since, "h") {
				if hours, err := time.ParseDuration(since); err == nil {
					startTime = now.Add(-hours)
					endTime = now
				} else {
					return nil, fmt.Errorf("invalid duration: %s", since)
				}
			} else if strings.HasSuffix(since, "d") {
				days := strings.TrimSuffix(since, "d")
				var d int
				if _, err := fmt.Sscanf(days, "%d", &d); err == nil {
					startTime = now.AddDate(0, 0, -d)
					endTime = now
				} else {
					return nil, fmt.Errorf("invalid duration: %s", since)
				}
			} else {
				return nil, fmt.Errorf("unrecognized time format: %s", since)
			}
		}
	}

	return &log.MetricsFilter{
		StartDate: startTime,
		EndDate:   endTime,
	}, nil
}