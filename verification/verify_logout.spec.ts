import { test, expect } from '@playwright/test';

test('Verify logout button in layout', async ({ page }) => {
  // Mock the /api/me endpoint to simulate a logged-in user
  await page.route('/api/me', async route => {
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({
        user: {
          id: '123',
          email: 'test@example.com',
          name: 'Test User',
          avatar_url: ''
        },
        organizations: [
          {
            id: 'org1',
            name: 'Test Org',
            slug: 'test-org',
            role: 'owner'
          }
        ]
      })
    });
  });

  // Mock tickets to avoid errors on dashboard
  await page.route('/api/tickets?organization_id=org1', async route => {
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify([])
    });
  });

  // Mock logout endpoint
  await page.route('/auth/logout', async route => {
    await route.fulfill({
      status: 200
    });
  });

  // Navigate to the dashboard
  await page.goto('http://localhost:5173/');

  // Wait for the layout to load and the user email to be visible
  await expect(page.getByText('test@example.com')).toBeVisible();

  // Check for the "OpsDeck" logo/text
  await expect(page.getByText('OpsDeck')).toBeVisible();

  // Check for the "Logout" button
  const logoutButton = page.getByRole('button', { name: 'Logout' });
  await expect(logoutButton).toBeVisible();

  // Take a screenshot of the dashboard with the header
  await page.screenshot({ path: 'verification/dashboard_header.png' });

  // Click logout
  await logoutButton.click();

  // After logout, the user state should be cleared and we should see the login screen
  // The login screen has "Login with Google" button
  await expect(page.getByRole('button', { name: 'Login with Google' })).toBeVisible();

  // Take a screenshot of the login screen
  await page.screenshot({ path: 'verification/login_screen_after_logout.png' });
});
