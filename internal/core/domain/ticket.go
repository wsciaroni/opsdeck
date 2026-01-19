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
	ID             uuid.UUID
	OrganizationID uuid.UUID
	Title          string
	Description    string
	Location       string
	StatusID       string
	PriorityID     string
	ReporterID     uuid.UUID
	AssigneeUserID *uuid.UUID
	CreatedAt      time.Time
	UpdatedAt      time.Time
	CompletedAt    *time.Time
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
