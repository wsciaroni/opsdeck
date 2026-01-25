## 2026-01-21 - Session Hijacking via Raw UUID Cookies
**Vulnerability:** The application used raw User UUIDs as session cookies (session_id). This allowed attackers to impersonate any user by guessing or obtaining their UUID and setting it as the cookie value.
**Learning:** MVP/Prototype code often uses simplified authentication ("For now, return user ID") which becomes a critical vulnerability if not replaced before production or wider testing. The AuthMiddleware blindly trusted the cookie value.
**Prevention:** Always cryptographically sign session tokens (JWT, signed cookies) or use a server-side session store with random high-entropy tokens. Never trust client-side identifiers for authentication without verification.

## 2024-05-23 - [SQL Injection and CSRF Analysis]
**Vulnerability:**
1.  **CSRF (High/Critical):** The application relies on `session_id` cookie for authentication (`internal/adapter/web/middleware/auth.go`) but does not implement any CSRF protection (middleware or token verification).
2.  **SQL Injection (Safe):** The SQL construction in `internal/adapter/storage/postgres/ticket_repo.go` uses `fmt.Sprintf` only for placeholders (e.g., `$1`, `$2`), which is safe as long as the arguments are passed separately to `Query/QueryRow`. The repository code looks clean in this regard.

**Learning:** Go's `database/sql` and `pgx` encourage parameterized queries, but manual query building with `fmt.Sprintf` can be risky if not done carefully. In this case, it is done correctly for dynamic filtering.

**Prevention:** To prevent CSRF, we should implement the "Double Submit Cookie" pattern or use a Synchronizer Token Pattern. Since the frontend is React, we can have the backend set a `X-CSRF-Token` cookie (httpOnly=false) and require the frontend to read it and send it back in a header `X-CSRF-Token`. Or, more simply for this exercise, we can add basic Security Headers as a quick win if CSRF is too complex for "one small fix".

## 2026-01-22 - OAuth CSRF Protection via State Parameter
**Vulnerability:** The OAuth flow used a hardcoded state parameter ("state-random-string"), making it vulnerable to CSRF/Account Takeover attacks where an attacker could force a user to log in to the attacker's account.
**Learning:** OAuth `state` parameter must be a unique, unguessable, cryptographically secure random string bound to the user's session (e.g., via a cookie) and verified in the callback.
**Prevention:** Implemented a `generateState` helper using `crypto/rand` and an `oauth_state` cookie (HTTPOnly, Secure, Lax) to store and verify the state during the OAuth dance.

## 2026-01-22 - CSV Injection (Formula Injection) in Exports
**Vulnerability:** The ticket export functionality (`ExportTickets`) directly included user-controlled input (`Title`, `Description`) in CSV files. If these fields started with `=`, `+`, `-`, or `@`, they would be executed as formulas by spreadsheet software (Excel, Sheets), potentially leading to command execution or data exfiltration.
**Learning:** CSV is not just text; it's a file format that spreadsheet software interprets. Any user input going into a CSV must be sanitized.
**Prevention:** Prepend a single quote `'` to any field starting with the dangerous characters (`=`, `+`, `-`, `@`) to force the spreadsheet software to treat it as a literal string.
