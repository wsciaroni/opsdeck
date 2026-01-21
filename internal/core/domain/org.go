package domain

import (
	"time"

	"github.com/google/uuid"
)

// Organization represents a workspace or tenant.
type Organization struct {
	ID               uuid.UUID `json:"id"`
	Name             string    `json:"name"`
	Slug             string    `json:"slug"`
	ShareLinkEnabled bool      `json:"share_link_enabled"`
	ShareLinkToken   *string   `json:"share_link_token"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// OrganizationMember represents the link between a user and an organization.
type OrganizationMember struct {
	OrganizationID uuid.UUID `json:"organization_id"`
	UserID         uuid.UUID `json:"user_id"`
	Role           string    `json:"role"`
	JoinedAt       time.Time `json:"joined_at"`
}

// UserMembership represents a user's membership in an organization, including the organization details and their role.
type UserMembership struct {
	Organization
	Role string `json:"role"`
}

// Member represents a user in an organization with their role.
type Member struct {
	UserID    uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	AvatarURL string    `json:"avatar_url"`
	Role      string    `json:"role"`
}
