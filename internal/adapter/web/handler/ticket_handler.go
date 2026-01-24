package handler

import (
	"encoding/csv"
	"encoding/json"
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

type TicketHandler struct {
	service  port.TicketService
	orgRepo  port.OrganizationRepository
	userRepo port.UserRepository
	logger   *slog.Logger
}

type TicketDetailResponse struct {
	*domain.Ticket
	ReporterName string `json:"reporter_name"`
	AssigneeName string `json:"assignee_name"`
}

type CreateTicketRequest struct {
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Priority       string    `json:"priority_id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	Location       string    `json:"location"`
	Sensitive      bool      `json:"sensitive"`
}

type CreatePublicTicketRequest struct {
	Token       string `json:"token"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority_id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
}

type UpdateTicketRequest struct {
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	Priority    *string    `json:"priority_id"`
	Status      *string    `json:"status_id"`
	AssigneeID  *uuid.UUID `json:"assignee_id"`
	Location    *string    `json:"location"`
	Sensitive   *bool      `json:"sensitive"`
}

func NewTicketHandler(service port.TicketService, orgRepo port.OrganizationRepository, userRepo port.UserRepository, logger *slog.Logger) *TicketHandler {
	return &TicketHandler{
		service:  service,
		orgRepo:  orgRepo,
		userRepo: userRepo,
		logger:   logger,
	}
}

func (h *TicketHandler) GetTicket(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "ticketID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid ticket ID", http.StatusBadRequest)
		return
	}

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

	reporterName := ""
	reporter, err := h.userRepo.GetByID(r.Context(), ticket.ReporterID)
	if err != nil {
		h.logger.Error("failed to get reporter", "error", err)
		// continue without reporter name
	} else if reporter != nil {
		reporterName = reporter.Name
	}

	assigneeName := ""
	if ticket.AssigneeUserID != nil {
		assignee, err := h.userRepo.GetByID(r.Context(), *ticket.AssigneeUserID)
		if err != nil {
			h.logger.Error("failed to get assignee", "error", err)
			// continue without assignee name
		} else if assignee != nil {
			assigneeName = assignee.Name
		}
	}

	resp := TicketDetailResponse{
		Ticket:       ticket,
		ReporterName: reporterName,
		AssigneeName: assigneeName,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.Error("failed to encode response", "error", err)
	}
}

func (h *TicketHandler) CreatePublicTicket(w http.ResponseWriter, r *http.Request) {
	var req CreatePublicTicketRequest

	contentType := r.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "multipart/form-data") {
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}
		req.Token = r.FormValue("token")
		req.Title = r.FormValue("title")
		req.Description = r.FormValue("description")
		req.Name = r.FormValue("name")
		req.Email = r.FormValue("email")
		req.Priority = r.FormValue("priority_id")
	} else {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
	}

	// 1. Validate Token & Org
	org, err := h.orgRepo.GetByShareToken(r.Context(), req.Token)
	if err != nil {
		h.logger.Error("failed to get organization by token", "error", err)
		// Assuming error means not found or db error.
		// If org not found by token, it returns sql.ErrNoRows which might be wrapped.
		// For security, just say forbidden or invalid.
		http.Error(w, "Invalid token", http.StatusForbidden)
		return
	}

	if org == nil || !org.ShareLinkEnabled {
		http.Error(w, "Share link disabled", http.StatusForbidden)
		return
	}

	// 2. Find or Create User
	user, err := h.userRepo.GetByEmail(r.Context(), req.Email)
	if err != nil || user == nil {
		// Try to create user
		// Ideally we check if error is "not found".
		// But for now let's assume if error or user is nil, user doesn't exist or db error.
		// We'll try to create.
		newUser := &domain.User{
			Email: req.Email,
			Name:  req.Name,
			Role:  domain.RolePublic,
		}
		if err := h.userRepo.Create(r.Context(), newUser); err != nil {
			// If create fails, maybe race condition or other error
			// If it's a constraint violation (email exists), we should probably retry get?
			// But for now, let's log and fail.
			h.logger.Error("failed to create public user", "error", err)
			http.Error(w, "Failed to process user", http.StatusInternalServerError)
			return
		}
		user = newUser
	}

	// 3. Create Ticket
	cmd := port.CreateTicketCmd{
		OrganizationID: org.ID,
		ReporterID:     user.ID,
		Title:          req.Title,
		Description:    req.Description,
		PriorityID:     req.Priority,
		// Location? Not in public form?
	}

	ticket, err := h.service.CreateTicket(r.Context(), cmd)
	if err != nil {
		h.logger.Error("failed to create public ticket", "error", err)
		http.Error(w, "Failed to create ticket", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(ticket); err != nil {
		h.logger.Error("failed to encode response", "error", err)
	}
}

func (h *TicketHandler) ExportTickets(w http.ResponseWriter, r *http.Request) {
	// 1. Verify Authentication & Admin Role
	user := middleware.GetUser(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if user.Role != domain.RoleAdmin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// 2. Fetch User Memberships
	memberships, err := h.orgRepo.ListByUser(r.Context(), user.ID)
	if err != nil {
		h.logger.Error("failed to list user memberships", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if len(memberships) == 0 {
		// User has no organizations, return empty list (or forbidden? Empty is fine)
		// We'll proceed with an empty list filter which should return nothing if we handle it correctly in logic
		// But repo logic says "if OrganizationIDs not empty, filter by it". If empty, it might fall back to "WHERE 1=1" if we are not careful.
		// So we should handle this case specifically.
		h.writeEmptyCSV(w)
		return
	}

	// 3. Parse Filters
	var filter port.TicketFilter

	orgIDStr := r.URL.Query().Get("organization_id")
	if orgIDStr != "" {
		parsed, err := uuid.Parse(orgIDStr)
		if err != nil {
			http.Error(w, "Invalid organization_id", http.StatusBadRequest)
			return
		}

		// Verify membership
		isMember := false
		for _, m := range memberships {
			if m.ID == parsed {
				isMember = true
				break
			}
		}

		if !isMember {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		filter.OrganizationID = &parsed
	} else {
		// No specific org requested, filter by ALL memberships
		orgIDs := make([]uuid.UUID, len(memberships))
		for i, m := range memberships {
			orgIDs[i] = m.ID
		}
		filter.OrganizationIDs = orgIDs
	}

	// 4. Fetch Tickets
	tickets, err := h.service.ListTickets(r.Context(), filter)
	if err != nil {
		h.logger.Error("failed to list tickets for export", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// 5. Stream CSV Response
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=\"tickets.csv\"")

	// We are buffering the CSV write, but we could also write directly to w.
	// encoding/csv writes to an io.Writer.
	writer := csv.NewWriter(w)
	defer writer.Flush()

	// Write Header
	header := []string{"ID", "Organization ID", "Title", "Status", "Priority", "Reporter ID", "Created At", "Description"}
	if err := writer.Write(header); err != nil {
		h.logger.Error("failed to write csv header", "error", err)
		return
	}

	// Write Rows
	for _, t := range tickets {
		row := []string{
			t.ID.String(),
			t.OrganizationID.String(),
			t.Title,
			t.StatusID,
			t.PriorityID,
			t.ReporterID.String(),
			t.CreatedAt.Format(time.RFC3339),
			t.Description,
		}
		if err := writer.Write(row); err != nil {
			h.logger.Error("failed to write csv row", "error", err)
			return
		}
	}
}

func (h *TicketHandler) writeEmptyCSV(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=\"tickets.csv\"")
	writer := csv.NewWriter(w)
	defer writer.Flush()
	header := []string{"ID", "Organization ID", "Title", "Status", "Priority", "Reporter ID", "Created At", "Description"}
	if err := writer.Write(header); err != nil {
		h.logger.Error("failed to write csv header", "error", err)
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
		Sensitive:      req.Sensitive,
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
		OrganizationID:     &orgID,
		ExcludeDescription: true,
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

	if search := r.URL.Query().Get("search"); search != "" {
		filter.Keyword = &search
	}

	tickets, err := h.service.ListTickets(r.Context(), filter)
	if err != nil {
		h.logger.Error("failed to list tickets", "error", err)
		http.Error(w, "Failed to list tickets", http.StatusInternalServerError)
		return
	}

	if len(tickets) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode([]TicketDetailResponse{}); err != nil {
			h.logger.Error("failed to encode response", "error", err)
		}
		return
	}

	members, err := h.orgRepo.ListMembers(r.Context(), orgID)
	if err != nil {
		h.logger.Error("failed to list organization members", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	memberMap := make(map[uuid.UUID]string)
	for _, m := range members {
		memberMap[m.UserID] = m.Name
	}

	response := make([]TicketDetailResponse, len(tickets))
	for i := range tickets {
		t := &tickets[i]
		assigneeName := ""
		if t.AssigneeUserID != nil {
			assigneeName = memberMap[*t.AssigneeUserID]
		}
		reporterName := memberMap[t.ReporterID]

		response[i] = TicketDetailResponse{
			Ticket:       t,
			AssigneeName: assigneeName,
			ReporterName: reporterName,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
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
		Sensitive:      req.Sensitive,
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
