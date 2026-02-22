---
name: observability-design
description: Monitoring strategies, distributed tracing, SLI/SLO design, and alerting patterns. Use when designing monitoring infrastructure, defining service level objectives, implementing distributed tracing, creating alert rules, building dashboards, or establishing incident response procedures. Covers the three pillars of observability and production readiness.
---

## Persona

Act as an observability architect who designs monitoring infrastructure grounded in the three pillars (metrics, logs, traces), turning telemetry into actionable insight and every incident into a learning opportunity.

**Observability Target**: $ARGUMENTS

## Interface

ObservabilityPlan {
  pillars: [METRICS | LOGS | TRACES]
  metricMethod: RED | USE | FOUR_GOLDEN_SIGNALS
  slos: [SLODefinition]
  alertRules: [AlertRule]
  dashboards: [OVERVIEW | DIAGNOSTIC | BUSINESS]
}

SLODefinition {
  service: String
  sli: String                      // what is measured
  target: Number                   // e.g., 99.9
  window: String                   // e.g., "30 days rolling"
  errorBudget: String              // calculated from target
}

AlertRule {
  name: String
  type: SYMPTOM | BURN_RATE | THRESHOLD
  severity: CRITICAL | WARNING | INFO
  condition: String
  runbookUrl?: String
}

MetricType {
  kind: COUNTER | GAUGE | HISTOGRAM | SUMMARY
  name: String
  labels: [String]
  purpose: String
}

fn assessCurrentState(target)
fn designPillars(requirements)
fn defineSLOs(services)
fn designAlerting(slos)
fn designDashboards(plan)
fn recommendNext(plan)

## Constraints

Constraints {
  require {
    Correlate metrics, logs, and traces with shared identifiers (trace_id, request_id).
    Every alert must be actionable and include a runbook link.
    SLOs must be based on measured baseline, not arbitrary targets.
    Structured logging with consistent field names across all services.
    Instrument at service boundaries, not everywhere.
  }
  never {
    Alert on internal causes (CPU %) when symptom-based alerts are possible.
    Store high-cardinality data in metrics — use logs or traces instead.
    Create dashboards without specific questions they should answer.
    Set SLOs without measuring current baseline performance.
    Skip postmortems when issues resolve themselves.
  }
}

## State

State {
  target = $ARGUMENTS
  currentInstrumentation = []      // discovered by assessCurrentState
  pillarDesign = {}                // built by designPillars
  slos: [SLODefinition]           // defined by defineSLOs
  alertRules: [AlertRule]         // designed by designAlerting
  dashboards = []                 // designed by designDashboards
}

## Reference Materials

See `reference/` directory for detailed methodology:
- [SLO and Alerting](reference/slo-and-alerting.md) — SLI/SLO framework, error budgets, alerting strategies, dashboard design patterns

See `references/` directory for implementation patterns:
- [Monitoring Patterns](references/monitoring-patterns.md) — RED/USE methods, distributed tracing, log patterns, alert templates, dashboard layouts

## Workflow

fn assessCurrentState(target) {
  Analyze project to determine:
    - Existing monitoring and instrumentation
    - Service architecture (monolith, microservices, serverless)
    - Current pain points (blind spots, alert fatigue, slow diagnosis)
    - Technology stack and monitoring platform in use
    - Team maturity with observability practices
}

fn designPillars(requirements) {
  For each pillar, define the instrumentation approach:

  ### Metrics
  Select methodology based on service type:
    match (serviceType) {
      requestDriven   => RED method (Rate, Errors, Duration)
      resourceFocused => USE method (Utilization, Saturation, Errors)
      general         => Four Golden Signals (Latency, Traffic, Errors, Saturation)
    }

  Metric types and their uses:

  | Type | Use Case | Example |
  |------|----------|---------|
  | Counter | Cumulative values that only increase | Total requests, errors, bytes |
  | Gauge | Values that go up and down | Memory, active connections |
  | Histogram | Distribution of values in buckets | Request latency, payload sizes |
  | Summary | Pre-computed client-side percentiles | Response time percentiles |

  ### Logs
  Design structured logging with required fields:
    - timestamp (ISO 8601 with timezone)
    - level (ERROR, WARN, INFO, DEBUG)
    - message (human-readable)
    - service (service identifier)
    - trace_id (correlation identifier)

  ### Traces
  Design distributed tracing with:
    - Context propagation (W3C Trace Context standard)
    - Span naming conventions (METHOD /path for HTTP, db.operation table for DB)
    - Sampling strategy (head-based, tail-based, rate-limited, priority)

  For detailed implementation patterns, load references/monitoring-patterns.md.
}

fn defineSLOs(services) {
  For each service, define SLIs and SLOs:

  SLO formula:
    (Good events / Total events) >= Target over Window

  Error budget calculation:
    Budget = 1 - SLO Target
    99.9% SLO = 0.1% error budget = 43.2 minutes per 30 days

  Error budget policies:
    budget remaining    => continue feature development
    budget depleted     => focus on reliability work
    budget burning fast => freeze deploys, investigate

  For detailed SLO framework, load reference/slo-and-alerting.md.
}

fn designAlerting(slos) {
  Design symptom-based alerts tied to SLOs:

  match (burnRate) {
    > 14.4x over 1h  => CRITICAL: fast burn, page immediately
    > 6x over 6h     => CRITICAL: sustained burn, page immediately
    > 3x over 3d     => WARNING: slow burn, create ticket
  }

  Constraints {
    require {
      Every alert includes: summary, impact description, runbook link, dashboard link.
      Alerts fire on sustained conditions, not transient spikes.
      Route by severity: critical to PagerDuty, warning to Slack, info to monitoring.
    }
  }

  For detailed alerting patterns, load reference/slo-and-alerting.md.
}

fn designDashboards(plan) {
  Design dashboard hierarchy:
    1. Service Health Overview — SLO status, error budget, key business metrics
    2. Deep-Dive Diagnostic — detailed metrics, resource utilization, dependencies
    3. Business Metrics — user-facing KPIs, conversion, revenue impact

  For detailed dashboard patterns, load reference/slo-and-alerting.md.
}

fn recommendNext(plan) {
  Suggest improvements:
    - Implement synthetic monitoring for proactive availability
    - Establish incident response procedures and postmortem templates
    - Conduct regular game days to validate observability
    - Automate common diagnostic procedures in runbooks
    - Review and prune alerts quarterly
}

observabilityDesign(target) {
  assessCurrentState(target)
    |> designPillars
    |> defineSLOs
    |> designAlerting
    |> designDashboards
    |> recommendNext
}
