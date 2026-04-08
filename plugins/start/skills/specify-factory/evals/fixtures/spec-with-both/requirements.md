---
title: "Query Validation API"
status: completed
---

# Product Requirements

## Problem
Users need to validate SQL queries before executing them against production databases.

## Features

### Must Have
- F1: POST endpoint to validate SQL syntax
- F2: Detect SQL injection patterns and reject dangerous queries
- F3: Extract table and column references from valid queries
- F4: Return structured JSON response with validation results

### Should Have
- F5: Compute query complexity score (joins, subqueries, aggregations)
- F6: Rate limiting to prevent abuse

## Acceptance Criteria
- AC-1: Valid SELECT queries return 200 with parsed structure
- AC-2: Injection attempts return 400 with clear error
- AC-3: Empty/null input returns 422 with validation message
- AC-4: Rate-limited requests return 429 with retry-after header
