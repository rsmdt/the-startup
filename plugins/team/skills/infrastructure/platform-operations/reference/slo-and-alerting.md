# Reference: SLO Framework and Alerting Strategies

Detailed SLI/SLO/SLA framework, error budget calculations, alerting patterns, and dashboard design.

---

## SLI/SLO/SLA Framework

### Service Level Indicators (SLIs)

Quantitative measures of service behavior from the user perspective.

**Common SLI categories:**

| Category | Measures | Example SLI |
|----------|----------|-------------|
| Availability | Service is responding | % of successful requests |
| Latency | Response speed | % of requests < 200ms |
| Throughput | Capacity | Requests processed per second |
| Error Rate | Correctness | % of requests without errors |
| Freshness | Data currency | % of data < 1 minute old |

**SLI specification:**
```
SLI: Request Latency
Definition: Time from request received to response sent
Measurement: Server-side histogram at p50, p95, p99
Exclusions: Health checks, internal tooling
Data source: Application metrics
```

### Service Level Objectives (SLOs)

Target reliability levels for SLIs over a time window.

**SLO formula:**
```
SLO = (Good events / Total events) >= Target over Window

Example:
99.9% of requests complete successfully in < 200ms
measured over a 30-day rolling window
```

**Setting SLO targets:**
- Start with current baseline performance
- Consider user expectations and business impact
- Balance reliability investment against feature velocity
- Document error budget policy

### Error Budgets

The allowed amount of unreliability within an SLO.

**Calculation:**
```
Error Budget = 1 - SLO Target

99.9% SLO = 0.1% error budget
= 43.2 minutes downtime per 30 days
= 8.64 seconds per day
```

**Error budget policies:**
- Budget remaining: Continue feature development
- Budget depleted: Focus on reliability work
- Budget burning fast: Freeze deploys, investigate

---

## Alerting Strategies

### Symptom-Based Alerts

Alert on user-visible symptoms, not internal causes.

**Good alerts:**
- Error rate exceeds threshold (users experiencing failures)
- Latency SLO at risk (users experiencing slowness)
- Queue depth growing (backlog affecting users)

**Poor alerts:**
- CPU at 80% (may not affect users)
- Pod restarted (self-healing, may not affect users)
- Disk at 70% (not yet impacting service)

### Multi-Window, Multi-Burn-Rate Alerts

Detect fast burns quickly, slow burns before budget depletion.

**Configuration:**
```
Fast burn: 14.4x burn rate over 1 hour
  - Fires in 1 hour if issue persists
  - Catches severe incidents quickly

Slow burn: 3x burn rate over 3 days
  - Fires before 30-day budget depletes
  - Catches gradual degradation
```

### Alert Fatigue Prevention

**Strategies:**
- Alert only on actionable issues
- Consolidate related alerts
- Set meaningful thresholds (not arbitrary)
- Require sustained condition before firing
- Include runbook links in every alert
- Review and prune alerts quarterly

**Alert quality checklist:**
- Can someone take action right now?
- Is the severity appropriate?
- Does it include enough context?
- Is there a runbook linked?
- Has it fired false positives recently?

---

## Dashboard Design

### Hierarchy of Dashboards

**Service Health Overview:**
- High-level SLO status
- Error budget consumption
- Key business metrics
- Designed for quick triage

**Deep-Dive Diagnostic:**
- Detailed metrics breakdown
- Resource utilization
- Dependency health
- Designed for investigation

**Business Metrics:**
- User-facing KPIs
- Conversion and engagement
- Revenue impact
- Designed for stakeholders

### Dashboard Principles

- Answer specific questions, not show all data
- Use consistent color coding (green=good, red=bad)
- Show time ranges appropriate to the metric
- Include context (deployments, incidents) on graphs
- Mobile-responsive for on-call use
- Provide drill-down paths to detailed views

### Essential Panels

| Panel | Purpose | Audience |
|-------|---------|----------|
| SLO Status | Current reliability vs target | Everyone |
| Error Budget | Remaining budget and burn rate | Engineering |
| Request Rate | Traffic patterns and anomalies | Operations |
| Latency Distribution | p50, p95, p99 over time | Engineering |
| Error Breakdown | Errors by type and endpoint | Engineering |
| Dependency Health | Status of upstream services | Operations |
