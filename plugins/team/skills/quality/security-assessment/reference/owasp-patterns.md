# OWASP Top 10 Review Patterns

Systematic patterns for identifying the most critical web application security risks.

## A01: Broken Access Control

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

## A02: Cryptographic Failures

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

## A03: Injection

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

## A04: Insecure Design

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

## A05: Security Misconfiguration

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

## A06: Vulnerable Components

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

## A07: Authentication Failures

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

## A08: Integrity Failures

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

## A09: Logging and Monitoring Failures

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

## A10: SSRF

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
