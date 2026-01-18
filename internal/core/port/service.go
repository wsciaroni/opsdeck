package port

import (
	"context"

	"github.com/wsciaroni/opsdeck/internal/core/domain"
)

// AuthService defines the interface for authentication operations.
type AuthService interface {
	GetLoginURL() string
	LoginFromProvider(ctx context.Context, code string) (*domain.User, error)
}
