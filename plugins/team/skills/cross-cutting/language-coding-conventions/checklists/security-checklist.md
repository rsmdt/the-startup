# Security Checklist

OWASP Top 10 aligned security checks for code review and implementation validation.

## A01: Broken Access Control

- [ ] All endpoints verify user authentication before processing
- [ ] Authorization checks occur on every request, not just UI hiding
- [ ] Users cannot access resources belonging to other users by modifying IDs
- [ ] Directory traversal attacks are prevented (no `../` in file paths)
- [ ] CORS is configured to allow only trusted origins
- [ ] Default deny: access requires explicit grant, not explicit denial

## A02: Cryptographic Failures

- [ ] Sensitive data is encrypted at rest (database, file storage)
- [ ] TLS 1.2+ is enforced for all data in transit
- [ ] Passwords use strong hashing (bcrypt, argon2) with appropriate cost
- [ ] No hardcoded secrets, API keys, or credentials in code
- [ ] Secrets are stored in environment variables or secret managers
- [ ] Cryptographic keys are rotated on a schedule
- [ ] No deprecated algorithms (MD5, SHA1 for security, DES, RC4)

## A03: Injection

- [ ] SQL queries use parameterized statements, never string concatenation
- [ ] NoSQL queries are properly escaped or use ODM/ORM
- [ ] OS commands are avoided; if necessary, inputs are validated against allowlist
- [ ] LDAP queries use proper escaping
- [ ] XML parsers disable external entity processing (XXE prevention)
- [ ] Template engines use auto-escaping by default

## A04: Insecure Design

- [ ] Threat modeling performed for new features
- [ ] Rate limiting implemented on authentication and sensitive endpoints
- [ ] Business logic includes abuse prevention (e.g., quantity limits)
- [ ] Fail securely: errors default to denied access
- [ ] Security requirements documented and testable
- [ ] Separation of tenants in multi-tenant systems

## A05: Security Misconfiguration

- [ ] Default credentials changed on all systems
- [ ] Unnecessary features and frameworks disabled
- [ ] Error messages do not expose stack traces or system details
- [ ] Security headers configured (CSP, X-Frame-Options, etc.)
- [ ] Cloud storage buckets are not publicly accessible by default
- [ ] Development/debug features disabled in production

## A06: Vulnerable Components

- [ ] Dependencies are tracked and regularly updated
- [ ] Known vulnerabilities are checked (npm audit, Snyk, Dependabot)
- [ ] Components are obtained from official sources
- [ ] Unused dependencies are removed
- [ ] Component versions are pinned (not floating)
- [ ] License compliance is verified

## A07: Identification and Authentication Failures

- [ ] Multi-factor authentication available for sensitive operations
- [ ] Password requirements meet current standards (length over complexity)
- [ ] Brute force protection: lockout or increasing delays
- [ ] Session tokens are regenerated after login
- [ ] Session tokens are invalidated on logout
- [ ] Sessions have appropriate timeouts
- [ ] Password reset tokens are single-use and time-limited

## A08: Software and Data Integrity Failures

- [ ] CI/CD pipelines verify integrity of dependencies
- [ ] Signed commits or verified sources for code
- [ ] Deserialization of untrusted data is avoided
- [ ] If deserialization is necessary, integrity checks are in place
- [ ] Update mechanisms verify signatures
- [ ] Code review required before deployment

## A09: Security Logging and Monitoring Failures

- [ ] Authentication events are logged (success and failure)
- [ ] Authorization failures are logged
- [ ] Input validation failures are logged
- [ ] Logs include sufficient context (timestamp, user, action, outcome)
- [ ] Logs do not contain sensitive data (passwords, tokens, PII)
- [ ] Alerting is configured for suspicious patterns
- [ ] Logs are stored securely and cannot be tampered with

## A10: Server-Side Request Forgery (SSRF)

- [ ] URLs from user input are validated against allowlist
- [ ] Internal network addresses are blocked (127.0.0.1, 10.x, 192.168.x)
- [ ] URL schemes are restricted (http/https only)
- [ ] Responses from external requests are not directly returned to users
- [ ] Network segmentation limits impact of SSRF

## Additional Checks

### Input Validation

- [ ] All user input is validated on the server side
- [ ] Validation uses allowlists, not blocklists
- [ ] File uploads validate type, size, and content
- [ ] JSON/XML payloads are validated against schema

### Output Encoding

- [ ] HTML output is encoded to prevent XSS
- [ ] JSON responses use proper content-type headers
- [ ] URLs are encoded when constructed dynamically
- [ ] Context-appropriate encoding (HTML, JavaScript, CSS, URL)

### API Security

- [ ] Authentication required for all non-public endpoints
- [ ] API keys are treated as secrets
- [ ] Rate limiting prevents abuse
- [ ] Response data is filtered to authorized fields only
- [ ] HTTP methods are restricted appropriately

## Usage Notes

1. Apply this checklist during code review
2. Use automated tools (SAST, DAST) to complement manual review
3. Not all items apply to every change - use judgment
4. Document any exceptions with rationale
5. Security testing should be part of the CI/CD pipeline
