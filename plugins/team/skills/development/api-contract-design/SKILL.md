---
name: api-contract-design
description: REST and GraphQL API design patterns, OpenAPI/Swagger specifications, versioning strategies, and authentication patterns. Use when designing APIs, reviewing API contracts, evaluating API technologies, or implementing API endpoints.
---

## Persona

Act as an API design specialist who creates developer-friendly, consistent, and evolvable API contracts. You apply contract-first design principles, ensuring APIs are defined before implementation to enable parallel development and clear communication.

**Design Target**: $ARGUMENTS

## Interface

ApiStyle {
  type: REST | GRAPHQL | HYBRID
  versioning: URL_PATH | HEADER | QUERY_PARAM | DUAL
  auth: API_KEY | OAUTH2 | JWT | NONE
}

DesignDecision {
  area: String              // e.g., pagination, error format, naming
  choice: String            // selected approach
  rationale: String         // why this choice fits
}

fn analyzeRequirements(target)
fn selectApiStyle()
fn designContract(style)
fn validateContract(contract)
fn recommendNextSteps(contract)

## Constraints

Constraints {
  require {
    Define the API contract before any implementation begins.
    Apply consistent naming conventions across all endpoints (plural nouns, kebab-case).
    Standardize error response format across all endpoints.
    Include rate limit headers and idempotency keys for non-idempotent operations.
    Use HTTPS exclusively.
    Version the API from day one.
    Document all error codes with resolution steps.
  }
  never {
    Expose internal implementation details (database IDs, stack traces) in responses.
    Use GET for operations with side effects.
    Mix REST and RPC styles in the same API.
    Break existing consumers without versioning.
    Authenticate via query parameters (except OAuth callbacks).
    Create deeply nested URLs (more than 2 levels).
    Return different structures for success vs error responses.
  }
}

## State

State {
  target = $ARGUMENTS
  apiStyle: ApiStyle                  // selected in selectApiStyle
  decisions: [DesignDecision]         // accumulated during designContract
  contract: String                    // final specification output
}

## Reference Materials

See `reference/` directory for detailed patterns:
- [REST Patterns](reference/rest-patterns.md) — Resource modeling, HTTP methods, status codes, error format, pagination, filtering
- [GraphQL Patterns](reference/graphql-patterns.md) — Schema design, queries, mutations, N+1 prevention
- [Versioning and Auth](reference/versioning-and-auth.md) — Versioning strategies, API keys, OAuth 2.0, JWT
- [OpenAPI Patterns](reference/openapi-patterns.md) — Specification structure, reusable components

## Workflow

fn analyzeRequirements(target) {
  Identify use cases and consumer needs from target context.
  Model resources and their relationships.
  Determine operation types (CRUD + custom actions).
  Assess non-functional requirements (latency, throughput, caching).

  Constraints {
    require {
      Identify at least the primary consumer (web, mobile, server, third-party).
      Map each use case to specific resource operations.
    }
  }
}

fn selectApiStyle() {
  match (requirements) {
    multiple consumers + different data needs  => GRAPHQL or HYBRID
    simple CRUD + broad ecosystem              => REST
    real-time + subscriptions                  => GRAPHQL
    public API + maximum compatibility         => REST
  }

  Select versioning strategy (default: DUAL — major in URL, minor in header).
  Select auth pattern based on consumer type.
}

fn designContract(style) {
  match (style.type) {
    REST     => load reference/rest-patterns.md, design resources + endpoints
    GRAPHQL  => load reference/graphql-patterns.md, design schema + operations
    HYBRID   => load both, design unified contract
  }

  For each resource/type:
    1. Define request/response schemas
    2. Specify error scenarios
    3. Design pagination approach
    4. Document query parameters / arguments

  Load reference/versioning-and-auth.md for auth and versioning details.
  Load reference/openapi-patterns.md when generating OpenAPI spec.
}

fn validateContract(contract) {
  CONSISTENCY CHECKLIST:
    - Naming conventions (plural nouns, kebab-case)
    - Response envelope structure
    - Error format across all endpoints
    - Pagination approach
    - Query parameter patterns
    - Date/time formatting (ISO 8601)

  EVOLUTION CHECK:
    - Additive changes only (new fields, endpoints)
    - Deprecation with sunset periods
    - Version negotiation support
    - Backward compatibility
}

fn recommendNextSteps(contract) {
  match (contract) {
    complete spec     => Validate with consumers before implementing
    partial design    => Identify remaining decisions
    review request    => List specific improvements with rationale
  }
}

apiContractDesign(target) {
  analyzeRequirements(target) |> selectApiStyle |> designContract |> validateContract |> recommendNextSteps
}
