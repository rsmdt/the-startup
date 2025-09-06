# Implementation Plan

## Phase 1: Context Ingestion

*[INSTRUCTION: This phase is MANDATORY. No coding until all items are checked. This instruction should not appear in final document.]*

### Required Reading
- [ ] **BRD**: `docs/specs/[ID]-[feature-name]/BRD.md` - Business Requirements (if exists)
- [ ] **PRD**: `docs/specs/[ID]-[feature-name]/PRD.md` - Product Requirements (if exists)
- [ ] **SDD**: `docs/specs/[ID]-[feature-name]/SDD.md` - Solution Design

### Key Design Decisions
*[Extract from SDD]*
- [ ] [Critical decision 1]
- [ ] [Critical decision 2]

### Validation Gate
**DO NOT PROCEED until all above items are checked**

## Phase Structure

*[INSTRUCTION: Phase 1 is always context ingestion. Subsequent phases should be adapted based on feature complexity. Always end with validation phase. This instruction should not appear in final document.]*

- [ ] **Phase X**: [Descriptive Phase Name]
    - [ ] [Specific task with clear completion criteria] [activity: ...]
    - [ ] [Another related task] [activity: ...]
    - [ ] **Review** [specific review agent and area] [activity: ...]
    - [ ] **Validate** [specific validation command or check] [activity: ...]

- [ ] **Phase Y**: [Phase with Parallel Tasks] [complexity: high]
  - [ ] [Component/Module A] [`parallel: true`]
    - [ ] [Specific implementation task] [activity: ...]
    - [ ] [Related task for this component] [activity: ...]
    - [ ] **Validate** [component-specific test command] [activity: ...]
    - [ ] **Review** [if needed based on complexity] [activity: ...]

  - [ ] [Component/Module B] [`parallel: true`]
    - [ ] [Specific implementation task] [activity: ...]
    - [ ] [Related task for this component] [activity: ...]
    - [ ] **Validate** [component-specific test command] [activity: ...]
    - [ ] **Review** [if needed based on complexity] [activity: ...]

**Final Phase: Validation & Cleanup**
- [ ] Run full test suite [activity: ...]
- [ ] Verify acceptance criteria [activity: ...]
- [ ] Update documentation [activity: ...]
- [ ] **Final Check**: [Build/deploy command] [activity: ...]

## Optional Metadata

*[INSTRUCTION: Use on tasks when coordination is critical. Do not add or invent additional metadata]*

- `[parallel: true]` - Tasks that can run concurrently
- `[complexity: high/medium/low]` - Needs careful review
- `[activity: api-design [, business-logic, ...]]` - Activity hint for agent selection
