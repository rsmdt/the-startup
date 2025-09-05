package log

import (
	"encoding/json"
	"time"
)

// HookPayload represents the data structure sent by Claude Code hooks
// This matches the actual format received from PreToolUse and PostToolUse events
type HookPayload struct {
	// Core fields from Claude Code
	HookEventName  string          `json:"hook_event_name"`  // "PreToolUse" or "PostToolUse"
	ToolName       string          `json:"tool_name"`         // Name of the tool being invoked
	SessionID      string          `json:"session_id"`        // Claude Code session identifier
	TranscriptPath string          `json:"transcript_path"`   // Path to the transcript file
	CWD            string          `json:"cwd"`               // Current working directory
	Timestamp      string          `json:"timestamp"`         // ISO 8601 timestamp from Claude
	
	// Tool invocation data
	ToolInput      json.RawMessage `json:"tool_input"`        // Input parameters to the tool
	ToolResponse   json.RawMessage `json:"tool_response"`     // Output from the tool (PostToolUse only)
	ToolOutput     json.RawMessage `json:"output"`            // Alternative field for tool output
	
	// Error information (PostToolUse only)
	Error          string          `json:"error"`             // Error message if tool failed
	ErrorType      string          `json:"error_type"`        // Type/category of error
	
	// Additional metadata
	RequestID      string          `json:"request_id"`        // Unique identifier for this request
	ParentID       string          `json:"parent_id"`         // ID of parent request (for nested calls)
}

// MetricsEntry represents a single tool invocation metric stored in JSONL logs
// This structure is optimized for metrics collection and analysis
type MetricsEntry struct {
	// Core identification
	ToolName       string          `json:"tool_name"`                // Name of the tool
	SessionID      string          `json:"session_id"`               // Claude Code session ID
	TranscriptPath string          `json:"transcript_path,omitempty"` // Path to transcript file
	CWD            string          `json:"cwd,omitempty"`            // Working directory at invocation
	HookEvent      string          `json:"hook_event"`               // "PreToolUse" or "PostToolUse"
	
	// Timing information
	Timestamp      time.Time       `json:"timestamp"`                // UTC timestamp of the event
	DurationMs     *int64          `json:"duration_ms,omitempty"`    // Calculated duration (PostToolUse)
	
	// Tool invocation data
	ToolInput      json.RawMessage `json:"tool_input"`               // Raw input parameters
	ToolOutput     json.RawMessage `json:"tool_output,omitempty"`    // Raw output (PostToolUse)
	
	// Correlation and status
	ToolID         string          `json:"tool_id"`                  // Generated ID to correlate pre/post
	Success        *bool           `json:"success,omitempty"`        // Whether the tool succeeded
	
	// Error tracking
	Error          string          `json:"error,omitempty"`          // Error message if failed
	ErrorType      string          `json:"error_type,omitempty"`     // Category of error
}

// ToolStatistics represents aggregated statistics for a specific tool
type ToolStatistics struct {
	Name            string                 `json:"name"`              // Tool name
	TotalCalls      int                    `json:"total_calls"`       // Total number of invocations
	SuccessCount    int                    `json:"success_count"`     // Number of successful calls
	FailureCount    int                    `json:"failure_count"`     // Number of failed calls
	TotalDurationMs int64                  `json:"total_duration_ms"` // Total time spent
	AvgDurationMs   float64                `json:"avg_duration_ms"`   // Average execution time
	MinDurationMs   int64                  `json:"min_duration_ms"`   // Fastest execution
	MaxDurationMs   int64                  `json:"max_duration_ms"`   // Slowest execution
	ErrorTypes      map[string]int         `json:"error_types"`       // Count by error type
	LastUsed        time.Time              `json:"last_used"`         // Most recent usage
}

// MetricsSummary represents an aggregated view of metrics over a time period
type MetricsSummary struct {
	Period          TimePeriod             `json:"period"`            // Time range for this summary
	TotalCalls      int                    `json:"total_calls"`       // Total tool invocations
	UniqueSessions  int                    `json:"unique_sessions"`   // Number of unique sessions
	SuccessRate     float64                `json:"success_rate"`      // Overall success percentage
	ToolStats       map[string]*ToolStatistics `json:"tool_stats"`   // Per-tool statistics
	TopErrors       []ErrorPattern         `json:"top_errors"`        // Most common errors
	HourlyActivity  []HourlyStats          `json:"hourly_activity"`   // Activity by hour
}

// TimePeriod represents a time range for metrics analysis
type TimePeriod struct {
	Start time.Time `json:"start"` // Start of the period
	End   time.Time `json:"end"`   // End of the period
}

// String returns a human-readable representation of the time period
func (tp TimePeriod) String() string {
	duration := tp.End.Sub(tp.Start)
	if duration <= 24*time.Hour {
		return tp.Start.Format("2006-01-02")
	}
	return tp.Start.Format("2006-01-02") + " to " + tp.End.Format("2006-01-02")
}

// ErrorPattern represents a common error pattern in tool executions
type ErrorPattern struct {
	ErrorType    string `json:"error_type"`    // Type/category of error
	ErrorMessage string `json:"error_message"` // Common error message
	Count        int    `json:"count"`         // Number of occurrences
	Tools        []string `json:"tools"`       // Tools that experienced this error
}

// HourlyStats represents tool usage statistics for a specific hour
type HourlyStats struct {
	Hour         time.Time `json:"hour"`          // Start of the hour
	TotalCalls   int       `json:"total_calls"`   // Calls in this hour
	SuccessCount int       `json:"success_count"` // Successful calls
	FailureCount int       `json:"failure_count"` // Failed calls
	UniqueTools  int       `json:"unique_tools"`  // Number of different tools used
}

// MetricsFilter represents criteria for filtering metrics during analysis
type MetricsFilter struct {
	StartDate      time.Time `json:"start_date"`      // Start of date range
	EndDate        time.Time `json:"end_date"`        // End of date range
	ToolNames      []string  `json:"tool_names"`      // Filter by specific tools
	SessionIDs     []string  `json:"session_ids"`     // Filter by sessions
	SuccessOnly    bool      `json:"success_only"`    // Only successful invocations
	FailuresOnly   bool      `json:"failures_only"`   // Only failed invocations
	MinDurationMs  *int64    `json:"min_duration_ms"` // Minimum duration threshold
	MaxDurationMs  *int64    `json:"max_duration_ms"` // Maximum duration threshold
}

// Matches checks if a MetricsEntry matches the filter criteria
func (f MetricsFilter) Matches(entry MetricsEntry) bool {
	// Check date range
	if entry.Timestamp.Before(f.StartDate) || entry.Timestamp.After(f.EndDate) {
		return false
	}
	
	// Check tool name filter
	if len(f.ToolNames) > 0 {
		found := false
		for _, name := range f.ToolNames {
			if entry.ToolName == name {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	
	// Check session filter
	if len(f.SessionIDs) > 0 {
		found := false
		for _, id := range f.SessionIDs {
			if entry.SessionID == id {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	
	// Check success/failure filters
	if f.SuccessOnly && entry.Success != nil && !*entry.Success {
		return false
	}
	if f.FailuresOnly && entry.Success != nil && *entry.Success {
		return false
	}
	
	// Check duration filters
	if entry.DurationMs != nil {
		if f.MinDurationMs != nil && *entry.DurationMs < *f.MinDurationMs {
			return false
		}
		if f.MaxDurationMs != nil && *entry.DurationMs > *f.MaxDurationMs {
			return false
		}
	}
	
	return true
}