# Palette's Journal

## 2025-05-18 - Icon-Only Button Accessibility
**Learning:** Found a critical accessibility gap where navigation buttons (like "Back") were implemented as icon-only buttons without `aria-label` or `title`. This makes them invisible to screen reader users and confusing for mouse users who rely on tooltips.
**Action:** Always verify icon-only buttons have accessible names (`aria-label`) and visual tooltips (`title`). Ensure `aria-hidden="true"` is applied to the decorative icon itself.
