# Analysis Perspectives

Perspective definitions, activation rules, and focus area mapping.

---

## Perspectives

| Perspective | Intent | What to Discover |
|-------------|--------|------------------|
| 📋 **Business** | Understand domain logic | Business rules, validation logic, workflows, state machines, domain entities |
| 🏗️ **Technical** | Map architecture | Design patterns, conventions, module structure, dependency patterns |
| 🔐 **Security** | Identify security model | Auth flows, authorization rules, data protection, input validation |
| ⚡ **Performance** | Find optimization opportunities | Bottlenecks, caching patterns, query patterns, resource usage |
| 🔌 **Integration** | Map external boundaries | External APIs, webhooks, data flows, third-party services |

## Focus Area Mapping

| Input | Perspectives to Launch |
|-------|----------------------|
| "business" or "domain" | 📋 Business |
| "technical" or "architecture" | 🏗️ Technical |
| "security" | 🔐 Security |
| "performance" | ⚡ Performance |
| "integration" or "api" | 🔌 Integration |
| Empty or broad request | All relevant perspectives |

## Perspective-Specific Agent Focus

| Perspective | Agent Focus | Output Location |
|-------------|-------------|-----------------|
| 📋 Business | Find domain rules, identify workflows, map entities and state machines | `docs/domain/` |
| 🏗️ Technical | Map patterns, note conventions, document module structure and dependencies | `docs/patterns/` |
| 🔐 Security | Trace auth flows, document sensitive paths, identify protection mechanisms | `docs/research/` |
| ⚡ Performance | Find hot paths, caching opportunities, expensive operations, resource usage | `docs/research/` |
| 🔌 Integration | Map external APIs, trace data flows, document third-party service contracts | `docs/interfaces/` |
