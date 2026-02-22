# Secure Coding Practices

Practical secure coding patterns with examples for common vulnerability categories.

## Input Validation

Always validate on the server side, regardless of client validation:

```javascript
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

## Output Encoding

Context-appropriate encoding prevents injection:

- HTML context: Encode `<`, `>`, `&`, `"`, `'`
- JavaScript context: Use JSON.stringify or hex encoding
- URL context: Use encodeURIComponent
- SQL context: Use parameterized queries (never encode manually)

## Secrets Management

Never commit secrets to source control:

```javascript
// Bad: Hardcoded secret
const apiKey = "sk-1234567890abcdef";

// Good: Environment variable
const apiKey = process.env.API_KEY;
if (!apiKey) {
  throw new ConfigurationError('API_KEY not configured');
}
```

## Error Handling for Security

Separate internal logging from user-facing errors:

```javascript
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
