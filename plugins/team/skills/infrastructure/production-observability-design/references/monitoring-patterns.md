# Reference: Monitoring Patterns

## Overview

Implementation patterns for metrics collection, distributed tracing, alerting configuration, and SLO management across different technology stacks and monitoring platforms.

## Quick Reference

| Pattern | Use Case | Complexity |
|---------|----------|------------|
| RED Method | Request-driven services | Low |
| USE Method | Resource monitoring | Low |
| Four Golden Signals | General service health | Low |
| Multi-window alerting | SLO-based alerting | Medium |
| Trace context propagation | Distributed tracing | Medium |
| Exemplars | Metric-to-trace correlation | Medium |
| Synthetic monitoring | Proactive availability | Medium |

## Metric Collection Patterns

### RED Method (Request-driven services)

Monitor three key metrics for every service:

```
Rate:     requests per second
Errors:   failed requests per second
Duration: distribution of request latency
```

**When to use:** Microservices, APIs, web applications

**Implementation:**
```
# Prometheus-style metrics
http_requests_total{method, endpoint, status}
http_request_duration_seconds{method, endpoint}
http_request_errors_total{method, endpoint, error_type}
```

### USE Method (Resources)

Monitor three aspects for every resource:

```
Utilization: percentage of resource capacity in use
Saturation:  degree to which resource has queued work
Errors:      count of error events
```

**When to use:** CPU, memory, disk, network, database connections

**Example metrics:**
```
CPU:
  Utilization: cpu_usage_percent
  Saturation:  load_average, runnable_processes
  Errors:      hardware_errors (rare)

Memory:
  Utilization: memory_used_percent
  Saturation:  swap_usage, oom_kills
  Errors:      memory_errors (rare)

Disk:
  Utilization: disk_used_percent
  Saturation:  io_queue_depth, io_wait
  Errors:      disk_errors
```

### Four Golden Signals (Google SRE)

```
Latency:   time to service a request
Traffic:   demand on the system
Errors:    rate of failed requests
Saturation: how full the service is
```

**Relationship to RED/USE:**
- Latency = Duration from RED
- Traffic = Rate from RED
- Errors = Errors from both methods
- Saturation = Saturation from USE

## Distributed Tracing Patterns

### Context Propagation

Pass trace context through all service calls:

```
Trace Context Headers:
- traceparent: version-trace_id-parent_id-flags
- tracestate: vendor-specific data

Example:
traceparent: 00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01
tracestate: congo=t61rcWkgMzE
```

**Propagation checklist:**
- HTTP headers (W3C Trace Context standard)
- Message queue metadata
- Database query comments
- Background job context
- Cross-region calls

### Span Best Practices

```
Span naming:
- HTTP: METHOD /path (GET /users/{id})
- Database: db.operation table (SELECT users)
- Queue: queue.operation name (PUBLISH orders)

Required attributes:
- service.name
- operation.name
- span.kind (server, client, producer, consumer)

Optional context:
- user.id (sanitized)
- request.id
- error.message (on failure)
```

### Sampling Strategies

| Strategy | Description | Use Case |
|----------|-------------|----------|
| Head-based | Decide at trace start | High throughput |
| Tail-based | Decide after trace complete | Error investigation |
| Rate-limited | Fixed traces per second | Cost control |
| Priority | Sample important paths more | Critical flows |

**Tail-based sampling rules:**
```
Always sample if:
- Error occurred in any span
- Latency exceeded threshold
- Important user or endpoint
- Random selection (baseline)
```

## SLO Implementation Patterns

### SLO Configuration

```
Service: payment-api
SLO: 99.9% availability over 30 days

SLI specification:
  Numerator: successful requests (2xx, 3xx)
  Denominator: all requests (excluding health checks)
  Measurement: server-side metrics

Alert thresholds:
  - 1h window, 2% budget burn: page immediately
  - 6h window, 5% budget burn: page immediately
  - 3d window, 10% budget burn: create ticket
```

### Multi-Window, Multi-Burn-Rate Alert Config

```
# Fast burn alert (catches acute issues)
alert: SLOFastBurn
expr: |
  (
    sum(rate(http_request_errors_total[1h]))
    /
    sum(rate(http_requests_total[1h]))
  ) > (14.4 * 0.001)
  AND
  (
    sum(rate(http_request_errors_total[5m]))
    /
    sum(rate(http_requests_total[5m]))
  ) > (14.4 * 0.001)
severity: critical

# Slow burn alert (catches gradual degradation)
alert: SLOSlowBurn
expr: |
  (
    sum(rate(http_request_errors_total[6h]))
    /
    sum(rate(http_requests_total[6h]))
  ) > (6 * 0.001)
  AND
  (
    sum(rate(http_request_errors_total[3d]))
    /
    sum(rate(http_requests_total[3d]))
  ) > (1 * 0.001)
severity: warning
```

### Error Budget Calculation

```
Monthly error budget for 99.9% SLO:
  Budget = (1 - 0.999) * 30 days * 24 hours * 60 minutes
  Budget = 0.001 * 43200 minutes = 43.2 minutes

Current consumption:
  Consumed = (bad_minutes / total_minutes) * 100
  Remaining = Budget - Consumed

Burn rate:
  Burn rate = (actual_error_rate / error_budget_rate)
  If burn rate > 1: consuming budget faster than allowed
```

## Alerting Patterns

### Alert Template

```
Alert: ServiceErrorRateHigh
Summary: {{ $labels.service }} error rate above threshold
Description: |
  Service {{ $labels.service }} is experiencing {{ $value }}% error rate
  over the last 5 minutes. Normal baseline is < 0.1%.

  Impact: Users may be experiencing failures on {{ $labels.endpoint }}

  Runbook: https://runbooks.example.com/service-errors
  Dashboard: https://grafana.example.com/d/service-overview

Labels:
  severity: critical
  team: platform

Annotations:
  dashboard: https://grafana.example.com/d/service-overview
  runbook: https://runbooks.example.com/service-errors
```

### Alert Routing

```
Routes by severity:
  critical:
    - PagerDuty (wake people up)
    - Slack #incidents (visibility)
  warning:
    - Slack #alerts (triage queue)
    - Create ticket (tracking)
  info:
    - Slack #monitoring (awareness)
    - No action required

Routes by team:
  team=platform -> platform-oncall
  team=payments -> payments-oncall
  team=frontend -> frontend-oncall
```

### Silence and Maintenance Windows

```
Planned maintenance:
  - Schedule silence before maintenance
  - Include reason and expected duration
  - Notify affected teams
  - Remove silence after completion

Incident response:
  - Silence noisy alerts during investigation
  - Document reason for silence
  - Set expiration (max 4 hours)
  - Review silences daily
```

## Log Patterns

### Structured Log Format

```json
{
  "timestamp": "2024-01-15T10:30:00.123Z",
  "level": "ERROR",
  "service": "payment-api",
  "version": "1.2.3",
  "environment": "production",
  "trace_id": "4bf92f3577b34da6a3ce929d0e0e4736",
  "span_id": "00f067aa0ba902b7",
  "message": "Payment processing failed",
  "error": {
    "type": "PaymentDeclined",
    "message": "Insufficient funds",
    "code": "INSUFFICIENT_FUNDS"
  },
  "context": {
    "user_id": "usr_abc123",
    "request_id": "req_xyz789",
    "amount_cents": 5000,
    "currency": "USD"
  },
  "duration_ms": 234
}
```

### Log Level Guidelines

| Level | When to Use | Alertable |
|-------|-------------|-----------|
| FATAL | Application cannot continue | Yes |
| ERROR | Operation failed, user impacted | Yes |
| WARN | Degraded but functional | Monitor rate |
| INFO | Significant business events | No |
| DEBUG | Detailed troubleshooting | No |

### Log Retention Strategy

```
Retention tiers:
  Hot (searchable):     7 days
  Warm (queryable):     30 days
  Cold (archived):      90 days
  Deleted:              after 90 days

Exception:
  Security logs:        1 year minimum
  Compliance logs:      per regulatory requirement
  Error logs:           30 days hot
```

## Dashboard Patterns

### Service Overview Dashboard

```
Row 1: Key SLOs
  - Availability SLO gauge (current vs target)
  - Error budget remaining gauge
  - SLO trend over time

Row 2: Traffic and Errors
  - Request rate graph
  - Error rate graph (by type)
  - Success/failure ratio

Row 3: Latency
  - p50, p95, p99 latency over time
  - Latency distribution histogram
  - Slow endpoints table

Row 4: Dependencies
  - Downstream service health
  - Database connection pool
  - Cache hit rate
```

### Incident Investigation Dashboard

```
Row 1: Timeline
  - Events overlay (deploys, alerts, incidents)
  - Error rate spike correlation
  - Change log entries

Row 2: Deep Dive
  - Error breakdown by type
  - Affected endpoints
  - Affected users (count, not identifiable)

Row 3: Resources
  - CPU, memory, disk usage
  - Pod/instance health
  - Network I/O

Row 4: Traces
  - Slow trace examples
  - Error trace examples
  - Trace search link
```

## Synthetic Monitoring Patterns

### Availability Checks

```
Check types:
  - HTTP endpoint availability
  - SSL certificate expiration
  - DNS resolution
  - TCP port connectivity

Configuration:
  Frequency: 1 minute
  Locations: multiple regions
  Timeout: 10 seconds
  Retries: 1 before alerting
```

### User Journey Checks

```
Critical paths to monitor:
  - Login flow
  - Checkout process
  - Core API operations
  - Key user workflows

Implementation:
  - Headless browser for UI
  - API calls for backend
  - Assert on response content
  - Measure each step timing
```

## Incident Response Patterns

### Severity Levels

| Level | Impact | Response | Example |
|-------|--------|----------|---------|
| SEV1 | Complete outage | All hands, exec notification | Site down |
| SEV2 | Significant degradation | On-call + backup | 50% errors |
| SEV3 | Minor impact | On-call investigation | Slow responses |
| SEV4 | No user impact | Next business day | Warning alerts |

### Postmortem Template

```
Incident Summary:
  - Duration: Start to resolution
  - Impact: Users/requests affected
  - Root cause: One sentence

Timeline:
  - Detection time and method
  - Key investigation steps
  - Resolution actions

Root Cause Analysis:
  - What failed and why
  - Contributing factors
  - Why detection took X minutes

Action Items:
  - Prevent recurrence (P0)
  - Improve detection (P1)
  - Improve recovery (P2)

Lessons Learned:
  - What went well
  - What could improve
  - Follow-up items
```

## External Resources

- [Google SRE Book - Monitoring](https://sre.google/sre-book/monitoring-distributed-systems/)
- [OpenTelemetry Specification](https://opentelemetry.io/docs/specs/)
- [Prometheus Best Practices](https://prometheus.io/docs/practices/)
- [Datadog Monitoring Guide](https://www.datadoghq.com/blog/monitoring-101-collecting-data/)
