## 2025-05-14 - Navigation State Feedback in Desktop Apps
**Learning:** In desktop utility applications with sidebar navigation, users can easily lose track of the current active section if visual feedback (like highlighting) is missing. Highlighting should not only change the background but also reset text color to maintain contrast.
**Action:** Always store navigation elements in a registry (e.g., a dictionary) to allow stateful updates when switching views. Use theme-default colors for active states and transparent/dimmed styles for inactive ones to create clear visual hierarchy.
