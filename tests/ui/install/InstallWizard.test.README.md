# InstallWizard Test Suite

## Overview

This test suite provides comprehensive coverage for the `InstallWizard` component, following TDD principles. The tests are written in the **Red phase** - they are designed to fail until the implementation is complete.

## Test Coverage

### 1. State Machine Transitions (5 tests)
Tests the core state machine flow:
- **Startup Path → Claude Path → File Selection → Complete**
- Validates each state transition
- Ensures proper progression through the wizard

**References:**
- PRD lines 162-175 (Install Command requirements)
- PRD lines 285-307 (User flow)

### 2. --local Flag Behavior (4 tests)
Tests non-interactive mode:
- Skips TUI when `--local` flag is provided
- Uses defaults: `./.the-startup` and `~/.claude`
- Proceeds directly to installation

**References:**
- PRD line 172 (`--local` flag specification)
- PRD line 309 (Business Rule 1)

### 3. --yes Flag Behavior (5 tests)
Tests auto-confirm mode:
- Auto-confirms all prompts with recommended paths
- Auto-selects all file categories
- Completes without user interaction
- Works in combination with `--local` flag

**References:**
- PRD line 173 (`--yes` flag specification)
- PRD line 310 (Business Rule 2)

### 4. Error Handling and Recovery (7 tests)
Tests all error scenarios from SDD:
- Invalid path errors
- Permission denied errors
- Disk full errors
- Settings merge failures
- Asset copy failures
- Error recovery flow

**References:**
- SDD lines 899-921 (Error Handling section)
- PRD lines 316-322 (Edge Cases)

### 5. Ctrl+C Cancellation and Rollback (6 tests)
Tests cancellation behavior:
- Detects Ctrl+C interruption
- Triggers rollback mechanism
- Deletes partial installation files
- Deletes incomplete lock file
- Shows cancellation message
- Exits gracefully

**References:**
- PRD line 321 (Edge Case 6: Ctrl+C handling)
- SDD line 306 (Rollback mechanism)

### 6. Business Rules (6 tests)
Tests PRD business rules enforcement:
- **Rule 1:** `--local` uses `./.the-startup` and `~/.claude` defaults
- **Rule 2:** `--yes` auto-confirms with recommended settings
- **Rule 3:** Merge with existing files (don't overwrite)
- **Rule 4:** Merge hooks in settings.json
- **Rule 5:** Detect reinstall from lock file, compare checksums

**References:**
- PRD lines 309-314 (Business Rules section)

### 7. Integration with Sub-components (5 tests)
Tests integration with existing components:
- PathSelector for startup path
- PathSelector for Claude path
- FileTree for file selection
- Complete for success screen
- Installer.install() integration

**Components:**
- `/src/ui/install/PathSelector.tsx`
- `/src/ui/install/FileTree.tsx`
- `/src/ui/install/Complete.tsx`
- `/src/core/installer/Installer.ts`

### 8. Progress Indication (3 tests)
Tests progress reporting for long operations:
- Shows progress for operations > 5 seconds
- Displays current installation stage
- Shows file count during copying

**References:**
- SDD lines 1079 (Progress requirements)

### 9. Component Pattern (5 tests)
Tests Ink component compliance:
- Functional component pattern
- Valid React element
- Uses `useState` for state management
- Uses `useInput` for keyboard handling
- Uses `useApp` for exit control

**References:**
- SDD lines 1073-1086 (Component Structure pattern)

### 10. Props (4 tests)
Tests component API:
- Required props: `options`, `onComplete`
- Optional `local` flag
- Optional `yes` flag
- `onComplete` callback behavior

### 11. Edge Cases (5 tests)
Tests special scenarios:
- Empty options object
- Missing Claude directory (PRD line 317)
- Malformed settings.json (PRD line 318)
- Disk full error (PRD line 322)
- Offline mode (PRD line 319)

### 12. Accessibility (3 tests)
Tests UI accessibility:
- Keyboard navigation hints
- Clear state indicators
- Visual feedback for actions

## Total Test Count: 58 tests

## Test Execution

Currently in **Red Phase** - all tests fail because `InstallWizard.tsx` doesn't exist:

```bash
npm test -- tests/ui/install/InstallWizard.test.tsx
```

Expected output:
```
Error: Failed to load url ../../../src/ui/install/InstallWizard
Does the file exist?
```

## Next Steps (Green Phase)

1. **Create** `/src/ui/install/InstallWizard.tsx`
2. **Implement** the state machine with 4 states
3. **Integrate** existing components (PathSelector, FileTree, Complete)
4. **Implement** flag behaviors (`--local`, `--yes`)
5. **Implement** error handling and recovery
6. **Implement** Ctrl+C cancellation and rollback
7. **Enforce** all business rules
8. **Add** progress indication
9. **Run tests** until they pass

## Implementation Checklist

Based on test requirements, the implementation must:

- [ ] Define state machine enum: `StartupPath | ClaudePath | FileSelection | Complete | Error`
- [ ] Use `useState` to manage current state
- [ ] Use `useInput` to detect Ctrl+C
- [ ] Use `useApp` to exit on cancellation
- [ ] Render appropriate component for each state
- [ ] Handle `--local` flag (skip prompts, use defaults)
- [ ] Handle `--yes` flag (auto-confirm all)
- [ ] Handle flag combinations (`--local --yes`)
- [ ] Implement error display with recovery
- [ ] Implement rollback on Ctrl+C
- [ ] Call `Installer.install()` with correct options
- [ ] Show progress for long operations
- [ ] Enforce all 5 business rules
- [ ] Handle all edge cases from PRD

## Key Files Reference

**Implementation file to create:**
- `/src/ui/install/InstallWizard.tsx`

**Dependencies (already implemented):**
- `/src/ui/install/PathSelector.tsx` - Path selection component
- `/src/ui/install/FileTree.tsx` - File tree selection
- `/src/ui/install/Complete.tsx` - Success screen
- `/src/core/installer/Installer.ts` - Installation logic
- `/src/core/types/config.ts` - Type definitions

**Specification references:**
- `/docs/specs/004-typescript-npm-package-migration/PRD.md` (lines 162-175, 285-322)
- `/docs/specs/004-typescript-npm-package-migration/SDD.md` (lines 899-921, 1073-1086)
- `/docs/specs/004-typescript-npm-package-migration/PLAN.md` (lines 306-316)

## Test Philosophy

These tests follow TDD best practices:

1. **Red Phase:** Tests fail because implementation doesn't exist ✓
2. **Green Phase:** Implement minimal code to make tests pass (next step)
3. **Refactor Phase:** Improve code while keeping tests green (after Green)

All tests are:
- **Specific:** Test one behavior per test
- **Independent:** No test depends on another
- **Repeatable:** Same result every time
- **Self-validating:** Clear pass/fail
- **Timely:** Written before implementation

## Notes for Implementation

1. **State Management:** Use single `useState` for state machine, not multiple booleans
2. **Error Handling:** Use ErrorDisplay component from shared components
3. **Progress:** Use Spinner component from shared components
4. **Validation:** Reuse PathValidationResult type from PathSelector
5. **Installation:** Use Installer class with dependency injection for testability
6. **Rollback:** Track installed files in state for cleanup
7. **Flags:** Check flags early, skip states accordingly
8. **Cancellation:** Use try-catch with cleanup in finally block
