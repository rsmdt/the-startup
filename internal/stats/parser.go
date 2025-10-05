package stats

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

// JSONLParser implements the Parser interface for streaming JSONL log parsing
type JSONLParser struct {
	options ParseOptions
	filter  FilterOptions
	mu      sync.RWMutex
}

// NewJSONLParser creates a new JSONL parser with default options
func NewJSONLParser() *JSONLParser {
	return &JSONLParser{
		options: ParseOptions{
			BufferSize:     1024 * 1024, // 1MB buffer for large Claude log lines
			ParseTools:     true,
			SkipSystemLogs: false,
			StrictMode:     false,
			CollectErrors:  false,
			MaxMemoryMB:    100,
		},
	}
}

// SetOptions updates parser configuration
func (p *JSONLParser) SetOptions(opts ParseOptions) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.options = opts
}

// SetFilter updates filtering criteria
func (p *JSONLParser) SetFilter(filter FilterOptions) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.filter = filter
}

// ParseStream processes a log file and returns channels of parsed entries and errors
func (p *JSONLParser) ParseStream(reader io.Reader) (<-chan ClaudeLogEntry, <-chan error) {
	p.mu.RLock()
	opts := p.options
	filter := p.filter
	p.mu.RUnlock()

	entries := make(chan ClaudeLogEntry, 100)
	errors := make(chan error, 10)

	go func() {
		defer close(entries)
		defer close(errors)

		scanner := bufio.NewScanner(reader)
		// Set buffer size for large files
		buf := make([]byte, 0, opts.BufferSize)
		scanner.Buffer(buf, opts.BufferSize*2) // Allow lines up to 2x buffer size

		lineNum := 0
		processedCount := 0
		eligibleCount := 0

		for scanner.Scan() {
			lineNum++
			line := scanner.Bytes()

			// Skip empty lines
			if len(line) == 0 {
				continue
			}

			// Parse the entry
			entry, err := p.parseEntry(line, lineNum, opts)
			if err != nil {
				if opts.StrictMode {
					errors <- fmt.Errorf("line %d: %w", lineNum, err)
					return
				}
				if opts.CollectErrors {
					errors <- fmt.Errorf("line %d: %w", lineNum, err)
				}
				continue
			}

			// Apply filters
			if p.shouldSkip(entry, filter) {
				continue
			}

			// Apply sampling if configured
			if filter.SampleRate > 0 && filter.SampleRate < 1.0 {
				// Use modulo to sample approximately the right percentage
				sampleInterval := int(1.0 / filter.SampleRate)
				if sampleInterval > 0 && eligibleCount%sampleInterval != 0 {
					eligibleCount++
					continue
				}
			}
			eligibleCount++

			// Check max entries limit
			if filter.MaxEntries > 0 && processedCount >= filter.MaxEntries {
				break
			}

			processedCount++
			entries <- *entry
		}

		if err := scanner.Err(); err != nil {
			errors <- fmt.Errorf("scanner error: %w", err)
		}
	}()

	return entries, errors
}

// ParseFile opens and parses a specific log file
func (p *JSONLParser) ParseFile(filePath string) (<-chan ClaudeLogEntry, <-chan error) {
	entries := make(chan ClaudeLogEntry, 100)
	errors := make(chan error, 10)

	go func() {
		defer close(entries)
		defer close(errors)

		file, err := os.Open(filePath)
		if err != nil {
			errors <- fmt.Errorf("failed to open file: %w", err)
			return
		}
		defer file.Close()

		// Use ParseStream for the actual parsing
		streamEntries, streamErrors := p.ParseStream(file)

		// Forward errors in the same goroutine to avoid race conditions
		// Process both channels in parallel using select
		for streamEntries != nil || streamErrors != nil {
			select {
			case entry, ok := <-streamEntries:
				if !ok {
					streamEntries = nil
				} else {
					entries <- entry
				}
			case err, ok := <-streamErrors:
				if !ok {
					streamErrors = nil
				} else if err != nil {
					errors <- err
				}
			}
		}
	}()

	return entries, errors
}

// parseEntry parses a single JSONL line into a ClaudeLogEntry
func (p *JSONLParser) parseEntry(line []byte, lineNum int, opts ParseOptions) (*ClaudeLogEntry, error) {
	// Parse the raw JSON to understand the structure
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(line, &raw); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Extract the type field
	var entryType string
	if typeRaw, ok := raw["type"]; ok {
		if err := json.Unmarshal(typeRaw, &entryType); err != nil {
			return nil, fmt.Errorf("failed to parse type: %w", err)
		}
	} else {
		return nil, fmt.Errorf("missing type field")
	}

	// Create the entry with the type
	entry := &ClaudeLogEntry{
		Type: entryType,
	}

	// Extract timestamp if present
	if timestampRaw, ok := raw["timestamp"]; ok {
		var timestamp string
		if err := json.Unmarshal(timestampRaw, &timestamp); err == nil {
			if t, err := time.Parse(time.RFC3339, timestamp); err == nil {
				entry.Timestamp = t
			}
		}
	}

	// Extract session ID if present
	if sessionRaw, ok := raw["sessionId"]; ok {
		json.Unmarshal(sessionRaw, &entry.SessionID)
	}

	// Parse the content based on the type
	switch entry.Type {
	case "user":
		// For user messages, the content is in the "message" field
		if msgRaw, ok := raw["message"]; ok {
			entry.Content = msgRaw
			
			// Parse the message field
			var msg struct {
				Role    string          `json:"role"`
				Content json.RawMessage `json:"content"`
			}
			if err := json.Unmarshal(msgRaw, &msg); err != nil {
				return nil, fmt.Errorf("failed to parse user message structure: %w", err)
			}
			
			// Create UserMessage
			userMsg := UserMessage{
				Role: msg.Role,
				Timestamp: entry.Timestamp,
			}
			
			// The content can be either a string (simple message) or an array of content blocks
			var contentStr string
			if err := json.Unmarshal(msg.Content, &contentStr); err == nil {
				// Simple string content
				userMsg.Text = contentStr
			} else {
				// Try as array of content blocks
				var contentArray []map[string]interface{}
				if err := json.Unmarshal(msg.Content, &contentArray); err == nil {
					// Process content blocks
					for _, block := range contentArray {
						if blockType, ok := block["type"].(string); ok {
							if blockType == "text" {
								if text, ok := block["text"].(string); ok {
									userMsg.Text += text
								}
							} else if blockType == "tool_result" {
								userMsg.IsToolResult = true
								if toolUseID, ok := block["tool_use_id"].(string); ok {
									userMsg.ToolUseID = toolUseID
								}
								if output, ok := block["output"].(string); ok {
									userMsg.Output = output
								}
								// Extract result field for Task tool results
								if result, ok := block["result"]; ok {
									if resultBytes, err := json.Marshal(result); err == nil {
										userMsg.Result = resultBytes
									}
								}
								// Extract tool_name if present
								if toolName, ok := block["tool_name"].(string); ok {
									userMsg.ToolName = toolName
								}
								if isError, ok := block["is_error"].(bool); ok && isError {
									errorMsg := "Tool execution failed"
									userMsg.Error = &errorMsg
								}
							}
						}
					}
				}
			}

			entry.User = &userMsg

			// Detect agent usage for user messages (e.g., Tool results with subagent_type)
			if agent, confidence := DetectAgent(*entry); agent != "" {
				entry.User.DetectedAgent = agent
				entry.User.DetectionConfidence = confidence
			}
		}

	case "assistant":
		// For assistant messages, the content is in the "message" field
		if msgRaw, ok := raw["message"]; ok {
			entry.Content = msgRaw
			
			// Parse the message field
			var msg struct {
				ID      string          `json:"id"`
				Type    string          `json:"type"`
				Role    string          `json:"role"`
				Model   string          `json:"model"`
				Content json.RawMessage `json:"content"`
				Usage   struct {
					InputTokens  int `json:"input_tokens"`
					OutputTokens int `json:"output_tokens"`
				} `json:"usage"`
			}
			if err := json.Unmarshal(msgRaw, &msg); err != nil {
				return nil, fmt.Errorf("failed to parse assistant message structure: %w", err)
			}
			
			// Create AssistantMessage
			assistantMsg := AssistantMessage{
				ID:           msg.ID,
				Role:         msg.Role,
				Model:        msg.Model,
				Timestamp:    entry.Timestamp,
				InputTokens:  msg.Usage.InputTokens,
				OutputTokens: msg.Usage.OutputTokens,
			}
			
			// Parse content array
			var contentArray []map[string]interface{}
			if err := json.Unmarshal(msg.Content, &contentArray); err == nil {
				for _, block := range contentArray {
					if blockType, ok := block["type"].(string); ok {
						if blockType == "text" {
							if text, ok := block["text"].(string); ok {
								assistantMsg.Text += text
							}
						} else if blockType == "tool_use" {
							// Parse tool use
							toolUse := ToolUse{}
							if id, ok := block["id"].(string); ok {
								toolUse.ID = id
							}
							if name, ok := block["name"].(string); ok {
								toolUse.Name = name
							}
							if input, ok := block["input"]; ok {
								// Convert input to JSON for storage
								if inputBytes, err := json.Marshal(input); err == nil {
									toolUse.Parameters = inputBytes
								}
							}
							assistantMsg.ToolUses = append(assistantMsg.ToolUses, toolUse)
						}
					}
				}
			}
			
			// Parse tool parameters if enabled
			if opts.ParseTools && len(assistantMsg.ToolUses) > 0 {
				for i := range assistantMsg.ToolUses {
					p.parseToolParameters(&assistantMsg.ToolUses[i])
				}
			}
			
			// Extract commands from text
			if assistantMsg.Text != "" {
				p.extractCommands(&assistantMsg)
			}

			entry.Assistant = &assistantMsg

			// Detect agent usage after entry is fully populated
			if agent, confidence := DetectAgent(*entry); agent != "" {
				entry.Assistant.DetectedAgent = agent
				entry.Assistant.DetectionConfidence = confidence
			}
		}

	case "system":
		if opts.SkipSystemLogs {
			return entry, nil
		}
		
		// For system messages, check if there's a message field
		if msgRaw, ok := raw["message"]; ok {
			entry.Content = msgRaw
			
			// Parse the message field
			var msg struct {
				Role    string          `json:"role"`
				Content json.RawMessage `json:"content"`
			}
			if err := json.Unmarshal(msgRaw, &msg); err == nil {
				systemMsg := SystemMessage{
					Role:      msg.Role,
					Timestamp: entry.Timestamp,
				}
				
				// Try to get text content
				var contentStr string
				if err := json.Unmarshal(msg.Content, &contentStr); err == nil {
					systemMsg.Text = contentStr
				}

				entry.System = &systemMsg

				// Detect agent usage for system messages
				if agent, confidence := DetectAgent(*entry); agent != "" {
					entry.System.DetectedAgent = agent
					entry.System.DetectionConfidence = confidence
				}
			}
		}

	case "summary":
		// For summary messages, extract summary text if available
		summaryMsg := SummaryMessage{
			Timestamp: entry.Timestamp,
		}
		
		if summaryRaw, ok := raw["summary"]; ok {
			var summaryText string
			if err := json.Unmarshal(summaryRaw, &summaryText); err == nil {
				summaryMsg.Summary = summaryText
			}
		}
		
		entry.Summary = &summaryMsg

	default:
		// Unknown type, but don't fail - just leave content as raw
		if opts.StrictMode {
			return nil, fmt.Errorf("unknown message type: %s", entry.Type)
		}
	}

	return entry, nil
}

// parseToolParameters parses tool-specific parameters
func (p *JSONLParser) parseToolParameters(tool *ToolUse) {
	if len(tool.Parameters) == 0 {
		return
	}

	switch tool.Name {
	case "Bash", "bash":
		var params BashParameters
		if err := json.Unmarshal(tool.Parameters, &params); err == nil {
			tool.BashParams = &params
		}

	case "Read", "read":
		var params ReadParameters
		if err := json.Unmarshal(tool.Parameters, &params); err == nil {
			tool.ReadParams = &params
		}

	case "Edit", "edit":
		var params EditParameters
		if err := json.Unmarshal(tool.Parameters, &params); err == nil {
			tool.EditParams = &params
		}

	case "Write", "write":
		var params WriteParameters
		if err := json.Unmarshal(tool.Parameters, &params); err == nil {
			tool.WriteParams = &params
		}

	case "WebSearch", "Grep", "Glob", "search":
		var params SearchParameters
		if err := json.Unmarshal(tool.Parameters, &params); err == nil {
			tool.SearchParams = &params
		}
	}
}

// extractCommands extracts command tags from assistant messages
func (p *JSONLParser) extractCommands(msg *AssistantMessage) {
	// Look for command tags in the format <command-name>...</command-name>
	// Common patterns include <bash>...</bash>, <python>...</python>, etc.
	// Note: Go regex doesn't support backreferences, so we do a simpler pattern
	commandPattern := regexp.MustCompile(`<([a-zA-Z_-]+)>([^<]*)</([a-zA-Z_-]+)>`)
	matches := commandPattern.FindAllStringSubmatch(msg.Text, -1)
	
	for _, match := range matches {
		if len(match) >= 4 && match[1] == match[3] {
			// match[1] is the opening tag
			// match[2] is the command content
			// match[3] is the closing tag
			// Store this information somewhere if needed for analysis
			// For now, we just parse it without storing
		}
	}
}

// shouldSkip determines if an entry should be skipped based on filters
func (p *JSONLParser) shouldSkip(entry *ClaudeLogEntry, filter FilterOptions) bool {
	// Time filters
	if filter.StartTime != nil && entry.Timestamp.Before(*filter.StartTime) {
		return true
	}
	if filter.EndTime != nil && entry.Timestamp.After(*filter.EndTime) {
		return true
	}

	// Session filters
	if len(filter.SessionIDs) > 0 {
		found := false
		for _, sessionID := range filter.SessionIDs {
			if entry.SessionID == sessionID {
				found = true
				break
			}
		}
		if !found {
			return true
		}
	}

	// Tool filters (only apply to assistant messages with tools)
	if entry.Assistant != nil && len(entry.Assistant.ToolUses) > 0 {
		for _, tool := range entry.Assistant.ToolUses {
			// Check include list
			if len(filter.IncludeTools) > 0 {
				found := false
				for _, includeTool := range filter.IncludeTools {
					if strings.EqualFold(tool.Name, includeTool) {
						found = true
						break
					}
				}
				if !found {
					return true
				}
			}

			// Check exclude list
			for _, excludeTool := range filter.ExcludeTools {
				if strings.EqualFold(tool.Name, excludeTool) {
					return true
				}
			}
		}
	}

	// Status filters for tool results
	if entry.User != nil && entry.User.IsToolResult {
		if filter.SuccessOnly && entry.User.Error != nil {
			return true
		}
		if filter.FailuresOnly && entry.User.Error == nil {
			return true
		}
	}

	return false
}