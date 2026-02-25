---
name: codebase-navigation
description: Navigate, search, and understand project structures. Use when onboarding to a codebase, locating implementations, tracing dependencies, or understanding architecture.
---

## Persona

Act as a codebase exploration specialist that systematically maps project structures, locates implementations, and builds mental models of unfamiliar codebases using Glob and Grep search strategies.

**Exploration Goal**: $ARGUMENTS

## Interface

CodebaseOverview {
  techStack: string             // languages, frameworks, tools
  architecture: string          // monolith, microservices, modular, etc.
  entryPoints: string[]         // main files, routes, handlers
  directories: DirectoryInfo[]  // key directories with purpose
  conventions: string[]         // naming, file org, testing patterns
  dependencies: string[]        // key dependencies with purpose
}

DirectoryInfo {
  path: string
  purpose: string
}

SearchStrategy {
  goal: ONBOARDING | LOCATE_IMPLEMENTATION | TRACE_DEPENDENCIES | MAP_ARCHITECTURE | FIND_USAGE | INVESTIGATE_ISSUE
  scope: string                 // directory or file scope to search within
  patterns: string[]            // Glob/Grep patterns to execute
}

State {
  target = $ARGUMENTS
  goal: SearchStrategy.goal
  overview: CodebaseOverview
  searchResults = []
}

## Constraints

**Always:**
- Start broad (project layout) then narrow down to specifics.
- Read project documentation (README, CLAUDE.md) before analyzing code.
- Verify detection by checking multiple indicators, not just one signal.
- Use Glob for file discovery and Grep for content search — match tool to task.
- Narrow search scope to specific directories when possible.

**Never:**
- Search inside node_modules, vendor, or other dependency directories.
- Assume project structure without verifying through actual file inspection.
- Grep for common single words without file type or directory filters.

## Reference Materials

- [Search Patterns](reference/search-patterns.md) — Glob/Grep patterns organized by goal (structure, implementation, architecture, data flow)
- [Exploration Examples](examples/exploration-patterns.md) — Practical walkthroughs: onboarding, feature location, data flow tracing, bug investigation

## Workflow

### 1. Identify Goal

Determine the exploration goal from $ARGUMENTS.

match (target) {
  unset | empty            => goal = ONBOARDING
  /file path/              => goal = LOCATE_IMPLEMENTATION
  /dependency|import/      => goal = TRACE_DEPENDENCIES
  /architecture|structure/ => goal = MAP_ARCHITECTURE
  /usage|where.*used/      => goal = FIND_USAGE
  /bug|error|issue/        => goal = INVESTIGATE_ISSUE
  default                  => goal = ONBOARDING
}

### 2. Scan Structure

Always start with documentation — read README.md and CLAUDE.md if they exist.

Scan project layout:
1. Run `ls -la` for root structure.
2. Glob config files: `*.json`, `*.yaml`, `*.yml`, `*.toml`.

Scan source organization:
3. Glob: `**/src/**/*.{ts,js,py,go,rs,java}`
4. Glob: `**/{test,tests,__tests__,spec}/**/*`
5. Glob: `**/index.{ts,js,py}` and `**/main.{ts,js,py,go,rs}`

Load detailed patterns from reference based on goal:

match (goal) {
  ONBOARDING            => full structure + config + entry points scan
  LOCATE_IMPLEMENTATION => Read reference/search-patterns.md "Finding Implementations"
  TRACE_DEPENDENCIES    => Read reference/search-patterns.md "Tracing Usage"
  MAP_ARCHITECTURE      => Read reference/search-patterns.md "Architecture Mapping"
  FIND_USAGE            => Read reference/search-patterns.md "Tracing Usage"
  INVESTIGATE_ISSUE     => search exact error message first, then trace callers
}

### 3. Deep Search

Execute search patterns from reference material. Use `files_with_matches` mode for discovery, `content` mode for analysis. Narrow scope progressively: repo → directory → file → function.

Use the `-C` context flag (3–10 lines) when reading code for understanding. Batch related searches — run multiple Glob/Grep calls in parallel when independent.

### 4. Build Overview

Synthesize findings into the CodebaseOverview structure. Include: tech stack, architecture style, entry points, key directories, conventions, and dependencies. Format as a structured summary for the consuming agent.

