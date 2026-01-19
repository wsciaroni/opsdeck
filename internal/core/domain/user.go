package domain

import (
	"time"

	"github.com/google/uuid"
)

// Role represents the role of a user in the system.
type Role string

const (
	RolePublic  Role = "public"
	RoleStaff   Role = "staff"
	RoleManager Role = "manager"
	RoleAdmin   Role = "admin"
)

// User represents a user in the system.
type User struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Role      Role      `json:"role"`
	AvatarURL string    `json:"avatar_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
