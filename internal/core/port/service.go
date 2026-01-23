package port

import (
	"context"

	"github.com/wsciaroni/opsdeck/internal/core/domain"
)

// AuthService defines the interface for authentication operations.
type AuthService interface {
	GetLoginURL(state string) string
	LoginFromProvider(ctx context.Context, code string) (*domain.User, error)
	CreateSession(ctx context.Context, user *domain.User) (string, error)
}
