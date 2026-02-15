<<<<<<< HEAD
# Palette's Journal

## 2025-05-18 - Icon-Only Button Accessibility
**Learning:** Found a critical accessibility gap where navigation buttons (like "Back") were implemented as icon-only buttons without `aria-label` or `title`. This makes them invisible to screen reader users and confusing for mouse users who rely on tooltips.
**Action:** Always verify icon-only buttons have accessible names (`aria-label`) and visual tooltips (`title`). Ensure `aria-hidden="true"` is applied to the decorative icon itself.

## 2025-05-18 - Input Feedback
**Learning:** File inputs that only show a count (e.g., "2 files selected") obscure critical information. Users need to verify *which* files they attached before submitting.
**Action:** Always display a list of selected filenames for file inputs. Combine this with loading spinners on submit buttons to provide complete feedback on the action lifecycle.

## 2026-01-24 - Accessible Clickable Table Rows
**Learning:** Clickable table rows (`<tr onClick={...} />`) are inaccessible to keyboard users by default, breaking the navigation flow. Adding `tabIndex="0"`, `role="button"` (or just handling interactions), and `onKeyDown` handlers for Enter/Space is essential for keyboard accessibility.
**Action:** When implementing clickable rows, ensure they are focusable (`tabIndex="0"`), have visible focus states (`focus:ring`), respond to keyboard events (`Enter`/`Space`), and provide an accessible name via `aria-label`.

## 2026-01-25 - Async Button Feedback
**Learning:** Adding immediate visual feedback (spinner + text change) to submit buttons significantly improves perceived performance and prevents double-submissions.
**Action:** Use the `Loader2` icon and `mutation.isPending` state for all async form submission buttons.

## 2026-01-26 - Skip to Content Link
**Learning:** Single Page Applications (SPAs) often neglect the "Skip to Content" link because navigation is handled by JS, but keyboard users still need a way to bypass repetitive header navigation on every page load.
**Action:** Always include a visually hidden, focusable anchor tag (`href="#main-content"`) at the top of the `Layout` component and ensure the `<main>` element has `id="main-content"` and `tabIndex={-1}` for focus management.

## 2026-02-05 - Disabled State Clarity
**Learning:** Simply disabling a button isn't enough; visual cues like reduced opacity (`opacity-50`) and cursor change (`cursor-not-allowed`) are crucial for communicating that an action is temporarily unavailable (e.g., during submission).
**Action:** Always pair `disabled={isPending}` with `disabled:opacity-50 disabled:cursor-not-allowed` utility classes to provide clear visual feedback.

## 2026-05-18 - Keyboard Shortcuts
**Learning:** Adding keyboard shortcuts (like 'c' for Create) significantly speeds up power user workflows, but they must be implemented carefully to avoid triggering during normal typing (e.g., in inputs/textareas).
**Action:** Always wrap global keydown listeners in checks for `tagName` (INPUT, TEXTAREA) and `isContentEditable`. Add tooltips or hints to UI elements to help users discover these shortcuts.

## 2026-05-18 - Search Input Patterns
**Learning:** Search inputs without a "Clear" button force users to manually delete text, which is tedious on mobile. Also, inputs often lack labels when designed with just icons.
**Action:** Always include a hidden `<label>` for accessibility and a "Clear" (X) button that appears when text is present. Position the button absolutely within the input wrapper.

## 2026-05-19 - File Size Formatting
**Learning:** Displaying raw bytes (e.g., "1536 Bytes") or hardcoded KB values (e.g., "1536.00 KB") is user-hostile, especially on mobile where data context matters.
**Action:** Use a reusable `formatBytes` utility across the application to consistently display human-readable file sizes (e.g., "1.5 KB", "5 MB").

## 2026-05-20 - Managed File Uploads
**Learning:** Standard `<input type="file" multiple>` is user-hostile because selecting new files replaces the entire list, preventing users from adding/removing individual files incrementally.
**Action:** Always maintain a `files` array state separate from the input. Implement a "file list" UI with removal buttons and use the input only for adding files (clearing its value after selection to allow re-selection).

## 2026-05-24 - Text Contrast for Secondary Information
**Learning:** `text-gray-400` (Tailwind) is often used for timestamps or metadata but fails WCAG AA contrast ratios on white backgrounds, especially for small text (`text-xs`).
**Action:** Use `text-gray-500` for standard text and `text-gray-600` for small text (`text-xs`) to ensure readability while maintaining visual hierarchy.

## 2026-05-25 - Icon Contrast and Focus
**Learning:** Decorative and functional icons using `text-gray-400` often fail contrast requirements on white backgrounds. Also, small remove buttons in lists need clear focus rings to be usable by keyboard users.
**Action:** Use `text-gray-500` as the minimum darkness for icons on white backgrounds. Ensure all interactive list items (like remove buttons) have `focus-visible:ring-2` styles.

## 2026-05-26 - Keyboard Shortcut Hints
**Learning:** Adding visual keyboard shortcut hints (badges) significantly improves discoverability for power users, but they must be implemented with `aria-keyshortcuts` and `aria-hidden="true"` to ensure a good screen reader experience.
**Action:** Always pair visual shortcut hints (e.g., "C" badge) with `aria-keyshortcuts` attributes on the interactive element.

## 2026-05-27 - Accessible Drag & Drop File Upload
**Learning:** Native file inputs are inaccessible to drag-and-drop workflows and screen readers often struggle with nested label structures.
**Action:** Enhance file uploads by wrapping the input in a semantic, keyboard-accessible drag-and-drop zone with visual feedback (`isDragging`) and clear instruction text. Ensure the input has a unique `id` and is properly associated via `aria-labelledby`.
=======
## 2024-05-22 - [Icon-Only Button Accessibility Pattern]
**Learning:** Icon-only buttons (like Edit/Delete in lists) often lack `title` attributes, making them inaccessible to sighted users who rely on tooltips, and lack focus states for keyboard users.
**Action:** Always add `title="Action Name"` and `className="p-1 rounded-full focus:outline-none focus:ring-2 focus:ring-offset-2 ..."` to icon-only buttons.
>>>>>>> origin/master
