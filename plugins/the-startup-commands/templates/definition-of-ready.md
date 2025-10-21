# Definition of Ready

<!--
This file defines prerequisites that must be met BEFORE creating each document.

AUTOMATED VALIDATION:
- Checks run automatically before document creation
- Failures block progression with specific remediation guidance
- Manual verification follows automated checks

USAGE:
- Before PRD: Validate "Before Creating PRD" section
- Before SDD: Validate "Before Creating SDD" section
- Before PLAN: Validate "Before Creating PLAN" section
-->

---

## Before Creating PRD

### Prerequisite Validation

Run these automated checks BEFORE starting PRD creation:

#### Research Completeness

**Automated Checks**:
- [ ] Problem statement drafted (check: draft document exists or user provided description)
- [ ] At least 2 stakeholders identified (check: stakeholder list has ≥2 entries)
- [ ] User research conducted (check: interview notes, survey results, or analytics exist)

**Manual Verification** (after automated checks pass):
- [ ] Can you articulate the problem in one sentence without jargon?
- [ ] Do you know who is affected and how badly?
- [ ] Have you validated the problem with actual users (not just assumptions)?

#### Context Gathering

**Automated Checks**:
- [ ] Competitive analysis done (check: competitor analysis document or notes exist)
- [ ] Success metrics identified (check: at least 1 measurable metric defined)

**Manual Verification**:
- [ ] Do you understand how competitors solve this problem?
- [ ] Do you know what "success" looks like quantitatively?
- [ ] Have you identified what you won't build (scope boundaries)?

#### Minimum Information Threshold

**Automated Score**: [AUTOMATED: Calculate based on checks above]
- Critical items complete: [X/5] (must be 100%)
- Overall items complete: [X/7] (must be ≥[NEEDS CLARIFICATION: dor threshold]%)

**If threshold not met**: BLOCK with specific gaps and remediation steps

---

## Before Creating SDD

### Prerequisite Validation

Run these automated checks BEFORE starting SDD creation:

#### PRD Completeness

**Automated Checks**:
- [ ] PRD file exists (check: docs/specs/[ID]/PRD.md present)
- [ ] PRD has no `[NEEDS CLARIFICATION]` markers (check: grep returns 0 matches)
- [ ] All PRD sections non-empty (check: no sections with only placeholders)
- [ ] PRD validation checklist complete (check: all boxes marked `[x]`)

**Manual Verification** (after automated checks pass):
- [ ] Do you understand all requirements in PRD without questions?
- [ ] Can you explain each user story in your own words?
- [ ] Do you know what success criteria must be met?

#### Technical Feasibility

**Automated Checks**:
- [ ] Architecture approach documented (check: approach notes exist or user provided summary)
- [ ] Technical dependencies identified (check: dependency list has ≥1 entry or "none" explicitly stated)

**Manual Verification**:
- [ ] Have you researched how to implement this technically?
- [ ] Do you know what libraries/frameworks will be used?
- [ ] Have you identified technical risks or unknowns?
- [ ] Is there proof-of-concept evidence for risky parts?

#### Minimum Information Threshold

**Automated Score**: [AUTOMATED: Calculate based on checks above]
- Critical items complete: [X/4] (must be 100%)
- Overall items complete: [X/7] (must be ≥[NEEDS CLARIFICATION: dor threshold]%)

**If threshold not met**: BLOCK with specific gaps and remediation steps

---

## Before Creating PLAN

### Prerequisite Validation

Run these automated checks BEFORE starting PLAN creation:

#### SDD Completeness

**Automated Checks**:
- [ ] SDD file exists (check: docs/specs/[ID]/SDD.md present)
- [ ] SDD has no `[NEEDS CLARIFICATION]` markers (check: grep returns 0 matches)
- [ ] All SDD sections non-empty (check: no sections with only placeholders)
- [ ] SDD validation checklist complete (check: all boxes marked `[x]`)
- [ ] Architecture diagram present (check: Mermaid diagram or image link exists)
- [ ] Component responsibilities defined (check: "Component:" sections exist)

**Manual Verification** (after automated checks pass):
- [ ] Do you understand the architecture without asking questions?
- [ ] Can you explain how components interact?
- [ ] Do you know what each component is responsible for?

#### Implementation Readiness

**Automated Checks**:
- [ ] Development environment ready (check: project builds successfully)
- [ ] Dependencies available (check: dependency installation succeeds)
- [ ] Test framework configured (check: can run tests even if 0 tests exist)

**Manual Verification**:
- [ ] Do you know the order in which to build components?
- [ ] Have you identified tasks that can run in parallel?
- [ ] Do you know which external services need to be mocked?

#### Minimum Information Threshold

**Automated Score**: [AUTOMATED: Calculate based on checks above]
- Critical items complete: [X/6] (must be 100%)
- Overall items complete: [X/9] (must be ≥[NEEDS CLARIFICATION: dor threshold]%)

**If threshold not met**: BLOCK with specific gaps and remediation steps

---

## Automated Validation Rules

### How Validation Works

For each "Before Creating [Document]" section:

1. **Run Automated Checks**:
   - File existence checks (`test -f path`)
   - Content checks (`grep`, `wc -l`)
   - Command execution checks (build, test setup)
   - Capture pass/fail for each check

2. **Calculate Score**:
   ```
   critical_score = critical_passed / critical_total * 100
   overall_score = total_passed / total_checks * 100
   ```

3. **Enforcement Decision**:
   ```
   if critical_score < 100:
       BLOCK("Critical prerequisites incomplete")
   elif overall_score < threshold:
       BLOCK("Overall readiness below threshold")
   else:
       PROCEED to manual verification
   ```

4. **Manual Verification**:
   - Present manual checklist to user
   - User confirms each item
   - If user marks any item incomplete:
     - Capture which items
     - Provide remediation guidance
     - BLOCK until resolved

5. **Pass**: Only if automated + manual both complete

---

## Automated Check Examples

### Check: PRD has no [NEEDS CLARIFICATION] markers

```bash
grep -c "\[NEEDS CLARIFICATION" docs/specs/001/PRD.md
# Expected: 0
# If >0: FAIL with count and line numbers
```

### Check: All PRD sections non-empty

```bash
# Extract section headers
grep "^## " docs/specs/001/PRD.md

# For each section, check next line isn't [NEEDS CLARIFICATION]
# FAIL if any section has only placeholder
```

### Check: Development environment builds

```bash
[NEEDS CLARIFICATION: build command]
# Expected: exit code 0
# If !=0: FAIL with build output
```

---

## Failure Message Format

When validation fails, display:

```
❌ Definition of Ready: BLOCKED

Document: [PRD/SDD/PLAN]
Overall: [X/Y] checks passed ([Z]%)
Critical: [A/B] checks passed ([C]%)

⛔ Critical Failures:
  • PRD has 3 [NEEDS CLARIFICATION] markers (lines 45, 67, 89)
    Impact: Cannot write SDD without complete PRD
    Fix: Complete PRD sections at lines 45, 67, 89

⚠️ Non-Critical Failures:
  • Competitive analysis not found
    Impact: May miss important context
    Fix: Create docs/research/competitors.md

📊 Automated Checks:
  ✅ PRD file exists
  ✅ PRD validation checklist complete
  ❌ PRD has no [NEEDS CLARIFICATION] markers (found 3)
  ✅ All PRD sections non-empty
  ⚠️  Competitive analysis document exists (not found)

🔧 Next Actions:
  1. Edit docs/specs/001/PRD.md lines 45, 67, 89
  2. Remove [NEEDS CLARIFICATION] markers
  3. (Optional) Create competitor analysis
  4. Re-run: /s:specify [continue]

Cannot proceed until critical items resolved.
```

---

## Configuration

### Thresholds

```yaml
thresholds:
  critical: 100  # Critical items must be 100% complete
  overall: [NEEDS CLARIFICATION: dor threshold]   # Overall items threshold (default: 85%)
```

### Validation Level

```yaml
# [NEEDS CLARIFICATION: validation level]
# Options: "strict", "balanced", "advisory"

validation_level: balanced

# strict: Block on any failure (critical or non-critical)
# balanced: Block on critical failures, warn on non-critical
# advisory: Never block, only provide warnings
```

<!-- OPTIONAL: MINIMAL_DOCS -->
### Minimal Mode

For teams using minimal documentation approach:

- Reduce checklist items by 50%
- Focus only on critical prerequisites
- Skip competitive analysis, detailed research
- Faster progression, accept more unknowns
<!-- END OPTIONAL: MINIMAL_DOCS -->

---

## Customization Guide

### Adding Custom Checks

Add automated checks for your team's needs:

```markdown
#### Custom Check Name

**Automated Checks**:
- [ ] Your custom check (check: command to run)
  ```bash
  your-validation-script.sh
  # Expected: exit code 0
  ```

**Manual Verification**:
- [ ] Your manual verification question?
```

### Language-Specific Checks

Add checks specific to your stack:

```markdown
<!-- IF: LANGUAGE_PYTHON -->
#### Python Environment

**Automated Checks**:
- [ ] Virtual environment activated (check: `which python` points to venv)
- [ ] Requirements installed (check: `pip list` includes all deps)
<!-- END IF -->
```

### Adjusting Thresholds

Change overall threshold based on project maturity:

- **New projects**: 70% (more exploration, less certainty)
- **Standard projects**: 85% (balanced)
- **Critical systems**: 95% (high certainty required)

---

## FAQ

**Q: What if I don't have all prerequisites?**
A: That's the point! DOR prevents starting too early. Gather the missing information first.

**Q: Can I skip DOR validation?**
A: Not by default. If you need to bypass (emergency), temporarily set `validation_level: advisory`.

**Q: What's the difference between automated and manual checks?**
A: Automated runs first (fast, objective). Manual runs after (judgment-based, context-aware).

**Q: Why block on critical items at 100%?**
A: Critical items are must-haves. Starting without them guarantees rework later.

**Q: Can I customize which checks are critical?**
A: Yes! Edit this file, move checks between critical/non-critical sections.
