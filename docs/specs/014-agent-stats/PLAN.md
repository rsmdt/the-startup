# Implementation Plan

## Validation Checklist
- [x] Context Ingestion section complete with all required specs
- [x] Implementation phases logically organized
- [x] Each phase starts with test definition (TDD approach)
- [x] Dependencies between phases identified
- [x] Parallel execution marked where applicable
- [ ] Multi-component coordination identified (if applicable)
- [x] Final validation phase included
- [x] No placeholder content remains

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

- **PRD**: `docs/specs/014-agent-stats/PRD.md` - Product Requirements
- **SDD**: `docs/specs/014-agent-stats/SDD.md` - Solution Design

### Key Design Decisions

- **Multi-method agent detection with confidence scoring** - Task tool (100%), patterns (75-95%), context (60-80%) `[ref: SDD; lines: 477-509]`
- **Pipeline extension pattern** - Wrapper around existing aggregator to minimize changes `[ref: SDD; lines: 282-293]`
- **Streaming statistics with bounded memory** - Welford's algorithm and t-digest for O(1) memory `[ref: SDD; lines: 511-570]`
- **Progressive enhancement strategy** - CLI-first with optional TUI using Bubble Tea `[ref: SDD; lines: 333-340]`

### Implementation Context

- **Commands to run**:
  - Unit tests: `go test ./internal/stats/...` `[ref: SDD; lines: 249]`
  - Integration: `go test ./cmd/... -tags=integration` `[ref: SDD; lines: 250]`
  - Benchmarks: `go test -bench=. ./internal/stats/...` `[ref: SDD; lines: 251]`
  - Build: `go build -o the-startup` `[ref: SDD; lines: 261]`
  - Run: `./the-startup stats agents [--since 7d] [--format json]` `[ref: SDD; lines: 266-269]`
- **Patterns to follow**:
  - Existing aggregator pattern `[ref: internal/stats/aggregator.go]`
  - Display formatting from FormatToolLeaderboard `[ref: internal/stats/display.go]`
  - Bubble Tea model-update-view `[ref: internal/ui/installer.go]`
- **Interfaces to implement**:
  - AgentAnalyzer interface `[ref: SDD; lines: 349-360]`
  - AgentInvocation struct `[ref: SDD; lines: 363-372]`
  - GlobalAgentStats struct `[ref: SDD; lines: 375-384]`

---

## Implementation Phases

- [x] **Phase 1**: Foundation - Agent Data Structures and Types

    - [x] **Load Context**: Review existing stats infrastructure and data structures
        - [x] Read aggregator implementation `[ref: internal/stats/aggregator.go; lines: 96-115]`
        - [x] Read types definitions `[ref: internal/stats/types.go]`
        - [x] Review SDD interface specifications `[ref: SDD; lines: 349-384]`
    - [x] **Write Tests**: Define behavior for agent data structures
        - [x] Test AgentInvocation validation and initialization `[ref: SDD; lines: 363-372]` `[activity: test-implementation]`
        - [x] Test GlobalAgentStats with Welford's algorithm `[ref: SDD; lines: 375-384]` `[activity: test-implementation]`
        - [x] Test confidence score boundaries (0.0-1.0) `[ref: SDD; lines: 477-509]` `[activity: test-implementation]`
        - [x] Test AgentAnalyzer interface compliance `[ref: SDD; lines: 349-360]` `[activity: test-implementation]`
    - [x] **Implement**: Create agent-specific data structures
        - [x] Create internal/stats/agent_types.go with core structures `[activity: code-implementation]`
        - [x] Create internal/stats/agent_aggregator.go with interface `[activity: code-implementation]`
        - [x] Implement Welford's algorithm for streaming stats `[ref: SDD; lines: 511-543]` `[activity: code-implementation]`
    - [x] **Validate**: Ensure foundation is solid
        - [x] Run unit tests and achieve 100% coverage `[activity: run-tests]`
        - [x] Verify interfaces compatible with existing code `[activity: review-code]`
        - [x] Check no impact on existing stats functionality `[activity: integration-test]`

- [x] **Phase 2**: Multi-Method Agent Detection Engine

    - [x] **Load Context**: Understand detection requirements and patterns
        - [x] Review current agent detection `[ref: cmd/stats.go; lines: 254-271]`
        - [x] Study multi-method detection design `[ref: SDD; lines: 477-509]`
        - [x] Review pattern matching requirements `[ref: PRD; lines: 119-122]`
    - [x] **Write Tests**: Test all detection methods and edge cases
        - [x] Test Task tool subagent_type extraction (100% confidence) `[ref: SDD; lines: 478-485]` `[activity: test-implementation]`
        - [x] Test pattern matching (75-95% confidence) `[ref: SDD; lines: 487-497]` `[activity: test-implementation]`
        - [x] Test command context inference (60-80% confidence) `[ref: SDD; lines: 499-507]` `[activity: test-implementation]`
        - [x] Test corrupted log handling `[ref: PRD; lines: 127]` `[activity: test-implementation]`
    - [x] **Implement**: Build multi-method detection system
        - [x] Create internal/stats/agent_detector.go `[activity: code-implementation]`
        - [x] Implement DetectAgent with confidence scoring `[ref: SDD; lines: 477-509]` `[activity: code-implementation]`
        - [x] Add detection hook to parser `[ref: internal/stats/parser.go]` `[activity: code-implementation]`
    - [x] **Validate**: Verify detection accuracy
        - [x] Achieve 99.5% accuracy for Task tool detection `[activity: run-tests]`
        - [x] Verify 90%+ accuracy for pattern matching `[activity: run-tests]`
        - [x] Performance impact <10% on processing `[activity: benchmark]`

- [x] **Phase 3**: Agent-Specific Aggregation Logic

    - [x] **Load Context**: Study existing aggregation and statistical algorithms
        - [x] Review Welford's algorithm implementation `[ref: SDD; lines: 511-543]`
        - [x] Understand t-digest for percentiles `[ref: SDD; lines: 544-570]`
        - [x] Study memory management requirements `[ref: SDD; lines: 828-831]`
    - [x] **Write Tests**: Test statistical accuracy and performance
        - [x] Test Welford's mean/variance calculations `[ref: SDD; lines: 519-543]` `[activity: test-implementation]`
        - [x] Test t-digest percentile accuracy (p50, p95, p99) `[ref: SDD; lines: 544-570]` `[activity: test-implementation]`
        - [x] Test LRU cache eviction under pressure `[activity: test-implementation]`
        - [x] Benchmark with 1M+ log entries `[activity: benchmark]`
    - [x] **Implement**: Build aggregation with streaming statistics
        - [x] Enhance agent_aggregator.go with Welford's `[ref: SDD; lines: 511-543]` `[activity: code-implementation]`
        - [x] Integrate t-digest for percentiles `[ref: SDD; lines: 544-570]` `[activity: code-implementation]`
        - [x] Implement LRU cache (max 1000 agents) `[activity: code-implementation]`
        - [x] Add thread-safe operations with sync.RWMutex `[activity: code-implementation]`
    - [x] **Validate**: Ensure performance and accuracy
        - [x] Response time <200ms for leaderboard `[ref: PRD; lines: 137]` `[activity: benchmark]`
        - [x] Memory usage <100MB for 1M entries `[ref: SDD; lines: 839]` `[activity: benchmark]`
        - [x] All race detector tests pass `[activity: run-tests]`

- [x] **Phase 4**: CLI Command Integration

    - [x] **Load Context**: Understand existing command structure
        - [x] Review stats command implementation `[ref: cmd/stats.go]`
        - [x] Study cobra command patterns `[ref: cmd/]`
        - [x] Review output format handling `[ref: internal/stats/display.go]`
    - [x] **Write Tests**: Test command integration and flags
        - [x] Test agent subcommand registration `[activity: test-implementation]`
        - [x] Test --since flag filtering `[ref: PRD; lines: 74]` `[activity: test-implementation]`
        - [x] Test --format flag (table, json, csv) `[ref: PRD; lines: 74]` `[activity: test-implementation]`
        - [x] Test empty results handling `[ref: PRD; lines: 124]` `[activity: test-implementation]`
    - [x] **Implement**: Add agent subcommand to stats
        - [x] Create cmd/stats_agent.go with cobra command `[activity: code-implementation]`
        - [x] Parse flags and call aggregator `[activity: code-implementation]`
        - [x] Update cmd/stats.go to register subcommand `[activity: code-implementation]`
    - [x] **Validate**: Verify CLI integration
        - [x] Agent subcommand appears in help `[activity: manual-test]`
        - [x] All existing commands still work `[activity: integration-test]`
        - [x] Integration tests pass `[activity: run-tests]`

- [x] **Phase 5**: Display Formatting and Output

    - [x] **Load Context**: Study existing display patterns
        - [x] Review FormatToolLeaderboard implementation `[ref: internal/stats/display.go]`
        - [x] Understand table/JSON/CSV formatting `[ref: internal/stats/]`
        - [x] Review sparkline generation patterns `[ref: PRD; lines: 73]`
    - [x] **Write Tests**: Test all output formats
        - [x] Test table formatting with tabwriter `[activity: test-implementation]`
        - [x] Test JSON marshaling correctness `[activity: test-implementation]`
        - [x] Test CSV header and row generation `[activity: test-implementation]`
        - [x] Test sparkline visualization `[ref: PRD; lines: 73]` `[activity: test-implementation]`
    - [x] **Implement**: Create display formatters
        - [x] Create internal/stats/agent_display.go `[activity: code-implementation]`
        - [x] Implement FormatAgentTable with rankings `[ref: PRD; lines: 71-74]` `[activity: code-implementation]`
        - [x] Add JSON and CSV formatters `[activity: code-implementation]`
        - [x] Generate sparklines for trends `[ref: PRD; lines: 73]` `[activity: code-implementation]`
    - [x] **Validate**: Ensure display quality
        - [x] Output matches existing format style `[activity: manual-test]`
        - [x] All formats produce valid output `[activity: run-tests]`
        - [x] Terminal colors work correctly `[activity: manual-test]`

- [x] **Phase 6**: Basic Delegation Pattern Analysis

    - [x] **Load Context**: Understand delegation tracking requirements
        - [x] Review delegation pattern requirements `[ref: PRD; lines: 83-91]`
        - [x] Study session statistics structure `[ref: internal/stats/types.go]`
        - [x] Review co-occurrence analysis needs `[ref: PRD; lines: 87-89]`
    - [x] **Write Tests**: Test delegation pattern tracking
        - [x] Test delegation chain construction `[ref: PRD; lines: 89]` `[activity: test-implementation]`
        - [x] Test agent co-occurrence tracking `[ref: PRD; lines: 87]` `[activity: test-implementation]`
        - [x] Test session boundary detection `[ref: PRD; lines: 90]` `[activity: test-implementation]`
        - [x] Test circular delegation handling `[activity: test-implementation]`
    - [x] **Implement**: Build delegation analysis
        - [x] Track agent transitions in sessions `[ref: PRD; lines: 87-90]` `[activity: code-implementation]`
        - [x] Build delegation graphs `[ref: PRD; lines: 89]` `[activity: code-implementation]`
        - [x] Calculate co-occurrence frequencies `[ref: PRD; lines: 88]` `[activity: code-implementation]`
    - [x] **Validate**: Verify pattern detection
        - [x] Accurately tracks multi-agent workflows `[activity: run-tests]`
        - [x] Handles complex patterns correctly `[activity: integration-test]`
        - [x] Performance remains acceptable `[activity: benchmark]`

- [x] **Phase 7**: Interactive TUI Mode (Optional)

    - [x] **Load Context**: Study Bubble Tea patterns
        - [x] Review existing TUI implementation `[ref: internal/ui/installer.go]`
        - [x] Study Bubble Tea model-update-view `[ref: SDD; lines: 74-76]`
        - [x] Review interactive requirements `[ref: PRD; lines: 114-116]`
    - [x] **Write Tests**: Test TUI components
        - [x] Test model initialization `[activity: test-implementation]`
        - [x] Test update message handling `[activity: test-implementation]`
        - [x] Test view rendering `[activity: test-implementation]`
        - [x] Test keyboard navigation `[activity: test-implementation]`
    - [x] **Implement**: Build interactive dashboard
        - [x] Create internal/ui/agent/dashboard_model.go `[activity: code-implementation]`
        - [x] Create leaderboard_component.go `[activity: code-implementation]`
        - [x] Add --interactive flag handling `[ref: PRD; lines: 114]` `[activity: code-implementation]`
    - [x] **Validate**: Ensure TUI quality
        - [x] TUI launches with --interactive flag `[activity: manual-test]`
        - [x] Navigation works smoothly `[activity: manual-test]`
        - [x] Graceful fallback to CLI `[ref: PRD; lines: 126]` `[activity: manual-test]`

- [x] **Integration & End-to-End Validation**
    - [x] All unit tests passing across all components `[activity: run-tests]`
    - [x] Integration tests for complete agent analytics flow `[activity: integration-test]`
    - [x] End-to-end tests with real JSONL files `[activity: integration-test]`
    - [x] Performance tests meet requirements `[ref: SDD; lines: 839-841]` `[activity: benchmark]`
    - [x] Security validation passes `[ref: SDD; lines: 841]` `[activity: security-test]`
    - [x] Acceptance criteria verified against PRD `[ref: PRD; lines: 71-91]` `[activity: business-acceptance]`
    - [x] Test coverage >90% for new code `[activity: coverage-test]`
    - [x] Documentation updated for new commands `[activity: documentation]`
    - [x] Build and deployment verification `[activity: build-test]`
    - [x] All PRD requirements implemented `[ref: PRD; lines: 64-104]` `[activity: business-acceptance]`
    - [x] Implementation follows SDD design `[ref: SDD]` `[activity: architecture-review]`