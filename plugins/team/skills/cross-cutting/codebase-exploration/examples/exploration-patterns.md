# Exploration Patterns Examples

Practical examples demonstrating codebase exploration strategies for common scenarios.

## Example 1: Onboarding to a New Project

### Context

You have been assigned to work on an unfamiliar codebase. You need to understand its structure, tech stack, and conventions before making changes.

### Pattern

```bash
# Step 1: Read existing documentation first
Read: README.md
Read: CLAUDE.md (if exists)

# Step 2: Identify tech stack from config files
Glob: package.json | requirements.txt | go.mod | Cargo.toml

# Step 3: Map source structure
Glob: **/src/**/*.{ts,js,py,go}

# Step 4: Find entry points
Glob: **/index.{ts,js} | **/main.{ts,js,py,go}

# Step 5: Locate tests to understand expected behaviors
Glob: **/*.test.{ts,js} | **/*.spec.{ts,js} | **/test_*.py
```

### Explanation

1. **Documentation first** - README and CLAUDE.md often contain project-specific conventions and commands
2. **Config files reveal stack** - package.json shows Node.js, go.mod shows Go, etc.
3. **Source mapping** - understand directory organization
4. **Entry points** - find where execution begins
5. **Tests as documentation** - tests reveal expected behaviors and usage patterns

### Variations

- For monorepos: Start with `Glob: **/package.json` to find all packages
- For microservices: Look for `Glob: **/Dockerfile` to identify service boundaries
- For legacy code: Check `Glob: **/*.sql` for database schema clues

### Anti-Patterns

- Skipping documentation and assuming structure
- Searching entire repo including node_modules/vendor
- Making changes before understanding conventions

---

## Example 2: Finding Where a Feature is Implemented

### Context

You need to modify user authentication but do not know where the code lives.

### Pattern

```
# Step 1: Search for obvious terms
Grep: (auth|login|authenticate)
  glob: **/*.{ts,js,py}
  output_mode: files_with_matches

# Step 2: Find route definitions
Grep: (login|signin|auth).*route
  glob: **/routes/**/*

# Step 3: Locate service/handler
Grep: (AuthService|AuthHandler|authenticate)
  output_mode: content
  -C: 3

# Step 4: Trace imports to find related files
Grep: import.*from.*auth
  output_mode: files_with_matches
```

### Explanation

1. **Broad search first** - Find all files mentioning authentication
2. **Route definitions** - Locate API entry points
3. **Service layer** - Find business logic implementations
4. **Import tracing** - Discover related modules and dependencies

### Variations

- For frontend: `Grep: (useAuth|AuthContext|AuthProvider)`
- For database: `Grep: (users|sessions).*table|schema`
- For middleware: `Grep: (middleware|interceptor).*auth`

### Anti-Patterns

- Searching for single common word like "user" (too many results)
- Not narrowing scope after initial discovery
- Ignoring test files (they often show usage patterns)

---

## Example 3: Understanding Data Flow

### Context

You need to understand how data moves from API request to database and back.

### Pattern

```
# Step 1: Find API routes/handlers
Grep: (app\.(get|post)|router\.(get|post)|@(Get|Post))
  glob: **/routes/**/* | **/controllers/**/*
  output_mode: content
  -C: 5

# Step 2: Find service layer calls
Grep: (Service|Repository)\.(create|find|update|delete)
  output_mode: content
  -C: 3

# Step 3: Find database operations
Grep: (prisma|sequelize|mongoose|typeorm)\.\w+\.(find|create|update)
  output_mode: content

# Step 4: Find data transformations
Grep: (map|transform|serialize|dto)
  glob: **/{dto,mapper,transformer}/**/*
```

### Explanation

1. **Start at API boundary** - Where requests enter the system
2. **Follow to services** - Business logic layer
3. **Track to database** - Data persistence layer
4. **Find transformations** - How data changes between layers

### Variations

- For event-driven: `Grep: (emit|publish|subscribe|on\()`
- For GraphQL: `Grep: (Query|Mutation|Resolver)`
- For message queues: `Grep: (queue|broker|consume|produce)`

### Anti-Patterns

- Starting from database and working up (harder to follow)
- Ignoring middleware transformations
- Missing async/await patterns that indicate I/O boundaries

---

## Example 4: Mapping Architecture Boundaries

### Context

You need to understand the high-level architecture to plan a refactoring effort.

### Pattern

```
# Step 1: Find module/package boundaries
Glob: **/package.json | **/go.mod | **/__init__.py
  (for monorepos: shows internal packages)

# Step 2: Identify layers
Glob: **/{controllers,handlers,routes}/**/*
Glob: **/{services,usecases,domain}/**/*
Glob: **/{repositories,dal,data}/**/*

# Step 3: Find external integrations
Grep: (axios|fetch|http\.|request\()
  output_mode: files_with_matches
Glob: **/{clients,integrations,adapters}/**/*

# Step 4: Map cross-cutting concerns
Glob: **/{middleware,interceptors,guards}/**/*
Grep: (logger|cache|metric|trace)
```

### Explanation

1. **Package boundaries** - Physical separation of concerns
2. **Architectural layers** - Presentation, business, data
3. **External boundaries** - Where system meets outside world
4. **Cross-cutting** - Shared functionality across layers

### Variations

- For microservices: Look for `Glob: **/Dockerfile` and service definitions
- For modular monolith: Check for internal API contracts
- For plugin architecture: `Grep: (plugin|extension|addon)`

### Anti-Patterns

- Assuming standard layer names (every project is different)
- Missing hidden dependencies through global state
- Not checking for circular dependencies

---

## Example 5: Investigating a Bug

### Context

Users report an error message: "Invalid token format". You need to find where this error originates.

### Pattern

```
# Step 1: Search for exact error message
Grep: Invalid token format
  output_mode: content
  -C: 10

# Step 2: Find related error handling
Grep: (throw|raise|Error).*token
  output_mode: content
  -C: 5

# Step 3: Trace the code path
# (after finding file, search for function callers)
Grep: validateToken|parseToken
  output_mode: files_with_matches

# Step 4: Check test files for expected behavior
Grep: Invalid token
  glob: **/*.test.{ts,js} | **/test_*.py
  output_mode: content
```

### Explanation

1. **Exact match first** - Find the error source directly
2. **Related errors** - Understand error handling context
3. **Caller analysis** - Who invokes this code
4. **Test context** - What is the expected behavior

### Variations

- For stack traces: Search for function names in the trace
- For database errors: Check migration files and schema definitions
- For runtime errors: Look for environment variable usage

### Anti-Patterns

- Guessing the source without searching
- Not checking test files for expected behavior
- Fixing symptom without understanding root cause

---

## Quick Reference: Search Strategy by Goal

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
