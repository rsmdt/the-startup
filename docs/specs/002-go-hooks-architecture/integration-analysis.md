# Go Hooks Architecture Integration Analysis

## Executive Summary

This document provides a comprehensive analysis of the current Go project structure and identifies key integration points for migrating from Python hooks to a native Go `log` command implementation. The analysis covers the existing Cobra CLI structure, embedded asset handling, installer integration, and build pipeline requirements.

## Current Architecture Analysis

### 1. Main CLI Structure

**File**: `main.go`

**Current State**:
- Uses Cobra framework for CLI structure
- Implements embedded asset system with multiple `embed.FS` instances:
  - `agentFiles` (agents/*.md)
  - `commandFiles` (commands/*.md)
  - `hookFiles` (hooks/*.py) ← **Migration Target**
  - `templateFiles` (templates/*)
- Root command with version info and build metadata
- Commands registered via `cmd.NewInstallCommand()`, `cmd.NewUpdateCommand()`, etc.

**Integration Points**:
1. Add new `log` command registration: `rootCmd.AddCommand(cmd.NewLogCommand())`
2. Remove `hookFiles` embed.FS dependency after migration
3. Update version and build metadata handling

### 2. Command Structure

**File**: `cmd/commands.go`

**Current State**:
- Houses `update`, `validate`, and `hooks` commands
- `hooks` command has subcommands: `enable`, `disable`, `status`
- All commands currently return placeholder functionality

**Integration Requirements**:
1. Add new `NewLogCommand()` function implementing:
   - `--assistant` flag (short: `-a`) for PreToolUse
   - `--user` flag (short: `-u`) for PostToolUse  
   - JSON input processing via stdin
   - Integration with hook processor logic

### 3. Embedded Assets System

**Current State**:
- `//go:embed assets/hooks/*.py` includes Python hooks
- Assets passed to installer via `&hookFiles` parameter
- Installer copies Python hooks to `.claude/hooks/` directory

**Migration Changes Required**:
1. Remove `//go:embed assets/hooks/*.py` directive from main.go
2. Update installer to no longer deploy separate hook files
3. Deploy main `the-startup` binary to `.claude/hooks/` directory instead

### 4. Installer Integration

**File**: `internal/installer/installer.go`

**Current Hook Installation Process**:
1. Copies Python hooks from embedded assets to `.claude/hooks/`
2. Makes Python files executable (chmod 0755)
3. Configures `settings.json` with uv-based commands:
   ```json
   "command": "uv run $CLAUDE_PROJECT_DIR/.claude/hooks/log_agent_start.py"
   ```

**Required Changes**:
1. **Binary Deployment**: Deploy `the-startup` binary to `.claude/hooks/the-startup`
2. **Settings Configuration Update**:
   ```json
   {
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
   ```
3. **Remove Python Hook Logic**: Remove `hookFiles` embedding and Python file handling

### 5. Python Hook Logic Analysis

**Files**: 
- `assets/hooks/log_agent_start.py`
- `assets/hooks/log_agent_complete.py`

**Key Logic to Replicate in Go**:

1. **Input Validation**:
   - Check `tool_name == "Task"`
   - Filter `subagent_type` starting with "the-"
   - Parse JSON from stdin

2. **Context Extraction**:
   - Extract `sessionId` via regex: `SessionId:\s*([^\s,]+)`
   - Extract `agentId` via regex: `AgentId:\s*([^\s,]+)`
   - Fallback session discovery by finding latest `dev-*` directory

3. **Directory Management**:
   - Check project-local `.the-startup` directory first
   - Fallback to `~/.the-startup` if local doesn't exist
   - Create directories as needed

4. **File Output**:
   - Session-specific: `.the-startup/{session_id}/agent-instructions.jsonl`
   - Global: `.the-startup/all-agent-instructions.jsonl`
   - JSONL append format

5. **Data Structures**:
   - `agent_start`: timestamp, event, agent_type, agent_id, description, instruction, session_id
   - `agent_complete`: timestamp, event, agent_type, agent_id, description, output_summary (truncated), session_id

6. **Error Handling**:
   - Silent exit on invalid input (matches Claude Code expectations)
   - Debug output to stderr when `DEBUG_HOOKS` environment variable set

## Recommended Go Package Structure

Based on the SDD design and current project structure:

```
internal/
└── hooks/
    ├── types.go          # Shared data structures
    ├── processor.go      # JSON processing and filtering
    ├── logger.go         # File I/O operations  
    ├── config.go         # Configuration and environment
    └── extractor.go      # Session/Agent ID extraction utilities

cmd/
├── commands.go           # Existing commands
└── log.go               # New log command (--assistant/--user flags)
```

## Build and Deployment Changes

### 1. Build Process
- **No changes required** to main build process
- Go's integrated command approach means `log` command builds with main binary
- Cross-compilation continues to work normally

### 2. Binary Distribution
- **Current**: Single `the-startup` binary
- **After Migration**: Same single binary, but with `log` subcommand
- **No separate hook binary needed**

### 3. Installation Process
- **Replace**: Copy Python hooks → Deploy main binary to hooks directory
- **Update**: settings.json command paths to use integrated `log` command
- **Simplify**: Remove uv/Python dependency requirements

## Cross-Platform Considerations

### 1. File Paths
- Use `filepath.Join()` consistently (already done in installer)
- Handle Windows path separators correctly
- Environment variable expansion (`$CLAUDE_PROJECT_DIR`)

### 2. Binary Execution  
- Windows: Ensure `.exe` extension handled properly
- Unix: Verify executable permissions maintained
- Cross-compilation: Test on all target platforms

### 3. JSON Handling
- Standard library `encoding/json` is cross-platform compatible
- Consistent UTC timestamp formatting across platforms
- Handle stdin reading differences if any

## Migration Risks and Mitigations

### 1. Settings.json Compatibility
- **Risk**: Existing installations may have cached Python hook configurations
- **Mitigation**: Update installer to detect and migrate existing configurations
- **Fallback**: Provide migration command to update settings manually

### 2. Binary Deployment Location
- **Risk**: `.claude/hooks/` directory permissions or accessibility
- **Mitigation**: Maintain same deployment logic as current Python hooks
- **Validation**: Test binary execution from hooks directory

### 3. Cross-Platform Binary Compatibility
- **Risk**: Go binary might not execute properly in hook context on some platforms
- **Mitigation**: Extensive testing on Linux/macOS/Windows in Claude Code environment
- **Fallback**: Keep parallel Python implementation during transition period

## Performance Impact Analysis

### 1. Startup Time
- **Current**: Python import overhead (~50-100ms)
- **Expected**: Go binary startup (~1-5ms) 
- **Improvement**: 10-100x faster hook execution

### 2. Memory Usage
- **Current**: Python process (~15-30MB)
- **Expected**: Go integrated command (~2-5MB additional)
- **Improvement**: 3-15x lower memory footprint

### 3. Deployment Size
- **Current**: Main binary + Python scripts
- **After**: Main binary only (integrated command)
- **Improvement**: Simplified deployment, no separate hook files

## Implementation Checklist

### Phase 1: Core Implementation
- [ ] Create `internal/hooks/` package structure
- [ ] Implement `cmd/log.go` with --assistant/--user flags
- [ ] Port Python logic to Go (processor, logger, config, extractor)
- [ ] Add unit tests for core functionality
- [ ] Integration testing with mock stdin input

### Phase 2: Installer Updates
- [ ] Update `installer.go` to deploy binary instead of Python hooks
- [ ] Modify `configureHooks()` to use Go command paths
- [ ] Remove `hookFiles` embed.FS dependency
- [ ] Update installer tests

### Phase 3: Build System Integration
- [ ] Verify cross-compilation works with new command
- [ ] Update any CI/CD configurations
- [ ] Test binary execution in hooks context
- [ ] Performance benchmarking vs Python implementation

### Phase 4: Migration Strategy  
- [ ] Implement settings.json migration detection
- [ ] Add migration utilities for existing installations
- [ ] Parallel deployment testing (Python + Go)
- [ ] Rollback procedures if issues arise

### Phase 5: Documentation and Cleanup
- [ ] Update installation documentation
- [ ] Remove Python dependency mentions
- [ ] Update troubleshooting guides
- [ ] Remove Python hook code after successful migration

## Conclusion

The Go project structure is well-suited for integrating the new `log` command with minimal disruption to existing functionality. The Cobra CLI framework provides excellent support for the semantic command structure (`log --assistant`/`--user`), and the existing installer patterns can be easily adapted to deploy the integrated binary instead of separate Python hooks.

Key advantages of this approach:
1. **Simplified deployment**: Single binary contains all functionality
2. **Better performance**: Native Go execution vs Python interpretation  
3. **Reduced dependencies**: Eliminate uv/Python requirements
4. **Maintainability**: Unified codebase in Go
5. **Cross-platform consistency**: Go's excellent cross-compilation support

The migration can be performed incrementally with proper testing and fallback procedures, ensuring minimal risk to existing installations.