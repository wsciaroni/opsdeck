package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wsciaroni/opsdeck/internal/core/domain"
)

type CommentRepository struct {
	db *pgxpool.Pool
}

func NewCommentRepository(db *pgxpool.Pool) *CommentRepository {
	return &CommentRepository{db: db}
}

func (r *CommentRepository) Create(ctx context.Context, comment *domain.Comment) error {
	query := `
		INSERT INTO comments (ticket_id, user_id, body, sensitive)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	err := r.db.QueryRow(ctx, query,
		comment.TicketID,
		comment.UserID,
		comment.Body,
		comment.Sensitive,
	).Scan(&comment.ID, &comment.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create comment: %w", err)
	}

	return nil
}

func (r *CommentRepository) ListByTicket(ctx context.Context, ticketID uuid.UUID) ([]domain.Comment, error) {
	query := `
		SELECT id, ticket_id, user_id, body, sensitive, created_at
		FROM comments
		WHERE ticket_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(ctx, query, ticketID)
	if err != nil {
		return nil, fmt.Errorf("failed to list comments: %w", err)
	}
	defer rows.Close()

	var comments []domain.Comment
	for rows.Next() {
		var c domain.Comment
		err := rows.Scan(
			&c.ID,
			&c.TicketID,
			&c.UserID,
			&c.Body,
			&c.Sensitive,
			&c.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}
		comments = append(comments, c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return comments, nil
}
