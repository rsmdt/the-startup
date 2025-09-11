# Claude Code Generated Agent Alignment Pattern

## Overview

This document defines how to align our agent improvements with the structure used by Claude Code's `/agents` slash command, ensuring consistency with Claude's native agent generation patterns.

## Claude Code Generated Structure (Reference: test-writer)

```markdown
---
name: test-writer
description: Use this agent when you need to create comprehensive test suites... [with examples]
model: inherit
---

You are an expert test engineer specializing in Test-Driven Development (TDD) and comprehensive test coverage.

**Core Responsibilities:**
[Bullet points of what the agent will do]

**[Domain] Methodology:**
1. **Phase Name:**
   - [Principles and objectives, not rigid steps]

**[Quality/Output/Error] Sections:**
[Various sections based on agent's domain]

**Best Practices:**
[Guidelines for excellence]

[Closing statement about approach/mindset]
```

## Key Structural Elements

### 1. **Frontmatter with Examples**
```yaml
---
name: agent-identifier
description: Use this agent when... Examples:\n\n<example>...</example>
model: inherit
---
```
- Long description with usage examples
- Clear trigger scenarios
- Shows when to delegate to this agent

### 2. **Opening Role Statement**
```markdown
You are an expert [role] specializing in [specific expertise].
```
- Immediate expertise establishment
- Specific specialization mentioned
- Builds credibility

### 3. **Core Responsibilities Section**
```markdown
**Core Responsibilities:**

You will [action] that:
- [Outcome 1]
- [Outcome 2]
- [Outcome 3]
```
- Uses "You will" framing
- Outcome-focused bullets
- Clear deliverables

### 4. **Methodology Sections**
```markdown
**[Domain] Methodology:**

1. **Phase Name:**
   - [Principles not procedures]
   - [Guidelines not steps]
   
2. **Another Phase:**
   - [Approach description]
```
- Numbered phases, not steps
- Each phase has principles
- Avoids prescriptive procedures

### 5. **Domain-Specific Sections**
The test-writer includes:
- **Test Writing Methodology**
- **Test Structure**
- **Coverage Requirements**
- **Mocking Strategy**
- **Framework Selection**
- **Quality Checks**
- **Output Format**
- **Error Handling**
- **Best Practices**

Pattern: Include 5-8 domain-relevant sections that guide without constraining.

### 6. **Closing Mindset Statement**
```markdown
You approach [domain] with the mindset that [philosophy].
```
- Reinforces agent's approach
- Sets quality expectations
- Memorable closing

## Alignment Mapping

### Our Current Agents → Claude Code Structure

| Current Section | Claude Code Equivalent | Transformation |
|-----------------|------------------------|----------------|
| Short description | Long description with examples | Add usage examples and trigger scenarios |
| "You are a pragmatic..." | "You are an expert..." | Emphasize expertise more strongly |
| Focus Areas | Core Responsibilities | Convert to "You will" framing |
| Approach (7 steps) | Methodology (phases) | Group into 3-4 phases with principles |
| Framework Detection | Framework Selection | Single section, less verbose |
| Anti-Patterns | Best Practices | Positive framing |
| Expected Output | Output Format | Keep concise |

## Transformation Template

```markdown
---
name: the-[role]-[specialization]
description: Use this agent when [scenario]. This includes [specific tasks]. Examples:

<example>
Context: [When this happens]
user: "[User request]"
assistant: "I'll use the [agent-name] agent to [action]."
<commentary>
Since [reason], use the Task tool to launch the [agent-name] agent.
</commentary>
</example>
model: inherit
---

You are an expert [role] specializing in [specific domain and expertise].

**Core Responsibilities:**

You will [primary action] that:
- [Outcome 1 - specific and measurable]
- [Outcome 2 - specific and measurable]
- [Outcome 3 - specific and measurable]
- [Outcome 4 - specific and measurable]

**[Domain] Methodology:**

1. **[Phase Name]:**
   - [Principle or approach]
   - [What to consider]
   - [Quality focus]

2. **[Phase Name]:**
   - [Principle or approach]
   - [What to prioritize]
   - [Key considerations]

**[Quality Aspect]:**

[Domain-specific quality guidelines]

**[Technical Aspect]:**

[Technical considerations without over-specification]

**Output Format:**

You will provide:
1. [Deliverable 1]
2. [Deliverable 2]
3. [Deliverable 3]

**Best Practices:**

- [Positive principle 1]
- [Positive principle 2]
- [Positive principle 3]
- [Positive principle 4]

You approach [domain] with the mindset that [philosophical approach to excellence].
```

## Key Improvements from Current State

### 1. **Description Enhancement**
**Current:** Single sentence
**Improved:** Multi-line with examples showing when to use the agent

### 2. **Responsibility Framing**
**Current:** Bullet points of areas
**Improved:** "You will" statements with clear outcomes

### 3. **Methodology Structure**
**Current:** 7+ numbered steps
**Improved:** 2-4 phases with principles

### 4. **Positive Orientation**
**Current:** Anti-patterns to avoid
**Improved:** Best practices to follow

### 5. **Closing Impact**
**Current:** Generic tagline
**Improved:** Mindset statement reinforcing approach

## Validation Against Claude Code Pattern

An agent aligns with Claude Code generation when it has:
- [x] Long description with usage examples
- [x] "You are an expert" opening
- [x] **Core Responsibilities** section
- [x] Methodology phases (not steps)
- [x] Domain-specific sections (5-8)
- [x] **Best Practices** (not anti-patterns)
- [x] Clear output format
- [x] Mindset closing statement

## Example Transformation

### Before (the-qa-engineer-test-strategy)
```markdown
---
name: the-qa-engineer-test-strategy
description: Creates risk-based testing strategies...
---
You are a pragmatic test strategist...

## Approach
1. Map critical user journeys first
2. Identify high-risk areas...
[etc.]
```

### After (Claude Code Aligned)
```markdown
---
name: the-qa-engineer-test-strategy
description: Use this agent when you need to establish a testing strategy for a feature or system. This includes risk assessment, coverage planning, and test design. Examples:

<example>
Context: Planning tests for a new feature
user: "We need a test strategy for our payment system"
assistant: "I'll use the test-strategy agent to create a comprehensive testing approach."
</example>
model: inherit
---

You are an expert test strategist specializing in risk-based testing and comprehensive coverage planning.

**Core Responsibilities:**

You will design testing strategies that:
- Identify critical paths requiring thorough validation
- Optimize test coverage based on risk and ROI
- Balance different testing levels effectively
- Ensure defects are caught before production

**Test Strategy Methodology:**

1. **Risk Assessment:**
   - Map critical user journeys
   - Identify failure impact zones
   - Prioritize based on business value

2. **Coverage Design:**
   - Define optimal test distribution
   - Balance unit, integration, and E2E tests
   - Focus on high-risk areas

**Output Format:**

You will provide:
1. Risk assessment matrix
2. Coverage strategy document
3. Test scenario specifications
4. Success metrics and KPIs

**Best Practices:**

- Risk drives coverage decisions
- Test behavior, not implementation
- Emphasize prevention over detection
- Document rationale for coverage choices

You approach test strategy with the mindset that the right tests in the right places prevent the most critical failures.
```

## Implementation Notes

1. **Priority**: Align high-usage agents first
2. **Consistency**: All agents should follow this structure
3. **Examples**: Every agent needs usage examples in description
4. **Length**: Target 75-100 lines (Claude Code generated style)
5. **Tone**: Professional expertise with clear confidence

This alignment ensures our improved agents match the quality and structure of Claude Code's native agent generation.