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

## Context Ingestion

*GATE: You MUST fully read all files mentioned in this section before starting any implementation.*

### Specification

- **PRD**: `docs/specs/015-interactive-stats-dashboard/PRD.md` - Product Requirements
- **SDD**: `docs/specs/015-interactive-stats-dashboard/SDD.md` - Solution Design

### Key Design Decisions

- **Complete CLI Replacement**: Replace stats command entirely with interactive dashboard (no --interactive flag)
- **Static Data Loading**: Use existing stats.ParseJSONLFiles() with re-read on filter changes (no streaming)
- **BubbleTea Architecture**: Component-based panels with Tab navigation and focus management
- **Time-based Filtering**: Support 1w, 1m, 3m, 6m, 1y, all lookback periods as primary filtering
- **Fullscreen Layout**: 4-panel dashboard (timeline, stats, delegation, co-occurrence) similar to btop

### Implementation Context

- **Commands to run**:
  - `go test ./internal/ui/agent/...` - Test dashboard components
  - `go test ./internal/stats/...` - Test stats engine integration
  - `go build -o the-startup` - Build binary
  - `./the-startup stats` - Test new dashboard experience
- **Patterns to follow**:
  - `docs/patterns/tui-dashboard-architecture.md` - BubbleTea component patterns
  - Existing installer patterns in `internal/ui/installer/` - State management and navigation
- **Interfaces to implement**:
  - `docs/interfaces/new-stats-command.md` - New command interface
  - `docs/interfaces/dashboard-data-api.md` - Data layer specifications

---

## Implementation Phases

**Current State Assessment**: The dashboard foundation is 75% complete with existing BubbleTea infrastructure, delegation analysis, and agent statistics already implemented. Focus on completing timeline visualization, filtering controls, and CLI replacement.

- [x] **Phase 1**: Foundation & Core Dashboard Replacement

    - [x] **Load Context**: Core dashboard replacement requirements
        - [x] Read SDD Section "Solution Strategy" `[ref: SDD; lines: 255-263]`
        - [x] Read new command interface `[ref: docs/interfaces/new-stats-command.md; lines: 1-50]`
        - [x] Analyze existing installer patterns `[ref: internal/ui/installer/; activity: code-analysis]`
    - [x] **Write Tests**: Dashboard launch and CLI replacement behavior
        - [x] Test dashboard launches without flags when running `stats` command `[activity: test-writer]`
        - [x] Test fullscreen layout renders correctly with all panels `[activity: test-writer]`
        - [x] Test default 1-month data loading works `[activity: test-writer]`
        - [x] Test exit functionality (q, Ctrl+C, ESC) `[activity: test-writer]`
    - [x] **Implement**: Replace cmd/stats.go with dashboard launcher
        - [x] Replace CLI stats command with dashboard launcher `[activity: software-engineer]`
        - [x] Create main dashboard model with 4-panel layout `[activity: component-architecture]`
        - [x] Implement basic panel interface and navigation `[activity: component-architecture]`
        - [x] Add data loading using existing stats engine `[activity: software-engineer]`
    - [x] **Validate**: Foundation working correctly
        - [x] Run tests to ensure dashboard launches `[activity: run-tests]`
        - [x] Verify CLI replacement works as expected `[activity: business-acceptance]`
        - [x] Check code follows existing patterns `[activity: review-code]`

- [x] **Phase 2**: Panel Implementation and Navigation

    - [x] **Stats Panel Implementation** `[parallel: true]` `[component: stats-panel]`
        - [x] **Load Context**: Agent statistics requirements
            - [x] Read existing leaderboard component `[ref: internal/ui/agent/leaderboard_component.go]`
            - [x] Read agent stats data structures `[ref: internal/stats/agent_detector.go]`
        - [x] **Write Tests**: Agent leaderboard with sorting and selection
            - [x] Test agent leaderboard rendering with real data `[activity: test-writer]`
            - [x] Test sorting by usage count, success rate, duration `[activity: test-writer]`
            - [x] Test agent selection and highlighting `[activity: test-writer]`
        - [x] **Implement**: Enhance existing stats panel
            - [x] Add sorting controls to leaderboard `[activity: component-architecture]`
            - [x] Implement agent selection and highlighting `[activity: software-engineer]`
            - [x] Add success rate and duration metrics display `[activity: software-engineer]`
        - [x] **Validate**: Stats panel functionality
            - [x] Verify sorting works correctly `[activity: run-tests]`
            - [x] Check selection state management `[activity: review-code]`

    - [x] **Delegation Panel Implementation** `[parallel: true]` `[component: delegation-panel]`
        - [x] **Load Context**: Delegation pattern requirements
            - [x] Read delegation analysis logic `[ref: internal/stats/delegation.go]`
            - [x] Read delegation data structures `[ref: internal/stats/delegation_test.go]`
        - [x] **Write Tests**: Delegation flow visualization
            - [x] Test delegation pattern list rendering `[activity: test-writer]`
            - [x] Test pattern frequency and success rate display `[activity: test-writer]`
            - [x] Test delegation chain navigation `[activity: test-writer]`
        - [x] **Implement**: Create delegation panel from scratch
            - [x] Build delegation pattern list component `[activity: component-architecture]`
            - [x] Add flow visualization (text-based) `[activity: software-engineer]`
            - [x] Implement pattern filtering and sorting `[activity: software-engineer]`
        - [x] **Validate**: Delegation panel functionality
            - [x] Test delegation patterns display correctly `[activity: run-tests]`
            - [x] Verify pattern analysis accuracy `[activity: business-acceptance]`

    - [x] **Panel Navigation System**
        - [x] **Load Context**: Navigation patterns
            - [x] Read installer navigation patterns `[ref: internal/ui/installer/models.go]`
            - [x] Read dashboard state management `[ref: internal/ui/agent/state.go]`
        - [x] **Write Tests**: Tab navigation and focus management
            - [x] Test Tab/Shift-Tab panel switching `[activity: test-writer]`
            - [x] Test panel focus indicators `[activity: test-writer]`
            - [x] Test keyboard shortcuts (1-4 for panel jumping) `[activity: test-writer]`
        - [x] **Implement**: Enhanced navigation system
            - [x] Implement Tab-based panel navigation `[activity: component-architecture]`
            - [x] Add visual focus indicators `[activity: software-engineer]`
            - [x] Create keyboard shortcut system `[activity: software-engineer]`
        - [x] **Validate**: Navigation works smoothly
            - [x] Test all navigation methods work `[activity: run-tests]`
            - [x] Verify focus management is consistent `[activity: review-code]`

- [x] **Phase 3**: Timeline and Co-occurrence Panels

    - [x] **Timeline Panel** `[component: timeline-panel]`
        - [x] **Load Context**: Timeline visualization requirements
            - [x] Read message timeline requirements `[ref: docs/interfaces/dashboard-data-api.md; lines: 50-80]`
        - [x] **Write Tests**: Timeline display and navigation
            - [x] Test message timeline rendering `[activity: test-writer]`
            - [x] Test time range navigation `[activity: test-writer]`
        - [x] **Implement**: Create timeline component
            - [x] Build horizontal timeline component `[activity: component-architecture]`
            - [x] Add agent usage visualization over time `[activity: software-engineer]`
        - [x] **Validate**: Timeline functionality
            - [x] Test timeline accuracy with real data `[activity: run-tests]`

    - [x] **Co-occurrence Panel** `[component: cooccurrence-panel]`
        - [x] **Load Context**: Agent relationship analysis
            - [x] Read co-occurrence data structures `[ref: docs/interfaces/dashboard-data-api.md; lines: 90-120]`
        - [x] **Write Tests**: Agent relationship display
            - [x] Test co-occurrence matrix rendering `[activity: test-writer]`
            - [x] Test correlation calculations `[activity: test-writer]`
        - [x] **Implement**: Agent relationship visualization
            - [x] Create co-occurrence matrix component `[activity: component-architecture]`
            - [x] Add correlation analysis display `[activity: software-engineer]`
        - [x] **Validate**: Co-occurrence analysis accuracy
            - [x] Verify relationship calculations `[activity: run-tests]`

- [x] **Phase 4**: Filtering and Polish

    - [x] **Time-based Filtering System**
        - [x] **Load Context**: Filter requirements
            - [x] Read filter specifications `[ref: docs/interfaces/new-stats-command.md; lines: 80-120]`
        - [x] **Write Tests**: Time filter functionality
            - [x] Test 1w, 1m, 3m, 6m, 1y, all filters `[activity: test-writer]`
            - [x] Test data refresh on filter changes `[activity: test-writer]`
        - [x] **Implement**: Filter controls and data refresh
            - [x] Add time range filter controls `[activity: component-architecture]`
            - [x] Implement data refresh on filter changes `[activity: software-engineer]`
        - [x] **Validate**: Filtering works correctly
            - [x] Test all time ranges load correct data `[activity: run-tests]`

    - [x] **Help System and Polish**
        - [x] **Write Tests**: Help and polish features
            - [x] Test help overlay display and navigation `[activity: test-writer]`
            - [x] Test all keyboard shortcuts work `[activity: test-writer]`
        - [x] **Implement**: Enhanced user experience
            - [x] Add help system with keyboard shortcuts `[activity: software-engineer]`
            - [x] Polish visual styling and layout `[activity: component-architecture]`
        - [x] **Validate**: Professional user experience
            - [x] Verify help system is comprehensive `[activity: business-acceptance]`

- [x] **Integration & End-to-End Validation**
    - [x] All unit tests passing for each panel component
    - [x] Integration tests for panel coordination and data flow
    - [x] End-to-end tests for complete user workflows (launch → navigate → filter → exit)
    - [x] CLI replacement validation - stats command launches dashboard `[ref: docs/interfaces/new-stats-command.md]`
    - [x] Data loading performance with large datasets acceptable
    - [x] All keyboard navigation and shortcuts functional
    - [x] Panel layout responsive to different terminal sizes
    - [x] Test coverage meets 80%+ standard for all components
    - [x] Documentation updated for breaking CLI changes
    - [x] Build verification - binary works without --interactive flag
    - [x] All dashboard panels display correct data from stats engine
    - [x] Implementation follows SDD component architecture `[ref: SDD Building Block View]`
    - [x] Timeline graph renders horizontally with User on top, Assistant on bottom
    - [x] Graph uses braille characters for smooth btop-style visualization
    - [x] Timeline panel height increased to 25% for better visibility
