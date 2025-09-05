# JSONL Metrics Collection Pattern

## Context

Claude Code provides hook mechanisms that allow shell commands to execute at specific lifecycle events. This pattern leverages PreToolUse and PostToolUse hooks to capture comprehensive metrics about all tool invocations, storing them in a simple JSONL format for analysis.

## Problem

Need to track tool usage statistics including:
- Which tools are invoked and how frequently
- Duration of tool executions
- Success/failure rates
- Tool usage patterns per session
- Error patterns and common failures

Without disrupting Claude Code's workflow or adding performance overhead.

## Solution

Implement a lightweight JSONL-based metrics collection system that:

1. **Captures all tool invocations** via Claude Code hooks
2. **Stores metrics in simple JSONL format** (one line per event)
3. **Tracks timing using PreToolUse and PostToolUse** events
4. **Provides aggregation commands** for analysis

### Architecture

```
Claude Code
    ├── PreToolUse Hook → the-startup metrics --pre
    │                          ↓
    │                    [Create entry with start_time]
    │                          ↓
    │                    metrics.jsonl (append)
    │
    ├── PostToolUse Hook → the-startup metrics --post
                               ↓
                         [Update entry with end_time]
                               ↓
                         metrics.jsonl (append)
```

### JSONL Schema

Each line represents a single tool invocation:

```json
{
  "tool_id": "unique_identifier",
  "tool_name": "Edit",
  "start_time": "2025-09-03T21:05:00.123Z",
  "end_time": "2025-09-03T21:05:00.456Z",
  "duration_ms": 333,
  "session_id": "dev-abc123",
  "status": "success",
  "tool_input": {...},
  "tool_output": {...}
}
```

### Hook Configuration

In `settings.json`:

```json
{
  "hooks": {
    "PreToolUse": [
      {
        "matcher": "*",
        "hooks": [
          {
            "type": "command",
            "command": "{{STARTUP_PATH}}/bin/the-startup metrics --pre"
          }
        ]
      }
    ],
    "PostToolUse": [
      {
        "matcher": "*",
        "hooks": [
          {
            "type": "command",
            "command": "{{STARTUP_PATH}}/bin/the-startup metrics --post"
          }
        ]
      }
    ]
  }
}
```

### Implementation Details

1. **Pre-hook processing**: Create new entry with status="pending"
2. **Post-hook processing**: Find matching entry, update with results
3. **Tool ID generation**: `{tool}_{timestamp}_{session_prefix}`
4. **Concurrent writes**: Use file locking (single writer)
5. **Error handling**: Always exit(0) to not disrupt Claude Code

## Examples

### Writing Metrics (Hook Handler)

```go
func HandlePreToolUse(input io.Reader) error {
    var hookData HookData
    if err := json.NewDecoder(input).Decode(&hookData); err != nil {
        return nil // Silent failure
    }
    
    entry := MetricsEntry{
        ToolID:     generateToolID(hookData),
        ToolName:   hookData.ToolName,
        StartTime:  time.Now().UTC(),
        SessionID:  hookData.SessionID,
        ToolInput:  hookData.ToolInput,
        Status:     "pending",
    }
    
    return appendToMetricsFile(entry)
}
```

### Reading Metrics (Analysis)

```go
func StreamMetrics(since time.Time) (chan MetricsEntry, error) {
    file, err := os.Open(metricsPath)
    if err != nil {
        return nil, err
    }
    
    entries := make(chan MetricsEntry)
    go func() {
        defer close(entries)
        scanner := bufio.NewScanner(file)
        for scanner.Scan() {
            var entry MetricsEntry
            if json.Unmarshal(scanner.Bytes(), &entry) == nil {
                if entry.StartTime.After(since) {
                    entries <- entry
                }
            }
        }
    }()
    
    return entries, nil
}
```

## When to Use

Use this pattern when:
- You need to track tool usage metrics
- Performance monitoring is required
- Error pattern analysis is needed
- Usage statistics are valuable
- Audit logging is required

## Benefits

1. **Simple**: JSONL format is easy to write and parse
2. **Performance**: Minimal overhead, async writes
3. **Flexible**: Can add fields without breaking parsers
4. **Streamable**: Can process large files efficiently
5. **Debuggable**: Human-readable format

## Trade-offs

1. **Storage**: No automatic rotation/cleanup (manual management)
2. **Querying**: No indexed queries (sequential scan)
3. **Correlation**: Manual matching of pre/post events
4. **Persistence**: File-based, no database features

## Related Patterns

- Session-based Context: Uses session IDs for correlation
- Agent ID Extraction: Similar hook-based data extraction