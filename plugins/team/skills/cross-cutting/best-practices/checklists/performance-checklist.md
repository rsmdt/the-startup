# Performance Checklist

Actionable performance optimization checks for frontend, backend, and database operations.

## Frontend Performance

### Initial Load

- [ ] Critical CSS is inlined or loaded first
- [ ] JavaScript is deferred or loaded async where possible
- [ ] Images are lazy-loaded below the fold
- [ ] Web fonts use `font-display: swap` to prevent blocking
- [ ] Bundle size is monitored and kept under budget
- [ ] Code splitting is used for route-based chunks
- [ ] Tree shaking removes unused code

### Runtime Performance

- [ ] Long lists use virtualization (render only visible items)
- [ ] Expensive calculations are memoized
- [ ] React/Vue components avoid unnecessary re-renders
- [ ] Event handlers are debounced or throttled where appropriate
- [ ] Animations use CSS transforms, not layout-triggering properties
- [ ] Large data processing uses Web Workers

### Assets

- [ ] Images are properly sized (not scaling down large images)
- [ ] Modern image formats used (WebP, AVIF with fallbacks)
- [ ] Images are compressed without quality loss
- [ ] SVGs are optimized and used for icons
- [ ] Fonts are subsetted to include only needed characters
- [ ] Assets are served from CDN

### Caching

- [ ] Static assets have cache-control headers with long max-age
- [ ] Asset filenames include content hash for cache busting
- [ ] Service worker caches critical resources
- [ ] API responses include appropriate cache headers

## Backend Performance

### Request Handling

- [ ] Endpoints respond within acceptable latency targets
- [ ] Heavy operations are processed asynchronously
- [ ] Request payload size limits are enforced
- [ ] Response payloads include only necessary data
- [ ] Pagination is used for list endpoints
- [ ] Compression (gzip, brotli) is enabled

### Computation

- [ ] CPU-intensive tasks are offloaded to worker processes
- [ ] Algorithms have acceptable time complexity for data size
- [ ] Loops avoid unnecessary iterations
- [ ] String concatenation in loops uses efficient methods
- [ ] Regular expressions are optimized and avoid catastrophic backtracking

### Caching

- [ ] Frequently accessed data is cached in memory (Redis, Memcached)
- [ ] Cache invalidation strategy is defined
- [ ] Cache hit rates are monitored
- [ ] Cache keys are designed to avoid collisions
- [ ] TTLs are set appropriately for data freshness requirements

### External Services

- [ ] API calls are batched where possible
- [ ] Circuit breakers prevent cascading failures
- [ ] Timeouts are set on all external calls
- [ ] Retries use exponential backoff
- [ ] Connection pools are sized appropriately

## Database Performance

### Query Optimization

- [ ] Queries retrieve only needed columns (no SELECT *)
- [ ] Indexes exist for frequently queried columns
- [ ] Composite indexes match query patterns
- [ ] N+1 query problems are eliminated (use JOINs or batch loading)
- [ ] Query execution plans are reviewed for slow queries
- [ ] LIMIT is used for large result sets

### Schema Design

- [ ] Tables are normalized appropriately (avoid over-normalization)
- [ ] Data types are sized correctly (not oversized)
- [ ] Foreign keys have indexes
- [ ] Frequently accessed fields are not in separate tables requiring JOINs

### Connection Management

- [ ] Connection pooling is configured
- [ ] Pool size matches expected concurrency
- [ ] Connections are released promptly after use
- [ ] Long-running transactions are avoided
- [ ] Read replicas are used for read-heavy workloads

### Data Lifecycle

- [ ] Old data is archived or deleted according to policy
- [ ] Large tables have partitioning strategy
- [ ] Bulk operations are batched to avoid lock contention
- [ ] Indexes are rebuilt periodically if needed

## API Performance

### Design

- [ ] Endpoints support field selection (sparse fieldsets)
- [ ] Batch endpoints available for multiple-item operations
- [ ] GraphQL queries have depth/complexity limits
- [ ] Large responses support streaming or pagination

### Network

- [ ] HTTP/2 or HTTP/3 is enabled
- [ ] Keep-alive connections are used
- [ ] Response compression is enabled
- [ ] CDN caching is used for cacheable responses

## Monitoring and Measurement

### Metrics

- [ ] Response time percentiles are tracked (p50, p95, p99)
- [ ] Error rates are monitored
- [ ] Resource utilization is tracked (CPU, memory, connections)
- [ ] Throughput is measured under load

### Alerting

- [ ] Latency degradation triggers alerts
- [ ] Resource exhaustion triggers alerts
- [ ] Error rate spikes trigger alerts

### Profiling

- [ ] Production profiling tools are available
- [ ] Slow query logging is enabled
- [ ] Application performance monitoring (APM) is configured

## Usage Notes

1. Measure before optimizing - identify actual bottlenecks
2. Set performance budgets and enforce them in CI
3. Not every optimization applies to every system
4. Document performance-critical paths
5. Load test under realistic conditions
6. Performance regression tests should be part of CI/CD
