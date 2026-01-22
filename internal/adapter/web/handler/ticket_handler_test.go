package handler_test

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
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

// Mocks
type MockTicketService struct {
	mock.Mock
}

func (m *MockTicketService) CreateTicket(ctx context.Context, cmd port.CreateTicketCmd) (*domain.Ticket, error) {
	args := m.Called(ctx, cmd)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Ticket), args.Error(1)
}

func (m *MockTicketService) UpdateTicket(ctx context.Context, id uuid.UUID, cmd port.UpdateTicketCmd) (*domain.Ticket, error) {
	args := m.Called(ctx, id, cmd)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Ticket), args.Error(1)
}

func (m *MockTicketService) ListTickets(ctx context.Context, filter port.TicketFilter) ([]domain.Ticket, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Ticket), args.Error(1)
}

func (m *MockTicketService) GetTicket(ctx context.Context, id uuid.UUID) (*domain.Ticket, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Ticket), args.Error(1)
}

type MockOrgRepo struct {
	mock.Mock
}

func (m *MockOrgRepo) Create(ctx context.Context, org *domain.Organization) error {
	return m.Called(ctx, org).Error(0)
}

func (m *MockOrgRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Organization, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Organization), args.Error(1)
}

func (m *MockOrgRepo) GetByShareToken(ctx context.Context, token string) (*domain.Organization, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Organization), args.Error(1)
}

func (m *MockOrgRepo) Update(ctx context.Context, org *domain.Organization) error {
	return m.Called(ctx, org).Error(0)
}

func (m *MockOrgRepo) ListByUser(ctx context.Context, userID uuid.UUID) ([]domain.UserMembership, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.UserMembership), args.Error(1)
}

func (m *MockOrgRepo) AddMember(ctx context.Context, orgID, userID uuid.UUID, role string) error {
	return m.Called(ctx, orgID, userID, role).Error(0)
}

func (m *MockOrgRepo) ListMembers(ctx context.Context, orgID uuid.UUID) ([]domain.Member, error) {
	args := m.Called(ctx, orgID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Member), args.Error(1)
}

func (m *MockOrgRepo) RemoveMember(ctx context.Context, orgID, userID uuid.UUID) error {
	return m.Called(ctx, orgID, userID).Error(0)
}

func (m *MockOrgRepo) UpdateMemberRole(ctx context.Context, orgID, userID uuid.UUID, role string) error {
	return m.Called(ctx, orgID, userID, role).Error(0)
}

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) Create(ctx context.Context, user *domain.User) error {
	return m.Called(ctx, user).Error(0)
}
func (m *MockUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}
func (m *MockUserRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}
func (m *MockUserRepo) Update(ctx context.Context, user *domain.User) error {
	return m.Called(ctx, user).Error(0)
}

func TestExportTickets(t *testing.T) {
	mockService := new(MockTicketService)
	mockOrgRepo := new(MockOrgRepo)
	mockUserRepo := new(MockUserRepo)
	h := handler.NewTicketHandler(mockService, mockOrgRepo, mockUserRepo, nil)

	r := chi.NewRouter()
	r.Get("/admin/export/tickets", h.ExportTickets)

	t.Run("Success - Admin exports tickets for their orgs", func(t *testing.T) {
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

		tickets := []domain.Ticket{
			{
				ID:             uuid.New(),
				Title:          "Ticket 1",
				OrganizationID: orgID,
				ReporterID:     uuid.New(),
				StatusID:       "open",
				PriorityID:     "high",
				CreatedAt:      time.Now(),
				Description:    "Description 1",
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

		// Header + 1 row
		assert.Equal(t, 2, len(records))
		assert.Equal(t, "ID", records[0][0])
		assert.Equal(t, tickets[0].Title, records[1][2])
	})

	t.Run("Forbidden - Non-admin user", func(t *testing.T) {
		regularUser := &domain.User{
			ID:   uuid.New(),
			Role: domain.RoleManager,
		}

		req := httptest.NewRequest("GET", "/admin/export/tickets", nil)
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, regularUser)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}

func TestCreatePublicTicket(t *testing.T) {
	t.Run("Success - Create Ticket for existing user", func(t *testing.T) {
		mockService := new(MockTicketService)
		mockOrgRepo := new(MockOrgRepo)
		mockUserRepo := new(MockUserRepo)
		h := handler.NewTicketHandler(mockService, mockOrgRepo, mockUserRepo, nil)
		r := chi.NewRouter()
		r.Post("/public/tickets", h.CreatePublicTicket)

		token := "valid-token"
		orgID := uuid.New()
		org := &domain.Organization{
			ID:               orgID,
			ShareLinkEnabled: true,
			ShareLinkToken:   &token,
		}
		userID := uuid.New()
		user := &domain.User{
			ID:    userID,
			Email: "test@example.com",
		}
		ticket := &domain.Ticket{
			ID:    uuid.New(),
			Title: "New Public Ticket",
		}

		reqBody := map[string]string{
			"token":       token,
			"title":       "New Public Ticket",
			"description": "Desc",
			"name":        "Tester",
			"email":       "test@example.com",
			"priority_id": "low",
		}
		bodyBytes, _ := json.Marshal(reqBody)

		mockOrgRepo.On("GetByShareToken", mock.Anything, token).Return(org, nil)
		mockUserRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(user, nil)
		mockService.On("CreateTicket", mock.Anything, mock.MatchedBy(func(cmd port.CreateTicketCmd) bool {
			return cmd.OrganizationID == orgID && cmd.ReporterID == userID && cmd.Title == "New Public Ticket"
		})).Return(ticket, nil)

		req := httptest.NewRequest("POST", "/public/tickets", bytes.NewReader(bodyBytes))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("Success - Create Ticket for new user", func(t *testing.T) {
		mockService := new(MockTicketService)
		mockOrgRepo := new(MockOrgRepo)
		mockUserRepo := new(MockUserRepo)
		h := handler.NewTicketHandler(mockService, mockOrgRepo, mockUserRepo, nil)
		r := chi.NewRouter()
		r.Post("/public/tickets", h.CreatePublicTicket)

		token := "valid-token"
		orgID := uuid.New()
		org := &domain.Organization{
			ID:               orgID,
			ShareLinkEnabled: true,
			ShareLinkToken:   &token,
		}

		reqBody := map[string]string{
			"token":       token,
			"title":       "New Public Ticket",
			"description": "Desc",
			"name":        "New User",
			"email":       "new@example.com",
			"priority_id": "low",
		}
		bodyBytes, _ := json.Marshal(reqBody)

		mockOrgRepo.On("GetByShareToken", mock.Anything, token).Return(org, nil)
		// User not found
		mockUserRepo.On("GetByEmail", mock.Anything, "new@example.com").Return(nil, nil)
		// Create User
		newUserID := uuid.New()
		mockUserRepo.On("Create", mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
			return u.Email == "new@example.com" && u.Name == "New User"
		})).Return(nil).Run(func(args mock.Arguments) {
			u := args.Get(1).(*domain.User)
			u.ID = newUserID
		})
		// Create Ticket
		mockService.On("CreateTicket", mock.Anything, mock.MatchedBy(func(cmd port.CreateTicketCmd) bool {
			return cmd.Title == "New Public Ticket"
		})).Return(&domain.Ticket{ID: uuid.New()}, nil).Run(func(args mock.Arguments) {
			cmd := args.Get(1).(port.CreateTicketCmd)
			assert.Equal(t, orgID, cmd.OrganizationID)
			assert.Equal(t, newUserID, cmd.ReporterID)
		})

		req := httptest.NewRequest("POST", "/public/tickets", bytes.NewReader(bodyBytes))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("Forbidden - Disabled Share Link", func(t *testing.T) {
		mockService := new(MockTicketService)
		mockOrgRepo := new(MockOrgRepo)
		mockUserRepo := new(MockUserRepo)
		h := handler.NewTicketHandler(mockService, mockOrgRepo, mockUserRepo, nil)
		r := chi.NewRouter()
		r.Post("/public/tickets", h.CreatePublicTicket)

		token := "disabled-token"
		org := &domain.Organization{
			ShareLinkEnabled: false,
		}
		reqBody := map[string]string{"token": token}
		bodyBytes, _ := json.Marshal(reqBody)

		mockOrgRepo.On("GetByShareToken", mock.Anything, token).Return(org, nil)

		req := httptest.NewRequest("POST", "/public/tickets", bytes.NewReader(bodyBytes))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}
