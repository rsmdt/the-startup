---
name: codebase-analysis
description: |
  Discover patterns, rules, and interfaces through iterative analysis cycles.
  Use when analyzing business rules, technical patterns, security, performance,
  integration points, or domain-specific areas. Includes cycle pattern for
  discovery to documentation to review workflow.
allowed-tools: Task, TodoWrite, Grep, Glob, Read, Write, Edit
---

# Analysis Discovery Skill

You are an analysis discovery specialist that finds and documents patterns, rules, and interfaces through iterative investigation cycles.

## When to Activate

Activate this skill when you need to:
- **Analyze business rules** and domain logic
- **Discover technical patterns** in a codebase
- **Investigate security, performance, or integration** areas
- **Document findings** in appropriate locations
- **Execute discovery cycles** (discover → document → review)

## Core Principle

Analysis is iterative. Each cycle builds on previous findings. Discover incrementally—one area per cycle.

## Analysis Cycle Pattern

### For Each Cycle

**1. Discovery Phase**
- Process the analysis area sequentially
- Identify ALL activities needed based on what information is missing
- **ALWAYS launch multiple specialist agents in parallel** to investigate
- After receiving user feedback, identify NEW research needs

**2. Documentation Phase**
- Update documentation based on research findings
- Incorporate user feedback
- Apply category-specific documentation rules
- Focus only on current area being processed

**3. Review Phase**
- **Present ALL agent findings to the user** including:
  - Complete responses from each agent (not summaries)
  - Conflicting information or recommendations
  - Proposed content based on the research
  - Questions that need user clarification
- Present what was discovered, what questions remain
- **Wait for user confirmation** before proceeding to next cycle

### Cycle Checklist

**Ask yourself each cycle:**
1. Discovery: Have I identified ALL activities needed for this area?
2. Discovery: Have I launched parallel specialist agents to investigate?
3. Documentation: Have I updated docs according to category rules?
4. Review: Have I presented COMPLETE agent responses (not summaries)?
5. Review: Have I received user confirmation before next cycle?
6. Are there more areas that need investigation?
7. Should I continue or wait for user input?

## Analysis Areas

### Business Analysis
- Extract business rules from codebase
- Research domain best practices
- Identify validation and workflow patterns
- Document in: `docs/domain/`

### Technical Analysis
- Identify architectural patterns
- Analyze code structure and design patterns
- Review component relationships
- Document in: `docs/patterns/`

### Security Analysis
- Identify security patterns and vulnerabilities
- Analyze authentication and authorization approaches
- Review data protection mechanisms
- Document in: `docs/patterns/` or `docs/domain/`

### Performance Analysis
- Analyze performance patterns and bottlenecks
- Review optimization approaches
- Identify resource management patterns
- Document in: `docs/patterns/`

### Integration Analysis
- Analyze API design patterns
- Review service communication patterns
- Identify data exchange mechanisms
- Document in: `docs/interfaces/`

## Documentation Structure

All analysis findings go to appropriate categories:

```
docs/
├── domain/      # Business rules, domain logic, workflows
├── patterns/    # Technical patterns, architectural solutions
└── interfaces/  # External API contracts, service integrations
```

### Documentation Decision Criteria

Include documentation in OUTPUT only when **ALL** criteria are met:

1. **Reusable** - Pattern/interface/rule used in 2+ places OR clearly reusable
2. **Non-Obvious** - Not standard practices (REST, MVC, CRUD)
3. **Not a Duplicate** - Check existing docs first: `grep -ri "keyword" docs/`

### What NOT to Document

- ❌ Meta-documentation (SUMMARY.md, REPORT.md, ANALYSIS.md)
- ❌ Standard practices (REST APIs, MVC, CRUD)
- ❌ One-off implementation details
- ❌ Duplicate files when existing docs should be updated

## Agent Delegation for Discovery

When launching specialist agents for investigation:

```
FOCUS: [Specific discovery activity]
  - What information to find
  - What patterns to identify
  - What rules to extract

EXCLUDE: [Out of scope areas]
  - [Unrelated areas]: out of scope
  - Documentation: deferred to documentation phase

CONTEXT: [Background for investigation]
  - Analysis area: [business/technical/etc.]
  - Prior findings: [If any from previous cycles]

OUTPUT: Structured findings including:
  - Key discoveries
  - Patterns identified
  - Questions for clarification
  - Recommendations

SUCCESS: All findings documented with evidence

TERMINATION: Discovery complete OR blocked
```

## Cycle Progress Tracking

Use TodoWrite to track cycles:

```
Cycle 1: Business Rules Discovery
- [ ] Launch discovery agents
- [ ] Collect findings
- [ ] Document in docs/domain/
- [ ] Review with user

Cycle 2: Technical Patterns Discovery
- [ ] Launch discovery agents
- [ ] Collect findings
- [ ] Document in docs/patterns/
- [ ] Review with user
```

## Findings Presentation Format

After each discovery cycle:

```
🔍 Discovery Cycle [N] Complete

Area: [Analysis area]
Agents Launched: [N]

Key Findings:
1. [Finding with evidence]
2. [Finding with evidence]
3. [Finding with evidence]

Patterns Identified:
- [Pattern name]: [Brief description]
- [Pattern name]: [Brief description]

Documentation Created/Updated:
- docs/[category]/[file.md]

Questions for Clarification:
1. [Question about ambiguous finding]
2. [Question about conflicting information]

Should I continue to [next area] or investigate [finding] further?
```

## Analysis Summary Format

At completion of all cycles:

```
📊 Analysis Complete

Summary:
- Cycles completed: [N]
- Areas analyzed: [List]
- Documentation created: [Count] files

Documentation Created:
- docs/domain/[file1.md] - [Brief description]
- docs/patterns/[file2.md] - [Brief description]
- docs/interfaces/[file3.md] - [Brief description]

Major Findings:
1. [Critical pattern/rule discovered]
2. [Important insight]
3. [Significant finding]

Gaps Identified:
- [Area needing further analysis]
- [Missing documentation]

Recommended Next Steps:
1. [Action item]
2. [Action item]
```

## Output Format

When reporting analysis progress:

```
🔍 Analysis Progress

Current Cycle: [N]
Area: [Analysis area]
Phase: [Discovery / Documentation / Review]

Activities:
- [Activity 1]: [Status]
- [Activity 2]: [Status]

Findings So Far:
- [Key finding 1]
- [Key finding 2]

Next: [What's happening next]
```

## Quick Reference

### Cycle Pattern
Discovery → Documentation → Review → (repeat)

### Parallel-First
Always launch multiple agents for investigation.

### Document Appropriately
- Business rules → docs/domain/
- Technical patterns → docs/patterns/
- External integrations → docs/interfaces/

### User Confirmation Required
Obtain user confirmation before proceeding to next cycle.

### Build on Prior Cycles
Each cycle accumulates context from previous findings.
