
import subprocess
import time
from playwright.sync_api import sync_playwright

def verify_ticket_detail():
    with sync_playwright() as p:
        browser = p.chromium.launch(headless=True)
        page = browser.new_page()

        # Mock API responses
        # 1. Mock /api/me (Authenticated User)
        page.route("**/api/me", lambda route: route.fulfill(
            status=200,
            content_type="application/json",
            body='{"user": {"id": "user1", "email": "test@opsdeck.dev", "name": "Test User"}, "organizations": [{"id": "org1", "name": "Test Org"}]}'
        ))

        # 2. Mock /api/tickets/ticket1 (Ticket Detail)
        page.route("**/api/tickets/ticket1", lambda route: route.fulfill(
            status=200,
            content_type="application/json",
            body='{"id": "ticket1", "title": "Fix Login Bug", "description": "Users cannot login on Safari.", "status_id": "new", "priority_id": "high", "reporter_name": "Alice Reporter", "created_at": "2023-10-27T10:00:00Z"}'
        ))

        # 3. Mock /api/tickets?organization_id=org1 (Dashboard List)
        page.route("**/api/tickets?organization_id=org1", lambda route: route.fulfill(
            status=200,
            content_type="application/json",
            body='[{"id": "ticket1", "title": "Fix Login Bug", "status_id": "new", "priority_id": "high", "created_at": "2023-10-27T10:00:00Z"}]'
        ))

        try:
            # Go to the Ticket Detail page directly (client-side routing should handle this)
            # Note: We need the frontend server running.
            # Assuming it's running on localhost:5173 (Vite default)
            page.goto("http://localhost:5173/tickets/ticket1")

            # Wait for content to load
            page.wait_for_selector("text=Fix Login Bug")

            # Verify details
            assert page.is_visible("text=Alice Reporter")
            assert page.is_visible("text=Users cannot login on Safari.")

            # Take screenshot
            page.screenshot(path="verification/ticket_detail.png")
            print("Screenshot saved to verification/ticket_detail.png")

        except Exception as e:
            print(f"Error: {e}")
        finally:
            browser.close()

if __name__ == "__main__":
    verify_ticket_detail()
