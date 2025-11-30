# Template: REST API Specification

## Purpose

Use this template to document REST API endpoints with consistent structure. Suitable for OpenAPI-style documentation, API design reviews, and developer reference guides.

## Template

```markdown
# [API Name] REST API

## Overview

[Brief description of the API's purpose and target consumers]

**Base URL:** `https://api.example.com/v1`

**Authentication:** [Bearer token / API key / OAuth 2.0]

**Content Type:** `application/json`

---

## Authentication

### [Auth Method Name]

[Description of authentication mechanism]

**Header Format:**
```
Authorization: Bearer <access_token>
```

**Error Responses:**
| Status | Code | Description |
|--------|------|-------------|
| 401 | UNAUTHORIZED | Missing or invalid token |
| 403 | FORBIDDEN | Valid token but insufficient permissions |

---

## Resources

### [Resource Name]

[Brief description of what this resource represents]

#### List [Resources]

Retrieves a paginated list of [resources].

**Endpoint:** `GET /[resources]`

**Query Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| page | integer | No | Page number (default: 1) |
| per_page | integer | No | Items per page (default: 20, max: 100) |
| sort | string | No | Sort field (prefix with - for descending) |
| [filter] | [type] | No | [Filter description] |

**Response:** `200 OK`
```json
{
  "data": [
    {
      "id": "uuid",
      "type": "[resource]",
      "attributes": {
        "[field]": "[value]"
      },
      "relationships": {
        "[related]": {
          "id": "uuid",
          "type": "[related_type]"
        }
      }
    }
  ],
  "meta": {
    "total": 100,
    "page": 1,
    "per_page": 20,
    "total_pages": 5
  },
  "links": {
    "self": "/[resources]?page=1",
    "next": "/[resources]?page=2",
    "last": "/[resources]?page=5"
  }
}
```

---

#### Get [Resource]

Retrieves a single [resource] by ID.

**Endpoint:** `GET /[resources]/{id}`

**Path Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| id | uuid | Unique identifier of the [resource] |

**Response:** `200 OK`
```json
{
  "data": {
    "id": "uuid",
    "type": "[resource]",
    "attributes": {
      "[field]": "[value]"
    }
  }
}
```

**Error Responses:**
| Status | Code | Description |
|--------|------|-------------|
| 404 | NOT_FOUND | [Resource] with specified ID not found |

---

#### Create [Resource]

Creates a new [resource].

**Endpoint:** `POST /[resources]`

**Request Body:**
```json
{
  "[field]": "[value]",
  "[required_field]": "[value]"
}
```

**Field Validation:**
| Field | Type | Required | Constraints |
|-------|------|----------|-------------|
| [field] | [type] | Yes/No | [min/max, format, enum values] |

**Response:** `201 Created`

**Headers:**
```
Location: /[resources]/{id}
```

```json
{
  "data": {
    "id": "uuid",
    "type": "[resource]",
    "attributes": {
      "[field]": "[value]"
    }
  }
}
```

**Error Responses:**
| Status | Code | Description |
|--------|------|-------------|
| 400 | BAD_REQUEST | Malformed request body |
| 422 | VALIDATION_ERROR | Field validation failed |
| 409 | CONFLICT | [Resource] already exists |

---

#### Update [Resource]

Partially updates an existing [resource].

**Endpoint:** `PATCH /[resources]/{id}`

**Request Body:**
```json
{
  "[field]": "[new_value]"
}
```

**Response:** `200 OK`
```json
{
  "data": {
    "id": "uuid",
    "type": "[resource]",
    "attributes": {
      "[field]": "[new_value]"
    }
  }
}
```

**Error Responses:**
| Status | Code | Description |
|--------|------|-------------|
| 404 | NOT_FOUND | [Resource] not found |
| 422 | VALIDATION_ERROR | Field validation failed |
| 409 | CONFLICT | Version conflict (if using optimistic locking) |

---

#### Delete [Resource]

Deletes a [resource].

**Endpoint:** `DELETE /[resources]/{id}`

**Response:** `204 No Content`

**Error Responses:**
| Status | Code | Description |
|--------|------|-------------|
| 404 | NOT_FOUND | [Resource] not found |
| 409 | CONFLICT | Cannot delete due to [dependencies] |

---

## Custom Actions

### [Action Name]

[Description of the action]

**Endpoint:** `POST /[resources]/{id}/[action]`

**Request Body:**
```json
{
  "[param]": "[value]"
}
```

**Response:** `200 OK` or `202 Accepted`
```json
{
  "data": {
    "[result_field]": "[value]"
  }
}
```

---

## Error Format

All errors follow this consistent format:

```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable error message",
    "details": [
      {
        "field": "field_name",
        "code": "FIELD_ERROR_CODE",
        "message": "Field-specific error message"
      }
    ],
    "request_id": "req_abc123",
    "documentation_url": "https://docs.example.com/errors/ERROR_CODE"
  }
}
```

### Error Codes

| Code | HTTP Status | Description | Resolution |
|------|-------------|-------------|------------|
| VALIDATION_ERROR | 422 | Request validation failed | Check field constraints |
| NOT_FOUND | 404 | Resource not found | Verify resource ID |
| UNAUTHORIZED | 401 | Authentication required | Provide valid credentials |
| FORBIDDEN | 403 | Permission denied | Check user permissions |
| CONFLICT | 409 | State conflict | Resolve conflict and retry |
| RATE_LIMITED | 429 | Too many requests | Wait and retry |
| INTERNAL_ERROR | 500 | Server error | Contact support |

---

## Rate Limiting

Rate limits are applied per API key/user.

**Headers:**
```
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 999
X-RateLimit-Reset: 1705320000
```

**Default Limits:**
| Tier | Requests/Hour | Burst |
|------|---------------|-------|
| Free | 100 | 10/min |
| Basic | 1,000 | 100/min |
| Pro | 10,000 | 1,000/min |

---

## Webhooks (if applicable)

### Event Types

| Event | Description | Payload |
|-------|-------------|---------|
| [resource].created | New [resource] created | Full [resource] object |
| [resource].updated | [Resource] modified | Changed fields |
| [resource].deleted | [Resource] removed | ID only |

### Payload Format

```json
{
  "id": "evt_abc123",
  "type": "[resource].created",
  "created_at": "2025-01-15T10:30:00Z",
  "data": {
    "object": { }
  }
}
```

---

## Versioning

Current version: `v1`

**Version Header:**
```
API-Version: 2025-01-15
```

**Deprecation Notice:**
When endpoints are deprecated, the following header is included:
```
Deprecation: true
Sunset: Sat, 15 Jan 2026 00:00:00 GMT
Link: <https://api.example.com/v2/[resource]>; rel="successor-version"
```

---

## SDKs and Examples

### cURL

```bash
curl -X GET "https://api.example.com/v1/[resources]" \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json"
```

### JavaScript

```javascript
const response = await fetch('https://api.example.com/v1/[resources]', {
  headers: {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  }
});
const data = await response.json();
```

### Python

```python
import requests

response = requests.get(
    'https://api.example.com/v1/[resources]',
    headers={
        'Authorization': f'Bearer {token}',
        'Content-Type': 'application/json'
    }
)
data = response.json()
```
```

## Usage Instructions

1. Copy the template above into your API documentation
2. Replace all `[bracketed]` placeholders with actual values
3. Remove sections that don't apply to your API
4. Add additional resources following the same pattern
5. Include real example values in JSON snippets
6. Document all error codes your API can return
7. Update rate limits to match your actual implementation

## Customization Tips

- **For internal APIs**: Simplify authentication section, skip SDK examples
- **For public APIs**: Add more detail to examples, include sandbox/testing info
- **For GraphQL companion**: Focus on REST-specific operations, link to GraphQL docs
- **For microservices**: Add service discovery and health check endpoints

## Checklist

Before publishing:
- [ ] All placeholder text replaced with actual content
- [ ] All endpoints tested and documented accurately
- [ ] Error codes match actual implementation
- [ ] Rate limits documented correctly
- [ ] Authentication flows verified
- [ ] Example requests/responses validated
- [ ] Version information current
- [ ] Contact/support information included
