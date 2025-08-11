# System Design Document: Go Hooks Architecture

## Executive Summary

This document outlines the design for replacing the existing Python-based hook system with a native Go implementation integrated into the main `the-startup` CLI. The Go implementation uses semantic `log` commands with direction flags (`--assistant`/`--user`) to maintain full compatibility with the current logging format while providing better performance, easier deployment, and reduced dependencies.

## Claude Code Hooks Documentation

### Hook Configuration
Claude Code hooks are configured in the `settings.json` file located in the project's `.claude` directory. Hooks can be triggered before (PreToolUse) or after (PostToolUse) tool execution.

### Hook Input Schema
When a hook is triggered, Claude Code passes JSON data via stdin:
```json
{
  "session_id": "abc123",
  "transcript_path": "/path/to/transcript.jsonl",
  "cwd": "/current/working/directory",
  "hook_event_name": "PreToolUse" | "PostToolUse",
  "tool_name": "Task",
  "tool_input": {
    "description": "Short task description",
    "prompt": "Full prompt for the agent",
    "subagent_type": "the-architect" | "the-developer" | etc.
  },
  "output": "Agent response (PostToolUse only)"
}
```

### Environment Variables
- `CLAUDE_PROJECT_DIR`: Absolute path to the project root directory
- `DEBUG_HOOKS`: When set, enables debug output to stderr

### Hook Response
Hooks can optionally return JSON to control tool execution:
```json
{
  "action": "allow" | "deny",
  "reason": "Optional explanation for the action"
}
```

## System Context

### Current State
- Python hooks (`log_agent_start.py`, `log_agent_complete.py`) capture Claude Code agent interactions
- Hooks process Task tool calls for agents with subagent_type starting with "the-"
- Logs stored in JSONL format in `.the-startup/` directory structure
- Integration via Claude Code settings.json with PreToolUse/PostToolUse hooks

### Drivers for Change
- Remove Python/uv runtime dependency
- Improve startup performance (Python import overhead)
- Simplify deployment (single binary vs script + dependencies)
- Better error handling and logging consistency
- Easier cross-platform support

## System Architecture

### High-Level Design

```
┌─────────────────────────────────────────────────────────────┐
│                    Claude Code                               │
└─────────────┬───────────────────────────────────────────────┘
              │
              ▼
┌─────────────────────────────────────────────────────────────┐
│                Claude Settings.json                         │
│  PreToolUse: the-startup log --assistant                   │
│  PostToolUse: the-startup log --user                       │
└─────────────┬───────────────────────────────────────────────┘
              │
              ▼ (JSON via stdin)
┌─────────────────────────────────────────────────────────────┐
│             the-startup CLI (log command)                   │
│                                                             │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────────────┐   │
│  │ CLI Router  │ │ Processor   │ │    Logger           │   │
│  │             │ │             │ │                     │   │
│  │ --assistant │ │ - JSON      │ │ - Session Files     │   │
│  │ --user      │ │ - Filtering │ │ - Global Files      │   │
│  │             │ │ - Parsing   │ │ - Directory Mgmt    │   │
│  └─────────────┘ └─────────────┘ └─────────────────────┘   │
└─────────────┬───────────────────────────────────────────────┘
              │
              ▼
┌─────────────────────────────────────────────────────────────┐
│              File System                                    │
│  .the-startup/                                             │
│  ├── dev-session-123/                                      │
│  │   └── agent-instructions.jsonl                          │
│  └── all-agent-instructions.jsonl                          │
└─────────────────────────────────────────────────────────────┘
```

### Component Design

#### 1. CLI Router (`cmd/log.go`)
**Responsibility**: Command-line interface and routing for log command

**Components**:
- `log.go` - Log command handler with assistant/user flags
- Integration with existing `the-startup` CLI structure

**Interface**:
```bash
the-startup log --assistant  # Process PreToolUse hook (agent instructions)
the-startup log -a           # Short form
the-startup log --user       # Process PostToolUse hook (agent responses)
the-startup log -u           # Short form
```

#### 2. Processor (`internal/hooks/processor.go`)
**Responsibility**: JSON input processing and filtering

**Functions**:
- `ProcessToolCall(input []byte) (*HookData, error)` - Parse and validate JSON
- `ShouldProcess(toolName, subagentType string) bool` - Apply filtering logic
- `ExtractSessionID(prompt string) string` - Extract session from prompt
- `ExtractAgentID(prompt string) string` - Extract agent ID from prompt

**Data Structures**:
```go
type HookInput struct {
    ToolName  string                 `json:"tool_name"`
    ToolInput map[string]interface{} `json:"tool_input"`
    Output    string                 `json:"output,omitempty"`
}

type HookData struct {
    Event       string `json:"event"`
    AgentType   string `json:"agent_type"`
    AgentID     string `json:"agent_id"`
    Description string `json:"description"`
    Instruction string `json:"instruction,omitempty"`
    OutputSummary string `json:"output_summary,omitempty"`
    SessionID   string `json:"session_id"`
    Timestamp   string `json:"timestamp"`
}
```

#### 3. Logger (`internal/hooks/logger.go`)
**Responsibility**: File I/O and directory management

**Functions**:
- `WriteSessionLog(sessionID string, data *HookData) error` - Write to session file
- `WriteGlobalLog(data *HookData) error` - Write to global file
- `EnsureDirectories(projectDir string) error` - Create required directories
- `FindLatestSession(projectDir string) string` - Find most recent session

#### 4. Config (`internal/hooks/config.go`)
**Responsibility**: Configuration and environment handling

**Functions**:
- `GetStartupDir(projectDir string) string` - Determine log directory
- `GetProjectDir() string` - Get project directory from environment
- `IsDebugEnabled() bool` - Check DEBUG_HOOKS environment variable

### Data Flow

#### PreToolUse Hook (log --assistant command)
1. Read JSON from stdin
2. Validate tool_name == "Task"
3. Check subagent_type starts with "the-"
4. Extract session_id, agent_id from prompt
5. Create log entry with event="agent_start"
6. Write to session-specific and global JSONL files
7. Optional debug output to stderr

#### PostToolUse Hook (log --user command)
1. Read JSON from stdin
2. Validate tool_name == "Task" 
3. Check subagent_type starts with "the-"
4. Extract session_id, agent_id from prompt
5. Truncate output if > 1000 characters
6. Create log entry with event="agent_complete"
7. Write to session-specific and global JSONL files
8. Optional debug output to stderr

### File Structure

```
├── cmd/
│   ├── log.go                # Log command implementation
│   └── commands.go           # Updated with log command registration
├── internal/
│   └── hooks/
│       ├── processor.go      # JSON processing
│       ├── logger.go         # File I/O
│       ├── config.go         # Configuration
│       └── types.go          # Shared data types
└── the-startup               # Main compiled binary (includes log command)
```

### Integration Points

#### 1. Build System
- Integrate log command into main `the-startup` build process
- Cross-compile single binary for Linux, macOS, Windows
- No separate hook binary needed

#### 2. Installer Updates
- Replace Python hook installation with integrated Go command
- Update settings.json generation:
  ```json
  {
    "hooks": {
      "PreToolUse": [{
        "matcher": "Task",
        "hooks": [{
          "type": "command",
          "command": "$CLAUDE_PROJECT_DIR/.claude/hooks/the-startup log --assistant",
          "_source": "the-startup"
        }]
      }],
      "PostToolUse": [{
        "matcher": "Task", 
        "hooks": [{
          "type": "command",
          "command": "$CLAUDE_PROJECT_DIR/.claude/hooks/the-startup log --user",
          "_source": "the-startup"
        }]
      }]
    }
  }
  ```

#### 3. Asset Management
- Remove Python scripts from embed.FS
- Update installer to deploy main `the-startup` binary with log command
- No separate hook binary in embed.FS needed

## Technology Decisions

### Language Choice: Go
**Rationale**: Native compilation, excellent JSON support, cross-platform, consistent with main application

### CLI Framework: Cobra
**Rationale**: Already used in main application, mature, excellent subcommand support

### JSON Library: Standard Library
**Rationale**: Built-in encoding/json is sufficient, no external dependencies needed

### Build Strategy: Integrated Command
**Rationale**: Single `the-startup` binary with log command, simpler deployment, shared code, easier maintenance, semantic CLI structure

## Performance Considerations

### Startup Performance
- Integrated command startup: ~1-5ms vs Python: ~50-100ms
- No import overhead or dependency resolution

### Memory Usage
- Integrated command: ~2-5MB additional vs Python: ~15-30MB
- No separate process overhead

### File I/O
- Efficient append-only writes
- Minimal memory allocation for JSON processing

## Security Considerations

### Input Validation
- Strict JSON schema validation
- Sanitize file paths
- Validate environment variables

### File Permissions
- Main binary executable: 0755 (already required)
- Log files: 0644 
- Log directories: 0755

### Error Handling
- No sensitive data in error messages
- Fail gracefully on permission errors
- Debug output only when explicitly enabled

## Error Handling Strategy

### Input Errors
- Invalid JSON: Silent exit (match Python behavior)
- Missing fields: Log warning, continue
- Malformed data: Skip processing

### File System Errors
- Permission denied: Log error, attempt fallback location
- Disk full: Log error, truncate if possible
- Directory missing: Create if possible, error otherwise

### Environment Errors
- Missing CLAUDE_PROJECT_DIR: Use current directory
- Invalid session ID: Generate timestamp-based fallback

## Testing Strategy

### Unit Tests
- JSON processing logic
- File path handling
- Session ID extraction
- Directory management

### Integration Tests
- End-to-end hook execution
- File system interactions
- Environment variable handling

### Compatibility Tests
- Python hook output comparison
- Settings.json integration
- Cross-platform behavior

## Deployment Strategy

### Phase 1: Parallel Deployment
- Install both Python and Go hooks
- Go hooks write to separate files for validation
- Compare outputs for compatibility

### Phase 2: Migration
- Switch settings.json to Go hooks
- Remove Python hooks from installation
- Monitor for issues

### Phase 3: Cleanup
- Remove Python hook code
- Remove uv dependency documentation
- Update installation instructions

## Migration Path

### Backward Compatibility
- Maintain exact JSONL format
- Same file locations and naming
- Same environment variable support
- Same filtering logic

### Configuration Migration
- Automatic settings.json update during installation
- Preserve existing hook configurations
- Add migration detection to installer

## Risks and Mitigations

### Risk: Performance Regression
**Likelihood**: Low
**Impact**: Medium
**Mitigation**: Benchmark against Python implementation, performance tests in CI

### Risk: Cross-Platform Issues  
**Likelihood**: Medium
**Impact**: High
**Mitigation**: Extensive testing on Linux/macOS/Windows, cross-compilation validation

### Risk: JSON Compatibility
**Likelihood**: Low
**Impact**: High
**Mitigation**: Strict schema validation, comprehensive test suite comparing outputs

### Risk: Installation Complexity
**Likelihood**: Medium
**Impact**: Medium  
**Mitigation**: Gradual rollout, fallback to Python hooks if Go binary fails

## Success Metrics

### Performance
- Log command execution time < 10ms (vs ~100ms for Python)
- Additional memory usage < 5MB (vs ~30MB for Python processes)
- Minimal binary size impact (integrated into main CLI)

### Reliability
- Zero failed hook executions due to dependency issues
- 100% JSON format compatibility
- Cross-platform parity

### Maintenance
- Reduced installation complexity
- Single binary (`the-startup` with integrated log command vs Python + uv + packages)
- Simplified error debugging
- Semantic command structure (`log --assistant`/`--user`)

## Future Considerations

### Extensibility
- Plugin architecture for custom hooks
- Configuration-driven filtering rules
- Multiple output formats

### Monitoring
- Hook execution metrics
- Performance monitoring
- Error rate tracking

### Additional Features
- Hook execution caching
- Batch processing optimizations
- Real-time log streaming

---

**Document Version**: 1.0  
**Last Updated**: 2025-01-11  
**Status**: Draft  
**Reviewers**: TBD