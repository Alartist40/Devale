# DevAle v2 (Go Rewrite)

DevAle is a high-performance system utility for Windows, rebuilt from the ground up in Go and TypeScript. It provides a robust 6-phase automated repair sequence, a guided IT support module, an app store (Winget), and an integrated terminal.

## Key Features
- **🚀 Efficiency:** Rebuilt in Go for near-instant execution and minimal resource usage.
- **🛡️ 6-Phase Repair Orchestrator:**
  1. DISM Health Check & Repair
  2. CHKDSK Scheduling (Auto-Resume after restart)
  3. System File Checker (SFC)
  4. Component Store Cleanup
  5. WMI Repository Repair
  6. AppX (Windows Store) Re-registration
- **📦 Neural App Store:** Curated selection of essential apps via Winget.
- **📟 Integrated Terminal:** Real-time log streaming and interactive shell.
- **🎨 Cyberpunk Aesthetic:** Sleek, minimal, high-tech dark theme.
- **🔐 Security:** Automatically requests Administrative privileges for system tasks.

## Architecture
- **Backend:** Go (Golang)
- **Frontend:** Wails Framework (Svelte + TypeScript)
- **State:** JSON-based persistence with `schtasks` for post-reboot recovery.

## Development Setup
1. **Prerequisites:**
   - Go 1.20+
   - Node.js & pnpm
   - Wails CLI (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`)

2. **Run in Dev Mode:**
   ```bash
   wails dev
   ```

3. **Build for Production:**
   ```bash
   wails build
   ```

## License
MIT License
