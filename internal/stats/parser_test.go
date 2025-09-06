package stats

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

// Sample JSONL data for testing
const sampleJSONL = `{"type":"user","session_id":"sess_001","timestamp":"2025-01-01T10:00:00Z","content":{"text":"Hello, can you help me?"}}
{"type":"assistant","session_id":"sess_001","timestamp":"2025-01-01T10:00:05Z","content":{"text":"I'll help you with that.","tool_uses":[{"id":"tool_001","name":"Bash","parameters":{"command":"ls -la","description":"List files"}}]}}
{"type":"user","session_id":"sess_001","timestamp":"2025-01-01T10:00:10Z","content":{"tool_use_id":"tool_001","tool_name":"Bash","result":{"output":"file1.txt\nfile2.txt"}}}
{"type":"system","session_id":"sess_001","timestamp":"2025-01-01T10:00:15Z","content":{"event":"hook_executed","message":"PreToolUse hook executed"}}
{"type":"summary","session_id":"sess_001","timestamp":"2025-01-01T10:00:20Z","content":{"duration":20000000000,"tools_used":1,"tokens_used":{"input":100,"output":50,"total":150},"tool_breakdown":{"Bash":1}}}
{"type":"assistant","session_id":"sess_001","timestamp":"2025-01-01T10:00:25Z","content":{"text":"Let me read that file.","tool_uses":[{"id":"tool_002","name":"Read","parameters":{"file_path":"/tmp/test.txt","limit":100}}]}}
{"type":"user","session_id":"sess_001","timestamp":"2025-01-01T10:00:30Z","content":{"tool_use_id":"tool_002","tool_name":"Read","error":{"code":"FILE_NOT_FOUND","message":"File not found"}}}
`

// Malformed JSONL for error handling tests
const malformedJSONL = `{"type":"user","session_id":"sess_002","timestamp":"2025-01-01T11:00:00Z","content":{"text":"Valid line"}}
{invalid json here}
{"type":"assistant","session_id":"sess_002","timestamp":"2025-01-01T11:00:05Z","content":{"text":"Another valid line"}}
{"type": "broken", incomplete...
{"type":"system","session_id":"sess_002","timestamp":"2025-01-01T11:00:10Z","content":{"event":"test"}}
`

func TestNewJSONLParser(t *testing.T) {
	parser := NewJSONLParser()
	if parser == nil {
		t.Fatal("NewJSONLParser returned nil")
	}
	
	// Check default options
	if parser.options.BufferSize != 64*1024 {
		t.Errorf("Expected default buffer size of 64KB, got %d", parser.options.BufferSize)
	}
	if !parser.options.ParseTools {
		t.Error("Expected ParseTools to be true by default")
	}
	if parser.options.StrictMode {
		t.Error("Expected StrictMode to be false by default")
	}
}

func TestParseStream(t *testing.T) {
	parser := NewJSONLParser()
	reader := strings.NewReader(sampleJSONL)
	
	entries, errors := parser.ParseStream(reader)
	
	var entryCount int
	var errorCount int
	entryTypes := make(map[string]int)
	
	// Collect all entries
	go func() {
		for err := range errors {
			errorCount++
			t.Logf("Error: %v", err)
		}
	}()
	
	for entry := range entries {
		entryCount++
		entryTypes[entry.Type]++
		
		// Validate each entry type
		switch entry.Type {
		case "user":
			if entry.User == nil {
				t.Errorf("User entry missing User field")
			}
		case "assistant":
			if entry.Assistant == nil {
				t.Errorf("Assistant entry missing Assistant field")
			}
		case "system":
			if entry.System == nil {
				t.Errorf("System entry missing System field")
			}
		case "summary":
			if entry.Summary == nil {
				t.Errorf("Summary entry missing Summary field")
			}
		}
	}
	
	// Verify counts
	if entryCount != 7 {
		t.Errorf("Expected 7 entries, got %d", entryCount)
	}
	if errorCount != 0 {
		t.Errorf("Expected 0 errors, got %d", errorCount)
	}
	
	// Verify type distribution
	expectedTypes := map[string]int{
		"user":      3,
		"assistant": 2,
		"system":    1,
		"summary":   1,
	}
	
	for typ, count := range expectedTypes {
		if entryTypes[typ] != count {
			t.Errorf("Expected %d %s entries, got %d", count, typ, entryTypes[typ])
		}
	}
}

func TestParseToolUses(t *testing.T) {
	parser := NewJSONLParser()
	reader := strings.NewReader(sampleJSONL)
	
	entries, _ := parser.ParseStream(reader)
	
	toolUseCount := 0
	for entry := range entries {
		if entry.Assistant != nil && len(entry.Assistant.ToolUses) > 0 {
			for _, tool := range entry.Assistant.ToolUses {
				toolUseCount++
				
				// Check that parameters are parsed
				switch tool.Name {
				case "Bash":
					if tool.BashParams == nil {
						t.Error("Bash parameters not parsed")
					} else if tool.BashParams.Command != "ls -la" {
						t.Errorf("Expected command 'ls -la', got '%s'", tool.BashParams.Command)
					}
				case "Read":
					if tool.ReadParams == nil {
						t.Error("Read parameters not parsed")
					} else if tool.ReadParams.FilePath != "/tmp/test.txt" {
						t.Errorf("Expected file path '/tmp/test.txt', got '%s'", tool.ReadParams.FilePath)
					}
				}
			}
		}
	}
	
	if toolUseCount != 2 {
		t.Errorf("Expected 2 tool uses, got %d", toolUseCount)
	}
}

func TestParseToolResults(t *testing.T) {
	parser := NewJSONLParser()
	reader := strings.NewReader(sampleJSONL)
	
	entries, _ := parser.ParseStream(reader)
	
	toolResultCount := 0
	errorCount := 0
	
	for entry := range entries {
		if entry.User != nil && entry.User.IsToolResult {
			toolResultCount++
			if entry.User.Error != nil {
				errorCount++
			}
		}
	}
	
	if toolResultCount != 2 {
		t.Errorf("Expected 2 tool results, got %d", toolResultCount)
	}
	if errorCount != 1 {
		t.Errorf("Expected 1 tool error, got %d", errorCount)
	}
}

func TestParseSummary(t *testing.T) {
	parser := NewJSONLParser()
	reader := strings.NewReader(sampleJSONL)
	
	entries, _ := parser.ParseStream(reader)
	
	var summary *SummaryMessage
	for entry := range entries {
		if entry.Summary != nil {
			summary = entry.Summary
			break
		}
	}
	
	if summary == nil {
		t.Fatal("No summary message found")
	}
	
	if summary.TokensUsed.Total != 150 {
		t.Errorf("Expected total tokens 150, got %d", summary.TokensUsed.Total)
	}
	if summary.ToolsUsed != 1 {
		t.Errorf("Expected 1 tool used, got %d", summary.ToolsUsed)
	}
	if summary.ToolBreakdown["Bash"] != 1 {
		t.Errorf("Expected Bash used 1 time, got %d", summary.ToolBreakdown["Bash"])
	}
}

func TestMalformedJSON(t *testing.T) {
	parser := NewJSONLParser()
	reader := strings.NewReader(malformedJSONL)
	
	// Set options to collect errors but not fail
	parser.SetOptions(ParseOptions{
		BufferSize:    64 * 1024,
		ParseTools:    true,
		StrictMode:    false,
		CollectErrors: true,
	})
	
	entries, errors := parser.ParseStream(reader)
	
	entryCount := 0
	errorCount := 0
	
	// Collect errors in a separate goroutine
	errChan := make(chan int)
	go func() {
		count := 0
		for err := range errors {
			count++
			t.Logf("Collected error: %v", err)
		}
		errChan <- count
	}()
	
	// Collect entries
	for entry := range entries {
		entryCount++
		t.Logf("Parsed entry type: %s", entry.Type)
	}
	
	errorCount = <-errChan
	
	// Should parse valid lines and skip malformed ones
	if entryCount != 3 {
		t.Errorf("Expected 3 valid entries, got %d", entryCount)
	}
	if errorCount != 2 {
		t.Errorf("Expected 2 errors collected, got %d", errorCount)
	}
}

func TestStrictMode(t *testing.T) {
	parser := NewJSONLParser()
	reader := strings.NewReader(malformedJSONL)
	
	// Enable strict mode
	parser.SetOptions(ParseOptions{
		BufferSize: 64 * 1024,
		StrictMode: true,
	})
	
	entries, errors := parser.ParseStream(reader)
	
	entryCount := 0
	errorCount := 0
	
	// Collect errors
	go func() {
		for err := range errors {
			errorCount++
			t.Logf("Strict mode error: %v", err)
		}
	}()
	
	// Should stop at first error in strict mode
	for range entries {
		entryCount++
	}
	
	// In strict mode, parsing should stop at first error
	if entryCount > 1 {
		t.Errorf("Strict mode should stop at first error, but got %d entries", entryCount)
	}
}

func TestFilters(t *testing.T) {
	parser := NewJSONLParser()
	
	// Test time filter
	startTime := time.Date(2025, 1, 1, 10, 0, 10, 0, time.UTC)
	endTime := time.Date(2025, 1, 1, 10, 0, 20, 0, time.UTC)
	
	parser.SetFilter(FilterOptions{
		StartTime: &startTime,
		EndTime:   &endTime,
	})
	
	reader := strings.NewReader(sampleJSONL)
	entries, _ := parser.ParseStream(reader)
	
	count := 0
	for entry := range entries {
		count++
		// Verify all entries are within time range
		if entry.Timestamp.Before(startTime) || entry.Timestamp.After(endTime) {
			t.Errorf("Entry timestamp %v outside filter range", entry.Timestamp)
		}
	}
	
	// Should filter out entries outside the time range
	if count != 3 {
		t.Errorf("Expected 3 entries within time range, got %d", count)
	}
}

func TestToolFilters(t *testing.T) {
	parser := NewJSONLParser()
	
	// Test include tools filter
	parser.SetFilter(FilterOptions{
		IncludeTools: []string{"Bash"},
	})
	
	reader := strings.NewReader(sampleJSONL)
	entries, _ := parser.ParseStream(reader)
	
	toolCount := 0
	for entry := range entries {
		if entry.Assistant != nil && len(entry.Assistant.ToolUses) > 0 {
			for _, tool := range entry.Assistant.ToolUses {
				if tool.Name == "Bash" {
					toolCount++
				} else {
					t.Errorf("Unexpected tool %s when filtering for Bash only", tool.Name)
				}
			}
		}
	}
	
	if toolCount != 1 {
		t.Errorf("Expected 1 Bash tool use, got %d", toolCount)
	}
}

func TestSessionFilter(t *testing.T) {
	// Create JSONL with multiple sessions
	multiSessionJSONL := `{"type":"user","session_id":"sess_001","timestamp":"2025-01-01T10:00:00Z","content":{"text":"Session 1"}}
{"type":"user","session_id":"sess_002","timestamp":"2025-01-01T10:00:01Z","content":{"text":"Session 2"}}
{"type":"user","session_id":"sess_001","timestamp":"2025-01-01T10:00:02Z","content":{"text":"Session 1 again"}}
{"type":"user","session_id":"sess_003","timestamp":"2025-01-01T10:00:03Z","content":{"text":"Session 3"}}
`
	
	parser := NewJSONLParser()
	parser.SetFilter(FilterOptions{
		SessionIDs: []string{"sess_001", "sess_003"},
	})
	
	reader := strings.NewReader(multiSessionJSONL)
	entries, _ := parser.ParseStream(reader)
	
	sessionCounts := make(map[string]int)
	for entry := range entries {
		sessionCounts[entry.SessionID]++
	}
	
	if sessionCounts["sess_001"] != 2 {
		t.Errorf("Expected 2 entries for sess_001, got %d", sessionCounts["sess_001"])
	}
	if sessionCounts["sess_002"] != 0 {
		t.Errorf("Expected 0 entries for sess_002, got %d", sessionCounts["sess_002"])
	}
	if sessionCounts["sess_003"] != 1 {
		t.Errorf("Expected 1 entry for sess_003, got %d", sessionCounts["sess_003"])
	}
}

func TestLargeFileStreaming(t *testing.T) {
	// Generate a large JSONL stream (simulate 100MB+ file)
	var buf bytes.Buffer
	entriesPerMB := 1000 // Approximate number of entries per MB
	totalEntries := 100 * entriesPerMB
	
	for i := 0; i < totalEntries; i++ {
		entry := map[string]interface{}{
			"type":       "user",
			"session_id": fmt.Sprintf("sess_%03d", i%10),
			"timestamp":  time.Now().Add(time.Duration(i) * time.Second).Format(time.RFC3339),
			"content": map[string]interface{}{
				"text": fmt.Sprintf("This is message number %d with some padding to make it larger %s", i, strings.Repeat("x", 500)),
			},
		}
		data, _ := json.Marshal(entry)
		buf.Write(data)
		buf.WriteByte('\n')
	}
	
	parser := NewJSONLParser()
	parser.SetOptions(ParseOptions{
		BufferSize: 128 * 1024, // 128KB buffer for better performance
		ParseTools: false,       // Skip tool parsing for speed
	})
	
	reader := bytes.NewReader(buf.Bytes())
	entries, errors := parser.ParseStream(reader)
	
	// Process in streaming fashion
	count := 0
	errorCount := 0
	
	// Handle errors in separate goroutine
	go func() {
		for err := range errors {
			errorCount++
			t.Logf("Error processing large file: %v", err)
		}
	}()
	
	start := time.Now()
	for range entries {
		count++
		if count%10000 == 0 {
			t.Logf("Processed %d entries...", count)
		}
	}
	duration := time.Since(start)
	
	if count != totalEntries {
		t.Errorf("Expected %d entries, got %d", totalEntries, count)
	}
	if errorCount != 0 {
		t.Errorf("Expected 0 errors, got %d", errorCount)
	}
	
	t.Logf("Processed %d entries in %v (%.0f entries/sec)", count, duration, float64(count)/duration.Seconds())
}

func TestMaxEntriesLimit(t *testing.T) {
	parser := NewJSONLParser()
	parser.SetFilter(FilterOptions{
		MaxEntries: 3,
	})
	
	reader := strings.NewReader(sampleJSONL)
	entries, _ := parser.ParseStream(reader)
	
	count := 0
	for range entries {
		count++
	}
	
	if count != 3 {
		t.Errorf("Expected max 3 entries, got %d", count)
	}
}

func TestSampling(t *testing.T) {
	// Create 100 entries
	var jsonl strings.Builder
	for i := 0; i < 100; i++ {
		jsonl.WriteString(fmt.Sprintf(`{"type":"user","session_id":"sess_001","timestamp":"2025-01-01T10:00:%02dZ","content":{"text":"Message %d"}}`, i%60, i))
		jsonl.WriteByte('\n')
	}
	
	parser := NewJSONLParser()
	parser.SetFilter(FilterOptions{
		SampleRate: 0.1, // Sample 10% of entries
	})
	
	reader := strings.NewReader(jsonl.String())
	entries, _ := parser.ParseStream(reader)
	
	count := 0
	for range entries {
		count++
	}
	
	// With 10% sampling of 100 entries, we expect around 10 entries
	// Allow some variance due to sampling algorithm
	if count < 8 || count > 12 {
		t.Errorf("Expected approximately 10 entries with 10%% sampling, got %d", count)
	}
}

func TestEmptyLines(t *testing.T) {
	// JSONL with empty lines
	jsonlWithEmpty := `{"type":"user","session_id":"sess_001","timestamp":"2025-01-01T10:00:00Z","content":{"text":"Line 1"}}

{"type":"user","session_id":"sess_001","timestamp":"2025-01-01T10:00:01Z","content":{"text":"Line 2"}}


{"type":"user","session_id":"sess_001","timestamp":"2025-01-01T10:00:02Z","content":{"text":"Line 3"}}
`
	
	parser := NewJSONLParser()
	reader := strings.NewReader(jsonlWithEmpty)
	entries, _ := parser.ParseStream(reader)
	
	count := 0
	for range entries {
		count++
	}
	
	if count != 3 {
		t.Errorf("Expected 3 entries (ignoring empty lines), got %d", count)
	}
}

// Benchmark for performance testing
func BenchmarkParseStream(b *testing.B) {
	// Generate sample data
	var buf bytes.Buffer
	for i := 0; i < 1000; i++ {
		entry := map[string]interface{}{
			"type":       "assistant",
			"session_id": "bench_sess",
			"timestamp":  time.Now().Format(time.RFC3339),
			"content": map[string]interface{}{
				"text": "Benchmark message",
				"tool_uses": []map[string]interface{}{
					{
						"id":   fmt.Sprintf("tool_%d", i),
						"name": "Bash",
						"parameters": map[string]interface{}{
							"command":     "echo test",
							"description": "Test command",
						},
					},
				},
			},
		}
		data, _ := json.Marshal(entry)
		buf.Write(data)
		buf.WriteByte('\n')
	}
	
	data := buf.Bytes()
	parser := NewJSONLParser()
	
	b.ResetTimer()
	b.SetBytes(int64(len(data)))
	
	for i := 0; i < b.N; i++ {
		reader := bytes.NewReader(data)
		entries, _ := parser.ParseStream(reader)
		
		// Consume all entries
		for range entries {
		}
	}
}

func BenchmarkParseStreamLarge(b *testing.B) {
	// Generate 10MB of sample data
	var buf bytes.Buffer
	targetSize := 10 * 1024 * 1024 // 10MB
	entryNum := 0
	
	for buf.Len() < targetSize {
		entry := map[string]interface{}{
			"type":       "user",
			"session_id": fmt.Sprintf("sess_%d", entryNum%100),
			"timestamp":  time.Now().Add(time.Duration(entryNum) * time.Second).Format(time.RFC3339),
			"content": map[string]interface{}{
				"text": fmt.Sprintf("Message %d with padding: %s", entryNum, strings.Repeat("x", 1000)),
			},
		}
		data, _ := json.Marshal(entry)
		buf.Write(data)
		buf.WriteByte('\n')
		entryNum++
	}
	
	data := buf.Bytes()
	parser := NewJSONLParser()
	parser.SetOptions(ParseOptions{
		BufferSize: 256 * 1024, // 256KB buffer for large files
		ParseTools: false,
	})
	
	b.ResetTimer()
	b.SetBytes(int64(len(data)))
	
	for i := 0; i < b.N; i++ {
		reader := bytes.NewReader(data)
		entries, _ := parser.ParseStream(reader)
		
		// Consume all entries
		count := 0
		for range entries {
			count++
		}
		
		b.Logf("Processed %d entries from %d bytes", count, len(data))
	}
}

// Test concurrent parsing of multiple files
func TestConcurrentParsing(t *testing.T) {
	// Create multiple parsers for concurrent use
	numParsers := 5
	
	// Each parser processes the same data concurrently
	processParser := func(id int) (int, error) {
		parser := NewJSONLParser()
		reader := strings.NewReader(sampleJSONL)
		entries, errors := parser.ParseStream(reader)
		
		count := 0
		errCount := 0
		
		// Collect errors
		go func() {
			for range errors {
				errCount++
			}
		}()
		
		for range entries {
			count++
		}
		
		if errCount > 0 {
			return count, fmt.Errorf("parser %d had %d errors", id, errCount)
		}
		return count, nil
	}
	
	// Run parsers concurrently
	results := make(chan int, numParsers)
	errorsChan := make(chan error, numParsers)
	
	for i := 0; i < numParsers; i++ {
		go func(id int) {
			count, err := processParser(id)
			if err != nil {
				errorsChan <- err
			} else {
				results <- count
			}
		}(i)
	}
	
	// Collect results
	totalCount := 0
	for i := 0; i < numParsers; i++ {
		select {
		case count := <-results:
			totalCount += count
		case err := <-errorsChan:
			t.Errorf("Concurrent parsing error: %v", err)
		}
	}
	
	expectedTotal := 7 * numParsers // 7 entries per parser
	if totalCount != expectedTotal {
		t.Errorf("Expected total count %d from %d parsers, got %d", expectedTotal, numParsers, totalCount)
	}
}

// Test parsing with different buffer sizes
func TestBufferSizes(t *testing.T) {
	bufferSizes := []int{
		1024,       // 1KB
		16 * 1024,  // 16KB
		64 * 1024,  // 64KB
		256 * 1024, // 256KB
	}
	
	for _, bufSize := range bufferSizes {
		t.Run(fmt.Sprintf("BufferSize_%d", bufSize), func(t *testing.T) {
			parser := NewJSONLParser()
			parser.SetOptions(ParseOptions{
				BufferSize: bufSize,
				ParseTools: true,
			})
			
			reader := strings.NewReader(sampleJSONL)
			entries, errors := parser.ParseStream(reader)
			
			count := 0
			errCount := 0
			
			go func() {
				for range errors {
					errCount++
				}
			}()
			
			for range entries {
				count++
			}
			
			if count != 7 {
				t.Errorf("With buffer size %d, expected 7 entries, got %d", bufSize, count)
			}
			if errCount != 0 {
				t.Errorf("With buffer size %d, expected 0 errors, got %d", bufSize, errCount)
			}
		})
	}
}

// Helper function to test ParseFile method
func TestParseFile(t *testing.T) {
	// Create a temporary file with test data
	tmpFile := t.TempDir() + "/test.jsonl"
	err := os.WriteFile(tmpFile, []byte(sampleJSONL), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	
	parser := NewJSONLParser()
	entries, errors := parser.ParseFile(tmpFile)
	
	count := 0
	errCount := 0
	
	go func() {
		for err := range errors {
			errCount++
			t.Logf("File parse error: %v", err)
		}
	}()
	
	for entry := range entries {
		count++
		if entry.SessionID != "sess_001" {
			t.Errorf("Unexpected session ID: %s", entry.SessionID)
		}
	}
	
	if count != 7 {
		t.Errorf("Expected 7 entries from file, got %d", count)
	}
	if errCount != 0 {
		t.Errorf("Expected 0 errors from file, got %d", errCount)
	}
}

// Test ParseFile with non-existent file
func TestParseFileNotFound(t *testing.T) {
	parser := NewJSONLParser()
	entries, errors := parser.ParseFile("/non/existent/file.jsonl")
	
	errCount := 0
	entryCount := 0
	errChan := make(chan int)
	
	// Errors should be reported
	go func() {
		count := 0
		for err := range errors {
			count++
			if !strings.Contains(err.Error(), "failed to open file") {
				t.Errorf("Expected 'failed to open file' error, got: %v", err)
			}
		}
		errChan <- count
	}()
	
	for range entries {
		entryCount++
	}
	
	errCount = <-errChan
	
	if entryCount != 0 {
		t.Errorf("Expected 0 entries for non-existent file, got %d", entryCount)
	}
	if errCount == 0 {
		t.Error("Expected error for non-existent file")
	}
}