# Product Requirement Document (PRD) - DevAle v2 (Go Rewrite)

## 1. Overview
DevAle v2 is a total rewrite of the original Python-based DevAle system utility. The goal is to move from a high-dependency, slower Python environment to a high-performance, single-binary Go environment with a modern TypeScript-based GUI (Wails + Svelte).

## 2. Goals & Principles
- **Efficiency:** Fast execution and accurate results.
- **Portability:** Single EXE with no external runtime dependencies.
- **TypeScript First:** UI logic written in TypeScript for type safety.
- **Sleek Aesthetic:** Modern, minimal, cyberpunk-inspired UI.
- **Offline Reliability:** Core maintenance features work without internet.
- **Admin First:** The application must run with elevated privileges to perform system repairs.

## 3. Architecture
- **Backend:** Go (Golang)
- **Frontend:** Wails Framework (Svelte + TypeScript)
- **Communication:** Wails Bridge (Binding Go methods to TS)
- **Security:** Windows Manifest set to `requireAdministrator`.

## 4. Key Features

### 4.1. Repair Orchestrator (The "Panic" Button)
Porting the sophisticated logic from the TBOK repair script:
- **Phase 1: DISM Health Check** (CheckHealth, ScanHealth, RestoreHealth).
- **Phase 2: CHKDSK Scheduling** (Schedule disk check for next boot).
- **Phase 3: SFC Scan** (System File Checker).
- **Phase 4: Component Store Cleanup** (StartComponentCleanup).
- **Phase 5: WMI Repair** (Aggressive re-registration of DLLs and MOFs).
- **Phase 6: AppX Re-registration** (Fixing Windows Store apps).

### 4.2. Reboot & Resume Mechanism
- **State Persistence:** Save the current phase and repair status to a local JSON/enc file.
- **Auto-Resume:** On restart, if a repair is in progress, the app should auto-start via a temporary Scheduled Task or `RunOnce` registry key.

### 4.3. IT Support Module
- Guided troubleshooting for non-tech users.
- Automated storage analysis and "PC Slow" fixes.

### 4.4. App Store (Winget)
- Use the existing `applications.json` to provide a curated list of apps.
- Batch installation via `winget`.

### 4.5. Integrated Terminal
- Interactive shell for manual command execution.
- Real-time output streaming for all automated actions.

## 5. UI/UX Requirements
- **Theme:** Dark mode, neon cyan/magenta accents.
- **Layout:** Sidebar nav, main content, fixed terminal area at the bottom.
- **Accessibility:** ARIA labels for buttons, clear visual feedback for active states.

## 6. Delivery
- A single `DevAle.exe` produced via `wails build`.
- Manifest configured for UAC elevation.
