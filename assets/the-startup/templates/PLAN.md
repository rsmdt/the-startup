# Implementation Plan

## Phase 1: Context Ingestion

*[INSTRUCTION: This phase is MANDATORY. No coding until all items are checked. This instruction should not appear in final document.]*

### Required Reading
- [ ] **BRD**: `docs/specs/[ID]-[feature-name]/BRD.md` - Business context (if exists)
- [ ] **PRD**: `docs/specs/[ID]-[feature-name]/PRD.md` - Product requirements (if exists)
- [ ] **SDD**: `docs/specs/[ID]-[feature-name]/SDD.md` - Technical design

### Key Design Decisions
*[Extract from SDD]*
- [ ] [Critical decision 1]
- [ ] [Critical decision 2]

### Validation Gate
**DO NOT PROCEED until all above items are checked**

## Phase Structure

*[INSTRUCTION: Phase 1 is always context ingestion. Subsequent phases should be adapted based on feature complexity. Always end with validation phase. This instruction should not appear in final document.]*

- [ ] **Phase X**: [Descriptive Phase Name]
    - [ ] [Specific task with clear completion criteria]
    - [ ] [Another related task]
    - [ ] **Review** [specific review agent and area]
    - [ ] **Validate** [specific validation command or check]

- [ ] **Phase Y**: [Phase with Parallel Tasks]
  - [ ] [Component/Module A] [`parallel: true`]
    - [ ] [Specific implementation task]
    - [ ] [Related task for this component]
    - [ ] **Validate** [component-specific test command]
    - [ ] **Review** [if needed based on complexity]

  - [ ] [Component/Module B] [`parallel: true`]
    - [ ] [Specific implementation task]
    - [ ] [Related task for this component]
    - [ ] **Validate** [component-specific test command]
    - [ ] **Review** [if needed based on complexity]

**Final Phase: Validation & Cleanup**
- [ ] Run full test suite
- [ ] Verify acceptance criteria
- [ ] Update documentation
- [ ] **Final Check**: [Build/deploy command]

## Minimal Metadata

*[INSTRUCTION: Use only when coordination is critical]*

- `[parallel: true]` - Tasks that can run concurrently
- `[complexity: high]` - Needs careful review
