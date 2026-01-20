package handler

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log/slog"
	"net/http"
	"regexp"
	"strings"

	"github.com/wsciaroni/opsdeck/internal/adapter/web/middleware"
	"github.com/wsciaroni/opsdeck/internal/core/domain"
	"github.com/wsciaroni/opsdeck/internal/core/port"
)

type OrgHandler struct {
	orgRepo port.OrganizationRepository
	logger  *slog.Logger
}

func NewOrgHandler(orgRepo port.OrganizationRepository, logger *slog.Logger) *OrgHandler {
	return &OrgHandler{
		orgRepo: orgRepo,
		logger:  logger,
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
