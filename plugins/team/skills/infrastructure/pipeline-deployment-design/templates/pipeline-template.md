# Template: CI/CD Pipeline

## Purpose

Complete pipeline template covering build, test, security scanning, and multi-environment deployment. Use this as a starting point for new projects or when modernizing existing pipelines.

## GitHub Actions Template

```yaml
# .github/workflows/ci-cd.yml
# Complete CI/CD Pipeline Template
#
# Features:
# - Build and test with caching
# - Security scanning (SAST, dependencies)
# - Multi-environment deployment (staging, production)
# - Manual approval for production
# - Automated rollback support

name: CI/CD Pipeline

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]
  workflow_dispatch:
    inputs:
      environment:
        description: 'Deploy to environment'
        required: true
        type: choice
        options:
          - staging
          - production

# Prevent concurrent deployments to same environment
concurrency:
  group: deploy-${{ github.ref }}
  cancel-in-progress: false

env:
  NODE_VERSION: '20'
  REGISTRY: ghcr.io
  IMAGE_NAME: ${{ github.repository }}

jobs:
  # ============================================
  # BUILD STAGE
  # ============================================
  build:
    name: Build
    runs-on: ubuntu-latest
    outputs:
      version: ${{ steps.version.outputs.version }}
    steps:
      - name: Checkout
        uses: actions/checkout@v4

      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: ${{ env.NODE_VERSION }}
          cache: 'npm'

      - name: Install dependencies
        run: npm ci

      - name: Build
        run: npm run build

      - name: Generate version
        id: version
        run: echo "version=${{ github.sha }}-$(date +%Y%m%d%H%M%S)" >> $GITHUB_OUTPUT

      - name: Upload build artifact
        uses: actions/upload-artifact@v4
        with:
          name: build-${{ github.sha }}
          path: dist/
          retention-days: 7

  # ============================================
  # TEST STAGE
  # ============================================
  test-unit:
    name: Unit Tests
    needs: build
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v4

      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: ${{ env.NODE_VERSION }}
          cache: 'npm'

      - name: Install dependencies
        run: npm ci

      - name: Run unit tests
        run: npm run test:unit -- --coverage

      - name: Upload coverage
        uses: actions/upload-artifact@v4
        with:
          name: coverage-unit
          path: coverage/

  test-integration:
    name: Integration Tests
    needs: build
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_USER: test
          POSTGRES_PASSWORD: test
          POSTGRES_DB: test
        ports:
          - 5432:5432
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
    steps:
      - name: Checkout
        uses: actions/checkout@v4

      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: ${{ env.NODE_VERSION }}
          cache: 'npm'

      - name: Install dependencies
        run: npm ci

      - name: Run integration tests
        run: npm run test:integration
        env:
          DATABASE_URL: postgresql://test:test@localhost:5432/test

  test-e2e:
    name: E2E Tests
    needs: build
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v4

      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: ${{ env.NODE_VERSION }}
          cache: 'npm'

      - name: Install dependencies
        run: npm ci

      - name: Download build
        uses: actions/download-artifact@v4
        with:
          name: build-${{ github.sha }}
          path: dist/

      - name: Run E2E tests
        run: npm run test:e2e

  # ============================================
  # ANALYZE STAGE
  # ============================================
  lint:
    name: Lint
    needs: build
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v4

      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: ${{ env.NODE_VERSION }}
          cache: 'npm'

      - name: Install dependencies
        run: npm ci

      - name: Run linter
        run: npm run lint

  security-sast:
    name: SAST Scan
    needs: build
    runs-on: ubuntu-latest
    permissions:
      security-events: write
    steps:
      - name: Checkout
        uses: actions/checkout@v4

      - name: Run CodeQL
        uses: github/codeql-action/init@v3
        with:
          languages: javascript

      - name: Autobuild
        uses: github/codeql-action/autobuild@v3

      - name: Perform CodeQL Analysis
        uses: github/codeql-action/analyze@v3

  security-dependencies:
    name: Dependency Scan
    needs: build
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v4

      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: ${{ env.NODE_VERSION }}
          cache: 'npm'

      - name: Install dependencies
        run: npm ci

      - name: Audit dependencies
        run: npm audit --audit-level=high

  # ============================================
  # PACKAGE STAGE
  # ============================================
  package:
    name: Package
    needs: [test-unit, test-integration, lint, security-sast, security-dependencies]
    runs-on: ubuntu-latest
    permissions:
      contents: read
      packages: write
    outputs:
      image-tag: ${{ steps.meta.outputs.tags }}
    steps:
      - name: Checkout
        uses: actions/checkout@v4

      - name: Download build
        uses: actions/download-artifact@v4
        with:
          name: build-${{ github.sha }}
          path: dist/

      - name: Setup Docker Buildx
        uses: docker/setup-buildx-action@v3

      - name: Login to Registry
        uses: docker/login-action@v3
        with:
          registry: ${{ env.REGISTRY }}
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}

      - name: Extract metadata
        id: meta
        uses: docker/metadata-action@v5
        with:
          images: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}
          tags: |
            type=sha,prefix=
            type=ref,event=branch
            type=semver,pattern={{version}}

      - name: Build and push
        uses: docker/build-push-action@v5
        with:
          context: .
          push: true
          tags: ${{ steps.meta.outputs.tags }}
          labels: ${{ steps.meta.outputs.labels }}
          cache-from: type=gha
          cache-to: type=gha,mode=max

  # ============================================
  # DEPLOY STAGING
  # ============================================
  deploy-staging:
    name: Deploy Staging
    needs: package
    if: github.ref == 'refs/heads/main'
    runs-on: ubuntu-latest
    environment:
      name: staging
      url: https://staging.example.com
    steps:
      - name: Checkout
        uses: actions/checkout@v4

      - name: Deploy to staging
        run: |
          echo "Deploying to staging..."
          # Replace with actual deployment command
          # kubectl set image deployment/app app=${{ needs.package.outputs.image-tag }}
          # OR
          # aws ecs update-service --cluster staging --service app --force-new-deployment

      - name: Wait for deployment
        run: |
          echo "Waiting for deployment to complete..."
          # kubectl rollout status deployment/app --timeout=300s
          sleep 30

      - name: Run smoke tests
        run: |
          echo "Running smoke tests..."
          # curl -f https://staging.example.com/health || exit 1

  # ============================================
  # DEPLOY PRODUCTION
  # ============================================
  deploy-production:
    name: Deploy Production
    needs: deploy-staging
    if: github.ref == 'refs/heads/main'
    runs-on: ubuntu-latest
    environment:
      name: production
      url: https://example.com
    steps:
      - name: Checkout
        uses: actions/checkout@v4

      - name: Deploy to production
        run: |
          echo "Deploying to production..."
          # Replace with actual deployment command

      - name: Wait for deployment
        run: |
          echo "Waiting for deployment to complete..."
          sleep 30

      - name: Run smoke tests
        run: |
          echo "Running smoke tests..."
          # curl -f https://example.com/health || exit 1

      - name: Notify success
        if: success()
        run: |
          echo "Deployment successful!"
          # Send Slack notification, etc.

  # ============================================
  # ROLLBACK (Manual Trigger)
  # ============================================
  rollback:
    name: Rollback Production
    if: github.event_name == 'workflow_dispatch' && failure()
    needs: deploy-production
    runs-on: ubuntu-latest
    environment:
      name: production
    steps:
      - name: Rollback deployment
        run: |
          echo "Rolling back production..."
          # kubectl rollout undo deployment/app
          # OR
          # aws ecs update-service --cluster production --service app --task-definition previous-version
```

## GitLab CI Template

```yaml
# .gitlab-ci.yml
# Complete CI/CD Pipeline Template
#
# Features:
# - Build and test with caching
# - Security scanning (SAST, dependencies)
# - Multi-environment deployment (staging, production)
# - Manual approval for production
# - Automated rollback support

stages:
  - build
  - test
  - analyze
  - package
  - deploy
  - verify

variables:
  NODE_VERSION: "20"
  DOCKER_TLS_CERTDIR: "/certs"

default:
  image: node:${NODE_VERSION}
  cache:
    key: ${CI_COMMIT_REF_SLUG}
    paths:
      - node_modules/
      - .npm/

# ============================================
# BUILD STAGE
# ============================================
build:
  stage: build
  script:
    - npm ci --cache .npm --prefer-offline
    - npm run build
  artifacts:
    paths:
      - dist/
    expire_in: 1 day

# ============================================
# TEST STAGE
# ============================================
test:unit:
  stage: test
  needs: [build]
  script:
    - npm ci --cache .npm --prefer-offline
    - npm run test:unit -- --coverage
  coverage: '/All files[^|]*\|[^|]*\s+([\d\.]+)/'
  artifacts:
    reports:
      coverage_report:
        coverage_format: cobertura
        path: coverage/cobertura-coverage.xml
    paths:
      - coverage/

test:integration:
  stage: test
  needs: [build]
  services:
    - name: postgres:15
      alias: database
  variables:
    POSTGRES_USER: test
    POSTGRES_PASSWORD: test
    POSTGRES_DB: test
    DATABASE_URL: postgresql://test:test@database:5432/test
  script:
    - npm ci --cache .npm --prefer-offline
    - npm run test:integration

test:e2e:
  stage: test
  needs: [build]
  script:
    - npm ci --cache .npm --prefer-offline
    - npm run test:e2e

# ============================================
# ANALYZE STAGE
# ============================================
lint:
  stage: analyze
  needs: [build]
  script:
    - npm ci --cache .npm --prefer-offline
    - npm run lint
  allow_failure: false

# Include GitLab security templates
include:
  - template: Security/SAST.gitlab-ci.yml
  - template: Security/Dependency-Scanning.gitlab-ci.yml
  - template: Security/Secret-Detection.gitlab-ci.yml

sast:
  stage: analyze
  needs: []

dependency_scanning:
  stage: analyze
  needs: []

secret_detection:
  stage: analyze
  needs: []

# ============================================
# PACKAGE STAGE
# ============================================
package:docker:
  stage: package
  needs:
    - job: build
      artifacts: true
    - job: test:unit
    - job: test:integration
    - job: lint
  image: docker:24
  services:
    - docker:24-dind
  variables:
    DOCKER_HOST: tcp://docker:2376
  script:
    - docker login -u $CI_REGISTRY_USER -p $CI_REGISTRY_PASSWORD $CI_REGISTRY
    - docker build -t $CI_REGISTRY_IMAGE:$CI_COMMIT_SHA -t $CI_REGISTRY_IMAGE:latest .
    - docker push $CI_REGISTRY_IMAGE:$CI_COMMIT_SHA
    - docker push $CI_REGISTRY_IMAGE:latest
  only:
    - main
    - tags

# ============================================
# DEPLOY STAGING
# ============================================
deploy:staging:
  stage: deploy
  needs: [package:docker]
  environment:
    name: staging
    url: https://staging.example.com
  script:
    - echo "Deploying to staging..."
    # Replace with actual deployment commands
    # - kubectl set image deployment/app app=$CI_REGISTRY_IMAGE:$CI_COMMIT_SHA
  only:
    - main

verify:staging:
  stage: verify
  needs: [deploy:staging]
  environment:
    name: staging
    url: https://staging.example.com
  script:
    - echo "Running smoke tests on staging..."
    # - curl -f https://staging.example.com/health
  only:
    - main

# ============================================
# DEPLOY PRODUCTION
# ============================================
deploy:production:
  stage: deploy
  needs: [verify:staging]
  environment:
    name: production
    url: https://example.com
  script:
    - echo "Deploying to production..."
    # Replace with actual deployment commands
  when: manual
  only:
    - main

verify:production:
  stage: verify
  needs: [deploy:production]
  environment:
    name: production
    url: https://example.com
  script:
    - echo "Running smoke tests on production..."
    # - curl -f https://example.com/health
  only:
    - main

# ============================================
# ROLLBACK (Manual)
# ============================================
rollback:production:
  stage: deploy
  environment:
    name: production
    url: https://example.com
  script:
    - echo "Rolling back production..."
    # - kubectl rollout undo deployment/app
  when: manual
  only:
    - main
```

## Usage Instructions

1. Copy the appropriate template (GitHub Actions or GitLab CI)
2. Replace placeholder deployment commands with actual commands for your platform
3. Configure environment protection rules in your repository settings
4. Set up required secrets:
   - Container registry credentials
   - Deployment credentials
   - Notification service tokens (Slack, etc.)
5. Customize test commands to match your project setup
6. Adjust Docker build context and configuration as needed
7. Configure environment URLs

## Customization Points

| Section | What to Customize |
|---------|-------------------|
| `NODE_VERSION` | Match your project's Node.js version |
| Database service | Replace with your database (MySQL, MongoDB, etc.) |
| `npm run` commands | Match your package.json scripts |
| Registry | Use your container registry (ECR, GCR, Docker Hub) |
| Deployment commands | Replace with kubectl, aws, gcloud, etc. |
| Environment URLs | Set actual staging/production URLs |
| Smoke tests | Add actual health check endpoints |

## Environment Variables Required

### GitHub Actions
- `GITHUB_TOKEN` - Automatically provided
- Deployment secrets configured per environment

### GitLab CI
- `CI_REGISTRY_USER` - Automatically provided
- `CI_REGISTRY_PASSWORD` - Automatically provided
- Deployment variables in CI/CD settings

## Examples

See CI/CD implementations in production projects for real-world examples with:
- Kubernetes deployments
- AWS ECS/Fargate
- Vercel/Netlify
- Cloud Run/App Engine
