# [Service Name] Integration

> **Category:** External Interface
> **Service:** [Service Name]
> **Last Updated:** [Date]
> **Status:** [Active/Deprecated/Planned]

## Overview

**Service:** [Full service name]
**Provider:** [Company/organization]
**Purpose:** [What this integration accomplishes]
**Documentation:** [Link to official API docs]

## Authentication

### Method

[OAuth 2.0 / API Key / Basic Auth / JWT / etc.]

### Credentials Management

**Location:** [Where credentials are stored]
**Environment Variables:**
```bash
SERVICE_API_KEY=xxx
SERVICE_SECRET=xxx
SERVICE_ENDPOINT=https://...
```

**Rotation Policy:** [How often credentials change]

### Authentication Example

```[language]
// Example of authentication setup
[Code snippet]
```

## API Endpoints Used

### Endpoint 1: [Name]

**URL:** `[METHOD] /path/to/endpoint`
**Purpose:** [What this endpoint does]

**Request:**
```json
{
  "field1": "value",
  "field2": "value"
}
```

**Response:**
```json
{
  "status": "success",
  "data": { }
}
```

**Error Handling:**
- `400`: [How we handle]
- `401`: [How we handle]
- `500`: [How we handle]

### Endpoint 2: [Name]

**URL:** `[METHOD] /path/to/endpoint`
**Purpose:** [What this endpoint does]

[Same structure as above]

## Webhooks (if applicable)

### Webhook 1: [Event Name]

**Event Type:** `[event.type]`
**Trigger:** [When this fires]
**URL:** `[Your webhook endpoint]`

**Payload:**
```json
{
  "event": "type",
  "data": { }
}
```

**Signature Verification:**
```[language]
// How to verify webhook authenticity
[Code snippet]
```

**Handling:**
```[language]
// How we process this webhook
[Code snippet]
```

## Rate Limits

- **Requests per second:** [Limit]
- **Requests per day:** [Limit]
- **Burst limit:** [Limit]

**Handling Strategy:** [How we respect limits]
```[language]
// Rate limiting implementation
[Code snippet]
```

## Data Mapping

### Our Model → Service Model

| Our Field | Service Field | Transformation |
|-----------|---------------|----------------|
| `userId` | `external_id` | String conversion |
| `email` | `email_address` | Direct mapping |
| `amount` | `total_cents` | Multiply by 100 |

### Service Model → Our Model

| Service Field | Our Field | Transformation |
|---------------|-----------|----------------|
| `id` | `externalId` | Direct mapping |
| `status` | `state` | Enum mapping |
| `created_at` | `createdAt` | ISO 8601 parse |

## Error Handling

### Common Errors

**Error 1: [Name/Code]**
- **Cause:** [What triggers this]
- **Recovery:** [How we handle it]
- **Retry:** [Yes/No, strategy]

**Error 2: [Name/Code]**
- **Cause:** [What triggers this]
- **Recovery:** [How we handle it]
- **Retry:** [Yes/No, strategy]

### Retry Strategy

```[language]
// Exponential backoff implementation
[Code snippet]
```

## Testing

### Test Credentials

**Sandbox URL:** `https://sandbox.service.com`
**Test API Key:** `[Where to get it]`

### Mock Server

**Location:** `tests/mocks/[service]-mock.ts`
**Usage:**
```[language]
// How to use mock in tests
[Code snippet]
```

### Integration Tests

```[language]
// Example integration test
[Code snippet]
```

## Monitoring

### Health Checks

**Endpoint:** `[Service status endpoint]`
**Frequency:** [How often we check]

### Metrics to Track

- Request success rate
- Response time (p50, p95, p99)
- Error rate by type
- Rate limit proximity

### Alerts

- **Critical:** [Conditions that trigger urgent alerts]
- **Warning:** [Conditions that trigger warnings]

## Security Considerations

- [Security consideration 1]
- [Security consideration 2]
- [Security consideration 3]

## Compliance

**Data Handling:**
- PII fields: [List]
- Retention policy: [Duration]
- Geographic restrictions: [Any]

**Regulations:**
- GDPR: [Compliance notes]
- CCPA: [Compliance notes]
- Other: [Relevant regulations]

## Cost Considerations

**Pricing Model:** [How service charges]
**Cost per request:** [Estimate]
**Monthly estimate:** [Based on usage]

## Migration/Upgrade Path

**Current Version:** [Version]
**Upgrade Available:** [Yes/No, version]
**Breaking Changes:** [List if applicable]

**Migration Steps:**
1. [Step 1]
2. [Step 2]
3. [Step 3]

## Related Documentation

- **Patterns:** [Pattern Doc](../patterns/doc.md) - [How we use this service]
- **Domain:** [Domain Doc](../domain/doc.md) - [Business rules related to this]
- **Specifications:** [Spec](../specs/NNN-name/SDD.md) - [Implementation details]

## External Resources

- [Official API documentation]
- [Status page]
- [Developer community/forum]
- [SDK/library used]

## Contact

**Support:** [How to get help]
**Account Manager:** [If applicable]
**Escalation:** [Critical issue contact]

## Version History

| Date | Change | Author |
|------|--------|--------|
| [Date] | Initial integration | [Name/Tool] |
| [Date] | Updated to v2 API | [Name/Tool] |
