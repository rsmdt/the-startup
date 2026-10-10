# Project Constitution

Rules every change must satisfy. **L1** rules are blocking and autofixable, **L2** rules are blocking and need a manual fix, **L3** rules are advisory.

## Security

### No hardcoded secrets
- **Level:** L1
- **Scope:** `src/**/*.js`
- **Rule:** Credentials, API keys, and tokens are never string literals in source; read them from `process.env`.
- **Detect:** an assignment of a string literal to an identifier matching `/(secret|token|api_?key|password)/i`
- **Autofix:** replace the literal with `process.env.<UPPER_SNAKE_NAME>`

## Architecture

### Services do not import from routes
- **Level:** L2
- **Scope:** `src/services/**`
- **Rule:** Modules under `src/services/` never `require` anything from `src/routes/`.

## Style

### Prefer early returns
- **Level:** L3
- **Scope:** `src/**/*.js`
- **Rule:** Prefer guard clauses over nested `if` blocks.
