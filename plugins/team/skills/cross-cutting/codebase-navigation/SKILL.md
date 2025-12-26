---
name: codebase-navigation
description: Navigate, search, and understand project structures. Use when onboarding to a codebase, locating implementations, tracing dependencies, or understanding architecture. Provides patterns for file searching with Glob, code searching with Grep, and systematic architecture analysis.
---

# Codebase Exploration

Systematic patterns for navigating and understanding codebases efficiently.

## When to Use

- **Onboarding to a new codebase** - Understanding project structure and conventions
- **Locating specific implementations** - Finding where functionality lives
- **Tracing dependencies** - Understanding how components connect
- **Architecture analysis** - Mapping system structure and boundaries
- **Finding usage patterns** - Discovering how APIs or functions are used
- **Investigating issues** - Tracing code paths for debugging

## Quick Structure Analysis

Start broad, then narrow down. This three-step pattern works for any codebase.

### Step 1: Project Layout

```bash
# Understand top-level structure
ls -la

# Find configuration files (reveals tech stack)
ls -la *.json *.yaml *.yml *.toml 2>/dev/null

# Check for documentation
ls -la README* CLAUDE.md docs/ 2>/dev/null
```

### Step 2: Source Organization

```
# Find source directories
Glob: **/src/**/*.{ts,js,py,go,rs,java}

# Find test directories
Glob: **/{test,tests,__tests__,spec}/**/*

# Find entry points
Glob: **/index.{ts,js,py} | **/main.{ts,js,py,go,rs}
```

### Step 3: Configuration Discovery

```
# Package/dependency files
Glob: **/package.json | **/requirements.txt | **/go.mod | **/Cargo.toml

# Build configuration
Glob: **/{tsconfig,vite.config,webpack.config,jest.config}.*

# Environment/deployment
Glob: **/{.env*,docker-compose*,Dockerfile}
```

## Deep Search Strategies

### Finding Implementations

When you need to locate where something is implemented:

```
# Find function/class definitions
Grep: (function|class|interface|type)\s+TargetName

# Find exports
Grep: export\s+(default\s+)?(function|class|const)\s+TargetName

# Find specific patterns (adjust for language)
Grep: def target_name  # Python
Grep: func TargetName  # Go
Grep: fn target_name   # Rust
```

### Tracing Usage

When you need to find where something is used:

```
# Find imports of a module
Grep: import.*from\s+['"].*target-module

# Find function calls
Grep: targetFunction\(

# Find references (broad search)
Grep: TargetName
```

### Architecture Mapping

When you need to understand system structure:

```
# Find all route definitions
Grep: (app\.(get|post|put|delete)|router\.)

# Find database models/schemas
Grep: (Schema|Model|Entity|Table)\s*\(
Glob: **/{models,entities,schemas}/**/*

# Find service boundaries
Glob: **/{services,controllers,handlers}/**/*
Grep: (class|interface)\s+\w+Service
```

## Exploration Patterns by Goal

### Goal: Understand Entry Points

```
# Web application routes
Grep: (Route|path|endpoint)
Glob: **/routes/**/* | **/*router*

# CLI commands
Grep: (command|program\.)
Glob: **/cli/**/* | **/commands/**/*

# Event handlers
Grep: (on|handle|subscribe)\s*\(
```

### Goal: Find Configuration

```
# Environment variables
Grep: (process\.env|os\.environ|env\.)

# Feature flags
Grep: (feature|flag|toggle)

# Constants/config objects
Grep: (const|let)\s+(CONFIG|config|settings)
Glob: **/{config,constants}/**/*
```

### Goal: Understand Data Flow

```
# Database queries
Grep: (SELECT|INSERT|UPDATE|DELETE|find|create|update)
Grep: (prisma|sequelize|typeorm|mongoose)\.

# API calls
Grep: (fetch|axios|http\.|request\()

# State management
Grep: (useState|useReducer|createStore|createSlice)
```

## Best Practices

### Search Efficiently

1. **Start with Glob** for file discovery - faster than grep for locating files
2. **Use Grep** for content search - supports regex and context
3. **Narrow scope** - search in specific directories when possible
4. **Check output modes** - use `files_with_matches` for discovery, `content` for analysis

### Build Mental Models

1. **Map the layers** - presentation, business logic, data access
2. **Identify patterns** - repository, service, controller, etc.
3. **Note conventions** - naming, file organization, code style
4. **Document boundaries** - where modules connect and separate

### Avoid Common Pitfalls

- **Do not** search entire node_modules/vendor directories
- **Do not** assume structure without verifying
- **Do not** skip reading project documentation (README, CLAUDE.md)
- **Do not** grep for common words without filtering (use glob filters)

## Output Format

After exploration, summarize findings:

```
## Codebase Overview

**Tech Stack:** [Languages, frameworks, tools]
**Architecture:** [Monolith, microservices, modular, etc.]
**Entry Points:** [Main files, routes, handlers]

## Key Directories

- `src/` - [Purpose]
- `lib/` - [Purpose]
- `tests/` - [Purpose]

## Conventions Observed

- Naming: [Pattern]
- File organization: [Pattern]
- Testing: [Pattern]

## Dependencies

- [Key dependency]: [Purpose]
- [Key dependency]: [Purpose]
```

## References

- [Exploration Patterns Examples](examples/exploration-patterns.md) - Detailed practical examples
