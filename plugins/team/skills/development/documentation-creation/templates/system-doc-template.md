# Template: System Documentation

## Purpose

Use this template to document a system's architecture, components, and operations. System documentation helps new team members onboard, supports operational decisions, and preserves critical knowledge about how the system works and why.

## Template

```markdown
# [System Name]

**Version:** [Document version, e.g., 1.0]
**Last Updated:** [YYYY-MM-DD]
**Owner:** [Team or individual responsible]

## Overview

### Purpose

[1-2 paragraphs explaining what this system does and why it exists.
Focus on the business problem it solves, not implementation details.]

### Key Capabilities

- [Capability 1: What the system enables users/other systems to do]
- [Capability 2]
- [Capability 3]

### Scope

**In Scope:**
- [What this system is responsible for]

**Out of Scope:**
- [What this system explicitly does NOT handle]
- [Common misconceptions about system boundaries]

---

## Architecture

### System Context

[Diagram showing this system and its external interactions]

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Users     │────▶│   System    │────▶│  External   │
│             │     │             │     │  Service    │
└─────────────┘     └─────────────┘     └─────────────┘
                           │
                           ▼
                    ┌─────────────┐
                    │  Database   │
                    └─────────────┘
```

[Description of each external entity and its relationship to the system]

### Component Architecture

[Diagram showing major internal components]

```
┌──────────────────────────────────────────────────┐
│                    System                         │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐       │
│  │   API    │  │  Worker  │  │ Scheduler│       │
│  │  Layer   │  │  Service │  │          │       │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘       │
│       │             │             │              │
│       └─────────────┼─────────────┘              │
│                     ▼                            │
│             ┌──────────────┐                     │
│             │    Domain    │                     │
│             │    Layer     │                     │
│             └──────────────┘                     │
└──────────────────────────────────────────────────┘
```

### Component Descriptions

| Component | Responsibility | Technology |
|-----------|---------------|------------|
| [API Layer] | [Handles HTTP requests, authentication] | [Node.js, Express] |
| [Worker Service] | [Processes async jobs from queue] | [Node.js, Bull] |
| [Domain Layer] | [Business logic, validations] | [TypeScript] |

---

## Data Architecture

### Data Stores

| Store | Type | Purpose | Retention |
|-------|------|---------|-----------|
| [Primary DB] | [PostgreSQL] | [Transactional data] | [Indefinite] |
| [Cache] | [Redis] | [Session, hot data] | [24 hours] |
| [Search] | [Elasticsearch] | [Full-text search] | [90 days] |

### Key Entities

[Diagram or table showing main data entities and relationships]

```
┌─────────┐       ┌─────────┐       ┌─────────┐
│  User   │──────▶│  Order  │──────▶│  Item   │
└─────────┘  1:N  └─────────┘  1:N  └─────────┘
```

### Data Flows

[Describe how data moves through the system for key operations]

**[Flow Name, e.g., Order Creation]:**
1. User submits order via API
2. API validates request and creates order record
3. Order event published to message queue
4. Worker processes order, updates inventory
5. Notification sent to user

---

## Integration Points

### Upstream Dependencies

[Services this system depends on]

| Service | Purpose | Protocol | SLA |
|---------|---------|----------|-----|
| [Auth Service] | [User authentication] | [REST/OAuth2] | [99.9%] |
| [Payment Gateway] | [Process payments] | [REST] | [99.95%] |

### Downstream Consumers

[Services that depend on this system]

| Consumer | Purpose | Protocol | Contract |
|----------|---------|----------|----------|
| [Mobile App] | [User interface] | [REST API] | [API v2] |
| [Analytics] | [Business metrics] | [Event stream] | [Event schema v1] |

### API Contracts

[Link to API documentation or summary of key endpoints]

---

## Deployment

### Infrastructure

[Diagram showing deployment architecture]

```
┌─────────────────────────────────────────────┐
│                   AWS                        │
│  ┌─────────┐    ┌──────────────────────┐    │
│  │   ALB   │───▶│    ECS Cluster       │    │
│  └─────────┘    │  ┌────┐ ┌────┐ ┌────┐│    │
│                 │  │ API│ │ API│ │ API││    │
│                 │  └────┘ └────┘ └────┘│    │
│                 └──────────────────────┘    │
│                          │                   │
│                 ┌────────┴────────┐         │
│                 ▼                 ▼          │
│          ┌──────────┐     ┌──────────┐      │
│          │   RDS    │     │   Redis  │      │
│          └──────────┘     └──────────┘      │
└─────────────────────────────────────────────┘
```

### Environments

| Environment | Purpose | URL | Notes |
|-------------|---------|-----|-------|
| Development | [Local testing] | [localhost:3000] | [Local Docker] |
| Staging | [Pre-prod testing] | [staging.example.com] | [Production-like] |
| Production | [Live traffic] | [api.example.com] | [Multi-AZ] |

### Configuration

[List key configuration variables and their purpose]

| Variable | Purpose | Default | Notes |
|----------|---------|---------|-------|
| `DATABASE_URL` | [DB connection] | [Required] | [No default] |
| `CACHE_TTL` | [Cache duration] | [3600] | [Seconds] |

---

## Operations

### Monitoring

| Metric | Description | Alert Threshold |
|--------|-------------|-----------------|
| [Request latency p95] | [API response time] | [> 500ms] |
| [Error rate] | [5xx responses / total] | [> 1%] |
| [Queue depth] | [Pending jobs] | [> 1000] |

### Health Checks

| Endpoint | Purpose | Expected Response |
|----------|---------|-------------------|
| `/health` | [Basic liveness] | [200 OK] |
| `/health/ready` | [Full readiness] | [200 with deps status] |

### Runbooks

- [Link to deployment runbook]
- [Link to incident response runbook]
- [Link to scaling runbook]

### On-Call

| Severity | Response Time | Escalation |
|----------|---------------|------------|
| P1 (Outage) | [15 minutes] | [Page on-call, then manager] |
| P2 (Degraded) | [1 hour] | [Page on-call] |
| P3 (Minor) | [Next business day] | [Slack channel] |

---

## Security

### Authentication

[How users/services authenticate]

### Authorization

[How permissions are determined]

### Data Classification

| Data Type | Classification | Handling |
|-----------|---------------|----------|
| [User PII] | [Confidential] | [Encrypted at rest, masked in logs] |
| [Public content] | [Public] | [No restrictions] |

### Compliance

[List relevant compliance requirements: SOC2, GDPR, HIPAA, etc.]

---

## Development

### Getting Started

[Quick start for local development]

```bash
# Clone and setup
git clone [repo-url]
cd [repo-name]
npm install

# Run locally
docker-compose up -d
npm run dev
```

### Testing

[How to run tests]

```bash
npm test              # Unit tests
npm run test:e2e      # Integration tests
npm run test:coverage # With coverage
```

### Contributing

[Link to contribution guidelines]

---

## Decision Log

[Link to ADRs or list key architectural decisions]

| ADR | Decision | Date |
|-----|----------|------|
| [ADR-001] | [Use PostgreSQL for persistence] | [2024-01-15] |
| [ADR-002] | [Adopt event-driven architecture] | [2024-02-20] |

---

## Appendix

### Glossary

| Term | Definition |
|------|------------|
| [Term 1] | [Definition] |
| [Term 2] | [Definition] |

### Related Documentation

- [Link to API docs]
- [Link to runbooks]
- [Link to design docs]

### Change Log

| Date | Version | Changes |
|------|---------|---------|
| [YYYY-MM-DD] | [1.0] | [Initial version] |
```

## Usage Instructions

1. Copy the template into a new file, typically `docs/system/[system-name].md`
2. Replace all bracketed placeholders with actual content
3. Remove sections that don't apply to your system
4. Add sections as needed for your specific context
5. Update the "Last Updated" date whenever you make changes
6. Review quarterly to ensure accuracy

## Section Selection Guide

**For a new/small system, start with:**
- Overview
- Architecture (system context)
- Deployment (basic)
- Development (getting started)

**Add as the system grows:**
- Component architecture
- Data architecture
- Integration points
- Operations (monitoring, runbooks)

**For production-critical systems, include:**
- All sections
- Detailed runbooks
- Security section
- On-call procedures

## Tips for Effective System Documentation

1. **Start with the overview**: A reader should understand what the system does after reading just the overview
2. **Use diagrams liberally**: Visual representations communicate architecture faster than text
3. **Keep it current**: Outdated docs are worse than no docs; schedule regular reviews
4. **Link, don't duplicate**: Reference other docs rather than copying content that will drift
5. **Include the "why"**: Not just what exists, but why decisions were made
6. **Test your getting started**: Have a new team member follow it and note confusion points

## Diagram Conventions

Use consistent shapes and arrows:

```
┌─────────┐  Rectangle: Component, Service, System
└─────────┘

┌─────────┐
│         │  Solid arrow: Synchronous call, direct dependency
└────┬────┘
     │
     ▼
┌─────────┐
└─────────┘

┌─────────┐
│         │  Dashed arrow: Asynchronous, event-based
└────┬────┘
     ┆
     ▼
┌─────────┐
└─────────┘

(─────────)  Cylinder: Database, persistent storage

[─────────]  Cloud: External service, third-party
```
