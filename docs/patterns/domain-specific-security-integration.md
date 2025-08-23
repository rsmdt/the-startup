# Domain-Specific Security Integration Pattern

## Context
Security practices need to be embedded across multiple Claude Code agents, but generic security advice is less effective than domain-specific guidance tailored to each agent's expertise area.

## Problem
- Generic security practices don't address domain-specific vulnerabilities
- Security guidance conflicts with pragmatic development workflows
- Duplication of generic security advice across multiple agents
- Tool-specific security recommendations reduce agent flexibility

## Solution
Distribute domain-specific security practices across relevant agents based on their expertise areas, using tool-agnostic language and actionable guidance.

### Security Distribution Strategy

**Backend Engineers**: API security, input validation, SQL injection prevention
```markdown
## Focus Areas
+ **API Security**: Validate all inputs at boundaries, use parameterized queries
+ **Authentication**: Implement token validation before business logic

## Approach  
+ Validate inputs at API boundaries - sanitize, type-check, size-limit
+ Use parameterized queries for all database operations

## Anti-Patterns to Avoid
+ Trusting client-provided data without server-side validation
+ Exposing database errors directly to API responses
```

**Frontend Engineers**: XSS prevention, client-side security limitations
```markdown
## Focus Areas
+ **Client Security**: XSS prevention, secure data handling, HTTPS enforcement

## Approach
+ Sanitize all user inputs before rendering to DOM
+ Use Content Security Policy headers to prevent script injection

## Anti-Patterns to Avoid  
+ Relying solely on client-side validation
+ Storing sensitive tokens in localStorage
```

**DevOps Engineers**: Secrets management, infrastructure security
```markdown
## Focus Areas
+ **Secrets Management**: Use environment variables or secret management systems
+ **Infrastructure Security**: Network segmentation, least-privilege access

## Approach
+ Implement secret management from deployment planning phase
+ Apply security scanning to all container images

## Anti-Patterns to Avoid
+ Storing secrets in deployment scripts or environment files
+ Running containers as root without necessity
```

### Tool-Agnostic Security Language

**Instead of**: "Use HashiCorp Vault for secret management"
**Use**: "Use dedicated secret management systems with automated rotation"

**Instead of**: "Configure AWS Security Groups"  
**Use**: "Implement network segmentation with security controls"

**Instead of**: "Use ESLint security plugins"
**Use**: "Apply security linting tools to detect common vulnerabilities"

## Examples

**Domain-Specific Integration** (the-data-engineer):
```markdown
## Focus Areas
- **Data Security**: Encrypt sensitive data at rest and in transit
- **Access Controls**: Role-based database access with query auditing
- **Privacy Compliance**: PII identification and data retention policies

## Approach
1. Implement data classification schemes for sensitive information
2. Apply encryption to sensitive data columns by default  
3. Audit all data access patterns and query activities
4. Design pipelines with data privacy controls built-in

## Anti-Patterns to Avoid
- Storing sensitive data in plain text format
- Using shared database accounts across applications
- Copying production data to development without sanitization
```

## When to Use
- Agent deals with user data, external systems, or sensitive operations
- Security practices need to be actionable within agent's domain expertise
- Generic security advice needs domain-specific adaptation
- Multiple agents need coordinated but specialized security guidance

## Benefits
- **Relevance**: Security practices directly applicable to agent's domain
- **Actionability**: Specific guidance agents can immediately apply
- **Coverage**: Comprehensive security across all relevant development areas
- **Consistency**: Coordinated security approach without duplication

## Implementation Notes
- Each security-relevant agent includes 1-2 domain-specific security practices
- Security practices embedded in existing agent structure (not separate section)
- Use tool-agnostic language that works across technology stacks
- Ensure security practices complement rather than conflict with pragmatic development
- Coordinate security handoff points between agents (e.g., frontend/backend API security)