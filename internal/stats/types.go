// Package stats provides data structures and interfaces for analyzing Claude Code's
// native JSONL log format. It defines types for parsing log entries, correlating
// tool invocations with results, aggregating statistics, and displaying insights.
package stats

import (
	"encoding/json"
	"io"
	"time"
)

// ClaudeLogEntry represents a single line from Claude Code's JSONL transcript
// This is the root structure for parsing each JSONL line
type ClaudeLogEntry struct {
	Type       string          `json:"type"`       // "user", "assistant", "system", "summary"
	SessionID  string          `json:"session_id"` // Claude Code session identifier
	Timestamp  time.Time       `json:"timestamp"`  // When this entry was created
	Content    json.RawMessage `json:"content"`    // Raw content, parsed based on Type
	
	// Parsed content based on Type
	User      *UserMessage      `json:"-"`
	Assistant *AssistantMessage `json:"-"`
	System    *SystemMessage    `json:"-"`
	Summary   *SummaryMessage   `json:"-"`
}

// CommandInfo represents extracted command information from assistant messages
// Commands are identified by <command-name> tags in Claude's responses
type CommandInfo struct {
	Name      string    `json:"name"`       // Command name from <command-name> tags
	Args      string    `json:"args"`       // Command arguments
	Timestamp time.Time `json:"timestamp"`  // When command was issued
	SessionID string    `json:"session_id"` // Session containing this command
}

// ConversationMessage represents a user or assistant message in the conversation
// Used for analyzing conversation patterns and content
type ConversationMessage struct {
	Role          string    `json:"role"`            // "user" or "assistant"
	Content       string    `json:"content"`         // Text content
	Timestamp     time.Time `json:"timestamp"`       // When message was sent
	TokenCount    int       `json:"token_count"`     // Token usage for this message
	HasCodeBlocks bool      `json:"has_code_blocks"` // Whether message contains code
	HasTools      bool      `json:"has_tools"`       // Whether message uses tools
}

// UserMessage represents user input or tool results
type UserMessage struct {
	// Role from the message
	Role string `json:"role,omitempty"`
	
	// Direct user input
	Text string `json:"text,omitempty"`
	
	// Tool result response
	ToolUseID string          `json:"tool_use_id,omitempty"` // Correlates with tool invocation
	ToolName  string          `json:"tool_name,omitempty"`   // Name of tool that was executed
	Output    string          `json:"output,omitempty"`      // Tool execution output
	Result    json.RawMessage `json:"result,omitempty"`      // Tool execution result
	Error     *string         `json:"error,omitempty"`       // Tool execution error
	
	// Agent detection (populated during processing)
	DetectedAgent       string  `json:"detected_agent,omitempty"`
	DetectionConfidence float64 `json:"detection_confidence,omitempty"`

	// Metadata
	IsToolResult bool      `json:"is_tool_result,omitempty"`
	Timestamp    time.Time `json:"timestamp"`
}

// AssistantMessage represents Claude's responses and tool invocations
type AssistantMessage struct {
	// Message ID and role
	ID   string `json:"id,omitempty"`
	Role string `json:"role,omitempty"`
	
	// Direct assistant response
	Text string `json:"text,omitempty"`
	
	// Tool invocations
	ToolUses []ToolUse `json:"tool_uses,omitempty"`
	
	// Model and token information
	Model        string `json:"model,omitempty"`
	InputTokens  int    `json:"input_tokens,omitempty"`
	OutputTokens int    `json:"output_tokens,omitempty"`

	// Agent detection (populated during processing)
	DetectedAgent       string  `json:"detected_agent,omitempty"`
	DetectionConfidence float64 `json:"detection_confidence,omitempty"`

	// Metadata
	ModelID    string    `json:"model_id,omitempty"`
	Timestamp  time.Time `json:"timestamp"`
}

// ToolUse represents a single tool invocation by the assistant
type ToolUse struct {
	ID         string          `json:"id"`         // Unique identifier for this invocation
	Name       string          `json:"name"`       // Tool name (e.g., "Bash", "Read", "Edit")
	Parameters json.RawMessage `json:"parameters"` // Tool-specific parameters
	
	// Parsed parameters for common tools (populated during processing)
	BashParams   *BashParameters   `json:"-"`
	ReadParams   *ReadParameters   `json:"-"`
	EditParams   *EditParameters   `json:"-"`
	WriteParams  *WriteParameters  `json:"-"`
	SearchParams *SearchParameters `json:"-"`
}

// SystemMessage represents system-level events and notifications
type SystemMessage struct {
	Role      string          `json:"role,omitempty"`   // Role from the message
	Text      string          `json:"text,omitempty"`   // Text content
	Event     string          `json:"event"`            // Event type
	Message   string          `json:"message"`          // System message
	Metadata  json.RawMessage `json:"metadata"`         // Additional event data

	// Agent detection (populated during processing)
	DetectedAgent       string  `json:"detected_agent,omitempty"`
	DetectionConfidence float64 `json:"detection_confidence,omitempty"`

	Timestamp time.Time `json:"timestamp"`
}

// SummaryMessage represents session summaries and statistics
type SummaryMessage struct {
	Summary       string            `json:"summary,omitempty"`
	Duration      time.Duration     `json:"duration"`
	ToolsUsed     int               `json:"tools_used"`
	TokensUsed    TokenUsage        `json:"tokens_used"`
	ToolBreakdown map[string]int    `json:"tool_breakdown"`
	Errors        []ToolError       `json:"errors,omitempty"`
	Timestamp     time.Time         `json:"timestamp"`
}

// TokenUsage tracks token consumption
type TokenUsage struct {
	Input  int `json:"input"`
	Output int `json:"output"`
	Total  int `json:"total"`
}

// ToolError represents an error from tool execution
type ToolError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Error implements the error interface for ToolError
func (e *ToolError) Error() string {
	if e.Details != "" {
		return e.Code + ": " + e.Message + " - " + e.Details
	}
	return e.Code + ": " + e.Message
}

// Common tool parameter structures
type BashParameters struct {
	Command     string `json:"command"`
	Description string `json:"description,omitempty"`
	Timeout     int    `json:"timeout,omitempty"`
	Background  bool   `json:"run_in_background,omitempty"`
}

type ReadParameters struct {
	FilePath string `json:"file_path"`
	Offset   int    `json:"offset,omitempty"`
	Limit    int    `json:"limit,omitempty"`
}

type EditParameters struct {
	FilePath   string `json:"file_path"`
	OldString  string `json:"old_string"`
	NewString  string `json:"new_string"`
	ReplaceAll bool   `json:"replace_all,omitempty"`
}

type WriteParameters struct {
	FilePath string `json:"file_path"`
	Content  string `json:"content"`
}

type SearchParameters struct {
	Query   string   `json:"query"`
	Path    string   `json:"path,omitempty"`
	Pattern string   `json:"pattern,omitempty"`
	Glob    string   `json:"glob,omitempty"`
}

// ToolInvocation represents a complete tool execution cycle
// This correlates a tool invocation with its result
type ToolInvocation struct {
	// Invocation details
	ID           string          `json:"id"`            // Tool use ID for correlation
	Name         string          `json:"name"`          // Tool name
	Parameters   json.RawMessage `json:"parameters"`    // Input parameters
	InvokedAt    time.Time       `json:"invoked_at"`    // When tool was invoked
	
	// Result details
	CompletedAt  *time.Time      `json:"completed_at,omitempty"`  // When result was received
	DurationMs   *int64          `json:"duration_ms,omitempty"`   // Execution duration
	Result       json.RawMessage `json:"result,omitempty"`        // Tool output
	Error        *ToolError      `json:"error,omitempty"`         // Error if failed
	Success      bool            `json:"success"`                 // Whether execution succeeded
	
	// Context
	SessionID    string          `json:"session_id"`
	MessageIndex int             `json:"message_index"` // Position in conversation
}

// SessionStatistics aggregates statistics for a single session
type SessionStatistics struct {
	SessionID      string                       `json:"session_id"`
	StartTime      time.Time                    `json:"start_time"`
	EndTime        *time.Time                   `json:"end_time,omitempty"`
	Duration       time.Duration                `json:"duration"`
	
	// Message counts
	UserMessages   int                          `json:"user_messages"`
	AssistantMessages int                       `json:"assistant_messages"`
	SystemMessages int                          `json:"system_messages"`
	
	// Tool usage
	ToolInvocations []ToolInvocation            `json:"tool_invocations"`
	ToolStats      map[string]*ToolSessionStats `json:"tool_stats"`
	TotalToolCalls int                          `json:"total_tool_calls"`
	SuccessfulCalls int                         `json:"successful_calls"`
	FailedCalls    int                          `json:"failed_calls"`
	
	// Token usage
	TotalTokens    TokenUsage                   `json:"total_tokens"`
	
	// Error tracking
	Errors         []ErrorOccurrence            `json:"errors,omitempty"`
	ErrorRate      float64                      `json:"error_rate"`
}

// ToolSessionStats tracks per-tool statistics within a session
type ToolSessionStats struct {
	Name           string        `json:"name"`
	CallCount      int           `json:"call_count"`
	SuccessCount   int           `json:"success_count"`
	FailureCount   int           `json:"failure_count"`
	TotalDurationMs int64        `json:"total_duration_ms"`
	AvgDurationMs  float64       `json:"avg_duration_ms"`
	MinDurationMs  *int64        `json:"min_duration_ms,omitempty"`
	MaxDurationMs  *int64        `json:"max_duration_ms,omitempty"`
	LastUsed       time.Time     `json:"last_used"`
}

// ErrorOccurrence represents a single error event
type ErrorOccurrence struct {
	Timestamp  time.Time  `json:"timestamp"`
	ToolName   string     `json:"tool_name"`
	ToolUseID  string     `json:"tool_use_id"`
	Error      ToolError  `json:"error"`
	Context    string     `json:"context,omitempty"` // Additional context
}

// AggregatedStatistics represents statistics across multiple sessions
type AggregatedStatistics struct {
	// Time range
	Period         TimePeriod                    `json:"period"`
	
	// Session overview
	TotalSessions  int                           `json:"total_sessions"`
	ActiveSessions []string                      `json:"active_sessions"`
	AvgSessionDuration time.Duration             `json:"avg_session_duration"`
	
	// Tool usage across all sessions
	GlobalToolStats map[string]*GlobalToolStats  `json:"global_tool_stats"`
	TotalToolCalls int                           `json:"total_tool_calls"`
	OverallSuccessRate float64                   `json:"overall_success_rate"`
	
	// Error analysis
	ErrorPatterns  []ErrorPattern                `json:"error_patterns"`
	TopErrors      []ErrorFrequency              `json:"top_errors"`
	
	// Temporal patterns
	HourlyActivity []HourlyActivity              `json:"hourly_activity"`
	DailyActivity  []DailyActivity               `json:"daily_activity"`
	
	// Performance metrics
	PerformanceMetrics PerformanceMetrics        `json:"performance_metrics"`
}

// OverallStats represents the top-level statistics summary
// This is a simplified view for quick insights
type OverallStats struct {
	Tools              map[string]ToolStats     `json:"tools"`               // Tool usage statistics
	Commands           map[string]int           `json:"commands"`            // Command usage counts
	Sessions           []SessionStats           `json:"sessions"`            // Session summaries
	TotalConversations int                      `json:"total_conversations"` // Total conversation count
	TotalTokensUsed    int                      `json:"total_tokens_used"`   // Total token consumption
	TimeRange          TimeRange                `json:"time_range"`          // Analysis time range
	MostActiveHours    []int                    `json:"most_active_hours"`   // Peak usage hours (0-23)
	ErrorPatterns      []ErrorInfo              `json:"error_patterns"`      // Common error patterns
}

// SessionStats represents a simplified session statistics summary
// Used in OverallStats for session summaries
type SessionStats struct {
	SessionID       string        `json:"session_id"`
	StartTime       time.Time     `json:"start_time"`
	EndTime         time.Time     `json:"end_time"`
	Duration        time.Duration `json:"duration"`
	TotalMessages   int           `json:"total_messages"`
	ToolsUsed       int           `json:"tools_used"`
	CommandsUsed    int           `json:"commands_used"`
	TokensUsed      int           `json:"tokens_used"`
	ErrorCount      int           `json:"error_count"`
}

// ToolStats represents simplified tool statistics
// Used in OverallStats for tool summaries  
type ToolStats struct {
	Name         string        `json:"name"`
	Count        int           `json:"count"`
	TotalDuration time.Duration `json:"total_duration"`
	AvgDuration  time.Duration `json:"avg_duration"`
	MinDuration  time.Duration `json:"min_duration"`
	MaxDuration  time.Duration `json:"max_duration"`
	ErrorCount   int           `json:"error_count"`
	SuccessRate  float64       `json:"success_rate"`
	LastUsed     time.Time     `json:"last_used"`
}

// TimeRange represents a time period for analysis
type TimeRange struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// ErrorInfo represents common error information
type ErrorInfo struct {
	Pattern   string `json:"pattern"`    // Error pattern or message
	Count     int    `json:"count"`      // Occurrence count
	ToolName  string `json:"tool_name"`  // Associated tool
	FirstSeen time.Time `json:"first_seen"`
	LastSeen  time.Time `json:"last_seen"`
}

// GlobalToolStats aggregates tool statistics across all sessions
type GlobalToolStats struct {
	Name              string         `json:"name"`
	TotalCalls        int            `json:"total_calls"`
	UniqueUsers       int            `json:"unique_users"` // Unique session count
	SuccessCount      int            `json:"success_count"`
	FailureCount      int            `json:"failure_count"`
	SuccessRate       float64        `json:"success_rate"`
	
	// Duration statistics
	TotalDurationMs   int64          `json:"total_duration_ms"`
	AvgDurationMs     float64        `json:"avg_duration_ms"`
	MedianDurationMs  float64        `json:"median_duration_ms"`
	P95DurationMs     float64        `json:"p95_duration_ms"`
	P99DurationMs     float64        `json:"p99_duration_ms"`
	MinDurationMs     int64          `json:"min_duration_ms"`
	MaxDurationMs     int64          `json:"max_duration_ms"`
	
	// Error breakdown
	ErrorTypes        map[string]int `json:"error_types"`
	CommonErrors      []string       `json:"common_errors"`
	
	// Usage patterns
	PeakHour          int            `json:"peak_hour"`   // Hour of day with most usage (0-23)
	LastUsed          time.Time      `json:"last_used"`
	FirstUsed         time.Time      `json:"first_used"`
}

// TimePeriod represents a time range for analysis
type TimePeriod struct {
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
	Duration time.Duration `json:"duration"`
}

// ErrorPattern identifies common error patterns
type ErrorPattern struct {
	Pattern      string   `json:"pattern"`      // Error pattern/signature
	Count        int      `json:"count"`        // Occurrence count
	Tools        []string `json:"tools"`        // Affected tools
	FirstSeen    time.Time `json:"first_seen"`
	LastSeen     time.Time `json:"last_seen"`
	ExampleError ToolError `json:"example_error"`
}

// ErrorFrequency tracks error occurrence frequency
type ErrorFrequency struct {
	ErrorCode    string  `json:"error_code"`
	ErrorMessage string  `json:"error_message"`
	Count        int     `json:"count"`
	Percentage   float64 `json:"percentage"` // Percentage of total errors
	ToolName     string  `json:"tool_name"`
}

// HourlyActivity represents activity metrics for a specific hour
type HourlyActivity struct {
	Hour           time.Time `json:"hour"`
	ToolCalls      int       `json:"tool_calls"`
	UniqueTools    int       `json:"unique_tools"`
	UniqueSessions int       `json:"unique_sessions"`
	SuccessCount   int       `json:"success_count"`
	FailureCount   int       `json:"failure_count"`
	AvgDurationMs  float64   `json:"avg_duration_ms"`
}

// DailyActivity represents activity metrics for a specific day
type DailyActivity struct {
	Date           time.Time `json:"date"`
	Sessions       int       `json:"sessions"`
	ToolCalls      int       `json:"tool_calls"`
	UniqueTools    int       `json:"unique_tools"`
	SuccessRate    float64   `json:"success_rate"`
	TotalDurationMs int64    `json:"total_duration_ms"`
	Errors         int       `json:"errors"`
}

// PerformanceMetrics tracks overall performance indicators
type PerformanceMetrics struct {
	// Response times
	AvgResponseTimeMs   float64 `json:"avg_response_time_ms"`
	P50ResponseTimeMs   float64 `json:"p50_response_time_ms"`
	P95ResponseTimeMs   float64 `json:"p95_response_time_ms"`
	P99ResponseTimeMs   float64 `json:"p99_response_time_ms"`
	
	// Throughput
	ToolCallsPerMinute  float64 `json:"tool_calls_per_minute"`
	ToolCallsPerHour    float64 `json:"tool_calls_per_hour"`
	
	// Reliability
	OverallSuccessRate  float64 `json:"overall_success_rate"`
	ConsecutiveFailures int     `json:"consecutive_failures_max"`
	
	// Resource usage
	PeakConcurrency     int     `json:"peak_concurrency"`
	AvgConcurrency      float64 `json:"avg_concurrency"`
}

// StreamingStats provides memory-efficient statistics during streaming
type StreamingStats struct {
	// Counters (minimal memory footprint)
	TotalEntries    int64
	ProcessedCount  int64
	SkippedCount    int64
	ErrorCount      int64
	
	// Running statistics (updated incrementally)
	ToolCounters    map[string]*StreamingToolCounter
	SessionSet      map[string]bool // Track unique sessions
	
	// Time bounds
	EarliestEntry   *time.Time
	LatestEntry     *time.Time
	
	// Memory management
	MaxBufferSize   int  // Maximum entries to buffer
	FlushThreshold  int  // Flush to aggregation when reaching threshold
}

// StreamingToolCounter tracks tool statistics during streaming
type StreamingToolCounter struct {
	Count          int64
	SuccessCount   int64
	FailureCount   int64
	TotalDurationMs int64
	
	// For calculating statistics without storing all values
	MinDurationMs  *int64
	MaxDurationMs  *int64
	
	// For variance calculation (Welford's algorithm)
	M2             float64 // Sum of squared differences from mean
	Mean           float64 // Running mean
}

// FilterOptions provides filtering capabilities for log analysis
type FilterOptions struct {
	// Time filters
	StartTime      *time.Time
	EndTime        *time.Time
	
	// Tool filters
	IncludeTools   []string // If specified, only these tools
	ExcludeTools   []string // If specified, exclude these tools
	
	// Session filters
	SessionIDs     []string // If specified, only these sessions
	
	// Status filters
	SuccessOnly    bool
	FailuresOnly   bool
	
	// Performance filters
	MinDurationMs  *int64
	MaxDurationMs  *int64
	
	// Sampling
	SampleRate     float64 // 0.0 to 1.0, 1.0 means process all
	MaxEntries     int     // Maximum entries to process
}

// ParseOptions configures how JSONL entries are parsed
type ParseOptions struct {
	// Performance options
	BufferSize      int  // Buffer size for reading
	ParseTools      bool // Whether to parse tool parameters
	SkipSystemLogs  bool // Skip system messages for performance
	
	// Error handling
	StrictMode      bool // Fail on parse errors vs skip
	CollectErrors   bool // Collect parsing errors for reporting
	
	// Memory management
	MaxMemoryMB     int  // Maximum memory to use for buffering
}

// Discovery interface for finding Claude Code log files
type Discovery interface {
	// FindLogFiles discovers Claude Code JSONL log files
	// projectPath: path to the project directory (e.g., ~/.claude/projects/[project])
	// options: filtering options for time range, sessions, etc.
	// Returns list of log file paths to process
	FindLogFiles(projectPath string, options FilterOptions) ([]string, error)
	
	// GetCurrentProject attempts to detect the current project from working directory
	// Returns the project path if found, empty string otherwise
	GetCurrentProject() string
	
	// ValidateProjectPath checks if a path contains Claude Code logs
	ValidateProjectPath(path string) bool
}

// Parser interface for streaming JSONL log parsing
type Parser interface {
	// ParseStream processes a log file and returns a channel of parsed entries
	// reader: the input stream to parse
	// Returns a channel that emits ClaudeLogEntry items
	ParseStream(reader io.Reader) (<-chan ClaudeLogEntry, <-chan error)
	
	// ParseFile opens and parses a specific log file
	// filePath: path to the JSONL file
	// Returns a channel of parsed entries and error channel
	ParseFile(filePath string) (<-chan ClaudeLogEntry, <-chan error)
	
	// SetOptions updates parser configuration
	SetOptions(opts ParseOptions)
	
	// SetFilter updates filtering criteria
	SetFilter(filter FilterOptions)
}

// Aggregator interface for calculating statistics from log entries
type Aggregator interface {
	// ProcessEntry adds a single log entry to aggregation
	ProcessEntry(entry ClaudeLogEntry) error
	
	// ProcessInvocation adds a complete tool invocation
	ProcessInvocation(invocation ToolInvocation) error
	
	// GetSessionStats returns statistics for a specific session
	GetSessionStats(sessionID string) (*SessionStatistics, error)
	
	// GetOverallStats returns aggregated statistics across all sessions
	GetOverallStats() (*AggregatedStatistics, error)
	
	// GetToolStats returns statistics for a specific tool
	GetToolStats(toolName string) (*GlobalToolStats, error)
	
	// Reset clears all accumulated statistics
	Reset()
	
	// Merge combines statistics from another aggregator
	Merge(other Aggregator) error
}

// Display interface for formatting and presenting statistics
type Display interface {
	// FormatStats formats aggregated statistics for display
	// stats: the statistics to format
	// Returns formatted string output
	FormatStats(stats *AggregatedStatistics) string
	
	// FormatSession formats session-specific statistics
	FormatSession(session *SessionStatistics) string
	
	// FormatToolLeaderboard creates a tool usage leaderboard
	// toolStats: map of tool names to their statistics
	// limit: maximum number of tools to display (0 for all)
	FormatToolLeaderboard(toolStats map[string]*GlobalToolStats, limit int) string
	
	// FormatErrors formats error patterns and frequencies
	FormatErrors(errors []ErrorPattern, frequencies []ErrorFrequency) string
	
	// FormatTimeline formats temporal activity patterns
	FormatTimeline(hourly []HourlyActivity, daily []DailyActivity) string
	
	// SetOutputFormat configures the output format (text, json, csv, etc.)
	SetOutputFormat(format string) error
	
	// Write outputs formatted content to a writer
	Write(w io.Writer, content string) error
}