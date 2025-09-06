package stats

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
	"text/tabwriter"
	"time"
)

// DisplayFormatter implements the Display interface for formatting statistics
type DisplayFormatter struct {
	outputFormat string
}

// NewDisplayFormatter creates a new display formatter with default table format
func NewDisplayFormatter() *DisplayFormatter {
	return &DisplayFormatter{
		outputFormat: "table",
	}
}

// SetOutputFormat configures the output format (table, json, csv)
func (d *DisplayFormatter) SetOutputFormat(format string) error {
	switch format {
	case "table", "json", "csv":
		d.outputFormat = format
		return nil
	default:
		return fmt.Errorf("unsupported output format: %s", format)
	}
}

// Write outputs formatted content to a writer
func (d *DisplayFormatter) Write(w io.Writer, content string) error {
	_, err := w.Write([]byte(content))
	return err
}

// FormatStats formats aggregated statistics for display
func (d *DisplayFormatter) FormatStats(stats *AggregatedStatistics) string {
	if stats == nil {
		return ""
	}

	switch d.outputFormat {
	case "json":
		return d.formatStatsJSON(stats)
	case "csv":
		return d.formatStatsCSV(stats)
	default:
		return d.formatStatsTable(stats)
	}
}

// FormatSession formats session-specific statistics
func (d *DisplayFormatter) FormatSession(session *SessionStatistics) string {
	if session == nil {
		return ""
	}

	switch d.outputFormat {
	case "json":
		return d.formatSessionJSON(session)
	case "csv":
		return d.formatSessionCSV(session)
	default:
		return d.formatSessionTable(session)
	}
}

// FormatToolLeaderboard creates a tool usage leaderboard
func (d *DisplayFormatter) FormatToolLeaderboard(toolStats map[string]*GlobalToolStats, limit int) string {
	if toolStats == nil || len(toolStats) == 0 {
		return ""
	}

	switch d.outputFormat {
	case "json":
		return d.formatToolLeaderboardJSON(toolStats, limit)
	case "csv":
		return d.formatToolLeaderboardCSV(toolStats, limit)
	default:
		return d.formatToolLeaderboardTable(toolStats, limit)
	}
}

// FormatErrors formats error patterns and frequencies
func (d *DisplayFormatter) FormatErrors(errors []ErrorPattern, frequencies []ErrorFrequency) string {
	if len(errors) == 0 && len(frequencies) == 0 {
		return ""
	}

	switch d.outputFormat {
	case "json":
		return d.formatErrorsJSON(errors, frequencies)
	case "csv":
		return d.formatErrorsCSV(errors, frequencies)
	default:
		return d.formatErrorsTable(errors, frequencies)
	}
}

// FormatTimeline formats temporal activity patterns
func (d *DisplayFormatter) FormatTimeline(hourly []HourlyActivity, daily []DailyActivity) string {
	if len(hourly) == 0 && len(daily) == 0 {
		return ""
	}

	switch d.outputFormat {
	case "json":
		return d.formatTimelineJSON(hourly, daily)
	case "csv":
		return d.formatTimelineCSV(hourly, daily)
	default:
		return d.formatTimelineTable(hourly, daily)
	}
}

// Table formatting methods

func (d *DisplayFormatter) formatStatsTable(stats *AggregatedStatistics) string {
	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', tabwriter.AlignRight)

	// Header
	fmt.Fprintln(&buf, "╔════════════════════════════════════════════════════════════════╗")
	fmt.Fprintf(&buf, "║  AGGREGATED STATISTICS - %s to %s  ║\n", 
		stats.Period.Start.Format("2006-01-02 15:04"),
		stats.Period.End.Format("2006-01-02 15:04"))
	fmt.Fprintln(&buf, "╚════════════════════════════════════════════════════════════════╝")
	fmt.Fprintln(&buf)

	// Session overview
	fmt.Fprintln(&buf, "SESSION OVERVIEW")
	fmt.Fprintln(&buf, "────────────────")
	fmt.Fprintf(w, "Total Sessions:\t%d\n", stats.TotalSessions)
	fmt.Fprintf(w, "Active Sessions:\t%d\n", len(stats.ActiveSessions))
	fmt.Fprintf(w, "Avg Duration:\t%s\n", formatDuration(stats.AvgSessionDuration))
	fmt.Fprintln(w)

	// Tool usage summary
	fmt.Fprintln(&buf, "TOOL USAGE SUMMARY")
	fmt.Fprintln(&buf, "──────────────────")
	fmt.Fprintf(w, "Total Tool Calls:\t%d\n", stats.TotalToolCalls)
	fmt.Fprintf(w, "Success Rate:\t%.1f%%\n", stats.OverallSuccessRate*100)
	fmt.Fprintln(w)

	// Performance metrics
	if stats.PerformanceMetrics.AvgResponseTimeMs > 0 {
		fmt.Fprintln(&buf, "PERFORMANCE METRICS")
		fmt.Fprintln(&buf, "───────────────────")
		fmt.Fprintf(w, "Avg Response Time:\t%.0fms\n", stats.PerformanceMetrics.AvgResponseTimeMs)
		fmt.Fprintf(w, "P50 Response Time:\t%.0fms\n", stats.PerformanceMetrics.P50ResponseTimeMs)
		fmt.Fprintf(w, "P95 Response Time:\t%.0fms\n", stats.PerformanceMetrics.P95ResponseTimeMs)
		fmt.Fprintf(w, "P99 Response Time:\t%.0fms\n", stats.PerformanceMetrics.P99ResponseTimeMs)
		fmt.Fprintf(w, "Throughput:\t%.1f calls/min\n", stats.PerformanceMetrics.ToolCallsPerMinute)
		fmt.Fprintln(w)
	}

	w.Flush()
	return buf.String()
}

func (d *DisplayFormatter) formatSessionTable(session *SessionStatistics) string {
	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', tabwriter.AlignRight)

	// Header
	fmt.Fprintln(&buf, "╔════════════════════════════════════════════════════════════════╗")
	fmt.Fprintf(&buf, "║  SESSION: %s  ║\n", session.SessionID)
	fmt.Fprintln(&buf, "╚════════════════════════════════════════════════════════════════╝")
	fmt.Fprintln(&buf)

	// Time information
	fmt.Fprintln(&buf, "TIME INFORMATION")
	fmt.Fprintln(&buf, "────────────────")
	fmt.Fprintf(w, "Start Time:\t%s\n", session.StartTime.Format("2006-01-02 15:04:05"))
	if session.EndTime != nil {
		fmt.Fprintf(w, "End Time:\t%s\n", session.EndTime.Format("2006-01-02 15:04:05"))
	}
	fmt.Fprintf(w, "Duration:\t%s\n", formatDuration(session.Duration))
	fmt.Fprintln(w)

	// Message counts
	fmt.Fprintln(&buf, "MESSAGE COUNTS")
	fmt.Fprintln(&buf, "──────────────")
	fmt.Fprintf(w, "User Messages:\t%d\n", session.UserMessages)
	fmt.Fprintf(w, "Assistant Messages:\t%d\n", session.AssistantMessages)
	fmt.Fprintf(w, "System Messages:\t%d\n", session.SystemMessages)
	fmt.Fprintln(w)

	// Tool usage
	fmt.Fprintln(&buf, "TOOL USAGE")
	fmt.Fprintln(&buf, "──────────")
	fmt.Fprintf(w, "Total Calls:\t%d\n", session.TotalToolCalls)
	fmt.Fprintf(w, "Successful:\t%d\n", session.SuccessfulCalls)
	fmt.Fprintf(w, "Failed:\t%d\n", session.FailedCalls)
	fmt.Fprintf(w, "Error Rate:\t%.1f%%\n", session.ErrorRate*100)
	fmt.Fprintln(w)

	// Token usage
	fmt.Fprintln(&buf, "TOKEN USAGE")
	fmt.Fprintln(&buf, "───────────")
	fmt.Fprintf(w, "Input Tokens:\t%d\n", session.TotalTokens.Input)
	fmt.Fprintf(w, "Output Tokens:\t%d\n", session.TotalTokens.Output)
	fmt.Fprintf(w, "Total Tokens:\t%d\n", session.TotalTokens.Total)
	fmt.Fprintln(w)

	// Tool breakdown
	if len(session.ToolStats) > 0 {
		fmt.Fprintln(&buf, "TOOL BREAKDOWN")
		fmt.Fprintln(&buf, "──────────────")
		fmt.Fprintln(w, "Tool\tCalls\tSuccess\tFail\tAvg Duration")
		fmt.Fprintln(w, "────\t─────\t───────\t────\t────────────")
		
		// Sort tools by call count
		tools := make([]string, 0, len(session.ToolStats))
		for name := range session.ToolStats {
			tools = append(tools, name)
		}
		sort.Slice(tools, func(i, j int) bool {
			return session.ToolStats[tools[i]].CallCount > session.ToolStats[tools[j]].CallCount
		})
		
		for _, name := range tools {
			stats := session.ToolStats[name]
			fmt.Fprintf(w, "%s\t%d\t%d\t%d\t%.0fms\n",
				name, stats.CallCount, stats.SuccessCount, stats.FailureCount, stats.AvgDurationMs)
		}
	}

	w.Flush()
	return buf.String()
}

func (d *DisplayFormatter) formatToolLeaderboardTable(toolStats map[string]*GlobalToolStats, limit int) string {
	var buf bytes.Buffer
	
	// Sort tools by total calls
	type toolEntry struct {
		name  string
		stats *GlobalToolStats
	}
	
	entries := make([]toolEntry, 0, len(toolStats))
	for name, stats := range toolStats {
		entries = append(entries, toolEntry{name, stats})
	}
	
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].stats.TotalCalls > entries[j].stats.TotalCalls
	})
	
	// Apply limit
	if limit > 0 && limit < len(entries) {
		entries = entries[:limit]
	}
	
	// Calculate total for percentages
	totalCalls := 0
	for _, e := range entries {
		totalCalls += e.stats.TotalCalls
	}
	
	// Header
	fmt.Fprintln(&buf, "╔════════════════════════════════════════════════════════════════╗")
	fmt.Fprintln(&buf, "║                    TOOL USAGE LEADERBOARD                     ║")
	fmt.Fprintln(&buf, "╚════════════════════════════════════════════════════════════════╝")
	fmt.Fprintln(&buf)
	
	w := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Rank\tTool\tCalls\t%\tSuccess Rate\tAvg Duration\tSparkline")
	fmt.Fprintln(w, "────\t────\t─────\t─\t────────────\t────────────\t─────────")
	
	for i, entry := range entries {
		percentage := float64(entry.stats.TotalCalls) * 100 / float64(totalCalls)
		sparkline := generateSparkline(entry.stats.TotalCalls, entries[0].stats.TotalCalls, 10)
		
		fmt.Fprintf(w, "%d\t%s\t%d\t%.1f%%\t%.1f%%\t%.0fms\t%s\n",
			i+1,
			entry.name,
			entry.stats.TotalCalls,
			percentage,
			entry.stats.SuccessRate*100,
			entry.stats.AvgDurationMs,
			sparkline)
	}
	
	w.Flush()
	return buf.String()
}

func (d *DisplayFormatter) formatErrorsTable(errors []ErrorPattern, frequencies []ErrorFrequency) string {
	var buf bytes.Buffer
	
	// Header
	fmt.Fprintln(&buf, "╔════════════════════════════════════════════════════════════════╗")
	fmt.Fprintln(&buf, "║                     ERROR ANALYSIS                            ║")
	fmt.Fprintln(&buf, "╚════════════════════════════════════════════════════════════════╝")
	fmt.Fprintln(&buf)
	
	// Error patterns
	if len(errors) > 0 {
		fmt.Fprintln(&buf, "ERROR PATTERNS")
		fmt.Fprintln(&buf, "──────────────")
		
		w := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "Pattern\tCount\tTools\tFirst Seen\tLast Seen")
		fmt.Fprintln(w, "───────\t─────\t─────\t──────────\t─────────")
		
		for _, pattern := range errors {
			fmt.Fprintf(w, "%s\t%d\t%s\t%s\t%s\n",
				truncateString(pattern.Pattern, 30),
				pattern.Count,
				strings.Join(pattern.Tools, ", "),
				pattern.FirstSeen.Format("01-02 15:04"),
				pattern.LastSeen.Format("01-02 15:04"))
		}
		w.Flush()
		fmt.Fprintln(&buf)
	}
	
	// Error frequencies
	if len(frequencies) > 0 {
		fmt.Fprintln(&buf, "TOP ERRORS BY FREQUENCY")
		fmt.Fprintln(&buf, "───────────────────────")
		
		w := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "Tool\tError\tCount\tPercentage")
		fmt.Fprintln(w, "────\t─────\t─────\t──────────")
		
		for _, freq := range frequencies {
			fmt.Fprintf(w, "%s\t%s\t%d\t%.1f%%\n",
				freq.ToolName,
				truncateString(freq.ErrorMessage, 40),
				freq.Count,
				freq.Percentage)
		}
		w.Flush()
	}
	
	return buf.String()
}

func (d *DisplayFormatter) formatTimelineTable(hourly []HourlyActivity, daily []DailyActivity) string {
	var buf bytes.Buffer
	
	// Header
	fmt.Fprintln(&buf, "╔════════════════════════════════════════════════════════════════╗")
	fmt.Fprintln(&buf, "║                    TEMPORAL ACTIVITY                          ║")
	fmt.Fprintln(&buf, "╚════════════════════════════════════════════════════════════════╝")
	fmt.Fprintln(&buf)
	
	// Daily activity
	if len(daily) > 0 {
		fmt.Fprintln(&buf, "DAILY ACTIVITY")
		fmt.Fprintln(&buf, "──────────────")
		
		w := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "Date\tSessions\tTool Calls\tUnique Tools\tSuccess Rate\tErrors")
		fmt.Fprintln(w, "────\t────────\t──────────\t────────────\t────────────\t──────")
		
		for _, day := range daily {
			fmt.Fprintf(w, "%s\t%d\t%d\t%d\t%.1f%%\t%d\n",
				day.Date.Format("2006-01-02"),
				day.Sessions,
				day.ToolCalls,
				day.UniqueTools,
				day.SuccessRate*100,
				day.Errors)
		}
		w.Flush()
		fmt.Fprintln(&buf)
	}
	
	// Hourly activity (show last 24 hours or less)
	if len(hourly) > 0 {
		fmt.Fprintln(&buf, "HOURLY ACTIVITY (Recent)")
		fmt.Fprintln(&buf, "────────────────────────")
		
		// Limit to last 24 entries
		start := 0
		if len(hourly) > 24 {
			start = len(hourly) - 24
		}
		
		w := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "Hour\tCalls\tUnique Tools\tSessions\tSuccess\tFail\tAvg Duration")
		fmt.Fprintln(w, "────\t─────\t────────────\t────────\t───────\t────\t────────────")
		
		for _, hour := range hourly[start:] {
			fmt.Fprintf(w, "%s\t%d\t%d\t%d\t%d\t%d\t%.0fms\n",
				hour.Hour.Format("15:04"),
				hour.ToolCalls,
				hour.UniqueTools,
				hour.UniqueSessions,
				hour.SuccessCount,
				hour.FailureCount,
				hour.AvgDurationMs)
		}
		w.Flush()
	}
	
	return buf.String()
}

// JSON formatting methods

func (d *DisplayFormatter) formatStatsJSON(stats *AggregatedStatistics) string {
	data, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		return fmt.Sprintf("Error formatting JSON: %v", err)
	}
	return string(data)
}

func (d *DisplayFormatter) formatSessionJSON(session *SessionStatistics) string {
	data, err := json.MarshalIndent(session, "", "  ")
	if err != nil {
		return fmt.Sprintf("Error formatting JSON: %v", err)
	}
	return string(data)
}

func (d *DisplayFormatter) formatToolLeaderboardJSON(toolStats map[string]*GlobalToolStats, limit int) string {
	// Convert to sorted slice
	type toolEntry struct {
		Rank  int              `json:"rank"`
		Name  string           `json:"name"`
		Stats *GlobalToolStats `json:"stats"`
	}
	
	entries := make([]toolEntry, 0, len(toolStats))
	for name, stats := range toolStats {
		entries = append(entries, toolEntry{Name: name, Stats: stats})
	}
	
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Stats.TotalCalls > entries[j].Stats.TotalCalls
	})
	
	// Apply limit and set ranks
	if limit > 0 && limit < len(entries) {
		entries = entries[:limit]
	}
	
	for i := range entries {
		entries[i].Rank = i + 1
	}
	
	data, err := json.MarshalIndent(map[string]interface{}{
		"leaderboard": entries,
		"total_tools": len(toolStats),
		"showing":     len(entries),
	}, "", "  ")
	
	if err != nil {
		return fmt.Sprintf("Error formatting JSON: %v", err)
	}
	return string(data)
}

func (d *DisplayFormatter) formatErrorsJSON(errors []ErrorPattern, frequencies []ErrorFrequency) string {
	data, err := json.MarshalIndent(map[string]interface{}{
		"patterns":    errors,
		"frequencies": frequencies,
	}, "", "  ")
	
	if err != nil {
		return fmt.Sprintf("Error formatting JSON: %v", err)
	}
	return string(data)
}

func (d *DisplayFormatter) formatTimelineJSON(hourly []HourlyActivity, daily []DailyActivity) string {
	data, err := json.MarshalIndent(map[string]interface{}{
		"hourly": hourly,
		"daily":  daily,
	}, "", "  ")
	
	if err != nil {
		return fmt.Sprintf("Error formatting JSON: %v", err)
	}
	return string(data)
}

// CSV formatting methods

func (d *DisplayFormatter) formatStatsCSV(stats *AggregatedStatistics) string {
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	
	// Headers
	w.Write([]string{"Metric", "Value"})
	
	// Data rows
	w.Write([]string{"Period Start", stats.Period.Start.Format(time.RFC3339)})
	w.Write([]string{"Period End", stats.Period.End.Format(time.RFC3339)})
	w.Write([]string{"Total Sessions", fmt.Sprintf("%d", stats.TotalSessions)})
	w.Write([]string{"Active Sessions", fmt.Sprintf("%d", len(stats.ActiveSessions))})
	w.Write([]string{"Avg Session Duration", stats.AvgSessionDuration.String()})
	w.Write([]string{"Total Tool Calls", fmt.Sprintf("%d", stats.TotalToolCalls)})
	w.Write([]string{"Overall Success Rate", fmt.Sprintf("%.2f", stats.OverallSuccessRate)})
	
	if stats.PerformanceMetrics.AvgResponseTimeMs > 0 {
		w.Write([]string{"Avg Response Time (ms)", fmt.Sprintf("%.0f", stats.PerformanceMetrics.AvgResponseTimeMs)})
		w.Write([]string{"P50 Response Time (ms)", fmt.Sprintf("%.0f", stats.PerformanceMetrics.P50ResponseTimeMs)})
		w.Write([]string{"P95 Response Time (ms)", fmt.Sprintf("%.0f", stats.PerformanceMetrics.P95ResponseTimeMs)})
		w.Write([]string{"P99 Response Time (ms)", fmt.Sprintf("%.0f", stats.PerformanceMetrics.P99ResponseTimeMs)})
		w.Write([]string{"Tool Calls Per Minute", fmt.Sprintf("%.2f", stats.PerformanceMetrics.ToolCallsPerMinute)})
	}
	
	w.Flush()
	return buf.String()
}

func (d *DisplayFormatter) formatSessionCSV(session *SessionStatistics) string {
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	
	// Headers
	w.Write([]string{"Metric", "Value"})
	
	// Data rows
	w.Write([]string{"Session ID", session.SessionID})
	w.Write([]string{"Start Time", session.StartTime.Format(time.RFC3339)})
	if session.EndTime != nil {
		w.Write([]string{"End Time", session.EndTime.Format(time.RFC3339)})
	}
	w.Write([]string{"Duration", session.Duration.String()})
	w.Write([]string{"User Messages", fmt.Sprintf("%d", session.UserMessages)})
	w.Write([]string{"Assistant Messages", fmt.Sprintf("%d", session.AssistantMessages)})
	w.Write([]string{"System Messages", fmt.Sprintf("%d", session.SystemMessages)})
	w.Write([]string{"Total Tool Calls", fmt.Sprintf("%d", session.TotalToolCalls)})
	w.Write([]string{"Successful Calls", fmt.Sprintf("%d", session.SuccessfulCalls)})
	w.Write([]string{"Failed Calls", fmt.Sprintf("%d", session.FailedCalls)})
	w.Write([]string{"Error Rate", fmt.Sprintf("%.2f", session.ErrorRate)})
	w.Write([]string{"Input Tokens", fmt.Sprintf("%d", session.TotalTokens.Input)})
	w.Write([]string{"Output Tokens", fmt.Sprintf("%d", session.TotalTokens.Output)})
	w.Write([]string{"Total Tokens", fmt.Sprintf("%d", session.TotalTokens.Total)})
	
	w.Flush()
	
	// Add tool breakdown as separate CSV section
	if len(session.ToolStats) > 0 {
		buf.WriteString("\n\nTool Breakdown\n")
		w = csv.NewWriter(&buf)
		w.Write([]string{"Tool", "Calls", "Success", "Failed", "Avg Duration (ms)"})
		
		for name, stats := range session.ToolStats {
			w.Write([]string{
				name,
				fmt.Sprintf("%d", stats.CallCount),
				fmt.Sprintf("%d", stats.SuccessCount),
				fmt.Sprintf("%d", stats.FailureCount),
				fmt.Sprintf("%.0f", stats.AvgDurationMs),
			})
		}
		w.Flush()
	}
	
	return buf.String()
}

func (d *DisplayFormatter) formatToolLeaderboardCSV(toolStats map[string]*GlobalToolStats, limit int) string {
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	
	// Sort tools
	type toolEntry struct {
		name  string
		stats *GlobalToolStats
	}
	
	entries := make([]toolEntry, 0, len(toolStats))
	for name, stats := range toolStats {
		entries = append(entries, toolEntry{name, stats})
	}
	
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].stats.TotalCalls > entries[j].stats.TotalCalls
	})
	
	// Apply limit
	if limit > 0 && limit < len(entries) {
		entries = entries[:limit]
	}
	
	// Headers
	w.Write([]string{"Rank", "Tool", "Total Calls", "Success Rate", "Avg Duration (ms)", "Min Duration (ms)", "Max Duration (ms)", "P95 Duration (ms)", "P99 Duration (ms)"})
	
	// Data rows
	for i, entry := range entries {
		w.Write([]string{
			fmt.Sprintf("%d", i+1),
			entry.name,
			fmt.Sprintf("%d", entry.stats.TotalCalls),
			fmt.Sprintf("%.2f", entry.stats.SuccessRate),
			fmt.Sprintf("%.0f", entry.stats.AvgDurationMs),
			fmt.Sprintf("%d", entry.stats.MinDurationMs),
			fmt.Sprintf("%d", entry.stats.MaxDurationMs),
			fmt.Sprintf("%.0f", entry.stats.P95DurationMs),
			fmt.Sprintf("%.0f", entry.stats.P99DurationMs),
		})
	}
	
	w.Flush()
	return buf.String()
}

func (d *DisplayFormatter) formatErrorsCSV(errors []ErrorPattern, frequencies []ErrorFrequency) string {
	var buf bytes.Buffer
	
	// Error patterns
	if len(errors) > 0 {
		buf.WriteString("Error Patterns\n")
		w := csv.NewWriter(&buf)
		w.Write([]string{"Pattern", "Count", "Tools", "First Seen", "Last Seen"})
		
		for _, pattern := range errors {
			w.Write([]string{
				pattern.Pattern,
				fmt.Sprintf("%d", pattern.Count),
				strings.Join(pattern.Tools, "; "),
				pattern.FirstSeen.Format(time.RFC3339),
				pattern.LastSeen.Format(time.RFC3339),
			})
		}
		w.Flush()
	}
	
	// Error frequencies
	if len(frequencies) > 0 {
		if len(errors) > 0 {
			buf.WriteString("\n\n")
		}
		buf.WriteString("Error Frequencies\n")
		w := csv.NewWriter(&buf)
		w.Write([]string{"Tool", "Error Code", "Error Message", "Count", "Percentage"})
		
		for _, freq := range frequencies {
			w.Write([]string{
				freq.ToolName,
				freq.ErrorCode,
				freq.ErrorMessage,
				fmt.Sprintf("%d", freq.Count),
				fmt.Sprintf("%.2f", freq.Percentage),
			})
		}
		w.Flush()
	}
	
	return buf.String()
}

func (d *DisplayFormatter) formatTimelineCSV(hourly []HourlyActivity, daily []DailyActivity) string {
	var buf bytes.Buffer
	
	// Daily activity
	if len(daily) > 0 {
		buf.WriteString("Daily Activity\n")
		w := csv.NewWriter(&buf)
		w.Write([]string{"Date", "Sessions", "Tool Calls", "Unique Tools", "Success Rate", "Total Duration (ms)", "Errors"})
		
		for _, day := range daily {
			w.Write([]string{
				day.Date.Format("2006-01-02"),
				fmt.Sprintf("%d", day.Sessions),
				fmt.Sprintf("%d", day.ToolCalls),
				fmt.Sprintf("%d", day.UniqueTools),
				fmt.Sprintf("%.2f", day.SuccessRate),
				fmt.Sprintf("%d", day.TotalDurationMs),
				fmt.Sprintf("%d", day.Errors),
			})
		}
		w.Flush()
	}
	
	// Hourly activity
	if len(hourly) > 0 {
		if len(daily) > 0 {
			buf.WriteString("\n\n")
		}
		buf.WriteString("Hourly Activity\n")
		w := csv.NewWriter(&buf)
		w.Write([]string{"Hour", "Tool Calls", "Unique Tools", "Unique Sessions", "Success Count", "Failure Count", "Avg Duration (ms)"})
		
		for _, hour := range hourly {
			w.Write([]string{
				hour.Hour.Format("2006-01-02 15:04"),
				fmt.Sprintf("%d", hour.ToolCalls),
				fmt.Sprintf("%d", hour.UniqueTools),
				fmt.Sprintf("%d", hour.UniqueSessions),
				fmt.Sprintf("%d", hour.SuccessCount),
				fmt.Sprintf("%d", hour.FailureCount),
				fmt.Sprintf("%.0f", hour.AvgDurationMs),
			})
		}
		w.Flush()
	}
	
	return buf.String()
}

// Helper functions

func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%.0fs", d.Seconds())
	}
	if d < time.Hour {
		return fmt.Sprintf("%.0fm", d.Minutes())
	}
	return fmt.Sprintf("%.1fh", d.Hours())
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func generateSparkline(value, max, width int) string {
	if max == 0 {
		return strings.Repeat("─", width)
	}
	
	chars := []string{"▁", "▂", "▃", "▄", "▅", "▆", "▇", "█"}
	ratio := float64(value) / float64(max)
	filled := int(ratio * float64(width))
	
	var sparkline strings.Builder
	for i := 0; i < width; i++ {
		if i < filled {
			// Calculate which character to use based on the fractional part
			charIndex := int(ratio * float64(len(chars)-1))
			if charIndex >= len(chars) {
				charIndex = len(chars) - 1
			}
			sparkline.WriteString(chars[charIndex])
		} else {
			sparkline.WriteString("─")
		}
	}
	
	return sparkline.String()
}

// CommandLeaderboard formats a leaderboard for commands (used by agents)
func (d *DisplayFormatter) FormatCommandLeaderboard(commands map[string]int) string {
	if commands == nil || len(commands) == 0 {
		return ""
	}

	switch d.outputFormat {
	case "json":
		return d.formatCommandLeaderboardJSON(commands)
	case "csv":
		return d.formatCommandLeaderboardCSV(commands)
	default:
		return d.formatCommandLeaderboardTable(commands)
	}
}

func (d *DisplayFormatter) formatCommandLeaderboardTable(commands map[string]int) string {
	var buf bytes.Buffer
	
	// Sort commands by count
	type cmdEntry struct {
		name  string
		count int
	}
	
	entries := make([]cmdEntry, 0, len(commands))
	totalCount := 0
	for name, count := range commands {
		entries = append(entries, cmdEntry{name, count})
		totalCount += count
	}
	
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].count > entries[j].count
	})
	
	// Header
	fmt.Fprintln(&buf, "╔════════════════════════════════════════════════════════════════╗")
	fmt.Fprintln(&buf, "║                   COMMAND/AGENT LEADERBOARD                   ║")
	fmt.Fprintln(&buf, "╚════════════════════════════════════════════════════════════════╝")
	fmt.Fprintln(&buf)
	
	w := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Rank\tCommand/Agent\tUsage\tPercentage\tSparkline")
	fmt.Fprintln(w, "────\t─────────────\t─────\t──────────\t─────────")
	
	maxCount := 0
	if len(entries) > 0 {
		maxCount = entries[0].count
	}
	
	for i, entry := range entries {
		percentage := float64(entry.count) * 100 / float64(totalCount)
		sparkline := generateSparkline(entry.count, maxCount, 15)
		
		fmt.Fprintf(w, "%d\t%s\t%d\t%.1f%%\t%s\n",
			i+1,
			entry.name,
			entry.count,
			percentage,
			sparkline)
	}
	
	w.Flush()
	return buf.String()
}

func (d *DisplayFormatter) formatCommandLeaderboardJSON(commands map[string]int) string {
	type cmdEntry struct {
		Rank  int    `json:"rank"`
		Name  string `json:"name"`
		Count int    `json:"count"`
	}
	
	entries := make([]cmdEntry, 0, len(commands))
	for name, count := range commands {
		entries = append(entries, cmdEntry{Name: name, Count: count})
	}
	
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Count > entries[j].Count
	})
	
	for i := range entries {
		entries[i].Rank = i + 1
	}
	
	data, err := json.MarshalIndent(map[string]interface{}{
		"leaderboard":    entries,
		"total_commands": len(commands),
	}, "", "  ")
	
	if err != nil {
		return fmt.Sprintf("Error formatting JSON: %v", err)
	}
	return string(data)
}

func (d *DisplayFormatter) formatCommandLeaderboardCSV(commands map[string]int) string {
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	
	// Sort commands
	type cmdEntry struct {
		name  string
		count int
	}
	
	entries := make([]cmdEntry, 0, len(commands))
	for name, count := range commands {
		entries = append(entries, cmdEntry{name, count})
	}
	
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].count > entries[j].count
	})
	
	// Headers
	w.Write([]string{"Rank", "Command/Agent", "Usage Count"})
	
	// Data rows
	for i, entry := range entries {
		w.Write([]string{
			fmt.Sprintf("%d", i+1),
			entry.name,
			fmt.Sprintf("%d", entry.count),
		})
	}
	
	w.Flush()
	return buf.String()
}