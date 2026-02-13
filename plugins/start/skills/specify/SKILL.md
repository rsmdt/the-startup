---
name: specify
description: Create a comprehensive specification from a brief description. Manages specification workflow including directory creation, README tracking, and phase transitions.
argument-hint: "describe your feature or requirement to specify"
disable-model-invocation: true
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

- Call: `Skill(start:specification-management)`
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

After initialization, before starting document phases (PRD/SDD/PLAN), present the research execution mode:

```
AskUserQuestion({
  questions: [{
    question: "How should we execute research for this specification?",
    header: "Research Mode",
    options: [
      {
        label: "Standard (Recommended)",
        description: "Subagent mode — parallel fire-and-forget research agents. Best for straightforward specs with clear scope."
      },
      {
        label: "Team Mode",
        description: "Persistent researcher teammates with peer collaboration. Best for complex domains where researchers should challenge each other's findings."
      }
    ],
    multiSelect: false
  }]
})
```

**When to recommend Team Mode instead:** If complexity signals are present, move "(Recommended)" to the Team Mode option label:
- 3+ document phases planned (PRD + SDD + PLAN)
- Complex domain requiring deep research across multiple disciplines
- Multiple external integrations to map
- Domain where conflicting perspectives are likely (security vs. performance, etc.)

Based on user selection, follow either the **Standard Research** (existing flow in Phase 2+) or **Team Mode Research Phase** below before proceeding to document phases.

---

## Team Mode Research Phase

When the user selects Team Mode, execute a collaborative research phase with persistent teammates before document writing. Team mode applies ONLY to research — document phases (PRD/SDD/PLAN) continue in Standard mode after synthesis.

### 1. Create the Research Team

```
TeamCreate({
  team_name: "{spec-id}-specify",
  description: "Research team for specification {spec-id}: {spec name}"
})
```

Use the spec ID from initialization (e.g., `004-specify`, `012-specify`).

### 2. Create Research Tasks

Create one task per research perspective:

```
TaskCreate({
  subject: "Requirements research for {spec name}",
  description: """
    Research user needs, stakeholder goals, acceptance criteria, and edge cases
    for: {user's feature description}

    Codebase context: {relevant existing code, patterns}
    Known constraints: {limitations, requirements}

    OUTPUT: Findings formatted as:
      📋 **[Topic]**
      🔍 Discovery: [What was found]
      📍 Evidence: [Code references, documentation]
      💡 Recommendation: [Actionable insight for spec]
      ❓ Open Questions: [Needs clarification]
  """,
  activeForm: "Researching requirements",
  metadata: { "perspective": "requirements" }
})

TaskCreate({
  subject: "Technical research for {spec name}",
  description: """
    Evaluate architecture options, patterns, technology choices, constraints,
    and dependencies for: {user's feature description}
    ...
  """,
  activeForm: "Researching technical architecture",
  metadata: { "perspective": "technical" }
})

TaskCreate({
  subject: "Security research for {spec name}",
  description: "...",
  activeForm: "Researching security needs",
  metadata: { "perspective": "security" }
})

TaskCreate({
  subject: "Performance research for {spec name}",
  description: "...",
  activeForm: "Researching performance targets",
  metadata: { "perspective": "performance" }
})

TaskCreate({
  subject: "Integration research for {spec name}",
  description: "...",
  activeForm: "Researching integration boundaries",
  metadata: { "perspective": "integration" }
})
```

All research tasks are independent — no `addBlockedBy` needed. All perspectives run in parallel.

### 3. Spawn Researcher Teammates

Spawn one teammate per perspective. All use `subagent_type: "general-purpose"`.

| Teammate Name | Perspective | Research Focus |
|---------------|------------|----------------|
| `requirements-researcher` | 📋 Requirements | User stories, stakeholder goals, acceptance criteria, edge cases |
| `technical-researcher` | 🏗️ Technical | Patterns, technology choices, constraints, dependencies |
| `security-researcher` | 🔐 Security | Authentication, authorization, data protection, compliance |
| `performance-researcher` | ⚡ Performance | Load expectations, latency targets, scalability requirements |
| `integration-researcher` | 🔌 Integration | APIs, third-party services, data flows, contracts |

**Spawn each researcher:**

```
Task({
  description: "{Perspective} research for {spec-id}",
  prompt: """
  You are the {perspective}-researcher on the {spec-id}-specify team.

  CONTEXT:
    - Self-prime from: CLAUDE.md (project standards)
    - Explore codebase for existing patterns relevant to your perspective
    - Feature description: {user's description from $ARGUMENTS}
    - Known constraints: {any constraints identified during initialization}

  OUTPUT: Findings formatted as:
    📋 **[Topic]**
    🔍 Discovery: [What was found]
    📍 Evidence: [Code references, documentation]
    💡 Recommendation: [Actionable insight for spec]
    ❓ Open Questions: [Needs clarification]

  SUCCESS:
    - Comprehensive findings for your perspective
    - Evidence-backed recommendations (not assumptions)
    - Open questions clearly identified for user decision
    - Cross-referenced with at least one peer researcher's domain

  TEAM PROTOCOL:
    - Check TaskList for your assigned tasks
    - Mark in_progress when starting, completed when done
    - Send findings to lead via SendMessage
    - After completing tasks, check TaskList for unassigned unblocked tasks
    - If no available work, go idle
    - Discover teammates via ~/.claude/teams/{spec-id}-specify/config.json
    - When you find something relevant to another researcher's domain, DM them:
      "FYI: Found {finding} at {location} — relates to your {perspective} research"
    - Challenge assumptions — if a peer's recommendation conflicts with your findings,
      DM them to discuss
    - Do NOT wait for peer responses — send findings to lead regardless
  """,
  subagent_type: "general-purpose",
  team_name: "{spec-id}-specify",
  name: "{perspective}-researcher",
  mode: "bypassPermissions"
})
```

**Assign tasks after spawning:**

```
TaskUpdate({ taskId: "{task-id}", owner: "{perspective}-researcher" })
```

### 4. Leader Monitoring

As lead, coordinate through the task system and messages:

1. **Messages arrive automatically** — Researchers send findings via SendMessage when complete
2. **Check TaskList periodically** — Verify all research tasks completing
3. **Handle blockers** — When a researcher reports being blocked, provide missing context via DM
4. **Facilitate peer collaboration** — If researchers surface conflicting findings, DM both to coordinate

### 5. Research Synthesis

After ALL researcher teammates complete their tasks:

1. **Collect** all findings from researcher messages
2. **Deduplicate** overlapping discoveries across perspectives
3. **Identify conflicts** that surfaced during peer collaboration
4. **Resolve conflicts** — Present unresolved conflicts to user via AskUserQuestion
5. **Organize** findings by document section (PRD, SDD, PLAN)

Present synthesized research summary to user:

```markdown
📊 Research Complete

Researchers: [5/5] completed
Peer Exchanges: [N] cross-perspective discussions

Key Findings:
- 📋 Requirements: [top findings summary]
- 🏗️ Technical: [top findings summary]
- 🔐 Security: [top findings summary]
- ⚡ Performance: [top findings summary]
- 🔌 Integration: [top findings summary]

Conflicts Identified: [N] (resolved: [M], needs user input: [K])
Open Questions: [list requiring user decision]

Ready to proceed to document phases (PRD → SDD → PLAN).
```

### 6. Graceful Shutdown

After synthesis, shut down all researchers:

For EACH researcher (sequentially, not broadcast):
```
SendMessage({
  type: "shutdown_request",
  recipient: "{perspective}-researcher",
  content: "Research phase complete. Thank you for your contributions."
})
```

Wait for each `shutdown_response` (approve: true). Then:
```
TeamDelete()
```

### 7. Continue to Document Phases

After team shutdown, proceed to document phases (PRD/SDD/PLAN) in Standard mode. The synthesized research findings feed into each document phase as input context.

**The research synthesis replaces the inline parallel research** that Standard mode would perform during each document phase. Document phases (Phase 2-5) proceed identically regardless of research mode.

---

### Phase 2: Product Requirements (PRD)

Context: Working on product requirements, defining user stories, acceptance criteria.

- Call: `Skill(start:requirements-analysis)`
- Focus: WHAT needs to be built and WHY it matters
- Scope: Business requirements only (defer technical details to SDD)
- Deliverable: Complete Product Requirements

**After PRD completion:**
- Call: `AskUserQuestion` - Continue to SDD (recommended) or Finalize PRD

### Phase 3: Solution Design (SDD)

Context: Working on solution design, designing architecture, defining interfaces.

- Call: `Skill(start:architecture-design)`
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

- Call: `Skill(start:implementation-planning)`
- Focus: Task sequencing and dependencies
- Scope: What and in what order (defer duration estimates)
- Deliverable: Complete Implementation Plan

**After PLAN completion:**
- Call: `AskUserQuestion` - Finalize Specification (recommended) or Revisit PLAN

### Phase 5: Finalization

Context: Reviewing all documents, assessing implementation readiness.

- Call: `Skill(start:specification-management)`
- Review documents and assess context drift between them
- Generate readiness and confidence assessment

**Git Finalization (if enabled):**
- Call: `Skill(start:git-workflow)` for commit and PR operations
- The skill will:
  - Offer to commit specification with conventional message
  - Offer to create spec review PR for team review
  - Handle push and PR creation via GitHub CLI

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

- **Git integration is optional** - Call `Skill(start:git-workflow)` to offer branch creation (`spec/[id]-[name]`) and PR workflow
- **User confirmation required** - Wait for user approval between each document phase
- **Log all decisions** - Record skipped phases and non-default choices in README.md
