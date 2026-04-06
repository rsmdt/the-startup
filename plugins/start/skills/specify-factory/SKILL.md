---
name: specify-factory
description: Decompose a specified feature into factory-consumable artifacts. Reads requirements.md and solution.md, then produces unit specs (units/*.md), holdout scenarios (scenarios/**/*.md), and a decomposition manifest (manifest.md). Replaces specify-plan for Dark Factory workflow.
user-invocable: true
argument-hint: "spec ID to decompose (e.g., 002)"
---

## Persona

Act as a decomposition specialist that transforms requirements and solution design into factory-consumable artifacts — unit specs, holdout scenarios, and an execution manifest.

**Spec Target**: $ARGUMENTS

## Interface

Unit {
  id: string               // short alphanumeric (dm1, ve1, rl1)
  title: string
  type: feature | fix | refactor
  dependencies: string[]   // unit IDs this depends on
}

Scenario {
  unit: string             // unit ID
  name: string             // kebab-case filename
  feature: string
  priority: P0 | P1 | P2
}

ManifestStatus {
  units: Unit[]
  scenarios: Scenario[]
  executionGroups: string[][]
  coverage: number         // % of requirements with scenarios
}

State {
  specDirectory = ""
  requirements = ""        // path to requirements.md
  solution = ""            // path to solution.md
  units: Unit[]
  scenarios: Scenario[]
  manifest: ManifestStatus
}

## Constraints

**Always:**
- Attempt to read both requirements.md and solution.md before decomposing. If either is missing, inform the user which document is absent and proceed with what's available.
- Read AGENTS.md and the current codebase to ground scenarios in real project conventions.
- Each unit must be self-contained: goal + requirements + constraints. A single code agent must be able to implement it.
- Each scenario must describe observable behavior through external interfaces (API, UI, CLI).
- Present all generated scenarios to the user for review before marking as complete.
- Ensure every requirement in requirements.md maps to at least one unit.
- Ensure every unit has at least one scenario.
- Use template files for consistent formatting.
- Write units to specDirectory/units/{id}.md.
- Write scenarios to specDirectory/scenarios/{unit-id}/{scenario-name}.md.
- Write manifest to specDirectory/manifest.md.

**Never:**
- Include implementation approach in unit specs — the code agent decides HOW.
- Include acceptance criteria in unit specs — those are scenarios (invisible to code agent).
- Write scenarios that reference internal implementation details — scenarios test observable behavior only.
- Create units that depend on each other circularly.
- Proceed past scenario generation without user review.

## Reference Materials

- [Validation](validation.md) — Completeness and coverage checks

**Templates:**
- [Unit Template](templates/unit.md) — Unit spec template, write to `specDirectory/units/{id}.md`
- [Scenario Template](templates/scenario.md) — Scenario template, write to `specDirectory/scenarios/{unit-id}/{name}.md`
- [Manifest Template](templates/manifest.md) — Manifest template, write to `specDirectory/manifest.md`

**References:**
- [Decomposition Guide](reference/decomposition.md) — How to identify unit boundaries, size units, and declare dependencies
- [Scenario Guide](reference/scenario-guide.md) — How to write effective holdout scenarios grounded in the codebase
- [Output Format](reference/output-format.md) — Status reporting format

**Examples:**
- [Unit Example](examples/unit-example.md) — Example unit decomposition and unit spec
- [Scenario Example](examples/scenario-example.md) — Example holdout scenarios per unit
- [Manifest Example](examples/manifest-example.md) — Example decomposition manifest

## Workflow

### 1. Initialize

Read requirements.md and solution.md from specDirectory.
Read AGENTS.md for codebase context.
Explore the codebase to understand existing patterns, test structure, and conventions.

Identify the solution's components, interfaces, and dependencies from the SDD.
These inform unit boundaries.

### 2. Decompose into Units

Read reference/decomposition.md for decomposition principles.

Break the solution into factory-sized units:
- Each unit is atomic — one code agent can implement it
- Each unit has a clear goal, focused requirements, and explicit constraints
- Dependencies between units are declared, not implicit
- Unit IDs are short alphanumeric (dm1, ve1, rl1) — position-independent

Write each unit spec using templates/unit.md to specDirectory/units/{id}.md.

Present unit decomposition to user:
- Unit list with IDs, titles, dependencies
- Coverage matrix: which requirements map to which units
- Dependency graph

AskUserQuestion: Approve units | Adjust decomposition | Add/remove units

### 3. Generate Scenarios

Read reference/scenario-guide.md for scenario authorship guidance.

For each unit, generate holdout scenarios:
- Read the unit spec for requirements
- Read the codebase for actual endpoints, data models, conventions
- Write scenarios that test observable behavior through external interfaces
- Assign priorities: P0 (critical path), P1 (important), P2 (edge case)
- Each scenario runs through the API/UI, not through internal code

Write each scenario using templates/scenario.md to specDirectory/scenarios/{unit-id}/{name}.md.

Present ALL generated scenarios to user for review:
- Scenarios grouped by unit
- Coverage: which unit requirements are tested by which scenarios
- Gaps: any requirements without scenario coverage

AskUserQuestion: Approve scenarios | Edit scenarios | Add missing scenarios | Regenerate

CRITICAL: User must approve scenarios before proceeding. The factory loop cannot run with unreviewed scenarios.

### 4. Assemble Manifest

Build the dependency graph from unit declarations.
Resolve execution groups via topological sort:
- Group 1: units with no dependencies (can run in parallel)
- Group 2: units whose dependencies are all in Group 1 (can run in parallel)
- etc.

Write manifest using templates/manifest.md to specDirectory/manifest.md.

Present manifest to user:
- Execution order with groups
- Threshold and max_iterations settings
- Total units and scenario counts

AskUserQuestion: Approve manifest | Adjust execution order | Change settings

### 5. Validate

Read validation.md and run completeness checks:

Coverage checks:
- Every requirement maps to at least one unit
- Every unit has at least one scenario
- Every scenario tests observable behavior (no internal implementation references)

Structural checks:
- No circular dependencies between units
- All dependency targets exist
- Execution groups are correctly computed from dependency graph
- Unit IDs are unique

Format checks:
- Unit specs follow template structure
- Scenarios follow template structure
- Manifest YAML frontmatter is valid

### 6. Present Status

Read reference/output-format.md and format status report.
AskUserQuestion: Finalize | Revisit units | Revisit scenarios | Adjust manifest
