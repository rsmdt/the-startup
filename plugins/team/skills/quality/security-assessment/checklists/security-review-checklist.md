# Security Review Checklist

Comprehensive checklist for security-focused code review, architecture assessment, and infrastructure validation. Use this checklist systematically during security assessments.

## Threat Modeling

### STRIDE Analysis

- [ ] Spoofing threats identified and mitigated
- [ ] Tampering threats identified and mitigated
- [ ] Repudiation threats identified and mitigated
- [ ] Information disclosure threats identified and mitigated
- [ ] Denial of service threats identified and mitigated
- [ ] Elevation of privilege threats identified and mitigated

### Attack Surface

- [ ] All entry points documented (APIs, UIs, file uploads, webhooks)
- [ ] Trust boundaries identified between components
- [ ] Data flows mapped with sensitivity classifications
- [ ] Third-party integrations assessed for risk
- [ ] Unused features and endpoints disabled or removed

### Threat Actors

- [ ] Relevant threat actors identified (external, internal, privileged)
- [ ] Attacker capabilities and motivations considered
- [ ] High-value targets identified and prioritized
- [ ] Attack scenarios documented for critical paths

## Authentication

### Credential Handling

- [ ] Passwords hashed with strong algorithm (bcrypt, argon2, scrypt)
- [ ] Appropriate cost factor configured for hashing
- [ ] Password requirements enforce length over complexity
- [ ] No maximum password length restrictions below 64 characters
- [ ] Passwords not logged or exposed in error messages

### Session Management

- [ ] Session tokens generated with cryptographic randomness
- [ ] Session tokens have sufficient entropy (128+ bits)
- [ ] Sessions invalidated on logout
- [ ] Sessions invalidated on password change
- [ ] Session timeout configured appropriately
- [ ] Session tokens regenerated after authentication

### Multi-Factor Authentication

- [ ] MFA available for sensitive operations
- [ ] MFA recovery options are secure
- [ ] MFA bypass conditions documented and minimized
- [ ] TOTP implementation uses standard algorithms

### Brute Force Protection

- [ ] Account lockout or progressive delays implemented
- [ ] Rate limiting on authentication endpoints
- [ ] Lockout notifications sent to account owner
- [ ] CAPTCHA or proof-of-work for suspicious activity

## Authorization

### Access Control

- [ ] Authorization checks performed on every request
- [ ] Authorization enforced server-side, not client-side
- [ ] Default deny: access requires explicit grant
- [ ] Principle of least privilege applied to all roles
- [ ] Separation of duties for sensitive operations

### Resource Access

- [ ] Users cannot access other users' resources by ID manipulation
- [ ] Direct object references validated against ownership
- [ ] Batch operations validate authorization for all items
- [ ] File access restricted to authorized directories
- [ ] API responses filtered to authorized data only

### Administrative Functions

- [ ] Administrative endpoints require elevated authentication
- [ ] Administrative actions logged with user attribution
- [ ] Privilege escalation paths audited and minimized
- [ ] Role assignments require appropriate authorization

## Input Validation

### General Input Handling

- [ ] All input validated server-side
- [ ] Validation uses allowlist approach (not blocklist)
- [ ] Input length limits enforced
- [ ] Input type validation enforced
- [ ] Validation errors do not reveal internal details

### Injection Prevention

- [ ] SQL queries use parameterized statements
- [ ] NoSQL queries properly escaped or use ODM
- [ ] OS commands avoided; if necessary, input validated against allowlist
- [ ] LDAP queries properly escaped
- [ ] XPath queries properly escaped
- [ ] Template injection prevented

### File Handling

- [ ] File uploads validate type by content, not extension
- [ ] File uploads have size limits
- [ ] Uploaded files stored outside web root
- [ ] Uploaded filenames sanitized
- [ ] Path traversal prevented in file operations
- [ ] File downloads validate user authorization

### XML Processing

- [ ] External entity processing disabled (XXE prevention)
- [ ] DTD processing disabled if not required
- [ ] XML parser configured with resource limits

## Cryptography

### Data Protection

- [ ] Sensitive data encrypted at rest
- [ ] TLS 1.2+ enforced for data in transit
- [ ] Certificate validation not disabled
- [ ] Strong cipher suites configured

### Algorithm Selection

- [ ] No deprecated algorithms (MD5, SHA1, DES, RC4)
- [ ] Appropriate key sizes (AES-256, RSA-2048+)
- [ ] Cryptographic library from reputable source
- [ ] Random number generation uses secure source

### Key Management

- [ ] Keys stored in secure location (HSM, secret manager)
- [ ] Keys not hardcoded in source code
- [ ] Key rotation process defined
- [ ] Key backup and recovery procedures in place

## Error Handling and Logging

### Error Messages

- [ ] Error messages do not expose stack traces
- [ ] Error messages do not reveal system architecture
- [ ] Error messages do not expose sensitive data
- [ ] Generic messages for authentication failures

### Security Logging

- [ ] Authentication events logged (success and failure)
- [ ] Authorization failures logged
- [ ] Input validation failures logged
- [ ] Administrative actions logged
- [ ] Security-relevant configuration changes logged

### Log Protection

- [ ] Logs do not contain passwords or tokens
- [ ] Logs do not contain full credit card numbers
- [ ] Logs do not contain other sensitive PII inappropriately
- [ ] Log integrity protected (append-only, signed)
- [ ] Log access restricted and audited

### Alerting

- [ ] Alerts configured for repeated authentication failures
- [ ] Alerts configured for anomalous access patterns
- [ ] Alerts configured for privilege escalation attempts
- [ ] Alert fatigue considered (appropriate thresholds)

## API Security

### Endpoint Security

- [ ] All non-public endpoints require authentication
- [ ] CORS configured to allow only trusted origins
- [ ] HTTP methods restricted appropriately
- [ ] Sensitive operations use POST/PUT/DELETE, not GET
- [ ] API versioning strategy does not expose deprecated endpoints

### Rate Limiting

- [ ] Rate limits configured per user/IP
- [ ] Rate limits appropriate for endpoint cost
- [ ] Rate limit responses include retry-after header
- [ ] Rate limit bypass not possible through header manipulation

### Data Exposure

- [ ] API responses do not include unnecessary fields
- [ ] Pagination prevents mass data extraction
- [ ] Bulk endpoints have appropriate limits
- [ ] Error responses do not expose internal details

## Infrastructure Security

### Network Configuration

- [ ] Internal services not exposed to public internet
- [ ] Network segmentation limits lateral movement
- [ ] Firewall rules follow least privilege
- [ ] Egress traffic restricted to known destinations

### Container Security

- [ ] Base images from trusted sources
- [ ] Images scanned for vulnerabilities
- [ ] Containers run as non-root
- [ ] Container capabilities minimized
- [ ] Read-only root filesystem where possible

### Cloud Security

- [ ] IAM follows principle of least privilege
- [ ] Service accounts have minimal permissions
- [ ] Cloud storage buckets not publicly accessible
- [ ] Cloud audit logging enabled
- [ ] Infrastructure changes require review

### Secrets Management

- [ ] Secrets stored in secret manager (not config files)
- [ ] Secrets injected at runtime, not build time
- [ ] Secret access audited
- [ ] Secret rotation automated where possible

## Dependency Management

### Vulnerability Management

- [ ] Dependencies tracked in manifest file
- [ ] Dependency versions pinned
- [ ] Known vulnerabilities checked regularly
- [ ] Process for responding to critical vulnerabilities
- [ ] Automated scanning in CI/CD pipeline

### Supply Chain Security

- [ ] Dependencies from official sources
- [ ] Package integrity verified (checksums, signatures)
- [ ] Typosquatting risk considered
- [ ] Transitive dependencies reviewed

## CI/CD Security

### Pipeline Security

- [ ] CI/CD credentials stored securely
- [ ] Pipeline definitions reviewed for security
- [ ] Build artifacts signed or checksummed
- [ ] Deployment requires review/approval

### Code Review

- [ ] Security review required for sensitive changes
- [ ] Automated security scanning (SAST) in pipeline
- [ ] Dynamic security testing (DAST) for web applications
- [ ] Secrets scanning prevents credential commits

## Usage Notes

1. Apply this checklist during security-focused code reviews
2. Not all items apply to every review - use judgment based on context
3. Document exceptions with security rationale
4. Update this checklist as new threats emerge
5. Combine with automated tools for comprehensive coverage
6. Prioritize findings by risk (likelihood x impact)
7. Track remediation and verify fixes

## Risk Rating Guide

When documenting findings, use this risk rating:

- **Critical**: Immediate exploitation possible, severe impact
- **High**: Likely exploitation, significant impact
- **Medium**: Exploitation requires conditions, moderate impact
- **Low**: Exploitation difficult, limited impact
- **Informational**: Best practice deviation, no direct risk
