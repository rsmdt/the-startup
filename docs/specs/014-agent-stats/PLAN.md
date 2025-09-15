# Implementation Plan - Agent Analytics Feature

## Overview
This implementation plan breaks down the agent analytics feature into logical phases following Test-Driven Development (TDD) principles. Each phase builds upon the previous one, with clear dependencies and testing requirements.

## Phase Structure
Each phase follows the Red-Green-Refactor TDD cycle:
- **RED**: Write failing tests that define the desired behavior
- **GREEN**: Write minimal code to make tests pass
- **REFACTOR**: Improve code structure while keeping tests green

---

## PHASE 1: Foundation - Agent Data Structures and Types
**Duration**: 1-2 days
**Dependency**: None (can start immediately)

### Objectives
- Define core agent data structures
- Establish type system for agent analytics
- Create interfaces for agent processing

### TDD Test Requirements
```go
// Test files to create first (RED phase)
internal/stats/agent_types_test.go
- Test AgentInvocation structure validation
- Test GlobalAgentStats initialization
- Test confidence score boundaries (0.0-1.0)
- Test detection method enumeration
- Test agent type constants

internal/stats/agent_aggregator_test.go
- Test AgentAnalyzer interface compliance
- Test empty aggregator state
- Test aggregator initialization
```

### Implementation Tasks (GREEN phase)
1. Create `internal/stats/agent_types.go`
   - Define `AgentInvocation` struct with fields: AgentType, SessionID, Timestamp, DurationMs, Success, Confidence, DetectionMethod
   - Define `GlobalAgentStats` struct with Welford's algorithm fields
   - Define `DelegationChain` struct for tracking agent-to-agent calls
   - Define detection method constants (TaskTool, PatternMatch, CommandContext)

2. Create `internal/stats/agent_aggregator.go`
   - Define `AgentAnalyzer` interface extending existing Aggregator
   - Implement basic struct that embeds existing aggregator
   - Add stubs for ProcessAgentInvocation, GetAgentStats, GetAllAgentStats

### Refactor Tasks
- Extract common statistical fields into reusable types
- Ensure consistency with existing tool stats structures
- Add comprehensive documentation

### Success Criteria
- All unit tests pass
- Interfaces compatible with existing aggregator
- Zero impact on existing stats functionality

---

## PHASE 2: Multi-Method Agent Detection Engine
**Duration**: 2-3 days
**Dependency**: Phase 1 (data structures)

### Objectives
- Implement multi-method agent detection with confidence scoring
- Extract agents from Task tool parameters
- Pattern matching for agent mentions
- Command context inference

### TDD Test Requirements
```go
// Test files to create first (RED phase)
internal/stats/agent_detector_test.go
- Test Task tool subagent_type extraction (100% confidence)
- Test pattern matching detection (75-95% confidence)
- Test command context inference (60-80% confidence)
- Test confidence score calculations
- Test edge cases: malformed logs, missing fields, empty content
- Test detection priority (Task > Pattern > Context)

// Test fixtures to create
testdata/agent_logs/
- task_tool_invocations.jsonl (valid Task tool with subagent_type)
- pattern_mentions.jsonl (logs with "delegating to" patterns)
- command_context.jsonl (logs with command-based agent inference)
- corrupted_entries.jsonl (malformed/incomplete logs)
```

### Implementation Tasks (GREEN phase)
1. Create `internal/stats/agent_detector.go`
   - Implement `DetectAgent(entry LogEntry) (agent string, confidence float64)`
   - Method 1: `extractSubagentType()` for Task tool parameters
   - Method 2: `matchAgentPatterns()` with regex patterns
   - Method 3: `inferFromCommand()` for command context
   - Confidence scoring logic based on detection method

2. Update `internal/stats/parser.go`
   - Add agent detection hook in ParseEntry
   - Pass detected agents to aggregator
   - Maintain backward compatibility

### Refactor Tasks
- Extract regex patterns to configuration
- Optimize pattern matching performance
- Cache compiled regex patterns

### Success Criteria
- 99.5% accuracy for Task tool detection
- 90%+ accuracy for pattern matching
- All corrupted log entries handled gracefully
- Performance impact <10% on log processing

---

## PHASE 3: Agent-Specific Aggregation Logic
**Duration**: 2-3 days
**Dependency**: Phase 2 (detection engine)

### Objectives
- Implement Welford's algorithm for agent statistics
- Add t-digest for percentile calculations
- Memory-bounded caching with LRU eviction
- Integration with existing aggregator

### TDD Test Requirements
```go
// Test files to enhance (RED phase)
internal/stats/agent_aggregator_test.go
- Test Welford's algorithm mean/variance calculations
- Test t-digest percentile accuracy (p50, p95, p99)
- Test LRU cache eviction under memory pressure
- Test concurrent agent processing (race conditions)
- Test aggregation with 1M+ entries
- Test memory usage stays under 100MB

// Benchmark tests
internal/stats/agent_aggregator_bench_test.go
- Benchmark ProcessAgentInvocation throughput
- Benchmark GetAgentStats retrieval time
- Benchmark memory allocation patterns
```

### Implementation Tasks (GREEN phase)
1. Enhance `internal/stats/agent_aggregator.go`
   - Implement `ProcessAgentInvocation()` with Welford's updates
   - Add t-digest integration for percentiles
   - Implement LRU cache for agent stats (max 1000 agents)
   - Add delegation chain tracking
   - Thread-safe operations with sync.RWMutex

2. Create memory management utilities
   - Bounded map with eviction policy
   - Memory pressure detection
   - Graceful degradation to sampling mode

### Refactor Tasks
- Extract statistical algorithms to reusable components
- Optimize hot paths identified by benchmarks
- Improve concurrency with fine-grained locking

### Success Criteria
- <200ms response for basic leaderboard
- Memory usage <100MB for 1M entries
- Zero data loss during processing
- All race detector tests pass

---

## PHASE 4: CLI Command Integration
**Duration**: 1-2 days
**Dependency**: Phase 3 (aggregation)

### Objectives
- Add agent subcommand to stats command
- Integrate with existing filter flags
- Support multiple output formats
- Maintain backward compatibility

### TDD Test Requirements
```go
// Test files to create (RED phase)
cmd/stats_agent_test.go
- Test agent subcommand registration
- Test --since flag filtering
- Test --format flag (table, json, csv)
- Test --global flag for all projects
- Test empty results handling
- Test help text and examples

// Integration tests
cmd/integration_test.go
- Test end-to-end agent stats flow
- Test with real JSONL files
- Test performance with large datasets
```

### Implementation Tasks (GREEN phase)
1. Create `cmd/stats_agent.go`
   - Define agent subcommand with Cobra
   - Parse command flags (--since, --format, --global)
   - Call agent aggregator
   - Pass results to display formatter

2. Update `cmd/stats.go`
   - Register agent subcommand
   - Add agent examples to help text
   - Ensure flag consistency

### Refactor Tasks
- Extract common flag parsing logic
- Unify error handling patterns
- Consolidate help text generation

### Success Criteria
- Agent subcommand appears in help
- All existing stats commands still work
- Consistent flag behavior across subcommands
- Integration tests pass

---

## PHASE 5: Display Formatting and Output
**Duration**: 1-2 days
**Dependency**: Phase 4 (CLI integration)

### Objectives
- Format agent leaderboard table
- JSON/CSV export formats
- Sparkline visualizations
- Success rate indicators

### TDD Test Requirements
```go
// Test files to create (RED phase)
internal/stats/agent_display_test.go
- Test table formatting with various data
- Test JSON marshaling correctness
- Test CSV header and row generation
- Test sparkline generation for trends
- Test empty data handling
- Test sorting by different columns
- Test truncation for long agent names
```

### Implementation Tasks (GREEN phase)
1. Create `internal/stats/agent_display.go`
   - Implement `FormatAgentTable()` with tabwriter
   - Implement `FormatAgentJSON()` with proper structure
   - Implement `FormatAgentCSV()` with headers
   - Add sparkline generation for usage trends
   - Color coding for success rates (green/yellow/red)

2. Integration with display pipeline
   - Hook into existing display formatter
   - Maintain consistent styling
   - Support terminal width detection

### Refactor Tasks
- Extract common formatting utilities
- Optimize string building performance
- Improve terminal compatibility

### Success Criteria
- Output matches existing stats format style
- All formats produce valid output
- Terminal colors work across platforms
- Handles edge cases gracefully

---

## PHASE 6: Basic Delegation Pattern Analysis
**Duration**: 2 days
**Dependency**: Phase 5 (display)

### Objectives
- Track agent co-occurrence patterns
- Identify delegation chains
- Display common agent combinations
- Session-level analysis

### TDD Test Requirements
```go
// Test files to enhance (RED phase)
internal/stats/agent_aggregator_test.go
- Test delegation chain construction
- Test agent co-occurrence tracking
- Test session boundary detection
- Test circular delegation handling
- Test maximum chain depth limits

internal/stats/agent_display_test.go
- Test delegation flow formatting
- Test co-occurrence table display
```

### Implementation Tasks (GREEN phase)
1. Enhance agent aggregator
   - Track agent transitions within sessions
   - Build delegation graphs
   - Calculate co-occurrence frequencies
   - Detect delegation patterns

2. Display delegation insights
   - Format delegation chains (A → B → C)
   - Show top agent combinations
   - Display delegation success rates

### Refactor Tasks
- Optimize graph traversal algorithms
- Improve memory usage for large chains
- Extract pattern detection logic

### Success Criteria
- Accurately tracks multi-agent workflows
- Handles complex delegation patterns
- Performance remains acceptable
- Clear visualization of patterns

---

## PHASE 7: Interactive TUI Mode (Optional Enhancement)
**Duration**: 3-4 days
**Dependency**: Phase 6 (complete CLI)

### Objectives
- Bubble Tea interactive dashboard
- Real-time sorting and filtering
- Keyboard navigation
- Progressive enhancement from CLI

### TDD Test Requirements
```go
// Test files to create (RED phase)
internal/ui/agent/dashboard_model_test.go
- Test model initialization
- Test update message handling
- Test view rendering
- Test keyboard navigation
- Test sorting state changes
- Test filter application

internal/ui/agent/leaderboard_component_test.go
- Test component rendering
- Test selection handling
- Test scroll behavior
```

### Implementation Tasks (GREEN phase)
1. Create TUI components
   - `internal/ui/agent/dashboard_model.go` - Main controller
   - `internal/ui/agent/leaderboard_component.go` - Agent ranking
   - `internal/ui/agent/timeline_component.go` - Temporal view
   - `internal/ui/agent/theme.go` - Consistent styling

2. Wire up interactive mode
   - Add --interactive flag handling
   - Launch Bubble Tea when flag present
   - Graceful fallback on terminal issues

### Refactor Tasks
- Extract reusable TUI components
- Optimize rendering performance
- Improve keyboard shortcut consistency

### Success Criteria
- TUI launches with --interactive flag
- All navigation works smoothly
- Graceful degradation to CLI
- <500ms response time for interactions

---

## PHASE 8: Performance Optimization and Polish
**Duration**: 1-2 days
**Dependency**: Phase 7 (feature complete)

### Objectives
- Profile and optimize hot paths
- Reduce memory allocations
- Improve startup time
- Polish user experience

### TDD Test Requirements
```go
// Benchmark improvements
*_bench_test.go files
- Achieve <200ms p95 response time
- Memory usage <100MB for 1M entries
- <10% processing overhead

// Stress tests
internal/stats/stress_test.go
- Test with 10M log entries
- Test with 1000+ unique agents
- Test memory bounds under pressure
```

### Implementation Tasks
1. Performance profiling
   - CPU profiling with pprof
   - Memory profiling and leak detection
   - Benchmark critical paths

2. Optimizations
   - Reduce allocations in hot loops
   - Implement object pooling if needed
   - Optimize regex compilation
   - Improve cache hit rates

### Refactor Tasks
- Simplify complex algorithms
- Remove redundant operations
- Improve code clarity

### Success Criteria
- All performance targets met
- No memory leaks detected
- Smooth user experience
- Clean, maintainable code

---

## Testing Strategy Summary

### Unit Test Coverage Requirements
- Core logic: 90%+ coverage
- Agent detection: All methods tested
- Aggregation: Statistical accuracy verified
- Display: All formats validated
- Edge cases: Comprehensive coverage

### Integration Test Scenarios
1. End-to-end agent analytics flow
2. Large dataset processing (1M+ entries)
3. Corrupted log handling
4. Multi-project aggregation
5. Time filtering accuracy

### Performance Test Targets
- Response time: <200ms (p95)
- Memory usage: <100MB (1M entries)
- Processing overhead: <10%
- Startup time: <100ms

### Test Data Requirements
- Sample JSONL files with various agent patterns
- Corrupted/malformed log entries
- Large dataset generators
- Multi-session test scenarios

---

## Risk Mitigation During Implementation

### Technical Risks
1. **Log format changes**: Implement version detection and adapters
2. **Memory pressure**: Bounded structures with graceful degradation
3. **Performance regression**: Continuous benchmarking at each phase
4. **Detection accuracy**: Multiple detection methods with confidence scoring

### Process Risks
1. **Scope creep**: Strict phase boundaries, defer enhancements
2. **Integration issues**: Incremental integration with existing code
3. **Test coverage gaps**: TDD ensures comprehensive testing
4. **Backward compatibility**: All changes behind feature flags initially

---

## Rollout Strategy

### Phase Rollout
1. **Alpha**: Phases 1-3 complete, internal testing only
2. **Beta**: Phases 4-6 complete, limited user testing with feedback
3. **GA**: Phase 7-8 complete, full release with documentation

### Feature Flags
- `STARTUP_AGENT_ANALYTICS_ENABLED`: Master feature toggle
- `STARTUP_AGENT_TUI_ENABLED`: Interactive mode toggle
- `STARTUP_AGENT_ADVANCED_PATTERNS`: Advanced detection methods

### Documentation Requirements
- Update CLI help text at each phase
- Add examples to README
- Create agent analytics guide
- Document performance characteristics

---

## Success Metrics

### Implementation Success
- All tests passing (unit, integration, performance)
- Performance targets achieved
- Zero regression in existing functionality
- Code coverage >90% for new code

### User Success (Post-Launch)
- 25% adoption within 30 days
- <3 critical bugs in first week
- Positive user feedback on usability
- Measurable improvement in agent delegation effectiveness

---

## Dependencies and Prerequisites

### Before Starting Implementation
1. Review existing stats codebase thoroughly
2. Set up test data and fixtures
3. Configure benchmark baselines
4. Ensure development environment ready

### External Dependencies
- Bubble Tea library (for TUI, Phase 7)
- T-digest library (for percentiles, Phase 3)
- No other external dependencies required

### Knowledge Requirements
- Welford's algorithm understanding
- T-digest percentile approximation
- Go concurrency patterns
- Bubble Tea TUI framework (Phase 7 only)