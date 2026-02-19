## 2024-05-23 - React.memo Requires Stable Callbacks
**Learning:** Memoizing a child component (`TicketList`, `TicketBoard`) is ineffective if the parent component (`Dashboard`) passes a new inline function reference (e.g., `onOpenNewTicket={() => setIsModalOpen(true)}`) on every render.
**Action:** Always wrap callback props in `useCallback` in the parent component when passing them to memoized children to ensure prop stability and preventing unnecessary re-renders.

## 2024-05-24 - Batching DB Calls
**Learning:** Batching multiple `GetByID` calls into a single `GetByIDs` call reduces database roundtrips and is a simple, effective optimization, especially when the repository already supports batch fetching.
**Action:** Look for sequential `GetByID` calls in handlers and refactor them to use batch methods.

## 2024-05-25 - Expensive Date Formatting in Render Loop
**Learning:** Repeatedly calling `new Date().toLocaleDateString()` inside a large list render loop is expensive because it parses the date string and creates a new formatter instance every time.
**Action:** Use a cached `Intl.DateTimeFormat` instance outside the component or in a utility function to format dates efficiently.

## 2024-05-26 - Memoizing Filter Checkboxes
**Learning:** Rendering checkboxes in a loop with inline arrow functions causes unnecessary re-renders of all items when one changes. Extracting to a memoized component and using stable callbacks prevents this.
**Action:** Always extract list items with interactivity to memoized components and ensure callbacks are stable.
