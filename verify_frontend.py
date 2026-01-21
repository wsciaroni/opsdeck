from playwright.sync_api import sync_playwright

def verify_team_management():
    with sync_playwright() as p:
        browser = p.chromium.launch(headless=True)
        # Set screen size to ensure layout is visible
        context = browser.new_context(viewport={"width": 1280, "height": 720})
        page = context.new_page()

        # Define API responses for Auth and Org
        user = {
            "id": "u1",
            "email": "owner@example.com",
            "name": "Owner User",
            "avatar_url": "",
            "role": "owner",
            "created_at": "2023-01-01T00:00:00Z",
            "updated_at": "2023-01-01T00:00:00Z"
        }

        org = {
            "id": "o1",
            "name": "Test Org",
            "slug": "test-org",
            "role": "owner",
            "created_at": "2023-01-01T00:00:00Z",
            "updated_at": "2023-01-01T00:00:00Z"
        }

        # Mock API routes
        # 1. /api/me
        page.route("**/api/me", lambda route: route.fulfill(
            status=200,
            content_type="application/json",
            body='{"user": ' + str(user).replace("'", '"') + ', "organizations": [' + str(org).replace("'", '"') + ']}'
        ))

        # 2. /api/organizations/o1/members
        members = [
            {"id": "u1", "email": "owner@example.com", "name": "Owner User", "role": "owner", "avatar_url": ""},
            {"id": "u2", "email": "member@example.com", "name": "Team Member", "role": "member", "avatar_url": ""}
        ]
        page.route("**/api/organizations/o1/members", lambda route: route.fulfill(
            status=200,
            content_type="application/json",
            body=str(members).replace("'", '"')
        ))

        # 3. Add Member Mock
        page.route("**/api/organizations/o1/members", lambda route: route.fulfill(
            status=201,
            body=""
        ) if route.request.method == "POST" else route.continue_())

        try:
            # Use 5175
            print("Navigating to page...")
            page.goto("http://localhost:5175/organizations/o1/settings/team")

            # Wait for content
            print("Waiting for selector...")
            page.wait_for_selector("h1:has-text('Team Members')", timeout=15000)

            # Verify list shows members
            page.wait_for_selector("text=owner@example.com")
            page.wait_for_selector("text=member@example.com")

            # Take screenshot
            page.screenshot(path="verification_team_settings.png")
            print("Screenshot saved to verification_team_settings.png")

        except Exception as e:
            print(f"Verification failed: {e}")
            page.screenshot(path="verification_failure.png")
        finally:
            browser.close()

if __name__ == "__main__":
    verify_team_management()
