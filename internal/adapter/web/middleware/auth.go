package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/wsciaroni/opsdeck/internal/core/domain"
	"github.com/wsciaroni/opsdeck/internal/core/port"
)

type contextKey string

const UserContextKey contextKey = "user"

type AuthMiddleware struct {
	userRepo port.UserRepository
	logger   *slog.Logger
}

func NewAuthMiddleware(userRepo port.UserRepository, logger *slog.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		userRepo: userRepo,
		logger:   logger,
	}
}

func (m *AuthMiddleware) Protect(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			// No session cookie
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userID, err := uuid.Parse(cookie.Value)
		if err != nil {
			// Invalid UUID
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		user, err := m.userRepo.GetByID(r.Context(), userID)
		if err != nil {
			m.logger.Error("failed to get user", "error", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if user == nil {
			// User not found
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUser retrieves the user from the context.
func GetUser(ctx context.Context) *domain.User {
	user, ok := ctx.Value(UserContextKey).(*domain.User)
	if !ok {
		return nil
	}
	return user
}
