# Task Definition of Done (TASK-DOD)

<!--
This file defines completion criteria for IMPLEMENTATION TASKS in PLAN.md.

WHEN TO USE:
- After completing each MAIN TASK during /s:implement
- TASK-DOD validates implementation quality
- DOR/DOD validate document creation quality

AUTOMATED VALIDATION:
- Build success checks
- Test execution and coverage
- Linting and formatting
- TDD cycle validation (if enabled)
-->

---

## Task Completion Criteria

### Build Success

**Automated Checks**:
- [ ] Project builds without errors
  ```bash
  [NEEDS CLARIFICATION: build command]
  # Expected: exit code 0
  ```

**Manual Verification**:
- [ ] Does the build complete without warnings?
- [ ] Are there any deprecated API usages?

### Test Execution

**Automated Checks**:
- [ ] All tests pass
  ```bash
  [NEEDS CLARIFICATION: test command]
  # Expected: exit code 0, all tests passing
  ```

**Manual Verification**:
- [ ] Are the new tests meaningful and testing the right behavior?
- [ ] Do the tests cover edge cases (null, empty, boundaries)?

### Code Coverage

**Automated Checks**:
- [ ] Coverage meets threshold
  ```bash
  [NEEDS CLARIFICATION: coverage command]
  # Expected: coverage ≥ [NEEDS CLARIFICATION: coverage target]%
  ```

**Manual Verification**:
- [ ] Is the coverage from meaningful tests (not just mocks)?
- [ ] Are critical paths covered by tests?

### Code Quality

**Automated Checks**:
- [ ] Linting passes
  ```bash
  [NEEDS CLARIFICATION: lint command]
  # Expected: exit code 0, no lint errors
  ```

- [ ] Code is formatted
  ```bash
  [NEEDS CLARIFICATION: format command]
  # Expected: exit code 0, no formatting changes needed
  ```

**Manual Verification**:
- [ ] Is the code readable and maintainable?
- [ ] Are variable names clear and intention-revealing?
- [ ] Are functions small and focused (single responsibility)?

<!-- OPTIONAL: TDD -->
### TDD Cycle Validation

**Automated Checks**:
- [ ] RED phase documented (failing test exists)
  ```bash
  # Check for test file creation timestamp before implementation
  # Verify initial test run failed
  ```

- [ ] GREEN phase documented (tests now pass)
  ```bash
  # Verify current test run succeeds
  # Implementation file timestamp after test file
  ```

**Manual Verification**:
- [ ] Was the test written BEFORE the implementation?
- [ ] Does the test fail for the right reason (not syntax error)?
- [ ] Did you write minimal code to make the test pass?
<!-- END OPTIONAL: TDD -->

---

## Phase-Specific Validation

### Prime Phase (Red - Write Tests)

**Automated Checks**:
- [ ] Test files exist for component
- [ ] Tests are initially failing (RED state)

**Manual Verification**:
- [ ] Do the tests describe the expected behavior clearly?
- [ ] Are all acceptance criteria covered by tests?

### Test Phase (Still Red - Verify Tests Fail)

**Automated Checks**:
- [ ] Tests still fail (confirming RED state)
- [ ] Test failures are meaningful (not syntax errors)

**Manual Verification**:
- [ ] Do the failures indicate what needs to be implemented?
- [ ] Are the test descriptions accurate?

### Implement Phase (Green - Make Tests Pass)

**Automated Checks**:
- [ ] All tests now pass (GREEN state)
- [ ] No skipped or ignored tests

**Manual Verification**:
- [ ] Did you implement the simplest solution that works?
- [ ] Are there any hardcoded values that should be parameters?

### Validate Phase (Refactor - Improve Code)

**Automated Checks**:
- [ ] Tests still pass after refactoring
- [ ] Code coverage maintained or improved
- [ ] No new lint warnings introduced

**Manual Verification**:
- [ ] Is the code more readable than before?
- [ ] Did you remove duplication?
- [ ] Are abstractions justified (not premature)?

---

## Integration Validation

### Component Integration

**Automated Checks**:
- [ ] Integration tests pass
  ```bash
  [NEEDS CLARIFICATION: integration test command]
  # Expected: all integration tests passing
  ```

**Manual Verification**:
- [ ] Do components interact correctly?
- [ ] Are interfaces used properly?
- [ ] Is error handling complete?

### End-to-End Validation

**Automated Checks**:
- [ ] E2E tests pass (if applicable)
  ```bash
  [NEEDS CLARIFICATION: e2e test command]
  # Expected: user journeys work end-to-end
  ```

**Manual Verification**:
- [ ] Can you manually test the feature?
- [ ] Does it match the PRD acceptance criteria?

---

<!-- OPTIONAL: SECURITY_CHECKS -->
## Security Validation

**Automated Checks**:
- [ ] Security scan passes
  ```bash
  [NEEDS CLARIFICATION: security scan command]
  # Expected: no high or critical vulnerabilities
  ```

**Manual Verification**:
- [ ] Are inputs validated and sanitized?
- [ ] Are secrets handled securely (not hardcoded)?
- [ ] Is authentication/authorization implemented correctly?
<!-- END OPTIONAL: SECURITY_CHECKS -->

---

<!-- OPTIONAL: PERFORMANCE_CHECKS -->
## Performance Validation

**Automated Checks**:
- [ ] Performance benchmarks pass
  ```bash
  [NEEDS CLARIFICATION: benchmark command]
  # Expected: within acceptable performance thresholds
  ```

**Manual Verification**:
- [ ] Are there any obvious performance bottlenecks?
- [ ] Is resource usage (memory, CPU) acceptable?
<!-- END OPTIONAL: PERFORMANCE_CHECKS -->

---

## Minimum Quality Threshold

**Automated Score**: [AUTOMATED: Calculate based on checks above]
- Critical items complete: [X/Y] (must be 100%)
- Overall items complete: [X/Z] (must be ≥[NEEDS CLARIFICATION: task dod threshold]%)

**If threshold not met**: BLOCK task completion with specific failures

---

## Failure Message Format

```
❌ Task Definition of Done: BLOCKED

Task: [Task Name]
Overall: [X/Y] checks passed ([Z]%)
Critical: [A/B] checks passed ([C]%)

⛔ Critical Failures:
  • Tests failing (3 failures)
    Impact: Feature is broken
    Fix: Run tests and fix failures

⚠️ Non-Critical Failures:
  • Coverage below threshold (72% < 80%)
    Impact: Insufficient test coverage
    Fix: Add tests for uncovered code paths

📊 Automated Checks:
  ✅ Build succeeds
  ❌ All tests pass (3 failures)
  ⚠️  Coverage meets threshold (72% < 80%)
  ✅ Linting passes
  ✅ Code is formatted

🔧 Next Actions:
  1. Fix failing tests
  2. Add tests to reach 80% coverage
  3. Re-run validation

Cannot mark task complete until critical items resolved.
```

---

## Configuration

### Thresholds

```yaml
thresholds:
  critical: 100  # Critical items must be 100% complete
  overall: [NEEDS CLARIFICATION: task dod threshold]  # Overall threshold (default: 85%)
```

### Commands

```yaml
commands:
  build: [NEEDS CLARIFICATION: build command]
  test: [NEEDS CLARIFICATION: test command]
  coverage: [NEEDS CLARIFICATION: coverage command]
  lint: [NEEDS CLARIFICATION: lint command]
  format: [NEEDS CLARIFICATION: format command]

# Optional commands
  integration_test: [NEEDS CLARIFICATION: integration test command]
  e2e_test: [NEEDS CLARIFICATION: e2e test command]
  security_scan: [NEEDS CLARIFICATION: security scan command]
  benchmark: [NEEDS CLARIFICATION: benchmark command]
```

### Coverage Target

```yaml
coverage:
  target: [NEEDS CLARIFICATION: coverage target]  # Default: 80%
```

---

## Customization Guide

### Adding Custom Checks

```markdown
#### Custom Check Name

**Automated Checks**:
- [ ] Your custom check (check: command to run)
  ```bash
  your-validation-command
  # Expected: exit code 0
  ```

**Manual Verification**:
- [ ] Your manual verification question?
```

### Language-Specific Checks

```markdown
<!-- IF: LANGUAGE_GO -->
#### Go-Specific Checks

**Automated Checks**:
- [ ] No race conditions (check: go test -race ./...)
- [ ] Vet passes (check: go vet ./...)
<!-- END IF -->

<!-- IF: LANGUAGE_JAVASCRIPT -->
#### JavaScript-Specific Checks

**Automated Checks**:
- [ ] TypeScript compilation (check: tsc --noEmit)
- [ ] Bundle size acceptable (check: bundlesize)
<!-- END IF -->
```

---

## FAQ

**Q: What if tests pass but manually testing reveals issues?**
A: Manual verification is part of TASK-DOD. Don't mark complete until both automated and manual checks pass.

**Q: Can I skip TDD cycle validation?**
A: Yes, if TDD is disabled in your workflow. Remove or comment out the TDD section.

**Q: What if I can't reach coverage threshold?**
A: Either add more tests or adjust the threshold if the uncovered code is truly untestable (logs, error messages).

**Q: How do I know if refactoring is needed?**
A: Look for code smells: duplication, long functions, unclear names, tight coupling.

**Q: What's the difference between TASK-DOD and DOD.md?**
A: TASK-DOD validates individual implementation tasks. DOD.md validates complete documents (PRD, SDD, PLAN).
