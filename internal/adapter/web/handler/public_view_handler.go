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

	// Hydrate with Assignee names? Maybe overkill for list, but good for completeness.
	// For now, let's just return tickets as is.
    // Wait, domain.Ticket doesn't have assignee name, just ID.
    // If I want to show names, I need to fetch users.

    // Collect user IDs (AssigneeID)
    userIDs := make(map[uuid.UUID]bool)
    for _, t := range tickets {
        if t.AssigneeUserID != nil {
            userIDs[*t.AssigneeUserID] = true
        }
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

    // Enhance response
    type TicketResponse struct {
        domain.Ticket
        AssigneeName string `json:"assignee_name,omitempty"`
    }

    respList := make([]TicketResponse, 0, len(tickets))
    for _, t := range tickets {
        tr := TicketResponse{Ticket: t}
        if t.AssigneeUserID != nil {
            if u, ok := users[*t.AssigneeUserID]; ok {
                tr.AssigneeName = u.Name
            }
        }
        respList = append(respList, tr)
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

    // Hydrate Reporter and Assignee names
    var reporterName, assigneeName string

    reporter, err := h.userRepo.GetByID(r.Context(), ticket.ReporterID)
    if err == nil && reporter != nil {
        reporterName = reporter.Name
    }

    if ticket.AssigneeUserID != nil {
        assignee, err := h.userRepo.GetByID(r.Context(), *ticket.AssigneeUserID)
        if err == nil && assignee != nil {
            assigneeName = assignee.Name
        }
    }

    resp := struct {
        domain.Ticket
        ReporterName string `json:"reporter_name"`
        AssigneeName string `json:"assignee_name,omitempty"`
    }{
        Ticket: *ticket,
        ReporterName: reporterName,
        AssigneeName: assigneeName,
    }

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
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
