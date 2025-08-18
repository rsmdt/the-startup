# Implementation Plan

## Context Documents

### Automatic Context Discovery

The orchestrator will search for and analyze the following documents in priority order:

1. **Primary Specifications** (same directory as PLAN.md):
   - **SDD.md** ✓: Technical architecture for output-style integration
   - **PRD.md** ✗: Not created (Level 2 complexity)
   - **BRD.md** ✗: Not created (Level 2 complexity)

2. **Secondary Context**:
   - Agent definitions in `assets/agents/*.md`
   - Claude Code documentation on output-styles
   - Pattern documentation at `docs/patterns/agent-personality-prompt-engineering.md`

3. **Codebase Patterns**:
   - Installation flow in `cmd/install.go`
   - Hook processing in `cmd/log.go`
   - Settings template in `assets/settings.json`

### Context Synthesis

**Available Specifications:**
- ✓ SDD.md: `/Users/irudi/Code/personal/the-startup/docs/specs/007-agent-output-style-personalization/SDD.md`
- ✓ Pattern Documentation: `/Users/irudi/Code/personal/the-startup/docs/patterns/agent-personality-prompt-engineering.md`
- ✗ PRD.md: Not required for Level 2 specification
- ✗ BRD.md: Not required for Level 2 specification

**Key Implementation Requirements:**

From SDD.md:
- Hybrid meta-style architecture with context detection
- Style generator interface for creating output-styles
- 5-minute context caching for performance
- Installation integration for automatic style generation
- Personality injection through Task tool prompts

From Pattern Documentation:
- Four agent archetype patterns (Technical, Business, Creative, Process)
- Four-layer personality architecture
- Visual signatures for identity anchoring
- Dynamic personality scaling based on context

From Claude Code Documentation:
- Output-styles stored in ~/.claude/output-styles or .claude/output-styles
- Styles referenced in settings.local.json
- Styles directly modify system prompt
- Static per session (no dynamic switching)

**Implementation Priorities Based on Context:**
1. Create style generator with meta-style template
2. Enhance agent definitions with stronger personality markers
3. Integrate style generation into installation process
4. Implement context detection logic in meta-style
5. Add performance monitoring and caching

## Phase 1: Foundation - Style Generator Infrastructure

- [ ] **Read and analyze existing codebase structure** [`complexity: low`]
  - Review internal/installer/installer.go for integration points
  - Understand asset embedding and file copying patterns
  - Identify settings.json modification approach

- [ ] **Create style generator package structure** [`complexity: medium`] [`agent: the-developer`]
  - Create internal/styles/ directory
  - Implement generator.go with StyleGenerator interface
  - Create templates.go for style templates
  - Add meta.go for meta-style logic

- [ ] **Design meta-style template** [`complexity: high`] [`review: prompt-engineering, personality`]
  - Create assets/templates/meta-style.md
  - Implement context detection patterns
  - Add personality overlay system
  - Include signature-based identity anchoring

- [ ] **Validation**: `go fmt ./internal/styles/... && go vet ./internal/styles/...`

## Phase 2: Agent Personality Enhancement

- [ ] **Analyze existing agent definitions** [`complexity: medium`] [`parallel: true`]
  - Review all 17 agent markdown files
  - Identify current personality patterns
  - Map agents to archetypes (Technical/Business/Creative/Process)

- [ ] **Enhance agent personalities with stronger markers** [`complexity: high`] [`review: personality, consistency`] [`agent: the-prompt-engineer`]
  - Add visual signatures to each agent
  - Strengthen voice characteristics
  - Define behavioral directives
  - Implement micro-personality expressions

- [ ] **Create individual style templates** [`complexity: medium`] [`parallel: true`]
  - Generate style for each archetype
  - Test personality differentiation
  - Validate Claude Code compatibility

- [ ] **Validation**: Manual review of enhanced personalities for consistency

## Phase 3: Style Generation Implementation

- [ ] **Implement style generator core logic** [`complexity: high`] [`review: architecture`] [`agent: the-developer`]
  - GenerateMetaStyle function
  - GenerateIndividualStyle function
  - Template processing with placeholders
  - Personality profile mapping

- [ ] **Add context detection logic** [`complexity: high`] [`review: performance`]
  - Parse for subagent_type parameters
  - Implement signature pattern detection
  - Add vocabulary analysis
  - Create 5-minute context cache

- [ ] **Create style installation logic** [`complexity: medium`]
  - InstallStyles function
  - Directory creation and validation
  - File writing with proper permissions
  - Error handling and rollback

- [ ] **Validation**: `go test ./internal/styles/...`

## Phase 4: Installation Integration

- [ ] **Modify installer to generate styles** [`complexity: medium`] [`review: integration`] [`agent: the-developer`]
  - Add style generation step to install flow
  - Create output-styles directory structure
  - Generate meta-style from agents
  - Optional: Generate individual styles

- [ ] **Update settings.json template** [`complexity: low`]
  - Add output_style configuration
  - Set meta-style as default
  - Include style directory path

- [ ] **Implement settings.json modification** [`complexity: medium`] [`review: error-handling`]
  - Parse existing settings
  - Merge output-style configuration
  - Preserve user customizations
  - Handle malformed JSON gracefully

- [ ] **Validation**: `./the-startup install && ls ~/.claude/output-styles/`

## Phase 5: Hook Processing Enhancement

- [ ] **Add personality tracking to hook processor** [`complexity: medium`] [`agent: the-developer`]
  - Log active agent personalities
  - Track personality switches
  - Monitor consistency metrics

- [ ] **Implement performance monitoring** [`complexity: medium`] [`review: performance`]
  - Measure style application latency
  - Track cache hit/miss rates
  - Log context detection accuracy

- [ ] **Add debug mode for personality system** [`complexity: low`]
  - DEBUG_PERSONALITY environment variable
  - Log detection decisions
  - Output personality overlays

- [ ] **Validation**: `DEBUG_HOOKS=1 ./the-startup log --assistant < test-input.json`

## Phase 6: Testing and Documentation

- [ ] **Create comprehensive test suite** [`complexity: high`] [`review: coverage`] [`agent: the-tester`]
  - Unit tests for style generator
  - Integration tests for installation
  - Personality consistency tests
  - Performance benchmarks

- [ ] **Write user documentation** [`complexity: medium`] [`agent: the-technical-writer`]
  - How to use output-styles
  - Switching between styles
  - Customizing personalities
  - Troubleshooting guide

- [ ] **Create example configurations** [`complexity: low`]
  - Sample meta-style variations
  - Individual agent style examples
  - Custom personality templates

- [ ] **Validation**: `go test -cover ./...`

## Phase 7: Integration and Final Validation

- [ ] **Full system integration test** [`complexity: high`] [`review: integration, performance`] [`validation: manual`]
  - Fresh installation test
  - Style generation verification
  - Personality application test
  - Context switching validation

- [ ] **Performance validation** [`complexity: medium`] [`validation: performance`]
  - Measure style application overhead
  - Verify <100ms latency target
  - Test cache effectiveness
  - Check memory usage

- [ ] **Compatibility testing** [`complexity: medium`] [`review: compatibility`]
  - Test with Claude Code 0.8.0+
  - Verify all features remain functional
  - Check for prompt conflicts
  - Validate with different commands

- [ ] **Final validation**: Full installation and usage test
  ```bash
  # Clean installation
  rm -rf ~/.claude/output-styles
  go build -o the-startup
  ./the-startup install
  
  # Verify styles
  ls ~/.claude/output-styles/
  cat ~/.claude/settings.local.json | grep output_style
  
  # Test with Claude Code
  # Manual: Use /output-style command
  # Manual: Invoke different agents and verify personalities
  ```

## Validation Checklist

### Code Quality
- [ ] Go formatting: `go fmt ./...`
- [ ] Go vetting: `go vet ./...`
- [ ] Test coverage: `go test -cover ./...`
- [ ] Build verification: `go build -o the-startup`

### Functionality
- [ ] Meta-style generates correctly
- [ ] Individual styles create properly
- [ ] Installation integrates styles
- [ ] Settings.json updates correctly
- [ ] Personalities apply consistently

### Integration
- [ ] Installer workflow unchanged for users
- [ ] Hooks continue processing normally
- [ ] Claude Code features remain functional
- [ ] Style switching works as expected

### Performance
- [ ] Style application <100ms
- [ ] Context cache working (5-minute window)
- [ ] No memory leaks in caching
- [ ] Minimal overhead for detection

### User Experience
- [ ] Each agent feels distinct
- [ ] Personalities remain consistent
- [ ] Transitions are clean
- [ ] Fallback behavior works

## Anti-Patterns to Avoid

### Architecture Anti-Patterns
- ❌ Creating complex dynamic switching mechanisms (output-styles are static)
- ❌ Modifying Claude Code's core behavior
- ❌ Over-engineering the detection logic
- ❌ Breaking existing agent functionality

### Integration Anti-Patterns
- ❌ Hardcoding style paths or names
- ❌ Assuming specific Claude Code versions
- ❌ Overwriting user's existing styles
- ❌ Forcing style usage without opt-out

### Personality Anti-Patterns
- ❌ Making personalities too extreme or distracting
- ❌ Losing core functionality for personality
- ❌ Creating inconsistent personality signals
- ❌ Mixing personality archetypes confusingly

### Testing Anti-Patterns
- ❌ Testing personality "feelings" instead of markers
- ❌ Skipping integration tests with Claude Code
- ❌ Not testing fallback scenarios
- ❌ Ignoring performance requirements

### Process Anti-Patterns
- ❌ Implementing without testing detection logic
- ❌ Assuming all agents need the same intensity
- ❌ Skipping user documentation
- ❌ Deploying without compatibility testing