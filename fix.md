# Fixes and Recommendations - DevAle v2

## Implemented Fixes

### 1. WMI Phase Resilience
- **Issue:** If the WMI repair phase was interrupted while the service was disabled, it could leave the system in a broken state.
- **Fix:** Added a `defer` mechanism in `runner.go` that ensures `winmgmt` is re-enabled and started regardless of whether the phase completes, fails, or is cancelled.

### 2. Categorized App Store
- **Issue:** Applications were hardcoded in the frontend, making it difficult to maintain and lacking organization.
- **Fix:**
    - Moved application management to `backend/runner/applications.json`.
    - Implemented categories (Browsers, Developer Tools, Game Development, Utilities).
    - Updated the Svelte UI to dynamically render categories and apps.

### 3. Testing Infrastructure
- **Improvement:** Introduced a `Commander` interface and `MockCommander`.
- **Value:** Allows developers to test Windows-specific logic on non-Windows systems (like this Linux devbox) and simulate complex failure scenarios (like CHKDSK's exit status 3).

### 4. Progress Persistence
- **Improvement:** Expanded the `State` struct in `persistence/state.go`.
- **Value:** Now supports tracking the last executed step and potentially log history, allowing for much more granular "Resume" functionality in future updates.

## Recommendations for Future Improvement

### 1. Command Security (Security)
- **Problem:** Currently, the app is vulnerable to command injection because it passes raw strings to `cmd /c`.
- **Recommendation:** Implement a whitelist of allowed commands or use a structured argument passing system (e.g., `exec.Command("winget", "install", id)`) instead of raw string interpolation.

### 2. Restart Cancellation Detection (UX)
- **Problem:** If a user cancels the 10-second restart timer via `shutdown /a`, the app remains in a "Waiting for Restart" state.
- **Recommendation:** Implement a polling mechanism or a "Check for Pending Restart" function that verifies if the CHKDSK task is still scheduled and if the system actually rebooted (e.g., by checking system uptime).

### 3. Granular Step Resumption (Architecture)
- **Problem:** If Phase 5 (WMI) fails on step 4, the app currently restarts from the beginning of Phase 5.
- **Recommendation:** Use the new `LastStep` field in the persistence state to skip already completed steps within a phase.

### 4. Admin Privilege Check (Robustness)
- **Problem:** Although the manifest requests Admin, the code doesn't explicitly verify it at runtime.
- **Recommendation:** Add a check on startup using `os.Open("\\\\.\\PHYSICALDRIVE0")` or similar Windows-specific calls to confirm true Administrator elevation and show a clear error if missing.
