# Claude Code Hooks Interface

## Overview

Claude Code hooks are shell commands that execute at specific points in Claude Code's lifecycle, enabling external programs to observe and react to tool invocations.

## Hook Event Types

### PreToolUse
Executed before a tool is invoked. Can block tool execution by returning non-zero exit code.

### PostToolUse
Executed after a tool completes. Receives tool response in addition to input.

### Other Events
- `UserPromptSubmit`: When user submits a prompt
- `Notification`: When Claude Code sends notifications
- `Stop`: When Claude Code finishes responding
- `SubagentStop`: When subagent tasks complete
- `PreCompact`: Before compact operation
- `SessionStart`: Starting/resuming a session
- `SessionEnd`: When session ends

## Data Format

Hook data is passed as JSON to the command's stdin:

### PreToolUse Payload
```json
{
  "session_id": "dev-20250903-abc123",
  "transcript_path": "/Users/user/.claude/projects/abc/transcript.jsonl",
  "cwd": "/Users/user/project",
  "hook_event_name": "PreToolUse",
  "tool_name": "Edit",
  "tool_input": {
    "file_path": "/path/to/file.txt",
    "old_string": "foo",
    "new_string": "bar"
  }
}
```

### PostToolUse Payload
```json
{
  "session_id": "dev-20250903-abc123",
  "transcript_path": "/Users/user/.claude/projects/abc/transcript.jsonl",
  "cwd": "/Users/user/project",
  "hook_event_name": "PostToolUse",
  "tool_name": "Edit",
  "tool_input": {
    "file_path": "/path/to/file.txt",
    "old_string": "foo",
    "new_string": "bar"
  },
  "tool_response": {
    "success": true,
    "message": "Successfully replaced 1 occurrence"
  }
}
```

## Configuration

Hooks are configured in Claude Code's `settings.json`:

```json
{
  "hooks": {
    "PreToolUse": [
      {
        "matcher": "*",
        "hooks": [
          {
            "type": "command",
            "command": "/path/to/hook-handler"
          }
        ]
      }
    ],
    "PostToolUse": [
      {
        "matcher": "Edit|Write|MultiEdit",
        "hooks": [
          {
            "type": "command",
            "command": "/path/to/file-change-handler"
          }
        ]
      }
    ]
  }
}
```

### Matcher Patterns
- `"*"`: Matches all tools
- `"Edit"`: Matches specific tool
- `"Edit|Write"`: Matches multiple tools (pipe-separated)

## Environment Variables

Hooks receive these environment variables:
- `CLAUDE_SESSION_ID`: Current session identifier
- `CLAUDE_PROJECT_ROOT`: Project root directory (if available)
- Standard shell environment from user's session

## Implementation Requirements

### Exit Codes
- **0**: Success, allow operation to continue
- **Non-zero** (PreToolUse only): Block the tool execution
- **Non-zero** (PostToolUse): Logged but doesn't affect Claude Code

### Error Handling
```go
func main() {
    defer func() {
        if r := recover(); r != nil {
            // Always exit 0 to not disrupt Claude Code
            os.Exit(0)
        }
    }()
    
    if err := processHook(os.Stdin); err != nil {
        // Log error if DEBUG enabled, but exit 0
        if os.Getenv("DEBUG_HOOKS") != "" {
            fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        }
        os.Exit(0)
    }
}
```

### Performance Considerations
- Hooks should complete quickly (< 1 second ideal)
- Long-running operations should be async
- Avoid blocking I/O operations
- Use buffering for file writes

## Available Tools

Common tools that trigger hooks:
- **File Operations**: Read, Write, Edit, MultiEdit
- **Search**: Grep, Glob, WebSearch
- **Execution**: Bash, Task
- **Web**: WebFetch
- **Development**: NotebookEdit, TodoWrite
- **MCP Tools**: Various mcp__* prefixed tools

## Examples

### Simple Logger
```bash
#!/bin/bash
# Log all tool invocations
cat >> ~/.claude-tools.log
```

### Python Hook Handler
```python
#!/usr/bin/env python3
import json
import sys
from datetime import datetime

try:
    data = json.load(sys.stdin)
    
    # Log to JSONL
    with open('/tmp/claude-metrics.jsonl', 'a') as f:
        data['timestamp'] = datetime.utcnow().isoformat()
        json.dump(data, f)
        f.write('\n')
        
except Exception:
    pass  # Silent failure

sys.exit(0)  # Always exit 0
```

### Go Hook Handler
```go
package main

import (
    "encoding/json"
    "os"
    "time"
)

type HookData struct {
    SessionID   string          `json:"session_id"`
    ToolName    string          `json:"tool_name"`
    ToolInput   json.RawMessage `json:"tool_input"`
    ToolOutput  json.RawMessage `json:"tool_response,omitempty"`
    HookEvent   string          `json:"hook_event_name"`
}

func main() {
    var data HookData
    if err := json.NewDecoder(os.Stdin).Decode(&data); err != nil {
        os.Exit(0)
    }
    
    // Process the hook data
    processMetrics(data)
    
    os.Exit(0)
}
```

## Security Considerations

⚠️ **Warning**: Hooks run with your environment's credentials and can access your filesystem.

- Always review hook scripts before installation
- Use absolute paths in commands
- Validate and sanitize input data
- Avoid executing user-provided strings
- Consider using restricted shells or containers

## Debugging

Enable debug output with environment variable:
```bash
export DEBUG_HOOKS=1
```

Test hooks manually:
```bash
echo '{"tool_name":"Test","session_id":"test"}' | /path/to/hook-handler
```

## Best Practices

1. **Fail silently**: Never let hook errors disrupt Claude Code
2. **Be fast**: Complete quickly to avoid slowing down tools
3. **Use absolute paths**: Avoid PATH dependencies
4. **Buffer writes**: Don't write to disk on every invocation
5. **Validate input**: Check JSON structure before processing
6. **Log errors separately**: Use stderr or separate log files
7. **Handle missing fields**: Not all fields present in all events