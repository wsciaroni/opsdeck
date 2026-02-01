from playwright.sync_api import sync_playwright
import os
import re

def run():
    with sync_playwright() as p:
        browser = p.chromium.launch(headless=True)
        context = browser.new_context(viewport={"width": 1280, "height": 720})
        page = context.new_page()

        page.on("console", lambda msg: print(f"PAGE LOG: {msg.text}"))
        page.on("pageerror", lambda err: print(f"PAGE ERROR: {err}"))

        # Mock API
        # Use regex to be safer
        page.route(re.compile(r".*/api/me$"), lambda route: route.fulfill(
            status=200,
            content_type="application/json",
            body='{"user": {"id": "u1", "name": "Test User", "email": "test@example.com", "role": "admin"}, "organizations": [{"id": "o1", "name": "Test Org", "role": "admin"}]}'
        ))

        page.route(re.compile(r".*/api/tickets(\?.*)?$"), lambda route: route.fulfill(
            status=200,
            content_type="application/json",
            body='[]'
        ))

        print("Navigating...")
        try:
            page.goto("http://localhost:5173/")

            # Wait for "New Ticket" button
            print("Waiting for New Ticket button...")
            page.wait_for_selector("button:has-text('New Ticket')", timeout=10000)

            # Click "New Ticket"
            page.click("button:has-text('New Ticket')")

            # Wait for modal
            print("Waiting for modal...")
            page.wait_for_selector("h3:has-text('Create New Ticket')")

            # Create a dummy file
            with open("test_file.txt", "w") as f:
                f.write("test content")

            # Upload file
            print("Uploading file...")
            page.set_input_files("input[type='file']", "test_file.txt")

            # Check if file name appears
            print("Verifying file appears...")
            page.wait_for_selector("text=test_file.txt")

            # Check if remove button appears
            # Look for button with aria-label starting with "Remove"
            remove_btn = page.locator("button[aria-label^='Remove']")
            if remove_btn.count() > 0:
                print("Remove button found!")
            else:
                print("Remove button NOT found!")

            # Take screenshot BEFORE removal
            os.makedirs("verification", exist_ok=True)
            page.screenshot(path="verification/before_remove.png")

            # Click remove
            print("Clicking remove...")
            remove_btn.first.click()

            # Verify file is gone
            # wait for it to disappear
            try:
                page.wait_for_selector("text=test_file.txt", state="hidden", timeout=5000)
                print("File removed successfully!")
            except:
                print("File NOT removed (timeout)!")

            # Take screenshot AFTER removal
            page.screenshot(path="verification/after_remove.png")

        except Exception as e:
            print(f"Error: {e}")
            os.makedirs("verification", exist_ok=True)
            page.screenshot(path="verification/error.png")
        finally:
            browser.close()

if __name__ == "__main__":
    run()
