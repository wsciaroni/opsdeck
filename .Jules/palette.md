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
