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
