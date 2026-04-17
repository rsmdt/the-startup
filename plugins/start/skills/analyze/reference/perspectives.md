# Analysis Perspectives

Perspective definitions, focus-area mapping, recommended agents, and depth expectations.

---

## Perspectives, Recommended Agents, and Doc Locations

Each perspective has a lens, a recommended subagent (the one closest to the expertise), and a canonical output location under `docs/`.

| Perspective | Lens — what to discover | Recommended Agent | Output Location |
|-------------|-------------------------|-------------------|-----------------|
| 📋 **Business** | Domain rules, validation logic, workflows, state machines, entities | `Explore` (with business-logic brief) | `docs/domain/` |
| 🏗️ **Technical** | Design patterns, conventions, module structure, dependency flow | `team:the-architect:design-system` · fallback `Explore` | `docs/patterns/` |
| 🔐 **Security** | Auth flows, authorization rules, data protection, input validation, threat model | `team:the-architect:review-security` | `docs/research/` |
| ⚡ **Performance** | Hot paths, caching, query patterns, resource usage, bottlenecks | `team:the-developer:optimize-performance` | `docs/research/` |
| 🔌 **Integration** | External APIs, webhooks, data flows, third-party contracts, failure modes | `team:the-architect:review-compatibility` · fallback `Explore` | `docs/interfaces/` |
| 💾 **Data** | Data models, schemas, relationships, migrations, indexing, access patterns | `Explore` (with data-model brief) | `docs/patterns/` |

**Agent selection rule:** use the recommended specialist when the analysis is about its core expertise (e.g., threat modeling → `review-security`, root-cause performance → `optimize-performance`). For pure discovery and mapping, `Explore` is usually faster and more focused than a decision-framing agent. When in doubt, prefer `Explore` with a perspective-specific depth brief drawn from the expectations below.

## Focus-Area Mapping

Resolve $ARGUMENTS to a perspective set:

| Input | Perspectives to Launch |
|-------|----------------------|
| "business" or "domain" | 📋 Business |
| "technical" or "architecture" | 🏗️ Technical |
| "security" | 🔐 Security |
| "performance" | ⚡ Performance |
| "integration" or "api" | 🔌 Integration |
| "data" or "schema" or "database" | 💾 Data |
| Empty or broad request | All relevant perspectives |

## Depth Expectations Per Perspective

Agents must go beyond identification. Each perspective has specific depth requirements — use these to brief the subagent and to judge whether a finding is complete.

### 📋 Business
- **Trace complete workflows** end-to-end — every branch, guard condition, and error path.
- **Map state machines** with all transitions, including what triggers them and what side effects occur.
- **Identify implicit rules** — business logic scattered across services, middleware, or validation layers that isn't centralized.

### 🏗️ Technical
- **Explain WHY patterns are used**, not just that they exist — what problem does the pattern solve here?
- **Trace dependency chains** — what depends on what, and what breaks if you change something?
- **Identify architectural boundaries** — where are the seams, what crosses them, what shouldn't?

### 🔐 Security
- **Trace the full auth flow** — from request entry to authorization decision to response, including token lifecycle.
- **Map trust boundaries** — where is input trusted vs validated, where does privilege escalation happen?
- **Identify the threat model** — what is the code protecting against, and where are the gaps?

### ⚡ Performance
- **Profile the actual hot paths** — don't guess, trace execution with evidence.
- **Quantify where possible** — O(n) vs O(n²), cache hit rates, query counts per request.
- **Explain the mechanism** — why is it slow, what's the bottleneck, what would fix it at the root?

### 🔌 Integration
- **Map the full contract** — request/response shapes, error codes, retry behavior, rate limits.
- **Trace data transformations** — what comes in, how it's mapped, what goes out.
- **Document failure modes** — what happens when the external service is down, slow, or returns unexpected data?

### 💾 Data
- **Map relationships completely** — not just "has a foreign key" but the full cardinality, cascade behavior, and query patterns.
- **Trace migration history** — how did the schema evolve, are there remnants of old designs?
- **Identify access patterns** — which queries hit which indexes, where are full table scans hiding?
