import unittest
from unittest.mock import MagicMock, patch
import sys
import os

# Create a mock for customtkinter
mock_ctk = MagicMock()

class MockBase:
    def __init__(self, *args, **kwargs): pass
    def grid(self, *args, **kwargs): pass
    def grid_columnconfigure(self, *args, **kwargs): pass
    def grid_rowconfigure(self, *args, **kwargs): pass
    def grid_forget(self, *args, **kwargs): pass
    def grid_propagate(self, *args, **kwargs): pass
    def pack(self, *args, **kwargs): pass
    def pack_forget(self, *args, **kwargs): pass
    def configure(self, *args, **kwargs): pass
    def bind(self, *args, **kwargs): pass
    def add(self, *args, **kwargs): return MockBase()
    def get(self, *args, **kwargs): return ""
    def insert(self, *args, **kwargs): pass
    def see(self, *args, **kwargs): pass
    def delete(self, *args, **kwargs): pass
    def after(self, *args, **kwargs): pass

class MockCTk(MockBase):
    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.title = MagicMock()
        self.geometry = MagicMock()
    def mainloop(self): pass

mock_ctk.CTk = MockCTk
mock_ctk.CTkFrame = MockBase
mock_ctk.CTkLabel = MockBase
mock_ctk.CTkButton = MockBase
mock_ctk.CTkTextbox = MockBase
mock_ctk.CTkEntry = MockBase
mock_ctk.CTkOptionMenu = MockBase
mock_ctk.CTkTabview = MockBase
mock_ctk.CTkCheckBox = MockBase
mock_ctk.CTkFont = MagicMock
mock_ctk.BooleanVar = MagicMock
mock_ctk.set_appearance_mode = MagicMock
mock_ctk.set_default_color_theme = MagicMock

# Patch sys.modules to use our mock
sys.modules['customtkinter'] = mock_ctk

# Add the app directory to sys.path if not already there
app_dir = os.path.join(os.getcwd(), "DevAle official")
if app_dir not in sys.path:
    sys.path.append(app_dir)

# Now we can import things that use customtkinter
from gui.app import DevAleGUI, NAV_ITEMS

class TestDevAleGUIStructure(unittest.TestCase):
    @patch('threading.Thread')
    def setUp(self, mock_thread):
        self.runner = MagicMock()
        # Ensure we use our MockCTk
        with patch('gui.app.ctk.CTk', MockCTk):
            try:
                self.gui = DevAleGUI(self.runner)
            except Exception as e:
                self.fail(f"DevAleGUI failed to initialize: {e}")

    def test_frames_initialization(self):
        """Verify that all frames defined in NAV_ITEMS are initialized"""
        for item in NAV_ITEMS:
            name = item["name"]
            self.assertIn(name, self.gui.frames)
            self.assertIsNotNone(self.gui.frames[name])

    def test_navigation_buttons_initialization(self):
        """Verify that all navigation buttons are created and stored in self.nav_buttons"""
        self.assertTrue(hasattr(self.gui, 'nav_buttons'))
        for item in NAV_ITEMS:
            name = item["name"]
            self.assertIn(name, self.gui.nav_buttons)

if __name__ == "__main__":
    unittest.main()
