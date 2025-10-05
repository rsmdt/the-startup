package stats

import (
	"encoding/json"
	"strings"
)

// AgentDetection represents a detected agent with confidence score
type AgentDetection struct {
	Agent      string  `json:"agent"`
	Confidence float64 `json:"confidence"`
	Method     string  `json:"method"` // "task_tool", "pattern", "context"
}

// TaskParameters represents parameters for the Task tool
type TaskParameters struct {
	Instructions string `json:"instructions"`
	SubagentType string `json:"subagent_type"`
}


// DetectAgent performs multi-method agent detection with confidence scoring
func DetectAgent(entry ClaudeLogEntry) (string, float64) {
	detection := detectAgentWithDetails(entry)
	return detection.Agent, detection.Confidence
}

// detectAgentWithDetails performs agent detection and returns full details
func detectAgentWithDetails(entry ClaudeLogEntry) AgentDetection {
	// Method 1: Direct Task tool detection (100% confidence)
	// This is the only reliable method - we only detect actual Task tool invocations
	if agent, ok := detectFromTaskTool(entry); ok {
		return AgentDetection{
			Agent:      agent,
			Confidence: 1.0,
			Method:     "task_tool",
		}
	}

	// Disabled pattern matching as it produces false positives
	// Only rely on actual Task tool invocations with subagent_type

	return AgentDetection{Agent: "", Confidence: 0.0}
}

// detectFromTaskTool checks for Task tool with subagent_type parameter
func detectFromTaskTool(entry ClaudeLogEntry) (string, bool) {
	// Check assistant messages for tool uses
	if entry.Type == "assistant" && entry.Assistant != nil {
		for _, toolUse := range entry.Assistant.ToolUses {
			if toolUse.Name == "Task" {
				agent := extractSubagentType(toolUse.Parameters)
				if agent != "" {
					return agent, true
				}
			}
		}
	}

	// Check user messages for tool results
	if entry.Type == "user" && entry.User != nil && entry.User.ToolName == "Task" {
		// Look for agent mentions in the result
		if entry.User.Result != nil {
			var result map[string]interface{}
			if err := json.Unmarshal(entry.User.Result, &result); err == nil {
				if subagent, ok := result["subagent_type"].(string); ok && subagent != "" {
					return subagent, true
				}
			}
		}
	}

	return "", false
}

// extractSubagentType extracts the subagent_type from Task tool parameters
func extractSubagentType(params json.RawMessage) string {
	if params == nil {
		return ""
	}

	var taskParams TaskParameters
	if err := json.Unmarshal(params, &taskParams); err != nil {
		return ""
	}

	return taskParams.SubagentType
}


// normalizeAgentName just cleans up the agent name without being clever
func normalizeAgentName(name string) string {
	// Just trim whitespace and return as-is
	return strings.TrimSpace(name)
}

// ValidateLogEntry checks if a log entry is valid for processing
func ValidateLogEntry(entry ClaudeLogEntry) error {
	// Basic validation - entry should have a type
	if entry.Type == "" {
		return &ToolError{
			Code:    "INVALID_ENTRY",
			Message: "Log entry missing type field",
		}
	}

	// Type-specific validation
	switch entry.Type {
	case "assistant":
		if entry.Assistant == nil {
			return &ToolError{
				Code:    "INVALID_ASSISTANT",
				Message: "Assistant entry missing content",
			}
		}
	case "user":
		if entry.User == nil {
			return &ToolError{
				Code:    "INVALID_USER",
				Message: "User entry missing content",
			}
		}
	case "system":
		if entry.System == nil {
			return &ToolError{
				Code:    "INVALID_SYSTEM",
				Message: "System entry missing content",
			}
		}
	}

	return nil
}