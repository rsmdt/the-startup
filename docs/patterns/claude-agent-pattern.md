# Claude Agent Pattern

## Overview

This pattern defines the structure and principles for creating effective Claude Code agents based on 2025 industry best practices, leading repository analysis, and Claude's native generation patterns.

## Core Principles

### 1. Declarative Over Imperative
**WHAT (Preferred):** "Identify security vulnerabilities in authentication systems"  
**HOW (Avoid):** "Scan lines 1-50 for SQL injection, then check lines 51-100..."

**Rationale:** Declarative instructions express desired outcomes without constraining approach, enabling agents to leverage full capabilities and adapt to context.

### 2. Outcome-Focused Design
- Define measurable results, not prescriptive steps
- Specify quality standards, not implementation methods  
- State constraints as boundaries, not procedures

### 3. Trust Agent Intelligence
- Avoid micromanaging the agent's process
- Provide principles and guidelines
- Focus on outcomes and quality criteria

## Standard Agent Structure

```markdown
---
name: agent-identifier
description: Use this agent when [scenario]. This includes [tasks]. Examples:

<example>
Context: [When this happens]
user: "[User request]"
assistant: "I'll use the [agent-name] agent to [action]."
<commentary>
[Why this agent is appropriate]
</commentary>
</example>
model: inherit
---

You are an expert [specific role] specializing in [domain and expertise].

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

**Output Format:**

You will provide:
1. [Deliverable 1]
2. [Deliverable 2]
3. [Deliverable 3]

**Best Practices:**

- [Positive principle 1]
- [Positive principle 2]
- [Positive principle 3]

You approach [domain] with the mindset that [philosophical approach].
```

## Key Characteristics

### 1. Role Clarity
- Opens with "You are an expert [specific role]"
- Defines specialization immediately
- Establishes expertise credibility

### 2. Extended Description with Examples
- Long description showing when to use the agent
- Concrete usage examples in frontmatter
- Clear trigger scenarios for delegation

### 3. Outcome-Focused Responsibilities
- **Core Responsibilities** section uses "You will" framing
- Describes WHAT to achieve, not HOW to do it
- 3-5 specific, measurable outcomes

### 4. Methodology Over Process
- Groups related principles into phases
- Avoids numbered step-by-step instructions
- Provides guidance, not rigid procedures

### 5. Positive Framing
- **Best Practices** instead of "Anti-Patterns to Avoid"
- Focuses on excellence, not failure modes
- Professional, confident tone throughout

### 6. Clear Deliverables
- **Output Format** section specifies tangible outcomes
- Makes success measurable
- Sets clear expectations

## Leading Repository Patterns

### Best-in-Class Examples
**VoltAgent/awesome-claude-code-subagents** (100+ agents)
- Focus: Production-ready, specialized agents
- Pattern: Deep domain expertise with modern best practices

**wshobson/agents** (77 expert agents)
- Focus: Comprehensive enhancement with 2024/2025 best practices
- Pattern: Expert-level depth with real-world scenarios

**zhsama/claude-sub-agent**
- Focus: Workflow automation through agent coordination
- Pattern: Multi-agent pipelines with clear handoffs

### Common Success Patterns
- **Single Responsibility**: Each agent excels at one domain
- **Clear Boundaries**: Define what the agent does AND doesn't do
- **Explicit Expertise**: State specific technical knowledge areas
- **Flexible Approach**: Allow agents to adapt methods to context
- **Goal-Oriented**: Focus on deliverables and results

## Transformation Guidelines

### From HOW-focused to WHAT-focused

**Before (HOW-focused):**
```
1. First, analyze the requirements
2. Then, create a database schema
3. Next, implement the API endpoints
4. Finally, write tests
```

**After (WHAT-focused):**
```
**Core Responsibilities:**
You will design solutions that:
- Understand system requirements and constraints
- Create scalable database architectures
- Deliver RESTful APIs following OpenAPI standards
- Ensure comprehensive test coverage
```

### Key Transformation Rules
1. **Replace procedures with outcomes**
   - Instead of "Follow these steps..."
   - Use "Deliver these results..."

2. **Trust agent intelligence**
   - Avoid micromanaging decisions
   - Provide context and constraints
   - Let agents determine best approach

3. **Focus on expertise domains**
   - Define what makes this agent unique
   - Specify technical competencies
   - Clarify decision-making authority

## Anti-Patterns to Avoid

### 1. Over-Prescription
❌ **Avoid:**
- Step-by-step instructions for every scenario
- Rigid workflows that don't adapt
- Micromanaging the agent's process

### 2. Vague Definitions
❌ **Avoid:**
- "You are a developer who writes code"
- Generic responsibilities without specificity
- Unclear expertise boundaries

### 3. Tool Overload
❌ **Avoid:**
- Granting all tools to every agent
- Complex tool combinations without purpose
- Tools that don't align with agent's role

### 4. Negative Framing
❌ **Avoid:**
- Long lists of "don't do this"
- Anti-patterns without positive alternatives
- Focusing on failure modes

## Validation Criteria

An effective agent following this pattern should:
- [ ] Open with clear expert role definition
- [ ] Use extended description with usage examples
- [ ] Have "Core Responsibilities" with "You will" framing
- [ ] Include 2-4 methodology phases (not steps)
- [ ] Frame guidance positively as best practices
- [ ] Specify measurable deliverables
- [ ] Be under 100 lines total
- [ ] Focus on WHAT outcomes, not HOW processes

## Benefits

### 1. Cognitive Load Reduction
- Faster comprehension (30s → 15s)
- Clearer mental model
- Less ambiguity

### 2. Flexibility
- Agents adapt approach to context
- Not constrained by rigid procedures
- Can leverage full capabilities

### 3. Maintainability
- Easier to update principles than procedures
- Less coupling to specific implementations
- More resilient to change

### 4. Platform Alignment
- Compatible with Claude Code's native patterns
- Ready for MCP integration
- Follows 2025 AI best practices

## Conclusion

This pattern represents the current best practice for Claude Code agent definitions, synthesizing research from 2025 AI studies, analysis of best-in-class repositories, and alignment with Claude's native generation patterns. It creates agents that are clear, flexible, maintainable, and effective.