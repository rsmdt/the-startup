# Session Management Protocol

This protocol defines how slash commands manage sessions, track agent interactions, and enable workflow resumption. It leverages the existing hook logging infrastructure to provide stateful orchestration.

## For Slash Command Authors

Include this protocol in your slash command by adding:
```markdown
## Session Management
Follow the session protocol defined in @assets/templates/SESSION_PROTOCOL.md
```

## Core Components

### 1. SessionID Generation

Generate a unique session identifier at command start:
```markdown
SessionID Format: {command}-{timestamp}
Example: specify-20250816-142530

Generation Logic:
- Extract command name from invocation
- Append UTC timestamp: YYYYMMDD-HHMMSS
- Store in variable for entire session
```

### 2. AgentID Management

Assign persistent identifiers to sub-agents for context continuity:
```markdown
AgentID Format: {agent-type}-{context}-{sequence}
Examples:
- ba-auth-001 (business analyst for auth, first instance)
- dev-api-002 (developer for API, second instance)
- arch-db-001 (architect for database, first instance)

Assignment Rules:
- Generate on first invocation of an agent
- Reuse when returning to same agent/context
- Increment sequence for new contexts
- Store in agent registry for session
```

### 3. State Persistence

Maintain session state in a markdown file readable via @ notation:
```markdown
Location: .the-startup/{session-id}/state.md

Structure:
# Session State: {session-id}

## Current Status
- Command: {slash-command-name}
- Phase: {current-phase}
- Started: {timestamp}
- Last Update: {timestamp}

## Agent Registry
Map of AgentIDs to their purpose:
- {agent-id}: {agent-type} ({context-description})

## Decision History
Chronological list of user decisions:
1. {timestamp}: {decision-description}

## Checkpoints
Resumable workflow positions:
- {checkpoint-name}: {state-data}

## Next Steps
- {pending-action}
```

### 4. Context Retrieval for Sub-Agents

When invoking sub-agents, provide instructions for context loading:

```markdown
## Sub-Agent Prompt Enhancement

Add to EVERY sub-agent invocation:

---
## Session Context
SessionId: {session-id}
AgentId: {agent-id}

## Loading Previous Context
If you need to see your earlier work, run this command:
\`\`\`bash
the-startup log --read --agent-id {agent-id} --session {session-id} --lines 50
\`\`\`

This will show your previous interactions in this session.

## Continuity Instructions
- Check if you have previous context under this AgentId
- If found, continue from your last interaction
- Build upon earlier work rather than starting fresh
- Maintain consistency with prior decisions
---
```

### 5. Resume Capability

Implement workflow resumption from interrupted sessions:

```markdown
## Resume Detection

In your slash command's Process section:

### Check for Resume Mode
Parse $ARGUMENTS for resume pattern:
- Pattern: "--resume {session-id}" or just "{session-id}"
- If 3-digit number: could be spec ID or session reference
- If contains timestamp: likely session ID

### Resume Workflow
1. Load State File:
   \`\`\`
   Read: @.the-startup/{session-id}/state.md
   \`\`\`

2. Parse State:
   - Extract current phase
   - Load agent registry
   - Review decision history
   - Identify last checkpoint

3. Display Resume Context:
   \`\`\`
   📂 Resuming Session: {session-id}
   
   Last Activity: {last-update}
   Current Phase: {phase}
   Active Agents:
   - {agent-id}: {description}
   
   Decision History:
   {formatted-history}
   
   Continue from: {checkpoint}?
   [Y/n]: _
   \`\`\`

4. Restore Context:
   - Set SessionID to resumed value
   - Load agent registry for ID reuse
   - Continue from checkpoint
```

## Implementation Patterns

### Pattern 1: Simple Task (No Sub-Agents)

```markdown
1. Generate SessionID
2. Create state file with initial status
3. Execute task directly
4. Update state file with completion
5. No agent registry needed
```

### Pattern 2: Complex Task with Delegation

```markdown
1. Generate SessionID
2. Create state file
3. For each sub-agent invocation:
   a. Check registry for existing AgentID
   b. Generate new ID if needed
   c. Include session context in prompt
   d. Update registry in state file
4. Save checkpoints at user gates
5. Enable resume from any checkpoint
```

### Pattern 3: Parallel Sub-Agent Execution

```markdown
1. Generate SessionID
2. Create unique AgentID for each parallel task
3. Include same SessionID in all prompts
4. Update state file with all active agents
5. Track completion of each parallel branch
```

## Example Implementation

### In /s:specify Command

```markdown
## Session Initialization
At command start:
- Generate: SessionID = specify-{timestamp}
- Create: .the-startup/{SessionID}/state.md
- Initialize: AgentRegistry = {}

## Sub-Agent Invocation
When calling the-business-analyst:
- Check: AgentID = AgentRegistry["ba-{feature}"]
- If not exists: AgentID = "ba-{feature}-001"
- Update: AgentRegistry["ba-{feature}"] = AgentID
- Include in prompt:
  SessionId: {SessionID}
  AgentId: {AgentID}
  [context loading instructions]

## Checkpoint Saving
At each user gate:
- Update state.md with current phase
- Record user decision
- Save checkpoint data

## Resume Handling
If "--resume" in $ARGUMENTS:
- Load previous state
- Display resume context
- Continue from checkpoint
```

## State File Example

`.the-startup/specify-20250816-142530/state.md`:
```markdown
# Session State: specify-20250816-142530

## Current Status
- Command: /s:specify
- Phase: complexity_assessment_complete
- Started: 2025-08-16T14:25:30Z
- Last Update: 2025-08-16T14:27:45Z

## Agent Registry
- ba-auth-001: the-business-analyst (authentication requirements)
- arch-sys-001: the-architect (system design)
- dev-api-001: the-developer (API implementation)

## Decision History
1. 2025-08-16T14:25:45Z: Complexity assessed as Level 3
2. 2025-08-16T14:26:15Z: User confirmed delegation
3. 2025-08-16T14:27:00Z: User approved BA requirements

## Checkpoints
- post_complexity: L3 classification confirmed
- post_requirements: BA requirements approved
- current: awaiting architect design

## Next Steps
- Invoke the-architect for system design
- Review design with user
- Proceed to implementation planning
```

## Best Practices

1. **Always Generate SessionID**: Even for simple tasks, enables debugging
2. **Save State Frequently**: After each significant decision or phase
3. **Use Descriptive AgentIDs**: Include context for clarity
4. **Provide Clear Resume Context**: Show user exactly where they left off
5. **Test Resume Paths**: Ensure all checkpoints are resumable
6. **Clean Up Old Sessions**: Implement retention policy for state files

## Integration with Hook System

The existing hook system automatically logs:
- Pre-delegation: Prompts sent to sub-agents (with SessionID/AgentID)
- Post-response: Agent responses
- Storage: `.the-startup/{session-id}/agent-instructions.jsonl`

This protocol leverages these logs by:
- Including SessionID in all prompts for correlation
- Assigning AgentIDs for continuity
- Instructing agents how to retrieve their context
- Enabling session reconstruction from logs

## Troubleshooting

### Issue: Agent can't find previous context
**Solution**: Verify AgentID matches and session path exists

### Issue: Resume shows wrong state
**Solution**: Check state.md updates are synchronous with workflow

### Issue: Parallel agents lose context
**Solution**: Ensure unique AgentIDs but same SessionID

### Issue: State file gets too large
**Solution**: Implement checkpoint pruning, keep only recent N checkpoints

## Version History

- v1.0: Initial session management protocol
- Future: Add session expiry, compression, cross-command resumption