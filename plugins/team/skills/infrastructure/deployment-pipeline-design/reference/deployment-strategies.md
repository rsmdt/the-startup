# Reference: Deployment Strategies

Detailed implementation patterns for each deployment strategy including architecture diagrams, configuration parameters, and rollback procedures.

---

## Blue-Green Deployment

Two identical production environments where traffic switches instantly.

```
                    Load Balancer
                         |
            +------------+------------+
            |                         |
        [Blue v1.0]              [Green v1.1]
         (active)                 (standby)
```

**When to Use:**
- Zero-downtime requirements
- Need instant rollback capability
- Sufficient infrastructure budget for duplicate environments

**Implementation Steps:**
1. Deploy new version to inactive environment (Green)
2. Run smoke tests against Green
3. Switch load balancer to Green
4. Monitor for issues
5. Keep Blue running for quick rollback
6. After confidence period, Blue becomes next deployment target

**Rollback:** Switch load balancer back to Blue (seconds)

---

## Canary Deployment

Gradually shift traffic from old version to new version.

```
Traffic Distribution Over Time:

T0:  [====== v1.0 100% ======]
T1:  [=== v1.0 95% ===][v1.1 5%]
T2:  [== v1.0 75% ==][= v1.1 25% =]
T3:  [= v1.0 50% =][== v1.1 50% ==]
T4:  [====== v1.1 100% ======]
```

**When to Use:**
- High-risk deployments
- Need to validate with real traffic
- Want gradual rollout with monitoring

**Traffic Progression (Example):**
1. 5% for 15 minutes - validate basic functionality
2. 25% for 30 minutes - monitor error rates
3. 50% for 1 hour - check performance metrics
4. 100% - full rollout

**Rollback Triggers:**
- Error rate exceeds baseline + threshold
- Latency exceeds acceptable limits
- Health check failures

---

## Rolling Deployment

Replace instances incrementally, one batch at a time.

```
Instance Pool (5 instances):

T0: [v1.0] [v1.0] [v1.0] [v1.0] [v1.0]
T1: [v1.1] [v1.0] [v1.0] [v1.0] [v1.0]
T2: [v1.1] [v1.1] [v1.0] [v1.0] [v1.0]
T3: [v1.1] [v1.1] [v1.1] [v1.0] [v1.0]
T4: [v1.1] [v1.1] [v1.1] [v1.1] [v1.0]
T5: [v1.1] [v1.1] [v1.1] [v1.1] [v1.1]
```

**When to Use:**
- Limited infrastructure resources
- Can tolerate mixed versions during deployment
- Stateless applications

**Configuration Parameters:**
- `maxUnavailable`: How many instances can be down simultaneously
- `maxSurge`: How many extra instances during deployment
- `minReadySeconds`: Wait time before considering instance healthy

---

## Feature Flags

Decouple deployment from release - deploy code without activating features.

```javascript
if (featureFlags.isEnabled('new-checkout', user)) {
  return newCheckoutFlow(cart);
} else {
  return legacyCheckoutFlow(cart);
}
```

**When to Use:**
- Long-running feature development
- A/B testing requirements
- Gradual feature rollouts
- Kill switch for problematic features

**Rollback:** Disable flag (no deployment required)

---

## Strategy Selection Matrix

| Factor | Blue-Green | Canary | Rolling | Feature Flag |
|--------|-----------|--------|---------|-------------|
| Zero downtime | Yes | Yes | Partial | Yes |
| Instant rollback | Yes | Partial | No | Yes |
| Infrastructure cost | High | Medium | Low | Low |
| Real traffic validation | No | Yes | No | Yes |
| Mixed version tolerance | No | Yes | Yes | Yes |
| Complexity | Medium | High | Low | Medium |
