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

**Critical: Identify the command for EACH category.** E2E tests often use a different runner than unit tests:

| Category | Common Separate Runners |
|----------|------------------------|
| Unit/Integration | `vitest`, `jest`, `pytest`, `go test` |
| E2E (browser) | `playwright test`, `cypress run`, `selenium` |
| E2E (API) | `supertest` (usually runs via unit runner), `hurl`, `bruno` |

If E2E config files exist (`playwright.config.*`, `cypress.config.*`) but no E2E script is in `package.json`, flag this — the tests exist but may not be wired up. Record each command separately in State.

## Step 4: Check for Related Commands

Look for additional quality commands that should pass alongside tests:
- Lint: `npm run lint`, `ruff check`, `cargo clippy`
- Type check: `npm run typecheck`, `mypy`, `cargo check`
- Format check: `npm run format:check`, `ruff format --check`

## Step 5: Capture Cost Signals

These signals are static (no test execution required) and inform whether an audit will pay off. Capture them whenever target is `audit`, and opportunistically during full-suite runs.

### Test-to-source ratio per directory

Compute `test_LOC / source_LOC` for each major area (per feature, per module, per service). Heuristics:

| Ratio | Signal |
|---|---|
| < 0.3 : 1 | Likely under-tested; check for gaps in critical-path code |
| 0.3 – 0.7 : 1 | Healthy zone for most layers |
| > 0.7 : 1 | Audit hotspot — combine with mock density to decide if over-tested |
| > 1.0 : 1 | Strong signal of over-testing or excessive boilerplate; audit recommended |

Ratio alone is not a verdict — security-critical code legitimately tests heavily, and parametrized tests inflate LOC without inflating duplication. Use the ratio to direct attention, not to judge.

**Adjust the bands for typed vs dynamic codebases.** In statically-typed codebases (TypeScript, Rust, Go, Kotlin, Java with strict checks), the type system absorbs a class of tests that dynamic-language projects need to write — prop pass-through, field-presence checks, lookup-key coverage, signature contracts. Expect typed projects to sit ~0.2–0.3 lower on this scale: 0.4 : 1 is healthy for TypeScript where 0.6 : 1 would be healthy for Python. A typed codebase at > 0.7 : 1 with high mock density is almost always over-testing the type system; flag it as a hotspot. See `test-design-rules.md` § "The type system is part of your coverage."

### Mock density per test file

Count occurrences of mocking primitives per file: `mock`, `Mock`, `MagicMock`, `AsyncMock`, `stub`, `spy`, `jest.fn`, `vi.fn`, `mocker`, `patch`, etc. (Adapt to the language.) Files with high mock density and no corresponding `Fake` / `InMemory` class are candidates for the "rewrite with fakes or delete" review per reference/test-design-rules.md.

Pair this with content inspection: `grep` for `assert_called_with`, `assert_awaited`, `toHaveBeenCalled`, etc. — the call-sequence-assertion frequency is a stronger signal than mock count alone.

### Per-category runtime

Once tests have been run, record duration per category. Slow suites have higher deletion ROI: every removed redundant test pays back on every CI run. Flag categories whose runtime exceeds 2× the size-weighted expectation as candidates for `audit`.

### Tautology / wrapper signal

Identify candidate tautology files cheaply:
- Test files under ~60 lines whose corresponding source file is also small (< 30 lines)
- Test files where the imports are exclusively from one source module (no cross-module logic)
- Test names that mirror function names with no scenario suffix (e.g., `test_format_key`, `test_get_id`) — often signal "test the function exists" rather than "test the behavior"

These are heuristics, not verdicts. Confirm by reading the file pair before recommending deletion.

### Recording cost signals

In State, populate `costSignals`:

```
costSignals: {
  runtimePerCategory: { unit: "12s", integration: "1m32s", e2e: "4m21s" },
  mockDensity:        { "tests/unit/.../test_x.py": 47, ... },
  testToSourceRatio:  { "src/features/foo/": 0.92, "src/core/redis/": 1.4, ... },
  hotspots:           [ "src/core/redis/", "src/features/subscription/" ]
}
```

Hotspots are directories meeting two of: ratio > 0.7, high mock density, no fake-based tests, runtime in top quartile.
