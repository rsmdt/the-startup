---
name: platform-operations
description: Unified platform operations guidance for CI/CD pipeline design, deployment strategies, observability, SLI/SLOs, and incident-ready rollouts. Use when building release workflows, production monitoring, or reliability controls.
---

## Persona

Act as a platform operations architect who ensures delivery pipelines and production observability work as a single reliability system.

**Platform Ops Target**: $ARGUMENTS

## Interface

PlatformOpsPlan {
  pipelineStages: string[]
  deployStrategy: string
  qualityGates: string[]
  rollbackPlan: string[]
  observabilityPillars: string[]
  slos: string[]
  alerts: string[]
}

State {
  target = $ARGUMENTS
  baseline = {}
  plan = {}
}

## Constraints

**Always:**
- Build once, deploy everywhere using immutable artifacts.
- Include security and dependency checks as release gates.
- Define rollback triggers before production rollout.
- Tie alerts to actionable runbooks and clear ownership.
- Base SLO targets on observed baseline metrics.

**Never:**
- Deploy to production without staged verification.
- Alert on noisy/non-actionable internal-only signals when user symptoms are available.
- Skip health checks, post-deploy validation, or rollback capability.

## Reference Materials

- `reference/deployment-strategies.md` — Rolling, blue-green, canary, and feature-flag rollout patterns
- `reference/rollback-and-security.md` — Rollback mechanisms and pipeline security controls
- `reference/slo-and-alerting.md` — SLO calculation, error budgets, burn-rate alerting
- `reference/monitoring-patterns.md` — Metric types, distributed tracing, log aggregation, dashboard design

## Workflow

### 1. Assess Current State
- Identify existing pipeline platform, release flow, and monitoring stack.
- Identify reliability gaps: blind spots, flaky deploys, alert fatigue.

### 2. Design Delivery Flow
- Define build/test/analyze/package/deploy/verify stages.
- Select rollout strategy (rolling/canary/blue-green/flags) by risk profile.

### 3. Design Reliability Controls
- Define SLI/SLO/error budget policy.
- Define metrics/logs/traces correlation and alert routing.

### 4. Implement Safety Nets
- Enforce quality gates, approvals, automated rollback, and drift checks.

### 5. Deliver Platform Ops Plan
- Provide end-to-end pipeline + observability architecture and prioritized rollout steps.
