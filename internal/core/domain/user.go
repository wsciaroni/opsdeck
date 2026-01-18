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
	ID        uuid.UUID
	Email     string
	Name      string
	Role      Role
	AvatarURL string
	CreatedAt time.Time
	UpdatedAt time.Time
}
