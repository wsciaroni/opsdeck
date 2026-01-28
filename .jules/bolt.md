## 2024-05-23 - React.memo Requires Stable Callbacks
**Learning:** Memoizing a child component (`TicketList`, `TicketBoard`) is ineffective if the parent component (`Dashboard`) passes a new inline function reference (e.g., `onOpenNewTicket={() => setIsModalOpen(true)}`) on every render.
**Action:** Always wrap callback props in `useCallback` in the parent component when passing them to memoized children to ensure prop stability and preventing unnecessary re-renders.
