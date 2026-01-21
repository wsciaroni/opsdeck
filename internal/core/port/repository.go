package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/wsciaroni/opsdeck/internal/core/domain"
)

// UserRepository defines the interface for interacting with user data.
type UserRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Create(ctx context.Context, user *domain.User) error
	Update(ctx context.Context, user *domain.User) error
}

// OrganizationRepository defines the interface for interacting with organization data.
type OrganizationRepository interface {
	Create(ctx context.Context, org *domain.Organization) error
	AddMember(ctx context.Context, orgID uuid.UUID, userID uuid.UUID, role string) error
	ListByUser(ctx context.Context, userID uuid.UUID) ([]domain.UserMembership, error)
	ListMembers(ctx context.Context, orgID uuid.UUID) ([]domain.Member, error)
	RemoveMember(ctx context.Context, orgID uuid.UUID, userID uuid.UUID) error
}
