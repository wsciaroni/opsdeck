## 2024-05-22 - [Icon-Only Button Accessibility Pattern]
**Learning:** Icon-only buttons (like Edit/Delete in lists) often lack `title` attributes, making them inaccessible to sighted users who rely on tooltips, and lack focus states for keyboard users.
**Action:** Always add `title="Action Name"` and `className="p-1 rounded-full focus:outline-none focus:ring-2 focus:ring-offset-2 ..."` to icon-only buttons.

## 2024-05-24 - [Semantic Links for Card Navigation]
**Learning:** Complex cards acting as navigation targets are often implemented as clickable divs (`role="button"`), breaking standard browser behaviors like "Open in new tab".
**Action:** Wrap the entire card content in a `<Link>` or `<a>` tag with `block` display, instead of using `onClick` handlers on a `div`.
