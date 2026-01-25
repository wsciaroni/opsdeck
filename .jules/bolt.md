## 2026-01-24 - Frontend List Performance
**Learning:** The Dashboard component frequently re-renders (e.g. on modal open), passing new function references to children like TicketBoard. Since TicketBoard renders many items, this caused O(N) re-renders of all ticket cards.
**Action:** Always memoize list item components (like `TicketCard`) in heavy dashboard views to isolate them from parent state changes.

## 2026-01-24 - Unstable Callback References in Lists
**Learning:** In `TicketList`, passing inline callbacks (e.g., `onClick={() => navigate(...)`) to children breaks `React.memo` optimizations because the function reference changes on every render.
**Action:** Move navigation logic or stable handlers inside the child component, or use `useCallback` for props passed to memoized list items.
