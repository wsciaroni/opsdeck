from playwright.sync_api import Page, expect, sync_playwright
import time

def test_profile_image(page: Page):
    # Mock data
    user_with_avatar = {
        "user": {
            "id": "123",
            "email": "test@example.com",
            "name": "Test User",
            "role": "staff",
            "avatar_url": "https://avatars.githubusercontent.com/u/1?v=4",
            "created_at": "2023-01-01T00:00:00Z",
            "updated_at": "2023-01-01T00:00:00Z"
        },
        "organizations": []
    }

    user_without_avatar = {
        "user": {
            "id": "123",
            "email": "test@example.com",
            "name": "Test User No Avatar",
            "role": "staff",
            "avatar_url": "",
            "created_at": "2023-01-01T00:00:00Z",
            "updated_at": "2023-01-01T00:00:00Z"
        },
        "organizations": []
    }

    # Test Case 1: With Avatar
    print("Testing with avatar...")
    page.route("**/api/me", lambda route: route.fulfill(json=user_with_avatar))

    page.goto("http://localhost:5173/")

    profile_button = page.get_by_role("button", name="Open user menu")
    expect(profile_button).to_be_visible()
    time.sleep(1)

    images = profile_button.locator("img").count()
    if images > 0:
        print("PASS: Image found in profile button for user with avatar")
    else:
        print("FAIL: No image found in profile button for user with avatar")

    page.screenshot(path="/home/jules/verification/with_avatar.png")

    # Test Case 2: Without Avatar
    print("Testing without avatar...")
    # Clear cookies/storage to force re-fetch or just reload with new mock?
    # Playwright's route override should work if we navigate again?
    # But fetching usually happens on mount.

    # We unroute and route again, or just override.
    # To be safe, let's create a new context or page, or just override the route and reload.

    page.unroute("**/api/me")
    page.route("**/api/me", lambda route: route.fulfill(json=user_without_avatar))

    # Reload page to trigger fetchMe
    page.reload()

    profile_button = page.get_by_role("button", name="Open user menu")
    expect(profile_button).to_be_visible()
    time.sleep(1)

    images = profile_button.locator("img").count()
    if images == 0:
        # Check if fallback icon is present. We can check if it contains the svg.
        # But simpler is "no img tag".
        print("PASS: No image found in profile button for user without avatar (Fallback used)")
    else:
        print("FAIL: Image found in profile button for user without avatar")

    page.screenshot(path="/home/jules/verification/without_avatar.png")

if __name__ == "__main__":
    with sync_playwright() as p:
        browser = p.chromium.launch(headless=True)
        page = browser.new_page()
        try:
            test_profile_image(page)
        except Exception as e:
            print(f"Error: {e}")
        finally:
            browser.close()
