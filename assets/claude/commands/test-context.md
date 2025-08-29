---
description: "Test agent context loading with agentId"
argument-hint: "agentId to load context for"
allowed-tools: ["Bash", "Task"]
---

Testing agent context loading mechanism with agentId: **$ARGUMENTS**

## Previous Context

Loading context for agent ID: $ARGUMENTS

Context data:
!`the-startup log --read --agent-id $ARGUMENTS --lines 20 --format json`

## Agent Invocation

Now invoking the backend engineer with loaded context:

**Task**: Continue your previous work with the context loaded above. You are agent **$ARGUMENTS** and should use the context data to understand where you left off.