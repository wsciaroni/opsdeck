package notification

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
)

type EmailNotificationService struct {
	logger *slog.Logger
}

func NewEmailNotificationService(logger *slog.Logger) *EmailNotificationService {
	return &EmailNotificationService{
		logger: logger,
	}
}

func (s *EmailNotificationService) NotifyUserAddedToOrg(ctx context.Context, email string, orgName string) error {
	s.logger.Info("Sending email: User added to organization", "email", email, "org", orgName)
	return nil
}

func (s *EmailNotificationService) NotifyTicketAssigned(ctx context.Context, email string, ticketTitle string, ticketID uuid.UUID) error {
	s.logger.Info("Sending email: Ticket assigned", "email", email, "ticket_title", ticketTitle, "ticket_id", ticketID)
	return nil
}
