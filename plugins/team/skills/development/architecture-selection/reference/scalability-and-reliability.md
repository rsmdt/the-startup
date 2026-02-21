# Scalability and Reliability Patterns

Horizontal scaling, caching, database scaling, and reliability patterns.

---

## Horizontal Scaling

Add more instances of the same component.

```
                    Load Balancer
                         │
         ┌───────────────┼───────────────┐
         │               │               │
         ▼               ▼               ▼
    ┌─────────┐     ┌─────────┐     ┌─────────┐
    │Instance │     │Instance │     │Instance │
    │    1    │     │    2    │     │    3    │
    └─────────┘     └─────────┘     └─────────┘

Requirements:
- Stateless services
- Shared session storage
- Database can handle connections
```

## Caching

Reduce load on slow resources.

```
┌─────────────────────────────────────────────────────┐
│                  Caching Layers                      │
├─────────────────────────────────────────────────────┤
│                                                     │
│  Browser Cache → CDN → App Cache → Database Cache  │
│                                                     │
│  Examples:                                          │
│  - Browser: Static assets, API responses           │
│  - CDN: Static content, cached API responses       │
│  - App: Redis/Memcached for sessions, computed data│
│  - Database: Query cache, connection pooling       │
│                                                     │
└─────────────────────────────────────────────────────┘

Cache Invalidation Strategies:
- TTL (Time to Live): Simplest, eventual consistency
- Write-through: Update cache on write
- Write-behind: Async update for performance
- Cache-aside: App manages cache explicitly
```

## Database Scaling

| Strategy | Use Case | Trade-off |
|----------|----------|-----------|
| **Read Replicas** | Read-heavy workloads | Replication lag |
| **Sharding** | Large datasets | Query complexity |
| **Partitioning** | Time-series data | Partition management |
| **CQRS** | Different read/write patterns | System complexity |

## Reliability Patterns

| Pattern | Purpose | Implementation |
|---------|---------|----------------|
| **Circuit Breaker** | Prevent cascade failures | Fail fast after threshold |
| **Bulkhead** | Isolate failures | Separate thread pools |
| **Retry** | Handle transient failures | Exponential backoff |
| **Timeout** | Bound wait times | Don't wait forever |
| **Rate Limiting** | Prevent overload | Throttle requests |

```
Circuit Breaker States:

    ┌────────┐
    │ CLOSED │ ──── Failure Threshold ──► ┌────────┐
    │(normal)│                            │  OPEN  │
    └────────┘                            │(failing│
         ▲                                └────┬───┘
         │                                     │
    Success                              Timeout
    Threshold                                  │
         │                                     ▼
    ┌────┴────┐                          ┌─────────┐
    │HALF-OPEN│ ◄─── Test Request ────── │         │
    └─────────┘                          └─────────┘
```
