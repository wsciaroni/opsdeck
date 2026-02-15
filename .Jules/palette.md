## 2024-05-22 - [Icon-Only Button Accessibility Pattern]
**Learning:** Icon-only buttons (like Edit/Delete in lists) often lack `title` attributes, making them inaccessible to sighted users who rely on tooltips, and lack focus states for keyboard users.
**Action:** Always add `title="Action Name"` and `className="p-1 rounded-full focus:outline-none focus:ring-2 focus:ring-offset-2 ..."` to icon-only buttons.

## 2025-02-15 - [Semantic Navigation with Link]
**Learning:** Using `div` with `role="button"` for navigation elements requires manual keyboard handling (`Enter`/`Space`) and `e.preventDefault()`, and loses native browser features like 'Open in new tab'.
**Action:** Always prefer native `<Link>` (or `<a>`) components for navigation-based cards or list items, applying `block` display style if necessary to maintain layout.
