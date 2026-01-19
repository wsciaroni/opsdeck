package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/wsciaroni/opsdeck/internal/core/domain"
)

// CreateTicketCmd defines the command to create a new ticket.
type CreateTicketCmd struct {
	OrganizationID uuid.UUID
	ReporterID     uuid.UUID
	Title          string
	Description    string
	Location       string
	PriorityID     string
}

// UpdateTicketCmd defines the command to update an existing ticket.
type UpdateTicketCmd struct {
	StatusID       *string
	PriorityID     *string
	AssigneeUserID *uuid.UUID
	Title          *string
	Description    *string
	Location       *string
}

// TicketService defines the interface for ticket business logic.
type TicketService interface {
	CreateTicket(ctx context.Context, cmd CreateTicketCmd) (*domain.Ticket, error)
	UpdateTicket(ctx context.Context, id uuid.UUID, cmd UpdateTicketCmd) (*domain.Ticket, error)
	ListTickets(ctx context.Context, filter TicketFilter) ([]domain.Ticket, error)
	GetTicket(ctx context.Context, id uuid.UUID) (*domain.Ticket, error)
}
