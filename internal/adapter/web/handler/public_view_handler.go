package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/wsciaroni/opsdeck/internal/core/domain"
	"github.com/wsciaroni/opsdeck/internal/core/port"
	"github.com/wsciaroni/opsdeck/internal/core/service"
)

type PublicViewHandler struct {
	orgRepo        port.OrganizationRepository
	ticketService  *service.TicketService
	commentService port.CommentService
	userRepo       port.UserRepository
	logger         *slog.Logger
}

func NewPublicViewHandler(
	orgRepo port.OrganizationRepository,
	ticketService *service.TicketService,
	commentService port.CommentService,
	userRepo port.UserRepository,
	logger *slog.Logger,
) *PublicViewHandler {
	return &PublicViewHandler{
		orgRepo:        orgRepo,
		ticketService:  ticketService,
		commentService: commentService,
		userRepo:       userRepo,
		logger:         logger,
	}
}

func (h *PublicViewHandler) GetOrganization(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	if token == "" {
		http.Error(w, "Token is required", http.StatusBadRequest)
		return
	}

	org, err := h.orgRepo.GetByPublicViewToken(r.Context(), token)
	if err != nil {
		h.logger.Error("Failed to get organization by public view token", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if org == nil || !org.PublicViewEnabled {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	// Only return necessary public fields
	resp := struct {
		ID   uuid.UUID `json:"id"`
		Name string    `json:"name"`
		Slug string    `json:"slug"`
	}{
		ID:   org.ID,
		Name: org.Name,
		Slug: org.Slug,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.Error("Failed to encode response", "error", err)
	}
}

func (h *PublicViewHandler) ListTickets(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	if token == "" {
		http.Error(w, "Token is required", http.StatusBadRequest)
		return
	}

	org, err := h.orgRepo.GetByPublicViewToken(r.Context(), token)
	if err != nil {
		h.logger.Error("Failed to get organization by public view token", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if org == nil || !org.PublicViewEnabled {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	sensitive := false
	filter := port.TicketFilter{
		OrganizationID: &org.ID,
		Sensitive:      &sensitive,
	}

	tickets, err := h.ticketService.ListTickets(r.Context(), filter)
	if err != nil {
		h.logger.Error("Failed to list tickets", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Sanitize tickets (hide sensitive details like user IDs)
	respList := make([]domain.Ticket, 0, len(tickets))
	for _, t := range tickets {
		ticketCopy := t
		ticketCopy.ReporterID = uuid.Nil
		ticketCopy.AssigneeUserID = nil
		respList = append(respList, ticketCopy)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(respList); err != nil {
		h.logger.Error("Failed to encode response", "error", err)
	}
}

func (h *PublicViewHandler) GetTicket(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	ticketIDStr := chi.URLParam(r, "ticketID")
	ticketID, err := uuid.Parse(ticketIDStr)
	if err != nil {
		http.Error(w, "Invalid ticket ID", http.StatusBadRequest)
		return
	}

	org, err := h.orgRepo.GetByPublicViewToken(r.Context(), token)
	if err != nil {
		h.logger.Error("Failed to get organization by public view token", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if org == nil || !org.PublicViewEnabled {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	ticket, err := h.ticketService.GetTicket(r.Context(), ticketID)
	if err != nil {
		h.logger.Error("Failed to get ticket", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if ticket == nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	// Ensure ticket belongs to org and is not sensitive
	if ticket.OrganizationID != org.ID || ticket.Sensitive {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	// Sanitize ticket (hide sensitive details like user IDs)
	ticketCopy := *ticket
	ticketCopy.ReporterID = uuid.Nil
	ticketCopy.AssigneeUserID = nil

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(ticketCopy); err != nil {
		h.logger.Error("Failed to encode response", "error", err)
	}
}

func (h *PublicViewHandler) ListComments(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	ticketIDStr := chi.URLParam(r, "ticketID")
	ticketID, err := uuid.Parse(ticketIDStr)
	if err != nil {
		http.Error(w, "Invalid ticket ID", http.StatusBadRequest)
		return
	}

	org, err := h.orgRepo.GetByPublicViewToken(r.Context(), token)
	if err != nil {
		h.logger.Error("Failed to get organization by public view token", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if org == nil || !org.PublicViewEnabled {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	// Check ticket access first to ensure it's not sensitive and belongs to org
	ticket, err := h.ticketService.GetTicket(r.Context(), ticketID)
	if err != nil {
		h.logger.Error("Failed to get ticket", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if ticket == nil || ticket.OrganizationID != org.ID || ticket.Sensitive {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	comments, err := h.commentService.ListComments(r.Context(), ticketID, false)
	if err != nil {
		h.logger.Error("Failed to list comments", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Hydrate users
	userIDs := make(map[uuid.UUID]bool)
	for _, c := range comments {
		userIDs[c.UserID] = true
	}

	users := make(map[uuid.UUID]*domain.User)
	for uid := range userIDs {
		u, err := h.userRepo.GetByID(r.Context(), uid)
		if err != nil {
			h.logger.Error("Failed to fetch user", "userID", uid, "error", err)
			continue
		}
		if u != nil {
			users[uid] = u
		}
	}

	respList := make([]CommentResponse, 0, len(comments))
	for _, c := range comments {
		u, exists := users[c.UserID]
		userSummary := UserSummary{ID: c.UserID, Name: "Unknown", AvatarURL: ""}
		if exists {
			userSummary.Name = u.Name
			userSummary.AvatarURL = u.AvatarURL
		}

		respList = append(respList, CommentResponse{
			ID:        c.ID,
			Body:      c.Body,
			Sensitive: c.Sensitive,
			CreatedAt: c.CreatedAt,
			User:      userSummary,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(respList); err != nil {
		h.logger.Error("Failed to encode response", "error", err)
	}
}
