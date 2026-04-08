---
title: "Query Validation API"
status: completed
---

# Solution Design

## Architecture
REST API with Express.js, following existing controller patterns in AGENTS.md.

## Components
- ValidationController: HTTP endpoint handling
- SQLParser: Syntax validation and structure extraction using JSQLParser
- InjectionDetector: Pattern matching for SQL injection
- ComplexityScorer: Calculates query complexity metrics
- RateLimiter: Token bucket rate limiting per API key

## ADRs
- ADR-1: Use JSQLParser (existing dependency) for parsing
- ADR-2: Rate limiting at application level, not reverse proxy
