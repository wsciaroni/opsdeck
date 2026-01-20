package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/wsciaroni/opsdeck/internal/core/domain"
)

// CreateCommentCmd defines the input for creating a comment.
type CreateCommentCmd struct {
	TicketID uuid.UUID
	UserID   uuid.UUID
	Body     string
}

// CommentService defines the interface for comment business logic.
type CommentService interface {
	CreateComment(ctx context.Context, cmd CreateCommentCmd) (*domain.Comment, error)
	ListComments(ctx context.Context, ticketID uuid.UUID) ([]domain.Comment, error)
}
