# API Versioning and Authentication Patterns

Versioning strategies and authentication/authorization patterns for APIs.

---

## Versioning Strategies

### URL Path Versioning

```
GET /v1/users
GET /v2/users

PROS:
- Explicit and visible
- Easy to route in infrastructure
- Clear in logs and monitoring

CONS:
- URL pollution
- Harder to deprecate gracefully
```

### Header Versioning

```
GET /users
Accept: application/vnd.api+json; version=2

PROS:
- Clean URLs
- Content negotiation friendly
- Easier partial versioning

CONS:
- Less visible
- Harder to test in browser
```

### Query Parameter Versioning

```
GET /users?api-version=2025-01-15

PROS:
- Easy to test
- Visible in URLs
- Date-based versions are intuitive

CONS:
- Clutters query strings
- Easy to forget
```

### Recommended: Dual Approach

```
1. Major versions in URL path: /v1/, /v2/
2. Minor versions via header: API-Version: 2025-01-15
3. Default to latest minor within major
4. Sunset headers for deprecation warnings
```

---

## Authentication Patterns

### API Keys

```
USAGE: Server-to-server, rate limiting, analytics
TRANSPORT: Header (Authorization: ApiKey xxx) or query param

SECURITY:
- Rotate keys regularly
- Different keys for environments
- Scope keys to specific operations
- Never expose in client-side code
```

### OAuth 2.0 / OIDC

```
FLOWS:
- Authorization Code + PKCE: Web apps, mobile apps
- Client Credentials: Server-to-server
- Device Code: CLI tools, smart TVs

TOKEN HANDLING:
- Short-lived access tokens (15-60 min)
- Refresh tokens for session extension
- Token introspection for validation
- Token revocation endpoint
```

### JWT Best Practices

```
CLAIMS:
{
  "iss": "https://auth.example.com",
  "sub": "user_123",
  "aud": "api.example.com",
  "exp": 1705320000,
  "iat": 1705316400,
  "scope": "read:users write:users"
}

SECURITY:
- Use asymmetric keys (RS256, ES256)
- Validate all claims
- Check token expiration
- Verify audience matches
- Keep tokens stateless when possible
```
