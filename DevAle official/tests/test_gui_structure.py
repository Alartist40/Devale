import unittest
from unittest.mock import MagicMock, patch
import sys
import os

# Add the project root to sys.path
sys.path.append(os.path.join(os.getcwd(), "DevAle official"))

# Mock customtkinter components
class MockCTk:
    def __init__(self, *args, **kwargs): pass
    def grid_columnconfigure(self, *args, **kwargs): pass
    def grid_rowconfigure(self, *args, **kwargs): pass
    def after(self, *args, **kwargs): pass
    def title(self, *args, **kwargs): pass
    def geometry(self, *args, **kwargs): pass
    def mainloop(self, *args, **kwargs): pass
    def bind(self, *args, **kwargs): pass

class MockFrame:
    def __init__(self, *args, **kwargs): pass
    def grid(self, *args, **kwargs): pass
    def grid_columnconfigure(self, *args, **kwargs): pass
    def grid_rowconfigure(self, *args, **kwargs): pass
    def grid_propagate(self, *args, **kwargs): pass
    def pack(self, *args, **kwargs): pass
    def grid_forget(self, *args, **kwargs): pass

class MockButton:
    def __init__(self, *args, **kwargs): pass
    def grid(self, *args, **kwargs): pass
    def pack(self, *args, **kwargs): pass

class MockLabel:
    def __init__(self, *args, **kwargs): pass
    def grid(self, *args, **kwargs): pass
    def pack(self, *args, **kwargs): pass

class MockTextbox:
    def __init__(self, *args, **kwargs): pass
    def pack(self, *args, **kwargs): pass
    def configure(self, *args, **kwargs): pass
    def insert(self, *args, **kwargs): pass
    def see(self, *args, **kwargs): pass

class MockEntry:
    def __init__(self, *args, **kwargs): pass
    def pack(self, *args, **kwargs): pass
    def bind(self, *args, **kwargs): pass
    def get(self, *args, **kwargs): return ""
    def delete(self, *args, **kwargs): pass

mock_ctk = MagicMock()
mock_ctk.CTk = MockCTk
mock_ctk.CTkFrame = MockFrame
mock_ctk.CTkButton = MockButton
mock_ctk.CTkLabel = MockLabel
mock_ctk.CTkTextbox = MockTextbox
mock_ctk.CTkEntry = MockEntry
mock_ctk.CTkFont = MagicMock
mock_ctk.CTkOptionMenu = MagicMock
mock_ctk.set_appearance_mode = MagicMock
mock_ctk.set_default_color_theme = MagicMock

sys.modules['customtkinter'] = mock_ctk

# Mock the frames to avoid importing them and their dependencies
sys.modules['gui.frames.home'] = MagicMock()
sys.modules['gui.frames.diagnose'] = MagicMock()
sys.modules['gui.frames.tools'] = MagicMock()
sys.modules['gui.frames.install'] = MagicMock()

# Now we can import DevAleGUI
from gui.app import DevAleGUI

class TestGUIStructure(unittest.TestCase):
    def setUp(self):
        self.runner = MagicMock()

    def test_gui_initialization(self):
        # This will test if the constructor runs without errors with mocked components
        app = DevAleGUI(self.runner)
        self.assertIsNotNone(app.frames)

        expected_items = ["Home", "Diagnose", "Tools", "Install"]

        # Check if the expected frames are there
        for name in expected_items:
            self.assertIn(name, app.frames)

        # Check if the expected buttons are stored for navigation access
        self.assertIsNotNone(app.nav_buttons)
        for name in expected_items:
            self.assertIn(name, app.nav_buttons)

if __name__ == "__main__":
    unittest.main()
