# Plugin Testing Strategy

**Goal:** Validate the Claude Code plugin works correctly before publishing

**Challenge:** We're building features (hooks) that don't exist in the current CLI, and need to test them in the actual Claude Code environment.

---

## Testing Approach: Incremental Validation

### Phase 1: Local Development Testing (Before Plugin)

**What:** Test individual components in isolation before migrating to plugin structure

#### 1.1 Test @ Notation with Existing Commands
- [ ] Create a test command in `.claude/commands/test-ref.md`
- [ ] Add `@rules/agent-delegation.md` reference
- [ ] Run the command in Claude Code
- [ ] Verify the rules file content is included
- [ ] **Success Criteria:** Claude reads and includes the referenced file

**Test Command Example:**
```markdown
---
description: Test @ notation file references
---

# Test Command

This should include the agent delegation rules:

@rules/agent-delegation.md

Did you see the rules above?
```

#### 1.2 Test Agent Frontmatter
- [ ] Copy one agent to `.claude/agents/test-agent.md`
- [ ] Launch it via `/agents`
- [ ] Verify it loads without errors
- [ ] Test basic functionality
- [ ] **Success Criteria:** Agent launches and executes correctly

#### 1.3 Develop Hooks Locally (BEFORE Plugin Migration)
- [ ] Create `hooks/` directory in project
- [ ] Develop `hooks/welcome.sh` with test output
- [ ] Test script execution manually: `bash hooks/welcome.sh`
- [ ] Verify JSON output is valid
- [ ] **Success Criteria:** Script runs and outputs valid JSON

**Manual Hook Test:**
```bash
# Test welcome.sh locally
bash hooks/welcome.sh

# Should output valid JSON:
# {
#   "systemMessage": "✓ Welcome",
#   "hookSpecificOutput": { ... }
# }
```

#### 1.4 Develop Spec Script Locally
- [ ] Create `scripts/spec.sh` in project
- [ ] Test script execution: `bash scripts/spec.sh test-feature`
- [ ] Verify directory creation works
- [ ] Verify TOML generation works
- [ ] Verify auto-increment works
- [ ] **Success Criteria:** Script creates correct directory structure

---

### Phase 2: Plugin Structure Testing (Migration Complete)

**What:** Test the actual plugin structure and installation

#### 2.1 Create Local Test Repository
- [ ] Create a separate test directory: `the-agentic-startup-test/`
- [ ] Copy plugin structure from main repo
- [ ] Initialize as Git repository
- [ ] Commit all files
- [ ] **Success Criteria:** Git repo ready for plugin installation

**Setup:**
```bash
# Create test plugin repo
mkdir -p ~/test-plugins/the-agentic-startup
cd ~/test-plugins/the-agentic-startup
git init
# Copy plugin files...
git add .
git commit -m "Initial plugin structure"
```

#### 2.2 Test Plugin Installation
- [ ] Install plugin from local path: `/plugin install ~/test-plugins/the-agentic-startup`
- [ ] Verify installation completes without errors
- [ ] Check Claude Code recognizes the plugin
- [ ] **Success Criteria:** Plugin installs successfully

#### 2.3 Verify Plugin Discovery
- [ ] Run `/agents` - verify all 50 agents appear
- [ ] Run `/help` - verify all 6 commands appear
- [ ] Check agent list for correct names
- [ ] **Success Criteria:** All components discovered

---

### Phase 3: Component Functionality Testing

**What:** Test each component works correctly in plugin context

#### 3.1 Test @ References in Commands
- [ ] Run `/s:specify test-feature`
- [ ] Verify Claude includes rules content (check response references agent delegation)
- [ ] Run `/s:analyze test-area`
- [ ] Verify rules loaded correctly
- [ ] **Success Criteria:** Commands execute and rules content is included

**Validation:**
- Look for agent delegation patterns in Claude's response
- Verify Claude mentions "FOCUS" and "EXCLUDE" (from rules)
- Check if cycle pattern steps are followed

#### 3.2 Test Hooks Configuration

**SessionStart Hook (Welcome Banner):**
- [ ] Create `hooks/hooks.json` with SessionStart config
- [ ] Test hook independently: Manually add to `~/.claude/settings.json`
```json
{
  "hooks": {
    "SessionStart": [{
      "type": "command",
      "command": "/path/to/hooks/welcome.sh"
    }]
  }
}
```
- [ ] Start new Claude Code session
- [ ] Verify banner appears in Claude's first response
- [ ] Check flag file created: `~/.the-startup/.plugin-initialized`
- [ ] Start second session - verify banner doesn't repeat
- [ ] **Success Criteria:** Banner shows once on first session

**UserPromptSubmit Hook (Statusline):**
- [ ] Test statusline script: `echo '{"workingDirectory":"/path"}' | bash hooks/statusline.sh`
- [ ] Verify git branch output
- [ ] Measure execution time: `time bash hooks/statusline.sh`
- [ ] Verify <10ms execution
- [ ] Add to settings.json manually for testing
- [ ] Verify statusline appears in Claude Code UI
- [ ] **Success Criteria:** Statusline displays git branch quickly

**Note:** Hooks in plugin context will use `${CLAUDE_PLUGIN_ROOT}` which gets resolved by Claude Code automatically. Test with absolute paths first, then confirm plugin resolution works.

#### 3.3 Test Spec Command
- [ ] Run `/s:spec test-feature`
- [ ] Verify script invocation works
- [ ] Check `docs/specs/spec-001-test-feature/` created
- [ ] Verify TOML file generated
- [ ] Verify templates copied
- [ ] Run `/s:spec another-feature`
- [ ] Verify spec-002 created (auto-increment)
- [ ] **Success Criteria:** Spec directories created with correct numbering

#### 3.4 Test Agent Execution
- [ ] Launch `/agents` and select `the-chief`
- [ ] Provide a test task
- [ ] Verify agent responds appropriately
- [ ] Test agent delegation (chief calling another agent)
- [ ] Test parallel agent execution
- [ ] **Success Criteria:** Agents work identically to CLI version

#### 3.5 Test Templates Access
- [ ] Run `/s:init`
- [ ] Verify command can access `templates/` directory
- [ ] Check if templates are read correctly
- [ ] **Success Criteria:** Init command accesses templates successfully

---

### Phase 4: Integration Testing

**What:** Test complete workflows end-to-end

#### 4.1 Complete Feature Specification Workflow
- [ ] Run `/s:specify new-login-feature`
- [ ] Verify rules loaded (agent delegation, cycle pattern)
- [ ] Verify agents launched appropriately
- [ ] Verify TodoWrite tracking works
- [ ] Check output organization
- [ ] **Success Criteria:** Complete spec workflow executes correctly

#### 4.2 Complete Spec Generation Workflow
- [ ] Run `/s:spec login-feature`
- [ ] Verify directory created
- [ ] Edit specification.toml
- [ ] Run `/s:implement spec-001-login-feature`
- [ ] Verify implementation workflow works
- [ ] **Success Criteria:** Full spec → implement workflow works

#### 4.3 Hook Integration Testing
- [ ] Start fresh session (SessionStart)
- [ ] Verify welcome banner
- [ ] Run several commands (UserPromptSubmit)
- [ ] Verify statusline updates
- [ ] Check performance doesn't degrade
- [ ] **Success Criteria:** Hooks work smoothly with normal workflow

---

### Phase 5: Edge Cases & Error Handling

#### 5.1 Test File Reference Edge Cases
- [ ] Test @ reference to non-existent file
- [ ] Test @ reference with wrong path
- [ ] Verify appropriate error messages
- [ ] **Success Criteria:** Graceful failure with clear errors

#### 5.2 Test Hook Failures
- [ ] Break welcome.sh (invalid JSON)
- [ ] Verify Claude Code handles gracefully
- [ ] Test hook timeout scenarios
- [ ] **Success Criteria:** Hooks don't break Claude Code on failure

#### 5.3 Test Spec Script Edge Cases
- [ ] Run spec with invalid feature name
- [ ] Test when specs/ directory doesn't exist
- [ ] Test with existing spec number
- [ ] **Success Criteria:** Script handles errors gracefully

---

## Testing Tools & Utilities

### Manual Testing Checklist
Create a testing checklist document to track:
- [ ] Create `TESTING-CHECKLIST.md` with all test cases
- [ ] Mark each test as PASS/FAIL/SKIP
- [ ] Document any failures with details
- [ ] Track fixes and re-tests

### Test Fixtures
Prepare test data:
- [ ] Sample feature descriptions for `/s:specify`
- [ ] Sample spec directories for `/s:implement`
- [ ] Sample analysis targets for `/s:analyze`
- [ ] Expected outputs for validation

### Debug Mode
Add debug output to scripts:
- [ ] Add `set -x` to bash scripts for debugging
- [ ] Add logging to spec.sh for troubleshooting
- [ ] Create verbose mode for hooks

---

## What We're Testing vs What We're Not

### ✅ Testing (In Scope)
- Plugin installs correctly
- Agents load and execute
- Commands execute with @ references
- Hooks fire at correct times
- Spec script creates directories
- Templates are accessible
- File references resolve correctly

### ❌ Not Testing (Out of Scope)
- Claude's LLM capabilities (assume they work)
- TodoWrite functionality (Claude Code built-in)
- Agent reasoning quality (subjective)
- Full cross-platform testing (focus on primary platform first)

---

## Testing Timeline

### Before Migration (1 day)
- Test @ notation with existing commands
- Develop and test hooks locally
- Develop and test spec script locally
- Test agent frontmatter compatibility

### During Migration (ongoing)
- Test each phase as completed
- Validate checkpoints
- Fix issues before moving forward

### After Migration (1 day)
- Full plugin installation test
- Complete functionality testing
- Integration workflows
- Edge case validation

---

## Rollback Plan

If testing reveals critical issues:

1. **Minor Issues:** Fix in plugin repo, test again
2. **Major Issues:** Document blockers, keep CLI as primary
3. **Show Stoppers:** Revert to CLI, document findings

---

## Success Criteria Summary

**Plugin is ready for release when:**
- [ ] Installs via `/plugin install` without errors
- [ ] All 50 agents discoverable and functional
- [ ] All 6 commands execute correctly
- [ ] @ references work for rules files
- [ ] Hooks execute without breaking Claude Code
- [ ] Spec command creates directories correctly
- [ ] No regressions from CLI version
- [ ] Documentation complete and accurate

---

## Current Testing Status

- [ ] Phase 1: Local Development Testing
- [ ] Phase 2: Plugin Structure Testing
- [ ] Phase 3: Component Functionality Testing
- [ ] Phase 4: Integration Testing
- [ ] Phase 5: Edge Cases Testing

**Next:** Start Phase 1 - Test @ notation and develop hooks locally BEFORE migration
