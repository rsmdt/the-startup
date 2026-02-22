# Optimization Patterns

Proven patterns for improving performance, organized by effort level.

## Quick Wins

1. **Enable caching** - Application, CDN, database query cache
2. **Add indexes** - For slow queries identified in profiling
3. **Compression** - Gzip/Brotli for responses
4. **Connection pooling** - Reduce connection overhead
5. **Batch operations** - Reduce round-trips

## Algorithmic Improvements

1. **Reduce complexity** - O(n^2) to O(n log n)
2. **Lazy evaluation** - Defer work until needed
3. **Memoization** - Cache computed results
4. **Pagination** - Limit data processed at once

## Architectural Changes

1. **Horizontal scaling** - Add more instances
2. **Async processing** - Queue background work
3. **Read replicas** - Distribute read load
4. **Caching layers** - Redis, Memcached
5. **CDN** - Edge caching for static content

## Capacity Planning

### Baseline Establishment

Measure current capacity under production load:

1. **Peak load metrics** - Maximum concurrent users, requests/sec
2. **Resource headroom** - How close to limits at peak
3. **Scaling patterns** - Linear, sub-linear, or super-linear

### Load Testing Approach

1. **Establish baseline** - Current performance at normal load
2. **Ramp testing** - Gradually increase load to find limits
3. **Stress testing** - Push beyond limits to understand failure modes
4. **Soak testing** - Sustained load to find memory leaks, degradation

### Capacity Metrics

| Metric | What It Tells You |
|--------|-------------------|
| Throughput at saturation | Maximum system capacity |
| Latency at 80% load | Performance before degradation |
| Error rate under stress | Failure patterns |
| Recovery time | How quickly system returns to normal |

### Growth Planning

Required Capacity = (Current Load x Growth Factor) + Safety Margin

Example:
- Current: 1000 req/sec
- Expected growth: 50% per year
- Safety margin: 30%
- Year 1 need = (1000 x 1.5) x 1.3 = 1950 req/sec
