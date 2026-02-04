## 2026-01-07 - [Subprocess Resource Management]
**Learning:** Using `subprocess.Popen` with `text=False` and `bufsize=1` triggers a `RuntimeWarning` in Python 3.12+ as line-buffering is only supported in text mode. Additionally, neglecting to explicitly close `process.stdout` when manually iterating over it can lead to `ResourceWarning` about unclosed files.
**Action:** Always use `bufsize=-1` (default) for binary mode subprocesses and ensure streams are closed using a `finally` block or context manager.

## 2026-01-07 - [Process Execution Abstraction]
**Learning:** Orchestrator methods (like `run_step`) become bloated when they mix high-level configuration (explanation lookup, dry-run checks) with low-level process management (subprocess creation, threading, timeout handling).
**Action:** Extract low-level subprocess management into dedicated private methods to reduce cyclomatic complexity and improve modularity.

## 2026-01-07 - [Data-Driven Navigation]
**Learning:** Hardcoded GUI navigation (buttons, frames) increases redundancy and makes scaling difficult. Centralizing UI structure into a configuration list (e.g., `NAV_ITEMS`) allows for dynamic generation of UI elements, reducing boilerplate and ensuring consistent layout behavior (like spacers).
**Action:** Use data-driven patterns for UI components that share a common structure or behavior to improve maintainability and scalability.

## 2026-01-07 - [GUI Structural Testing]
**Learning:** Testing GUI applications in headless environments requires extensive mocking of the GUI library (like `customtkinter`). A robust strategy involves creating base mock classes with common methods (`grid`, `pack`, `configure`) to avoid `AttributeError` and `StopIteration` during complex widget initialization.
**Action:** Use minimal mock classes that provide necessary GUI method stubs when performing structural verification of GUI classes.
