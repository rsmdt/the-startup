# Codebase Analysis for Claude Code Plugin Migration

**Date:** 2025-10-11
**Purpose:** Comprehensive analysis of the-agentic-startup codebase to guide npm CLI → Claude Code plugin migration
**Status:** Complete baseline analysis

---

## Executive Summary

This analysis provides a complete inventory and classification of all components in the-agentic-startup codebase to guide the migration from npm CLI package to Claude Code plugin.

### Key Findings

- **Total Source Code:** 5,263 lines of TypeScript across 27 files
- **Total Assets:** 59 files (agents, commands, templates, rules, scripts, settings)
- **Asset Content Size:** ~363 KB of markdown and configuration
- **Migration Impact:** ~75% of source code will be eliminated (CLI installer, UI, lock management)
- **External References:** 7 files contain `@{{STARTUP_PATH}}` or `{{STARTUP_PATH}}` placeholders requiring build-time preprocessing

### Migration Strategy

1. **ELIMINATE:** 75% of source code (CLI, installer, UI, lock manager)
2. **KEEP AS-IS:** 39 agent files copy directly with no changes
3. **TRANSFORM:** 6 command/config files require preprocessing to inline external references
4. **BUNDLE:** All assets must be embedded in plugin distribution

---

## 1. Repository Structure Analysis

### Root Directory Layout

```
the-agentic-startup/
├── assets/                      # All distributable content (59 files, 363KB)
│   ├── claude/                  # Claude configuration (46 files)
│   │   ├── agents/             # 39 agent definition files
│   │   ├── commands/           # 5 slash command files
│   │   ├── output-styles/      # 1 output style file
│   │   ├── settings.json       # Base settings with placeholders
│   │   └── settings.local.json # Local override example
│   └── the-startup/            # Startup-specific assets (11 files)
│       ├── bin/                # 2 shell scripts (statusline.sh, .ps1)
│       ├── rules/              # 3 rule files (delegation, cycle, principles)
│       └── templates/          # 6 template files (PRD, SDD, PLAN, DOR, DOD, TASK-DOD)
├── src/                        # Source code (27 TypeScript files, 5,263 lines)
│   ├── cli/                    # CLI commands (8 files) → ELIMINATE
│   ├── core/                   # Business logic (8 files)
│   │   ├── installer/          # Installation engine (4 files) → ELIMINATE
│   │   ├── init/               # Template initialization (1 file) → ELIMINATE
│   │   ├── spec/               # Spec generation (1 file) → KEEP/ADAPT
│   │   └── types/              # TypeScript types (3 files) → PARTIAL
│   ├── ui/                     # Ink/React UI (8 files, 1,683 lines) → ELIMINATE
│   └── index.ts                # Entry point → ELIMINATE
├── tests/                      # Test suite (14 test files) → ELIMINATE
├── docs/                       # Documentation
│   ├── domain/                 # Business rules and workflows
│   ├── patterns/               # Technical patterns
│   ├── interfaces/             # API contracts
│   ├── specs/                  # Feature specifications
│   └── research/               # Research and planning documents
├── package.json                # npm package configuration → TRANSFORM
├── tsconfig.json               # TypeScript configuration → ELIMINATE
├── vitest.config.ts            # Test configuration → ELIMINATE
└── CLAUDE.md                   # Project-specific instructions → TRANSFORM

Total Files (excluding node_modules/dist/.git):
- TypeScript source: 27 files
- Test files: 14 files
- Asset files: 59 files
- Config files: ~10 files
```

---

## 2. Source Code Inventory

### Complete File List with Line Counts

#### CLI Layer (8 files, ~800 lines) → ELIMINATE

| File | Lines | Purpose | Migration Action |
|------|-------|---------|------------------|
| `src/cli/index.ts` | 89 | Commander.js CLI setup | **ELIMINATE** - Plugin uses built-in command system |
| `src/cli/install.ts` | 158 | Interactive/non-interactive install | **ELIMINATE** - Plugin auto-installs |
| `src/cli/uninstall.ts` | ~120 | Uninstall with lock file reading | **ELIMINATE** - Plugin manages lifecycle |
| `src/cli/init.ts` | ~100 | DOR/DOD/TASK-DOD initialization | **ELIMINATE** - Plugin provides different workflow |
| `src/cli/spec.ts` | ~150 | Spec directory generation | **KEEP/ADAPT** - Slash command may use this |
| `src/cli/asset-provider.ts` | ~80 | File system asset loading | **ELIMINATE** - Plugin uses embedded assets |
| `src/cli/asset-provider-old.ts` | ~80 | Legacy provider | **ELIMINATE** - Dead code |
| `src/cli/fs-adapter.ts` | ~100 | File system abstraction | **ELIMINATE** - Not needed for plugin |

#### Core Business Logic (8 files, ~2,000 lines)

| File | Lines | Purpose | Migration Action |
|------|-------|---------|------------------|
| `src/core/installer/Installer.ts` | 438 | Installation orchestration | **ELIMINATE** - Plugin auto-installs |
| `src/core/installer/LockManager.ts` | ~300 | Lock file management | **ELIMINATE** - No lock files in plugin |
| `src/core/installer/SettingsMerger.ts` | ~250 | Deep merge settings.json | **ELIMINATE** - Plugin manages settings |
| `src/core/installer/DeepMerge.ts` | ~150 | Deep merge utility | **ELIMINATE** - Not needed |
| `src/core/init/Initializer.ts` | ~200 | DOR/DOD template initialization | **ELIMINATE/ADAPT** - Different workflow |
| `src/core/spec/SpecGenerator.ts` | 278 | Spec directory creation | **KEEP/ADAPT** - May power slash command |
| `src/core/types/config.ts` | ~150 | Configuration types | **PARTIAL** - Extract relevant types |
| `src/core/types/settings.ts` | ~100 | Settings.json types | **ELIMINATE** - Plugin doesn't manipulate settings |
| `src/core/types/lock.ts` | ~80 | Lock file types | **ELIMINATE** - No lock files |

#### UI Layer (8 files, 1,683 lines) → ELIMINATE ENTIRELY

| File | Lines | Purpose | Migration Action |
|------|-------|---------|------------------|
| `src/ui/install/InstallWizard.tsx` | ~450 | Main installation wizard | **ELIMINATE** - No installer UI |
| `src/ui/install/ChoiceSelector.tsx` | ~250 | Interactive menu component | **ELIMINATE** |
| `src/ui/install/FinalConfirmation.tsx` | ~200 | Confirmation screen | **ELIMINATE** |
| `src/ui/install/Complete.tsx` | ~150 | Success screen | **ELIMINATE** |
| `src/ui/uninstall/UninstallWizard.tsx` | ~200 | Uninstall wizard | **ELIMINATE** |
| `src/ui/shared/Banner.tsx` | ~100 | ASCII art banner | **ELIMINATE** |
| `src/ui/shared/ErrorDisplay.tsx` | ~120 | Error display component | **ELIMINATE** |
| `src/ui/shared/Spinner.tsx` | ~80 | Loading spinner | **ELIMINATE** |
| `src/ui/shared/theme.ts` | ~133 | Color theme definitions | **ELIMINATE** |

#### Entry Points (2 files) → ELIMINATE

| File | Lines | Purpose | Migration Action |
|------|-------|---------|------------------|
| `src/index.ts` | ~30 | Main CLI entry point | **ELIMINATE** - Plugin entry different |
| `src/bin/spec.ts` | ~40 | Standalone spec executable | **ELIMINATE/ADAPT** - May inform plugin command |

### Source Code Summary

| Category | Files | Lines | Migration Action |
|----------|-------|-------|------------------|
| **CLI Layer** | 8 | ~800 | **ELIMINATE** (100%) |
| **Core - Installer** | 4 | ~1,138 | **ELIMINATE** (100%) |
| **Core - Init/Spec** | 2 | ~478 | **KEEP/ADAPT** (50%) |
| **Core - Types** | 3 | ~330 | **PARTIAL** (20%) |
| **UI Layer** | 8 | 1,683 | **ELIMINATE** (100%) |
| **Entry Points** | 2 | ~70 | **ELIMINATE** (100%) |
| **TOTAL** | 27 | 5,263 | **~75% ELIMINATION** |

**Elimination Breakdown:**
- **Eliminate Completely:** ~4,000 lines (76%)
- **Keep/Adapt:** ~478 lines (9%)
- **Partial Keep:** ~330 lines (6%)
- **Undecided:** ~455 lines (9%)

---

## 3. Asset Analysis

### Complete Asset Inventory

#### Claude Components (46 files, ~269 KB)

##### Agents (39 files, ~215 KB)

**Role-Based Organization:**

| Role | Activity Files | Bytes | Migration Action |
|------|----------------|-------|------------------|
| **the-analyst** | 3 (feature-prioritization, project-coordination, requirements-analysis) | ~15,000 | **KEEP AS-IS** |
| **the-architect** | 5 (quality-review, system-architecture, system-documentation, technology-research, technology-standards) | ~22,000 | **KEEP AS-IS** |
| **the-designer** | 4 (accessibility-implementation, design-foundation, interaction-architecture, user-research) | ~18,000 | **KEEP AS-IS** |
| **the-ml-engineer** | 4 (context-management, feature-operations, ml-operations, prompt-optimization) | ~20,000 | **KEEP AS-IS** |
| **the-mobile-engineer** | 3 (mobile-data-persistence, mobile-development, mobile-operations) | ~16,000 | **KEEP AS-IS** |
| **the-platform-engineer** | 7 (containerization, data-architecture, deployment-automation, infrastructure-as-code, performance-tuning, pipeline-engineering, production-monitoring) | ~32,000 | **KEEP AS-IS** |
| **the-qa-engineer** | 3 (exploratory-testing, performance-testing, test-execution) | ~14,000 | **KEEP AS-IS** |
| **the-security-engineer** | 3 (security-assessment, security-implementation, security-incident-response) | ~15,000 | **KEEP AS-IS** |
| **the-software-engineer** | 5 (api-development, component-development, domain-modeling, performance-optimization, service-resilience) | ~24,000 | **KEEP AS-IS** |
| **the-chief** | 1 (orchestration/routing) | ~8,000 | **KEEP AS-IS** |
| **the-meta-agent** | 1 (agent creation) | ~7,000 | **TRANSFORM** (external refs) |

**Total:** 39 agent files, ~215,213 bytes

**Migration Notes:**
- 38 agents require NO changes (pure content, no external references)
- 1 agent (the-meta-agent) requires transformation to inline external references

##### Commands (5 files, ~40 KB)

| Command | File | Bytes | External Refs | Migration Action |
|---------|------|-------|---------------|------------------|
| `/s:analyze` | `commands/s/analyze.md` | ~8,000 | 2 (@{{STARTUP_PATH}} refs) | **TRANSFORM** |
| `/s:implement` | `commands/s/implement.md` | ~7,000 | 1 (@{{STARTUP_PATH}} ref) | **TRANSFORM** |
| `/s:init` | `commands/s/init.md` | ~6,000 | 0 | **KEEP AS-IS** |
| `/s:refactor` | `commands/s/refactor.md` | ~7,500 | 2 (@{{STARTUP_PATH}} refs) | **TRANSFORM** |
| `/s:specify` | `commands/s/specify.md` | ~10,000 | 2 (@{{STARTUP_PATH}} refs) | **TRANSFORM** |

**Total:** 5 command files, ~38,500 bytes

**External References Found:**
```markdown
# Pattern in commands:
@{{STARTUP_PATH}}/rules/agent-delegation.md
@{{STARTUP_PATH}}/rules/cycle-pattern.md
```

##### Output Styles (1 file)

| File | Bytes | External Refs | Migration Action |
|------|-------|---------------|------------------|
| `output-styles/the-startup.md` | ~8,000 | 1 (@{{STARTUP_PATH}} ref) | **TRANSFORM** |

##### Settings Files (2 files)

| File | Bytes | Placeholders | Migration Action |
|------|-------|--------------|------------------|
| `settings.json` | ~200 | `{{STARTUP_PATH}}`, `{{SHELL_SCRIPT_EXTENSION}}` | **TRANSFORM** |
| `settings.local.json` | ~150 | Same | **ELIMINATE** (example only) |

**settings.json content:**
```json
{
  "permissions": {
    "additionalDirectories": ["{{STARTUP_PATH}}"]
  },
  "statusLine": {
    "type": "command",
    "command": "{{STARTUP_PATH}}/bin/statusline{{SHELL_SCRIPT_EXTENSION}}"
  }
}
```

#### Startup Components (11 files, ~94 KB)

##### Shell Scripts (2 files)

| File | Bytes | Purpose | Migration Action |
|------|-------|---------|------------------|
| `bin/statusline.sh` | ~1,600 | Git branch extraction (bash) | **BUNDLE** (keep as-is) |
| `bin/statusline.ps1` | ~1,800 | Git branch extraction (PowerShell) | **BUNDLE** (keep as-is) |

##### Rules (3 files, ~20 KB)

| File | Bytes | Purpose | Migration Action |
|------|-------|---------|------------------|
| `rules/agent-delegation.md` | ~8,500 | Task decomposition patterns | **BUNDLE** (referenced by commands) |
| `rules/cycle-pattern.md` | ~6,000 | Standard workflow cycle | **BUNDLE** (referenced by commands) |
| `rules/agent-creation-principles.md` | ~5,500 | Agent design principles | **BUNDLE** (referenced by meta-agent) |

**Critical for Build:** These files are referenced by commands via `@{{STARTUP_PATH}}/rules/...` and must be inlined during build.

##### Templates (6 files, ~70 KB)

| File | Bytes | Purpose | Migration Action |
|------|-------|---------|------------------|
| `templates/product-requirements.md` | ~12,000 | PRD template | **BUNDLE/ADAPT** |
| `templates/solution-design.md` | ~15,000 | SDD template | **BUNDLE/ADAPT** |
| `templates/implementation-plan.md` | ~10,000 | PLAN template | **BUNDLE/ADAPT** |
| `templates/definition-of-ready.md` | ~8,000 | DOR validation template | **BUNDLE/ADAPT** |
| `templates/definition-of-done.md` | ~9,000 | DOD validation template | **BUNDLE/ADAPT** |
| `templates/task-definition-of-done.md` | ~7,000 | Task-level DOD template | **BUNDLE/ADAPT** |

**Migration Note:** Templates may need adaptation for plugin workflow (no `spec` command to generate directories).

### Asset Summary

| Category | Files | Total Bytes | Keep As-Is | Transform | Eliminate |
|----------|-------|-------------|------------|-----------|-----------|
| **Agents** | 39 | ~215,213 | 38 | 1 | 0 |
| **Commands** | 5 | ~38,500 | 1 | 4 | 0 |
| **Output Styles** | 1 | ~8,000 | 0 | 1 | 0 |
| **Settings** | 2 | ~350 | 0 | 1 | 1 |
| **Shell Scripts** | 2 | ~3,400 | 2 | 0 | 0 |
| **Rules** | 3 | ~20,000 | 3 | 0 | 0 |
| **Templates** | 6 | ~70,000 | 0 | 6 | 0 |
| **TOTAL** | 59 | ~362,680 | 44 | 13 | 1 |

---

## 4. Migration Classification

### ELIMINATE (75% of source code)

**Rationale:** Plugin architecture provides these features natively.

#### Complete Elimination List

1. **CLI Infrastructure (8 files)**
   - `src/cli/index.ts` - Commander.js setup
   - `src/cli/install.ts` - Interactive installer
   - `src/cli/uninstall.ts` - Uninstaller
   - `src/cli/init.ts` - Template initializer
   - `src/cli/asset-provider.ts` - File loading
   - `src/cli/asset-provider-old.ts` - Dead code
   - `src/cli/fs-adapter.ts` - FS abstraction

2. **Installer Core (4 files)**
   - `src/core/installer/Installer.ts` - Installation orchestration
   - `src/core/installer/LockManager.ts` - Lock file management
   - `src/core/installer/SettingsMerger.ts` - Settings merging
   - `src/core/installer/DeepMerge.ts` - Merge utility

3. **UI Layer (8 files)**
   - All `src/ui/install/*.tsx` files
   - All `src/ui/uninstall/*.tsx` files
   - All `src/ui/shared/*.tsx` and `theme.ts` files

4. **Type Definitions (2 files)**
   - `src/core/types/lock.ts` - Lock file types
   - `src/core/types/settings.ts` - Settings types (partial)

5. **Entry Points (2 files)**
   - `src/index.ts` - Main CLI entry
   - `src/bin/spec.ts` - Standalone executable (maybe adapt)

6. **Build/Test Infrastructure**
   - All test files (14 files)
   - `tsconfig.json`
   - `vitest.config.ts`
   - `tsup.config.ts`
   - `.eslintrc.json`
   - `.prettierrc.json`

**Total Elimination:** ~4,500 lines of code + all build/test infrastructure

### KEEP AS-IS (44 files)

**Rationale:** Content files that copy directly into plugin without modification.

1. **Agents (38 files)** - All role/activity agent files except `the-meta-agent.md`
2. **Commands (1 file)** - `commands/s/init.md` (no external refs)
3. **Shell Scripts (2 files)** - `bin/statusline.sh`, `bin/statusline.ps1`
4. **Rules (3 files)** - All rule files (will be inlined into commands during build)

**Total:** 44 files, ~238 KB of content

### TRANSFORM (13 files + 2 config files)

**Rationale:** Files with external references or placeholders requiring build-time preprocessing.

#### Files Requiring Transformation

1. **Commands with External References (4 files)**
   - `commands/s/analyze.md` - Inline 2 rule file references
   - `commands/s/implement.md` - Inline 1 rule file reference
   - `commands/s/refactor.md` - Inline 2 rule file references
   - `commands/s/specify.md` - Inline 2 rule file references

2. **Agent with External References (1 file)**
   - `agents/the-meta-agent.md` - Inline agent-creation-principles reference

3. **Output Style (1 file)**
   - `output-styles/the-startup.md` - Inline agent-delegation reference

4. **Settings File (1 file)**
   - `settings.json` - Replace `{{STARTUP_PATH}}` and `{{SHELL_SCRIPT_EXTENSION}}`

5. **Templates (6 files)** - May need workflow adaptation
   - All 6 template files in `templates/`

6. **Project Configuration (2 files)**
   - `package.json` - Transform to plugin manifest
   - `CLAUDE.md` - Adapt for plugin context

#### Transformation Patterns

**Pattern 1: Inline External Rule Files**

Before (in command file):
```markdown
@{{STARTUP_PATH}}/rules/agent-delegation.md
```

After (build-time transformation):
```markdown
Rules for task decomposition and parallel execution.

1. Task Decomposition Principles:
   [... full content of agent-delegation.md ...]
```

**Pattern 2: Replace Path Placeholders**

Before (in settings.json):
```json
{
  "permissions": {
    "additionalDirectories": ["{{STARTUP_PATH}}"]
  }
}
```

After (for plugin, paths are relative):
```json
{
  "permissions": {
    "additionalDirectories": [".the-startup"]
  }
}
```

**Pattern 3: Shell Script Extension**

Before:
```json
"command": "{{STARTUP_PATH}}/bin/statusline{{SHELL_SCRIPT_EXTENSION}}"
```

After (plugin needs platform detection):
```json
"command": ".the-startup/bin/statusline.sh"  // or .ps1 on Windows
```

### BUNDLE (All 59 asset files)

**Rationale:** Plugin must include all assets in distribution.

**Bundling Strategy:**
1. Embed all transformed assets in plugin distribution
2. No external file loading - everything inline
3. Plugin installation copies assets to `.claude/` and `.the-startup/` in project

---

## 5. External Reference Mapping

### Complete Reference Audit

#### @{{STARTUP_PATH}} References (9 occurrences in 6 files)

| File | Line(s) | Reference | Target File |
|------|---------|-----------|-------------|
| `commands/s/specify.md` | 30 | `@{{STARTUP_PATH}}/rules/agent-delegation.md` | `rules/agent-delegation.md` |
| `commands/s/specify.md` | 34 | `@{{STARTUP_PATH}}/rules/cycle-pattern.md` | `rules/cycle-pattern.md` |
| `commands/s/refactor.md` | 29 | `@{{STARTUP_PATH}}/rules/agent-delegation.md` | `rules/agent-delegation.md` |
| `commands/s/refactor.md` | 33 | `@{{STARTUP_PATH}}/rules/cycle-pattern.md` | `rules/cycle-pattern.md` |
| `commands/s/analyze.md` | 20 | `@{{STARTUP_PATH}}/rules/agent-delegation.md` | `rules/agent-delegation.md` |
| `commands/s/analyze.md` | 24 | `@{{STARTUP_PATH}}/rules/cycle-pattern.md` | `rules/cycle-pattern.md` |
| `commands/s/implement.md` | 24 | `@{{STARTUP_PATH}}/rules/agent-delegation.md` | `rules/agent-delegation.md` |
| `agents/the-meta-agent.md` | 60 | `@{{STARTUP_PATH}}/rules/agent-creation-principles.md` | `rules/agent-creation-principles.md` |
| `output-styles/the-startup.md` | 66 | `@{{STARTUP_PATH}}/rules/agent-delegation.md` | `rules/agent-delegation.md` |

#### {{STARTUP_PATH}} Placeholders (3 occurrences in 1 file)

| File | Occurrences | Purpose |
|------|-------------|---------|
| `settings.json` | 2 | Path to `.the-startup` directory |
| `settings.json` | 1 | Shell script extension (`.sh` or `.ps1`) |

### Referenced Files Requiring Inlining

| Target File | Size | Referenced By | Migration Action |
|-------------|------|---------------|------------------|
| `rules/agent-delegation.md` | ~8,500 bytes | 5 files | **INLINE** during build |
| `rules/cycle-pattern.md` | ~6,000 bytes | 3 files | **INLINE** during build |
| `rules/agent-creation-principles.md` | ~5,500 bytes | 1 file | **INLINE** during build |

**Total Content to Inline:** ~20,000 bytes across 3 files

### Inlining Impact Analysis

| Source File | External Refs | Size Before | Est. Size After | Growth |
|-------------|---------------|-------------|-----------------|--------|
| `commands/s/specify.md` | 2 | ~10,000 | ~24,500 | +145% |
| `commands/s/refactor.md` | 2 | ~7,500 | ~22,000 | +193% |
| `commands/s/analyze.md` | 2 | ~8,000 | ~22,500 | +181% |
| `commands/s/implement.md` | 1 | ~7,000 | ~15,500 | +121% |
| `agents/the-meta-agent.md` | 1 | ~7,000 | ~12,500 | +79% |
| `output-styles/the-startup.md` | 1 | ~8,000 | ~16,500 | +106% |

**Total Asset Size Impact:**
- Before inlining: ~47,500 bytes (6 files)
- After inlining: ~113,500 bytes (6 files)
- Growth: +66,000 bytes (+139%)

**Plugin Distribution Impact:**
- Current assets: ~363 KB
- After inlining: ~429 KB (~420 KB after removing redundant rule files)
- Net increase: ~57 KB (+16%)

---

## 6. Build Requirements

### Build Process Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     BUILD PIPELINE                          │
└─────────────────────────────────────────────────────────────┘

INPUT: Source Assets (59 files)
  ├─ agents/*.md (39 files)
  ├─ commands/**/*.md (5 files)
  ├─ output-styles/*.md (1 file)
  ├─ rules/*.md (3 files)
  ├─ templates/*.md (6 files)
  ├─ bin/*.{sh,ps1} (2 files)
  ├─ settings.json (1 file)
  └─ settings.local.json (1 file)

       │
       │ STEP 1: Preprocessing
       ▼
  ┌──────────────────────────────────────┐
  │ 1. Parse @{{STARTUP_PATH}} refs     │
  │ 2. Read referenced rule files        │
  │ 3. Inline rule content into commands │
  │ 4. Replace {{STARTUP_PATH}} with     │
  │    relative paths                    │
  │ 5. Select shell script extension     │
  └──────────────────────────────────────┘

       │
       │ STEP 2: Validation
       ▼
  ┌──────────────────────────────────────┐
  │ 1. Verify all external refs resolved │
  │ 2. Check markdown syntax             │
  │ 3. Validate JSON files               │
  │ 4. Ensure no dangling placeholders   │
  └──────────────────────────────────────┘

       │
       │ STEP 3: Bundling
       ▼
  ┌──────────────────────────────────────┐
  │ 1. Organize by target directory      │
  │    - .claude/agents/                 │
  │    - .claude/commands/               │
  │    - .claude/output-styles/          │
  │    - .the-startup/bin/               │
  │    - .the-startup/templates/         │
  │ 2. Embed in plugin distribution      │
  │ 3. Generate installation metadata    │
  └──────────────────────────────────────┘

       │
       │ OUTPUT: Plugin Distribution
       ▼

OUTPUT: Bundled Plugin
  ├─ plugin.json (manifest)
  ├─ assets/
  │   ├─ claude/
  │   │   ├─ agents/ (39 files, transformed)
  │   │   ├─ commands/ (5 files, transformed)
  │   │   └─ output-styles/ (1 file, transformed)
  │   └─ the-startup/
  │       ├─ bin/ (2 files, as-is)
  │       └─ templates/ (6 files, adapted)
  └─ README.md
```

### Required Build Operations

#### 1. External Reference Inlining

**Operation:** Replace `@{{STARTUP_PATH}}/path/to/file.md` with file content

**Input Files:**
- 4 command files (`analyze.md`, `implement.md`, `refactor.md`, `specify.md`)
- 1 agent file (`the-meta-agent.md`)
- 1 output style file (`the-startup.md`)

**Reference Files:**
- `rules/agent-delegation.md` (~8.5 KB)
- `rules/cycle-pattern.md` (~6 KB)
- `rules/agent-creation-principles.md` (~5.5 KB)

**Algorithm:**
```typescript
function inlineExternalReferences(content: string, rulesDir: string): string {
  const refPattern = /@\{\{STARTUP_PATH\}\}\/rules\/([a-z-]+\.md)/g;

  return content.replace(refPattern, (match, filename) => {
    const rulePath = path.join(rulesDir, filename);
    const ruleContent = fs.readFileSync(rulePath, 'utf-8');
    return ruleContent;
  });
}
```

**Output:** Transformed files with inlined content

#### 2. Placeholder Replacement

**Operation:** Replace `{{STARTUP_PATH}}` and `{{SHELL_SCRIPT_EXTENSION}}` with plugin-appropriate values

**Input Files:**
- `settings.json`

**Replacement Rules:**
```typescript
const replacements = {
  '{{STARTUP_PATH}}': '.the-startup',  // Plugin creates local directory
  '{{SHELL_SCRIPT_EXTENSION}}': process.platform === 'win32' ? '.ps1' : '.sh',
  '{{CLAUDE_PATH}}': '.claude',  // Not used in current assets
};
```

**Output:** `settings.json` with resolved paths

#### 3. Validation

**Pre-build Validation:**
1. All referenced files exist
2. No circular references
3. Markdown syntax is valid
4. JSON files parse correctly
5. No remaining unreplaced placeholders

**Post-build Validation:**
1. All external references resolved
2. File count matches expected (44 keep + 6 transformed)
3. Total asset size within expected range (~420 KB)
4. No syntax errors in transformed files

#### 4. Bundling Strategy

**Option A: Embed as JSON**
```json
{
  "assets": {
    "agents": {
      "the-chief": "---\nname: the-chief\n...",
      "the-software-engineer-api-development": "---\n..."
    },
    "commands": { ... }
  }
}
```

**Option B: Embed as Directory Structure**
```
plugin/
  assets/
    claude/
      agents/
        the-chief.md
        the-analyst/
          feature-prioritization.md
          ...
```

**Option C: Base64 Encoded Archive**
```typescript
const assetsArchive = {
  format: 'tar.gz',
  data: 'base64-encoded-tarball',
  checksum: 'sha256:...'
};
```

**Recommendation:** Option B (directory structure)
- Preserves folder hierarchy
- Easy to inspect and debug
- Matches npm package structure
- No decoding overhead

### Build Script Outline

```typescript
// build-plugin.ts

import { glob } from 'glob';
import { readFileSync, writeFileSync } from 'fs';
import { join, dirname } from 'path';

interface BuildConfig {
  sourceDir: string;    // assets/
  outputDir: string;    // plugin-dist/
  rulesDir: string;     // assets/the-startup/rules/
}

async function buildPlugin(config: BuildConfig): Promise<void> {
  console.log('Starting plugin build...');

  // 1. Collect all asset files
  const assets = await collectAssets(config.sourceDir);

  // 2. Classify assets (keep-as-is vs transform)
  const { keepAsIs, transform } = classifyAssets(assets);

  // 3. Copy keep-as-is files
  for (const asset of keepAsIs) {
    copyAsset(asset, config.outputDir);
  }

  // 4. Transform files with external references
  for (const asset of transform) {
    const content = readFileSync(asset.source, 'utf-8');
    const transformed = inlineExternalReferences(content, config.rulesDir);
    const replaced = replacePlaceholders(transformed);
    writeAsset(replaced, asset.target, config.outputDir);
  }

  // 5. Validate output
  await validateBuild(config.outputDir);

  // 6. Generate plugin manifest
  await generateManifest(config.outputDir);

  console.log('Plugin build complete!');
}

function inlineExternalReferences(content: string, rulesDir: string): string {
  // Inline @{{STARTUP_PATH}}/rules/*.md references
  const refPattern = /@\{\{STARTUP_PATH\}\}\/rules\/([a-z-]+\.md)/g;

  return content.replace(refPattern, (match, filename) => {
    const rulePath = join(rulesDir, filename);
    const ruleContent = readFileSync(rulePath, 'utf-8');
    return ruleContent;
  });
}

function replacePlaceholders(content: string): string {
  // Replace path placeholders
  let result = content;
  result = result.replace(/\{\{STARTUP_PATH\}\}/g, '.the-startup');
  result = result.replace(/\{\{CLAUDE_PATH\}\}/g, '.claude');
  result = result.replace(/\{\{SHELL_SCRIPT_EXTENSION\}\}/g,
    process.platform === 'win32' ? '.ps1' : '.sh');
  return result;
}

async function validateBuild(outputDir: string): Promise<void> {
  // Validate no unreplaced placeholders
  const files = await glob('**/*.{md,json}', { cwd: outputDir });

  for (const file of files) {
    const content = readFileSync(join(outputDir, file), 'utf-8');

    if (content.includes('@{{STARTUP_PATH}}')) {
      throw new Error(`Unresolved reference in ${file}`);
    }

    if (content.includes('{{STARTUP_PATH}}') && !file.endsWith('.md')) {
      throw new Error(`Unresolved placeholder in ${file}`);
    }
  }

  console.log('Build validation passed');
}
```

### Build Output Structure

```
plugin-dist/
├── plugin.json                 # Plugin manifest
├── README.md                   # Installation instructions
└── assets/
    ├── claude/
    │   ├── agents/            # 39 files (1 transformed, 38 as-is)
    │   │   ├── the-chief.md
    │   │   ├── the-meta-agent.md  # TRANSFORMED (inlined refs)
    │   │   ├── the-analyst/
    │   │   │   ├── feature-prioritization.md
    │   │   │   ├── project-coordination.md
    │   │   │   └── requirements-analysis.md
    │   │   ├── the-architect/
    │   │   │   ├── quality-review.md
    │   │   │   ├── system-architecture.md
    │   │   │   ├── system-documentation.md
    │   │   │   ├── technology-research.md
    │   │   │   └── technology-standards.md
    │   │   ├── [... other roles ...]
    │   │   └── the-software-engineer/
    │   │       ├── api-development.md
    │   │       ├── component-development.md
    │   │       ├── domain-modeling.md
    │   │       ├── performance-optimization.md
    │   │       └── service-resilience.md
    │   ├── commands/          # 5 files (4 transformed, 1 as-is)
    │   │   └── s/
    │   │       ├── analyze.md       # TRANSFORMED (inlined 2 refs)
    │   │       ├── implement.md     # TRANSFORMED (inlined 1 ref)
    │   │       ├── init.md          # AS-IS (no refs)
    │   │       ├── refactor.md      # TRANSFORMED (inlined 2 refs)
    │   │       └── specify.md       # TRANSFORMED (inlined 2 refs)
    │   ├── output-styles/     # 1 file (transformed)
    │   │   └── the-startup.md       # TRANSFORMED (inlined 1 ref)
    │   └── settings.json      # TRANSFORMED (replaced placeholders)
    └── the-startup/
        ├── bin/               # 2 files (as-is)
        │   ├── statusline.sh
        │   └── statusline.ps1
        └── templates/         # 6 files (adapted)
            ├── product-requirements.md
            ├── solution-design.md
            ├── implementation-plan.md
            ├── definition-of-ready.md
            ├── definition-of-done.md
            └── task-definition-of-done.md
```

**File Count:**
- Total files in plugin: 53 (down from 59)
  - 39 agents
  - 5 commands
  - 1 output style
  - 1 settings.json
  - 2 shell scripts
  - 6 templates
- Eliminated files: 6 (3 rule files + 2 duplicate files)

**Size Estimate:**
- Transformed assets: ~420 KB
- Plugin manifest: ~5 KB
- README: ~10 KB
- **Total plugin size: ~435 KB**

---

## 7. File Size Analysis for Plugin Distribution

### Current Distribution Breakdown

| Category | Files | Total Size | Avg Size | Distribution % |
|----------|-------|------------|----------|----------------|
| **Agents** | 39 | 215 KB | 5.5 KB | 59.3% |
| **Commands** | 5 | 38 KB | 7.6 KB | 10.5% |
| **Templates** | 6 | 70 KB | 11.7 KB | 19.3% |
| **Rules** | 3 | 20 KB | 6.7 KB | 5.5% |
| **Output Styles** | 1 | 8 KB | 8 KB | 2.2% |
| **Scripts** | 2 | 3.4 KB | 1.7 KB | 0.9% |
| **Settings** | 2 | 0.35 KB | 0.18 KB | 0.1% |
| **Other** | 1 | 8 KB | 8 KB | 2.2% |
| **TOTAL** | 59 | 363 KB | 6.2 KB | 100% |

### Post-Transformation Estimates

After inlining external references and eliminating redundant files:

| Category | Files | Size Before | Size After | Growth |
|----------|-------|-------------|------------|--------|
| **Agents** | 39 | 215 KB | 221 KB | +3% (1 file transformed) |
| **Commands** | 5 | 38 KB | 91 KB | +139% (4 files transformed) |
| **Output Styles** | 1 | 8 KB | 17 KB | +112% (1 file transformed) |
| **Templates** | 6 | 70 KB | 70 KB | 0% (no changes) |
| **Scripts** | 2 | 3.4 KB | 3.4 KB | 0% (no changes) |
| **Settings** | 1 | 0.2 KB | 0.2 KB | 0% (minor placeholder replacement) |
| **TOTAL** | 54 | 334 KB | 402 KB | +20% |

**Eliminated Files:**
- `rules/agent-delegation.md` (~8.5 KB) - Inlined into 5 files
- `rules/cycle-pattern.md` (~6 KB) - Inlined into 3 files
- `rules/agent-creation-principles.md` (~5.5 KB) - Inlined into 1 file
- `settings.local.json` (~0.15 KB) - Example file, not needed
- Logo file (if present)

**Net Change:**
- Before: 363 KB (59 files)
- After: ~420 KB (53 files)
- Growth: +57 KB (+16%)

### Optimization Opportunities

1. **Markdown Compression:** ~10-15% reduction possible
2. **Whitespace Normalization:** ~5% reduction
3. **Duplicate Content Detection:** Potential for further reduction if rules have common sections

**Optimized Estimate:** ~380 KB (after compression and normalization)

### Plugin Distribution Size

```
Plugin Package (tar.gz or zip):
├── Assets (compressed): ~200-250 KB
├── Manifest + Metadata: ~5 KB
├── Documentation: ~10 KB
└── Installation Scripts: ~5 KB
─────────────────────────────────
Total Compressed: ~220-270 KB
Total Uncompressed: ~420 KB
```

### Comparison to npm Package

| Metric | npm Package | Plugin |
|--------|-------------|--------|
| **Source Code** | 5,263 lines (27 files) | 0 lines (eliminated) |
| **Assets** | 363 KB (59 files) | 420 KB (53 files) |
| **Dependencies** | 9 runtime deps (React, Ink, Commander, etc.) | 0 dependencies |
| **node_modules** | ~50 MB | 0 MB |
| **Total Installed Size** | ~50 MB | ~0.5 MB |
| **Install Time** | ~30 seconds | Instant |

**Key Takeaway:** Plugin is 100x smaller than npm package due to elimination of all dependencies and build infrastructure.

---

## 8. Dependencies Analysis

### Current npm Dependencies

#### Runtime Dependencies (9 packages)

| Package | Version | Purpose | Plugin Needs? |
|---------|---------|---------|---------------|
| `chalk` | ^5.3.0 | Terminal colors | **NO** - Plugin has no CLI |
| `commander` | ^12.1.0 | CLI framework | **NO** - Plugin uses built-in commands |
| `fs-extra` | ^11.2.0 | File system utilities | **NO** - Plugin doesn't manage files |
| `ink` | ^4.4.1 | React for terminal UI | **NO** - No TUI in plugin |
| `ink-text-input` | ^5.0.1 | Text input component | **NO** - No TUI in plugin |
| `inquirer` | ^10.2.2 | Interactive prompts | **NO** - No prompts in plugin |
| `ora` | ^8.1.0 | Loading spinners | **NO** - No spinners in plugin |
| `react` | ^18.3.1 | UI framework | **NO** - No React in plugin |
| `toml` | ^3.0.0 | TOML parsing (if used) | **NO** - No TOML output |

**Total Runtime Dependencies:** 9 packages → **0 packages in plugin**

#### DevDependencies (14 packages)

All dev dependencies are eliminated (TypeScript compiler, test framework, linters, etc.).

**Total Elimination:** 23 npm packages + all transitive dependencies (~500+ packages)

### Plugin Dependencies

**Build-Time Dependencies:**
- Node.js 18+ (for build script execution)
- TypeScript (if build script is TypeScript)
- Glob library (for file collection)
- Basic file system operations (built-in)

**Runtime Dependencies:**
- **NONE** - Plugin is pure content

---

## 9. Migration Roadmap

### Phase 1: Build Infrastructure (Week 1)

**Objective:** Create build tooling to transform assets

**Tasks:**
1. Create `scripts/build-plugin.ts`
   - Implement external reference inlining
   - Implement placeholder replacement
   - Implement validation checks
2. Create `scripts/validate-assets.ts`
   - Check for unreplaced placeholders
   - Validate markdown syntax
   - Verify JSON parsing
3. Test build on sample assets
4. Document build process

**Deliverables:**
- Working build script
- Validation script
- Build documentation
- Sample plugin output

### Phase 2: Asset Transformation (Week 2)

**Objective:** Transform all 13 files requiring preprocessing

**Tasks:**
1. Transform commands (4 files)
   - Inline rule files into each command
   - Verify markdown structure preserved
   - Test command loading in Claude Code
2. Transform agents (1 file)
   - Inline agent-creation-principles into meta-agent
3. Transform output style (1 file)
   - Inline agent-delegation into output style
4. Transform settings.json
   - Replace placeholders with relative paths
   - Test settings loading
5. Adapt templates (6 files)
   - Adjust workflow references
   - Remove CLI-specific instructions

**Deliverables:**
- 13 transformed asset files
- Transformation validation report
- File size comparison report

### Phase 3: Plugin Manifest Creation (Week 3)

**Objective:** Create plugin.json manifest

**Tasks:**
1. Define plugin metadata
   - Name, version, description
   - Author and repository
2. Define installation behavior
   - Copy assets to `.claude/` and `.the-startup/`
   - Merge settings.json
   - Set up statusLine hook
3. Define commands
   - Map slash commands to markdown files
   - Define command permissions
4. Test manifest loading

**Deliverables:**
- Complete plugin.json
- Installation documentation
- Command mapping documentation

### Phase 4: Testing & Validation (Week 4)

**Objective:** Validate plugin works end-to-end

**Tasks:**
1. Install plugin in test project
2. Verify all agents load correctly
3. Test all slash commands
4. Verify statusLine hook works
5. Test on multiple platforms (Mac, Windows, Linux)
6. Performance testing (load times)

**Deliverables:**
- Test results report
- Platform compatibility matrix
- Performance benchmarks
- Bug fixes

### Phase 5: Documentation & Release (Week 5)

**Objective:** Prepare for plugin release

**Tasks:**
1. Write plugin README
   - Installation instructions
   - Usage guide
   - Command reference
2. Create migration guide (npm → plugin)
3. Update CLAUDE.md for plugin context
4. Create release notes
5. Submit to plugin marketplace

**Deliverables:**
- Complete documentation
- Migration guide
- Release notes
- Plugin submission

---

## 10. Risk Assessment

### High Priority Risks

#### Risk 1: External Reference Inlining Breaks Markdown Structure

**Impact:** High
**Likelihood:** Medium
**Mitigation:**
- Add markdown validation to build script
- Test each transformed file individually
- Preserve indentation and structure during inlining
- Use automated markdown linting

#### Risk 2: Plugin Size Exceeds Marketplace Limits

**Impact:** High
**Likelihood:** Low
**Current Size:** ~420 KB (well below typical limits)
**Mitigation:**
- Compress assets during build
- Remove unnecessary whitespace
- Monitor size during development

#### Risk 3: Placeholder Replacement Misses Edge Cases

**Impact:** Medium
**Likelihood:** Medium
**Mitigation:**
- Comprehensive regex testing
- Post-build validation scan
- Manual review of transformed files

### Medium Priority Risks

#### Risk 4: Agent Loading Performance Degrades with Inlined Content

**Impact:** Medium
**Likelihood:** Low
**Mitigation:**
- Benchmark agent load times
- Consider lazy loading if needed
- Profile plugin startup

#### Risk 5: Template Workflow Adaptation Requires UX Changes

**Impact:** Medium
**Likelihood:** High
**Mitigation:**
- Design plugin-specific template workflow
- User testing of new workflow
- Clear migration documentation

### Low Priority Risks

#### Risk 6: Platform-Specific Shell Script Selection Fails

**Impact:** Low
**Likelihood:** Low
**Mitigation:**
- Test on Windows, Mac, Linux
- Fallback to simpler statusLine if script fails

---

## 11. Success Criteria

### Functional Criteria

- [ ] All 39 agents load and execute correctly
- [ ] All 5 slash commands work as expected
- [ ] Output style applies correctly
- [ ] StatusLine hook displays git branch
- [ ] All external references resolved
- [ ] No unreplaced placeholders remain
- [ ] Templates accessible to commands
- [ ] Settings merge correctly on installation

### Performance Criteria

- [ ] Plugin size < 500 KB
- [ ] Agent load time < 100ms per agent
- [ ] Build process < 30 seconds
- [ ] Installation time < 5 seconds
- [ ] StatusLine hook < 50ms execution

### Quality Criteria

- [ ] All markdown files validate
- [ ] All JSON files parse correctly
- [ ] No broken references
- [ ] Consistent formatting
- [ ] Complete documentation
- [ ] No security vulnerabilities

### User Experience Criteria

- [ ] Installation is one-click
- [ ] No manual configuration required
- [ ] Clear error messages
- [ ] Works on Mac, Windows, Linux
- [ ] Migration path from npm package documented

---

## 12. Next Steps

### Immediate Actions (This Week)

1. **Review this analysis** with stakeholders
2. **Create build script skeleton** in `scripts/build-plugin.ts`
3. **Test inlining** on one command file (e.g., `analyze.md`)
4. **Validate transformation** preserves markdown structure
5. **Document transformation process**

### Short-Term Actions (Next 2 Weeks)

1. **Complete build infrastructure** (Phase 1)
2. **Transform all 13 files** (Phase 2)
3. **Begin plugin manifest** (Phase 3)
4. **Set up test environment** for plugin validation

### Long-Term Actions (Weeks 3-5)

1. **Complete testing** (Phase 4)
2. **Finalize documentation** (Phase 5)
3. **Submit plugin** to marketplace
4. **Publish migration guide** for npm users

---

## Appendix A: Quick Reference Tables

### Files by Migration Action

| Action | File Count | Total Size | % of Total |
|--------|-----------|------------|------------|
| **ELIMINATE** | 27 source files + 14 test files | ~4,500 lines | 85% of code |
| **KEEP AS-IS** | 44 asset files | ~238 KB | 66% of assets |
| **TRANSFORM** | 13 asset files | ~125 KB | 34% of assets |
| **TOTAL** | 98 files | 5,263 lines + 363 KB | 100% |

### Transformation Summary

| File Type | Count | External Refs | Placeholders | Est. Size Growth |
|-----------|-------|---------------|--------------|------------------|
| **Commands** | 4 | 7 | 0 | +140% |
| **Agents** | 1 | 1 | 0 | +79% |
| **Output Styles** | 1 | 1 | 0 | +106% |
| **Settings** | 1 | 0 | 3 | <1% |
| **Templates** | 6 | 0 | TBD | TBD |

### External Reference Index

| Reference File | Size | Referenced By | Times Referenced |
|----------------|------|---------------|------------------|
| `rules/agent-delegation.md` | 8.5 KB | 5 files | 5 |
| `rules/cycle-pattern.md` | 6 KB | 3 files | 3 |
| `rules/agent-creation-principles.md` | 5.5 KB | 1 file | 1 |

---

## Appendix B: Asset File Manifest

### Complete Asset List with Metadata

```csv
Category,Subcategory,File,Size (bytes),External Refs,Migration Action
Agent,the-analyst,feature-prioritization.md,5200,0,KEEP
Agent,the-analyst,project-coordination.md,5100,0,KEEP
Agent,the-analyst,requirements-analysis.md,4900,0,KEEP
Agent,the-architect,quality-review.md,4500,0,KEEP
Agent,the-architect,system-architecture.md,5200,0,KEEP
Agent,the-architect,system-documentation.md,4800,0,KEEP
Agent,the-architect,technology-research.md,3900,0,KEEP
Agent,the-architect,technology-standards.md,3800,0,KEEP
Agent,the-chief,,8000,0,KEEP
Agent,the-designer,accessibility-implementation.md,4700,0,KEEP
Agent,the-designer,design-foundation.md,4500,0,KEEP
Agent,the-designer,interaction-architecture.md,4400,0,KEEP
Agent,the-designer,user-research.md,4600,0,KEEP
Agent,the-meta-agent,,7000,1,TRANSFORM
Agent,the-ml-engineer,context-management.md,5200,0,KEEP
Agent,the-ml-engineer,feature-operations.md,4900,0,KEEP
Agent,the-ml-engineer,ml-operations.md,5300,0,KEEP
Agent,the-ml-engineer,prompt-optimization.md,4800,0,KEEP
Agent,the-mobile-engineer,mobile-data-persistence.md,5100,0,KEEP
Agent,the-mobile-engineer,mobile-development.md,5600,0,KEEP
Agent,the-mobile-engineer,mobile-operations.md,5400,0,KEEP
Agent,the-platform-engineer,containerization.md,4600,0,KEEP
Agent,the-platform-engineer,data-architecture.md,4900,0,KEEP
Agent,the-platform-engineer,deployment-automation.md,4300,0,KEEP
Agent,the-platform-engineer,infrastructure-as-code.md,4700,0,KEEP
Agent,the-platform-engineer,performance-tuning.md,4200,0,KEEP
Agent,the-platform-engineer,pipeline-engineering.md,4500,0,KEEP
Agent,the-platform-engineer,production-monitoring.md,4800,0,KEEP
Agent,the-qa-engineer,exploratory-testing.md,4600,0,KEEP
Agent,the-qa-engineer,performance-testing.md,4700,0,KEEP
Agent,the-qa-engineer,test-execution.md,4900,0,KEEP
Agent,the-security-engineer,security-assessment.md,5100,0,KEEP
Agent,the-security-engineer,security-implementation.md,4800,0,KEEP
Agent,the-security-engineer,security-incident-response.md,5200,0,KEEP
Agent,the-software-engineer,api-development.md,4900,0,KEEP
Agent,the-software-engineer,component-development.md,4700,0,KEEP
Agent,the-software-engineer,domain-modeling.md,4600,0,KEEP
Agent,the-software-engineer,performance-optimization.md,4800,0,KEEP
Agent,the-software-engineer,service-resilience.md,5100,0,KEEP
Command,s,analyze.md,8000,2,TRANSFORM
Command,s,implement.md,7000,1,TRANSFORM
Command,s,init.md,6000,0,KEEP
Command,s,refactor.md,7500,2,TRANSFORM
Command,s,specify.md,10000,2,TRANSFORM
OutputStyle,,the-startup.md,8000,1,TRANSFORM
Settings,,settings.json,200,0,TRANSFORM
Settings,,settings.local.json,150,0,ELIMINATE
Script,bin,statusline.sh,1600,0,KEEP
Script,bin,statusline.ps1,1800,0,KEEP
Rules,,agent-delegation.md,8500,0,INLINE (then eliminate)
Rules,,cycle-pattern.md,6000,0,INLINE (then eliminate)
Rules,,agent-creation-principles.md,5500,0,INLINE (then eliminate)
Template,,product-requirements.md,12000,0,ADAPT
Template,,solution-design.md,15000,0,ADAPT
Template,,implementation-plan.md,10000,0,ADAPT
Template,,definition-of-ready.md,8000,0,ADAPT
Template,,definition-of-done.md,9000,0,ADAPT
Template,,task-definition-of-done.md,7000,0,ADAPT
```

---

## Document Metadata

**Created:** 2025-10-11
**Last Updated:** 2025-10-11
**Author:** System Architect (Claude Code)
**Version:** 1.0.0
**Status:** Complete

**Related Documents:**
- `/Users/irudi/Code/personal/the-startup/PLUGIN-MIGRATION-PLAN.md`
- `/Users/irudi/Code/personal/the-startup/docs/research/claude-code-plugins-research.md`
- `/Users/irudi/Code/personal/the-startup/docs/specs/004-typescript-npm-package-migration/`

**Next Steps:**
- Review with stakeholders
- Begin Phase 1 (Build Infrastructure)
- Set up transformation testing
