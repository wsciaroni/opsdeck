## 2024-05-22 - [Icon-Only Button Accessibility Pattern]
**Learning:** Icon-only buttons (like Edit/Delete in lists) often lack `title` attributes, making them inaccessible to sighted users who rely on tooltips, and lack focus states for keyboard users.
**Action:** Always add `title="Action Name"` and `className="p-1 rounded-full focus:outline-none focus:ring-2 focus:ring-offset-2 ..."` to icon-only buttons.

## 2024-05-23 - [Semantic Board Accessibility]
**Learning:** Kanban boards often use nested `div`s which provide no structure for screen readers. Using `h3` for column headers and `ul`/`li` for card lists provides critical navigation context ("List of 5 items").
**Action:** Refactor board components to use semantic HTML hierarchy (`h2/h3` > `ul` > `li`) and use `Link` components instead of clickable `div`s for cards.
