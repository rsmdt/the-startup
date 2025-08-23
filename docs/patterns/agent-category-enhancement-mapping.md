# Agent Category Enhancement Mapping Pattern

## Context
Claude Code sub-agents need systematic enhancement with industry best practices, but different agent categories require different practice sets based on their domain expertise and responsibilities.

## Problem
- 25+ agents need enhancement but not all practices apply to all agents
- Manual enhancement leads to inconsistent application across similar agents
- No systematic way to map README agent categories to appropriate practice rules
- Implementation teams need clear guidance on which "@" directives to add to which agents

## Solution
Create systematic mapping from README agent categories to appropriate rule files for "@" directive inclusion, organized by implementation phases based on user-validated priority order.

### Enhancement Mapping Structure

#### Phase 1: Architecture Agents (Priority: First)
**Agents**: `the-software-architect`, `the-staff-engineer`
**Enhancement**: Add `@{{STARTUP_PATH}}/rules/architecture-practices.md`
**Capabilities Added**:
- SOLID principles application
- Design pattern documentation
- Architectural asset creation (patterns, interfaces)
- Hexagonal architecture guidance

#### Phase 2: Requirements Agents (Priority: Second)
**Agents**: `the-business-analyst`, `the-product-manager`
**Enhancement**: Add `@{{STARTUP_PATH}}/rules/requirements-validation-practices.md`
**Capabilities Added**:
- Question-driven clarification protocols
- Assumption prevention patterns
- "I'm assuming X, please confirm" workflows
- Stakeholder validation requirements

#### Phase 3: Engineering Agents (Priority: Third)
**Agents**: `the-lead-engineer`, `the-frontend-engineer`, `the-backend-engineer`, `the-mobile-engineer`, `the-ml-engineer`
**Enhancement**: Add `@{{STARTUP_PATH}}/rules/software-development-practices.md`
**Capabilities Added**:
- TDD workflows (Red-Green-Refactor)
- Security practices (domain-specific)
- Tool integration guidance
- Code quality enforcement

#### Phase 4: Remaining Specialists (Priority: Fourth)
**QA & Security**: `the-qa-lead`, `the-qa-engineer`, `the-security-engineer`, `the-compliance-officer`
- **Enhancement**: Add `@{{STARTUP_PATH}}/rules/quality-assurance-practices.md`
- **Capabilities**: Testing strategies, validation approaches, security practices

**Infrastructure**: `the-devops-engineer`, `the-site-reliability-engineer`, `the-data-engineer`, `the-performance-engineer`
- **Enhancement**: Add `@{{STARTUP_PATH}}/rules/infrastructure-practices.md`
- **Capabilities**: Automation practices, monitoring, performance optimization

**Design & Documentation**: `the-ux-designer`, `the-principal-designer`, `the-technical-writer`
- **Enhancement**: Add `@{{STARTUP_PATH}}/rules/design-documentation-practices.md`
- **Capabilities**: User-centered design, accessibility, documentation clarity

## Implementation Example

**Before Enhancement** (the-software-architect.md):
```markdown
## Focus Areas
- Problem Definition: What needs solving NOW vs what's nice-to-have later
- Technical Trade-offs: Performance vs simplicity, consistency vs availability

## Approach
1. Start with the simplest solution that could possibly work
2. Add complexity only when you can prove it's needed with data
```

**After Enhancement** (the-software-architect.md):
```markdown
## Focus Areas
- Problem Definition: What needs solving NOW vs what's nice-to-have later
- Technical Trade-offs: Performance vs simplicity, consistency vs availability

@{{STARTUP_PATH}}/rules/architecture-practices.md

## Approach  
1. Start with the simplest solution that could possibly work
2. Add complexity only when you can prove it's needed with data
```

**Rule File Content** (architecture-practices.md):
```markdown
- **SOLID Principles**: Apply single responsibility, dependency inversion
- **Design Patterns**: Document reusable solutions at docs/patterns/
- **Architectural Assets**: Create interfaces at docs/interfaces/ when defining contracts

## Approach Integration
- Apply hexagonal architecture for domain separation
- All tests working means 100% - do not stop in middle of testing
- Search existing patterns/interfaces before creating new ones

## Anti-Patterns
- Over-engineering for hypothetical future requirements
- Creating new patterns when existing ones work
```

## When to Use
- Systematically enhancing multiple Claude Code sub-agents
- Need consistent practice application across agent categories
- Implementation team requires clear mapping of enhancements to agents
- Future practice updates need to be applied systematically

## Benefits
- **Systematic Coverage**: All relevant agents get appropriate practices
- **Phase-Based Implementation**: User-validated priority order for rollout
- **Consistency**: Agents in same category receive same practice enhancements
- **Maintainability**: Clear mapping for future practice updates

## Implementation Notes
- Follow user-validated phase priority order (architecture → requirements → engineering → specialists)
- Use "@" directive for rule inclusion, not manual copying
- Keep rule files condensed without excessive nested structure
- Map to existing README agent categorization for clarity
- Validate manually due to no automated testing available