from playwright.sync_api import sync_playwright, expect

def run(playwright):
    browser = playwright.chromium.launch(headless=True)
    page = browser.new_page()

    # Mock /api/me
    page.route("**/api/me", lambda route: route.fulfill(
        status=200,
        content_type="application/json",
        body='{"user": {"id": "123e4567-e89b-12d3-a456-426614174000", "email": "test@example.com", "name": "Test User", "avatar_url": "http://broken.com/image.png", "role": "admin"}, "organizations": []}'
    ))

    # Mock broken image
    page.route("http://broken.com/image.png", lambda route: route.fulfill(status=404))

    # Navigate to app
    page.goto("http://localhost:5173/")

    # Wait for the avatar in the header.
    user_menu_button = page.get_by_role("button", name="Open user menu")

    # Inside it, there should be an img.
    avatar_img = user_menu_button.locator("img")

    # Expect it to eventually point to ui-avatars with INITIALS "TU" (Test User)
    expect(avatar_img).to_have_attribute("src", "https://ui-avatars.com/api/?name=TU&background=random")

    # Take screenshot
    page.screenshot(path="/home/jules/verification/verification.png")

    browser.close()

with sync_playwright() as playwright:
    run(playwright)
