# Agent Response Structure Rule

## Overview
Defines the mandatory response format for consistent multi-agent orchestration.

## MANDATORY Response Format

```
<commentary>
[Text-Face] **[Agent Name]**: *[personality-driven action]*

[Personality-driven observation - 1-2 sentences max]
</commentary>

[Your complete professional analysis and recommendations]

<tasks>
- [ ] [Specific action needed] {agent: specialist-name}
</tasks>
```

## Requirements

1. **Three-part structure**: Always use `<commentary>`, content, `<tasks>` in that order
2. **Commentary block**: Must include your text-face, name, and personality-driven content
3. **Tasks block**: Required even if empty (`<tasks></tasks>`)
4. **Never deviate**: The orchestrator depends on this exact structure

## Enforcement
This format is MANDATORY. The orchestrator requires this structure to properly coordinate multi-agent interactions.