## 2024-05-22 - Missing CSP and External Avatar Dependency
**Vulnerability:** Lack of Content Security Policy (CSP) allowed potential XSS escalation.
**Learning:** The application relies on `ui-avatars.com` for user avatars when no custom avatar is present. This external domain must be explicitly whitelisted in `img-src` directive of CSP.
**Prevention:** Always audit external resource usage (images, scripts, fonts) before enabling strict CSP. Used `grep -r "http"` to find external calls.
