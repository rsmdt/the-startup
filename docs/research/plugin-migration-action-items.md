# Plugin Migration Action Items

**Based on:** claude-code-plugin-file-references-and-scripts.md
**Date:** 2025-10-11

---

## Overview

This document outlines specific action items for migrating the current CLI to a Claude Code plugin based on research findings.

---

## Action Items

### 1. File Reference Strategy

**Decision:** Inline all referenced content during build

**Current State:**
```markdown
<!-- In command files -->
@{{STARTUP_PATH}}/rules/agent-delegation.md
@{{STARTUP_PATH}}/templates/PRD.md
```

**Target State:**
```markdown
<!-- Content inlined during build -->
## Agent Delegation Rules

[Full content from agent-delegation.md]

## PRD Template

[Full content from PRD.md]
```

**Implementation:**
- [ ] Create build script to process command/agent markdown files
- [ ] Replace `@{{STARTUP_PATH}}/path` with file contents
- [ ] Maintain source files separately for maintainability
- [ ] Add source map comments for debugging

**Build Process:**
```typescript
// src/build/inliner.ts
export function inlineReferences(content: string, basePath: string): string {
  return content.replace(
    /@\{\{STARTUP_PATH\}\}\/(.+?)\.md/g,
    (match, relativePath) => {
      const fullPath = path.join(basePath, relativePath + '.md');
      const fileContent = fs.readFileSync(fullPath, 'utf-8');
      return `\n<!-- START: ${relativePath}.md -->\n${fileContent}\n<!-- END: ${relativePath}.md -->\n`;
    }
  );
}
```

---

### 2. Directory Structure Changes

**Decision:** Use `hooks/` directory per official convention

**Current Structure:**
```
assets/the-startup/
├── bin/
│   ├── statusline.sh
│   └── statusline.ps1
├── rules/
├── templates/
└── ...
```

**Target Structure:**
```
plugin-root/
├── .claude-plugin/
│   └── plugin.json
├── commands/
│   ├── prd-create.md
│   ├── prd-execute.md
│   ├── spec.md
│   └── ...
├── agents/
│   ├── the-chief.md
│   ├── the-analyst.md
│   └── ...
├── hooks/
│   ├── hooks.json
│   ├── statusline.sh
│   └── statusline.ps1
└── README.md
```

**Implementation:**
- [ ] Rename `assets/the-startup/bin/` → `hooks/`
- [ ] Move statusline scripts to `hooks/`
- [ ] Create `hooks/hooks.json` configuration
- [ ] Update build process to copy to correct locations
- [ ] Inline `rules/` and `templates/` content into commands/agents

**hooks.json:**
```json
{
  "description": "The Agentic Startup workflow hooks",
  "statusLine": {
    "type": "command",
    "command": "${CLAUDE_PLUGIN_ROOT}/hooks/statusline.sh"
  }
}
```

---

### 3. Spec Command Extraction

**Decision:** Create separate npm package for spec command

**Current State:**
- Spec command in CLI: `the-agentic-startup spec <feature>`
- Complex logic: directory numbering, TOML generation, template instantiation

**Target State:**
- Separate package: `the-agentic-startup-spec`
- Standalone CLI tool
- Can be used independently or with plugin

**Implementation:**
- [ ] Create new package: `packages/spec/` or separate repo
- [ ] Extract spec generation logic from `src/core/spec/`
- [ ] Create CLI entry point with commander
- [ ] Add TOML generation library (e.g., `@iarna/toml`)
- [ ] Implement rich CLI features (inquirer, chalk, ora)
- [ ] Publish to npm as separate package
- [ ] Update plugin command to document external tool

**Package Structure:**
```
the-agentic-startup-spec/
├── src/
│   ├── cli.ts              # CLI entry point
│   ├── generator.ts        # Spec directory generation
│   ├── toml.ts             # TOML frontmatter generation
│   ├── templates.ts        # Template management
│   └── types.ts            # TypeScript types
├── bin/
│   └── spec.js             # Executable (built from src/cli.ts)
├── templates/              # Embedded spec templates
│   ├── SPEC.md.template
│   └── PLAN.md.template
├── package.json
└── README.md
```

**Plugin Integration:**
```markdown
<!-- commands/spec.md -->
---
description: Create specification directory with auto-incrementing numbers
---

# Create Specification Directory

This command requires the `the-agentic-startup-spec` CLI tool.

## Installation

If not already installed:
```bash
npm install -g the-agentic-startup-spec
```

## Usage

To create a new specification directory:
```bash
the-agentic-startup-spec $ARGUMENTS
```

This tool ensures:
- Correct auto-incrementing directory numbers (spec-001, spec-002, etc.)
- Valid TOML frontmatter generation
- Proper template instantiation with all placeholders replaced

## Alternative

You can also ask me to create the spec directory structure, but using the CLI tool guarantees:
- Atomic directory creation
- Correct TOML syntax
- Consistent numbering even with concurrent operations
```

---

### 4. Plugin Configuration

**Decision:** Create comprehensive plugin.json

**Implementation:**
- [ ] Create `.claude-plugin/plugin.json`
- [ ] Define all commands with relative paths
- [ ] Define all agents with relative paths
- [ ] Configure hooks reference
- [ ] Add metadata (author, version, description)

**plugin.json:**
```json
{
  "name": "the-agentic-startup",
  "version": "1.0.0",
  "description": "Complete agentic development framework with PRD/SDD/PLAN workflows, 11 specialized roles, and agent delegation system",
  "author": {
    "name": "Your Name",
    "email": "your.email@example.com",
    "url": "https://github.com/yourusername"
  },
  "homepage": "https://github.com/yourusername/the-startup",
  "repository": "https://github.com/yourusername/the-startup",
  "license": "MIT",
  "keywords": [
    "workflow",
    "agents",
    "prd",
    "sdd",
    "plan",
    "specification",
    "development"
  ],
  "commands": [
    "./commands/prd-create.md",
    "./commands/prd-execute.md",
    "./commands/spec.md",
    "./commands/init.md"
  ],
  "agents": [
    "./agents/the-chief.md",
    "./agents/the-analyst.md",
    "./agents/the-architect.md",
    "./agents/the-researcher.md",
    "./agents/the-developer.md",
    "./agents/the-tester.md",
    "./agents/the-documenter.md",
    "./agents/the-reviewer.md",
    "./agents/the-coordinator.md",
    "./agents/the-optimizer.md",
    "./agents/the-monitor.md"
  ],
  "hooks": "./hooks/hooks.json"
}
```

---

### 5. Build Process Updates

**Current Build:**
```bash
npm run build  # TypeScript → JavaScript
```

**Target Build:**
```bash
npm run build  # TypeScript → JavaScript + Asset Processing
```

**Implementation:**
- [ ] Create `src/build/` directory
- [ ] Add `inliner.ts` for content inlining
- [ ] Add `plugin-builder.ts` for plugin packaging
- [ ] Update `package.json` scripts
- [ ] Create `dist-plugin/` output directory

**Build Script (package.json):**
```json
{
  "scripts": {
    "build": "npm run build:ts && npm run build:plugin",
    "build:ts": "tsc",
    "build:plugin": "tsx src/build/plugin-builder.ts",
    "build:plugin:watch": "tsx --watch src/build/plugin-builder.ts"
  }
}
```

**Plugin Builder:**
```typescript
// src/build/plugin-builder.ts
import fs from 'fs-extra';
import path from 'path';
import { inlineReferences } from './inliner.js';

const PLUGIN_ROOT = path.join(process.cwd(), 'dist-plugin');
const ASSETS_ROOT = path.join(process.cwd(), 'assets/the-startup');

async function buildPlugin() {
  // Clean output
  await fs.remove(PLUGIN_ROOT);
  await fs.ensureDir(PLUGIN_ROOT);

  // Copy and process commands
  const commandsDir = path.join(ASSETS_ROOT, 'commands');
  const commands = await fs.readdir(commandsDir);

  for (const file of commands) {
    const content = await fs.readFile(path.join(commandsDir, file), 'utf-8');
    const inlined = inlineReferences(content, ASSETS_ROOT);
    await fs.writeFile(path.join(PLUGIN_ROOT, 'commands', file), inlined);
  }

  // Copy and process agents
  const agentsDir = path.join(ASSETS_ROOT, 'agents');
  const agents = await fs.readdir(agentsDir);

  for (const file of agents) {
    const content = await fs.readFile(path.join(agentsDir, file), 'utf-8');
    const inlined = inlineReferences(content, ASSETS_ROOT);
    await fs.writeFile(path.join(PLUGIN_ROOT, 'agents', file), inlined);
  }

  // Copy hooks
  await fs.copy(
    path.join(ASSETS_ROOT, 'hooks'),
    path.join(PLUGIN_ROOT, 'hooks')
  );

  // Copy plugin metadata
  await fs.copy(
    path.join(process.cwd(), '.claude-plugin'),
    path.join(PLUGIN_ROOT, '.claude-plugin')
  );

  console.log('✅ Plugin built successfully:', PLUGIN_ROOT);
}

buildPlugin().catch(console.error);
```

---

### 6. Package Distribution

**Current Distribution:**
- npm package with CLI
- Assets in `assets/` directory
- Installation via CLI command

**Target Distribution:**

**Option A: Dual Package**
- Main package: `the-agentic-startup` (CLI for installation)
- Plugin package: `@the-agentic-startup/plugin` (pure plugin)
- Spec package: `the-agentic-startup-spec` (separate CLI)

**Option B: Plugin-Only**
- Remove CLI installer completely
- Distribute as pure plugin
- Use Claude Code's plugin system for installation
- Spec as separate package

**Recommendation: Option A (Dual Package)**

**Implementation:**
- [ ] Maintain CLI for backward compatibility
- [ ] Create plugin build output in `dist-plugin/`
- [ ] Publish plugin to plugin marketplace
- [ ] Create separate spec package
- [ ] Update installation docs

**Marketplace Configuration:**
```json
// .claude-plugin/marketplace.json
{
  "name": "the-agentic-startup-marketplace",
  "owner": {
    "name": "Your Name",
    "email": "your.email@example.com",
    "url": "https://github.com/yourusername"
  },
  "metadata": {
    "description": "The Agentic Startup plugin marketplace",
    "version": "1.0.0"
  },
  "plugins": [
    {
      "name": "the-agentic-startup",
      "source": "./",
      "description": "Complete agentic development framework",
      "version": "1.0.0",
      "author": {
        "name": "Your Name",
        "url": "https://github.com/yourusername"
      },
      "homepage": "https://github.com/yourusername/the-startup",
      "repository": "https://github.com/yourusername/the-startup",
      "license": "MIT",
      "keywords": ["workflow", "agents", "prd", "sdd"],
      "category": "workflows",
      "strict": false,
      "commands": [
        "./commands/prd-create.md",
        "./commands/prd-execute.md",
        "./commands/spec.md",
        "./commands/init.md"
      ],
      "agents": [
        "./agents/the-chief.md",
        "./agents/the-analyst.md",
        "./agents/the-architect.md",
        "./agents/the-researcher.md",
        "./agents/the-developer.md",
        "./agents/the-tester.md",
        "./agents/the-documenter.md",
        "./agents/the-reviewer.md",
        "./agents/the-coordinator.md",
        "./agents/the-optimizer.md",
        "./agents/the-monitor.md"
      ]
    }
  ]
}
```

---

## Priority Order

### Phase 1: Core Migration (Week 1)
1. ✅ Research completed
2. [ ] Create `.claude-plugin/plugin.json`
3. [ ] Rename `bin/` → `hooks/`
4. [ ] Create `hooks/hooks.json`
5. [ ] Build basic plugin structure

### Phase 2: Content Processing (Week 2)
1. [ ] Implement `inliner.ts`
2. [ ] Process all commands with inlining
3. [ ] Process all agents with inlining
4. [ ] Create `plugin-builder.ts`
5. [ ] Test plugin locally

### Phase 3: Spec Extraction (Week 3)
1. [ ] Create `the-agentic-startup-spec` package
2. [ ] Extract spec generation logic
3. [ ] Implement CLI interface
4. [ ] Add TOML generation
5. [ ] Publish to npm
6. [ ] Update plugin command to reference CLI

### Phase 4: Testing & Distribution (Week 4)
1. [ ] Test plugin installation via `/plugin` command
2. [ ] Verify all commands work
3. [ ] Verify all agents work
4. [ ] Test statusline hook
5. [ ] Create marketplace.json
6. [ ] Publish plugin
7. [ ] Update documentation

---

## Success Criteria

- [ ] Plugin installs via `/plugin` command
- [ ] All 4 commands accessible via `/` prefix
- [ ] All 11 agents accessible via Task tool
- [ ] Statusline works with hook
- [ ] Spec CLI works independently
- [ ] No runtime file reference issues
- [ ] All content properly inlined
- [ ] Documentation updated

---

## Migration Checklist

### Pre-Migration
- [x] Research plugin architecture
- [x] Identify file reference mechanisms
- [x] Determine script directory convention
- [x] Decide on spec command approach
- [ ] Review current CLI structure
- [ ] Audit all file references in commands/agents

### During Migration
- [ ] Create plugin structure
- [ ] Implement build process
- [ ] Inline all references
- [ ] Rename directories
- [ ] Configure hooks
- [ ] Extract spec command
- [ ] Update package.json
- [ ] Test locally

### Post-Migration
- [ ] Test all workflows
- [ ] Verify agent delegation
- [ ] Check statusline
- [ ] Publish packages
- [ ] Update documentation
- [ ] Announce migration
- [ ] Monitor for issues

---

## Risk Mitigation

### Risk: Content Inlining Increases File Size

**Mitigation:**
- Monitor command/agent file sizes
- Consider creating multiple smaller commands if needed
- Use markdown comments to document sources

### Risk: Spec CLI Not Installed

**Mitigation:**
- Clear installation instructions in command
- Consider auto-installation check in hook
- Provide fallback instructions for manual creation

### Risk: Breaking Changes for Existing Users

**Mitigation:**
- Maintain CLI installer in parallel
- Provide migration guide
- Version bump to 2.0.0
- Deprecation warnings

### Risk: Plugin Marketplace Approval

**Mitigation:**
- Follow all official guidelines
- Test thoroughly before submission
- Have community beta testers
- Prepare for feedback/changes

---

## Next Steps

1. **Review this document** with stakeholders
2. **Prioritize action items** based on impact
3. **Create GitHub issues** for tracking
4. **Start Phase 1** implementation
5. **Set up project board** for tracking progress

---

## References

- Main Research: `claude-code-plugin-file-references-and-scripts.md`
- Official Plugins: https://github.com/anthropics/claude-code/tree/main/plugins
- Plugin Docs: https://docs.claude.com/en/docs/claude-code/plugins-reference
