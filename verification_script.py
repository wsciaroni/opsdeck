from playwright.sync_api import Page, expect, sync_playwright

def test_ticket_board_links(page: Page):
    page.on("console", lambda msg: print(f"Console: {msg.text}"))
    page.on("pageerror", lambda exc: print(f"Page error: {exc}"))

    # Go to the test verification page
    try:
        page.goto("http://localhost:5173/test-verification", timeout=60000)
    except Exception as e:
        print(f"Navigation failed: {e}")
        return

    # Wait for something to ensure page loaded
    try:
        page.wait_for_selector("text=Ticket Board Verification", timeout=5000)
    except Exception as e:
        print("Header not found, taking screenshot")
        page.screenshot(path="debug_verification.png")
        raise e

    # Check if the ticket card is a link
    ticket1_link = page.get_by_role("link", name="Verification Ticket 1")
    expect(ticket1_link).to_be_visible()

    # Get the href attribute
    href = ticket1_link.get_attribute("href")
    print(f"Ticket 1 href: {href}")

    # Verify the href is correct (relative or absolute)
    assert "tickets/1" in href

    # Take a screenshot
    page.screenshot(path="verification.png")

if __name__ == "__main__":
    with sync_playwright() as p:
        browser = p.chromium.launch(headless=True)
        page = browser.new_page()
        try:
            test_ticket_board_links(page)
        except Exception as e:
            print(f"Test failed: {e}")
        finally:
            browser.close()
