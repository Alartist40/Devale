from playwright.sync_api import sync_playwright

def verify_new_features():
    with sync_playwright() as p:
        browser = p.chromium.launch(headless=True)
        page = browser.new_page()
        page.set_viewport_size({"width": 1280, "height": 800})

        page.goto("http://localhost:5178")
        page.wait_for_selector(".sidebar")

        # 1. Check App Store browsers
        page.click("text=APP STORE")
        page.wait_for_timeout(500)
        page.screenshot(path="verify_appstore_browsers.png")

        # 2. Check Tools (MRT, TaskMgr)
        page.click("text=TOOLS")
        page.wait_for_timeout(500)
        page.screenshot(path="verify_tools_mrt.png")

        # 3. Check Stop Button (Simulation - not easy to trigger without actual repair)
        # But we can check if layout doesn't break

        browser.close()

if __name__ == "__main__":
    verify_new_features()
