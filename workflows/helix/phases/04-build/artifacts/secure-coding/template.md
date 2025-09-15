# Secure Coding Checklist

**Project**: [Project Name]
**Date**: [Creation Date]
**Developer**: [Name]

## Input Validation and Sanitization

### Server-Side Validation
- [ ] All user inputs validated on server side
- [ ] Input validation uses allowlist approach
- [ ] Data type, length, and format validation implemented
- [ ] Business logic validation applied
- [ ] File upload restrictions and validation

### Output Encoding
- [ ] HTML output encoding implemented
- [ ] JavaScript output encoding for dynamic content
- [ ] URL encoding for URL parameters
- [ ] SQL parameterized queries used exclusively
- [ ] LDAP injection prevention measures

### Error Handling
- [ ] Generic error messages for users
- [ ] Detailed errors logged securely (not exposed)
- [ ] Stack traces and system information hidden
- [ ] Failed operations don't reveal sensitive information

## Authentication and Session Management

### Password Security
- [ ] Passwords never stored in plain text
- [ ] Strong password hashing (bcrypt, Argon2)
- [ ] Salt used for password hashing
- [ ] Password complexity requirements enforced
- [ ] Password history prevention implemented

### Session Security
- [ ] Secure session token generation
- [ ] HttpOnly flag set on session cookies
- [ ] Secure flag set for HTTPS-only cookies
- [ ] SameSite attribute configured appropriately
- [ ] Session timeout implemented
- [ ] Session invalidation on logout
- [ ] Concurrent session limits enforced

### Multi-Factor Authentication
- [ ] MFA implemented for sensitive operations
- [ ] TOTP/HOTP algorithms used correctly
- [ ] Backup codes provided and secured
- [ ] MFA bypass protections implemented

## Access Control

### Authorization Checks
- [ ] Authorization checked on every request
- [ ] Default deny policy implemented
- [ ] Least privilege principle applied
- [ ] Resource-level access controls implemented
- [ ] Indirect object references used (no direct DB IDs)

### Role-Based Access Control
- [ ] Roles and permissions properly defined
- [ ] Role inheritance implemented correctly
- [ ] Permission checks at appropriate granularity
- [ ] Administrative functions properly protected

## Data Protection

### Encryption Implementation
- [ ] Strong encryption algorithms used (AES-256)
- [ ] Proper key management implemented
- [ ] Encryption keys stored securely
- [ ] Key rotation procedures implemented
- [ ] Encrypted data properly initialized (IV/nonce)

### Sensitive Data Handling
- [ ] Sensitive data identified and classified
- [ ] Data minimization principles applied
- [ ] Secure data transmission (TLS 1.3+)
- [ ] Sensitive data not logged or cached
- [ ] Secure data disposal implemented

### Database Security
- [ ] Database connections encrypted
- [ ] Principle of least privilege for DB access
- [ ] Stored procedures used where appropriate
- [ ] Database credentials secured
- [ ] SQL injection prevention verified

## Communication Security

### TLS/HTTPS Configuration
- [ ] TLS 1.3 or TLS 1.2 minimum
- [ ] Strong cipher suites configured
- [ ] Certificate validation implemented
- [ ] HSTS headers configured
- [ ] Certificate pinning implemented (where appropriate)

### API Security
- [ ] API authentication implemented
- [ ] Rate limiting configured
- [ ] Input validation for API parameters
- [ ] Proper HTTP status codes returned
- [ ] CORS configured securely

## Logging and Monitoring

### Security Logging
- [ ] Authentication events logged
- [ ] Authorization failures logged
- [ ] Administrative actions logged
- [ ] Security-relevant events logged
- [ ] Log integrity protection implemented

### Monitoring Implementation
- [ ] Real-time security event monitoring
- [ ] Anomaly detection configured
- [ ] Alert thresholds properly set
- [ ] Incident response triggers configured

## Code Quality and Security

### Secure Coding Practices
- [ ] Code review completed with security focus
- [ ] Static analysis tools results addressed
- [ ] Third-party dependencies scanned for vulnerabilities
- [ ] Security-focused unit tests written
- [ ] Integration tests include security scenarios

### Configuration Security
- [ ] Secure defaults configured
- [ ] Unnecessary features disabled
- [ ] Debug information disabled in production
- [ ] Configuration files secured
- [ ] Environment-specific configurations isolated

## Secrets Management

### API Keys and Secrets
- [ ] No hardcoded secrets in code
- [ ] External secret management system used
- [ ] Secret rotation procedures implemented
- [ ] Secrets encrypted at rest
- [ ] Access to secrets logged and monitored

### Certificates and Keys
- [ ] Private keys properly protected
- [ ] Certificate validation implemented
- [ ] Key storage follows security best practices
- [ ] Certificate lifecycle management implemented

## Cloud Security (if applicable)

### Cloud Configuration
- [ ] Cloud security best practices followed
- [ ] IAM roles and policies properly configured
- [ ] Network security groups configured
- [ ] Storage encryption enabled
- [ ] Logging and monitoring enabled

### Container Security (if applicable)
- [ ] Base images from trusted sources
- [ ] Container images scanned for vulnerabilities
- [ ] Runtime security configured
- [ ] Secrets management for containers
- [ ] Network policies configured

## Compliance Requirements

### Data Privacy
- [ ] GDPR requirements implemented (if applicable)
- [ ] Data subject rights functionality implemented
- [ ] Privacy by design principles followed
- [ ] Consent management implemented

### Industry Standards
- [ ] Industry-specific requirements addressed
- [ ] Regulatory compliance validated
- [ ] Audit trail requirements met
- [ ] Documentation requirements fulfilled

## Security Testing Integration

### Automated Security Tests
- [ ] SAST integration in CI/CD pipeline
- [ ] Dependency vulnerability scanning enabled
- [ ] Security unit tests implemented
- [ ] Automated security regression tests

### Manual Security Testing
- [ ] Security code review completed
- [ ] Manual security testing performed
- [ ] Penetration testing preparation completed

## Final Security Checklist

### Pre-Deployment Security Review
- [ ] All security requirements implemented
- [ ] Security test suite passes completely
- [ ] Vulnerability scan results acceptable
- [ ] Security architecture review completed
- [ ] Security documentation updated

### Security Sign-off
- [ ] Developer security checklist completed
- [ ] Security champion review completed
- [ ] Security team approval obtained
- [ ] Technical lead sign-off received

## Approval and Sign-off

| Role | Name | Date |
|------|------|------|
| Developer | [Name] | |
| Security Champion | [Name] | |
| Technical Lead | [Name] | |

---
*Document Version: 1.0*