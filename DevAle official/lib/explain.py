
EXPLAIN = {
    "clean_temp": "Cleaning temporary files to free up space.",
    "clean_cache": "Removing application caches.",
    "scan_virus": "Running a quick virus scan.",
    "scan_sfc": "Checking system files for corruption (SFC).",
    "scan_dism": "Repairing Windows image health (DISM).",
    "restore_point": "Creating a system restore point backup.",
    "disable_telemetry": "Disabling Windows telemetry.",
    "show_extensions": "Showing file extensions in Explorer.",
    "show_hidden": "Showing hidden files in Explorer.",
    "restart_explorer": "Restarting Windows Explorer to apply changes.",
    "update_security": "Checking for security updates.",
    "flush_dns": "Flushing DNS cache to fix network issues.",
    "reset_network": "Resetting network adapters.",
    "optimize_drives": "Optimizing storage drives.",
}

def get_explanation(tag):
    return EXPLAIN.get(tag, f"Running action: {tag}")
