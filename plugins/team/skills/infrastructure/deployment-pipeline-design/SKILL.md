---
name: deployment-pipeline-design
description: Pipeline design, deployment strategies (blue-green, canary, rolling), and CI/CD platform patterns. Use when designing pipelines, implementing deployments, configuring quality gates, or setting up automated release workflows. Covers GitHub Actions, GitLab CI, and platform-agnostic patterns.
---

# CI/CD Patterns

A comprehensive skill for designing and implementing continuous integration and deployment pipelines. Covers pipeline architecture, deployment strategies, quality gates, and platform-specific patterns for GitHub Actions and GitLab CI.

## When to Use

- Designing new CI/CD pipelines from scratch
- Implementing deployment strategies (blue-green, canary, rolling)
- Setting up quality gates and approval workflows
- Configuring GitHub Actions or GitLab CI pipelines
- Implementing automated rollback mechanisms
- Creating multi-environment deployment workflows
- Integrating security scanning into pipelines

## Pipeline Architecture

### Pipeline Stages

A well-designed pipeline follows these stages in order:

```
Build -> Test -> Analyze -> Package -> Deploy -> Verify
```

**Stage Breakdown:**

| Stage | Purpose | Failure Action |
|-------|---------|----------------|
| Build | Compile code, resolve dependencies | Fail fast, notify developer |
| Test | Unit tests, integration tests | Block deployment |
| Analyze | SAST, linting, code coverage | Block or warn based on threshold |
| Package | Create artifacts, container images | Fail fast |
| Deploy | Push to environment | Rollback on failure |
| Verify | Smoke tests, health checks | Trigger rollback |

### Pipeline Design Principles

1. **Fail Fast**: Run quick checks (lint, unit tests) before slow ones
2. **Parallel Execution**: Run independent jobs concurrently
3. **Artifact Caching**: Cache dependencies between runs
4. **Immutable Artifacts**: Build once, deploy everywhere
5. **Environment Parity**: Dev, staging, and prod should be identical

## Deployment Strategies

### Blue-Green Deployment

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

### Canary Deployment

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

### Rolling Deployment

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

### Feature Flags

Decouple deployment from release - deploy code without activating features.

```
Code deployed with feature flag:

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

## Quality Gates

### Required Gates

Every pipeline should include these gates:

| Gate | Threshold | Block Deploy? |
|------|-----------|---------------|
| Unit Tests | 100% pass | Yes |
| Integration Tests | 100% pass | Yes |
| Code Coverage | >= 80% | Yes |
| Security Scan (Critical) | 0 findings | Yes |
| Security Scan (High) | 0 new findings | Configurable |
| Dependency Vulnerabilities | 0 critical | Yes |

### Manual Approval Gates

Use for production deployments:

```yaml
# Conceptual flow
stages:
  - test
  - deploy-staging
  - approval        # Manual gate
  - deploy-prod
  - verify
```

**Approval Requirements:**
- At least 2 approvers for production
- No self-approval allowed
- Time-boxed approval windows
- Audit trail of approvals

## GitHub Actions Patterns

### Basic Pipeline Structure

```yaml
name: CI/CD Pipeline

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

concurrency:
  group: ${{ github.workflow }}-${{ github.ref }}
  cancel-in-progress: true

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: '20'
          cache: 'npm'
      - run: npm ci
      - run: npm run build
      - uses: actions/upload-artifact@v4
        with:
          name: build
          path: dist/

  test:
    needs: build
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/download-artifact@v4
        with:
          name: build
          path: dist/
      - run: npm ci
      - run: npm test

  deploy-staging:
    needs: test
    if: github.ref == 'refs/heads/main'
    environment: staging
    runs-on: ubuntu-latest
    steps:
      - uses: actions/download-artifact@v4
        with:
          name: build
      - run: ./deploy.sh staging

  deploy-prod:
    needs: deploy-staging
    if: github.ref == 'refs/heads/main'
    environment: production
    runs-on: ubuntu-latest
    steps:
      - uses: actions/download-artifact@v4
        with:
          name: build
      - run: ./deploy.sh production
```

### Matrix Builds

Run tests across multiple configurations:

```yaml
jobs:
  test:
    strategy:
      fail-fast: false
      matrix:
        os: [ubuntu-latest, windows-latest, macos-latest]
        node: [18, 20, 22]
    runs-on: ${{ matrix.os }}
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: ${{ matrix.node }}
      - run: npm ci
      - run: npm test
```

### Reusable Workflows

Create reusable workflow in `.github/workflows/deploy-reusable.yml`:

```yaml
name: Reusable Deploy

on:
  workflow_call:
    inputs:
      environment:
        required: true
        type: string
    secrets:
      DEPLOY_KEY:
        required: true

jobs:
  deploy:
    environment: ${{ inputs.environment }}
    runs-on: ubuntu-latest
    steps:
      - run: ./deploy.sh ${{ inputs.environment }}
        env:
          DEPLOY_KEY: ${{ secrets.DEPLOY_KEY }}
```

Call from another workflow:

```yaml
jobs:
  deploy-staging:
    uses: ./.github/workflows/deploy-reusable.yml
    with:
      environment: staging
    secrets:
      DEPLOY_KEY: ${{ secrets.STAGING_DEPLOY_KEY }}
```

### Environment Protection Rules

Configure in repository settings:

- Required reviewers for production
- Wait timer (e.g., 15 minutes before prod deploy)
- Restrict to specific branches
- Required status checks

## GitLab CI Patterns

### Basic Pipeline Structure

```yaml
stages:
  - build
  - test
  - deploy

variables:
  NODE_VERSION: "20"

default:
  image: node:${NODE_VERSION}
  cache:
    paths:
      - node_modules/

build:
  stage: build
  script:
    - npm ci
    - npm run build
  artifacts:
    paths:
      - dist/
    expire_in: 1 hour

test:unit:
  stage: test
  script:
    - npm ci
    - npm run test:unit
  coverage: '/Coverage: \d+\.\d+%/'

test:integration:
  stage: test
  services:
    - postgres:15
  variables:
    POSTGRES_DB: test
    POSTGRES_USER: test
    POSTGRES_PASSWORD: test
  script:
    - npm ci
    - npm run test:integration

deploy:staging:
  stage: deploy
  environment:
    name: staging
    url: https://staging.example.com
  script:
    - ./deploy.sh staging
  only:
    - main

deploy:production:
  stage: deploy
  environment:
    name: production
    url: https://example.com
  script:
    - ./deploy.sh production
  when: manual
  only:
    - main
```

### Pipeline Rules

```yaml
deploy:production:
  rules:
    - if: $CI_COMMIT_BRANCH == "main"
      when: manual
    - if: $CI_COMMIT_TAG
      when: on_success
    - when: never
```

### Include Templates

```yaml
include:
  - template: Security/SAST.gitlab-ci.yml
  - template: Security/Dependency-Scanning.gitlab-ci.yml
  - local: .gitlab/ci/deploy.yml
  - project: 'devops/ci-templates'
    ref: main
    file: '/templates/docker-build.yml'
```

### Dynamic Environments

```yaml
deploy:review:
  stage: deploy
  environment:
    name: review/$CI_COMMIT_REF_SLUG
    url: https://$CI_COMMIT_REF_SLUG.review.example.com
    on_stop: stop:review
  script:
    - ./deploy.sh review
  only:
    - merge_requests

stop:review:
  stage: deploy
  environment:
    name: review/$CI_COMMIT_REF_SLUG
    action: stop
  script:
    - ./teardown.sh review
  when: manual
  only:
    - merge_requests
```

## Rollback Mechanisms

### Automated Rollback Triggers

```yaml
# Conceptual rollback configuration
rollback:
  triggers:
    - metric: error_rate
      threshold: 5%
      window: 5m
    - metric: latency_p99
      threshold: 2000ms
      window: 5m
    - metric: health_check_failures
      threshold: 3
      window: 1m
  action:
    type: previous_version
    notify:
      - slack: #deployments
      - pagerduty: on-call
```

### Database Migration Rollback

1. **Forward-only migrations** (preferred):
   - Never use destructive operations (DROP, DELETE)
   - Add new columns as nullable
   - Use feature flags to switch behavior
   - Clean up old columns in later release

2. **Rollback migrations**:
   - Every migration must have a corresponding rollback
   - Test rollbacks in staging before production
   - Keep rollback window defined (e.g., 24 hours)

### Artifact-Based Rollback

```yaml
rollback:production:
  stage: deploy
  environment:
    name: production
  script:
    - PREVIOUS_VERSION=$(get-previous-version.sh)
    - ./deploy.sh production $PREVIOUS_VERSION
  when: manual
  only:
    - main
```

## Security Integration

### SAST/DAST Integration

```yaml
security:sast:
  stage: analyze
  image: security-scanner:latest
  script:
    - sast-scan --format sarif --output sast-results.sarif
  artifacts:
    reports:
      sast: sast-results.sarif

security:dependency:
  stage: analyze
  script:
    - npm audit --audit-level=high
    - trivy fs --security-checks vuln .
```

### Secret Scanning

- Never commit secrets to repository
- Use environment secrets or vault integration
- Scan for exposed secrets in pre-commit hooks
- Rotate secrets immediately if exposed

## Best Practices

### Pipeline Design

- Keep pipelines under 15 minutes for main branch
- Use caching aggressively for dependencies
- Run expensive tests in parallel
- Fail fast with quick checks first
- Use artifacts to avoid rebuilding

### Deployment Safety

- Always have a rollback plan
- Deploy to staging before production
- Use feature flags for risky changes
- Monitor deployments in real-time
- Document deployment procedures

### Quality Assurance

- Enforce code coverage thresholds
- Block deployments on security findings
- Require peer approval for production
- Maintain environment parity
- Test rollback procedures regularly

### Observability

- Log all deployment events
- Track deployment frequency and lead time
- Monitor change failure rate
- Measure mean time to recovery
- Alert on deployment failures

## References

- `templates/pipeline-template.md` - Complete pipeline template with all stages
