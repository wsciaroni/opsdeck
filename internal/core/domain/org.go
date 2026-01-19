package domain

import (
	"time"

	"github.com/google/uuid"
)

// Organization represents a workspace or tenant.
type Organization struct {
	ID        uuid.UUID
	Name      string
	Slug      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// OrganizationMember represents the link between a user and an organization.
type OrganizationMember struct {
	OrganizationID uuid.UUID
	UserID         uuid.UUID
	Role           string
	JoinedAt       time.Time
}

// UserMembership represents a user's membership in an organization, including the organization details and their role.
type UserMembership struct {
	Organization
	Role string
}
