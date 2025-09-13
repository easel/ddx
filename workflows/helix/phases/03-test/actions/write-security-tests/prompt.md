# Write Security Tests Prompt

Create comprehensive security tests that validate authentication, authorization, data protection, and vulnerability prevention mechanisms. These tests ensure the application is resistant to common attacks and follows security best practices.

## Purpose

Security tests verify:
- Authentication mechanisms work correctly
- Authorization boundaries are enforced
- Input validation prevents injection attacks
- Sensitive data is properly protected
- Security headers and policies are implemented
- Known vulnerabilities are addressed

## Security Test Categories

### Authentication Testing
Verify identity verification mechanisms:
- Login/logout functionality
- Password strength requirements
- Multi-factor authentication (MFA)
- Session management
- Account lockout policies
- Password reset flows

### Authorization Testing
Ensure proper access control:
- Role-based access control (RBAC)
- Resource-level permissions
- API endpoint authorization
- Cross-tenant data isolation
- Privilege escalation prevention
- Admin function protection

### Input Validation Testing
Prevent injection attacks:
- SQL injection prevention
- Cross-Site Scripting (XSS) prevention
- Command injection prevention
- Path traversal prevention
- XML/JSON injection prevention
- File upload validation

### Data Protection Testing
Verify sensitive data handling:
- Encryption at rest
- Encryption in transit
- PII masking/redaction
- Secure cookie handling
- Token security
- Secret management

## Authentication Tests

### Login Security
```javascript
describe('Authentication Security', () => {
  it('should enforce password complexity requirements', async () => {
    const weakPasswords = [
      'password',
      '12345678',
      'qwerty123',
      'Password',  // No numbers
      'Pass123',   // Too short
    ];

    for (const password of weakPasswords) {
      const response = await request(app)
        .post('/api/auth/register')
        .send({
          email: 'test@example.com',
          password: password
        });

      expect(response.status).toBe(400);
      expect(response.body.error).toContain('password');
    }
  });

  it('should prevent brute force attacks', async () => {
    const attempts = [];

    // Make 5 failed login attempts
    for (let i = 0; i < 5; i++) {
      attempts.push(
        request(app)
          .post('/api/auth/login')
          .send({
            email: 'user@example.com',
            password: 'wrongpassword'
          })
      );
    }

    const responses = await Promise.all(attempts);

    // After 5 attempts, account should be locked
    const lastResponse = responses[4];
    expect(lastResponse.status).toBe(429); // Too Many Requests
    expect(lastResponse.body.error).toContain('locked');
  });

  it('should invalidate sessions on logout', async () => {
    // Login
    const loginRes = await request(app)
      .post('/api/auth/login')
      .send({
        email: 'user@example.com',
        password: 'ValidPass123!'
      });

    const token = loginRes.body.token;

    // Logout
    await request(app)
      .post('/api/auth/logout')
      .set('Authorization', `Bearer ${token}`);

    // Try to use the token
    const response = await request(app)
      .get('/api/profile')
      .set('Authorization', `Bearer ${token}`);

    expect(response.status).toBe(401);
  });
});
```

### Session Management
```javascript
describe('Session Security', () => {
  it('should expire sessions after inactivity', async () => {
    const token = await getAuthToken();

    // Fast-forward time
    jest.advanceTimersByTime(31 * 60 * 1000); // 31 minutes

    const response = await request(app)
      .get('/api/profile')
      .set('Authorization', `Bearer ${token}`);

    expect(response.status).toBe(401);
    expect(response.body.error).toContain('expired');
  });

  it('should prevent session fixation', async () => {
    const oldSessionId = await getSessionId();

    // Login
    await request(app)
      .post('/api/auth/login')
      .set('Cookie', `sessionId=${oldSessionId}`)
      .send({
        email: 'user@example.com',
        password: 'ValidPass123!'
      });

    const newSessionId = await getSessionId();

    // Session ID should change after login
    expect(newSessionId).not.toBe(oldSessionId);
  });
});
```

## Authorization Tests

### Access Control Testing
```javascript
describe('Authorization Security', () => {
  it('should prevent unauthorized access to admin endpoints', async () => {
    const userToken = await getUserToken(); // Regular user token

    const response = await request(app)
      .get('/api/admin/users')
      .set('Authorization', `Bearer ${userToken}`);

    expect(response.status).toBe(403); // Forbidden
    expect(response.body.error).toContain('insufficient permissions');
  });

  it('should prevent cross-tenant data access', async () => {
    const tenant1Token = await getTenantToken('tenant1');
    const tenant2DataId = 'data-from-tenant2';

    const response = await request(app)
      .get(`/api/data/${tenant2DataId}`)
      .set('Authorization', `Bearer ${tenant1Token}`);

    expect(response.status).toBe(404); // Not found (hidden for security)
  });

  it('should enforce field-level permissions', async () => {
    const limitedUserToken = await getLimitedUserToken();

    const response = await request(app)
      .get('/api/users/123')
      .set('Authorization', `Bearer ${limitedUserToken}`);

    expect(response.status).toBe(200);
    expect(response.body).not.toHaveProperty('ssn');
    expect(response.body).not.toHaveProperty('salary');
    expect(response.body).toHaveProperty('name'); // Allowed field
  });
});
```

### Privilege Escalation Testing
```javascript
describe('Privilege Escalation Prevention', () => {
  it('should prevent role manipulation', async () => {
    const userToken = await getUserToken();

    const response = await request(app)
      .patch('/api/profile')
      .set('Authorization', `Bearer ${userToken}`)
      .send({
        name: 'Updated Name',
        role: 'admin' // Attempting to escalate
      });

    expect(response.status).toBe(200);

    // Verify role was not changed
    const profile = await request(app)
      .get('/api/profile')
      .set('Authorization', `Bearer ${userToken}`);

    expect(profile.body.role).toBe('user');
    expect(profile.body.role).not.toBe('admin');
  });
});
```

## Input Validation Tests

### SQL Injection Prevention
```javascript
describe('SQL Injection Prevention', () => {
  it('should prevent SQL injection in search queries', async () => {
    const maliciousInputs = [
      "'; DROP TABLE users; --",
      "1' OR '1'='1",
      "admin'--",
      "' UNION SELECT * FROM passwords--"
    ];

    for (const input of maliciousInputs) {
      const response = await request(app)
        .get('/api/search')
        .query({ q: input });

      expect(response.status).not.toBe(500); // No server error
      expect(response.body).not.toContain('SQL');
      expect(response.body).not.toContain('syntax error');
    }

    // Verify table still exists
    const validResponse = await request(app)
      .get('/api/users')
      .set('Authorization', `Bearer ${adminToken}`);

    expect(validResponse.status).toBe(200);
  });

  it('should use parameterized queries', async () => {
    // This test checks the implementation
    const queryFunction = require('../src/database').query;
    const queryCall = jest.spyOn(queryFunction, 'query');

    await request(app)
      .get('/api/users/123');

    // Verify parameterized query was used
    expect(queryCall).toHaveBeenCalledWith(
      expect.stringContaining('$1'),
      expect.arrayContaining([123])
    );
  });
});
```

### XSS Prevention
```javascript
describe('XSS Prevention', () => {
  it('should sanitize user input in responses', async () => {
    const xssPayloads = [
      '<script>alert("XSS")</script>',
      '<img src=x onerror="alert(\'XSS\')">',
      'javascript:alert("XSS")',
      '<iframe src="javascript:alert(\'XSS\')"></iframe>'
    ];

    for (const payload of xssPayloads) {
      const response = await request(app)
        .post('/api/comments')
        .send({
          content: payload
        });

      // Check response doesn't contain executable script
      expect(response.body.content).not.toContain('<script>');
      expect(response.body.content).not.toContain('javascript:');
      expect(response.body.content).not.toContain('onerror=');
    }
  });

  it('should set proper Content-Security-Policy headers', async () => {
    const response = await request(app).get('/');

    expect(response.headers['content-security-policy']).toContain("default-src 'self'");
    expect(response.headers['content-security-policy']).toContain("script-src 'self'");
    expect(response.headers['content-security-policy']).not.toContain("'unsafe-inline'");
  });
});
```

### File Upload Security
```javascript
describe('File Upload Security', () => {
  it('should validate file types', async () => {
    const maliciousFiles = [
      { name: 'shell.php', content: '<?php system($_GET["cmd"]); ?>' },
      { name: 'script.exe', content: Buffer.from([0x4D, 0x5A]) }, // EXE header
      { name: 'image.jpg.html', content: '<script>alert("XSS")</script>' }
    ];

    for (const file of maliciousFiles) {
      const response = await request(app)
        .post('/api/upload')
        .attach('file', Buffer.from(file.content), file.name);

      expect(response.status).toBe(400);
      expect(response.body.error).toContain('file type not allowed');
    }
  });

  it('should prevent path traversal in file operations', async () => {
    const pathTraversalInputs = [
      '../../../etc/passwd',
      '..\\..\\..\\windows\\system32\\config\\sam',
      'uploads/../../../config/database.yml'
    ];

    for (const input of pathTraversalInputs) {
      const response = await request(app)
        .get('/api/files')
        .query({ path: input });

      expect(response.status).toBe(400);
      expect(response.body.error).toContain('invalid path');
    }
  });
});
```

## Data Protection Tests

### Encryption Testing
```javascript
describe('Data Encryption', () => {
  it('should encrypt sensitive data at rest', async () => {
    // Create user with sensitive data
    await request(app)
      .post('/api/users')
      .send({
        name: 'Test User',
        ssn: '123-45-6789',
        creditCard: '4111111111111111'
      });

    // Query database directly
    const userData = await db.query('SELECT * FROM users WHERE name = $1', ['Test User']);

    // Sensitive fields should be encrypted
    expect(userData.ssn).not.toBe('123-45-6789');
    expect(userData.creditCard).not.toBe('4111111111111111');
    expect(userData.ssn).toMatch(/^encrypted:/);
  });

  it('should use HTTPS for all endpoints', async () => {
    const response = await request(app)
      .get('/api/profile')
      .set('X-Forwarded-Proto', 'http');

    // Should redirect to HTTPS
    expect(response.status).toBe(301);
    expect(response.headers.location).toMatch(/^https:/);
  });

  it('should set secure cookie flags', async () => {
    const response = await request(app)
      .post('/api/auth/login')
      .send({
        email: 'user@example.com',
        password: 'ValidPass123!'
      });

    const cookies = response.headers['set-cookie'];
    expect(cookies[0]).toContain('Secure');
    expect(cookies[0]).toContain('HttpOnly');
    expect(cookies[0]).toContain('SameSite=Strict');
  });
});
```

## Security Headers Testing

### HTTP Security Headers
```javascript
describe('Security Headers', () => {
  it('should set all required security headers', async () => {
    const response = await request(app).get('/');

    // Check security headers
    expect(response.headers['x-content-type-options']).toBe('nosniff');
    expect(response.headers['x-frame-options']).toBe('DENY');
    expect(response.headers['x-xss-protection']).toBe('1; mode=block');
    expect(response.headers['strict-transport-security']).toContain('max-age=');
    expect(response.headers['referrer-policy']).toBe('strict-origin-when-cross-origin');
    expect(response.headers['permissions-policy']).toBeDefined();
  });

  it('should not expose sensitive information in headers', async () => {
    const response = await request(app).get('/');

    // Should not expose server details
    expect(response.headers['server']).not.toContain('Express');
    expect(response.headers['x-powered-by']).toBeUndefined();
  });
});
```

## CSRF Protection Testing

### Cross-Site Request Forgery Prevention
```javascript
describe('CSRF Protection', () => {
  it('should require CSRF token for state-changing operations', async () => {
    const response = await request(app)
      .post('/api/transfer')
      .send({
        amount: 1000,
        recipient: 'attacker@example.com'
      });

    expect(response.status).toBe(403);
    expect(response.body.error).toContain('CSRF token');
  });

  it('should validate CSRF token correctly', async () => {
    // Get CSRF token
    const tokenResponse = await request(app)
      .get('/api/csrf-token');

    const csrfToken = tokenResponse.body.token;

    // Use token in request
    const response = await request(app)
      .post('/api/transfer')
      .set('X-CSRF-Token', csrfToken)
      .send({
        amount: 100,
        recipient: 'valid@example.com'
      });

    expect(response.status).toBe(200);
  });
});
```

## API Security Testing

### Rate Limiting
```javascript
describe('Rate Limiting', () => {
  it('should enforce rate limits', async () => {
    const requests = [];

    // Make 100 requests rapidly
    for (let i = 0; i < 100; i++) {
      requests.push(
        request(app).get('/api/products')
      );
    }

    const responses = await Promise.all(requests);

    // Some requests should be rate limited
    const rateLimited = responses.filter(r => r.status === 429);
    expect(rateLimited.length).toBeGreaterThan(0);

    // Check rate limit headers
    const limitedResponse = rateLimited[0];
    expect(limitedResponse.headers['x-ratelimit-limit']).toBeDefined();
    expect(limitedResponse.headers['x-ratelimit-remaining']).toBe('0');
    expect(limitedResponse.headers['x-ratelimit-reset']).toBeDefined();
  });
});
```

## Vulnerability Scanning

### Dependency Security
```javascript
describe('Dependency Security', () => {
  it('should not have known vulnerabilities', async () => {
    const { exec } = require('child_process');
    const { promisify } = require('util');
    const execAsync = promisify(exec);

    // Run security audit
    const { stdout } = await execAsync('npm audit --json');
    const audit = JSON.parse(stdout);

    expect(audit.metadata.vulnerabilities.high).toBe(0);
    expect(audit.metadata.vulnerabilities.critical).toBe(0);
  });
});
```

## Security Test Configuration

### Test Environment Setup
```javascript
// security.test.config.js
module.exports = {
  testEnvironment: {
    // Use test database
    DATABASE_URL: 'postgresql://test:test@localhost:5432/security_test',

    // Enable all security features
    ENFORCE_HTTPS: true,
    ENABLE_RATE_LIMIT: true,
    ENABLE_CSRF: true,
    SESSION_TIMEOUT: 1800000, // 30 minutes

    // Use test secrets
    JWT_SECRET: 'test-jwt-secret-key',
    ENCRYPTION_KEY: 'test-encryption-key'
  },

  setupFilesAfterEnv: ['./security-test-setup.js'],

  // Security test utilities
  globals: {
    securityUtils: {
      generateToken: (payload) => jwt.sign(payload, process.env.JWT_SECRET),
      hashPassword: (password) => bcrypt.hashSync(password, 10),
      encrypt: (data) => crypto.encrypt(data, process.env.ENCRYPTION_KEY)
    }
  }
};
```

## Quality Checklist

Before security tests are complete:
- [ ] Authentication mechanisms are tested
- [ ] Authorization boundaries are validated
- [ ] Input validation prevents injections
- [ ] XSS prevention is verified
- [ ] CSRF protection is tested
- [ ] File upload security is validated
- [ ] Encryption at rest/transit is verified
- [ ] Security headers are checked
- [ ] Rate limiting is tested
- [ ] Session management is secure
- [ ] No known vulnerabilities exist
- [ ] Error messages don't leak information

## Best Practices

### DO
- ✅ Test with realistic attack payloads
- ✅ Verify both positive and negative cases
- ✅ Test authorization at multiple levels
- ✅ Include security regression tests
- ✅ Automate security scanning
- ✅ Test in production-like environments
- ✅ Keep security tests updated
- ✅ Document security assumptions

### DON'T
- ❌ Test against production data
- ❌ Hardcode secrets in tests
- ❌ Skip edge cases
- ❌ Ignore failing security tests
- ❌ Test only happy paths
- ❌ Assume frameworks are secure
- ❌ Delay security testing
- ❌ Share test credentials

Remember: Security is not a feature, it's a requirement. Every endpoint, every input, every piece of data must be protected. Test security early and continuously.