# Example: Sequential Build Delegation

This example shows how to coordinate dependent implementation tasks that must execute sequentially.

## Scenario

**User Request:** "Implement JWT authentication for our REST API"

## Task Decomposition

```
Original Task: Implement JWT authentication for REST API

Activities Identified:
1. Design authentication database schema
   - Expertise: Database design
   - Output: Schema design document
   - Dependencies: None

2. Create database migration
   - Expertise: Database implementation
   - Output: Migration files
   - Dependencies: Activity 1 (schema design)

3. Implement authentication middleware
   - Expertise: Backend development
   - Output: JWT middleware code
   - Dependencies: Activity 2 (tables must exist)

4. Create auth endpoints (login/logout)
   - Expertise: Backend development
   - Output: Auth routes and controllers
   - Dependencies: Activity 3 (middleware needed)

5. Add tests for auth flow
   - Expertise: Test automation
   - Output: Integration tests
   - Dependencies: Activity 4 (endpoints must work)

Execution Strategy: Sequential
Reasoning: Each activity depends on the previous one's output. No parallelization possible in this chain.

Dependency Chain: 1 → 2 → 3 → 4 → 5
```

## Agent Prompts Generated

### Agent 1: Schema Design (First)

```
DISCOVERY_FIRST: Before starting, understand the environment:
    - Check existing database structure: ls migrations/ or db/schema/
    - Identify database system: PostgreSQL, MySQL, SQLite?
    - Review existing table patterns: user-related tables
    - Note naming conventions: snake_case, camelCase, PascalCase?

FOCUS: Design database schema for JWT authentication
    - Users table (if not exists) with email, password_hash
    - Sessions table for active JWT tokens
    - Include created_at, updated_at timestamps
    - Design appropriate indexes for performance
    - Plan foreign key relationships

EXCLUDE:
    - Don't create the migration yet (next task)
    - Don't implement OAuth tables (separate feature)
    - Don't modify existing user table if it exists
    - Don't add two-factor auth tables (not in scope)

CONTEXT: Building JWT authentication for REST API.
    - From security requirements: bcrypt hashing, cost factor 12
    - Session expiry: 24 hours
    - Email must be unique
    - Follow project database conventions exactly

OUTPUT: Schema design document at docs/patterns/auth-database-schema.md
    - Table definitions with column types
    - Indexes and constraints
    - Foreign key relationships
    - Example data

SUCCESS: Schema designed and documented
    - All required fields included
    - Performance indexes identified
    - Follows project conventions
    - Ready for migration implementation

TERMINATION:
    - Design complete and documented
    - Blocked by missing existing schema info
```

### Agent 2: Database Migration (After Agent 1)

```
DISCOVERY_FIRST: Before starting, understand the environment:
    - Check migration system: Knex, Sequelize, TypeORM, raw SQL?
    - Find migration directory location
    - Review existing migration file format
    - Note up/down pattern used

FOCUS: Create database migration for authentication tables
    - Implement users table from schema design
    - Implement sessions table from schema design
    - Add all indexes from schema design
    - Create both up (create) and down (drop) migrations
    - Follow migration naming conventions

EXCLUDE:
    - Don't run the migration yet (separate step)
    - Don't seed data (separate concern)
    - Don't modify existing migrations
    - Don't add tables not in schema design

CONTEXT: Implementing schema designed in previous step.
    - Schema document: docs/patterns/auth-database-schema.md
    - Tables: users (email, password_hash, created_at, updated_at)
    - Tables: sessions (id, user_id, token_hash, expires_at, created_at)
    - Indexes: users.email (unique), sessions.token_hash, sessions.user_id

OUTPUT: Migration file at [DISCOVERED_MIGRATION_PATH]/[timestamp]_create_auth_tables.js
    - Up migration creates tables and indexes
    - Down migration drops tables cleanly
    - Follows project migration format

SUCCESS: Migration created and ready to run
    - Matches schema design exactly
    - Both up and down work correctly
    - Follows project patterns
    - No syntax errors

TERMINATION:
    - Migration file created successfully
    - Blocked by unclear migration system
    - Migration format doesn't match project (document issue)
```

### Agent 3: JWT Middleware (After Agent 2)

```
DISCOVERY_FIRST: Before starting, understand the environment:
    - Find existing middleware location: src/middleware/ or app/middleware/?
    - Check JWT library in use: jsonwebtoken, jose, other?
    - Review existing middleware patterns
    - Note error handling style used

FOCUS: Implement JWT authentication middleware
    - Verify JWT token from Authorization header
    - Decode and validate token
    - Check token against sessions table (not blacklisted)
    - Attach user object to request
    - Handle missing/invalid/expired tokens appropriately

EXCLUDE:
    - Don't implement login/logout endpoints (next task)
    - Don't implement token refresh (not in scope)
    - Don't add rate limiting (separate concern)
    - Don't implement permission checking (just authentication)

CONTEXT: JWT middleware for REST API authentication.
    - JWT secret: from environment variable JWT_SECRET
    - Token expiry: 24 hours
    - Sessions table: check if token_hash exists and not expired
    - Error responses: 401 for invalid/missing token
    - Follow project error handling patterns

OUTPUT: Middleware file at [DISCOVERED_LOCATION]/auth.middleware.ts
    - verifyJWT function (main middleware)
    - Helper functions (decode, validate, etc.)
    - Error handling for all cases
    - TypeScript types if applicable

SUCCESS: Middleware implemented and ready to use
    - Verifies JWT correctly
    - Checks session validity
    - Attaches user to request
    - Handles all error cases
    - Follows project patterns

TERMINATION:
    - Middleware complete
    - Blocked by missing JWT library
    - Database connection pattern unclear
```

### Agent 4: Auth Endpoints (After Agent 3)

```
DISCOVERY_FIRST: Before starting, understand the environment:
    - Find routes directory: src/routes/ or app/routes/?
    - Check router library: Express, Fastify, Koa?
    - Review existing route patterns
    - Note controller/handler organization

FOCUS: Implement authentication endpoints
    - POST /api/auth/login (email, password → JWT token)
    - POST /api/auth/logout (invalidate current session)
    - Use auth middleware from previous step
    - Create sessions table entry on login
    - Remove sessions table entry on logout

EXCLUDE:
    - Don't implement registration (separate feature)
    - Don't implement password reset (separate feature)
    - Don't add OAuth endpoints (separate feature)
    - Don't implement GET /me endpoint (separate task)

CONTEXT: Auth endpoints using JWT middleware implemented in previous step.
    - Middleware location: [from Agent 3 OUTPUT]
    - Login validates: email format, password against bcrypt hash
    - Login creates: session entry, JWT token with user_id claim
    - Logout requires: valid JWT (use middleware)
    - Logout invalidates: session entry in database

OUTPUT:
    - Routes file at [DISCOVERED_LOCATION]/auth.routes.ts
    - Controller file at [DISCOVERED_LOCATION]/auth.controller.ts
    - Follow project route organization pattern

SUCCESS: Endpoints implemented and functional
    - Login validates credentials correctly
    - Login returns valid JWT
    - Logout requires authentication
    - Logout invalidates session
    - All errors handled appropriately
    - Follows project patterns

TERMINATION:
    - Endpoints complete
    - Blocked by unclear auth middleware interface
    - bcrypt usage pattern unclear
```

### Agent 5: Integration Tests (After Agent 4)

```
DISCOVERY_FIRST: Before starting, understand the environment:
    - Find test directory: tests/, __tests__/, spec/?
    - Identify test framework: Jest, Mocha, Vitest?
    - Review existing test patterns
    - Note test database setup approach

FOCUS: Create integration tests for authentication flow
    - Test successful login (valid credentials → JWT returned)
    - Test failed login (invalid credentials → 401)
    - Test logout (valid JWT → session invalidated)
    - Test protected route with valid JWT (→ 200)
    - Test protected route with invalid JWT (→ 401)
    - Test protected route with no JWT (→ 401)

EXCLUDE:
    - Don't test registration (not implemented)
    - Don't test password reset (not implemented)
    - Don't unit test internal functions (integration tests only)
    - Don't test OAuth flows (not implemented)

CONTEXT: Testing JWT authentication implemented in previous steps.
    - Endpoints: POST /api/auth/login, POST /api/auth/logout
    - Middleware: JWT verification from auth.middleware.ts
    - Test database: Use test database, clean between tests
    - Follow project testing patterns exactly

OUTPUT: Test file at [DISCOVERED_LOCATION]/auth.integration.test.ts
    - All test cases listed in FOCUS
    - Setup and teardown (database, test user)
    - Clear test descriptions
    - Assertions verify correct behavior

SUCCESS: Complete test coverage for auth flow
    - All tests pass
    - Tests are independent (can run in any order)
    - Database cleanup works correctly
    - Follows project test patterns
    - Coverage includes happy path and error cases

TERMINATION:
    - Tests complete and passing
    - Blocked by test framework unclear
    - Test database setup pattern unclear
```

## Execution Flow

### Phase 1: Schema Design (Agent 1)

```
🚀 Launching Agent 1: Schema Design

Status: RUNNING
Time: 15 minutes

Agent 1: COMPLETE ✅
Output: docs/patterns/auth-database-schema.md

Deliverables:
- Users table design (5 columns + indexes)
- Sessions table design (5 columns + indexes)
- Foreign key relationship documented
- Performance considerations noted

Validation: ✅ Ready for migration implementation
```

### Phase 2: Database Migration (Agent 2)

```
🚀 Launching Agent 2: Database Migration
Dependencies: Agent 1 complete ✅

Context provided:
- Schema document from Agent 1
- Table definitions and indexes

Status: RUNNING
Time: 20 minutes

Agent 2: COMPLETE ✅
Output: migrations/20250124120000_create_auth_tables.js

Deliverables:
- Up migration creates both tables
- Down migration drops cleanly
- All indexes included
- Follows Knex migration pattern

Validation: ✅ Ready for middleware implementation
```

### Phase 3: JWT Middleware (Agent 3)

```
🚀 Launching Agent 3: JWT Middleware
Dependencies: Agent 2 complete ✅

Context provided:
- Migration created tables
- Database schema known

Status: RUNNING
Time: 30 minutes

Agent 3: COMPLETE ✅
Output: src/middleware/auth.middleware.ts

Deliverables:
- verifyJWT middleware function
- Token validation logic
- Session checking against database
- Error handling for all cases

Validation: ✅ Ready for endpoint implementation
```

### Phase 4: Auth Endpoints (Agent 4)

```
🚀 Launching Agent 4: Auth Endpoints
Dependencies: Agent 3 complete ✅

Context provided:
- Middleware from Agent 3
- Function signatures and usage

Status: RUNNING
Time: 35 minutes

Agent 4: COMPLETE ✅
Output:
- src/routes/auth.routes.ts
- src/controllers/auth.controller.ts

Deliverables:
- POST /api/auth/login endpoint
- POST /api/auth/logout endpoint
- Integration with middleware
- Error responses

Validation: ✅ Ready for testing
```

### Phase 5: Integration Tests (Agent 5)

```
🚀 Launching Agent 5: Integration Tests
Dependencies: Agent 4 complete ✅

Context provided:
- All endpoints from Agent 4
- Middleware from Agent 3
- Expected behavior

Status: RUNNING
Time: 25 minutes

Agent 5: COMPLETE ✅
Output: tests/integration/auth.test.ts

Deliverables:
- 6 integration tests
- Database setup/teardown
- All tests passing ✅
- 95% coverage of auth flow

Validation: ✅ Feature complete and tested
```

## Results

### Total Time: 125 minutes (sequential)

**Sequential necessary:** Each task depends on previous
**No parallelization possible** in this dependency chain

### Context Accumulation

Each agent received growing context:

- Agent 1: Fresh start
- Agent 2: Agent 1's schema design
- Agent 3: Agent 2's migration + Agent 1's schema
- Agent 4: Agent 3's middleware + all prior context
- Agent 5: All previous implementations

### Deliverables

```
📁 Project structure:
├── docs/patterns/
│   └── auth-database-schema.md (Agent 1)
├── migrations/
│   └── 20250124120000_create_auth_tables.js (Agent 2)
├── src/
│   ├── middleware/
│   │   └── auth.middleware.ts (Agent 3)
│   ├── routes/
│   │   └── auth.routes.ts (Agent 4)
│   └── controllers/
│       └── auth.controller.ts (Agent 4)
└── tests/integration/
    └── auth.test.ts (Agent 5)
```

## Lessons Learned

### What Worked Well

✅ **Clear dependency chain:** Each agent knew exactly what it needed
✅ **Context accumulation:** Prior outputs informed each subsequent agent
✅ **DISCOVERY_FIRST:** Ensured consistency with project patterns
✅ **Validation at each step:** Caught issues before they propagated

### Challenges Encountered

⚠️ **Agent 2 Issue:** Initial migration didn't match project format
- **Solution:** Retry with more specific Knex pattern in CONTEXT
- **Lesson:** DISCOVERY_FIRST examples critical

⚠️ **Agent 3 Issue:** Used wrong JWT library (jose instead of jsonwebtoken)
- **Solution:** More explicit in EXCLUDE and CONTEXT
- **Lesson:** Specify exact libraries when project uses specific ones

⚠️ **Agent 5 Issue:** Tests didn't clean up database properly
- **Solution:** Retry with explicit teardown requirements
- **Lesson:** Test isolation must be explicit in SUCCESS criteria

### Improvements for Next Time

1. **Specify exact libraries** in CONTEXT (don't assume agent will discover)
2. **Include output of previous agent** verbatim in next agent's CONTEXT
3. **Validation step between agents** to catch issues before next dependency
4. **Checkpoint approach:** Allow user to review after each agent before launching next

## Reusable Template

This pattern works for any sequential build:

```
1. Identify dependency chain (what depends on what)
2. Order activities by dependencies
3. Each agent's CONTEXT includes prior agent outputs
4. Launch sequentially, validate each before next
5. Accumulate context as you go
```

**Use when:**
- Building implementation layers (DB → Logic → API → UI)
- Pipeline-style workflows (Design → Build → Test → Deploy)
- Learning workflows (Research → Design → Implement → Validate)
- Any task where B genuinely needs A's output

## Comparison: What If We Tried Parallel?

**Attempted parallel (would fail):**

```
Agent 2 (Migration): Needs Agent 1's schema design → BLOCKED
Agent 3 (Middleware): Needs Agent 2's tables to exist → BLOCKED
Agent 4 (Endpoints): Needs Agent 3's middleware → BLOCKED
Agent 5 (Tests): Needs Agent 4's endpoints → BLOCKED

Result: All agents blocked or produce incorrect results
```

**Lesson:** Don't force parallelization when dependencies exist. Sequential is correct here.
