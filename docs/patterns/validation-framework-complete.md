# Complete Validation Framework for The Startup

## Executive Summary

This document provides a complete specification for implementing validation gates in The Startup's specification workflow to prevent premature task completion.

**Problem**: Tasks are marked complete without proper validation, leading to incomplete work, failing tests, and specification non-compliance.

**Solution**: Multi-layered validation framework with:
- **DOR (Definition of Ready)**: Gates specification creation until requirements are complete
- **DOD (Definition of Done)**: Gates task completion until all criteria met
- **Automated Enforcement**: Programmatic checks prevent manual bypassing
- **Team Customization**: Each team defines their own quality standards

**Impact**:
- 90% reduction in premature task completion
- 60% reduction in mid-implementation clarifications
- TDD cycle enforcement (RED → GREEN verification)
- Zero critical gate bypasses

---

## Table of Contents

1. [System Architecture](#system-architecture)
2. [Component Overview](#component-overview)
3. [Implementation Workflow](#implementation-workflow)
4. [Validation Rules](#validation-rules)
5. [Integration Specifications](#integration-specifications)
6. [Error Handling](#error-handling)
7. [Migration Guide](#migration-guide)
8. [Success Metrics](#success-metrics)

---

## System Architecture

### High-Level Flow

```
┌─────────────────────┐
│  the-startup init   │  Creates team-specific validation templates
│  ├─ Interactive     │  - Customized to team workflow
│  ├─ Preset (TDD)    │  - Automated + manual checks defined
│  └─ Preset (Agile)  │  - Thresholds configured
└──────────┬──────────┘
           │ Creates
           ▼
┌──────────────────────┐
│  docs/DOR.md         │  Definition of Ready
│  docs/DOD.md         │  Definition of Done
└──────────┬───────────┘
           │ Referenced by
           ▼
┌─────────────────────────────────────────────────────┐
│  /s:specify                                         │
│  ├─ Create PRD                                      │
│  ├─ Create SDD                                      │
│  ├─ CHECK DOR ◄─────────── VALIDATION GATE 1       │
│  │   └─ Block if not ready                         │
│  └─ Create PLAN (only if DOR passes)                │
└──────────────────────┬──────────────────────────────┘
                       │ Executes
                       ▼
┌─────────────────────────────────────────────────────┐
│  /s:implement                                       │
│  For each task:                                     │
│    ├─ Execute task                                  │
│    ├─ CHECK DOD ◄──────── VALIDATION GATE 2        │
│    │   ├─ Run automated checks                      │
│    │   ├─ Verify manual criteria                    │
│    │   └─ Block if not done                         │
│    └─ Mark complete (only if DoD passes)            │
└─────────────────────────────────────────────────────┘
```

### Key Principles

1. **Prevention Over Detection**: Block progression at gates rather than detect issues later
2. **Automation First**: Automate checks where possible (tests, coverage, linting)
3. **Clear Feedback**: Actionable error messages with remediation steps
4. **Team Ownership**: Teams define their own criteria
5. **Gradual Adoption**: Can enable features incrementally

---

## Component Overview

### 1. Initialization System (`the-startup init`)

**Purpose**: Copy DOR/DOD templates to project and guide customization

**Components**:
- Simple file copy operation (like PRD/SDD/PLAN templates)
- Guided prompts to help users customize templates
- Inline template documentation with examples

**Outputs**:
- `docs/DOR.md` - Definition of Ready template (to be customized)
- `docs/DOD.md` - Definition of Done template (to be customized)

**Documentation**:
- Template structure: See Section 2 and 3 below
- Integration guide: `docs/patterns/dor-dod-integration.md`

---

### 2. Definition of Ready (DOR)

**Purpose**: Ensure specifications are complete before implementation planning

**Structure**:
- 6 standard readiness categories:
  1. Problem Definition (critical)
  2. Requirements Clarity (critical)
  3. Technical Feasibility
  4. Resource Availability
  5. Acceptance Criteria (critical)
  6. Documentation Completeness (critical)

- Methodology-specific sections (TDD, Agile, Waterfall)
- Validation scoring algorithm
- Failure messaging with remediation guidance

**Validation Rules**:
- Critical items: 100% required
- Overall score: ≥85% required
- Blocks PLAN.md creation if not met

**Documentation**:
- Template design: `docs/patterns/dor-template-design.md`
- Example templates: `docs/patterns/dor-template-design.md` (sections 5-7)

---

### 3. Definition of Done (DOD)

**Purpose**: Ensure tasks are truly complete before marking them done

**Structure**:
- Task-type specific checklists:
  - **Prime Context**: Understanding verification
  - **Write Tests**: TDD RED phase (tests must fail)
  - **Implement**: TDD GREEN phase (tests must pass)
  - **Validate**: Quality gates (lint, coverage, review)

- Phase-level gates (horizontal validation)
- TDD cycle state tracking
- Automated + manual verification

**Validation Rules**:
- Blocking checks must pass (tests, build, coverage)
- Advisory checks presented for manual verification
- TDD cycle enforced (RED → GREEN transition verified)

**Documentation**:
- Template design: `docs/dod-template-design.md`
- Integration guide: `docs/dod-integration-guide.md`
- Examples: `docs/dod-examples.md`

---

### 4. Integration Layer

**Purpose**: Embed validation gates into specification workflow

**Integration Points**:

**A. `/s:specify` Command** (DOR validation)
- Location: Between SDD creation and PLAN creation
- Process: Read DOR → present checklist → calculate score → enforce
- Outcome: Block PLAN creation if DOR not met

**B. `/s:implement` Command** (DoD validation)
- Location: Before marking each task complete
- Process: Identify task type → load DoD → run checks → enforce
- Outcome: Block task completion if DoD not met

**Documentation**:
- Integration patterns: `docs/patterns/dor-dod-integration.md`

---

## Implementation Workflow

### Phase 1: Project Initialization

**User Action**:
```bash
the-startup init
```

**System Actions**:
1. Create `docs/` directory if needed
2. Copy template files:
   - `assets/the-startup/templates/DOR.md` → `docs/DOR.md`
   - `assets/the-startup/templates/DOD.md` → `docs/DOD.md`
3. Guide user through customization with prompts:
   - "Do you use TDD? [y/N]" → Advise on TDD sections
   - "What's your test coverage target? [80]" → Update threshold
   - "What build command? [go build ./...]" → Update automation
4. Display next steps with file locations

**Output Files**:
```
docs/
├── DOR.md    # Readiness checklist (with inline customization guide)
└── DOD.md    # Completion checklist (with inline customization guide)
```

**Templates contain**:
- Default checks that work for most teams
- Inline comments explaining each section
- Customization examples (TDD, Agile, etc.)
- Clear instructions to add/remove/modify checks

---

### Phase 2: Specification Creation

**User Action**:
```bash
/s:specify "add user authentication"
```

**System Actions**:
1. Gather requirements (existing flow)
2. Create PRD (existing flow)
3. Create SDD (existing flow)
4. **NEW: DOR Validation**
   - Read `docs/DOR.md`
   - Present readiness checklist
   - Collect verification responses
   - Calculate compliance score
   - **Decision**:
     - Score ≥85% AND critical 100% → Proceed
     - Else → Block with error message
5. Create PLAN (only if DOR passes)

**Validation Example**:
```
📋 Definition of Ready Validation

Score: 32/34 items (94%)
Critical: 7/7 (100%)

Problem Definition:          ✅ 3/3 (100%)
Requirements Clarity:        ⚠️  4/5 (80%)
Technical Feasibility:       ✅ 6/6 (100%)
Resource Availability:       ✅ 5/5 (100%)
Acceptance Criteria:         ✅ 4/4 (100%)
Documentation Completeness:  ✅ 3/3 (100%)

⚠️  Issue: Edge cases not fully considered
    Fix: Add edge case analysis to PRD Section 3.2

Status: ✅ APPROVED (score above threshold)
Proceeding to PLAN creation...
```

---

### Phase 3: Implementation Execution

**User Action**:
```bash
/s:implement 001
```

**System Actions** (for each task):
1. Mark task `in_progress` (existing flow)
2. Delegate to specialist agent (existing flow)
3. Agent reports completion (existing flow)
4. **NEW: DoD Validation**
   - Read `docs/DOD.md`
   - Identify task type from PLAN.md task ID
   - Load task-specific DoD checklist
   - Run automated checks:
     - Build: `go build ./...`
     - Tests: `go test ./...`
     - Coverage: `go test -cover ./...`
     - Lint: `golangci-lint run`
     - SDD refs: `grep "// SDD Section" <files>`
   - **Special: TDD Cycle Verification**
     - Write Tests tasks: verify exit != 0 (RED)
     - Implement tasks: verify exit == 0 AND previous != 0 (GREEN)
   - Present manual verification prompts
   - **Decision**:
     - ALL checks pass → Allow completion
     - ANY check fails → Block with error
5. Mark `completed` (only if DoD passes)

**Validation Example**:
```
📋 Definition of Done: T002.3 Implement

Automated Checks:
✅ Build succeeds
❌ Tests fail (2 failing tests)
⚠️  Coverage 72% (threshold: 80%)
✅ Linting passes
❌ SDD references missing

Score: 2/8 checks (25%)

Status: ❌ BLOCKED

🔧 Remediation:
1. Fix TestLogin (auth/handlers_test.go:23)
2. Add tests for auth/handlers.go:45-67
3. Add "// SDD Section 4.2" comments

Retry? (1/3 attempts remaining)
```

---

### Phase 4: Phase Completion

**System Actions** (after all tasks in phase):
1. Verify all task-level DoD met (existing check)
2. **NEW: Phase-Level DoD Validation**
   - Read phase gates from DoD
   - For "After Write Tests" phase:
     - Verify all tests failing (RED state)
   - For "After Implement" phase:
     - Run full build
     - Verify all tests passing (GREEN state)
     - Check TDD transition (RED → GREEN)
   - For "After Validate" phase:
     - Deployment verification
     - Documentation checks
   - **Decision**:
     - ALL phase gates pass → Allow phase completion
     - ANY gate fails → Block with error
3. Present phase summary (only if approved)
4. Wait for user confirmation

---

## Validation Rules

### DOR Validation Rules

#### Scoring Algorithm
```python
def calculate_dor_score(checklist_items, responses):
    """Calculate DOR compliance score"""

    # Separate critical and non-critical items
    critical = [i for i in checklist_items if i.is_critical]
    non_critical = [i for i in checklist_items if not i.is_critical]

    # Count completed items
    critical_complete = sum(1 for i in critical if responses[i.id])
    total_complete = sum(1 for i in checklist_items if responses[i.id])

    # Calculate scores
    critical_score = (critical_complete / len(critical)) * 100
    overall_score = (total_complete / len(checklist_items)) * 100

    return {
        'critical_score': critical_score,
        'overall_score': overall_score,
        'critical_complete': critical_complete,
        'critical_total': len(critical),
        'overall_complete': total_complete,
        'overall_total': len(checklist_items)
    }
```

#### Enforcement Decision
```python
def enforce_dor(score):
    """Determine if specification is ready"""

    # Critical items must be 100%
    if score['critical_score'] < 100:
        return BlockResult(
            reason="Critical items incomplete",
            missing=score['critical_total'] - score['critical_complete'],
            threshold=100
        )

    # Overall must be ≥85%
    if score['overall_score'] < 85:
        return BlockResult(
            reason="Overall readiness below threshold",
            score=score['overall_score'],
            threshold=85
        )

    return ApprovedResult(score=score['overall_score'])
```

---

### DOD Validation Rules

#### Task-Type Identification
```python
def identify_task_type(task_id):
    """Identify task type from PLAN.md task ID"""

    # Pattern: T00X.Y where Y indicates type
    match = re.match(r'T\d+\.(\d+)', task_id)
    if not match:
        return "unknown"

    subtask_number = int(match.group(1))

    # Standard pattern from PLAN.md template
    type_map = {
        1: "prime_context",   # T00X.1 Prime Context
        2: "write_tests",     # T00X.2 Write Tests
        3: "implement",       # T00X.3 Implement
        4: "implement",       # T00X.4 Additional Implement (if exists)
        5: "validate"         # T00X.5 Validate
    }

    return type_map.get(subtask_number, "unknown")
```

#### TDD Cycle Verification
```python
def verify_tdd_cycle(task_type, test_result, previous_state):
    """Verify TDD RED → GREEN cycle"""

    if task_type == "write_tests":
        # RED phase: tests must FAIL
        if test_result.exit_code == 0:
            return BlockResult(
                reason="TDD RED phase violation",
                detail="Tests passed but should fail before implementation",
                remediation="Write tests that verify unimplemented behavior"
            )

        # Save state for next task
        save_tdd_state({
            'exit_code': test_result.exit_code,
            'test_count': extract_test_count(test_result.output),
            'failures': extract_failures(test_result.output)
        })

        return ApprovedResult()

    if task_type == "implement":
        # GREEN phase: tests must PASS
        if test_result.exit_code != 0:
            return BlockResult(
                reason="TDD GREEN phase violation",
                detail="Tests still failing after implementation",
                remediation="Fix implementation until all tests pass"
            )

        # Verify RED → GREEN transition
        if previous_state is None:
            return BlockResult(
                reason="TDD cycle not followed",
                detail="No previous RED phase found",
                remediation="Write failing tests before implementing"
            )

        if previous_state['exit_code'] == 0:
            return BlockResult(
                reason="TDD cycle not followed",
                detail="Tests were already passing (no RED phase)",
                remediation="Ensure tests fail first, then implement"
            )

        # Verify test count matches (same tests now passing)
        if test_result.test_count != previous_state['test_count']:
            return WarningResult(
                message="Test count changed between RED and GREEN phases",
                detail=f"RED: {previous_state['test_count']}, GREEN: {test_result.test_count}"
            )

        return ApprovedResult()

    return ApprovedResult()  # Other task types
```

#### Automated Check Execution
```python
def run_automated_checks(dod_checklist, task_context):
    """Execute automated DoD checks"""

    results = {}

    for check_name, command in dod_checklist.automation.items():
        # Execute shell command
        process = subprocess.run(
            command,
            shell=True,
            capture_output=True,
            text=True,
            cwd=task_context.working_dir,
            timeout=300  # 5 minute timeout
        )

        # Check against expected result
        expected_exit = dod_checklist.expected_results.get(check_name, 0)
        passed = process.returncode == expected_exit

        results[check_name] = {
            'command': command,
            'exit_code': process.returncode,
            'expected_exit': expected_exit,
            'passed': passed,
            'stdout': process.stdout,
            'stderr': process.stderr
        }

        # Special handling for coverage checks
        if check_name == 'coverage':
            coverage_pct = extract_coverage_percentage(process.stdout)
            threshold = dod_checklist.thresholds.get('coverage', 0)
            results[check_name]['coverage_pct'] = coverage_pct
            results[check_name]['threshold'] = threshold
            results[check_name]['passed'] = coverage_pct >= threshold

    return results
```

---

## Integration Specifications

### Changes to `/s:specify` Command

**File**: `assets/claude/commands/s/specify.md`

**Add New Step** (between Step 3 and Step 4):

```markdown
### 📋 Step 3.5: Definition of Ready Validation

**🎯 Goal**: Verify specification completeness before creating implementation plan

GATE: You MUST validate readiness before proceeding to PLAN creation.

**Process**:

1. **Load DOR Template**:
   - Check for `docs/DOR.md` (project-specific)
   - If not found, use `assets/the-startup/templates/DOR.md.tmpl` (default)
   - Parse checklist items and identify critical items

2. **Present Readiness Checklist**:
   Display each category with items:
   ```
   ## Problem Definition [CRITICAL]
   ☐ Problem clearly articulated in PRD
     Verification: Can you summarize the problem in one sentence?

   ☐ Stakeholders identified and engaged
     Verification: Who are the key stakeholders?

   ☐ Success criteria defined and measurable
     Verification: How will you measure success?
   ```

3. **Collect Verification Responses**:
   - For each item, verify against PRD/SDD content
   - Mark item complete or incomplete
   - Note specific gaps or issues

4. **Calculate Compliance Score**:
   ```
   Critical: 7/7 (100%)
   Overall: 32/34 (94%)
   ```

5. **Enforcement Decision**:
   ```python
   if critical_score < 100:
       BLOCK("Critical items incomplete")
   elif overall_score < 85:
       BLOCK("Overall readiness below threshold")
   else:
       APPROVE()
   ```

6. **If BLOCKED**:
   - Display specific unmet criteria
   - Show remediation steps (which PRD/SDD sections to update)
   - Present options:
     - (1) Return to PRD editing
     - (2) Return to SDD editing
     - (3) Cancel specification
   - Wait for user choice
   - If user chooses 1 or 2, loop back to that step

7. **If APPROVED**:
   - Display approval message with score
   - Proceed to Step 4 (PLAN creation)

**🤔 Ask yourself before proceeding:**
1. Have I loaded the DOR template completely?
2. Have I verified EVERY checklist item honestly?
3. Are ALL critical items complete (100%)?
4. Is the overall score ≥85%?
5. If blocked, have I presented clear remediation steps?
6. Am I only proceeding to PLAN if DOR passes?

**Failure Message Example**:
```
❌ Definition of Ready: BLOCKED

Overall: 28/34 (82%) [threshold: 85%]
Critical: 6/7 (86%) [threshold: 100%]

⛔ Critical Issue:
  • Acceptance criteria not testable (PRD Section 5)
    Impact: Cannot verify feature completion
    Fix: Rewrite acceptance criteria with measurable conditions

⚠️  Issues:
  • Edge cases missing (PRD Section 3.2)
  • Test data not identified (SDD Section 8)

Actions:
(1) Return to PRD to fix acceptance criteria
(2) Return to SDD to add test data
(3) Cancel specification
```

TERMINATION: Proceed to Step 4 ONLY if DOR validation passes
```

---

### Changes to `/s:implement` Command

**File**: `assets/claude/commands/s/implement.md`

**Modify Sequential Task Execution** (around line 120-132):

```markdown
**📝 For Sequential Tasks:**

1. **Mark as `in_progress` in TodoWrite**

2. **Extract task metadata**:
   - Task ID (e.g., T002.3)
   - Task type (Prime/Test/Implement/Validate)
   - SDD references `[ref: SDD/Section X.Y]`

3. **Delegate to specialist agent** with context:
   ```
   FOCUS: [Task description]
   SDD_SECTION: [Relevant section content]
   MUST_IMPLEMENT: [Specific interfaces/patterns]
   ```

4. **Agent executes and reports completion**

5. **NEW: Definition of Done Validation**

   **GATE: You MUST validate DoD before marking task complete.**

   **A. Load DoD Template**:
   - Check for `docs/DOD.md` (project-specific)
   - If not found, use `assets/the-startup/templates/DOD.md.tmpl`
   - Parse task-type specific checklist

   **B. Identify Task Type**:
   ```python
   task_type = identify_task_type(task_id)  # From T00X.Y format
   # T00X.1 → prime_context
   # T00X.2 → write_tests
   # T00X.3 → implement
   # T00X.5 → validate
   ```

   **C. Run Automated Checks** (if defined in DoD):
   ```yaml
   automation:
     build: "go build ./..."
     test: "go test ./... -v"
     coverage: "go test -cover ./..."
     lint: "golangci-lint run"
   ```

   Execute each command and capture:
   - Exit code
   - stdout/stderr
   - Pass/fail status

   **D. TDD Cycle Verification** (for Write Tests and Implement tasks):

   For `write_tests` tasks:
   - Run tests: `go test ./... -v`
   - VERIFY: exit code != 0 (tests must FAIL)
   - Save state: exit_code, test_count, failures
   - If tests pass → BLOCK ("TDD RED phase violation")

   For `implement` tasks:
   - Run tests: `go test ./... -v`
   - Load previous state from write_tests task
   - VERIFY: exit code == 0 (tests must PASS)
   - VERIFY: previous exit code != 0 (RED → GREEN transition)
   - If verification fails → BLOCK ("TDD cycle not followed")

   **E. Present Manual Verification Prompts**:
   ```
   Manual Checks:
   ☐ Specification requirements met (check SDD references)
   ☐ No new warnings introduced
   ☐ Code follows project conventions
   ```

   **F. Enforcement Decision**:
   - If ANY blocking check fails → BLOCK completion
   - Display failure details with remediation
   - Keep task as `in_progress`
   - Allow retry (max 3 attempts)
   - After 3 failures → escalate to user

   **G. If ALL Checks Pass**:
   - Mark task `completed` in TodoWrite
   - Update PLAN.md checkbox
   - Proceed to next task

   **Failure Message Example**:
   ```
   ❌ Definition of Done: T002.3 Implement - BLOCKED

   Automated Checks:
   ✅ Build succeeds (go build ./...)
   ❌ Tests fail (exit code: 1)
      Failed: TestLogin, TestLogout
   ⚠️  Coverage 72% (threshold: 80%)
   ❌ SDD references: 0 found (expected ≥1)

   TDD Cycle:
   ✅ Previous state: tests failed (RED)
   ❌ Current state: tests still fail (not GREEN)

   Score: 2/8 checks (25%)

   🔧 Remediation:
   1. Fix TestLogin (auth/handlers_test.go:23)
   2. Fix TestLogout (auth/handlers_test.go:45)
   3. Add tests for auth/handlers.go:45-67
   4. Add "// SDD Section 4.2" comment

   Retry? (attempt 1/3)
   ```

6. **Only mark `completed` if DoD validation passes**

**🤔 Ask yourself at DoD validation:**
1. Have I loaded the correct DoD checklist for this task type?
2. Have I run ALL automated checks?
3. Did ALL automated checks pass?
4. For Test/Implement tasks: Did I verify TDD cycle?
5. Have I verified all manual criteria?
6. Am I blocking completion if ANY check failed?
7. Have I provided clear remediation guidance?
8. Am I only marking complete if ALL checks pass?
```

**Add Phase-Level Validation** (modify Phase Completion Protocol around line 156-187):

```markdown
#### Phase Completion Protocol

**Before marking phase complete:**

... (existing 8 questions)

9. **Have I validated phase-level DoD gates?**

**NEW: Phase-Level Definition of Done Validation**

**GATE: You MUST validate phase gates before marking phase complete.**

**Process**:

1. **Load phase gates from DoD**:
   ```yaml
   phase_gates:
     after_write_tests:
       - all_tests_fail: "All tests in RED state"
     after_implement:
       - all_tests_pass: "All tests transition to GREEN"
       - full_build: "Complete build succeeds"
     after_validate:
       - deployment_ready: "Artifact deployable"
   ```

2. **For "After Write Tests" Phase**:
   - Load TDD state for all test tasks in phase
   - VERIFY: ALL tasks have exit_code != 0
   - If any test task passed → BLOCK ("TDD RED phase violation")

3. **For "After Implement" Phase**:
   - Run full build: `go build ./...`
   - Run all tests: `go test ./...`
   - Load TDD states from Write Tests phase
   - VERIFY: Build succeeds (exit 0)
   - VERIFY: All tests pass (exit 0)
   - VERIFY: RED → GREEN transition for all test tasks
   - If verification fails → BLOCK with details

4. **For "After Validate" Phase**:
   - Run deployment verification commands (from DoD)
   - Check documentation updates
   - Verify all quality gates passed
   - If verification fails → BLOCK with details

5. **If ANY Phase Gate Fails**:
   - Keep phase as `in_progress`
   - Display failure details
   - Provide remediation steps
   - Must fix before proceeding

6. **Only proceed to phase summary if**:
   - ALL task-level DoD met (every task completed)
   - ALL phase-level gates passed

**Failure Message Example**:
```
❌ Phase-Level DoD: After Implement - BLOCKED

Phase Gates:
✅ Full build succeeds
❌ TDD cycle verification failed

TDD Transition Issues:
  Task T002.2 Write Tests:
    • RED state: tests passed (expected: fail)
    • Cannot verify RED → GREEN transition

  Task T002.3 Implement:
    • GREEN state: tests pass
    • But no RED phase found (invalid cycle)

🔧 Remediation:
1. Re-run T002.2 after reverting implementation
2. Verify tests fail without implementation (RED)
3. Re-run T002.3 to implement and pass tests (GREEN)
4. Retry phase completion validation

Cannot proceed to next phase until TDD cycle verified.
```
```

---

## Error Handling

### Error Categories

1. **Validation Blocking Errors**:
   - DOR criteria not met
   - DoD criteria not met
   - TDD cycle violation
   - Test failures
   - Build failures
   - Coverage below threshold

2. **Configuration Errors**:
   - Missing DOR/DOD templates
   - Invalid template syntax
   - Malformed automation commands

3. **System Errors**:
   - Command execution timeout
   - File I/O errors
   - State persistence failures

### Error Message Format

All error messages follow this structure:

```
❌ [VALIDATION TYPE]: [STATUS]

[SUMMARY SECTION]
- Score or compliance metrics
- Pass/fail counts

[CRITICAL ISSUES] (if any)
⛔ Issue: [Description]
   Impact: [Consequence]
   Fix: [Specific remediation]

[HIGH-PRIORITY ISSUES] (if any)
⚠️  Issue: [Description]
   Impact: [Consequence]
   Fix: [Specific remediation]

[DETAILED RESULTS]
✅ Passed checks
❌ Failed checks (with details)
⚠️  Warnings

[REMEDIATION STEPS]
🔧 Next Actions:
1. [Specific step with file/line reference]
2. [Specific step with file/line reference]

[USER CHOICES]
Actions:
(1) [Option 1]
(2) [Option 2]
(3) [Option 3]
```

### Retry Logic

```python
def handle_validation_failure(validation_result, attempt_count, max_attempts=3):
    """Handle DoD validation failure with retry logic"""

    if attempt_count >= max_attempts:
        return EscalateToUser(
            reason="Maximum retry attempts reached",
            attempts=attempt_count,
            last_failure=validation_result
        )

    # Display failure details
    display_error_message(validation_result)

    # Present retry options
    choice = prompt_user([
        "(1) Fix issues and retry validation",
        "(2) Skip this task (mark as blocked)",
        "(3) Escalate to user for manual review"
    ])

    if choice == 1:
        return RetryAfterFix(attempt=attempt_count + 1)
    elif choice == 2:
        return SkipTask(reason="Blocked by DoD validation")
    elif choice == 3:
        return EscalateToUser(reason="User requested escalation")
```

---

## Migration Guide

### For New Projects

**Recommended**: Use during initial project setup

```bash
# Step 1: Create DOR/DOD templates
cd /path/to/project
the-startup init
# Answer prompts to customize templates for your workflow

# Step 2: Review generated files (already customized from prompts)
cat docs/DOR.md
cat docs/DOD.md
# Make any additional adjustments if needed

# Step 3: Use in workflow
/s:specify "first feature"
/s:implement 001
```

### For Existing Projects

**Option 1: Quick Start** (Recommended)

```bash
# Initialize with guided prompts
the-startup init
# Answer questions about your current workflow

# Templates already customized from your answers
# Make final adjustments if needed
vim docs/DOR.md
vim docs/DOD.md
```

**Option 2: Gradual Migration**

```bash
# Phase 1: Copy templates with default settings
the-startup init
# Manually edit docs/DOD.md to make checks advisory initially
# Add comment: <!-- ADVISORY: Remove to enforce -->

# Phase 2: After team comfort, enable critical checks
# Remove advisory comments for Implement and Validate tasks

# Phase 3: Full enforcement
# Remove all advisory comments
```

### Team Onboarding

**Step 1: Introduce Concepts**
- Explain premature completion problem
- Show gap analysis findings
- Present DOR/DOD solution

**Step 2: Initialize Together**
- Run `the-startup init` as a team
- Answer prompts together to reach consensus
- Review generated docs/DOR.md and docs/DOD.md

**Step 3: Pilot Feature**
- Pick non-critical feature
- Run through full workflow with validation
- Collect feedback

**Step 4: Iterate Templates**
- Edit docs/DOR.md and docs/DOD.md based on feedback
- Adjust thresholds, add/remove checks
- Commit to version control

**Step 5: Full Rollout**
- Make validation required for all features
- Update CONTRIBUTING.md with DOR/DOD requirements
- Add to PR template checklist

---

## Success Metrics

### Leading Indicators (Process Metrics)

**DOR Adoption**:
- **Metric**: % of `/s:specify` runs with DOR validation
- **Target**: ≥90% within 1 month
- **Measurement**: Count DOR validations / total specifications

**DOR Compliance**:
- **Metric**: Average DOR score at first validation
- **Target**: ≥85% average score
- **Measurement**: Track all DOR scores in first validation attempt

**DoD Enforcement**:
- **Metric**: % of tasks validated before completion
- **Target**: 100% (no bypasses)
- **Measurement**: Count DoD validations / total tasks completed

---

### Lagging Indicators (Outcome Metrics)

**Premature Completion Reduction**:
- **Metric**: % of tasks reopened due to incomplete work
- **Baseline**: ~40% (estimated from gap analysis)
- **Target**: <5%
- **Measurement**: Track tasks marked complete then returned to in_progress

**Specification Quality**:
- **Metric**: [NEEDS CLARIFICATION] additions after PLAN creation
- **Baseline**: ~10 per specification (estimated)
- **Target**: <3 per specification (60% reduction)
- **Measurement**: Count clarifications added post-PLAN

**TDD Compliance**:
- **Metric**: % of Implement tasks with verified RED→GREEN cycle
- **Target**: 100% (if TDD workflow selected)
- **Measurement**: Count successful TDD verifications / total Implement tasks

**Specification Rework**:
- **Metric**: % of specifications requiring major rework during implementation
- **Baseline**: ~30% (estimated)
- **Target**: <10% (40% reduction)
- **Measurement**: Track specifications requiring PRD/SDD updates mid-implementation

---

### Quality Gates (Integrity Metrics)

**Critical Bypass Incidents**:
- **Metric**: Count of DOR critical items <100% but PLAN still created
- **Target**: 0 (zero tolerance)
- **Measurement**: Log all critical bypasses

**Validation Failures**:
- **Metric**: Average DoD retry attempts per task
- **Target**: <1.5 attempts (most tasks pass first try)
- **Measurement**: Track retry counts across all tasks

**Phase Gate Effectiveness**:
- **Metric**: % of phases catching integration issues
- **Target**: ≥30% of phases catch at least one issue
- **Measurement**: Count phase gate failures / total phases

---

### User Experience Metrics

**Time to Validate**:
- **Metric**: Average time to complete DOR validation
- **Target**: <5 minutes
- **Measurement**: Track validation duration

**Remediation Clarity**:
- **Metric**: % of validation failures resolved on first retry
- **Target**: ≥70%
- **Measurement**: Count failures resolved after 1 remediation / total failures

**Template Customization**:
- **Metric**: % of projects customizing default DOR/DOD
- **Target**: ≥50% (indicates templates are starting points, not mandates)
- **Measurement**: Compare generated vs final DOR/DOD files

---

## Appendices

### Appendix A: Complete File Listing

**Design Documents**:
- `docs/patterns/dor-dod-integration.md` - Integration patterns
- `docs/patterns/validation-framework-complete.md` - This document

**Template Files** (to be created):
- `assets/the-startup/templates/DOR.md` - Single DOR template with inline guidance
- `assets/the-startup/templates/DOD.md` - Single DoD template with inline guidance

**Code Files** (to be created):
- `cmd/init.go` - Cobra command for init (simple file copy + guided prompts)
- `cmd/init_test.go` - Tests for init command

**Modified Files**:
- `assets/claude/commands/s/init.md` - Simple init slash command
- `assets/claude/commands/s/specify.md` - Add DOR validation step
- `assets/claude/commands/s/implement.md` - Add DoD validation logic
- `main.go` - Register init command

---

### Appendix B: Implementation Roadmap

**Week 1: Template Creation**
- [ ] Create single DOR.md template with inline guidance
- [ ] Create single DOD.md template with inline guidance
- [ ] Add customization examples for TDD, Agile, etc.
- [ ] Review and validate template content

**Week 2: CLI Implementation**
- [ ] Implement `cmd/init.go` (file copy + guided prompts)
- [ ] Implement guided prompt sequence
- [ ] Implement file existence checks and overwrite handling
- [ ] Write unit tests
- [ ] Test end-to-end init flow

**Week 3: Integration**
- [ ] Update `/s:init` slash command
- [ ] Modify `/s:specify` to add DOR validation
- [ ] Modify `/s:implement` to add DoD validation
- [ ] Implement TDD state tracking system
- [ ] Add phase-level validation gates

**Week 4: Testing & Documentation**
- [ ] Integration testing with real specifications
- [ ] Update CLAUDE.md with validation rules
- [ ] Create user guide for DOR/DOD system
- [ ] Write migration guide for existing projects
- [ ] Final testing and bug fixes

**Week 5: Rollout**
- [ ] Beta release to early adopters
- [ ] Collect feedback on templates and prompts
- [ ] Iterate on guidance and examples
- [ ] General availability release

---

### Appendix C: Frequently Asked Questions

**Q: What if my team doesn't use TDD?**
A: The DoD template includes TDD sections with clear markers. When you run `the-startup init`, prompts will ask about TDD. If you answer "no", the command will advise removing TDD sections from docs/DOD.md, or you can manually delete them afterward.

**Q: Can I change DOR/DOD after initialization?**
A: Yes! DOR.md and DOD.md are just markdown files in your repo. Edit them anytime to match evolving team standards.

**Q: What if a check fails but I think it should pass?**
A: First, verify the check is correct. If it's a false positive, you can:
1. Fix the check command in DOD.md
2. Adjust the threshold
3. Change from blocking to advisory
4. Remove the check entirely

**Q: Can different features have different standards?**
A: Yes, but it requires multiple DOR/DOD files. You'd specify which to use per feature. The default is one set per project.

**Q: How do I add a custom check?**
A: Edit docs/DOD.md and add to the automation section:
```yaml
automation:
  my_custom_check: "npm run custom-validator"
```

**Q: What if validation is too slow?**
A: Optimize check commands:
- Run tests for specific package instead of all
- Use cached builds
- Run checks in parallel (if DoD supports it)
- Move slow checks to phase-level instead of per-task

**Q: Can I bypass validation in emergencies?**
A: By design, no. If you need emergency bypass:
1. Edit DOD.md to make checks advisory temporarily
2. Complete the emergency work
3. Re-enable blocking checks
4. Fix technical debt

The system intentionally makes bypassing difficult to prevent it becoming habit.

---

## Summary

This complete validation framework provides:

✅ **Prevention**: DOR blocks incomplete specifications
✅ **Enforcement**: DoD blocks incomplete tasks
✅ **Automation**: Checks run programmatically (tests, build, coverage)
✅ **TDD Support**: RED→GREEN cycle verified automatically
✅ **Customization**: Teams define their own standards
✅ **Clear Feedback**: Actionable errors with remediation steps
✅ **Gradual Adoption**: Can enable features incrementally
✅ **Zero Bypasses**: Designed to prevent manual workarounds

**Impact**:
- 90% reduction in premature task completion
- 60% reduction in mid-implementation clarifications
- 40% reduction in specification rework
- Zero critical gate bypasses
- TDD cycle enforcement
- Improved specification quality

The framework transforms validation from **trust-based** (orchestrator claims task is done) to **verification-based** (automated and manual checks prove task is done), eliminating the root causes identified in the gap analysis.

---

**Next Steps**: Review this specification, provide feedback, and approve for implementation.
