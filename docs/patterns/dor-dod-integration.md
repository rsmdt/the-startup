# DOR/DOD Integration with Specification Workflow

## Overview

This document describes how Definition of Ready (DOR) and Definition of Done (DOD) templates integrate with the specification workflow (`/s:specify` and `/s:implement` commands) to prevent premature task completion.

## System Architecture

```
┌─────────────────────────────────────────────────────────────┐
│  User runs: the-startup init                                │
│  Creates: docs/DOR.md, docs/DOD.md                          │
└───────────────────────┬─────────────────────────────────────┘
                        │
                        ▼
┌─────────────────────────────────────────────────────────────┐
│  /s:specify: Create specification (PRD, SDD, PLAN)          │
│  ├─ Step 1: Gather requirements                             │
│  ├─ Step 2: Create PRD                                      │
│  ├─ Step 3: Create SDD                                      │
│  ├─ Step 4: DOR VALIDATION GATE ◄─── NEW                    │
│  │   ├─ Read docs/DOR.md                                    │
│  │   ├─ Present readiness checklist                         │
│  │   ├─ Calculate compliance score                          │
│  │   └─ BLOCK if score < 85% or critical items incomplete   │
│  └─ Step 5: Create PLAN (only if DOR passes)                │
└───────────────────────┬─────────────────────────────────────┘
                        │
                        ▼
┌─────────────────────────────────────────────────────────────┐
│  /s:implement: Execute implementation plan                  │
│  For each task:                                             │
│    ├─ Task starts → TodoWrite: pending → in_progress        │
│    ├─ Agent executes task                                   │
│    ├─ Agent reports completion                              │
│    ├─ DOD VALIDATION GATE ◄─── NEW                          │
│    │   ├─ Read docs/DOD.md                                  │
│    │   ├─ Identify task type (Prime/Test/Implement/Validate)│
│    │   ├─ Present task-specific DoD checklist               │
│    │   ├─ Run automated checks (go test, go build, etc.)    │
│    │   ├─ Verify manual criteria                            │
│    │   └─ BLOCK TodoWrite "completed" if any check fails    │
│    └─ TodoWrite: in_progress → completed (only if DoD passes)│
└─────────────────────────────────────────────────────────────┘
```

## Component Interactions

### 1. Initialization Phase

**Command**: `the-startup init`

**Process**:
1. User runs interactive wizard or selects preset
2. CLI generates customized DOR.md and DOD.md
3. Files written to `docs/` directory
4. Team can customize templates further

**Output Files**:
- `docs/DOR.md` - Readiness checklist for specifications
- `docs/DOD.md` - Completion checklist for tasks

---

### 2. Specification Creation (`/s:specify`)

**DOR Integration Point**: Between SDD completion and PLAN creation

#### Current Flow (No Validation)
```
/s:specify "add authentication"
  → Gather requirements
  → Create PRD
  → Create SDD
  → Create PLAN  ← No validation gate
```

#### Enhanced Flow (With DOR)
```
/s:specify "add authentication"
  → Gather requirements
  → Create PRD
  → Create SDD
  → READ docs/DOR.md
  → VALIDATE readiness
      ├─ Present checklist to orchestrator
      ├─ Orchestrator verifies each item
      ├─ Calculate score: 32/34 items = 94%
      └─ Decision:
          • Score ≥85% AND critical items 100% → PROCEED to PLAN
          • Else → BLOCK with failure message
  → Create PLAN (if approved)
```

#### Implementation in `/s:specify`

**File**: `assets/claude/commands/s/specify.md`

**Add New Step Between SDD and PLAN**:

```markdown
### 📋 Step 3.5: Definition of Ready Validation

**🎯 Goal**: Verify specification meets readiness criteria before creating implementation plan

**Process**:
1. Check for DOR template:
   - If `docs/DOR.md` exists → use project-specific DOR
   - Else → use default DOR from assets/the-startup/templates/DOR.md.tmpl

2. Read DOR checklist and extract:
   - All checklist items `- [ ] Item description`
   - Critical items (marked with `[CRITICAL]` tag)
   - Validation questions for each category

3. Present DOR to orchestrator:
   ```
   📋 Definition of Ready Validation

   Before creating the implementation plan, verify readiness:

   ## Problem Definition (3/3 complete)
   ✅ Problem clearly articulated
   ✅ Stakeholders identified
   ✅ Success criteria defined

   ## Requirements Clarity (4/5 complete) ⚠️
   ✅ Functional requirements listed
   ✅ Non-functional requirements specified
   ⚠️  Edge cases not fully considered
   ✅ Acceptance criteria clear

   ... (other categories)

   Score: 32/34 items (94%)
   Critical: 7/7 (100%)

   Status: ✅ READY TO PROCEED
   ```

4. Enforcement logic:
   - If critical items < 100% → BLOCK with error message
   - If overall score < 85% → BLOCK with warning
   - Else → PROCEED to PLAN creation

5. If BLOCKED:
   - Display specific unmet criteria
   - Show which PRD/SDD sections need updates
   - Ask: "Return to PRD (1), SDD (2), or cancel (3)?"
   - Allow iteration back to earlier steps

**🤔 Ask yourself before proceeding:**
1. Have I read the DOR template completely?
2. Have I verified each checklist item honestly?
3. Are ALL critical items complete?
4. Is the overall score ≥85%?
5. If blocked, have I identified specific gaps?
6. Am I about to create PLAN only if DOR passes?
```

**Validation Algorithm**:

```python
def validate_dor(prd_content, sdd_content, dor_template):
    """Validates specification against DOR criteria"""

    # Parse DOR checklist
    items = parse_dor_checklist(dor_template)
    critical_items = [i for i in items if i.is_critical]

    # Present each item for verification
    results = []
    for item in items:
        # Show item to orchestrator
        response = verify_item(item, prd_content, sdd_content)
        results.append(response)

    # Calculate scores
    total_score = sum(results) / len(items) * 100
    critical_score = sum(r for i, r in zip(items, results) if i.is_critical)
    critical_total = len(critical_items)

    # Enforcement decision
    if critical_total > 0 and critical_score < critical_total:
        return BlockResult(
            reason="Critical items incomplete",
            missing=critical_items[critical_score:],
            remediation="Complete all critical items before proceeding"
        )

    if total_score < 85:
        return BlockResult(
            reason="Overall readiness below threshold",
            score=total_score,
            threshold=85,
            remediation="Address flagged items in PRD/SDD"
        )

    return ApprovedResult(score=total_score)
```

---

### 3. Implementation Execution (`/s:implement`)

**DOD Integration Point**: Before marking each task `completed`

#### Current Flow (No Validation)
```
Task: "T001.3 Implement authentication"
  → Mark in_progress
  → Agent implements
  → Agent reports done
  → Mark completed  ← No validation gate
```

#### Enhanced Flow (With DOD)
```
Task: "T001.3 Implement authentication"
  → Mark in_progress
  → Agent implements
  → Agent reports done
  → READ docs/DOD.md
  → IDENTIFY task type: "Implement"
  → VALIDATE against DoD
      ├─ Automated checks:
      │   ├─ go build ./... (exit 0?)
      │   ├─ go test ./... (exit 0?)
      │   ├─ grep "// SDD Section" (references present?)
      │   └─ go test -cover (≥ threshold?)
      ├─ Manual verification:
      │   ├─ Spec requirements met?
      │   └─ No new warnings?
      └─ Decision:
          • ALL checks pass → ALLOW completion
          • ANY check fails → BLOCK with failure details
  → Mark completed (if approved)
```

#### Implementation in `/s:implement`

**File**: `assets/claude/commands/s/implement.md`

**Modify Task Completion Logic (Step 2 - Phase Execution)**:

**Current** (lines 100-132):
```markdown
**📝 For Sequential Tasks:**
- Execute one at a time
- Mark as `in_progress` in TodoWrite
- Delegate to specialist agent
- After completion, mark `completed` in TodoWrite
```

**Enhanced**:
```markdown
**📝 For Sequential Tasks:**
- Execute one at a time
- Mark as `in_progress` in TodoWrite
- Extract SDD references from task
- Delegate to specialist agent with specification context
- **After agent reports completion → DOD VALIDATION**

  **🛡️ Definition of Done Validation:**
  1. Read DOD template:
     - If `docs/DOD.md` exists → use project-specific DoD
     - Else → use default DoD from assets

  2. Identify task type from PLAN.md task ID:
     - `T00X.1` = Prime Context
     - `T00X.2` = Write Tests
     - `T00X.3` = Implement
     - `T00X.5` = Validate

  3. Load task-type specific DoD checklist

  4. Run automated checks (if defined):
     ```yaml
     automation:
       build: "go build ./..."
       test: "go test ./... -v"
       coverage: "go test -cover ./..."
       lint: "golangci-lint run"
     ```
     - Execute each command
     - Capture exit codes and output
     - Mark check as PASS/FAIL

  5. Present manual verification prompts:
     ```
     📋 Definition of Done: Implement Task

     Automated Checks:
     ✅ Build succeeds (go build ./...)
     ✅ Tests pass (go test ./...)
     ⚠️  Coverage 75% (threshold: 80%)
     ✅ Linting passes

     Manual Verification:
     ? Specification requirements met (check SDD references)
     ? No new warnings introduced

     Score: 4/6 checks (67%)
     Status: ❌ BLOCKED (coverage below threshold)
     ```

  6. Enforcement decision:
     - If ANY blocking check fails → BLOCK completion
     - Display failure details with remediation steps
     - Keep task as `in_progress`
     - Retry after fixes (max 3 attempts)
     - After 3 failed attempts → escalate to user

  7. Only if ALL checks pass:
     - Mark `completed` in TodoWrite
     - Update PLAN.md checkbox
     - Proceed to next task

**🤔 Ask yourself at DoD validation:**
1. Have I run all automated checks?
2. Did ALL automated checks pass?
3. Have I verified all manual criteria?
4. Is this task truly complete per DoD definition?
5. Am I blocking completion if any check failed?
6. Have I provided clear remediation guidance?
```

**DoD Validation Algorithm**:

```python
def validate_dod(task_id, task_type, dod_template, tdd_state):
    """Validates task against DoD criteria before marking complete"""

    # Load task-type specific checklist
    checklist = load_dod_for_task_type(dod_template, task_type)

    # Run automated checks
    auto_results = {}
    for check_name, command in checklist.automation.items():
        result = run_shell_command(command)
        auto_results[check_name] = {
            'passed': result.exit_code == 0,
            'output': result.stdout,
            'error': result.stderr
        }

    # Special: TDD cycle verification for Write Tests and Implement
    if task_type == "Write Tests":
        # Tests should FAIL (RED phase)
        test_result = auto_results['test']
        if test_result['passed']:  # Tests passed when should fail!
            return BlockResult(
                reason="TDD RED phase violation",
                detail="Tests passed but should fail before implementation",
                remediation="Ensure tests verify behavior not yet implemented"
            )
        # Store state for next task
        save_tdd_state(task_id, {
            'exit_code': test_result['exit_code'],
            'test_count': extract_test_count(test_result['output'])
        })

    if task_type == "Implement":
        # Tests should PASS (GREEN phase) AND previously failed
        test_result = auto_results['test']
        previous_state = load_tdd_state(previous_task_id(task_id))

        if not test_result['passed']:
            return BlockResult(
                reason="TDD GREEN phase violation",
                detail="Tests still failing after implementation",
                remediation="Fix implementation until tests pass"
            )

        if previous_state['exit_code'] == 0:
            return BlockResult(
                reason="TDD cycle not followed",
                detail="Tests were already passing (no RED phase)",
                remediation="Write failing tests first, then implement"
            )

    # Check all automated results
    failed_auto = [name for name, result in auto_results.items()
                   if not result['passed']]
    if failed_auto:
        return BlockResult(
            reason="Automated checks failed",
            failed_checks=failed_auto,
            details=auto_results,
            remediation=generate_remediation(failed_auto, auto_results)
        )

    # Present manual verification prompts
    manual_results = []
    for manual_check in checklist.manual:
        response = prompt_manual_verification(manual_check)
        manual_results.append(response)

    failed_manual = [check for check, result in zip(checklist.manual, manual_results)
                     if not result]
    if failed_manual:
        return BlockResult(
            reason="Manual verification failed",
            failed_checks=failed_manual,
            remediation="Address flagged criteria before marking complete"
        )

    # All checks passed
    return ApprovedResult(
        auto_passed=len(auto_results),
        manual_passed=len(manual_results)
    )
```

---

### 4. Phase-Level DoD Validation

**Integration Point**: Phase completion (after all tasks in phase)

**Add to Phase Completion Protocol** (implement.md lines 156-187):

```markdown
#### Phase Completion Protocol

**Before marking phase complete:**

... (existing 8 questions)

9. **Have I validated phase-level DoD criteria?**

**📋 Phase-Level Definition of Done:**

In addition to individual task DoD, verify phase-level criteria:

1. Read phase-specific DoD section from docs/DOD.md:
   ```yaml
   phase_gates:
     after_write_tests:
       - all_tests_fail: "All tests should be in RED state"
       - test_coverage: "Test coverage includes edge cases"

     after_implement:
       - full_build: "Complete build succeeds"
       - all_tests_pass: "All tests transition from RED to GREEN"
       - integration: "Component integration verified"

     after_validate:
       - deployment_ready: "Artifact deployable"
       - documentation: "README and API docs updated"
   ```

2. Run phase-gate checks:
   - For "After Write Tests" phase:
     - Verify ALL tests in phase are failing (RED state)
     - Check TDD state tracking: all test tasks have exit_code != 0

   - For "After Implement" phase:
     - Run full build: `go build ./...`
     - Verify ALL tests now pass: `go test ./...`
     - Check TDD transition: RED → GREEN for all test tasks

   - For "After Validate" phase:
     - Run deployment verification
     - Check documentation updates
     - Verify all quality gates passed

3. Enforcement:
   - If phase gate fails → BLOCK phase completion
   - Keep phase as in_progress
   - Provide specific remediation
   - Must fix before proceeding to next phase

**Only proceed to phase summary if:**
- ALL task-level DoD criteria met (every task completed)
- ALL phase-level DoD criteria met (horizontal validation passed)
```

---

## Configuration Schema

### DOR.md Structure

```markdown
# Definition of Ready

## Configuration

<!-- validation:
  critical_threshold: 100
  overall_threshold: 85
  categories: 6
-->

## Problem Definition [CRITICAL]

- [ ] Problem clearly articulated in PRD
- [ ] Stakeholders identified and engaged
- [ ] Success criteria defined and measurable

**Validation Question**: Can you summarize the problem and explain how success will be measured?

## Requirements Clarity [CRITICAL]

- [ ] Functional requirements listed in PRD
- [ ] Non-functional requirements specified (performance, security, etc.)
- [ ] Edge cases and error scenarios considered
- [ ] Acceptance criteria clear and testable

**Validation Question**: Are there any [NEEDS CLARIFICATION] markers remaining in PRD/SDD?

... (other categories)
```

### DOD.md Structure

```yaml
---
# DoD Configuration
version: 1.0
workflow: tdd-strict

automation:
  build: "go build ./..."
  test: "go test ./... -v"
  coverage: "go test -cover ./..."
  lint: "golangci-lint run"

thresholds:
  coverage: 80
  complexity: 15

tdd:
  enforce_cycle: true
  require_red_phase: true
---

# Definition of Done

## Prime Context Tasks

- [ ] All referenced files read completely (not skimmed)
- [ ] SDD sections understood (can summarize key decisions)
- [ ] Interface contracts identified and documented

**Automated**: None
**Manual**: Comprehension verification

## Write Tests Tasks

- [ ] Tests written for specified behavior
- [ ] Tests currently FAIL (TDD red phase) ← Automated
- [ ] Failure messages match expected behavior
- [ ] Test coverage includes edge cases

**Automated**:
- `go test ./... -v` → exit code != 0 (must fail)

**Manual**:
- Verify tests cover edge cases
- Confirm failure messages are descriptive

## Implement Tasks

- [ ] Code compiles/builds successfully ← Automated
- [ ] Tests now PASS (TDD green phase) ← Automated
- [ ] Specification requirements met (SDD references present) ← Automated
- [ ] Test coverage meets threshold ← Automated
- [ ] No new warnings/errors introduced ← Manual

**Automated**:
- `go build ./...` → exit 0
- `go test ./...` → exit 0 AND previous state was != 0
- `grep -r "// SDD Section" <files>` → references found
- `go test -cover ./...` → ≥ 80%

**Manual**:
- Verify no new warnings in build output
- Confirm spec requirements implemented

## Validate Tasks

- [ ] Linting passed ← Automated
- [ ] Type checking passed ← Automated
- [ ] Code review completed ← Manual
- [ ] Integration tests pass ← Automated
- [ ] Specification compliance verified ← Manual

**Automated**:
- `golangci-lint run` → exit 0
- `go test -tags=integration ./...` → exit 0

**Manual**:
- Code review approval obtained
- Spec compliance confirmed

## Phase-Level Gates

### After Write Tests Phase
- [ ] All tests in phase are failing (RED state)
- [ ] Test coverage includes integration scenarios

### After Implement Phase
- [ ] Full build succeeds
- [ ] All tests transition RED → GREEN
- [ ] Component integration verified

### After Validate Phase
- [ ] Deployment artifact created
- [ ] Documentation updated (README, API docs)
- [ ] All quality standards met
```

---

## Error Handling

### DOR Validation Failure

**Scenario**: Specification doesn't meet readiness criteria

**Error Message Format**:
```
❌ Definition of Ready: BLOCKED

Overall Score: 78% (threshold: 85%)
Critical Items: 6/7 (threshold: 100%)

⛔ Critical Issues:
  • Edge cases not considered in PRD Section 3.2
    Impact: Implementation may miss important scenarios
    Fix: Add edge case analysis to PRD Section 3.2

⚠️  High-Priority Issues:
  • Test data not identified
    Impact: Testing may be delayed
    Fix: Add test data section to SDD Section 8

📊 Current State:
  Problem Definition:     ✅ 3/3 (100%)
  Requirements Clarity:   ⚠️  4/5 (80%)
  Technical Feasibility:  ✅ 6/6 (100%)
  Resource Availability:  ⚠️  5/7 (71%)
  Acceptance Criteria:    ✅ 4/4 (100%)
  Documentation:          ❌ 0/2 (0%)

🔧 Next Actions:
  1. Return to PRD to address edge cases (Section 3.2)
  2. Update SDD with test data (Section 8)
  3. Complete documentation (add [NEEDS CLARIFICATION] resolutions)

Choose action:
  (1) Return to PRD editing
  (2) Return to SDD editing
  (3) Cancel specification
```

### DOD Validation Failure

**Scenario**: Task doesn't meet completion criteria

**Error Message Format**:
```
❌ Definition of Done: BLOCKED

Task: T002.3 Implement user authentication
Type: Implement

Automated Checks:
  ✅ Build succeeds (go build ./...)
  ❌ Tests fail (go test ./...)
     Exit code: 1
     Failed tests:
       • TestLogin: expected 200, got 401
       • TestLogout: session not cleared
  ⚠️  Coverage 72% (threshold: 80%)
     Missing coverage:
       • auth/handlers.go: lines 45-67
       • auth/session.go: lines 23-29
  ✅ Linting passes (golangci-lint run)
  ❌ SDD references missing
     Expected: "// SDD Section 4.2" comments
     Found: 0 references

Manual Verification:
  ⚠️  Specification requirements: Not verified
  ⚠️  No new warnings: Not verified

Score: 2/8 checks (25%)

🔧 Remediation Steps:

1. Fix failing tests:
   File: auth/handlers_test.go:23
   Issue: Login returns 401 instead of expected 200
   Action: Check authentication logic in auth/handlers.go:45

2. Increase test coverage to 80%:
   Add tests for:
   • auth/handlers.go: lines 45-67 (password validation)
   • auth/session.go: lines 23-29 (session cleanup)

3. Add SDD reference comments:
   Add to auth/handlers.go:
   // SDD Section 4.2: Authentication Flow
   // Implements login endpoint as specified

Retry count: 1/3

Actions:
  (1) Fix issues and retry validation
  (2) Skip this task (mark as blocked)
  (3) Escalate to user for guidance
```

---

## Success Metrics

### DOR Effectiveness
- **Target**: 60% reduction in mid-implementation clarifications
- **Measure**: Track [NEEDS CLARIFICATION] additions after PLAN creation
- **Success**: <3 clarifications per specification on average

### DOD Effectiveness
- **Target**: 90% reduction in premature task completion
- **Measure**: Tasks reopened due to incomplete work
- **Success**: <5% task reopening rate

### Workflow Adoption
- **Target**: 90% of specifications use DOR
- **Measure**: Percentage of `/s:specify` runs with DOR validation
- **Success**: DOR validation runs in 90%+ of specifications

### Enforcement Integrity
- **Target**: Zero critical bypass incidents
- **Measure**: DOR critical items <100% but PLAN still created
- **Success**: 0 bypasses of critical gates

---

## Migration Path

### Existing Projects Without DOR/DOD

**Option 1: Initialize with Preset**
```bash
# Quick start with TDD preset
the-startup init --preset=tdd

# Review and customize
vim docs/DOR.md
vim docs/DOD.md
```

**Option 2: Interactive Setup**
```bash
# Full wizard for custom workflow
the-startup init
```

**Option 3: Manual Creation**
```bash
# Copy templates and customize manually
cp assets/the-startup/templates/DOR.md.tmpl docs/DOR.md
cp assets/the-startup/templates/DOD.md.tmpl docs/DOD.md
# Edit to match team workflow
```

### Gradual Rollout Strategy

**Phase 1: DOR Only**
- Enable DOR validation in `/s:specify`
- Keep DOD optional (warnings only)
- Collect feedback on DOR effectiveness

**Phase 2: DOD for Critical Tasks**
- Enable DOD for "Implement" and "Validate" tasks
- Keep "Prime" and "Write Tests" optional
- Monitor task completion quality

**Phase 3: Full Enforcement**
- Enable DOD for all task types
- Enable phase-level gates
- Achieve full TDD cycle enforcement

---

## Summary

This integration provides:

1. **Preventive Gates**: DOR blocks incomplete specifications
2. **Completion Gates**: DOD blocks incomplete tasks
3. **TDD Enforcement**: RED-GREEN cycle programmatically verified
4. **Customization**: Teams define their own criteria
5. **Automation**: Checks run automatically where possible
6. **Clear Feedback**: Actionable error messages guide fixes
7. **Gradual Adoption**: Can enable features incrementally

The system transforms validation from trust-based to verification-based, eliminating the root causes of premature completion identified in the gap analysis.
