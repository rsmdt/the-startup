#!/usr/bin/env python3
# /// script
# requires-python = ">=3.8"
# ///
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

def get_startup_dir(project_dir):
    """Get the startup directory for logs."""
    # First check project local directory
    local_startup = Path(project_dir) / '.the-startup'
    if local_startup.exists():
        return local_startup
    
    # Otherwise use home directory
    home_startup = Path.home() / '.the-startup'
    home_startup.mkdir(exist_ok=True, parents=True)
    return home_startup

def find_latest_session(project_dir):
    """Find the most recent session directory."""
    startup_dir = get_startup_dir(project_dir)
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

def extract_commentary(output):
    """Extract commentary section from agent output."""
    if not isinstance(output, str):
        return None
    
    # Look for <commentary> tags
    import re
    pattern = r'<commentary>(.*?)</commentary>'
    match = re.search(pattern, output, re.DOTALL)
    
    if match:
        return match.group(1).strip()
    return None

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
        
        # Extract commentary if present
        commentary = extract_commentary(output)
        
        # Create log entry
        log_entry = {
            'timestamp': datetime.utcnow().isoformat() + 'Z',
            'event': 'agent_complete',
            'agent_type': subagent_type,
            'agent_id': agent_id,
            'description': description,
            'output_summary': truncate_output(output),
            'commentary': commentary,  # Add commentary to log
            'session_id': session_id
        }
        
        # Get startup directory
        startup_dir = get_startup_dir(project_dir)
        
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
        
        # Display commentary in chat if present
        if commentary:
            # Print to stdout to display in chat
            print(f"\n{'='*60}", file=sys.stdout)
            print(f"📝 {subagent_type} Commentary:", file=sys.stdout)
            print(f"{'='*60}", file=sys.stdout)
            print(commentary, file=sys.stdout)
            print(f"{'='*60}\n", file=sys.stdout)
        
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