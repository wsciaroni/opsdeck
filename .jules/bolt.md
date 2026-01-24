## 2026-01-24 - Frontend List Performance
**Learning:** The Dashboard component frequently re-renders (e.g. on modal open), passing new function references to children like TicketBoard. Since TicketBoard renders many items, this caused O(N) re-renders of all ticket cards.
**Action:** Always memoize list item components (like `TicketCard`) in heavy dashboard views to isolate them from parent state changes.
