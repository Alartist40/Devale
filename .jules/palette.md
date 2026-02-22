## 2025-01-24 - Sidebar Active State Highlighting
**Learning:** In multi-frame CustomTkinter applications, providing visual feedback for the active navigation section is crucial for user orientation. Using `fg_color=None` and `text_color=None` allows widgets to revert to theme defaults, which provides high contrast for the active state without hardcoding colors.
**Action:** Store navigation buttons in a dictionary during initialization to allow easy programmatic state updates in the frame-switching logic.
