## 2024-05-22 - [Icon-Only Button Accessibility Pattern]
**Learning:** Icon-only buttons (like Edit/Delete in lists) often lack `title` attributes, making them inaccessible to sighted users who rely on tooltips, and lack focus states for keyboard users.
**Action:** Always add `title="Action Name"` and `className="p-1 rounded-full focus:outline-none focus:ring-2 focus:ring-offset-2 ..."` to icon-only buttons.

## 2024-10-27 - [Kanban Board Semantics]
**Learning:** Kanban boards are often implemented as nested divs, making them opaque to screen readers. Using `<section>`, `<h2>` for columns, and `<ul>`/`<li>` for cards dramatically improves navigability.
**Action:** Check board-like components for semantic structure and use aria-labels for counts.
