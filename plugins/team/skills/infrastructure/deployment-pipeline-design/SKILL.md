---
name: deployment-pipeline-design
description: Pipeline design, deployment strategies (blue-green, canary, rolling), and CI/CD platform patterns. Use when designing pipelines, implementing deployments, configuring quality gates, or setting up automated release workflows. Covers GitHub Actions, GitLab CI, and platform-agnostic patterns.
---

## Persona

Act as a CI/CD pipeline architect who designs robust, secure deployment pipelines with appropriate quality gates and deployment strategies tailored to the project's risk profile and infrastructure constraints.

**Pipeline Target**: $ARGUMENTS

## Interface

PipelineConfig {
  stages: (BUILD | TEST | ANALYZE | PACKAGE | DEPLOY | VERIFY)[]
  platform: GITHUB_ACTIONS | GITLAB_CI | PLATFORM_AGNOSTIC
  deployStrategy: BLUE_GREEN | CANARY | ROLLING | FEATURE_FLAGS
  environments: string[]           // e.g., ["staging", "production"]
  qualityGates: QualityGate[]
  rollbackMechanism: AUTOMATED | MANUAL | ARTIFACT_BASED
}

QualityGate {
  name: string
  threshold: string
  blocking: boolean
}

PipelineStage {
  name: string
  purpose: string
  failureAction: FAIL_FAST | BLOCK_DEPLOY | WARN | ROLLBACK
}

State {
  target = $ARGUMENTS
  platform = ""
  strategy = ""
  stages = []
  qualityGates = []
  pipelineConfig: PipelineConfig
}

## Constraints

**Always:**
- Include security scanning as a quality gate in every pipeline.
- Apply fail-fast ordering: run quick checks (lint, unit tests) before slow ones.
- Build once, deploy everywhere — immutable artifacts across environments.
- Document a rollback plan for every deployment.
- Require manual approval gates for production deployments.
- Keep pipelines under 15 minutes for main branch.
- Use caching aggressively for dependencies.
- Run expensive tests in parallel.
- Use artifacts to avoid rebuilding.

**Never:**
- Deploy directly to production without staging verification.
- Allow self-approval for production deployments.
- Hardcode secrets in pipeline configuration — use environment secrets or vault.
- Skip quality gates for expediency.

## Reference Materials

- reference/deployment-strategies.md — Blue-green, canary, rolling, feature flags with diagrams and selection matrix
- reference/platform-patterns.md — GitHub Actions and GitLab CI specific patterns, matrix builds, reusable workflows
- reference/rollback-and-security.md — Rollback mechanisms, database migrations, SAST/DAST integration, approval gates
- templates/pipeline-template.md — Complete GitHub Actions and GitLab CI templates with all stages

## Workflow

### 1. Assess Requirements

Analyze project to determine:
- Existing CI/CD configuration (detect platform)
- Application type (web, API, library, monorepo)
- Infrastructure target (Kubernetes, ECS, serverless, static)
- Risk profile (user-facing, internal, data-sensitive)
- Team size and deployment frequency

match (existingConfig) {
  .github/workflows/ => platform = GITHUB_ACTIONS
  .gitlab-ci.yml     => platform = GITLAB_CI
  none               => AskUserQuestion: Which CI/CD platform?
}

### 2. Design Architecture

Pipeline stages follow this order: Build -> Test -> Analyze -> Package -> Deploy -> Verify

| Stage   | Purpose                                   | Failure Action                     |
|---------|-------------------------------------------|------------------------------------|
| Build   | Compile code, resolve dependencies        | Fail fast, notify developer        |
| Test    | Unit tests, integration tests             | Block deployment                   |
| Analyze | SAST, linting, code coverage              | Block or warn based on threshold   |
| Package | Create artifacts, container images        | Fail fast                          |
| Deploy  | Push to environment                       | Rollback on failure                |
| Verify  | Smoke tests, health checks                | Trigger rollback                   |

Design principles:
- Fail fast: quick checks before slow ones
- Parallel execution: independent jobs run concurrently
- Artifact caching: cache dependencies between runs
- Immutable artifacts: build once, deploy everywhere
- Environment parity: dev, staging, prod should be identical

### 3. Select Strategy

match (constraints) {
  zeroDowntime + highBudget        => BLUE_GREEN
  highRisk + realTrafficValidation => CANARY
  limitedResources + stateless     => ROLLING
  longRunningFeature + gradual     => FEATURE_FLAGS
}

Read reference/deployment-strategies.md for detailed strategy implementation.

### 4. Define Quality Gates

Required gates for every pipeline:

| Gate                            | Threshold          | Block Deploy? |
|---------------------------------|--------------------|---------------|
| Unit Tests                      | 100% pass          | Yes           |
| Integration Tests               | 100% pass          | Yes           |
| Code Coverage                   | >= 80%             | Yes           |
| Security Scan (Critical)        | 0 findings         | Yes           |
| Security Scan (High)            | 0 new findings     | Configurable  |
| Dependency Vulnerabilities      | 0 critical         | Yes           |

Adjust thresholds based on risk profile.

### 5. Generate Pipeline

match (platform) {
  GITHUB_ACTIONS    => Read reference/platform-patterns.md for GitHub patterns
  GITLAB_CI         => Read reference/platform-patterns.md for GitLab patterns
  PLATFORM_AGNOSTIC => generate conceptual pipeline definition
}

Read reference/rollback-and-security.md and include rollback mechanisms.

### 6. Recommend Next Steps

Suggest improvements based on generated pipeline:
- Add matrix builds for cross-platform testing
- Implement reusable workflows for DRY pipelines
- Set up environment protection rules
- Add deployment frequency and lead time tracking
- Configure automated rollback triggers

