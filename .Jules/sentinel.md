## 2024-05-22 - Missing CSP and External Avatar Dependency
**Vulnerability:** Lack of Content Security Policy (CSP) allowed potential XSS escalation.
**Learning:** The application relies on `ui-avatars.com` for user avatars when no custom avatar is present. This external domain must be explicitly whitelisted in `img-src` directive of CSP.
**Prevention:** Always audit external resource usage (images, scripts, fonts) before enabling strict CSP. Used `grep -r "http"` to find external calls.

## 2024-05-24 - Rate Limiter Bypass via Map Reset
**Vulnerability:** The `RateLimiter` middleware completely reset the client tracking map when it reached `maxClients`, allowing attackers to bypass rate limits by flooding the system with unique IPs to trigger a reset.
**Learning:** Simple `if len >= max { map = make() }` logic creates a fail-open condition under load, neutralizing protection against distributed attacks.
**Prevention:** Use random eviction or LRU for bounded maps instead of full reset. In Go, `for k := range m { delete(m, k); break }` provides O(1) random eviction.
