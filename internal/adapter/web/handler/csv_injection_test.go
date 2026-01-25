package handler_test

import (
	"context"
	"encoding/csv"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wsciaroni/opsdeck/internal/adapter/web/handler"
	"github.com/wsciaroni/opsdeck/internal/adapter/web/middleware"
	"github.com/wsciaroni/opsdeck/internal/core/domain"
	"github.com/wsciaroni/opsdeck/internal/core/port"
)

func TestExportTickets_CSVInjection(t *testing.T) {
	mockService := new(MockTicketService)
	mockOrgRepo := new(MockOrgRepo)
	mockUserRepo := new(MockUserRepo)
	h := handler.NewTicketHandler(mockService, mockOrgRepo, mockUserRepo, nil)

	r := chi.NewRouter()
	r.Get("/admin/export/tickets", h.ExportTickets)

	adminUser := &domain.User{
		ID:    uuid.New(),
		Role:  domain.RoleAdmin,
		Email: "admin@example.com",
	}

	orgID := uuid.New()
	memberships := []domain.UserMembership{
		{
			Organization: domain.Organization{ID: orgID},
			Role:         "admin",
		},
	}

	// Payload that triggers CSV injection (Formula Injection)
	maliciousTitle := "=cmd|' /C calc'!A0"
	maliciousDesc := "+SUM(1+1)*cmd|' /C calc'!A0"

	tickets := []domain.Ticket{
		{
			ID:             uuid.New(),
			Title:          maliciousTitle,
			OrganizationID: orgID,
			ReporterID:     uuid.New(),
			StatusID:       "open",
			PriorityID:     "high",
			CreatedAt:      time.Now(),
			Description:    maliciousDesc,
		},
	}

	// Expect checking memberships
	mockOrgRepo.On("ListByUser", mock.Anything, adminUser.ID).Return(memberships, nil)

	// Expect ListTickets with OrganizationIDs set
	mockService.On("ListTickets", mock.Anything, port.TicketFilter{
		OrganizationIDs: []uuid.UUID{orgID},
	}).Return(tickets, nil)

	req := httptest.NewRequest("GET", "/admin/export/tickets", nil)
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, adminUser)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "text/csv", w.Header().Get("Content-Type"))

	// Parse CSV
	reader := csv.NewReader(strings.NewReader(w.Body.String()))
	records, err := reader.ReadAll()
	assert.NoError(t, err)

	// Verify the raw output contains the malicious payload unmodified
	// This confirms the vulnerability exists (if we haven't fixed it yet)
	// We want to eventually assert that it IS sanitized (starts with ')
	// But for now, let's just inspect it.

	// Row 1 is header, Row 2 is data
	titleCell := records[1][2]
	descCell := records[1][7]

	// It should now be sanitized (prefixed with single quote)
	assert.Equal(t, "'"+maliciousTitle, titleCell, "Title should be sanitized")
	assert.Equal(t, "'"+maliciousDesc, descCell, "Description should be sanitized")
}
