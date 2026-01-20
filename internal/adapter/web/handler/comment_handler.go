package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/wsciaroni/opsdeck/internal/adapter/web/middleware"
	"github.com/wsciaroni/opsdeck/internal/core/domain"
	"github.com/wsciaroni/opsdeck/internal/core/port"
	"github.com/wsciaroni/opsdeck/internal/core/service"
)

type CommentHandler struct {
	commentService port.CommentService
	ticketService  *service.TicketService
	userRepo       port.UserRepository
	orgRepo        port.OrganizationRepository
	logger         *slog.Logger
}

func NewCommentHandler(
	commentService port.CommentService,
	ticketService *service.TicketService,
	userRepo port.UserRepository,
	orgRepo port.OrganizationRepository,
	logger *slog.Logger,
) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
		ticketService:  ticketService,
		userRepo:       userRepo,
		orgRepo:        orgRepo,
		logger:         logger,
	}
}

type CreateCommentRequest struct {
	Body string `json:"body"`
}

type UserSummary struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	AvatarURL string    `json:"avatar_url"`
}

type CommentResponse struct {
	ID        uuid.UUID   `json:"id"`
	Body      string      `json:"body"`
	CreatedAt time.Time   `json:"created_at"`
	User      UserSummary `json:"user"`
}

func (h *CommentHandler) Create(w http.ResponseWriter, r *http.Request) {
	ticketIDStr := chi.URLParam(r, "ticketID")
	ticketID, err := uuid.Parse(ticketIDStr)
	if err != nil {
		http.Error(w, "Invalid ticket ID", http.StatusBadRequest)
		return
	}

	user := middleware.GetUser(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Check access to ticket
	ticket, err := h.ticketService.GetTicket(r.Context(), ticketID)
	if err != nil {
		h.logger.Error("Failed to get ticket", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if ticket == nil {
		http.Error(w, "Ticket not found", http.StatusNotFound)
		return
	}

	// Verify user is member of organization
	if err := h.checkOrgAccess(r.Context(), user.ID, ticket.OrganizationID); err != nil {
		h.logger.Warn("Unauthorized access attempt to ticket comment", "user_id", user.ID, "ticket_id", ticketID, "error", err)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var req CreateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	cmd := port.CreateCommentCmd{
		TicketID: ticketID,
		UserID:   user.ID,
		Body:     req.Body,
	}

	comment, err := h.commentService.CreateComment(r.Context(), cmd)
	if err != nil {
		h.logger.Error("Failed to create comment", "error", err)
		http.Error(w, "Failed to create comment", http.StatusInternalServerError)
		return
	}

	// We need to return the User details too for the UI to update immediately
	resp := CommentResponse{
		ID:        comment.ID,
		Body:      comment.Body,
		CreatedAt: comment.CreatedAt,
		User: UserSummary{
			ID:        user.ID,
			Name:      user.Name,
			AvatarURL: user.AvatarURL,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.Error("Failed to encode response", "error", err)
	}
}

func (h *CommentHandler) List(w http.ResponseWriter, r *http.Request) {
	ticketIDStr := chi.URLParam(r, "ticketID")
	ticketID, err := uuid.Parse(ticketIDStr)
	if err != nil {
		http.Error(w, "Invalid ticket ID", http.StatusBadRequest)
		return
	}

	user := middleware.GetUser(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Check access to ticket
	ticket, err := h.ticketService.GetTicket(r.Context(), ticketID)
	if err != nil {
		h.logger.Error("Failed to get ticket", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if ticket == nil {
		http.Error(w, "Ticket not found", http.StatusNotFound)
		return
	}

	// Verify user is member of organization
	if err := h.checkOrgAccess(r.Context(), user.ID, ticket.OrganizationID); err != nil {
		h.logger.Warn("Unauthorized access attempt to list comments", "user_id", user.ID, "ticket_id", ticketID, "error", err)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	comments, err := h.commentService.ListComments(r.Context(), ticketID)
	if err != nil {
		h.logger.Error("Failed to list comments", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Fetch users
	// Collect User IDs
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
			CreatedAt: c.CreatedAt,
			User:      userSummary,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(respList); err != nil {
		h.logger.Error("Failed to encode response", "error", err)
	}
}

func (h *CommentHandler) checkOrgAccess(ctx context.Context, userID, orgID uuid.UUID) error {
	memberships, err := h.orgRepo.ListByUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to list user memberships: %w", err)
	}

	for _, m := range memberships {
		// UserMembership embeds Organization, so we access OrganizationID via m.ID (since UserMembership embeds Organization struct)
		// Wait, UserMembership struct is:
		// type UserMembership struct {
		// 	Organization
		// 	Role string
		// }
		// So m.ID is the Organization ID.
		if m.ID == orgID {
			return nil
		}
	}

	return fmt.Errorf("user is not a member of organization %s", orgID)
}
