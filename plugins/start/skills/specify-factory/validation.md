# Specify-Factory Validation Checklist

Use this checklist to validate factory artifacts before marking the spec as ready.

## Coverage Checks

- [ ] **Every requirement has a unit** — All requirements in requirements.md map to at least one unit
- [ ] **Every unit has scenarios** — Each unit in units/ has at least one scenario in scenarios/{id}/
- [ ] **P0 coverage** — Every unit has at least one P0 (critical path) scenario
- [ ] **No orphan scenarios** — Every scenario references a valid unit ID

## Unit Quality

- [ ] **Self-contained** — Each unit has goal, requirements, and constraints
- [ ] **Right-sized** — Each unit can be implemented by a single code agent in one session
- [ ] **No implementation prescription** — Units describe WHAT, not HOW
- [ ] **No acceptance criteria in units** — Those are scenarios, not unit content
- [ ] **Dependencies are explicit** — Every dependency references an existing unit ID
- [ ] **No circular dependencies** — Dependency graph is a DAG

## Scenario Quality

- [ ] **Observable behavior only** — Scenarios test through external interfaces (API, UI, CLI)
- [ ] **Codebase-grounded** — Scenarios reference actual endpoints, data shapes, conventions
- [ ] **Unambiguous pass/fail** — Each scenario's expected outcome is clear and testable
- [ ] **No implementation details** — Scenarios don't reference class names, function signatures, file paths
- [ ] **User reviewed** — All scenarios have been presented to and approved by the user

## E2E Stub Quality (when generated)

- [ ] **Framework detected** — Test framework identified from project configuration
- [ ] **One stub per scenario** — Every scenario has a corresponding executable test stub
- [ ] **Uses project conventions** — Stubs use the project's assertion library, HTTP client, naming patterns
- [ ] **Marked pending** — All stubs are skipped/pending until evaluation phase
- [ ] **Tests external interface** — Stubs make HTTP/browser/CLI calls, never import internal modules

## Manifest Quality

- [ ] **All units listed** — Every unit in units/ appears in the manifest
- [ ] **Dependencies match** — Manifest dependency clauses match unit frontmatter
- [ ] **Execution order valid** — Groups computed correctly from dependency graph
- [ ] **No duplicate unit IDs** — Each ID appears exactly once
- [ ] **YAML frontmatter valid** — threshold, max_iterations, status fields present

## Completion Criteria

✅ **Factory artifacts are ready when:**
- All checklist items pass
- User has reviewed and approved scenarios
- Manifest execution order is confirmed
- A factory loop could execute these artifacts without further clarification
