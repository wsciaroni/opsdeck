from playwright.sync_api import sync_playwright, Page, expect
import time

def run(playwright):
    browser = playwright.chromium.launch(headless=True)
    context = browser.new_context(viewport={'width': 1280, 'height': 720})
    page = context.new_page()

    def handle_api_me(route):
        print(f"Intercepted ME: {route.request.url}")
        if "/src/" in route.request.url:
            route.continue_()
            return
        route.fulfill(
            status=200,
            content_type="application/json",
            body='''{
                "user": {
                    "id": "u1",
                    "email": "test@example.com",
                    "name": "Test User",
                    "role": "admin",
                    "avatar_url": "",
                    "created_at": "2023-01-01T00:00:00Z",
                    "updated_at": "2023-01-01T00:00:00Z"
                },
                "organizations": [
                    {
                        "id": "org1",
                        "name": "Test Org",
                        "slug": "test-org",
                        "role": "admin",
                        "created_at": "2023-01-01T00:00:00Z",
                        "updated_at": "2023-01-01T00:00:00Z"
                    }
                ]
            }'''
        )

    def handle_api_tickets(route):
        print(f"Intercepted TICKETS: {route.request.url}")
        if "/src/" in route.request.url:
            route.continue_()
            return
        route.fulfill(
            status=200,
            content_type="application/json",
            body='''[
                {
                    "id": "t1",
                    "organization_id": "org1",
                    "title": "Task 1",
                    "description": "Desc 1",
                    "status_id": "new",
                    "priority_id": "high",
                    "reporter_id": "u1",
                    "assignee_user_id": null,
                    "sensitive": false,
                    "created_at": "2023-01-01T00:00:00Z",
                    "updated_at": "2023-01-01T00:00:00Z",
                    "completed_at": null
                },
                {
                    "id": "t2",
                    "organization_id": "org1",
                    "title": "Task 2",
                    "description": "Desc 2",
                    "status_id": "done",
                    "priority_id": "medium",
                    "reporter_id": "u1",
                    "assignee_user_id": null,
                    "sensitive": false,
                    "created_at": "2023-01-01T00:00:00Z",
                    "updated_at": "2023-01-01T00:00:00Z",
                    "completed_at": null
                }
            ]'''
        )

    def handle_api_members(route):
        print(f"Intercepted MEMBERS: {route.request.url}")
        if "/src/" in route.request.url:
            route.continue_()
            return
        route.fulfill(
            status=200,
            content_type="application/json",
            body="[]"
        )

    # Mock /api/me
    page.route("**/*api/me", handle_api_me)

    # Mock /api/tickets
    page.route("**/*api/tickets*", handle_api_tickets)

    # Mock /api/organizations/*/members
    page.route("**/*api/organizations/*/members*", handle_api_members)

    page.add_init_script("localStorage.setItem('dashboard_view_mode', 'board')")

    page.on("console", lambda msg: print(f"CONSOLE: {msg.text}"))
    page.on("pageerror", lambda err: print(f"PAGE ERROR: {err}"))
    page.on("requestfailed", lambda request: print(f"FAILED: {request.url} {request.failure}"))

    page.goto("http://localhost:5173/")

    try:
        page.wait_for_selector("text=New", timeout=5000)
    except:
        print(f"Timeout. Current URL: {page.url}")
        print("Page Content Snippet:")
        print(page.content()[:500])

        if page.locator("text=Sign in").is_visible():
            print("On Login page.")

        page.screenshot(path="/home/jules/verification/error.png")
        return

    # 1. Verify default view
    print("Verifying default view...")
    if page.is_visible("text=Done"):
        print("FAIL: Done column is visible in default view")
    else:
        print("PASS: Done column is hidden in default view")

    page.screenshot(path="/home/jules/verification/board_default.png")

    # 2. Filter to show 'Done'
    print("Filtering to show Done...")
    page.click("button:has-text('Filter')")
    page.wait_for_selector("text=Status")
    page.click("text=Done")
    page.mouse.click(0, 0)
    time.sleep(1)

    if page.is_visible("text=Done"):
        print("PASS: Done column is visible after filtering")
    else:
        print("FAIL: Done column is hidden after filtering")

    page.screenshot(path="/home/jules/verification/board_filtered.png")

    new_col = page.locator("text=New").first
    column = new_col.locator("xpath=../..")
    width = column.evaluate("el => el.getBoundingClientRect().width")
    print(f"Column width: {width}px")

    browser.close()

with sync_playwright() as playwright:
    run(playwright)
