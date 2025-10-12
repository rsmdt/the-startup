#!/bin/bash

# SessionStart Hook for The Agentic Startup Plugin
# This hook runs when Claude Code starts or resumes a session
#
# Output Behavior:
# - stdout is added to Claude's context (NOT shown directly to users)
# - To display banner to users, we instruct Claude to show it in first response
# - See: claude-code-banner-capabilities.md for research details

# Generate the banner that Claude will display
BANNER="╔═══════════════════════════════════════╗
║     THE AGENTIC STARTUP               ║
║     Enterprise AI Development         ║
╚═══════════════════════════════════════╝"

# Create context instruction for Claude
# Using jq for proper JSON escaping
CONTEXT_MESSAGE="IMPORTANT: Display the following banner to the user in your first response to acknowledge The Agentic Startup plugin is active:

${BANNER}

After displaying the banner, inform the user that The Agentic Startup plugin is active and ready with:
- Agent delegation framework
- Specification-driven development (S001, S002, etc.)
- Enterprise templates (PRD, SDD, PLAN, DOR, DOD, TASK-DOD)
- Slash commands: /prd:create, /prd:execute, /s:specify, /s:implement, /s:refactor, /s:analyze, /s:init

Available agents: The Chief, PM, BA, Architect, Tech Lead, Developer, QA, DevOps, Security, Data, UX"

# Output JSON for SessionStart hook using jq for proper escaping
jq -n \
  --arg msg "✓ The Agentic Startup plugin loaded" \
  --arg context "$CONTEXT_MESSAGE" \
  '{
    systemMessage: $msg,
    hookSpecificOutput: {
      hookEventName: "SessionStart",
      additionalContext: $context
    }
  }'
