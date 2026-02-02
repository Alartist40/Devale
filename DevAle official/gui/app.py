import customtkinter as ctk
from .frames.home import HomeFrame
from .frames.diagnose import DiagnoseFrame
from .frames.tools import ToolsFrame
from .frames.install import InstallFrame
import os
import sys

class DevAleGUI(ctk.CTk):
    def __init__(self, runner):
        super().__init__()
        
        self.runner = runner
        self.title("DevAle 1.0")
        self.geometry("900x650")

        # Navigation Configuration - Architectural Rationale: Centralized
        # configuration for easy maintenance and scalability.
        self.NAV_ITEMS = [
            ("Home", "🏠 Home", HomeFrame),
            ("Diagnose", "🏥 Diagnose", DiagnoseFrame),
            ("Tools", "🛠 Tools", ToolsFrame),
            ("Install", "📦 Install", InstallFrame),
        ]
        
        # Grid layout 1x2 (Sidebar + Content)
        self.grid_columnconfigure(1, weight=1)
        self.grid_rowconfigure(0, weight=1) # Content
        self.grid_rowconfigure(1, weight=0) # Console (fixed height)
        
        self.setup_sidebar()
        
        # Frames container
        self.frames = {}
        self.current_frame = None
        
        # Initialize Frames dynamically
        for name, _, frame_class in self.NAV_ITEMS:
            self.frames[name] = frame_class(self, self.runner)
        
        # Console / Mini Terminal
        self.console_frame = ctk.CTkFrame(self, height=150, corner_radius=0)
        self.console_frame.grid(row=1, column=1, sticky="nsew", padx=0, pady=0)
        self.console_frame.grid_propagate(False) # Fixed height
        
        self.console_label = ctk.CTkLabel(self.console_frame, text="System Console", font=ctk.CTkFont(size=12, weight="bold"))
        self.console_label.pack(side="top", anchor="w", padx=5, pady=2)
        
        self.console_log = ctk.CTkTextbox(self.console_frame, font=ctk.CTkFont(family="Consolas", size=12))
        self.console_log.pack(side="top", fill="both", expand=True, padx=5, pady=(5, 0))
        self.console_log.configure(state="disabled")

        # Input Entry
        self.input_entry = ctk.CTkEntry(self.console_frame, font=ctk.CTkFont(family="Consolas", size=12), placeholder_text="Type command here...")
        self.input_entry.pack(side="bottom", fill="x", padx=5, pady=5)
        self.input_entry.bind("<Return>", self.run_custom_command)

        # Connect Logger
        self.runner.set_logger(self.log)
        self.log(f"DevAle Initialized. OS: {os.name} | Platform: {sys.platform}")

        self.last_log_line = ""

        self.show_frame("Home")

    def run_custom_command(self, event):
        cmd = self.input_entry.get()
        if cmd.strip():
            self.input_entry.delete(0, 'end')
            self.runner.run_interactive_cmd(cmd)

    def log(self, message):
        """Thread-safe logging with filtering"""
        def _write():
            try:
                # Spam filter: Don't print "Verification x%" lines repeatedly
                msg_str = str(message)
                if "Verification" in msg_str and "%" in msg_str:
                    # Optional: update last line instead of appending? 
                    # For now, just skip every other one or ignore if identical to last?
                    # Simply reducing spam: if it contains "Verification", skip it entirely or print only 10% increments?
                    # Let's just skip it to keep console clean for now, as user requested "no loop".
                    # Better: Print only 0%, 50%, 100%? 
                    # User said "don't want it to print out all the time creating a loop".
                    return

                self.console_log.configure(state="normal")
                self.console_log.insert("end", f"> {message}\n")
                self.console_log.see("end")
                self.console_log.configure(state="disabled")
            except Exception as e:
                # Silently catch GUI errors to prevent crashes during logging
                # but use explicit Exception to comply with Ruff BLE001.
                pass
        
        self.after(0, _write)

    def setup_sidebar(self):
        self.sidebar = ctk.CTkFrame(self, width=200, corner_radius=0)
        self.sidebar.grid(row=0, column=0, rowspan=2, sticky="nsew") # Span both content and console rows
        
        title_label = ctk.CTkLabel(self.sidebar, text="DevAle", font=ctk.CTkFont(size=20, weight="bold"))
        title_label.grid(row=0, column=0, padx=20, pady=(20, 10))

        # Navigation Buttons - Architectural Rationale: Data-driven generation
        # reduces code duplication and improves maintainability.
        self.nav_buttons = {}
        for i, (name, label, _) in enumerate(self.NAV_ITEMS, start=1):
            btn = ctk.CTkButton(
                self.sidebar, text=label,
                command=lambda n=name: self.show_frame(n),
                fg_color="transparent", text_color=("gray10", "#DCE4EE"),
                border_width=2, border_color=("gray", "gray")
            )
            btn.grid(row=i, column=0, padx=20, pady=10, sticky="ew")
            self.nav_buttons[name] = btn

        # Spacer row with weight 1 to push following items to the bottom
        spacer_row = len(self.NAV_ITEMS) + 1
        self.sidebar.grid_rowconfigure(spacer_row, weight=1)
        
        # Appearance Mode
        self.appearance_mode_label = ctk.CTkLabel(self.sidebar, text="Theme:", anchor="w")
        self.appearance_mode_label.grid(row=spacer_row + 1, column=0, padx=20, pady=(10, 0))
        self.appearance_mode_optionemenu = ctk.CTkOptionMenu(
            self.sidebar, values=["System", "Light", "Dark"],
            command=self.change_appearance_mode_event
        )
        self.appearance_mode_optionemenu.grid(row=spacer_row + 2, column=0, padx=20, pady=(10, 20))

    def change_appearance_mode_event(self, new_appearance_mode: str):
        ctk.set_appearance_mode(new_appearance_mode)
        
    def show_frame(self, name):
        if self.current_frame:
            self.current_frame.grid_forget()
            
        frame = self.frames[name]
        frame.grid(row=0, column=1, sticky="nsew", padx=20, pady=20)
        self.current_frame = frame
