## 2024-05-22 - Missing CSP and External Avatar Dependency
**Vulnerability:** Lack of Content Security Policy (CSP) allowed potential XSS escalation.
**Learning:** The application relies on `ui-avatars.com` for user avatars when no custom avatar is present. This external domain must be explicitly whitelisted in `img-src` directive of CSP.
**Prevention:** Always audit external resource usage (images, scripts, fonts) before enabling strict CSP. Used `grep -r "http"` to find external calls.

## 2024-05-24 - Removed External Avatar Dependency
**Vulnerability:** Usage of `ui-avatars.com` leaked user initials and names to a third-party service, and required weakening CSP.
**Learning:** Sending user data (even just names/initials) to external services for avatars is a privacy risk and expands the attack surface.
**Prevention:** Implemented local avatar generation using `Avatar.tsx` to generate colored initials deterministically. This allowed removing `ui-avatars.com` from CSP `img-src`.
