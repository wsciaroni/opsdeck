package port

import (
	"context"

	"github.com/google/uuid"
)

// NotificationService defines the interface for sending notifications.
type NotificationService interface {
	NotifyUserAddedToOrg(ctx context.Context, email string, orgName string) error
	NotifyTicketAssigned(ctx context.Context, email string, ticketTitle string, ticketID uuid.UUID) error
}
