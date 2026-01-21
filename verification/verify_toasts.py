import os
import json
from playwright.sync_api import sync_playwright, expect

def run(playwright):
    browser = playwright.chromium.launch(headless=True)
    context = browser.new_context()
    page = context.new_page()

    # Log console messages
    page.on("console", lambda msg: print(f"Browser Console: {msg.text}"))

    # Mock API responses
    # Mock /api/me
    page.route("**/api/me", lambda route: route.fulfill(
        status=200,
        content_type="application/json",
        body=json.dumps({
            "user": {
                "id": "u1",
                "email": "test@example.com",
                "name": "Test User",
                "avatar_url": "https://ui-avatars.com/api/?name=Test+User"
            },
            "organizations": [
                {
                    "id": "o1",
                    "name": "Test Org",
                    "slug": "test-org",
                    "role": "owner"
                }
            ]
        })
    ))

    # Mock /api/tickets (GET) - Empty State
    # Note: verify query params if strictly needed, or just match the path base
    def handle_tickets(route):
        if route.request.method == "GET":
             route.fulfill(
                status=200,
                content_type="application/json",
                body=json.dumps([])
            )
        else:
            route.continue_()

    page.route("**/api/tickets?organization_id=*", handle_tickets)

    # Mock /api/tickets (POST) - Create Ticket
    page.route("**/api/tickets", lambda route: route.fulfill(
        status=201,
        content_type="application/json",
        body=json.dumps({
            "id": "t1",
            "title": "New Ticket",
            "description": "Description",
            "status_id": "new",
            "priority_id": "medium",
            "organization_id": "o1",
            "assignee_user_id": None,
            "created_at": "2023-01-01T00:00:00Z",
            "updated_at": "2023-01-01T00:00:00Z"
        })
    ))

    # 1. Navigate to Dashboard (expecting empty state)
    print("Navigating to Dashboard...")
    page.goto("http://localhost:5173")

    # Wait for empty state to appear
    print("Checking for Empty State...")
    try:
        expect(page.get_by_role("heading", name="No tickets found")).to_be_visible(timeout=5000)
    except Exception as e:
        print("Empty State not found. Taking screenshot.")
        page.screenshot(path="verification/failed_empty_state.png")
        raise e

    # Take screenshot of Empty State
    page.screenshot(path="verification/dashboard_empty_state.png")
    print("Screenshot saved: verification/dashboard_empty_state.png")

    # 2. Test Toast Notification (Create Ticket)
    print("Testing Ticket Creation Toast...")

    # Click New Ticket button in Empty State (it's the second one, or use a more specific selector)
    # The header button is still there, plus the empty state button.
    # Let's click the one inside the empty state.
    page.locator(".flex-col button").click()

    # Fill form
    page.fill("input[name='title']", "My First Ticket")
    page.fill("textarea[name='description']", "This is a test ticket.")

    # Update mock to return 1 ticket for subsequent calls (re-fetch after creation)
    page.unroute("**/api/tickets?organization_id=*")
    page.route("**/api/tickets?organization_id=*", lambda route: route.fulfill(
        status=200,
        content_type="application/json",
        body=json.dumps([
             {
                "id": "t1",
                "title": "My First Ticket",
                "description": "This is a test ticket.",
                "status_id": "new",
                "priority_id": "medium",
                "organization_id": "o1",
                "assignee_user_id": None,
                "created_at": "2023-01-01T00:00:00Z",
                "updated_at": "2023-01-01T00:00:00Z"
            }
        ])
    ))

    page.get_by_role("button", name="Create").click()

    # Check for Toast
    print("Waiting for toast...")
    try:
        expect(page.get_by_text("Ticket created!")).to_be_visible()
    except Exception as e:
        print("Toast not found. Taking screenshot.")
        page.screenshot(path="verification/failed_toast.png")
        raise e

    # Take screenshot of Toast
    page.screenshot(path="verification/dashboard_toast.png")
    print("Screenshot saved: verification/dashboard_toast.png")

    browser.close()

with sync_playwright() as playwright:
    run(playwright)
