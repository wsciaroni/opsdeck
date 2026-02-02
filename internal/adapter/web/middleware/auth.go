package middleware

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/wsciaroni/opsdeck/internal/core/domain"
	"github.com/wsciaroni/opsdeck/internal/core/port"
)

type contextKey string

const UserContextKey contextKey = "user"

type AuthMiddleware struct {
	userRepo port.UserRepository
	logger   *slog.Logger
	secret   []byte
}

func NewAuthMiddleware(userRepo port.UserRepository, logger *slog.Logger, secret string) *AuthMiddleware {
	return &AuthMiddleware{
		userRepo: userRepo,
		logger:   logger,
		secret:   []byte(secret),
	}
}

// SignSessionID signs the session ID with the secret and expiration.
func SignSessionID(id string, expiresAt time.Time, secret []byte) string {
	expStr := strconv.FormatInt(expiresAt.Unix(), 10)
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(id + "." + expStr))
	signature := hex.EncodeToString(mac.Sum(nil))
	return id + "." + expStr + "." + signature
}

func (m *AuthMiddleware) Protect(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			// No session cookie
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(cookie.Value, ".")
		if len(parts) != 3 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		id := parts[0]
		expStr := parts[1]
		signature := parts[2]

		mac := hmac.New(sha256.New, m.secret)
		mac.Write([]byte(id + "." + expStr))
		expectedSignature := hex.EncodeToString(mac.Sum(nil))

		if !hmac.Equal([]byte(signature), []byte(expectedSignature)) {
			m.logger.Warn("invalid session signature", "cookie", cookie.Value)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		exp, err := strconv.ParseInt(expStr, 10, 64)
		if err != nil {
			m.logger.Warn("invalid session expiration format", "cookie", cookie.Value)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if time.Now().Unix() > exp {
			m.logger.Warn("session expired", "cookie", cookie.Value)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userID, err := uuid.Parse(id)
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
