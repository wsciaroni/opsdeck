package domain

import (
	"time"

	"github.com/google/uuid"
)

// Ticket Priorities
const (
	TicketPriorityLow      = "low"
	TicketPriorityMedium   = "medium"
	TicketPriorityHigh     = "high"
	TicketPriorityCritical = "critical"
)

// Ticket Statuses
const (
	TicketStatusNew        = "new"
	TicketStatusInProgress = "in_progress"
	TicketStatusOnHold     = "on_hold"
	TicketStatusDone       = "done"
	TicketStatusCanceled   = "canceled"
)

// Ticket represents a support ticket in the system.
type Ticket struct {
	ID             uuid.UUID  `json:"id"`
	OrganizationID uuid.UUID  `json:"organization_id"`
	Title          string     `json:"title"`
	Description    string     `json:"description"`
	Location       string     `json:"location"`
	StatusID       string     `json:"status_id"`
	PriorityID     string     `json:"priority_id"`
	ReporterID     uuid.UUID  `json:"reporter_id"`
	AssigneeUserID *uuid.UUID `json:"assignee_user_id"`
	Sensitive      bool       `json:"sensitive"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	CompletedAt    *time.Time `json:"completed_at"`
}

// TicketStatus represents a valid status for a ticket.
type TicketStatus struct {
	ID    string
	Label string
	Level int
}

// TicketPriority represents a valid priority for a ticket.
type TicketPriority struct {
	ID    string
	Label string
	Level int
}
