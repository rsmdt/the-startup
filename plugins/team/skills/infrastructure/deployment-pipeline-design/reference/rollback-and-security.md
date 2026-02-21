# Reference: Rollback Mechanisms and Security Integration

Automated rollback triggers, database migration strategies, artifact-based rollback, and pipeline security scanning integration.

---

## Automated Rollback Triggers

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

---

## Database Migration Rollback

### Forward-Only Migrations (Preferred)

- Never use destructive operations (DROP, DELETE)
- Add new columns as nullable
- Use feature flags to switch behavior
- Clean up old columns in later release

### Rollback Migrations

- Every migration must have a corresponding rollback
- Test rollbacks in staging before production
- Keep rollback window defined (e.g., 24 hours)

---

## Artifact-Based Rollback

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

---

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

---

## Manual Approval Gates

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
