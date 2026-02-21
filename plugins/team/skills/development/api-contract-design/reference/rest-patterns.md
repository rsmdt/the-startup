# REST API Patterns

Detailed REST patterns for resource modeling, HTTP methods, status codes, error handling, pagination, and filtering.

---

## Resource Modeling

Resources represent business entities. URLs should reflect the resource hierarchy.

```
GOOD:
GET    /users                    # List users
POST   /users                    # Create user
GET    /users/{id}               # Get user
PATCH  /users/{id}               # Partial update
DELETE /users/{id}               # Delete user
GET    /users/{id}/orders        # User's orders (sub-resource)

AVOID:
GET    /getUsers                 # Verbs in URLs
POST   /createNewUser            # Redundant verbs
GET    /user-list                # Inconsistent naming
POST   /users/{id}/delete        # Wrong HTTP method
```

## HTTP Method Semantics

| Method | Usage | Idempotent | Safe |
|--------|-------|------------|------|
| GET | Retrieve resource(s) | Yes | Yes |
| POST | Create resource, trigger action | No | No |
| PUT | Replace entire resource | Yes | No |
| PATCH | Partial update | Yes | No |
| DELETE | Remove resource | Yes | No |
| OPTIONS | CORS preflight, capability discovery | Yes | Yes |

## Status Code Selection

```
SUCCESS:
200 OK           - Successful GET, PUT, PATCH, DELETE
201 Created      - Successful POST (include Location header)
202 Accepted     - Async operation started
204 No Content   - Success with no response body

CLIENT ERRORS:
400 Bad Request  - Malformed request, validation failure
401 Unauthorized - Missing or invalid authentication
403 Forbidden    - Authenticated but not authorized
404 Not Found    - Resource doesn't exist
409 Conflict     - State conflict (duplicate, version mismatch)
422 Unprocessable- Semantically invalid (business rule violation)
429 Too Many     - Rate limit exceeded

SERVER ERRORS:
500 Internal     - Unexpected server error
502 Bad Gateway  - Upstream service failure
503 Unavailable  - Temporary overload or maintenance
504 Gateway Timeout - Upstream timeout
```

## Error Response Format

Standardize error responses across all endpoints:

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Request validation failed",
    "details": [
      {
        "field": "email",
        "code": "INVALID_FORMAT",
        "message": "Email must be a valid email address"
      }
    ],
    "requestId": "req_abc123",
    "timestamp": "2025-01-15T10:30:00Z",
    "documentation": "https://api.example.com/docs/errors#VALIDATION_ERROR"
  }
}
```

## Pagination Patterns

### Offset-Based (Simple, not for large datasets)

```
GET /users?offset=20&limit=10

Response:
{
  "data": [...],
  "pagination": {
    "total": 150,
    "offset": 20,
    "limit": 10,
    "hasMore": true
  }
}
```

### Cursor-Based (Recommended for large datasets)

```
GET /users?cursor=eyJpZCI6MTAwfQ&limit=10

Response:
{
  "data": [...],
  "pagination": {
    "nextCursor": "eyJpZCI6MTEwfQ",
    "prevCursor": "eyJpZCI6OTB9",
    "hasMore": true
  }
}
```

## Filtering and Sorting

```
FILTERING:
GET /users?status=active                    # Exact match
GET /users?created_after=2025-01-01         # Date range
GET /users?role=admin,moderator             # Multiple values
GET /users?search=john                      # Full-text search

SORTING:
GET /users?sort=created_at                  # Ascending (default)
GET /users?sort=-created_at                 # Descending (prefix -)
GET /users?sort=status,-created_at          # Multiple fields

FIELD SELECTION:
GET /users?fields=id,name,email             # Sparse fieldsets
GET /users?expand=organization              # Include related
```
