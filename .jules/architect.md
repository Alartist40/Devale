## 2026-01-07 - [Subprocess Resource Management]
**Learning:** Using `subprocess.Popen` with `text=False` and `bufsize=1` triggers a `RuntimeWarning` in Python 3.12+ as line-buffering is only supported in text mode. Additionally, neglecting to explicitly close `process.stdout` when manually iterating over it can lead to `ResourceWarning` about unclosed files.
**Action:** Always use `bufsize=-1` (default) for binary mode subprocesses and ensure streams are closed using a `finally` block or context manager.

## 2026-01-07 - [Process Execution Abstraction]
**Learning:** Orchestrator methods (like `run_step`) become bloated when they mix high-level configuration (explanation lookup, dry-run checks) with low-level process management (subprocess creation, threading, timeout handling).
**Action:** Extract low-level subprocess management into dedicated private methods to reduce cyclomatic complexity and improve modularity.

## 2026-01-07 - [Data-Driven GUI Navigation]
**Learning:** Hardcoding GUI navigation buttons and frame initialization leads to significant code duplication and makes scaling difficult. A data-driven approach using a configuration list (e.g., `NAV_ITEMS`) allows for dynamic generation of both UI components and backing frames, improving maintainability and reducing the risk of inconsistent layout logic.
**Action:** Use a centralized configuration list for navigation items and loop over it to initialize frames and create sidebar buttons.
