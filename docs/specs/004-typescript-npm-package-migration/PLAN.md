# Implementation Plan

## Validation Checklist
- [ ] Context Ingestion section complete with all required specs
- [ ] Implementation phases logically organized
- [ ] Each phase starts with test definition (TDD approach)
- [ ] Dependencies between phases identified
- [ ] Parallel execution marked where applicable
- [ ] Multi-component coordination identified (if applicable)
- [ ] Final validation phase included
- [ ] No placeholder content remains

## Specification Compliance Guidelines

### How to Ensure Specification Adherence

1. **Before Each Phase**: Complete the Pre-Implementation Specification Gate
2. **During Implementation**: Reference specific SDD sections in each task
3. **After Each Task**: Run Specification Compliance checks
4. **Phase Completion**: Verify all specification requirements are met

### Deviation Protocol

If implementation cannot follow specification exactly:
1. Document the deviation and reason
2. Get approval before proceeding
3. Update SDD if the deviation is an improvement
4. Never deviate without documentation

## Metadata Reference

- `[parallel: true]` - Tasks that can run concurrently
- `[component: component-name]` - For multi-component features
- `[ref: document/section; lines: 1, 2-3]` - Links to specifications, patterns, or interfaces and (if applicable) line(s)
- `[activity: type]` - Activity hint for specialist agent selection

---

## Context Priming

*GATE: You MUST fully read all files mentioned in this section before starting any implementation.*

**CRITICAL: Complete T000 (Pre-Implementation Discovery) FIRST** - This phase analyzes the existing Go codebase to extract undocumented behaviors, edge cases, and implementation patterns that must be preserved in the TypeScript migration.

**Specification**:

- `docs/specs/004-typescript-npm-package-migration/PRD.md` - Product Requirements
- `docs/specs/004-typescript-npm-package-migration/SDD.md` - Solution Design
- Existing Go source files (analyzed in T000) - Current implementation patterns and edge cases

**Key Design Decisions**:

1. **TUI Framework: Ink** (production-ready, 1.48M DL/week vs OpenTUI alpha v0.1.25)
2. **Build Tool: tsup** (TypeScript-first, dual ESM/CJS, automatic .d.ts generation)
3. **Lock File: SHA-256 checksums with backward compatibility** (enables idempotent reinstalls, auto-migrates legacy format)
4. **Statusline: Shell scripts (.sh/.ps1)** (achieves <10ms requirement, Node.js startup is 50-200ms)
5. **Asset Strategy: File copying** (not bundling - keeps package small at <5MB, files human-readable)
6. **Settings.json: Deep merge with backup** (preserve user customizations, atomic operation with rollback)
7. **Package Architecture: Layered (CLI → UI → Core → Lib)** (clean separation, TUI replaceable without touching business logic)

**Implementation Context**:

- Commands to run:
  ```bash
  # Development
  npm install                    # Install dependencies
  npm run dev                    # tsup watch mode
  npm run build                  # Build project (tsup outputs to dist/)

  # Testing
  npm test                       # vitest run
  npm run test:watch             # vitest watch
  npm run test:coverage          # vitest --coverage
  npm run test:integration       # Full install flow test

  # Code Quality
  npm run lint                   # eslint src/
  npm run typecheck              # tsc --noEmit
  npm run format                 # prettier --write src/

  # Shell Script Testing
  echo '{"model":"sonnet"}' | ./bin/statusline.sh    # Test Unix statusline
  time (echo '{}' | ./bin/statusline.sh)             # Benchmark statusline
  ```

- Patterns to follow:
  - TDD: Red-Green-Refactor (CLAUDE.md requirement)
  - Layered architecture: CLI → UI → Core → Lib (SDD Solution Strategy)
  - Ink component pattern (SDD Cross-Cutting Concepts, lines 1073-1086)
  - Deep merge algorithm for settings.json (SDD Complex Logic, lines 923-954)
  - Lock file backward compatibility with auto-migration (SDD Complex Logic, lines 956-1023)

- Interfaces to implement:
  - CLI Commands (Commander.js): install, uninstall, init, spec, statusline
  - Lock File format: v2 with SHA-256 checksums, v1 backward compat (SDD lines 513-549)
  - Settings.json hook structure with placeholder replacement (SDD lines 551-567)
  - Statusline JSON input/output format (SDD lines 218-222)
  - TOML output for spec --read command (PRD line 206)

---

## Implementation Phases

- [x] T000 Pre-Implementation Discovery & Codebase Analysis

    - [x] T000.1 Analyze Existing Go Implementation
        - [x] T000.1.1 Read main.go and go.mod to understand current architecture `[ref: SDD; lines: 56-67]` `[activity: requirements-analysis]`
        - [x] T000.1.2 Read cmd/*.go to extract CLI patterns and error messages `[ref: SDD; lines: 94-117]` `[activity: requirements-analysis]`
        - [x] T000.1.3 Read internal/ui/*.go to understand BubbleTea state machine patterns `[ref: SDD; lines: 119-135]` `[activity: requirements-analysis]`
        - [x] T000.1.4 Read internal/installer/*.go to extract business logic edge cases `[ref: SDD; lines: 137-151]` `[activity: requirements-analysis]`
        - [x] T000.1.5 Document undocumented behaviors and edge cases found in Go code `[activity: requirements-analysis]`

    - [x] T000.2 Validate Migration Assumptions
        - [x] T000.2.1 Verify all Go CLI flags have TypeScript equivalents `[ref: PRD; line: 344]` `[activity: requirements-analysis]`
        - [x] T000.2.2 Extract current error message texts for consistency `[activity: requirements-analysis]`
        - [x] T000.2.3 Identify Go-specific patterns that need Ink equivalents (BubbleTea → Ink) `[activity: requirements-analysis]`
        - [x] T000.2.4 Document TUI state transitions from Go implementation `[activity: requirements-analysis]`

    - [x] T000.3 Validate
        - [x] T000.3.1 All SDD Implementation Context files reviewed `[ref: SDD; lines: 56-151]` `[activity: review-code]`
        - [x] T000.3.2 Edge cases and gotchas documented for implementation reference `[activity: review-code]`
        - [x] T000.3.3 Migration risks identified and mitigation strategies documented `[activity: review-code]`

- [x] T001 Project Foundation & Type System

    - [x] T001.1 Prime Context
        - [x] T001.1.1 Read SDD Constraints section `[ref: SDD; lines: 23-50]`
        - [x] T001.1.2 Read SDD Build Constraints `[ref: SDD; lines: 41-45]`
        - [x] T001.1.3 Read SDD Type definitions `[ref: SDD; lines: 513-653]`
        - [x] T001.1.4 Read SDD Project Commands `[ref: SDD; lines: 298-341]`

    - [x] T001.2 Write Tests
        - [x] T001.2.1 Test vitest can discover test files `[activity: test-execution]`
        - [x] T001.2.2 Test tsup builds dual ESM/CJS output `[activity: test-execution]`
        - [x] T001.2.3 Test TypeScript strict mode compilation `[activity: test-execution]`

    - [x] T001.3 Implement Project Configuration
        - [x] T001.3.1 Create package.json with dependencies and bin field `[ref: PRD; lines: 155-161]` `[activity: component-development]`
        - [x] T001.3.2 Create tsconfig.json with strict mode enabled `[ref: SDD; lines: 42]` `[activity: component-development]`
        - [x] T001.3.3 Create tsup.config.ts for dual ESM/CJS build `[ref: SDD; lines: 43]` `[activity: component-development]`
        - [x] T001.3.4 Create vitest.config.ts with 90% coverage threshold `[ref: SDD; lines: 1241]` `[activity: component-development]`
        - [x] T001.3.5 Create directory structure: src/, tests/, bin/, assets/ `[ref: SDD; lines: 432-502]` `[activity: component-development]`

    - [x] T001.4 Implement Type Definitions
        - [x] T001.4.1 Create src/core/types/lock.ts with FileEntry and LockFile interfaces `[ref: SDD; lines: 514-549]` `[activity: domain-modeling]`
        - [x] T001.4.2 Create src/core/types/settings.ts with ClaudeSettings and PlaceholderMap interfaces `[ref: SDD; lines: 551-567]` `[activity: domain-modeling]`
        - [x] T001.4.3 Create src/core/types/config.ts with command option interfaces `[ref: SDD; lines: 588-645]` `[activity: domain-modeling]`
        - [x] T001.4.4 Create SpecNumbering interface in spec types `[ref: SDD; lines: 647-653]` `[activity: domain-modeling]`

    - [x] T001.5 Implement Test Infrastructure
        - [x] T001.5.1 Create tests/shared/testUtils.ts with common helpers `[activity: test-execution]`
        - [x] T001.5.2 Create tests/core/__mocks__/fs.ts for file system mocking `[activity: test-execution]`
        - [x] T001.5.3 Create test fixtures (sample lock files, settings.json) `[activity: test-execution]`

    - [x] T001.6 Validate
        - [x] T001.6.1 Run typecheck: npm run typecheck passes `[activity: run-tests]`
        - [x] T001.6.2 Run build: npm run build produces dist/ output `[activity: run-tests]`
        - [x] T001.6.3 Verify dual module output (ESM + CJS + .d.ts) `[activity: run-tests]`
        - [x] T001.6.4 Verify strict TypeScript compilation (no implicit any) `[ref: SDD; line: 42]` `[activity: review-code]`

- [x] T002 Core Foundation & Shell Scripts

    - [x] T002.1 LockManager Implementation `[parallel: true]` `[component: LockManager]`
        - [x] T002.1.1 Prime Context
            - [x] T002.1.1.1 Read Lock File algorithm `[ref: SDD; lines: 956-1023]`
            - [x] T002.1.1.2 Read Lock File structure `[ref: SDD; lines: 514-549]`
        - [x] T002.1.2 Write Tests
            - [x] T002.1.2.1 Test reading legacy lock file (string[] format) `[activity: test-execution]`
            - [x] T002.1.2.2 Test reading new lock file (FileEntry[] with checksums) `[activity: test-execution]`
            - [x] T002.1.2.3 Test auto-migration from v1 to v2 format `[activity: test-execution]`
            - [x] T002.1.2.4 Test writing lock file with SHA-256 checksums `[activity: test-execution]`
            - [x] T002.1.2.5 Test idempotent reinstall via checksum comparison `[activity: test-execution]`
        - [x] T002.1.3 Implement src/core/installer/LockManager.ts `[ref: SDD; lines: 956-1023]` `[activity: domain-modeling]`
        - [x] T002.1.4 Validate
            - [x] T002.1.4.1 LockManager tests pass with 90%+ coverage `[activity: run-tests]`
            - [x] T002.1.4.2 Backward compatibility verified `[ref: SDD; lines: 49, 1116]` `[activity: business-acceptance]`

    - [x] T002.2 SettingsMerger Implementation `[parallel: true]` `[component: SettingsMerger]`
        - [x] T002.2.1 Prime Context
            - [x] T002.2.1.1 Read Settings merge algorithm `[ref: SDD; lines: 923-954]`
            - [x] T002.2.1.2 Read Settings structure `[ref: SDD; lines: 551-567]`
        - [x] T002.2.2 Write Tests
            - [x] T002.2.2.1 Test deep merge preserves user hooks `[activity: test-execution]`
            - [x] T002.2.2.2 Test new hooks are added without overwriting `[activity: test-execution]`
            - [x] T002.2.2.3 Test placeholder replacement ({{STARTUP_PATH}}, {{CLAUDE_PATH}}) `[activity: test-execution]`
            - [x] T002.2.2.4 Test backup creation before modification `[activity: test-execution]`
            - [x] T002.2.2.5 Test rollback on merge failure `[activity: test-execution]`
            - [x] T002.2.2.6 Test settings.json creation from scratch (file doesn't exist) `[ref: SDD; line: 242]` `[activity: test-execution]`
            - [x] T002.2.2.7 Test JSON validation produces clear error messages `[ref: SDD; line: 1159]` `[activity: test-execution]`
            - [x] T002.2.2.8 Test placeholder replacement failure with clear error message `[ref: SDD; line: 904]` `[activity: test-execution]`
        - [x] T002.2.3 Implement src/core/installer/SettingsMerger.ts `[ref: SDD; lines: 923-954]` `[activity: domain-modeling]`
        - [x] T002.2.4 Validate
            - [x] T002.2.4.1 SettingsMerger tests pass with 90%+ coverage `[activity: run-tests]`
            - [x] T002.2.4.2 User data preservation verified `[ref: PRD; line: 245]` `[activity: business-acceptance]`

    - [x] T002.3 Shell Scripts Implementation `[parallel: true]` `[component: Statusline]`
        - [x] T002.3.1 Prime Context
            - [x] T002.3.1.1 Read Statusline requirements `[ref: PRD; lines: 213-223]`
            - [x] T002.3.1.2 Read Statusline flow `[ref: SDD; lines: 799-820]`
        - [x] T002.3.2 Write Tests
            - [x] T002.3.2.1 Test Unix script parses JSON from stdin `[activity: test-execution]`
            - [x] T002.3.2.2 Test Windows script parses JSON from stdin using ConvertFrom-Json `[ref: SDD; line: 1190]` `[activity: test-execution]`
            - [x] T002.3.2.3 Test execution time < 10ms (benchmark) `[ref: PRD; line: 219]` `[activity: performance-testing]`
            - [x] T002.3.2.4 Test git branch detection (happy path) `[activity: test-execution]`
            - [x] T002.3.2.5 Test graceful degradation when git missing `[activity: test-execution]`
            - [x] T002.3.2.6 Test exact Claude Code JSON input format with schema validation `[ref: SDD; lines: 218-222]` `[activity: test-execution]`
            - [x] T002.3.2.7 Test git fallback strategy (.git/HEAD read fails → git command) `[ref: SDD; lines: 235-240]` `[activity: test-execution]`
            - [x] T002.3.2.8 Test home directory expansion (~) on macOS, Linux, Windows `[ref: SDD; line: 1188]` `[activity: test-execution]`
        - [x] T002.3.3 Implement bin/statusline.sh (bash/zsh compatible) `[ref: PRD; line: 217]` `[activity: component-development]`
        - [x] T002.3.4 Implement bin/statusline.ps1 (PowerShell compatible) `[ref: PRD; line: 218]` `[activity: component-development]`
        - [x] T002.3.5 Validate
            - [x] T002.3.5.1 Shell script tests pass on all platforms `[activity: run-tests]`
            - [x] T002.3.5.2 Performance benchmark: < 10ms verified `[ref: SDD; line: 1153]` `[activity: performance-testing]`

    - [x] T002.4 UI Shared Components `[parallel: true]` `[component: UI-Shared]`
        - [x] T002.4.1 Prime Context
            - [x] T002.4.1.1 Read Ink component pattern `[ref: SDD; lines: 1073-1086]`
            - [x] T002.4.1.2 Read UI theme structure `[ref: SDD; line: 454]`
        - [x] T002.4.2 Write Tests
            - [x] T002.4.2.1 Test Spinner renders correctly `[activity: test-execution]`
            - [x] T002.4.2.2 Test ErrorDisplay shows error message `[activity: test-execution]`
            - [x] T002.4.2.3 Test theme provides consistent colors `[activity: test-execution]`
        - [x] T002.4.3 Implement src/ui/shared/theme.ts (color scheme) `[ref: SDD; line: 454]` `[activity: component-development]`
        - [x] T002.4.4 Implement src/ui/shared/Spinner.tsx `[activity: component-development]`
        - [x] T002.4.5 Implement src/ui/shared/ErrorDisplay.tsx `[activity: component-development]`
        - [x] T002.4.6 Validate
            - [x] T002.4.6.1 Shared component tests pass `[activity: run-tests]`
            - [x] T002.4.6.2 Components follow Ink patterns `[ref: SDD; lines: 1073-1086]` `[activity: review-code]`

- [x] T003 Core Features & UI Components

    - [x] T003.1 Installer Core Implementation `[component: Installer]`
        - [x] T003.1.1 Prime Context
            - [x] T003.1.1.1 Read Installer interface `[ref: SDD; lines: 569-586]`
            - [x] T003.1.1.2 Read Install flow `[ref: SDD; lines: 677-718]`
            - [x] T003.1.1.3 Read Non-interactive flow `[ref: SDD; lines: 720-756]`
        - [x] T003.1.2 Write Tests
            - [x] T003.1.2.1 Test asset file copying to selected directories `[activity: test-execution]`
            - [x] T003.1.2.2 Test settings.json merge via SettingsMerger `[activity: test-execution]`
            - [x] T003.1.2.3 Test lock file creation via LockManager `[activity: test-execution]`
            - [x] T003.1.2.4 Test specific error: Invalid path with re-enter suggestion message `[ref: SDD; line: 900]` `[activity: test-execution]`
            - [x] T003.1.2.5 Test specific error: Permission denied with chmod/permissions suggestion `[ref: SDD; line: 901]` `[activity: test-execution]`
            - [x] T003.1.2.6 Test specific error: Disk full with space needed message `[ref: SDD; line: 922]` `[activity: test-execution]`
            - [x] T003.1.2.7 Test rollback on partial installation failure `[ref: SDD; line: 903]` `[activity: test-execution]`
            - [x] T003.1.2.8 Test progress indicators for long operations (>5s) `[ref: SDD; line: 1161]` `[activity: test-execution]`
            - [x] T003.1.2.9 Test case-sensitive vs case-insensitive path handling `[ref: SDD; lines: 37, 1189]` `[activity: test-execution]`
        - [x] T003.1.3 Implement src/core/installer/Installer.ts `[ref: SDD; lines: 569-586]` `[activity: domain-modeling]`
        - [x] T003.1.4 Validate
            - [x] T003.1.4.1 Installer tests pass with 90%+ coverage `[activity: run-tests]`
            - [x] T003.1.4.2 All error scenarios handled `[ref: SDD; lines: 899-921]` `[activity: business-acceptance]`

    - [x] T003.2 Initializer Implementation `[parallel: true]` `[component: Initializer]`
        - [x] T003.2.1 Prime Context
            - [x] T003.2.1.1 Read Init requirements `[ref: PRD; lines: 190-200]`
            - [x] T003.2.1.2 Read Init flow `[ref: SDD; lines: 822-868]`
            - [x] T003.2.1.3 Read Init interfaces `[ref: SDD; lines: 614-629]`
        - [x] T003.2.2 Write Tests
            - [x] T003.2.2.1 Test template copying to docs/ directory `[activity: test-execution]`
            - [x] T003.2.2.2 Test custom value replacement in templates `[activity: test-execution]`
            - [x] T003.2.2.3 Test --dry-run preview mode `[ref: PRD; line: 196]` `[activity: test-execution]`
            - [x] T003.2.2.4 Test --force overwrite mode `[ref: PRD; line: 197]` `[activity: test-execution]`
            - [x] T003.2.2.5 Test docs/ directory creation when missing `[activity: test-execution]`
        - [x] T003.2.3 Implement src/core/init/Initializer.ts `[ref: SDD; lines: 614-629]` `[activity: domain-modeling]`
        - [x] T003.2.4 Validate
            - [x] T003.2.4.1 Initializer tests pass with 90%+ coverage `[activity: run-tests]`
            - [x] T003.2.4.2 All PRD acceptance criteria met `[ref: PRD; lines: 192-199]` `[activity: business-acceptance]`

    - [x] T003.3 SpecGenerator Implementation `[parallel: true]` `[component: SpecGenerator]`
        - [x] T003.3.1 Prime Context
            - [x] T003.3.1.1 Read Spec requirements `[ref: PRD; lines: 201-211]`
            - [x] T003.3.1.2 Read Spec flow `[ref: SDD; lines: 870-895]`
            - [x] T003.3.1.3 Read Spec interfaces `[ref: SDD; lines: 631-653]`
        - [x] T003.3.2 Write Tests
            - [x] T003.3.2.1 Test spec directory creation with auto-incrementing ID `[activity: test-execution]`
            - [x] T003.3.2.2 Test TOML output format for --read flag `[ref: PRD; line: 206]` `[activity: test-execution]`
            - [x] T003.3.2.3 Test template generation for --add flag `[ref: PRD; line: 205]` `[activity: test-execution]`
            - [x] T003.3.2.4 Test spec ID parsing from directory names `[activity: test-execution]`
        - [x] T003.3.3 Implement src/core/spec/SpecGenerator.ts `[ref: SDD; lines: 631-653]` `[activity: domain-modeling]`
        - [x] T003.3.4 Validate
            - [x] T003.3.4.1 SpecGenerator tests pass with 90%+ coverage `[activity: run-tests]`
            - [x] T003.3.4.2 TOML output format verified `[ref: PRD; line: 206]` `[activity: business-acceptance]`

    - [x] T003.4 UI Interactive Components `[parallel: true]` `[component: UI-Interactive]`
        - [x] T003.4.1 Prime Context
            - [x] T003.4.1.1 Read Install TUI requirements `[ref: PRD; lines: 162-175]`
            - [x] T003.4.1.2 Read Install flow detailed `[ref: PRD; lines: 285-307]`
        - [x] T003.4.2 Write Tests
            - [x] T003.4.2.1 Test PathSelector path validation `[activity: test-execution]`
            - [x] T003.4.2.2 Test FileTree keyboard navigation (arrows, vim bindings) `[ref: PRD; line: 169]` `[activity: test-execution]`
            - [x] T003.4.2.3 Test FileTree space/enter selection `[activity: test-execution]`
            - [x] T003.4.2.4 Test Complete screen displays success message `[activity: test-execution]`
        - [x] T003.4.3 Implement src/ui/install/PathSelector.tsx `[activity: component-development]`
        - [x] T003.4.4 Implement src/ui/install/FileTree.tsx `[activity: component-development]`
        - [x] T003.4.5 Implement src/ui/install/Complete.tsx `[activity: component-development]`
        - [x] T003.4.6 Validate
            - [x] T003.4.6.1 Interactive component tests pass `[activity: run-tests]`
            - [x] T003.4.6.2 Keyboard navigation works as specified `[ref: PRD; line: 169]` `[activity: business-acceptance]`

- [ ] T004 UI Wizards

    - [ ] T004.1 InstallWizard Implementation `[parallel: true]` `[component: InstallWizard]`
        - [ ] T004.1.1 Prime Context
            - [ ] T004.1.1.1 Read Install wizard requirements `[ref: PRD; lines: 162-175]`
            - [ ] T004.1.1.2 Read Install user flow `[ref: PRD; lines: 285-307]`
            - [ ] T004.1.1.3 Read Install business rules `[ref: PRD; lines: 309-314]`
        - [ ] T004.1.2 Write Tests
            - [ ] T004.1.2.1 Test state transitions (Startup Path → Claude Path → File Selection → Complete) `[activity: test-execution]`
            - [ ] T004.1.2.2 Test --local flag skips TUI, uses defaults `[ref: PRD; line: 172]` `[activity: test-execution]`
            - [ ] T004.1.2.3 Test --yes flag auto-confirms prompts `[ref: PRD; line: 173]` `[activity: test-execution]`
            - [ ] T004.1.2.4 Test error display and recovery `[ref: PRD; lines: 316-322]` `[activity: test-execution]`
            - [ ] T004.1.2.5 Test Ctrl+C cancellation and rollback `[ref: PRD; line: 321]` `[activity: test-execution]`
        - [ ] T004.1.3 Implement src/ui/install/InstallWizard.tsx `[ref: SDD; lines: 447]` `[activity: component-development]`
        - [ ] T004.1.4 Validate
            - [ ] T004.1.4.1 InstallWizard tests pass `[activity: run-tests]`
            - [ ] T004.1.4.2 All edge cases handled `[ref: PRD; lines: 316-322]` `[activity: business-acceptance]`
            - [ ] T004.1.4.3 Business rules enforced `[ref: PRD; lines: 309-314]` `[activity: business-acceptance]`

    - [ ] T004.2 UninstallWizard Implementation `[parallel: true]` `[component: UninstallWizard]`
        - [ ] T004.2.1 Prime Context
            - [ ] T004.2.1.1 Read Uninstall requirements `[ref: PRD; lines: 177-189]`
            - [ ] T004.2.1.2 Read Uninstall flow `[ref: SDD; lines: 758-797]`
        - [ ] T004.2.2 Write Tests
            - [ ] T004.2.2.1 Test lock file reading and file list display `[activity: test-execution]`
            - [ ] T004.2.2.2 Test confirmation prompt `[activity: test-execution]`
            - [ ] T004.2.2.3 Test --keep-logs flag preserves logs directory `[ref: PRD; line: 186]` `[activity: test-execution]`
            - [ ] T004.2.2.4 Test --keep-settings flag skips settings restoration `[ref: PRD; line: 187]` `[activity: test-execution]`
            - [ ] T004.2.2.5 Test graceful handling of missing files `[ref: PRD; line: 188]` `[activity: test-execution]`
        - [ ] T004.2.3 Implement src/ui/uninstall/UninstallWizard.tsx `[ref: SDD; lines: 450]` `[activity: component-development]`
        - [ ] T004.2.4 Validate
            - [ ] T004.2.4.1 UninstallWizard tests pass `[activity: run-tests]`
            - [ ] T004.2.4.2 All PRD acceptance criteria met `[ref: PRD; lines: 179-188]` `[activity: business-acceptance]`

    - [ ] T004.3 Integration Test Harness `[component: Integration-Tests]`
        - [ ] T004.3.1 Create tests/integration/install-flow.test.ts (full install cycle) `[activity: test-execution]`
        - [ ] T004.3.2 Create tests/integration/uninstall-flow.test.ts (full uninstall cycle) `[activity: test-execution]`
        - [ ] T004.3.3 Create tests/integration/reinstall-flow.test.ts (idempotent reinstall) `[activity: test-execution]`
        - [ ] T004.3.4 Validate integration tests pass `[activity: run-tests]`

    - [ ] T004.4 Go-to-npm Migration Testing `[component: Migration]`
        - [ ] T004.4.1 Test upgrade from Go lock file (v1 string[]) to npm lock file (v2 with checksums) `[ref: SDD; lines: 49, 1116]` `[activity: test-execution]`
        - [ ] T004.4.2 Test detection of existing Go binary installation (if present) `[activity: test-execution]`
        - [ ] T004.4.3 Test settings.json preservation through Go→npm migration `[activity: test-execution]`
        - [ ] T004.4.4 Test all 55 assets migrate correctly from Go to npm installation `[activity: test-execution]`
        - [ ] T004.4.5 Validate migration guide instructions accuracy `[ref: PRD; line: 427]` `[activity: business-acceptance]`

- [ ] T005 CLI Commands

    - [ ] T005.1 Statusline CLI (Simplest) `[component: CLI-Statusline]`
        - [ ] T005.1.1 Prime Context
            - [ ] T005.1.1.1 Read Statusline requirements `[ref: PRD; lines: 213-223]`
        - [ ] T005.1.2 Write Tests
            - [ ] T005.1.2.1 Test CLI passes stdin to shell script `[activity: test-execution]`
            - [ ] T005.1.2.2 Test package.json bin field maps to correct script per platform `[ref: PRD; line: 222]` `[activity: test-execution]`
        - [ ] T005.1.3 Implement src/cli/statusline.ts `[ref: SDD; lines: 441]` `[activity: api-development]`
        - [ ] T005.1.4 Validate
            - [ ] T005.1.4.1 Statusline CLI tests pass `[activity: run-tests]`
            - [ ] T005.1.4.2 Cross-platform bin mapping works `[activity: business-acceptance]`

    - [ ] T005.2 Init and Spec CLIs `[parallel: true]`
        - [ ] T005.2.1 Init CLI `[component: CLI-Init]`
            - [ ] T005.2.1.1 Prime Context: Read Init requirements `[ref: PRD; lines: 190-200]`
            - [ ] T005.2.1.2 Write Tests
                - [ ] T005.2.1.2.1 Test Inquirer prompts launch correctly `[activity: test-execution]`
                - [ ] T005.2.1.2.2 Test --dry-run and --force flags work `[activity: test-execution]`
            - [ ] T005.2.1.3 Implement src/cli/init.ts `[ref: SDD; lines: 439]` `[activity: api-development]`
            - [ ] T005.2.1.4 Validate Init CLI tests pass `[activity: run-tests]`

        - [ ] T005.2.2 Spec CLI `[parallel: true]` `[component: CLI-Spec]`
            - [ ] T005.2.2.1 Prime Context: Read Spec requirements `[ref: PRD; lines: 201-211]`
            - [ ] T005.2.2.2 Write Tests
                - [ ] T005.2.2.2.1 Test spec creation with name argument `[activity: test-execution]`
                - [ ] T005.2.2.2.2 Test --add and --read flags work `[activity: test-execution]`
            - [ ] T005.2.2.3 Implement src/cli/spec.ts `[ref: SDD; lines: 440]` `[activity: api-development]`
            - [ ] T005.2.2.4 Validate Spec CLI tests pass `[activity: run-tests]`

    - [ ] T005.3 Install and Uninstall CLIs (Most Complex)
        - [ ] T005.3.1 Install CLI `[component: CLI-Install]`
            - [ ] T005.3.1.1 Prime Context
                - [ ] T005.3.1.1.1 Read Install requirements `[ref: PRD; lines: 162-175]`
                - [ ] T005.3.1.1.2 Read Install integration flow `[ref: SDD; lines: 677-718]`
            - [ ] T005.3.1.2 Write Tests
                - [ ] T005.3.1.2.1 Test TUI launch for interactive mode `[activity: test-execution]`
                - [ ] T005.3.1.2.2 Test direct Installer call for --local mode `[activity: test-execution]`
                - [ ] T005.3.1.2.3 Test flag combinations (--local --yes) `[activity: test-execution]`
            - [ ] T005.3.1.3 Implement src/cli/install.ts `[ref: SDD; lines: 438]` `[activity: api-development]`
            - [ ] T005.3.1.4 Validate Install CLI tests pass `[activity: run-tests]`

        - [ ] T005.3.2 Uninstall CLI `[component: CLI-Uninstall]`
            - [ ] T005.3.2.1 Prime Context: Read Uninstall requirements `[ref: PRD; lines: 177-189]`
            - [ ] T005.3.2.2 Write Tests
                - [ ] T005.3.2.2.1 Test TUI launch with file list from lock `[activity: test-execution]`
                - [ ] T005.3.2.2.2 Test --keep-logs and --keep-settings flags `[activity: test-execution]`
            - [ ] T005.3.2.3 Implement src/cli/uninstall.ts `[ref: SDD; lines: 438]` `[activity: api-development]`
            - [ ] T005.3.2.4 Validate Uninstall CLI tests pass `[activity: run-tests]`

    - [ ] T005.4 CLI Integration & Entry Points
        - [ ] T005.4.1 Implement src/cli/index.ts (Commander setup, register all commands) `[ref: SDD; lines: 437]` `[activity: api-development]`
        - [ ] T005.4.2 Implement src/index.ts (main entry point) `[ref: SDD; lines: 473]` `[activity: api-development]`
        - [ ] T005.4.3 Implement src/lib/index.ts (public API exports) `[ref: SDD; lines: 471-472]` `[activity: api-development]`
        - [ ] T005.4.4 Validate
            - [ ] T005.4.4.1 All CLI commands discoverable via --help `[activity: run-tests]`
            - [ ] T005.4.4.2 Entry point tests pass `[activity: run-tests]`

- [ ] T006 Integration & End-to-End Validation
    - [ ] T006.1 Unit Test Coverage
        - [ ] T006.1.1 Core business logic: 90%+ coverage verified `[ref: SDD; line: 1241]` `[activity: run-tests]`
        - [ ] T006.1.2 All error paths tested `[ref: SDD; line: 1244]` `[activity: run-tests]`
    - [ ] T006.2 Integration Tests
        - [ ] T006.2.1 Full install flow test passes (interactive mode) `[activity: run-tests]`
        - [ ] T006.2.2 Full install flow test passes (--local mode) `[activity: run-tests]`
        - [ ] T006.2.3 Full uninstall flow test passes `[activity: run-tests]`
        - [ ] T006.2.4 Idempotent reinstall test passes (checksum-based) `[activity: run-tests]`
        - [ ] T006.2.5 Settings.json merge conflict test passes `[ref: SDD; lines: 1224-1230]` `[activity: run-tests]`
    - [ ] T006.3 End-to-End User Flows
        - [ ] T006.3.1 Install happy path: All 55 assets copied correctly `[ref: PRD; lines: 155-161]` `[activity: business-acceptance]`
        - [ ] T006.3.2 Uninstall with lock file: All files removed `[ref: PRD; lines: 179-188]` `[activity: business-acceptance]`
        - [ ] T006.3.3 Init command: Templates created in docs/ `[ref: PRD; lines: 192-199]` `[activity: business-acceptance]`
        - [ ] T006.3.4 Spec command: Directory created with TOML output `[ref: PRD; lines: 204-210]` `[activity: business-acceptance]`
    - [ ] T006.4 Performance Tests
        - [ ] T006.4.1 Statusline execution < 10ms on slowest platform (macOS bash 3.2) `[ref: SDD; lines: 29, 1153]` `[activity: performance-testing]`
        - [ ] T006.4.2 Install completes in < 30 seconds for ALL 55 assets on slow disk I/O `[ref: SDD; line: 1041]` `[activity: performance-testing]`
        - [ ] T006.4.3 Uninstall completes in < 10 seconds with all 55 assets installed `[ref: SDD; line: 1042]` `[activity: performance-testing]`
        - [ ] T006.4.4 Package size < 5MB verified `[ref: SDD; lines: 32, 1156]` `[activity: performance-testing]`
        - [ ] T006.4.5 Statusline performance with large git repository (>1GB .git directory) `[activity: performance-testing]`
        - [ ] T006.4.6 Install performance under concurrent file system operations `[activity: performance-testing]`
    - [ ] T006.5 Cross-Platform Testing
        - [ ] T006.5.1 Test on macOS (darwin) `[ref: SDD; line: 36]` `[activity: exploratory-testing]`
        - [ ] T006.5.2 Test on Linux `[ref: SDD; line: 36]` `[activity: exploratory-testing]`
        - [ ] T006.5.3 Test on Windows (win32) `[ref: SDD; line: 36]` `[activity: exploratory-testing]`
        - [ ] T006.5.4 Shell scripts work on bash 3.2+, zsh, PowerShell 5.1+ `[ref: SDD; line: 29]` `[activity: exploratory-testing]`
    - [ ] T006.6 Security & Quality Validation
        - [ ] T006.6.1 Settings.json backup created before modification `[ref: SDD; line: 246]` `[activity: security-assessment]`
        - [ ] T006.6.2 Atomic operations verified (install/uninstall all-or-nothing) `[ref: SDD; line: 1163]` `[activity: security-assessment]`
        - [ ] T006.6.3 No hardcoded secrets in codebase `[activity: security-assessment]`
        - [ ] T006.6.4 npm audit passes (no vulnerabilities) `[activity: security-assessment]`
    - [ ] T006.7 Acceptance Criteria Verification
        - [ ] T006.7.1 All PRD Feature 1 acceptance criteria met (npm package distribution) `[ref: PRD; lines: 155-161]` `[activity: business-acceptance]`
        - [ ] T006.7.2 All PRD Feature 2 acceptance criteria met (install command) `[ref: PRD; lines: 164-175]` `[activity: business-acceptance]`
        - [ ] T006.7.3 All PRD Feature 3 acceptance criteria met (uninstall command) `[ref: PRD; lines: 180-188]` `[activity: business-acceptance]`
        - [ ] T006.7.4 All PRD Feature 4 acceptance criteria met (init command) `[ref: PRD; lines: 192-199]` `[activity: business-acceptance]`
        - [ ] T006.7.5 All PRD Feature 5 acceptance criteria met (spec command) `[ref: PRD; lines: 204-210]` `[activity: business-acceptance]`
        - [ ] T006.7.6 All PRD Feature 6 acceptance criteria met (statusline command) `[ref: PRD; lines: 215-222]` `[activity: business-acceptance]`
        - [ ] T006.7.7 All PRD Feature 7 acceptance criteria met (asset embedding) `[ref: PRD; lines: 228-235]` `[activity: business-acceptance]`
        - [ ] T006.7.8 All PRD Feature 8 acceptance criteria met (settings.json integration) `[ref: PRD; lines: 239-246]` `[activity: business-acceptance]`
        - [ ] T006.7.9 All PRD Feature 9 acceptance criteria met (lock file management) `[ref: PRD; lines: 249-256]` `[activity: business-acceptance]`
    - [ ] T006.8 Build & Deployment Verification
        - [ ] T006.8.1 npm run build produces clean dist/ output `[activity: run-tests]`
        - [ ] T006.8.2 Dual module output verified (ESM + CJS + .d.ts) `[ref: SDD; line: 1269]` `[activity: run-tests]`
        - [ ] T006.8.3 Package.json bin field correctly maps to CLI entry `[activity: run-tests]`
        - [ ] T006.8.4 npm pack verifies package contents `[ref: SDD; line: 338]` `[activity: run-tests]`
    - [ ] T006.9 Documentation & Assets
        - [ ] T006.9.1 README.md updated with TypeScript installation instructions `[activity: component-development]`
        - [ ] T006.9.2 All 55 assets present in dist/assets/ after build `[ref: PRD; lines: 228-235]` `[activity: business-acceptance]`
        - [ ] T006.9.3 Migration guide created (Go → TypeScript) `[activity: component-development]`
    - [ ] T006.10 SDD Compliance
        - [ ] T006.10.1 Implementation follows layered architecture (CLI → UI → Core → Lib) `[ref: SDD; lines: 346-352]` `[activity: review-code]`
        - [ ] T006.10.2 All architecture decisions implemented as designed `[ref: SDD; lines: 1099-1148]` `[activity: review-code]`
        - [ ] T006.10.3 Deep merge algorithm implemented correctly `[ref: SDD; lines: 923-954]` `[activity: review-code]`
        - [ ] T006.10.4 Lock file backward compatibility working `[ref: SDD; lines: 956-1023]` `[activity: review-code]`
    - [ ] T006.11 Final Checklist
        - [ ] T006.11.1 100% feature parity with Go version (excluding stats) `[ref: PRD; lines: 340-345]` `[activity: business-acceptance]`
        - [ ] T006.11.2 All CLI flags match Go version exactly `[ref: PRD; line: 344]` `[activity: business-acceptance]`
        - [ ] T006.11.3 No breaking changes to public interfaces `[ref: SDD; lines: 273-278]` `[activity: business-acceptance]`
        - [ ] T006.11.4 Package ready for npm publish `[activity: business-acceptance]`
