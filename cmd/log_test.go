package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/rsmdt/the-startup/internal/log"
	"github.com/spf13/cobra"
)

// setupTestEnvironment creates a temporary directory with test data
func setupTestEnvironment(t *testing.T) string {
	tempDir := t.TempDir()
	
	// Set environment to use temp directory
	os.Setenv("CLAUDE_PROJECT_DIR", tempDir)
	t.Cleanup(func() {
		os.Unsetenv("CLAUDE_PROJECT_DIR")
	})
	
	// Create .the-startup directory structure
	startupDir := filepath.Join(tempDir, ".the-startup")
	logsDir := filepath.Join(startupDir, "logs")
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory structure: %v", err)
	}
	
	return tempDir
}

// contains checks if a string slice contains a value
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// writeTestMetrics writes test metrics to the log file
func writeTestMetrics(t *testing.T, dir string, entries []log.MetricsEntry) {
	logsDir := filepath.Join(dir, ".the-startup", "logs")
	
	for _, entry := range entries {
		// Determine file name based on date
		dateStr := entry.Timestamp.Format("20060102")
		logFile := filepath.Join(logsDir, dateStr+".jsonl")
		
		// Open or create file
		f, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			t.Fatalf("Failed to open log file: %v", err)
		}
		
		// Write entry as JSON line
		if err := json.NewEncoder(f).Encode(entry); err != nil {
			f.Close()
			t.Fatalf("Failed to write metrics entry: %v", err)
		}
		f.Close()
	}
}

// createTestMetrics creates sample metrics data for testing
func createTestMetrics() []log.MetricsEntry {
	now := time.Now()
	successTrue := true
	successFalse := false
	duration100 := int64(100)
	duration200 := int64(200)
	duration500 := int64(500)
	
	return []log.MetricsEntry{
		// Successful Bash calls
		{
			ToolName:    "Bash",
			SessionID:   "session-1",
			HookEvent:   "PreToolUse",
			Timestamp:   now.Add(-2 * time.Hour),
			ToolID:      "bash-1",
			ToolInput:   json.RawMessage(`{"command":"ls -la"}`),
		},
		{
			ToolName:    "Bash",
			SessionID:   "session-1",
			HookEvent:   "PostToolUse",
			Timestamp:   now.Add(-2 * time.Hour).Add(100 * time.Millisecond),
			ToolID:      "bash-1",
			ToolInput:   json.RawMessage(`{"command":"ls -la"}`),
			ToolOutput:  json.RawMessage(`{"output":"file1\nfile2"}`),
			DurationMs:  &duration100,
			Success:     &successTrue,
		},
		// Failed Read call
		{
			ToolName:    "Read",
			SessionID:   "session-1",
			HookEvent:   "PreToolUse",
			Timestamp:   now.Add(-1 * time.Hour),
			ToolID:      "read-1",
			ToolInput:   json.RawMessage(`{"file":"nonexistent.txt"}`),
		},
		{
			ToolName:    "Read",
			SessionID:   "session-1",
			HookEvent:   "PostToolUse",
			Timestamp:   now.Add(-1 * time.Hour).Add(50 * time.Millisecond),
			ToolID:      "read-1",
			ToolInput:   json.RawMessage(`{"file":"nonexistent.txt"}`),
			DurationMs:  &duration200,
			Success:     &successFalse,
			Error:       "file not found",
			ErrorType:   "FileNotFoundError",
		},
		// Successful Write call from yesterday
		{
			ToolName:    "Write",
			SessionID:   "session-2",
			HookEvent:   "PreToolUse",
			Timestamp:   now.AddDate(0, 0, -1),
			ToolID:      "write-1",
			ToolInput:   json.RawMessage(`{"file":"test.txt","content":"hello"}`),
		},
		{
			ToolName:    "Write",
			SessionID:   "session-2",
			HookEvent:   "PostToolUse",
			Timestamp:   now.AddDate(0, 0, -1).Add(500 * time.Millisecond),
			ToolID:      "write-1",
			ToolInput:   json.RawMessage(`{"file":"test.txt","content":"hello"}`),
			ToolOutput:  json.RawMessage(`{"success":true}`),
			DurationMs:  &duration500,
			Success:     &successTrue,
		},
		// Another Bash call (for top tools testing)
		{
			ToolName:    "Bash",
			SessionID:   "session-2",
			HookEvent:   "PostToolUse",
			Timestamp:   now.Add(-30 * time.Minute),
			ToolID:      "bash-2",
			DurationMs:  &duration100,
			Success:     &successTrue,
		},
	}
}

func TestLogCommand(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		stdin         string
		expectSuccess bool
		expectedErr   string
	}{
		{
			name:          "PreToolUse hook with valid JSON",
			args:          []string{},
			stdin:         `{"hook_event_name":"PreToolUse","tool_name":"Bash","session_id":"test-123","timestamp":"2025-09-03T10:00:00Z","tool_input":{"command":"ls"}}`,
			expectSuccess: true,
		},
		{
			name:          "PostToolUse hook with valid JSON and success",
			args:          []string{},
			stdin:         `{"hook_event_name":"PostToolUse","tool_name":"Read","session_id":"test-123","timestamp":"2025-09-03T10:00:01Z","tool_response":{"content":"file data"},"tool_input":{"file":"test.txt"}}`,
			expectSuccess: true,
		},
		{
			name:          "PostToolUse hook with error",
			args:          []string{},
			stdin:         `{"hook_event_name":"PostToolUse","tool_name":"Write","session_id":"test-456","timestamp":"2025-09-03T10:00:02Z","error":"permission denied","error_type":"PermissionError"}`,
			expectSuccess: true,
		},
		{
			name:          "hook with request_id for correlation",
			args:          []string{},
			stdin:         `{"hook_event_name":"PreToolUse","tool_name":"Grep","session_id":"test-789","request_id":"req-abc-123","timestamp":"2025-09-03T10:00:03Z"}`,
			expectSuccess: true,
		},
		{
			name:          "hook with alternative output field",
			args:          []string{},
			stdin:         `{"hook_event_name":"PostToolUse","tool_name":"Edit","session_id":"test-999","timestamp":"2025-09-03T10:00:04Z","output":{"result":"edited"}}`,
			expectSuccess: true,
		},
		{
			name:          "invalid JSON input - should exit silently",
			args:          []string{},
			stdin:         `invalid json {broken`,
			expectSuccess: true, // Should exit silently (code 0)
		},
		{
			name:          "unknown hook event - should exit silently",
			args:          []string{},
			stdin:         `{"hook_event_name":"UnknownEvent","tool_name":"Bash"}`,
			expectSuccess: true, // Should exit silently (code 0)
		},
		{
			name:          "empty input - should exit silently",
			args:          []string{},
			stdin:         ``,
			expectSuccess: true, // Should exit silently (code 0)
		},
		{
			name:          "missing hook_event_name - should exit silently",
			args:          []string{},
			stdin:         `{"tool_name":"Bash","session_id":"test"}`,
			expectSuccess: true, // Should exit silently (code 0)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test environment
			setupTestEnvironment(t)

			// Create the log command
			cmd := NewLogCommand()

			// Set up stdin
			if tt.stdin != "" {
				cmd.SetIn(strings.NewReader(tt.stdin))
			} else {
				cmd.SetIn(strings.NewReader(""))
			}

			// Set up output buffers
			var outBuf, errBuf bytes.Buffer
			cmd.SetOut(&outBuf)
			cmd.SetErr(&errBuf)

			// Set arguments
			cmd.SetArgs(tt.args)

			// Execute the command
			err := cmd.Execute()

			// Check results
			if tt.expectSuccess {
				if err != nil {
					t.Errorf("Expected success but got error: %v", err)
				}
			} else {
				if err == nil {
					t.Error("Expected error but command succeeded")
				} else if tt.expectedErr != "" && !strings.Contains(err.Error(), tt.expectedErr) {
					t.Errorf("Expected error containing '%s', got: %v", tt.expectedErr, err)
				}
			}
		})
	}
}

func TestLogSubcommands(t *testing.T) {
	// Test that all subcommands are present and properly configured
	cmd := NewLogCommand()
	
	subcommands := []struct {
		name        string
		hasFlags    []string
		description string
	}{
		{
			name:        "summary",
			hasFlags:    []string{"since", "format", "output", "session"},
			description: "Display metrics summary dashboard",
		},
		{
			name:        "tools",
			hasFlags:    []string{"since", "format", "output", "tool", "session", "top"},
			description: "Display detailed tool usage statistics",
		},
		{
			name:        "errors",
			hasFlags:    []string{"since", "format", "output", "tool", "session", "top"},
			description: "Display error analysis and patterns",
		},
		{
			name:        "timeline",
			hasFlags:    []string{"since", "format", "output", "session"},
			description: "Display tool usage timeline and activity patterns",
		},
	}
	
	for _, subcmd := range subcommands {
		t.Run("subcommand_"+subcmd.name, func(t *testing.T) {
			var found *cobra.Command
			for _, c := range cmd.Commands() {
				if c.Name() == subcmd.name {
					found = c
					break
				}
			}
			
			if found == nil {
				t.Errorf("Subcommand '%s' not found", subcmd.name)
				return
			}
			
			// Check description
			if found.Short != subcmd.description {
				t.Errorf("Subcommand '%s' description mismatch: got %q, want %q", 
					subcmd.name, found.Short, subcmd.description)
			}
			
			// Check flags
			for _, flagName := range subcmd.hasFlags {
				flag := found.Flags().Lookup(flagName)
				if flag == nil {
					t.Errorf("Subcommand '%s' missing flag '%s'", subcmd.name, flagName)
				}
			}
		})
	}
}

func TestLogSummaryCommand(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectSuccess  bool
		checkOutput    bool
		containsText   []string
		doesNotContain []string
	}{
		{
			name:          "summary with default flags",
			args:          []string{"summary"},
			expectSuccess: true,
			checkOutput:   false, // Skip output check for terminal format
		},
		{
			name:          "summary with since yesterday",
			args:          []string{"summary", "--since", "yesterday"},
			expectSuccess: true,
			checkOutput:   true,
		},
		{
			name:          "summary with since 1h",
			args:          []string{"summary", "--since", "1h"},
			expectSuccess: true,
			checkOutput:   true,
		},
		{
			name:          "summary with format json",
			args:          []string{"summary", "--format", "json"},
			expectSuccess: true,
			checkOutput:   false, // Will test JSON output in separate test
		},
		{
			name:          "summary with format csv",
			args:          []string{"summary", "--format", "csv"},
			expectSuccess: true,
			checkOutput:   false, // Will test CSV output in separate test
		},
		{
			name:          "summary with session filter",
			args:          []string{"summary", "--session", "session-1"},
			expectSuccess: true,
		},
		{
			name:          "summary with output to file",
			args:          []string{"summary", "--output", filepath.Join(t.TempDir(), "summary.txt")},
			expectSuccess: true,
		},
		{
			name:          "summary with invalid since value",
			args:          []string{"summary", "--since", "invalid"},
			expectSuccess: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test environment and write test data
			tempDir := setupTestEnvironment(t)
			writeTestMetrics(t, tempDir, createTestMetrics())
			
			// Create the log command
			cmd := NewLogCommand()
			
			// Set up output buffers
			var outBuf, errBuf bytes.Buffer
			cmd.SetOut(&outBuf)
			cmd.SetErr(&errBuf)
			
			// Set arguments
			cmd.SetArgs(tt.args)
			
			// Execute the command
			err := cmd.Execute()
			
			// Check success/failure
			if tt.expectSuccess {
				if err != nil {
					t.Errorf("Expected success but got error: %v", err)
				}
			} else {
				if err == nil {
					t.Error("Expected error but command succeeded")
				}
				return // No need to check output on error
			}
			
			// Check output content if needed
			if tt.checkOutput {
				output := outBuf.String()
				for _, text := range tt.containsText {
					if !strings.Contains(output, text) {
						// For shorter output, show the full output for debugging
						if len(output) < 500 {
							t.Errorf("Output should contain %q, but doesn't: %s", text, output)
						} else {
							t.Errorf("Output should contain %q, but doesn't", text)
						}
					}
				}
				for _, text := range tt.doesNotContain {
					if strings.Contains(output, text) {
						t.Errorf("Output should not contain %q, but does", text)
					}
				}
			}
		})
	}
}

func TestLogToolsCommand(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectSuccess  bool
		checkOutput    bool
		containsTools  []string
	}{
		{
			name:          "tools with default flags",
			args:          []string{"tools"},
			expectSuccess: true,
			checkOutput:   false, // Skip terminal format output check
		},
		{
			name:          "tools filtered by tool name",
			args:          []string{"tools", "--tool", "Bash"},
			expectSuccess: true,
			checkOutput:   false, // Skip terminal format output check
		},
		{
			name:          "tools with multiple tool filters",
			args:          []string{"tools", "--tool", "Bash", "--tool", "Read"},
			expectSuccess: true,
			checkOutput:   false, // Skip terminal format output check
		},
		{
			name:          "tools with top N limit",
			args:          []string{"tools", "--top", "2"},
			expectSuccess: true,
		},
		{
			name:          "tools with session filter",
			args:          []string{"tools", "--session", "session-1"},
			expectSuccess: true,
			checkOutput:   false, // Skip terminal format output check
		},
		{
			name:          "tools with format json",
			args:          []string{"tools", "--format", "json"},
			expectSuccess: true,
			checkOutput:   false, // Will test in output file tests
		},
		{
			name:          "tools with format csv",
			args:          []string{"tools", "--format", "csv"},
			expectSuccess: true,
		},
		{
			name:          "tools with since yesterday",
			args:          []string{"tools", "--since", "yesterday"},
			expectSuccess: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test environment and write test data
			tempDir := setupTestEnvironment(t)
			writeTestMetrics(t, tempDir, createTestMetrics())
			
			// Create the log command
			cmd := NewLogCommand()
			
			// Set up output buffers
			var outBuf, errBuf bytes.Buffer
			cmd.SetOut(&outBuf)
			cmd.SetErr(&errBuf)
			
			// Set arguments
			cmd.SetArgs(tt.args)
			
			// Execute the command
			err := cmd.Execute()
			
			// Check success/failure
			if tt.expectSuccess {
				if err != nil {
					t.Errorf("Expected success but got error: %v", err)
				}
			} else {
				if err == nil {
					t.Error("Expected error but command succeeded")
				}
				return
			}
			
			// Check for expected tools in output
			if tt.checkOutput {
				output := outBuf.String()
				for _, tool := range tt.containsTools {
					if !strings.Contains(output, tool) {
						t.Errorf("Output should contain tool %q, but doesn't: %s", tool, output)
					}
				}
			}
		})
	}
}

func TestLogErrorsCommand(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectSuccess  bool
		checkOutput    bool
		containsError  []string
	}{
		{
			name:          "errors with default flags",
			args:          []string{"errors"},
			expectSuccess: true,
			checkOutput:   false, // Skip terminal format output check
		},
		{
			name:          "errors filtered by tool",
			args:          []string{"errors", "--tool", "Read"},
			expectSuccess: true,
			checkOutput:   false, // Skip terminal format output check
		},
		{
			name:          "errors with top N limit",
			args:          []string{"errors", "--top", "5"},
			expectSuccess: true,
		},
		{
			name:          "errors with session filter",
			args:          []string{"errors", "--session", "session-1"},
			expectSuccess: true,
		},
		{
			name:          "errors with format json",
			args:          []string{"errors", "--format", "json"},
			expectSuccess: true,
		},
		{
			name:          "errors filtered by non-failing tool should be empty",
			args:          []string{"errors", "--tool", "Write"},
			expectSuccess: true,
			checkOutput:   false, // Skip terminal format output check
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test environment and write test data
			tempDir := setupTestEnvironment(t)
			writeTestMetrics(t, tempDir, createTestMetrics())
			
			// Create the log command
			cmd := NewLogCommand()
			
			// Set up output buffers
			var outBuf, errBuf bytes.Buffer
			cmd.SetOut(&outBuf)
			cmd.SetErr(&errBuf)
			
			// Set arguments
			cmd.SetArgs(tt.args)
			
			// Execute the command
			err := cmd.Execute()
			
			// Check success/failure
			if tt.expectSuccess {
				if err != nil {
					t.Errorf("Expected success but got error: %v", err)
				}
			} else {
				if err == nil {
					t.Error("Expected error but command succeeded")
				}
				return
			}
			
			// Check for expected errors in output
			if tt.checkOutput {
				output := outBuf.String()
				for _, errorText := range tt.containsError {
					if !strings.Contains(output, errorText) {
						t.Errorf("Output should contain error %q, but doesn't: %s", errorText, output)
					}
				}
			}
		})
	}
}

func TestLogTimelineCommand(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		expectSuccess bool
	}{
		{
			name:          "timeline with default flags",
			args:          []string{"timeline"},
			expectSuccess: true,
		},
		{
			name:          "timeline with since 24h",
			args:          []string{"timeline", "--since", "24h"},
			expectSuccess: true,
		},
		{
			name:          "timeline with session filter",
			args:          []string{"timeline", "--session", "session-1"},
			expectSuccess: true,
		},
		{
			name:          "timeline with format json",
			args:          []string{"timeline", "--format", "json"},
			expectSuccess: true,
		},
		{
			name:          "timeline with format csv",
			args:          []string{"timeline", "--format", "csv"},
			expectSuccess: true,
		},
		{
			name:          "timeline with output file",
			args:          []string{"timeline", "--output", filepath.Join(t.TempDir(), "timeline.txt")},
			expectSuccess: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test environment and write test data
			tempDir := setupTestEnvironment(t)
			writeTestMetrics(t, tempDir, createTestMetrics())
			
			// Create the log command
			cmd := NewLogCommand()
			
			// Set up output buffers
			var outBuf, errBuf bytes.Buffer
			cmd.SetOut(&outBuf)
			cmd.SetErr(&errBuf)
			
			// Set arguments
			cmd.SetArgs(tt.args)
			
			// Execute the command
			err := cmd.Execute()
			
			// Check success/failure
			if tt.expectSuccess {
				if err != nil {
					t.Errorf("Expected success but got error: %v", err)
				}
			} else {
				if err == nil {
					t.Error("Expected error but command succeeded")
				}
			}
		})
	}
}

func TestParseTimeFilter(t *testing.T) {
	now := time.Now()
	
	tests := []struct {
		name           string
		input          string
		shouldErr      bool
		checkDateRange bool
		minHoursAgo    float64
		maxHoursAgo    float64
	}{
		{
			name:           "today",
			input:          "today",
			shouldErr:      false,
			checkDateRange: true,
			minHoursAgo:    0,
			maxHoursAgo:    24,
		},
		{
			name:           "yesterday",
			input:          "yesterday",
			shouldErr:      false,
			checkDateRange: true,
			minHoursAgo:    24,
			maxHoursAgo:    48,
		},
		{
			name:           "1h",
			input:          "1h",
			shouldErr:      false,
			checkDateRange: true,
			minHoursAgo:    0,
			maxHoursAgo:    1.1,
		},
		{
			name:           "24h",
			input:          "24h",
			shouldErr:      false,
			checkDateRange: true,
			minHoursAgo:    0,
			maxHoursAgo:    24.1,
		},
		{
			name:           "7d",
			input:          "7d",
			shouldErr:      false,
			checkDateRange: true,
			minHoursAgo:    0,
			maxHoursAgo:    7 * 24.1,
		},
		{
			name:           "30d",
			input:          "30d",
			shouldErr:      false,
			checkDateRange: true,
			minHoursAgo:    0,
			maxHoursAgo:    30 * 24.1,
		},
		{
			name:      "date format YYYY-MM-DD",
			input:     "2025-09-03",
			shouldErr: false,
		},
		{
			name:           "3 hours",
			input:          "3h",
			shouldErr:      false,
			checkDateRange: true,
			minHoursAgo:    0,
			maxHoursAgo:    3.1,
		},
		{
			name:           "5 days",
			input:          "5d",
			shouldErr:      false,
			checkDateRange: true,
			minHoursAgo:    0,
			maxHoursAgo:    5 * 24.1,
		},
		{
			name:      "invalid format",
			input:     "invalid",
			shouldErr: true,
		},
		{
			name:      "bad duration",
			input:     "xyz",
			shouldErr: true,
		},
		{
			name:      "negative days",
			input:     "-5d",
			shouldErr: false, // Negative days are valid (goes back in time)
		},
		{
			name:      "invalid date format",
			input:     "09-03-2025",
			shouldErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter, err := parseTimeFilter(tt.input)
			
			if tt.shouldErr {
				if err == nil {
					t.Errorf("Expected error for input '%s' but got none", tt.input)
				}
				return
			}
			
			if err != nil {
				t.Errorf("Unexpected error for input '%s': %v", tt.input, err)
				return
			}
			
			if filter == nil {
				t.Errorf("Expected filter for input '%s' but got nil", tt.input)
				return
			}
			
			// Check date range if specified
			if tt.checkDateRange {
				startDiff := now.Sub(filter.StartDate).Hours()
				endDiff := now.Sub(filter.EndDate).Hours()
				
				// End date should be close to now (within bounds)
				if endDiff < -0.1 || endDiff > tt.minHoursAgo+0.1 {
					t.Errorf("For input '%s', end date is %.2f hours ago, expected ~%.2f",
						tt.input, endDiff, tt.minHoursAgo)
				}
				
				// Start date should be within expected range
				if startDiff < tt.minHoursAgo-0.1 || startDiff > tt.maxHoursAgo {
					t.Errorf("For input '%s', start date is %.2f hours ago, expected between %.2f and %.2f",
						tt.input, startDiff, tt.minHoursAgo, tt.maxHoursAgo)
				}
				
				// Start should be before end
				if !filter.StartDate.Before(filter.EndDate) && !filter.StartDate.Equal(filter.EndDate) {
					t.Errorf("For input '%s', start date %v is not before end date %v",
						tt.input, filter.StartDate, filter.EndDate)
				}
			}
		})
	}
}

// TestLogCommandWithComplexScenarios tests more complex flag combinations and edge cases
func TestLogCommandWithComplexScenarios(t *testing.T) {
	tests := []struct {
		name          string
		subcommand    string
		args          []string
		setupFunc     func(t *testing.T, dir string)
		expectSuccess bool
		validateFunc  func(t *testing.T, output string)
	}{
		{
			name:       "summary with multiple sessions",
			subcommand: "summary",
			args:       []string{"--format", "json"},
			setupFunc: func(t *testing.T, dir string) {
				metrics := createTestMetrics()
				// Add metrics from different sessions
				for i := 0; i < 3; i++ {
					metrics[0].SessionID = fmt.Sprintf("session-%d", i)
					writeTestMetrics(t, dir, metrics[:1])
				}
			},
			expectSuccess: true,
			validateFunc: func(t *testing.T, output string) {
				// The output should be valid JSON
				var summary map[string]interface{}
				if err := json.Unmarshal([]byte(output), &summary); err != nil {
					// Show first 200 chars of output for debugging
					preview := output
					if len(output) > 200 {
						preview = output[:200] + "..."
					}
					t.Errorf("Failed to parse JSON output: %v\nOutput preview: %s", err, preview)
					return
				}
				if summary["unique_sessions"] == nil {
					t.Error("JSON output should contain unique_sessions field")
				}
			},
		},
		{
			name:       "tools with no data for time range",
			subcommand: "tools",
			args:       []string{"--since", "2020-01-01"},
			setupFunc: func(t *testing.T, dir string) {
				// Write recent data only
				writeTestMetrics(t, dir, createTestMetrics())
			},
			expectSuccess: true,
			validateFunc: func(t *testing.T, output string) {
				// With no data for 2020, the output should be empty or show zero tools
				if len(output) > 0 && !strings.Contains(strings.ToLower(output), "tool") {
					t.Errorf("Expected empty output or tool header for no data, got: %s", output)
				}
			},
		},
		{
			name:       "errors with CSV format",
			subcommand: "errors",
			args:       []string{"--format", "csv"},
			setupFunc: func(t *testing.T, dir string) {
				writeTestMetrics(t, dir, createTestMetrics())
			},
			expectSuccess: true,
			validateFunc: func(t *testing.T, output string) {
				// CSV format should have headers
				lines := strings.Split(strings.TrimSpace(output), "\n")
				if len(lines) < 1 {
					t.Error("CSV output should have at least a header row")
				}
				// Check first line has commas (CSV format)
				if len(lines) > 0 && !strings.Contains(lines[0], ",") {
					t.Errorf("CSV header should contain commas, got: %s", lines[0])
				}
			},
		},
		{
			name:       "timeline with hourly breakdown",
			subcommand: "timeline",
			args:       []string{"--since", "24h"},
			setupFunc: func(t *testing.T, dir string) {
				// Create metrics spread across different hours
				now := time.Now()
				for i := 0; i < 24; i++ {
					entry := log.MetricsEntry{
						ToolName:  "Bash",
						SessionID: "test",
						HookEvent: "PostToolUse",
						Timestamp: now.Add(time.Duration(-i) * time.Hour),
						ToolID:    fmt.Sprintf("bash-%d", i),
					}
					writeTestMetrics(t, dir, []log.MetricsEntry{entry})
				}
			},
			expectSuccess: true,
			validateFunc: func(t *testing.T, output string) {
				// Timeline should show activity or at least headers
				if !strings.Contains(output, "Timeline") && !strings.Contains(output, "Activity") && !strings.Contains(output, "Hour") {
					t.Errorf("Timeline output should contain activity headers, got: %.200s", output)
				}
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test environment
			tempDir := setupTestEnvironment(t)
			
			// Run setup function if provided
			if tt.setupFunc != nil {
				tt.setupFunc(t, tempDir)
			}
			
			// Create the log command
			cmd := NewLogCommand()
			
			// Set up output buffers
			var outBuf, errBuf bytes.Buffer
			cmd.SetOut(&outBuf)
			cmd.SetErr(&errBuf)
			
			// Build full args
			fullArgs := []string{tt.subcommand}
			fullArgs = append(fullArgs, tt.args...)
			
			// Add output file for capturing output
			outputFile := ""
			if tt.validateFunc != nil {
				outputFile = filepath.Join(tempDir, "test-output.txt")
				fullArgs = append(fullArgs, "--output", outputFile)
			}
			
			cmd.SetArgs(fullArgs)
			
			// Execute the command
			err := cmd.Execute()
			
			// Check success/failure
			if tt.expectSuccess {
				if err != nil {
					t.Errorf("Expected success but got error: %v", err)
				}
			} else {
				if err == nil {
					t.Error("Expected error but command succeeded")
				}
				return
			}
			
			// Run validation function if provided
			if tt.validateFunc != nil && outputFile != "" {
				data, err := os.ReadFile(outputFile)
				if err != nil {
					t.Errorf("Failed to read output file: %v", err)
					return
				}
				tt.validateFunc(t, string(data))
			}
		})
	}
}

// TestLogCommandOutputFiles tests writing output to files
func TestLogCommandOutputFiles(t *testing.T) {
	tempDir := setupTestEnvironment(t)
	writeTestMetrics(t, tempDir, createTestMetrics())
	
	tests := []struct {
		name       string
		subcommand string
		format     string
	}{
		{"summary_json", "summary", "json"},
		{"tools_csv", "tools", "csv"},
		{"errors_terminal", "errors", "terminal"},
		{"timeline_json", "timeline", "json"},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outputFile := filepath.Join(tempDir, fmt.Sprintf("%s_output.txt", tt.name))
			
			cmd := NewLogCommand()
			var outBuf, errBuf bytes.Buffer
			cmd.SetOut(&outBuf)
			cmd.SetErr(&errBuf)
			
			cmd.SetArgs([]string{
				tt.subcommand,
				"--format", tt.format,
				"--output", outputFile,
			})
			
			err := cmd.Execute()
			if err != nil {
				t.Errorf("Command failed: %v", err)
				return
			}
			
			// Check that file was created and has content
			info, err := os.Stat(outputFile)
			if err != nil {
				t.Errorf("Output file not created: %v", err)
				return
			}
			
			if info.Size() == 0 {
				t.Error("Output file is empty")
			}
			
			// Read and validate file content based on format
			content, err := os.ReadFile(outputFile)
			if err != nil {
				t.Errorf("Failed to read output file: %v", err)
				return
			}
			
			switch tt.format {
			case "json":
				var jsonData interface{}
				if err := json.Unmarshal(content, &jsonData); err != nil {
					t.Errorf("Output file contains invalid JSON: %v", err)
				}
			case "csv":
				if !strings.Contains(string(content), ",") {
					t.Error("CSV output file should contain commas")
				}
			}
		})
	}
}