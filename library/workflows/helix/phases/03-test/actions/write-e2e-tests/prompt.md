# Write End-to-End Tests Prompt

Create end-to-end tests that validate complete user journeys and critical business flows through the entire application stack. These tests verify that the system works correctly from the user's perspective.

## Test Output Location

Generate end-to-end tests in: `tests/e2e/`

Organize tests by journey type:
- `tests/e2e/critical/` - Critical path tests
- `tests/e2e/features/` - Feature-specific flows
- `tests/e2e/regression/` - Regression test suites
- `tests/e2e/performance/` - Performance validation

## Purpose

End-to-end tests ensure:
- Critical user paths work completely
- UI, backend, and database integrate correctly
- Business workflows execute properly
- User experience meets requirements
- System behaves correctly under real-world usage

## Test Scope

### Critical User Journeys
Focus on paths that:
- Generate revenue (checkout, subscription)
- Acquire users (registration, onboarding)
- Retain users (core features)
- Handle sensitive data (authentication, payments)
- Represent common use cases

### Cross-Browser Testing
Verify functionality across:
- Chrome (latest)
- Firefox (latest)
- Safari (latest)
- Edge (latest)
- Mobile browsers (iOS Safari, Chrome Android)

## Test Implementation

### Browser Automation Setup
```javascript
// Using Playwright for cross-browser testing
import { test, expect } from '@playwright/test';

test.describe('User Registration Flow', () => {
  test('should complete registration and access dashboard', async ({ page }) => {
    // Navigate to registration
    await page.goto('/register');

    // Fill registration form
    await page.fill('[name="email"]', 'newuser@example.com');
    await page.fill('[name="password"]', 'SecurePass123!');
    await page.fill('[name="confirmPassword"]', 'SecurePass123!');

    // Submit and wait for navigation
    await Promise.all([
      page.waitForNavigation(),
      page.click('[type="submit"]')
    ]);

    // Verify redirect to dashboard
    expect(page.url()).toContain('/dashboard');

    // Verify welcome message
    await expect(page.locator('.welcome-message')).toContainText('Welcome');

    // Verify user can access features
    await page.click('[data-test="create-project"]');
    await expect(page.locator('.project-form')).toBeVisible();
  });
});
```

### Mobile Testing
```javascript
test('should work on mobile devices', async ({ browser }) => {
  // Create mobile context
  const context = await browser.newContext({
    viewport: { width: 375, height: 667 },
    userAgent: 'Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X)',
    hasTouch: true
  });

  const page = await context.newPage();
  await page.goto('/');

  // Test mobile-specific interactions
  await page.tap('[data-test="mobile-menu"]');
  await expect(page.locator('.mobile-nav')).toBeVisible();

  // Test touch gestures
  await page.locator('.carousel').swipe({ direction: 'left' });
  await expect(page.locator('.slide-2')).toBeVisible();
});
```

## Critical Path Testing

### E-Commerce Checkout Flow
```javascript
test('complete purchase flow', async ({ page }) => {
  // Login
  await page.goto('/login');
  await page.fill('[name="email"]', 'testuser@example.com');
  await page.fill('[name="password"]', 'password123');
  await page.click('[type="submit"]');

  // Browse and add to cart
  await page.goto('/products');
  await page.click('[data-product-id="123"] .add-to-cart');
  await expect(page.locator('.cart-count')).toHaveText('1');

  // Checkout
  await page.goto('/checkout');

  // Fill shipping info
  await page.fill('[name="shipping.address"]', '123 Test St');
  await page.fill('[name="shipping.city"]', 'Test City');
  await page.selectOption('[name="shipping.state"]', 'CA');
  await page.fill('[name="shipping.zip"]', '90210');

  // Fill payment info (using test card)
  await page.frameLocator('#card-frame').locator('[name="cardnumber"]').fill('4242424242424242');
  await page.frameLocator('#card-frame').locator('[name="exp-date"]').fill('12/25');
  await page.frameLocator('#card-frame').locator('[name="cvc"]').fill('123');

  // Complete order
  await page.click('[data-test="place-order"]');

  // Verify confirmation
  await page.waitForSelector('.order-confirmation');
  const orderNumber = await page.locator('.order-number').textContent();
  expect(orderNumber).toMatch(/^ORD-\d{6}$/);

  // Verify email was sent (check test email service)
  const emails = await getTestEmails('testuser@example.com');
  expect(emails[0].subject).toContain('Order Confirmation');
});
```

### Authentication Flow
```javascript
test('authentication and authorization', async ({ page }) => {
  // Test login
  await page.goto('/login');
  await page.fill('[name="email"]', 'admin@example.com');
  await page.fill('[name="password"]', 'AdminPass123!');
  await page.click('[type="submit"]');

  // Verify admin access
  await page.goto('/admin');
  await expect(page.locator('.admin-dashboard')).toBeVisible();

  // Test logout
  await page.click('[data-test="logout"]');
  await expect(page).toHaveURL('/');

  // Verify cannot access admin after logout
  await page.goto('/admin');
  await expect(page).toHaveURL('/login?redirect=/admin');
});
```

## Data Management

### Test Data Setup
```javascript
// Reset and seed data before each test
test.beforeEach(async () => {
  await resetDatabase();
  await seedTestData({
    users: [
      { email: 'testuser@example.com', role: 'user' },
      { email: 'admin@example.com', role: 'admin' }
    ],
    products: [
      { id: '123', name: 'Test Product', price: 99.99 }
    ]
  });
});

// Clean up after tests
test.afterEach(async () => {
  await cleanupTestData();
});
```

## Performance Monitoring

### Page Load Performance
```javascript
test('should load quickly', async ({ page }) => {
  const metrics = [];

  // Collect performance metrics
  page.on('load', async () => {
    const perf = await page.evaluate(() => ({
      domContentLoaded: performance.timing.domContentLoadedEventEnd - performance.timing.navigationStart,
      loadComplete: performance.timing.loadEventEnd - performance.timing.navigationStart,
      firstPaint: performance.getEntriesByName('first-paint')[0]?.startTime,
      firstContentfulPaint: performance.getEntriesByName('first-contentful-paint')[0]?.startTime
    }));
    metrics.push(perf);
  });

  await page.goto('/');

  // Assert performance thresholds
  expect(metrics[0].firstContentfulPaint).toBeLessThan(1500);
  expect(metrics[0].loadComplete).toBeLessThan(3000);
});
```

## Visual Regression Testing

### Screenshot Comparisons
```javascript
test('visual consistency', async ({ page }) => {
  await page.goto('/');

  // Take screenshot
  await expect(page).toHaveScreenshot('homepage.png', {
    fullPage: true,
    animations: 'disabled'
  });

  // Test responsive design
  await page.setViewportSize({ width: 768, height: 1024 });
  await expect(page).toHaveScreenshot('homepage-tablet.png');

  await page.setViewportSize({ width: 375, height: 667 });
  await expect(page).toHaveScreenshot('homepage-mobile.png');
});
```

## Accessibility Testing

### WCAG Compliance
```javascript
test('should be accessible', async ({ page }) => {
  await page.goto('/');

  // Run accessibility checks
  const violations = await page.evaluate(() => {
    return window.axe.run();
  });

  expect(violations.violations).toHaveLength(0);

  // Test keyboard navigation
  await page.keyboard.press('Tab');
  const focusedElement = await page.evaluate(() => document.activeElement.tagName);
  expect(focusedElement).not.toBe('BODY');

  // Test screen reader compatibility
  const ariaLabels = await page.$$eval('[aria-label]', elements =>
    elements.map(el => el.getAttribute('aria-label'))
  );
  expect(ariaLabels.every(label => label && label.length > 0)).toBe(true);
});
```

## Test Organization

### Page Object Model
```javascript
// pages/LoginPage.js
class LoginPage {
  constructor(page) {
    this.page = page;
    this.emailInput = page.locator('[name="email"]');
    this.passwordInput = page.locator('[name="password"]');
    this.submitButton = page.locator('[type="submit"]');
    this.errorMessage = page.locator('.error-message');
  }

  async navigate() {
    await this.page.goto('/login');
  }

  async login(email, password) {
    await this.emailInput.fill(email);
    await this.passwordInput.fill(password);
    await this.submitButton.click();
  }

  async getErrorMessage() {
    return await this.errorMessage.textContent();
  }
}

// Usage in tests
test('should show error for invalid credentials', async ({ page }) => {
  const loginPage = new LoginPage(page);
  await loginPage.navigate();
  await loginPage.login('invalid@example.com', 'wrongpassword');

  const error = await loginPage.getErrorMessage();
  expect(error).toBe('Invalid email or password');
});
```

## Configuration

### Test Environment Setup
```javascript
// playwright.config.js
module.exports = {
  testDir: './tests/e2e',
  timeout: 30000,
  retries: 2,
  workers: 4,

  use: {
    baseURL: process.env.E2E_BASE_URL || 'http://localhost:3000',
    screenshot: 'only-on-failure',
    video: 'retain-on-failure',
    trace: 'on-first-retry'
  },

  projects: [
    {
      name: 'Chrome',
      use: { ...devices['Desktop Chrome'] }
    },
    {
      name: 'Firefox',
      use: { ...devices['Desktop Firefox'] }
    },
    {
      name: 'Safari',
      use: { ...devices['Desktop Safari'] }
    },
    {
      name: 'Mobile Chrome',
      use: { ...devices['Pixel 5'] }
    },
    {
      name: 'Mobile Safari',
      use: { ...devices['iPhone 12'] }
    }
  ]
};
```

## Quality Checklist

Before E2E tests are complete:
- [ ] All critical user paths are tested
- [ ] Tests work across all target browsers
- [ ] Mobile experience is validated
- [ ] Performance thresholds are verified
- [ ] Accessibility standards are met
- [ ] Visual regression tests are in place
- [ ] Tests are stable (no flakiness)
- [ ] Test data is properly managed
- [ ] Error scenarios are covered
- [ ] Tests can run in parallel

## Best Practices

### DO
- ✅ Test from the user's perspective
- ✅ Focus on critical business flows
- ✅ Use Page Object Model for maintainability
- ✅ Wait for elements explicitly
- ✅ Test across browsers and devices
- ✅ Include accessibility checks
- ✅ Monitor performance metrics
- ✅ Use test IDs for reliable selection

### DON'T
- ❌ Test implementation details
- ❌ Use brittle selectors (classes that change)
- ❌ Hardcode wait times
- ❌ Share state between tests
- ❌ Test third-party services
- ❌ Ignore flaky tests
- ❌ Run all E2E tests on every commit
- ❌ Test what unit/integration tests cover

Remember: E2E tests are expensive but valuable. Focus on critical paths that would cause significant business impact if broken.