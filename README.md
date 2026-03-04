# DevAle v2 - Go Rewrite

DevAle v2 is a modern, high-performance system maintenance utility rewritten from the ground up. It transitions from a Python-based architecture to a single-binary Go backend with a TypeScript/Svelte frontend using the Wails framework.

## 🚀 Key Features
- **6-Phase Repair Orchestrator:** Automated system restoration (DISM, CHKDSK, SFC, Component Cleanup, WMI Repair, AppX Re-registration).
- **Reboot & Resume:** Intelligent persistence that handles system restarts and resumes maintenance automatically.
- **Neumorphic Dashboard:** A sleek, soft-UI design with full support for Light and Dark modes.
- **Hardware Telemetry:** Real-time monitoring of CPU, RAM, GPU, and Storage health.
- **Integrated Terminal:** Real-time logging and manual command execution.
- **App Store:** One-click installations using Microsoft Winget.
- **Zero Dependencies:** Compiled to a single standalone executable.

## 🏗️ Architecture
- **Backend:** Go (Golang) for system-level execution and subprocess management.
- **Frontend:** TypeScript, Svelte, and CSS3 (Neumorphic Design).
- **Communication:** Wails runtime bridge for asynchronous IPC.
- **Security:** Integrated Windows Manifest for mandatory UAC elevation.

## 🛠️ Setup & Build

### ⚡ Easy Way (Windows)
Just run the included batch script. It will check for Go and Node.js, install Wails if missing, and build the final EXE for you:
```cmd
easy_build.bat
```

### 🛠️ Manual Way
#### Prerequisites
- [Go](https://go.dev/) (1.21+)
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

## 📂 Project Structure
- `main.go`: Application entry point and Wails configuration.
- `app.go`: Go-to-Frontend bindings and business logic.
- `runner.go`: Subprocess orchestration and repair phase logic.
- `sysinfo.go`: Hardware telemetry collection.
- `state.go`: Persistence layer for repair status.
- `frontend/`: Svelte source code, Neumorphic styles, and TypeScript components.

## 📜 License
MIT License. Created with ❤️ for I.T. Support efficiency.
