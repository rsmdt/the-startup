package log

import (
	"sort"
	"strings"
	"sync"
	"time"
)

// eventPair holds correlated Pre and Post tool events
type eventPair struct {
	preEvent  *MetricsEntry
	postEvent *MetricsEntry
}

// AggregateMetrics processes a stream of metrics entries to produce aggregated statistics
// It correlates Pre/Post events, calculates durations, and computes various statistics
func AggregateMetrics(entries <-chan MetricsEntry, filter *MetricsFilter) *MetricsSummary {
	// Initialize tracking structures
	toolStats := make(map[string]*ToolStatistics)
	sessions := make(map[string]bool)
	errorCounts := make(map[string]map[string]int) // errorType -> errorMessage -> count
	errorTools := make(map[string]map[string]bool)  // errorType -> tools set
	hourlyData := make(map[time.Time]*HourlyStats)
	
	// Correlation tracking: toolID -> eventPair
	correlationMap := make(map[string]*eventPair)
	var correlationMutex sync.Mutex
	
	// Track time bounds for the period
	var minTime, maxTime time.Time
	firstEntry := true
	
	// Process entries from channel
	for entry := range entries {
		// Apply filter if provided
		if filter != nil && !filter.Matches(entry) {
			continue
		}
		
		// Update time bounds
		if firstEntry || entry.Timestamp.Before(minTime) {
			minTime = entry.Timestamp
		}
		if firstEntry || entry.Timestamp.After(maxTime) {
			maxTime = entry.Timestamp
			firstEntry = false
		}
		
		// Track unique sessions
		if entry.SessionID != "" {
			sessions[entry.SessionID] = true
		}
		
		// Initialize tool statistics if needed
		if _, exists := toolStats[entry.ToolName]; !exists {
			toolStats[entry.ToolName] = &ToolStatistics{
				Name:          entry.ToolName,
				ErrorTypes:    make(map[string]int),
				MinDurationMs: int64(^uint64(0) >> 1), // Max int64 initially
			}
		}
		
		// Update hourly statistics
		hour := entry.Timestamp.Truncate(time.Hour)
		if _, exists := hourlyData[hour]; !exists {
			hourlyData[hour] = &HourlyStats{
				Hour: hour,
			}
		}
		
		// Process based on event type
		correlationMutex.Lock()
		if entry.HookEvent == "PreToolUse" {
			// Store pre event for correlation
			if entry.ToolID != "" {
				if pair, exists := correlationMap[entry.ToolID]; exists {
					pair.preEvent = &entry
				} else {
					correlationMap[entry.ToolID] = &eventPair{
						preEvent: &entry,
					}
				}
			}
			
			// Update tool statistics for pre-event
			stats := toolStats[entry.ToolName]
			stats.TotalCalls++
			stats.LastUsed = entry.Timestamp
			
			// Update hourly data
			hourlyData[hour].TotalCalls++
			
		} else if entry.HookEvent == "PostToolUse" {
			// Correlate with pre event if available
			var duration *int64
			if entry.ToolID != "" {
				if pair, exists := correlationMap[entry.ToolID]; exists {
					pair.postEvent = &entry
					
					// Calculate duration if we have both events
					if pair.preEvent != nil {
						durationMs := entry.Timestamp.Sub(pair.preEvent.Timestamp).Milliseconds()
						if durationMs >= 0 { // Sanity check
							duration = &durationMs
							
							// Update duration in the entry for consistency
							entry.DurationMs = duration
						}
					}
				} else {
					correlationMap[entry.ToolID] = &eventPair{
						postEvent: &entry,
					}
				}
			}
			
			// Use provided duration if calculation failed
			if duration == nil && entry.DurationMs != nil {
				duration = entry.DurationMs
			}
			
			// Update tool statistics
			stats := toolStats[entry.ToolName]
			
			// Update success/failure counts
			if entry.Success != nil {
				if *entry.Success {
					stats.SuccessCount++
					hourlyData[hour].SuccessCount++
				} else {
					stats.FailureCount++
					hourlyData[hour].FailureCount++
					
					// Track error patterns
					if entry.ErrorType != "" {
						stats.ErrorTypes[entry.ErrorType]++
						
						// Track global error patterns
						if errorCounts[entry.ErrorType] == nil {
							errorCounts[entry.ErrorType] = make(map[string]int)
							errorTools[entry.ErrorType] = make(map[string]bool)
						}
						
						// Normalize error message for grouping
						errorMsg := normalizeErrorMessage(entry.Error)
						errorCounts[entry.ErrorType][errorMsg]++
						errorTools[entry.ErrorType][entry.ToolName] = true
					}
				}
			}
			
			// Update duration statistics
			if duration != nil && *duration > 0 {
				stats.TotalDurationMs += *duration
				
				if *duration < stats.MinDurationMs {
					stats.MinDurationMs = *duration
				}
				if *duration > stats.MaxDurationMs {
					stats.MaxDurationMs = *duration
				}
			}
		}
		correlationMutex.Unlock()
	}
	
	// Post-process statistics
	totalCalls := 0
	totalSuccesses := 0
	totalFailures := 0
	
	for _, stats := range toolStats {
		// Calculate average duration
		if stats.SuccessCount+stats.FailureCount > 0 {
			stats.AvgDurationMs = float64(stats.TotalDurationMs) / float64(stats.SuccessCount+stats.FailureCount)
		}
		
		// Fix min duration if no durations were recorded
		if stats.MinDurationMs == int64(^uint64(0)>>1) {
			stats.MinDurationMs = 0
		}
		
		// Accumulate totals
		totalCalls += stats.TotalCalls
		totalSuccesses += stats.SuccessCount
		totalFailures += stats.FailureCount
	}
	
	// Calculate overall success rate
	var successRate float64
	if totalSuccesses+totalFailures > 0 {
		successRate = float64(totalSuccesses) / float64(totalSuccesses+totalFailures) * 100
	}
	
	// Build top errors list
	topErrors := buildTopErrors(errorCounts, errorTools)
	
	// Build hourly activity list
	hourlyActivity := buildHourlyActivity(hourlyData, toolStats)
	
	// Create summary
	return &MetricsSummary{
		Period: TimePeriod{
			Start: minTime,
			End:   maxTime,
		},
		TotalCalls:     totalCalls,
		UniqueSessions: len(sessions),
		SuccessRate:    successRate,
		ToolStats:      toolStats,
		TopErrors:      topErrors,
		HourlyActivity: hourlyActivity,
	}
}

// normalizeErrorMessage extracts a normalized version of error messages for grouping
// It removes specific values like file paths, IDs, etc. to group similar errors
func normalizeErrorMessage(errorMsg string) string {
	if errorMsg == "" {
		return "unknown error"
	}
	
	// Truncate very long messages
	if len(errorMsg) > 200 {
		errorMsg = errorMsg[:200] + "..."
	}
	
	// Remove file paths (Unix and Windows)
	errorMsg = strings.ReplaceAll(errorMsg, "\\", "/")
	parts := strings.Fields(errorMsg)
	normalized := make([]string, 0, len(parts))
	
	for _, part := range parts {
		// Skip parts that look like paths or IDs
		if strings.Contains(part, "/") && len(part) > 10 {
			normalized = append(normalized, "<path>")
		} else if len(part) > 20 && !strings.Contains(part, " ") {
			// Likely an ID or hash
			normalized = append(normalized, "<id>")
		} else {
			normalized = append(normalized, part)
		}
	}
	
	result := strings.Join(normalized, " ")
	if result == "" {
		return "unknown error"
	}
	return result
}

// buildTopErrors creates a sorted list of most common error patterns
func buildTopErrors(errorCounts map[string]map[string]int, errorTools map[string]map[string]bool) []ErrorPattern {
	var patterns []ErrorPattern
	
	for errorType, messages := range errorCounts {
		for message, count := range messages {
			// Convert tool set to slice
			tools := make([]string, 0, len(errorTools[errorType]))
			for tool := range errorTools[errorType] {
				tools = append(tools, tool)
			}
			sort.Strings(tools)
			
			patterns = append(patterns, ErrorPattern{
				ErrorType:    errorType,
				ErrorMessage: message,
				Count:        count,
				Tools:        tools,
			})
		}
	}
	
	// Sort by count descending, then by error type
	sort.Slice(patterns, func(i, j int) bool {
		if patterns[i].Count != patterns[j].Count {
			return patterns[i].Count > patterns[j].Count
		}
		return patterns[i].ErrorType < patterns[j].ErrorType
	})
	
	// Return top 10 errors
	if len(patterns) > 10 {
		patterns = patterns[:10]
	}
	
	return patterns
}

// buildHourlyActivity creates hourly statistics from the aggregated data
func buildHourlyActivity(hourlyData map[time.Time]*HourlyStats, toolStats map[string]*ToolStatistics) []HourlyStats {
	// Count unique tools per hour
	hourlyTools := make(map[time.Time]map[string]bool)
	
	for _, stats := range toolStats {
		hour := stats.LastUsed.Truncate(time.Hour)
		if _, exists := hourlyTools[hour]; !exists {
			hourlyTools[hour] = make(map[string]bool)
		}
		hourlyTools[hour][stats.Name] = true
	}
	
	// Build sorted list
	var hours []time.Time
	for hour := range hourlyData {
		hours = append(hours, hour)
	}
	sort.Slice(hours, func(i, j int) bool {
		return hours[i].Before(hours[j])
	})
	
	// Create result
	result := make([]HourlyStats, 0, len(hours))
	for _, hour := range hours {
		stats := hourlyData[hour]
		if tools, exists := hourlyTools[hour]; exists {
			stats.UniqueTools = len(tools)
		}
		result = append(result, *stats)
	}
	
	return result
}

// StreamAggregateMetrics is a variant that processes entries as they arrive
// It's useful for real-time monitoring or when dealing with very large datasets
func StreamAggregateMetrics(entries <-chan MetricsEntry, results chan<- *MetricsSummary, updateInterval time.Duration, filter *MetricsFilter) {
	ticker := time.NewTicker(updateInterval)
	defer ticker.Stop()
	
	// Buffer to accumulate entries
	buffer := make([]MetricsEntry, 0, 1000)
	var bufferMutex sync.Mutex
	
	// Background goroutine to process entries
	go func() {
		for entry := range entries {
			bufferMutex.Lock()
			buffer = append(buffer, entry)
			bufferMutex.Unlock()
		}
	}()
	
	// Periodically aggregate and send results
	for range ticker.C {
		bufferMutex.Lock()
		if len(buffer) == 0 {
			bufferMutex.Unlock()
			continue
		}
		
		// Create channel for current batch
		batchChan := make(chan MetricsEntry, len(buffer))
		for _, entry := range buffer {
			batchChan <- entry
		}
		close(batchChan)
		
		// Clear buffer
		buffer = buffer[:0]
		bufferMutex.Unlock()
		
		// Aggregate and send result
		summary := AggregateMetrics(batchChan, filter)
		select {
		case results <- summary:
		default:
			// Don't block if receiver is not ready
		}
	}
}