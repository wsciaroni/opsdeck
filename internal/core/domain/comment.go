package domain

import (
	"time"

	"github.com/google/uuid"
)

// Comment represents a comment on a ticket.
type Comment struct {
	ID        uuid.UUID `json:"id"`
	TicketID  uuid.UUID `json:"ticket_id"`
	UserID    uuid.UUID `json:"user_id"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}
