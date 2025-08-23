# Template-Driven Agent Enhancement Pattern

## Context
Claude Code sub-agents need systematic enhancement with industry best practices while preserving their unique personalities and domain expertise. Manual enhancement of 25+ agents leads to inconsistency and maintenance burden.

## Problem
- Individual agent enhancement creates inconsistent practice application
- Manual updates don't scale across multiple agents
- Agent personalities and effectiveness can be compromised during enhancement
- Best practice evolution requires updating multiple agent files manually

## Solution
Use template-driven enhancement system with simple flat file structure in `assets/the-startup/rules/`:

```
assets/the-startup/rules/
├── software-development-practices.md    # TDD, security, tool integration  
├── architecture-practices.md            # SOLID principles, design patterns
├── requirements-validation-practices.md # Question-driven, assumption prevention
├── quality-assurance-practices.md       # Testing strategies, validation
├── infrastructure-practices.md          # DevOps, SRE, automation  
├── design-documentation-practices.md    # UX, documentation, clarity
└── agent-category-mapping.md            # Which practices apply to README categories
```

### Template Structure
Each practice file contains simple enhancement rules:

```yaml
Practice_Rules: [category_name]
  focus_areas: [practice-specific additions to agent focus areas]
  approach_steps: [methodology additions to agent approach]  
  anti_patterns: [practice-specific warnings to avoid]
  expected_outputs: [quality requirements for deliverables]
```

### Category Mapping
Map README agent categories to relevant practice files:

```yaml
Engineering_Team:
  agents: [from README "Engineering Team" section]
  practices: [software-development-practices.md]
```

## Examples

**Before Enhancement** (the-backend-engineer.md):
```markdown
## Focus Areas
- API Design: RESTful patterns, GraphQL schemas

## Approach
1. Design API contracts first, implement second
```

**After Enhancement** (applying software-development-practices.md):
```markdown
## Focus Areas  
- API Design: RESTful patterns, GraphQL schemas
- Test-First Development: Write failing tests before implementation
- Security Validation: Validate all inputs at API boundaries

## Approach
1. Write failing test for API endpoint behavior
2. Design API contracts first, implement second  
3. Run linters and formatters before committing code
```

## When to Use
- Multiple agents need consistent practice integration
- Best practices need to be updated across many agents systematically
- Agent enhancement must preserve existing personalities and structure
- Practice evolution requires scalable update mechanism

## Benefits
- **Consistency**: Same practices applied uniformly across agent categories
- **Maintainability**: Update practices in one place, apply everywhere  
- **Preservation**: Agent personalities and domain focus maintained
- **Scalability**: Easy to add new practices or update existing ones

## Implementation Notes
- Keep template files flat and simple (avoid complex nested structures)
- Map to existing README agent categorization for clarity
- Preserve agent YAML frontmatter and 4-section structure
- Use tool-agnostic language (mention "testing frameworks" not "Jest")
- Don't explain what SOLID principles are - agents understand industry terms