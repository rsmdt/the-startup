# Start Plugin - The Agentic Startup

**Workflow orchestration plugin for spec-driven development in Claude Code.**

The `start` plugin provides nine user-invocable workflow skills, five autonomous skills, and two output styles to transform how you build software with Claude Code.

**📖 For quick start, workflow guide, and skill selection, see the [main README](../../README.md).**

---

## Table of Contents

- [User-Invocable Skills](#user-invocable-skills) — specify, implement, validate, review, document, analyze, refactor, debug, constitution
- [Autonomous Skills](#autonomous-skills) — 5 context-activated skills
- [Documentation Structure](#-documentation-structure) — specs, domain, patterns, interfaces
- [Output Styles](#-output-styles) — The Startup, The ScaleUp
- [Typical Development Workflow](#typical-development-workflow) — primary and maintenance flows
- [Skills in Action](#skills-in-action) — real-world examples
- [Templates](#templates) — PRD, SDD, PLAN, DOR, DOD
- [Philosophy](#philosophy) — spec-driven development principles

---

## User-Invocable Skills

These skills are invoked by the user via slash commands (e.g., `/start:specify`). Unlike autonomous skills which activate automatically based on context, user-invocable skills wait for explicit invocation.

### `/start:specify <description>`

Create comprehensive specifications from brief descriptions through deep research and specialist agent coordination.

**Purpose:** Transform ideas into implementation-ready specifications with product requirements, solution design, and implementation plan documents

**Usage:**
```bash
/start:specify Build a real-time notification system with WebSocket support
/start:specify 001  # Resume existing specification work
```

**Key Features:**
- **Auto-incrementing Spec IDs** - Automatically creates numbered directories (001, 002, etc.)
- **Resume Capability** - Can resume work on existing specifications by ID
- **Pattern Documentation** - Automatically documents discovered patterns in `docs/patterns/`
- **Interface Documentation** - Captures external service contracts in `docs/interfaces/`
- **Domain Documentation** - Records business rules in `docs/domain/`
- **Confidence Scoring** - Provides implementation readiness assessment with risk analysis
- **Quality Gates** - Requires user approval between major phases

<details>
<summary><strong>View Details</strong></summary>

**What you get:** Three comprehensive documents in `docs/specs/[NNN]-[name]/`:

- **product-requirements.md** - User stories, feature specifications, success criteria, non-functional requirements
- **solution-design.md** - Technical architecture, system components, data models, technology stack, security and performance considerations
- **implementation-plan.md** - Phased task breakdown, dependencies, acceptance criteria, risk assessment

```mermaid
flowchart TD
    A([Your Feature Idea]) --> |initialize| B{Check<br>Existing}
    B --> |exists| C[Review and Refine]
    C --> END[🚀 Ready for /start:implement 001]
    B --> |new| D[📄 **Requirements Gathering**<br/>Create *product-requirements.md* if needed]
    D --> E[📄 **Technical Research**<br/>Create *solution-design.md* if needed, document patterns, interfaces]
    E --> F[📄 **Implementation Planning**<br/>Create *implementation-plan.md*]
    F --> END
```

</details>

---

### `/start:implement <spec-id>`

Execute implementation plans phase-by-phase with parallel specialist agents and continuous validation.

**Purpose:** Transform validated specifications into working code with quality gates and progress tracking

**Usage:**
```bash
/start:implement 001
/start:implement path/to/custom/implementation-plan.md
```

**Key Features:**
- **Parallel Execution** - Multiple agents work simultaneously within phases
- **Sequential Phases** - Phases execute in order with validation gates
- **Rollback on Failure** - Automatic reversion if tests fail
- **Specification Compliance** - Continuous validation against product-requirements.md/solution-design.md
- **Pattern Recognition** - Documents implementation patterns discovered
- **Real-time Updates** - TodoWrite tracking shows live progress
- **Custom Plans** - Can implement any implementation-plan.md file, not just specs

<details>
<summary><strong>View Details</strong></summary>

Loads implementation-plan.md and executes phase-by-phase with approval gates between phases. Multiple specialist agents work in parallel within each phase when tasks are independent. All changes are validated against acceptance criteria and tests run after each task.

```mermaid
flowchart TD
    A([📄 *implementation-plan.md*]) --> |load| B[**Initialize Plan**<br/>Parse phases & tasks]
    B --> |approve| C{Phases<br>Remaining?}
    C --> |yes| D[**Execute Phase N**<br/>⚡ *Parallel agent execution*<br/>✓ *Run tests after each task*]
    D --> |validate| E[**Phase Review**<br/>Check test results<br/>Review changes]
    E --> |continue| C
    C --> |no| F[**Final Validation**<br/>Run full test suite<br/>Verify all requirements]
    F --> END[✅ **Implementation Complete**]
```

</details>

---

### `/start:validate <target>`

Validate specifications, implementations, or understanding through intelligent context detection and the 3 Cs framework (Completeness, Consistency, Correctness).

**Purpose:** Quality gate that works at any lifecycle stage - during specification, before implementation, or after completion

**Usage:**
```bash
/start:validate 001                                          # Validate spec by ID
/start:validate docs/specs/001/solution-design.md            # Validate specific file
/start:validate Check the auth implementation against SDD    # Compare implementation to spec
/start:validate Is my caching approach correct?              # Validate understanding
```

**Key Features:**
- **Intelligent Mode Detection** - Automatically determines validation type from input
- **The 3 Cs Framework** - Checks Completeness, Consistency, and Correctness
- **Ambiguity Detection** - Scans for vague language ("should", "various", "etc.")
- **Cross-Document Traceability** - Verifies PRD→SDD→PLAN alignment
- **Advisory Only** - Provides recommendations without blocking
- **Comparison Validation** - Compares implementations against specifications
- **Understanding Validation** - Confirms correctness of approach or design

<details>
<summary><strong>View Details</strong></summary>

**Four validation modes** automatically detected from input:

| Input Type | Mode | What Gets Validated |
|------------|------|---------------------|
| Spec ID (`005`) | Specification | Full spec quality and readiness |
| File path (`src/auth.ts`) | File | Individual file quality |
| "Check X against Y" | Comparison | Implementation vs specification |
| Freeform text | Understanding | Approach correctness |

**The 3 Cs Framework:**

1. **Completeness** - All sections filled, no `[NEEDS CLARIFICATION]` markers, checklists complete
2. **Consistency** - Cross-document traceability, terminology alignment, no contradictions
3. **Correctness** - ADRs confirmed, dependencies valid, acceptance criteria testable

```mermaid
flowchart TD
    A([Validation Request]) --> |parse| B{Detect Mode}
    B --> |spec ID| C[**Specification Validation**<br/>3 Cs + Ambiguity + Readiness]
    B --> |file path| D[**File Validation**<br/>Quality + Completeness]
    B --> |"against"| E[**Comparison Validation**<br/>Implementation vs Spec]
    B --> |freeform| F[**Understanding Validation**<br/>Approach Correctness]
    C --> G[**Report**<br/>📊 Findings + 💡 Recommendations]
    D --> G
    E --> G
    F --> G
```

</details>

---

### `/start:review [target]`

Multi-agent code review with security, performance, quality, and test coverage specialists running in parallel.

**Purpose:** Comprehensive code review before merging, with specialized agents examining different concerns simultaneously

**Usage:**
```bash
/start:review                                    # Review current PR/staged changes
/start:review --pr 123                           # Review specific PR
/start:review --branch feature/auth              # Review branch changes
/start:review src/auth/ src/users/               # Review specific files/directories
```

**Key Features:**
- **4 Parallel Specialists** - Security, Performance, Quality, and Test agents review simultaneously
- **Target Auto-Detection** - Automatically detects PR, staged changes, or branch diffs
- **Confidence Scoring** - Each finding includes confidence level (HIGH/MEDIUM/LOW)
- **PR Integration** - Posts comments directly to GitHub PRs via `gh` CLI
- **Severity Classification** - CRITICAL, HIGH, MEDIUM, LOW findings

<details>
<summary><strong>View Details</strong></summary>

**Review Agents:**

| Agent | Focus Areas |
|-------|-------------|
| **Security** | SQL injection, XSS, hardcoded secrets, auth bypasses, input validation |
| **Performance** | N+1 queries, missing indexes, memory leaks, inefficient algorithms |
| **Quality** | Code complexity, naming, SOLID principles, error handling, duplication |
| **Tests** | Coverage gaps, missing edge cases, test quality, assertion completeness |

**Output Modes:**

- **PR Mode**: Posts inline comments to GitHub PR
- **Local Mode**: Generates detailed findings report
- **Both include**: Confidence scores, code locations, suggested fixes

```mermaid
flowchart TD
    A([Review Request]) --> |detect| B{Detect Target}
    B --> |PR| C[Load PR Diff]
    B --> |staged| D[Load Staged Changes]
    B --> |branch| E[Load Branch Diff]
    B --> |files| F[Load File Contents]
    C --> G[**Launch 4 Parallel Agents**<br/>🔒 Security<br/>⚡ Performance<br/>✨ Quality<br/>🧪 Tests]
    D --> G
    E --> G
    F --> G
    G --> |merge| H[**Consolidate Findings**<br/>Deduplicate<br/>Rank by severity]
    H --> I{Post to PR?}
    I --> |yes| J[Post Comments via gh CLI]
    I --> |no| K[Generate Report]
    J --> END[✅ Review Complete]
    K --> END
```

</details>

---

### `/start:document [target]`

Generate and sync documentation including API docs, READMEs, JSDoc comments, and documentation audits.

**Purpose:** Keep documentation current with code, generate missing docs, and identify staleness

**Usage:**
```bash
/start:document src/api/                         # Generate API documentation
/start:document --mode readme                    # Update project README
/start:document --mode code src/utils/           # Add JSDoc to code files
/start:document --mode audit                     # Audit documentation coverage
/start:document --mode module src/auth/          # Document entire module
```

**Key Features:**
- **5 Documentation Modes** - Code, API, README, Audit, Module
- **Staleness Detection** - Identifies outdated documentation
- **Coverage Metrics** - Reports documentation completeness percentage
- **OpenAPI Generation** - Creates OpenAPI/Swagger specs from API code
- **Multi-Agent Parallel** - Multiple documentation agents work simultaneously

<details>
<summary><strong>View Details</strong></summary>

**Documentation Modes:**

| Mode | Output | Use Case |
|------|--------|----------|
| `code` | JSDoc/TSDoc comments | Adding inline documentation |
| `api` | OpenAPI spec, endpoint docs | API documentation |
| `readme` | README.md updates | Project documentation |
| `audit` | Coverage report | Finding documentation gaps |
| `module` | Complete module docs | Full module documentation |

**Staleness Detection:**

The skill automatically detects when documentation is outdated by:
- Comparing doc timestamps to code changes
- Checking if documented APIs still match implementation
- Identifying undocumented new exports

```mermaid
flowchart TD
    A([Document Request]) --> |detect| B{Detect Mode}
    B --> |code| C[**JSDoc Generation**<br/>Functions, types, exports]
    B --> |api| D[**API Documentation**<br/>OpenAPI, endpoints, schemas]
    B --> |readme| E[**README Update**<br/>Features, usage, examples]
    B --> |audit| F[**Coverage Audit**<br/>Find gaps, staleness]
    B --> |module| G[**Full Module Docs**<br/>All of the above]
    C --> H[**Generate Documentation**<br/>⚡ Parallel agents when possible]
    D --> H
    E --> H
    F --> H
    G --> H
    H --> I[**Sync & Report**<br/>📊 Coverage metrics<br/>⚠️ Staleness warnings]
    I --> END[✅ Documentation Complete]
```

</details>

---

### `/start:analyze <area>`

Discover and document business rules, technical patterns, and system interfaces through iterative exploration.

**Purpose:** Extract organizational knowledge from existing codebase and create reusable documentation

**Usage:**
```bash
/start:analyze security patterns in authentication
/start:analyze business rules for user permissions
/start:analyze technical patterns in our microservices architecture
```

<details>
<summary><strong>View Details</strong></summary>

Uses cyclical discovery-documentation-review workflow to extract organizational knowledge. Specialist agents explore the codebase to identify patterns, rules, and interfaces across business, technical, security, performance, integration, data, testing, and deployment areas. Documentation is automatically organized into `docs/domain/`, `docs/patterns/`, and `docs/interfaces/` directories.

```mermaid
flowchart TD
    A([Analysis Request]) --> |initialize| B[**Scope Definition**<br/>Clarify analysis area<br/>Set cycle plan]
    B --> |start cycle| C[**Discovery Phase**<br/>⚡ *Specialist analysis*<br/>🔍 *Pattern identification*]
    C --> |document| D[**Documentation Phase**<br/>📄 *Create domain docs*<br/>📄 *Create pattern docs*<br/>📄 *Create interface docs*]
    D --> |review| E[**Review & Validation**<br/>Check completeness<br/>Identify gaps]
    E --> |continue?| F{More Cycles<br>Needed?}
    F --> |yes| C
    F --> |no| G[**Final Summary**<br/>📊 *Analysis report*<br/>🎯 *Recommendations*<br/>📋 *Next steps*]
    G --> END[✅ **Analysis Complete**]
```

</details>

---

### `/start:refactor <description>`

Improve code quality while strictly preserving all existing behavior through test-validated incremental changes.

**Purpose:** Safe, systematic refactoring with automatic rollback on test failures

**Usage:**
```bash
/start:refactor Simplify the authentication middleware for better testability
/start:refactor Improve the WebSocket connection manager
```

<details>
<summary><strong>View Details</strong></summary>

Strictly preserves behavior through test-validated incremental changes. All tests must pass before refactoring begins and after each change. Automatic rollback on test failures. For simple refactorings, applies changes directly with continuous validation. For complex refactorings, creates specification documents and defers to `/start:implement` for planned execution.

```mermaid
flowchart TD
    A([Refactoring Request]) --> |analyze| B[**Goal Clarification**<br/>Define objectives<br/>Analyze codebase]
    B --> |assess| C{**Complexity<br>Check**}
    C --> |simple| D[**Direct Refactoring**<br/>✓ *Run tests first*<br/>🔧 *Apply changes*<br/>✓ *Validate each step*]
    D --> |review| E[**Specialist Review**<br/>Code quality check<br/>Performance impact]
    E --> DONE[✅ **Refactoring Complete**]
    C --> |complex| F[**Create Specification**<br/>📄 *Generate solution-design.md*<br/>📄 *Generate implementation-plan.md*<br/>Document approach]
    F --> |defer| G[🚀 **Ready for /start:implement**<br/>Execute via planned phases]
```

</details>

---

### `/start:debug <description>`

Diagnose and resolve bugs through conversational investigation with systematic root cause analysis.

**Purpose:** Natural language debugging partner that helps identify and fix issues through dialogue, not rigid procedures

**Usage:**
```bash
/start:debug The API returns 500 errors when uploading large files
/start:debug Tests are failing intermittently on CI but pass locally
/start:debug Users report slow page loads after the latest deployment
```

**Key Features:**
- **Conversational Flow** - Natural dialogue, not rigid checklists or procedures
- **Progressive Disclosure** - Starts with summary, reveals details on request
- **Observable Actions Only** - Reports only what was actually checked and found
- **User-Driven** - Proposes next steps, lets user guide the direction
- **Hypothesis Tracking** - Forms and tests ranked hypotheses systematically
- **Evidence-Based** - Never fabricates reasoning; all conclusions backed by evidence

<details>
<summary><strong>View Details</strong></summary>

Uses a conversational approach through five natural phases: understand the problem, narrow it down, find the root cause, fix and verify, wrap up. The debugger reports only observable actions ("I checked X and found Y") and never fabricates reasoning. Users can ask "what did you check?" at any point and receive honest, verifiable answers.

```mermaid
flowchart TD
    A([Bug Description]) --> |understand| B[**Phase 1: Understand**<br/>Reproduce issue<br/>Gather context]
    B --> |isolate| C[**Phase 2: Narrow Down**<br/>Form hypotheses<br/>Binary search]
    C --> |investigate| D[**Phase 3: Root Cause**<br/>Test hypotheses<br/>Find evidence]
    D --> |fix| E[**Phase 4: Fix & Verify**<br/>Propose fix<br/>Run tests]
    E --> |close| F[**Phase 5: Wrap Up**<br/>Summarize if needed<br/>Suggest follow-ups]
    F --> END[✅ **Bug Resolved**]

    D --> |stuck| G{Need More<br/>Context?}
    G --> |yes| B
    G --> |different angle| C
```

**The Four Commandments:**
1. **Conversational, not procedural** - It's a dialogue, not a checklist
2. **Observable only** - "I looked at X and found Y" not "This is probably..."
3. **Progressive disclosure** - Start brief, expand on request
4. **User in control** - "Want me to...?" not "I will now..."

</details>

---

### `/start:constitution [focus-areas]`

Create or update a project constitution with governance rules through discovery-based pattern analysis.

**Purpose:** Establish checkable project rules that are enforced during implementation, review, and validation

**Usage:**
```bash
/start:constitution                                    # Create new constitution via codebase discovery
/start:constitution "security and testing"             # Focus on specific areas
/start:constitution "Add API patterns"                 # Update existing constitution
```

**Key Features:**
- **Discovery-Based Rules** - Analyzes actual codebase patterns, never assumes frameworks
- **L1/L2/L3 Level System** - L1 (blocking + autofix), L2 (blocking, manual), L3 (advisory)
- **Three-Layer Enforcement** - Checked during specify (SDD), implement, and review
- **Pattern + Check Rules** - Supports regex patterns and semantic LLM-interpreted checks
- **Graceful Degradation** - System works normally if no constitution exists

<details>
<summary><strong>View Details</strong></summary>

Creates `CONSTITUTION.md` at project root (like README, LICENSE, CODE_OF_CONDUCT). The constitution defines checkable guardrails that detect violations during development.

**Key Distinction:**
- **CLAUDE.md** = Project description, AI guidance ("Use React with TypeScript")
- **CONSTITUTION.md** = Checkable rules that catch violations ("No barrel exports")

**Level Definitions:**

| Level | Name | Blocking | Autofix | Use Case |
|-------|------|----------|---------|----------|
| **L1** | Must | ✅ Yes | ✅ AI auto-corrects | Security, correctness, critical architecture |
| **L2** | Should | ✅ Yes | ❌ No | Important rules requiring human judgment |
| **L3** | May | ❌ No | ❌ No | Style preferences, suggestions |

**Rule Format Example:**

```markdown
### No Hardcoded Secrets

\```yaml
level: L1
pattern: "(api_key|secret|password)\\s*[:=]\\s*['\"][^'\"]{8,}['\"]"
scope: "**/*.{ts,js}"
exclude: "**/*.test.*, .env.example"
message: Hardcoded secret detected. Use environment variables.
\```

Secrets must never be committed to source control.
```

**Three-Layer Enforcement:**

| Phase | Command | Enforcement |
|-------|---------|-------------|
| **Planning** | `/start:specify` (SDD) | SDD must not violate constitutional principles |
| **Task** | `/start:implement` | Task ordering respects constitutional priorities |
| **Implementation** | `/start:implement` | Generated code checked; L1/L2 violations block completion |

```mermaid
flowchart TD
    A([/start:constitution]) --> |check| B{Constitution<br>Exists?}
    B --> |no| C[**Discovery Phase**<br/>Explore codebase patterns]
    C --> D[**Rule Generation**<br/>Create L1/L2/L3 rules]
    D --> E[**User Confirmation**<br/>Present proposed rules]
    E --> F[**Write CONSTITUTION.md**<br/>At project root]
    B --> |yes| G{Update or<br>Validate?}
    G --> |update| H[Add new rules<br/>Focus on specified areas]
    H --> E
    G --> |validate| I[Run /start:validate constitution]
    F --> END[✅ Constitution Created]
```

</details>

---

### Installation

Install The Agentic Startup framework using the one-line installer:

```bash
curl -fsSL https://raw.githubusercontent.com/rsmdt/the-startup/main/install.sh | sh
```

**What it does:**
- Adds the `rsmdt/the-startup` marketplace
- Installs `team@the-startup` and `start@the-startup` plugins
- Configures `start:The Startup` as the default output style
- Optionally installs the git-aware statusline

**Flags:**
- `--yes` - Skip all confirmation prompts
- `--no-statusline` - Skip statusline installation
- `--help` - Show usage information

**Note:** Output styles are available immediately via `/output-style` - no additional setup required.

---

## Autonomous Skills

The `start` plugin includes five autonomous skills that activate automatically based on context. You never need to explicitly invoke them — they work when needed.

### Specification Skills

| Skill | Purpose |
|-------|---------|
| `specify-meta` | Spec directory creation, README tracking, phase transitions |
| `specify-requirements` | PRD template, validation, requirements gathering |
| `specify-solution` | SDD template, architecture design, ADR management |
| `specify-plan` | PLAN template, task sequencing, dependency mapping |

### Methodology Skills

| Skill | Purpose |
|-------|---------|
| `writing-skills` | Skill authoring, auditing, and verification methodology |

---

## 🏗️ Documentation Structure

The plugin encourages structured knowledge management:

```
docs/
├── specs/
│   └── [3-digit-number]-[feature-name]/
│       ├── product-requirements.md         # What to build
│       ├── solution-design.md              # How to build it
│       └── implementation-plan.md          # Implementation tasks
│
├── domain/                                  # Business rules
│   ├── user-permissions.md
│   ├── order-workflow.md
│   └── pricing-rules.md
│
├── patterns/                                # Technical patterns
│   ├── authentication-flow.md
│   ├── caching-strategy.md
│   └── error-handling.md
│
└── interfaces/                              # External integrations
    ├── stripe-payments.md
    ├── sendgrid-webhooks.md
    └── oauth-providers.md
```

### Auto-Documentation

The `document` skill's Capture perspective automatically creates files in the correct location when patterns, interfaces, or domain rules are discovered during:
- Specification creation (`/start:specify`)
- Implementation (`/start:implement`)
- Analysis (`/start:analyze`)

### Deduplication

The capture workflow always checks existing documentation before creating new files, preventing duplicates.

---

## 🎨 Output Styles

The `start` plugin includes two output styles for different working preferences. Both share the same operational fundamentals (verification, code ownership, scope management) but express them differently.

**Activate via:** `/output-style start:The Startup` or `/output-style start:The ScaleUp`

---

### The Startup 🚀

**High-energy execution with structured momentum.**

| Aspect | Description |
|--------|-------------|
| **Vibe** | Demo day energy, Y Combinator intensity |
| **Voice** | "Let's deliver this NOW!", "BOOM! That's what I'm talking about!" |
| **Mantra** | "Done is better than perfect, but quality is non-negotiable" |

**Personality:**
- **The Visionary Leader** - "We'll figure it out" - execute fast, iterate faster
- **The Rally Captain** - Turn challenges into team victories
- **The Orchestrator** - Run parallel execution like a conductor
- **The Pragmatist** - MVP today beats perfect next quarter

**Best for:**
- Fast-paced development sprints
- High-energy execution mode
- When you want momentum and celebration

---

### The ScaleUp 📈

**Calm confidence with educational depth.**

| Aspect | Description |
|--------|-------------|
| **Vibe** | Professional craft, engineering excellence |
| **Voice** | "We've solved harder problems. Here's the approach.", "This decision matters because..." |
| **Mantra** | "Sustainable speed at scale. We move fast, but we don't break things" |

**Personality:**
- **The Seasoned Leader** - We've been through the fire. Now we build to last.
- **The Strategist** - Think two steps ahead. Today's shortcut is tomorrow's outage.
- **The Multiplier** - Your job is to make the whole team better, not just ship code.
- **The Guardian** - Reliability isn't optional. Customers trust us with their business.

**Unique feature - Educational Insights:**

The ScaleUp provides contextual explanations as it works:

> I've added the retry logic to the API client:
> ```typescript
> await retry(fetchUser, { maxAttempts: 3, backoff: 'exponential' });
> ```
> 💡 *Insight: I used exponential backoff here because this endpoint has rate limiting. The existing `src/utils/retry.ts` helper already implements this pattern - I'm reusing it rather than adding a new dependency.*

**Best for:**
- Learning while building
- Understanding codebase patterns
- When you want explanations with your code
- Onboarding to unfamiliar codebases

---

### Comparison

| Dimension | The Startup | The ScaleUp |
|-----------|-------------|-------------|
| **Energy** | High-octane, celebratory | Calm, measured |
| **Explanations** | Minimal - ships fast | Educational insights included |
| **Failures** | "That didn't work. Moving on." | "Here's what failed and why..." |
| **Closing thought** | "What did we deliver?" | "Can the team maintain this without me?" |

---

## Typical Development Workflow

### Setup (Optional)

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                             PROJECT SETUP                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│   /start:constitution ────► Create project governance rules                 │
│        │                    L1/L2/L3 rules auto-enforced in BUILD flow      │
│        │                    CONSTITUTION.md at project root                  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Primary Workflow: Specify → Validate → Implement → Review

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          PRIMARY DEVELOPMENT FLOW                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│   /start:specify ──► /start:validate ──► /start:implement ──► /start:review │
│        │                   │                    │                   │        │
│   Create specs      Check quality        Execute plan        Code review    │
│   PRD + SDD + PLAN  3 Cs framework      Phase-by-phase     Security + Perf │
│        │                   │                    │                   │        │
│   ↳ Constitution     ↳ Constitution      ↳ Constitution      ↳ Constitution │
│     checked on SDD     mode available      + drift enforced    compliance   │
│                                                                              │
│   Optional: /start:document after implementation for documentation sync     │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

*If `CONSTITUTION.md` exists, rules are automatically checked at each stage.*

**1. Create Specification**

```bash
/start:specify Add real-time notification system with WebSocket support
```

**What happens:**
- Creates `docs/specs/001-notification-system/`
- Generates product-requirements.md, solution-design.md, implementation-plan.md
- Documents discovered patterns/interfaces
- Optional: Creates `spec/001-notification-system` git branch

**2. Validate Before Implementation (Recommended)**

```bash
/start:validate 001
```

**What happens:**
- Checks completeness, consistency, correctness (3 Cs)
- Detects ambiguities and vague language
- Security scanning for common vulnerabilities
- Verifies cross-document traceability

**3. Execute Implementation**

```bash
/start:implement 001
```

**What happens:**
- Optional: Creates `feature/001-notification-system` git branch
- Executes phases sequentially with user approval
- Parallel agent coordination within phases
- Continuous test validation
- Optional: Creates PR at completion

**4. Review Before Merge**

```bash
/start:review
```

**What happens:**
- 4 parallel specialists review (Security, Performance, Quality, Tests)
- Posts findings to PR if applicable
- Generates consolidated report with severity rankings

**5. Generate Documentation (Optional)**

```bash
/start:document src/notifications/
```

**What happens:**
- Adds JSDoc/TSDoc comments
- Updates README if needed
- Reports documentation coverage

---

### Maintenance Workflows

**Understand Existing Code**

```bash
/start:analyze security patterns in authentication
```

Documents findings in `docs/patterns/`, `docs/domain/`, `docs/interfaces/`

**Refactoring**

```bash
/start:refactor Restructure the authentication module for better testability
```

For architectural changes - creates specs, plans migration, handles breaking changes.

**Fix Bugs**

```bash
/start:debug The notification system stops working after 100 concurrent users
```

Conversational investigation with observable evidence and user-driven direction.

**Audit Documentation**

```bash
/start:document --mode audit
```

Reports documentation coverage and identifies stale or missing docs.

---

## Skills in Action

### Example: Knowledge Capture via Document Skill

**Scenario:** During implementation, an agent discovers a pattern

```
Agent output: "I implemented a retry mechanism with exponential backoff for API calls"
```

**What happens automatically:**
1. Document skill's Capture perspective activates
2. Checks `docs/patterns/` for existing retry patterns
3. Not found → Creates `docs/patterns/api-retry-strategy.md`
4. Uses pattern template
5. Reports: "📝 Created docs/patterns/api-retry-strategy.md"

**You didn't have to:** Manually request documentation or specify the path

---

---

## Templates

Rich templates for structured documentation, co-located with their skills:

```
plugins/start/skills/
├── specify-requirements/template.md   # Product requirements structure
├── specify-solution/template.md       # Solution design structure
├── specify-plan/template.md           # Implementation plan structure
└── document/templates/                # Knowledge capture templates
    ├── domain-template.md             # Business rules
    ├── pattern-template.md            # Technical patterns
    └── interface-template.md          # External integrations
```

**Usage:** Automatically used by `/start:specify` and `/start:document` when creating documentation

---

## Philosophy

### Spec-Driven Development

**The problem:** Features built without clear requirements lead to technical debt, rework, and inconsistency.

**Our approach:**

1. **Specify First** - Comprehensive specifications before code
2. **Review & Refine** - Validate with stakeholders
3. **Implement with Confidence** - Execute validated plans
4. **Document & Learn** - Capture patterns for reuse

### Core Principles

- **Measure twice, cut once** - Upfront planning saves time
- **Documentation as code** - Specs evolve with your codebase
- **Parallel execution** - Multiple specialists, clear boundaries
- **Quality gates** - DOR/DOD enforce standards
- **Progressive disclosure** - Load details only when needed

---

## Further Reading

- **[Main README](../../README.md)** - Project overview and installation
- **[Claude Code Documentation](https://docs.claude.com/claude-code)** - Official Claude Code docs
- **[Claude Code Skills Guide](https://docs.claude.com/claude-code/skills)** - How to create skills
