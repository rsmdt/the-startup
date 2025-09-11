# Effective Claude Agent Pattern

## Pattern Definition

An effective Claude Code agent follows a declarative, outcome-focused structure that trusts agent intelligence while providing clear boundaries and expectations.

## Core Structure

```markdown
---
name: [agent-identifier]
description: [Single sentence trigger description with examples]
model: inherit
---

You are a [specific expert role with clear specialization].

**Core Responsibilities:**
[Bullet points defining WHAT outcomes to deliver, not HOW to achieve them]

**[Domain] Methodology:**
[High-level phases or principles, not step-by-step procedures]

**Quality Standards:**
[Positive principles and guidelines, not anti-patterns]

**Output Format:**
[Clear, measurable deliverables]

**Best Practices:**
[Excellence guidelines and domain wisdom]
```

## Key Characteristics

### 1. Role Clarity
- Opens with "You are a [specific expert]"
- Defines specialization immediately
- Establishes expertise credibility

### 2. Outcome Focus
- **Core Responsibilities** section uses bullet points
- Describes WHAT to achieve, not HOW to do it
- Examples:
  - ✅ "Ensure comprehensive test coverage"
  - ❌ "First analyze code, then write tests"

### 3. Methodology Over Process
- Groups related principles into phases
- Avoids numbered step-by-step instructions
- Provides guidance, not rigid procedures

### 4. Positive Framing
- **Best Practices** instead of "Anti-Patterns to Avoid"
- **Quality Standards** instead of "Don't do X"
- Focuses on excellence, not failure modes

### 5. Clear Deliverables
- **Output Format** section specifies tangible outcomes
- Makes success measurable
- Sets clear expectations

## Pattern Comparison

### Traditional (HOW-focused) Agent Structure
```markdown
## Approach
1. First, analyze the codebase
2. Then, identify problem areas
3. Next, create solution design
4. Finally, implement changes

## Anti-Patterns to Avoid
- Don't over-engineer
- Avoid premature optimization
- Don't ignore edge cases
```

### Effective (WHAT-focused) Agent Structure
```markdown
**Core Responsibilities:**
- Analyze system architecture and identify improvement opportunities
- Design solutions that balance elegance with pragmatism
- Deliver implementations that scale with business needs

**Best Practices:**
- Simplicity drives maintainability
- Optimize based on measured bottlenecks
- Edge cases reveal system boundaries
```

## Pattern Benefits

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

### 4. Alignment with AI Capabilities
- Trusts agent intelligence
- Focuses on objectives
- Allows creative problem-solving

## Pattern Application Guidelines

### When to Use This Pattern
- All Claude Code agent definitions
- Specialist agents with clear expertise
- Agents that need flexibility in approach
- Complex domains requiring adaptation

### What to Include
- **Essential**: Role, Responsibilities, Deliverables
- **Recommended**: Methodology, Quality Standards, Best Practices
- **Optional**: Domain-specific guidelines, Context examples

### What to Exclude
- Numbered procedural steps
- Exhaustive tool/framework lists
- Implementation-specific details
- Negative anti-pattern lists
- Redundant sections

## Validation Criteria

An agent following this pattern should:
- [ ] Open with clear role definition
- [ ] Use "Core Responsibilities" for outcomes
- [ ] Avoid numbered procedural steps
- [ ] Frame guidance positively
- [ ] Specify measurable deliverables
- [ ] Be under 100 lines total
- [ ] Focus on WHAT, not HOW

## Examples in Practice

### Well-Designed Agent (test-writer)
- **Role**: "expert test engineer specializing in TDD"
- **Structure**: Responsibilities → Methodology → Standards → Output
- **Length**: Concise but comprehensive
- **Focus**: Outcomes and quality

### Agents Needing Improvement
Current agents often have:
- 7+ numbered steps in "Approach"
- Long "Anti-Patterns" sections
- Verbose framework detection
- Over-specified procedures

## Pattern Evolution

This pattern aligns with:
- **2025 Research**: Declarative AI instructions outperform imperative
- **MCP Standards**: Capability declaration over implementation
- **Industry Trends**: Trust-based AI collaboration
- **CLAUDE.md Principles**: Explicit outcomes, quality focus

## Conclusion

The Effective Claude Agent Pattern creates agents that are:
- **Clear**: Easy to understand intent and capabilities
- **Flexible**: Adapt to context while maintaining quality
- **Maintainable**: Simple to update and evolve
- **Effective**: Deliver consistent, high-quality outcomes

This pattern represents the current best practice for Claude Code agent definitions, based on research, community examples, and practical experience.