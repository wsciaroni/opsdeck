## 2024-05-23 - React.memo Requires Stable Callbacks
**Learning:** Memoizing a child component (`TicketList`, `TicketBoard`) is ineffective if the parent component (`Dashboard`) passes a new inline function reference (e.g., `onOpenNewTicket={() => setIsModalOpen(true)}`) on every render.
**Action:** Always wrap callback props in `useCallback` in the parent component when passing them to memoized children to ensure prop stability and preventing unnecessary re-renders.

## 2024-05-24 - Batching DB Calls
**Learning:** Batching multiple `GetByID` calls into a single `GetByIDs` call reduces database roundtrips and is a simple, effective optimization, especially when the repository already supports batch fetching.
**Action:** Look for sequential `GetByID` calls in handlers and refactor them to use batch methods.

## 2026-02-11 - Playwright API Mocks & Vite MIME Types
**Learning:** When mocking API routes (e.g., `/api/tickets`) in Playwright with Vite, use specific regex patterns (e.g., `re.compile(r".*/api/tickets(\?.*)?$")`) instead of globs like `**/api/tickets*`. Broad globs can accidentally intercept requests for source files (e.g., `src/api/tickets.ts`), causing Vite to fail with MIME type errors.
**Action:** Always use strict regex for API mocks in Playwright.
