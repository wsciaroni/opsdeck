from playwright.sync_api import sync_playwright, expect

def run(playwright):
    browser = playwright.chromium.launch(headless=True)
    context = browser.new_context()
    page = context.new_page()

    page.on("console", lambda msg: print(f"Console: {msg.text}"))
    page.on("pageerror", lambda err: print(f"Page Error: {err}"))

    def handle_route(route):
        url = route.request.url

        # Match /api/me (Auth)
        if "/api/me" in url and "/src/" not in url:
            print(f"  -> Mocking Auth (Me): {url}")
            route.fulfill(
                status=200,
                content_type="application/json",
                body='{"user": {"id": "user1", "email": "test@example.com", "name": "Test User", "role": "admin", "avatar_url": "", "created_at": "2023-01-01T00:00:00Z", "updated_at": "2023-01-01T00:00:00Z"}, "organizations": [{"id": "org1", "name": "Test Org", "slug": "test-org", "role": "admin", "created_at": "2023-01-01T00:00:00Z", "updated_at": "2023-01-01T00:00:00Z"}]}'
            )
            return

        # /api/tickets (list)
        if "/api/tickets" in url and "/src/" not in url:
            print(f"  -> Mocking Tickets: {url}")
            route.fulfill(
                status=200,
                content_type="application/json",
                body='[]'
            )
            return

        route.continue_()

    page.route("**/*", handle_route)

    # Navigate to dashboard
    page.goto("http://localhost:5173/")

    print(f"URL after load: {page.url}")

    try:
        # Wait for dashboard to load (look for "Tickets" header)
        expect(page.get_by_role("heading", name="Tickets")).to_be_visible(timeout=5000)

        # Click the "Filter" button to open popover
        page.get_by_role("button", name="Filter").click()

        # Wait for popover content to appear
        expect(page.get_by_text("Status", exact=True)).to_be_visible()
        expect(page.get_by_text("Priority", exact=True)).to_be_visible()

        # Take screenshot
        page.screenshot(path="verification/filter_popover.png")
    except Exception as e:
        print(f"Error: {e}")
        page.screenshot(path="verification/error.png")
        print(f"Error URL: {page.url}")

    browser.close()

with sync_playwright() as playwright:
    run(playwright)
