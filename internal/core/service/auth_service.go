package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/wsciaroni/opsdeck/internal/core/domain"
	"github.com/wsciaroni/opsdeck/internal/core/port"
)

// AuthService implements port.AuthService.
type AuthService struct {
	repo   port.UserRepository
	oidc   port.OIDCProvider
	logger *slog.Logger
}

// NewAuthService creates a new AuthService.
func NewAuthService(repo port.UserRepository, oidc port.OIDCProvider, logger *slog.Logger) *AuthService {
	return &AuthService{
		repo:   repo,
		oidc:   oidc,
		logger: logger,
	}
}

// GetLoginURL delegates to the OIDC provider to get the authentication URL.
func (s *AuthService) GetLoginURL() string {
	// For this MVP phase, we use a static state string.
	return s.oidc.AuthCodeURL("state-random-string")
}

// LoginFromProvider authenticates a user using an authorization code.
func (s *AuthService) LoginFromProvider(ctx context.Context, code string) (*domain.User, error) {
	// Step 1: Exchange code for user info
	userInfo, err := s.oidc.ExchangeCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	// Step 2: Check if user exists
	user, err := s.repo.GetByEmail(ctx, userInfo.Email)
	if err != nil {
		// Assuming repository returns nil, nil for not found as per guidelines.
		// If it returns a real error (like DB connection fail), we should fail.
		// However, I need to know if GetByEmail returns error on "not found" or just nil, nil.
		// The memory says: "Repository retrieval methods must return nil, nil instead of an error when a record is not found".
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}

	// Step 3: Auto-Provisioning
	if user == nil {
		// New User
		newUser := &domain.User{
			ID:        uuid.New(),
			Email:     userInfo.Email,
			Name:      userInfo.Name,
			Role:      domain.RolePublic,
			AvatarURL: userInfo.AvatarURL,
		}

		if err := s.repo.Create(ctx, newUser); err != nil {
			return nil, fmt.Errorf("failed to create new user: %w", err)
		}
		s.logger.Info("provisioning new user", "user_id", newUser.ID, "email", newUser.Email)
		return newUser, nil
	}

	// Existing User
	updated := false
	if user.Name != userInfo.Name {
		user.Name = userInfo.Name
		updated = true
	}
	if user.AvatarURL != userInfo.AvatarURL {
		user.AvatarURL = userInfo.AvatarURL
		updated = true
	}

	if updated {
		if err := s.repo.Update(ctx, user); err != nil {
			return nil, fmt.Errorf("failed to update user details: %w", err)
		}
	}

	return user, nil
}

// CreateSession creates a new session for the user and returns the session token.
func (s *AuthService) CreateSession(ctx context.Context, user *domain.User) (string, error) {
	// For now, we simply return the user ID as the session token.
	return user.ID.String(), nil
}
