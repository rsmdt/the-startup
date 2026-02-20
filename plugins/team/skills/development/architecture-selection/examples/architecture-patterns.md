# Architecture Patterns: Decision Guide

Practical catalog of architectural patterns with trade-offs. Use this alongside the pattern selection table in SKILL.md. Each entry answers the two questions that matter most: when to reach for this pattern, and when to avoid it.

---

## Monolith

A single deployable unit. All modules share the same process, memory space, and database.

**Reach for this when:**
- Team is under 10 developers sharing the same codebase
- Domain is not yet well understood — premature decomposition locks in wrong boundaries
- Time to first working product is the dominant constraint
- Infrastructure expertise is limited or unavailable
- Transactions must be ACID across what would otherwise be service boundaries

**Do not use when:**
- Multiple teams need to deploy independently without coordinating releases
- Different parts of the system have wildly different scaling profiles (e.g., file uploads vs. real-time feeds)
- You need polyglot persistence (each domain using its most appropriate database)
- The codebase is already too large to hold in one developer's head

**Scaling characteristics:**
- Scale by cloning the entire application behind a load balancer
- All modules scale together regardless of where the bottleneck actually is
- Database becomes the bottleneck first; address with read replicas and caching before reaching for microservices
- Practical ceiling: tens of millions of requests per day with a well-tuned monolith on modern hardware

**Team size fit:** 1–10 developers

---

## Modular Monolith

A monolith with enforced internal module boundaries. Modules own their data and expose explicit interfaces to each other. No shared tables. No calling into another module's internals.

**Reach for this when:**
- You want the operational simplicity of a monolith but are building toward future service extraction
- Domain boundaries are becoming clearer but you are not ready to pay the operational cost of microservices
- Team is growing (10–20 developers) and you need to reduce coordination overhead without splitting deployment
- You have been burned by a big ball of mud monolith and want structure without the distributed systems tax

**Do not use when:**
- Teams genuinely need independent deployment pipelines today
- Module isolation is not enforced by tooling — without enforcement, boundaries erode under deadline pressure
- Scaling requirements are already differentiated between modules

**Scaling characteristics:**
- Same ceiling as a monolith: scale by replication
- Module isolation makes it significantly easier to extract a service later when a specific module becomes the bottleneck
- The database is still shared at the infrastructure level even if modules have schema ownership

**Team size fit:** 5–25 developers

---

## Microservices

Independent services each owning a bounded context, deployed and scaled separately.

**Reach for this when:**
- Multiple autonomous teams need to ship independently without release coordination
- Scaling needs are genuinely differentiated: the checkout service needs 50x more capacity than the admin dashboard
- Different services have legitimately different technology requirements (streaming processing vs. CRUD vs. ML inference)
- The domain is well understood and bounded contexts are stable — wrong boundaries are expensive to fix later
- High availability is a hard requirement: a failure in recommendations must not take down checkout

**Do not use when:**
- Team is under 20 developers — you will spend more time on distributed systems plumbing than on product
- Domain boundaries are still being discovered — wait until the seams are clear
- The team lacks experience with distributed systems, container orchestration, service meshes, and observability
- Transactions must be strongly consistent across what would be service boundaries — distributed sagas are a significant complexity investment
- You are starting a new product — the monolith-first rule exists for good reason

**Scaling characteristics:**
- Scale individual services independently based on measured demand
- Each service can adopt the database technology that fits its access patterns
- Horizontal scaling at the service level; the API gateway and message bus become new bottlenecks to manage
- Requires investment in service discovery, load balancing, distributed tracing, and centralized logging

**Team size fit:** 20+ developers, organized into product teams aligned to service ownership

---

## Event-Driven Architecture

Services communicate by publishing and subscribing to events on a message broker. No direct synchronous calls between services for primary workflows.

**Reach for this when:**
- Workflows span multiple services and you cannot afford to couple their availability
- You need a durable audit trail of everything that has happened in the system
- Processing is inherently asynchronous: order fulfillment, email delivery, report generation
- You need to replay historical events to rebuild state or feed new downstream consumers
- Fan-out is natural: one event (OrderPlaced) needs to trigger actions in multiple independent systems

**Do not use when:**
- The user needs a synchronous response to complete their action — eventual consistency is a user experience trade-off, not just a technical one
- The team has no experience operating a message broker under production load
- The domain has complex ordering requirements that are difficult to enforce across independent consumers
- Debugging and distributed tracing tooling is not in place — event-driven systems are significantly harder to debug without it

**Scaling characteristics:**
- Producers and consumers scale independently
- The message broker (Kafka, RabbitMQ, SQS) becomes the critical scaling and reliability component
- Consumer groups allow parallel processing of event streams
- Backpressure is managed through queue depth monitoring and consumer scaling policies

**Team size fit:** Any size, but requires operational maturity. Do not introduce event-driven patterns without investing in observability first.

---

## CQRS (Command Query Responsibility Segregation)

Separate models for write operations (commands) and read operations (queries). Often paired with event sourcing but does not require it.

**Reach for this when:**
- Read and write access patterns are fundamentally different: complex reporting queries on the same data used for transactional writes
- The read model needs to be denormalized for performance but the write model needs normalization for consistency
- Read throughput is orders of magnitude higher than write throughput
- You need multiple specialized read models from the same underlying data (e.g., search index, reporting database, and API response shape)

**Do not use when:**
- The domain is simple CRUD — CQRS adds complexity without benefit
- The team is not prepared to handle eventual consistency between the write and read sides
- You are in early product discovery — CQRS optimizes for a usage pattern you may not have confirmed yet
- Event sourcing is not in scope and the synchronization strategy between write and read stores has not been designed

**Scaling characteristics:**
- Read and write sides scale independently
- Read models can be replicated aggressively since they are derived, not authoritative
- Write side remains the consistency bottleneck; scale cautiously and measure before sharding
- Projection lag (time between a write and the read model reflecting it) must be measured and communicated to users where it matters

**Team size fit:** 5–50 developers, in systems where read/write asymmetry has been measured

---

## Hexagonal Architecture (Ports and Adapters)

The application core contains pure business logic with no framework or infrastructure dependencies. All external concerns (HTTP, databases, message queues, external APIs) connect through defined ports implemented by adapters.

**Reach for this when:**
- Business logic must be testable in isolation without spinning up databases or HTTP servers
- You need to swap infrastructure implementations: test database vs. production, email adapter vs. SMS adapter
- The domain is complex and protecting it from framework churn is a priority
- The team is practicing DDD and needs a clear boundary between domain and infrastructure

**Do not use when:**
- The application is primarily CRUD with little business logic — the overhead of ports and adapters is not justified
- The team is unfamiliar with the pattern — misapplied hexagonal architecture produces more abstractions than value
- Rapid prototyping is the goal — the pattern slows initial development in exchange for long-term maintainability

**Scaling characteristics:**
- The pattern is architectural, not a scaling strategy
- Infrastructure adapters can be swapped to use scaling-oriented implementations (e.g., swap in-memory queue adapter for SQS) without changing the domain
- Enables independent testing of scaling scenarios at the adapter level

**Team size fit:** Any size. Most valuable on teams of 5+ working on complex domains. Not worth the overhead for scripts or simple CRUD services.

---

## Serverless

Business logic runs as short-lived functions invoked by events or HTTP requests. The platform manages all server provisioning, scaling, and availability.

**Reach for this when:**
- Workload is irregular or spiky — pay for actual execution time, not idle capacity
- Operations are event-triggered: file uploaded, webhook received, scheduled job
- Functions are short-lived (under 15 minutes) and stateless
- Time to production is the dominant constraint and the team cannot afford infrastructure management
- Cost optimization for low-traffic workloads is a priority

**Do not use when:**
- Operations are long-running or require persistent in-memory state
- Cold start latency is unacceptable for user-facing requests (sub-100ms response time requirements)
- Local development and testing parity with production is critical — serverless local emulation is imperfect
- Vendor lock-in is a hard constraint — serverless functions are deeply coupled to platform-specific runtimes and event formats
- The system has complex inter-function orchestration — function chaining becomes a debugging and reliability problem

**Scaling characteristics:**
- Scales to zero when idle — no cost for unused capacity
- Scales to thousands of concurrent invocations automatically
- Concurrency limits are platform-enforced and can cause throttling at high traffic without reserved concurrency configuration
- The database connection problem: functions scale to thousands; databases have connection limits. Use connection poolers (RDS Proxy, PgBouncer) or serverless-native databases (DynamoDB, Aurora Serverless)

**Team size fit:** 1–20 developers, particularly effective for small teams that cannot staff platform/infrastructure roles

---

## Layered Architecture (N-Tier)

Code organized into horizontal layers: Presentation, Application/Business Logic, Data Access. Each layer only communicates with the layer directly below it.

**Reach for this when:**
- Building a traditional web application where the layered model is well understood by the team
- Enforcing separation between UI, business logic, and data access is the primary goal
- The domain is not complex enough to warrant DDD but you still want structure
- The team is early-career or coming from frameworks (Spring, ASP.NET, Django) that naturally express this structure

**Do not use when:**
- Business logic needs to be tested in isolation without the data access layer — strict layering often causes domain logic to be expressed in terms of database rows rather than business concepts
- The architecture needs to evolve toward hexagonal or DDD — layered architecture is harder to refactor than ports and adapters
- The "layer" boundaries become a ritual rather than a structural constraint — teams often skip layers or create anemic domain models where business logic leaks into controllers

**Scaling characteristics:**
- The same scaling constraints as a monolith apply
- Each tier can in principle be scaled independently (stateless presentation tier, application server pool, database)
- In practice, the data layer is the constraint and the pattern does not provide meaningful tools for addressing it

**Team size fit:** 1–15 developers, particularly teams starting with a framework-centric approach

---

## Pattern Comparison at a Glance

| Pattern | Operational Complexity | Time to First Feature | Long-Term Maintainability | Scaling Ceiling |
|---|---|---|---|---|
| Monolith | Low | Fast | Medium | Medium |
| Modular Monolith | Low | Medium | High | Medium |
| Microservices | High | Slow | High (with discipline) | High |
| Event-Driven | High | Medium | High (with tooling) | High |
| CQRS | Medium | Slow | High | High |
| Hexagonal | Low | Medium | High | N/A (architectural) |
| Serverless | Low | Fast | Medium | High |
| Layered | Low | Fast | Low-Medium | Medium |

---

## The Default Path

When in doubt, follow this sequence:

1. Start with a **modular monolith** using hexagonal architecture internally
2. Extract services only when a specific module has a scaling or deployment independence requirement that is measured, not anticipated
3. Introduce event-driven communication at service boundaries when synchronous coupling becomes a reliability risk
4. Add CQRS to specific modules when read/write asymmetry is measured and causing problems

Premature decomposition is the most common architecture mistake. The cost of splitting a well-structured monolith is low. The cost of merging a poorly-bounded microservices system is very high.
