# Search Patterns Reference

Detailed search strategies organized by goal. Load when the agent needs specific search patterns for a task.

---

## Structure Analysis

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
Glob: **/src/**/*.{ts,js,py,go,rs,java}    # Source directories
Glob: **/{test,tests,__tests__,spec}/**/*   # Test directories
Glob: **/index.{ts,js,py} | **/main.{ts,js,py,go,rs}  # Entry points
```

### Step 3: Configuration Discovery

```
Glob: **/package.json | **/requirements.txt | **/go.mod | **/Cargo.toml
Glob: **/{tsconfig,vite.config,webpack.config,jest.config}.*
Glob: **/{.env*,docker-compose*,Dockerfile}
```

---

## Finding Implementations

```
# Function/class definitions
Grep: (function|class|interface|type)\s+TargetName

# Exports
Grep: export\s+(default\s+)?(function|class|const)\s+TargetName

# Language-specific
Grep: def target_name     # Python
Grep: func TargetName     # Go
Grep: fn target_name      # Rust
```

## Tracing Usage

```
# Imports of a module
Grep: import.*from\s+['"].*target-module

# Function calls
Grep: targetFunction\(

# Broad references
Grep: TargetName
```

## Architecture Mapping

```
# Route definitions
Grep: (app\.(get|post|put|delete)|router\.)

# Database models/schemas
Grep: (Schema|Model|Entity|Table)\s*\(
Glob: **/{models,entities,schemas}/**/*

# Service boundaries
Glob: **/{services,controllers,handlers}/**/*
Grep: (class|interface)\s+\w+Service
```

---

## Patterns by Goal

### Entry Points

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

### Configuration

```
# Environment variables
Grep: (process\.env|os\.environ|env\.)

# Feature flags
Grep: (feature|flag|toggle)

# Constants/config objects
Grep: (const|let)\s+(CONFIG|config|settings)
Glob: **/{config,constants}/**/*
```

### Data Flow

```
# Database queries
Grep: (SELECT|INSERT|UPDATE|DELETE|find|create|update)
Grep: (prisma|sequelize|typeorm|mongoose)\.

# API calls
Grep: (fetch|axios|http\.|request\()

# State management
Grep: (useState|useReducer|createStore|createSlice)
```

---

## Quick Reference

| Goal | Primary Tool | Pattern |
|------|--------------|---------|
| Find file by name | Glob | `**/target-name*` |
| Find file by content | Grep | `pattern` with `files_with_matches` |
| Understand function | Grep | Function name with `-C: 10` for context |
| Find all usages | Grep | Call pattern with `files_with_matches` |
| Map directory structure | Glob | `**/src/**/*` |
| Find configuration | Glob | `**/*.{json,yaml,toml,env}` |
| Trace dependencies | Grep | Import/require patterns |
| Find tests | Glob | `**/*.test.* | **/*.spec.* | **/test_*` |
