# DevAle v2 - Go Rewrite

DevAle v2 is a modern, high-performance system maintenance utility rewritten from the ground up. It transitions from a Python-based architecture to a single-binary Go backend with a TypeScript/Svelte frontend using the Wails framework.

## 🚀 Key Features
- **6-Phase Repair Orchestrator:** Automated system restoration (DISM, CHKDSK, SFC, Component Cleanup, WMI Repair, AppX Re-registration).
- **Intelligent Step-Level Resumption:** Skip already-completed steps within a phase to save time.
- **Safety Backups:** Automatic System Restore Point creation before starting repairs.
- **Advanced Diagnostics:** Real-time CPU/RAM progress bars, detailed Partition Map, and Disk Health monitoring (WMIC).
- **Pre-Flight Dashboard:** Real-time monitoring of Battery, Pending Updates, and Network Ping.
- **Power-User Kit:** Integrated Driver Audit, Startup Manager, and Disk Management shortcuts.
- **Neumorphic Dashboard:** A sleek, "Soft-UI" design with Full Light/Dark mode and pulsating repair animations.
- **Categorized App Store:** One-click installations using Microsoft Winget for Browsers, Dev Tools, and Gaming.
- **Integrated Terminal:** Real-time color-coded logging and manual command execution.
- **Zero Dependencies:** Compiled to a single standalone executable.

## 🏗️ Architecture
- **Backend:** Go (Golang) for system-level execution and subprocess management.
- **Frontend:** TypeScript, Svelte, and CSS3 (Neumorphic Design).
- **Communication:** Wails runtime bridge for asynchronous IPC and event streaming.
- **Security:** Integrated Windows Manifest for mandatory UAC elevation and internal command trust boundaries.

## 🛠️ Setup & Build

### ⚡ Easy Way (Windows)
Just run the included batch script. It will check for Go and Node.js, install Wails if missing, and build the final EXE for you:
```cmd
easy_build.bat
```

### 🛠️ Manual Way
#### Prerequisites
- [Go](https://go.dev/) (1.24+)
- [Node.js](https://nodejs.org/) & [NPM](https://www.npmjs.com/)
- [Wails CLI](https://wails.io/docs/gettingstarted/installation)

#### Development
```bash
wails dev
```

#### Production Build
To build a standalone Windows executable:
```bash
wails build -platform windows/amd64
```
The output binary will be located in `build/bin/devale-v2.exe`.

**Note for Developers:** This build includes the latest fixes for command cancellation and UI state persistence.

## 📂 Project Structure
- `main.go`: Application entry point and Wails configuration.
- `app.go`: Go-to-Frontend bindings and business logic.
- `backend/runner/`: Subprocess orchestration, repair phase logic, and application management.
- `backend/sysinfo/`: Hardware telemetry and system health collection.
- `backend/persistence/`: Persistence layer for repair status.
- `frontend/`: Svelte source code, Neumorphic styles, and TypeScript components.

## 📜 License
MIT License. Created with ❤️ for I.T. Support efficiency.
