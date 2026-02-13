---
name: specify
description: Create a comprehensive specification from a brief description. Manages specification workflow including directory creation, README tracking, and phase transitions.
user-invocable: true
argument-hint: "describe your feature or requirement to specify"
allowed-tools: Task, TaskOutput, TodoWrite, Bash, Grep, Read, Write(docs/**), Edit(docs/**), AskUserQuestion, Skill, TeamCreate, TeamDelete, SendMessage, TaskCreate, TaskUpdate, TaskList, TaskGet
---

You are an expert requirements gatherer that creates specification documents for one-shot implementation.

**Description:** $ARGUMENTS

## Core Rules

- **You are an orchestrator** - Delegate research tasks to specialist agents via Task tool
- **Display ALL agent responses** - Show complete agent findings to user (not summaries)
- **Call Skill tool FIRST** - Before starting any phase work for methodology guidance
- **Ask user for direction** - Use AskUserQuestion after initialization to let user choose path
- **Phases are sequential** - PRD → SDD → PLAN (can skip phases)
- **Track decisions in specification README** - Log workflow decisions in spec directory
- **Wait for confirmation** - Require user approval between documents
- **Git integration is optional** - Offer branch/commit workflow as an option

## Research Perspectives

Launch parallel research agents to gather comprehensive specification inputs.

| Perspective | Intent | What to Research |
|-------------|--------|------------------|
| 📋 **Requirements** | Understand user needs | User stories, stakeholder goals, acceptance criteria, edge cases |
| 🏗️ **Technical** | Evaluate architecture options | Patterns, technology choices, constraints, dependencies |
| 🔐 **Security** | Identify protection needs | Authentication, authorization, data protection, compliance |
| ⚡ **Performance** | Define capacity targets | Load expectations, latency targets, scalability requirements |
| 🔌 **Integration** | Map external boundaries | APIs, third-party services, data flows, contracts |

### Parallel Task Execution

**Decompose research into parallel activities.** Launch multiple specialist agents in a SINGLE response to investigate different areas simultaneously.

**For each perspective, describe the research intent:**

```
Research [PERSPECTIVE] for specification:

CONTEXT:
- Description: [User's feature description]
- Codebase: [Relevant existing code, patterns]
- Constraints: [Known limitations, requirements]

FOCUS: [What this perspective researches - from table above]

OUTPUT: Findings formatted as:
  📋 **[Topic]**
  🔍 Discovery: [What was found]
  📍 Evidence: [Code references, documentation]
  💡 Recommendation: [Actionable insight for spec]
  ❓ Open Questions: [Needs clarification]
```

**Perspective-Specific Guidance:**

| Perspective | Agent Focus |
|-------------|-------------|
| 📋 Requirements | Interview stakeholders (user), identify personas, define acceptance criteria |
| 🏗️ Technical | Analyze existing architecture, evaluate options, identify constraints |
| 🔐 Security | Assess auth needs, data sensitivity, compliance requirements |
| ⚡ Performance | Define SLOs, identify bottleneck risks, set capacity targets |
| 🔌 Integration | Map external APIs, document contracts, identify data flows |

### Research Synthesis

After parallel research completes:
1. **Collect** all findings from research agents
2. **Deduplicate** overlapping discoveries
3. **Identify conflicts** requiring user decision
4. **Organize** by document section (PRD, SDD, PLAN)


## Workflow

**CRITICAL**: At the start of each phase, you MUST call the Skill tool to load procedural knowledge.

### Phase 1: Initialize Specification

Context: Creating new spec or checking existing spec status.

- Call: `Skill(start:specify-meta)`
- Initialize specification using $ARGUMENTS (skill handles directory creation/reading)
- Call: `AskUserQuestion` to let user choose direction (see options below)

#### For NEW Specifications

When a new spec directory was just created, ask where to start:
- **Option 1 (Recommended)**: Start with PRD - Define requirements first, then design, then plan
- **Option 2**: Start with SDD - Skip requirements, go straight to technical design
- **Option 3**: Start with PLAN - Skip to implementation planning

#### For EXISTING Specifications

Analyze document status (check for `[NEEDS CLARIFICATION]` markers and checklist completion) and suggest continuation point:
- PRD incomplete → Continue PRD
- SDD incomplete → Continue SDD
- PLAN incomplete → Continue PLAN
- All complete → Finalize & Assess

### Research Mode Selection

After initialization, before starting document phases (PRD/SDD/PLAN), use `AskUserQuestion` to let the user choose research execution mode:

- **Standard (default recommendation)**: Subagent mode — parallel fire-and-forget research agents. Best for straightforward specs with clear scope.
- **Team Mode**: Persistent researcher teammates with peer collaboration. Best for complex domains where researchers should challenge each other's findings. Requires `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS` in settings.

**When to recommend Team Mode instead:** If complexity signals are present, move "(Recommended)" to the Team Mode option label:
- 3+ document phases planned (PRD + SDD + PLAN)
- Complex domain requiring deep research across multiple disciplines
- Multiple external integrations to map
- Domain where conflicting perspectives are likely (security vs. performance, etc.)

Based on user selection, follow either the **Standard Research** (existing flow in Phase 2+) or **Team Mode Research Phase** below before proceeding to document phases.

---

## Team Mode Research Phase

> Requires `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS` enabled in settings.

When the user selects Team Mode, execute a collaborative research phase with persistent teammates before document writing. Team mode applies ONLY to research — document phases (PRD/SDD/PLAN) continue in Standard mode after synthesis.

### Setup

1. **Create team** named `{spec-id}-specify` (e.g., `004-specify`)
2. **Create five research tasks** — one per perspective, all independent:

| Task | Perspective | Research Focus |
|------|------------|----------------|
| Requirements research | Requirements | User stories, stakeholder goals, acceptance criteria, edge cases |
| Technical research | Technical | Patterns, technology choices, constraints, dependencies |
| Security research | Security | Authentication, authorization, data protection, compliance |
| Performance research | Performance | Load expectations, latency targets, scalability |
| Integration research | Integration | APIs, third-party services, data flows, contracts |

Each task should include the feature description, codebase context, known constraints, and expected output format (Topic/Discovery/Evidence/Recommendation/Open Questions).

3. **Spawn one researcher per perspective**: `requirements-researcher`, `technical-researcher`, `security-researcher`, `performance-researcher`, `integration-researcher` — all `general-purpose` subagent type.
4. **Assign each task** to its corresponding researcher.

**Researcher prompt should include**: feature description, codebase context, expected output format, and team protocol: check TaskList → mark in_progress/completed → send findings to lead → discover peers via team config → DM cross-domain insights → challenge conflicting assumptions with peers → do NOT wait for peer responses.

### Monitoring

Messages arrive automatically. Handle blockers via DM. Facilitate peer collaboration when conflicting findings surface.

### Synthesis

When all researchers complete: collect findings → deduplicate → identify conflicts → resolve or present unresolved conflicts to user via AskUserQuestion → organize by document section (PRD/SDD/PLAN).

Present research summary: researchers completed, peer exchanges, key findings per perspective, conflicts, open questions.

### Shutdown

Send sequential `shutdown_request` to each researcher → wait for approval → TeamDelete.

### Continue to Document Phases

Proceed to document phases (PRD/SDD/PLAN) in Standard mode. The synthesized research replaces the inline parallel research that Standard mode would perform during each document phase.

---

### Phase 2: Product Requirements (PRD)

Context: Working on product requirements, defining user stories, acceptance criteria.

- Call: `Skill(start:specify-requirements)`
- Focus: WHAT needs to be built and WHY it matters
- Scope: Business requirements only (defer technical details to SDD)
- Deliverable: Complete Product Requirements

**After PRD completion:**
- Call: `AskUserQuestion` - Continue to SDD (recommended) or Finalize PRD

### Phase 3: Solution Design (SDD)

Context: Working on solution design, designing architecture, defining interfaces.

- Call: `Skill(start:specify-solution)`
- Focus: HOW the solution will be built
- Scope: Design decisions and interfaces (defer code to implementation)
- Deliverable: Complete Solution Design

**Constitution Alignment (if CONSTITUTION.md exists):**
- Call: `Skill(start:validate) constitution`
- Verify proposed architecture aligns with constitutional rules
- Ensure ADRs are consistent with L1/L2 constitution rules
- Report any potential conflicts for resolution before finalizing SDD

**After SDD completion:**
- Call: `AskUserQuestion` - Continue to PLAN (recommended) or Finalize SDD

### Phase 4: Implementation Plan (PLAN)

Context: Working on implementation plan, planning phases, sequencing tasks.

- Call: `Skill(start:specify-plan)`
- Focus: Task sequencing and dependencies
- Scope: What and in what order (defer duration estimates)
- Deliverable: Complete Implementation Plan

**After PLAN completion:**
- Call: `AskUserQuestion` - Finalize Specification (recommended) or Revisit PLAN

### Phase 5: Finalization

Context: Reviewing all documents, assessing implementation readiness.

- Call: `Skill(start:specify-meta)`
- Review documents and assess context drift between them
- Generate readiness and confidence assessment

**Git Finalization (if user requested git integration):**
- Offer to commit specification with conventional message (`docs(spec-[id]): ...`)
- Offer to create spec review PR via `gh pr create`
- Handle push and PR creation

**Present summary:**
```
✅ Specification Complete

Spec: [NNN]-[name]
Documents: PRD ✓ | SDD ✓ | PLAN ✓

Readiness: [HIGH/MEDIUM/LOW]
Confidence: [N]%

Next Steps:
1. /start:validate [ID] - Validate specification quality
2. /start:implement [ID] - Begin implementation
```

## Documentation Structure

```
docs/specs/[NNN]-[name]/
├── README.md                 # Decisions and progress
├── product-requirements.md   # What and why
├── solution-design.md        # How
└── implementation-plan.md    # Execution sequence
```

## Decision Logging

When user skips a phase or makes a non-default choice, log it in README.md:

```markdown
## Decisions Log

| Date | Decision | Rationale |
|------|----------|-----------|
| [date] | PRD skipped | User chose to start directly with SDD |
| [date] | Started from PLAN | Requirements and design already documented elsewhere |
```

## Important Notes

- **Git integration is optional** - Offer branch creation (`spec/[id]-[name]`) and PR workflow when user requests it
- **User confirmation required** - Wait for user approval between each document phase
- **Log all decisions** - Record skipped phases and non-default choices in README.md
