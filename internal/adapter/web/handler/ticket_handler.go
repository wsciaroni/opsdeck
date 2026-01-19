package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/wsciaroni/opsdeck/internal/adapter/web/middleware"
	"github.com/wsciaroni/opsdeck/internal/core/port"
)

type TicketHandler struct {
	service port.TicketService
	orgRepo port.OrganizationRepository
	logger  *slog.Logger
}

type CreateTicketRequest struct {
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Priority       string    `json:"priority_id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	Location       string    `json:"location"`
}

type UpdateTicketRequest struct {
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	Priority    *string    `json:"priority"`
	Status      *string    `json:"status"`
	AssigneeID  *uuid.UUID `json:"assignee_id"`
	Location    *string    `json:"location"`
}

func NewTicketHandler(service port.TicketService, orgRepo port.OrganizationRepository, logger *slog.Logger) *TicketHandler {
	return &TicketHandler{
		service: service,
		orgRepo: orgRepo,
		logger:  logger,
	}
}

func (h *TicketHandler) CreateTicket(w http.ResponseWriter, r *http.Request) {
	var req CreateTicketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user := middleware.GetUser(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Security Check: Verify user belongs to the organization
	memberships, err := h.orgRepo.ListByUser(r.Context(), user.ID)
	if err != nil {
		h.logger.Error("failed to list user memberships", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	isMember := false
	for _, m := range memberships {
		if m.ID == req.OrganizationID {
			isMember = true
			break
		}
	}

	if !isMember {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	cmd := port.CreateTicketCmd{
		OrganizationID: req.OrganizationID,
		ReporterID:     user.ID,
		Title:          req.Title,
		Description:    req.Description,
		Location:       req.Location,
		PriorityID:     req.Priority,
	}

	ticket, err := h.service.CreateTicket(r.Context(), cmd)
	if err != nil {
		h.logger.Error("failed to create ticket", "error", err)
		http.Error(w, "Failed to create ticket", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(ticket); err != nil {
		h.logger.Error("failed to encode response", "error", err)
	}
}

func (h *TicketHandler) ListTickets(w http.ResponseWriter, r *http.Request) {
	orgIDStr := r.URL.Query().Get("organization_id")
	if orgIDStr == "" {
		http.Error(w, "organization_id is required", http.StatusBadRequest)
		return
	}

	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		http.Error(w, "Invalid organization_id", http.StatusBadRequest)
		return
	}

	user := middleware.GetUser(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Security Check: Verify user belongs to the organization
	memberships, err := h.orgRepo.ListByUser(r.Context(), user.ID)
	if err != nil {
		h.logger.Error("failed to list user memberships", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	isMember := false
	for _, m := range memberships {
		if m.ID == orgID {
			isMember = true
			break
		}
	}

	if !isMember {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	filter := port.TicketFilter{
		OrganizationID: orgID,
	}

	if status := r.URL.Query().Get("status"); status != "" {
		filter.StatusID = &status
	}

	if assigneeIDStr := r.URL.Query().Get("assignee_id"); assigneeIDStr != "" {
		assigneeID, err := uuid.Parse(assigneeIDStr)
		if err != nil {
			http.Error(w, "Invalid assignee_id", http.StatusBadRequest)
			return
		}
		filter.AssigneeID = &assigneeID
	}

	tickets, err := h.service.ListTickets(r.Context(), filter)
	if err != nil {
		h.logger.Error("failed to list tickets", "error", err)
		http.Error(w, "Failed to list tickets", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(tickets); err != nil {
		h.logger.Error("failed to encode response", "error", err)
	}
}

func (h *TicketHandler) UpdateTicket(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "ticketID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid ticket ID", http.StatusBadRequest)
		return
	}

	var req UpdateTicketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Security Check: Get existing ticket and verify user membership
	ticket, err := h.service.GetTicket(r.Context(), id)
	if err != nil {
		h.logger.Error("failed to get ticket", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if ticket == nil {
		http.Error(w, "Ticket not found", http.StatusNotFound)
		return
	}

	user := middleware.GetUser(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	memberships, err := h.orgRepo.ListByUser(r.Context(), user.ID)
	if err != nil {
		h.logger.Error("failed to list user memberships", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	isMember := false
	for _, m := range memberships {
		if m.ID == ticket.OrganizationID {
			isMember = true
			break
		}
	}

	if !isMember {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	cmd := port.UpdateTicketCmd{
		StatusID:       req.Status,
		PriorityID:     req.Priority,
		AssigneeUserID: req.AssigneeID,
		Title:          req.Title,
		Description:    req.Description,
		Location:       req.Location,
	}

	updatedTicket, err := h.service.UpdateTicket(r.Context(), id, cmd)
	if err != nil {
		h.logger.Error("failed to update ticket", "error", err)
		// Check for specific validation errors if needed, for now 500
		http.Error(w, "Failed to update ticket", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(updatedTicket); err != nil {
		h.logger.Error("failed to encode response", "error", err)
	}
}
