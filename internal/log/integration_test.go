package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

// TestEndToEndHookProcessing tests the complete flow from hook input to metrics output
func TestEndToEndHookProcessing(t *testing.T) {
	// Setup test environment
	tempDir := t.TempDir()
	os.Setenv("THE_STARTUP_PATH", tempDir)
	defer os.Unsetenv("THE_STARTUP_PATH")

	// Create logs directory
	logsDir := filepath.Join(tempDir, "logs")
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Test data: PreToolUse followed by PostToolUse
	// Use a fixed base time to ensure timestamps are within the same second
	baseTime := time.Now().UTC().Truncate(time.Second)
	
	preHookData := HookPayload{
		HookEventName: "PreToolUse",
		ToolName:      "Edit",
		SessionID:     "test-session-123",
		Timestamp:     baseTime.Format(time.RFC3339),
		CWD:           "/test/dir",
		ToolInput:     json.RawMessage(`{"file": "test.go", "content": "hello world"}`),
	}

	postHookData := HookPayload{
		HookEventName: "PostToolUse",
		ToolName:      "Edit",
		SessionID:     "test-session-123",
		Timestamp:     baseTime.Add(500 * time.Millisecond).Format(time.RFC3339),
		CWD:           "/test/dir",
		ToolInput:     json.RawMessage(`{"file": "test.go", "content": "hello world"}`),
		ToolResponse:  json.RawMessage(`{"success": true, "lines_changed": 5}`),
	}

	// Process PreToolUse event
	preHookJSON, _ := json.Marshal(preHookData)
	if err := ProcessHook(bytes.NewReader(preHookJSON)); err != nil {
		t.Errorf("ProcessHook failed for PreToolUse: %v", err)
	}

	// Process PostToolUse event
	postHookJSON, _ := json.Marshal(postHookData)
	if err := ProcessHook(bytes.NewReader(postHookJSON)); err != nil {
		t.Errorf("ProcessHook failed for PostToolUse: %v", err)
	}

	// Verify metrics were written to daily file
	today := time.Now().UTC().Format("20060102")
	metricsFile := filepath.Join(logsDir, fmt.Sprintf("%s.jsonl", today))
	if _, err := os.Stat(metricsFile); os.IsNotExist(err) {
		t.Fatal("Metrics file was not created")
	}

	// Read metrics back using StreamMetrics
	filter := MetricsFilter{
		StartDate: time.Now().UTC().Truncate(24 * time.Hour),
		EndDate:   time.Now().UTC().Add(24 * time.Hour),
	}

	metricsChan, err := StreamMetrics(filter)
	if err != nil {
		t.Fatal(err)
	}

	var entries []MetricsEntry
	for entry := range metricsChan {
		entries = append(entries, entry)
	}

	// Verify we got both events
	if len(entries) != 2 {
		t.Fatalf("Expected 2 metrics entries, got %d", len(entries))
	}

	// Verify Pre event
	if entries[0].HookEvent != "PreToolUse" {
		t.Errorf("Expected first event to be PreToolUse, got %s", entries[0].HookEvent)
	}
	if entries[0].ToolName != "Edit" {
		t.Errorf("Expected tool name Edit, got %s", entries[0].ToolName)
	}

	// Verify Post event
	if entries[1].HookEvent != "PostToolUse" {
		t.Errorf("Expected second event to be PostToolUse, got %s", entries[1].HookEvent)
	}
	if entries[1].Success == nil || !*entries[1].Success {
		t.Error("Expected success to be true for PostToolUse")
	}

	// Both events should have the same ToolID for correlation
	if entries[0].ToolID != entries[1].ToolID {
		t.Errorf("ToolIDs don't match for correlation: %s vs %s", 
			entries[0].ToolID, entries[1].ToolID)
	}

	// Test aggregation on the collected metrics
	aggregateChan := make(chan MetricsEntry, len(entries))
	for _, entry := range entries {
		aggregateChan <- entry
	}
	close(aggregateChan)

	summary := AggregateMetrics(aggregateChan, nil)
	
	// Verify aggregated statistics
	if summary.TotalCalls != 1 {
		t.Errorf("Expected 1 total call, got %d", summary.TotalCalls)
	}
	if summary.UniqueSessions != 1 {
		t.Errorf("Expected 1 unique session, got %d", summary.UniqueSessions)
	}
	if summary.SuccessRate != 100.0 {
		t.Errorf("Expected 100%% success rate, got %f", summary.SuccessRate)
	}

	// Check tool-specific stats
	editStats, exists := summary.ToolStats["Edit"]
	if !exists {
		t.Fatal("Edit tool stats not found")
	}
	if editStats.TotalCalls != 1 {
		t.Errorf("Expected 1 call for Edit tool, got %d", editStats.TotalCalls)
	}
	if editStats.SuccessCount != 1 {
		t.Errorf("Expected 1 success for Edit tool, got %d", editStats.SuccessCount)
	}
}

// TestConcurrentWrites tests that concurrent hook processing doesn't corrupt data
func TestConcurrentWrites(t *testing.T) {
	// Setup test environment
	tempDir := t.TempDir()
	os.Setenv("THE_STARTUP_PATH", tempDir)
	defer os.Unsetenv("THE_STARTUP_PATH")

	// Create logs directory
	logsDir := filepath.Join(tempDir, "logs")
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Number of concurrent writers and events per writer
	numWriters := 10
	eventsPerWriter := 100
	totalEvents := numWriters * eventsPerWriter * 2 // Pre and Post for each

	var wg sync.WaitGroup
	wg.Add(numWriters)

	// Launch concurrent writers
	for i := 0; i < numWriters; i++ {
		go func(writerID int) {
			defer wg.Done()

			sessionID := fmt.Sprintf("session-%d", writerID)
			
			for j := 0; j < eventsPerWriter; j++ {
				// Create unique tool name for tracking
				toolName := fmt.Sprintf("Tool%d_%d", writerID, j)
				timestamp := time.Now().UTC()

				// Send PreToolUse
				preHook := HookPayload{
					HookEventName: "PreToolUse",
					ToolName:      toolName,
					SessionID:     sessionID,
					Timestamp:     timestamp.Format(time.RFC3339),
					ToolInput:     json.RawMessage(fmt.Sprintf(`{"id": %d}`, j)),
				}
				
				preJSON, _ := json.Marshal(preHook)
				ProcessHook(bytes.NewReader(preJSON))

				// Small delay to simulate processing time
				time.Sleep(time.Millisecond)

				// Send PostToolUse
				postHook := HookPayload{
					HookEventName: "PostToolUse",
					ToolName:      toolName,
					SessionID:     sessionID,
					Timestamp:     timestamp.Add(10 * time.Millisecond).Format(time.RFC3339),
					ToolInput:     json.RawMessage(fmt.Sprintf(`{"id": %d}`, j)),
					ToolResponse:  json.RawMessage(`{"success": true}`),
				}
				
				postJSON, _ := json.Marshal(postHook)
				ProcessHook(bytes.NewReader(postJSON))
			}
		}(i)
	}

	// Wait for all writers to complete
	wg.Wait()

	// Give a moment for final writes to flush
	time.Sleep(100 * time.Millisecond)

	// Read all metrics back
	filter := MetricsFilter{
		StartDate: time.Now().UTC().Truncate(24 * time.Hour),
		EndDate:   time.Now().UTC().Add(24 * time.Hour),
	}

	metricsChan, err := StreamMetrics(filter)
	if err != nil {
		t.Fatal(err)
	}

	// Collect all entries and verify count
	entriesMap := make(map[string]*MetricsEntry)
	for entry := range metricsChan {
		// Use a combination of tool name and event type as key
		key := fmt.Sprintf("%s-%s-%s", entry.ToolName, entry.HookEvent, entry.ToolID)
		entriesMap[key] = &entry
	}

	actualEvents := len(entriesMap)
	if actualEvents != totalEvents {
		t.Errorf("Expected %d events, got %d", totalEvents, actualEvents)
	}

	// Verify data integrity - check for proper JSON structure
	today := time.Now().UTC().Format("20060102")
	metricsFile := filepath.Join(logsDir, fmt.Sprintf("%s.jsonl", today))
	
	file, err := os.Open(metricsFile)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	lineCount := 0
	for decoder.More() {
		var entry MetricsEntry
		if err := decoder.Decode(&entry); err != nil {
			t.Errorf("Failed to decode line %d: %v", lineCount+1, err)
		}
		lineCount++
	}

	if lineCount != totalEvents {
		t.Errorf("File contains %d lines, expected %d", lineCount, totalEvents)
	}
}

// TestMetricsPipelineWithFiltering tests the complete pipeline with various filters
func TestMetricsPipelineWithFiltering(t *testing.T) {
	// Setup test environment
	tempDir := t.TempDir()
	os.Setenv("THE_STARTUP_PATH", tempDir)
	defer os.Unsetenv("THE_STARTUP_PATH")

	// Create logs directory
	logsDir := filepath.Join(tempDir, "logs")
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Generate test data with various tools and outcomes
	tools := []string{"Edit", "Write", "Read", "Bash", "Search"}
	sessions := []string{"session-1", "session-2", "session-3"}
	
	for _, tool := range tools {
		for i, session := range sessions {
			// Some tools will fail for certain sessions
			shouldFail := (tool == "Bash" && i == 0) || (tool == "Write" && i == 1)
			
			preHook := HookPayload{
				HookEventName: "PreToolUse",
				ToolName:      tool,
				SessionID:     session,
				Timestamp:     time.Now().UTC().Format(time.RFC3339),
				ToolInput:     json.RawMessage(`{"test": true}`),
			}
			
			preJSON, _ := json.Marshal(preHook)
			ProcessHook(bytes.NewReader(preJSON))

			// Post event with varying durations
			duration := time.Duration((i+1)*100) * time.Millisecond
			postHook := HookPayload{
				HookEventName: "PostToolUse",
				ToolName:      tool,
				SessionID:     session,
				Timestamp:     time.Now().UTC().Add(duration).Format(time.RFC3339),
				ToolInput:     json.RawMessage(`{"test": true}`),
			}

			if shouldFail {
				postHook.Error = "Test error"
				postHook.ErrorType = "TestError"
				postHook.ToolResponse = json.RawMessage(`{"error": "Test error"}`)
			} else {
				postHook.ToolResponse = json.RawMessage(`{"success": true}`)
			}

			postJSON, _ := json.Marshal(postHook)
			ProcessHook(bytes.NewReader(postJSON))
		}
	}

	// Test 1: Filter by tool name
	t.Run("FilterByTool", func(t *testing.T) {
		filter := MetricsFilter{
			StartDate: time.Now().UTC().Truncate(24 * time.Hour),
			EndDate:   time.Now().UTC().Add(24 * time.Hour),
			ToolNames: []string{"Edit", "Write"},
		}

		metricsChan, err := StreamMetrics(filter)
		if err != nil {
			t.Fatal(err)
		}

		count := 0
		for entry := range metricsChan {
			if entry.ToolName != "Edit" && entry.ToolName != "Write" {
				t.Errorf("Unexpected tool in filtered results: %s", entry.ToolName)
			}
			count++
		}

		// 2 tools * 3 sessions * 2 events (pre/post) = 12
		if count != 12 {
			t.Errorf("Expected 12 filtered entries, got %d", count)
		}
	})

	// Test 2: Filter by session
	t.Run("FilterBySession", func(t *testing.T) {
		filter := MetricsFilter{
			StartDate:  time.Now().UTC().Truncate(24 * time.Hour),
			EndDate:    time.Now().UTC().Add(24 * time.Hour),
			SessionIDs: []string{"session-1"},
		}

		metricsChan, err := StreamMetrics(filter)
		if err != nil {
			t.Fatal(err)
		}

		count := 0
		for entry := range metricsChan {
			if entry.SessionID != "session-1" {
				t.Errorf("Unexpected session in filtered results: %s", entry.SessionID)
			}
			count++
		}

		// 5 tools * 1 session * 2 events = 10
		if count != 10 {
			t.Errorf("Expected 10 filtered entries, got %d", count)
		}
	})

	// Test 3: Aggregate and verify error tracking
	t.Run("AggregateWithErrors", func(t *testing.T) {
		filter := MetricsFilter{
			StartDate: time.Now().UTC().Truncate(24 * time.Hour),
			EndDate:   time.Now().UTC().Add(24 * time.Hour),
		}

		metricsChan, err := StreamMetrics(filter)
		if err != nil {
			t.Fatal(err)
		}

		// Collect entries for aggregation
		var entries []MetricsEntry
		for entry := range metricsChan {
			entries = append(entries, entry)
		}

		// Create channel for aggregation
		aggregateChan := make(chan MetricsEntry, len(entries))
		for _, entry := range entries {
			aggregateChan <- entry
		}
		close(aggregateChan)

		summary := AggregateMetrics(aggregateChan, nil)

		// Verify Bash tool has failures
		bashStats, exists := summary.ToolStats["Bash"]
		if !exists {
			t.Fatal("Bash tool stats not found")
		}
		if bashStats.FailureCount != 1 {
			t.Errorf("Expected 1 failure for Bash, got %d", bashStats.FailureCount)
		}

		// Verify Write tool has failures
		writeStats, exists := summary.ToolStats["Write"]
		if !exists {
			t.Fatal("Write tool stats not found")
		}
		if writeStats.FailureCount != 1 {
			t.Errorf("Expected 1 failure for Write, got %d", writeStats.FailureCount)
		}

		// Check error patterns
		if len(summary.TopErrors) == 0 {
			t.Error("Expected error patterns to be tracked")
		}
	})
}

// TestMultipleDayMetrics tests handling metrics across multiple days
func TestMultipleDayMetrics(t *testing.T) {
	// Setup test environment
	tempDir := t.TempDir()
	os.Setenv("THE_STARTUP_PATH", tempDir)
	defer os.Unsetenv("THE_STARTUP_PATH")

	// Create logs directory
	logsDir := filepath.Join(tempDir, "logs")
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create metrics for multiple days
	now := time.Now().UTC()
	dates := []time.Time{
		now.AddDate(0, 0, -2), // 2 days ago
		now.AddDate(0, 0, -1), // yesterday
		now,                    // today
	}

	for dayIndex, date := range dates {
		// Create entries directly in the appropriate daily file
		for i := 0; i < 5; i++ {
			entry := MetricsEntry{
				ToolName:  fmt.Sprintf("Tool%d", i),
				HookEvent: "PreToolUse",
				Timestamp: date,
				SessionID: fmt.Sprintf("session-day%d", dayIndex),
				ToolID:    fmt.Sprintf("tool_%d_%d", dayIndex, i),
				ToolInput: json.RawMessage(`{"test": true}`),
			}

			// Manually set the timestamp for the file name
			oldTimestamp := entry.Timestamp
			entry.Timestamp = date
			AppendMetrics(&entry)
			entry.Timestamp = oldTimestamp
		}
	}

	// Test reading metrics from date range
	t.Run("DateRangeQuery", func(t *testing.T) {
		// Query for yesterday's metrics
		yesterday := now.AddDate(0, 0, -1)
		filter := MetricsFilter{
			StartDate: yesterday.Truncate(24 * time.Hour),
			EndDate:   yesterday.Truncate(24 * time.Hour).Add(23*time.Hour + 59*time.Minute + 59*time.Second),
		}

		metricsChan, err := StreamMetrics(filter)
		if err != nil {
			t.Fatal(err)
		}

		count := 0
		for entry := range metricsChan {
			// Verify entries are from yesterday
			entryDate := entry.Timestamp.Truncate(24 * time.Hour)
			expectedDate := yesterday.Truncate(24 * time.Hour)
			if !entryDate.Equal(expectedDate) {
				t.Errorf("Entry date %v doesn't match expected %v", entryDate, expectedDate)
			}
			count++
		}

		if count != 5 {
			t.Errorf("Expected 5 entries for yesterday, got %d", count)
		}
	})

	// Test listing available dates
	t.Run("ListAvailableDates", func(t *testing.T) {
		availableDates, err := ListAvailableDates()
		if err != nil {
			t.Fatal(err)
		}

		if len(availableDates) != 3 {
			t.Errorf("Expected 3 available dates, got %d", len(availableDates))
		}

		// Verify dates are in ascending order
		for i := 1; i < len(availableDates); i++ {
			if !availableDates[i].After(availableDates[i-1]) {
				t.Error("Dates are not in ascending order")
			}
		}
	})

	// Test GetLatestMetrics
	t.Run("GetLatestMetrics", func(t *testing.T) {
		metricsChan, err := GetLatestMetrics()
		if err != nil {
			t.Fatal(err)
		}

		count := 0
		for entry := range metricsChan {
			// All entries should be from today
			entryDate := entry.Timestamp.Truncate(24 * time.Hour)
			todayDate := now.Truncate(24 * time.Hour)
			if !entryDate.Equal(todayDate) {
				t.Errorf("Latest metrics contains entry from %v, expected today %v", 
					entryDate, todayDate)
			}
			count++
		}

		if count != 5 {
			t.Errorf("Expected 5 entries for today, got %d", count)
		}
	})
}

// TestCorrelationAcrossEvents tests that Pre and Post events are properly correlated
func TestCorrelationAcrossEvents(t *testing.T) {
	// Setup test environment
	tempDir := t.TempDir()
	os.Setenv("THE_STARTUP_PATH", tempDir)
	defer os.Unsetenv("THE_STARTUP_PATH")

	// Create logs directory
	logsDir := filepath.Join(tempDir, "logs")
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Generate correlated Pre/Post events
	// Use truncated time to ensure events stay within the same second
	baseTime := time.Now().UTC().Truncate(time.Second)
	
	// Event 1: Pre without Post (incomplete)
	pre1 := HookPayload{
		HookEventName: "PreToolUse",
		ToolName:      "Edit",
		SessionID:     "session-1",
		Timestamp:     baseTime.Format(time.RFC3339),
		ToolInput:     json.RawMessage(`{"file": "test1.go"}`),
	}
	pre1JSON, _ := json.Marshal(pre1)
	ProcessHook(bytes.NewReader(pre1JSON))

	// Event 2: Both Pre and Post (use same second for correlation)
	event2Time := baseTime.Add(2 * time.Second)
	pre2 := HookPayload{
		HookEventName: "PreToolUse",
		ToolName:      "Write",
		SessionID:     "session-1",
		Timestamp:     event2Time.Format(time.RFC3339),
		ToolInput:     json.RawMessage(`{"file": "test2.go"}`),
	}
	pre2JSON, _ := json.Marshal(pre2)
	ProcessHook(bytes.NewReader(pre2JSON))

	post2 := HookPayload{
		HookEventName: "PostToolUse",
		ToolName:      "Write",
		SessionID:     "session-1",
		Timestamp:     event2Time.Add(500*time.Millisecond).Format(time.RFC3339),
		ToolInput:     json.RawMessage(`{"file": "test2.go"}`),
		ToolResponse:  json.RawMessage(`{"success": true, "bytes_written": 1024}`),
	}
	post2JSON, _ := json.Marshal(post2)
	ProcessHook(bytes.NewReader(post2JSON))

	// Event 3: Post without Pre (orphaned)
	event3Time := baseTime.Add(4 * time.Second)
	post3 := HookPayload{
		HookEventName: "PostToolUse",
		ToolName:      "Read",
		SessionID:     "session-1",
		Timestamp:     event3Time.Format(time.RFC3339),
		ToolInput:     json.RawMessage(`{"file": "test3.go"}`),
		ToolResponse:  json.RawMessage(`{"success": true, "content": "data"}`),
		Error:         "File not found",
		ErrorType:     "FileError",
	}
	post3JSON, _ := json.Marshal(post3)
	ProcessHook(bytes.NewReader(post3JSON))

	// Read all metrics
	filter := MetricsFilter{
		StartDate: baseTime.Truncate(24 * time.Hour),
		EndDate:   baseTime.Add(24 * time.Hour),
	}

	metricsChan, err := StreamMetrics(filter)
	if err != nil {
		t.Fatal(err)
	}

	// Group by ToolID for correlation check
	toolIDMap := make(map[string][]*MetricsEntry)
	for entry := range metricsChan {
		e := entry // Capture for pointer
		toolIDMap[entry.ToolID] = append(toolIDMap[entry.ToolID], &e)
	}

	// Check correlations
	correlatedPairs := 0
	orphanedPre := 0
	orphanedPost := 0

	for _, events := range toolIDMap {
		hasPre := false
		hasPost := false
		
		for _, event := range events {
			if event.HookEvent == "PreToolUse" {
				hasPre = true
			} else if event.HookEvent == "PostToolUse" {
				hasPost = true
			}
		}

		if hasPre && hasPost {
			correlatedPairs++
		} else if hasPre && !hasPost {
			orphanedPre++
		} else if !hasPre && hasPost {
			orphanedPost++
		}
	}

	if correlatedPairs != 1 {
		t.Errorf("Expected 1 correlated pair, got %d", correlatedPairs)
	}
	if orphanedPre != 1 {
		t.Errorf("Expected 1 orphaned Pre event, got %d", orphanedPre)
	}
	if orphanedPost != 1 {
		t.Errorf("Expected 1 orphaned Post event, got %d", orphanedPost)
	}

	// Test aggregation handles partial correlations correctly
	metricsChan2, _ := StreamMetrics(filter)
	entries := make([]MetricsEntry, 0)
	for entry := range metricsChan2 {
		entries = append(entries, entry)
	}

	aggregateChan := make(chan MetricsEntry, len(entries))
	for _, entry := range entries {
		aggregateChan <- entry
	}
	close(aggregateChan)

	summary := AggregateMetrics(aggregateChan, nil)
	
	// Write tool should have complete stats (both Pre and Post)
	writeStats := summary.ToolStats["Write"]
	if writeStats.TotalCalls != 1 {
		t.Errorf("Expected 1 call for Write, got %d", writeStats.TotalCalls)
	}
	if writeStats.SuccessCount != 1 {
		t.Errorf("Expected 1 success for Write, got %d", writeStats.SuccessCount)
	}

	// Edit tool should have incomplete stats (only Pre)
	editStats := summary.ToolStats["Edit"]
	if editStats.TotalCalls != 1 {
		t.Errorf("Expected 1 call for Edit, got %d", editStats.TotalCalls)
	}
	if editStats.SuccessCount != 0 {
		t.Errorf("Expected 0 successes for Edit (no Post event), got %d", editStats.SuccessCount)
	}
}

// TestSuccessDetection tests various scenarios for detecting tool success/failure
func TestSuccessDetection(t *testing.T) {
	// Setup test environment
	tempDir := t.TempDir()
	os.Setenv("THE_STARTUP_PATH", tempDir)
	defer os.Unsetenv("THE_STARTUP_PATH")

	// Create logs directory
	logsDir := filepath.Join(tempDir, "logs")
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		name           string
		toolOutput     json.RawMessage
		error          string
		errorType      string
		expectedSuccess *bool
	}{
		{
			name:           "ExplicitSuccess",
			toolOutput:     json.RawMessage(`{"success": true, "data": "test"}`),
			expectedSuccess: boolPtr(true),
		},
		{
			name:           "ExplicitFailure",
			toolOutput:     json.RawMessage(`{"success": false, "error": "failed"}`),
			expectedSuccess: boolPtr(false),
		},
		{
			name:           "ErrorField",
			toolOutput:     json.RawMessage(`{"error": "something went wrong"}`),
			expectedSuccess: boolPtr(false),
		},
		{
			name:           "StatusSuccess",
			toolOutput:     json.RawMessage(`{"status": "success", "result": "done"}`),
			expectedSuccess: boolPtr(true),
		},
		{
			name:           "StatusOK",
			toolOutput:     json.RawMessage(`{"status": "OK"}`),
			expectedSuccess: boolPtr(true),
		},
		{
			name:           "StatusCompleted",
			toolOutput:     json.RawMessage(`{"status": "completed"}`),
			expectedSuccess: boolPtr(true),
		},
		{
			name:           "ExplicitError",
			error:          "Tool execution failed",
			errorType:      "ExecutionError",
			expectedSuccess: boolPtr(false),
		},
		{
			name:            "NoOutputNoError",
			toolOutput:      nil,
			expectedSuccess: nil, // Can't determine
		},
		{
			name:           "OutputWithoutStatus",
			toolOutput:     json.RawMessage(`{"data": "some data", "count": 42}`),
			expectedSuccess: boolPtr(true), // Has output, no error = success
		},
	}

	for i, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hook := HookPayload{
				HookEventName: "PostToolUse",
				ToolName:      fmt.Sprintf("Tool_%d", i),
				SessionID:     "test-session",
				Timestamp:     time.Now().UTC().Format(time.RFC3339),
				ToolInput:     json.RawMessage(`{"test": true}`),
				ToolResponse:  tc.toolOutput,  // Use ToolResponse instead of ToolOutput
				Error:         tc.error,
				ErrorType:     tc.errorType,
			}

			hookJSON, _ := json.Marshal(hook)
			ProcessHook(bytes.NewReader(hookJSON))

			// Read back the metric
			filter := MetricsFilter{
				StartDate: time.Now().UTC().Truncate(24 * time.Hour),
				EndDate:   time.Now().UTC().Add(24 * time.Hour),
				ToolNames: []string{fmt.Sprintf("Tool_%d", i)},
			}

			metricsChan, err := StreamMetrics(filter)
			if err != nil {
				t.Fatal(err)
			}

			var entry *MetricsEntry
			for e := range metricsChan {
				if e.HookEvent == "PostToolUse" {
					entry = &e
					break
				}
			}

			if entry == nil {
				t.Fatal("Could not find PostToolUse entry")
			}

			// Debug output
			if tc.name == "NoOutputNoError" {
				t.Logf("Test %s: ToolOutput=%s (len=%d), Success=%v", tc.name, string(entry.ToolOutput), len(entry.ToolOutput), entry.Success)
			}

			// Compare success detection
			if tc.expectedSuccess == nil {
				if entry.Success != nil {
					t.Errorf("Expected nil success, got %v", *entry.Success)
				}
			} else {
				if entry.Success == nil {
					t.Errorf("Expected success=%v, got nil", *tc.expectedSuccess)
				} else if *entry.Success != *tc.expectedSuccess {
					t.Errorf("Expected success=%v, got %v", *tc.expectedSuccess, *entry.Success)
				}
			}
		})
	}
}

// TestStreamingAggregation tests real-time streaming aggregation
func TestStreamingAggregation(t *testing.T) {
	// Setup test environment
	tempDir := t.TempDir()
	os.Setenv("THE_STARTUP_PATH", tempDir)
	defer os.Unsetenv("THE_STARTUP_PATH")

	// Create logs directory
	logsDir := filepath.Join(tempDir, "logs")
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create a stream of metrics
	inputChan := make(chan MetricsEntry, 100)
	resultsChan := make(chan *MetricsSummary, 10)
	
	// Start streaming aggregation
	go StreamAggregateMetrics(inputChan, resultsChan, 100*time.Millisecond, nil)

	// Send metrics in batches
	go func() {
		for i := 0; i < 50; i++ {
			entry := MetricsEntry{
				ToolName:  fmt.Sprintf("Tool%d", i%5),
				HookEvent: "PreToolUse",
				Timestamp: time.Now().UTC(),
				SessionID: fmt.Sprintf("session-%d", i%3),
				ToolID:    fmt.Sprintf("tool_%d", i),
				ToolInput: json.RawMessage(`{"batch": 1}`),
			}
			inputChan <- entry

			// Send corresponding Post event
			postEntry := MetricsEntry{
				ToolName:   fmt.Sprintf("Tool%d", i%5),
				HookEvent:  "PostToolUse", 
				Timestamp:  time.Now().UTC().Add(50 * time.Millisecond),
				SessionID:  fmt.Sprintf("session-%d", i%3),
				ToolID:     fmt.Sprintf("tool_%d", i),
				ToolInput:  json.RawMessage(`{"batch": 1}`),
				ToolOutput: json.RawMessage(`{"success": true}`),
				Success:    boolPtr(true),
				DurationMs: int64Ptr(50),
			}
			inputChan <- postEntry
		}
		close(inputChan)
	}()

	// Collect results
	summaries := make([]*MetricsSummary, 0)
	timeout := time.After(2 * time.Second)
	
	for {
		select {
		case summary, ok := <-resultsChan:
			if !ok {
				goto done
			}
			if summary != nil {
				summaries = append(summaries, summary)
			}
		case <-timeout:
			goto done
		}
	}

done:
	// Verify we received at least one summary
	if len(summaries) == 0 {
		t.Fatal("No summaries received from streaming aggregation")
	}

	// Check last summary has expected data
	lastSummary := summaries[len(summaries)-1]
	if lastSummary.TotalCalls == 0 {
		t.Error("Last summary has no calls")
	}
	if lastSummary.UniqueSessions != 3 {
		t.Errorf("Expected 3 unique sessions, got %d", lastSummary.UniqueSessions)
	}
	if len(lastSummary.ToolStats) != 5 {
		t.Errorf("Expected 5 tools in stats, got %d", len(lastSummary.ToolStats))
	}
}

// TestLargeScaleProcessing tests the pipeline with a large volume of data
func TestLargeScaleProcessing(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping large scale test in short mode")
	}

	// Setup test environment
	tempDir := t.TempDir()
	os.Setenv("THE_STARTUP_PATH", tempDir)
	defer os.Unsetenv("THE_STARTUP_PATH")

	// Create logs directory
	logsDir := filepath.Join(tempDir, "logs")
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Generate a large number of events
	numEvents := 10000
	tools := []string{"Edit", "Write", "Read", "Bash", "Search", "Test", "Build", "Run"}
	sessions := make([]string, 50)
	for i := range sessions {
		sessions[i] = fmt.Sprintf("session-%d", i)
	}

	startTime := time.Now()
	
	// Use goroutines to speed up generation
	var wg sync.WaitGroup
	batchSize := 1000
	numBatches := numEvents / batchSize

	for batch := 0; batch < numBatches; batch++ {
		wg.Add(1)
		go func(batchNum int) {
			defer wg.Done()
			
			for i := 0; i < batchSize; i++ {
				eventNum := batchNum*batchSize + i
				tool := tools[eventNum%len(tools)]
				session := sessions[eventNum%len(sessions)]
				
				// Create Pre event
				preHook := HookPayload{
					HookEventName: "PreToolUse",
					ToolName:      tool,
					SessionID:     session,
					Timestamp:     time.Now().UTC().Format(time.RFC3339),
					ToolInput:     json.RawMessage(fmt.Sprintf(`{"event": %d}`, eventNum)),
				}
				preJSON, _ := json.Marshal(preHook)
				ProcessHook(bytes.NewReader(preJSON))

				// Create Post event with varying outcomes
				postHook := HookPayload{
					HookEventName: "PostToolUse",
					ToolName:      tool,
					SessionID:     session,
					Timestamp:     time.Now().UTC().Add(time.Duration(eventNum%100) * time.Millisecond).Format(time.RFC3339),
					ToolInput:     json.RawMessage(fmt.Sprintf(`{"event": %d}`, eventNum)),
				}

				// 90% success rate
				if eventNum%10 != 0 {
					postHook.ToolResponse = json.RawMessage(`{"success": true}`)
				} else {
					postHook.Error = fmt.Sprintf("Error for event %d", eventNum)
					postHook.ErrorType = "TestError"
				}

				postJSON, _ := json.Marshal(postHook)
				ProcessHook(bytes.NewReader(postJSON))
			}
		}(batch)
	}

	wg.Wait()
	processingTime := time.Since(startTime)

	t.Logf("Generated %d events in %v", numEvents*2, processingTime)

	// Read and aggregate all metrics
	aggregateStart := time.Now()
	
	filter := MetricsFilter{
		StartDate: time.Now().UTC().Truncate(24 * time.Hour),
		EndDate:   time.Now().UTC().Add(24 * time.Hour),
	}

	metricsChan, err := StreamMetrics(filter)
	if err != nil {
		t.Fatal(err)
	}

	// Count metrics while streaming to aggregator
	count := 0
	aggregateChan := make(chan MetricsEntry, 100)
	
	go func() {
		for entry := range metricsChan {
			count++
			aggregateChan <- entry
		}
		close(aggregateChan)
	}()

	summary := AggregateMetrics(aggregateChan, nil)
	aggregateTime := time.Since(aggregateStart)

	t.Logf("Aggregated %d metrics in %v", count, aggregateTime)

	// Verify results
	expectedEvents := numEvents * 2 // Pre and Post for each
	if count != expectedEvents {
		t.Errorf("Expected %d events, got %d", expectedEvents, count)
	}

	if summary.TotalCalls != numEvents {
		t.Errorf("Expected %d total calls, got %d", numEvents, summary.TotalCalls)
	}

	if summary.UniqueSessions != len(sessions) {
		t.Errorf("Expected %d unique sessions, got %d", len(sessions), summary.UniqueSessions)
	}

	// Check success rate (should be around 90%)
	expectedSuccessRate := 90.0
	tolerance := 1.0
	if summary.SuccessRate < expectedSuccessRate-tolerance || 
	   summary.SuccessRate > expectedSuccessRate+tolerance {
		t.Errorf("Expected success rate around %.1f%%, got %.1f%%", 
			expectedSuccessRate, summary.SuccessRate)
	}

	// Verify all tools are represented
	if len(summary.ToolStats) != len(tools) {
		t.Errorf("Expected %d tools in stats, got %d", len(tools), len(summary.ToolStats))
	}

	// Performance check
	if processingTime > 10*time.Second {
		t.Logf("Warning: Processing took longer than expected: %v", processingTime)
	}
	if aggregateTime > 5*time.Second {
		t.Logf("Warning: Aggregation took longer than expected: %v", aggregateTime)
	}
}

// Helper functions
func int64Ptr(i int64) *int64 {
	return &i
}