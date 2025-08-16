package log

// HookInput represents the JSON structure passed from Claude Code hooks
// This matches the input schema documented in the SDD
type HookInput struct {
	SessionID      string                 `json:"session_id"`
	TranscriptPath string                 `json:"transcript_path"`
	Cwd            string                 `json:"cwd"`
	HookEventName  string                 `json:"hook_event_name"`
	ToolName       string                 `json:"tool_name"`
	ToolInput      map[string]interface{} `json:"tool_input"`
	Output         string                 `json:"output,omitempty"` // Only present in PostToolUse
}

// HookData represents the simplified JSONL log entry structure
type HookData struct {
	Role      string `json:"role"`      // "user" for agent_start, "assistant" for agent_complete
	Content   string `json:"content"`   // prompt for user, response for assistant
	Timestamp string `json:"timestamp"` // UTC ISO format with Z suffix

	// Internal fields not serialized to JSON
	SessionID string `json:"-"` // Used internally for file routing
	AgentID   string `json:"-"` // Used internally for file naming
}

// LogFlags represents command-line flags for log command
type LogFlags struct {
	// Write mode flags
	Assistant bool // --assistant, -a
	User      bool // --user, -u

	// Read mode flags
	Read            bool   // --read, -r
	AgentID         string // --agent-id
	Lines           int    // --lines
	Session         string // --session
	Format          string // --format
	IncludeMetadata bool   // --include-metadata
}
