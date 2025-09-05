package log

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// StreamMetrics reads metrics entries from daily JSONL files and streams them through a channel
// The channel is closed when all matching files have been processed
// Malformed lines are silently skipped
func StreamMetrics(filter MetricsFilter) (<-chan MetricsEntry, error) {
	// Get the startup path
	startupPath, err := getStartupPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get startup path: %w", err)
	}

	logsDir := filepath.Join(startupPath, "logs")

	// Check if logs directory exists
	if _, err := os.Stat(logsDir); os.IsNotExist(err) {
		// Return empty channel if no logs directory
		ch := make(chan MetricsEntry)
		close(ch)
		return ch, nil
	}

	// Get list of files to read based on date range
	files, err := getFilesInDateRange(logsDir, filter.StartDate, filter.EndDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get files in date range: %w", err)
	}

	// Create channel with buffer for better performance
	ch := make(chan MetricsEntry, 100)

	// Start goroutine to stream metrics
	go func() {
		defer close(ch)

		for _, file := range files {
			streamFileMetrics(file, filter, ch)
		}
	}()

	return ch, nil
}

// GetTodayMetrics returns a channel streaming all metrics from today
func GetTodayMetrics() (<-chan MetricsEntry, error) {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24*time.Hour).Add(-time.Nanosecond)

	filter := MetricsFilter{
		StartDate: startOfDay,
		EndDate:   endOfDay,
	}

	return StreamMetrics(filter)
}

// GetYesterdayMetrics returns a channel streaming all metrics from yesterday
func GetYesterdayMetrics() (<-chan MetricsEntry, error) {
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	startOfDay := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, yesterday.Location())
	endOfDay := startOfDay.Add(24*time.Hour).Add(-time.Nanosecond)

	filter := MetricsFilter{
		StartDate: startOfDay,
		EndDate:   endOfDay,
	}

	return StreamMetrics(filter)
}

// GetMetricsForDateRange returns a channel streaming metrics for a specific date range
func GetMetricsForDateRange(startDate, endDate time.Time) (<-chan MetricsEntry, error) {
	// Ensure end date includes the full day
	endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, endDate.Location())

	filter := MetricsFilter{
		StartDate: startDate,
		EndDate:   endDate,
	}

	return StreamMetrics(filter)
}

// GetMetricsForTool returns a channel streaming metrics for a specific tool
func GetMetricsForTool(toolName string, startDate, endDate time.Time) (<-chan MetricsEntry, error) {
	filter := MetricsFilter{
		StartDate: startDate,
		EndDate:   endDate,
		ToolNames: []string{toolName},
	}

	return StreamMetrics(filter)
}

// GetMetricsForSession returns a channel streaming metrics for a specific session
func GetMetricsForSession(sessionID string) (<-chan MetricsEntry, error) {
	// For session queries, scan all available files
	// Sessions typically span a single day but we'll scan broadly to be safe
	now := time.Now()
	thirtyDaysAgo := now.AddDate(0, 0, -30)

	filter := MetricsFilter{
		StartDate:  thirtyDaysAgo,
		EndDate:    now,
		SessionIDs: []string{sessionID},
	}

	return StreamMetrics(filter)
}

// getFilesInDateRange returns a sorted list of JSONL files that fall within the date range
func getFilesInDateRange(logsDir string, startDate, endDate time.Time) ([]string, error) {
	// Read all files in logs directory
	entries, err := os.ReadDir(logsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read logs directory: %w", err)
	}

	var files []string

	// Iterate through each day in the range
	current := startDate
	for !current.After(endDate) {
		filename := fmt.Sprintf("%s.jsonl", current.Format("20060102"))
		filePath := filepath.Join(logsDir, filename)

		// Check if file exists
		if fileExists(filePath, entries) {
			files = append(files, filePath)
		}

		// Move to next day
		current = current.AddDate(0, 0, 1)
	}

	// Sort files by name (which sorts by date given our naming scheme)
	sort.Strings(files)

	return files, nil
}

// fileExists checks if a file exists in the directory entries
func fileExists(filePath string, entries []os.DirEntry) bool {
	filename := filepath.Base(filePath)
	for _, entry := range entries {
		if !entry.IsDir() && entry.Name() == filename {
			return true
		}
	}
	return false
}

// streamFileMetrics reads a single JSONL file and streams matching entries to the channel
func streamFileMetrics(filePath string, filter MetricsFilter, ch chan<- MetricsEntry) {
	file, err := os.Open(filePath)
	if err != nil {
		// Silent failure - skip this file
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			// Skip malformed lines
			continue
		}

		// Skip empty lines
		line = trimLineEnding(line)
		if len(line) == 0 {
			continue
		}

		// Parse JSON
		var entry MetricsEntry
		if err := json.Unmarshal(line, &entry); err != nil {
			// Skip malformed JSON lines
			continue
		}

		// Check if entry matches filter
		if filter.Matches(entry) {
			ch <- entry
		}
	}
}

// trimLineEnding removes trailing newline and carriage return characters
func trimLineEnding(line []byte) []byte {
	// Remove trailing \n
	if len(line) > 0 && line[len(line)-1] == '\n' {
		line = line[:len(line)-1]
	}
	// Remove trailing \r (for Windows-style line endings)
	if len(line) > 0 && line[len(line)-1] == '\r' {
		line = line[:len(line)-1]
	}
	return line
}

// ListAvailableDates returns a list of dates for which log files exist
func ListAvailableDates() ([]time.Time, error) {
	startupPath, err := getStartupPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get startup path: %w", err)
	}

	logsDir := filepath.Join(startupPath, "logs")

	// Check if logs directory exists
	if _, err := os.Stat(logsDir); os.IsNotExist(err) {
		return []time.Time{}, nil
	}

	entries, err := os.ReadDir(logsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read logs directory: %w", err)
	}

	var dates []time.Time

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// Check if file matches our pattern (YYYYMMDD.jsonl)
		name := entry.Name()
		if !strings.HasSuffix(name, ".jsonl") {
			continue
		}

		// Extract date from filename
		dateStr := strings.TrimSuffix(name, ".jsonl")
		if len(dateStr) != 8 {
			continue
		}

		// Parse date
		date, err := time.Parse("20060102", dateStr)
		if err != nil {
			continue
		}

		dates = append(dates, date)
	}

	// Sort dates
	sort.Slice(dates, func(i, j int) bool {
		return dates[i].Before(dates[j])
	})

	return dates, nil
}

// GetLatestMetrics returns a channel streaming metrics from the most recent log file
func GetLatestMetrics() (<-chan MetricsEntry, error) {
	dates, err := ListAvailableDates()
	if err != nil {
		return nil, err
	}

	if len(dates) == 0 {
		// Return empty channel if no logs
		ch := make(chan MetricsEntry)
		close(ch)
		return ch, nil
	}

	// Get the latest date
	latestDate := dates[len(dates)-1]
	
	// Stream metrics for that date
	return GetMetricsForDateRange(latestDate, latestDate)
}

// GetMetricsCount returns the total number of metrics entries matching the filter
// This is useful for progress indicators and pagination
func GetMetricsCount(filter MetricsFilter) (int, error) {
	count := 0
	
	ch, err := StreamMetrics(filter)
	if err != nil {
		return 0, err
	}

	for range ch {
		count++
	}

	return count, nil
}