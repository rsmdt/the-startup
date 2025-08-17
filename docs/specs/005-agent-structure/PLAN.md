# Implementation Plan

## Phase 1: Context Analysis and Preparation
- [x] Read and analyze PRD.md for requirements
- [x] Read and analyze SDD.md for technical specifications
- [x] Review all 14 existing agent definitions in assets/agents/
- [x] Identify agents needing updates based on quality ratings
- [x] Review the agent-template.md for standardization pattern
- [x] **Validation**: Confirm all documentation reviewed and understood

## Phase 2: Fix Critical Issues in Existing Agents
- [x] Fix the-product-manager.md - Remove duplicate content (lines 72-94)
- [x] Add "use PROACTIVELY" or "MUST BE USED" patterns to all agent descriptions
- [x] Add `tools: inherit` field to all agents missing it
- [x] Ensure all agents have 3 usage examples in description
- [x] Verify all agents have "Previous Conversation History" section
- [x] **Validation**: `go fmt ./... && go vet ./...` - Ensure no syntax errors

## Phase 3: Enhance Lower-Rated Agents (7/10 ratings)
- [x] the-developer: Add error handling patterns, API design guidance
- [x] the-technical-writer: Add detailed process, template references
- [x] the-security-engineer: Add OWASP checklist, specific tools
- [x] the-devops-engineer: Add CI/CD tools (GitHub Actions, Terraform)
- [x] the-data-engineer: Add database patterns, migration strategies
- [x] the-tester: Add testing frameworks (Jest, Pytest), automation patterns
- [x] **Validation**: Review each agent follows template structure

## Phase 4: Implement New Agents
- [x] Install the-lead-developer.md in assets/agents/
- [x] Install the-ux-designer.md in assets/agents/
- [x] Install the-compliance-officer.md in assets/agents/
- [x] Verify new agents follow agent-template.md pattern
- [x] Ensure new agents have unique emoji personalities
- [x] **Validation**: `go test ./internal/log/...` - Verify agent filtering works

## Phase 5: Testing and Validation

### Agent Invocation Testing (Based on Research)
- [ ] Create test cases for "the-" prefix filtering in processor.go
- [ ] Test ShouldProcess function with new agent names
- [ ] Verify hook logging captures new agent invocations
- [ ] **Validation**: `go test -v ./internal/log/processor_test.go`

### Manual Testing (LLM-as-Judge + Human Validation)
Per research: "LLM-as-a-judge worked well... but human evaluation was essential" (Anthropic, 2025)
- [ ] Test each agent with sample prompts to verify invocation
- [ ] Verify "use PROACTIVELY" patterns trigger correctly
- [ ] Check task handoff blocks work between agents
- [ ] Test edge cases: vague requests, conflicting requirements
- [ ] **Validation**: Document test results in test-results.md

### Integration Testing
- [ ] Run `./the-startup install` to deploy updated agents
- [ ] Test with Claude Code to verify agent selection
- [ ] Verify session/agent ID tracking in logs
- [ ] Check `.the-startup/*/agent-instructions.jsonl` for captures
- [ ] **Validation**: `./the-startup validate agents` (when implemented)

## Phase 6: Documentation Updates
- [x] Update README.md with new agent descriptions
- [x] Create agent capability matrix showing all 17 agents
- [x] Add "Agent Discovery Guide" section
- [x] Document common agent collaboration patterns
- [x] Include training examples for new agents
- [x] Add troubleshooting section for agent invocation issues
- [x] **Validation**: Review documentation for completeness

## Phase 7: Final Validation
- [ ] Run full test suite: `go test ./...`
- [ ] Build binary: `go build -o the-startup`
- [ ] Test installation: `./the-startup install`
- [ ] Verify all agents appear in `.claude/agents/`
- [ ] Check agent quality ratings improved (target: all 8+/10)
- [ ] Confirm invocation rate improvement (target: 90%+)
- [ ] **Final Validation**: All acceptance criteria from PRD met

## Validation Checklist

### Code Quality
- [ ] All Go code formatted: `go fmt ./...`
- [ ] No vet issues: `go vet ./...`
- [ ] All tests passing: `go test ./...`
- [ ] Binary builds successfully: `go build -o the-startup`

### Agent Structure Compliance
- [ ] All 17 agents have YAML frontmatter with name, description, tools
- [ ] All agents have "use PROACTIVELY" or "MUST BE USED" in description
- [ ] All agents have exactly 3 usage examples
- [ ] All agents have unique emoji and personality
- [ ] All agents have Previous Conversation History section
- [ ] All agents follow output format with tasks blocks

### Functionality Verification
- [ ] Agent filtering works for "the-" prefix
- [ ] New agents properly logged by hook system
- [ ] Task handoffs include agent assignments
- [ ] Template placeholders ({{STARTUP_PATH}}) work
- [ ] No sub-agents calling sub-agents (removed from all)

### Documentation Completeness
- [ ] README.md includes all 17 agents
- [ ] Agent capability matrix created
- [ ] Discovery guide helps users find right agent
- [ ] Training examples provided
- [ ] Troubleshooting section included

## Anti-Patterns to Avoid

### Architecture Anti-Patterns
- ❌ Creating new architectural patterns when established ones exist
- ❌ Allowing sub-agents to invoke other sub-agents
- ❌ Adding agents without clear role distinction
- ❌ Modifying core hook processing logic
- ❌ Breaking "the-" prefix convention

### Agent Design Anti-Patterns
- ❌ Vague trigger descriptions without examples
- ❌ Missing personality or inconsistent emoji
- ❌ Overlapping responsibilities between agents
- ❌ Technology-specific implementations in architect
- ❌ Missing task handoff blocks

### Testing Anti-Patterns
- ❌ Testing implementation details instead of behavior
- ❌ Skipping manual validation of agent invocation
- ❌ Not testing with actual Claude Code
- ❌ Ignoring edge cases identified in SDD
- ❌ Assuming agents will be invoked without "PROACTIVELY"

### Process Anti-Patterns
- ❌ Skipping validation steps to move faster
- ❌ Implementing without testing agent invocation
- ❌ Adding backwards compatibility (forward-only change)
- ❌ Creating agents without real-world role mapping
- ❌ Ignoring community-reported issues