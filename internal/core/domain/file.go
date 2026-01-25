package domain

import (
	"time"

	"github.com/google/uuid"
)

// File represents an uploaded file associated with a ticket.
type File struct {
	ID          uuid.UUID `json:"id"`
	TicketID    uuid.UUID `json:"ticket_id"`
	Filename    string    `json:"filename"`
	ContentType string    `json:"content_type"`
	Size        int64     `json:"size"`
	Data        []byte    `json:"-"` // Don't expose data in JSON responses by default
	CreatedAt   time.Time `json:"created_at"`
}
