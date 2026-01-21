package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/wsciaroni/opsdeck/internal/core/domain"
	"github.com/wsciaroni/opsdeck/internal/core/port"
)

// TicketService implements business logic for ticket management.
type TicketService struct {
	repo port.TicketRepository
}

// NewTicketService creates a new TicketService.
func NewTicketService(repo port.TicketRepository) *TicketService {
	return &TicketService{repo: repo}
}

// GetTicket retrieves a ticket by its ID.
func (s *TicketService) GetTicket(ctx context.Context, id uuid.UUID) (*domain.Ticket, error) {
	return s.repo.GetByID(ctx, id)
}

// CreateTicket creates a new ticket.
func (s *TicketService) CreateTicket(ctx context.Context, cmd port.CreateTicketCmd) (*domain.Ticket, error) {
	if cmd.Title == "" {
		return nil, fmt.Errorf("title is required")
	}

	if !isValidPriority(cmd.PriorityID) {
		return nil, fmt.Errorf("invalid priority: %s", cmd.PriorityID)
	}

	ticket := &domain.Ticket{
		OrganizationID: cmd.OrganizationID,
		ReporterID:     cmd.ReporterID,
		Title:          cmd.Title,
		Description:    cmd.Description,
		Location:       cmd.Location,
		StatusID:       domain.TicketStatusNew,
		PriorityID:     cmd.PriorityID,
	}

	if err := s.repo.Create(ctx, ticket); err != nil {
		return nil, fmt.Errorf("failed to create ticket: %w", err)
	}

	return ticket, nil
}

// UpdateTicket updates an existing ticket.
func (s *TicketService) UpdateTicket(ctx context.Context, id uuid.UUID, cmd port.UpdateTicketCmd) (*domain.Ticket, error) {
	ticket, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get ticket: %w", err)
	}
	if ticket == nil {
		return nil, fmt.Errorf("ticket not found")
	}

	oldStatus := ticket.StatusID

	if cmd.Title != nil {
		ticket.Title = *cmd.Title
	}
	if cmd.Description != nil {
		ticket.Description = *cmd.Description
	}
	if cmd.Location != nil {
		ticket.Location = *cmd.Location
	}
	if cmd.PriorityID != nil {
		if !isValidPriority(*cmd.PriorityID) {
			return nil, fmt.Errorf("invalid priority: %s", *cmd.PriorityID)
		}
		ticket.PriorityID = *cmd.PriorityID
	}
	if cmd.AssigneeUserID != nil {
		if *cmd.AssigneeUserID == uuid.Nil {
			ticket.AssigneeUserID = nil
		} else {
			ticket.AssigneeUserID = cmd.AssigneeUserID
		}
	}

	if cmd.StatusID != nil {
		if !isValidStatus(*cmd.StatusID) {
			return nil, fmt.Errorf("invalid status: %s", *cmd.StatusID)
		}
		ticket.StatusID = *cmd.StatusID
	}

	// Handle CompletedAt Logic
	if oldStatus != ticket.StatusID {
		if ticket.StatusID == domain.TicketStatusDone || ticket.StatusID == domain.TicketStatusCanceled {
			now := time.Now()
			ticket.CompletedAt = &now
		} else if oldStatus == domain.TicketStatusDone && ticket.StatusID == domain.TicketStatusNew {
			ticket.CompletedAt = nil
		}
	}

	ticket.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, ticket); err != nil {
		return nil, fmt.Errorf("failed to update ticket: %w", err)
	}

	return ticket, nil
}

// ListTickets lists tickets based on the filter.
func (s *TicketService) ListTickets(ctx context.Context, filter port.TicketFilter) ([]domain.Ticket, error) {
	// If OrganizationID is nil, we assume the caller (Handler) has already authorized the user to see all tickets
	// So we don't enforce it here anymore, as TicketRepository now supports nil.
	return s.repo.List(ctx, filter)
}

func isValidPriority(p string) bool {
	switch p {
	case domain.TicketPriorityLow, domain.TicketPriorityMedium, domain.TicketPriorityHigh, domain.TicketPriorityCritical:
		return true
	default:
		return false
	}
}

func isValidStatus(s string) bool {
	switch s {
	case domain.TicketStatusNew, domain.TicketStatusInProgress, domain.TicketStatusOnHold, domain.TicketStatusDone, domain.TicketStatusCanceled:
		return true
	default:
		return false
	}
}
