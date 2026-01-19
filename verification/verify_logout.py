from playwright.sync_api import sync_playwright
import json

def verify_logout():
    with sync_playwright() as p:
        browser = p.chromium.launch(headless=True)
        page = browser.new_page()

        # Mock /api/me to simulate logged-in user
        def handle_me(route):
            route.fulfill(
                status=200,
                content_type="application/json",
                body=json.dumps({
                    "user": {
                        "id": "123",
                        "email": "test@example.com",
                        "name": "Test User",
                        "avatar_url": ""
                    },
                    "organizations": [
                        {
                            "id": "org1",
                            "name": "Test Org",
                            "slug": "test-org",
                            "role": "owner"
                        }
                    ]
                })
            )

        page.route("**/api/me", handle_me)

        # Mock /api/tickets to avoid errors
        def handle_tickets(route):
            route.fulfill(
                status=200,
                content_type="application/json",
                body=json.dumps([])
            )

        page.route("**/api/tickets?organization_id=org1", handle_tickets)

        # Mock /auth/logout
        def handle_logout(route):
            route.fulfill(status=200)

        page.route("**/auth/logout", handle_logout)

        # Go to app
        page.goto("http://localhost:5173/")

        # Verify header elements
        print("Checking for user email...")
        page.get_by_text("test@example.com").wait_for()

        print("Checking for OpsDeck logo...")
        page.get_by_text("OpsDeck").first.wait_for()

        print("Checking for Logout button...")
        logout_btn = page.get_by_role("button", name="Logout")
        logout_btn.wait_for()

        # Screenshot dashboard
        page.screenshot(path="verification/dashboard_header.png")
        print("Dashboard screenshot taken.")

        # Click logout
        logout_btn.click()

        # Verify we are back at login screen (mocked /me will still return user, but app state should clear momentarily?
        # WAIT: In the actual app, logout clears the state.
        # But if the page reloads, it fetches /me again.
        # The Logout function in AuthContext clears state but does NOT force reload page, it just sets user to null.
        # So we should see "Login with Google".

        print("Checking for Login button after logout...")
        login_btn = page.get_by_role("button", name="Login with Google")
        login_btn.wait_for()

        page.screenshot(path="verification/login_screen_after_logout.png")
        print("Logout verification complete.")

        browser.close()

if __name__ == "__main__":
    verify_logout()
