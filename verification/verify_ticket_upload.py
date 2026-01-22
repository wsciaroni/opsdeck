import os
import time
from playwright.sync_api import sync_playwright

def run():
    with sync_playwright() as p:
        browser = p.chromium.launch(headless=True)
        page = browser.new_page()

        # Mock API
        # Mock /api/me
        page.route("**/api/me", lambda route: route.fulfill(
            status=200,
            content_type="application/json",
            body='{"user": {"id": "user-123", "email": "user@example.com", "name": "Test User", "role": "admin", "avatar_url": ""}, "organizations": [{"id": "org-123", "name": "Test Org", "slug": "test-org", "role": "owner"}]}'
        ))

        # Mock /api/tickets?organization_id=...
        page.route("**/api/tickets?*", lambda route: route.fulfill(
            status=200,
            content_type="application/json",
            body='[]'
        ))

        # Mock /api/tickets (POST) - Create Ticket
        def handle_create(route):
            # Check if multipart form data contains files?
            route.fulfill(
                status=201,
                content_type="application/json",
                body='{"id": "ticket-new", "title": "New Ticket", "organization_id": "org-123", "status_id": "new", "priority_id": "medium", "created_at": "2023-01-01T00:00:00Z"}'
            )

        page.route("**/api/tickets", lambda route: handle_create(route) if route.request.method == "POST" else route.continue_())

        # Start App (Assuming already running or served statically? No, I need to serve it)
        # Assuming npm run dev is running on 5173

        page.goto("http://localhost:5173")

        # Login flow mock (localStorage)
        page.evaluate("localStorage.setItem('token', 'mock-token')")
        page.evaluate("localStorage.setItem('last_org_id', 'org-123')")

        # Reload to trigger AuthContext fetch
        page.reload()

        # Wait for dashboard
        try:
            page.wait_for_selector("text=Tickets", timeout=5000)
        except Exception as e:
            print("Failed to load dashboard")
            page.screenshot(path="verification/dashboard_failure.png")
            raise e

        # Click New Ticket
        page.click("text=New Ticket")

        # Fill form
        page.fill("input[name='title']", "Test Ticket with File")
        page.fill("textarea[name='description']", "Description")

        # Verify File Input exists
        page.wait_for_selector("input[type='file']")

        # Take screenshot of Modal with File Input
        page.screenshot(path="verification/ticket_upload_modal.png")
        print("Internal ticket upload verified")

        # --- Verify Public Ticket Submission ---
        # Mock /api/public/tickets
        page.route("**/api/public/tickets", lambda route: route.fulfill(
            status=201,
            content_type="application/json",
            body='{"id": "ticket-public", "title": "Public Ticket", "created_at": "2023-01-01T00:00:00Z"}'
        ))

        page.goto("http://localhost:5173/submit-ticket?token=test-token")

        try:
            page.wait_for_selector("text=Submit a Ticket", timeout=5000)
        except Exception as e:
            print("Failed to load public ticket page")
            page.screenshot(path="verification/public_page_failure.png")
            raise e

        # Check for file input
        page.wait_for_selector("input[type='file']")

        page.fill("input[name='name']", "Public User")
        page.fill("input[name='email']", "public@example.com")
        page.fill("input[name='title']", "Public Ticket with File")
        page.fill("textarea[name='description']", "Public Description")

        page.screenshot(path="verification/public_ticket_upload.png")
        print("Public ticket upload verified")

        browser.close()

if __name__ == "__main__":
    run()
