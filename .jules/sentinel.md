## 2026-01-21 - Session Hijacking via Raw UUID Cookies
**Vulnerability:** The application used raw User UUIDs as session cookies (session_id). This allowed attackers to impersonate any user by guessing or obtaining their UUID and setting it as the cookie value.
**Learning:** MVP/Prototype code often uses simplified authentication ("For now, return user ID") which becomes a critical vulnerability if not replaced before production or wider testing. The AuthMiddleware blindly trusted the cookie value.
**Prevention:** Always cryptographically sign session tokens (JWT, signed cookies) or use a server-side session store with random high-entropy tokens. Never trust client-side identifiers for authentication without verification.
