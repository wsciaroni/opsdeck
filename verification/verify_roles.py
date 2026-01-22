
import asyncio
from playwright.async_api import async_playwright, expect

async def run():
    async with async_playwright() as p:
        browser = await p.chromium.launch(headless=True)
        context = await browser.new_context()
        page = await context.new_page()

        # Mocking API responses
        await page.route("**/api/me", lambda route: route.fulfill(
            status=200,
            content_type="application/json",
            body='{"user": {"id": "user1", "email": "admin@example.com", "name": "Admin User", "avatar_url": ""}, "organizations": [{"id": "org1", "name": "Org 1", "slug": "org1", "created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z"}]}'
        ))

        await page.route("**/api/organizations/*/members", lambda route: route.fulfill(
            status=200,
            content_type="application/json",
            body='[{"id": "user1", "email": "admin@example.com", "name": "Admin User", "avatar_url": "", "role": "owner"}, {"id": "user2", "email": "member@example.com", "name": "Member User", "avatar_url": "", "role": "member"}]'
        ))

        await page.route("**/api/organizations/*/share", lambda route: route.fulfill(
             status=200,
             content_type="application/json",
             body='{"share_link_enabled": true, "share_link_token": "token123"}'
        ))

        try:
            print("Navigating to page...")
            await page.goto("http://localhost:5173/organizations/org1/settings/team")

            print("Waiting for page content...")
            await page.screenshot(path="verification/debug_load.png")

            # Wait for members to load
            print("Waiting for Admin User text...")
            await expect(page.get_by_text("Admin User")).to_be_visible()

            print("Waiting for Member User text...")
            await expect(page.get_by_text("Member User")).to_be_visible()

            # Check if the role menu is present for the member
            row = page.get_by_role("listitem").filter(has_text="Member User")

            # Find the role button in that row. It should say "member"
            # Being specific about text match
            role_button = row.get_by_role("button", name="member", exact=True)
            await expect(role_button).to_be_visible()

            # Click the role button to open the menu
            await role_button.click()

            # Check if options are visible
            await expect(page.get_by_role("menuitem", name="Admin")).to_be_visible()
            await expect(page.get_by_role("menuitem", name="Owner")).to_be_visible()

            # Screenshot the open menu
            await page.screenshot(path="verification/role_menu.png")
            print("Screenshot saved to verification/role_menu.png")

        except Exception as e:
            print(f"Error: {e}")
            await page.screenshot(path="verification/error.png")

        await browser.close()

if __name__ == "__main__":
    asyncio.run(run())
