# Changelog

## [2.0.0] - 2025-03-04
### Added
- **Complete Rewrite**: Migrated entire codebase from Python to **Go** for maximum performance.
- **Wails Framework**: Integrated Wails v2 for native Windows GUI capabilities.
- **Cyberpunk UI**: New Svelte/TypeScript frontend with a sleek modern aesthetic.
- **6-Phase Repair Orchestrator**: Ported and enhanced 400+ lines of batch script logic into the Go backend.
- **Persistence Layer**: Implemented phase tracking and `schtasks` integration for resuming repairs after system reboots.
- **Admin Elevation**: Configured application manifest to automatically request UAC administrator rights.
- **Integrated Terminal**: Added a real-time command output stream and interactive shell.
- **Hardware Telemetry**: Real-time system specs (CPU, GPU, RAM, Disk) using `ghw`.
- **Winget App Store**: Curated app installation module.

### Changed
- Improved binary size to ~11MB.
- Switched to `pnpm` for lightning-fast frontend builds.
- Moved from CustomTkinter to a web-based UI stack for better compatibility and design flexibility.

### Fixed
- Eliminated "dependency hell" by bundling everything into a single, standalone EXE.
- Fixed encoding issues in command output streaming (UTF-8/CP850 support).
