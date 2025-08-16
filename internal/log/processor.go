package log

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"time"
)

var (
	sessionIDRegex = regexp.MustCompile(`\bSessionId\s*:\s*([^\s,]+)`)
	agentIDRegex   = regexp.MustCompile(`\bAgentId\s*:\s*([^\s,]+)`)
)

// ProcessToolCall reads JSON from stdin and processes it according to the hook logic
// Returns HookData if processing should continue, nil if it should be skipped
func ProcessToolCall(input io.Reader, isPostHook bool) (*HookData, error) {
	// Read JSON input
	inputBytes, err := io.ReadAll(input)
	if err != nil {
		return nil, fmt.Errorf("failed to read input: %w", err)
	}

	var hookInput HookInput
	if err := json.Unmarshal(inputBytes, &hookInput); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Apply filtering logic
	if !ShouldProcess(hookInput.ToolName, hookInput.ToolInput) {
		return nil, nil // Skip processing
	}

	// Extract tool input fields
	subagentType, _ := hookInput.ToolInput["subagent_type"].(string)
	prompt, _ := hookInput.ToolInput["prompt"].(string)

	// Extract session and agent IDs
	sessionID := ExtractSessionID(prompt)

	// If no session ID found in prompt, try to find latest session
	if sessionID == "" {
		projectDir := GetProjectDir()
		sessionID = FindLatestSession(projectDir)
	}

	// Use enhanced AgentID extraction with fallback generation
	agentID := ExtractOrGenerateAgentID(prompt, subagentType, sessionID)

	// Create simplified log entry based on hook type
	hookData := &HookData{
		Timestamp: time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
		SessionID: sessionID, // Internal use only
		AgentID:   agentID,   // Internal use only
	}

	if isPostHook {
		hookData.Role = "assistant"
		hookData.Content = hookInput.Output // Full output, not truncated
	} else {
		hookData.Role = "user"
		hookData.Content = prompt // Full prompt sent to agent
	}

	return hookData, nil
}

// ShouldProcess determines if a hook should process this tool call
// Only processes Task tools with subagent_type starting with "the-"
func ShouldProcess(toolName string, toolInput map[string]interface{}) bool {
	// Check if this is a Task tool call
	if toolName != "Task" {
		return false
	}

	// Check subagent_type
	subagentType, ok := toolInput["subagent_type"].(string)
	if !ok {
		return false
	}

	// Only process agents starting with "the-"
	return strings.HasPrefix(subagentType, "the-")
}

// ExtractSessionID extracts session ID from prompt using regex
func ExtractSessionID(prompt string) string {
	matches := sessionIDRegex.FindStringSubmatch(prompt)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// ExtractAgentID extracts agent ID from prompt using regex
func ExtractAgentID(prompt string) string {
	matches := agentIDRegex.FindStringSubmatch(prompt)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// DebugLog outputs debug information to stderr if DEBUG_HOOKS is set
func DebugLog(format string, args ...interface{}) {
	if IsDebugEnabled() {
		fmt.Fprintf(os.Stderr, "[HOOK] "+format+"\n", args...)
	}
}

// DebugError outputs error information to stderr if DEBUG_HOOKS is set
func DebugError(err error) {
	if IsDebugEnabled() {
		fmt.Fprintf(os.Stderr, "[HOOK ERROR] %v\n", err)
	}
}
