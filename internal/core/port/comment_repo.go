package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/wsciaroni/opsdeck/internal/core/domain"
)

// CommentRepository defines the interface for interacting with comment data.
type CommentRepository interface {
	Create(ctx context.Context, comment *domain.Comment) error
	ListByTicket(ctx context.Context, ticketID uuid.UUID) ([]domain.Comment, error)
}
