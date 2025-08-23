# Architectural Asset Management Pattern

## Context
Architecture agents (the-software-architect, the-staff-engineer, the-business-analyst) design solutions and define contracts that need to be captured as reusable documentation. These architectural decisions should be preserved for future reference and reuse.

## Problem
- Architectural patterns get lost when only implemented without documentation
- Service interfaces are redefined repeatedly instead of reusing existing contracts
- No systematic approach for capturing reusable design decisions
- Duplication of similar patterns and interfaces across projects
- Architectural knowledge exists only in individual agent responses

## Solution
Architecture agents systematically create and maintain architectural assets in two categories:

### 1. Patterns (Reusable Design Solutions)
**Create at**: `docs/patterns/[descriptive-name].md`
**When to Create**: Solution will be reused across multiple features or components
**Required Sections**:
- **Context**: When and why this pattern applies
- **Problem**: What specific problem this pattern solves  
- **Solution**: How the pattern addresses the problem
- **Implementation Example**: Concrete code or configuration example

### 2. Interfaces (Service Contracts)
**Create at**: `docs/interfaces/[service-name].md`  
**When to Create**: Defining contracts between services, systems, or external APIs
**Required Sections**:
- **Authentication**: How to authenticate with the service
- **Rate Limits**: Request limits and throttling policies
- **Data Formats**: Request/response schemas and content types
- **Examples**: Sample requests and responses

### 3. Deduplication Protocol
**Before Creating**: Always search existing `docs/patterns/` and `docs/interfaces/`
**Prefer Enhancement**: Update existing documentation over creating duplicates
**Naming Convention**: Use descriptive, searchable kebab-case names
**Cross-Reference**: Link between related patterns and interfaces

## Examples

**Pattern Creation** (the-software-architect):
```markdown
# API Error Handling Pattern

## Context
APIs need consistent error responses across all services to enable predictable client-side error handling.

## Problem  
Different services return errors in different formats, making client-side error handling complex and inconsistent.

## Solution
Standardize all API errors using this response format:
```json
{
  "error": {
    "code": "VALIDATION_FAILED",
    "message": "Request validation failed",
    "details": {
      "field": "email",
      "reason": "Invalid email format"
    }
  }
}
```

## Implementation Example
```javascript
// Error handler middleware
function errorHandler(err, req, res, next) {
  const standardError = {
    error: {
      code: err.code || 'INTERNAL_ERROR',
      message: err.message,
      details: err.details || {}
    }
  };
  res.status(err.statusCode || 500).json(standardError);
}
```

**Interface Creation** (the-staff-engineer):
```markdown
# User Authentication Service Interface

## Authentication
Bearer token required in Authorization header:
```
Authorization: Bearer <jwt-token>
```

## Rate Limits
- 100 requests per minute per API key
- 429 status code returned when limit exceeded
- Retry after 60 seconds

## Data Formats
All requests/responses use `application/json`

### Login Request
```json
{
  "email": "user@example.com",
  "password": "securePassword123"
}
```

### Login Response  
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_at": "2024-01-15T10:30:00Z",
  "user": {
    "id": "usr_123",
    "email": "user@example.com"
  }
}
```

## Examples
[Additional request/response examples for common operations]
```

**Enhancement vs Creation** (the-business-analyst):
```markdown
# Before creating new pattern, search existing:
Found: docs/patterns/api-validation.md (similar to what I need)

# Instead of creating docs/patterns/input-validation.md:
# Update existing docs/patterns/api-validation.md with additional validation rules
# Add cross-reference to related security patterns
```

## When to Use
- Architecture agent designs a solution that could be reused elsewhere
- Service contracts need to be defined for inter-system communication  
- Design decisions need to be captured for future reference
- Multiple teams need consistent implementation guidance

## Benefits
- **Reusability**: Architectural decisions captured for reuse across features
- **Consistency**: Standard approach for similar problems across codebase
- **Knowledge Preservation**: Design rationale and context preserved
- **Reduced Duplication**: Existing assets enhanced rather than recreated

## Implementation Notes
- Search existing patterns/interfaces BEFORE creating new ones
- Use descriptive, searchable names (avoid generic terms like "helper" or "util")
- Include concrete examples, not just abstract descriptions
- Cross-reference related patterns and interfaces
- Keep documentation up-to-date as implementations evolve
- Focus on reusable aspects, not project-specific details