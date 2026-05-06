# Test Design Rules

Language-agnostic rules for deciding what deserves a test, what shape it should take, and what to do with tests that produce no signal. Apply these during `audit` mode and whenever a failing test makes you ask "should this test even exist?"

A test produces **signal** when it can fail for a reason a caller's test wouldn't already catch. A test produces **noise** when it pins implementation details, restates the framework's contract, or only verifies that mocks were called. Signal earns the test its keep; noise costs maintenance without buying coverage.

---

## When TO write a test (positive criteria)

### 1. Pure logic with real branches

Domain calculations, state transitions, parsing, math, formatting that has rules (not just templating). Inputs in, outputs out, no I/O.

- Status computation across date ranges
- Quota math with overflow / clamping
- Retry / backoff schedules
- Permission resolution given a role and resource
- Sanitizer regex behavior

These tests are cheap, refactor-safe, and catch real bugs.

### 2. Security and correctness gates

Code where one missed branch is catastrophic. Tests are usually small and dense.

- Input sanitizers and validators
- Path traversal and absolute-path rejection
- Authentication and authorization checks
- Exception → status code / error code mapping
- Cryptographic helpers (constant-time compare, nonce generation)

The cost-of-miss justifies the test even if the function is small.

### 3. Stateful coordinators — using fakes, not mocks

Multi-step orchestrators where the assertion is the *outcome* of N steps: "after these calls, the state is X, these events were emitted, this side effect happened."

Build in-memory functional fakes for collaborators (a `FakeRepository` with a dict, a `FakeEventBus` with a list, a `FakeBlobStore` with a map). Run the real coordinator end-to-end against the fakes. Assert outcomes, not call sequences.

This is the highest-value pattern in any test suite: it exercises real logic without integration overhead, survives refactors, and catches the bugs mock-heavy tests miss.

### 4. Edge cases an integration test wouldn't naturally hit

Error paths, timeouts, malformed inputs, concurrency edge cases, configuration errors. These are where unit tests beat integration tests.

---

## When NOT to write a test (negative criteria)

### 1. Trivial wrappers and formatters

A function that returns a literal, a template, or the result of a single library call has no failure mode worth testing. Any caller's test catches breakage.

- `return f"prefix:{x}"` style key formatters
- `return self._client.do(x)` thin pass-throughs
- Constant-returning functions
- Single-call delegations

Cost: maintenance overhead. Benefit: zero — the test asserts what the implementation literally is.

### 2. Framework / library behavior

Schema libraries, ORMs, routers, and serialization libraries are well-tested by their authors. Your test should not re-verify their contract.

- "This required field is required" (the schema lib enforces this)
- "This default value is the default" (the language enforces this)
- "This roundtrip preserves data" (the serializer's job)
- "This route is reachable" (the framework's job)
- "This ORM session method was called" (the ORM's job)

Test your custom validators and your business logic. Don't test the runtime.

### 3. Repository / data-access layers with mocked clients

If the test mocks the database session / driver / client, no real query runs. The test verifies that you call the mock's methods in the right order — not that the SQL is correct, the schema matches, or the data persists.

These tests pass while real bugs (broken query, wrong column, missing index, transaction ordering) ship to production. Integration tests against a real (or in-memory) store are the spec for this layer. Either delete the unit test or replace it with a fake-backed test that exercises real logic.

### 4. Identity-mapped service tests

Pattern: service calls `dep.method()`, test mocks `dep` to return X, assertion checks the service returned X. The service does no transformation worth verifying — the test only proves the mock works.

If the service has real logic (transforms, filters, validates), test that logic. If it doesn't, the test is noise.

### 5. Call-sequence assertions as the primary check

Tests whose primary assertion is `mock.assert_called_with(...)` or `mock.call_count == N` verify the implementation, not the behavior. Refactoring to call a different (still correct) method breaks the test without revealing a bug.

Replace with outcome assertions: state changed, event emitted, exception raised, side effect happened.

### 6. API endpoints already covered end-to-end

If an integration test exercises the route through a real HTTP client with real auth and asserts the response, a unit test that mocks every service the route calls and asserts the mocks were called adds no signal. Delete the unit test or narrow it to specifically what integration can't reach (e.g., a particular middleware error mapping).

### 7. Tests pinning settings / constants

A test that asserts `settings.timeout == 30` against `timeout: int = 30` is tautological. The constant is the spec. Configuration tests are valuable only when there's logic — env-var precedence, validation, derived values.

### 8. Pure presentational components

Components (UI nodes, view templates, render-only widgets) with no state, no hooks beyond context-as-data, and no conditional rendering driven by runtime values are not unit-testable in any useful sense. They render markup. The type system enforces props; a parent's behavior test exercises rendering. Don't unit-test them in isolation.

Symptoms: source is < 50 lines, no internal state, no branching beyond a className ternary. Test renders the component and asserts that prop values appear in the output. Both the render and the assertion are restating the JSX/template.

### 9. DOM structure, CSS classes, data-attributes

Tests that assert specific CSS class names (`expect(el).toHaveClass("bg-secondary/10")`), `data-*` attributes (`toHaveAttribute("data-step", "2")`), or rigid DOM nesting pin the implementation, not the behavior. A purely visual refactor breaks the test without changing what the user experiences.

For UI tests, prefer accessibility-tree accessors — `getByRole`, `getByLabelText`, `getByText`, `findByDisplayValue`. They document user intent, survive style refactors, and double as a sanity check on accessibility. If you find yourself reaching for `querySelector` or `getByTestId` for a styling element, ask whether `getByRole` would work and refactor the component if not.

---

## Universal principles

### Mocks for boundaries you don't own; fakes for code you do

- **External services** (payment APIs, LLM providers, third-party SDKs): mocks are appropriate — you can't run their code in tests, and contract tests cover the boundary.
- **Your own code** (your repositories, your services, your event bus, your blob store): build fakes. A fake is a working in-memory implementation of the same interface. Mocks of your own code lock tests to implementation details and rot into call-sequence verification.

### Mock at the foreign boundary, not at your dependency boundary

When testing code that wraps a library (a hook around a query/cache library, a store around a state library, a service around an SDK), the right mock point is the *network or SDK boundary the data ultimately crosses* — not the library wrapper you happened to use.

Mocking the wrapper (`useQuery` to return `{data: X, isLoading: false}`, the store's `getState` to return `Y`) turns the test into "did I read what the mock said." Mocking the foreign boundary (`fetch`, the SDK client, the storage backend) and rendering with the real wrapper lets the library run for real, so the test exercises the code you actually wrote on top of it — caching behavior, derived state, error mapping, retries.

If a test passes when you replace the wrapper-mock's return value with constants, the test was reading the mock, not exercising the code.

### The type system is part of your coverage

In a statically-typed codebase, the type checker fails before any test runs. A test that would only fail when the type checker would also fail adds no signal — it just costs maintenance.

This excludes a large category of tests in typed languages:
- Prop / argument-type mismatches (the compiler catches them)
- Required-field omissions on schemas/dataclasses
- Function signature changes
- Return-shape changes
- Lookup-table key presence (when keys are typed as enums or string literals)
- "Did the data flow through" tests where the data shape is fully typed

Test what types can't enforce: runtime values, branching logic, async flows, side effects, error paths, business rules. Expect a noticeably lower test-to-source LOC ratio in typed codebases (e.g. TypeScript, Rust, Go, Kotlin) than in dynamic ones — what looks like under-testing is the type system absorbing tests you'd otherwise need to write.

### Assert outcomes, not call sequences

A test that passes because a mock was called proves nothing about correctness. A test that passes because state changed, an event was emitted, or an error was raised proves something a caller cares about. Always prefer the latter.

### A test must be able to fail for a reason a caller's test wouldn't

If the test would only fail when a caller's test would also fail, the test is redundant. Delete it and let the caller's test cover both.

### Coverage of trivial branches is not coverage

Branch / line coverage measures how much code ran during tests, not how much *behavior* was verified. 100% coverage of a tautology suite is worse than 70% coverage of an outcome suite. Track coverage of branches with failure modes, not all branches.

### Deletion is a valid fix

When a failing test pins an implementation detail that should never have been pinned, the correct response is to delete the test, not to update it to match the new implementation. Updating preserves the noise; deleting restores the signal-to-noise ratio. See `failure-investigation.md` (NOISE_TEST category).

---

## Decision flow

When writing or auditing a test, ask in order:

1. **What can this test fail to catch?** If "nothing a caller's test wouldn't also catch" — don't write it / delete it.
2. **What does this test pin?** If "the implementation, not the behavior" — rewrite to assert outcomes, or delete.
3. **What does it mock?** If "code I own" — build a fake instead.
4. **Where does this code live in the cost-of-miss spectrum?** Security / domain logic / coordinators → test it. Trivial wrapper / framework re-verification → don't.
5. **Is integration coverage already exercising this path?** If yes, narrow the unit test to what integration can't reach, or delete.

Apply these rules opinionatedly. A green suite of noise tests is a regression in signal-to-noise, not a deliverable.
