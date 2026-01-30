## 2026-01-07 - [Subprocess Resource Management]
**Learning:** Using `subprocess.Popen` with `text=False` and `bufsize=1` triggers a `RuntimeWarning` in Python 3.12+ as line-buffering is only supported in text mode. Additionally, neglecting to explicitly close `process.stdout` when manually iterating over it can lead to `ResourceWarning` about unclosed files.
**Action:** Always use `bufsize=-1` (default) for binary mode subprocesses and ensure streams are closed using a `finally` block or context manager.
