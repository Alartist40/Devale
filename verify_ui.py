from playwright.sync_api import sync_playwright
import time

def verify_ui_improvements():
    with sync_playwright() as p:
        browser = p.chromium.launch(headless=True)
        page = browser.new_page()
        # Navigate to the local dev server (assuming it's running on 34115 from previous steps)
        # In this environment, I'll just use the file path to the built frontend if possible,
        # but usually wails dev serves it. I'll try the last known port.
        try:
            page.goto("http://localhost:34115")
            time.sleep(2)

            # Inject some logs to trigger scrollbar
            page.evaluate("""
                const logs = document.getElementById('terminal-logs');
                for(let i=0; i<50; i++) {
                    const div = document.createElement('div');
                    div.textContent = 'LOG LINE ' + i;
                    logs.appendChild(div);
                }
                logs.scrollTop = logs.scrollHeight;
            """)
            time.sleep(1)

            page.screenshot(path="verify_terminal_scrollbar.png")
            print("Screenshot saved to verify_terminal_scrollbar.png")

        except Exception as e:
            print(f"Error: {e}")
        finally:
            browser.close()

if __name__ == "__main__":
    verify_ui_improvements()
