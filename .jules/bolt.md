## 2024-05-23 - React.memo Requires Stable Callbacks
**Learning:** Memoizing a child component (`TicketList`, `TicketBoard`) is ineffective if the parent component (`Dashboard`) passes a new inline function reference (e.g., `onOpenNewTicket={() => setIsModalOpen(true)}`) on every render.
**Action:** Always wrap callback props in `useCallback` in the parent component when passing them to memoized children to ensure prop stability and preventing unnecessary re-renders.

## 2024-05-24 - Batching DB Calls
**Learning:** Batching multiple `GetByID` calls into a single `GetByIDs` call reduces database roundtrips and is a simple, effective optimization, especially when the repository already supports batch fetching.
**Action:** Look for sequential `GetByID` calls in handlers and refactor them to use batch methods.
