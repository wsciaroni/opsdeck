## 2024-05-22 - Missing CSP and External Avatar Dependency
**Vulnerability:** Lack of Content Security Policy (CSP) allowed potential XSS escalation.
**Learning:** The application relies on `ui-avatars.com` for user avatars when no custom avatar is present. This external domain must be explicitly whitelisted in `img-src` directive of CSP.
**Prevention:** Always audit external resource usage (images, scripts, fonts) before enabling strict CSP. Used `grep -r "http"` to find external calls.

## 2024-05-24 - Rate Limiter Bypass via Map Reset
**Vulnerability:** The `RateLimiter` middleware completely reset the client tracking map when it reached `maxClients`, allowing attackers to bypass rate limits by flooding the system with unique IPs to trigger a reset.
**Learning:** Simple `if len >= max { map = make() }` logic creates a fail-open condition under load, neutralizing protection against distributed attacks.
**Prevention:** Use random eviction or LRU for bounded maps instead of full reset. In Go, `for k := range m { delete(m, k); break }` provides O(1) random eviction.

## 2025-02-23 - Privilege Escalation via Generic Membership Check
**Vulnerability:** ScheduledTaskHandler only verified that a user was a 'member' of the organization to perform administrative actions (Create/Update/Delete), allowing any member to modify critical schedules.
**Learning:** Generic `verifyMembership` checks are insufficient for sensitive operations. Always verify specific roles (e.g., 'owner', 'admin') or permissions.
**Prevention:** Use `GetMemberRole` for explicit role verification instead of iterating `ListByUser`.

## 2025-02-23 - DoS via Excessive Body Limit on JSON Endpoints
**Vulnerability:** `TicketHandler` used a single `MaxRequestSize` (32MB) designed for file uploads across all endpoints, allowing attackers to exhaust memory by sending large JSON payloads to `UpdateTicket` (which only requires small text updates).
**Learning:** Validating input length (e.g., `len(title) > 200`) *after* decoding the full request body is insufficient for DoS protection because the memory is already allocated.
**Prevention:** Enforce strict body size limits (e.g., 1MB via `MaxJSONBodySize`) for JSON endpoints using `http.MaxBytesReader` *before* decoding, and reserve larger limits only for `multipart/form-data` requests.
