# Analysis Perspectives

Perspective definitions, activation rules, focus area mapping, and depth expectations.

---

## Perspectives

| Perspective | Intent | What to Discover |
|-------------|--------|------------------|
| 📋 **Business** | Understand domain logic | Business rules, validation logic, workflows, state machines, domain entities |
| 🏗️ **Technical** | Map architecture | Design patterns, conventions, module structure, dependency patterns |
| 🔐 **Security** | Identify security model | Auth flows, authorization rules, data protection, input validation |
| ⚡ **Performance** | Find optimization opportunities | Bottlenecks, caching patterns, query patterns, resource usage |
| 🔌 **Integration** | Map external boundaries | External APIs, webhooks, data flows, third-party services |
| 💾 **Data** | Map persistence layer | Data models, schemas, relationships, migrations, storage patterns, indexing strategies |

## Focus Area Mapping

| Input | Perspectives to Launch |
|-------|----------------------|
| "business" or "domain" | 📋 Business |
| "technical" or "architecture" | 🏗️ Technical |
| "security" | 🔐 Security |
| "performance" | ⚡ Performance |
| "integration" or "api" | 🔌 Integration |
| "data" or "schema" or "database" | 💾 Data |
| Empty or broad request | All relevant perspectives |

## Perspective-Specific Agent Focus

| Perspective | Agent Focus | Output Location |
|-------------|-------------|-----------------|
| 📋 Business | Find domain rules, identify workflows, map entities and state machines | `docs/domain/` |
| 🏗️ Technical | Map patterns, note conventions, document module structure and dependencies | `docs/patterns/` |
| 🔐 Security | Trace auth flows, document sensitive paths, identify protection mechanisms | `docs/research/` |
| ⚡ Performance | Find hot paths, caching opportunities, expensive operations, resource usage | `docs/research/` |
| 🔌 Integration | Map external APIs, trace data flows, document third-party service contracts | `docs/interfaces/` |
| 💾 Data | Map data models, trace schema relationships, document storage patterns and migration history | `docs/patterns/` |

## Depth Expectations Per Perspective

Agents must go beyond identification. Each perspective has specific depth requirements:

### 📋 Business
- **Trace complete workflows** end-to-end — every branch, guard condition, and error path
- **Map state machines** with all transitions, including what triggers them and what side effects occur
- **Identify implicit rules** — business logic scattered across services, middleware, or validation layers that isn't centralized

### 🏗️ Technical
- **Explain WHY patterns are used**, not just that they exist — what problem does the pattern solve here?
- **Trace dependency chains** — what depends on what, and what breaks if you change something?
- **Identify architectural boundaries** — where are the seams, what crosses them, what shouldn't?

### 🔐 Security
- **Trace the full auth flow** — from request entry to authorization decision to response, including token lifecycle
- **Map trust boundaries** — where is input trusted vs validated, where does privilege escalation happen?
- **Identify the threat model** — what is the code protecting against, and where are the gaps?

### ⚡ Performance
- **Profile the actual hot paths** — don't guess, trace the execution with evidence
- **Quantify where possible** — O(n) vs O(n²), cache hit rates, query counts per request
- **Explain the mechanism** — why is it slow, what's the bottleneck, what would fix it at the root?

### 🔌 Integration
- **Map the full contract** — request/response shapes, error codes, retry behavior, rate limits
- **Trace data transformations** — what comes in, how it's mapped, what goes out
- **Document failure modes** — what happens when the external service is down, slow, or returns unexpected data?

### 💾 Data
- **Map relationships completely** — not just "has a foreign key" but the full cardinality, cascade behavior, and query patterns
- **Trace migration history** — how did the schema evolve, are there remnants of old designs?
- **Identify access patterns** — which queries hit which indexes, where are full table scans hiding?
