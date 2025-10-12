#!/bin/bash

# Advanced SessionStart Hook for The Agentic Startup Plugin
# Provides context-aware session initialization with state detection

# Read hook payload from stdin (optional - for future enhancements)
PAYLOAD=$(cat)

# Extract session source if available (startup, resume, clear)
# SESSION_SOURCE=$(echo "$PAYLOAD" | jq -r '.source // "unknown"')

# Define the banner
BANNER="╔═══════════════════════════════════════╗
║     THE AGENTIC STARTUP               ║
║     Enterprise AI Development         ║
╚═══════════════════════════════════════╝"

# Check if project has The Agentic Startup installed
PROJECT_INSTALLED=false
if [ -d ".the-startup" ] && [ -f ".the-startup/.lock" ]; then
    PROJECT_INSTALLED=true
fi

# Build context message based on installation state
if [ "$PROJECT_INSTALLED" = true ]; then
    # Project has The Agentic Startup installed
    CONTEXT_MESSAGE="IMPORTANT: Display the following banner in your first response:

${BANNER}

The Agentic Startup plugin is ACTIVE in this project.

Available capabilities:
- **Agent Delegation**: 11 specialized agents (Chief, PM, BA, Architect, Tech Lead, Developer, QA, DevOps, Security, Data, UX)
- **Specification System**: Auto-incrementing specs (S001, S002...) with TOML metadata
- **Enterprise Templates**: PRD, SDD, PLAN, DOR, DOD, TASK-DOD
- **Slash Commands**:
  - /prd:create <description> - Create Product Requirement Document
  - /prd:execute <PRD-ID> - Execute PRD implementation
  - /s:specify <description> - Create numbered specification
  - /s:implement <spec-ID> - Execute specification plan
  - /s:refactor <description> - Refactor code with spec
  - /s:analyze <area> - Analyze and document patterns
  - /s:init - Initialize validation templates

After showing the banner, acknowledge the plugin is active and ready to assist with enterprise AI development."

    SYSTEM_MSG="✓ The Agentic Startup plugin active (project installed)"
else
    # Plugin available but not installed in project
    CONTEXT_MESSAGE="The Agentic Startup plugin is available but not yet installed in this project.

If the user wants to use The Agentic Startup patterns, suggest:
1. Run 'the-agentic-startup install' to install in this project
2. Or run 'npx the-agentic-startup install' if not globally installed

The plugin provides:
- Agent delegation framework (11 specialized agents)
- Specification-driven development
- Enterprise templates and workflows"

    SYSTEM_MSG="⚠ The Agentic Startup plugin available (not installed in project)"
fi

# Output JSON for SessionStart hook using jq for proper escaping
jq -n \
  --arg msg "$SYSTEM_MSG" \
  --arg context "$CONTEXT_MESSAGE" \
  '{
    systemMessage: $msg,
    hookSpecificOutput: {
      hookEventName: "SessionStart",
      additionalContext: $context
    }
  }'
