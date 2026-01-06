import sys
import os

def get_resource_path(relative_path):
    """ Get absolute path to resource, works for dev and for PyInstaller """
    try:
        # PyInstaller creates a temp folder and stores path in _MEIPASS
        base_path = sys._MEIPASS
    except Exception:
        base_path = os.path.abspath(".")
    
    # If we are in dev mode, sometimes we are in root, sometimes in my/
    # We want to find where 'commands' is.
    
    # Check if base_path has the file, if not try to adjust
    path = os.path.join(base_path, relative_path)
    
    if not os.path.exists(path):
        # Try finding it in 'my/' if we are in dev root
        alt_path = os.path.join(base_path, 'my', relative_path)
        if os.path.exists(alt_path):
            return alt_path
            
    return path
