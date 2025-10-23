# Definition of Done

<!--
This file defines completion criteria that must be met AFTER creating each document.

AUTOMATED VALIDATION:
- Structural checks (sections complete, no placeholders)
- Logical flow checks (SCQA validation)
- Coverage checks (MECE validation)
- Consistency checks (cross-document validation)

USAGE:
- After PRD: Validate "PRD Completion" section
- After SDD: Validate "SDD Completion" section
- After PLAN: Validate "PLAN Completion" section
-->

---

## PRD Completion

### Document Completeness

**Automated Checks**:
- [ ] All required sections present (check: grep for section headers matches template)
- [ ] No `[NEEDS CLARIFICATION]` markers remaining (check: grep returns 0)
- [ ] Validation checklist at top complete (check: all items marked `[x]`)
- [ ] Minimum word count met (check: wc -w ≥500 for PRD)
- [ ] All MoSCoW categories addressed (check: Must/Should/Could/Won't sections exist)

**Manual Verification** (after automated checks pass):
- [ ] Is every section meaningful (not just placeholder text)?
- [ ] Would a new team member understand this PRD?

### SCQA Validation

<!-- OPTIONAL: SCQA_ALL_DOCS -->
**Situation-Complication-Question-Answer** logical flow validation:

**Automated Checks**:
- [ ] **Situation**: "Context" or "Background" section exists and >100 words
- [ ] **Complication**: "Problem Statement" section exists and identifies pain points
- [ ] **Question**: "Goals" or "Objectives" section exists
- [ ] **Answer**: "Solution Overview" or "Value Proposition" section exists

**Logic Flow Check**:
```python
# Automated: Check that sections follow logical order
sections = extract_sections(prd)
required_order = ["Context/Background", "Problem", "Solution"]
validate_order(sections, required_order)
# FAIL if sections out of order or missing
```

**Manual Verification**:
- [ ] Does the Context flow naturally into the Problem?
- [ ] Does the Problem make the reader ask "What should we do?"?
- [ ] Does the Solution answer that question?
- [ ] Is there a clear cause-and-effect chain?
<!-- END OPTIONAL: SCQA_ALL_DOCS -->

<!-- OPTIONAL: SCQA_PRD_ONLY -->
**SCQA validation for PRD only** (see above for checks)
<!-- END OPTIONAL: SCQA_PRD_ONLY -->

### MECE Validation

<!-- OPTIONAL: MECE_ALL_DOCS -->
**Mutually Exclusive, Collectively Exhaustive** validation:

**Automated Checks - Mutually Exclusive (no overlap)**:
- [ ] No duplicate user stories (check: fuzzy match story descriptions, flag >80% similarity)
- [ ] No contradicting requirements (check: parse requirements, flag conflicting statements like "must be fast" and "performance not critical")
- [ ] Feature list has no redundancy (check: TF-IDF similarity between features <0.7)

**Automated Checks - Collectively Exhaustive (nothing missing)**:
- [ ] All user personas have at least 1 user journey (check: persona count ≤ journey count)
- [ ] Every feature has acceptance criteria (check: each "Feature:" section has "Acceptance Criteria:" subsection)
- [ ] Every metric has tracking method defined (check: metrics table has "How to Measure" column)
- [ ] All edge cases covered (check: "Edge Cases" or "Error Scenarios" section exists)

**Manual Verification**:
- [ ] Are there any overlapping features that should be combined?
- [ ] Is anything missing from the user journey?
- [ ] Have you covered all critical user paths?
<!-- END OPTIONAL: MECE_ALL_DOCS -->

<!-- OPTIONAL: MECE_PRD_ONLY -->
**MECE validation for PRD only** (see above for checks)
<!-- END OPTIONAL: MECE_PRD_ONLY -->

### Consistency Validation

<!-- IF: CONSISTENCY_AUTOMATED -->
**Automated Cross-Section Consistency**:
- [ ] Personas mentioned in journeys actually defined (check: extract persona names from journeys, verify in persona section)
- [ ] Metrics align with goals (check: each goal has corresponding metric)
- [ ] Features support stated objectives (check: features reference objectives in rationale)

**Manual Verification**:
- [ ] Do the personas align with the target market?
- [ ] Do the features actually solve the stated problem?
<!-- END IF -->

### Minimum Quality Threshold

**Automated Score**: [AUTOMATED: Calculate based on checks above]
- Critical items complete: [X/Y] (must be 100%)
- Overall items complete: [X/Z] (must be ≥[NEEDS CLARIFICATION: dod threshold]%)

**If threshold not met**: BLOCK with specific failures and remediation steps

---

## SDD Completion

### Document Completeness

**Automated Checks**:
- [ ] All required sections present (check: grep for section headers matches template)
- [ ] No `[NEEDS CLARIFICATION]` markers remaining (check: grep returns 0)
- [ ] Validation checklist at top complete (check: all items marked `[x]`)
- [ ] Architecture diagram present (check: mermaid code block or image link exists)
- [ ] Minimum word count met (check: wc -w ≥800 for SDD)

**Manual Verification**:
- [ ] Would a developer understand how to implement from this?
- [ ] Are all architecture decisions justified?

### SCQA Validation

<!-- OPTIONAL: SCQA_ALL_DOCS -->
**Situation-Complication-Question-Answer** for technical decisions:

**Automated Checks**:
- [ ] **Situation**: "Constraints" or "Requirements" section exists
- [ ] **Complication**: Technical challenges identified
- [ ] **Question**: "How will we implement?" implied in design sections
- [ ] **Answer**: Architecture decisions documented with rationale

**Decision Rationale Check**:
```python
# Automated: Check that architecture decisions have rationale
decisions = extract_decisions(sdd)
for decision in decisions:
    if not has_rationale(decision):
        FAIL(f"Decision '{decision}' missing rationale")
```

**Manual Verification**:
- [ ] Is it clear WHY each architecture decision was made?
- [ ] Are trade-offs explained (what was considered and rejected)?
<!-- END OPTIONAL: SCQA_ALL_DOCS -->

### MECE Validation

<!-- OPTIONAL: MECE_ALL_DOCS -->
**Mutually Exclusive, Collectively Exhaustive** for architecture:

**Automated Checks - Mutually Exclusive (no component overlap)**:
- [ ] Components have distinct responsibilities (check: extract component descriptions, flag >70% similarity)
- [ ] No duplicate APIs or interfaces (check: endpoint paths unique)
- [ ] Database tables have clear ownership (check: each table assigned to one component)

**Automated Checks - Collectively Exhaustive (all PRD requirements covered)**:
```python
# Automated: Parse PRD requirements and map to SDD components
prd_requirements = extract_requirements("docs/specs/[ID]/PRD.md")
sdd_coverage = extract_requirement_mappings("docs/specs/[ID]/SDD.md")

for req in prd_requirements:
    if req not in sdd_coverage:
        FAIL(f"PRD requirement '{req}' not addressed in SDD")
```
- [ ] Every PRD requirement mapped to SDD component
- [ ] All interfaces between components defined
- [ ] All data flows documented

**Manual Verification**:
- [ ] Is there component responsibility overlap that should be resolved?
- [ ] Did you miss any PRD requirements in the design?
- [ ] Are all integration points covered?
<!-- END OPTIONAL: MECE_ALL_DOCS -->

### Consistency Validation

<!-- IF: CONSISTENCY_AUTOMATED -->
**Automated Cross-Document Consistency**:
```python
# Check PRD requirements match SDD implementation
prd_features = extract_features("docs/specs/[ID]/PRD.md")
sdd_components = extract_components("docs/specs/[ID]/SDD.md")

# Every feature should map to at least one component
for feature in prd_features:
    if not has_component_mapping(feature, sdd_components):
        FAIL(f"PRD feature '{feature}' not implemented in SDD")
```
- [ ] PRD requirements align with SDD design
- [ ] Technical constraints match PRD assumptions
- [ ] Success metrics have implementation plan

**Manual Verification**:
- [ ] Does the architecture actually solve the problem from PRD?
- [ ] Are there any PRD requirements that can't be met with this design?
<!-- END IF -->

### Minimum Quality Threshold

**Automated Score**: [AUTOMATED: Calculate based on checks above]
- Critical items complete: [X/Y] (must be 100%)
- Overall items complete: [X/Z] (must be ≥[NEEDS CLARIFICATION: dod threshold]%)

**If threshold not met**: BLOCK with specific failures and remediation steps

---

## PLAN Completion

### Document Completeness

**Automated Checks**:
- [ ] All required sections present (check: grep for section headers matches template)
- [ ] No `[NEEDS CLARIFICATION]` markers remaining (check: grep returns 0)
- [ ] Validation checklist at top complete (check: all items marked `[x]`)
- [ ] All phases follow TDD structure (check: each phase has Prime→Test→Implement→Validate)
- [ ] Dependencies tagged (check: `[ref:` markers present linking to SDD sections)
- [ ] Integration phase defined (check: final phase includes E2E tests)

**Manual Verification**:
- [ ] Could a developer execute this plan independently?
- [ ] Are task descriptions clear and actionable?

### MECE Validation

<!-- OPTIONAL: MECE_ALL_DOCS -->
**Mutually Exclusive, Collectively Exhaustive** for implementation:

**Automated Checks - Mutually Exclusive (no duplicate work)**:
```python
# Check for task overlap
tasks = extract_tasks("docs/specs/[ID]/PLAN.md")
for i, task1 in enumerate(tasks):
    for task2 in tasks[i+1:]:
        similarity = calculate_similarity(task1.description, task2.description)
        if similarity > 0.75:
            FAIL(f"Tasks '{task1.id}' and '{task2.id}' appear to overlap")
```
- [ ] No duplicate tasks (check: similarity score between tasks <75%)
- [ ] Clear responsibility boundaries (check: each task has distinct `[activity: ]` hint)

**Automated Checks - Collectively Exhaustive (all SDD components covered)**:
```python
# Check SDD component coverage
sdd_components = extract_components("docs/specs/[ID]/SDD.md")
plan_tasks = extract_tasks("docs/specs/[ID]/PLAN.md")

for component in sdd_components:
    if not has_implementation_tasks(component, plan_tasks):
        FAIL(f"SDD component '{component}' has no implementation tasks")
```
- [ ] Every SDD component has implementation tasks
- [ ] All SDD interfaces have integration tasks
- [ ] E2E tests cover all user journeys from PRD

**Manual Verification**:
- [ ] Is any work duplicated across tasks?
- [ ] Did you miss implementing any SDD components?
- [ ] Do the tasks cover the full user journey end-to-end?
<!-- END OPTIONAL: MECE_ALL_DOCS -->

### Consistency Validation

<!-- IF: CONSISTENCY_AUTOMATED -->
**Automated Cross-Document Consistency**:
```python
# Verify PLAN tasks implement SDD design
sdd_components = extract_components("docs/specs/[ID]/SDD.md")
plan_tasks = extract_tasks("docs/specs/[ID]/PLAN.md")

# Check that implementation follows architecture
for task in plan_tasks:
    if task.references_sdd:
        sdd_ref = parse_sdd_reference(task)
        if not sdd_ref_exists(sdd_ref, sdd_components):
            FAIL(f"Task '{task.id}' references non-existent SDD section")
```
- [ ] Tasks reference actual SDD sections (check: `[ref: SDD/Section X.Y]` points to existing section)
- [ ] Implementation order respects dependencies (check: no circular dependencies)
- [ ] Test tasks align with PRD acceptance criteria (check: test descriptions reference PRD criteria)

**Manual Verification**:
- [ ] Does the implementation plan match the SDD architecture?
- [ ] Are tests verifying the right acceptance criteria from PRD?
<!-- END IF -->

### Minimum Quality Threshold

**Automated Score**: [AUTOMATED: Calculate based on checks above]
- Critical items complete: [X/Y] (must be 100%)
- Overall items complete: [X/Z] (must be ≥[NEEDS CLARIFICATION: dod threshold]%)

**If threshold not met**: BLOCK with specific failures and remediation steps

---

## Automated Validation Rules

### Validation Execution Flow

For each "Document Completion" section:

1. **Run Automated Checks** (structural, SCQA, MECE, consistency)
2. **Calculate Scores** (critical and overall percentages)
3. **Enforce Decision** (block if below threshold)
4. **Manual Verification** (present checklist to user)
5. **Final Decision** (pass only if automated + manual both complete)

### Failure Message Format

```
❌ Definition of Done: BLOCKED

Document: SDD
Overall: 18/22 checks passed (82%)
Critical: 5/6 checks passed (83%)

⛔ Critical Failures:
  • MECE: PRD requirement "User authentication" not covered in SDD
    Impact: Feature will be missing from implementation
    Fix: Add authentication component to SDD Section 4

⚠️ SCQA Failures:
  • Architecture decision "Use PostgreSQL" missing rationale
    Impact: Team may not understand why this choice was made
    Fix: Add rationale to SDD Section 5.2

⚠️ MECE Failures:
  • Components "AuthService" and "UserService" have 78% description overlap
    Impact: Unclear responsibility boundaries
    Fix: Clarify distinct responsibilities or merge components

📊 Automated Checks:
  ✅ All sections present
  ✅ No [NEEDS CLARIFICATION] markers
  ✅ Validation checklist complete
  ✅ Architecture diagram present
  ❌ All PRD requirements covered (missing 1)
  ⚠️  Architecture decisions have rationale (1 missing)
  ⚠️  Components have distinct responsibilities (1 overlap)

🔧 Next Actions:
  1. Add authentication component to SDD Section 4
  2. Add PostgreSQL rationale to Section 5.2
  3. Clarify AuthService vs UserService boundaries
  4. Re-run: /s:specify [validate]

Cannot proceed until critical items resolved.
```

---

## Configuration

### Thresholds

```yaml
thresholds:
  critical: 100  # Critical items must be 100% complete
  overall: [NEEDS CLARIFICATION: dod threshold]   # Overall items threshold (default: 85%)
```

### Validation Features

```yaml
# [NEEDS CLARIFICATION: validation features]
features:
  scqa: [NEEDS CLARIFICATION: enable scqa]     # Enable SCQA validation (true/false)
  scqa_scope: [NEEDS CLARIFICATION: scqa scope] # "all", "prd-sdd", "prd-only"
  mece: [NEEDS CLARIFICATION: enable mece]     # Enable MECE validation (true/false)
  mece_scope: [NEEDS CLARIFICATION: mece scope] # "all", "prd-sdd", "prd-only"
  consistency: [NEEDS CLARIFICATION: consistency check] # "automated", "manual", "disabled"
```

### Validation Level

```yaml
# [NEEDS CLARIFICATION: validation level]
validation_level: balanced  # "strict", "balanced", "advisory"
```

---

## Customization Guide

### Adjusting Similarity Thresholds

Change overlap detection sensitivity:

```python
# For MECE duplicate detection
similarity_threshold = 0.75  # Default
# Lower (0.6) = more strict, flag more potential duplicates
# Higher (0.85) = more lenient, only flag obvious duplicates
```

### Adding Custom SCQA Checks

```markdown
**Automated Checks**:
- [ ] Custom SCQA element (check: your validation logic)
```

### Language-Specific Quality Checks

Add checks for your stack:

```markdown
<!-- IF: LANGUAGE_GO -->
#### Go-Specific Checks
- [ ] All exported types documented (check: `go doc` output coverage)
- [ ] Error handling follows Go conventions
<!-- END IF -->
```

---

## FAQ

**Q: What if SCQA/MECE checks fail but document seems fine?**
A: Automated checks can have false positives. Review the failure, and if it's incorrect, you can temporarily disable that specific check.

**Q: How strict are the similarity thresholds?**
A: Default 75% catches obvious duplicates. Adjust based on your tolerance for overlap.

**Q: Can I skip SCQA/MECE validation?**
A: Yes, set `features.scqa: false` and `features.mece: false` in configuration.

**Q: What's the difference between critical and overall thresholds?**
A: Critical = must-have requirements (always 100%). Overall = all checks combined (default 85%, adjustable).

**Q: How does consistency validation work across documents?**
A: It parses all docs, extracts key entities (requirements, components, tasks), and verifies mappings exist.
