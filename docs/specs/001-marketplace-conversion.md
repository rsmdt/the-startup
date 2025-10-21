# Marketplace Conversion Plan

**Spec ID**: 001
**Feature**: Convert monolithic plugin to 4-plugin marketplace
**Date Created**: 2025-10-21
**Last Updated**: 2025-10-21
**Status**: 🟡 In Progress

---

## 📊 Overall Progress

**Completion**: 50% (3/6 phases completed)

| Phase | Status | Completion |
|-------|--------|------------|
| Phase 1: Repository Restructuring | ✅ Completed | 5/5 tasks |
| Phase 2: Fix References & Manifests | ✅ Completed | 6/6 tasks |
| Phase 3: Marketplace Manifest | ✅ Completed | 3/3 tasks |
| Phase 4: Local Testing | ⬜ Not Started | 0/13 tests |
| Phase 5: Documentation | ⬜ Not Started | 0/5 tasks |
| Phase 6: Git & Publishing | ⬜ Not Started | 0/5 tasks |

**Legend**: ⬜ Not Started | 🟡 In Progress | ✅ Completed

---

## 🚀 Quick Resume

**Current Phase**: Phase 4 - Local Testing
**Next Action**: Add marketplace and install plugins locally
**Estimated Time Remaining**: 5.5 hours

**To Resume Work**:
1. Open this file and check "Current Phase" above
2. Jump to the phase section below
3. Start with the first unchecked task
4. Update status indicators as you complete tasks
5. Update "Overall Progress" table when phase completes

---

## 📝 How to Track Progress

As you work through this plan, update the following:

### When Starting a Phase
1. Change phase **Status** from `⬜ Not Started` to `🟡 In Progress`
2. Update "Overall Progress" table at the top
3. Update "Quick Resume" section with current phase

### During a Phase
1. Check off tasks as you complete them: `- [ ]` → `- [x]`
2. Update **Completion** count (e.g., `0/5 tasks` → `2/5 tasks`)
3. Check validation items as they pass

### When Completing a Phase
1. Verify all tasks are checked
2. Verify all validation items pass
3. Change phase **Status** from `🟡 In Progress` to `✅ Completed`
4. Update "Overall Progress" table: increment completed phases
5. Update "Quick Resume" to next phase
6. Update **Last Updated** date at the top

### Example Progress Update

Before:
```markdown
**Status**: ⬜ Not Started
**Completion**: 0/5 tasks
- [ ] **Task 1.1**: Create directories
```

After completing Task 1.1:
```markdown
**Status**: 🟡 In Progress
**Completion**: 1/5 tasks
- [x] **Task 1.1**: Create directories
```

After completing all tasks:
```markdown
**Status**: ✅ Completed
**Completion**: 5/5 tasks
- [x] **Task 1.1**: Create directories
```

---

## Overview

Convert `the-startup` from a single monolithic plugin into a marketplace with 4 modular plugins organized by component type:
- `the-startup-commands` - Workflow orchestration commands
- `the-startup-agents` - 50 specialist agents
- `the-startup-output-style` - Custom Claude Code output style
- `the-startup-statusline` - Productivity hooks (welcome, statusline)

This enables users to install only the components they need while maintaining full functionality when all plugins are installed together.

---

## Current State Analysis

### Current Plugin Structure

```
the-startup/
├── .claude-plugin/
│   ├── plugin.json
│   └── marketplace.json
├── agents/                    # 50 agents (2 top-level + 9 role directories)
│   ├── the-chief.md
│   ├── the-meta-agent.md
│   ├── the-analyst/           # 3 agents
│   ├── the-architect/         # 5 agents
│   ├── the-designer/          # 4 agents
│   ├── the-ml-engineer/       # 4 agents
│   ├── the-mobile-engineer/   # 3 agents
│   ├── the-platform-engineer/ # 7 agents
│   ├── the-qa-engineer/       # 3 agents
│   ├── the-security-engineer/ # 3 agents
│   └── the-software-engineer/ # 5 agents
├── commands/
│   └── s/
│       ├── analyze.md         # References @rules/
│       ├── specify.md         # References @rules/
│       ├── implement.md       # References @rules/
│       ├── refactor.md        # References @rules/
│       ├── init.md            # References templates/
│       └── spec.md            # Uses ${CLAUDE_PLUGIN_ROOT}/scripts/spec.py
├── rules/                     # Referenced by commands via @rules/
│   ├── agent-delegation.md
│   ├── cycle-pattern.md
│   └── agent-creation-principles.md
├── templates/                 # Used by spec.py and init.md
│   ├── product-requirements.md
│   ├── solution-design.md
│   ├── implementation-plan.md
│   ├── definition-of-ready.md
│   ├── definition-of-done.md
│   └── task-definition-of-done.md
├── scripts/                   # Used by spec.md
│   ├── spec.py
│   ├── spec.sh
│   ├── statusline.sh
│   └── statusline.ps1
├── hooks/
│   ├── hooks.json
│   └── welcome.sh
└── output-styles/
    └── the-startup.md
```

### Dependency Analysis

**Commands Dependencies**:
- `analyze.md`, `specify.md`, `implement.md`, `refactor.md` → `@rules/agent-delegation.md`, `@rules/cycle-pattern.md`
- `init.md` → `templates/` directory
- `spec.md` → `${CLAUDE_PLUGIN_ROOT}/scripts/spec.py` → `templates/` directory
- All commands delegate to agents (functional dependency, not technical)

**Output Style Dependencies**:
- `the-startup.md` (line 66) → `@{{STARTUP_PATH}}/rules/agent-delegation.md`
- This placeholder needs updating to `@rules/agent-delegation.md`
- Requires `rules/agent-delegation.md` (225 lines) in the same plugin

**Agents Dependencies**:
- No @ references or file dependencies
- `the-meta-agent.md` mentions `rules/agent-creation-principles.md` in documentation (text only, not an @ reference)
- Completely independent

**Hooks Dependencies**:
- No @ references or file dependencies
- Completely independent

**Technical Constraints**:
- `@rules/` references MUST be in the same plugin (@ references are plugin-scoped)
- `${CLAUDE_PLUGIN_ROOT}/scripts/` references MUST be in the same plugin
- `spec.py` references `templates/` at plugin root (MUST be in same plugin)

**Rules Duplication Decision**:
- `agent-delegation.md` needed by both `commands` plugin (4 commands) AND `output-style` plugin (1 reference)
- **Solution**: Duplicate `agent-delegation.md` in both plugins
- Other rules (`cycle-pattern.md`, `agent-creation-principles.md`) only needed in commands plugin

---

## Target Marketplace Structure

### Repository Layout

```
the-startup/
├── .claude-plugin/
│   └── marketplace.json           # Marketplace catalog
├── plugins/
│   ├── the-startup-commands/
│   │   ├── .claude-plugin/
│   │   │   └── plugin.json
│   │   ├── commands/
│   │   │   └── s/                 # All 6 commands
│   │   ├── rules/                 # 3 rules (needed by 4 commands)
│   │   │   ├── agent-delegation.md
│   │   │   ├── cycle-pattern.md
│   │   │   └── agent-creation-principles.md
│   │   ├── scripts/               # Required by spec.md
│   │   │   ├── spec.py
│   │   │   └── spec.sh
│   │   └── templates/             # Required by spec.py and init.md
│   │       ├── product-requirements.md
│   │       ├── solution-design.md
│   │       ├── implementation-plan.md
│   │       ├── definition-of-ready.md
│   │       ├── definition-of-done.md
│   │       └── task-definition-of-done.md
│   ├── the-startup-agents/
│   │   ├── .claude-plugin/
│   │   │   └── plugin.json
│   │   └── agents/                # All 50 agents
│   │       ├── the-chief.md
│   │       ├── the-meta-agent.md
│   │       ├── the-analyst/
│   │       ├── the-architect/
│   │       ├── the-designer/
│   │       ├── the-ml-engineer/
│   │       ├── the-mobile-engineer/
│   │       ├── the-platform-engineer/
│   │       ├── the-qa-engineer/
│   │       ├── the-security-engineer/
│   │       └── the-software-engineer/
│   ├── the-startup-output-style/
│   │   ├── .claude-plugin/
│   │   │   └── plugin.json
│   │   ├── rules/                 # 1 rule (duplicated from commands)
│   │   │   └── agent-delegation.md
│   │   └── output-styles/
│   │       └── the-startup.md
│   └── the-startup-statusline/
│       ├── .claude-plugin/
│       │   └── plugin.json
│       └── hooks/
│           ├── hooks.json
│           ├── welcome.sh
│           ├── statusline.sh
│           └── statusline.ps1
└── README.md
```

---

## Plugin Specifications

### Plugin 1: the-startup-commands

**Purpose**: Workflow orchestration commands for specification, implementation, and refactoring

**Contains**:
- Commands: `/s:analyze`, `/s:specify`, `/s:implement`, `/s:refactor`, `/s:init`, `/s:spec`
- Rules: `agent-delegation.md`, `cycle-pattern.md`, `agent-creation-principles.md`
- Scripts: `spec.py`, `spec.sh`
- Templates: 6 template files

**plugin.json**:
```json
{
  "name": "the-startup-commands",
  "version": "2.0.0",
  "description": "Workflow orchestration commands for agentic software development",
  "author": {
    "name": "Rudolf S."
  },
  "homepage": "https://github.com/rsmdt/the-startup",
  "repository": "https://github.com/rsmdt/the-startup",
  "license": "MIT",
  "keywords": [
    "workflow",
    "commands",
    "specification",
    "implementation",
    "refactoring"
  ],
  "commands": ["commands/"]
}
```

**Functional Note**: Commands delegate to agents but don't technically require them. Users could use commands with their own custom agents or the standard Claude Code Task tool.

---

### Plugin 2: the-startup-agents

**Purpose**: 50 specialized agents across 9 professional roles

**Contains**:
- Orchestration: The Chief (routing), The Meta Agent (generation)
- Specialist roles: Analyst, Architect, Designer, ML Engineer, Mobile Engineer, Platform Engineer, QA Engineer, Security Engineer, Software Engineer

**plugin.json**:
```json
{
  "name": "the-startup-agents",
  "version": "2.0.0",
  "description": "50 specialized agents for software development across 9 professional roles",
  "author": {
    "name": "Rudolf S."
  },
  "homepage": "https://github.com/rsmdt/the-startup",
  "repository": "https://github.com/rsmdt/the-startup",
  "license": "MIT",
  "keywords": [
    "agents",
    "specialists",
    "roles",
    "delegation"
  ],
  "agents": ["agents/"]
}
```

**Agent Count by Role**:
- The Analyst: 3 agents
- The Architect: 5 agents
- The Designer: 4 agents
- The ML Engineer: 4 agents
- The Mobile Engineer: 3 agents
- The Platform Engineer: 7 agents
- The QA Engineer: 3 agents
- The Security Engineer: 3 agents
- The Software Engineer: 5 agents
- Orchestration: 2 agents (The Chief, The Meta Agent)

---

### Plugin 3: the-startup-output-style

**Purpose**: Custom Claude Code output style for The Startup methodology

**Contains**:
- Output style: `the-startup.md`
- Rules: `agent-delegation.md` (duplicated from commands plugin - required by @ reference in output style)

**Note**: The output style references `@rules/agent-delegation.md` on line 66. This file must be present in the same plugin for the @ reference to resolve.

**plugin.json**:
```json
{
  "name": "the-startup-output-style",
  "version": "2.0.0",
  "description": "Custom Claude Code output style for The Startup methodology",
  "author": {
    "name": "Rudolf S."
  },
  "homepage": "https://github.com/rsmdt/the-startup",
  "repository": "https://github.com/rsmdt/the-startup",
  "license": "MIT",
  "keywords": [
    "output-style",
    "formatting",
    "ux"
  ],
  "outputStyles": ["output-styles/"]
}
```

---

### Plugin 4: the-startup-statusline

**Purpose**: Productivity hooks for welcome banner and statusline

**Contains**:
- Hooks: SessionStart (welcome banner), UserPromptSubmit (statusline)
- Scripts: `welcome.sh`, `statusline.sh`, `statusline.ps1`

**plugin.json**:
```json
{
  "name": "the-startup-statusline",
  "version": "2.0.0",
  "description": "Productivity hooks: welcome banner and git statusline",
  "author": {
    "name": "Rudolf S."
  },
  "homepage": "https://github.com/rsmdt/the-startup",
  "repository": "https://github.com/rsmdt/the-startup",
  "license": "MIT",
  "keywords": [
    "hooks",
    "statusline",
    "productivity",
    "ux"
  ],
  "hooks": "hooks/hooks.json"
}
```

---

## Marketplace Configuration

**marketplace.json**:
```json
{
  "name": "the-startup-marketplace",
  "owner": {
    "name": "Rudolf S.",
    "email": "rudolf@example.com"
  },
  "metadata": {
    "description": "Comprehensive agentic software development framework for Claude Code",
    "version": "2.0.0",
    "homepage": "https://github.com/rsmdt/the-startup"
  },
  "plugins": [
    {
      "name": "the-startup-commands",
      "source": "./plugins/the-startup-commands",
      "description": "Workflow orchestration commands for agentic software development",
      "version": "2.0.0",
      "category": "workflow",
      "keywords": ["workflow", "commands", "specification", "implementation"]
    },
    {
      "name": "the-startup-agents",
      "source": "./plugins/the-startup-agents",
      "description": "50 specialized agents for software development across 9 professional roles",
      "version": "2.0.0",
      "category": "agents",
      "keywords": ["agents", "specialists", "roles", "delegation"]
    },
    {
      "name": "the-startup-output-style",
      "source": "./plugins/the-startup-output-style",
      "description": "Custom Claude Code output style for The Startup methodology",
      "version": "2.0.0",
      "category": "productivity",
      "keywords": ["output-style", "formatting", "ux"]
    },
    {
      "name": "the-startup-statusline",
      "source": "./plugins/the-startup-statusline",
      "description": "Productivity hooks: welcome banner and git statusline",
      "version": "2.0.0",
      "category": "productivity",
      "keywords": ["hooks", "statusline", "productivity"]
    }
  ]
}
```

---

## Migration Plan

---

### Phase 1: Repository Restructuring

**Status**: ✅ Completed
**Duration**: 2 hours
**Completion**: 5/5 tasks

#### Tasks

- [x] **Task 1.1**: Create `plugins/` directory structure
- [x] **Task 1.2**: Create subdirectories for each plugin with `.claude-plugin/` directories
- [x] **Task 1.3**: Move files to respective plugin directories
- [x] **Task 1.4**: Duplicate `agent-delegation.md` to output-style plugin
- [x] **Task 1.5**: Verify directory structure and clean up empty directories

#### Detailed Steps

```bash
# Create plugin directories
mkdir -p plugins/the-startup-commands/.claude-plugin
mkdir -p plugins/the-startup-agents/.claude-plugin
mkdir -p plugins/the-startup-output-style/.claude-plugin
mkdir -p plugins/the-startup-statusline/.claude-plugin

# Move commands plugin files
mv commands plugins/the-startup-commands/
mv rules plugins/the-startup-commands/
mv templates plugins/the-startup-commands/
mkdir plugins/the-startup-commands/scripts
mv scripts/spec.py plugins/the-startup-commands/scripts/
mv scripts/spec.sh plugins/the-startup-commands/scripts/

# Move agents plugin files
mv agents plugins/the-startup-agents/

# Move output-style plugin files
mv output-styles plugins/the-startup-output-style/

# Duplicate agent-delegation.md to output-style plugin (needed by @ reference)
mkdir -p plugins/the-startup-output-style/rules
cp plugins/the-startup-commands/rules/agent-delegation.md plugins/the-startup-output-style/rules/

# Move statusline plugin files
mkdir -p plugins/the-startup-statusline/hooks
mv hooks/hooks.json plugins/the-startup-statusline/hooks/
mv hooks/welcome.sh plugins/the-startup-statusline/hooks/
mv scripts/statusline.sh plugins/the-startup-statusline/hooks/
mv scripts/statusline.ps1 plugins/the-startup-statusline/hooks/

# Clean up empty directories
rmdir hooks scripts
```

#### Validation Checklist

- [x] All component directories moved to plugin subdirectories
- [x] No orphaned files in root
- [x] Each plugin has `.claude-plugin/` directory
- [x] `agent-delegation.md` exists in both `the-startup-commands/rules/` and `the-startup-output-style/rules/`
- [x] Both copies of `agent-delegation.md` are identical

**✅ Phase 1 Complete When**: All tasks checked, all validations passed

---

### Phase 2: Fix File References and Create Plugin Manifests

**Status**: ✅ Completed
**Duration**: 1.5 hours
**Completion**: 6/6 tasks

#### Tasks

- [x] **Task 2.1**: Fix output-style placeholder reference (`@{{STARTUP_PATH}}/rules/` → `@rules/`)
- [x] **Task 2.2**: Create `plugin.json` for `the-startup-commands`
- [x] **Task 2.3**: Create `plugin.json` for `the-startup-agents`
- [x] **Task 2.4**: Create `plugin.json` for `the-startup-output-style`
- [x] **Task 2.5**: Create `plugin.json` for `the-startup-statusline`
- [x] **Task 2.6**: Validate all manifests

#### Detailed Steps

**Step 2.1: Fix Output Style Placeholder**

Update `plugins/the-startup-output-style/output-styles/the-startup.md` line 66:
- Change: `@{{STARTUP_PATH}}/rules/agent-delegation.md`
- To: `@rules/agent-delegation.md`

```bash
# Using sed to fix the placeholder (macOS)
sed -i '' 's|@{{STARTUP_PATH}}/rules/agent-delegation.md|@rules/agent-delegation.md|g' \
  plugins/the-startup-output-style/output-styles/the-startup.md

# Or using your preferred text editor
# Line 66: @{{STARTUP_PATH}}/rules/agent-delegation.md → @rules/agent-delegation.md
```

#### Validation Checklist

- [x] All manifests have required fields (`name`, `version`, `description`, `author`, etc.)
- [x] Output style references `@rules/agent-delegation.md` (no `{{STARTUP_PATH}}` placeholder)
- [x] `agent-delegation.md` is duplicated and identical in both plugins
- [x] Component paths are correct and relative (`commands/`, `agents/`, etc.)
- [x] All versions use semantic versioning (2.0.0)
- [x] Keywords are descriptive and relevant
- [x] No validation errors from `claude plugin validate`

**Note on Rules Duplication**:
- `agent-delegation.md` is intentionally duplicated in both `commands` and `output-style` plugins
- This is required because @ references are plugin-scoped
- The files should be kept in sync if rules are updated
- Consider this an acceptable trade-off for plugin modularity

**✅ Phase 2 Complete When**: All tasks checked, all validations passed, no placeholder references remain

---

### Phase 3: Marketplace Manifest Creation

**Status**: ✅ Completed
**Duration**: 30 minutes
**Completion**: 3/3 tasks

#### Tasks

- [x] **Task 3.1**: Update `.claude-plugin/marketplace.json` with all 4 plugins
- [x] **Task 3.2**: Configure metadata, owner, and plugin details
- [x] **Task 3.3**: Validate marketplace.json syntax and structure

#### Validation Checklist

- [x] All 4 plugins listed in marketplace.json
- [x] Source paths point to correct directories (`./plugins/plugin-name`)
- [x] Metadata is complete (name, description, version, homepage)
- [x] Owner information is present
- [x] JSON syntax is valid (no trailing commas, proper quotes)
- [x] Plugin versions match individual plugin.json files

**✅ Phase 3 Complete When**: All tasks checked, all validations passed

---

### Phase 4: Local Testing

**Status**: ⬜ Not Started
**Duration**: 2 hours
**Completion**: 0/13 tests

#### Tasks

- [ ] **Task 4.1**: Add local marketplace to Claude Code
- [ ] **Task 4.2**: Install all 4 plugins from marketplace
- [ ] **Task 4.3**: Verify component visibility (commands, agents)
- [ ] **Task 4.4**: Test workflow commands functionality
- [ ] **Task 4.5**: Test spec commands functionality
- [ ] **Task 4.6**: Test output style activation and @ references
- [ ] **Task 4.7**: Test hooks functionality

#### Test Commands
```bash
# Add local marketplace
/plugin marketplace add /Users/irudi/Code/personal/the-startup

# Install plugins
/plugin install the-startup-commands@the-startup-marketplace
/plugin install the-startup-agents@the-startup-marketplace
/plugin install the-startup-output-style@the-startup-marketplace
/plugin install the-startup-statusline@the-startup-marketplace

# Verify components
/help s:
/agents
```

#### Test Cases Checklist

**Component Visibility**:
- [ ] All 6 commands appear in `/help s:`
- [ ] All 50 agents are available in `/agents`

**Workflow Commands** (the-startup-commands):
- [ ] `/s:analyze` command works and delegates to agents
- [ ] `/s:specify` command works and creates spec documents
- [ ] `/s:implement` command works with phase execution
- [ ] `/s:refactor` command works

**Spec Commands** (the-startup-commands):
- [ ] `/s:init` command initializes templates
- [ ] `/s:spec` command creates numbered directories

**File References** (the-startup-commands):
- [ ] `@rules/` references resolve correctly in commands (4 commands reference agent-delegation.md and cycle-pattern.md)
- [ ] `${CLAUDE_PLUGIN_ROOT}/scripts/spec.py` executes correctly

**Output Style** (the-startup-output-style):
- [ ] Output style can be activated with `/settings add "outputStyle": "the-startup"`
- [ ] Output style `@rules/agent-delegation.md` reference resolves correctly (no errors when activated)

**Hooks** (the-startup-statusline):
- [ ] SessionStart hook shows welcome banner (first run)
- [ ] UserPromptSubmit hook shows statusline (if enabled)

**✅ Phase 4 Complete When**: All test cases pass, no errors in logs

---

### Phase 5: Documentation Updates

**Status**: ⬜ Not Started
**Duration**: 2 hours
**Completion**: 0/5 tasks

#### Tasks

- [ ] **Task 5.1**: Update root README.md with marketplace installation instructions
- [ ] **Task 5.2**: Create plugin-specific READMEs for each plugin
- [ ] **Task 5.3**: Update CHANGELOG.md with v2.0.0 entry
- [ ] **Task 5.4**: Document installation patterns (full, selective, migration)
- [ ] **Task 5.5**: Add troubleshooting section

**README.md Structure**:
```markdown
# The Startup - Marketplace

## Installation

### Install All Plugins (Recommended)
```bash
/plugin marketplace add rsmdt/the-startup
/plugin install the-startup-commands@the-startup-marketplace
/plugin install the-startup-agents@the-startup-marketplace
/plugin install the-startup-output-style@the-startup-marketplace
/plugin install the-startup-statusline@the-startup-marketplace
```

### Install Selectively
Choose only the plugins you need:
- **commands**: Core workflow orchestration
- **agents**: 50 specialist agents
- **output-style**: Custom output formatting
- **statusline**: Optional productivity hooks
```

#### Validation Checklist

- [ ] Installation instructions are clear and accurate
- [ ] Each of 4 plugins is documented with purpose and contents
- [ ] Usage examples provided for all commands
- [ ] Migration guide from monolithic plugin is complete
- [ ] Troubleshooting section addresses common issues
- [ ] README renders correctly on GitHub

**✅ Phase 5 Complete When**: All tasks checked, all validations passed, documentation is complete

---

### Phase 6: Git and Publishing

**Status**: ⬜ Not Started
**Duration**: 1 hour
**Completion**: 0/5 tasks

#### Tasks

- [ ] **Task 6.1**: Commit marketplace changes to `plugin` branch
- [ ] **Task 6.2**: Review changes and create pull request to main branch
- [ ] **Task 6.3**: Tag release as `v2.0.0`
- [ ] **Task 6.4**: Push to GitHub
- [ ] **Task 6.5**: Test installation from GitHub repository

**Git Workflow**:
```bash
# Review changes
git status
git diff

# Commit marketplace structure
git add .
git commit -m "feat: Convert to marketplace with 4 modular plugins

- the-startup-commands: Workflow orchestration
- the-startup-agents: 50 specialist agents
- the-startup-output-style: Custom output style
- the-startup-statusline: Productivity hooks

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>"

# Push to plugin branch
git push origin plugin

# After merge to main, tag release
git tag v2.0.0
git push origin v2.0.0
```

#### Validation Checklist

- [ ] Changes committed to plugin branch with descriptive commit message
- [ ] Pull request created and reviewed (no breaking changes)
- [ ] GitHub installation works: `/plugin marketplace add rsmdt/the-startup`
- [ ] All 4 plugins install successfully from GitHub
- [ ] Functionality verified post-installation (run Phase 4 tests again)
- [ ] Tag `v2.0.0` created and pushed

**✅ Phase 6 Complete When**: All tasks checked, all validations passed, marketplace live on GitHub

---

## User Installation Patterns

### Pattern 1: Full Installation (Recommended)

Install all plugins for complete functionality:

```bash
/plugin marketplace add rsmdt/the-startup
/plugin install the-startup-commands@the-startup-marketplace
/plugin install the-startup-agents@the-startup-marketplace
/plugin install the-startup-output-style@the-startup-marketplace
/plugin install the-startup-statusline@the-startup-marketplace
```

**Use Case**: Users who want the complete Agentic Startup experience

---

### Pattern 2: Core Workflow Only

Install commands and agents for workflow functionality:

```bash
/plugin marketplace add rsmdt/the-startup
/plugin install the-startup-commands@the-startup-marketplace
/plugin install the-startup-agents@the-startup-marketplace
```

**Use Case**: Users who want workflow automation without UX enhancements

---

### Pattern 3: Agents Only

Install just the specialized agents:

```bash
/plugin marketplace add rsmdt/the-startup
/plugin install the-startup-agents@the-startup-marketplace
```

**Use Case**: Users who want specialist agents for manual delegation, without workflow commands

---

### Pattern 4: Minimal UX

Install just productivity enhancements:

```bash
/plugin marketplace add rsmdt/the-startup
/plugin install the-startup-output-style@the-startup-marketplace
/plugin install the-startup-statusline@the-startup-marketplace
```

**Use Case**: Users who want UX improvements but will use their own agents/workflows

---

## Migration from Monolithic Plugin

For users who installed the monolithic `the-startup` plugin:

**Step 1: Remove old plugin**
```bash
/plugin uninstall the-startup
```

**Step 2: Install marketplace**
```bash
/plugin marketplace add rsmdt/the-startup
```

**Step 3: Install desired plugins**
```bash
/plugin install the-startup-commands@the-startup-marketplace
/plugin install the-startup-agents@the-startup-marketplace
/plugin install the-startup-output-style@the-startup-marketplace
/plugin install the-startup-statusline@the-startup-marketplace
```

**No Breaking Changes**: All functionality is preserved, just organized differently.

---

## 🎯 Success Criteria

**Status**: ⬜ Not Met

### Functional Requirements (6/6 Required)

- [ ] All 6 commands work correctly (analyze, specify, implement, refactor, init, spec)
- [ ] All 50 agents are available and functional
- [ ] Output style can be activated and renders correctly
- [ ] Hooks trigger correctly (welcome banner, statusline)
- [ ] `@rules/` references resolve in both commands and output-style
- [ ] `${CLAUDE_PLUGIN_ROOT}` paths work for scripts

### Non-Functional Requirements (6/6 Required)

- [ ] Installation is straightforward and documented
- [ ] Documentation is clear and comprehensive
- [ ] Plugins can be installed independently
- [ ] Plugins work together when all installed
- [ ] No performance degradation from monolithic version
- [ ] No functionality lost from monolithic version

### Quality Gates (6/6 Required)

- [ ] All plugin manifests validate (4 plugins)
- [ ] Marketplace manifest validates
- [ ] Local testing passes all 13 test cases
- [ ] GitHub installation works from remote repository
- [ ] README documentation is complete and accurate
- [ ] Migration guide is clear and tested

**✅ Project Complete When**: All success criteria met (18/18 checkboxes)

---

## Rollback Plan

If marketplace conversion fails:

**Step 1: Keep plugin branch**
- Do not merge to main until testing complete
- plugin branch contains marketplace structure
- main branch retains monolithic plugin

**Step 2: Tag fallback version**
```bash
git tag v1.9.9 main
```

**Step 3: Document issues**
- Create GitHub issues for any blockers
- Document what failed and why
- Plan fixes before retry

**Step 4: Communicate with users**
- If already published, notify users to use v1.9.9
- Provide timeline for marketplace completion

---

## ⏱️ Timeline Summary

| Phase | Duration | Status | Deliverable |
|-------|----------|--------|-------------|
| 1. Repository Restructuring | 2 hours | ⬜ Not Started | Plugin directories created, files moved, rules duplicated |
| 2. Fix File References & Manifests | 1.5 hours | ⬜ Not Started | Placeholders fixed, all plugin.json files created |
| 3. Marketplace Manifest | 30 min | ⬜ Not Started | marketplace.json configured |
| 4. Local Testing | 2 hours | ⬜ Not Started | All functionality validated |
| 5. Documentation Updates | 2 hours | ⬜ Not Started | README and docs updated |
| 6. Git and Publishing | 1 hour | ⬜ Not Started | Committed, tagged, published |

**Total Estimated Time**: 9 hours
**Actual Time Spent**: 0 hours (Update as you work)

**Note**: Update the Status column as phases progress (⬜ → 🟡 → ✅)

---

## Notes

**Why This Split?**
- Component-type organization is intuitive
- Users can install only what they need
- Each plugin has clear, single responsibility
- Easier to maintain and version independently

**Technical Constraints Respected**:
- Commands + rules + scripts + templates together (@ references and ${CLAUDE_PLUGIN_ROOT})
- Output style + rules/agent-delegation.md together (@ reference requires same plugin)
- Agents standalone (no dependencies)
- Hooks standalone (no dependencies)
- Rules duplication: `agent-delegation.md` exists in both `commands` and `output-style` plugins (required for @ references to work)

**Future Enhancements**:
- Version plugins independently
- Add more agents without affecting commands
- Update hooks without full reinstall
- Community contributions to individual plugins

**Maintenance Notes**:
- When updating `agent-delegation.md`, remember to update BOTH copies:
  - `plugins/the-startup-commands/rules/agent-delegation.md`
  - `plugins/the-startup-output-style/rules/agent-delegation.md`
- Consider adding a pre-commit hook to verify both files are identical
- Other rules files (`cycle-pattern.md`, `agent-creation-principles.md`) only exist in commands plugin

---

## 📋 Quick Reference

### Key Checkpoints

1. **Phase 1 Complete**: Files organized in `plugins/` directories, rules duplicated
2. **Phase 2 Complete**: All plugin.json files created, placeholders fixed
3. **Phase 3 Complete**: marketplace.json configured
4. **Phase 4 Complete**: All 13 tests pass locally
5. **Phase 5 Complete**: Documentation updated and complete
6. **Phase 6 Complete**: Published to GitHub, installation works remotely

### Critical Files to Track

- `plugins/the-startup-commands/rules/agent-delegation.md` (225 lines - duplicated)
- `plugins/the-startup-output-style/rules/agent-delegation.md` (225 lines - duplicated)
- `plugins/the-startup-output-style/output-styles/the-startup.md` (line 66 - placeholder fix)
- `.claude-plugin/marketplace.json` (4 plugins listed)

### Task Counts

- Phase 1: 5 tasks
- Phase 2: 6 tasks
- Phase 3: 3 tasks
- Phase 4: 13 test cases
- Phase 5: 5 tasks
- Phase 6: 5 tasks
- **Total**: 37 checkable items

### Progress Tracking Tips

1. Update "Last Updated" date at the top when you make changes
2. Keep "Overall Progress" table synchronized with phase statuses
3. Update "Quick Resume" section after completing each phase
4. Mark Success Criteria checkboxes only after Phase 4 and Phase 6
5. Update "Actual Time Spent" in Timeline Summary as you work

---

## 🏁 Final Checklist

Before closing this spec, verify:

- [ ] All 6 phases show ✅ Completed status
- [ ] Overall Progress shows 100% (6/6 phases completed)
- [ ] All 37 tasks/tests are checked
- [ ] Success Criteria shows 18/18 checkboxes complete
- [ ] Timeline Summary "Actual Time Spent" is updated
- [ ] Status at top is updated to ✅ Completed
- [ ] GitHub marketplace installation works
- [ ] All 4 plugins functional when installed together

**When all items above are checked, the marketplace conversion is complete!** 🎉
