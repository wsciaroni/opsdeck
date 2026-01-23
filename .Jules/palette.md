# Palette's Journal

## 2025-05-18 - Icon-Only Button Accessibility
**Learning:** Found a critical accessibility gap where navigation buttons (like "Back") were implemented as icon-only buttons without `aria-label` or `title`. This makes them invisible to screen reader users and confusing for mouse users who rely on tooltips.
**Action:** Always verify icon-only buttons have accessible names (`aria-label`) and visual tooltips (`title`). Ensure `aria-hidden="true"` is applied to the decorative icon itself.

## 2025-05-18 - Input Feedback
**Learning:** File inputs that only show a count (e.g., "2 files selected") obscure critical information. Users need to verify *which* files they attached before submitting.
**Action:** Always display a list of selected filenames for file inputs. Combine this with loading spinners on submit buttons to provide complete feedback on the action lifecycle.
