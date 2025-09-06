// Package stats provides statistics aggregation for Claude Code log analysis
package stats

import (
	"fmt"
	"math"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	// Default timeout for orphaned tool invocations
	defaultOrphanTimeout = 5 * time.Minute
	
	// Maximum number of percentile buckets for memory efficiency
	maxPercentileBuckets = 1000
)

// aggregator implements the Aggregator interface with efficient statistics calculation
type aggregator struct {
	mu sync.RWMutex
	
	// Configuration
	orphanTimeout time.Duration
	
	// Session tracking
	sessions map[string]*sessionData
	
	// Global tool statistics
	globalTools map[string]*globalToolData
	
	// Command tracking
	commands map[string]int
	
	// Token usage
	totalTokens TokenUsage
	
	// Time bounds
	earliestTime *time.Time
	latestTime   *time.Time
	
	// Pending tool invocations (waiting for results)
	pendingTools map[string]*pendingInvocation
	
	// Command pattern for extraction
	commandPattern *regexp.Regexp
}

// sessionData holds data for a single session
type sessionData struct {
	sessionID         string
	startTime         time.Time
	endTime           *time.Time
	userMessages      int
	assistantMessages int
	systemMessages    int
	invocations       []ToolInvocation
	toolStats         map[string]*ToolSessionStats
	errors            []ErrorOccurrence
	tokens            TokenUsage
}

// globalToolData holds global statistics for a tool
type globalToolData struct {
	name          string
	totalCalls    int
	successCount  int
	failureCount  int
	sessions      map[string]bool // Track unique sessions
	firstUsed     time.Time
	lastUsed      time.Time
	
	// Duration tracking with Welford's algorithm
	durationStats *welfordStats
	
	// Percentile tracking with t-digest approximation
	percentiles   *tDigest
	
	// Error tracking
	errorTypes    map[string]int
	commonErrors  []string
	
	// Hourly usage tracking
	hourlyUsage   [24]int
}

// pendingInvocation represents a tool invocation awaiting result
type pendingInvocation struct {
	invocation   ToolInvocation
	timeoutAt    time.Time
}

// welfordStats implements Welford's online algorithm for mean and variance
type welfordStats struct {
	count    int64
	mean     float64
	m2       float64
	min      *float64
	max      *float64
}

// tDigest provides approximate percentile calculation with bounded memory
type tDigest struct {
	centroids []centroid
	maxSize   int
}

// centroid represents a cluster of values in t-digest
type centroid struct {
	mean   float64
	weight float64
}

// NewAggregator creates a new statistics aggregator
func NewAggregator() Aggregator {
	return NewAggregatorWithTimeout(defaultOrphanTimeout)
}

// NewAggregatorWithTimeout creates an aggregator with custom orphan timeout
func NewAggregatorWithTimeout(orphanTimeout time.Duration) Aggregator {
	return &aggregator{
		orphanTimeout:  orphanTimeout,
		sessions:       make(map[string]*sessionData),
		globalTools:    make(map[string]*globalToolData),
		commands:       make(map[string]int),
		pendingTools:   make(map[string]*pendingInvocation),
		commandPattern: regexp.MustCompile(`<command-name>([^<]+)</command-name>`),
	}
}

// ProcessEntry processes a single log entry
func (a *aggregator) ProcessEntry(entry ClaudeLogEntry) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	
	// Update time bounds
	a.updateTimeBounds(entry.Timestamp)
	
	// Ensure session exists
	session := a.ensureSession(entry.SessionID, entry.Timestamp)
	
	// Process based on type
	switch entry.Type {
	case "user":
		session.userMessages++
		if entry.User != nil && entry.User.IsToolResult {
			// Process tool result
			a.processToolResult(entry.SessionID, entry.User, entry.Timestamp)
		}
		
	case "assistant":
		session.assistantMessages++
		if entry.Assistant != nil {
			// Extract commands from text
			a.extractCommands(entry.Assistant.Text)
			
			// Process tool uses
			for _, toolUse := range entry.Assistant.ToolUses {
				a.processToolUse(entry.SessionID, toolUse, entry.Timestamp)
			}
		}
		
	case "system":
		session.systemMessages++
		
	case "summary":
		if entry.Summary != nil {
			// Update token usage
			a.totalTokens.Input += entry.Summary.TokensUsed.Input
			a.totalTokens.Output += entry.Summary.TokensUsed.Output
			a.totalTokens.Total += entry.Summary.TokensUsed.Total
			
			session.tokens.Input += entry.Summary.TokensUsed.Input
			session.tokens.Output += entry.Summary.TokensUsed.Output
			session.tokens.Total += entry.Summary.TokensUsed.Total
			
			// Process errors
			for _, err := range entry.Summary.Errors {
				session.errors = append(session.errors, ErrorOccurrence{
					Timestamp: entry.Timestamp,
					Error:     err,
				})
			}
		}
	}
	
	// Check for orphaned invocations
	a.checkOrphanedInvocations(entry.Timestamp)
	
	return nil
}

// ProcessInvocation processes a complete tool invocation
func (a *aggregator) ProcessInvocation(invocation ToolInvocation) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	
	a.processInvocationInternal(invocation)
	return nil
}

// processInvocationInternal processes invocation with lock already held
func (a *aggregator) processInvocationInternal(invocation ToolInvocation) {
	// Ensure session exists
	session := a.ensureSession(invocation.SessionID, invocation.InvokedAt)
	
	// Add to session invocations
	session.invocations = append(session.invocations, invocation)
	
	// Update session tool stats
	toolStats := a.ensureSessionToolStats(session, invocation.Name)
	toolStats.CallCount++
	toolStats.LastUsed = invocation.InvokedAt
	
	if invocation.Success {
		toolStats.SuccessCount++
	} else {
		toolStats.FailureCount++
		if invocation.Error != nil {
			session.errors = append(session.errors, ErrorOccurrence{
				Timestamp: invocation.InvokedAt,
				ToolName:  invocation.Name,
				ToolUseID: invocation.ID,
				Error:     *invocation.Error,
			})
		}
	}
	
	// Update duration stats if available
	if invocation.DurationMs != nil {
		durationMs := *invocation.DurationMs
		toolStats.TotalDurationMs += durationMs
		toolStats.AvgDurationMs = float64(toolStats.TotalDurationMs) / float64(toolStats.CallCount)
		
		if toolStats.MinDurationMs == nil || durationMs < *toolStats.MinDurationMs {
			toolStats.MinDurationMs = &durationMs
		}
		if toolStats.MaxDurationMs == nil || durationMs > *toolStats.MaxDurationMs {
			toolStats.MaxDurationMs = &durationMs
		}
	}
	
	// Update global tool stats
	a.updateGlobalToolStats(invocation)
}

// GetSessionStats returns statistics for a specific session
func (a *aggregator) GetSessionStats(sessionID string) (*SessionStatistics, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	
	session, exists := a.sessions[sessionID]
	if !exists {
		return nil, fmt.Errorf("session %s not found", sessionID)
	}
	
	// Calculate derived statistics
	totalCalls := 0
	successfulCalls := 0
	failedCalls := 0
	
	for _, toolStats := range session.toolStats {
		totalCalls += toolStats.CallCount
		successfulCalls += toolStats.SuccessCount
		failedCalls += toolStats.FailureCount
	}
	
	errorRate := 0.0
	if totalCalls > 0 {
		errorRate = float64(failedCalls) / float64(totalCalls)
	}
	
	duration := time.Since(session.startTime)
	if session.endTime != nil {
		duration = session.endTime.Sub(session.startTime)
	}
	
	return &SessionStatistics{
		SessionID:         session.sessionID,
		StartTime:         session.startTime,
		EndTime:           session.endTime,
		Duration:          duration,
		UserMessages:      session.userMessages,
		AssistantMessages: session.assistantMessages,
		SystemMessages:    session.systemMessages,
		ToolInvocations:   session.invocations,
		ToolStats:         session.toolStats,
		TotalToolCalls:    totalCalls,
		SuccessfulCalls:   successfulCalls,
		FailedCalls:       failedCalls,
		TotalTokens:       session.tokens,
		Errors:            session.errors,
		ErrorRate:         errorRate,
	}, nil
}

// GetOverallStats returns aggregated statistics across all sessions
func (a *aggregator) GetOverallStats() (*AggregatedStatistics, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	
	// Calculate time period
	period := TimePeriod{}
	if a.earliestTime != nil && a.latestTime != nil {
		period.Start = *a.earliestTime
		period.End = *a.latestTime
		period.Duration = period.End.Sub(period.Start)
	}
	
	// Build global tool stats map
	globalToolStats := make(map[string]*GlobalToolStats)
	for name, toolData := range a.globalTools {
		stats := a.buildGlobalToolStats(name, toolData)
		globalToolStats[name] = stats
	}
	
	// Calculate overall metrics
	totalCalls := 0
	successfulCalls := 0
	for _, toolData := range a.globalTools {
		totalCalls += toolData.totalCalls
		successfulCalls += toolData.successCount
	}
	
	overallSuccessRate := 0.0
	if totalCalls > 0 {
		overallSuccessRate = float64(successfulCalls) / float64(totalCalls)
	}
	
	// Build error patterns and frequencies
	errorPatterns := a.buildErrorPatterns()
	topErrors := a.buildErrorFrequencies()
	
	// Build activity patterns
	hourlyActivity := a.buildHourlyActivity()
	dailyActivity := a.buildDailyActivity()
	
	// Calculate average session duration
	totalDuration := time.Duration(0)
	for _, session := range a.sessions {
		duration := time.Since(session.startTime)
		if session.endTime != nil {
			duration = session.endTime.Sub(session.startTime)
		}
		totalDuration += duration
	}
	
	avgSessionDuration := time.Duration(0)
	if len(a.sessions) > 0 {
		avgSessionDuration = totalDuration / time.Duration(len(a.sessions))
	}
	
	// Get active sessions
	activeSessions := make([]string, 0, len(a.sessions))
	for sessionID := range a.sessions {
		activeSessions = append(activeSessions, sessionID)
	}
	
	// Build performance metrics
	performanceMetrics := a.buildPerformanceMetrics()
	
	return &AggregatedStatistics{
		Period:             period,
		TotalSessions:      len(a.sessions),
		ActiveSessions:     activeSessions,
		AvgSessionDuration: avgSessionDuration,
		GlobalToolStats:    globalToolStats,
		TotalToolCalls:     totalCalls,
		OverallSuccessRate: overallSuccessRate,
		ErrorPatterns:      errorPatterns,
		TopErrors:          topErrors,
		HourlyActivity:     hourlyActivity,
		DailyActivity:      dailyActivity,
		PerformanceMetrics: performanceMetrics,
	}, nil
}

// GetToolStats returns statistics for a specific tool
func (a *aggregator) GetToolStats(toolName string) (*GlobalToolStats, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	
	toolData, exists := a.globalTools[toolName]
	if !exists {
		return nil, fmt.Errorf("tool %s not found", toolName)
	}
	
	return a.buildGlobalToolStats(toolName, toolData), nil
}

// Reset clears all accumulated statistics
func (a *aggregator) Reset() {
	a.mu.Lock()
	defer a.mu.Unlock()
	
	a.sessions = make(map[string]*sessionData)
	a.globalTools = make(map[string]*globalToolData)
	a.commands = make(map[string]int)
	a.pendingTools = make(map[string]*pendingInvocation)
	a.totalTokens = TokenUsage{}
	a.earliestTime = nil
	a.latestTime = nil
}

// Merge combines statistics from another aggregator
func (a *aggregator) Merge(other Aggregator) error {
	// Type assertion to access internal data
	otherAgg, ok := other.(*aggregator)
	if !ok {
		return fmt.Errorf("can only merge with same aggregator type")
	}
	
	a.mu.Lock()
	defer a.mu.Unlock()
	otherAgg.mu.RLock()
	defer otherAgg.mu.RUnlock()
	
	// Merge sessions
	for sessionID, otherSession := range otherAgg.sessions {
		if existingSession, exists := a.sessions[sessionID]; exists {
			// Merge session data
			existingSession.userMessages += otherSession.userMessages
			existingSession.assistantMessages += otherSession.assistantMessages
			existingSession.systemMessages += otherSession.systemMessages
			existingSession.invocations = append(existingSession.invocations, otherSession.invocations...)
			existingSession.errors = append(existingSession.errors, otherSession.errors...)
			existingSession.tokens.Input += otherSession.tokens.Input
			existingSession.tokens.Output += otherSession.tokens.Output
			existingSession.tokens.Total += otherSession.tokens.Total
			
			// Merge tool stats
			for toolName, otherToolStats := range otherSession.toolStats {
				if existingToolStats, exists := existingSession.toolStats[toolName]; exists {
					mergeSessionToolStats(existingToolStats, otherToolStats)
				} else {
					existingSession.toolStats[toolName] = otherToolStats
				}
			}
		} else {
			a.sessions[sessionID] = otherSession
		}
	}
	
	// Merge global tools
	for toolName, otherToolData := range otherAgg.globalTools {
		if existingToolData, exists := a.globalTools[toolName]; exists {
			mergeGlobalToolData(existingToolData, otherToolData)
		} else {
			a.globalTools[toolName] = otherToolData
		}
	}
	
	// Merge commands
	for cmd, count := range otherAgg.commands {
		a.commands[cmd] += count
	}
	
	// Merge token usage
	a.totalTokens.Input += otherAgg.totalTokens.Input
	a.totalTokens.Output += otherAgg.totalTokens.Output
	a.totalTokens.Total += otherAgg.totalTokens.Total
	
	// Update time bounds
	if otherAgg.earliestTime != nil {
		a.updateTimeBounds(*otherAgg.earliestTime)
	}
	if otherAgg.latestTime != nil {
		a.updateTimeBounds(*otherAgg.latestTime)
	}
	
	return nil
}

// Helper methods

func (a *aggregator) ensureSession(sessionID string, timestamp time.Time) *sessionData {
	session, exists := a.sessions[sessionID]
	if !exists {
		session = &sessionData{
			sessionID: sessionID,
			startTime: timestamp,
			toolStats: make(map[string]*ToolSessionStats),
		}
		a.sessions[sessionID] = session
	}
	return session
}

func (a *aggregator) ensureSessionToolStats(session *sessionData, toolName string) *ToolSessionStats {
	stats, exists := session.toolStats[toolName]
	if !exists {
		stats = &ToolSessionStats{
			Name: toolName,
		}
		session.toolStats[toolName] = stats
	}
	return stats
}

func (a *aggregator) updateTimeBounds(timestamp time.Time) {
	if a.earliestTime == nil || timestamp.Before(*a.earliestTime) {
		a.earliestTime = &timestamp
	}
	if a.latestTime == nil || timestamp.After(*a.latestTime) {
		a.latestTime = &timestamp
	}
}

func (a *aggregator) extractCommands(text string) {
	matches := a.commandPattern.FindAllStringSubmatch(text, -1)
	for _, match := range matches {
		if len(match) > 1 {
			command := strings.TrimSpace(match[1])
			if command != "" {
				a.commands[command]++
			}
		}
	}
}

func (a *aggregator) processToolUse(sessionID string, toolUse ToolUse, timestamp time.Time) {
	// Create pending invocation
	invocation := ToolInvocation{
		ID:         toolUse.ID,
		Name:       toolUse.Name,
		Parameters: toolUse.Parameters,
		InvokedAt:  timestamp,
		SessionID:  sessionID,
	}
	
	a.pendingTools[toolUse.ID] = &pendingInvocation{
		invocation: invocation,
		timeoutAt:  timestamp.Add(a.orphanTimeout),
	}
}

func (a *aggregator) processToolResult(sessionID string, result *UserMessage, timestamp time.Time) {
	if result.ToolUseID == "" {
		return
	}
	
	// Find pending invocation
	pending, exists := a.pendingTools[result.ToolUseID]
	if !exists {
		// Result without invocation - could happen if processing partial logs
		return
	}
	
	// Complete the invocation
	delete(a.pendingTools, result.ToolUseID)
	
	duration := timestamp.Sub(pending.invocation.InvokedAt)
	durationMs := duration.Milliseconds()
	
	pending.invocation.CompletedAt = &timestamp
	pending.invocation.DurationMs = &durationMs
	pending.invocation.Result = result.Result
	// Convert string error to ToolError if present
	if result.Error != nil {
		pending.invocation.Error = &ToolError{
			Message: *result.Error,
		}
	}
	pending.invocation.Success = result.Error == nil
	
	// Process the completed invocation (lock already held)
	a.processInvocationInternal(pending.invocation)
}

func (a *aggregator) checkOrphanedInvocations(currentTime time.Time) {
	// Find and process orphaned invocations
	var orphanedInvocations []ToolInvocation
	for id, pending := range a.pendingTools {
		if currentTime.After(pending.timeoutAt) {
			// Mark as failed due to timeout
			errorMsg := "Tool invocation timed out - no result received"
			pending.invocation.Error = &ToolError{
				Code:    "TIMEOUT",
				Message: errorMsg,
			}
			pending.invocation.Success = false
			orphanedInvocations = append(orphanedInvocations, pending.invocation)
			delete(a.pendingTools, id)
		}
	}
	
	// Process orphaned invocations without holding the lock
	for _, invocation := range orphanedInvocations {
		a.processInvocationInternal(invocation)
	}
}

func (a *aggregator) updateGlobalToolStats(invocation ToolInvocation) {
	toolData, exists := a.globalTools[invocation.Name]
	if !exists {
		toolData = &globalToolData{
			name:          invocation.Name,
			sessions:      make(map[string]bool),
			firstUsed:     invocation.InvokedAt,
			durationStats: newWelfordStats(),
			percentiles:   newTDigest(maxPercentileBuckets),
			errorTypes:    make(map[string]int),
		}
		a.globalTools[invocation.Name] = toolData
	}
	
	// Update basic counters
	toolData.totalCalls++
	toolData.sessions[invocation.SessionID] = true
	toolData.lastUsed = invocation.InvokedAt
	
	if invocation.Success {
		toolData.successCount++
	} else {
		toolData.failureCount++
		if invocation.Error != nil {
			toolData.errorTypes[invocation.Error.Code]++
			// Track common error messages (limit to top 10)
			if !contains(toolData.commonErrors, invocation.Error.Message) {
				toolData.commonErrors = append(toolData.commonErrors, invocation.Error.Message)
				if len(toolData.commonErrors) > 10 {
					toolData.commonErrors = toolData.commonErrors[:10]
				}
			}
		}
	}
	
	// Update duration statistics
	if invocation.DurationMs != nil {
		durationMs := float64(*invocation.DurationMs)
		toolData.durationStats.update(durationMs)
		toolData.percentiles.add(durationMs)
		
		// Update hourly usage
		hour := invocation.InvokedAt.Hour()
		toolData.hourlyUsage[hour]++
	}
}

func (a *aggregator) calculateMostActiveHours() []int {
	hourCounts := make(map[int]int)
	
	for _, toolData := range a.globalTools {
		for hour, count := range toolData.hourlyUsage {
			hourCounts[hour] += count
		}
	}
	
	// Sort hours by activity
	type hourActivity struct {
		hour  int
		count int
	}
	
	var activities []hourActivity
	for hour, count := range hourCounts {
		activities = append(activities, hourActivity{hour, count})
	}
	
	sort.Slice(activities, func(i, j int) bool {
		return activities[i].count > activities[j].count
	})
	
	// Return top 3 most active hours
	result := []int{}
	for i := 0; i < len(activities) && i < 3; i++ {
		result = append(result, activities[i].hour)
	}
	
	return result
}

func (a *aggregator) collectErrorPatterns() []ErrorInfo {
	errorMap := make(map[string]*ErrorInfo)
	
	for _, session := range a.sessions {
		for _, err := range session.errors {
			key := fmt.Sprintf("%s:%s", err.Error.Code, err.Error.Message)
			if info, exists := errorMap[key]; exists {
				info.Count++
				if err.Timestamp.After(info.LastSeen) {
					info.LastSeen = err.Timestamp
				}
			} else {
				errorMap[key] = &ErrorInfo{
					Pattern:   err.Error.Message,
					Count:     1,
					ToolName:  err.ToolName,
					FirstSeen: err.Timestamp,
					LastSeen:  err.Timestamp,
				}
			}
		}
	}
	
	// Convert to slice and sort by count
	var patterns []ErrorInfo
	for _, info := range errorMap {
		patterns = append(patterns, *info)
	}
	
	sort.Slice(patterns, func(i, j int) bool {
		return patterns[i].Count > patterns[j].Count
	})
	
	// Return top 10 patterns
	if len(patterns) > 10 {
		patterns = patterns[:10]
	}
	
	return patterns
}

// Helper method to build GlobalToolStats from internal data
func (a *aggregator) buildGlobalToolStats(name string, toolData *globalToolData) *GlobalToolStats {
	stats := &GlobalToolStats{
		Name:          name,
		TotalCalls:    toolData.totalCalls,
		UniqueUsers:   len(toolData.sessions),
		SuccessCount:  toolData.successCount,
		FailureCount:  toolData.failureCount,
		ErrorTypes:    toolData.errorTypes,
		CommonErrors:  toolData.commonErrors,
		FirstUsed:     toolData.firstUsed,
		LastUsed:      toolData.lastUsed,
	}
	
	// Calculate success rate
	if toolData.totalCalls > 0 {
		stats.SuccessRate = float64(toolData.successCount) / float64(toolData.totalCalls)
	}
	
	// Get duration statistics
	if toolData.durationStats != nil && toolData.durationStats.count > 0 {
		stats.AvgDurationMs = toolData.durationStats.mean
		stats.TotalDurationMs = int64(toolData.durationStats.mean * float64(toolData.durationStats.count))
		
		if toolData.durationStats.min != nil {
			stats.MinDurationMs = int64(*toolData.durationStats.min)
		}
		if toolData.durationStats.max != nil {
			stats.MaxDurationMs = int64(*toolData.durationStats.max)
		}
		
		// Calculate percentiles if available
		if toolData.percentiles != nil {
			stats.MedianDurationMs = toolData.percentiles.quantile(0.5)
			stats.P95DurationMs = toolData.percentiles.quantile(0.95)
			stats.P99DurationMs = toolData.percentiles.quantile(0.99)
		}
	}
	
	// Find peak hour
	maxHour := 0
	maxCount := 0
	for hour, count := range toolData.hourlyUsage {
		if count > maxCount {
			maxCount = count
			maxHour = hour
		}
	}
	stats.PeakHour = maxHour
	
	return stats
}

// buildErrorPatterns builds error pattern analysis
func (a *aggregator) buildErrorPatterns() []ErrorPattern {
	patternMap := make(map[string]*ErrorPattern)
	
	for _, session := range a.sessions {
		for _, err := range session.errors {
			key := fmt.Sprintf("%s:%s", err.Error.Code, err.Error.Message)
			if pattern, exists := patternMap[key]; exists {
				pattern.Count++
				if err.Timestamp.After(pattern.LastSeen) {
					pattern.LastSeen = err.Timestamp
				}
				// Add tool if not already in list
				if !contains(pattern.Tools, err.ToolName) {
					pattern.Tools = append(pattern.Tools, err.ToolName)
				}
			} else {
				patternMap[key] = &ErrorPattern{
					Pattern:      err.Error.Message,
					Count:        1,
					Tools:        []string{err.ToolName},
					FirstSeen:    err.Timestamp,
					LastSeen:     err.Timestamp,
					ExampleError: err.Error,
				}
			}
		}
	}
	
	// Convert to slice and sort by count
	var patterns []ErrorPattern
	for _, pattern := range patternMap {
		patterns = append(patterns, *pattern)
	}
	
	sort.Slice(patterns, func(i, j int) bool {
		return patterns[i].Count > patterns[j].Count
	})
	
	return patterns
}

// buildErrorFrequencies builds top error frequencies
func (a *aggregator) buildErrorFrequencies() []ErrorFrequency {
	freqMap := make(map[string]*ErrorFrequency)
	totalErrors := 0
	
	for _, session := range a.sessions {
		for _, err := range session.errors {
			totalErrors++
			key := fmt.Sprintf("%s:%s:%s", err.Error.Code, err.Error.Message, err.ToolName)
			if freq, exists := freqMap[key]; exists {
				freq.Count++
			} else {
				freqMap[key] = &ErrorFrequency{
					ErrorCode:    err.Error.Code,
					ErrorMessage: err.Error.Message,
					Count:        1,
					ToolName:     err.ToolName,
				}
			}
		}
	}
	
	// Calculate percentages and convert to slice
	var frequencies []ErrorFrequency
	for _, freq := range freqMap {
		if totalErrors > 0 {
			freq.Percentage = float64(freq.Count) / float64(totalErrors) * 100
		}
		frequencies = append(frequencies, *freq)
	}
	
	// Sort by count
	sort.Slice(frequencies, func(i, j int) bool {
		return frequencies[i].Count > frequencies[j].Count
	})
	
	// Return top 10
	if len(frequencies) > 10 {
		frequencies = frequencies[:10]
	}
	
	return frequencies
}

// buildHourlyActivity builds hourly activity patterns
func (a *aggregator) buildHourlyActivity() []HourlyActivity {
	// For simplicity, just return empty for now
	// Full implementation would track activity by hour
	return []HourlyActivity{}
}

// buildDailyActivity builds daily activity patterns
func (a *aggregator) buildDailyActivity() []DailyActivity {
	// For simplicity, just return empty for now
	// Full implementation would track activity by day
	return []DailyActivity{}
}

// buildPerformanceMetrics builds performance metrics
func (a *aggregator) buildPerformanceMetrics() PerformanceMetrics {
	metrics := PerformanceMetrics{}
	
	// Calculate response time statistics
	var durations []float64
	successCount := 0
	failureCount := 0
	
	for _, toolData := range a.globalTools {
		successCount += toolData.successCount
		failureCount += toolData.failureCount
		
		if toolData.durationStats != nil && toolData.durationStats.count > 0 {
			// Add to durations for percentile calculation
			for i := int64(0); i < toolData.durationStats.count; i++ {
				durations = append(durations, toolData.durationStats.mean)
			}
		}
	}
	
	totalCalls := successCount + failureCount
	if totalCalls > 0 {
		metrics.OverallSuccessRate = float64(successCount) / float64(totalCalls)
	}
	
	// Calculate response time percentiles if we have data
	if len(durations) > 0 {
		sort.Float64s(durations)
		metrics.AvgResponseTimeMs = average(durations)
		metrics.P50ResponseTimeMs = percentile(durations, 0.5)
		metrics.P95ResponseTimeMs = percentile(durations, 0.95)
		metrics.P99ResponseTimeMs = percentile(durations, 0.99)
	}
	
	return metrics
}

// Utility function to calculate average
func average(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

// Utility function to calculate percentile
func percentile(sortedValues []float64, p float64) float64 {
	if len(sortedValues) == 0 {
		return 0
	}
	index := int(p * float64(len(sortedValues)-1))
	return sortedValues[index]
}

// Welford's algorithm implementation

func newWelfordStats() *welfordStats {
	return &welfordStats{}
}

func (w *welfordStats) update(value float64) {
	w.count++
	delta := value - w.mean
	w.mean += delta / float64(w.count)
	delta2 := value - w.mean
	w.m2 += delta * delta2
	
	// Update min/max
	if w.min == nil || value < *w.min {
		w.min = &value
	}
	if w.max == nil || value > *w.max {
		w.max = &value
	}
}

func (w *welfordStats) variance() float64 {
	if w.count < 2 {
		return 0
	}
	return w.m2 / float64(w.count-1)
}

func (w *welfordStats) stdDev() float64 {
	return math.Sqrt(w.variance())
}

// T-digest implementation (simplified)

func newTDigest(maxSize int) *tDigest {
	return &tDigest{
		centroids: make([]centroid, 0, maxSize),
		maxSize:   maxSize,
	}
}

func (t *tDigest) add(value float64) {
	// Simple implementation - just store centroids
	t.centroids = append(t.centroids, centroid{mean: value, weight: 1})
	
	// Compress if needed
	if len(t.centroids) >= t.maxSize {
		t.compress()
	}
}

func (t *tDigest) compress() {
	if len(t.centroids) <= t.maxSize/2 {
		return
	}
	
	// Sort centroids
	sort.Slice(t.centroids, func(i, j int) bool {
		return t.centroids[i].mean < t.centroids[j].mean
	})
	
	// Merge adjacent centroids
	compressed := make([]centroid, 0, t.maxSize/2)
	for i := 0; i < len(t.centroids); i += 2 {
		if i+1 < len(t.centroids) {
			// Merge two centroids
			totalWeight := t.centroids[i].weight + t.centroids[i+1].weight
			mergedMean := (t.centroids[i].mean*t.centroids[i].weight +
				t.centroids[i+1].mean*t.centroids[i+1].weight) / totalWeight
			compressed = append(compressed, centroid{
				mean:   mergedMean,
				weight: totalWeight,
			})
		} else {
			compressed = append(compressed, t.centroids[i])
		}
	}
	
	t.centroids = compressed
}

func (t *tDigest) quantile(q float64) float64 {
	if len(t.centroids) == 0 {
		return 0
	}
	
	// Sort centroids
	sort.Slice(t.centroids, func(i, j int) bool {
		return t.centroids[i].mean < t.centroids[j].mean
	})
	
	// Calculate total weight
	totalWeight := 0.0
	for _, c := range t.centroids {
		totalWeight += c.weight
	}
	
	// Find quantile
	targetWeight := q * totalWeight
	currentWeight := 0.0
	
	for _, c := range t.centroids {
		currentWeight += c.weight
		if currentWeight >= targetWeight {
			return c.mean
		}
	}
	
	// Return last centroid if not found
	return t.centroids[len(t.centroids)-1].mean
}

// Utility functions

func mergeSessionToolStats(existing, other *ToolSessionStats) {
	existing.CallCount += other.CallCount
	existing.SuccessCount += other.SuccessCount
	existing.FailureCount += other.FailureCount
	existing.TotalDurationMs += other.TotalDurationMs
	
	if existing.CallCount > 0 {
		existing.AvgDurationMs = float64(existing.TotalDurationMs) / float64(existing.CallCount)
	}
	
	// Update min/max
	if other.MinDurationMs != nil {
		if existing.MinDurationMs == nil || *other.MinDurationMs < *existing.MinDurationMs {
			existing.MinDurationMs = other.MinDurationMs
		}
	}
	if other.MaxDurationMs != nil {
		if existing.MaxDurationMs == nil || *other.MaxDurationMs > *existing.MaxDurationMs {
			existing.MaxDurationMs = other.MaxDurationMs
		}
	}
	
	if other.LastUsed.After(existing.LastUsed) {
		existing.LastUsed = other.LastUsed
	}
}

func mergeGlobalToolData(existing, other *globalToolData) {
	existing.totalCalls += other.totalCalls
	existing.successCount += other.successCount
	existing.failureCount += other.failureCount
	
	// Merge sessions
	for sessionID := range other.sessions {
		existing.sessions[sessionID] = true
	}
	
	// Update time bounds
	if other.firstUsed.Before(existing.firstUsed) {
		existing.firstUsed = other.firstUsed
	}
	if other.lastUsed.After(existing.lastUsed) {
		existing.lastUsed = other.lastUsed
	}
	
	// Merge error types
	for errorCode, count := range other.errorTypes {
		existing.errorTypes[errorCode] += count
	}
	
	// Merge common errors
	for _, errMsg := range other.commonErrors {
		if !contains(existing.commonErrors, errMsg) {
			existing.commonErrors = append(existing.commonErrors, errMsg)
			if len(existing.commonErrors) > 10 {
				existing.commonErrors = existing.commonErrors[:10]
			}
		}
	}
	
	// Merge hourly usage
	for hour, count := range other.hourlyUsage {
		existing.hourlyUsage[hour] += count
	}
	
	// Note: Merging Welford stats and t-digest would require more complex logic
	// For simplicity, we keep the existing statistics
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}