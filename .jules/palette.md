## 2025-05-14 - [Sidebar Navigation Highlighting]
**Learning:** In multi-frame CustomTkinter applications, providing visual feedback for the active navigation state is crucial for user orientation. Storing button references in a dictionary keyed by frame name allows for efficient programmatic state management. Using class constants for design tokens (like inactive colors) prevents style duplication and makes the UI easier to maintain.
**Action:** Always implement active state tracking in sidebar navigation by storing button references and using a centralized styling method or constants.

## 2025-05-15 - [Console Management and Cognitive Load]
**Learning:** In utility applications that stream a lot of system output, the console can quickly become cluttered, leading to high cognitive load and difficulty in spotting new information. Providing a "Clear" action allows users to reset their workspace and focus on current operations.
**Action:** Always include a way to clear persistent logs or consoles, especially when they are used for real-time feedback. Ensure the clear action is easily accessible but does not distract from the main output.
