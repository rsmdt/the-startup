# Discovery Perspectives

Perspective definitions and focus area mapping for constitution rule discovery.

---

## Perspectives

| Perspective | Intent | What to Discover |
|-------------|--------|------------------|
| 🔐 **Security** | Identify security patterns and risks | Auth methods, secret handling, input validation, injection prevention, CORS |
| 🏗️ **Architecture** | Understand structural patterns | Layer structure, module boundaries, API patterns, data flow, dependencies |
| 📝 **Code Quality** | Find coding conventions | Naming conventions, import patterns, error handling, logging, code organization |
| 🧪 **Testing** | Discover test practices | Test framework, file patterns, coverage requirements, mocking approaches |
| 📦 **Dependencies** | Discover package governance | License restrictions, version pinning strategy, prohibited packages, lockfile requirements, private registry configuration |
| ⚡ **Performance** | Discover performance constraints | Bundle size budgets, response time targets, query count limits, caching requirements, lazy loading mandates |

## Focus Area Mapping

| Input | Discovery Perspectives |
|-------|----------------------|
| "security" | 🔐 Security |
| "testing" | 🧪 Testing |
| "architecture" | 🏗️ Architecture |
| "code quality" | 📝 Code Quality |
| "dependencies" or "packages" | 📦 Dependencies |
| "performance" | ⚡ Performance |
| Empty or "all" | All perspectives |
| Framework-specific | Relevant subset based on framework |

## Framework-Specific Interpretation

| Input | Discovery Focus |
|-------|-----------------|
| "React" | Hooks, components, state management |
| "Next.js" | Pages, API routes, SSR patterns |
| "monorepo" | Package boundaries, shared code |
| "API" | Endpoints, validation, error handling |
