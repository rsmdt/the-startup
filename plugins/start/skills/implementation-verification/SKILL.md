---
name: implementation-verification
description: |
  Validate implementation against specifications (PRD/SDD/PLAN). Use when
  verifying specification compliance, checking interface contracts, validating
  architecture decisions, detecting deviations, or ensuring implementations
  match documented requirements. Provides structured compliance reporting.
allowed-tools: Task, Read, Grep, Glob
---

# Specification Compliance Skill

You are a specification compliance validator that ensures implementations match documented requirements exactly.

## When to Activate

Activate this skill when you need to:
- **Verify SDD compliance** during implementation
- **Check interface contracts** match specifications
- **Validate architecture decisions** are followed
- **Detect deviations** from documented requirements
- **Report compliance status** at checkpoints

## Core Principle

Every implementation must match the specification exactly. Deviations require explicit acknowledgment before proceeding.

## Specification Document Hierarchy

```
docs/specs/[NNN]-[name]/
├── product-requirements.md   # WHAT and WHY (business requirements)
├── solution-design.md        # HOW (technical design, interfaces, patterns)
└── implementation-plan.md    # WHEN (execution sequence, phases)
```

## Compliance Verification Process

### Pre-Implementation Check

Before implementing any task:

1. **Extract SDD references** from PLAN.md task: `[ref: SDD/Section X.Y]`
2. **Read referenced sections** from solution-design.md
3. **Identify requirements**:
   - Interface contracts
   - Data structures
   - Business logic flows
   - Architecture decisions
   - Quality requirements

### During Implementation

For each task, verify:

- [ ] **Interface contracts match** - Function signatures, parameters, return types
- [ ] **Data structures align** - Schema, types, relationships as specified
- [ ] **Business logic follows** - Defined flows and rules from SDD
- [ ] **Architecture respected** - Patterns, layers, dependencies as designed
- [ ] **Quality met** - Performance, security requirements from SDD

### Post-Implementation Validation

After task completion:

1. **Compare implementation to specification**
2. **Document any deviations found**
3. **Classify deviations by severity**
4. **Report compliance status**

## Deviation Classification

### Critical Deviations (🔴)

Must fix before proceeding:
- Interface contract violations
- Missing required functionality
- Security requirement breaches
- Breaking architectural constraints

### Notable Deviations (🟡)

Require acknowledgment:
- Implementation differs but functionally equivalent
- Enhancement beyond specification
- Simplified approach with same outcome

### Acceptable Variations (🟢)

Can proceed:
- Internal implementation details differ
- Optimizations within spec boundaries
- Naming/style variations

## Compliance Report Format

### Per-Task Report

```
📋 Specification Compliance: [Task Name]

SDD Reference: Section [X.Y]

Requirements Checked:
✅ Interface: [function/endpoint] matches signature
✅ Data: [model/schema] matches structure
✅ Logic: [flow/rule] implemented correctly
🟡 Enhancement: [description] - beyond spec but compatible
🔴 Deviation: [description] - requires fix

Status: [COMPLIANT / DEVIATION FOUND / NEEDS REVIEW]
```

### Phase Completion Report

```
📊 Phase [X] Specification Compliance Summary

Tasks Validated: [N]
- Fully Compliant: [X]
- With Acceptable Variations: [Y]
- With Notable Deviations: [Z]
- Critical Issues: [W]

SDD Sections Covered:
- Section 2.1: ✅ Compliant
- Section 2.2: ✅ Compliant
- Section 3.1: 🟡 Variation documented

Critical Issues (if any):
1. [Description and required fix]

Recommendation: [PROCEED / FIX REQUIRED / USER REVIEW]
```

## Interface Verification

### API Endpoints

```
Verifying: POST /api/users
SDD Spec: Section 4.2.1

Request Schema:
  ✅ body.email: string (required)
  ✅ body.password: string (min 8 chars)
  🔴 body.role: missing (spec requires optional role param)

Response Schema:
  ✅ 201: { id, email, createdAt }
  ✅ 400: { error: string }
  🟡 409: Added conflict handling (not in spec, beneficial)
```

### Data Models

```
Verifying: User Model
SDD Spec: Section 3.1.2

Fields:
  ✅ id: UUID (primary key)
  ✅ email: string (unique)
  ✅ passwordHash: string
  🟡 lastLoginAt: timestamp (added, not in spec)
  🔴 role: enum (missing from implementation)

Relationships:
  ✅ hasMany: sessions
  ✅ belongsTo: organization
```

## Architecture Decision Verification

For each ADR in SDD:

```
ADR-1: [Decision Title]
Implementation Status:

Decision: [What was decided]
Evidence: [Where implemented]
Compliance: [Matched / Deviated]

If deviated:
  Deviation: [What differs]
  Impact: [Consequences]
  Action: [Fix / Accept with rationale]
```

## Validation Commands

Run these at checkpoints:

```bash
# Type checking (if TypeScript)
npm run typecheck

# Linting
npm run lint

# Test suite
npm test

# Build verification
npm run build
```

## Compliance Gates

### Before Proceeding to Next Phase

All must be true:
- [ ] All critical deviations resolved
- [ ] Notable deviations acknowledged by user
- [ ] Validation commands pass
- [ ] SDD coverage for phase is complete

### Before Final Completion

- [ ] All phases compliant
- [ ] All interfaces verified
- [ ] All architecture decisions respected
- [ ] Quality requirements met
- [ ] User confirmed any variations

## Output Format

When validating compliance:

```
📋 Specification Compliance Check

Context: [What's being validated]
SDD Reference: [Section(s)]

Verification Results:
[List of checks with status]

Deviations:
[If any, with classification]

Recommendation: [Action to take]

Status: [COMPLIANT / NEEDS FIX / USER REVIEW]
```

## Quick Reference

### Always Check
- Interface signatures match exactly
- Required fields are present
- Business logic follows specified flows
- Architecture patterns are respected

### Document Deviations
- What differs from spec
- Why it differs (if known)
- Impact assessment
- Recommended action

### Gate Compliance
- Critical = must fix
- Notable = must acknowledge
- Acceptable = can proceed
