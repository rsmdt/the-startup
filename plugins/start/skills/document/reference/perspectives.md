# Documentation Perspectives

Perspective definitions, target mapping, and documentation standards.

---

## Perspectives

| Perspective | Intent | What to Document |
|-------------|--------|------------------|
| 📖 **Code** | Make code self-explanatory | Functions, classes, interfaces, types with JSDoc/TSDoc/docstrings |
| 🔌 **API** | Enable integration | Endpoints, request/response schemas, authentication, error codes, OpenAPI spec |
| 📘 **README** | Enable quick start | Features, installation, configuration, usage examples, troubleshooting |
| 📊 **Audit** | Identify documentation gaps | Coverage metrics, stale docs, missing documentation, prioritized backlog *(meta-action: informs which other perspectives to run, not a parallel work stream)* |
| 🗂️ **Capture** | Preserve discoveries | Business rules → `docs/domain/`, technical patterns → `docs/patterns/`, external integrations → `docs/interfaces/` |
| 🏛️ **Architecture** | Document system design decisions | ADRs for key decisions, module/component overviews, data flow diagrams, integration topology, technology rationale |

## Target Mapping

| Target | Perspectives to Launch |
|--------|----------------------|
| File/Directory | 📖 Code |
| `api` | 🔌 API + 📖 Code (for handlers) |
| `readme` | 📘 README |
| `audit` | 📊 Audit (all areas) |
| `capture` or pattern/rule/interface discovery | 🗂️ Capture |
| `architecture` or `adr` | 🏛️ Architecture |
| `all` or empty | All applicable perspectives |

## Documentation Standards

Every documented element should have:
1. **Summary** — One-line description
2. **Parameters** — All inputs with types and descriptions
3. **Returns** — Output type and description
4. **Throws/Raises** — Possible errors
5. **Example** — Usage example (for public APIs)
