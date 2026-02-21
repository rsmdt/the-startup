# Test Discovery Protocol

Before running tests, understand the project's test infrastructure. Discover in this order.

---

## Step 1: Identify Test Runner & Configuration

Search for test configuration files:

| File | Runner | Ecosystem |
|------|--------|-----------|
| `package.json` (scripts.test) | npm/yarn/pnpm/bun | Node.js |
| `jest.config.*` | Jest | Node.js |
| `vitest.config.*` | Vitest | Node.js |
| `.mocharc.*` | Mocha | Node.js |
| `playwright.config.*` | Playwright | Node.js (E2E) |
| `cypress.config.*` | Cypress | Node.js (E2E) |
| `pytest.ini`, `pyproject.toml`, `setup.cfg` | pytest | Python |
| `Cargo.toml` | cargo test | Rust |
| `go.mod` | go test | Go |
| `build.gradle*`, `pom.xml` | JUnit/TestNG | Java |
| `Makefile` (test target) | make | Any |
| `Taskfile.yml` (test task) | task | Any |
| `.github/workflows/*` | CI config | Any (check for test commands) |

## Step 2: Locate Test Files

Common patterns:
- `**/*.test.{ts,tsx,js,jsx}` — Co-located tests
- `**/*.spec.{ts,tsx,js,jsx}` — Co-located specs
- `__tests__/**/*` — Test directories
- `tests/**/*` — Top-level test dir
- `test/**/*` — Alternative test dir
- `*_test.go` — Go tests
- `test_*.py`, `*_test.py` — Python tests
- `**/*_test.rs` — Rust tests

## Step 3: Assess Test Suite Scope

Count and categorize:
- **Unit tests** — isolated component/function tests
- **Integration tests** — cross-module/service tests
- **E2E tests** — browser/API end-to-end tests
- **Other** — snapshot, performance, accessibility tests

## Step 4: Check for Related Commands

Look for additional quality commands that should pass alongside tests:
- Lint: `npm run lint`, `ruff check`, `cargo clippy`
- Type check: `npm run typecheck`, `mypy`, `cargo check`
- Format check: `npm run format:check`, `ruff format --check`
