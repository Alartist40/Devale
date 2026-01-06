import logging
import os
import yaml
import subprocess
import threading
from .explain import get_explanation
from .helpers import get_resource_path

class CommandRunner:
    def __init__(self, dry_run=False):
        self.dry_run = dry_run
        self.logger = logging.getLogger("DevAleRunner")
        self.log_callback = None

    def set_logger(self, callback):
        self.log_callback = callback

    def log(self, msg):
        self.logger.info(msg)
        if self.log_callback:
            self.log_callback(msg)

    def run_as_admin(self, cmd_str):
        """Runs a command as administrator using ShellExecute"""
        import ctypes
        import sys
        
        self.log(f"ELEVATING: {cmd_str}")
        try:
             # This will trigger UAC and run the command in a NEW window
             # /k keeps it open, /c closes. For background tasks like Restore Point, /c is better
             # but user might want to see output.
             # We'll use /c and rely on the command itself to pause if needed, or just let it run.
             # For Restore Point: `wmic /Namespace:\\root\default Path SystemRestore Call CreateRestorePoint "DevAle Manual Point", 100, 7`
             
             params = f"/c {cmd_str}"
             ctypes.windll.shell32.ShellExecuteW(None, "runas", "cmd.exe", params, None, 1)
             self.log("Admin request sent.")
             return True
        except Exception as e:
            self.log(f"Elevation Failed: {e}")
            return False

    def run_interactive_cmd(self, cmd_str):
        """Run a raw command string from the console input"""
        step = {'cmd': cmd_str, 'friendly': f"User Command: {cmd_str}", 'timeout': 300}
        # Run in thread to not block UI
        threading.Thread(target=self.run_step, args=(step,), daemon=True).start()

    def load_commands(self, file_path):
        # Resolve full path
        full_path = get_resource_path(file_path)
        if not os.path.exists(full_path):
            self.logger.error(f"Cannot find command file: {full_path} (orig: {file_path})")
            return []
            
        with open(full_path, 'r') as f:
            return yaml.safe_load(f)

    def run_step(self, step, progress_callback=None):
        cmd = step.get('cmd')
        tag = step.get('tag', 'unknown')
        friendly = step.get('friendly')
        timeout = step.get('timeout', 60)
        
        explanation = friendly if friendly else get_explanation(tag)
        
        if progress_callback:
            progress_callback(explanation)
            
        self.log(f"EXEC: {explanation}")
        
        if self.dry_run:
            self.log(f"[DRY] {cmd}")
            return True, "Dry run"

        try:
            # Prepare command
            shell_cmd = cmd
            if isinstance(cmd, list):
                # Join simple list commands for shell execution if needed, or pass list
                # heuristic: prefer string for complex shell pipes
                 shell_cmd = " ".join(cmd)
            
            self.log(f"CMD: {shell_cmd}")

            process = subprocess.Popen(
                shell_cmd,
                shell=True,
                stdout=subprocess.PIPE,
                stderr=subprocess.STDOUT, # Merge stderr to stdout
                text=False, # Read bytes to handle encoding manually
                bufsize=1   # Line buffering
            )

            captured_output = []
            
            # Stream output
            for line in iter(process.stdout.readline, b''):
                try:
                    # Windows usually uses CP437 or CP850 for console output, but Python defaults to utf-8
                    # We try utf-8 first, then cp850, then replace
                    decoded_line = line.decode('utf-8').strip()
                except UnicodeDecodeError:
                    try:
                        decoded_line = line.decode('cp850').strip()
                    except:
                        decoded_line = line.decode('utf-8', errors='replace').strip()
                
                if decoded_line:
                    self.log(decoded_line)
                    captured_output.append(decoded_line)

            process.wait(timeout=timeout)
            
            if process.returncode == 0:
                self.log("SUCCESS")
                return True, "\n".join(captured_output)
            else:
                self.log(f"FAILED (Code {process.returncode})")
                return False, "\n".join(captured_output)

        except subprocess.TimeoutExpired:
            self.log("TIMEOUT")
            return False, "Command timed out"
        except Exception as e:
            self.log(f"ERROR: {e}")
            return False, str(e)

    def run_recipe(self, recipe_path, progress_callback=None):
        steps = self.load_commands(recipe_path)
        if not steps:
             return

        total = len(steps)
        
        for i, step in enumerate(steps):
             if progress_callback:
                 # Pass a tuple or percentage? Let's just pass the text for now
                 # Caller can handle progress bar math if they know total, 
                 # but for simplicity let's just pass the text.
                 pass
             
             success, output = self.run_step(step, progress_callback)
             if not success:
                 self.logger.error(f"Step failed: {output}")
                 # Depends if we want to stop or continue. 
                 # For panic button, maybe continue best effort?
                 continue

    def run_native_cmd(self, cmd_list):
        """Run a command using list args and shell=False (Legacy compatibility)"""
        self.log(f"CMD (Native): {' '.join(cmd_list)}")
        try:
            # shell=False means we execute the executable directly
            # capture_output=True equivalent to pipe stdout/stderr
            # Use STARTUPINFO to hide window (safer than CREATE_NO_WINDOW)
            startupinfo = subprocess.STARTUPINFO()
            startupinfo.dwFlags |= subprocess.STARTF_USESHOWWINDOW
            startupinfo.wShowWindow = subprocess.SW_HIDE
            
            process = subprocess.run(
                cmd_list,
                capture_output=True,
                text=True,
                startupinfo=startupinfo,
                timeout=300
            ) 
            
            if process.returncode == 0:
                self.log("SUCCESS")
                self.log(process.stdout)
                return True, process.stdout
            else:
                self.log(f"FAILED (Code {process.returncode})")
                self.log(process.stderr)
                return False, process.stderr
                
        except Exception as e:
            self.log(f"ERROR: {e}")
            return False, str(e)
