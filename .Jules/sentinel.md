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

## 2025-05-24 - DoS via Large File Uploads Before Authorization
**Vulnerability:** `CreatePublicTicket` and `CreateTicket` parsed multipart request bodies (up to 32MB) before validating authorization, allowing unauthorized users to exhaust server resources.
**Learning:** `ParseMultipartForm` reads the entire body (if within limit) before returning control. Relying on body parameters for authorization in file upload handlers is inherently vulnerable to DoS.
**Prevention:** Require authorization tokens or IDs in URL query parameters for multipart endpoints to enable early rejection (fail-fast) before processing large payloads. Limit JSON bodies independently (e.g., 1MB vs 32MB for files).
