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
- **Categorized App Store:** Transitioned to a data-driven `applications.json` with categories (Browsers, Developer Tools, Game Development, Utilities).
- **Repair Resilience:** Implemented `defer` blocks in critical repair phases (WMI) to ensure system services are re-enabled if the app is interrupted.

### Changed
- **Persistence:** Enhanced state tracking to include `LastStep` and `Logs` for more granular progress monitoring.

### Fixed
- **Phase 5 Reliability:** Ensured WMI services are always re-enabled even if intermediate repair steps fail or are cancelled.
