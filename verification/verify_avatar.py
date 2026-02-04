from playwright.sync_api import sync_playwright
import json

def verify_avatar():
    with sync_playwright() as p:
        browser = p.chromium.launch(headless=True)
        context = browser.new_context(viewport={"width": 1280, "height": 720})
        page = context.new_page()

        # Debug console
        page.on("console", lambda msg: print(f"CONSOLE: {msg.text}"))
        page.on("pageerror", lambda err: print(f"PAGE ERROR: {err}"))

        # Mock /api/me to avoid login redirect
        user = {
            "id": "u1",
            "email": "user@example.com",
            "name": "Test User",
            "role": "admin",
            "created_at": "2023-01-01T00:00:00Z"
        }
        page.route("**/api/me", lambda route: route.fulfill(
            status=200,
            content_type="application/json",
            body=json.dumps({"user": user})
        ))

        # Mock Ticket details
        ticket = {
            "id": "t1",
            "organization_id": "o1",
            "title": "Test Ticket",
            "description": "Description",
            "status_id": "open",
            "priority_id": "low",
            "reporter_id": "u1",
            "created_at": "2023-01-01T00:00:00Z",
            "updated_at": "2023-01-01T00:00:00Z",
            "reporter_name": "Test User",
            "files": []
        }
        page.route("**/api/tickets/t1", lambda route: route.fulfill(
            status=200,
            content_type="application/json",
            body=json.dumps(ticket)
        ))

        # Mock Members (empty list)
        page.route("**/api/organizations/o1/members", lambda route: route.fulfill(
            status=200,
            content_type="application/json",
            body=json.dumps([])
        ))

        # Mock Comments
        comments = [
            {
                "id": "c1",
                "body": "Comment with no avatar",
                "sensitive": False,
                "created_at": "2023-01-01T00:00:00Z",
                "user": {
                    "id": "u1",
                    "name": "John Doe",
                    "avatar_url": ""
                }
            },
            {
                "id": "c2",
                "body": "Comment with avatar",
                "sensitive": False,
                "created_at": "2023-01-01T00:00:00Z",
                "user": {
                    "id": "u2",
                    "name": "Jane Smith",
                    "avatar_url": "https://example.com/avatar.jpg"
                }
            }
        ]
        page.route("**/api/tickets/t1/comments", lambda route: route.fulfill(
            status=200,
            content_type="application/json",
            body=json.dumps(comments)
        ))

        try:
            print("Navigating to ticket page...")
            page.goto("http://localhost:5173/tickets/t1")

            # Wait for any text to ensure page loaded
            page.wait_for_load_state("networkidle")

            print("Waiting for comments...")
            page.wait_for_selector("text=Comment with no avatar", timeout=10000)

            # Check for Initials "JD" for John Doe
            page.wait_for_selector("text=JD")
            print("Found initials 'JD'")

            # Check that there is NO img with ui-avatars.com
            images = page.locator("img").all()
            for img in images:
                src = img.get_attribute("src")
                if src and "ui-avatars.com" in src:
                    raise Exception(f"Found ui-avatars.com image: {src}")
            print("No ui-avatars.com images found.")

            page.screenshot(path="verification_avatar.png")
            print("Screenshot saved to verification_avatar.png")

        except Exception as e:
            print(f"Verification failed: {e}")
            page.screenshot(path="verification_failure.png")
            raise e
        finally:
            browser.close()

if __name__ == "__main__":
    verify_avatar()
