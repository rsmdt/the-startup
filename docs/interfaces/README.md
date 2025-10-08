# Interface Documentation

This directory contains API contracts, service integrations, and external system interface specifications.

## Purpose

Interface documentation captures **how the system communicates** with external services and systems. This includes:

- **API Contracts**: REST/GraphQL endpoint specifications
- **Service Integrations**: Third-party service integration patterns
- **Data Exchange**: Message formats and data transformation rules
- **External Dependencies**: External system requirements and constraints

## When to Create Interface Documentation

Create interface documentation when using `/s:analyze` to discover:
- External API integrations (payment gateways, auth providers, etc.)
- Service-to-service communication patterns
- Data exchange formats and protocols
- External system dependencies

## File Naming Convention

Use descriptive names that identify the external system or API:
- `stripe-payment-api.md` - Stripe payment integration contract
- `auth0-authentication.md` - Auth0 integration specification
- `sendgrid-email-service.md` - SendGrid email API contract
- `aws-s3-storage.md` - AWS S3 storage interface

## Structure

Each interface document should include:
1. **Overview**: What external system and integration purpose
2. **API Contract**: Endpoints, methods, request/response formats
3. **Authentication**: Auth mechanism and credentials management
4. **Error Handling**: Error codes and retry strategies
5. **Rate Limits**: Throttling and quota constraints
6. **Examples**: Sample requests and responses
7. **Related Patterns**: Links to implementation patterns in `docs/patterns/`

## Related Directories

- `../domain/` - Business rules that drive integration requirements
- `../patterns/` - Technical patterns for implementing integrations
