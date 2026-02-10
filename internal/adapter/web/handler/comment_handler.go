package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/wsciaroni/opsdeck/internal/adapter/web/middleware"
	"github.com/wsciaroni/opsdeck/internal/core/domain"
	"github.com/wsciaroni/opsdeck/internal/core/port"
)

type CommentHandler struct {
	commentService port.CommentService
	ticketService  port.TicketService
	userRepo       port.UserRepository
	orgRepo        port.OrganizationRepository
	logger         *slog.Logger
}

func NewCommentHandler(
	commentService port.CommentService,
	ticketService port.TicketService,
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
	Body      string `json:"body"`
	Sensitive bool   `json:"sensitive"`
}

type UserSummary struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	AvatarURL string    `json:"avatar_url"`
}

type CommentResponse struct {
	ID        uuid.UUID   `json:"id"`
	Body      string      `json:"body"`
	Sensitive bool        `json:"sensitive"`
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

	// Limit request size to 1MB to prevent DoS
	r.Body = http.MaxBytesReader(w, r.Body, MaxJSONBodySize)

	var req CreateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		if strings.Contains(err.Error(), "request body too large") {
			http.Error(w, "Request too large", http.StatusRequestEntityTooLarge)
			return
		}
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if len(req.Body) > 5000 {
		http.Error(w, "Comment too long", http.StatusBadRequest)
		return
	}

	cmd := port.CreateCommentCmd{
		TicketID:  ticketID,
		UserID:    user.ID,
		Body:      req.Body,
		Sensitive: req.Sensitive,
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
		Sensitive: comment.Sensitive,
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

	comments, err := h.commentService.ListComments(r.Context(), ticketID, true)
	if err != nil {
		h.logger.Error("Failed to list comments", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Fetch users
	// Collect unique User IDs
	uniqueUserIDs := make([]uuid.UUID, 0)
	seenUserIDs := make(map[uuid.UUID]bool)
	for _, c := range comments {
		if !seenUserIDs[c.UserID] {
			uniqueUserIDs = append(uniqueUserIDs, c.UserID)
			seenUserIDs[c.UserID] = true
		}
	}

	users := make(map[uuid.UUID]*domain.User)
	if len(uniqueUserIDs) > 0 {
		fetchedUsers, err := h.userRepo.GetByIDs(r.Context(), uniqueUserIDs)
		if err != nil {
			h.logger.Error("Failed to fetch users", "error", err)
			// Continue with partial/empty users map
		} else {
			for i := range fetchedUsers {
				u := &fetchedUsers[i]
				users[u.ID] = u
			}
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

func (h *CommentHandler) checkOrgAccess(ctx context.Context, userID, orgID uuid.UUID) error {
	role, err := h.orgRepo.GetMemberRole(ctx, orgID, userID)
	if err != nil {
		return fmt.Errorf("failed to get member role: %w", err)
	}
	if role == "" {
		return fmt.Errorf("user is not a member of organization %s", orgID)
	}
	return nil
}
