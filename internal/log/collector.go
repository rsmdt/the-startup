package log

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"
)

// ProcessHook reads hook data from the provided reader (typically stdin),
// auto-detects the event type via hook_event_name field, generates tool IDs
// for correlation, and extracts success status from tool outputs.
// Returns nil on any error (silent failure for hook processing).
func ProcessHook(input io.Reader) error {
	// Parse hook data from input
	var hookData HookPayload
	if err := json.NewDecoder(input).Decode(&hookData); err != nil {
		return nil // Silent failure
	}

	// Auto-detect and handle event type based on hook_event_name field
	switch hookData.HookEventName {
	case "PreToolUse":
		return processPreToolUse(hookData)
	case "PostToolUse":
		return processPostToolUse(hookData)
	default:
		// Ignore other event types
		return nil
	}
}

// processPreToolUse handles PreToolUse events by creating a metrics entry
// with a generated tool ID for later correlation with PostToolUse events.
func processPreToolUse(hookData HookPayload) error {
	// Generate unique tool ID for correlation
	toolID := generateToolID(hookData)

	// Parse timestamp from hook data or use current time
	timestamp := parseTimestamp(hookData.Timestamp)

	entry := MetricsEntry{
		ToolID:         toolID,
		ToolName:       hookData.ToolName,
		HookEvent:      "PreToolUse",
		Timestamp:      timestamp,
		SessionID:      hookData.SessionID,
		TranscriptPath: hookData.TranscriptPath,
		CWD:            hookData.CWD,
		ToolInput:      hookData.ToolInput,
	}

	// Write the metrics entry to disk
	return AppendMetrics(&entry)
}

// processPostToolUse handles PostToolUse events by creating a metrics entry
// with the same tool ID as the corresponding PreToolUse event, and extracting
// the success status from the tool response.
func processPostToolUse(hookData HookPayload) error {
	// Generate the same tool ID for correlation with PreToolUse
	toolID := generateToolID(hookData)

	// Parse timestamp from hook data or use current time
	timestamp := parseTimestamp(hookData.Timestamp)

	// Extract success status from tool response
	success := extractSuccess(hookData)

	entry := MetricsEntry{
		ToolID:         toolID,
		ToolName:       hookData.ToolName,
		HookEvent:      "PostToolUse",
		Timestamp:      timestamp,
		SessionID:      hookData.SessionID,
		TranscriptPath: hookData.TranscriptPath,
		CWD:            hookData.CWD,
		ToolInput:      hookData.ToolInput,
		ToolOutput:     coalesceOutput(hookData),
		Success:        success,
		Error:          hookData.Error,
		ErrorType:      hookData.ErrorType,
	}

	// Write the metrics entry to disk
	return AppendMetrics(&entry)
}

// generateToolID creates a unique identifier for correlating Pre and Post tool events.
// It uses the RequestID if available (preferred for exact correlation), 
// or generates a deterministic ID based on the tool invocation parameters.
func generateToolID(hookData HookPayload) string {
	// If RequestID is provided, use it directly for perfect correlation
	if hookData.RequestID != "" {
		return hookData.RequestID
	}
	
	// Otherwise, generate a deterministic ID based on the invocation parameters
	// This uses a hash of the tool input to ensure Pre and Post events with the same
	// input get the same ID
	inputHash := ""
	if len(hookData.ToolInput) > 0 {
		// Use a simple hash of first 16 chars of input for uniqueness
		inputStr := string(hookData.ToolInput)
		if len(inputStr) > 16 {
			inputStr = inputStr[:16]
		}
		// Simple hash: just use the input directly (since it's deterministic)
		inputHash = fmt.Sprintf("%x", inputStr)
	}
	
	// Get session prefix (first 8 chars of session ID)
	sessionPrefix := hookData.SessionID
	if len(sessionPrefix) > 8 {
		sessionPrefix = sessionPrefix[:8]
	}

	// For PreToolUse, use the timestamp as-is
	// For PostToolUse, try to find a matching Pre by truncating to second precision
	timestamp := parseTimestamp(hookData.Timestamp)
	// Truncate to second precision for better correlation
	timeStr := timestamp.Truncate(time.Second).Format("20060102T150405Z")

	// Generate ID in format: toolname_session_time_inputhash
	return fmt.Sprintf("%s_%s_%s_%s",
		strings.ToLower(hookData.ToolName),
		sessionPrefix,
		timeStr,
		inputHash[:min(8, len(inputHash))])
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// extractSuccess determines if a tool execution was successful based on
// the presence of errors and the content of the tool response.
func extractSuccess(hookData HookPayload) *bool {
	// If there's an explicit error, the tool failed
	if hookData.Error != "" || hookData.ErrorType != "" {
		success := false
		return &success
	}

	// Get the tool output (could be in ToolResponse or ToolOutput field)
	output := coalesceOutput(hookData)
	
	// If we have output, check for error indicators
	// Note: JSON null is represented as the string "null" (4 bytes)
	if len(output) > 0 && string(output) != "null" {
		// Try to parse the output as JSON to look for error fields
		var outputData map[string]interface{}
		if err := json.Unmarshal(output, &outputData); err == nil {
			// Check for common error field patterns
			if errorField, ok := outputData["error"]; ok {
				if errorField != nil && errorField != false && errorField != "" {
					success := false
					return &success
				}
			}
			
			// Check for explicit success field
			if successField, ok := outputData["success"]; ok {
				if boolVal, ok := successField.(bool); ok {
					return &boolVal
				}
			}

			// Check for status field
			if statusField, ok := outputData["status"]; ok {
				if strVal, ok := statusField.(string); ok {
					successVal := strings.ToLower(strVal) == "success" || 
					              strings.ToLower(strVal) == "ok" ||
					              strings.ToLower(strVal) == "completed"
					return &successVal
				}
			}
		}

		// If output exists and no errors found, assume success
		success := true
		return &success
	}

	// For PostToolUse without output or error, we can't determine success
	return nil
}

// coalesceOutput returns the tool output from either ToolResponse or ToolOutput field
func coalesceOutput(hookData HookPayload) json.RawMessage {
	if len(hookData.ToolResponse) > 0 {
		return hookData.ToolResponse
	}
	return hookData.ToolOutput
}

// parseTimestamp attempts to parse the timestamp string from hook data.
// Falls back to current time if parsing fails.
func parseTimestamp(timestampStr string) time.Time {
	if timestampStr == "" {
		return time.Now().UTC()
	}

	// Try common timestamp formats
	formats := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02T15:04:05.999Z07:00",
		"2006-01-02 15:04:05",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, timestampStr); err == nil {
			return t.UTC()
		}
	}

	// If all parsing attempts fail, use current time
	return time.Now().UTC()
}