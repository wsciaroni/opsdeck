package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/wsciaroni/opsdeck/internal/core/domain"
	"github.com/wsciaroni/opsdeck/internal/core/port"
)

// CommentService implements business logic for comments.
type CommentService struct {
	repo port.CommentRepository
}

// NewCommentService creates a new CommentService.
func NewCommentService(repo port.CommentRepository) *CommentService {
	return &CommentService{repo: repo}
}

// CreateComment creates a new comment.
func (s *CommentService) CreateComment(ctx context.Context, cmd port.CreateCommentCmd) (*domain.Comment, error) {
	if cmd.Body == "" {
		return nil, fmt.Errorf("comment body cannot be empty")
	}
	if cmd.TicketID == uuid.Nil {
		return nil, fmt.Errorf("ticket_id is required")
	}
	if cmd.UserID == uuid.Nil {
		return nil, fmt.Errorf("user_id is required")
	}

	comment := &domain.Comment{
		TicketID:  cmd.TicketID,
		UserID:    cmd.UserID,
		Body:      cmd.Body,
		Sensitive: cmd.Sensitive,
	}

	if err := s.repo.Create(ctx, comment); err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	return comment, nil
}

// ListComments lists comments for a ticket.
func (s *CommentService) ListComments(ctx context.Context, ticketID uuid.UUID, includeSensitive bool) ([]domain.Comment, error) {
	return s.repo.ListByTicket(ctx, ticketID, includeSensitive)
}
