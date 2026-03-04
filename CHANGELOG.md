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
