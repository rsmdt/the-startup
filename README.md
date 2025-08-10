# The Startup - Agent System for Development Tools

A comprehensive agent system for enhancing development workflows with specialized AI agents, hooks, and commands.

## Quick Installation

### One-Line Install

Install The Startup agent system with a single command:

```bash
curl -LsSf https://raw.githubusercontent.com/the-startup/the-startup/main/install.sh | sh
```

### Update Existing Installation

```bash
curl -LsSf https://raw.githubusercontent.com/the-startup/the-startup/main/install.sh | sh -s -- --update
```

The installer will:
1. Check prerequisites (curl, Python, uv)
2. Let you choose between:
   - Global installation: `~/.config/the-startup/` (recommended)
   - Local installation: `./.the-startup` (project-specific)
3. Select which components to install (agents, hooks, commands)
4. Configure everything automatically:
   - Installs files to the chosen location (e.g., `~/.config/the-startup/`)
   - Creates references in `~/.claude/` that point to the installed files using `@path` syntax
5. Create a lock file at `~/.config/the-startup/the-startup.lock` to track your installation

### What Gets Installed

- **12+ Specialized Agents**: From architecture to security to testing
  - Installed to: `~/.config/the-startup/agents/*.md`
  - Referenced from: `~/.claude/agents/` using `@~/.config/the-startup/agents/[agent-name].md`
- **Automated Hooks**: Log and track agent interactions
  - Installed to: `~/.config/the-startup/hooks/*.py`
  - Copied to: `~/.claude/hooks/` (made executable)
- **Custom Commands**: Enhanced development workflows
  - Installed to: `~/.config/the-startup/commands/*.md`
  - Copied to: `~/.claude/commands/`

## Features

### 🤖 Specialized Agents

Each agent is an expert in their domain:

- **the-chief**: Routes requests to the right specialist
- **the-architect**: Technical design and system architecture
- **the-business-analyst**: Requirements gathering and analysis
- **the-developer**: Implementation with TDD and clean code
- **the-product-manager**: Product specifications and roadmaps
- **the-project-manager**: Task coordination and progress tracking
- **the-security-engineer**: Security assessments and compliance
- **the-site-reliability-engineer**: Debugging and performance
- **the-data-engineer**: Database optimization and data modeling
- **the-devops-engineer**: Deployment automation and CI/CD
- **the-technical-writer**: Documentation and guides
- **the-tester**: Comprehensive testing strategies

### 🎯 Smart Commands

- **develop**: Orchestrates multiple agents for feature development
- **start**: Quick project initialization

### 🔍 Intelligent Hooks

Automatically track:
- Agent instructions and responses
- Session management
- Task progress
- Debug information

## System Requirements

- **curl**: For downloading files
- **Python 3.8+**: For hook scripts
- **uv** (optional): For Python package management
- **gum** (optional): For enhanced UI experience
- **jq** (optional): For JSON processing

## Installation Options

### Interactive Installation

The installer provides an interactive experience:

1. **Tool Selection**: Currently supports Claude Code (more tools coming)
2. **Location Choice**: 
   - Global: `~/.config/the-startup/` (recommended, works across all projects)
   - Local: `./.the-startup` (current project only)
   - Custom: Choose your own installation path
3. **Component Selection**: Choose which agents, hooks, and commands to install
4. **Automatic Configuration**: 
   - Installs all files to the chosen location
   - Creates references in `~/.claude/` using the `@path/to/file` syntax
   - Configures hooks in settings automatically

### Manual Installation

If you prefer manual setup:

```bash
# Clone the repository
git clone https://github.com/the-startup/the-startup.git
cd the-startup

# Copy files to your Claude directory
cp -r agents ~/.claude/
cp -r hooks ~/.claude/
cp -r commands ~/.claude/
cp -r rules ~/.claude/

# Make hooks executable
chmod +x ~/.claude/hooks/*.py
```

## Configuration

### Lock File

The installer creates a lock file at your installation path (e.g., `~/.config/the-startup/the-startup.lock`) to track:
- Installed version
- Installation type and path
- Installed files with checksums
- Installation date
- Component locations

### Settings

Hooks are configured in `.claude/settings.json`:

```json
{
  "hooks": {
    "PreToolUse": [{
      "matcher": "Task",
      "hooks": [{
        "type": "command",
        "command": "uv run $CLAUDE_PROJECT_DIR/.claude/hooks/log_agent_start.py",
        "_source": "the-startup"
      }]
    }]
  }
}
```

## Usage

### Using Agents

After installation, agents are available in Claude Code through reference files:

```bash
# Start a new development task
/develop "Create a user authentication system"

# The chief will analyze and route to appropriate specialists
# Agents work together to deliver comprehensive solutions
```

**How it works:**
- Agents are installed to `~/.config/the-startup/agents/`
- Claude Code reads references from `~/.claude/agents/`
- Each reference contains `@~/.config/the-startup/agents/[agent-name].md`
- This allows centralized updates while maintaining Claude compatibility

### Understanding the @path Syntax

The `@path/to/file` syntax is Claude Code's reference mechanism that allows files to point to other files:

**Example Reference File** (`~/.claude/agents/the-chief.md`):
```
@~/.config/the-startup/agents/the-chief.md
```

**Dynamic Path Resolution:**
- Agent and command files use `{{INSTALL_PATH}}` placeholders
- The installer replaces these with your actual installation path
- This allows flexible installation locations (global, local, or custom)

**Example Placeholder Usage in Agent Files:**
```markdown
# Before installation (in source)
Follow instructions in @{{INSTALL_PATH}}/rules/context-management.md

# After installation (with default path)
Follow instructions in @~/.config/the-startup/rules/context-management.md

# After installation (with custom path)
Follow instructions in @/my/custom/path/rules/context-management.md
```

**Benefits:**
- **Centralized Management**: Update agents in one location
- **Flexible Installation**: Choose any installation directory
- **Version Control**: Keep agents in the git repository, install them anywhere
- **Clean Separation**: Claude's working directory stays uncluttered
- **Easy Updates**: Pull latest changes and reinstall to update all projects

**File Structure After Installation:**
```
~/.config/the-startup/          # Main installation directory
├── agents/                     # Agent source files
│   ├── the-architect.md
│   ├── the-chief.md
│   └── ...
├── commands/                   # Command source files
│   ├── develop.md
│   └── create.md
├── hooks/                      # Hook scripts
│   └── *.py
├── templates/                  # Document templates
│   └── *.md
└── the-startup.lock           # Installation manifest

~/.claude/                      # Claude Code directory
├── agents/                     # Reference files (contain @paths)
│   ├── the-architect.md       # Contains: @~/.config/the-startup/agents/the-architect.md
│   ├── the-chief.md           # Contains: @~/.config/the-startup/agents/the-chief.md
│   └── ...
├── commands/                   # Copied command files (with placeholders replaced)
│   ├── develop.md
│   └── create.md
├── hooks/                      # Copied hook scripts (executable)
│   └── *.py
└── templates/                  # Reference files (contain @paths)
    ├── PRD.md                  # Contains: @~/.config/the-startup/templates/PRD.md
    ├── BRD.md                  # Contains: @~/.config/the-startup/templates/BRD.md
    └── ...
```

### Viewing Logs

Agent interactions are logged to:
- `~/.the-startup/all-agent-instructions.jsonl` - Global log
- `~/.the-startup/[session-id]/agent-instructions.jsonl` - Session-specific logs

View logs with:
```bash
# View recent entries
tail -f ~/.the-startup/all-agent-instructions.jsonl | jq .

# Search for specific agent
grep "the-developer" ~/.the-startup/all-agent-instructions.jsonl | jq .
```

## Claude Code Hooks Documentation

### Overview

Claude Code hooks are configurable scripts that intercept and respond to specific events during Claude Code execution. They enable you to customize behavior, log actions, validate inputs, and integrate with external systems without modifying the core Claude Code functionality.

### Key Benefits
- **Non-intrusive monitoring**: Log and track actions without modifying existing code
- **Safety controls**: Prevent dangerous operations before they execute
- **Automation**: Trigger follow-up actions automatically
- **Integration**: Connect Claude Code with external tools and services

### Security Considerations
⚠️ **Important**: Hooks execute shell commands automatically with your current environment credentials. Always:
- Review hook scripts before enabling them
- Validate and sanitize any inputs
- Use absolute paths to prevent command injection
- Test hooks in a safe environment first

## Available Hook Types

### 1. PreToolUse
**When**: Before a tool is executed  
**Can block**: Yes (return non-zero exit code)  
**Receives**: Tool name and input parameters  
**Use cases**: Input validation, logging, preventing dangerous operations

### 2. PostToolUse  
**When**: After a tool completes successfully  
**Can block**: No  
**Receives**: Tool name, input parameters, and output  
**Use cases**: Result logging, code formatting, triggering follow-ups

### 3. UserPromptSubmit
**When**: User submits a prompt to Claude  
**Can block**: Yes  
**Receives**: User's message  
**Use cases**: Request logging, input validation, context injection

### 4. Stop
**When**: Main agent finishes responding  
**Can block**: No  
**Receives**: Session information  
**Use cases**: Cleanup, final logging, summary generation

### 5. SubagentStop
**When**: A subagent (Task tool) finishes  
**Can block**: No  
**Receives**: Subagent information and results  
**Use cases**: Subagent metrics, result collection

### 6. PreCompact
**When**: Before context compaction (conversation too long)  
**Can block**: No  
**Receives**: Current context information  
**Use cases**: Save important context, create summaries

### 7. SessionStart
**When**: Starting new or resuming session  
**Can block**: No  
**Receives**: Session information  
**Use cases**: Initialize logging, setup environment

### 8. Notification
**When**: Tool permissions or idle periods  
**Can block**: No  
**Receives**: Notification type and details  
**Use cases**: Desktop alerts, status updates

## Configuration Structure

Hooks are configured in `.claude/settings.local.json`:

```json
{
  "hooks": {
    "EventName": [
      {
        "matcher": "RegexPattern",  // Optional: filter which tools/events
        "hooks": [
          {
            "type": "command",
            "command": "path/to/script.sh or python3 script.py"
          }
        ]
      }
    ]
  }
}
```

### Matcher Patterns
- Use regex to filter specific tools or patterns
- Example: `"^the-.*"` matches all agents starting with "the-"
- Omit matcher to apply to all events of that type

## Implementation Example: Agent Instruction Logging

This example demonstrates intercepting Task tool calls for agents starting with "the-" to log their instructions and responses.

### Problem Statement
We want to:
1. Capture exact instructions sent to each agent (the-chief, the-developer, etc.)
2. Log when each agent completes their task
3. Organize logs by session for debugging and analysis

### Solution Architecture

```
Task Tool Called → PreToolUse Hook → Log Instruction → Agent Executes → PostToolUse Hook → Log Response
```

### Step 1: Create Hook Directory

```bash
mkdir -p .claude/hooks
```

### Step 2: Create Start Hook Script

**File**: `.claude/hooks/log_agent_start.py`

```python
#!/usr/bin/env python3
"""
Log agent Task instructions before execution.
Only captures agents with subagent_type starting with "the-"
"""

import json
import sys
import os
import re
from datetime import datetime
from pathlib import Path

def extract_session_id(prompt):
    """Extract sessionId from prompt text."""
    match = re.search(r'SessionId:\s*([^\s,]+)', prompt)
    return match.group(1) if match else None

def extract_agent_id(prompt):
    """Extract agentId from prompt text."""
    match = re.search(r'AgentId:\s*([^\s,]+)', prompt)
    return match.group(1) if match else None

def find_latest_session(project_dir):
    """Find the most recent session directory."""
    startup_dir = Path(project_dir) / '.the-startup'
    if not startup_dir.exists():
        return None
    
    session_dirs = [d for d in startup_dir.iterdir() 
                   if d.is_dir() and d.name.startswith('dev-')]
    
    if not session_dirs:
        return None
    
    latest = max(session_dirs, key=lambda d: d.stat().st_mtime)
    return latest.name

def main():
    try:
        # Read JSON input from stdin
        input_data = json.loads(sys.stdin.read())
        
        # Check if this is a Task tool call
        if input_data.get('tool_name') != 'Task':
            sys.exit(0)
        
        tool_input = input_data.get('tool_input', {})
        subagent_type = tool_input.get('subagent_type', '')
        
        # Only process agents starting with "the-"
        if not subagent_type.startswith('the-'):
            sys.exit(0)
        
        prompt = tool_input.get('prompt', '')
        description = tool_input.get('description', '')
        
        # Get project directory
        project_dir = os.environ.get('CLAUDE_PROJECT_DIR', '.')
        
        # Extract context
        session_id = extract_session_id(prompt)
        agent_id = extract_agent_id(prompt)
        
        # Find session if not in prompt
        if not session_id:
            session_id = find_latest_session(project_dir)
        
        # Create log entry
        log_entry = {
            'timestamp': datetime.utcnow().isoformat() + 'Z',
            'event': 'agent_start',
            'agent_type': subagent_type,
            'agent_id': agent_id,
            'description': description,
            'instruction': prompt,
            'session_id': session_id
        }
        
        # Ensure directories exist
        startup_dir = Path(project_dir) / '.the-startup'
        startup_dir.mkdir(exist_ok=True)
        
        # Write to session-specific file
        if session_id:
            session_dir = startup_dir / session_id
            session_dir.mkdir(exist_ok=True)
            
            context_file = session_dir / 'agent-instructions.jsonl'
            with open(context_file, 'a') as f:
                json.dump(log_entry, f)
                f.write('\n')
        
        # Write to global log
        global_log = startup_dir / 'all-agent-instructions.jsonl'
        with open(global_log, 'a') as f:
            json.dump(log_entry, f)
            f.write('\n')
        
        # Debug output if enabled
        if os.environ.get('DEBUG_HOOKS'):
            print(f"[HOOK] Agent starting: {subagent_type} (session: {session_id}, agent: {agent_id})", 
                  file=sys.stderr)
        
    except Exception as e:
        if os.environ.get('DEBUG_HOOKS'):
            print(f"[HOOK ERROR] {e}", file=sys.stderr)
        sys.exit(0)

if __name__ == '__main__':
    main()
```

### Step 3: Create Completion Hook Script

**File**: `.claude/hooks/log_agent_complete.py`

```python
#!/usr/bin/env python3
"""
Log agent Task completion after execution.
Only captures agents with subagent_type starting with "the-"
"""

import json
import sys
import os
import re
from datetime import datetime
from pathlib import Path

def extract_session_id(prompt):
    """Extract sessionId from prompt text."""
    match = re.search(r'SessionId:\s*([^\s,]+)', prompt)
    return match.group(1) if match else None

def extract_agent_id(prompt):
    """Extract agentId from prompt text."""
    match = re.search(r'AgentId:\s*([^\s,]+)', prompt)
    return match.group(1) if match else None

def find_latest_session(project_dir):
    """Find the most recent session directory."""
    startup_dir = Path(project_dir) / '.the-startup'
    if not startup_dir.exists():
        return None
    
    session_dirs = [d for d in startup_dir.iterdir() 
                   if d.is_dir() and d.name.startswith('dev-')]
    
    if not session_dirs:
        return None
    
    latest = max(session_dirs, key=lambda d: d.stat().st_mtime)
    return latest.name

def truncate_output(output, max_length=1000):
    """Truncate output if too long."""
    if isinstance(output, str) and len(output) > max_length:
        return output[:max_length] + f"... [truncated {len(output) - max_length} chars]"
    return output

def main():
    try:
        # Read JSON input from stdin
        input_data = json.loads(sys.stdin.read())
        
        # Check if this is a Task tool call
        if input_data.get('tool_name') != 'Task':
            sys.exit(0)
        
        tool_input = input_data.get('tool_input', {})
        subagent_type = tool_input.get('subagent_type', '')
        
        # Only process agents starting with "the-"
        if not subagent_type.startswith('the-'):
            sys.exit(0)
        
        prompt = tool_input.get('prompt', '')
        description = tool_input.get('description', '')
        output = input_data.get('output', '')
        
        # Get project directory
        project_dir = os.environ.get('CLAUDE_PROJECT_DIR', '.')
        
        # Extract context
        session_id = extract_session_id(prompt)
        agent_id = extract_agent_id(prompt)
        
        # Find session if not in prompt
        if not session_id:
            session_id = find_latest_session(project_dir)
        
        # Create log entry
        log_entry = {
            'timestamp': datetime.utcnow().isoformat() + 'Z',
            'event': 'agent_complete',
            'agent_type': subagent_type,
            'agent_id': agent_id,
            'description': description,
            'output_summary': truncate_output(output),
            'session_id': session_id
        }
        
        # Ensure directories exist
        startup_dir = Path(project_dir) / '.the-startup'
        startup_dir.mkdir(exist_ok=True)
        
        # Write to session-specific file
        if session_id:
            session_dir = startup_dir / session_id
            session_dir.mkdir(exist_ok=True)
            
            context_file = session_dir / 'agent-instructions.jsonl'
            with open(context_file, 'a') as f:
                json.dump(log_entry, f)
                f.write('\n')
        
        # Write to global log
        global_log = startup_dir / 'all-agent-instructions.jsonl'
        with open(global_log, 'a') as f:
            json.dump(log_entry, f)
            f.write('\n')
        
        # Debug output if enabled
        if os.environ.get('DEBUG_HOOKS'):
            print(f"[HOOK] Agent completed: {subagent_type} (session: {session_id}, agent: {agent_id})", 
                  file=sys.stderr)
        
    except Exception as e:
        if os.environ.get('DEBUG_HOOKS'):
            print(f"[HOOK ERROR] {e}", file=sys.stderr)
        sys.exit(0)

if __name__ == '__main__':
    main()
```

### Step 4: Update Settings Configuration

**File**: `.claude/settings.local.json`

```json
{
  "permissions": {
    "allow": [
      "WebFetch(domain:docs.anthropic.com)",
      "mcp__sequential-thinking__sequentialthinking",
      "Bash(mkdir:*)",
      "Bash(ls:*)",
      "Bash(find:*)",
      "Bash(mkdir:*)",
      "Bash(mv:*)"
    ],
    "deny": []
  },
  "hooks": {
    "PreToolUse": [
      {
        "matcher": "Task",
        "hooks": [
          {
            "type": "command",
            "command": "python3 $CLAUDE_PROJECT_DIR/.claude/hooks/log_agent_start.py"
          }
        ]
      }
    ],
    "PostToolUse": [
      {
        "matcher": "Task",
        "hooks": [
          {
            "type": "command",
            "command": "python3 $CLAUDE_PROJECT_DIR/.claude/hooks/log_agent_complete.py"
          }
        ]
      }
    ]
  }
}
```

## Quick Installation

### Prerequisites
- **uv**: Python package manager (`curl -LsSf https://astral.sh/uv/install.sh | sh`)
- **Gum** (optional): For enhanced UI (`brew install gum` on macOS)

### Automated Setup

```bash
# Make installer executable
chmod +x .claude/scripts/setup.sh

# Run interactive setup
.claude/scripts/setup.sh
```

The installer will:
1. Check and install prerequisites
2. Configure hooks interactively
3. Set up logging directories
4. Test the installation
5. Provide usage instructions

### Hook Management

After installation, use the hook manager for ongoing maintenance:

```bash
# Make manager executable
chmod +x .claude/scripts/hook-manager.sh

# Run hook manager
.claude/scripts/hook-manager.sh

# Or use direct commands:
.claude/scripts/hook-manager.sh status  # Check hook status
.claude/scripts/hook-manager.sh logs    # View logs
.claude/scripts/hook-manager.sh toggle  # Enable/disable hooks
.claude/scripts/hook-manager.sh clear   # Clear logs
.claude/scripts/hook-manager.sh test    # Test hooks
```

### Manual Setup (Alternative)

If you prefer manual configuration:

```bash
# 1. Make scripts executable
chmod +x .claude/hooks/log_agent_start.py
chmod +x .claude/hooks/log_agent_complete.py
chmod +x .claude/scripts/*.sh

# 2. Ensure uv is installed
curl -LsSf https://astral.sh/uv/install.sh | sh

# 3. Hooks are already configured in settings.local.json
```

## Testing the Hooks

### 1. Enable Debug Output (Optional)
```bash
export DEBUG_HOOKS=1
```

### 2. Run a Command That Uses Agents
```bash
claude /develop "create a user authentication system"
```

### 3. Check the Logs

**Session-specific log:**
```bash
cat .the-startup/dev-*/agent-instructions.jsonl | jq .
```

**Global log:**
```bash
cat .the-startup/all-agent-instructions.jsonl | jq .
```

### Expected Output Format

```json
{
  "timestamp": "2025-01-10T15:30:45Z",
  "event": "agent_start",
  "agent_type": "the-business-analyst",
  "agent_id": "ba3k5m",
  "description": "Analyze requirements",
  "instruction": "Analyze user authentication requirements. SessionId: dev-20250110-153045-a7b9, AgentId: ba3k5m",
  "session_id": "dev-20250110-153045-a7b9"
}
```

```json
{
  "timestamp": "2025-01-10T15:31:20Z",
  "event": "agent_complete",
  "agent_type": "the-business-analyst",
  "agent_id": "ba3k5m",
  "description": "Analyze requirements",
  "output_summary": "Requirements analysis complete. Identified 5 key features...",
  "session_id": "dev-20250110-153045-a7b9"
}
```

## Troubleshooting

### Hooks Not Firing
1. Check settings file is valid JSON: `jq . .claude/settings.local.json`
2. Verify scripts are executable: `ls -la .claude/hooks/`
3. Enable debug output: `export DEBUG_HOOKS=1`
4. Check stderr output during execution

### Permission Errors
```bash
chmod +x .claude/hooks/*.py
```

### Python Not Found
Ensure Python 3 is installed and in PATH:
```bash
which python3
python3 --version
```

### Logs Not Created
1. Check `.the-startup/` directory exists and is writable
2. Verify `CLAUDE_PROJECT_DIR` environment variable is set correctly
3. Look for error messages in stderr when `DEBUG_HOOKS=1`

### Session ID Not Found
The hooks will attempt to:
1. Extract from prompt text (pattern: `SessionId: xxx`)
2. Find the latest session directory if not in prompt
3. Still log to global file even if session not found

## API Reference

### Hook Input Structure (stdin)
```json
{
  "tool_name": "Task",
  "tool_input": {
    "description": "Short description",
    "prompt": "Full instructions with context",
    "subagent_type": "the-agent-name"
  },
  "output": "Tool execution result (PostToolUse only)"
}
```

### Hook Exit Codes
- `0`: Success, continue execution
- `Non-zero`: (PreToolUse only) Block tool execution

### Environment Variables
- `CLAUDE_PROJECT_DIR`: Project root directory
- `DEBUG_HOOKS`: Enable debug output to stderr

## Advanced Usage

### Filtering Specific Agents
Modify the Python scripts to filter specific agents:
```python
ALLOWED_AGENTS = ['the-developer', 'the-architect']
if subagent_type not in ALLOWED_AGENTS:
    sys.exit(0)
```

### Adding Metrics
Extend log entries with performance metrics:
```python
log_entry['duration_ms'] = execution_time
log_entry['memory_usage'] = get_memory_usage()
```

### Integration with External Services
Send logs to external services:
```python
import requests
requests.post('https://api.logging-service.com/logs', json=log_entry)
```

## Best Practices

1. **Always handle exceptions** - Never let a hook crash block tool execution
2. **Keep hooks fast** - Long-running hooks slow down Claude Code
3. **Use structured logging** - JSON format enables easy parsing
4. **Rotate logs periodically** - Prevent unlimited growth
5. **Test in isolation** - Verify hooks work before enabling
6. **Document custom hooks** - Help future maintainers understand your setup

## Further Resources

- [Claude Code Documentation](https://docs.anthropic.com/en/docs/claude-code)
- [Hooks Guide](https://docs.anthropic.com/en/docs/claude-code/hooks-guide)
- [Hooks Reference](https://docs.anthropic.com/en/docs/claude-code/hooks)