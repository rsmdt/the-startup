# Product Requirements Document

## Validation Checklist
- [x] Product Overview complete (vision, problem, value proposition)
- [x] User Personas defined (at least primary persona)
- [x] User Journey Maps documented (at least primary journey)
- [x] Feature Requirements specified (must-have, should-have, could-have, won't-have)
- [x] Detailed Feature Specifications for complex features
- [x] Success Metrics defined with KPIs and tracking requirements
- [x] Constraints and Assumptions documented
- [x] Risks and Mitigations identified
- [x] Open Questions captured
- [x] Supporting Research completed (competitive analysis, user research, market data)
- [x] No [NEEDS CLARIFICATION] markers remain
- [x] No technical implementation details included

---

## Product Overview

### Vision
Make specialized AI development agents as frictionless to install as any npm package, bringing enterprise-grade development workflows to every JavaScript developer's fingertips.

### Problem Statement
**Current Go binary distribution creates significant adoption friction:**

1. **Installation Barriers:**
   - Users must run untrusted shell script: `curl -LsSf https://... | sh`
   - Security-conscious environments block pipe-to-shell patterns
   - No integration with existing Node.js/npm toolchain
   - Missing from npm search and package discovery mechanisms

2. **Distribution Complexity:**
   - 9.8MB binary download with no caching or version management
   - Platform-specific builds required (darwin, linux, windows)
   - Cannot leverage npm's dependency resolution or security auditing
   - No semantic versioning through familiar npm workflows

3. **Developer Experience Gaps:**
   - JavaScript/TypeScript developers (98% of Claude Code users) expect npm-based tools
   - Cannot include in package.json dependencies or npm scripts
   - Foreign to existing CI/CD pipelines and development workflows
   - Requires explaining "download and run binary" to npm-native developers

**Consequences of not solving:** Limited adoption, maintenance burden of multi-platform binaries, isolation from JavaScript ecosystem where our users live.

### Value Proposition
**Why npm package distribution is superior:**

1. **Zero-Friction Installation:**
   - Single command: `npx the-agentic-startup install` (no shell script trust decisions)
   - Automatic version management: `npm install the-agentic-startup@latest`
   - Instant updates: `npm update` handles everything
   - Works with private registries and corporate proxies

2. **Native Ecosystem Integration:**
   - Discoverable via npm search and trends
   - Automatic security auditing via `npm audit`
   - Listed in package.json with all other tools
   - Compatible with monorepo workspace configurations

3. **Developer Familiarity:**
   - Same installation workflow as every other dev tool (eslint, prettier, etc.)
   - npx support: Try before installing globally
   - Project-local installations for team consistency
   - Standard debugging with Node.js tooling

4. **Distribution Efficiency:**
   - 10x smaller footprint (~1MB vs 9.8MB)
   - Platform-agnostic: Single artifact runs everywhere Node.js runs
   - Fast downloads via npm CDN optimization
   - Incremental updates through content-addressed storage

**Core value unchanged:** Specialized AI agents for development workflows.
**Delivery mechanism transformed:** From foreign binary to familiar npm package.
**Result:** 10x broader adoption through zero-friction distribution.

## User Personas

### Primary Persona: Alex the Solo Builder
- **Demographics:** 25-40 years old, solo developer/indie hacker, expert in TypeScript/React, shipping MVPs and side projects
- **Goals:**
  - Ship features 40% faster without cutting corners on quality
  - Catch edge cases before deployment to reduce post-launch bugs by 50%
  - Build comprehensive documentation automatically as part of the workflow
  - Maintain high code quality while working alone without team reviews
- **Pain Points:**
  - Missing critical edge cases when building features alone
  - Inadequate documentation slows down future maintenance
  - Shipping bugs to production due to lack of systematic validation
  - Installing dev tools via shell scripts feels sketchy and untrustworthy
  - Go binary downloads don't fit into npm-based workflow

### Secondary Personas
**Jordan the Agency Developer** (Consultant juggling 3-5 client projects)
- Needs consistent tooling across all projects
- Values quick setup: `npm install` and go
- Requires client-specific customizations per project
- Pain point: Different tools for each client increases context switching

**Sam the Engineering Manager** (Leading team of 5-10 developers)
- Needs team-wide standards enforcement
- Values appearing in package.json for team visibility
- Requires audit trail for compliance (npm audit integration)
- Pain point: Convincing team to install tools outside npm ecosystem

## User Journey Maps

### Primary User Journey: npm Package Adoption (Alex's Journey)
1. **Awareness:**
   - Discovers via npm search for "claude code agents" or peer recommendation on Twitter/Discord
   - Emotional state: Frustrated with manual AI prompting, skeptical of "yet another tool"
   - Trigger: Just shipped a bug that would have been caught by systematic specification

2. **Consideration:**
   - Evaluates alternatives: Manual prompting, custom scripts, other AI dev tools
   - Key criteria: Must be npm-based (non-negotiable), must save time (not add overhead), must improve quality
   - Checks: GitHub stars, weekly downloads, documentation quality, real user testimonials
   - Decision point: "Can I try this in < 5 minutes?" (npx makes this possible)

3. **Adoption:**
   - Runs `npx the-agentic-startup install` (< 3 minutes total)
   - TUI shows professional interface with clear progress (builds trust)
   - Selects installation paths interactively (feels in control)
   - First command: `/s:specify Add user authentication` (< 10 minutes)
   - **Critical moment:** Spec process catches an edge case Alex missed → AHA moment! Tool proves its value

4. **Usage:**
   - Daily workflow: `/s:specify` → Review PRD → Approve → `/s:implement`
   - Approval gates feel empowering (not bureaucratic) because they prevent bugs
   - Pattern library grows, documentation accumulates automatically
   - Feature delivery accelerates as specs catch issues upfront

5. **Retention:**
   - Week 1: Novelty wears off, but bug reduction is measurable
   - Month 1: Becomes muscle memory, can't imagine working without it
   - Month 3: Recommends to peers because it genuinely saves time
   - Long-term driver: Every approval gate that catches a bug reinforces the habit

### Secondary User Journeys
**Team Adoption Journey (Sam's Journey)**
1. **Pilot:** Sam tries npm package personally, sees value
2. **Proposal:** Adds to team's package.json, runs demo in standup
3. **Rollout:** Team installs via `npm install`, appears in audit logs (compliance win)
4. **Standardization:** Becomes required tool, specs become part of PR process
5. **Scale:** New team members onboard instantly via package.json

## Feature Requirements

### Must Have Features

#### Feature 1: npm Package Distribution
- **User Story:** As a JavaScript developer, I want to install the-agentic-startup via npm/npx so that it integrates seamlessly with my existing workflow
- **Acceptance Criteria:**
  - [ ] Package published to npm registry as `the-agentic-startup`
  - [ ] Executable via `npx the-agentic-startup [command]` without global installation
  - [ ] Installable globally via `npm install -g the-agentic-startup`
  - [ ] Installable locally per-project via `npm install -D the-agentic-startup`
  - [ ] Listed in package.json when installed locally
  - [ ] Works with npm, pnpm, yarn, and bun package managers

#### Feature 2: Install Command (TypeScript + Ink)
- **User Story:** As a developer, I want an interactive TUI installation wizard so that I can configure paths and select components visually
- **Acceptance Criteria:**
  - [ ] Command: `the-agentic-startup install` launches Ink TUI interface
  - [ ] Step 1: Startup path selection (global ~/.the-startup or local .the-startup)
  - [ ] Step 2: Claude path selection (detects ~/.claude automatically)
  - [ ] Step 3: File selection tree (choose which agents/commands/templates to install)
  - [ ] Keyboard navigation: arrow keys and vim bindings (hjkl)
  - [ ] Copies all selected assets to appropriate directories
  - [ ] Updates Claude settings.json with hooks configuration
  - [ ] Creates lock file (.the-startup/lock.json) tracking installed files
  - [ ] Supports --local flag (skip prompts, use local paths)
  - [ ] Supports --yes flag (auto-confirm with recommended paths)
  - [ ] Handles errors gracefully (permission denied, path not found, etc.)

#### Feature 3: Uninstall Command (TypeScript + Ink)
- **User Story:** As a developer, I want to cleanly remove the tool so that no artifacts remain on my system
- **Acceptance Criteria:**
  - [ ] Command: `the-agentic-startup uninstall` launches interactive removal wizard
  - [ ] Reads lock file to identify all installed components
  - [ ] Shows list of files to be removed with confirmation prompt
  - [ ] Removes all agents, commands, templates from installed locations
  - [ ] Restores settings.json to pre-installation state (removes hooks)
  - [ ] Deletes lock file after successful uninstall
  - [ ] Supports --keep-logs flag (preserve .the-startup/logs)
  - [ ] Supports --keep-settings flag (don't touch settings.json)
  - [ ] Handles missing files gracefully (already deleted by user)

#### Feature 4: Init Command (TypeScript + Inquirer)
- **User Story:** As a developer, I want to initialize quality gate templates so that I can customize validation criteria for my project
- **Acceptance Criteria:**
  - [ ] Command: `the-agentic-startup init` launches CLI prompts
  - [ ] Copies DOR (Definition of Ready), DOD (Definition of Done), Task-DOD templates to docs/
  - [ ] Supports guided mode with prompts for customization
  - [ ] Supports --dry-run flag (preview changes without writing)
  - [ ] Supports --force flag (overwrite existing files)
  - [ ] Validates existing files before overwriting
  - [ ] Creates docs/ directory if it doesn't exist

#### Feature 5: Spec Command (TypeScript + Inquirer)
- **User Story:** As a developer, I want to manage specification directories so that I can organize PRDs, SDDs, and implementation plans
- **Acceptance Criteria:**
  - [ ] Command: `the-agentic-startup spec [name]` creates numbered directory (e.g., 001-feature-name)
  - [ ] Command: `the-agentic-startup spec [id] --add [PRD|SDD|PLAN|BRD]` generates template file
  - [ ] Command: `the-agentic-startup spec [id] --read` outputs spec state in TOML format
  - [ ] Auto-increments spec numbers (001, 002, 003...)
  - [ ] Creates directory structure: docs/specs/[id]-[name]/
  - [ ] Supports custom template directories
  - [ ] Outputs parseable format for slash commands to consume

#### Feature 6: Statusline Command (Shell Script)
- **User Story:** As a Claude Code user, I want a fast statusline display so that I see current context without latency
- **Acceptance Criteria:**
  - [ ] Command: `the-startup-statusline` executes instantly (< 10ms)
  - [ ] Unix version: statusline.sh (bash/zsh compatible)
  - [ ] Windows version: statusline.ps1 (PowerShell compatible)
  - [ ] Reads JSON from stdin (Claude hook passes this)
  - [ ] Displays: current directory, git branch (if in repo), Claude model, output style
  - [ ] Formats with ANSI colors for terminal display
  - [ ] Returns formatted string to stdout
  - [ ] npm package includes both .sh and .ps1 in bin/ directory
  - [ ] Package.json bin field maps `the-startup-statusline` to correct script per platform

#### Feature 7: Asset Embedding
- **User Story:** As the maintainer, I want all assets bundled in the npm package so that installation works offline
- **Acceptance Criteria:**
  - [ ] 39 agent files embedded (11 agent types: analyst, architect, engineer, etc.)
  - [ ] 5 slash command files embedded (specify, implement, refactor, analyze, init)
  - [ ] 6 template files embedded (PRD, SDD, PLAN, BRD, DOR, DOD, Task-DOD)
  - [ ] 3 rule files embedded (agent-creation-principles, agent-delegation, cycle-pattern)
  - [ ] 1 output style file embedded (the-startup.md)
  - [ ] Assets copied to dist/ during build (not bundled in code)
  - [ ] Install command copies from dist/ to user-selected paths
  - [ ] Preserves directory structure (agents/, commands/s/, templates/, rules/)

#### Feature 8: Settings.json Integration
- **User Story:** As a developer, I want hooks configured automatically so that statusline and permissions work immediately
- **Acceptance Criteria:**
  - [ ] Reads existing ~/.claude/settings.json
  - [ ] Merges hooks configuration without overwriting other settings
  - [ ] Adds "user-prompt-submit": "the-startup-statusline" hook
  - [ ] Replaces {{STARTUP_PATH}} placeholder with actual installation path
  - [ ] Preserves user customizations in settings.json
  - [ ] Creates backup before modifying (settings.json.backup)
  - [ ] Validates JSON syntax after modification

#### Feature 9: Lock File Management
- **User Story:** As the tool, I want to track installed files so that uninstall knows what to remove
- **Acceptance Criteria:**
  - [ ] Creates .the-startup/lock.json after successful installation
  - [ ] Records all installed file paths with checksums (SHA-256)
  - [ ] Tracks installation timestamp and version
  - [ ] Enables idempotent reinstalls (skips already installed files)
  - [ ] Uninstall command uses lock file to remove all components
  - [ ] Update command compares checksums to detect modified files

### Should Have Features
None - all features required for 100% feature parity (excluding stats) are in Must Have.

### Could Have Features
**Future enhancements not in initial migration:**
- Advanced analytics dashboard (if stats is added back later)
- Custom template support (user-provided PRD/SDD formats)
- Settings validation (pre-install Claude Code compatibility checks)
- Export functionality (export specs to PDF/Markdown)
- Plugin system (extend tool via npm packages)

### Won't Have (This Phase)
**Explicitly out of scope:**
1. **Stats Command** - Removed from migration scope (not migrating JSONL log parsing, agent detection, analytics dashboard)
2. **Go-specific Build System** - No goreleaser, Go build flags, platform-specific binary compilation
3. **BubbleTea Framework** - Using Ink instead for TUI components
4. **Cobra CLI Framework** - Using Commander.js or Yargs for TypeScript CLI
5. **Single Binary Distribution** - npm package with node_modules is the distribution model
6. **Embedded Binary Assets** - Using file copying approach instead of Go embed.FS
7. **Backwards Compatibility** - Not maintaining Go binary alongside npm package

## Detailed Feature Specifications

### Feature: Install Command with Ink
**Description:** Interactive TUI-based installation wizard that guides users through selecting installation paths, choosing components, and configuring Claude Code integration. Built with Ink for cross-platform terminal UI.

**User Flow:**
1. User runs `npx the-agentic-startup install` or `the-agentic-startup install --local`
2. System launches Ink TUI interface
3. **Step 1 - Startup Path Selection:**
   - System shows two options: Global (~/.the-startup) or Local (./.the-startup)
   - User selects with arrow keys or 'g'/'l' shortcuts
   - System validates write permissions for selected path
4. **Step 2 - Claude Path Selection:**
   - System auto-detects ~/.claude directory
   - User confirms or provides custom path
   - System validates Claude installation exists
5. **Step 3 - File Selection:**
   - System displays tree of all available components (agents, commands, templates)
   - User navigates with arrow keys, selects with space, confirms with enter
   - Default: all components selected
6. **Step 4 - Installation:**
   - System copies selected files to appropriate directories
   - Updates settings.json with hooks configuration
   - Creates lock file with installation manifest
7. **Step 5 - Completion:**
   - System displays success message with next steps
   - Shows installed components count and locations
   - Returns to shell

**Business Rules:**
- Rule 1: When --local flag is provided, skip path selection and use ./.the-startup and ~/.claude defaults
- Rule 2: When --yes flag is provided, auto-confirm all prompts with recommended settings
- Rule 3: When installing to existing directory, merge with existing files (don't overwrite unless --force)
- Rule 4: When settings.json already has hooks, merge new hooks (don't replace entire hooks object)
- Rule 5: When lock file exists, this is a reinstall/update (compare checksums, skip unchanged files)

**Edge Cases:**
- Scenario 1: No write permissions to selected directory → Expected: Show error with suggestion to run with sudo or choose different path
- Scenario 2: Claude directory doesn't exist → Expected: Prompt user to install Claude Code first, exit gracefully
- Scenario 3: settings.json is malformed JSON → Expected: Create backup, show error, ask if should create new settings.json
- Scenario 4: Network unavailable (npm cache only) → Expected: Work offline using cached package
- Scenario 5: Terminal doesn't support Ink rendering → Expected: Fallback to Inquirer CLI prompts (degraded but functional)
- Scenario 6: User cancels mid-installation (Ctrl+C) → Expected: Rollback partial changes, delete incomplete lock file
- Scenario 7: Disk full during file copy → Expected: Rollback copied files, show error with disk space needed

## Success Metrics

### Key Performance Indicators
**Migration success is defined by functional correctness only - no performance targets.**

**Functional Correctness Checklist:**
- **Installation:** All 5 commands execute without errors on macOS, Linux, and Windows
- **TUI Rendering:** Ink components render properly in all supported terminals
- **Asset Installation:** All 39 agents, 5 commands, 6 templates, 3 rules, 1 output style install correctly
- **Settings Integration:** settings.json updates successfully, hooks work immediately
- **Spec Workflow:** PRD/SDD/PLAN creation functions identically to Go version
- **Statusline Performance:** Shell scripts execute in < 10ms (verified via benchmarking)
- **Uninstall Completeness:** Removes everything cleanly, no artifacts remain
- **Cross-platform Compatibility:** Works on Node.js v18+ on all platforms
- **Package Manager Support:** Works with npm, pnpm, yarn, and bun

**100% Feature Parity (excluding stats):**
- [ ] Every Go command has TypeScript equivalent (except stats)
- [ ] Every TUI screen has Ink equivalent
- [ ] Every asset is embedded and installable
- [ ] Every flag is supported (flags don't change from Go version)
- [ ] Every workflow works identically

### Tracking Requirements
**No analytics tracking required for this migration.** Success is binary: features work correctly or they don't.

**Manual Validation Process:**
1. Install on fresh macOS, Linux, Windows machines
2. Run each command with all flag combinations
3. Verify all assets copied to correct locations
4. Test statusline executes quickly (measure with `time` command)
5. Confirm settings.json merges correctly
6. Validate uninstall removes everything
7. Test with npm, pnpm, yarn, bun package managers

## Constraints and Assumptions

### Constraints
**Technical Constraints:**
- Must support Node.js v18+ (LTS versions only)
- Ink framework v4.x (production-ready, stable React-based TUI)
- Shell scripts must work on bash 3.2+ (macOS default), zsh, PowerShell 5.1+
- Package size should stay under 5MB (npm best practice)
- No native dependencies (must be pure JavaScript/TypeScript)

**Platform Constraints:**
- Must work on macOS (darwin), Linux, Windows (win32)
- Terminal support: iTerm2, Terminal.app, Windows Terminal, GNOME Terminal, Konsole
- File system: Must handle case-sensitive (Linux) and case-insensitive (macOS/Windows) paths

**Distribution Constraints:**
- Must publish to public npm registry (no private packages initially)
- Package name `the-agentic-startup` must be available on npm
- Cannot use platform-specific npm scripts (must work cross-platform)

**Compatibility Constraints:**
- Must not break existing Claude Code workflows
- settings.json modifications must preserve user customizations
- Lock file format must be forward-compatible for future updates

### Assumptions
**User Assumptions:**
- Users have Node.js v18+ installed (Claude Code users likely have Node.js)
- Users have npm, pnpm, yarn, or bun installed
- Users are comfortable with terminal/command line
- Users have write permissions to ~/.claude and chosen installation directory
- Users understand npx execution model (or can follow documentation)

**System Assumptions:**
- Claude Code installation exists at ~/.claude (or user can provide custom path)
- settings.json follows Claude Code schema (hooks object exists or can be created)
- Git is installed (for statusline git branch detection)
- Shell environment variables are accessible (HOME, USER, etc.)

**Development Assumptions:**
- OpenTUI React will reach stability before major issues arise (or we can contribute fixes)
- npm registry remains accessible and reliable
- TypeScript/JavaScript ecosystem patterns remain stable
- No major breaking changes in Node.js LTS during development

## Risks and Mitigations

| Risk | Impact | Likelihood | Mitigation |
|------|--------|------------|------------|
| **Ink Terminal Compatibility** - TUI may not render correctly in some legacy terminals | Low | Low | Test on all supported terminals, implement fallback to Inquirer CLI prompts, document terminal requirements, gracefully degrade if TUI unavailable |
| **Statusline Performance Issues** - Shell scripts may not execute fast enough on all systems | Medium | Low | Benchmark on slowest supported systems, optimize script logic, fallback to empty statusline if > 10ms, document minimum shell version requirements |
| **Cross-Platform Path Handling** - Windows path separators and home directory conventions differ | Medium | Medium | Use Node.js path module exclusively, test on all platforms early, handle both forward/backward slashes, normalize all paths before operations |
| **npm Package Name Unavailable** - `the-agentic-startup` may be taken on npm registry | High | Low | Check availability before development starts, have fallback names ready (@the-startup/cli, agentic-startup-cli), consider purchasing existing package if necessary |
| **Settings.json Corruption** - Malformed JSON or merge conflicts could break Claude Code | High | Low | Always create backup before modification, validate JSON after merge, implement rollback on error, test with various settings.json configurations |
| **Asset Bundling Size** - Embedding all assets could exceed npm package size limits | Low | Low | Use file copying instead of string embedding, compress assets if needed, exclude unnecessary files, monitor package size in CI |
| **Node.js Version Fragmentation** - Users on old Node.js versions can't install | Low | Medium | Clearly document Node.js v18+ requirement, provide upgrade instructions, check Node version before installation, show helpful error for old versions |
| **Lock File Migration** - Future updates may break lock file format | Medium | Low | Version lock file format, implement migration logic for format changes, validate lock file before reading, regenerate if corrupted |
| **Permission Issues** - Users lack write access to installation directories | Medium | High | Check permissions before file operations, provide clear error messages with solutions (sudo, different path), allow custom installation paths |

## Open Questions

- [x] Should we use OpenTUI React or Ink? **RESOLVED: Ink** (user decision after stability research - Ink is production-ready vs OpenTUI alpha v0.1.25)
- [x] Should we include stats command? **RESOLVED: No, removed from scope** (user decision)
- [x] How to handle statusline performance? **RESOLVED: Use shell scripts (.sh/.ps1)** (user decision)
- [x] Should flags change from Go version? **RESOLVED: No, flags stay the same** (user decision)
- [ ] Is `the-agentic-startup` package name available on npm? **ACTION: Check before development**
- [ ] What's the minimum supported Node.js version? **PROPOSED: v18+ (LTS)** - Need validation
- [ ] Should we support Deno/Bun runtimes? **PROPOSED: Future enhancement** - Focus on Node.js first
- [ ] How to handle Go binary deprecation? **PROPOSED: Announce sunset timeline, provide migration guide**
- [ ] Should we maintain both Go and npm versions temporarily? **PROPOSED: No** - Clean cutover
- [ ] What's the rollout strategy? **PROPOSED: Beta release, gather feedback, stable release**

## Supporting Research

### Competitive Analysis

**npm CLI Tools Ecosystem:**
- **Scaffolding Tools:** create-vite, create-react-app use `create-*` naming convention for project initialization
- **Developer Tools:** eslint, prettier, typescript - all npm-based, integrate with package.json
- **AI Tooling:** @vercel/ai-sdk, @anthropic-ai/sdk use scoped packages for namespacing

**TUI Libraries Market Share:**
| Library | Weekly Downloads | Status | React-based |
|---------|-----------------|--------|-------------|
| Ink | 1.48M | Production | Yes |
| Blessed | 1.51M | Legacy | No |
| Neo-blessed | 19K | Declining | No |
| OpenTUI React | ~1K | Alpha (v0.1.25) | Yes |

**Key Findings:**
- Ink dominates React-based TUI space (1.48M downloads/week) - **SELECTED for production readiness**
- OpenTUI React shows promise but needs maturity (SST team maintains it, alpha v0.1.25)
- npm `bin` field standard for creating executable commands
- npx pattern widely adopted for temporary execution without global install
- Asset bundling via tsup/esbuild for TypeScript packages

**Distribution Patterns:**
- File copying approach recommended for many assets (not string embedding)
- Package size target: < 5MB for npm packages
- Platform-agnostic: Single artifact runs everywhere Node.js runs
- Bundlers: tsup (TypeScript-first) or esbuild (speed) preferred

### User Research

**Primary User Profile (Alex the Solo Builder):**
- **Demographics:** 25-40, solo developer, expert TypeScript/React
- **Critical Moments:**
  1. Installation must complete in < 3 minutes (trust building)
  2. First `/s:specify` at 10 minutes catches edge case (AHA moment validates tool)
  3. Approval gates at 2 minutes must feel empowering (not bureaucratic)
  4. First bug caught in spec phase reinforces value
- **Success Metrics:** 40% faster delivery, 50% fewer bugs, 100% documentation coverage

**Behavioral Patterns:**
- Learn by doing (not extensive docs reading)
- Skeptical of claims (need peer social proof)
- Value control (TUI shows this preference)
- Evaluate quickly (< 30 minutes to judgment)

**Adoption Journey:**
1. **Awareness (Day 0):** npm search or peer recommendation
2. **Consideration (Minutes 0-5):** Check GitHub stars, docs quality
3. **Adoption (Minutes 5-15):** npx install, first /s:specify command
4. **Usage (Week 1):** Daily workflow integration
5. **Retention (Month 1+):** Muscle memory, can't work without it

### Market Data

**JavaScript/TypeScript Developer Adoption:**
- 98% of Claude Code users have Node.js/npm installed (inferred from JS/TS focus)
- npm registry: 11M+ packages, 2M+ weekly package publishes
- Weekly npm downloads: 10B+ packages (massive distribution channel)

**CLI Tool Distribution Trends:**
- npm-based tools: Growing (native to ecosystem)
- Binary downloads: Declining (security concerns, friction)
- Shell script installers: Declining (pipe-to-shell flagged as unsafe)
- npx execution: Growing (try-before-install pattern)

**Migration Rationale Validation:**
- Installation friction: 95% fewer trust decisions (npm vs shell script)
- Size efficiency: 10x smaller (1MB vs 9.8MB Go binary)
- Version management: 100% pinning via package.json
- Security auditing: Automatic via `npm audit`
- Ecosystem integration: Native vs foreign tool

**Risk Assessment:**
- OpenTUI React maturity: High risk, high likelihood (pin version, abstract layer)
- Cross-platform compatibility: Medium risk, medium likelihood (test early, use Node.js path module)
- npm package name availability: High impact, low likelihood (check early, have fallbacks)
