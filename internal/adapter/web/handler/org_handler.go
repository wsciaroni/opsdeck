package handler

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log/slog"
	"net/http"
	"regexp"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/wsciaroni/opsdeck/internal/adapter/web/middleware"
	"github.com/wsciaroni/opsdeck/internal/core/domain"
	"github.com/wsciaroni/opsdeck/internal/core/port"
)

type OrgHandler struct {
	orgRepo  port.OrganizationRepository
	userRepo port.UserRepository
	logger   *slog.Logger
}

func NewOrgHandler(orgRepo port.OrganizationRepository, userRepo port.UserRepository, logger *slog.Logger) *OrgHandler {
	return &OrgHandler{
		orgRepo:  orgRepo,
		userRepo: userRepo,
		logger:   logger,
	}
}

type CreateOrgRequest struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

func (h *OrgHandler) CreateOrganization(w http.ResponseWriter, r *http.Request) {
	var req CreateOrgRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	// Generate Slug if not provided
	if req.Slug == "" {
		req.Slug = generateSlug(req.Name)
	}

	user := middleware.GetUser(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	org := &domain.Organization{
		Name: req.Name,
		Slug: req.Slug,
	}

	if err := h.orgRepo.Create(r.Context(), org); err != nil {
		h.logger.Error("failed to create organization", "error", err)
		http.Error(w, "Failed to create organization", http.StatusInternalServerError)
		return
	}

	// Add Creator as Owner
	if err := h.orgRepo.AddMember(r.Context(), org.ID, user.ID, "owner"); err != nil {
		h.logger.Error("failed to add owner to organization", "error", err)
		// Note: Ideally this should be transactional or we should roll back the org creation.
		// For MVP, we'll accept the risk of orphan orgs or handle it manually.
		http.Error(w, "Failed to assign ownership", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(org); err != nil {
		h.logger.Error("failed to encode response", "error", err)
	}
}

type AddMemberRequest struct {
	Email string `json:"email"`
}

func (h *OrgHandler) AddMember(w http.ResponseWriter, r *http.Request) {
	// Parse Org ID
	orgIDStr := chi.URLParam(r, "id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		http.Error(w, "Invalid Organization ID", http.StatusBadRequest)
		return
	}

	// Parse Request Body
	var req AddMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if req.Email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	// Check Auth
	currentUser := middleware.GetUser(r.Context())
	if currentUser == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Check Permissions (Must be owner or admin of the org)
	memberships, err := h.orgRepo.ListByUser(r.Context(), currentUser.ID)
	if err != nil {
		h.logger.Error("failed to list user memberships", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	isAuthorized := false
	for _, m := range memberships {
		if m.Organization.ID == orgID && (m.Role == "owner" || m.Role == "admin") {
			isAuthorized = true
			break
		}
	}

	if !isAuthorized {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Find User to Add
	userToAdd, err := h.userRepo.GetByEmail(r.Context(), req.Email)
	if err != nil {
		h.logger.Error("failed to get user by email", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if userToAdd == nil {
		http.Error(w, "User must register first", http.StatusNotFound)
		return
	}

	// Check if already a member?
	// For MVP, we'll let the DB constraint handle it or just add.
	// DB likely has (org_id, user_id) unique constraint.
	// If so, AddMember might fail.
	// Let's assume we can try to add.

	// Add Member (default role 'member')
	if err := h.orgRepo.AddMember(r.Context(), orgID, userToAdd.ID, "member"); err != nil {
		// Check for duplicate key error if possible, but for MVP generic error log is fine
		h.logger.Error("failed to add member", "error", err)
		http.Error(w, "Failed to add member", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *OrgHandler) ListMembers(w http.ResponseWriter, r *http.Request) {
	// Parse Org ID
	orgIDStr := chi.URLParam(r, "id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		http.Error(w, "Invalid Organization ID", http.StatusBadRequest)
		return
	}

	// Check Auth
	currentUser := middleware.GetUser(r.Context())
	if currentUser == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Check Membership (Any role can view)
	memberships, err := h.orgRepo.ListByUser(r.Context(), currentUser.ID)
	if err != nil {
		h.logger.Error("failed to list user memberships", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	isMember := false
	for _, m := range memberships {
		if m.Organization.ID == orgID {
			isMember = true
			break
		}
	}

	if !isMember {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	members, err := h.orgRepo.ListMembers(r.Context(), orgID)
	if err != nil {
		h.logger.Error("failed to list organization members", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(members); err != nil {
		h.logger.Error("failed to encode members response", "error", err)
	}
}

func (h *OrgHandler) RemoveMember(w http.ResponseWriter, r *http.Request) {
	// Parse Org ID
	orgIDStr := chi.URLParam(r, "id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		http.Error(w, "Invalid Organization ID", http.StatusBadRequest)
		return
	}

	// Parse User ID to remove
	userIDStr := chi.URLParam(r, "userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Invalid User ID", http.StatusBadRequest)
		return
	}

	// Check Auth
	currentUser := middleware.GetUser(r.Context())
	if currentUser == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Check Permissions
	// 1. Owner can remove anyone.
	// 2. User can remove themselves (leave org).
	memberships, err := h.orgRepo.ListByUser(r.Context(), currentUser.ID)
	if err != nil {
		h.logger.Error("failed to list user memberships", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var currentMemberRole string
	isMember := false
	for _, m := range memberships {
		if m.Organization.ID == orgID {
			isMember = true
			currentMemberRole = m.Role
			break
		}
	}

	if !isMember {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	canRemove := false
	if currentUser.ID == userID {
		canRemove = true // Removing self
	} else if currentMemberRole == "owner" {
		canRemove = true // Owner removing others
	}
	// Admin removing members? Not specified in MVP requirements, but usually admins can remove members.
	// Prompt says: "Check if requester is 'owner' (or removing themselves)."
	// I will stick to "owner" or "self".

	if !canRemove {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if err := h.orgRepo.RemoveMember(r.Context(), orgID, userID); err != nil {
		h.logger.Error("failed to remove member", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func generateSlug(name string) string {
	// Slugify: lowercase, replace spaces with dashes, remove non-alphanumeric
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	reg := regexp.MustCompile("[^a-z0-9-]+")
	slug = reg.ReplaceAllString(slug, "")

	// Append random hex
	bytes := make([]byte, 2) // 2 bytes = 4 hex chars
	if _, err := rand.Read(bytes); err != nil {
		// Fallback if random fails (unlikely)
		return slug + "-temp"
	}
	return slug + "-" + hex.EncodeToString(bytes)
}

func (h *OrgHandler) GetShareSettings(w http.ResponseWriter, r *http.Request) {
	orgIDStr := chi.URLParam(r, "id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		http.Error(w, "Invalid Organization ID", http.StatusBadRequest)
		return
	}

	currentUser := middleware.GetUser(r.Context())
	if currentUser == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if !h.isMember(r.Context(), orgID, currentUser.ID) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	org, err := h.orgRepo.GetByID(r.Context(), orgID)
	if err != nil {
		h.logger.Error("failed to get organization", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	resp := struct {
		ShareLinkEnabled bool    `json:"share_link_enabled"`
		ShareLinkToken   *string `json:"share_link_token"`
	}{
		ShareLinkEnabled: org.ShareLinkEnabled,
		ShareLinkToken:   org.ShareLinkToken,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.Error("failed to encode response", "error", err)
	}
}

type UpdateShareSettingsRequest struct {
	Enabled bool `json:"enabled"`
}

func (h *OrgHandler) UpdateShareSettings(w http.ResponseWriter, r *http.Request) {
	orgIDStr := chi.URLParam(r, "id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		http.Error(w, "Invalid Organization ID", http.StatusBadRequest)
		return
	}

	var req UpdateShareSettingsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	currentUser := middleware.GetUser(r.Context())
	if currentUser == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Only owner/admin can update settings
	if !h.isAdminOrOwner(r.Context(), orgID, currentUser.ID) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	org, err := h.orgRepo.GetByID(r.Context(), orgID)
	if err != nil {
		h.logger.Error("failed to get organization", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	org.ShareLinkEnabled = req.Enabled
	if req.Enabled && org.ShareLinkToken == nil {
		token := generateToken()
		org.ShareLinkToken = &token
	}

	if err := h.orgRepo.Update(r.Context(), org); err != nil {
		h.logger.Error("failed to update organization", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	resp := struct {
		ShareLinkEnabled bool    `json:"share_link_enabled"`
		ShareLinkToken   *string `json:"share_link_token"`
	}{
		ShareLinkEnabled: org.ShareLinkEnabled,
		ShareLinkToken:   org.ShareLinkToken,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.Error("failed to encode response", "error", err)
	}
}

func (h *OrgHandler) RegenerateShareToken(w http.ResponseWriter, r *http.Request) {
	orgIDStr := chi.URLParam(r, "id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		http.Error(w, "Invalid Organization ID", http.StatusBadRequest)
		return
	}

	currentUser := middleware.GetUser(r.Context())
	if currentUser == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Only owner/admin can update settings
	if !h.isAdminOrOwner(r.Context(), orgID, currentUser.ID) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	org, err := h.orgRepo.GetByID(r.Context(), orgID)
	if err != nil {
		h.logger.Error("failed to get organization", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	token := generateToken()
	org.ShareLinkToken = &token

	if err := h.orgRepo.Update(r.Context(), org); err != nil {
		h.logger.Error("failed to update organization", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	resp := struct {
		ShareLinkEnabled bool    `json:"share_link_enabled"`
		ShareLinkToken   *string `json:"share_link_token"`
	}{
		ShareLinkEnabled: org.ShareLinkEnabled,
		ShareLinkToken:   org.ShareLinkToken,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.Error("failed to encode response", "error", err)
	}
}

func (h *OrgHandler) isMember(ctx context.Context, orgID, userID uuid.UUID) bool {
	memberships, err := h.orgRepo.ListByUser(ctx, userID)
	if err != nil {
		return false
	}
	for _, m := range memberships {
		if m.Organization.ID == orgID {
			return true
		}
	}
	return false
}

func (h *OrgHandler) isAdminOrOwner(ctx context.Context, orgID, userID uuid.UUID) bool {
	memberships, err := h.orgRepo.ListByUser(ctx, userID)
	if err != nil {
		return false
	}
	for _, m := range memberships {
		if m.Organization.ID == orgID && (m.Role == "owner" || m.Role == "admin") {
			return true
		}
	}
	return false
}

func generateToken() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
