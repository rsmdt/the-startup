---
id: {{UNIT_ID}}
title: {{UNIT_TITLE}}
type: {{feature | fix | refactor}}
dependencies: [{{DEPENDENCY_IDS}}]
---
# {{UNIT_TITLE}}

## Goal
{{One-line description of the user-observable capability this unit delivers end-to-end. Name a behavior, not a layer.}}

## Requirements
{{Bulleted list of what this unit must do. Each requirement is testable through the unit's external interface.}}

## Constraints
{{Bulleted list of constraints: patterns to follow, dependencies to avoid, coverage expectations.}}
- This unit is a **vertical slice**: it delivers a complete, externally observable capability. Its holdout scenarios must pass when only this unit (plus its declared dependencies) is built. If a requirement here can't be tested without a sibling unit, it belongs in that unit — re-slice.
- Follow existing patterns (see AGENTS.md)
- Full test coverage for all requirements

## Interfaces
{{Optional. Only include if this unit defines or consumes interfaces shared with other units. If this is the first slice to need a shared entity/schema, it OWNS that foundation here; later units depend on this unit and extend it.}}
