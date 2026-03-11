# Changelog - DevAle v2

All notable changes to the DevAle project during the Go rewrite are documented here.

## [2.0.0] - Go Rewrite (v2-go-rewrite)

### Added
- **Core Engine:** Total rewrite in Go for 10x faster startup and zero dependency hell.
- **Repair System:** Ported and optimized 6-phase Windows repair logic from TBOK batch scripts.
- **Persistence:** Added JSON-based state tracking in `%APPDATA%` to handle reboot cycles.
- **Auto-Resume:** Implemented Windows Scheduled Task integration to auto-resume repairs after Phase 2 (CHKDSK).
- **GUI:** New modern interface built with Wails, Svelte, and TypeScript.
- **Neumorphic Design:** Implemented "Soft UI" aesthetics as requested.
- **Theme Engine:** Full Light and Dark mode toggle with persistent state.
- **Telemetry:** Hardware monitoring (CPU/RAM/GPU/Storage) using the `ghw` library.
- **Terminal:** Integrated command-line console with real-time output streaming and CP850-to-UTF8 decoding.
- **App Store:** Winget integration for common software installation.
- **UAC:** Mandatory administrator elevation via application manifest.

### Changed
- **Architecture:** Switched from `customtkinter` (Python) to `Wails` (Go/TS/Webkit).
- **Styling:** Moved from Cyberpunk/Neon to a sophisticated Neumorphic dashboard style.
- **Binary:** Output is now a single compressed 11MB EXE instead of a bulky PyInstaller folder.

### Fixed
- Fixed command output truncation in the console using a dedicated Go reader thread.
- Resolved encoding issues with Windows command results (UTF-8 conversion).
- Improved the reliability of WMI repairs by automating MOF recompilation.
- **Terminal UX:** Filtered out Windows shell headers (copyright, version) from the integrated terminal.
- **UX Refinement:** Replaced default system scrollbars with custom-themed Neumorphic scrollbars for a consistent "Cyberpunk-lite" feel.
- **Logging Optimization:** Suppressed verbose per-file command output (DLL/MOF/DISM progress) to prevent terminal "glitching" and replaced it with concise "Loading..." sequences.
- **Reliability:** Added manual "Reset Repair State" capability to recover from interrupted processes.
- **Bug Fix:** Fixed AppX re-registration command syntax error and ensured the "Panic Button" correctly resets its active state.

## [2.1.0] - Resilience & App Store Enhancements

### Added
- **Mock Windows Environment:** Introduced `Commander` interface and `MockCommander` for simulating Windows system responses in testing/CI.
- **Categorized App Store:** Transitioned to a data-driven `applications.json` with categories (Browsers, Developer Tools, Gaming, Utilities) and expanded software list (Unreal, Steam, etc.).
- **Safety Backups:** Integrated automatic Windows System Restore Point creation before starting repairs.
- **Security:** Added Administrator privilege verification on startup.
- **Diagnostics:** implemented Disk Health monitoring via WMIC and detailed Host/OS telemetry.
- **Terminal UX:** Overhauled with a typed, color-coded logging system (Phase, Error, Success, User, Highlight).
- **Partition Map:** Added visual storage breakdown for all system partitions.
- **Pre-Flight Dashboard:** Added real-time Battery, Pending Updates, and Network Ping indicators to the Home screen.
- **Power-User Kit:** Integrated Driver Audit, Startup Manager, and Disk Management shortcuts.
- **Log Export:** Added capability to export terminal history to `.txt`.
- **Animations:** Added pulsing glow animations to the Panic Button during active repairs.
- **God Mode:** Added Windows "Master Settings" (God Mode) shortcut for power-user control.
- **Video Compatibility:** Refactored window rendering to support screen capture and recording for presentations.

### Changed
- **Persistence:** Enhanced state tracking to include `LastStep` and `Logs` for more granular progress monitoring.
- **Architecture:** Implemented step-level resumption to skip completed commands within a phase.
- **Security:** Refactored command execution to include shell sanitization and trust boundaries.
- **Diagnostics UI:** Replaced static text with real-time CPU/RAM usage progress bars.

### Fixed
- **Phase 5 Reliability:** Ensured WMI services are always re-enabled even if intermediate repair steps fail or are cancelled.
- **Resource Management:** Added proactive context cancellation checks and null-safety to the command runner.

## [2.1.1] - Maintenance

### Changed
- **Build:** Rebuilt the production Windows executable (`DevAle_v2.exe`) to include all recent code updates and UI refinements.
