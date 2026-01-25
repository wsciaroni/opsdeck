from playwright.sync_api import sync_playwright

def verify_upload():
    with sync_playwright() as p:
        browser = p.chromium.launch(headless=True)
        page = browser.new_page()

        # Mock API responses
        page.route("**/api/me", lambda route: route.fulfill(
            status=200,
            content_type="application/json",
            body='{"user": {"id": "user1", "name": "Test User", "email": "test@example.com", "role": "admin"}, "organizations": [{"id": "org1", "name": "Test Org", "slug": "test-org", "role": "admin", "share_link_enabled": true}]}'
        ))

        page.route("**/api/tickets?*", lambda route: route.fulfill(
            status=200,
            content_type="application/json",
            body='[]'
        ))

        page.route("**/api/organizations/*/members", lambda route: route.fulfill(
            status=200,
            content_type="application/json",
            body='[{"user_id": "user1", "name": "Test User", "role": "admin"}]'
        ))

        # Navigate to Dashboard
        page.goto("http://localhost:5173/")

        # Wait for dashboard to load
        page.wait_for_selector("text=New Ticket", timeout=10000)

        # Click "New Ticket"
        page.get_by_role("button", name="New Ticket").first.click()

        # Wait for modal
        page.wait_for_selector("text=Create New Ticket")

        # Check for "Upload files" text or input
        page.wait_for_selector("text=Upload files")

        # Take screenshot
        page.screenshot(path="verification/upload_modal.png")

        browser.close()

if __name__ == "__main__":
    verify_upload()
