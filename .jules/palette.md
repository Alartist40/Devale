<<<<<<< palette/navigation-highlighting-17953901055037944356
## 2025-05-14 - Navigation State Feedback in Desktop Apps
**Learning:** In desktop utility applications with sidebar navigation, users can easily lose track of the current active section if visual feedback (like highlighting) is missing. Highlighting should not only change the background but also reset text color to maintain contrast.
**Action:** Always store navigation elements in a registry (e.g., a dictionary) to allow stateful updates when switching views. Use theme-default colors for active states and transparent/dimmed styles for inactive ones to create clear visual hierarchy.
=======
## 2025-05-14 - [Sidebar Navigation Highlighting]
**Learning:** In multi-frame CustomTkinter applications, providing visual feedback for the active navigation state is crucial for user orientation. Storing button references in a dictionary keyed by frame name allows for efficient programmatic state management. Using class constants for design tokens (like inactive colors) prevents style duplication and makes the UI easier to maintain.
**Action:** Always implement active state tracking in sidebar navigation by storing button references and using a centralized styling method or constants.
>>>>>>> main
