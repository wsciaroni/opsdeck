In the **Modular Monolith** (Hexagonal/Ports & Adapters) architecture, "Ports" are the interfaces that define the boundaries of your application core.

* **Primary Ports (Services):** How the outside world (HTTP Handlers) talks to your Core.
* **Secondary Ports (Repositories/Adapters):** How your Core talks to external tools (Database, Google OIDC).

By defining these interfaces *before* implementing the logic, we ensure that our business logic doesn't care if we swap Postgres for MySQL or Google for Azure AD later.

Here are the Go interfaces for the **Authentication & User Module**.

### 1. The Repository Port (Database Contract)

This interface defines exactly what data operations the core needs. It belongs in `internal/core/port/`.

**File:** `internal/core/port/repository.go`

```go
package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/wsciaroni/opsdeck/internal/core/domain"
)

// UserRepository defines how we persist User entities.
// Implementation will be in internal/adapter/storage/postgres/user_repo.go
type UserRepository interface {
	// GetByID retrieves a user by their system UUID.
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)

	// GetByEmail retrieves a user by their email address (used during login).
	GetByEmail(ctx context.Context, email string) (*domain.User, error)

	// Create saves a new user to the database.
	Create(ctx context.Context, user *domain.User) error

	// Update modifies an existing user (e.g., updating last login time or avatar).
	Update(ctx context.Context, user *domain.User) error
}

```

### 2. The OIDC Provider Port (External Auth Contract)

To keep your core testable, we don't want to hardcode the "Google" library inside the core logic. Instead, we define an interface for what we *need* from an identity provider.

**File:** `internal/core/port/oidc.go`

```go
package port

import (
	"context"
)

// UserInfo represents the standard data returned by an Identity Provider
type UserInfo struct {
	Email     string
	Name      string
	AvatarURL string
}

// OIDCProvider defines the interaction with an external Identity Provider (Google).
// Implementation will be in internal/adapter/auth/google/google_oidc.go
type OIDCProvider interface {
	// AuthCodeURL returns the URL to redirect the user to (Google Login Screen).
	AuthCodeURL(state string) string

	// ExchangeCode converts the temporary authorization code into User Info.
	ExchangeCode(ctx context.Context, code string) (*UserInfo, error)
}

```

### 3. The Service Port (Business Logic Contract)

This is the interface your HTTP handlers will call. It orchestrates the flow: "Take this code, talk to Google, then talk to the DB, and give me a User."

**File:** `internal/core/port/service.go`

```go
package port

import (
	"context"

	"github.com/wsciaroni/opsdeck/internal/core/domain"
)

// AuthService defines the business logic for authentication.
// Implementation will be in internal/core/service/auth_service.go
type AuthService interface {
	// GetLoginURL generates the redirect URL for the frontend button.
	GetLoginURL() string

	// LoginFromProvider handles the callback from Google.
	// It exchanges the code, finds or creates the user in the DB, and returns the domain User.
	LoginFromProvider(ctx context.Context, authCode string) (*domain.User, error)
	
	// CreateSession generates a secure session/cookie for the user.
	// (Returns the session token string)
	CreateSession(ctx context.Context, userID domain.User) (string, error)
}

```

### Why split it this way?

This structure gives you massive flexibility for the future:

1. **Mocking for Tests:** You can write a `MockOIDCProvider` that returns fake users ("Alice", "Bob") so you can test your login logic without actually clicking buttons in a browser or needing internet access.
2. **Swapping Providers:** If the church later decides to use Microsoft/Azure (common in non-profits), you simply write a new implementation of `OIDCProvider` without touching a single line of your core business logic or HTTP handlers.

## Possible AuthService Implementation

Here is the implementation of the **AuthService**. This is the "brain" of your authentication flow.

It resides in `internal/core/service/` and orchestrates the interaction between the outside world (Google) and your database.

### File: `internal/core/service/auth_service.go`

```go
package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/wsciaroni/opsdeck/internal/core/domain"
	"github.com/wsciaroni/opsdeck/internal/core/port"
)

// Ensure implementation matches interface
var _ port.AuthService = (*AuthService)(nil)

type AuthService struct {
	userRepo port.UserRepository
	oidc     port.OIDCProvider
	logger   *slog.Logger
}

// NewAuthService creates a new instance of the business logic
func NewAuthService(repo port.UserRepository, oidc port.OIDCProvider, logger *slog.Logger) *AuthService {
	return &AuthService{
		userRepo: repo,
		oidc:     oidc,
		logger:   logger,
	}
}

// GetLoginURL delegates to the OIDC provider to generate the Google link
func (s *AuthService) GetLoginURL() string {
	// In production, "state" should be a random string stored in a cookie to prevent CSRF.
	// For MVP, we use a static string or simple random generation.
	return s.oidc.AuthCodeURL("random-state-string")
}

// LoginFromProvider is the core logic: "Google said yes, now who is this?"
func (s *AuthService) LoginFromProvider(ctx context.Context, authCode string) (*domain.User, error) {
	// 1. Exchange the code with Google for user details
	userInfo, err := s.oidc.ExchangeCode(ctx, authCode)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange OIDC code: %w", err)
	}

	// 2. Check if user exists in our DB
	user, err := s.userRepo.GetByEmail(ctx, userInfo.Email)
	if err != nil {
		// Real error (DB down, etc.)
		return nil, fmt.Errorf("failed to lookup user: %w", err)
	}

	// 3. User Logic: Create vs Update
	if user == nil {
		// Case A: New User (Auto-Provisioning)
		s.logger.Info("provisioning new user", "email", userInfo.Email)
		
		newUser := &domain.User{
			ID:        uuid.New(),
			Email:     userInfo.Email,
			Name:      userInfo.Name,
			Role:      domain.RolePublic, // Default role as per FR-02
			AvatarURL: userInfo.AvatarURL,
		}

		if err := s.userRepo.Create(ctx, newUser); err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}
		return newUser, nil
	}

	// Case B: Existing User (Update details if changed)
	// We update the name/avatar in case they changed it on Google
	if user.Name != userInfo.Name || user.AvatarURL != userInfo.AvatarURL {
		user.Name = userInfo.Name
		user.AvatarURL = userInfo.AvatarURL
		if err := s.userRepo.Update(ctx, user); err != nil {
			s.logger.Warn("failed to update user details", "user_id", user.ID, "err", err)
			// Non-critical error, proceed with login
		}
	}

	return user, nil
}

// CreateSession would typically generate a JWT or Session Cookie.
// For the MVP, we will rely on the HTTP Handler layer (using alexedwards/scs) 
// to manage the actual cookie, so this method might just return the UserID 
// or be handled entirely in the adapter. 
func (s *AuthService) CreateSession(ctx context.Context, user domain.User) (string, error) {
	// Placeholder: In a Stateless JWT setup, this would sign a token.
	// Since we are using Session Cookies in the Handler, this logic 
	// is often minimal or moved to the handler.
	return user.ID.String(), nil
}

```

### Why this code is robust:

1. **Auto-Provisioning (FR-02):** It explicitly handles the `user == nil` case by creating a new user with `RolePublic`.
2. **Data Synchronization:** It updates the local user profile (`Name`, `Avatar`) if the external provider's data changes. This keeps your system fresh.
3. **Logging:** It uses structured logging (`slog`) to track when new users are provisioned, which is critical for the audit trail.

### Next Step in the Plan: Adapters

Now that we have the **Core Logic** (`AuthService`), it needs to talk to something real. It can't run yet because `port.UserRepository` and `port.OIDCProvider` are just empty interfaces.

We need to build the **Adapters** (Phase 1.1 / 1.2 continued):

1. **Postgres Adapter:** The actual SQL code to `SELECT` and `INSERT` users.
2. **Google Adapter:** The code that uses `coreos/go-oidc` to talk to Google.
