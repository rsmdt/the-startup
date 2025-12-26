---
name: security-assessment
description: Vulnerability review, OWASP patterns, secure coding practices, and threat modeling approaches. Use when reviewing code security, designing secure systems, performing threat analysis, or validating security implementations.
---

# Security Assessment

A specialized skill for systematic security evaluation of code, architecture, and infrastructure. Combines threat modeling methodologies with practical code review techniques to identify and remediate security vulnerabilities.

## When to Use

- Reviewing code changes for security vulnerabilities
- Designing new features with security requirements
- Performing threat analysis on system architecture
- Validating security controls in infrastructure
- Assessing third-party integrations and dependencies
- Preparing for security audits or compliance reviews

## Threat Modeling with STRIDE

STRIDE is a threat modeling framework that categorizes threats by their nature. Apply this model during architecture review and feature design.

### Spoofing (Authentication)

Threat: Attacker pretends to be another user or system.

Questions to ask:
- How do we verify the identity of users and systems?
- Can authentication tokens be stolen or forged?
- Are there any authentication bypass paths?

Mitigations:
- Strong authentication mechanisms (MFA)
- Secure token generation and validation
- Session management with proper invalidation

### Tampering (Integrity)

Threat: Attacker modifies data in transit or at rest.

Questions to ask:
- Can data be modified between components?
- Are database records protected from unauthorized changes?
- Can configuration files be altered?

Mitigations:
- Input validation at all boundaries
- Cryptographic signatures for critical data
- Database integrity constraints and audit logs

### Repudiation (Non-repudiation)

Threat: Attacker denies performing an action.

Questions to ask:
- Can we prove who performed an action?
- Are audit logs tamper-resistant?
- Is there sufficient logging for forensics?

Mitigations:
- Comprehensive audit logging
- Secure, immutable log storage
- Digital signatures for critical operations

### Information Disclosure (Confidentiality)

Threat: Attacker gains access to sensitive information.

Questions to ask:
- What sensitive data exists in this system?
- How is data protected at rest and in transit?
- Are error messages revealing internal details?

Mitigations:
- Encryption for sensitive data (TLS, AES)
- Proper access controls and authorization
- Sanitized error messages

### Denial of Service (Availability)

Threat: Attacker makes the system unavailable.

Questions to ask:
- What resources can be exhausted?
- Are there rate limits on expensive operations?
- How does the system handle malformed input?

Mitigations:
- Rate limiting and throttling
- Input validation and size limits
- Resource quotas and timeouts

### Elevation of Privilege (Authorization)

Threat: Attacker gains higher privileges than intended.

Questions to ask:
- Can users access resources beyond their role?
- Are privilege checks performed consistently?
- Can administrative functions be accessed by regular users?

Mitigations:
- Principle of least privilege
- Role-based access control (RBAC)
- Authorization checks at every layer

## OWASP Top 10 Review Patterns

Systematic patterns for identifying the most critical web application security risks.

### A01: Broken Access Control

Review pattern:
1. Identify all endpoints and their expected access levels
2. Trace authorization logic from request to resource
3. Test for horizontal privilege escalation (accessing other users' data)
4. Test for vertical privilege escalation (accessing admin functions)
5. Verify CORS configuration restricts origins appropriately

Red flags:
- Authorization based on client-side state
- Direct object references without ownership verification
- Missing authorization checks on API endpoints

### A02: Cryptographic Failures

Review pattern:
1. Map all sensitive data flows (credentials, PII, financial)
2. Verify encryption at rest and in transit
3. Check for hardcoded secrets in code or configuration
4. Review cryptographic algorithm choices
5. Verify key management practices

Red flags:
- Sensitive data in logs or error messages
- Deprecated algorithms (MD5, SHA1, DES)
- Secrets in source control

### A03: Injection

Review pattern:
1. Identify all user input entry points
2. Trace input flow to database queries, OS commands, LDAP
3. Verify parameterized queries or proper escaping
4. Check for dynamic code execution (eval, exec)
5. Review XML parsing for XXE vulnerabilities

Red flags:
- String concatenation in queries
- User input in system commands
- Disabled XML external entity protection

### A04: Insecure Design

Review pattern:
1. Verify threat modeling was performed
2. Check for abuse case handling (rate limits, quantity limits)
3. Review business logic for security assumptions
4. Assess multi-tenancy isolation
5. Verify secure defaults

Red flags:
- No rate limiting on authentication
- Trust assumptions without verification
- Security as an afterthought

### A05: Security Misconfiguration

Review pattern:
1. Review default configurations for security settings
2. Check for unnecessary features or services
3. Verify error handling does not expose details
4. Review security headers (CSP, HSTS, X-Frame-Options)
5. Check cloud resource permissions

Red flags:
- Debug mode in production
- Default credentials unchanged
- Overly permissive cloud IAM policies

### A06: Vulnerable Components

Review pattern:
1. Inventory all dependencies and their versions
2. Check for known vulnerabilities (CVE databases)
3. Verify dependencies from trusted sources
4. Review for unused dependencies
5. Check for version pinning

Red flags:
- Unpinned dependencies
- Known critical vulnerabilities
- Dependencies from unofficial sources

### A07: Authentication Failures

Review pattern:
1. Review password policy enforcement
2. Check session management implementation
3. Verify brute force protection
4. Review token generation and validation
5. Check credential storage mechanisms

Red flags:
- Weak password requirements
- Sessions that do not invalidate on logout
- Predictable session tokens

### A08: Integrity Failures

Review pattern:
1. Review CI/CD pipeline security
2. Check for unsigned code or dependencies
3. Review deserialization of untrusted data
4. Verify update mechanism security
5. Check for code review requirements

Red flags:
- Deserialization without integrity checks
- Unsigned updates or dependencies
- No code review before deployment

### A09: Logging and Monitoring Failures

Review pattern:
1. Verify authentication events are logged
2. Check for authorization failure logging
3. Review log content for sensitive data
4. Verify log integrity protection
5. Check alerting configuration

Red flags:
- Missing authentication failure logs
- Sensitive data in logs
- No alerting on suspicious patterns

### A10: SSRF

Review pattern:
1. Identify all server-side URL fetching
2. Verify URL validation against allowlist
3. Check for internal network blocking
4. Review URL scheme restrictions
5. Verify response handling

Red flags:
- User-controlled URLs without validation
- Internal addresses not blocked
- Raw responses returned to users

## Secure Coding Practices

### Input Validation

Always validate on the server side, regardless of client validation:

```
function validateInput(input) {
  // Type validation
  if (typeof input !== 'string') {
    throw new ValidationError('Input must be a string');
  }

  // Length validation
  if (input.length > MAX_LENGTH) {
    throw new ValidationError('Input exceeds maximum length');
  }

  // Format validation (allowlist approach)
  if (!ALLOWED_PATTERN.test(input)) {
    throw new ValidationError('Input contains invalid characters');
  }

  return sanitize(input);
}
```

### Output Encoding

Context-appropriate encoding prevents injection:

- HTML context: Encode `<`, `>`, `&`, `"`, `'`
- JavaScript context: Use JSON.stringify or hex encoding
- URL context: Use encodeURIComponent
- SQL context: Use parameterized queries (never encode manually)

### Secrets Management

Never commit secrets to source control:

```
// Bad: Hardcoded secret
const apiKey = "sk-1234567890abcdef";

// Good: Environment variable
const apiKey = process.env.API_KEY;
if (!apiKey) {
  throw new ConfigurationError('API_KEY not configured');
}
```

### Error Handling for Security

Separate internal logging from user-facing errors:

```
try {
  await processRequest(data);
} catch (error) {
  // Log full details internally
  logger.error('Request processing failed', {
    error: error.message,
    stack: error.stack,
    userId: user.id,
    requestId: request.id
  });

  // Return generic message to user
  throw new UserError('Unable to process request');
}
```

## Infrastructure Security Considerations

### Network Security

- Segment networks to limit blast radius
- Use private subnets for internal services
- Implement network policies in Kubernetes
- Restrict egress traffic to known destinations

### Container Security

- Use minimal base images (distroless, Alpine)
- Run as non-root user
- Set read-only root filesystem where possible
- Scan images for vulnerabilities
- Limit container capabilities

### Secrets in Infrastructure

- Use secret management services (Vault, AWS Secrets Manager)
- Inject secrets as environment variables, not files
- Rotate secrets regularly
- Audit secret access

### Cloud IAM

- Apply principle of least privilege
- Use service accounts with minimal permissions
- Audit IAM policies regularly
- Avoid using root/admin accounts for routine operations

## Code Review Security Focus Areas

Priority areas for security-focused code review:

1. **Authentication and session management** - Token generation, validation, session lifecycle
2. **Authorization checks** - Access control at all layers
3. **Input handling** - All user input paths
4. **Data exposure** - Logs, errors, API responses
5. **Cryptography usage** - Algorithm selection, key management
6. **Third-party integrations** - Data sharing, authentication
7. **Error handling** - Information leakage, fail-secure behavior

## Best Practices

- Perform threat modeling before implementation
- Apply defense in depth (multiple security layers)
- Assume breach: design for detection and containment
- Automate security testing in CI/CD
- Keep dependencies updated and audited
- Document security decisions and accepted risks
- Train developers on secure coding practices

## References

- `checklists/security-review-checklist.md` - Comprehensive security review checklist
