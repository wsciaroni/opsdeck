package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/wsciaroni/opsdeck/internal/core/domain"
	"github.com/wsciaroni/opsdeck/internal/core/port"
)

// AuthService implements port.AuthService.
type AuthService struct {
	repo    port.UserRepository
	orgRepo port.OrganizationRepository
	oidc    port.OIDCProvider
	logger  *slog.Logger
}

// NewAuthService creates a new AuthService.
func NewAuthService(repo port.UserRepository, orgRepo port.OrganizationRepository, oidc port.OIDCProvider, logger *slog.Logger) *AuthService {
	return &AuthService{
		repo:    repo,
		orgRepo: orgRepo,
		oidc:    oidc,
		logger:  logger,
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

		// Create Default Organization
		orgName := "Personal Workspace"
		orgSlug, err := generateSlug(orgName)
		if err != nil {
			return nil, fmt.Errorf("failed to generate org slug: %w", err)
		}

		newOrg := &domain.Organization{
			Name: orgName,
			Slug: orgSlug,
		}

		if err := s.orgRepo.Create(ctx, newOrg); err != nil {
			return nil, fmt.Errorf("failed to create default organization: %w", err)
		}

		// Add User as Owner
		if err := s.orgRepo.AddMember(ctx, newOrg.ID, newUser.ID, "owner"); err != nil {
			return nil, fmt.Errorf("failed to add user to organization: %w", err)
		}
		s.logger.Info("created default organization", "org_id", newOrg.ID, "user_id", newUser.ID)

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

func generateSlug(name string) (string, error) {
	// slugify(name) + "-" + randomHex(4)
	lowerName := strings.ToLower(name)
	reg, err := regexp.Compile("[^a-z0-9]+")
	if err != nil {
		return "", err
	}
	slugBase := reg.ReplaceAllString(lowerName, "-")
	slugBase = strings.Trim(slugBase, "-")

	bytes := make([]byte, 2) // 2 bytes = 4 hex chars
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	randomSuffix := hex.EncodeToString(bytes)

	return fmt.Sprintf("%s-%s", slugBase, randomSuffix), nil
}
