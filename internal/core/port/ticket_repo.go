package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/wsciaroni/opsdeck/internal/core/domain"
)

// TicketFilter defines criteria for listing tickets.
type TicketFilter struct {
	OrganizationID     *uuid.UUID
	OrganizationIDs    []uuid.UUID
	StatusIDs          []string
	PriorityIDs        []string
	AssigneeID         *uuid.UUID
	ReporterID         *uuid.UUID
	ExcludeDescription bool
	Sensitive          *bool
	Keyword            *string
}

// TicketRepository defines the interface for interacting with ticket data.
type TicketRepository interface {
	Create(ctx context.Context, ticket *domain.Ticket) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Ticket, error)
	List(ctx context.Context, filter TicketFilter) ([]domain.Ticket, error)
	Update(ctx context.Context, ticket *domain.Ticket) error

	AddFile(ctx context.Context, file *domain.File) error
	GetFile(ctx context.Context, id uuid.UUID) (*domain.File, error)
	ListFiles(ctx context.Context, ticketID uuid.UUID) ([]domain.File, error)
}
