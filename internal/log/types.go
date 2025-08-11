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

// HookData represents the JSONL log entry structure
// This matches the output format from the Python hooks
type HookData struct {
	Event         string `json:"event"`          // "agent_start" or "agent_complete"
	AgentType     string `json:"agent_type"`     // subagent_type from tool_input
	AgentID       string `json:"agent_id"`       // extracted from prompt
	Description   string `json:"description"`    // description from tool_input
	Instruction   string `json:"instruction,omitempty"`    // full prompt (agent_start only)
	OutputSummary string `json:"output_summary,omitempty"` // truncated output (agent_complete only)
	SessionID     string `json:"session_id"`     // extracted from prompt or found latest
	Timestamp     string `json:"timestamp"`      // UTC ISO format with Z suffix
}