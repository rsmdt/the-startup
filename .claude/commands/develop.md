---
allowed-tools: ["*"]
description: "Full-cycle software development with strategic planning and specialist coordination"
argument-hint: "describe your project or feature requirements"
model: claude-3-5-sonnet-20241022
---

# 🚀 Full Development Pipeline

I'll orchestrate a complete development process for: **$ARGUMENTS**

## Initial Assessment

First, let me determine if this is a simple request or requires strategic planning:

```bash
# Quick complexity check
if [[ "$ARGUMENTS" =~ ^(fix|update|change|modify|add).*(single|one|simple|small|minor|quick) ]] || [[ $(echo "$ARGUMENTS" | wc -w) -lt 10 ]]; then
    echo "SIMPLE_REQUEST=true"
else
    echo "SIMPLE_REQUEST=false"
fi
```

!if [ "$SIMPLE_REQUEST" = "false" ]; then

## Phase 1: Strategic Assessment

Let me engage our Chief Technology Officer for strategic analysis:

<Task description="Strategic analysis and planning" 
      prompt="Provide strategic analysis for this request: '$ARGUMENTS'. Assess complexity (Simple/Medium/Complex), identify risks and unknowns, and recommend which specialists are needed. Focus on strategic assessment only - no implementation details or timelines."
      subagent_type="the-chief" />

Based on the Chief's assessment, I'll now coordinate the appropriate specialists.

!else

Since this is a straightforward request, I'll proceed directly with implementation.

!fi

## Phase 2: Requirements Clarification

<Task description="Check requirement clarity"
      prompt="Quickly assess if these requirements are clear enough to proceed or need deeper analysis: '$ARGUMENTS'. If vague, we'll need the business analyst. Otherwise, we can proceed. Be concise."
      subagent_type="general-purpose" />

## Phase 3: Specialist Coordination

Now I'll engage the appropriate specialists based on the strategic assessment. I'll:

1. **Display each specialist's commentary exactly as provided** (preserving their personality)
2. **Extract and track all tasks** in a master todo list
3. **Coordinate handoffs** between specialists
4. **Ensure systematic execution** of all identified tasks

### Task Execution Protocol

For each specialist response, I will:
- Print their `<commentary>` block verbatim
- Extract tasks from `<tasks>` blocks
- Update the master todo list
- Execute dependent tasks in proper sequence
- Run parallel tasks simultaneously where indicated

Let me begin the specialist engagement based on the assessments above...

## Implementation Notes

The specific workflow will adapt based on:
- **Simple projects**: Direct to developer → tester
- **Medium projects**: Add architect for design decisions
- **Complex projects**: Full cycle with business analyst → product manager → architect → project manager → developer → tester

Additional specialists (security engineer, data engineer, DevOps, etc.) will be engaged as identified in the strategic assessment.

---

*This command orchestrates your entire development team, ensuring professional delivery from concept to implementation.*