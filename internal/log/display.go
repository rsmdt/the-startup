package log

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
	"time"
)

// DisplayOptions configures how metrics are displayed
type DisplayOptions struct {
	Format      string    // "terminal", "json", or "csv"
	Output      io.Writer // Where to write output (default: os.Stdout)
	TopN        int       // Number of top items to show (default: 10)
	ShowDetails bool      // Whether to show detailed information
	NoColor     bool      // Disable terminal colors
}

// DefaultDisplayOptions returns default display settings
func DefaultDisplayOptions() *DisplayOptions {
	return &DisplayOptions{
		Format:      "terminal",
		Output:      os.Stdout,
		TopN:        10,
		ShowDetails: false,
		NoColor:     false,
	}
}

// DisplaySummary shows the main metrics dashboard
func DisplaySummary(summary *MetricsSummary, opts *DisplayOptions) error {
	if opts == nil {
		opts = DefaultDisplayOptions()
	}

	switch opts.Format {
	case "json":
		return displaySummaryJSON(summary, opts.Output)
	case "csv":
		return displaySummaryCSV(summary, opts.Output)
	default:
		return displaySummaryTerminal(summary, opts)
	}
}

// displaySummaryTerminal renders the summary in terminal format
func displaySummaryTerminal(summary *MetricsSummary, opts *DisplayOptions) error {
	w := opts.Output

	// Header with box drawing
	fmt.Fprintln(w, "╔══════════════════════════════════════════════════════════════════════════╗")
	fmt.Fprintf(w, "║                    Metrics Summary - %s                    ║\n", 
		centerString(summary.Period.String(), 30))
	fmt.Fprintln(w, "╚══════════════════════════════════════════════════════════════════════════╝")
	fmt.Fprintln(w)

	// Overview section
	fmt.Fprintln(w, "📊 Overview")
	fmt.Fprintln(w, "───────────")
	fmt.Fprintf(w, "  Total Invocations:     %d\n", summary.TotalCalls)
	fmt.Fprintf(w, "  Unique Sessions:       %d\n", summary.UniqueSessions)
	fmt.Fprintf(w, "  Overall Success Rate:  %.1f%%\n", summary.SuccessRate)
	fmt.Fprintln(w)

	// Top tools table
	fmt.Fprintln(w, "🔧 Top Tools by Usage")
	fmt.Fprintln(w, "────────────────────")
	
	if err := displayToolsTable(summary.ToolStats, opts); err != nil {
		return err
	}

	// Top errors if present
	if len(summary.TopErrors) > 0 {
		fmt.Fprintln(w)
		fmt.Fprintln(w, "⚠️  Most Common Errors")
		fmt.Fprintln(w, "─────────────────────")
		
		if err := displayErrorsTable(summary.TopErrors, opts); err != nil {
			return err
		}
	}

	// Hourly activity summary if available
	if len(summary.HourlyActivity) > 0 && opts.ShowDetails {
		fmt.Fprintln(w)
		fmt.Fprintln(w, "📈 Activity Timeline")
		fmt.Fprintln(w, "───────────────────")
		
		if err := displayHourlyActivity(summary.HourlyActivity, opts); err != nil {
			return err
		}
	}

	return nil
}

// displaySummaryJSON outputs the summary as JSON
func displaySummaryJSON(summary *MetricsSummary, w io.Writer) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(summary)
}

// displaySummaryCSV outputs the summary as CSV
func displaySummaryCSV(summary *MetricsSummary, w io.Writer) error {
	csvWriter := csv.NewWriter(w)
	defer csvWriter.Flush()

	// Write header
	headers := []string{
		"Tool", "Total Calls", "Success Count", "Failure Count",
		"Success Rate (%)", "Avg Duration (ms)", "Min Duration (ms)", 
		"Max Duration (ms)", "Total Duration (ms)", "Last Used",
	}
	if err := csvWriter.Write(headers); err != nil {
		return err
	}

	// Sort tools by total calls
	tools := getTopToolsList(summary.ToolStats, len(summary.ToolStats))

	// Write tool data
	for _, stats := range tools {
		successRate := 0.0
		if stats.TotalCalls > 0 {
			successRate = float64(stats.SuccessCount) / float64(stats.TotalCalls) * 100
		}

		row := []string{
			stats.Name,
			fmt.Sprintf("%d", stats.TotalCalls),
			fmt.Sprintf("%d", stats.SuccessCount),
			fmt.Sprintf("%d", stats.FailureCount),
			fmt.Sprintf("%.1f", successRate),
			fmt.Sprintf("%.1f", stats.AvgDurationMs),
			fmt.Sprintf("%d", stats.MinDurationMs),
			fmt.Sprintf("%d", stats.MaxDurationMs),
			fmt.Sprintf("%d", stats.TotalDurationMs),
			stats.LastUsed.Format(time.RFC3339),
		}
		if err := csvWriter.Write(row); err != nil {
			return err
		}
	}

	return nil
}

// DisplayTools shows detailed tool statistics
func DisplayTools(summary *MetricsSummary, opts *DisplayOptions) error {
	if opts == nil {
		opts = DefaultDisplayOptions()
	}

	switch opts.Format {
	case "json":
		return displayToolsJSON(summary.ToolStats, opts.Output)
	case "csv":
		return displayToolsCSV(summary.ToolStats, opts.Output)
	default:
		return displayToolsTerminal(summary.ToolStats, opts)
	}
}

// displayToolsTable renders tool statistics as a table
func displayToolsTable(toolStats map[string]*ToolStatistics, opts *DisplayOptions) error {
	tools := getTopToolsList(toolStats, opts.TopN)
	
	// Create tabwriter for aligned columns
	tw := tabwriter.NewWriter(opts.Output, 0, 0, 2, ' ', 0)
	defer tw.Flush()

	// Table header
	fmt.Fprintln(tw, "Tool\tCalls\tSuccess\tAvg Time\tMin/Max\tErrors")
	fmt.Fprintln(tw, "────\t─────\t───────\t────────\t───────\t──────")

	// Table rows
	for _, stats := range tools {
		successRate := 0.0
		if stats.TotalCalls > 0 {
			successRate = float64(stats.SuccessCount) / float64(stats.TotalCalls) * 100
		}

		// Format duration times
		avgTime := formatDuration(stats.AvgDurationMs)
		minMax := fmt.Sprintf("%s/%s", 
			formatDuration(float64(stats.MinDurationMs)),
			formatDuration(float64(stats.MaxDurationMs)))

		fmt.Fprintf(tw, "%s\t%d\t%.1f%%\t%s\t%s\t%d\n",
			truncateString(stats.Name, 20),
			stats.TotalCalls,
			successRate,
			avgTime,
			minMax,
			stats.FailureCount)
	}

	return nil
}

// displayToolsTerminal shows detailed tool statistics in terminal
func displayToolsTerminal(toolStats map[string]*ToolStatistics, opts *DisplayOptions) error {
	w := opts.Output

	fmt.Fprintln(w, "╔══════════════════════════════════════════════════════════════════════════╗")
	fmt.Fprintln(w, "║                          Tool Statistics                                 ║")
	fmt.Fprintln(w, "╚══════════════════════════════════════════════════════════════════════════╝")
	fmt.Fprintln(w)

	return displayToolsTable(toolStats, opts)
}

// displayToolsJSON outputs tool statistics as JSON
func displayToolsJSON(toolStats map[string]*ToolStatistics, w io.Writer) error {
	// Convert map to slice for better JSON output
	tools := make([]*ToolStatistics, 0, len(toolStats))
	for _, stats := range toolStats {
		tools = append(tools, stats)
	}

	// Sort by total calls
	sort.Slice(tools, func(i, j int) bool {
		return tools[i].TotalCalls > tools[j].TotalCalls
	})

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(tools)
}

// displayToolsCSV outputs tool statistics as CSV
func displayToolsCSV(toolStats map[string]*ToolStatistics, w io.Writer) error {
	csvWriter := csv.NewWriter(w)
	defer csvWriter.Flush()

	// Write header
	headers := []string{
		"Tool", "Total Calls", "Success Count", "Failure Count",
		"Avg Duration (ms)", "Min Duration (ms)", "Max Duration (ms)",
		"Error Types", "Last Used",
	}
	if err := csvWriter.Write(headers); err != nil {
		return err
	}

	// Sort tools
	tools := getTopToolsList(toolStats, len(toolStats))

	// Write data
	for _, stats := range tools {
		errorTypes := make([]string, 0, len(stats.ErrorTypes))
		for errType, count := range stats.ErrorTypes {
			errorTypes = append(errorTypes, fmt.Sprintf("%s(%d)", errType, count))
		}

		row := []string{
			stats.Name,
			fmt.Sprintf("%d", stats.TotalCalls),
			fmt.Sprintf("%d", stats.SuccessCount),
			fmt.Sprintf("%d", stats.FailureCount),
			fmt.Sprintf("%.1f", stats.AvgDurationMs),
			fmt.Sprintf("%d", stats.MinDurationMs),
			fmt.Sprintf("%d", stats.MaxDurationMs),
			strings.Join(errorTypes, "; "),
			stats.LastUsed.Format(time.RFC3339),
		}
		if err := csvWriter.Write(row); err != nil {
			return err
		}
	}

	return nil
}

// DisplayTimeline shows activity over time
func DisplayTimeline(summary *MetricsSummary, opts *DisplayOptions) error {
	if opts == nil {
		opts = DefaultDisplayOptions()
	}

	switch opts.Format {
	case "json":
		return displayTimelineJSON(summary.HourlyActivity, opts.Output)
	case "csv":
		return displayTimelineCSV(summary.HourlyActivity, opts.Output)
	default:
		return displayTimelineTerminal(summary.HourlyActivity, opts)
	}
}

// displayTimelineTerminal shows timeline in terminal
func displayTimelineTerminal(hourlyActivity []HourlyStats, opts *DisplayOptions) error {
	w := opts.Output

	fmt.Fprintln(w, "╔══════════════════════════════════════════════════════════════════════════╗")
	fmt.Fprintln(w, "║                         Activity Timeline                                ║")
	fmt.Fprintln(w, "╚══════════════════════════════════════════════════════════════════════════╝")
	fmt.Fprintln(w)

	return displayHourlyActivity(hourlyActivity, opts)
}

// displayHourlyActivity renders hourly statistics
func displayHourlyActivity(hourlyActivity []HourlyStats, opts *DisplayOptions) error {
	tw := tabwriter.NewWriter(opts.Output, 0, 0, 2, ' ', 0)
	defer tw.Flush()

	// Find max calls for bar chart scaling
	maxCalls := 0
	for _, hour := range hourlyActivity {
		if hour.TotalCalls > maxCalls {
			maxCalls = hour.TotalCalls
		}
	}

	fmt.Fprintln(tw, "Hour\tCalls\tSuccess\tFailure\tTools\tActivity")
	fmt.Fprintln(tw, "────\t─────\t───────\t───────\t─────\t────────")

	for _, hour := range hourlyActivity {
		successRate := 0.0
		if hour.TotalCalls > 0 {
			successRate = float64(hour.SuccessCount) / float64(hour.TotalCalls) * 100
		}

		// Create activity bar
		barWidth := 20
		if maxCalls > 0 {
			barWidth = (hour.TotalCalls * 20) / maxCalls
		}
		bar := strings.Repeat("█", barWidth)

		fmt.Fprintf(tw, "%s\t%d\t%.1f%%\t%d\t%d\t%s\n",
			hour.Hour.Format("15:04"),
			hour.TotalCalls,
			successRate,
			hour.FailureCount,
			hour.UniqueTools,
			bar)
	}

	return nil
}

// displayTimelineJSON outputs timeline as JSON
func displayTimelineJSON(hourlyActivity []HourlyStats, w io.Writer) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(hourlyActivity)
}

// displayTimelineCSV outputs timeline as CSV
func displayTimelineCSV(hourlyActivity []HourlyStats, w io.Writer) error {
	csvWriter := csv.NewWriter(w)
	defer csvWriter.Flush()

	// Write header
	headers := []string{
		"Hour", "Total Calls", "Success Count", "Failure Count",
		"Success Rate (%)", "Unique Tools",
	}
	if err := csvWriter.Write(headers); err != nil {
		return err
	}

	// Write data
	for _, hour := range hourlyActivity {
		successRate := 0.0
		if hour.TotalCalls > 0 {
			successRate = float64(hour.SuccessCount) / float64(hour.TotalCalls) * 100
		}

		row := []string{
			hour.Hour.Format(time.RFC3339),
			fmt.Sprintf("%d", hour.TotalCalls),
			fmt.Sprintf("%d", hour.SuccessCount),
			fmt.Sprintf("%d", hour.FailureCount),
			fmt.Sprintf("%.1f", successRate),
			fmt.Sprintf("%d", hour.UniqueTools),
		}
		if err := csvWriter.Write(row); err != nil {
			return err
		}
	}

	return nil
}

// DisplayErrors shows error analysis
func DisplayErrors(summary *MetricsSummary, opts *DisplayOptions) error {
	if opts == nil {
		opts = DefaultDisplayOptions()
	}

	switch opts.Format {
	case "json":
		return displayErrorsJSON(summary.TopErrors, opts.Output)
	case "csv":
		return displayErrorsCSV(summary.TopErrors, opts.Output)
	default:
		return displayErrorsTerminal(summary.TopErrors, opts)
	}
}

// displayErrorsTerminal shows errors in terminal
func displayErrorsTerminal(topErrors []ErrorPattern, opts *DisplayOptions) error {
	w := opts.Output

	fmt.Fprintln(w, "╔══════════════════════════════════════════════════════════════════════════╗")
	fmt.Fprintln(w, "║                           Error Analysis                                 ║")
	fmt.Fprintln(w, "╚══════════════════════════════════════════════════════════════════════════╝")
	fmt.Fprintln(w)

	if len(topErrors) == 0 {
		fmt.Fprintln(w, "  No errors found in the selected time period.")
		return nil
	}

	return displayErrorsTable(topErrors, opts)
}

// displayErrorsTable renders error patterns as a table
func displayErrorsTable(topErrors []ErrorPattern, opts *DisplayOptions) error {
	tw := tabwriter.NewWriter(opts.Output, 0, 0, 2, ' ', 0)
	defer tw.Flush()

	fmt.Fprintln(tw, "Type\tMessage\tCount\tTools")
	fmt.Fprintln(tw, "────\t───────\t─────\t─────")

	limit := opts.TopN
	if limit > len(topErrors) || limit <= 0 {
		limit = len(topErrors)
	}

	for i := 0; i < limit; i++ {
		err := topErrors[i]
		tools := strings.Join(err.Tools, ", ")
		if len(tools) > 30 {
			tools = tools[:27] + "..."
		}

		message := truncateString(err.ErrorMessage, 40)

		fmt.Fprintf(tw, "%s\t%s\t%d\t%s\n",
			truncateString(err.ErrorType, 15),
			message,
			err.Count,
			tools)
	}

	return nil
}

// displayErrorsJSON outputs errors as JSON
func displayErrorsJSON(topErrors []ErrorPattern, w io.Writer) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(topErrors)
}

// displayErrorsCSV outputs errors as CSV
func displayErrorsCSV(topErrors []ErrorPattern, w io.Writer) error {
	csvWriter := csv.NewWriter(w)
	defer csvWriter.Flush()

	// Write header
	headers := []string{"Error Type", "Error Message", "Count", "Affected Tools"}
	if err := csvWriter.Write(headers); err != nil {
		return err
	}

	// Write data
	for _, err := range topErrors {
		row := []string{
			err.ErrorType,
			err.ErrorMessage,
			fmt.Sprintf("%d", err.Count),
			strings.Join(err.Tools, "; "),
		}
		if err := csvWriter.Write(row); err != nil {
			return err
		}
	}

	return nil
}

// Helper functions

// getTopToolsList returns a sorted list of top N tools by total calls
func getTopToolsList(toolStats map[string]*ToolStatistics, n int) []*ToolStatistics {
	tools := make([]*ToolStatistics, 0, len(toolStats))
	for _, stats := range toolStats {
		tools = append(tools, stats)
	}

	// Sort by total calls descending
	sort.Slice(tools, func(i, j int) bool {
		return tools[i].TotalCalls > tools[j].TotalCalls
	})

	// Return top N
	if n > 0 && n < len(tools) {
		tools = tools[:n]
	}

	return tools
}

// formatDuration formats milliseconds as a human-readable duration
func formatDuration(ms float64) string {
	if ms < 1000 {
		return fmt.Sprintf("%.0fms", ms)
	} else if ms < 60000 {
		return fmt.Sprintf("%.1fs", ms/1000)
	} else {
		minutes := int(ms / 60000)
		seconds := (ms - float64(minutes*60000)) / 1000
		return fmt.Sprintf("%dm%.0fs", minutes, seconds)
	}
}

// truncateString truncates a string to max length with ellipsis
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

// centerString centers a string within a given width
func centerString(s string, width int) string {
	if len(s) >= width {
		return s
	}
	padding := (width - len(s)) / 2
	return strings.Repeat(" ", padding) + s + strings.Repeat(" ", width-len(s)-padding)
}