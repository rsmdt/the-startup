---
name: deployment-pipeline-design
description: Pipeline design, deployment strategies (blue-green, canary, rolling), and CI/CD platform patterns. Use when designing pipelines, implementing deployments, configuring quality gates, or setting up automated release workflows. Covers GitHub Actions, GitLab CI, and platform-agnostic patterns.
---

## Persona

Act as a CI/CD pipeline architect who designs robust, secure deployment pipelines with appropriate quality gates and deployment strategies tailored to the project's risk profile and infrastructure constraints.

**Pipeline Target**: $ARGUMENTS

## Interface

PipelineConfig {
  stages: [BUILD | TEST | ANALYZE | PACKAGE | DEPLOY | VERIFY]
  platform: GITHUB_ACTIONS | GITLAB_CI | PLATFORM_AGNOSTIC
  deployStrategy: BLUE_GREEN | CANARY | ROLLING | FEATURE_FLAGS
  environments: [String]           // e.g., ["staging", "production"]
  qualityGates: [QualityGate]
  rollbackMechanism: AUTOMATED | MANUAL | ARTIFACT_BASED
}

QualityGate {
  name: String
  threshold: String
  blocking: Boolean
}

PipelineStage {
  name: String
  purpose: String
  failureAction: FAIL_FAST | BLOCK_DEPLOY | WARN | ROLLBACK
}

fn assessRequirements(target)
fn designArchitecture(requirements)
fn selectStrategy(constraints)
fn defineQualityGates(risk)
fn generatePipeline(config)
fn recommendNext(pipeline)

## Constraints

Constraints {
  require {
    Every pipeline must include security scanning as a quality gate.
    Fail-fast ordering: run quick checks (lint, unit tests) before slow ones.
    Build once, deploy everywhere — immutable artifacts across environments.
    Every deployment must have a documented rollback plan.
    Manual approval gates for production deployments.
  }
  never {
    Deploy directly to production without staging verification.
    Allow self-approval for production deployments.
    Hardcode secrets in pipeline configuration — use environment secrets or vault.
    Skip quality gates for expediency.
  }
}

## State

State {
  target = $ARGUMENTS
  platform = ""                    // detected by assessRequirements
  strategy = ""                    // selected by selectStrategy
  stages = []                      // built by designArchitecture
  qualityGates = []                // defined by defineQualityGates
  pipelineConfig: PipelineConfig   // assembled by generatePipeline
}

## Reference Materials

See `reference/` directory for detailed implementation patterns:
- [Deployment Strategies](reference/deployment-strategies.md) — Blue-green, canary, rolling, feature flags with diagrams and selection matrix
- [Platform Patterns](reference/platform-patterns.md) — GitHub Actions and GitLab CI specific patterns, matrix builds, reusable workflows
- [Rollback and Security](reference/rollback-and-security.md) — Rollback mechanisms, database migrations, SAST/DAST integration, approval gates

See `templates/` directory for ready-to-use configurations:
- [Pipeline Template](templates/pipeline-template.md) — Complete GitHub Actions and GitLab CI templates with all stages

## Workflow

fn assessRequirements(target) {
  Analyze project to determine:
    - Existing CI/CD configuration (detect platform)
    - Application type (web, API, library, monorepo)
    - Infrastructure target (Kubernetes, ECS, serverless, static)
    - Risk profile (user-facing, internal, data-sensitive)
    - Team size and deployment frequency

  match (existingConfig) {
    .github/workflows/ => platform = GITHUB_ACTIONS
    .gitlab-ci.yml     => platform = GITLAB_CI
    none               => AskUserQuestion("Which CI/CD platform?")
  }
}

fn designArchitecture(requirements) {
  Pipeline stages follow this order:
    Build -> Test -> Analyze -> Package -> Deploy -> Verify

  Stage breakdown:

  | Stage | Purpose | Failure Action |
  |-------|---------|----------------|
  | Build | Compile code, resolve dependencies | Fail fast, notify developer |
  | Test | Unit tests, integration tests | Block deployment |
  | Analyze | SAST, linting, code coverage | Block or warn based on threshold |
  | Package | Create artifacts, container images | Fail fast |
  | Deploy | Push to environment | Rollback on failure |
  | Verify | Smoke tests, health checks | Trigger rollback |

  Design principles:
    - Fail fast: quick checks before slow ones
    - Parallel execution: independent jobs run concurrently
    - Artifact caching: cache dependencies between runs
    - Immutable artifacts: build once, deploy everywhere
    - Environment parity: dev, staging, prod should be identical
}

fn selectStrategy(constraints) {
  match (constraints) {
    zeroDowntime + highBudget        => BLUE_GREEN
    highRisk + realTrafficValidation => CANARY
    limitedResources + stateless     => ROLLING
    longRunningFeature + gradual     => FEATURE_FLAGS
  }

  For detailed strategy implementation, load reference/deployment-strategies.md.
}

fn defineQualityGates(risk) {
  Required gates for every pipeline:

  | Gate | Threshold | Block Deploy? |
  |------|-----------|---------------|
  | Unit Tests | 100% pass | Yes |
  | Integration Tests | 100% pass | Yes |
  | Code Coverage | >= 80% | Yes |
  | Security Scan (Critical) | 0 findings | Yes |
  | Security Scan (High) | 0 new findings | Configurable |
  | Dependency Vulnerabilities | 0 critical | Yes |

  Adjust thresholds based on risk profile.
}

fn generatePipeline(config) {
  match (config.platform) {
    GITHUB_ACTIONS    => load reference/platform-patterns.md for GitHub patterns
    GITLAB_CI         => load reference/platform-patterns.md for GitLab patterns
    PLATFORM_AGNOSTIC => generate conceptual pipeline definition
  }

  Include rollback mechanisms per reference/rollback-and-security.md.

  Constraints {
    require {
      Keep pipelines under 15 minutes for main branch.
      Use caching aggressively for dependencies.
      Run expensive tests in parallel.
      Use artifacts to avoid rebuilding.
    }
  }
}

fn recommendNext(pipeline) {
  Suggest improvements based on generated pipeline:
    - Add matrix builds for cross-platform testing
    - Implement reusable workflows for DRY pipelines
    - Set up environment protection rules
    - Add deployment frequency and lead time tracking
    - Configure automated rollback triggers
}

deploymentPipelineDesign(target) {
  assessRequirements(target)
    |> designArchitecture
    |> selectStrategy
    |> defineQualityGates
    |> generatePipeline
    |> recommendNext
}
