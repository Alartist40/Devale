# Product Requirement Document (PRD) - DevAle v2 (Go Rewrite)

## 1. Overview
DevAle v2 is a total rewrite of the original Python-based DevAle system utility. The goal is to move from a high-dependency, slower Python environment to a high-performance, single-binary Go environment with a modern TypeScript-based GUI (Wails + Svelte).

## 2. Goals & Principles
- **Efficiency:** Fast execution and accurate results.
- **Portability:** Single EXE with no external runtime dependencies.
- **TypeScript First:** UI logic written in TypeScript for type safety.
- **Sleek Aesthetic:** Modern, minimal, Cyberpunk-inspired Neumorphic UI.
- **Offline Reliability:** Core maintenance features work without internet.
- **Admin First:** The application must run with elevated privileges to perform system repairs.

## 3. Architecture
- **Backend:** Go (Golang) for system-level execution and telemetry.
- **Frontend:** Wails Framework (Svelte + TypeScript) with Neumorphic CSS3.
- **Communication:** Wails Bridge (Binding Go methods to TS) with real-time event streaming.
- **Security:** Windows Manifest set to `requireAdministrator`.

## 4. Key Features

### 4.1. Repair Orchestrator (The "Panic" Button)
- **Phase 1: DISM Health Check** (CheckHealth, ScanHealth, RestoreHealth).
- **Phase 2: CHKDSK Scheduling** (Schedule disk check for next boot).
- **Phase 3: SFC Scan** (System File Checker).
- **Phase 4: Component Store Cleanup** (StartComponentCleanup).
- **Phase 5: WMI Repair** (Aggressive re-registration of DLLs and MOFs).
- **Phase 6: AppX Re-registration** (Fixing Windows Store apps).
- **Resumption:** Intelligent step-level progress tracking and resumption logic.
- **Safety:** Automatic System Restore Point creation before repairs.

### 4.2. Advanced Diagnostics & Telemetry
- **Hardware Monitoring:** CPU, RAM, GPU, and Storage health (via ghw).
- **Real-time Usage:** Visual progress bars for CPU and RAM consumption.
- **Partition Map:** Visual breakdown of all system drives and space usage.
- **Pre-Flight Checks:** Battery level, Pending Updates detection, and Network Ping status.

### 4.3. IT Power-User Module
- **Driver Audit:** Detect hardware devices with errors.
- **Startup Manager:** Audit programs starting with Windows.
- **Disk Management:** Direct shortcut to system disk management.
- **Log Export:** Save terminal history to `.txt` for professional reporting.

### 4.4. Categorized App Store (Winget)
- Curated, data-driven application list (`applications.json`).
- Categories: Browsers, Developer Tools, Gaming, Utilities.

### 4.5. Integrated Terminal
- Interactive shell for manual command execution with color-coded logging.
- Real-time output streaming for all automated actions.

## 5. UI/UX Requirements
- **Theme:** Light/Dark mode, Neumorphic "Soft UI" aesthetic with a minimalist branding approach (no "Cynapse" text).
- **Animations:** Pulsing "Panic" glow during active repairs.
- **Accessibility:** ARIA labels, clear visual feedback for active states.
- **Terminal:** Robust CP850-to-UTF8 encoding to prevent garbled text in the integrated console.
